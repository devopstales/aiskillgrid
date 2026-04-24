# Indexing and memory (hub tools)

Use these when the project or IDE has them enabled (merged `mcp.json`, `install.sh`, or local CLI). Do not assume a tool exists—check MCP/tool availability first.

## Engram (persistent memory MCP)

- **Typical tools:** `mem_context`, `mem_search`, `mem_get_observation`, `mem_save`, `mem_session_summary` (exact set depends on server config).
- **Role:** Cross-session facts, **Skillgrid** / OpenSpec artifact keys when the project uses **Engram** (e.g. `skillgrid-init/{project}`, `skillgrid/{change}/verify-report`), and decisions worth recalling.
- **Rule:** `mem_search` previews are short; use `mem_get_observation(id)` for full content before relying on an artifact.
- **When:** Session start, before claiming “we already decided…”, after locking a decision worth replaying.
- **Discipline:** `.agents/skills/memory-protocol/SKILL.md` (Engram memory protocol).

## Graphify

- **Role:** Repo knowledge graph (`graphify-out/`, e.g. `GRAPH_REPORT.md`) for communities / hot spots.
- **When:** Broad architecture or “how is this codebase structured?”—especially if `AGENTS.md` or project rules mention graphify.
- **Refresh:** After substantive code or layout changes, run `graphify update .` (see `.configs/AGENTS.md` and `.cursor/rules/graphify.md`).

## Code search (structural)

- **Role:** Find symbols, strings, and call sites with **ripgrep** (`rg`), IDE search, and LSP (go-to-definition, references)—not a separate index product in this hub.
- **When:** “Where is X implemented?”, impact analysis, or verifying patterns before edit. Prefer scoped search over reading entire directories.

## MCP in-memory server (`memory-npx`)

- **Config fragment:** `.configs/mcp/node/memory.json` (`@modelcontextprotocol/server-memory`)
- **Role:** Lightweight structured recall; use when enabled and appropriate to the task (not a substitute for Engram project memory).

## Per-change handoff (Skillgrid)

When the project uses **Skillgrid** and a named OpenSpec **change** is active, the repo may have **`.skillgrid/tasks/context_<change-id>.md`** (rolling handoff) and **`.skillgrid/tasks/research/<change-id>/`**. **Use these for subagent / parent sync** in the current effort; they are not a replacement for **Engram** (cross-session) or **graphify** (whole-repo map). See `docs/wokflow.md` — *Filesystem handoff*.

## Suggested order

1. Rules + `mem_context` / quick `mem_search` (if memory MCPs available)  
2. For an active **change** with a handoff file: read **`.skillgrid/tasks/context_<change-id>.md`** (when present)  
3. `graphify-out/` skim if present  
4. `rg` / IDE search (and LSP) for concrete code locations  
5. Read files / run tests as usual  
6. `mem_save` (or SDD persistence) for durable outputs the next session must see  

Narrative overview: [`docs/memory.md`](../../../docs/memory.md).
