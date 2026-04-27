---
name: skillgrid-checkpoints
description: >
  Manages Skillgrid checkpoint log entries for safe pause, resume, apply, validation, and cleanup.
  Trigger: Creating, verifying, listing, or cleaning Skillgrid checkpoints.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

## When to Use

Use this skill when a command needs to create or reason about `.skillgrid/tasks/checkpoints.log`.

## Critical Patterns

### Canonical Log

```text
.skillgrid/tasks/checkpoints.log
```

Checkpoints complement PRDs, OpenSpec changes, and handoff files. They do not create commits unless the user explicitly asks.

### When To Create

Create checkpoints:

- before `/skillgrid-apply` starts implementation
- before high-risk validation or cleanup
- before long pauses where resume state matters
- when the user explicitly asks for one

Recommended apply checkpoint:

```text
before-apply-<change-id>
```

### Entry Contents

Each checkpoint should record:

- name
- timestamp
- branch
- git SHA when available
- dirty status
- active PRD
- active OpenSpec change
- handoff file
- optional verification evidence

Use a stable line-oriented format so it can be appended and inspected without special tooling:

```text
<timestamp> name=<checkpoint-name> branch=<branch> sha=<git-sha-or-unknown> dirty=<yes|no> prd=<path-or-none> change=<change-id-or-none> context=<path-or-none> evidence="<short note>"
```

Example:

```text
2026-04-27T15:30:00Z name=before-apply-auth-foundation branch=feature/auth sha=abc123 dirty=yes prd=.skillgrid/prd/PRD01_auth-foundation.md change=auth-foundation context=.skillgrid/tasks/context_auth-foundation.md evidence="lint baseline passed"
```

### Verification Report Template

When verifying a checkpoint, report:

```markdown
## Checkpoint Verification: <checkpoint-name>

- **Recorded branch:** `<branch>`
- **Current branch:** `<branch>`
- **Recorded SHA:** `<sha>`
- **Current SHA:** `<sha>`
- **Dirty then:** `<yes|no>`
- **Dirty now:** `<yes|no>`
- **Context file:** `<path>`

## Drift

- <None, or list changed conditions>

## Recommendation

<Continue, refresh checkpoint, inspect drift, or pause.>
```

### Cleanup

During `/skillgrid-finish`, remove or archive checkpoints that belong only to the completed change. Keep unrelated checkpoints.

### Verification

Checkpoint verification compares current state with the recorded state and reports drift. It should not revert files.

## Commands

```bash
git status --short
git rev-parse --abbrev-ref HEAD
git rev-parse HEAD
```

## Resources

- Handoff: `skillgrid-filesystem-handoff`
- Persistence: `skillgrid-hybrid-persistence`
- Command source: `.cursor/commands/skillgrid-checkpoint.md`
