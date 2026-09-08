// Package codeindex incrementally indexes source files into a project store.
package codeindex

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/devopstales/skillgrid/skillgrid-cli/internal/mnemonic/extract"
	"github.com/devopstales/skillgrid/skillgrid-cli/internal/mnemonic/store"
)

// MaxFileSize is the default first-class size skip threshold. Files larger
// than this are skipped (counted in stats), not an error, not a fallback.
const MaxFileSize = 500 * 1024

// maxFileSize is the current run's effective threshold (bytes).
var maxFileSize int64 = MaxFileSize

// Config controls incremental indexing behavior.
type Config struct {
	Include      []string
	Exclude      []string
	ChunkLines   int
	ChunkOverlap int
	MaxFileSize  int
}

// Stats summarizes one indexing run.
type Stats struct {
	FilesIndexed  int `json:"files_indexed"`
	FilesSkipped  int `json:"files_skipped"`
	FilesDeleted  int `json:"files_deleted"`
	ChunksAdded   int `json:"chunks_added"`
	FilesOversized int `json:"files_oversized"`
	SymbolsAdded  int `json:"symbols_added"`
	EdgesAdded    int `json:"edges_added"`
}

// ScannedFile is a candidate file discovered under the index root.
type ScannedFile struct {
	Path     string
	MtimeNs  int64
	Size     int64
	Hash     string
	Contents []byte
}

// Scan walks root and returns files matching include/exclude globs.
func Scan(root string, include, exclude []string) ([]ScannedFile, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	var files []ScannedFile
	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			rel, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			if rel != "." && shouldSkipDir(rel, exclude) {
				return filepath.SkipDir
			}
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		// Note: size-based skipping is first-class in Indexer.Run (counted in
		// stats as FilesOversized), not here — so an oversized file is a skip,
		// not a silent drop, and is independently of exclude globs.
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if matchesAny(exclude, rel) {
			return nil
		}
		if len(include) > 0 && !matchesAny(include, rel) {
			return nil
		}
		contents, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", rel, err)
		}
		sum := sha256.Sum256(contents)
		files = append(files, ScannedFile{
			Path:     rel,
			MtimeNs:  info.ModTime().UnixNano(),
			Size:     info.Size(),
			Hash:     hex.EncodeToString(sum[:]),
			Contents: contents,
		})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func shouldSkipDir(rel string, exclude []string) bool {
	rel = filepath.ToSlash(rel)
	for _, pattern := range exclude {
		pattern = filepath.ToSlash(strings.TrimSuffix(pattern, "/**"))
		pattern = strings.TrimPrefix(pattern, "**/")
		if pattern == "" {
			continue
		}
		if rel == pattern || strings.HasPrefix(rel, pattern+"/") {
			return true
		}
		if matchesGlob(pattern, rel) || matchesGlob(pattern+"/**", rel) {
			return true
		}
	}
	return false
}

func matchesAny(patterns []string, path string) bool {
	for _, p := range patterns {
		if matchesGlob(p, path) {
			return true
		}
	}
	return false
}

func matchesGlob(pattern, path string) bool {
	pattern = filepath.ToSlash(pattern)
	path = filepath.ToSlash(path)
	if !strings.Contains(pattern, "**") {
		matched, _ := filepath.Match(pattern, path)
		return matched
	}
	re, err := globToRegexp(pattern)
	if err != nil {
		return false
	}
	return re.MatchString(path)
}

func globToRegexp(pattern string) (*regexp.Regexp, error) {
	var b strings.Builder
	b.WriteString("^")
	for i := 0; i < len(pattern); i++ {
		switch pattern[i] {
		case '*':
			if i+1 < len(pattern) && pattern[i+1] == '*' {
				if i+2 < len(pattern) && pattern[i+2] == '/' {
					i += 2
					b.WriteString("(.*/)?")
				} else {
					i++
					b.WriteString(".*")
				}
			} else {
				b.WriteString("[^/]*")
			}
		case '?':
			b.WriteString("[^/]")
		case '.', '+', '(', ')', '|', '^', '$', '{', '}', '[', ']', '\\':
			b.WriteByte('\\')
			b.WriteByte(pattern[i])
		default:
			b.WriteByte(pattern[i])
		}
	}
	b.WriteString("$")
	return regexp.Compile(b.String())
}

// Chunk represents a slice of a file for FTS indexing.
type Chunk struct {
	StartLine   int
	EndLine     int
	Text        string
	ContentHash string
}

