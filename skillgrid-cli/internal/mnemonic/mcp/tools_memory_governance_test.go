package mcp

import (
	"context"
	"encoding/json"
	"sort"
	"testing"

	"github.com/devopstales/skillgrid/skillgrid-cli/internal/mnemonic/service"
)

// memGovernanceFixture pins the project to a stable bucket so handlers that
// open the CWD project resolve to one store.
func memGovernanceFixture(t *testing.T) {
	t.Helper()
	dataDir := t.TempDir()
	t.Setenv("MNEMONIC_PROJECT", "govmcp-probe")
	svc := service.New(dataDir)
	SetService(svc)
	t.Cleanup(func() { SetService(nil) })
}

// TestMemGovernanceTools covers @step-01 RED "Mnemonic tool surface":
// mem_share + mem_governance are registered with distinct mem_* names, the
// existing 005 mem_* tools keep their names + required params unchanged, and
// bad governance args are rejected with a clear error (visibility unchanged).
func TestMemGovernanceTools(t *testing.T) {
	memGovernanceFixture(t)

	// New tools registered.
	tools := NewServer().ListTools()
	for _, name := range []string{"mem_share", "mem_governance"} {
		if _, ok := tools[name]; !ok {
			t.Errorf("expected tool %q to be registered", name)
		}
	}

	// Tool surface grows additively: 75 baseline + 2 governance = 77.
	if len(tools) != 77 {
		t.Errorf("expected 77 tools (75 baseline + 2 governance), got %d", len(tools))
	}

	// Existing 005 mem_* tools keep their names + required params unchanged.
	for name, wantRequired := range expectedMemToolSurface {
		st, ok := tools[name]
		if !ok {
			t.Errorf("005 mem tool %q is no longer registered", name)
			continue
		}
		got := append([]string(nil), st.Tool.InputSchema.Required...)
		want := append([]string(nil), wantRequired...)
		sort.Strings(got)
		sort.Strings(want)
		if len(got) != len(want) {
			t.Errorf("%q: required params changed: got %v, want %v", name, got, want)
			continue
		}
		for i := range want {
			if got[i] != want[i] {
				t.Errorf("%q: required param %d changed: got %q, want %q", name, i, got[i], want[i])
			}
		}
	}

	// Bad governance args rejected clearly:
	//  - mem_share with no id → validation error.
	res, err := handleMemShare(context.Background(), newCallTool("mem_share", map[string]any{
		"target": "team",
	}))
	if err != nil {
		t.Fatalf("handleMemShare dispatch: %v", err)
	}
	if !res.IsError {
		t.Errorf("mem_share without id should be a validation error, got: %s", callResultText(t, res))
	}
	//  - mem_share with a bad target → validation error.
	res, err = handleMemShare(context.Background(), newCallTool("mem_share", map[string]any{
		"id":     1,
		"target": "galaxy",
	}))
	if err != nil {
		t.Fatalf("handleMemShare dispatch (bad target): %v", err)
	}
	if !res.IsError {
		t.Errorf("mem_share with bad target should be a validation error, got: %s", callResultText(t, res))
	}
	//  - mem_governance with no id → validation error.
	res, err = handleMemGovernance(context.Background(), newCallTool("mem_governance", map[string]any{}))
	if err != nil {
		t.Fatalf("handleMemGovernance dispatch: %v", err)
	}
	if !res.IsError {
		t.Errorf("mem_governance without id should be a validation error, got: %s", callResultText(t, res))
	}
}

