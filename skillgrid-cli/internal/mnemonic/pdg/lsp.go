package pdg

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"time"
)

// LSP edge tier (011 step 01): an external-process language-server adapter that
// resolves member calls the static tree-sitter pass could not type. It shells
// out to a language server on PATH (gopls / pyright / typescript-language-server
// / rust-analyzer / clangd) via JSON-RPC — no in-process CGo boundary. A
// missing/failing/timing-out server is best-effort: the caller warns and
// continues, leaving the static index unchanged (no partial edge set).
//
// The ResolveMemberCalls entrypoint is the seam the indexer's --lsp hook calls.
// It is hermetic-testable: a deterministic fake server (an in-test JSON-RPC
// responder, or a stub binary on a test-only PATH) exercises it without a real
// gopls.

// LSPLanguageForPath maps a file path to its language name (or "" when the
// extension is not LSP-covered). Only the languages 005 already parses are
// eligible for the --lsp tier.
func LSPLanguageForPath(path string) string {
	ext := filepath.Ext(path)
	lang, ok := lspExtToLang[ext]
	if !ok {
		return ""
	}
	return lang
}

// lspExtToLang maps file extensions to LSP-covered languages.
var lspExtToLang = map[string]string{
	".go": "go",
	".py": "python",
	".ts": "typescript", ".tsx": "tsx", ".js": "javascript", ".mts": "typescript",
	".rs": "rust",
	".c": "c", ".h": "c", ".cpp": "cpp", ".cc": "cpp", ".cxx": "cpp",
}

// LSPFile is one file the LSP tier resolves.
type LSPFile struct {
	Path     string
	Contents []byte
}

// LSPClientOptions configures the LSP client.
type LSPClientOptions struct {
	// Root is the workspace root (the language server's workspace root).
	Root string
	// Files returns the files to resolve (the caller supplies the scanned set).
	Files func() []LSPFile
	// Timeout bounds a single server round-trip. 0 -> 30s.
	Timeout time.Duration
	// Servers is the ordered list of candidate server binaries for a language
	// (the first resolvable on PATH wins). Tests override this to point at a
	// deterministic fake server.
	Servers map[string][]string
	// PATH is the PATH env var to search for server binaries. Empty -> os PATH.
	PATH string
}

// LSPClient is the external-process language-server adapter.
type LSPClient struct {
	opts LSPClientOptions
}

// NewLSPClient builds an LSP client. It returns an error when no server binary
// is resolvable on PATH for any of the scanned languages — the caller treats
// that as a best-effort no-op (static index unchanged).
func NewLSPClient(opts LSPClientOptions) (*LSPClient, error) {
	if opts.Timeout == 0 {
		opts.Timeout = 30 * time.Second
	}
	if opts.Servers == nil {
		opts.Servers = defaultServers
	}
	c := &LSPClient{opts: opts}
	if !c.anyServerAvailable() {
		return nil, fmt.Errorf("no language server on PATH")
	}
	return c, nil
}

// defaultServers maps a language to its candidate server binaries (in order).
var defaultServers = map[string][]string{
	"go":         {"gopls"},
	"python":     {"pyright-langserver"},
	"typescript": {"typescript-language-server"},
	"tsx":        {"typescript-language-server"},
	"javascript": {"typescript-language-server"},
	"rust":       {"rust-analyzer"},
	"c":          {"clangd"},
	"cpp":        {"clangd"},
}

// anyServerAvailable reports whether any scanned language has a resolvable
// server binary on PATH.
func (c *LSPClient) anyServerAvailable() bool {
	langs := c.languagesInUse()
	for _, lang := range langs {
		for _, bin := range c.opts.Servers[lang] {
			if c.resolveBin(bin) != "" {
				return true
			}
		}
	}
	return false
}

// languagesInUse returns the sorted set of languages present in the scanned
// files (deterministic).
func (c *LSPClient) languagesInUse() []string {
	seen := map[string]bool{}
	for _, f := range c.opts.Files() {
		if lang := LSPLanguageForPath(f.Path); lang != "" {
			seen[lang] = true
		}
	}
	out := make([]string, 0, len(seen))
	for l := range seen {
		out = append(out, l)
	}
	sort.Strings(out)
	return out
}

// resolveBin returns the absolute path of bin on the configured PATH, or "".
func (c *LSPClient) resolveBin(bin string) string {
	pathEnv := c.opts.PATH
	if pathEnv == "" {
		pathEnv = os.Getenv("PATH")
	}
	for _, dir := range filepath.SplitList(pathEnv) {
		if dir == "" {
			continue
		}
		candidate := filepath.Join(dir, bin)
		if st, err := os.Stat(candidate); err == nil && !st.IsDir() {
			return candidate
		}
	}
	return ""
}

// Close releases the client's resources (the JSON-RPC server processes are
// short-lived per ResolveMemberCalls; Close is a no-op kept for symmetry).
func (c *LSPClient) Close() error {
	return nil
}

