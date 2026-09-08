---
name: mnemonic
description: "skillgrid-cli memory + code index + web cache — sessions, mem_save/search, code ladder, centralized search router"
trigger: /mnemonic
---

# /mnemonic

Operate Mnemonic, skillgrid's local-first persistent memory (SQLite + FTS5, single `skillgrid` binary).
Three things it does that a bare agent cannot:
1. **Persistent memory** — decisions, bugfixes, and discoveries in `~/.skillgrid/mnemonic/<project>.sqlite` survive across sessions and compactions. Recall weeks later without re-reading everything.
2. **Code index** — per-project chunk/symbol/edge graph with freshness tracking. Orient without dumping whole trees.
3. **Web cache** — Context7/Exa/DeepWiki/WebFetch snapshots with TTLs, so online research is reused instead of re-fetched.

## Usage

```
/mnemonic                                             # session status: mem_context + doctor + index status
/mnemonic <path>                                      # orient on a path: index status → map → memory search
/mnemonic save "<title>" --type decision --content "..."   # curated save (needs active session)
/mnemonic recall "<question>"                         # mem_context → mem_search → mem_get_observation
/mnemonic index [--dir .] [--watch]                   # incremental code index (opt-in live watcher)
/mnemonic search "<query>" [--corpus code|grep|mem|symbols|hybrid|semantic]
/mnemonic code <subcommand>                           # orient/graph: map, get-symbol, impact, callers/callees...
/mnemonic serve | mcp | setup | export | trail        # daemon, agent wiring, export, trails
```

## What You Must Do When Invoked

If no path was given, use `.` (current directory). Do not ask the user for a path.

Follow these steps in order. Do not skip steps.

### Step 1 — Ensure the binary and store exist

```bash
command -v skillgrid >/dev/null || go build -o /tmp/skillgrid-op ./skillgrid-cli 2>&1 | tail -3
skillgrid --help 2>&1 | head -20
ls "${SKILLGRID_MNEMONIC_DATA_DIR:-$HOME/.skillgrid/mnemonic}" 2>/dev/null || echo "no mnemonic store yet (created on first index/save)"
```

Facts (verified against `skillgrid-cli/`):
- Data dir: `~/.skillgrid/mnemonic`, override `SKILLGRID_MNEMONIC_DATA_DIR`. Per-project store `<project>.sqlite`.
- Daemons: `skillgrid serve` → HTTP API `http://127.0.0.1:7438` (viewer `/`, Swagger `/swagger-ui`); `skillgrid mcp` → stdio server `skillgrid-mnemonic` with `mem_*`, `code_*`, `web_*`, `team_*` tools.
- Config: `skillgrid-cli/config.d/indexing.yaml` (include go/ts/tsx/md; chunk 80 lines / overlap 10; web-cache cap 256 KB; TTLs context7 720h, exa 168h, deepwiki 336h, fetch 168h).
- If the binary is missing, say so and stop — do not invent memory contents.

### Step 2 — Open or recover the session

Session tools require a live session id (`sid`). `mem_save` without `session_id` fails.

```
mem_session_start(title: "mnemonic/<goal>")   # once per agent session, reuse sid for every mem_save
```

Recovery (never ask what you can search):
```
mem_context(limit: 5)                          # fast: recent session summaries
→ if empty: mem_search(query: "<phase keywords>") → mem_get_observation(id)  # previews are truncated; get_observation is the only full-content path
```

SDD naming (when working a change `<NNN-slug>`): `title` == `topic_key` == `sdd/<NNN-slug>/<artifact>` (`research|intent|change|tasks|spec|apply-progress|verification|archive-report|state`), `type: architecture`, `scope: project`. Same `topic_key` + `scope` → UPDATE, not INSERT.

### Step 3 — Recall before doing (mandatory before new work)

```
1. mem_context — checks recent history (cheap)
2. mem_search(query, match_mode: any|all, scope: project|user|global, all_projects: bool) — FTS5 over observations
3. mem_get_observation(id) — full untruncated body
```

Also search proactively when starting work that may overlap prior sessions, or when the user says any variant of "remember / recall / what did we do / how did we solve". For SDD changes, group searches first (STEP A), then retrievals (STEP B).

### Step 4 — Save during work (mandatory — do NOT wait to be asked)

Call `mem_save` IMMEDIATELY after: bug fix, architecture/decision, non-obvious discovery, config change, new pattern, user preference/correction/confirmation.

