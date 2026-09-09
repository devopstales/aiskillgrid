// Package community detects subsystems in the 005 symbols/edges graph:
// seeded Leiden partitioning (bluuewhale/loom), degree-ranked god nodes, and
// LLM-free community labels. The partition is advisory, never load-bearing.
package community

import (
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"
	"sort"
	"time"

	"github.com/bluuewhale/loom/graph"
)

// Options tunes the community pass.
type Options struct {
	// Seed is the RNG seed for the Leiden detector. A fixed non-zero seed makes
	// the partition reproducible; 0 lets loom run its multi-run best-Q mode.
	Seed int64
	// Resolution is the modularity resolution (default 1).
	Resolution float64
	// MaxIterations bounds the Leiden refinement loop (0 = unlimited).
	MaxIterations int
	// NumRuns selects loom's multi-run best-Q selection when Seed == 0.
	NumRuns int
	// Limit caps the number of god nodes attached per community (0 = 10).
	Limit int
}

// Community is one detected subsystem: its members (symbol ids) and its
// LLM-free label (derived from god nodes + paths, never fabricated).
type Community struct {
	ID       int      `json:"id"`
	Label    string   `json:"label"`
	Members  []int64  `json:"members"`
	GodNodes []string `json:"god_nodes,omitempty"`
}

// Result is the Detect output: the communities, the stable content-hash cache
// key, whether the pass was served from cache, and any non-fatal warnings.
type Result struct {
	Communities []Community `json:"communities"`
	CacheKey    string      `json:"cache_key"`
	FromCache   bool        `json:"from_cache"`
	Warning     string      `json:"warning,omitempty"`
	Modularity  float64     `json:"modularity,omitempty"`
	Updated     time.Time   `json:"updated_at"`
}

// nodeRef is a symbol node in the partition: its symbol id plus the registry
// key (uid preferred, else name) used to map it through loom.
type nodeRef struct {
	symbolID int64
	key      string
}

// loadGraph loads the 005 symbols/edges tables as a registry + undirected
// graph. Edges are treated undirected for community detection (a call
// relationship clusters the same either way). Only edges whose both endpoints
// are live symbols (to_id set) build connectivity; name-only dangling edges
// are ignored (they would fabricate nodes).
func loadGraph(db *sql.DB) (*graph.NodeRegistry, *graph.Graph, []nodeRef, error) {
	reg := graph.NewRegistry()
	g := graph.NewGraph(false)

	rows, err := db.Query(`SELECT id, uid, name FROM symbols ORDER BY id`)
	if err != nil {
		return nil, nil, nil, err
	}
	defer rows.Close()
	var nodes []nodeRef
	idToLoom := map[int64]graph.NodeID{}
	for rows.Next() {
		var id int64
		var uid, name sql.NullString
		if err := rows.Scan(&id, &uid, &name); err != nil {
			return nil, nil, nil, err
		}
		key := uid.String
		if key == "" {
			key = name.String
		}
		if key == "" {
			key = fmt.Sprintf("sym-%d", id)
		}
		nid := reg.Register(key)
		idToLoom[id] = nid
		g.AddNode(nid, 1.0)
		nodes = append(nodes, nodeRef{symbolID: id, key: key})
	}
	if err := rows.Err(); err != nil {
		return nil, nil, nil, err
	}

	eRows, err := db.Query(`SELECT from_id, to_id FROM edges WHERE to_id IS NOT NULL`)
	if err != nil {
		return nil, nil, nil, err
	}
	defer eRows.Close()
	seen := map[[2]graph.NodeID]bool{}
	for eRows.Next() {
		var fromID, toID int64
		if err := eRows.Scan(&fromID, &toID); err != nil {
			return nil, nil, nil, err
		}
		f, ok1 := idToLoom[fromID]
		to, ok2 := idToLoom[toID]
		if !ok1 || !ok2 || f == to {
			continue
		}
		pair := [2]graph.NodeID{f, to}
		if pair[1] < pair[0] {
			pair = [2]graph.NodeID{to, f}
		}
		if seen[pair] {
			continue
		}
		seen[pair] = true
		g.AddEdge(f, to, 1.0)
	}
	if err := eRows.Err(); err != nil {
		return nil, nil, nil, err
	}
	return reg, g, nodes, nil
}

