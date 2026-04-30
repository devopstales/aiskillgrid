---
name: /skillgrid-session
id: skillgrid-session
category: Workflow
description: Phase 0 — session charter, context budget, MCP/tool selection, checkpoints
allowed-tools: Read, Glob, Grep, Write, Task
argument-hint: "[optional: session goal or resume context]"
---

<objective>

You are executing **`/skillgrid-session`** for the Skillgrid workflow.

Restore enough project and change context for a fresh or resumed agent session without assuming chat history.

**Status on exit:** no PRD status change; recommend the next Skillgrid command.

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

- `skillgrid-hybrid-persistence`
- `skillgrid-filesystem-handoff`
- `skillgrid-questioning`
- `skillgrid-codebase-map`
- `ccc`

## Event Log Rule

For any identified Skillgrid change id, create `.skillgrid/tasks/events/` if needed and append short JSONL events to `.skillgrid/tasks/events/<change-id>.jsonl` when this command starts, completes, blocks, skips, dispatches/receives subagents, or changes workflow state. If a delegated agent cannot write, require it to return a suggested event object and append that event from the parent session before advancing.

## Steps

1. Read `.skillgrid/config.json` if present; report `ticketing.provider`, `artifactStore.mode`, and `prdWorkflow` status columns / phase mapping.
2. List active PRDs under `.skillgrid/prd/` and note their `Status:` values.
3. Check `.skillgrid/tasks/context_*.md` for resumable changes and last checkpoint references.
4. If Engram is available, search for project and active-change keys; load details only when needed.
5. Keep context minimal: load deeper skills, docs, and OpenSpec files only for the next substantive phase.

## Completion Report

Report loaded context, active changes, handoff files, memory lookups, and the recommended next command.

</process>
