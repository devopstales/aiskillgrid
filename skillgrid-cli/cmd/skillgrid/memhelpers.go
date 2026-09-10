package main

import (
	"context"
	"encoding/json"
	"io"
	"regexp"
	"strconv"
	"time"

	"github.com/devopstales/skillgrid/skillgrid-cli/internal/mnemonic/memory"
)

// memoryBudgetCfg is the CLI's tunable read-budget flags (change 013, step 03).
// Zero fields fall back to the memory package defaults.
type memoryBudgetCfg struct {
	Items     int
	Chars     int
	TimeoutNs int64
}

// memoryBudget adapts the CLI budget flags to the memory package's BudgetConfig.
func memoryBudget(c memoryBudgetCfg) memory.BudgetConfig {
	return memory.BudgetConfig{Items: c.Items, Chars: c.Chars, TimeoutNs: c.TimeoutNs}
}

// memoryShareInput adapts the CLI share flags to the memory package's ShareInput.
type memoryShareInput = memory.ShareInput

// hCtx returns the read context for a CLI memory read (the budget's context
// timeout is honored from the caller's deadline; the CLI uses Background).
func hCtx() context.Context {
	return context.Background()
}

func newJSONEncoder(w io.Writer) *json.Encoder {
	return json.NewEncoder(w)
}

var memWindowRe = regexp.MustCompile(`^(\d+)([smhd])$`)

// parseMemWindow parses a timeline window (e.g. "30m", "2h"); default 1h.
func parseMemWindow(s string) time.Duration {
	s = regexp.MustCompile(`\s+`).ReplaceAllString(s, "")
	if s == "" {
		return time.Hour
	}
	if m := memWindowRe.FindStringSubmatch(s); m != nil {
		n, _ := strconv.Atoi(m[1])
		switch m[2] {
		case "s":
			return time.Duration(n) * time.Second
		case "m":
			return time.Duration(n) * time.Minute
		case "h":
			return time.Duration(n) * time.Hour
		case "d":
			return time.Duration(n) * 24 * time.Hour
		}
	}
	return time.Hour
}
