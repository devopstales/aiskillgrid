package knowledge

import (
	"context"
	"database/sql"
	"path/filepath"
	"strings"
)

// FileInput is one scanned file for the knowledge pass: its index-relative
// path and its contents.
type FileInput struct {
	Path     string
	Contents []byte
}

// PassResult summarizes one knowledge pass run.
type PassResult struct {
	Docs   int
	Configs int
	SQL    int
}

// RunPasses runs the three knowledge extractors (doc + config + SQL) over the
// scanned files, persisting doc_nodes / config_nodes / sql_schema_nodes and
// their references / configures / reads / writes edges via st (the indexer's
// open transaction in production, so the pass is in the SAME tx as the 005
// extraction and a single rollback undoes it). Each extractor is non-fatal: a
// malformed file yields zero rows (the index continues, the bad part is
// skipped, the rest is indexed).
func RunPasses(ctx context.Context, st *Store, files []FileInput) (PassResult, error) {
	var res PassResult
	for _, f := range files {
		if err := ctx.Err(); err != nil {
			return res, err
		}
		switch {
		case isMarkdown(f.Path):
			doc := ExtractDoc(f.Path, f.Contents)
			if n, err := st.SaveDoc(ctx, f.Path, doc); err != nil {
				return res, err
			} else {
				res.Docs++
				_ = n
			}
		case isConfig(f.Path):
			cfg := ExtractConfig(f.Path, f.Contents)
			if n, err := st.SaveConfig(ctx, f.Path, cfg); err != nil {
				return res, err
			} else {
				res.Configs++
				_ = n
			}
		case isSQL(f.Path, f.Contents):
			sq := ExtractSQL(f.Path, f.Contents)
			if n, err := st.SaveSchema(ctx, f.Path, sq); err != nil {
				return res, err
			} else {
				res.SQL++
				_ = n
			}
		}
	}
	return res, nil
}

// isSQL reports whether a file should be parsed for SQL: a .sql file, or any
// file whose contents contain a SQL keyword (code-embedded SQL). A .sql file
// is always parsed (even with no DDL/DML, so it is a known schema source); a
// non-.sql file is parsed only when it contains SQL (so the Go extractor's
// code files that embed SQL are covered, and pure code files are not).
func isSQL(path string, contents []byte) bool {
	if filepath.Ext(path) == ".sql" || filepath.Ext(path) == ".SQL" {
		return true
	}
	s := string(contents)
	for _, kw := range []string{"SELECT", "INSERT", "CREATE TABLE", "UPDATE ", "DELETE FROM"} {
		if containsWord(s, kw) {
			return true
		}
	}
	return false
}

// containsWord reports whether s contains the keyword as a whole word
// (case-sensitive, to avoid matching inside identifiers).
func containsWord(s, kw string) bool {
	if kw == "" {
		return false
	}
	i := strings.Index(s, kw)
	for i >= 0 {
		if (i == 0 || !isIdentByte(s[i-1])) && (i+len(kw) == len(s) || !isIdentByte(s[i+len(kw)])) {
			return true
		}
		i = strings.Index(s[i+1:], kw)
		if i >= 0 {
			i += 1
		}
	}
	return false
}

func isIdentByte(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '_'
}

var _ = sql.ErrNoRows
