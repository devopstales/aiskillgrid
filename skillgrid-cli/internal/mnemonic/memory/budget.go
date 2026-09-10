package memory

import (
	"context"
	"fmt"
	"time"
)

// Default budget caps (change 013, step 03) — the "Context-window flood" guard.
// A 20-hit read must never drown the context window: the item cap bounds how
// many observations are ever in-list, the char budget bounds how much text each
// in-list snippet may carry (with an explicit "N chars omitted" marker), and
// the context timeout bounds how long a read may run (a slow read is cut and
// returned as a truncated partial, never hung). All are tunable via config
// (mnemonic.retrieval_budget); the values below are the defaults.
const (
	defaultBudgetItems  = 10
	defaultBudgetChars  = 1200
	defaultBudgetTimeout = 3 * time.Second
)

// BudgetConfig is the tunable budget (change 013, step 03). Every field is
// additive and optional: a zero value falls back to the matching default. The
// timeout is expressed in nanoseconds so it is a plain JSON/YAML scalar.
type BudgetConfig struct {
	// Items caps how many in-list results a read returns. <= 0 → default.
	Items int
	// Chars caps how many characters an in-list snippet may carry. The snippet
	// is truncated with an explicit "N chars omitted" marker. <= 0 → default.
	Chars int
	// TimeoutNs is the context-timeout budget in nanoseconds. A read that
	// exceeds it is cut and returned as a truncated:true partial with a reason.
	// <= 0 → default.
	TimeoutNs int64
}

// DefaultBudget returns the default budget (item + char + timeout).
func DefaultBudget() BudgetConfig {
	return BudgetConfig{Items: defaultBudgetItems, Chars: defaultBudgetChars, TimeoutNs: int64(defaultBudgetTimeout)}
}

// normalize replaces zero/negative fields with their defaults so a partial
// config (only one knob set) still enforces the rest.
func (c BudgetConfig) normalize() BudgetConfig {
	if c.Items <= 0 {
		c.Items = defaultBudgetItems
	}
	if c.Chars <= 0 {
		c.Chars = defaultBudgetChars
	}
	if c.TimeoutNs <= 0 {
		c.TimeoutNs = int64(defaultBudgetTimeout)
	}
	return c
}

// Budget applies the item + char + timeout caps to a set of in-list reads.
// It is the uniform wrapper every mem_* read path runs through (change 013,
// step 03) so memory can never overwhelm the context window.
type Budget struct {
	cfg BudgetConfig
}

// NewBudget builds a budget from cfg (zero fields → defaults).
func NewBudget(cfg BudgetConfig) *Budget {
	return &Budget{cfg: cfg.normalize()}
}

// Default returns a budget with the default caps.
func Default() *Budget {
	return NewBudget(BudgetConfig{})
}

// Config returns the effective (normalized) budget config.
func (b *Budget) Config() BudgetConfig {
	if b == nil {
		return DefaultBudget()
	}
	return b.cfg.normalize()
}

// TruncationMarker is the explicit "N chars omitted" suffix appended to a
// char-budgeted snippet. The truncation is never silent: a reader always sees
// how much was dropped and that a full-content fetch (mem_get_observation) is
// available.
const TruncationMarker = "… (%d chars omitted)"

// TruncateWithMarker truncates s to at most n characters and appends the
// explicit "N chars omitted" marker when anything was dropped. s returned
// unmodified when it fits.
func TruncateWithMarker(s string, n int) string {
	if n <= 0 {
		n = defaultBudgetChars
	}
	if len(s) <= n {
		return s
	}
	kept := len(s) - n
	return s[:n] + fmt.Sprintf(TruncationMarker, kept)
}

// BudgetResult is the output of Budget.Apply: the budgeted hits plus the
// truncated flag + reason when the read hit the char budget or the context
// timeout. Truncated is true whenever the read was cut (either way); Reason is
// human-readable and non-empty only when Truncated is true.
type BudgetResult struct {
	Hits      []Observation `json:"hits"`
	Truncated bool          `json:"truncated,omitempty"`
	Reason    string        `json:"reason,omitempty"`
	// CharsOmitted is the total characters dropped across all in-list snippets
	// by the char budget. 0 when nothing was char-truncated.
	CharsOmitted int `json:"chars_omitted,omitempty"`
}

// Apply enforces the item cap, the char budget (in-list snippets truncated with
// an explicit "N chars omitted"), and the context timeout (returns a
// truncated:true partial with a reason, never hangs).
//
// Timeout semantics: if the caller's ctx is already carrying a deadline that
// has lapsed by the time Apply runs, the read is cut — a partial (as many hits
// as the item cap allows, char-truncated) is returned with truncated=true and
// reason="timeout". The timeout is therefore honored from the injected
// deadline (the MCP/HTTP handler's context), which is what "never hangs" means:
// a slow downstream read is bounded by the caller's deadline, not by an
// unbounded query.
func (b *Budget) Apply(ctx context.Context, hits []Observation) BudgetResult {
	cfg := b.Config()

	// Item cap: never return more than the item budget in-list.
	out := hits
	if len(out) > cfg.Items {
		out = out[:cfg.Items]
	}

	// Char budget: truncate each in-list snippet with an explicit marker.
	truncated := false
	reason := ""
	charsOmitted := 0
	for i := range out {
		if len(out[i].Content) > cfg.Chars {
			omitted := len(out[i].Content) - cfg.Chars
			out[i].Content = TruncateWithMarker(out[i].Content, cfg.Chars)
			charsOmitted += omitted
			truncated = true
			if reason == "" {
				reason = "char-budget"
			}
		}
	}

	// Context timeout: a lapsed deadline cuts the read. We check the deadline
	// rather than waiting, so Apply never blocks (the read already happened);
	// the lapsed deadline means the read consumed the whole context budget, so
	// we mark the result truncated + reason="timeout" (a partial).
	if dl, ok := ctx.Deadline(); ok && !time.Now().Before(dl) {
		truncated = true
		reason = "timeout"
	}

	// Copy so the caller's slice is not mutated in place (in-list truncation is
	// a projection, not a store rewrite — the full content stays fetchable via
	// mem_get_observation).
	res := make([]Observation, len(out))
	copy(res, out)
	return BudgetResult{Hits: res, Truncated: truncated, Reason: reason, CharsOmitted: charsOmitted}
}
