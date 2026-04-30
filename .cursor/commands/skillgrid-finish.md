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

**Status on exit:** set the PRD and `INDEX.md` entry to `.skillgrid/config.json` `prdWorkflow.phaseStatusMap.finish` when the change is finished successfully (default: `done`).

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

- `ccc`
- `skillgrid-spec-artifacts`
- `skillgrid-project-docs`
- `skillgrid-checkpoints`
- `skillgrid-issue-creation`
- `skillgrid-hybrid-persistence`
- `git-master`, `git-workflow-and-versioning`, `documentation-and-adrs`, `openspec-sync-specs`, `openspec-archive-change` as needed

## Event Log Rule

For any identified Skillgrid change id, create `.skillgrid/tasks/events/` if needed and append short JSONL events to `.skillgrid/tasks/events/<change-id>.jsonl` when this command starts, completes, blocks, skips, dispatches/receives subagents, or changes workflow state. If a delegated agent cannot write, require it to return a suggested event object and append that event from the parent session before advancing.

## Steps

1. Re-read PRD, OpenSpec change, handoff, validation evidence, and project docs.
2. Sync delta specs to main specs when requested or required.
3. Refresh `.skillgrid/project/*` and `DESIGN.md` if the finished work changed durable project facts.
4. Archive OpenSpec changes when appropriate.
5. Update tracker state according to `.skillgrid/config.json`.
6. Clean up change-scoped checkpoints and previews when safe.
7. Save Engram closure summary when available:
   - Update `skillgrid/<change-id>/state` to the final phase/status and next action.
   - Save `skillgrid/<change-id>/archive` or equivalent closure notes with paths to final PRD, OpenSpec archive, validation evidence, and remaining risks.
   - If recovering prior memory during finish, use `mem_search` to find IDs and `mem_get_observation(id)` before relying on any recovered content.
8. If `.engram/` team sharing is used or requested, recommend `engram sync` after significant finish work; do not run it unless the user or project policy explicitly asks because `.engram/` may contain sensitive context.
9. Use `git-master` for branch, diff, commit, and PR hygiene. Create commits or PRs only when explicitly requested by the user.

## Completion Report

Report spec sync/archive result, docs alignment, tracker updates, checkpoint cleanup, Engram closure/state update or skip reason, optional `engram sync` guidance, PR/CI status, configured status set for `finish`, and remaining risks.

</process>
