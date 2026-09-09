package knowledge

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
)

// Store is the knowledge pass's persistence seam: it upserts knowledge nodes
// (doc/config/sql) and their confidence-labeled edges into the 005 graph. The
// indexer's same-transaction hook passes a Store backed by the open index tx,
// so a knowledge node is always resolvable and a single rollback undoes the
// whole pass (matching the 010 route-hook pattern).
type Store struct {
	db *sql.DB
}

// NewStore wraps a *sql.DB (the indexer's open transaction in production, a
// scratch store in tests).
func NewStore(db *sql.DB) *Store { return &Store{db: db} }

// docNodeID resolves (or lazily creates) the doc node id for a file path.
// It is idempotent: a doc node is one-per-file (015 idx_doc_nodes_file unique).
func (s *Store) docNodeID(path string) (int64, error) {
	var id int64
	err := s.db.QueryRow(`SELECT id FROM doc_nodes WHERE path = ?`, path).Scan(&id)
	if err == nil {
		return id, nil
	}
	if err != sql.ErrNoRows {
		return 0, err
	}
	// The doc's source file must exist for the cascade; resolve it by path.
	var fileID int64
	if ferr := s.db.QueryRow(`SELECT id FROM files WHERE path = ?`, path).Scan(&fileID); ferr != nil {
		return 0, ferr
	}
	res, err := s.db.Exec(`INSERT INTO doc_nodes (file_id, title, path) VALUES (?, ?, ?)`,
		fileID, DocTitle(path, nil), path)
	if err != nil {
		return 0, fmt.Errorf("insert doc node %s: %w", path, err)
	}
	return res.LastInsertId()
}

// SaveDoc upserts one doc node and its references edges (target-state: the
// doc's prior references edges are pruned first, then the freshly extracted
// set is written). A link target that does not resolve to an indexed doc node
// is dropped (drop-not-guess for docs, per 03.2 — a markdown link to a
// non-doc is not a doc->doc reference). Never errors on a bad link: it is
// simply skipped.
func (s *Store) SaveDoc(ctx context.Context, path string, doc *DocResult) (int, error) {
	if doc == nil || !doc.IsDoc {
		return 0, nil
	}
	nodeID, err := s.docNodeID(path)
	if err != nil {
		return 0, err
	}
	// Target-state: prune this doc's prior references edges, then upsert.
	if _, err := s.db.Exec(`
		DELETE FROM edges
		WHERE kind = 'references' AND from_id = ?`, nodeID); err != nil {
		return 0, err
	}
	stored := 0
	for _, link := range doc.Links {
		target := resolveDocTarget(path, link.Target)
		toID, ok := s.resolveDocNode(target)
		if !ok {
			continue // unresolvable doc reference → dropped (not fabricated)
		}
		if _, err := s.db.Exec(`
			INSERT INTO edges (kind, from_id, file_id, to_id, to_name, target_path, confidence, line)
			VALUES ('references', ?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT(kind, from_id, file_id, to_id, to_name, target_path, line) DO UPDATE SET
			  confidence = excluded.confidence`,
			nodeID, fileIDFor(s.db, path), toID, target, target, link.Confidence, link.Line); err != nil {
			return stored, fmt.Errorf("upsert doc reference %s: %w", link.Target, err)
		}
		stored++
	}
	return stored, nil
}

// resolveDocNode resolves a doc target path to its doc node id (creating the
// node if the source file is indexed but the node is not yet written). It
// returns ok=false when the target is not an indexed doc (no file row).
func (s *Store) resolveDocNode(target string) (int64, bool) {
	var id int64
	if err := s.db.QueryRow(`SELECT id FROM doc_nodes WHERE path = ?`, target).Scan(&id); err == nil {
		return id, true
	}
	if err := s.db.QueryRow(`SELECT id FROM files WHERE path = ?`, target).Scan(new(int64)); err != nil {
		return 0, false // the target file is not indexed → not a doc node
	}
	// The file is indexed but the doc node is not yet written (e.g. the target
	// doc is processed later in the run). Create it now so the edge resolves.
	var fileID int64
	if err := s.db.QueryRow(`SELECT id FROM files WHERE path = ?`, target).Scan(&fileID); err != nil {
		return 0, false
	}
	res, err := s.db.Exec(`INSERT INTO doc_nodes (file_id, title, path) VALUES (?, ?, ?)`,
		fileID, DocTitle(target, nil), target)
	if err != nil {
		return 0, false
	}
	id, _ = res.LastInsertId()
	return id, true
}

// resolveDocTarget normalizes a markdown link target relative to the source
// doc's directory into an index-relative path (slash-separated).
func resolveDocTarget(from, target string) string {
	base := filepath.Dir(from)
	joined := filepath.Join(base, target)
	return filepath.ToSlash(joined)
}

// fileIDFor resolves a file path to its files.id (0 when unknown — the edges
// table allows a null file_id for knowledge edges).
func fileIDFor(db *sql.DB, path string) int64 {
	var id int64
	_ = db.QueryRow(`SELECT id FROM files WHERE path = ?`, path).Scan(&id)
	return id
}
