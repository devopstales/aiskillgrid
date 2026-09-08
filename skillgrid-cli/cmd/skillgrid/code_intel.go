package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/devopstales/skillgrid/skillgrid-cli/internal/mnemonic/graph"
	"github.com/devopstales/skillgrid/skillgrid-cli/internal/mnemonic/service"
)

// runCodeIntel handles the `skillgrid orient` and `skillgrid grep` subcommands
// (CLI parity for the Tier-1 orientation and structural code_grep MCP tools).
func runCodeIntel(version string, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: skillgrid <orient|grep> ... (see --help)")
		os.Exit(2)
	}
	switch args[0] {
	case "orient":
		runOrient(version, args[1:])
	case "grep":
		runGrep(version, args[1:])
	case "callers", "callees", "dependents", "implementors", "hierarchy", "tests-for":
		runNeighbors(version, args[0], args[1:])
	case "path":
		runGraphPath(version, args[1:])
	case "explain":
		runExplain(version, args[1:])
	case "impact":
		runImpact(version, args[1:])
	case "explore":
		runExplore(version, args[1:])
	case "help", "-h", "--help":
		printCodeIntelUsage()
	default:
		fmt.Fprintf(os.Stderr, "error: unknown code_intel subcommand %q (see --help)\n", args[0])
		os.Exit(2)
	}
}

func printCodeIntelUsage() {
	fmt.Fprintln(os.Stderr, `usage: skillgrid <orient|grep> [flags]

  skillgrid orient SYMBOL  Tier-1 orientation for a symbol: signature, file TOC,
                           map, list, metadata, and linked rationale.
  skillgrid grep PATTERN [path]  Structural by-example grep (index-free):
                           matches the gotreesitter syntax tree, per language.

  callers|callees|dependents|implementors|hierarchy|tests-for SYMBOL
                           Graph neighbors (every edge confidence-labeled)
  path FROM TO            Shortest edge path, or where the graph stops
  explain SYMBOL          Symbol node + degree + connections ranked by degree
  impact SYMBOL           Risk-tiered blast radius (WILL BREAK / LIKELY AFFECTED)
  explore SYMBOL          Composite: source + call-flow + blast radius in one call

  graph flags:
    --json    Emit machine-readable JSON (default: human table)
  impact:
    --file F --uid U --kind K --min-confidence C --max-depth N
  grep:
    --json    Emit machine-readable JSON (default: human table)`)
}

// graphView maps a CLI subcommand to its graph view.
func graphView(sub string) string {
	views := map[string]string{
		"callers":      "callers",
		"callees":      "callees",
		"dependents":   "dependents",
		"implementors": "implementors",
		"hierarchy":    "hierarchy",
		"tests-for":    "tests_for",
	}
	if v, ok := views[sub]; ok {
		return v
	}
	return sub
}

// openGraphService resolves the project for the current directory.
func openGraphService() (*service.Service, string, error) {
	dataDir := envOr("SKILLGRID_MNEMONIC_DATA_DIR", "")
	svc, err := newMnemonicService(dataDir)
	if err != nil {
		return nil, "", err
	}
	projectID, err := svc.ResolveProject(".")
	if err != nil {
		return nil, "", err
	}
	return svc, projectID, nil
}

// printImpact renders a risk-tiered blast-radius result (or a candidate list).
func printImpact(out *service.ImpactResultDTO) {
	if out.Ambiguous {
		fmt.Printf("ambiguous target: %d candidates (narrow with --file/--uid/--kind)\n", len(out.Candidates))
		for i, c := range out.Candidates {
			fmt.Printf("  %d. %s (%s) %s:%d [uid %s]\n", i+1, c.Name, c.Kind, c.Path, c.StartLine, c.UID)
		}
		return
	}
	if out.Target == nil {
		fmt.Println("not found: no matching symbol")
		return
	}
	fmt.Printf("%s (%s) %s:%d\n", out.Target.Name, out.Target.Kind, out.Target.Path, out.Target.StartLine)
	fmt.Printf("WILL BREAK (%d):\n", len(out.WillBreak))
	for _, e := range out.WillBreak {
		fmt.Printf("  %s (%s) %s:%d [%s %s] depth %d\n", e.Symbol.Name, e.Symbol.Kind, e.Symbol.Path, e.Symbol.StartLine, e.Kind, e.Confidence, e.Depth)
	}
	fmt.Printf("LIKELY AFFECTED (%d):\n", len(out.Likely))
	for _, e := range out.Likely {
		fmt.Printf("  %s (%s) %s:%d [%s %s] depth %d\n", e.Symbol.Name, e.Symbol.Kind, e.Symbol.Path, e.Symbol.StartLine, e.Kind, e.Confidence, e.Depth)
	}
	if out.Excluded > 0 {
		fmt.Printf("excluded low-confidence hops: %d\n", out.Excluded)
	}
}

