---
name: skillgrid-filesystem-handoff
description: >
  Maintains Skillgrid per-change handoff files and research spill files for parent sessions and subagents.
  Trigger: Creating, resuming, delegating, or validating work for a Skillgrid change.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

## When to Use

Use this skill whenever a Skillgrid change needs durable in-repo state shared between the main session, subagents, and later resumed sessions.

## Critical Patterns

### Canonical Paths

```text
.skillgrid/tasks/context_<change-id>.md
.skillgrid/tasks/events/<change-id>.jsonl
.skillgrid/tasks/research/<change-id>/
```

The handoff file is the rolling state for one change. The event log is the append-only workflow timeline. The research directory holds long reports, scrapes, browser evidence, or subagent output.

Treat the handoff as the current-state file for the change. It should answer: where are we, what is blocked, what evidence exists, and what should happen next.

### Handoff Contents

Keep `context_<change-id>.md` concise and skimmable. Use this template:

```markdown
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

## HITL blockers

- [ ] `<owner>` <decision/action needed before AFK work can continue>

## AFK-ready work

- [ ] <slice/task ready for autonomous implementation>

## Research index

- `.skillgrid/tasks/research/<change-id>/<topic>.md` — <one-line finding>

## Last checkpoint

- `<checkpoint-name>` — <timestamp or git SHA>

## Next recommended action

<Next command or task, e.g. `/skillgrid-apply <change-id> task 2`.>
```

Do not turn the handoff into a raw transcript. Link to research files when details are long.

### Engram State Alignment

When Engram is available, mirror only the compact current-state index to `skillgrid/<change-id>/state`. The handoff remains the in-repo current-state file; the Engram snapshot helps a fresh session find the right files and last known phase after compaction.

Update the snapshot whenever the handoff phase, status, blocker list, active artifacts, or next recommended action changes. Do not copy the full handoff into Engram unless the project explicitly runs in `artifactStore.mode: engram`.

### Required Handoff Event Log

Append workflow events to:

```text
.skillgrid/tasks/events/<change-id>.jsonl
```

Create `.skillgrid/tasks/events/` before the first append if it does not exist. Every Skillgrid command, parent session, and write-capable Skillgrid subagent must append an event when it starts work, completes a phase, blocks, skips a phase, or changes workflow state. If a subagent is read-only or cannot write, it must return a suggested event object and the parent must append it before advancing.

Each line is one JSON object. Keep entries short and stable so the local dashboard can render them:

```json
{"time":"2026-04-28T08:00:00Z","changeId":"planning-dashboard","prd":"PRD02_dashboard.md","node":"apply","phase":"apply","status":"started","summary":"Started AFK slice 2","artifacts":[".skillgrid/tasks/context_planning-dashboard.md"]}
{"time":"2026-04-28T08:12:00Z","changeId":"planning-dashboard","node":"apply","phase":"apply","status":"blocked","summary":"Missing HITL design approval","blocker":"Approve option B preview","artifacts":[".skillgrid/preview/planning-dashboard-options.html"]}
```

Supported fields:

- `time`: ISO timestamp.
- `changeId`: OpenSpec/Skillgrid change id.
- `prd`: optional PRD filename.
- `node`: workflow node or command id, for example `plan`, `apply`, `test`, `validate`.
- `phase`: Skillgrid phase.
- `status`: `started`, `progress`, `passed`, `failed`, `blocked`, `completed`, or `skipped`.
- `summary`: one-line event summary.
- `blocker`: optional HITL or technical blocker.
- `artifacts`: optional paths to PRDs, handoff files, previews, reports, checks, or research.
- `external`: optional issue key or URL.
- `agent` / `subagent`: optional subagent name for the dashboard Subagents view.
- `role`: optional role such as `implementer`, `spec-reviewer`, `code-quality-reviewer`, `design-critic`, or `security-auditor`.
- `task`: optional short delegated task label.
- `output`: optional primary output path when it is not already listed in `artifacts`.

Do not replace `context_<change-id>.md` with events. The event log is the timeline; the handoff remains the current-state summary.

### Slice Completion Summary

After every completed vertical slice, delegated task, review loop, or blocked apply attempt, update the handoff with a compact reassessment:

```markdown
## Last slice summary

- **Slice/task:** <identifier or title>
- **Result:** completed | blocked | needs-review | reverted
- **Evidence:** <test command, preview, checkpoint, review report, or research file>
- **Changed assumptions:** <none, or what changed since planning>
- **Blockers:** <none, or HITL/technical blocker>
- **Next recommended action:** <next command or slice>
```

This replaces long transcript-style progress logs. Keep details in research files and link them from the summary.

### Research Spill Files

Use `.skillgrid/tasks/research/<change-id>/` for:

- external web research
- documentation extracts
- browser testing reports
- subagent reports
- large comparison tables
- security or performance findings that exceed a short summary

Each research file should include:

- topic
- date or “as of” note when relevant
- sources or local files inspected
- findings
- recommendation
- open questions

Template:

```markdown
# Research: <topic>

- **Change:** `<change-id>`
- **Date:** <YYYY-MM-DD>
- **Question:** <decision this research informs>
- **Method:** repo search | web search | docs lookup | browser test | subagent review

## Sources

- `<local path>` — <why inspected>
- [Source title](https://example.com) — accessed <date>

## Findings

- <Evidence-backed finding>

## Recommendation

<What should the parent session do next?>

## Open questions

- <Question and owner, or `None`>

## Handoff update

Add to `.skillgrid/tasks/context_<change-id>.md`:

- `<this file>` — <one-line finding>
```

### Subagent Contract

Every subagent prompt for a Skillgrid change should include:

- the handoff path
- the event log path
- the PRD path
- the OpenSpec change path when present
- the expected output file under `research/<change-id>/`
- a requirement to append a short event when write-capable, or return a suggested event object when read-only
- a requirement to return a short summary with file paths

After a subagent returns, the parent session must read the handoff and cited research files before editing product code or changing workflow state.

### Parent Session Rule

The parent session owns implementation decisions. Subagents gather evidence, critique, test, or draft bounded artifacts. The parent reconciles their output against PRD and OpenSpec artifacts.

## Commands

```bash
ls .skillgrid/tasks
ls .skillgrid/tasks/research/<change-id>
```

## Resources

- Parallel research: `skillgrid-parallel-research`
- Subagent orchestration: `skillgrid-subagent-orchestration`
- Persistence modes: `skillgrid-hybrid-persistence`
- Workflow overview: `docs/workflow.md`
