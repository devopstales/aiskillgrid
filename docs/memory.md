# Memory, indexing, and context (hub model)

This document explains **what we use for “memory” in the broad sense**: durable facts across sessions, orientation inside the **repository**, and optional ways to shrink noisy tool output. Install paths and MCP fragments are summarized in [`tools.md`](tools.md); agent behavior is also reflected in [`.configs/AGENTS.md`](../.configs/AGENTS.md).

---

## Mental model: three layers

| Layer | Question it answers | Planned tools |
|--------|---------------------|----------------|
| **Cross-session memory** | What did we decide, discover, or agree to last time? | **Engram** (primary) |
| **Repo index / map** | How is this codebase structured; where are the hot spots? | **graphify** (`graphify-out/`) |
| **Structural code search** | Where is this symbol, string, or pattern implemented? | **ripgrep** (`rg`), IDE search, **LSP** |
| **Per-change handoff (Skillgrid)** | For one OpenSpec **change id**, what is the rolling goal, state, and where are subagent reports? | **`.skillgrid/tasks/context_<change-id>.md`** and **`.skillgrid/tasks/research/<change-id>/`**, documented in `docs/wokflow.md` (*Filesystem handoff*) |

Those layers are **complementary**. Engram does not replace reading source; graphify does not replace exact search for a function name.

**Handoff vs Engram vs graphify:** The **context file** and **research/** spill are **on-disk, git-friendly state** for a single change—shared by the parent session and `Task` subagents without pasting long tool output into chat. **Engram** is for **cross-session** or **durable** observations (`mem_save`, `mem_search`) and compaction survival; the handoff can *point* to Engram topics but should not duplicate every memo. **graphify** remains the **structural** map of the repo (`graphify-out/`); the handoff **references** key paths or follow-up work, it does not replace a graph or `rg` for “where is X implemented.”

---

## Tools by purpose

### Engram — persistent agent memory

- **Purpose:** Store and retrieve **observations** that should survive chat compaction and new sessions: decisions, bugfix learnings, SDD artifact pointers, preferences, session summaries.
- **Interface:** MCP tools such as `mem_save`, `mem_search`, `mem_context`, `mem_get_observation`, `mem_session_summary` (exact set depends on [Engram](https://github.com/Gentleman-Programming/engram) server profile).
- **Hub wiring:** [`.configs/mcp/command/engram.json`](../.configs/mcp/command/engram.json); same merged MCP entry on **every** agent surface (Cursor, OpenCode, Copilot, etc.) so all agents share one store (default: `~/.engram/`).
- **Discipline:** [`.agents/skills/memory-protocol/SKILL.md`](../.agents/skills/memory-protocol/SKILL.md) (when to save, search, and close sessions). Frontmatter `name` is `engram-memory-protocol`; a duplicate tree exists at `engram-memory-protocol/`—**use `memory-protocol/`** in this hub (includes **Engram in this repository** wiring).
- **Teams:** Optional [git sync](https://github.com/Gentleman-Programming/engram) (`engram sync` / import) to share memory chunks via `.engram/` in the repo.

### Engram skills in `.agents/skills/`

Besides **memory protocol**, this hub vendors the **[Engram `skills/`](https://github.com/Gentleman-Programming/engram/tree/main/skills)** set under **`engram-*/`** (plus the duplicate `engram-memory-protocol/` folder). They cover maintainer workflow (triage, PRs, issues, commits), architecture and API guardrails, dashboard/TUI product rules, and more—many examples assume the **Engram** Go codebase; each skill’s **Engram in this repository** section (where present) explains how to map that to **aiskillgrid**.

**Full table:** [`skills.md`](skills.md) → **Engram skills (vendored)**.

**Often useful on any repo (with adaptation):** `engram-docs-alignment`, `engram-commit-hygiene`, `engram-cultural-norms`, `engram-backlog-triage`, `engram-architecture-guardrails` (principles), `engram-project-structure` (layering ideas), `engram-branch-pr`, `engram-issue-creation`, `engram-pr-review-deep`.

**Mostly Engram-product-specific:** `engram-business-rules`, `engram-server-api`, `engram-dashboard-htmx`, `engram-ui-elements`, `engram-visual-language`, `engram-tui-quality`, `engram-plugin-thin`, `gentleman-bubbletea` (under `engram-gentleman-bubbletea/`), `engram-testing-coverage`.

**Overlaps hub workflow:** `engram-sdd-flow` is a short phase list for Engram-product context; prefer hub **`skillgrid-*`** commands (see [`commands.md`](commands.md)) and optional **`openspec-*`** skills when this repo is the source of truth.

### graphify — repository knowledge graph

- **Purpose:** Generated **artifacts** under `graphify-out/` (e.g. `GRAPH_REPORT.md`, optional wiki index) so agents can orient on **architecture**, communities, and important nodes without embedding the whole tree in context.
- **When to refresh:** After substantive code or layout changes; convention in this hub is `graphify update .` (see [`.configs/AGENTS.md`](../.configs/AGENTS.md) and [`.cursor/rules/graphify.md`](../.cursor/rules/graphify.md)).
- **Commit policy:** Either commit `graphify-out/` for a shared snapshot or keep it local/regenerated after clone—choose per project and document it for contributors.

### Optional: context-mode

- **Purpose:** **Context window** optimization (e.g. sandboxing or shrinking bulky tool output)—helps token pressure, **not** repository indexing or long-term memory.
- **Status:** Listed for evaluation in [`tools.md`](tools.md); adopt if sessions routinely overflow context with tool dumps.

---

## Suggested order of use (agents)

1. **Rules + memory:** If Engram is connected, use `mem_context` / targeted `mem_search` before contradicting past decisions.
2. **Map:** If `graphify-out/` exists, skim `GRAPH_REPORT.md` (and wiki index if present) per `AGENTS.md`.
3. **Code:** Use `rg` / IDE search and LSP, then read files or run tests.
4. **Persist:** After meaningful decisions or session work, `mem_save` / `mem_session_summary` (and SDD persistence rules when in Engram or OpenSpec modes).

A shorter checklist for skills and personas: [`.agents/skills/references/indexing-and-memory.md`](../.agents/skills/references/indexing-and-memory.md) (if present in your checkout).

---

## Related documentation

- [`tools.md`](tools.md) — install, `uv` CLIs, Engram MCP subsection, optional tools
- [`skills.md`](skills.md) — **Engram skills (vendored)** table, `memory-protocol`, indexing references
- [`.configs/AGENTS.md`](../.configs/AGENTS.md) — graphify rules and project copy behavior
