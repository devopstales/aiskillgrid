package mcp

import (
	"sort"
	"testing"

	mcplib "github.com/mark3labs/mcp-go/mcp"
)

// TestCodeSearchSchemaStable locks the code_search tool's name and required
// `query` param schema so it cannot silently regress when new tools are added.
// Threat: Mnemonic tool surface. The existing four code_* tools must keep
// their name and required-param schema (additive fields only).
func TestCodeSearchSchemaStable(t *testing.T) {
	tool := codeSearchTool()
	if tool.Name != "code_search" {
		t.Fatalf("code_search tool name changed to %q", tool.Name)
	}
	schema := tool.InputSchema
	// `query` must be a required string parameter.
	props, ok := schema.Properties["query"]
	if !ok {
		t.Fatalf("code_search input schema is missing the required 'query' property: %+v", schema.Properties)
	}
	if _, isMap := props.(map[string]any); !isMap {
		t.Fatalf("code_search 'query' property is not a JSON schema object: %T", props)
	}
	if !containsString(schema.Required, "query") {
		t.Fatalf("code_search 'query' is not a required parameter; required=%v", schema.Required)
	}
	// The other three existing tools keep their names + required params.
	assertCodeToolStable(t, codeIndexTool(), "code_index", nil)
	assertCodeToolStable(t, codeStatusTool(), "code_status", nil)
	assertCodeToolStable(t, codeReadTool(), "code_read", []string{"path"})
}

// assertCodeToolStable locks a tool's name and its required parameters.
func assertCodeToolStable(t *testing.T, tool mcplib.Tool, name string, required []string) {
	t.Helper()
	if tool.Name != name {
		t.Errorf("tool name changed: got %q, want %q", tool.Name, name)
	}
	got := append([]string(nil), tool.InputSchema.Required...)
	sort.Strings(got)
	sort.Strings(required)
	if len(got) != len(required) {
		t.Errorf("%s required params = %v, want %v", name, got, required)
	}
	for i := range required {
		if got[i] != required[i] {
			t.Errorf("%s required params = %v, want %v", name, got, required)
			break
		}
	}
}

func containsString(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}
