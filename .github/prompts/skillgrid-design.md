---
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
- `.agents/skills/skillgrid-parallel-research/SKILL.md` — coordinate external research and long evidence capture.
- `.agents/skills/skillgrid-subagent-orchestration/SKILL.md` — dispatch bounded subagents with handoff paths and two-stage review.
- `.agents/skills/skillgrid-prd-artifacts/SKILL.md` — PRD numbering, `INDEX.md`, title blocks, and status lifecycle.
- `.agents/skills/skillgrid-spec-artifacts/SKILL.md` — PRD-to-OpenSpec artifacts and validation.
- `.agents/skills/skillgrid-vertical-slices/SKILL.md` — shippable slices, `[HITL]` / `[AFK]`, and testable increments.
- `.agents/skills/skillgrid-ui-design-artifacts/SKILL.md` — UI decisions, previews, `DESIGN.md`, and OpenSpec design constraints.
- `.agents/skills/skillgrid-issue-creation/SKILL.md` — local/GitHub/GitLab/Jira issue behavior from `.skillgrid/config.json`.
- `.agents/skills/skillgrid-hybrid-persistence/SKILL.md` — disk plus Engram persistence.
- `.agents/skills/skillgrid-filesystem-handoff/SKILL.md` — `context_<change-id>.md` and `research/<change-id>/`.
- `.agents/skills/skillgrid-openspec-config/SKILL.md` — `openspec/config.yaml` overlay rules.
- `.agents/skills/skillgrid-project-docs/SKILL.md` — `DESIGN.md` and `.skillgrid/project/*` docs.
- `.agents/skills/skillgrid-checkpoints/SKILL.md` — `.skillgrid/tasks/checkpoints.log`.


## Phase-Specific Skills

Load these first for this command:

- `skillgrid-ui-design-artifacts`
- `skillgrid-questioning`
- `skillgrid-parallel-research`
- `skillgrid-project-docs`
- UI skills such as `superdesign`, `impeccable`, `frontend-design`, or `frontend-ui-engineering` as needed

## Steps

1. Identify the surface, audience, and design decision to make.
2. Load existing `DESIGN.md`, relevant project docs, and UI source context.
3. Generate or compare design options when useful; keep long comparisons in `.skillgrid/preview/` or research files.
4. Record the chosen direction in `DESIGN.md`, PRD success criteria, or OpenSpec `design.md` as appropriate.
5. Call out accessibility, responsive, implementation, and testing implications.

## Completion Report

Report chosen direction, preview artifacts, docs updated, unresolved design risks, and the next command.

</process>
