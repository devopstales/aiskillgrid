---
description: Explore the problem and repo: OpenSpec list, PRD backfill, .skillgrid/project, AGENTS, graphify
allowed-tools: Read, Write, Glob, Grep, Bash, Task
argument-hint: "[optional: topic, change id, or area to explore]"
---

<objective>

You are executing **`/skillgrid-explore`** for the Skillgrid workflow.

Map the problem and repo without implementing product behavior. Exploration may update planning artifacts, PRDs, and project docs.

**Status on exit:** no mandatory status change; backfilled PRDs use `skillgrid-prd-artifacts` status rules.

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

- `skillgrid-codebase-map`
- `skillgrid-parallel-research`
- `skillgrid-prd-artifacts`
- `skillgrid-project-docs`
- `skillgrid-filesystem-handoff`

## Steps

1. Select the topic, change, or repo area to explore; ask only if the target is unclear.
2. Inventory OpenSpec changes, PRDs, project docs, handoff files, tests, and graphify output where present.
3. Backfill missing PRDs only when needed, following `skillgrid-prd-artifacts`.
4. Use parallel research/subagents for independent repo, docs, or web questions; spill long output to `.skillgrid/tasks/research/<change-id>/`.
5. Refresh `.skillgrid/project/*` and `DESIGN.md` only for stable discovered facts.
6. Stop when the scope is clear enough to recommend `/skillgrid-brainstorm` or `/skillgrid-plan`.

## Completion Report

Summarize explored areas, artifacts inspected or updated, important findings, open questions, and the recommended next command.

</process>
