package route

import (
	"strings"
	"testing"
)

// TestUnresolvedComputedDestinations covers @step-01 (Scenario: Computed or
// unserved destination stays unresolved): a router.push with a computed
// (non-literal) destination, and a navigates edge to a screen no route
// serves, are left unresolved — never fabricated into a stored edge.
func TestUnresolvedComputedDestinations(t *testing.T) {
	// A computed destination (a variable, not a string literal) must NOT
	// produce a navigates edge (it is unresolved, not fabricated).
	src := "function go() {\n  const dest = buildDest();\n  router.push(dest);\n}\n"
	b := Build("a.tsx", []byte(src), fileSyms("go"), known())
	for _, n := range b.Navs {
		t.Errorf("computed destination produced a fabricated navigates edge to %q (want none)", n.ToName)
	}

	// A literal destination is stored, even though no route serves that
	// screen — the screen key is the literal, never an invented match.
	src2 := "function go() {\n  router.push('/nowhere');\n}\n"
	b2 := Build("b.tsx", []byte(src2), fileSyms("go"), known())
	if len(b2.Navs) != 1 || b2.Navs[0].ToName != "/nowhere" {
		t.Errorf("literal destination should be stored as a navigates edge to /nowhere, got %v", navScreens(b2.Navs))
	}
}

// TestUnresolvedWebRouteNoRouteServes covers the web side: a route whose
// handler is a name that no route serves (unresolvable) still yields a route
// node, but no references edge is fabricated (drop-not-guess).
func TestUnresolvedWebRouteNoRouteServes(t *testing.T) {
	// `unserved` is a bare identifier (a reference) but no route serves it —
	// it is unresolvable, so the references edge is dropped, not fabricated.
	src := "app.get('/mystery', unserved);\n"
	b := Build("server.js", []byte(src), fileSyms("other"), dummyIndex{known: map[string]int64{}})
	found := false
	for _, n := range b.Nodes {
		if n.PathPattern == "/mystery" {
			found = true
			if n.HandlerSymbol != 0 {
				t.Errorf("unresolvable handler should be unresolved (HandlerSymbol=0), got %d", n.HandlerSymbol)
			}
		}
	}
	if !found {
		t.Errorf("expected a route node for /mystery, got none")
	}
	if b.Dropped != 1 {
		t.Errorf("expected the unresolvable handler ref to be dropped (warning), got %d drops", b.Dropped)
	}
}

// TestDropNotGuess covers @step-01 (Scenario: Ambiguous references are dropped
// not guessed): a reference with no same-file match, no explicit specifier,
// and no unique global/owner-qualified match is dropped at extraction (no edge
// stored) and reported as a warning (count + sample), not a silent discard.
// AMBIGUOUS is reserved for edges that WERE resolved via a heuristic (a
// markup-written navigates link) — such an edge is stored INFERRED/AMBIGUOUS,
// while an unresolvable one is absent.
func TestDropNotGuess(t *testing.T) {
	// A route whose handler has no same-file match and no unique global match
	// (the dummy index reports it ambiguous) → dropped, with a warning.
	src := "app.get('/x', someAmbiguousHandler);\n"
	idx := dummyIndex{known: map[string]int64{}} // nothing known → unresolvable
	b := Build("server.js", []byte(src), fileSyms("other"), idx)
	if b.Dropped != 1 {
		t.Fatalf("expected 1 dropped reference (drop-not-guess), got %d", b.Dropped)
	}
	if b.DropSample != "someAmbiguousHandler" {
		t.Errorf("drop warning sample = %q, want someAmbiguousHandler", b.DropSample)
	}
	// The route node is present but has no references edge (HandlerSymbol=0).
	for _, n := range b.Nodes {
		if n.PathPattern == "/x" {
			if n.HandlerSymbol != 0 {
				t.Errorf("ambiguous handler reference should be dropped (HandlerSymbol=0), got %d", n.HandlerSymbol)
			}
		}
	}

	// A resolvable reference (unique global match) → stored, NOT dropped.
	src2 := "app.get('/y', knownHandler);\n"
	b2 := Build("server2.js", []byte(src2), fileSyms("other"), known("knownHandler"))
	if b2.Dropped != 0 {
		t.Errorf("resolvable handler should not be dropped, got %d drops", b2.Dropped)
	}
	for _, n := range b2.Nodes {
		if n.PathPattern == "/y" && n.HandlerSymbol == 0 {
			t.Errorf("resolvable handler reference should be stored (HandlerSymbol!=0)")
		}
	}
}

// TestDropNotGuessWarningIsReported covers the "reported as a warning, not a
// silent discard" half of the policy: the drop count + a sample are surfaced.
func TestDropNotGuessWarningIsReported(t *testing.T) {
	src := "app.get('/a', h1);\napp.get('/b', h2);\n"
	b := Build("server.js", []byte(src), fileSyms("other"), dummyIndex{known: map[string]int64{}})
	if b.Dropped != 2 {
		t.Errorf("expected 2 dropped refs, got %d", b.Dropped)
	}
	if b.DropSample == "" {
		t.Errorf("expected a non-empty drop sample (warning), got empty")
	}
	if !strings.Contains(b.DropSample, "h") {
		t.Errorf("drop sample %q should name a dropped handler", b.DropSample)
	}
}

// TestMalformedFileFallsBackAndContinues covers @step-01 (Scenario: Malformed
// routing file falls back and index continues): a malformed routing file
// yields an empty (non-nil) result, never an error, so the index continues.
func TestMalformedFileFallsBackAndContinues(t *testing.T) {
	// Malformed python: a broken urls.py that is not valid Python.
	badPy := "urlpatterns = [\npath('/x/', views.Home.as_view(,\n"
	fr := ExtractFile("urls.py", []byte(badPy))
	if fr == nil {
		t.Fatalf("ExtractFile must not return nil for a malformed file")
	}
	// The extractor degrades (no routes) but does not panic / error.

	// An empty file yields an empty result, not nil.
	empty := ExtractFile("empty.js", []byte(""))
	if empty == nil {
		t.Fatalf("ExtractFile must not return nil for an empty file")
	}

	// A nil-ish (all whitespace) file also degrades.
	blank := ExtractFile("blank.ts", []byte("\n\n  \n"))
	if blank == nil {
		t.Fatalf("ExtractFile must not return nil for a blank file")
	}

	// Build with a malformed src still returns a non-nil Built (no route
	// nodes, no navigates) so the indexer's per-file pass continues.
	b := Build("weird.py", []byte("::::\n"), fileSyms("x"), known())
	if b == nil {
		t.Fatalf("Build must not return nil for a malformed file")
	}
}

// TestNoRoutesFrameworkReportsNone covers @step-01 (Scenario: Framework with
// no routes reports none): a recognized framework file with no route
// definitions yields zero route nodes and zero references edges (it reports
// "none", not a fabricated route).
func TestNoRoutesFrameworkReportsNone(t *testing.T) {
	// A Django urls.py with an empty urlpatterns list.
	src := "from . import views\n\nurlpatterns = []\n"
	b := Build("urls.py", []byte(src), fileSyms("views"), known())
	if len(b.Nodes) != 0 {
		t.Errorf("empty urlpatterns should report no routes, got %d nodes", len(b.Nodes))
	}
	// A Spring controller with no @*Mapping annotations.
	src2 := "package c;\n\npublic class Empty {\n  public void nothing() {}\n}\n"
	b2 := Build("Empty.java", []byte(src2), fileSyms("nothing"), known())
	if len(b2.Nodes) != 0 {
		t.Errorf("controller without mappings should report no routes, got %d nodes", len(b2.Nodes))
	}
}
