package mcp

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/devopstales/skillgrid/skillgrid-cli/internal/mnemonic/service"
)

// CodeToolsEnvVar is the config knob that re-enables the narrow code_* menu
// tools on the MCP surface (like CodeGraph's CODEGRAPH_MCP_TOOLS). By default
// the menu is unlisted; the composite code_explore stays listed.
const CodeToolsEnvVar = "SKILLGRID_MCP_CODE_TOOLS"

// menuCodeTools are the narrow graph/orientation tools demoted to unlisted-by-
// default on the MCP surface. They stay registered (callable) and the CLI
// exposes all of them.
var menuCodeTools = []string{
	"code_orient", "code_signature", "code_file_toc", "code_rationale",
	"code_get_callers", "code_get_callees", "code_get_dependents",
	"code_get_implementors", "code_get_hierarchy", "code_get_tests_for",
	"code_path", "code_explain",
}

// exploreInitializeGuidance is injected at MCP initialize (server
// instructions): steer agents to the composite tool and keep them off grep
// loops.
func exploreInitializeGuidance() string {
	return "skillgrid-mnemonic: code_explore is the primary code-intelligence tool. " +
		"Answer structural questions (callers, callees, call paths, blast radius) directly " +
		"with the graph and treat the source it returns as already read; don't re-verify " +
		"with grep. The narrow code_* menu tools are unlisted by default and re-enable via " +
		CodeToolsEnvVar + " (the CLI still exposes all of them)."
}

// codeExploreTool is the documented primary MCP tool: one call returns the
// relevant symbols' verbatim source grouped by file, the call paths between
// them (including INFERRED dynamic-dispatch hops), and a blast-radius summary.
func codeExploreTool() mcplib.Tool {
	return mcplib.NewTool("code_explore",
		mcplib.WithDescription("Primary code-intelligence tool. One call returns the relevant symbols' verbatim source grouped by file, the call paths between them (including INFERRED dynamic-dispatch hops), and a blast-radius summary. Answer structural questions directly with the graph; treat returned source as already read and don't re-verify with grep. Optional maxTokens bounds the response (truncated with an ellipsis, stays well-formed)."),
		mcplib.WithString("symbol", mcplib.Required(), mcplib.Description("Symbol name to explore (callers, callees, and source)")),
		mcplib.WithString("repo", mcplib.Description("Optional project/repo name; omitted when a single project is indexed or the MCP default/cwd resolves it")),
		mcplib.WithNumber("max_tokens", mcplib.Description("Optional response budget in tokens (deterministic ~4-bytes/token estimate); the formatted response is truncated with an ellipsis and stays valid when exceeded")),
	)
}

// registerExploreTools registers the composite primary tool.
func registerExploreTools(s *server.MCPServer) {
	s.AddTool(codeExploreTool(), handleCodeExplore)
}

// applyCodeToolFilter returns a ToolFilterFunc that hides the narrow menu
// tools unless CodeToolsEnvVar re-enables them. The composite tool and the
// four stable chunk tools are always listed.
func applyCodeToolFilter() server.ToolFilterFunc {
	return func(_ context.Context, tools []mcplib.Tool) []mcplib.Tool {
		if codeToolsEnabled() {
			return tools
		}
		hidden := map[string]bool{}
		for _, name := range menuCodeTools {
			hidden[name] = true
		}
		out := make([]mcplib.Tool, 0, len(tools))
		for _, t := range tools {
			if hidden[t.Name] {
				continue
			}
			out = append(out, t)
		}
		return out
	}
}

// codeToolsEnabled reports whether the narrow menu tools are re-enabled.
func codeToolsEnabled() bool {
	v := strings.TrimSpace(os.Getenv(CodeToolsEnvVar))
	return v == "all" || v == "1" || strings.EqualFold(v, "true") || strings.EqualFold(v, "yes")
}

