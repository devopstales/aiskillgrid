# Context: <change-id>

## Goal

<One sentence describing what this change is trying to ship.>

## Current state

- **Phase:** plan | breakdown | apply | test | validate | finish
- **Status:** active | blocked | paused | ready-for-next-command
- **Last updated:** <YYYY-MM-DD>
- **Engram state:** `skillgrid/<change-id>/state` if Engram is available

## Active artifacts

- PRD: `.skillgrid/prd/PRD<NN>_<slug>.md`
- OpenSpec: `openspec/changes/<change-id>/`
- Tasks: `openspec/changes/<change-id>/tasks.md`
- Event log: `.skillgrid/tasks/events/<change-id>.jsonl`
- Preview: `.skillgrid/preview/<file>` if any

## Decisions

- <Decision> — <why> — <where recorded>

## Failed approaches

- <What was attempted, why it failed, and what to avoid repeating>

## HITL blockers

- [ ] `<owner>` <decision/action needed before AFK work can continue>

## AFK-ready work

- [ ] <slice/task ready for autonomous implementation> — context packet: <paths> — verification: <command>

## Dependency waves

- **Wave 1:** <independent slices and file ownership>
- **Blocked by:** <relationships that must clear before later waves>

## Research index

- `.skillgrid/tasks/research/<change-id>/<topic>.md` — <one-line finding>

## Last checkpoint

- `<checkpoint-name>` — <timestamp or git SHA>

## Next steps

1. <Immediate action, e.g. `/sdd-verify`>
2. <Follow-up action>
3. <Verification or decision checkpoint>

## Next recommended action

<Single best next command or task, e.g. `/skillgrid-apply <change-id> task 2`.>
