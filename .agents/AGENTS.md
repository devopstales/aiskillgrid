<!-- skillgrid:global -->

# AGENTS.md

Global rules for all AI agents (Kilo, OpenCode, Cursor, Claude, Codex, Gemini). Applies to every project; project-specific rules in local `AGENTS.md`/`CLAUDE.md` take precedence when they conflict.

**Tradeoff:** These rules bias toward caution over speed. For trivial tasks, use judgment.

## If You Are an AI Agent

Be a tool that protects your human partner, not one that embarrasses them. Quality over volume:

- One problem per change. Never bundle unrelated edits.
- Solve a real problem your human experienced — not a theoretical one. If "it could cause issues" is your only justification, stop.
- Before contributing to an external repo: search existing (open and closed) PRs, read the PR template, identify yourself, show your human the full diff before submitting.
- Never fabricate: no invented claims, no "my review agent flagged this", no filler sentences.

## 1. Think Before Coding

Don't assume. Don't hide confusion. Surface tradeoffs.

- State assumptions explicitly; if uncertain, ask.
- If multiple interpretations exist, present them — don't pick silently.
- If a simpler approach exists, say so. Push back when warranted.
- If something is unclear, stop. Name what's confusing. Ask.
- **Before writing any code, check for applicable skills and invoke them first.** No exceptions for "simple" tasks. If a skill exists for the work, use it.

## 2. Simplicity First

Minimum code that solves the problem. Nothing speculative.

- No features beyond what was asked.
- No abstractions for single-use code.
- No "flexibility" or "configurability" that wasn't requested.
- No error handling for impossible scenarios.
- If you write 200 lines and it could be 50, rewrite it.

Ask yourself: "Would a senior engineer say this is overcomplicated?" If yes, simplify.

## 3. Surgical Changes

Touch only what you must. Clean up only your own mess.

- Don't "improve" adjacent code, comments, or formatting.
- Don't refactor things that aren't broken.
- Match existing style, even if you'd do it differently.
- Unrelated dead code: mention it, don't delete it.
- Remove imports/variables/functions that YOUR changes made unused; don't remove pre-existing dead code unless asked.

The test: every changed line should trace directly to the human's request.

## 4. Goal-Driven Execution

Define success criteria. Loop until verified.

- "Add validation" → "Write tests for invalid inputs, then make them pass"
- "Fix the bug" → "Write a test that reproduces it, then make it pass"
- Multi-step work: state a brief plan with a verification step per item:
  ```
  1. [step] -> verify: [check]
  2. [step] -> verify: [check]
  ```

## Codebase Rules

- Concise, direct responses. No preamble, no filler.
- Prefer editing existing files over creating new ones.
- Follow existing code style and conventions before introducing new patterns.
- Never commit secrets or keys.

## Verification

Verify work before claiming it is done — evidence before assertions:

- Run the relevant tests/build/lint and confirm the output passes.
- For CLI tools, run the command the user would run and check exit code + output.
- Prefer no-op confirmation: a change is done only when its success criterion (tests, command output, validation) has been observed.

## Environment

- Mnemonic persistent memory is available (skillgrid local-first store: `mem_*`, `code_*`, `web_*` over the `mnemonic` MCP): update it on decisions, bug fixes, and non-obvious discoveries.

<!-- /skillgrid:global -->

<!-- skillgrid:mnemonic-memory -->

## Mnemonic Persistent Memory

The full protocol — save triggers, search ladder, topic-key upserts, session close, compaction recovery — lives in the `mnemonic` skill (`.agents/skills/mnemonic/SKILL.md`) and is injected automatically at session start. Do not duplicate it here.

<!-- /skillgrid:mnemonic-memory -->

<!-- skillgrid:mnemonic-code-index -->

## Mnemonic Code Index

Per-project code search over the same SQLite store as memory (`~/.skillgrid/mnemonic/<project>.sqlite`). The full protocol lives in the `mnemonic-code-index` skill (`.agents/skills/mnemonic-code-index/SKILL.md`) and `.agents/skills/_shared/conventions/mnemonic-code-indexing.md` — do not duplicate it here.

Orientation ladder (canonical order):

```
1. code_status   → health + freshness (fresh|lag|empty|unknown); if not fresh → code_index
2. code_map      → structural overview, not a prose dump
3. symbols / outline → narrow the unit
4. code_related  → neighbors via edges
5. code_read     → exact slice for an already-narrowed path + range
```

Rules: never dump whole trees to orient; never read whole files speculatively; `rg`/`grep` for exact identifiers, index for exploration. CLI equivalent: `skillgrid index|search|export|serve`.

<!-- /skillgrid:mnemonic-code-index -->

<!-- skillgrid:mnemonic-search-cache -->

## Mnemonic Web Research Cache

Before calling a remote research MCP (Context7, Exa, DeepWiki, WebFetch):

1. `web_cache_lookup(source, ...)` — check the cache first.
2. On miss: call the remote MCP.
3. `web_cache_save(source, content, ...)` — persist immediately after (cap 256 KB per snapshot — summarize first).

"What did we find about X online?" → `web_cache_search(query)`, then `web_cache_get(id)`. Full protocol + TTLs in the `mnemonic` skill.

<!-- /skillgrid:mnemonic-search-cache -->

<!-- skillgrid:repo -->

## Repo structure

```
AGENTS.md                 # cross-agent SoT: SDD workflow, tracker, conventions
.agents/skills/           # skill sources (sdd-*, execution, mnemonic, …)
docs/skillgrid/           # SDD artifacts: config.yaml, changes/, archive/, glossary/
skillgrid-cli/            # Go binary (Mnemonic store, code index, web cache)
plugins/                  # opencode / kilo / cursor wiring
config.d/                 # indexing + cache config
scripts/                  # operator scripts
```

Project facts (stack, testing, tracker, rules) live in `docs/skillgrid/config.yaml` — never duplicate them here.

## Ticketing

Backlog.md. Tickets live under `.backlog/tasks/` and are managed **only** via the `backlog` CLI — never edit task markdown by hand. The full workflow (overview → task-creation / task-execution / task-finalization) is in root `AGENTS.md` under BACKLOG.MD GUIDELINES. SDD specs may file tickets via `issue-creation`.

## Scratch workspaces (`.skillgrid/sdd/`)

`subagent-execution` keeps per-change scratch under `.skillgrid/sdd/<NNN-slug>/`: task briefs, implementer reports, review diffs, `progress.md` ledger. Git-ignored, plan-scoped, deleted at finish — never commit it, never hand-edit it. Resolve the path via `scripts/sdd-workspace <tasks.md>` (see the `subagent-execution` skill).

Not to be confused with the `handoff` skill: a handoff brief peels an out-of-scope side problem for another session and lives in OS temp (`/tmp`), never in `.skillgrid/`.

<!-- /skillgrid:repo -->