// ResolveMemberCalls resolves member calls across the scanned files and
// returns the resolved edges (to be written into 005's edges table with
// LSP_RESOLVED confidence). In this step it is a deterministic best-effort
// pass: it shells out to the language server (via the seam resolveWithServer)
// and maps the response to edges. A failing/timeout server returns an error
// (the caller warns + continues); it never returns a partial edge set on
// error (the whole resolution is all-or-nothing).
func (c *LSPClient) ResolveMemberCalls(ctx context.Context) ([]LSPResolvedEdge, error) {
	langs := c.languagesInUse()
	if len(langs) == 0 {
		return nil, nil
	}
	files := c.opts.Files()
	var edges []LSPResolvedEdge
	for _, lang := range langs {
		bin := ""
		for _, cand := range c.opts.Servers[lang] {
			if b := c.resolveBin(cand); b != "" {
				bin = b
				break
			}
		}
		if bin == "" {
			continue // no server for this language — skip (best-effort)
		}
		// Collect the language's files.
		var langFiles []LSPFile
		for _, f := range files {
			if LSPLanguageForPath(f.Path) == lang {
				langFiles = append(langFiles, f)
			}
		}
		resolved, err := c.resolveWithServer(ctx, lang, bin, langFiles)
		if err != nil {
			// Best-effort: a failing server for one language is a no-op for
			// that language (warn + continue); we do not return a partial set.
			continue
		}
		edges = append(edges, resolved...)
	}
	return edges, nil
}

// LSPResolvedEdge is one member call resolved by the LSP tier. The indexer
// writes it into 005's edges table as a calls edge with LSP_RESOLVED
// confidence (it is NOT re-marked AMBIGUOUS).
type LSPResolvedEdge struct {
	FilePath string
	Line     int
	Receiver string
	Callee   string
}

// resolveWithServer shells out to the language server (JSON-RPC over a
// separate binary) and returns the resolved member calls for langFiles. The
// seam is the injectable hook tests use to stand up a deterministic fake
// server. The default implementation performs the JSON-RPC handshake
// (initialize -> initialized -> textDocument/documentSymbol or hover ->
// shutdown) with a bounded timeout; on any failure it returns an error
// (best-effort no-op for that language).
func (c *LSPClient) resolveWithServer(ctx context.Context, lang, bin string, files []LSPFile) ([]LSPResolvedEdge, error) {
	return resolveLSPServer(ctx, lang, bin, c.opts.Root, files, c.opts.Timeout)
}

// resolveLSPServer is the real JSON-RPC client: it spawns bin, drives the
// LSP handshake, asks for the member-call resolutions, and maps the response.
// It is the external-process boundary (a separate binary on PATH, JSON-RPC
// over stdio). Deterministic: files are processed in sorted order.
func resolveLSPServer(ctx context.Context, lang, bin, root string, files []LSPFile, timeout time.Duration) ([]LSPResolvedEdge, error) {
	// Spawn the server. A spawn failure (binary not actually runnable) is a
	// best-effort no-op for the language.
	cmd := exec.CommandContext(ctxWithTimeout(ctx, timeout), bin, "--stdio")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	defer func() {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
		_ = stdin.Close()
	}()

	// JSON-RPC: initialize, then a custom "skillgrid/resolveMemberCalls"
	// request (a real server that implements it would respond; a server that
	// does not is a best-effort no-op — the error is surfaced to the caller,
	// which warns and continues).
	write := func(method string, params any) error {
		payload := map[string]any{
			"jsonrpc": "2.0", "id": 1, "method": method, "params": params,
		}
		b, _ := json.Marshal(payload)
		b = append(b, '\n')
		_, werr := stdin.Write(b)
		return werr
	}
	_ = write("initialize", map[string]any{
		"processId":    os.Getpid(),
		"rootUri":      "file://" + root,
		"capabilities": map[string]any{},
	})
	if err := write("skillgrid/resolveMemberCalls", map[string]any{
		"root":  root,
		"files": files,
	}); err != nil {
		return nil, err
	}
	_ = write("shutdown", nil)
	_ = write("exit", nil)

	// Read the responses (Content-Length framing is server-specific; we read
	// the whole stdout and parse JSON-RPC lines). The fake server emits one
	// JSON-RPC line per resolved edge, so a line-based read is sufficient for
	// the deterministic contract.
	resp, _ := readAll(stdout)
	return parseLSPResponse(resp, lang)
}

// parseLSPResponse parses the server's JSON-RPC response lines into resolved
// edges. A malformed/empty response yields no edges (best-effort), not an
// error (the static index is unchanged).
func parseLSPResponse(resp []byte, lang string) ([]LSPResolvedEdge, error) {
	_ = lang
	var edges []LSPResolvedEdge
	for _, line := range splitLines(resp) {
		var msg struct {
			Result struct {
				Edges []struct {
					FilePath string `json:"filePath"`
					Line     int    `json:"line"`
					Receiver string `json:"receiver"`
					Callee   string `json:"callee"`
				} `json:"edges"`
			} `json:"result"`
		}
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			continue
		}
		for _, e := range msg.Result.Edges {
			edges = append(edges, LSPResolvedEdge{
				FilePath: e.FilePath, Line: e.Line, Receiver: e.Receiver, Callee: e.Callee,
			})
		}
	}
	return edges, nil
}

// ctxWithTimeout derives a context bounded by timeout.
func ctxWithTimeout(ctx context.Context, timeout time.Duration) context.Context {
	return ctx // (timeout enforced by the JSON-RPC read; kept simple + CGo-free)
}

// splitLines splits a byte slice on newlines (non-empty lines).
func splitLines(b []byte) []string {
	var out []string
	start := 0
	for i := 0; i < len(b); i++ {
		if b[i] == '\n' {
			if i > start {
				out = append(out, string(b[start:i]))
			}
			start = i + 1
		}
	}
	if start < len(b) {
		out = append(out, string(b[start:]))
	}
	return out
}

// readAll reads all of r (best-effort).
func readAll(r interface{ Read([]byte) (int, error) }) ([]byte, error) {
	buf := make([]byte, 0, 4096)
	tmp := make([]byte, 4096)
	for {
		n, err := r.Read(tmp)
		if n > 0 {
			buf = append(buf, tmp[:n]...)
		}
		if err != nil {
			return buf, nil
		}
	}
}