// exploreResult is the composite answer: source by file, call-flow between the
// returned symbols, and a blast-radius summary.
type exploreResult struct {
	Symbol     string                   `json:"symbol"`
	Source     map[string][]srcSpan     `json:"source"`
	CallFlow   []flowEdge               `json:"call_flow"`
	Impact     *service.ImpactResultDTO `json:"impact"`
	Truncated  bool                     `json:"truncated,omitempty"`
}

type srcSpan struct {
	Symbol   string `json:"symbol"`
	Start    int    `json:"start_line"`
	End      int    `json:"end_line"`
	Content  string `json:"content"`
}

type flowEdge struct {
	From       string `json:"from"`
	To         string `json:"to"`
	Kind       string `json:"kind"`
	Confidence string `json:"confidence"`
	Line       int    `json:"line"`
}

// handleCodeExplore runs the composite explore: resolve the symbol, gather its
// neighbors, pull verbatim source for the touched files, derive the call-flow
// (including INFERRED hops), and attach a blast-radius summary.
func handleCodeExplore(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	symbol, err := req.RequireString("symbol")
	if err != nil {
		return toolError(err)
	}
	_, _ = req.RequireString("repo") // optional; pool/cwd resolves it
	maxTokens := int(req.GetFloat("max_tokens", 0))

	svc, projectID, cleanup, err := openService()
	if err != nil {
		return toolError(err)
	}
	defer cleanup()

	res, err := svc.CodeExplore(ctx, projectID, symbol)
	if err != nil {
		return toolError(err)
	}

	// Build the composite answer.
	out := exploreResult{Symbol: symbol, Source: map[string][]srcSpan{}}
	for path, spans := range res.Source {
		for _, sp := range spans {
			out.Source[path] = append(out.Source[path], srcSpan{
				Symbol:  sp.Symbol,
				Start:   sp.StartLine,
				End:     sp.EndLine,
				Content: sp.Content,
			})
		}
	}
	for _, e := range res.Flow {
		out.CallFlow = append(out.CallFlow, flowEdge{
			From:       e.From,
			To:         e.To,
			Kind:       e.Kind,
			Confidence: e.Confidence,
			Line:       e.Line,
		})
	}
	out.Impact = res.Impact

	text, truncated := formatExplore(out, maxTokens)
	resText := exploreEnvelope{text, truncated}
	return JSONResult(resText)
}

// exploreEnvelope is the response body: the formatted text plus a truncation
// flag, so a budgeted response stays well-formed JSON.
type exploreEnvelope struct {
	Text      string `json:"text"`
	Truncated bool   `json:"truncated,omitempty"`
}

// formatExplore renders the composite answer as a single readable text block
// and, when maxTokens (a deterministic ~4-bytes/token budget) is exceeded,
// truncates it with an ellipsis so the response stays valid and bounded.
func formatExplore(out exploreResult, maxTokens int) (string, bool) {
	var b strings.Builder
	files := make([]string, 0, len(out.Source))
	for f := range out.Source {
		files = append(files, f)
	}
	sort.Strings(files)
	b.WriteString("## Source\n")
	for _, f := range files {
		b.WriteString("\n### " + f + "\n")
		for _, sp := range out.Source[f] {
			fmt.Fprintf(&b, "-- %s (%d-%d)\n%s\n", sp.Symbol, sp.Start, sp.End, sp.Content)
		}
	}
	if len(out.CallFlow) > 0 {
		b.WriteString("\n## Call flow\n")
		for _, e := range out.CallFlow {
			fmt.Fprintf(&b, "%s -> %s [%s %s] :%d\n", e.From, e.To, e.Kind, e.Confidence, e.Line)
		}
	}
	if out.Impact != nil {
		b.WriteString("\n## Blast radius\n")
		b.WriteString(out.Impact.Summary())
	}
	text := b.String()
	budget := maxTokens * 4
	if maxTokens > 0 && len(text) > budget {
		return text[:budget] + "\n…", true
	}
	return text, false
}
