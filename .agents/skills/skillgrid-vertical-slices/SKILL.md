---
name: skillgrid-vertical-slices
description: >
  Splits Skillgrid PRDs, specs, and tasks into small independently testable vertical slices.
  Trigger: Breaking down work, creating issues, planning implementation order, or deciding whether a PRD is too large.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

## When to Use

Use this skill during `/skillgrid-plan`, `/skillgrid-breakdown`, issue creation, or implementation planning when work must be split into reviewable, shippable increments.

## Critical Patterns

### Definition

A good Skillgrid vertical slice produces:

- one observable capability
- one testable path
- one clean checkpoint

Prefer slices that cross layers when needed:

```text
schema + API + UI + test evidence
```

Avoid plans that build all backend first, all frontend later, and only become usable at the end.

### Split PRDs Early

If a PRD contains multiple independent capabilities, split it into ordered PRDs:

```text
PRD01_auth-foundation.md
PRD02_invite-users.md
PRD03_billing-settings.md
```

Use `skillgrid-prd-artifacts` to renumber and update `INDEX.md`.

### Roadmap Hierarchy

Use this lightweight hierarchy when planning larger work:

```text
PRD sequence -> PRD slice -> tasks.md work items
```

- A sequence of ordered PRDs can represent a milestone or roadmap outcome.
- Each PRD should own one independently reviewable product slice.
- `openspec/changes/<change-id>/tasks.md` should contain task-level work for that slice.

Do not create a separate runtime or database to model this hierarchy. The ordered PRD index, OpenSpec change, and handoff file are enough.

### Task Tags

Use tags in `tasks.md` and PRD implementation tasks when the project follows HITL/AFK discipline:

- `[HITL]` human decision, credential, approval, or manual action required
- `[AFK]` autonomous work safe for an agent to execute from existing artifacts

Order `[HITL]` tasks before dependent `[AFK]` work.

### Done Criteria Per Slice

Each slice should state:

- user-visible result
- files or modules likely touched
- tests or verification command
- docs or artifact update
- checkpoint name when relevant
- rollback or validation risk if high

After completion, the handoff should record a slice summary: what changed, evidence, blockers, changed assumptions, and the next recommended slice.

### Slice Template

Use this when writing PRD implementation tasks, OpenSpec `tasks.md`, or issue bodies:

```markdown
### Slice <N>: <user-visible capability>

- **Tag:** `[AFK]` or `[HITL]`
- **Outcome:** <what a user/operator can observe after this slice>
- **Scope:** <small set of behavior included>
- **Out of scope:** <nearby behavior intentionally excluded>
- **Likely files:** `<path-or-module>`, `<path-or-module>`
- **Dependencies:** <previous slice, decision, credential, or none>
- **Verification:** `<command>` or <manual check>
- **Docs / artifacts:** <PRD/tasks/handoff/OpenSpec/design update>
- **Checkpoint:** `after-slice-<N>-<slug>` when needed

#### Checklist

- [ ] `[HITL]` <blocking human decision, if any>
- [ ] `[AFK]` <implementation step>
- [ ] `[AFK]` <test or verification step>
- [ ] `[AFK]` <artifact update step>
```

### PRD Split Template

When splitting a broad PRD, record the new sequence:

```markdown
## Decomposition

This work is split into ordered PRDs:

1. `PRD01_<first-slice>.md` — <first independently shippable result>
2. `PRD02_<second-slice>.md` — <next result, depends on PRD01>
3. `PRD03_<third-slice>.md` — <later result or optional extension>

The current PRD owns only item <N>.
```

### Issue Creation

External issues should map to slices, not arbitrary checklist fragments, unless the checklist item is a real blocker.

## Commands

```bash
openspec status --change "<change-id>" --json
```

## Resources

- Task breakdown: `planning-and-task-breakdown`
- PRD rules: `skillgrid-prd-artifacts`
- Issues: `skillgrid-issue-creation`
- Checkpoints: `skillgrid-checkpoints`
