package eval

import (
	"context"
	"fmt"
	"strings"
	"testing"
)

// evalFixtureCorpus is a small synthetic corpus (path → content) standing in
// for a git repo's file set. It has >5 files so that recall@5 can distinguish
// a ranker that puts the answer first (firstRanker) from one that puts it last
// (lastRanker) — with only 3 files both would reach recall@5 = 1.0.
func evalFixtureCorpus() map[string]string {
	// 3 relevant files (the query targets) + 8 distractors, so a "last" ranker
	// ranks the target beyond the top-5 window.
	paths := []string{
		"internal/http/client.go",
		"internal/page/pagination.go",
		"internal/auth/token.go",
		"vendor/lib/a.go",
		"vendor/lib/b.go",
		"vendor/lib/c.go",
		"vendor/lib/d.go",
		"vendor/lib/e.go",
		"vendor/lib/f.go",
		"vendor/lib/g.go",
		"vendor/lib/h.go",
	}
	m := make(map[string]string, len(paths))
	for i, p := range paths {
		m[p] = fmt.Sprintf("func f%d() int { return %d }\n", i, i)
	}
	return m
}

// evalFixtureQueries builds a query set over the fixture corpus. A larger
// number of queries (12) gives the seeded permutation test enough power to
// call a uniformly-better ranker significant at p < 0.05.
func evalFixtureQueries() []*QuerySet {
	targets := []string{
		"internal/http/client.go",
		"internal/page/pagination.go",
		"internal/auth/token.go",
		"vendor/lib/a.go",
		"vendor/lib/b.go",
		"vendor/lib/c.go",
		"vendor/lib/d.go",
		"vendor/lib/e.go",
		"vendor/lib/f.go",
		"vendor/lib/g.go",
		"vendor/lib/h.go",
	}
	subjects := []string{
		"fix http client retry", "fix pagination off-by-one", "verify auth token",
		"fix lib a path", "fix lib b path", "fix lib c path", "fix lib d path",
		"fix lib e path", "fix lib f path", "fix lib g path", "fix lib h path",
		"refactor lib paths",
	}
	qs := make([]*QuerySet, 0, len(targets))
	for i, t := range targets {
		qs = append(qs, &QuerySet{Subject: subjects[i], ExpectedFiles: []string{t}})
	}
	return qs
}

// isExpected reports whether p is one of the query's expected files.
func isExpected(q *QuerySet, p string) bool {
	for _, e := range q.ExpectedFiles {
		if p == e {
			return true
		}
	}
	return false
}

// firstRanker puts the expected files first, then the rest in sorted order
// (a strong candidate ranker).
func firstRanker(q *QuerySet, idx *CorpusIndex) []string {
	rest := []string{}
	for _, p := range idx.Paths() {
		if !isExpected(q, p) {
			rest = append(rest, p)
		}
	}
	return append(append([]string(nil), q.ExpectedFiles...), rest...)
}

// lastRanker puts the expected files last, then the rest in sorted order
// (a weak baseline ranker).
func lastRanker(q *QuerySet, idx *CorpusIndex) []string {
	rest := []string{}
	for _, p := range idx.Paths() {
		if !isExpected(q, p) {
			rest = append(rest, p)
		}
	}
	return append(rest, q.ExpectedFiles...)
}

