package main

import (
	"context"
	"os/exec"
	"strconv"
	"strings"
	"testing"

	"github.com/devopstales/skillgrid/skillgrid-cli/internal/mnemonic/memory"
	"github.com/devopstales/skillgrid/skillgrid-cli/internal/mnemonic/memory/layer"
	"github.com/devopstales/skillgrid/skillgrid-cli/internal/mnemonic/store"
)

// seedMemCLIPeriod creates a project store with a session, a saved observation,
// and a distilled L2/L3 layer set (so `mem layers` has something to inspect).
// It returns the project id, the observation id, and the session id.
func seedMemCLIPeriod(t *testing.T, dataDir, project string) (int64, string) {
	t.Helper()
	st, err := store.Open(dataDir, project)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer st.Close()
	mem := memory.New(st, project)
	ctx := context.Background()
	if _, err := st.DB.Exec(`
		INSERT INTO sessions (id, project, directory, started_at, status, title, summary)
		VALUES ('sess-memcli', ?, '/tmp', '2026-01-01T00:00:00Z', 'ended', 'mem cli session', ?)`,
		project, "## Goal\nmem cli parity\n\n## Key Learnings:\n- The CLI must expose parity with the memory tools for layers and governance"); err != nil {
		t.Fatalf("insert session: %v", err)
	}
	id, err := mem.Save(ctx, memory.SaveInput{
		SessionID: "sess-memcli",
		Type:      "decision",
		Title:     "mem cli parity probe",
		Content:   "CLI parity body for mem governance/share/layers",
	})
	if err != nil {
		t.Fatalf("save: %v", err)
	}
	// Distill the session so observation_layers + personas exist (mem layers).
	if _, err := layer.Distill(ctx, mem, "sess-memcli", layer.DistillOptions{}); err != nil {
		t.Fatalf("distill: %v", err)
	}
	return id, "sess-memcli"
}

func runMemCLI(t *testing.T, dataDir string, args ...string) string {
	t.Helper()
	cmdArgs := append([]string{"run", ".", "mem"}, args...)
	cmd := exec.Command("go", cmdArgs...)
	cmd.Dir = mustWD(t)
	cmd.Env = append(cmd.Environ(), "SKILLGRID_MNEMONIC_DATA_DIR="+dataDir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("mem %v: %v\n%s", args, err, out)
	}
	return string(out)
}

// TestMemCLIParity is 03.8 [AFK] — the CLI exposes parity with the memory
// tools: `mem layers`, `mem governance`, and `mem share` each route through the
// same seams as their MCP counterparts.
// Scenario: cli-parity-for-layer-governance-share.
func TestMemCLIParity(t *testing.T) {
	dataDir := t.TempDir()
	project := "memcli-parity"
	obsID, sessID := seedMemCLIPeriod(t, dataDir, project)

	// mem layers <session_id>: returns the L0→L1→L2→L3 chain.
	out := runMemCLI(t, dataDir, "layers", sessID, "--project", project, "--dir", dataDir)
	if !strings.Contains(out, "l0_session") || !strings.Contains(out, sessID) {
		t.Fatalf("mem layers missing the L0 session chain: %s", out)
	}
	// The distilled L2/L3 layers must be present (provenance-linked).
	if !strings.Contains(out, "L2") && !strings.Contains(out, "L3") {
		t.Fatalf("mem layers missing distilled L2/L3 layers: %s", out)
	}

	// mem governance <id>: returns the governed-asset view.
	out = runMemCLI(t, dataDir, "governance", strconv.FormatInt(obsID, 10), "--project", project, "--dir", dataDir)
	if !strings.Contains(out, "visibility") || !strings.Contains(out, "owner") {
		t.Fatalf("mem governance missing asset fields: %s", out)
	}

	// mem share <id> --target-visibility team: widens visibility to team.
	out = runMemCLI(t, dataDir, "share", strconv.FormatInt(obsID, 10),
		"--target-visibility", "team", "--project", project, "--dir", dataDir)
	if !strings.Contains(out, "team") {
		t.Fatalf("mem share did not report the team visibility: %s", out)
	}

	// Verify the share landed (governance now shows team).
	out = runMemCLI(t, dataDir, "governance", strconv.FormatInt(obsID, 10), "--project", project, "--dir", dataDir)
	if !strings.Contains(out, `"team"`) {
		t.Fatalf("after mem share, governance must show visibility team: %s", out)
	}
}

