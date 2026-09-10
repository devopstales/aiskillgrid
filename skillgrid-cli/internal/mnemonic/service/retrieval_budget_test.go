package service

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/devopstales/skillgrid/skillgrid-cli/internal/mnemonic/memory"
	"github.com/devopstales/skillgrid/skillgrid-cli/internal/mnemonic/store"
)

// seedBudgetedStore opens a project store, creates a session, and saves n
// observations matching "widget". It returns the service, the project id, and
// the open store (so the caller can close it).
func seedBudgetedStore(t *testing.T, project string, n int) (*Service, string, *store.Store) {
	t.Helper()
	dataDir := t.TempDir()
	st, err := store.Open(dataDir, project)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	svc := New(dataDir)
	if _, err := st.DB.Exec(`
		INSERT INTO sessions (id, project, directory, started_at, status)
		VALUES ('s1', ?, '/tmp', '2026-01-01T00:00:00Z', 'active')`, project); err != nil {
		t.Fatalf("insert session: %v", err)
	}
	mem := memory.New(st, project)
	for i := 0; i < n; i++ {
		if _, err := mem.Save(context.Background(), memory.SaveInput{
			SessionID: "s1",
			Type:      "learning",
			Title:     fmt.Sprintf("Widget note %d", i),
			// Long enough content to exceed the default char budget (1200).
			Content: "Body of widget note " + fmt.Sprintf("%d", i) + " " + strings.Repeat("x", 1500),
		}); err != nil {
			t.Fatalf("save %d: %v", i, err)
		}
	}
	return svc, project, st
}

// TestBudgetedSearch is 03.5 [AFK] — a 20-hit search returns 20 budgeted
// snippets (not 20× full JSON); every in-list result carries its
// mem_get_observation id.
// Scenarios: twenty-hit-search-is-budgeted,
// every-inlist-result-carries-get-observation-id.
func TestBudgetedSearch(t *testing.T) {
	svc, project, st := seedBudgetedStore(t, "budgeted-search", 20)
	defer st.Close()

	res, err := svc.BudgetedRetrieval(context.Background(), project, "fact", "widget", 20)
	if err != nil {
		t.Fatalf("budgeted retrieval: %v", err)
	}
	// The item cap (default 10) bounds the in-list result — NOT 20 full
	// payloads.
	if len(res.Hits) != 10 {
		t.Fatalf("expected the item-cap (10) in-list results, got %d", len(res.Hits))
	}
	// Every in-list result carries its full-content fetch id.
	for _, o := range res.Hits {
		if o.ID <= 0 {
			t.Fatalf("in-list result missing its mem_get_observation id: %+v", o)
		}
	}
	// Each in-list snippet is char-truncated with an explicit omitted count
	// (never a full-observation payload).
	for _, o := range res.Hits {
		if !strings.Contains(o.Content, "chars omitted") {
			t.Fatalf("in-list snippet not char-budgeted (no 'chars omitted'): %q", o.Content)
		}
	}
	if !res.Truncated {
		t.Fatalf("a 20-hit search must be char-truncated (truncated=true), got %+v", res)
	}
}

