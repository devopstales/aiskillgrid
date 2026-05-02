---
description: Combined gate — full review, full security, then user validation
allowed-tools: Read, Write, Glob, Grep, Bash, Task
argument-hint: "[change-id or release scope]"
---

<objective>

You are executing **`/skillgrid-validate`** for the Skillgrid workflow.

Run the combined review gate: spec compliance, code quality, security, and user sign-off.

**Status on exit:** set the PRD and `INDEX.md` entry to `.skillgrid/config.json` `prdWorkflow.phaseStatusMap.validate` only when validation passes or accepted risks are explicit (default: `devdone`).

</objective>

<process>

## Shared Skillgrid Skills

Before acting, load only the skills needed for the phase:

- `.agents/skills/skillgrid-questioning/SKILL.md` — ask only blocking questions and record answers.
- `.agents/skills/skillgrid-codebase-map/SKILL.md` — map repo structure, GitNexus output, tests, and conventions.
- `.agents/skills/ccc/SKILL.md` — CocoIndex Code: `ccc init`, `ccc index`, `ccc search`; optional MCP `cocoindex-code`.
- `.agents/skills/references/indexing-and-memory.md` — Engram, GitNexus (`.gitnexus/`, `npx -y gitnexus@1.3.11 analyze`), **ccc**, and MCP memory ordering.
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
- `skillgrid-subagent-orchestration`
- `skillgrid-filesystem-handoff`
- `skillgrid-hybrid-persistence`
- review/security/test skills as needed

## Event Log Rule

For any identified Skillgrid change id, create `.skillgrid/tasks/events/` if needed and append short JSONL events to `.skillgrid/tasks/events/<change-id>.jsonl` when this command starts, completes, blocks, skips, dispatches/receives subagents, or changes workflow state. If a delegated agent cannot write, require it to return a suggested event object and append that event from the parent session before advancing.

## Steps

1. Re-read PRD, OpenSpec artifacts, tasks, handoff, and evidence.
2. Verify implementation against PRD goals, delta specs, and task checkboxes.
3. Apply two-stage review: spec compliance first, then quality/security/maintainability.
4. Request or run independent review with fresh artifact context, not chat history; use parallel reviewers only for independent perspectives.
5. Evaluate incoming review feedback before applying it; clarify unclear items, push back on incorrect or out-of-scope suggestions, and fix accepted issues one at a time.
6. Resolve, accept, or track every blocking issue; do not proceed with open critical or important findings.
7. Ask for user sign-off when required.
8. Save a concise verify report and update handoff/memory.

## Completion Report

Report validation result, issues found/resolved/accepted, evidence paths, configured status set for `validate` if passed, and recommended `/skillgrid-finish`.

</process>
