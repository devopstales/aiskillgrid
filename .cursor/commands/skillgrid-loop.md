---
name: /skillgrid-loop
id: skillgrid-loop
category: Workflow
description: Loop through the next safe Skillgrid phase or AFK slice until blocked or complete
allowed-tools: Read, Write, Edit, Glob, Grep, Bash, Task
argument-hint: "[change-id; optional phase or task focus]"
---

<objective>

You are executing **`/skillgrid-loop`** for the Skillgrid workflow.

Continue the active Skillgrid change through the next safe command, phase, or `[AFK]` implementation slice while preserving explicit stop conditions.

**Status on exit:** do not invent a new PRD status. Preserve the status contract of the command or phase that actually advanced.

</objective>

<process>

## Shared Skillgrid Skills

Before acting, load only the skills needed for the phase:

- `.agents/skills/skillgrid-questioning/SKILL.md` — intent gate, blocking questions, and decision-tree interview mode.
- `.agents/skills/skillgrid-codebase-map/SKILL.md` — map repo structure, GitNexus output, tests, and conventions.
- `.agents/skills/ccc/SKILL.md` — CocoIndex Code: `ccc init`, `ccc index`, `ccc search`; optional MCP `cocoindex-code`.
- `.agents/skills/references/indexing-and-memory.md` — Engram, GitNexus (`.gitnexus/`, `npx -y gitnexus@1.3.11 analyze`), **ccc**, and MCP memory ordering.
- `.agents/skills/skillgrid-subagent-orchestration/SKILL.md` — bounded subagents, handoff paths, and two-stage review.
- `.agents/skills/skillgrid-vertical-slices/SKILL.md` — `[HITL]` / `[AFK]`, shippable slices, and testable increments.
- `.agents/skills/skillgrid-filesystem-handoff/SKILL.md` — `context_<change-id>.md`, `events/<change-id>.jsonl`, and research spill files.
- `.agents/skills/skillgrid-checkpoints/SKILL.md` — `.skillgrid/tasks/checkpoints.log`.
- `.agents/skills/skillgrid-hybrid-persistence/SKILL.md` — disk plus Engram persistence.

## Phase-Specific Skills

Load these first for this command:

- `skillgrid-questioning`
- `skillgrid-filesystem-handoff`
- `skillgrid-vertical-slices`
- `skillgrid-subagent-orchestration`
- `skillgrid-checkpoints`
- implementation, test, review, or security skills only when the next safe step requires them

## Event Log Rule

For any identified Skillgrid change id, create `.skillgrid/tasks/events/` if needed and append short JSONL events to `.skillgrid/tasks/events/<change-id>.jsonl` when this command starts, completes, blocks, skips, dispatches/receives subagents, or changes workflow state. If a delegated agent cannot write, require it to return a suggested event object and append that event from the parent session before advancing.

## Steps

1. Run the `skillgrid-questioning` intent gate and identify the active change id from arguments, `.skillgrid/tasks/context_*.md`, `.skillgrid/prd/INDEX.md`, `openspec/changes/`, or the full retrieved Engram observation for `skillgrid/<change-id>/state`.
2. Read the active PRD, OpenSpec artifacts, `tasks.md`, handoff file, latest checkpoint, and cited research before taking any phase action. If Engram lookup is needed, use `mem_search` only to find IDs and `mem_get_observation(id)` before relying on the content.
3. Determine the next safe action:
   - no PRD or unclear scope -> route to `/skillgrid-brainstorm` or `/skillgrid-plan`;
   - PRD exists but tasks are incomplete -> route to `/skillgrid-breakdown`;
   - next task is `[AFK]` and apply-ready -> execute one `/skillgrid-apply`-style slice;
   - implementation exists but evidence is missing -> route to `/skillgrid-test`;
   - evidence exists but sign-off is missing -> route to `/skillgrid-validate`;
   - validation passed -> route to `/skillgrid-finish`.
4. Before any apply-style action, enforce the context budget gate: the slice must have acceptance criteria, context packet, bounded file list, verification command, and split trigger. If not, route back to breakdown.
5. Advance only one phase transition or one small `[AFK]` slice at a time.
6. After each unit, update the handoff with result, evidence, changed assumptions, blockers, and the next recommended action.
7. Append an event under `.skillgrid/tasks/events/<change-id>.jsonl` when loop starts, advances, blocks, dispatches or receives a subagent, records a review, or completes.
8. When Engram is available, save or update the compact `skillgrid/<change-id>/state` snapshot with phase, status, artifact store, PRD/OpenSpec/handoff paths, blockers, next action, and `last_updated`.
9. Reassess before continuing. Do not continue just because there are unchecked tasks.

## Completion Report

Report the phase or slice advanced, evidence paths, handoff/event updates, Engram state snapshot update or skip reason, the stop condition if any, and the next recommended command.

</process>
