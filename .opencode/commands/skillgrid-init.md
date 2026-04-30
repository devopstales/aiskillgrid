---
name: /skillgrid-init
id: skillgrid-init
category: Workflow
description: Bootstrap workflow: structure, CocoIndex (ccc), graphify, hybrid persistence (OpenSpec + Engram)
allowed-tools: Read, Write, Glob, Grep, Bash, Task
argument-hint: "[optional: app purpose, stack, or brownfield notes]"
---

<objective>

You are executing **`/skillgrid-init`** for the Skillgrid workflow.

Bootstrap or align a repository for Skillgrid: config, artifact store, OpenSpec, project docs, design docs, and optional indexing.

**Status on exit:** no change PRD status unless creating initial planning artifacts.

</objective>

<process>

## Shared Skillgrid Skills

Before acting, load only the skills needed for the phase:

- `.agents/skills/skillgrid-questioning/SKILL.md` — ask only blocking questions and record answers.
- `.agents/skills/skillgrid-codebase-map/SKILL.md` — map repo structure, graphify output, tests, and conventions.
- `.agents/skills/ccc/SKILL.md` — CocoIndex Code: `ccc init`, `ccc index`, `ccc search`; optional MCP `cocoindex-code`.
- `.agents/skills/references/indexing-and-memory.md` — Engram, graphify (`graphify-out/`, `graphify update .`), **ccc**, and MCP memory ordering.
- `.agents/skills/skillgrid-parallel-research/SKILL.md` — coordinate external research and long evidence capture.
- `.agents/skills/skillgrid-subagent-orchestration/SKILL.md` — dispatch bounded subagents with handoff paths and two-stage review.
- `.agents/skills/skillgrid-prd-artifacts/SKILL.md` — PRD numbering, `INDEX.md`, title blocks, and status lifecycle.
- `.agents/skills/skillgrid-spec-artifacts/SKILL.md` — PRD-to-OpenSpec artifacts and validation.
- `.agents/skills/skillgrid-vertical-slices/SKILL.md` — shippable slices, `[HITL]` / `[AFK]`, and testable increments.
- `.agents/skills/skillgrid-ui-design-artifacts/SKILL.md` — UI decisions, previews, `DESIGN.md`, and OpenSpec design constraints.
- `.agents/skills/skillgrid-issue-creation/SKILL.md` — local/GitHub/GitLab/Jira issue behavior from `.skillgrid/config.json`.
- `.agents/skills/skillgrid-hybrid-persistence/SKILL.md` — disk plus Engram persistence.
- `.agents/skills/skillgrid-filesystem-handoff/SKILL.md` — `context_<change-id>.md`, `events/<change-id>.jsonl`, and `research/<change-id>/`.
- `.agents/skills/skillgrid-openspec-config/SKILL.md` — `openspec/config.yaml` overlay rules.
- `.agents/skills/skillgrid-project-docs/SKILL.md` — `DESIGN.md` and `.skillgrid/project/*` docs.
- `.agents/skills/skillgrid-checkpoints/SKILL.md` — `.skillgrid/tasks/checkpoints.log`.


## Phase-Specific Skills

Load these first for this command:

- `skillgrid-questioning`
- `skillgrid-codebase-map`
- `ccc`
- `skillgrid-openspec-config`
- `skillgrid-project-docs`
- `skillgrid-hybrid-persistence`

## Event Log Rule

For any identified Skillgrid change id, create `.skillgrid/tasks/events/` if needed and append short JSONL events to `.skillgrid/tasks/events/<change-id>.jsonl` when this command starts, completes, blocks, skips, dispatches/receives subagents, or changes workflow state. If a delegated agent cannot write, require it to return a suggested event object and append that event from the parent session before advancing.

## Steps

1. Ask only for blocking setup choices: ticketing provider, artifact-store mode, and PRD workflow source.
2. For PRD workflow source, offer:
   - `skillgrid-default` preset: `draft`, `todo`, `inprogress`, `devdone`, `done`.
   - provider preset for GitHub/GitLab/Jira when the project wants remote-like local columns.
   - provider import when credentials/tooling can discover project columns or statuses.
   - custom ordered statuses with phase-to-status mapping.
3. Create or merge `.skillgrid/config.json` with local-first defaults when missing, including `prdWorkflow.statuses`, `fallbackStatus`, and `phaseStatusMap`.
4. If provider workflow import fails or is unavailable, record the reason and ask for a preset or custom fallback; do not fail init solely because remote status discovery failed.
5. If artifact store includes OpenSpec, initialize or align `openspec/` and `openspec/config.yaml`.
6. Create or refresh `.skillgrid/project/` docs and root `DESIGN.md` using discovered repo facts.
7. **CocoIndex (`ccc`) (semantic index):** If the **`ccc`** CLI is installed, from the repository root ensure the project is initialized (**`ccc init`** when needed—see [`ccc` skill](../.agents/skills/ccc/SKILL.md)), then run **`ccc index`** so **`ccc search`** and the optional [**`cocoindex-code`** MCP](../.configs/mcp/command/cocoindex-code.json) have a fresh index. If the CLI is missing or the user opts out, skip and note that in the completion report; do not fail init.
8. **graphify (initial index):** If the `graphify` CLI is installed (e.g. `uv tool install graphifyy`), run **`graphify .`** from the repository root—or in a Chat session, **`/graphify .`**. That creates or refreshes **`graphify-out/`** (including `graph.json` for the optional [graphify MCP](../.configs/mcp/python/graphify.json) and `GRAPH_REPORT.md` for agents). If the CLI is missing or the user opts out, skip and note that in the completion report; do not fail init.
9. Record durable setup decisions through Engram when available and appropriate.
10. Recommend `/skillgrid-explore` for brownfield mapping or `/skillgrid-plan` for a clear first change.

## Completion Report

Report config choices, PRD workflow source/statuses/phase mapping, provider import result or fallback, created/aligned artifacts, project docs touched, **ccc** / **graphify** indexing (if run), memory saves, and recommended next command.

</process>
