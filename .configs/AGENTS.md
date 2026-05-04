# Application

# Behavioral guidelines

## 1. Think Before Coding

**Don't assume. Don't hide confusion. Surface tradeoffs.**

Before implementing:
- State your assumptions explicitly. If uncertain, ask.
- If multiple interpretations exist, present them - don't pick silently.
- If a simpler approach exists, say so. Push back when warranted.
- If something is unclear, stop. Name what's confusing. Ask.

## 2. Simplicity First

**Minimum code that solves the problem. Nothing speculative.**

- No features beyond what was asked.
- No abstractions for single-use code.
- No "flexibility" or "configurability" that wasn't requested.
- No error handling for impossible scenarios.
- If you write 200 lines and it could be 50, rewrite it.

Ask yourself: "Would a senior engineer say this is overcomplicated?" If yes, simplify.

## 3. Surgical Changes

**Touch only what you must. Clean up only your own mess.**

When editing existing code:
- Don't "improve" adjacent code, comments, or formatting.
- Don't refactor things that aren't broken.
- Match existing style, even if you'd do it differently.
- If you notice unrelated dead code, mention it - don't delete it.

When your changes create orphans:
- Remove imports/variables/functions that YOUR changes made unused.
- Don't remove pre-existing dead code unless asked.

The test: Every changed line should trace directly to the user's request.

## 4. Goal-Driven Execution

**Define success criteria. Loop until verified.**

Transform tasks into verifiable goals:
- "Add validation" → "Write tests for invalid inputs, then make them pass"
- "Fix the bug" → "Write a test that reproduces it, then make it pass"
- "Refactor X" → "Ensure tests pass before and after"

For multi-step tasks, state a brief plan:
```
1. [Step] → verify: [check]
2. [Step] → verify: [check]
3. [Step] → verify: [check]
```

Strong success criteria let you loop independently. Weak criteria ("make it work") require constant clarification.

<!-- skillgrid-gitnexus:start -->
## Skillgrid + GitNexus integration

For **all** Skillgrid workflows (brainstorm, design, plan, explore, apply, validate, finish, code review, security review, audits, codebase mapping, scans, fixers, and parallel subagents), use **GitNexus MCP** before grep-first exploration of application code. Use `repo: "aiskillgrid"` (or the repo id returned by `npx -y gitnexus@1.3.11 status`) for all queries. Subagents inherit this.

**CLI pin:** use **`npx -y gitnexus@1.3.11 …`** for analyze, status, serve, mcp, and other GitNexus commands unless a project file pins otherwise.

Rule file: **`.cursor/rules/skillgrid-gitnexus.md`** (mirrors under **`.agents/rules/`**, **`.opencode/rules/`**, **`.kilo/rules/`**). Workflow: **`.agents/workflows/gitnexus.md`**. Skills: **`.agents/skills/gitnexus-*/`**.
<!-- skillgrid-gitnexus:end -->

<!-- gitnexus:start -->
## GitNexus — code intelligence

This project is indexed locally with **GitNexus** (under `.gitnexus/` after `npx -y gitnexus@1.3.11 analyze`). Use the **aiskillgrid** repo id in examples below—if yours differs, substitute the id from **`npx -y gitnexus@1.3.11 status`** or **`gitnexus://repo/<id>/context`**.

> If any GitNexus tool warns the index is stale, run `npx -y gitnexus@1.3.11 analyze` in the repo root first.

### Always do

- **MUST run impact analysis before editing any symbol.** Before modifying a function, class, or method, run `gitnexus_impact({target: "symbolName", direction: "upstream"})` and report the blast radius (direct callers, affected processes, risk level) to the user.
- **MUST run `gitnexus_detect_changes()` before committing** to verify your changes only affect expected symbols and execution flows.
- **MUST warn the user** if impact analysis returns HIGH or CRITICAL risk before proceeding with edits.
- When exploring unfamiliar code, use `gitnexus_query({query: "concept"})` to find execution flows instead of grepping alone. It returns process-grouped results ranked by relevance.
- When you need full context on a specific symbol — callers, callees, which execution flows it participates in — use `gitnexus_context({name: "symbolName"})`.

### When debugging

1. `gitnexus_query({query: "<error or symptom>"})` — find execution flows related to the issue
2. `gitnexus_context({name: "<suspect function>"})` — see all callers, callees, and process participation
3. **READ** `gitnexus://repo/aiskillgrid/process/{processName}` — trace the full execution flow step by step
4. For regressions: `gitnexus_detect_changes({scope: "compare", base_ref: "main"})` — see what your branch changed

