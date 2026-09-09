// Package knowledge extracts the code's surrounding knowledge — markdown
// docs, config files, and SQL DDL — into the 005 graph as first-class nodes
// (doc_nodes / config_nodes / sql_schema_nodes) with confidence-labeled
// edges (references / configures / reads / writes). It is additive: it never
// rewrites a 005 symbol or edge, and a malformed file yields zero knowledge
// rows (the index continues, the bad part is skipped, the rest is indexed).
package knowledge

import (
	"hash/fnv"
	"path/filepath"
	"regexp"
	"strings"
)

// Edge kinds emitted by the knowledge extractors (distinct from 005's
// call/heritage kinds; references is shared with 010's route references but
// scoped here to doc node -> doc node).
const (
	KindReferences = "references"
	KindConfigures = "configures"
	KindReads      = "reads"
	KindWrites     = "writes"
)

// Confidence labels (mirrored from extract/route/graph).
//
//	EXTRACTED  — explicit syntax (a literal doc link, a config key naming a
//	            symbol by explicit spec, an SQL statement naming a table).
//	INFERRED   — convention-derived (a wikilink resolved by file convention).
//	AMBIGUOUS  — a resolved-but-inferred guess (an unresolved config reference
//	            kept, not dropped, per 03.5).
const (
	ConfidenceExtracted = "EXTRACTED"
	ConfidenceInferred  = "INFERRED"
	ConfidenceAmbiguous = "AMBIGUOUS"
)

// KindDoc / KindConfig / KindTable / KindColumn are the knowledge node kinds.
const (
	KindDoc    = "doc"
	KindConfig = "config"
	KindTable  = "table"
	KindColumn = "column"
)

// mdLink matches a markdown link [text](target). The target is a relative
// path (or a URL fragment); anchors (#...) are stripped before resolution.
var mdLink = regexp.MustCompile(`\[[^\]]*\]\(([^)\s]+)\)`)

// mdWiki matches a wikilink [[target]].
var mdWiki = regexp.MustCompile(`\[\[([^\[\]]+)\]\]`)

// DocLink is one markdown link / wikilink in a doc file: its target as
// written, the line it appears on, and its confidence.
type DocLink struct {
	Target     string
	Line       int
	Confidence string
}

// DocResult is the extraction result for one markdown file: the doc node
// (its title + path) and its outgoing references links.
type DocResult struct {
	Path   string
	Title  string
	Links  []DocLink
	IsDoc  bool
}

// DocTitle derives the doc node's title from the first H1 heading, falling
// back to the file name (without extension).
func DocTitle(path string, src []byte) string {
	for _, line := range strings.Split(string(src), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "# ") {
			return strings.TrimSpace(strings.TrimPrefix(trimmed, "# "))
		}
	}
	return strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
}

// ExtractDoc parses one markdown file into doc node + references links.
// It never errors and never returns nil: a malformed file (unparseable link)
// skips the bad link and keeps the rest (03.6). A non-markdown path returns
// an empty (IsDoc=false) result.
func ExtractDoc(path string, src []byte) *DocResult {
	if !isMarkdown(path) {
		return &DocResult{Path: path, Title: filepath.Base(path)}
	}
	res := &DocResult{Path: path, Title: DocTitle(path, src), IsDoc: true}
	for lineNo, line := range strings.Split(string(src), "\n") {
		for _, m := range mdLink.FindAllStringSubmatchIndex(line, -1) {
			target := line[m[2]:m[3]]
			if !isDocLinkTarget(target) {
				continue
			}
			res.Links = append(res.Links, DocLink{
				Target:     normalizeLinkTarget(target),
				Line:       lineNo + 1,
				Confidence: ConfidenceExtracted,
			})
		}
		for _, m := range mdWiki.FindAllStringSubmatchIndex(line, -1) {
			target := line[m[2]:m[3]]
			res.Links = append(res.Links, DocLink{
				Target:     normalizeWikiTarget(target),
				Line:       lineNo + 1,
				Confidence: ConfidenceInferred,
			})
		}
	}
	return res
}

// isMarkdown reports whether path is a markdown file.
func isMarkdown(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".md", ".markdown":
		return true
	}
	return false
}

// isDocLinkTarget reports whether a markdown link target is a doc reference
// (not an external URL, a bare anchor, or an image) — only doc links become
// references edges (03.2: links between doc nodes).
func isDocLinkTarget(target string) bool {
	if target == "" {
		return false
	}
	if strings.Contains(target, "://") {
		return false // external URL
	}
	if strings.HasPrefix(target, "#") {
		return false // in-page anchor, not a doc reference
	}
	return strings.HasSuffix(target, ".md") || strings.HasSuffix(target, ".markdown")
}

// normalizeLinkTarget strips an in-page anchor (#...) from a markdown link
// target so the reference resolves to the doc, not the fragment.
func normalizeLinkTarget(target string) string {
	if i := strings.Index(target, "#"); i >= 0 {
		target = target[:i]
	}
	return strings.TrimSpace(target)
}

// normalizeWikiTarget strips an alias ([[page|alias]] -> page) and an anchor
// from a wikilink target. Wikilinks name a doc by title/path; a bare name
// (no path separator, no .md suffix) is resolved to <name>.md by file
// convention so it links to the same-named doc node (03.2: wikilinks between
// doc nodes).
func normalizeWikiTarget(target string) string {
	if i := strings.Index(target, "|"); i >= 0 {
		target = target[:i]
	}
	if i := strings.Index(target, "#"); i >= 0 {
		target = target[:i]
	}
	target = strings.TrimSpace(target)
	if target == "" {
		return target
	}
	if strings.Contains(target, "/") || strings.Contains(target, "\\") {
		if !strings.HasSuffix(target, ".md") && !strings.HasSuffix(target, ".markdown") {
			target += ".md"
		}
		return target
	}
	if !strings.HasSuffix(target, ".md") && !strings.HasSuffix(target, ".markdown") {
		target += ".md"
	}
	return target
}

// docNodeID derives the doc node's deterministic id-key (the file path,
// normalized to slashes) — a doc node is per-file.
func docNodeKey(path string) string {
	return filepath.ToSlash(path)
}

// fnvHex returns the hex FNV-1a digest of s (pure Go, CGo-free) — used for
// deterministic knowledge node uids.
func fnvHex(s string) string {
	h := fnv.New64a()
	_, _ = h.Write([]byte(s))
	return strings.ToLower(strings.TrimLeft(
		func() string { return hex64(h.Sum64()) }(), "0"))
}

func hex64(v uint64) string {
	const digits = "0123456789abcdef"
	var b [16]byte
	for i := 15; i >= 0; i-- {
		b[i] = digits[v&0xf]
		v >>= 4
	}
	return string(b[:])
}