// contentKey summarizes the edge set so the cache key is stable across index
// runs with the same graph and changes when the graph changes.
func contentKey(db *sql.DB) (string, error) {
	var a, b string
	err := db.QueryRow(`
		SELECT COALESCE(MIN(from_id),0) || '-' || COALESCE(MIN(to_id),0) || '-' ||
		       COALESCE(MAX(from_id),0) || '-' || COALESCE(MAX(to_id),0) || '-' ||
		       COUNT(*),
		       COUNT(DISTINCT from_id || '-' || COALESCE(to_id, 0))
		FROM edges WHERE to_id IS NOT NULL`).Scan(&a, &b)
	if err == sql.ErrNoRows {
		return "empty", nil
	}
	if err != nil {
		return "", err
	}
	return a + "|" + b, nil
}

// Detect runs the community pass, consulting the community_meta cache first.
// A matching content-hash cache key returns the stored partition
// (FromCache=true) without re-running the detector; otherwise the seeded
// Leiden pass runs and the result is cached.
func Detect(ctx context.Context, db *sql.DB, opts Options) (*Result, error) {
	if opts.Seed == 0 {
		opts.Seed = 1
	}
	if opts.Resolution <= 0 {
		opts.Resolution = 1
	}
	if opts.Limit <= 0 {
		opts.Limit = 10
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	key, err := contentKey(db)
	if err != nil {
		return nil, fmt.Errorf("content key: %w", err)
	}
	cacheKey := hex.EncodeToString([]byte("comm:" + key))

	if cached, err := cachedResult(db, cacheKey); err == nil && cached != nil {
		cached.FromCache = true
		return cached, nil
	}

	res, err := detectCore(db, opts, cacheKey)
	if err != nil {
		return nil, err
	}
	if err := writeMeta(db, res); err != nil {
		return res, fmt.Errorf("detect ok but cache write failed: %w", err)
	}
	return res, nil
}

// detectCore runs the seeded Leiden pass, assigns labels, and writes the
// communities table (the 012 additive table; 005 symbols/edges untouched).
func detectCore(db *sql.DB, opts Options, cacheKey string) (*Result, error) {
	res := &Result{CacheKey: cacheKey, Updated: time.Now().UTC()}

	var symCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM symbols`).Scan(&symCount); err != nil {
		return nil, err
	}

	if symCount < 2 {
		// warn+continue: one trivial community, never a crash.
		res.Warning = fmt.Sprintf("graph has %d symbols (<2); producing a single trivial community", symCount)
		ids, err := allSymbolIDs(db)
		if err != nil {
			return nil, err
		}
		gods, label := labelForCommunity(db, 0, ids, opts)
		res.Communities = []Community{{ID: 0, Label: label, Members: ids, GodNodes: gods}}
		return res, nil
	}

	reg, g, nodes, err := loadGraph(db)
	if err != nil {
		return nil, fmt.Errorf("load graph: %w", err)
	}

	det := graph.NewLeiden(graph.LeidenOptions{
		Seed:          opts.Seed,
		Resolution:    opts.Resolution,
		MaxIterations: opts.MaxIterations,
		NumRuns:       opts.NumRuns,
	})
	part, err := det.Detect(g)
	if err != nil {
		return nil, fmt.Errorf("leiden: %w", err)
	}
	res.Modularity = part.Modularity

	// Invert the partition into communities: loom node -> member symbol ids.
	loomToSymbol := map[graph.NodeID]int64{}
	for _, n := range nodes {
		if id, ok := reg.ID(n.key); ok {
			loomToSymbol[id] = n.symbolID
		}
	}
	byComm := map[int][]int64{}
	for node, cid := range part.Partition {
		sid, ok := loomToSymbol[node]
		if !ok {
			continue
		}
		byComm[cid] = append(byComm[cid], sid)
	}
	var commIDs []int
	for cid := range byComm {
		commIDs = append(commIDs, cid)
	}
	sort.Ints(commIDs)

	communities := make([]Community, 0, len(byComm))
	for i, cid := range commIDs {
		members := byComm[cid]
		sort.Slice(members, func(a, b int) bool { return members[a] < members[b] })
		gods, label := labelForCommunity(db, i, members, opts)
		communities = append(communities, Community{ID: i, Label: label, Members: members, GodNodes: gods})
	}
	res.Communities = communities

	if err := writeCommunities(db, communities); err != nil {
		return nil, err
	}
	return res, nil
}

// allSymbolIDs returns every symbol id, ascending.
func allSymbolIDs(db *sql.DB) ([]int64, error) {
	rows, err := db.Query(`SELECT id FROM symbols ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

// writeCommunities replaces the communities table with the detected
// partition (idempotent).
func writeCommunities(db *sql.DB, communities []Community) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.Exec(`DELETE FROM communities`); err != nil {
		return err
	}
	for _, c := range communities {
		for _, m := range c.Members {
			if _, err := tx.Exec(`INSERT OR REPLACE INTO communities (id, symbol_id) VALUES (?, ?)`, c.ID, m); err != nil {
				return err
			}
		}
	}
	return tx.Commit()
}

// cachedResult rebuilds a Result from the stored communities table when the
// content-hash cache key matches. Returns nil (no error) on a cache miss.
func cachedResult(db *sql.DB, cacheKey string) (*Result, error) {
	var stored string
	err := db.QueryRow(`SELECT value FROM community_meta_cache WHERE key = 'cache_key'`).Scan(&stored)
	if err != nil || stored != cacheKey {
		return nil, nil
	}
	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM communities`).Scan(&count); err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, nil
	}
	res := &Result{CacheKey: cacheKey, Updated: time.Now().UTC()}
	rows, err := db.Query(`SELECT id, symbol_id FROM communities ORDER BY id, symbol_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	byID := map[int][]int64{}
	var order []int
	for rows.Next() {
		var id int
		var sid int64
		if err := rows.Scan(&id, &sid); err != nil {
			return nil, err
		}
		if _, ok := byID[id]; !ok {
			order = append(order, id)
		}
		byID[id] = append(byID[id], sid)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	sort.Ints(order)
	for i, id := range order {
		res.Communities = append(res.Communities, Community{
			ID:      i,
			Label:   communityLabel(db, id),
			Members: byID[id],
		})
	}
	return res, nil
}

// communityLabel reads a stored label (community_meta keyed by community id).
func communityLabel(db *sql.DB, id int) string {
	var label string
	_ = db.QueryRow(`SELECT label FROM community_meta WHERE id = ?`, id).Scan(&label)
	if label == "" {
		return "community-" + fmt.Sprintf("%d", id)
	}
	return label
}

// writeMeta caches the detected partition: one community_meta row per
// community (label + symbol count) and the content-hash cache key.
func writeMeta(db *sql.DB, res *Result) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.Exec(`DELETE FROM community_meta`); err != nil {
		return err
	}
	for _, c := range res.Communities {
		if _, err := tx.Exec(`INSERT INTO community_meta (id, label, symbol_count, cache_key) VALUES (?, ?, ?, ?)`,
			c.ID, c.Label, len(c.Members), res.CacheKey); err != nil {
			return err
		}
	}
	if _, err := tx.Exec(`CREATE TABLE IF NOT EXISTS community_meta_cache (key TEXT PRIMARY KEY, value TEXT)`); err != nil {
		return err
	}
	if _, err := tx.Exec(`INSERT INTO community_meta_cache (key, value) VALUES ('cache_key', ?) ON CONFLICT(key) DO UPDATE SET value = excluded.value`, res.CacheKey); err != nil {
		return err
	}
	return tx.Commit()
}

// labelForCommunity is a hook for tests; defaults to the LLM-free derivation
// in labels.go.
var labelForCommunity = labelForCommunityImpl
