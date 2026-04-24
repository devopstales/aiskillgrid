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

## graphify

This project uses **graphify** as the hub-supported codebase index (outputs under `graphify-out/`). Pair it with **`rg`/IDE search** and LSP for exact symbols.

Rules:
- Before answering architecture or codebase questions, read graphify-out/GRAPH_REPORT.md for god nodes and community structure
- If graphify-out/wiki/index.md exists, navigate it instead of reading raw files
- After modifying code files in this session, run `graphify update .` to keep the graph current (AST-only, no API cost)
- **Commit policy:** Either commit `graphify-out/` so all agents share the same snapshot, or keep it local/gitignored and regenerate after clone—pick one per repo; document in the root README or CONTRIBUTING if non-obvious

## Skillgrid: per-change handoff (filesystem context)

When the project uses **Skillgrid** + **OpenSpec**, keep **change-scoped** state on disk so the main session and **`Task` / subagents** stay aligned without pasting long tool output in chat. Full layout, handoff **template**, and subagent **contract** live in **`docs/wokflow.md`** (*Filesystem handoff* and *Parallel discovery*).

**Paths (`<change-id>` = directory name under `openspec/changes/<change-id>/`):**

- **Handoff:** `.skillgrid/tasks/context_<change-id>.md` — rolling goal, state, index of subagent work
- **Spill (long research / scrapes):** `.skillgrid/tasks/research/<change-id>/` — one markdown file per topic or run

**Rules:**

- **Parent session** does implementation (`/skillgrid-apply`) with full repo context. OpenSpec **`contextFiles`** and the handoff file both apply: read the handoff when it exists.
- **Subagents** (e.g. researcher, design critic, explore architect): **read** the handoff first, **write** long output under `tasks/research/…`, **update** the handoff, return a **short** reply with file paths.
- **Delegation:** If you use a subagent for a change, **put the handoff path in the `Task` prompt**; when it returns, **read the handoff and cited `research/…` files** before changing product code.
- **Durable memory** (decisions that must survive compactions) still goes through **Engram** / `mem_save` where configured; the handoff is **git-visible** state for the team and does not replace graphify or `rg` for “where is X in code.”

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

If **Engram** is enabled in your MCP config (see **`.configs/mcp/`** and **docs/tools.md**), use it for facts that should survive **compaction and new sessions**. It complements **on-disk** handoff and PRD/OpenSpec, not a replacement for reading the repo.

**When to `mem_save` (promptly, not only at end of task):** decisions, non-obvious **bugfix** learnings, discovered patterns, **config** or product preferences, and pointers to where the detail lives (stable **`topic_key`**, e.g. `skillgrid/{change}/…`).

**When to recall:** `mem_context` or **`mem_search`** at session start or before repeating similar work; if a hit is truncated, use **`mem_get_observation(id)`** before relying on it. On first message, if the user names a project area or problem, a quick **`mem_search`** on their keywords is appropriate when Engram is available.

**Session close:** Before you treat work as **done**, call **`mem_session_summary`** (goal, what changed, discoveries, next steps, relevant file paths) when the product supports it.

**Boundaries — handoff vs Engram:** The **per-change** handoff (`.skillgrid/tasks/context_…`) is for **this change** and **git** visibility. **Engram** is for **cross-session** and **durable** recall; avoid duplicating the same novel in both—**link** or one-line point from one to the other.

Full protocol: **`.agents/skills/memory-protocol/SKILL.md`**.