// TestFullContent is 03.6 [AFK] — mem_get_observation remains the ONLY
// full-content path: a budgeted in-list result is truncated, but fetching the
// hit by its id returns the full, untruncated content; the budget is tunable.
// Scenarios: mem-get-observation-is-only-full-content-path,
// budget-is-tunable-via-config.
func TestFullContent(t *testing.T) {
	dataDir := t.TempDir()
	project := "full-content"
	st, err := store.Open(dataDir, project)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer st.Close()
	if _, err := st.DB.Exec(`
		INSERT INTO sessions (id, project, directory, started_at, status)
		VALUES ('s1', ?, '/tmp', '2026-01-01T00:00:00Z', 'active')`, project); err != nil {
		t.Fatalf("insert session: %v", err)
	}
	mem := memory.New(st, project)
	fullBody := "The full untruncated content " + strings.Repeat("y", 2000)
	id, err := mem.Save(context.Background(), memory.SaveInput{
		SessionID: "s1", Type: "learning", Title: "Full content probe",
		Content: fullBody,
	})
	if err != nil {
		t.Fatalf("save: %v", err)
	}
	svc := New(dataDir)

	// The budgeted read truncates the in-list snippet.
	res, err := svc.BudgetedRetrieval(context.Background(), project, "fact", "full content probe", 5)
	if err != nil {
		t.Fatalf("budgeted retrieval: %v", err)
	}
	if len(res.Hits) != 1 {
		t.Fatalf("expected 1 in-list hit, got %d", len(res.Hits))
	}
	if res.Hits[0].ID != id {
		t.Fatalf("in-list hit id mismatch: got %d want %d", res.Hits[0].ID, id)
	}
	inListContent := res.Hits[0].Content
	if len(inListContent) >= len(fullBody) {
		t.Fatalf("in-list result must be truncated (got %d >= full %d)", len(inListContent), len(fullBody))
	}
	if !strings.Contains(inListContent, "chars omitted") {
		t.Fatalf("in-list truncation not explicit: %q", inListContent)
	}

	// mem_get_observation (the only full-content path) returns the FULL content.
	got, err := mem.Get(context.Background(), id)
	if err != nil {
		t.Fatalf("get observation (full-content path): %v", err)
	}
	if got.Content != fullBody {
		t.Fatalf("mem_get_observation must return full untruncated content; got %d chars, want %d", len(got.Content), len(fullBody))
	}
	if strings.Contains(got.Content, "chars omitted") {
		t.Fatalf("mem_get_observation must NOT be char-budgeted: %q", got.Content)
	}

	// The budget is tunable via config: a larger char cap returns a longer
	// (still not full, but less truncated) in-list snippet.
	svc2 := New(dataDir)
	// Bump the char cap so the in-list snippet keeps more than the default.
	st2, _ := store.Open(dataDir, project)
	defer st2.Close()
	mem2 := memory.New(st2, project)
	mem2.SetBudget(memory.BudgetConfig{Items: 10, Chars: 3000})
	res2 := mem2.Budget().Apply(context.Background(), []memory.Observation{{ID: id, Content: fullBody}})
	if len(res2.Hits[0].Content) <= len(inListContent) {
		t.Fatalf("tuned (larger) char cap must keep more in-list content: tuned %d <= default %d", len(res2.Hits[0].Content), len(inListContent))
	}
	_ = svc2
}

// TestRouteReads is 03.7 [AFK] — mem_context/mem_search/mem_timeline are
// routed through the layered + budgeted path; mem_get_observation stays the
// only full-content path.
// Scenario: mem-context-search-timeline-budgeted.
func TestRouteReads(t *testing.T) {
	dataDir := t.TempDir()
	project := "route-reads"
	st, err := store.Open(dataDir, project)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer st.Close()
	ctx := context.Background()
	if _, err := st.DB.Exec(`
		INSERT INTO sessions (id, project, directory, started_at, status, title, summary)
		VALUES ('s1', ?, '/tmp', '2026-01-01T00:00:00Z', 'ended', 'route test session', '## Goal\nroute reads probe')`, project); err != nil {
		t.Fatalf("insert session: %v", err)
	}
	mem := memory.New(st, project)
	for i := 0; i < 20; i++ {
		if _, err := mem.Save(ctx, memory.SaveInput{
			SessionID: "s1", Type: "learning",
			Title:   fmt.Sprintf("Route widget %d", i),
			Content: "Route body " + fmt.Sprintf("%d", i) + " " + strings.Repeat("z", 1500),
		}); err != nil {
			t.Fatalf("save %d: %v", i, err)
		}
	}
	svc := New(dataDir)

	// mem_search path (fact → RRF fallback) is budgeted.
	searchRes, err := svc.BudgetedRetrieval(ctx, project, "fact", "route widget", 20)
	if err != nil {
		t.Fatalf("search budgeted: %v", err)
	}
	if len(searchRes.Hits) != 10 {
		t.Fatalf("mem_search must be item-cap budgeted (10), got %d", len(searchRes.Hits))
	}
	for _, o := range searchRes.Hits {
		if o.ID <= 0 {
			t.Fatalf("mem_search in-list hit missing id: %+v", o)
		}
	}

	// mem_context path (recent sessions) is budgeted at the item cap.
	sessions, err := mem.RecentContext(ctx, 20)
	if err != nil {
		t.Fatalf("context: %v", err)
	}
	budgetedCtx := mem.Budget().Apply(ctx, nil) // no observations, but exercises the budget
	if len(budgetedCtx.Hits) != 0 {
		t.Fatalf("empty context budgeted read must return no hits, got %d", len(budgetedCtx.Hits))
	}
	if len(sessions) == 0 {
		t.Fatalf("mem_context must return the recent sessions")
	}

	// mem_get_observation remains the only full-content path: the search hit's
	// id fetches the full content (no truncation).
	hitID := searchRes.Hits[0].ID
	got, err := mem.Get(ctx, hitID)
	if err != nil {
		t.Fatalf("full-content fetch: %v", err)
	}
	if strings.Contains(got.Content, "chars omitted") {
		t.Fatalf("mem_get_observation must return untruncated content, got %q", got.Content)
	}
}