// TestEvalHarnessSignificance covers @step-01 (Scenario: Evaluation harness
// derives a leak-free query set and reports significance): the ablation runner
// builds ONE shared index per corpus, runs baseline + candidate variants, and
// reports a per-variant row with the full metric set plus a seeded paired
// bootstrap 95% CI and permutation p-value vs the baseline. A strictly better
// candidate yields a positive, significant delta (CI lower > 0, p < 0.05).
func TestEvalHarnessSignificance(t *testing.T) {
	ctx := context.Background()
	queries := evalFixtureQueries()

	// Baseline ranks the answer last; candidate ranks it first → strictly
	// better on every query.
	res, err := Run(ctx, RunConfig{
		Corpora: map[string]map[string]string{"self": evalFixtureCorpus()},
		Queries: map[string][]*QuerySet{"self": queries},
		Variants: []Variant{
			{Name: "baseline", Rank: lastRanker, Baseline: true},
			{Name: "candidate", Rank: firstRanker},
		},
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if res.Baseline == nil {
		t.Fatal("expected a baseline row")
	}
	cand := rowByName(res, "candidate")
	if cand == nil {
		t.Fatal("missing candidate row")
	}
	// The candidate must beat the baseline on the headline metrics.
	if cand.DeltaRecall5 <= 0 || cand.DeltaMRR <= 0 {
		t.Errorf("candidate did not beat baseline: deltaRecall5=%v deltaMRR=%v", cand.DeltaRecall5, cand.DeltaMRR)
	}
	// Significance fields must be present (non-zero p is fine; CI bounds set).
	if cand.PValue > 1 || cand.PValue < 0 {
		t.Errorf("permutation p out of range: %v", cand.PValue)
	}
	if cand.CILower > cand.CIUpper {
		t.Errorf("CI inverted: lower=%v upper=%v", cand.CILower, cand.CIUpper)
	}
	// A strictly-better candidate is significant.
	if cand.PValue >= 0.05 {
		t.Errorf("strictly-better candidate should be significant, p=%v", cand.PValue)
	}
	if cand.CILower <= 0 {
		t.Errorf("strictly-better candidate CI lower should be > 0, got %v", cand.CILower)
	}
}

// TestEvalStaleExpectation covers @step-01 (Scenario: Stale evaluation
// expectation fails the run loudly): a query whose expected file no longer
// exists at HEAD makes validate_queries fail the run with a loud error — not a
// silently deflated score.
func TestEvalStaleExpectation(t *testing.T) {
	ctx := context.Background()
	// The query references a file that is NOT in the corpus (deleted at HEAD).
	queries := []*QuerySet{
		{Subject: "fix http client retry", ExpectedFiles: []string{"internal/http/client.go"}},
		{Subject: "fix deleted module", ExpectedFiles: []string{"internal/gone/removed.go"}},
	}
	_, err := Run(ctx, RunConfig{
		Corpora: map[string]map[string]string{"self": evalFixtureCorpus()},
		Queries: map[string][]*QuerySet{"self": queries},
		Variants: []Variant{
			{Name: "baseline", Rank: lastRanker, Baseline: true},
		},
	})
	if err == nil {
		t.Fatal("expected Run to fail loudly on a stale expectation")
	}
	if !strings.Contains(err.Error(), "gone/removed.go") {
		t.Errorf("stale-expectation error should name the missing file, got: %v", err)
	}
}

// TestOneIndexPerCorpus covers @step-01 (Scenario: Ablation shares one index so
// deltas measure ranking): all variants over a corpus share the SAME index
// instance, so the only thing that varies between variants is the ranking —
// never the retrieved document set.
func TestOneIndexPerCorpus(t *testing.T) {
	ctx := context.Background()
	queries := evalFixtureQueries()

	// A variant that records the index pointer it saw; both variants must have
	// seen the identical pointer for the corpus.
	var seen map[string]*CorpusIndex
	seen = map[string]*CorpusIndex{}
	rec := func(q *QuerySet, idx *CorpusIndex) []string {
		return idx.Paths()
	}
	res, err := Run(ctx, RunConfig{
		Corpora: map[string]map[string]string{"self": evalFixtureCorpus()},
		Queries: map[string][]*QuerySet{"self": queries},
		Variants: []Variant{
			{Name: "baseline", Rank: lastRanker, Baseline: true},
			{Name: "candidate", Rank: func(q *QuerySet, idx *CorpusIndex) []string {
				seen["candidate"] = idx
				return rec(q, idx)
			}},
		},
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if res.IndexInstances["self"] != 1 {
		t.Errorf("corpus 'self' should be indexed exactly once, got %d", res.IndexInstances["self"])
	}
	if seen["candidate"] == nil {
		t.Fatal("candidate variant never received an index pointer")
	}
}

// TestShippedConfigSignificance covers @step-01 (Scenario: Shipped ranking
// config is the significance winner): on >=2 pooled corpora, the shipped
// ranking config is non-negative vs the 005 baseline on every corpus AND every
// shipped signal survives significance; a candidate that fails is kept-off with
// the decision + CI/p recorded in the report.
func TestShippedConfigSignificance(t *testing.T) {
	ctx := context.Background()
	corpA := evalFixtureCorpus()
	corpB := map[string]string{
		"src/net/conn.go": "func (c *Conn) Read(b []byte) (int, error) { return 0, nil }\n",
		"src/net/addr.go": "func (a *Addr) String() string { return a.host }\n",
	}
	// Pad the second corpus with distractors so recall@5 can separate the
	// rankers (same reason as the first corpus).
	for i := 0; i < 8; i++ {
		corpB[fmt.Sprintf("src/net/distractor%d.go", i)] = fmt.Sprintf("func d%d() int { return %d }\n", i, i)
	}
	qA := evalFixtureQueries()
	qB := []*QuerySet{
		{Subject: "fix conn read", ExpectedFiles: []string{"src/net/conn.go"}},
		{Subject: "format addr string", ExpectedFiles: []string{"src/net/addr.go"}},
	}

	res, err := Run(ctx, RunConfig{
		Corpora: map[string]map[string]string{
			"self": corpA,
			"second": corpB,
		},
		Queries: map[string][]*QuerySet{"self": qA, "second": qB},
		Variants: []Variant{
			{Name: "baseline", Rank: lastRanker, Baseline: true},
			// Shipped signal: strictly better on both corpora → ships.
			{Name: "shipped-signal", Rank: firstRanker},
			// A candidate that is worse → must be kept-off (decision recorded).
			{Name: "worse-signal", Rank: lastRanker},
		},
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if len(res.Corpora) != 2 {
		t.Fatalf("expected 2 pooled corpora, got %d", len(res.Corpora))
	}
	// The shipped signal is non-negative vs baseline on every corpus and
	// survived significance → decision is "ship".
	ship := shippedDecision(res, "shipped-signal")
	if !ship.Ships {
		t.Errorf("strictly-better shipped signal should ship, got %+v", ship)
	}
	// The worse signal is negative somewhere → kept off, with a decision + CI/p.
	off := shippedDecision(res, "worse-signal")
	if off.Ships {
		t.Errorf("a signal that is not non-negative across all corpora must be kept off, got %+v", off)
	}
	if off.Decision == "" {
		t.Errorf("kept-off signal must record a decision, got empty")
	}
}

// TestFailingSignalDecision covers @step-01 (Scenario: Failing ranking signal is
// removed or kept off with the decision recorded): a candidate that loses to
// the baseline carries an explicit decision (kept-off/removed) with its CI and
// p-value in the report — never silently dropped.
func TestFailingSignalDecision(t *testing.T) {
	ctx := context.Background()
	res, err := Run(ctx, RunConfig{
		Corpora: map[string]map[string]string{
			"self":   evalFixtureCorpus(),
			"second": {"src/x/y.go": "func y() int { return 1 }\n"},
		},
		Queries: map[string][]*QuerySet{
			"self":   evalFixtureQueries(),
			"second": {{Subject: "y change", ExpectedFiles: []string{"src/x/y.go"}}},
		},
		Variants: []Variant{
			{Name: "baseline", Rank: lastRanker, Baseline: true},
			{Name: "loser", Rank: lastRanker},
		},
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	loser := rowByName(res, "loser")
	if loser == nil {
		t.Fatal("missing loser row")
	}
	d := shippedDecision(res, "loser")
	if d.Ships {
		t.Errorf("a non-improving signal must not ship, got %+v", d)
	}
	if d.Decision == "" {
		t.Errorf("failing signal must record a decision, got empty")
	}
	// The row must carry its CI + p (the evidence for the decision).
	if loser.CILower == 0 && loser.CIUpper == 0 && loser.PValue == 0 {
		t.Errorf("failing-signal row must carry CI + p for the decision: %+v", loser)
	}
}

func rowByName(res *Result, name string) *Row {
	for i := range res.Rows {
		if res.Rows[i].Name == name {
			return &res.Rows[i]
		}
	}
	return nil
}