// printGraphStops renders a where-the-graph-stops answer.
func printGraphStops(g *graph.GraphStops) {
	fmt.Printf("where the graph stops: %s at line %d\n", g.DispatchKind, g.Line)
	for _, r := range g.Refused {
		if r.Path != "" {
			fmt.Printf("  refused: %s [%s] at %s:%d\n", r.Name, r.Confidence, r.Path, r.Line)
		} else {
			fmt.Printf("  refused: %s [%s] at line %d\n", r.Name, r.Confidence, r.Line)
		}
	}
}

// runOrient implements `skillgrid orient SYMBOL`.
func runOrient(version string, args []string) {
	fs := flag.NewFlagSet("orient", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	var jsonOut bool
	fs.BoolVar(&jsonOut, "json", false, "emit machine-readable JSON")
	fs.Usage = func() {
		fmt.Fprintln(fs.Output(), "usage: skillgrid orient SYMBOL [--json]")
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}
	if fs.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "error: orient requires exactly one SYMBOL argument")
		os.Exit(2)
	}
	symbol := fs.Arg(0)
	dataDir := envOr("SKILLGRID_MNEMONIC_DATA_DIR", "")
	svc, err := newMnemonicService(dataDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	_ = version
	projectID, err := svc.ResolveProject(".")
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	out, err := svc.OrientSymbol(ctx, projectID, symbol)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	if jsonOut {
		b, _ := json.MarshalIndent(out, "", "  ")
		fmt.Fprintln(os.Stdout, string(b))
		return
	}
	if !out.Found {
		fmt.Fprintf(os.Stderr, "not found: %s\n", out.Reason)
		return
	}
	s := out.Symbol
	fmt.Printf("%s (%s) %s:%d-%d\n", s["name"], s["kind"], s["path"], s["start_line"], s["end_line"])
	if out.Signature != "" {
		fmt.Printf("  %s\n", out.Signature)
	}
	fmt.Printf("file TOC (%d symbols):\n", len(out.FileTOC))
	for _, e := range out.FileTOC {
		mark := "  "
		if e["name"] == s["name"] {
			mark = "  * "
		}
		fmt.Printf("%s%s (%s) :%d-%d\n", mark, e["name"], e["kind"], e["start_line"], e["end_line"])
	}
	if len(out.Rationale) > 0 {
		fmt.Printf("rationale:\n")
		for _, r := range out.Rationale {
			fmt.Printf("  [%s] %s:%d %v\n", r["kind"], s["path"], r["line"], r["text"])
		}
	}
}

// runGrep implements `skillgrid grep PATTERN [path]`.
func runGrep(version string, args []string) {
	fs := flag.NewFlagSet("grep", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	var jsonOut bool
	fs.BoolVar(&jsonOut, "json", false, "emit machine-readable JSON")
	fs.Usage = func() {
		fmt.Fprintln(fs.Output(), "usage: skillgrid grep PATTERN [path] [--json]")
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}
	if fs.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "error: grep requires a PATTERN argument")
		os.Exit(2)
	}
	pattern := fs.Arg(0)
	path := "."
	if fs.NArg() >= 2 {
		path = fs.Arg(1)
	}
	dataDir := envOr("SKILLGRID_MNEMONIC_DATA_DIR", "")
	svc, err := newMnemonicService(dataDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	_ = version
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	res, err := svc.CodeGrep(ctx, path, pattern)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	if jsonOut {
		b, _ := json.MarshalIndent(res, "", "  ")
		fmt.Fprintln(os.Stdout, string(b))
		return
	}
	for _, n := range res.Notes {
		fmt.Fprintf(os.Stderr, "note: skipping %s files: %s\n", n.Language, n.Reason)
	}
	if len(res.Hits) == 0 && len(res.Notes) == 0 {
		fmt.Fprintln(os.Stderr, "no matches")
		return
	}
	for _, h := range res.Hits {
		fmt.Printf("%s:%d:%d  %s\n", h.Path, h.Line, h.Col, h.Text)
	}
}

