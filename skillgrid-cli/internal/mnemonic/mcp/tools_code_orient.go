package mcp

import (
	"context"

	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// registerOrientTools registers the Tier-1 orientation MCP tools. They resolve
// a symbol and return signature / file TOC / map / list / metadata plus linked
// rationale. Unknown symbols return not-found (no fabricated symbol).
func registerOrientTools(s *server.MCPServer) {
	tools := []struct {
		tool    mcplib.Tool
		handler server.ToolHandlerFunc
	}{
		{codeOrientTool(), handleCodeOrient},
		{codeSignatureTool(), handleCodeSignature},
		{codeFileTOCTool(), handleCodeFileTOC},
		{codeRationaleTool(), handleCodeRationale},
	}
	for _, entry := range tools {
		s.AddTool(entry.tool, entry.handler)
	}
}

func codeOrientTool() mcplib.Tool {
	return mcplib.NewTool("code_orient",
		mcplib.WithDescription("Tier-1 orientation for a resolved symbol: returns the symbol metadata, its signature, the file TOC (all symbols in the file), a list of symbols, and linked rationale. Use after code_search/code_orient to understand where a symbol lives and why it exists. Unknown symbols return not-found, not a fabricated symbol."),
		mcplib.WithString("symbol", mcplib.Required(), mcplib.Description("Symbol name (exact or identifier; camelCase/snake_case both resolve)")),
	)
}

func codeSignatureTool() mcplib.Tool {
	return mcplib.NewTool("code_signature",
		mcplib.WithDescription("Return just the signature and span of a resolved symbol. Unknown symbols return not-found."),
		mcplib.WithString("symbol", mcplib.Required(), mcplib.Description("Symbol name")),
	)
}

func codeFileTOCTool() mcplib.Tool {
	return mcplib.NewTool("code_file_toc",
		mcplib.WithDescription("Return the table of contents (every symbol, in line order) for the file that contains a resolved symbol. Unknown symbols return not-found."),
		mcplib.WithString("symbol", mcplib.Required(), mcplib.Description("Symbol name (resolves to its file)")),
	)
}

func codeRationaleTool() mcplib.Tool {
	return mcplib.NewTool("code_rationale",
		mcplib.WithDescription("Return the rationale comments (# NOTE: / # WHY: / ADR-RFC citations) linked to a resolved symbol's nearest enclosing definition. Unknown symbols return not-found."),
		mcplib.WithString("symbol", mcplib.Required(), mcplib.Description("Symbol name")),
	)
}

func handleCodeOrient(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	svc, projectID, cleanup, err := openService()
	if err != nil {
		return toolError(err)
	}
	defer cleanup()
	symbol, err := req.RequireString("symbol")
	if err != nil {
		return toolError(err)
	}
	out, err := svc.OrientSymbol(ctx, projectID, symbol)
	if err != nil {
		return toolError(err)
	}
	return JSONResult(out)
}

func handleCodeSignature(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	svc, projectID, cleanup, err := openService()
	if err != nil {
		return toolError(err)
	}
	defer cleanup()
	symbol, err := req.RequireString("symbol")
	if err != nil {
		return toolError(err)
	}
	out, err := svc.OrientSymbol(ctx, projectID, symbol)
	if err != nil {
		return toolError(err)
	}
	return JSONResult(map[string]any{
		"found":     out.Found,
		"symbol":    out.Symbol,
		"signature": out.Signature,
		"reason":    out.Reason,
	})
}

func handleCodeFileTOC(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	svc, projectID, cleanup, err := openService()
	if err != nil {
		return toolError(err)
	}
	defer cleanup()
	symbol, err := req.RequireString("symbol")
	if err != nil {
		return toolError(err)
	}
	out, err := svc.OrientSymbol(ctx, projectID, symbol)
	if err != nil {
		return toolError(err)
	}
	return JSONResult(map[string]any{
		"found":    out.Found,
		"file_toc": out.FileTOC,
		"reason":   out.Reason,
	})
}

func handleCodeRationale(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	svc, projectID, cleanup, err := openService()
	if err != nil {
		return toolError(err)
	}
	defer cleanup()
	symbol, err := req.RequireString("symbol")
	if err != nil {
		return toolError(err)
	}
	out, err := svc.OrientSymbol(ctx, projectID, symbol)
	if err != nil {
		return toolError(err)
	}
	return JSONResult(map[string]any{
		"found":     out.Found,
		"rationale": out.Rationale,
		"reason":    out.Reason,
	})
}
