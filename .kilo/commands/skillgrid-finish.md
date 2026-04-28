---
name: /skillgrid-finish
id: skillgrid-finish
category: Workflow
description: Close the change: optional spec sync, archive, git hygiene, ship checklist, PR
allowed-tools: Read, Write, Glob, Grep, Bash, Task
argument-hint: "[change-id or branch name]"
---

<objective>

You are executing **`/skillgrid-finish`** for the Skillgrid workflow.

Close a Skillgrid change: final docs alignment, optional spec sync/archive, tracker hygiene, PR handoff, and checkpoint cleanup.

**Status on exit:** set the PRD and `INDEX.md` entry to `done` when the change is finished successfully.

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
- `.agents/skills/skillgrid-filesystem-handoff/SKILL.md` — `context_<change-id>.md` and `research/<change-id>/`.
- `.agents/skills/skillgrid-openspec-config/SKILL.md` — `openspec/config.yaml` overlay rules.
- `.agents/skills/skillgrid-project-docs/SKILL.md` — `DESIGN.md` and `.skillgrid/project/*` docs.
- `.agents/skills/skillgrid-checkpoints/SKILL.md` — `.skillgrid/tasks/checkpoints.log`.


## Phase-Specific Skills

Load these first for this command:

- `ccc`
- `skillgrid-spec-artifacts`
- `skillgrid-project-docs`
- `skillgrid-checkpoints`
- `skillgrid-issue-creation`
- `skillgrid-hybrid-persistence`
- `git-workflow-and-versioning`, `documentation-and-adrs`, `openspec-sync-specs`, `openspec-archive-change` as needed

## Steps

1. Re-read PRD, OpenSpec change, handoff, validation evidence, and project docs.
2. Sync delta specs to main specs when requested or required.
3. Refresh `.skillgrid/project/*` and `DESIGN.md` if the finished work changed durable project facts.
4. Archive OpenSpec changes when appropriate.
5. Update tracker state according to `.skillgrid/config.json`.
6. Clean up change-scoped checkpoints and previews when safe.
7. Save Engram closure summary when available.
8. Prepare PR/merge/deploy guidance only when requested by the user.

## Completion Report

Report spec sync/archive result, docs alignment, tracker updates, checkpoint cleanup, PR/CI status, status set to `done`, and remaining risks.

</process>