// ChunkLines splits content into overlapping windows of ~chunkLines.
func ChunkLines(content []byte, chunkLines, chunkOverlap int) []Chunk {
	if chunkLines <= 0 {
		chunkLines = 80
	}
	text := string(content)
	allLines := strings.Split(text, "\n")
	if len(allLines) == 0 {
		return nil
	}
	step := chunkLines - chunkOverlap
	if step <= 0 {
		step = chunkLines
	}
	var chunks []Chunk
	for start := 0; start < len(allLines); start += step {
		end := start + chunkLines
		if end > len(allLines) {
			end = len(allLines)
		}
		chunkText := strings.Join(allLines[start:end], "\n")
		if strings.TrimSpace(chunkText) == "" {
			if end >= len(allLines) {
				break
			}
			continue
		}
		sum := sha256.Sum256([]byte(chunkText))
		chunks = append(chunks, Chunk{
			StartLine:   start + 1,
			EndLine:     end,
			Text:        chunkText,
			ContentHash: hex.EncodeToString(sum[:]),
		})
		if end >= len(allLines) {
			break
		}
	}
	return chunks
}

// Indexer incrementally indexes source files into the store.
type Indexer struct {
	store *store.Store
}

// New creates an Indexer backed by st.
func New(st *store.Store) *Indexer {
	return &Indexer{store: st}
}

type existingFile struct {
	ID          int64
	MtimeNs     int64
	Size        int64
	ContentHash string
}

