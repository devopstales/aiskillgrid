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
| `hybrid` or missing | Use on-disk Skillgrid/OpenSpec artifacts and Engram summaries. |
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
```

### Save Triggers

Save to Engram when:

- a plan or PRD is accepted
- tasks materially change
- a design or architecture decision is made
- validation/security findings are resolved or accepted
- a change is archived, shipped, discarded, or intentionally paused

Use `memory-protocol` for low-level Engram tool rules.

### Context Rot Guardrail

When chat context grows, prefer writing state to disk or Engram over restating it in chat.

Before continuing after a long pause or compaction:

1. Read `.skillgrid/config.json`.
2. Read relevant PRD and OpenSpec files.
3. Read `.skillgrid/tasks/context_<change-id>.md` if present.
4. Search Engram for the change key if Engram is available.
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
- Workflow overview: `docs/workflow.md`