// runNeighbors implements the graph neighbor subcommands (CLI parity for the
// Tier-2 code_get_* MCP tools).
func runNeighbors(version, sub string, args []string) {
	fs := flag.NewFlagSet(sub, flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	var jsonOut bool
	fs.BoolVar(&jsonOut, "json", false, "emit machine-readable JSON")
	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "usage: skillgrid %s SYMBOL [--json]\n", sub)
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}
	if fs.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "error: %s requires exactly one SYMBOL argument\n", sub)
		os.Exit(2)
	}
	symbol := fs.Arg(0)
	svc, projectID, err := openGraphService()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	_ = version
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	out, err := svc.GrabNeighbors(ctx, projectID, graphView(sub), symbol)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	if jsonOut {
		b, _ := json.MarshalIndent(out, "", "  ")
		fmt.Fprintln(os.Stdout, string(b))
		return
	}
	if out.Reason != "" {
		fmt.Fprintln(os.Stderr, "note:", out.Reason)
	}
	if len(out.Edges) == 0 {
		fmt.Fprintln(os.Stderr, "no edges")
		return
	}
	for _, e := range out.Edges {
		to := e.ToName
		if e.To.ID != 0 {
			to = e.To.Name
		}
		fmt.Printf("%s -> %s [%s %s] :%d\n", e.From.Name, to, e.Kind, e.Confidence, e.Line)
	}
}

// runGraphPath implements `skillgrid path FROM TO`.
func runGraphPath(version string, args []string) {
	fs := flag.NewFlagSet("path", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	var jsonOut bool
	fs.BoolVar(&jsonOut, "json", false, "emit machine-readable JSON")
	fs.Usage = func() {
		fmt.Fprintln(fs.Output(), "usage: skillgrid path FROM TO [--json]")
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}
	if fs.NArg() != 2 {
		fmt.Fprintln(os.Stderr, "error: path requires FROM and TO arguments")
		os.Exit(2)
	}
	from, to := fs.Arg(0), fs.Arg(1)
	svc, projectID, err := openGraphService()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	_ = version
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	out, err := svc.CodePath(ctx, projectID, from, to)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	if jsonOut {
		b, _ := json.MarshalIndent(out, "", "  ")
		fmt.Fprintln(os.Stdout, string(b))
		return
	}
	if out.Reason != "" {
		fmt.Fprintln(os.Stderr, "note:", out.Reason)
		return
	}
	if !out.Found {
		if out.GraphStops != nil {
			printGraphStops(out.GraphStops)
		} else {
			fmt.Fprintln(os.Stderr, "no path")
		}
		return
	}
	if len(out.Path) == 0 {
		fmt.Println(from, "==", to)
		return
	}
	prev := from
	for _, e := range out.Path {
		next := e.ToName
		if e.To.ID != 0 {
			next = e.To.Name
		}
		fmt.Printf("%s -> %s [%s %s] :%d\n", prev, next, e.Kind, e.Confidence, e.Line)
		prev = next
	}
}

// runExplain implements `skillgrid explain SYMBOL`.
func runExplain(version string, args []string) {
	fs := flag.NewFlagSet("explain", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	var jsonOut bool
	fs.BoolVar(&jsonOut, "json", false, "emit machine-readable JSON")
	fs.Usage = func() {
		fmt.Fprintln(fs.Output(), "usage: skillgrid explain SYMBOL [--json]")
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}
	if fs.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "error: explain requires exactly one SYMBOL argument")
		os.Exit(2)
	}
	symbol := fs.Arg(0)
	svc, projectID, err := openGraphService()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	_ = version
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	out, err := svc.CodeExplain(ctx, projectID, symbol)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	if jsonOut {
		b, _ := json.MarshalIndent(out, "", "  ")
		fmt.Fprintln(os.Stdout, string(b))
		return
	}
	if !out.Found {
		fmt.Fprintln(os.Stderr, "not found:", out.Reason)
		return
	}
	fmt.Printf("%s degree %d\n", out.Symbol.Name, out.Degree)
	for _, c := range out.Connections {
		name := c.Symbol.Name
		if name == "" {
			name = c.ToName
		}
		fmt.Printf("  %s [%s %s] :%d (degree %d)\n", name, c.Kind, c.Confidence, c.Line, c.Degree)
	}
}

