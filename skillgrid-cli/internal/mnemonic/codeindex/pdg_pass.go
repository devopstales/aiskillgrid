package codeindex

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/devopstales/skillgrid/skillgrid-cli/internal/mnemonic/pdg"
)

// pdgPass runs the opt-in per-function CFG + PDG pass over the indexed graph,
// on the pass DB (a fresh *sql.DB opened after the 005 tx commits — the store's
// single-connection pool cannot share that tx). It is gated on the --pdg flag;
// without it it never runs (the cfg/pdg tables stay empty). Advisory, never
// load-bearing: a per-function failure warns and continues.
func (idx *Indexer) pdgPass(ctx context.Context, passDB *sql.DB, scanned []ScannedFile) error {
	if tableMissing(passDB, "cfg_blocks") {
		return nil // pre-016 store — nothing to do
	}
	// Group files by path for source lookup.
	fileByPath := map[string][]byte{}
	for _, f := range scanned {
		fileByPath[f.Path] = f.Contents
	}
	// Load the function/method symbols (intradependent per function).
	syms, err := pdgFunctions(passDB)
	if err != nil {
		return err
	}
	if len(syms) == 0 {
		return nil
	}
	// Resolve call sites (member calls + their LSP_RESOLVED status) from 005's
	// edges table. The LSP tier (when run first) has already written
	// LSP_RESOLVED edges, so a member call with one is Resolved.
	callsBySymbol, err := pdgCallSites(passDB)
	if err != nil {
		return err
	}
	for _, s := range syms {
		if err := ctx.Err(); err != nil {
			return err
		}
		relPath, err := fileIDToPath(passDB, s.FileID)
		if err != nil {
			continue
		}
		src, ok := fileByPath[relPath]
		if !ok {
			// File not scanned this run (unchanged): read from disk.
			if b, rerr := os.ReadFile(filepath.Join(scanRoot, relPath)); rerr == nil {
				src = b
			} else {
				continue
			}
		}
		c, err := pdg.BuildCFG(src, s.Language, s.Name, s.StartLine)
		if err != nil {
			// Malformed function CFG: skip + continue (01.8).
			fmt.Fprintf(os.Stderr, "warn: pdg cfg %s (%s): %v\n", s.Name, relPath, err)
			continue
		}
		if c == nil {
			continue // body not locatable — skip
		}
		rows, err := pdg.Build(s.ID, c, callsBySymbol[s.ID])
		if err != nil {
			fmt.Fprintf(os.Stderr, "warn: pdg build %s: %v\n", s.Name, err)
			continue
		}
		if err := pdg.PersistBlocks(passDB, s.ID, c); err != nil {
			fmt.Fprintf(os.Stderr, "warn: pdg persist blocks %s: %v\n", s.Name, err)
			continue
		}
		if err := pdg.Persist(passDB, rows); err != nil {
			fmt.Fprintf(os.Stderr, "warn: pdg persist %s: %v\n", s.Name, err)
		}
	}
	return nil
}

// scanRoot is the root the pdg pass reads unchanged files from. It is set by
// Run before the pass; for a cold index it is the scan root, for an incremental
// one the scanned files are already in memory (fileByPath) so this is a
// fallback.
var scanRoot string

// pdgFunctions loads the function/method symbols that a CFG can be built for.
type pdgFunction struct {
	ID        int64
	FileID    int64
	Name      string
	Language  string
	StartLine int
	EndLine   int
}

func pdgFunctions(db *sql.DB) ([]pdgFunction, error) {
	rows, err := db.Query(`
		SELECT s.id, s.file_id, s.name, COALESCE(s.language,'go'), s.start_line, s.end_line
		FROM symbols s
		WHERE s.kind IN ('function','method')
		ORDER BY s.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []pdgFunction
	for rows.Next() {
		var f pdgFunction
		if err := rows.Scan(&f.ID, &f.FileID, &f.Name, &f.Language, &f.StartLine, &f.EndLine); err != nil {
			return nil, err
		}
		out = append(out, f)
	}
	return out, rows.Err()
}

// pdgCallSites resolves, per symbol, the call sites in its body and whether
// each is resolved (a calls edge with a non-NULL to_id, or an LSP_RESOLVED
// edge). Returns symbol id -> []pdg.CallSite.
func pdgCallSites(db *sql.DB) (map[int64][]pdg.CallSite, error) {
	// Join calls edges to their from-symbol; to_id NULL means unresolved.
	rows, err := db.Query(`
		SELECT e.from_id, e.line, COALESCE(s.name,''), e.to_id IS NOT NULL
		FROM edges e JOIN symbols s ON s.id = e.from_id
		WHERE e.kind = 'calls' AND e.to_name IS NOT NULL
		ORDER BY e.from_id, e.line`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[int64][]pdg.CallSite{}
	for rows.Next() {
		var symID int64
		var line int
		var name string
		var resolved bool
		if err := rows.Scan(&symID, &line, &name, &resolved); err != nil {
			return nil, err
		}
		out[symID] = append(out[symID], pdg.CallSite{
			Line:     line,
			Name:     name,
			Resolved: resolved,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	// Mark LSP_RESOLVED boundaries as resolved too (they were AMBIGUOUS in the
	// static pass). The edges table stores LSP_RESOLVED calls with to_id set,
	// so the to_id IS NOT NULL above already covers them.
	return out, nil
}

// fileIDToPath resolves a file id to its path.
func fileIDToPath(db *sql.DB, fileID int64) (string, error) {
	var path string
	err := db.QueryRow(`SELECT path FROM files WHERE id = ?`, fileID).Scan(&path)
	return path, err
}

// lspPass runs the independent opt-in language-server edge tier: it resolves
// member calls the static pass could not type and writes LSP_RESOLVED edges
// into 005's edges table. It is gated on the --lsp flag (independent of --pdg)
// and runs BEFORE the PDG pass so LSP_RESOLVED edges exist before the PDG
// derives dependences. Best-effort: an absent/failing/timing-out server warns
// and continues, leaving the static index unchanged (no partial edge set).
// See pdg/lsp.go for the external-process adapter.
func (idx *Indexer) lspPass(ctx context.Context, passDB *sql.DB, scanned []ScannedFile) error {
	if tableMissing(passDB, "edges") {
		return nil
	}
	client, err := pdg.NewLSPClient(pdg.LSPClientOptions{
		Root: scanRoot,
		Files: func() []pdg.LSPFile {
			var out []pdg.LSPFile
			for _, f := range scanned {
				if lang := pdg.LSPLanguageForPath(f.Path); lang != "" {
					out = append(out, pdg.LSPFile{Path: f.Path, Contents: f.Contents})
				}
			}
			return out
		},
	})
	if err != nil {
		// No resolvable server for the scanned languages: warn + continue,
		// static index unchanged (01.3 / 01.10).
		fmt.Fprintf(os.Stderr, "warn: lsp tier: %v (static index unchanged)\n", err)
		return nil
	}
	defer client.Close()
	if _, err := client.ResolveMemberCalls(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "warn: lsp tier: %v (static index unchanged)\n", err)
		return nil
	}
	return nil
}

// lspPass is a no-op for languages without a known server; the adapter in
// pdg/lsp.go resolves the server binary on PATH and shells out via JSON-RPC.

// (helper) sort is imported for deterministic ordering in call-site maps.
var _ = sort.Strings

// (helper) strings is imported for language-detection helpers.
var _ = strings.TrimSpace
