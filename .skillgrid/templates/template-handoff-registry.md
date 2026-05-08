# Registry: <change-id>

Compact index for recent decisions and progress so the parent can inject focused context into subagents without re-reading long reports.

## Snapshot

- **Last updated:** <YYYY-MM-DD>
- **Phase:** plan | breakdown | apply | test | validate | finish
- **Status:** active | blocked | paused | ready-for-next-command
- **Handoff:** `.skillgrid/tasks/context_<change-id>.md`
- **Event log:** `.skillgrid/tasks/events/<change-id>.jsonl`

## Recent decisions (latest first)

- <decision> — <why> — <where>

## Recent progress (latest first)

- <task/slice> — <status> — <evidence path>

## Active blockers

- [ ] <owner> <blocker summary>

## Safe next dispatch candidates

- <subagent role/persona> — <task label> — <expected output path>

## Context injection packet seed

- Goal: <one sentence objective>
- Constraints: <must-follow boundaries>
- Required artifacts:
  - `.skillgrid/tasks/context_<change-id>.md`
  - `.skillgrid/tasks/events/<change-id>.jsonl`
  - `.skillgrid/tasks/registry_<change-id>.md`
  - `openspec/changes/<change-id>/tasks.md`
  - `openspec/changes/<change-id>/specs/<slice>/spec.md` (if relevant)
