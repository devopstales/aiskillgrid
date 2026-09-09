package mcp

import (
	"context"
	"fmt"

	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/devopstales/skillgrid/skillgrid-cli/internal/mnemonic/service"
)

// openServiceForRepo opens the service + project for a request's optional
// repo param: an explicit repo wins, otherwise the cwd-resolved project.
// Returns a cleanup func to close the store handle.
func openServiceForRepo(repo string) (*service.Service, string, func(), error) {
	svc, err := rootService()
	if err != nil {
		return nil, "", nil, err
	}
	projectID, err := projectIDForPool(svc, repo)
	if err != nil {
		return nil, "", nil, err
	}
	_, cleanup, err := svc.OpenForDirectory(".")
	if err != nil {
		return nil, "", nil, err
	}
	return svc, projectID, cleanup, nil
}

// registerHybridTools registers the offline hybrid search surface:
// code_hybrid_search (FTS + deterministic signals + optional embeddings,
// per-signal provenance), code_semantic_search (vector leg; symbol-level hits
// named), and code_embedding_status (provider/model/coverage). All are
// distinct code_* tools — none clashes with memory semantic_search.
func registerHybridTools(s *server.MCPServer) {
	tools := []struct {
		tool    mcplib.Tool
		handler server.ToolHandlerFunc
	}{
		{codeHybridSearchTool(), handleCodeHybridSearch},
		{codeSemanticSearchTool(), handleCodeSemanticSearch},
		{codeEmbeddingStatusTool(), handleCodeEmbeddingStatus},
	}
	for _, entry := range tools {
		s.AddTool(entry.tool, entry.handler)
	}
}

func codeHybridSearchTool() mcplib.Tool {
	return mcplib.NewTool("code_hybrid_search",
		mcplib.WithDescription("Offline hybrid code search: fuses identifier/chunk FTS with deterministic signals (proximity, TF-IDF, type/API) and — when the embedder is available — semantic vectors via RRF. Every hit carries per-signal provenance. Embeddings are optional: a down/missing embedder degrades to FTS + signals, never a hard fail. Distinct from memory semantic_search."),
		mcplib.WithString("query", mcplib.Required(), mcplib.Description("Search query (identifiers or free text)")),
		mcplib.WithNumber("limit", mcplib.Description("Maximum hits (default 20)")),
		mcplib.WithString("repo", mcplib.Description("Optional project/repo name; omitted when the cwd resolves it")),
	)
}

func codeSemanticSearchTool() mcplib.Tool {
	return mcplib.NewTool("code_semantic_search",
		mcplib.WithDescription("Semantic (embedding) code search. Symbol-level hits return the named symbol with file + line; chunk-level hits return the line range for non-symbol code. Requires an active embedder; with none configured it returns an empty result (the hybrid tool keeps the FTS+signals floor)."),
		mcplib.WithString("query", mcplib.Required(), mcplib.Description("Free-text semantic query")),
		mcplib.WithNumber("limit", mcplib.Description("Maximum hits (default 20)")),
		mcplib.WithString("repo", mcplib.Description("Optional project/repo name; omitted when the cwd resolves it")),
	)
}

func codeEmbeddingStatusTool() mcplib.Tool {
	return mcplib.NewTool("code_embedding_status",
		mcplib.WithDescription("Embedding index status: active provider (onnx/external/off), model name, vector dimension, embedded symbol count, and indexed model (model-swap guard). Read-only."),
		mcplib.WithString("repo", mcplib.Description("Optional project/repo name; omitted when the cwd resolves it")),
	)
}

func handleCodeHybridSearch(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	query, err := req.RequireString("query")
	if err != nil {
		return toolError(err)
	}
	if query == "" {
		return toolError(fmt.Errorf("code_hybrid_search: 'query' is required and must be non-empty"))
	}
	limit := int(req.GetFloat("limit", 20))

	svc, projectID, cleanup, err := openServiceForRepo(req.GetString("repo", ""))
	if err != nil {
		return toolError(err)
	}
	defer cleanup()

	out, err := svc.CodeHybridSearch(ctx, projectID, query, limit)
	if err != nil {
		return toolError(err)
	}
	return JSONResult(out)
}

func handleCodeSemanticSearch(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	query, err := req.RequireString("query")
	if err != nil {
		return toolError(err)
	}
	if query == "" {
		return toolError(fmt.Errorf("code_semantic_search: 'query' is required and must be non-empty"))
	}
	limit := int(req.GetFloat("limit", 20))

	svc, projectID, cleanup, err := openServiceForRepo(req.GetString("repo", ""))
	if err != nil {
		return toolError(err)
	}
	defer cleanup()

	out, err := svc.CodeSemanticSearch(ctx, projectID, query, limit)
	if err != nil {
		return toolError(err)
	}
	return JSONResult(out)
}

func handleCodeEmbeddingStatus(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
	svc, projectID, cleanup, err := openServiceForRepo(req.GetString("repo", ""))
	if err != nil {
		return toolError(err)
	}
	defer cleanup()

	out, err := svc.CodeEmbeddingStatus(ctx, projectID)
	if err != nil {
		return toolError(err)
	}
	return JSONResult(out)
}
