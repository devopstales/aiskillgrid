package main

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/devopstales/skillgrid/skillgrid-cli/internal/mnemonic/eval"
)

// gitRun runs a git command in dir, fatal-ing the test on error.
func gitRun(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=test", "GIT_AUTHOR_EMAIL=test@example.com",
		"GIT_COMMITTER_NAME=test", "GIT_COMMITTER_EMAIL=test@example.com",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v (out %s)", args, err, out)
	}
}

// evalGitFixture builds a real git repo in a temp dir with a genuine code
// change + a noise change (a merge + a version bump), returning the repo root.
func evalGitFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	gitRun(t, root, "init", "-q")
	gitRun(t, root, "config", "user.email", "test@example.com")
	gitRun(t, root, "config", "user.name", "test")

	// Commit 1 (genuine): a real code change.
	if err := os.MkdirAll(filepath.Join(root, "internal", "http"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "internal", "http", "client.go"),
		[]byte("package http\n\nfunc NewClient() *Client { return &Client{} }\n"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	gitRun(t, root, "add", "-A")
	gitRun(t, root, "commit", "-q", "-m", "feat: add http client")

	// Commit 2 (genuine): a real code change in a second file.
	if err := os.WriteFile(filepath.Join(root, "internal", "http", "retry.go"),
		[]byte("package http\n\nfunc (c *Client) Retry(n int) bool { return n > 0 }\n"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	gitRun(t, root, "add", "-A")
	gitRun(t, root, "commit", "-q", "-m", "fix: retry backoff")

	// Commit 3 (noise): a version bump (go.mod) — must be dropped.
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module x\n\ngo 1.21\n"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	gitRun(t, root, "add", "-A")
	gitRun(t, root, "commit", "-q", "-m", "chore: bump version to 1.21")

	return root
}

// TestSkillgridEvalSelfCorpus covers @step-01 (Scenario: `skillgrid eval
// --corpus <self>`): the eval runner reads the git history of the self corpus,
// derives the leak-free query set, builds one index, runs the ablation, and
// emits a report with the baseline + the shipped config and its significance.
func TestSkillgridEvalSelfCorpus(t *testing.T) {
	root := evalGitFixture(t)
	// Derive commits from the real git history and confirm the noise drop.
	commits, err := readGitHistory(root)
	if err != nil {
		t.Fatalf("readGitHistory: %v", err)
	}
	if len(commits) < 2 {
		t.Fatalf("expected >=2 commits, got %d", len(commits))
	}
	queries := eval.DeriveQuerySet(commits)
	// The version-bump commit must be dropped; the two real commits kept.
	if len(queries) != 2 {
		t.Errorf("expected 2 leak-free queries (version bump dropped), got %d: %+v", len(queries), queries)
	}

	// The self corpus's file set (HEAD).
	files := loadCorpusFiles(root)
	if len(files) == 0 {
		t.Fatalf("expected corpus files, got none")
	}
	// Run the harness over the self corpus (one index, baseline + shipped).
	res, err := eval.Run(context.Background(), eval.RunConfig{
		Corpora: map[string]map[string]string{"self": files},
		Queries: map[string][]*eval.QuerySet{"self": queries},
		Variants: []eval.Variant{
			{Name: "baseline", Rank: func(q *eval.QuerySet, idx *eval.CorpusIndex) []string { return idx.Paths() }, Baseline: true},
		},
	})
	if err != nil {
		t.Fatalf("eval.Run: %v", err)
	}
	if res.IndexInstances["self"] != 1 {
		t.Errorf("self corpus should be indexed once, got %d", res.IndexInstances["self"])
	}
}

// TestSkillgridEvalUnknownCorpus covers @step-01 (bad community/eval args
// rejected): an unknown/invalid --corpus is rejected with a clear error, not an
// invented run.
func TestSkillgridEvalUnknownCorpus(t *testing.T) {
	_, err := parseEvalCorpora([]string{"self", "not-a-real-name"})
	if err == nil {
		t.Fatal("expected an error for an unknown corpus name without a path")
	}
	if !strings.Contains(err.Error(), "not-a-real-name") {
		t.Errorf("unknown-corpus error should name the corpus, got: %v", err)
	}

	// A corpus with a missing path is also rejected clearly.
	_, err = parseEvalCorpora([]string{"self", "ghost=/nonexistent/path/xyz"})
	if err == nil {
		t.Fatal("expected an error for a corpus path that does not exist")
	}
	if !strings.Contains(err.Error(), "does not exist") {
		t.Errorf("missing-path error should say the path does not exist, got: %v", err)
	}
}

// TestSkillgridEvalRequiresCorpus covers @step-01: the eval runner requires at
// least one --corpus (no invented default corpus).
func TestSkillgridEvalRequiresCorpus(t *testing.T) {
	if _, err := parseEvalCorpora(nil); err == nil {
		t.Error("expected an error when no --corpus is provided")
	}
}