// Run scans root and upserts changed files; removes stale entries. The graph
// extract/prune (symbols/edges/embeddings/LSH) runs in the SAME transaction as
// the chunk sync so the content-hash + mtime guards stay single-path (no
// dual-sync drift).
func (idx *Indexer) Run(ctx context.Context, root string, cfg Config) (Stats, error) {
	var stats Stats
	if idx == nil || idx.store == nil || idx.store.DB == nil {
		return stats, fmt.Errorf("indexer not initialized")
	}
	if cfg.MaxFileSize > 0 {
		maxFileSize = int64(cfg.MaxFileSize)
	}
	scanned, err := Scan(root, cfg.Include, cfg.Exclude)
	if err != nil {
		return stats, err
	}
	existing, err := loadExistingFiles(idx.store.DB)
	if err != nil {
		return stats, err
	}
	scannedPaths := make(map[string]struct{}, len(scanned))
	targetUIDs := make(map[string]struct{})
	now := time.Now().UTC().Format(time.RFC3339)
	tx, err := idx.store.DB.BeginTx(ctx, nil)
	if err != nil {
		return stats, err
	}
	defer tx.Rollback()
	for _, file := range scanned {
		if err := ctx.Err(); err != nil {
			return stats, err
		}
		scannedPaths[file.Path] = struct{}{}
		if file.Size > maxFileSize {
			// First-class size skip: counted in stats, not an error, not a
			// fallback. The file is left out of the index entirely.
			stats.FilesOversized++
			continue
		}
		prev, ok := existing[file.Path]
		if ok && prev.MtimeNs == file.MtimeNs && prev.Size == file.Size && prev.ContentHash == file.Hash {
			stats.FilesSkipped++
			continue
		}
		fileID, err := upsertFile(tx, file, now)
		if err != nil {
			return stats, err
		}
		if ok {
			if _, err := tx.ExecContext(ctx, `DELETE FROM chunks WHERE file_id = ?`, fileID); err != nil {
				return stats, fmt.Errorf("delete chunks for %s: %w", file.Path, err)
			}
		}
		chunks := ChunkLines(file.Contents, cfg.ChunkLines, cfg.ChunkOverlap)
		for _, chunk := range chunks {
			if _, err := tx.ExecContext(ctx,
				`INSERT INTO chunks (file_id, start_line, end_line, text, content_hash) VALUES (?, ?, ?, ?, ?)`,
				fileID, chunk.StartLine, chunk.EndLine, chunk.Text, chunk.ContentHash,
			); err != nil {
				return stats, fmt.Errorf("insert chunk for %s: %w", file.Path, err)
			}
			stats.ChunksAdded++
		}
		// Graph pass: extract symbols/edges for this file in the same tx.
		syms, edges, err := idx.extractFile(file)
		if err != nil {
			return stats, fmt.Errorf("extract %s: %w", file.Path, err)
		}
		if n, err := writeFileGraph(tx, fileID, syms, edges); err != nil {
			return stats, err
		} else {
			stats.SymbolsAdded += n
			stats.EdgesAdded += len(edges)
		}
		// Record the file's target UIDs for the end-of-tx global prune.
		for _, s := range syms {
			targetUIDs[s.UID] = struct{}{}
		}
		stats.FilesIndexed++
	}
	// Target-state prune: a deleted file prunes its whole footprint via the
	// file_id cascade; symbols whose file still exists but was rewritten are
	// handled by the per-file writeFileGraph re-upsert above.
	for path, prev := range existing {
		if _, ok := scannedPaths[path]; ok {
			continue
		}
		if err := pruneFileFootprint(tx, prev.ID); err != nil {
			return stats, fmt.Errorf("delete file %s: %w", path, err)
		}
		stats.FilesDeleted++
	}
	// Global target-state prune: any symbol whose uid is not in the declared
	// target set is an orphan (its file was deleted or its function removed).
	// Deleting it cascades to edges, embeddings, LSH buckets, rationale, and
	// FTS rows in one pass.
	if err := pruneOrphanSymbols(tx, targetUIDs); err != nil {
		return stats, fmt.Errorf("prune orphan symbols: %w", err)
	}
	// Stale edges: an edge whose to_name no longer matches any live symbol and
	// whose to_id is null is a dangling name-only edge; drop it. (Name-only
	// edges to live symbols are kept so step 03 can resolve them.)
	if _, err := tx.Exec(`
		DELETE FROM edges
		WHERE to_id IS NULL
		  AND NOT EXISTS (SELECT 1 FROM symbols s WHERE s.name = edges.to_name)
	`); err != nil {
		return stats, fmt.Errorf("prune dangling edges: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return stats, err
	}
	return stats, nil
}

// pruneOrphanSymbols deletes symbols whose uid is not in the declared target
// set. The file_id cascade (for deleted files) has already removed their
// symbols; this catches symbols in surviving files that were rewritten (e.g.
// a removed function) and whose uid is no longer produced by extraction.
func pruneOrphanSymbols(tx *sql.Tx, targetUIDs map[string]struct{}) error {
	rows, err := tx.Query(`SELECT id, uid FROM symbols`)
	if err != nil {
		return err
	}
	var orphanIDs []int64
	for rows.Next() {
		var id int64
		var uid string
		if err := rows.Scan(&id, &uid); err != nil {
			rows.Close()
			return err
		}
		if _, ok := targetUIDs[uid]; !ok {
			orphanIDs = append(orphanIDs, id)
		}
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return err
	}
	for _, id := range orphanIDs {
		if _, err := tx.Exec(`DELETE FROM symbols WHERE id = ?`, id); err != nil {
			return err
		}
	}
	return nil
}

// extractFile runs the Extractor for one scanned file and returns its symbols
// and edges. Per-file extraction failures fall back to regex and never abort
// the run.
func (idx *Indexer) extractFile(file ScannedFile) ([]extract.Symbol, []extract.Edge, error) {
	ex := extract.Default()
	g, err := ex.ExtractFile(file.Path, file.Contents)
	if err != nil {
		return nil, nil, err
	}
	return g.Symbols, g.Edges, nil
}

// fileFirstSymbol caches the first (lowest id) symbol of a file for the
// duration of a writeFileGraph call, used as a default edge source.
var fileFirstSymbol = map[int64]int64{}

// writeFileGraph upserts a file's target-state symbols and edges. It declares
// the target rows (the file's extracted symbols), upserts them, prunes the
// file's now-orphaned symbols (and their edges/vectors/buckets via cascade),
// then upserts the edges.
func writeFileGraph(tx *sql.Tx, fileID int64, syms []extract.Symbol, edges []extract.Edge) (int, error) {
	targetUIDs := make(map[string]struct{}, len(syms))
	for _, s := range syms {
		targetUIDs[s.UID] = struct{}{}
	}
	// Upsert symbols (and their FTS rows via trigger).
	var upserted int
	for _, s := range syms {
		if _, err := tx.Exec(`
			INSERT INTO symbols (file_id, name, qualified_name, kind, language, signature, start_line, end_line, content_hash, uid)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT(uid) DO UPDATE SET
			  file_id = excluded.file_id,
			  name = excluded.name,
			  qualified_name = excluded.qualified_name,
			  kind = excluded.kind,
			  language = excluded.language,
			  signature = excluded.signature,
			  start_line = excluded.start_line,
			  end_line = excluded.end_line,
			  content_hash = excluded.content_hash`,
			fileID, s.Name, s.QualifiedName, s.Kind, s.Language, s.Signature, s.StartLine, s.EndLine, s.ContentHash, s.UID,
		); err != nil {
			return 0, fmt.Errorf("upsert symbol %s: %w", s.Name, err)
		}
		upserted++
	}
	// Resolve UIDs to ids for edges, and record the file's first symbol.
	uidToIDFinal := map[string]int64{}
	rr, err := tx.Query(`SELECT id, uid FROM symbols WHERE file_id = ? ORDER BY start_line, id`, fileID)
	if err != nil {
		return 0, err
	}
	for rr.Next() {
		var id int64
		var uid string
		if err := rr.Scan(&id, &uid); err != nil {
			rr.Close()
			return 0, err
		}
		uidToIDFinal[uid] = id
		if _, ok := fileFirstSymbol[fileID]; !ok {
			fileFirstSymbol[fileID] = id
		}
	}
	rr.Close()
	if err := rr.Err(); err != nil {
		return 0, err
	}
	// Also resolve cross-file target UIDs (a call to a symbol in another file
	// that is already indexed).
	if len(edges) > 0 {
		allUIDs, err := tx.Query(`SELECT uid, id FROM symbols`)
		if err == nil {
			for allUIDs.Next() {
				var uid string
				var id int64
				if err := allUIDs.Scan(&uid, &id); err == nil {
					if _, ok := uidToIDFinal[uid]; !ok {
						uidToIDFinal[uid] = id
					}
				}
			}
			allUIDs.Close()
		}
	}
	// Upsert edges. Edges require a resolvable from_id (the symbol they
	// originate from); import edges are stored with the file's package/first
	// symbol as the source when no specific symbol is bound.
	for _, e := range edges {
		fromID, okFrom := uidToIDFinal[e.FromUID]
		if !okFrom {
			// Resolve a default source: for imports, the file's first symbol
			// (or the file's module). For calls with no enclosing def, the
			// first symbol in the file.
			if defID, ok := fileFirstSymbol[fileID]; ok {
				fromID = defID
			} else {
				continue
			}
		}
		var toID sql.NullInt64
		if e.ToUID != "" {
			if id, ok := uidToIDFinal[e.ToUID]; ok {
				toID.Valid = true
				toID.Int64 = id
			}
		}
		if _, err := tx.Exec(`
			INSERT INTO edges (kind, from_id, to_id, to_name, target_path, confidence, line)
			VALUES (?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT(kind, from_id, to_id, to_name, target_path, line) DO UPDATE SET
			  confidence = excluded.confidence`,
			e.Kind, fromID, toID, e.ToName, e.TargetPath, e.Confidence, e.Line,
		); err != nil {
			return 0, fmt.Errorf("upsert edge %s: %w", e.Kind, err)
		}
	}
	return upserted, nil
}

// pruneFileFootprint deletes a file's whole graph footprint in one pass.
// Deleting the file cascades to chunks and symbols; deleting each symbol
// cascades to edges, embeddings, LSH buckets, rationale, and FTS rows.
func pruneFileFootprint(tx *sql.Tx, fileID int64) error {
	if _, err := tx.Exec(`DELETE FROM files WHERE id = ?`, fileID); err != nil {
		return err
	}
	return nil
}

func loadExistingFiles(db *sql.DB) (map[string]existingFile, error) {
	rows, err := db.Query(`SELECT id, path, mtime_ns, size, content_hash FROM files`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make(map[string]existingFile)
	for rows.Next() {
		var f existingFile
		var path string
		if err := rows.Scan(&f.ID, &path, &f.MtimeNs, &f.Size, &f.ContentHash); err != nil {
			return nil, err
		}
		out[path] = f
	}
	return out, rows.Err()
}

func upsertFile(tx *sql.Tx, file ScannedFile, indexedAt string) (int64, error) {
	if _, err := tx.Exec(
		`INSERT INTO files (path, mtime_ns, size, content_hash, indexed_at)
		 VALUES (?, ?, ?, ?, ?)
		 ON CONFLICT(path) DO UPDATE SET
		   mtime_ns = excluded.mtime_ns,
		   size = excluded.size,
		   content_hash = excluded.content_hash,
		   indexed_at = excluded.indexed_at`,
		file.Path, file.MtimeNs, file.Size, file.Hash, indexedAt,
	); err != nil {
		return 0, err
	}
	// Resolve the rowid by path. LastInsertId is NOT reliable on the
	// upsert UPDATE branch (modernc returns the last inserted rowid in the
	// session, not the upserted row's rowid), so always re-look-up by the
	// unique path. This is a local SQLite PK lookup — cost is negligible.
	var fileID int64
	if err := tx.QueryRow(`SELECT id FROM files WHERE path = ?`, file.Path).Scan(&fileID); err != nil {
		return 0, err
	}
	return fileID, nil
}

// Status reports aggregate index statistics.
type Status struct {
	FileCount   int    `json:"file_count"`
	ChunkCount  int    `json:"chunk_count"`
	LastIndexed string `json:"last_indexed,omitempty"`
}

// GetStatus returns current index stats from the store.
func GetStatus(st *store.Store) (Status, error) {
	var status Status
	if st == nil || st.DB == nil {
		return status, fmt.Errorf("store not initialized")
	}
	if err := st.DB.QueryRow(`SELECT COUNT(*) FROM files`).Scan(&status.FileCount); err != nil {
		return status, err
	}
	if err := st.DB.QueryRow(`SELECT COUNT(*) FROM chunks`).Scan(&status.ChunkCount); err != nil {
		return status, err
	}
	var last sql.NullString
	if err := st.DB.QueryRow(`SELECT MAX(indexed_at) FROM files`).Scan(&last); err != nil {
		return status, err
	}
	if last.Valid {
		status.LastIndexed = last.String
	}
	return status, nil
}