```
mem_save(
  title: "Verb + what",                        # short, searchable, e.g. "Fixed N+1 query in UserList"
  type: decision|architecture|bugfix|pattern|config|discovery|learning|preference|convention,
  content: "**What**: one sentence\n**Why**: motive\n**Where**: files/paths\n**Learned**: gotchas (omit if none)",
  session_id: sid,                              # REQUIRED
  scope: project,                               # default; user|global when appropriate
  topic_key: "stable/key",                      # recommended for evolving topics (upsert)
)
```

Rules:
- Different topics must not share a `topic_key`. If unsure, call `mem_suggest_topic_key` first, then reuse it.
- Read current content with `mem_get_observation` before rewriting an evolved topic.
- Wrap secrets/tokens/PII in `<private>…</private>` — stripped before storage.
- Safety net: end task responses with a `## Key Learnings:` numbered list; passive capture persists it even if `mem_save` was missed.

### Step 5 — Code retrieval: Orientation Ladder + Search Router

Never dump whole directory trees. Prefer:

```
code_status → code_map → symbols / code_outline → code_related → code_read
```

```bash
skillgrid index --dir .                                   # reindex when status is not fresh
skillgrid search "AuthService"                            # intent router: daily path auto-picked, stamped route+provenance
skillgrid search --corpus symbols "AuthService"           # escape hatch: code|grep|mem|symbols|hybrid|semantic
skillgrid code map [--path-prefix DIR]
skillgrid code get-symbol --name Name
skillgrid code impact --name Name                          # blast-radius (additive; does not replace code_get_*)
skillgrid code get-callers|get-callees|get-dependents|get-implementors|get-tests-for|get-type-hierarchy --name Name
skillgrid code index-status
```

- **Freshness**: `fresh` | `lag` | `empty` | `unknown`. Compat `stale=true` means empty-only, not lag. Reindex when not `fresh`.
- **Router**: omit `--corpus` for the daily path (identifier→symbols, `func $NAME`→grep, decision/remember→mem, else code). `--corpus` is the advanced escape hatch, not six equal peers.
- **Provenance**: mem and code are never silently fused — mem results carry `observations`, code-family results never mix mem observations unlabeled.
- **Semantic** sets `degraded=true` explicitly when the Local Code Embedder is unavailable (never silent).

### Step 6 — Web research cache (before/after every remote lookup)

```
web_cache_lookup(source: context7|exa|deepwiki|fetch|manual, ...)  # BEFORE the remote call
→ miss: call Context7 / Exa / DeepWiki / WebFetch
→ web_cache_save(source, content, ...) within the same turn        # cap 256 KB — summarize first
"what did we find about X?" → web_cache_search(query) → web_cache_get(id)
```

### Step 7 — Close the session (mandatory before saying "done")

```
1. mem_session_summary(session_id: sid, summary: "## Goal… ## Discoveries… ## Accomplished… ## Next Steps… ## Relevant Files…")
2. mem_session_end(session_id: sid, summary: "one-line outcome")
```

After compaction / "FIRST ACTION REQUIRED": FIRST call `mem_session_summary` with the compacted content, then `mem_context`, only then continue.

## Subcommands (CLI reference)

| Command | Purpose |
|---|---|
| `skillgrid mcp [--debug]` | stdio MCP server (`mem_*`, `code_*`, `web_*`, `team_*`) |
| `skillgrid serve [--port 7438] [--bind 127.0.0.1] [--dir DIR]` | HTTP API + dashboard (Memory Visualization and Code Graph are read-only; mutate via viz is rejected) |
| `skillgrid index [--dir .] [--watch] [--embeddings]` | incremental code index; `--watch` opts into debounced live indexer |
| `skillgrid search [--corpus …] [--limit 20] [--dir .] <query>` | centralized retrieval with intent routing; JSON with `route` + `provenance` |
| `skillgrid code <get-symbol\|get-signature\|symbols-in-file\|map\|list-projects\|index-status\|impact\|get-callers\|…>` | orient/graph (search lives under `skillgrid search`) |
| `skillgrid setup <opencode\|kilocode\|cursor> [--dry-run]` | install `skillgrid-mnemonic` MCP + agent plugins |
| `skillgrid migrate [--tier]` | backfill tier sidecars |
| `skillgrid trail <recent\|show>` | inspect retrieval trails |
| `skillgrid export --project ID --out DIR` | Obsidian Markdown + viz JSON (allowed-root enforced) |

Full protocol refs: `.agents/skills/mnemonic-memory/SKILL.md`, `.agents/skills/mnemonic-code-index/SKILL.md`, `.agents/skills/_shared/conventions/mnemonic-memory.md`.
