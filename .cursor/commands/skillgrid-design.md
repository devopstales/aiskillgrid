---
name: /skillgrid-design
id: skillgrid-design
category: Design
description: On-demand design workshop — extract, browse, tune, design, audit
allowed-tools: Read, Write, Glob, Grep, Bash, Task, WebFetch
argument-hint: "[extract | browse | tune | design | audit]"
---

<objective>

You are executing **`/skillgrid-design`** for the Skillgrid workflow.

Run a focused design workshop for UI/UX direction and design-system artifacts.

**Status on exit:** no PRD status change unless the user asks to update PRD criteria.

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

- `skillgrid-codebase-map`
- `ccc`
- `skillgrid-ui-design-artifacts`
- `skillgrid-questioning`
- `skillgrid-parallel-research`
- `skillgrid-project-docs`
- UI skills such as `huashu-design`, `superdesign`, `impeccable`, `frontend-design`, or `frontend-ui-engineering` as needed

## Event Log Rule

For any identified Skillgrid change id, create `.skillgrid/tasks/events/` if needed and append short JSONL events to `.skillgrid/tasks/events/<change-id>.jsonl` when this command starts, completes, blocks, skips, dispatches/receives subagents, or changes workflow state. If a delegated agent cannot write, require it to return a suggested event object and append that event from the parent session before advancing.

## Steps

1. Identify the surface, audience, and design decision to make.
2. Load existing `DESIGN.md`, relevant project docs, and UI source context.
3. When the user asks for multiple UI options, A/B/C directions, or browser-viewable previews, run or offer `.skillgrid/scripts/preview.sh <change-or-surface-slug>` to scaffold a default HTML preview under `.skillgrid/preview/`. Use `.skillgrid/scripts/preview.sh --md <slug>` only when a text-only comparison is enough.
4. Generate or compare three differentiated design options when useful. Prefer `huashu-design` for HTML-based design variants and direction-advisor work; fill scaffolded preview HTMLs, attach SuperDesign exports, or keep long comparisons in `.skillgrid/preview/` or research files.
5. Record the chosen direction in `DESIGN.md`, PRD success criteria, or OpenSpec `design.md` as appropriate after the user picks or approves a direction.
6. Call out accessibility, responsive, implementation, and testing implications.

## Completion Report

Report chosen direction, preview artifacts, docs updated, unresolved design risks, and the next command. If an HTML preview was generated, include the exact file path and prompt the user to open it locally for review.

</process>