// runImpact implements `skillgrid impact SYMBOL`.
func runImpact(version string, args []string) {
	fs := flag.NewFlagSet("impact", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	var jsonOut bool
	var file, uid, kind, minConf string
	var maxDepth int
	fs.BoolVar(&jsonOut, "json", false, "emit machine-readable JSON")
	fs.StringVar(&file, "file", "", "narrow to a file path")
	fs.StringVar(&uid, "uid", "", "narrow to a symbol UID")
	fs.StringVar(&kind, "kind", "", "narrow to a symbol kind")
	fs.StringVar(&minConf, "min-confidence", "", "minimum edge confidence (EXTRACTED|INFERRED|AMBIGUOUS)")
	fs.IntVar(&maxDepth, "max-depth", 0, "bound the traversal depth")
	fs.Usage = func() {
		fmt.Fprintln(fs.Output(), "usage: skillgrid impact SYMBOL [--file F] [--uid U] [--kind K] [--min-confidence C] [--max-depth N] [--json]")
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}
	if fs.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "error: impact requires exactly one SYMBOL argument")
		os.Exit(2)
	}
	symbol := fs.Arg(0)
	svc, projectID, err := openGraphService()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	_ = version
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	out, err := svc.CodeImpact(ctx, projectID, symbol, service.ImpactOptions{
		File: file, UID: uid, Kind: kind, MinConfidence: minConf, MaxDepth: maxDepth,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	if jsonOut {
		b, _ := json.MarshalIndent(out, "", "  ")
		fmt.Fprintln(os.Stdout, string(b))
		return
	}
	printImpact(out)
}

// runExplore implements `skillgrid explore SYMBOL` (CLI parity for the
// composite code_explore).
func runExplore(version string, args []string) {
	fs := flag.NewFlagSet("explore", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	var jsonOut bool
	fs.BoolVar(&jsonOut, "json", false, "emit machine-readable JSON")
	fs.Usage = func() {
		fmt.Fprintln(fs.Output(), "usage: skillgrid explore SYMBOL [--json]")
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}
	if fs.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "error: explore requires exactly one SYMBOL argument")
		os.Exit(2)
	}
	symbol := fs.Arg(0)
	svc, projectID, err := openGraphService()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	_ = version
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	out, err := svc.CodeExplore(ctx, projectID, symbol)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	if jsonOut {
		b, _ := json.MarshalIndent(out, "", "  ")
		fmt.Fprintln(os.Stdout, string(b))
		return
	}
	for path, spans := range out.Source {
		fmt.Printf("== %s\n", path)
		for _, sp := range spans {
			fmt.Printf("-- %s (%d-%d)\n%s\n", sp.Symbol, sp.StartLine, sp.EndLine, sp.Content)
		}
	}
	if len(out.Flow) > 0 {
		fmt.Println("call flow:")
		for _, e := range out.Flow {
			fmt.Printf("  %s -> %s [%s %s] :%d\n", e.From, e.To, e.Kind, e.Confidence, e.Line)
		}
	}
	if out.Impact != nil {
		fmt.Println("blast radius:")
		printImpact(out.Impact)
	}
}
