---
name: mnemonic-code-index
description: "Mnemonic code index and retrieval — Orientation Ladder, Index Freshness, Search Intent Router, extractors, dialect v2, code_impact, export, and serve Memory Visualization. Use when exploring code, searching, or checking index health."
license: MIT
metadata:
  author: devopstales
  version: "2.0"
  part-of: skillgrid
  notes: "Documents the Mnemonic code_* MCP tools and skillgrid search/export/serve surfaces (010)."
---

# Mnemonic Code Index — Protocol

Local-first code intelligence in the same SQLite store as `mem_*` / `web_*`
(`~/.skillgrid/mnemonic/<project>.sqlite`). Shared contract:
[`.agents/skills/_shared/conventions/mnemonic-code-indexing.md`](../_shared/conventions/mnemonic-code-indexing.md).

## Orientation Ladder (primary nav)

Prefer this over dumping whole directory trees into chat:

```
code_status → code_map → symbols / code_outline → code_related → code_read
```

1. **`code_status`** — Index Freshness: `fresh` | `lag` | `empty` | `unknown`. Compat `stale=true` means **empty only**, not lag.
2. If not `fresh` → **`code_index`**, then continue.
3. **`code_map`** — indexed structural overview.
4. **symbols / `code_outline`** — narrow units (`code_search_symbols` / outline).
5. **`code_related`** — Edge neighbors.
6. **`code_read`** — exact narrowed slice only.

Narrative Codebase Map (`docs/skillgrid/codebase/`) is **secondary** orientation; do not replace this ladder with a tree dump.

## Search Intent Router

`skillgrid search <query>` uses the **Search Intent Router** (daily path). This is **not** “six equal corpora” as the only model.

| Daily path | When |
|---|---|
| `symbols` | Identifier-shaped query |
| `grep` | Structural / Matcher dialect (`func $NAME…`) |
| `mem` | Decision / remember / past-work language |
| `code` | Unclassified (`route=default`) |

Advanced escape hatch: `--corpus code|grep|mem|symbols|hybrid|semantic`. Keep provenance when multi-corpus; never silently fuse mem+code.

### Explicit semantic degrade

Semantic search sets `degraded=true` (+ reason + fallback) when the Local Code Embedder is unavailable. **Never silent** FTS-as-semantic.

## Extractors / dialect / impact

- Python tree-sitter Extractor: build with **`-tags treesitter`** (CGo only inside Extractor Adapters).
- Matcher **dialect v2** adds: `struct $NAME`, `method $NAME($ARGS)`, `const $NAME`, `enum $NAME` (plus v1 patterns). Invalid pattern → abort.
- **`code_impact`** — additive blast-radius Impact Analysis; does not replace `code_get_*`.

## Export and Memory Visualization

```bash
skillgrid export --project ID --out DIR [--allowed-root DIR]
skillgrid serve   # dashboard: Memory Visualization + Code Graph (read-only)
```

Export writes Obsidian Markdown + viz JSON under an allowed root. Serve dashboard viz is read-only — no mutate/delete via the viz client.

**Code Graph** (CodeGraph-shaped): dashboard tab with 3-pane callers | source | callees, blast-radius summary, symbol search, entry-points screen, and flow/path between two symbols — over Mnemonic Edges (`/graph/*`). Mutating graph viz APIs return 405.

## Foundation tools (stable)

Do not rename/remove: `code_status`, `code_index`, `code_search`, `code_read`, `code_get_*`. Lexical baseline remains `chunks` / `chunks_fts`.

## CLI cheat sheet

```bash
skillgrid index --dir .
skillgrid search "AuthService"
skillgrid search --corpus grep 'func $NAME($ARGS)'
skillgrid search --corpus semantic "how does indexing work"   # watch degraded
skillgrid export --project myproj --out ~/vault/mnemonic-export
skillgrid serve
```
