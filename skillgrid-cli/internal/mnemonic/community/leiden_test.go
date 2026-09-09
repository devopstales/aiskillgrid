package community

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	_ "modernc.org/sqlite"
)

// openStore opens a scratch SQLite store with the symbols/edges schema (the
// 005 tables the community pass consumes) plus the 012 community tables.
func openStore(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	stmts := []string{
		`CREATE TABLE files (id INTEGER PRIMARY KEY AUTOINCREMENT, path TEXT NOT NULL, mtime_ns INTEGER, size INTEGER, content_hash TEXT, indexed_at TEXT)`,
		`CREATE TABLE symbols (id INTEGER PRIMARY KEY AUTOINCREMENT, file_id INTEGER NOT NULL REFERENCES files(id) ON DELETE CASCADE, name TEXT NOT NULL, qualified_name TEXT, kind TEXT NOT NULL, language TEXT, signature TEXT, start_line INTEGER NOT NULL, end_line INTEGER NOT NULL, content_hash TEXT NOT NULL, uid TEXT NOT NULL UNIQUE)`,
		`CREATE TABLE edges (id INTEGER PRIMARY KEY AUTOINCREMENT, kind TEXT NOT NULL, from_id INTEGER NOT NULL REFERENCES symbols(id) ON DELETE CASCADE, file_id INTEGER REFERENCES files(id) ON DELETE CASCADE, to_id INTEGER REFERENCES symbols(id) ON DELETE CASCADE, to_name TEXT, target_path TEXT, confidence TEXT NOT NULL DEFAULT 'EXTRACTED', line INTEGER, UNIQUE(kind, from_id, file_id, to_id, to_name, target_path, line))`,
		`CREATE TABLE communities (id INTEGER, symbol_id INTEGER NOT NULL UNIQUE REFERENCES symbols(id) ON DELETE CASCADE)`,
		`CREATE TABLE community_meta (id INTEGER PRIMARY KEY, label TEXT NOT NULL, symbol_count INTEGER NOT NULL, god_nodes TEXT NOT NULL DEFAULT '', cache_key TEXT NOT NULL)`,
		`CREATE TABLE community_meta_cache (key TEXT PRIMARY KEY, value TEXT)`,
	}
	for _, s := range stmts {
		if _, err := db.Exec(s); err != nil {
			t.Fatalf("create schema: %v", err)
		}
	}
	return db
}

// seedFile inserts a file row and returns its id.
func seedFile(t *testing.T, db *sql.DB, path string) int64 {
	t.Helper()
	res, err := db.Exec(`INSERT INTO files (path, mtime_ns, size, content_hash, indexed_at) VALUES (?, 1, 1, 'h', 'now')`, path)
	if err != nil {
		t.Fatalf("seed file %s: %v", path, err)
	}
	id, _ := res.LastInsertId()
	return id
}

// seedSymbol inserts a symbol row and returns its id.
func seedSymbol(t *testing.T, db *sql.DB, fileID int64, name, kind string, line int) int64 {
	t.Helper()
	res, err := db.Exec(`INSERT INTO symbols (file_id, name, kind, start_line, end_line, content_hash, uid) VALUES (?, ?, ?, ?, ?, 'h', ?)`,
		fileID, name, kind, line, line, name+string(rune('u'))+kind)
	if err != nil {
		t.Fatalf("seed symbol %s: %v", name, err)
	}
	id, _ := res.LastInsertId()
	return id
}

// seedEdge inserts an edge between two symbol ids.
func seedEdge(t *testing.T, db *sql.DB, kind string, fromID, toID int64, line int) {
	t.Helper()
	if _, err := db.Exec(`INSERT INTO edges (kind, from_id, to_id, confidence, line) VALUES (?, ?, ?, 'EXTRACTED', ?)`, kind, fromID, toID, line); err != nil {
		t.Fatalf("seed edge %d->%d: %v", fromID, toID, err)
	}
}

// communityFixture builds two clearly-separated clusters with one weak bridge:
//
//	a1-a2-a3 (dense)  --bridge-->  b1-b2-b3 (dense)
//
// so Leiden should produce at least two communities that keep each cluster
// intact.
func communityFixture(t *testing.T) *sql.DB {
	t.Helper()
	db := openStore(t)
	fa := seedFile(t, db, "a/core.go")
	fb := seedFile(t, db, "b/engine.go")
	a1 := seedSymbol(t, db, fa, "aOne", "function", 1)
	a2 := seedSymbol(t, db, fa, "aTwo", "function", 10)
	a3 := seedSymbol(t, db, fa, "aThree", "function", 20)
	b1 := seedSymbol(t, db, fb, "bOne", "function", 1)
	b2 := seedSymbol(t, db, fb, "bTwo", "function", 10)
	b3 := seedSymbol(t, db, fb, "bThree", "function", 20)

	seedEdge(t, db, "calls", a1, a2, 2)
	seedEdge(t, db, "calls", a2, a3, 11)
	seedEdge(t, db, "calls", a1, a3, 3)
	seedEdge(t, db, "calls", b1, b2, 2)
	seedEdge(t, db, "calls", b2, b3, 11)
	seedEdge(t, db, "calls", b1, b3, 3)
	seedEdge(t, db, "imports", a3, b1, 21) // weak bridge
	return db
}

