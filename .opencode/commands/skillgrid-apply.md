---
name: /skillgrid-apply
id: skillgrid-apply
category: Workflow
description: Implement from tasks; OpenSpec apply loop; keep PRD and tasks.md in sync
allowed-tools: Read, Write, Edit, Glob, Grep, Bash, Task
argument-hint: "[change-id; optional task focus]"
---

<objective>

You are executing **`/skillgrid-apply`** for the Skillgrid workflow.

Implement from approved Skillgrid and OpenSpec tasks in small verified batches.

**Status on exit:** set the PRD and `INDEX.md` entry to `.skillgrid/config.json` `prdWorkflow.phaseStatusMap.apply` once implementation starts; keep it there until validation (default: `inprogress`).

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

- `skillgrid-checkpoints`
- `ccc`
- `skillgrid-spec-artifacts`
- `skillgrid-vertical-slices`
- `skillgrid-filesystem-handoff`
- `skillgrid-subagent-orchestration`
- `skillgrid-hybrid-persistence`
- implementation skills such as `test-driven-development`, `api-and-interface-design`, or UI skills as needed

## Steps

1. Confirm PRD, OpenSpec change, and `tasks.md` are apply-ready.
2. Create `before-apply-<change-id>` checkpoint before edits.
3. Read the handoff and any cited research.
4. Critically review the selected task before editing; stop if instructions, expected tests, or scope are unclear.
5. Implement one vertical slice or small task batch at a time.
6. For behavioral code, use TDD: write one failing test, verify RED, implement minimum code, verify GREEN, then refactor.
7. When delegating implementation or review, use `skillgrid-subagent-orchestration` templates: implementer, spec reviewer, then code quality reviewer.
8. Evaluate review feedback technically, send accepted required fixes back to the implementer, and re-review before marking any task complete.
9. Update task checkboxes and handoff state as work completes.
10. Pause on ambiguous requirements, failing verification, or HITL blockers.

## Completion Report

Report implemented slices, tests or checks run, task updates, handoff updates, blockers, and recommended `/skillgrid-test` or next apply slice.

</process>