// runMemCLIExpectError runs the mem CLI and returns (output, non-nil error);
// it does NOT t.Fatal on a non-zero exit (for negative-arg assertions).
func runMemCLIExpectError(t *testing.T, dataDir string, args ...string) (string, error) {
	t.Helper()
	cmdArgs := append([]string{"run", ".", "mem"}, args...)
	cmd := exec.Command("go", cmdArgs...)
	cmd.Dir = mustWD(t)
	cmd.Env = append(cmd.Environ(), "SKILLGRID_MNEMONIC_DATA_DIR="+dataDir)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// TestMemSearchModeFlag is 02.3 [AFK] — the CLI `mem search --mode` accepts
// trigram/prefix/phrase/all, routes trigram/prefix through the FTS match-mode
// read path (BudgetedRetrievalAsRootFTS), and rejects invalid modes with a
// clear error (exit 2, no store access).
func TestMemSearchModeFlag(t *testing.T) {
	dataDir := t.TempDir()
	project := "memcli-mode-flag"
	seedMemCLIPeriod(t, dataDir, project)

	// --mode trigram: the FTS fact leg runs in trigram mode (no error,
	// budgeted result shape).
	out := runMemCLI(t, dataDir, "search", "mem", "--mode", "trigram", "--project", project, "--dir", dataDir)
	if !strings.Contains(out, `"observations"`) {
		t.Fatalf("mem search --mode trigram should return the budgeted search shape, got: %s", out)
	}

	// --mode prefix: same read path, prefix FTS expansion.
	out = runMemCLI(t, dataDir, "search", "mem", "--mode", "prefix", "--project", project, "--dir", dataDir)
	if !strings.Contains(out, `"observations"`) {
		t.Fatalf("mem search --mode prefix should return the budgeted search shape, got: %s", out)
	}

	// --mode phrase: the default (phrase OR) alias.
	out = runMemCLI(t, dataDir, "search", "mem", "--mode", "phrase", "--project", project, "--dir", dataDir)
	if !strings.Contains(out, `"observations"`) {
		t.Fatalf("mem search --mode phrase should return the budgeted search shape, got: %s", out)
	}

	// --mode all: AND-joined phrases.
	out = runMemCLI(t, dataDir, "search", "mem", "--mode", "all", "--project", project, "--dir", dataDir)
	if !strings.Contains(out, `"observations"`) {
		t.Fatalf("mem search --mode all should return the budgeted search shape, got: %s", out)
	}

	// Invalid mode: rejected clearly with a non-zero exit, before any store
	// access (the flag parsing layer validates, so the error is ours).
	out, err := runMemCLIExpectError(t, dataDir, "search", "mem", "--mode", "bogus", "--project", project, "--dir", dataDir)
	if err == nil {
		t.Fatalf("mem search --mode bogus should fail, got: %s", out)
	}
	if !strings.Contains(out, "invalid match mode") {
		t.Fatalf("invalid --mode should be rejected clearly, got: %s", out)
	}
}

// TestMemCLIBudgetedContext is the finding-03.3 proof: the CLI `mem context`
// honors its --char budget flag (previously ignored). A session with a long
// summary is char-truncated with an explicit "N chars omitted" marker when a
// small --char cap is passed.
func TestMemCLIBudgetedContext(t *testing.T) {
	dataDir := t.TempDir()
	project := "memcli-budgeted-ctx"
	st, err := store.Open(dataDir, project)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	longSummary := "## Goal\nmem context budget probe\n\n## Key Learnings:\n- " + strings.Repeat("learn", 300)
	if _, err := st.DB.Exec(`
		INSERT INTO sessions (id, project, directory, started_at, status, title, summary)
		VALUES ('sess-ctx', ?, '/tmp', '2026-01-01T00:00:00Z', 'ended', 'ctx budget session', ?)`,
		project, longSummary); err != nil {
		t.Fatalf("insert session: %v", err)
	}
	st.Close()

	// A small --char cap must truncate the long summary (explicit marker).
	out := runMemCLI(t, dataDir, "context", "--project", project, "--char", "40", "--dir", dataDir)
	if !strings.Contains(out, "chars omitted") {
		t.Fatalf("mem context with --char 40 must char-truncate the summary (explicit 'chars omitted'), got: %s", out)
	}
	if !strings.Contains(out, `"truncated": true`) {
		t.Fatalf("mem context with a truncating --char must report truncated:true, got: %s", out)
	}
}

// TestMemCLIBadArgs proves the CLI rejects bad mem args clearly (parity with
// the MCP bad-args behavior).
func TestMemCLIBadArgs(t *testing.T) {
	dataDir := t.TempDir()
	project := "memcli-badargs"
	st, err := store.Open(dataDir, project)
	if err != nil {
		t.Fatal(err)
	}
	st.Close()

	// mem layers with no target → clear error (exit 2).
	cmd := exec.Command("go", "run", ".", "mem", "layers", "--project", project, "--dir", dataDir)
	cmd.Dir = mustWD(t)
	cmd.Env = append(cmd.Environ(), "SKILLGRID_MNEMONIC_DATA_DIR="+dataDir)
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("mem layers with no target should fail, got: %s", out)
	}
	if !strings.Contains(string(out), "requires") && !strings.Contains(string(out), "session_id") {
		t.Fatalf("bad mem layers args should be rejected clearly, got: %s", out)
	}

	// mem share with an invalid visibility → clear error.
	id := "1"
	cmd = exec.Command("go", "run", ".", "mem", "share", id, "--target-visibility", "bogus", "--project", project, "--dir", dataDir)
	cmd.Dir = mustWD(t)
	cmd.Env = append(cmd.Environ(), "SKILLGRID_MNEMONIC_DATA_DIR="+dataDir)
	out, err = cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("mem share with bad visibility should fail, got: %s", out)
	}
}
