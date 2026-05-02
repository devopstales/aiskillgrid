---
name: skillgrid-hybrid-persistence
description: >
  Coordinates Skillgrid disk artifacts and Engram memory according to .skillgrid/config.json artifact-store mode.
  Trigger: Initializing, planning, applying, validating, or finishing Skillgrid work that must survive context compaction.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

## When to Use

Use this skill whenever a Skillgrid command needs to decide where state belongs: PRDs, OpenSpec files, handoff files, checkpoints, Engram memory, or a combination of those.

## Critical Patterns

### Artifact Store Modes

Read `.skillgrid/config.json` when it exists.

| Mode | Meaning |
|---|---|
| `hybrid` or missing | Strongly recommended default. Use on-disk Skillgrid/OpenSpec artifacts and Engram summaries. |
| `openspec` | Use on-disk artifacts first; Engram is optional for durable decisions if available. |
| `engram` | Use Engram for durable summaries; do not assume `openspec/` exists. |

If `artifactStore` is missing, treat `hybrid` as the default and recommend `/skillgrid-init` to write the config.

### What Goes Where

| State | Primary home |
|---|---|
| Product intent | `.skillgrid/prd/PRD<NN>_<slug>.md` |
| Technical contract | `openspec/changes/<change-id>/` |
| Active session state | `.skillgrid/tasks/context_<change-id>.md` |
| Long research output | `.skillgrid/tasks/research/<change-id>/` |
| Named safety points | `.skillgrid/tasks/checkpoints.log` |
| Durable cross-session memory | Engram `mem_save` |

Do not paste long research or complete handoff files into Engram. Save concise durable facts with paths back to disk artifacts.

In `engram` mode, save concrete equivalents rather than vague summaries:

| Disk artifact | Engram topic key |
|---|---|
| `.skillgrid/prd/PRD<NN>_<slug>.md` | `skillgrid/<change-id>/prd` |
| `.skillgrid/prd/INDEX.md` | `skillgrid/index` |
| `openspec/changes/<change-id>/proposal.md` | `skillgrid/<change-id>/proposal` |
| `openspec/changes/<change-id>/specs/*/spec.md` | `skillgrid/<change-id>/spec` |
| `openspec/changes/<change-id>/tasks.md` | `skillgrid/<change-id>/tasks` |
| `.skillgrid/tasks/context_<change-id>.md` | `skillgrid/<change-id>/context` |
| `.skillgrid/tasks/research/<change-id>/...` | `skillgrid/<change-id>/research/<topic>` |
| `.skillgrid/tasks/events/<change-id>.jsonl` | `skillgrid/<change-id>/events` |
| `.skillgrid/tasks/checkpoints.log` | `skillgrid/<change-id>/checkpoint` |
| `.skillgrid/preview/*.html` / `DESIGN.md` changes | `skillgrid/<change-id>/design` |
| Compact resumable state | `skillgrid/<change-id>/state` |

### Engram Topic Keys

Use stable, predictable keys:

```text
skillgrid-init/<project-name>
skillgrid/<change-id>/plan
skillgrid/<change-id>/tasks
skillgrid/<change-id>/apply
skillgrid/<change-id>/verify-report
skillgrid/<change-id>/security
skillgrid/<change-id>/archive
skillgrid/<change-id>/state
```

### Compact State Snapshot

Use `skillgrid/<change-id>/state` as the small recovery index for a change. It is not the PRD, spec, tasks, or handoff itself; it points to them and records the minimum state needed after compaction or a cold start.

Save it when a Skillgrid command creates, advances, blocks, pauses, validates, or finishes a change. Keep the content short and structured:

```yaml
change: <change-id>
phase: plan | breakdown | apply | test | validate | finish
status: active | blocked | paused | ready-for-next-command | done
artifact_store: hybrid | openspec | engram
prd: .skillgrid/prd/PRD<NN>_<slug>.md
openspec: openspec/changes/<change-id>/
handoff: .skillgrid/tasks/context_<change-id>.md
event_log: .skillgrid/tasks/events/<change-id>.jsonl
blockers:
  - <short blocker or none>
next_action: <next command or slice>
last_updated: <ISO date>
```

In `hybrid` and `openspec` modes, the snapshot should point back to disk artifacts rather than duplicate their full content. In `engram` mode, use the snapshot as an index to the concrete Engram artifact equivalents listed above.

### Save Triggers

Save to Engram when:

- a plan or PRD is accepted
- tasks materially change
- a design or architecture decision is made
- validation/security findings are resolved or accepted
- a change is archived, shipped, discarded, or intentionally paused
- the compact `skillgrid/<change-id>/state` snapshot changes phase, status, blocker, or next action

Use `memory-protocol` for low-level Engram tool rules.

### Recovery Protocol

Engram search results are previews. Before relying on any recovered Skillgrid artifact or state snapshot:

1. Use `mem_search` to find candidate observations by stable topic key.
2. Use `mem_get_observation(id)` for the full content of each selected observation.
3. Reconcile recovered memory with `.skillgrid/config.json`, PRD files, OpenSpec files, and the handoff before acting.

Never treat a truncated `mem_search` snippet as authoritative for requirements, task status, blocker state, or user decisions.

### Context Rot Guardrail

When chat context grows, prefer writing state to disk or Engram over restating it in chat.

Before continuing after a long pause or compaction:

1. Read `.skillgrid/config.json`.
2. Read relevant PRD and OpenSpec files.
3. Read `.skillgrid/tasks/context_<change-id>.md` if present.
4. Search Engram for `skillgrid/<change-id>/state` and any needed artifact keys if Engram is available, then retrieve full observations by ID.
5. Continue from artifacts, not from memory of the conversation.

## Commands

```bash
openspec list --json
openspec status --change "<change-id>" --json
```

Engram access happens through MCP tools when configured.

## Resources

- Engram protocol: `memory-protocol`
- Handoff files: `skillgrid-filesystem-handoff`
- Checkpoints: `skillgrid-checkpoints`
- Workflow overview: `docs/02-workflow-usage.md`
