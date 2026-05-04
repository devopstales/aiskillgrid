# Indexing and memory (hub tools)

Use these when the project or IDE has them enabled (merged `mcp.json`, `install.sh`, or local CLI). Do not assume a tool exists—check MCP/tool availability first.

## Engram (persistent memory MCP)

- **Typical tools:** `mem_context`, `mem_search`, `mem_get_observation`, `mem_save`, `mem_session_summary` (exact set depends on server config).
- **Role:** Cross-session facts, **Skillgrid** / OpenSpec artifact keys when the project uses **Engram** (e.g. `skillgrid-init/{project}`, `skillgrid/{change}/state`, `skillgrid/{change}/verify-report`), and decisions worth recalling.
- **Rule:** `mem_search` previews are short; `mem_get_observation(id)` is mandatory before relying on recovered artifact, blocker, task, status, or decision content.
- **When:** Session start, before claiming “we already decided…”, after locking a decision worth replaying.
- **Discipline:** `.agents/skills/memory-protocol/SKILL.md` (Engram memory protocol).

### Skillgrid Recovery Keys

Use stable keys for change recovery:

```text
skillgrid-init/<project>
skillgrid/<change>/state
skillgrid/<change>/prd
skillgrid/<change>/proposal
skillgrid/<change>/spec
skillgrid/<change>/tasks
skillgrid/<change>/context
skillgrid/<change>/verify-report
skillgrid/<change>/archive
```

`skillgrid/<change>/state` is the compact index: phase, status, paths, blockers, next action, and last updated. It points to canonical disk artifacts in hybrid/openspec mode and to concrete Engram artifact observations in engram mode.

Recovery is always two-step:

1. `mem_search` the stable topic key to get candidate IDs.
2. `mem_get_observation(id)` for full untruncated content before acting.

If disk artifacts and Engram disagree, treat disk artifacts as canonical in `hybrid` / `openspec` mode, then update the state snapshot once the correct state is known.

## GitNexus

- **Role:** Repo knowledge graph and agent code-intelligence index (`.gitnexus/`, GitNexus MCP resources/tools) for communities, execution flows, impact analysis, and symbol context.
- **When:** Broad architecture or “how is this codebase structured?”—especially if `AGENTS.md` or project rules mention GitNexus.
- **Refresh:** After substantive code or layout changes, run `npx -y gitnexus@1.3.11 analyze` from the repo root. Use `npx -y gitnexus@1.3.11 analyze --force` only when a full rebuild is needed.
- **Init:** **`npx -y gitnexus@1.3.11 analyze`** from the repo root when building the graph the first time (e.g. **`/skillgrid-init`**).
- **MCP:** Configure **`npx -y gitnexus@1.3.11 mcp`** so agents can call GitNexus `query`, `context`, `impact`, `detect_changes`, and related graph tools.

## CocoIndex Code (`ccc`)

- **Skill:** `.agents/skills/ccc/SKILL.md`
- **Role:** Semantic search over the repo (`ccc search`); MCP fragment **`cocoindex-code`** runs `ccc mcp` when merged from `.configs/mcp/command/cocoindex-code.json`.
- **Init:** From project root, **`ccc init`** if not initialized, then **`ccc index`** (also part of **`/skillgrid-init`** when the CLI is installed). Re-run **`ccc index`** after large refactors or many new files.
- **When:** “Where is X implemented?” in natural language, impact analysis, finding patterns before edit—pair with **`rg`** for exact symbols.

## Code search (structural)

- **Role:** Find symbols, strings, and call sites with **ripgrep** (`rg`), IDE search, and LSP—exact matches, not semantic index results.
- **When:** Unique identifiers, known file paths, or verifying a hypothesis from **`ccc search`**.

## MCP in-memory server (`memory-npx`)

- **Config fragment:** `.configs/mcp/node/memory.json` (`@modelcontextprotocol/server-memory`)
- **Role:** Lightweight structured recall; use when enabled and appropriate (not a substitute for Engram project memory).

## Per-change handoff (Skillgrid)

When the project uses **Skillgrid** and a named OpenSpec **change** is active, the repo may have **`.skillgrid/tasks/context_<change-id>.md`** (rolling handoff) and **`.skillgrid/tasks/research/<change-id>/`**. Use these for subagent / parent sync; they are not a replacement for **Engram**, **GitNexus**, or **ccc**. See `docs/02-workflow-usage.md` — *Filesystem handoff*.

## Suggested order

1. Rules + `mem_context` / quick `mem_search` (if memory MCPs available; retrieve full observations before relying on hits)  
2. For an active **change** with a handoff file: read **`.skillgrid/tasks/context_<change-id>.md`** (when present)  
3. GitNexus MCP/resources or `.gitnexus/` status when present  
4. **`ccc search`** when the CocoIndex index is fresh (after **`ccc index`**)  
5. `rg` / IDE search (and LSP) for concrete code locations  
6. Read files / run tests as usual  
7. `mem_save` (or SDD persistence) for durable outputs the next session must see  

Narrative overview: [`docs/memory.md`](../../docs/memory.md).
