<!-- skillgrid-gitnexus v0.1.0 -->
---
description: Use GitNexus graph intelligence before grep or file exploration for any codebase question or Skillgrid workflow
alwaysApply: true
---

# Skillgrid + GitNexus — exploration first

For **any** question about how the codebase works — how a feature is implemented, where something is stored, how data flows, what calls what — use **GitNexus MCP** before **Grep**, **Glob**, or **Shell** file walks. Open source files to **confirm** details after the graph points you there. Fall back to filesystem search only when GitNexus returns no useful results or is unavailable.

**Repo argument:** use `repo: "aiskillgrid"` on **`gitnexus_query`**, **`gitnexus_context`**, **`gitnexus_impact`**, and **`gitnexus_detect_changes`** (unless `npx -y gitnexus@1.3.11 status` / MCP shows a different registered id — then use that consistently).

**CLI:** refresh a stale index with **`npx -y gitnexus@1.3.11 analyze`** from the repo root (see `.agents/workflows/gitnexus.md`).

## Priority: GitNexus → `.skillgrid/project/` → grep/glob/shell

1. **GitNexus** — execution flows, clusters, symbol context, impact.
2. **`.skillgrid/project/`** (`STRUCTURE.md`, `ARCHITECTURE.md`, `PROJECT.md`) — team-written summaries; skim first for orientation, then **confirm** in the graph.
3. **Grep / glob / shell** — exact symbols, tests, or when GitNexus is empty or down.

## Skillgrid workflows (mandatory)

Applies to **every** Skillgrid phase and **Task** / subagent that inspects application code — including **`/skillgrid-brainstorm`**, **`/skillgrid-design`**, **`/skillgrid-plan`**, **`/skillgrid-explore`**, **`/skillgrid-breakdown`**, **`/skillgrid-apply`**, **`/skillgrid-validate`**, **`/skillgrid-finish`**, **`/skillgrid-session`**, checkpoints, parallel research repo lanes, codebase mapping, code review / security / test personas, and any flow that answers “what calls this?” or “what breaks if…?”.

1. **Before** relying on grep/glob/shell for import trees, callers, callees, execution flows, or blast radius: use **`gitnexus_query`**, **`gitnexus_context`**, and/or **`gitnexus_impact`** with `repo: "aiskillgrid"` (or your registered id).
2. **`gitnexus_detect_changes`**: use the same `repo` (git-aware scope checks).

## General codebase questions

1. `gitnexus_query({query: "<concept>", repo: "aiskillgrid"})`
2. `gitnexus_context({name: "<symbol>", repo: "aiskillgrid"})` for symbols the query surfaces
3. Read **source files** only to verify implementation details
4. Fall back to Grep/Glob/Shell only if GitNexus returns **0** relevant processes or tools error out

## Fallback triggers

Use grep/glob/shell (and **`skillgrid-codebase-map`** / **`rg`** as documented) when:

- **`gitnexus_query`** returns no relevant processes
- **`gitnexus_impact`** returns 0 upstream — **retry once** (wrong symbol name or repo id), then fall back if still empty
- GitNexus warns the index is **stale** — run **`npx -y gitnexus@1.3.11 analyze`**, then retry
- **MCP** server is unavailable — note it and proceed with search; restore GitNexus when possible

## See also

- **`.kilo/rules/gitnexus.md`** — short index/stale-graph reminder
- **`.configs/AGENTS.md`** — full GitNexus MUST/NEVER and tool table
- **`gitnexus-exploring`** skill — step-by-step investigation workflow
