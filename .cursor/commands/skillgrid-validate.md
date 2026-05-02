---
name: /skillgrid-validate
id: skillgrid-validate
category: Workflow
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

1. Re-read PRD, OpenSpec artifacts, tasks, handoff, event log, review reports, and verification evidence.
2. Confirm no active context budget, HITL, scope, schema drift, security, credential, destructive-action, or release blocker remains unresolved.
3. Run spec compliance review first. Check PRD/OpenSpec/task traceability, acceptance criteria, explicit deferrals, and assigned slice boundaries. Answer: did we build the right thing and only the assigned thing?
4. If spec compliance finds required changes, stop validation, record findings, fix or assign focused follow-up work, rerun verification, and repeat spec compliance review.
5. Only after spec compliance passes, run code quality review. Check correctness, maintainability, architecture, security, performance, test quality, and local conventions. Answer: is the implementation good enough to keep?
6. If code quality finds required changes, stop validation, record findings, fix or assign focused follow-up work, rerun verification, and repeat code quality review.
7. Run or reconcile security and test evidence as needed for the change risk.
8. For user-facing behavior, collect UAT/manual QA notes and convert failed UAT into focused fix tasks or plans.
9. Request or run independent fresh-context review with artifact paths and pushed coding standards, not copied chat history. Use parallel reviewers only for independent perspectives.
10. Resolve, accept with rationale, or track every critical or important finding; do not proceed with open blocking issues.
11. Ask for user sign-off when required.
12. Save a concise verify report and update handoff, events, PRD status, and memory.

## Completion Report

Report validation result, issues found/resolved/accepted, evidence paths, configured status set for `validate` if passed, and recommended `/skillgrid-finish`.

</process>
