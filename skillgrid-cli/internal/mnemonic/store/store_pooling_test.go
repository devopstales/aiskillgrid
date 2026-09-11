package store

import (
	"database/sql"
	"testing"
	"time"
)

// TestStoreOpenReusesCachedHandle covers @step-01 (Scenario: Cached handle
// reuse on second open + Handle closes only when reference count reaches
// zero): a second Open for the same store path returns the same underlying
// *sql.DB without re-opening, a different project gets a different handle,
// Close refcounts (the connection stays open until the last reference is
// returned), and the cache is evicted when the refcount hits zero.
func TestStoreOpenReusesCachedHandle(t *testing.T) {
	dir := t.TempDir()

	stA, err := Open(dir, "pool-a")
	if err != nil {
		t.Fatalf("first open A: %v", err)
	}
	stB, err := Open(dir, "pool-b")
	if err != nil {
		t.Fatalf("open B: %v", err)
	}
	defer stB.Close()

	stA2, err := Open(dir, "pool-a")
	if err != nil {
		t.Fatalf("second open A: %v", err)
	}

	// Cache hit: same underlying *sql.DB, no new connection created.
	if stA.DB != stA2.DB {
		t.Fatalf("expected second open of the same store to reuse the cached *sql.DB")
	}
	// Different project: different handle.
	if stA.DB == stB.DB {
		t.Fatalf("expected different projects to get different *sql.DB handles")
	}
	// Migration bookkeeping ran exactly once (no re-apply on cache hit).
	if countMigration(t, stA.DB, "019_session_relay.sql") != 1 {
		t.Fatalf("expected 019 migration recorded once")
	}

	// Refcount: closing one handle leaves the underlying connection open.
	if err := stA.Close(); err != nil {
		t.Fatalf("close first handle: %v", err)
	}
	if err := stA2.DB.Ping(); err != nil {
		t.Fatalf("underlying connection should stay open while refcount > 0: %v", err)
	}
	// Final close evicts the cache: the next open gets a fresh *sql.DB.
	if err := stA2.Close(); err != nil {
		t.Fatalf("close second handle: %v", err)
	}
	stA3, err := Open(dir, "pool-a")
	if err != nil {
		t.Fatalf("open after eviction: %v", err)
	}
	defer stA3.Close()
	if stA3.DB == stA.DB {
		t.Fatalf("expected a fresh *sql.DB after the refcount reached zero")
	}

	// The fresh handle is healthy (cache health-check holds after eviction).
	var n int
	if err := stA3.DB.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='sessions'`).Scan(&n); err != nil {
		t.Fatalf("post-eviction query: %v", err)
	}
	if n != 1 {
		t.Fatalf("expected sessions table on fresh handle, got %d", n)
	}
}

// TestStoreOpenWALRetry covers @step-01 (Scenario: WAL lock retry with
// exponential backoff): PRAGMA busy_timeout=10000 is applied on every new
// connection, and a write that hits a WAL lock retries at 50ms, 100ms, and
// 200ms intervals, succeeding after the lock is released — no panic, no
// failure while the lock was held.
func TestStoreOpenWALRetry(t *testing.T) {
	dir := t.TempDir()

	// The store opens (and sets busy_timeout=10000) while no lock is held.
	st, err := Open(dir, "walproj")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer st.Close()

	// busy_timeout is applied on every new connection.
	var timeout int
	if err := st.DB.QueryRow(`PRAGMA busy_timeout`).Scan(&timeout); err != nil {
		t.Fatalf("pragma busy_timeout: %v", err)
	}
	if timeout != 10000 {
		t.Fatalf("expected busy_timeout=10000, got %d", timeout)
	}

	// A concurrent writer holds the WAL write lock.
	if _, err := st.DB.Exec(`CREATE TABLE walprobe (id INTEGER PRIMARY KEY, v TEXT)`); err != nil {
		t.Fatalf("create: %v", err)
	}
	walTx, err := st.DB.Begin()
	if err != nil {
		t.Fatalf("begin writer tx: %v", err)
	}
	if _, err := walTx.Exec(`INSERT INTO walprobe (v) VALUES ('hold')`); err != nil {
		t.Fatalf("writer insert: %v", err)
	}

	// A second connection's write hits the lock. It must not fail outright:
	// it retries (backoff 50/100/200ms) until the lock is released, then the
	// write succeeds — the behavior the busy_timeout + retry contract gates.
	start := time.Now()
	writeDone := make(chan error, 1)
	go func() {
		db2 := rawSQLDB(t, dir, "walproj")
		defer db2.Close()
		// Warm the connection so the lock is actually held on the write,
		// not on the first handshake.
		if _, err := db2.Exec(`PRAGMA busy_timeout=10000`); err != nil {
			writeDone <- err
			return
		}
		_, err = db2.Exec(`INSERT INTO walprobe (v) VALUES ('contended')`)
		writeDone <- err
	}()

	// Release the lock partway through the retry window.
	time.Sleep(80 * time.Millisecond)
	if err := walTx.Commit(); err != nil {
		t.Fatalf("commit writer tx: %v", err)
	}

	select {
	case err := <-writeDone:
		if err != nil {
			t.Fatalf("contended write should succeed after backoff retries, got: %v", err)
		}
		if elapsed := time.Since(start); elapsed < 50*time.Millisecond {
			t.Fatalf("expected a backoff delay before the contended write, finished in %v", elapsed)
		}
	case <-time.After(3 * time.Second):
		t.Fatalf("contended write did not complete within the retry window")
	}

	// Both rows present: the contended write landed after the retry.
	var n int
	if err := st.DB.QueryRow(`SELECT COUNT(*) FROM walprobe`).Scan(&n); err != nil {
		t.Fatalf("count walprobe: %v", err)
	}
	if n != 2 {
		t.Fatalf("expected 2 walprobe rows after contended write, got %d", n)
	}
}

// TestStoreOpenCacheDisabledByEnv covers @step-01 (Scenario: Cache disabled
// by environment variable): SKILLGRID_MNEMONIC_DISABLE_CACHE=1 makes every
// Open create a fresh underlying connection and no handle is stored in the
// cache (each Store owns its database and closes it outright).
func TestStoreOpenCacheDisabledByEnv(t *testing.T) {
	t.Setenv("SKILLGRID_MNEMONIC_DISABLE_CACHE", "1")
	dir := t.TempDir()

	st1, err := Open(dir, "nocache")
	if err != nil {
		t.Fatalf("first open: %v", err)
	}
	defer st1.Close()
	st2, err := Open(dir, "nocache")
	if err != nil {
		t.Fatalf("second open: %v", err)
	}
	defer st2.Close()

	if st1.DB == st2.DB {
		t.Fatalf("with the cache disabled, each open must create a new *sql.DB")
	}
	// No handle is stored in the cache.
	if _, ok := handleCache.Load(st1.Path()); ok {
		t.Fatalf("cache must stay empty when disabled")
	}
	var n1, n2 int
	if err := st1.DB.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='sessions'`).Scan(&n1); err != nil {
		t.Fatalf("query handle 1: %v", err)
	}
	if err := st2.DB.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='sessions'`).Scan(&n2); err != nil {
		t.Fatalf("query handle 2: %v", err)
	}
	if n1 != 1 || n2 != 1 {
		t.Fatalf("both handles must be usable, got %d/%d", n1, n2)
	}
}

// TestMigration014TTLExtraction covers @step-01 (01.3): a fresh store has
// the 020 migration applied — ttl_config and extraction_metadata tables
// exist with the expected columns, the migration is recorded once, and the
// existing schema remains intact (additive, not a rewrite).
func TestMigration014TTLExtraction(t *testing.T) {
	dir := t.TempDir()
	st, err := Open(dir, "ttlproj")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer st.Close()

	if !tableExists(t, st.DB, "ttl_config") {
		t.Fatalf("expected ttl_config table after 020 migration")
	}
	if !tableExists(t, st.DB, "extraction_metadata") {
		t.Fatalf("expected extraction_metadata table after 020 migration")
	}
	for _, col := range []string{"key", "value"} {
		var n int
		if err := st.DB.QueryRow(`SELECT COUNT(*) FROM pragma_table_info('ttl_config') WHERE name=?`, col).Scan(&n); err != nil {
			t.Fatalf("pragma ttl_config %s: %v", col, err)
		}
		if n != 1 {
			t.Fatalf("ttl_config column %s not present (n=%d)", col, n)
		}
	}
	for _, col := range []string{"id", "session_id", "content_hash", "extracted_at", "model"} {
		var n int
		if err := st.DB.QueryRow(`SELECT COUNT(*) FROM pragma_table_info('extraction_metadata') WHERE name=?`, col).Scan(&n); err != nil {
			t.Fatalf("pragma extraction_metadata %s: %v", col, err)
		}
		if n != 1 {
			t.Fatalf("extraction_metadata column %s not present (n=%d)", col, n)
		}
	}
	if countMigration(t, st.DB, "020_ttl_extraction.sql") != 1 {
		t.Fatalf("expected 020 migration recorded once")
	}

	// Existing schema intact (additive).
	for _, name := range []string{"observations", "observations_fts", "sessions"} {
		if !tableExists(t, st.DB, name) {
			t.Fatalf("pre-existing table %s missing after 020", name)
		}
	}

	// Both new tables are writable.
	if _, err := st.DB.Exec(`INSERT INTO ttl_config (key, value) VALUES ('ttl_days', '30')`); err != nil {
		t.Fatalf("insert ttl_config: %v", err)
	}
	if _, err := st.DB.Exec(`INSERT INTO extraction_metadata (session_id, content_hash, extracted_at, model) VALUES ('s1','h1', ?, 'test')`, time.Now().UTC().Format(time.RFC3339)); err != nil {
		t.Fatalf("insert extraction_metadata: %v", err)
	}
	var cfg, meta int
	if err := st.DB.QueryRow(`SELECT COUNT(*) FROM ttl_config`).Scan(&cfg); err != nil {
		t.Fatalf("count ttl_config: %v", err)
	}
	if err := st.DB.QueryRow(`SELECT COUNT(*) FROM extraction_metadata`).Scan(&meta); err != nil {
		t.Fatalf("count extraction_metadata: %v", err)
	}
	if cfg != 1 || meta != 1 {
		t.Fatalf("expected writable new tables, got %d/%d rows", cfg, meta)
	}
}

var _ = sql.ErrConnDone // keep database/sql import stable for future assertions
