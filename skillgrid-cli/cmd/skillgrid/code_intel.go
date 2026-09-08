package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

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

  orient:
    --json    Emit machine-readable JSON (default: human table)
  grep:
    --json    Emit machine-readable JSON (default: human table)`)
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

var _ service.Service
