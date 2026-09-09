package route

import (
	"sort"
	"strings"
)

// RouteNode is one route node to be upserted into the 005 symbols table.
type RouteNode struct {
	Name          string
	UID           string
	PathPattern   string
	Method        string
	Framework     string
	Language      string
	Line          int
	ContentHash   string
	HandlerName   string
	HandlerSymbol int64 // 0 when the handler could not be resolved
	HandlerUID    string
	HandlerExplicit bool
}

// NavigationNode is one navigation edge to be upserted into the 005 edges
// table (kind = navigates).
type NavigationNode struct {
	FromSymbol int64
	ToName     string // the screen (path-like)
	Confidence string
	Line       int
}

// Build is the per-framework node map. Given the raw extraction for a file,
// the symbols extracted by 005 (for the same file), and the symbol index, it
// produces the route nodes + references/navigates edges. This is the single
// place where the drop-rather-than-guess policy is applied to references.
func Build(path string, src []byte, fileSyms []FileSymbol, index SymbolIndex) *Built {
	fr := ExtractFile(path, src)
	if fr == nil {
		fr = &FileRoutes{}
	}
	sortRoutes(fr)

	built := &Built{
		Framework: fr.Framework,
		Nodes:     make([]RouteNode, 0, len(fr.Routes)),
		Navs:      make([]NavigationNode, 0, len(fr.Navigates)),
	}
	built.buildRoutes(path, src, fr, fileSyms, index)
	built.buildNavigates(path, src, fr, fileSyms, index)
	return built
}

// Built is the output of Build: the route nodes and navigates edges to store,
// plus the drop-not-guess warnings.
type Built struct {
	Framework string
	Nodes     []RouteNode
	Navs      []NavigationNode
	Dropped   int
	DropSample string
}

// FileSymbol is a minimal view of one 005-extracted symbol in the file, used
// to resolve route handlers. It is intentionally narrow so the route package
// does not import the extract package (no cycle).
type FileSymbol struct {
	Name string
	UID  string
	ID   int64
}

// SymbolIndex resolves handler names to their symbol IDs. It is backed by the
// 005 symbols table: same-file lookup first, then a unique global
// owner-qualified / bare-name match. Multiple matches are ambiguous and yield
// (0, "") so the caller can drop the reference.
type SymbolIndex interface {
	// ResolveHandler returns the symbol ID + UID for a handler name. The
	// second return is false when the reference is ambiguous (multiple global
	// matches) or unresolvable (no match) — the drop-not-guess case.
	ResolveHandler(name string) (id int64, uid string, ok bool)
}

// buildRoutes turns each raw Route into a RouteNode + (optionally) a
// references edge, applying the drop-not-guess policy to the handler
// reference.
func (b *Built) buildRoutes(path string, src []byte, fr *FileRoutes, fileSyms []FileSymbol, index SymbolIndex) {
	for _, r := range fr.Routes {
		name := routeNodeName(r.Framework, r.PathPattern, r.Method)
		uid := symbolUID(r.Framework, r.PathPattern, r.Method)
		node := RouteNode{
			Name:          name,
			UID:           uid,
			PathPattern:   r.PathPattern,
			Method:        r.Method,
			Framework:     r.Framework,
			Language:      strings.TrimPrefix(filepathExt(path), "."),
			Line:          r.Line,
			ContentHash:   uid,
			HandlerName:   r.Handler,
			HandlerExplicit: r.Explicit,
		}
		if node.Language == "" {
			node.Language = "unknown"
		}
		// Resolve the handler reference. The route node is always created (it
		// is the URL endpoint); the references edge is dropped when the
		// handler reference is ambiguous (drop-not-guess).
		if r.Handler != "" {
			id, uid2, ok := index.ResolveHandler(r.Handler)
			if !ok {
				// Drop-not-guess: no same-file match, no unique global match.
				b.recordDrop(r.Handler)
			} else {
				node.HandlerSymbol = id
				node.HandlerUID = uid2
			}
		}
		b.Nodes = append(b.Nodes, node)
	}
}

// buildNavigates turns each raw Navigation into a navigates edge. Literal,
// programmatic destinations are EXTRACTED; markup-written links are INFERRED.
// A destination that is computed (not a literal) is left unresolved — it is
// never stored with a fabricated screen, so it is dropped from Navs.
func (b *Built) buildNavigates(path string, src []byte, fr *FileRoutes, fileSyms []FileSymbol, index SymbolIndex) {
	// The sending function is the nearest enclosing function symbol; when the
	// navigation is at top level (module scope) the source is the file's
	// first symbol.
	var fromID int64
	if len(fileSyms) > 0 {
		fromID = fileSyms[0].ID
	}
	for _, n := range fr.Navigates {
		if n.Dest == "" {
			// Computed / no-literal destination: unresolved, not fabricated.
			continue
		}
		conf := ConfidenceExtracted
		if n.Markup {
			conf = ConfidenceInferred
		}
		b.Navs = append(b.Navs, NavigationNode{
			FromSymbol: fromID,
			ToName:     n.Screen,
			Confidence: conf,
			Line:       n.Line,
		})
	}
}

// recordDrop tallies a dropped reference for the warning.
func (b *Built) recordDrop(name string) {
	b.Dropped++
	if b.DropSample == "" {
		b.DropSample = name
	}
}

// routeNodeName derives a readable route node name (method + pattern).
func routeNodeName(framework, pattern, method string) string {
	if method == "" {
		return "route:" + pattern
	}
	return method + " " + pattern
}

// filepathExt returns the lowercased extension (without dot) of path.
func filepathExt(path string) string {
	idx := strings.LastIndex(path, ".")
	if idx < 0 {
		return ""
	}
	return strings.ToLower(path[idx+1:])
}

// SortedHandlers returns the sorted unique handler names referenced by a Built
// set (for deterministic drop samples).
func (b *Built) SortedHandlers() []string {
	seen := map[string]bool{}
	var out []string
	for _, n := range b.Nodes {
		if n.HandlerName != "" && !seen[n.HandlerName] {
			seen[n.HandlerName] = true
			out = append(out, n.HandlerName)
		}
	}
	sort.Strings(out)
	return out
}
