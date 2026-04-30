---
description: PRD implementation checklist synced with OpenSpec tasks.md
allowed-tools: Read, Write, Glob, Grep, Bash, Task
argument-hint: "[change-id or PRD slug]"
---

<objective>

You are executing **`/skillgrid-breakdown`** for the Skillgrid workflow.

Turn an accepted PRD/OpenSpec change into an implementation checklist synced with `tasks.md`.

**Status on exit:** set the PRD and `INDEX.md` entry to `.skillgrid/config.json` `prdWorkflow.phaseStatusMap.breakdown` when the checklist is approved (default: `todo`).

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
- `skillgrid-spec-artifacts`
- `skillgrid-prd-artifacts`
- `skillgrid-vertical-slices`
- `skillgrid-issue-creation`
- `skillgrid-filesystem-handoff`

## Event Log Rule

For any identified Skillgrid change id, create `.skillgrid/tasks/events/` if needed and append short JSONL events to `.skillgrid/tasks/events/<change-id>.jsonl` when this command starts, completes, blocks, skips, dispatches/receives subagents, or changes workflow state. If a delegated agent cannot write, require it to return a suggested event object and append that event from the parent session before advancing.

## Steps

1. Select the change and PRD; stop if either is missing.
2. Use `openspec status` and instructions to identify incomplete artifacts.
3. Build or refresh `openspec/changes/<change-id>/tasks.md` from PRD, proposal, specs, and design.
4. Keep the PRD Implementation tasks section and OpenSpec `tasks.md` aligned.
5. Tag tasks `[HITL]` or `[AFK]`, ordering human decisions before dependent autonomous work.
6. Create external slice issues only when configured.
7. Update the handoff with blockers, AFK-ready work, and next apply target.

## Completion Report

Report selected PRD/change, task count, HITL blockers, AFK-ready slices, configured status set for `breakdown`, and recommended `/skillgrid-apply`.

</process>
