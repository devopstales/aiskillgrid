package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/devopstales/skillgrid/skillgrid-cli/internal/mnemonic/eval"
)

// evalCorpus is one named corpus for the eval runner: its name (the pooled axis
// — "self" / a second-language name) and its root (the git repo / file tree it
// is derived from).
type evalCorpus struct {
	Name string
	Root string
}

// readGitHistory reads a repo's commit log (hash, subject, changed files) in
// reverse-chronological order. It shells out to git (the cmd layer may use
// os/exec; the eval package stays pure-Go and just consumes these rows).
func readGitHistory(root string) ([]eval.Commit, error) {
	cmd := exec.Command("git", "-C", root,
		"log", "--pretty=format:%H%x1f%s", "--name-only")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git log in %s: %w", root, err)
	}
	// Format: <hash>\x1f<subject>\n<file>\n<file>\n ... per commit, newest first.
	// Split on newlines; a line with \x1f starts a new commit.
	lines := strings.Split(string(out), "\n")
	var commits []eval.Commit
	for _, line := range lines {
		if line == "" {
			continue
		}
		if idx := strings.Index(line, "\x1f"); idx >= 0 {
			commits = append(commits, eval.Commit{
				Hash:    line[:idx],
				Subject: strings.TrimSpace(line[idx+1:]),
			})
			continue
		}
		if len(commits) > 0 {
			commits[len(commits)-1].Files = append(commits[len(commits)-1].Files, strings.TrimSpace(line))
		}
	}
	return commits, nil
}

// loadCorpusFiles returns the corpus's file set at HEAD (path → content),
// excluding benchmark scaffolding (01.18) and common non-code files.
func loadCorpusFiles(root string) map[string]string {
	out := map[string]string{}
	_ = filepath.WalkDir(root, func(p string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			name := d.Name()
			if name == ".git" || name == "vendor" || name == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}
		rel, rerr := filepath.Rel(root, p)
		if rerr != nil {
			return nil
		}
		rel = filepath.ToSlash(rel)
		if eval.ExcludedFromCorpus(rel) {
			return nil
		}
		// Only text-y code files (the eval grades code retrieval).
		if !isCodeFile(rel) {
			return nil
		}
		if b, rerr := os.ReadFile(p); rerr == nil {
			out[rel] = string(b)
		}
		return nil
	})
	return out
}

// isCodeFile reports whether a path is a code file the eval grades.
func isCodeFile(path string) bool {
	lower := strings.ToLower(path)
	for _, ext := range []string{".go", ".py", ".js", ".ts", ".java", ".rb", ".rs", ".c", ".cpp", ".h", ".hpp", ".cs", ".kt", ".swift"} {
		if strings.HasSuffix(lower, ext) {
			return true
		}
	}
	return false
}

// parseEvalCorpora parses the repeated --corpus flags. Each value is `name` or
// `name=path`. `self` (no path) resolves to the cwd's git root; any other name
// REQUIRES a path that exists — an unknown name or missing path is rejected
// clearly (no invented run).
func parseEvalCorpora(values []string) ([]evalCorpus, error) {
	if len(values) == 0 {
		return nil, fmt.Errorf("eval: at least one --corpus is required (e.g. --corpus self or --corpus second=/path/to/repo)")
	}
	cwd, _ := os.Getwd()
	var out []evalCorpus
	seen := map[string]bool{}
	for _, v := range values {
		name, path, hasPath := strings.Cut(v, "=")
		name = strings.TrimSpace(name)
		if name == "" {
			return nil, fmt.Errorf("eval: empty corpus name in %q", v)
		}
		if seen[name] {
			return nil, fmt.Errorf("eval: duplicate corpus name %q", name)
		}
		seen[name] = true
		switch {
		case name == "self" && !hasPath:
			root, err := gitRootOf(cwd)
			if err != nil {
				return nil, fmt.Errorf("eval: cannot resolve the self corpus (cwd is not in a git repo): %w", err)
			}
			out = append(out, evalCorpus{Name: name, Root: root})
		case hasPath:
			path = strings.TrimSpace(path)
			if fi, err := os.Stat(path); err != nil || !fi.IsDir() {
				return nil, fmt.Errorf("eval: corpus %q path %q does not exist or is not a directory", name, path)
			}
			out = append(out, evalCorpus{Name: name, Root: path})
		default:
			return nil, fmt.Errorf("eval: unknown corpus %q — use `self` or `name=path` (a path to a corpus root)", name)
		}
	}
	return out, nil
}