// TestMemGovernanceRoundTrip covers @step-01 (Scenario:
// mem-governance-surfaces-asset-fields): a mem_save → mem_update → mem_share →
// mem_governance round-trip surfaces owner, version history, status, usage,
// and visibility, with prior content recoverable.
func TestMemGovernanceRoundTrip(t *testing.T) {
	memGovernanceFixture(t)

	// Create a session first (mem_save requires a session_id that exists).
	startRes, err := handleMemSessionStart(context.Background(), newCallTool("mem_session_start", map[string]any{}))
	if err != nil {
		t.Fatalf("handleMemSessionStart dispatch: %v", err)
	}
	if startRes.IsError {
		t.Fatalf("mem_session_start errored: %s", callResultText(t, startRes))
	}
	var startOut struct {
		SessionID string `json:"session_id"`
	}
	if err := json.Unmarshal([]byte(callResultText(t, startRes)), &startOut); err != nil {
		t.Fatalf("unmarshal session start: %v (text %s)", err, callResultText(t, startRes))
	}

	// Save an observation.
	saveRes, err := handleMemSave(context.Background(), newCallTool("mem_save", map[string]any{
	"title":      "gov roundtrip",
	"type":       "decision",
	"content":    "original body",
	"session_id": startOut.SessionID,
	}))
	if err != nil {
		t.Fatalf("handleMemSave: %v", err)
	}
	if saveRes.IsError {
		t.Fatalf("mem_save errored: %s", callResultText(t, saveRes))
	}
	var saveOut struct {
		ID int64 `json:"id"`
	}
	if err := json.Unmarshal([]byte(callResultText(t, saveRes)), &saveOut); err != nil {
		t.Fatalf("unmarshal save: %v (text %s)", err, callResultText(t, saveRes))
	}

	// Update (appends a version).
	if _, err := handleMemUpdate(context.Background(), newCallTool("mem_update", map[string]any{
		"id": float64(saveOut.ID), "content": "updated body",
	})); err != nil {
		t.Fatalf("handleMemUpdate dispatch: %v", err)
	}
	// Share to team.
	shareRes, err := handleMemShare(context.Background(), newCallTool("mem_share", map[string]any{
		"id": float64(saveOut.ID), "target": "team",
	}))
	if err != nil {
		t.Fatalf("handleMemShare dispatch: %v", err)
	}
	if shareRes.IsError {
		t.Fatalf("mem_share errored: %s", callResultText(t, shareRes))
	}

	// Query governance.
	govRes, err := handleMemGovernance(context.Background(), newCallTool("mem_governance", map[string]any{
		"id": float64(saveOut.ID),
	}))
	if err != nil {
		t.Fatalf("handleMemGovernance dispatch: %v", err)
	}
	if govRes.IsError {
		t.Fatalf("mem_governance errored: %s", callResultText(t, govRes))
	}
	var govOut struct {
		ID             int64  `json:"id"`
		Owner          string `json:"owner"`
		Visibility     string `json:"visibility"`
		Status         string `json:"status"`
		RevisionCount  int    `json:"revision_count"`
		RetrievalUsage int    `json:"retrieval_usage"`
		Versions       []struct {
			Revision int    `json:"revision"`
			Content  string `json:"content"`
		} `json:"versions"`
	}
	if err := json.Unmarshal([]byte(callResultText(t, govRes)), &govOut); err != nil {
		t.Fatalf("unmarshal governance: %v (text %s)", err, callResultText(t, govRes))
	}
	if govOut.Visibility != "team" {
		t.Errorf("governance visibility = %q, want team", govOut.Visibility)
	}
	if govOut.Status != "active" {
		t.Errorf("governance status = %q, want active", govOut.Status)
	}
	if govOut.RevisionCount != 1 {
		t.Errorf("governance revision_count = %d, want 1", govOut.RevisionCount)
	}
	if len(govOut.Versions) != 1 || govOut.Versions[0].Content != "original body" {
		t.Errorf("governance must recover prior content, got %+v", govOut.Versions)
	}
	if govOut.Owner == "" {
		t.Errorf("governance owner must be set")
	}
}

// expectedMemToolSurface pins the 005 mem_* tool contract: each tool's name
// and its required parameter list. Later changes (governance, layers, budgets)
// must never rename a 005 mem_* tool or change its required params — this map
// is the baseline lock that fails when they do.
var expectedMemToolSurface = map[string][]string{
	"mem_save":            {"title", "type", "content", "session_id"},
	"mem_search":          {"query"},
	"mem_context":         {},
	"mem_get_observation": {"id"},
	"mem_timeline":        {"id"},
	"mem_update":          {"id"},
	"mem_delete":          {"id"},
	"mem_stats":           {},
	"mem_save_prompt":     {"content", "session_id"},
	"mem_current_project": {},
	"mem_doctor":          {},
	"mem_review":          {},
	"mem_judge":           {"src_id", "dst_id", "verdict"},
	"mem_compare":         {"src_id", "dst_id"},
	"mem_merge_projects":  {"source", "canonical"},
	"mem_session_start":   {},
	"mem_session_end":     {"session_id"},
	"mem_session_summary": {"session_id", "summary"},
	"mem_session_set_title": {"session_id", "title"},
	"mem_suggest_topic_key": {"type"},
	"mem_capture_passive":   {"content"},
	"mem_pin":               {"id"},
	"mem_unpin":             {"id"},
	"mem_unify":             {"canonical", "sources"},
}