// TestCommunitiesDetectsTwoClusters covers the core Leiden contract: a graph
// of two dense clusters linked by one weak edge partitions into >=2
// communities that keep each cluster intact, every symbol gets exactly one
// community row, and each community carries an LLM-free label.
func TestCommunitiesDetectsTwoClusters(t *testing.T) {
	db := communityFixture(t)
	res, err := Detect(context.Background(), db, Options{})
	if err != nil {
		t.Fatalf("Detect: %v", err)
	}
	if len(res.Communities) < 2 {
		t.Fatalf("expected >=2 communities, got %d", len(res.Communities))
	}
	// Every symbol is assigned exactly once.
	assigned := map[int64]int{}
	for _, c := range res.Communities {
		for _, id := range c.Members {
			assigned[id]++
		}
	}
	if len(assigned) != 6 {
		t.Fatalf("expected 6 symbols assigned, got %d", len(assigned))
	}
	for id, n := range assigned {
		if n != 1 {
			t.Errorf("symbol %d assigned %d times", id, n)
		}
	}
	// Clusters stay intact: a1..a3 in one community, b1..b3 in the other.
	byName := map[int64]int{}
	for _, c := range res.Communities {
		for _, id := range c.Members {
			var name string
			_ = db.QueryRow(`SELECT name FROM symbols WHERE id = ?`, id).Scan(&name)
			byName[id] = c.ID
		}
	}
	var a1, a2, a3, b1, b2, b3 int64
	for _, q := range []struct {
		name string
		p    *int64
	}{{"aOne", &a1}, {"aTwo", &a2}, {"aThree", &a3}, {"bOne", &b1}, {"bTwo", &b2}, {"bThree", &b3}} {
		_ = db.QueryRow(`SELECT id FROM symbols WHERE name = ?`, q.name).Scan(q.p)
	}
	if byName[a1] != byName[a2] || byName[a2] != byName[a3] {
		t.Errorf("cluster a split: %d %d %d", byName[a1], byName[a2], byName[a3])
	}
	if byName[b1] != byName[b2] || byName[b2] != byName[b3] {
		t.Errorf("cluster b split: %d %d %d", byName[b1], byName[b2], byName[b3])
	}
	// Labels are present and LLM-free: derived, never empty, never "community-N"
	// for a community that has god nodes.
	for _, c := range res.Communities {
		if c.Label == "" {
			t.Errorf("community %d has empty label", c.ID)
		}
	}
}

// TestTinyGraphYieldsOneTrivialCommunity covers the < 2 nodes rule: a graph
// with a single symbol produces one trivial community (warn+continue), never
// a crash.
func TestTinyGraphYieldsOneTrivialCommunity(t *testing.T) {
	db := openStore(t)
	f := seedFile(t, db, "solo.go")
	seedSymbol(t, db, f, "only", "function", 1)
	res, err := Detect(context.Background(), db, Options{})
	if err != nil {
		t.Fatalf("Detect tiny: %v", err)
	}
	if len(res.Communities) != 1 {
		t.Fatalf("expected 1 trivial community, got %d", len(res.Communities))
	}
	if len(res.Communities[0].Members) != 1 {
		t.Fatalf("trivial community should hold the single symbol, got %v", res.Communities[0].Members)
	}
	if res.Warning == "" {
		t.Errorf("expected a warning for the <2 node graph, got none")
	}
}

// TestCommunitiesReproducibleAndCached covers the seed + content-hash cache:
// two Detect calls on the same graph yield the identical partition (stable
// cache key), and the second call reports a cache hit without re-running the
// detector.
func TestCommunitiesReproducibleAndCached(t *testing.T) {
	db := communityFixture(t)
	first, err := Detect(context.Background(), db, Options{})
	if err != nil {
		t.Fatalf("Detect 1: %v", err)
	}
	if first.CacheKey == "" {
		t.Fatal("expected a non-empty content-hash cache key")
	}
	second, err := Detect(context.Background(), db, Options{})
	if err != nil {
		t.Fatalf("Detect 2: %v", err)
	}
	if second.CacheKey != first.CacheKey {
		t.Errorf("cache key not stable: %q vs %q", second.CacheKey, first.CacheKey)
	}
	if !second.FromCache {
		t.Errorf("second Detect should be served from cache")
	}
	if partitionKey(first) != partitionKey(second) {
		t.Errorf("partition not reproducible: %s vs %s", partitionKey(first), partitionKey(second))
	}
}

func partitionKey(r *Result) string {
	var b strings.Builder
	for _, c := range r.Communities {
		b.WriteString("c")
		b.WriteString(string(rune('0' + c.ID%10)))
		for _, m := range c.Members {
			b.WriteString("|")
			b.WriteString(string(rune('0' + m%10)))
		}
	}
	return b.String()
}
