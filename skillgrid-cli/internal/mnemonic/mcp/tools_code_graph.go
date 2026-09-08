package mcp

import (
	"context"
	"fmt"

	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/devopstales/skillgrid/skillgrid-cli/internal/mnemonic/service"
)

// registerGraphTools registers the Tier-2 graph MCP tools. They are part of
// the narrow menu (unlisted by default on the MCP surface; re-enabled via
// CodeToolsEnvVar) and are all exposed by the CLI.
func registerGraphTools(s *server.MCPServer) {
	tools := []struct {
		tool    mcplib.Tool
		handler server.ToolHandlerFunc
	}{
		{neighborTool("code_get_callers", "callers"), handleCodeNeighbors},
		{neighborTool("code_get_callees", "callees"), handleCodeNeighbors},
		{neighborTool("code_get_dependents", "dependents"), handleCodeNeighbors},
		{neighborTool("code_get_implementors", "implementors"), handleCodeNeighbors},
		{neighborTool("code_get_hierarchy", "hierarchy"), handleCodeNeighbors},
		{neighborTool("code_get_tests_for", "tests_for"), handleCodeNeighbors},
		{codePathTool(), handleCodePath},
		{codeExplainTool(), handleCodeExplain},
		{codeImpactTool(), handleCodeImpact},
	}
	for _, entry := range tools {
		s.AddTool(entry.tool, entry.handler)
	}
}

func neighborTool(name, view string) mcplib.Tool {
	desc := map[string]string{
		"code_get_callers":      "Symbols that call/import the given symbol (reverse edges). Every edge carries a confidence label; ambiguous name-only resolution returns AMBIGUOUS edges, never a silent drop.",
		"code_get_callees":      "Symbols the given symbol calls/imports (forward edges). Every edge carries a confidence label.",
		"code_get_dependents":   "Symbols that depend on the given symbol (calls/imports in either direction). Every edge carries a confidence label.",
		"code_get_implementors": "Symbols implementing the given interface/type.",
		"code_get_hierarchy":    "Structural hierarchy (extends/implements) around the given symbol.",
		"code_get_tests_for":    "Test relationships for the given symbol.",
	}[name]
	_ = view
	return mcplib.NewTool(name,
		mcplib.WithDescription(desc),
		mcplib.WithString("symbol", mcplib.Required(), mcplib.Description("Symbol name")),
	)
}

func codePathTool() mcplib.Tool {
	return mcplib.NewTool("code_path",
		mcplib.WithDescription("Shortest edge path between two symbols (each hop confidence-labeled). When no static path exists it returns a 'where the graph stops' answer: the dispatch kind that ended the flow, its line, and the refused name-only matches with their confidence. Never empty, never fabricated."),
		mcplib.WithString("from", mcplib.Required(), mcplib.Description("Source symbol name")),
		mcplib.WithString("to", mcplib.Required(), mcplib.Description("Destination symbol name")),
	)
}

func codeExplainTool() mcplib.Tool {
	return mcplib.NewTool("code_explain",
		mcplib.WithDescription("Explain a symbol: the node, its degree, and all its connections ranked by the neighbor's degree (hubs first)."),
		mcplib.WithString("symbol", mcplib.Required(), mcplib.Description("Symbol name")),
	)
}

func codeImpactTool() mcplib.Tool {
	return mcplib.NewTool("code_impact",
		mcplib.WithDescription("Risk-tiered blast radius for a symbol: depth 1 dependents are WILL BREAK, deeper are LIKELY AFFECTED. Every edge is confidence-tagged and min_confidence filters low-confidence hops. A target matching several symbols returns a ranked candidate list (never a silent pick); narrow with file/uid/kind."),
		mcplib.WithString("symbol", mcplib.Required(), mcplib.Description("Symbol name (or ambiguous name to disambiguate)")),
		mcplib.WithString("file", mcplib.Description("Narrow to a symbol in this file path")),
		mcplib.WithString("uid", mcplib.Description("Narrow to a symbol by UID")),
		mcplib.WithString("kind", mcplib.Description("Narrow to a symbol kind (function, method, class, ...)")),
		mcplib.WithString("min_confidence", mcplib.Description("Minimum edge confidence: EXTRACTED | INFERRED | AMBIGUOUS")),
		mcplib.WithNumber("max_depth", mcplib.Description("Bound the traversal depth (default 5)")),
	)
}

func neighborView(name string) (string, bool) {
	views := map[string]string{
		"code_get_callers":      "callers",
		"code_get_callees":      "callees",
		"code_get_dependents":   "dependents",
		"code_get_implementors": "implementors",
		"code_get_hierarchy":    "hierarchy",
		"code_get_tests_for":    "tests_for",
	}
	v, ok := views[name]
	return v, ok
}

func handleCodeNeighbors(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	symbol, err := req.RequireString("symbol")
	if err != nil {
		return toolError(err)
	}
	svc, projectID, cleanup, err := openService()
	if err != nil {
		return toolError(err)
	}
	defer cleanup()

	name := ""
	if n, _ := req.RequireString("symbol"); n != "" {
		name = n
	}
	toolName := req.Params.Name
	view, ok := neighborView(toolName)
	if !ok {
		return toolError(fmt.Errorf("unknown neighbor view for %s", toolName))
	}
	out, err := svc.GrabNeighbors(ctx, projectID, view, symbol)
	_ = name
	if err != nil {
		return toolError(err)
	}
	return JSONResult(out)
}

func handleCodePath(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	from, err := req.RequireString("from")
	if err != nil {
		return toolError(err)
	}
	to, err := req.RequireString("to")
	if err != nil {
		return toolError(err)
	}
	svc, projectID, cleanup, err := openService()
	if err != nil {
		return toolError(err)
	}
	defer cleanup()
	out, err := svc.CodePath(ctx, projectID, from, to)
	if err != nil {
		return toolError(err)
	}
	return JSONResult(out)
}

func handleCodeExplain(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	symbol, err := req.RequireString("symbol")
	if err != nil {
		return toolError(err)
	}
	svc, projectID, cleanup, err := openService()
	if err != nil {
		return toolError(err)
	}
	defer cleanup()
	out, err := svc.CodeExplain(ctx, projectID, symbol)
	if err != nil {
		return toolError(err)
	}
	return JSONResult(out)
}

func handleCodeImpact(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	symbol, err := req.RequireString("symbol")
	if err != nil {
		return toolError(err)
	}
	file, _ := req.RequireString("file")
	uid, _ := req.RequireString("uid")
	kind, _ := req.RequireString("kind")
	minConf, _ := req.RequireString("min_confidence")
	maxDepth := int(req.GetFloat("max_depth", 0))

	svc, projectID, cleanup, err := openService()
	if err != nil {
		return toolError(err)
	}
	defer cleanup()
	opts := service.ImpactOptions{
		File:          file,
		UID:           uid,
		Kind:          kind,
		MinConfidence: minConf,
		MaxDepth:      maxDepth,
	}
	out, err := svc.CodeImpact(ctx, projectID, symbol, opts)
	if err != nil {
		return toolError(err)
	}
	return JSONResult(out)
}