// gitRootOf returns the git root of dir (or an error when not in a repo).
func gitRootOf(dir string) (string, error) {
	cmd := exec.Command("git", "-C", dir, "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

// runEval handles `skillgrid eval --corpus <name> [--corpus <name>] [--json]`.
// It derives a leak-free query set from each corpus's git history, builds ONE
// shared index per corpus, runs the ablation (baseline + the shipped ranker
// config), and prints the report with significance. Unknown corpora are
// rejected clearly (01.21).
func runEval(version string, args []string) {
	_ = version
	fs := flag.NewFlagSet("eval", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	var corpora []string
	var jsonOut bool
	fs.Var((*stringSliceFlag)(&corpora), "corpus", "a corpus: `self` or `name=path` (repeatable)")
	fs.BoolVar(&jsonOut, "json", false, "emit machine-readable JSON")
	fs.Usage = func() {
		fmt.Fprintln(fs.Output(), `usage: skillgrid eval --corpus <self|name=path> [--corpus ...] [--json]

Runs the retrieval-eval ablation over the given corpora. The query set is
derived from each corpus's git history (leak-free: merges/reverts/releases/
bumps/formatting/changelog/benchmark commits are dropped). One shared index per
corpus is built so deltas measure ranking only. The report carries the baseline
plus the shipped ranker config with a seeded paired bootstrap 95% CI and a
permutation p-value vs the baseline.

Examples:
  skillgrid eval --corpus self
  skillgrid eval --corpus self --corpus second=/path/to/other-repo`)
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}
	corps, err := parseEvalCorpora(corpora)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(2)
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	fileSets := make(map[string]map[string]string)
	querySets := make(map[string][]*eval.QuerySet)
	for _, c := range corps {
		commits, herr := readGitHistory(c.Root)
		if herr != nil {
			fmt.Fprintf(os.Stderr, "error: corpus %q: %v\n", c.Name, herr)
			os.Exit(1)
		}
		fileSets[c.Name] = loadCorpusFiles(c.Root)
		querySets[c.Name] = eval.DeriveQuerySet(commits)
	}

	res, err := eval.Run(ctx, eval.RunConfig{
		Corpora: fileSets,
		Queries: querySets,
		Variants: []eval.Variant{
			{Name: "baseline", Rank: identityRank, Baseline: true},
			{Name: "shipped", Rank: identityRank},
		},
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	if jsonOut {
		b, _ := json.MarshalIndent(res, "", "  ")
		fmt.Fprintln(os.Stdout, string(b))
		return
	}
	printEvalReport(res, querySets)
}

// identityRank is the shipped ranker config for the CLI: an identity ranking
// over the corpus (the 005 baseline behavior). The CLI's purpose is to RUN the
// harness and report significance; the shipped config and its decision are
// produced by the eval package's significance gate.
func identityRank(q *eval.QuerySet, idx *eval.CorpusIndex) []string {
	return idx.Paths()
}

// stringSliceFlag is a repeatable string flag for --corpus.
type stringSliceFlag []string

func (s *stringSliceFlag) String() string { return strings.Join(*s, ",") }
func (s *stringSliceFlag) Set(v string) error {
	*s = append(*s, v)
	return nil
}

// printEvalReport renders the eval ablation report (human table).
func printEvalReport(res *eval.Result, querySets map[string][]*eval.QuerySet) {
	fmt.Printf("skillgrid eval — corpora: %s (queries: %d)\n",
		strings.Join(res.Corpora, ", "), totalQueries(querySets))
	fmt.Println()
	if res.Baseline != nil {
		fmt.Printf("baseline : recall@5=%.3f recall@10=%.3f MRR=%.3f nDCG@10=%.3f dup%%=%.3f tokens=%d\n",
			res.Baseline.Recall5, res.Baseline.Recall10, res.Baseline.MRR, res.Baseline.NDGC10, res.Baseline.DupPercent, res.Baseline.Tokens)
	}
	for _, row := range res.Rows {
		d := res.Decisions[row.Name]
		fmt.Printf("%-10s: recall@5=%.3f (Δ%.3f) MRR=%.3f (Δ%.3f) CI[%.3f,%.3f] p=%.4f → %s\n",
			row.Name, row.Recall5, row.DeltaRecall5, row.MRR, row.DeltaMRR,
			row.CILower, row.CIUpper, row.PValue, d.Decision)
	}
	// One-index-per-corpus invariant (01.19).
	for name, n := range res.IndexInstances {
		if n != 1 {
			fmt.Printf("  note: corpus %q was indexed %d times (expected 1)\n", name, n)
		}
	}
}

func totalQueries(qs map[string][]*eval.QuerySet) int {
	n := 0
	for _, v := range qs {
		n += len(v)
	}
	return n
}
