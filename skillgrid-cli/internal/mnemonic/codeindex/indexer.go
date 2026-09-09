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

	"github.com/devopstales/skillgrid/skillgrid-cli/internal/mnemonic/embedder"
	"github.com/devopstales/skillgrid/skillgrid-cli/internal/mnemonic/extract"
	"github.com/devopstales/skillgrid/skillgrid-cli/internal/mnemonic/memory"
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
	FilesIndexed   int `json:"files_indexed"`
	FilesSkipped   int `json:"files_skipped"`
	FilesDeleted   int `json:"files_deleted"`
	ChunksAdded    int `json:"chunks_added"`
	FilesOversized int `json:"files_oversized"`
	SymbolsAdded   int `json:"symbols_added"`
	EdgesAdded     int `json:"edges_added"`
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
	emb   embedder.Embedder
}

// New creates an Indexer backed by st.
func New(st *store.Store) *Indexer {
	return &Indexer{store: st}
}

// WithEmbedder sets the optional embedder used by the eager dual-granularity
// embedding pass. A nil embedder (the default) skips the pass entirely.
func (idx *Indexer) WithEmbedder(e embedder.Embedder) *Indexer {
	idx.emb = e
	return idx
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
	// skippedFileIDs holds the file IDs of unchanged files that were skipped
	// this run. Their symbols are already correct in the DB, so their UIDs
	// must be added to targetUIDs to prevent the orphan prune from deleting them.
	skippedFileIDs := make([]int64, 0, len(existing))
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
			skippedFileIDs = append(skippedFileIDs, prev.ID)
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
		// Graph pass: extract symbols/edges/rationale for this file in the
		// same tx.
		syms, edges, rationale, err := idx.extractFile(file)
		if err != nil {
			return stats, fmt.Errorf("extract %s: %w", file.Path, err)
		}
		if n, err := writeFileGraph(tx, fileID, syms, edges, rationale); err != nil {
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
	// Seed the orphan-prune target set with UIDs from unchanged (skipped)
	// files. Their symbols are already correct in the DB, but their UIDs never
	// made it into targetUIDs (only re-indexed files contribute above), so
	// without this the orphan prune would delete them.
	if err := seedTargetUIDsFromSkippedFiles(tx, skippedFileIDs, targetUIDs); err != nil {
		return stats, fmt.Errorf("seed target uids: %w", err)
	}
	// Global target-state prune: any symbol whose uid is not in the declared
	// target set is an orphan (its file was deleted or its function removed).
	// Deleting it cascades to edges, embeddings, LSH buckets, rationale, and
	// FTS rows in one pass.
	if err := pruneOrphanSymbols(tx, targetUIDs); err != nil {
		return stats, fmt.Errorf("prune orphan symbols: %w", err)
	}
	// Eager dual-granularity embedding: symbol-level vectors (function/type
	// signatures) upserted in the same tx so the content-hash + mtime guard
	// stays single-path. Batched + resumable: only embeds symbols that are new
	// or whose content changed. A down embedder warns and continues — the FTS
	// floor still works.
	if idx.emb != nil {
		if err := idx.embedPass(ctx, tx); err != nil {
			fmt.Fprintf(os.Stderr, "warn: embed pass: %v\n", err)
		}
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
// seedTargetUIDsFromSkippedFiles adds the UIDs of all symbols belonging to
// skipped (unchanged) files to targetUIDs. This ensures the orphan prune does
// not delete symbols from files that were not re-indexed this run.
func seedTargetUIDsFromSkippedFiles(tx *sql.Tx, fileIDs []int64, targetUIDs map[string]struct{}) error {
	if len(fileIDs) == 0 {
		return nil
	}
	placeholders := make([]string, len(fileIDs))
	args := make([]interface{}, len(fileIDs))
	for i, id := range fileIDs {
		placeholders[i] = "?"
		args[i] = id
	}
	query := fmt.Sprintf(`SELECT uid FROM symbols WHERE file_id IN (%s)`, strings.Join(placeholders, ","))
	rows, err := tx.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var uid string
		if err := rows.Scan(&uid); err != nil {
			return err
		}
		targetUIDs[uid] = struct{}{}
	}
	return rows.Err()
}

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

// embedPass runs the eager dual-granularity embedding pass. It embeds
// symbol-level vectors (name + signature) and upserts only rows that are new
// or whose content changed. The embedding_model guard in embed_meta ensures a
// model swap triggers a full re-embed.
func (idx *Indexer) embedPass(ctx context.Context, tx *sql.Tx) error {
	model := idx.emb.Model()
	dim := idx.emb.Dimension()
	if model == "" || dim <= 0 {
		return nil
	}
	now := time.Now().UTC().Format(time.RFC3339)
	// Model-swap guard: if the indexed model differs, clear all vectors.
	var indexedModel string
	_ = tx.QueryRow(`SELECT value FROM embed_meta WHERE key = 'embedding_model'`).Scan(&indexedModel)
	if indexedModel != "" && indexedModel != model {
		if _, err := tx.Exec(`DELETE FROM embeddings`); err != nil {
			return fmt.Errorf("clear stale embeddings: %w", err)
		}
		if _, err := tx.Exec(`INSERT OR REPLACE INTO embed_meta (key, value) VALUES ('embedding_model', ?)`, model); err != nil {
			return err
		}
	} else if indexedModel == "" {
		if _, err := tx.Exec(`INSERT OR REPLACE INTO embed_meta (key, value) VALUES ('embedding_model', ?)`, model); err != nil {
			return err
		}
	}
	// Symbol-level embedding: embed each symbol's name + signature.
	rows, err := tx.Query(`
		SELECT s.id, s.name, s.signature
		FROM symbols s
		ORDER BY s.id
	`)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		if err := ctx.Err(); err != nil {
			return err
		}
		var symID int64
		var name, sig string
		if err := rows.Scan(&symID, &name, &sig); err != nil {
			return err
		}
		// Skip if already embedded with the current model (model-swap clears
		// all rows above, so an existing row is by definition current-model).
		var existingModel string
		err := tx.QueryRow(`SELECT model FROM embeddings WHERE symbol_id = ?`, symID).Scan(&existingModel)
		if err == nil && existingModel == model {
			continue
		}
		text := name
		if sig != "" {
			text += "\n" + sig
		}
		vec, err := idx.emb.Embed(ctx, text)
		if err != nil {
			continue
		}
		if len(vec.Data) != dim {
			continue
		}
		blob := memory.EncodeVector(vec)
		if _, err := tx.Exec(`
			INSERT INTO embeddings (symbol_id, model, dim, vector, updated_at)
			VALUES (?, ?, ?, ?, ?)
			ON CONFLICT(symbol_id) DO UPDATE SET
			  model = excluded.model,
			  dim = excluded.dim,
			  vector = excluded.vector,
			  updated_at = excluded.updated_at
		`, symID, model, dim, blob, now); err != nil {
			return err
		}
	}
	return rows.Err()
}

// extractFile runs the Extractor for one scanned file and returns its symbols,
// edges, and rationale nodes. Per-file extraction failures fall back to regex
// and never abort the run.
func (idx *Indexer) extractFile(file ScannedFile) ([]extract.Symbol, []extract.Edge, []extract.Rationale, error) {
	ex := extract.Default()
	g, err := ex.ExtractFile(file.Path, file.Contents)
	if err != nil {
		return nil, nil, nil, err
	}
	return g.Symbols, g.Edges, g.Rationales, nil
}

// fileFirstSymbol caches the first (lowest id) symbol of a file for the
// duration of a writeFileGraph call, used as a default edge source.
var fileFirstSymbol = map[int64]int64{}

// writeFileGraph upserts a file's target-state symbols, edges, and rationale.
// It declares the target rows (the file's extracted symbols), upserts them,
// prunes the file's now-orphaned symbols (and their edges/vectors/buckets via
// cascade), then upserts the edges and rationale nodes.
func writeFileGraph(tx *sql.Tx, fileID int64, syms []extract.Symbol, edges []extract.Edge, rationale []extract.Rationale) (int, error) {
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
	// Target-state: drop this file's existing outgoing edges before upserting
	// the freshly extracted set. Pruning by file_id (not by the from-symbol's
	// id) is essential: the edges table has no per-symbol cascade that can
	// reach a file's edges once its symbols are re-keyed on rewrite, so a
	// re-index would otherwise leave stale rows behind.
	if _, err := tx.Exec(`DELETE FROM edges WHERE file_id = ?`, fileID); err != nil {
		return 0, fmt.Errorf("prune edges for file %d: %w", fileID, err)
	}
	// Upsert edges. Edges require a resolvable from_id (the symbol they
	// originate from); import edges are stored with the file's package/first
	// symbol as the source when no specific symbol is bound.
	for _, e := range edges {
		fromID, okFrom := uidToIDFinal[e.FromUID]
		if !okFrom {
			// The extractor binds call/heritage edges to their enclosing
			// definition, but a name-only edge (or one the extractor left
			// unbound) must still resolve to a real source: the file's first
			// symbol. Edges carry a NOT-NULL from_id, so an unresolved source
			// must fall back here rather than silently store a null source.
			defID, ok := fileFirstSymbol[fileID]
			if !ok {
				continue
			}
			fromID = defID
		}
		var toID sql.NullInt64
		if e.ToUID != "" {
			if id, ok := uidToIDFinal[e.ToUID]; ok {
				toID.Valid = true
				toID.Int64 = id
			}
		}
		if _, err := tx.Exec(`
			INSERT INTO edges (kind, from_id, file_id, to_id, to_name, target_path, confidence, line)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)
			ON CONFLICT(kind, from_id, file_id, to_id, to_name, target_path, line) DO UPDATE SET
			  confidence = excluded.confidence`,
			e.Kind, fromID, fileID, toID, e.ToName, e.TargetPath, e.Confidence, e.Line,
		); err != nil {
			return 0, fmt.Errorf("upsert edge %s: %w", e.Kind, err)
		}
	}
	// Upsert rationale nodes. A rationale links to the nearest enclosing
	// symbol (by UID, resolved to an id); a rationale with no resolvable
	// symbol is dropped (no fabricated link).
	if len(rationale) > 0 {
		// Remove this file's existing rationale rows first (target-state: the
		// re-extracted set is the full target; orphans are pruned). Rationale
		// rows reference the file's symbols; symbol id changes on rewrite.
		if _, err := tx.Exec(`
			DELETE FROM rationale WHERE symbol_id IN (SELECT id FROM symbols WHERE file_id = ?)
		`, fileID); err != nil {
			return 0, fmt.Errorf("prune rationale for %d: %w", fileID, err)
		}
		for _, r := range rationale {
			if r.SymbolUID == "" {
				continue
			}
			symID, ok := uidToIDFinal[r.SymbolUID]
			if !ok {
				continue
			}
			if _, err := tx.Exec(`
				INSERT INTO rationale (symbol_id, text, kind, line)
				VALUES (?, ?, ?, ?)
			`, symID, r.Text, r.Kind, r.Line); err != nil {
				return 0, fmt.Errorf("insert rationale: %w", err)
			}
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