### When refactoring

- **Renaming**: MUST use `gitnexus_rename({symbol_name: "old", new_name: "new", dry_run: true})` first. Review the preview — graph edits are safe, text_search edits need manual review. Then run with `dry_run: false`.
- **Extracting/splitting**: MUST run `gitnexus_context({name: "target"})` to see all incoming/outgoing refs, then `gitnexus_impact({target: "target", direction: "upstream"})` to find all external callers before moving code.
- After any refactor: run `gitnexus_detect_changes({scope: "all"})` to verify only expected files changed.

### Never do

- NEVER edit a function, class, or method without first running `gitnexus_impact` on it.
- NEVER ignore HIGH or CRITICAL risk warnings from impact analysis.
- NEVER rename symbols with find-and-replace — use `gitnexus_rename`, which understands the call graph.
- NEVER commit changes without running `gitnexus_detect_changes()` to check affected scope.

### Tools quick reference

| Tool | When to use | Command |
|------|-------------|---------|
| `query` | Find code by concept | `gitnexus_query({query: "auth validation"})` |
| `context` | 360-degree view of one symbol | `gitnexus_context({name: "validateUser"})` |
| `impact` | Blast radius before editing | `gitnexus_impact({target: "X", direction: "upstream"})` |
| `detect_changes` | Pre-commit scope check | `gitnexus_detect_changes({scope: "staged"})` |
| `rename` | Safe multi-file rename | `gitnexus_rename({symbol_name: "old", new_name: "new", dry_run: true})` |
| `cypher` | Custom graph queries | `gitnexus_cypher({query: "MATCH ..."})` |

### Impact risk levels

| Depth | Meaning | Action |
|-------|---------|--------|
| d=1 | WILL BREAK — direct callers/importers | MUST update these |
| d=2 | LIKELY AFFECTED — indirect deps | Should test |
| d=3 | MAY NEED TESTING — transitive | Test if critical path |

### Resources

| Resource | Use for |
|----------|---------|
| `gitnexus://repo/aiskillgrid/context` | Codebase overview, check index freshness |
| `gitnexus://repo/aiskillgrid/clusters` | All functional areas |
| `gitnexus://repo/aiskillgrid/processes` | All execution flows |
| `gitnexus://repo/aiskillgrid/process/{name}` | Step-by-step execution trace |

### Self-check before finishing

Before completing any code modification task, verify:

1. `gitnexus_impact` was run for all modified symbols
2. No HIGH/CRITICAL risk warnings were ignored
3. `gitnexus_detect_changes()` confirms changes match expected scope
4. All d=1 (WILL BREAK) dependents were updated

### Keeping the index fresh

After committing code changes, the GitNexus index becomes stale. Re-run analyze to update it:

```bash
npx -y gitnexus@1.3.11 analyze
```

If the index previously included embeddings, preserve them by adding `--embeddings`:

```bash
npx -y gitnexus@1.3.11 analyze --embeddings
```

To check whether embeddings exist, inspect **`.gitnexus/meta.json`** — the `stats.embeddings` field shows the count (0 means no embeddings). **Running analyze without `--embeddings` will delete any previously generated embeddings.**

> Claude Code users: A PostToolUse hook can handle re-indexing automatically after `git commit` and `git merge` when configured.

### CLI

- Re-index: `npx -y gitnexus@1.3.11 analyze`
- Check freshness: `npx -y gitnexus@1.3.11 status`
- Generate docs: `npx -y gitnexus@1.3.11 wiki`

<!-- gitnexus:end -->

## Skillgrid session bootstrap

For Skillgrid + OpenSpec work, read when present (before deep implementation):

1. Root **`DESIGN.md`**
2. **`.skillgrid/project/PROJECT.md`**, **`ARCHITECTURE.md`**, **`STRUCTURE.md`**
3. **`.skillgrid/prd/INDEX.md`** — dependency-ordered PRDs and the **Execution snapshot** at the top (current phase, active change/slice, discovered work)
4. Active **`.skillgrid/tasks/context_<change-id>.md`** and **`openspec/changes/<change-id>/`** (`proposal.md`, `design.md`, `tasks.md`)
5. The active vertical slice: **`openspec/changes/<change-id>/specs/<slice-slug>/spec.md`**
6. **`.skillgrid/project/SKILL_REGISTRY.md`** only when choosing skills or subagents (compact index; do not paste the whole file into chat)

After meaningful work: refresh the INDEX snapshot, `tasks.md`, slice specs, and handoff per phase (`skillgrid-prd-artifacts`, `skillgrid-spec-artifacts`).

**File blanks:** copy from **`.skillgrid/templates/`** (`template-*.md`; see **`docs/03-skillgrid-logic.md`**). Architectural decisions: **`.skillgrid/templates/template-adr.md`** → new `NNNN-*.md` under **`.skillgrid/adr/`** (that folder holds ADRs only).

## Skillgrid: per-change handoff (filesystem context)

When the project uses **Skillgrid** + **OpenSpec**, keep **change-scoped** state on disk so the main session and **`Task` / subagents** stay aligned without pasting long tool output in chat. Full layout, handoff **template**, and subagent **contract** live in **`docs/02-workflow-usage.md`** (*Filesystem handoff* and *Parallel discovery*).

**Paths (`<change-id>` = directory name under `openspec/changes/<change-id>/`):**

- **Handoff:** `.skillgrid/tasks/context_<change-id>.md` — rolling goal, state, index of subagent work
- **Spill (long research / scrapes):** `.skillgrid/tasks/research/<change-id>/` — one markdown file per topic or run

**Rules:**

- **Parent session** does implementation (`/skillgrid-apply`) with full repo context. OpenSpec **`contextFiles`** and the handoff file both apply: read the handoff when it exists.
- **Subagents** (e.g. researcher, design critic, explore architect): **read** the handoff first, **write** long output under `tasks/research/…`, **update** the handoff, return a **short** reply with file paths.
- **Delegation:** If you use a subagent for a change, **put the handoff path in the `Task` prompt**; when it returns, **read the handoff and cited `research/…` files** before changing product code.
- **Durable memory** (decisions that must survive compactions) still goes through **Engram** / `mem_save` where configured; the handoff is **git-visible** state for the team and does not replace **GitNexus**, `rg`, or IDE search for “where is X in code.”

## Web research and scrape

When the session uses **MCP** web tools (e.g. **Firecrawl**, **Brave** search, **Exa**, built-in `WebFetch` / `WebSearch`):

- **Cite** — Factual claims from the web need **source** (title, URL, access date or “as of” when it matters). Do not pass scraped text off as your own training.
- **Spill** — Long pages, full scrape output, or multi-page research belong in a **file** (e.g. **`.skillgrid/tasks/research/<change-id>/…md`** on Skillgrid projects, or a short team path). Keep the chat turn to a **summary + links/paths** so context stays small.
- **Tool choice** — Use **search** to discover URLs; use **scrape** for one or a few known URLs; use **map/crawl** only when the task needs broad URL discovery or a section of a site. Prefer an **interact** / browser flow over blind scraping when the content is behind JS, login, or multi-step UI (and stop if a human is required).
- **Secrets and policy** — Never commit or paste **API keys, cookies, session headers, or paywalled** content from scrapes into the repo. If a site blocks automation or TOS is unclear, say so and switch to public docs, official API docs, or **Context7** when available.
- **Failures** — If a tool errors or is unavailable, state it and fall back to another allowed tool or a smaller scope. Do not burn retries on the same URL without a new hypothesis.
- **Rate and scope** — Stay proportional to the question; no bulk crawling “just in case.”

Deeper patterns live under **`.agents/skills/`** (e.g. `firecrawl-scrape`, `exa-search`, `deep-research`).

## Engram (persistent memory MCP)

If **Engram** is enabled in your MCP config (see **`.configs/mcp/`** and **docs/01-installation.md**), use it for facts that should survive **compaction and new sessions**. It complements **on-disk** handoff and PRD/OpenSpec, not a replacement for reading the repo.

**When to `mem_save` (promptly, not only at end of task):** decisions, non-obvious **bugfix** learnings, discovered patterns, **config** or product preferences, and pointers to where the detail lives (stable **`topic_key`**, e.g. `skillgrid/{change}/…`).

**When to recall:** `mem_context` or **`mem_search`** at session start or before repeating similar work; if a hit is truncated, use **`mem_get_observation(id)`** before relying on it. On first message, if the user names a project area or problem, a quick **`mem_search`** on their keywords is appropriate when Engram is available.

**Session close:** Before you treat work as **done**, call **`mem_session_summary`** (goal, what changed, discoveries, next steps, relevant file paths) when the product supports it.

**Boundaries — handoff vs Engram:** The **per-change** handoff (`.skillgrid/tasks/context_…`) is for **this change** and **git** visibility. **Engram** is for **cross-session** and **durable** recall; avoid duplicating the same novel in both—**link** or one-line point from one to the other.

Full protocol: **`.agents/skills/memory-protocol/SKILL.md`**.