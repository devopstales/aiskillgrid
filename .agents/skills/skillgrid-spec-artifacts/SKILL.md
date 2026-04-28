---
name: skillgrid-spec-artifacts
description: >
  Bridges Skillgrid PRD intent to OpenSpec proposal, design, delta specs, tasks, and validation artifacts.
  Trigger: Creating or updating OpenSpec artifacts from a Skillgrid PRD or checking PRD/spec alignment.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

## When to Use

Use this skill when a Skillgrid phase needs to scaffold or refresh `openspec/changes/<change-id>/`, derive a change id from a PRD, align OpenSpec artifacts with product intent, or verify that `tasks.md` traces back to the PRD and delta specs.

## Critical Patterns

### Responsibilities

This skill owns the Skillgrid bridge into OpenSpec:

- derive or confirm `change-id`
- create or locate `openspec/changes/<change-id>/`
- align `proposal.md` with PRD problem, scope, and goals
- align `design.md` with implementation approach and UI constraints
- align delta specs with measurable requirements and scenarios
- align `tasks.md` with PRD implementation tasks
- run or request OpenSpec validation before implementation or archival

Do not duplicate the full OpenSpec phase skills. Reuse:

- `openspec-propose`
- `openspec-explore`
- `openspec-apply-change`
- `openspec-verify-change`
- `openspec-sync-specs`
- `openspec-archive-change`

### Artifact Flow

```text
PRD intent
  -> openspec/changes/<change-id>/proposal.md
  -> openspec/changes/<change-id>/design.md
  -> openspec/changes/<change-id>/specs/**/*
  -> openspec/changes/<change-id>/tasks.md
  -> validate/apply/archive
```

### Change Id

- Prefer the existing change id if the PRD already links one.
- Otherwise derive a kebab-case id from the PRD title or user request.
- Do not silently create a second change for the same PRD.
- If two artifacts disagree on the change id, stop and reconcile links before generating more content.

### OpenSpec CLI Loop

When the project uses OpenSpec on disk:

```bash
openspec list --json
openspec status --change "<change-id>" --json
openspec instructions <artifact-id> --change "<change-id>" --json
openspec validate --change "<change-id>"
```

Use `openspec instructions` output as constraints. Do not paste raw `context`, `rules`, or schema blocks into generated artifact files.

### Alignment Rules

- `proposal.md` says what changes and why.
- `design.md` says how the system changes and records tradeoffs.
- Delta specs define requirements and scenarios that can be verified.
- `tasks.md` is the implementation checklist and should be compatible with `skillgrid-vertical-slices`.
- PRD implementation tasks and OpenSpec `tasks.md` should not diverge; if one changes materially, refresh the other.

### Tracer-Bullet Breakdown

When converting a PRD, proposal, or design into `tasks.md`, draft independently grabbable vertical slices before writing a flat checklist.

Work from existing conversation context and artifacts first. If the source is a remote issue, fetch it with comments through the configured provider workflow, then convert it into Skillgrid artifacts instead of treating the remote issue as the source of truth.

Each slice should be a tracer bullet:

- a narrow but complete path through every layer needed for the behavior
- demoable or verifiable on its own
- thin enough to review independently
- tagged `[AFK]` when an agent can implement it from existing artifacts
- tagged `[HITL]` only when it truly needs human input, credentials, approval, architecture choice, or design review

Avoid horizontal slices such as "build all backend", "build all frontend", or "write all tests". Prefer many thin slices over a few thick ones.

Before creating external issues, present the proposed slice breakdown for approval:

- **Title:** short user-visible capability
- **Type:** `[HITL]` or `[AFK]`
- **Blocked by:** dependency slice, decision, credential, or `None`
- **User stories covered:** PRD story numbers when available
- **Acceptance criteria:** observable checks for this slice

Ask whether the granularity, dependencies, and HITL/AFK tags are correct. Iterate until approved. Then, if `.skillgrid/config.json` enables a remote tracker and the user wants tracker-backed work, hand off to `skillgrid-issue-creation`. Create/link issues in dependency order so blocker references are real.

### Proposal Template

Use when OpenSpec instructions do not provide a stronger template:

```markdown
# Proposal: <change-id>

## Summary

<One paragraph: what changes and why.>

## Problem

<Current behavior or missing capability from the PRD.>

## Goals

- <Goal copied or refined from PRD>

## Non-goals

- <Explicitly excluded behavior>

## Scope

### In

- <Included behavior>

### Out

- <Excluded behavior>

## User impact

<Who experiences the change and what they can do afterward.>

## Related artifacts

- PRD: `.skillgrid/prd/PRD<NN>_<slug>.md`
- Handoff: `.skillgrid/tasks/context_<change-id>.md`
```

### Design Template

```markdown
# Design: <change-id>

## Context

<Relevant codebase, product, and technical context.>

## Approach

<Chosen implementation approach.>

## Boundaries

- <Module/service/UI boundary>

## Data and contracts

<API, schema, event, file, or persistence contracts that change.>

## UI / UX constraints

<Design constraints, states, accessibility, responsive behavior, and links to `DESIGN.md` or previews.>

## Testing strategy

<Unit, integration, E2E, browser, manual, or security checks.>

## Risks and mitigations

| Risk | Mitigation |
|---|---|
| <risk> | <mitigation> |

## Alternatives considered

- <Alternative> — <why rejected>
```

### Delta Spec Template

Use OpenSpec's required capability path and heading conventions when available. A generic fallback:

```markdown
# <Capability Name>

## ADDED Requirements

### Requirement: <Requirement name>

The system SHALL <observable behavior>.

#### Scenario: <Scenario name>

- **GIVEN** <initial state>
- **WHEN** <actor action or system event>
- **THEN** <expected observable result>

## MODIFIED Requirements

### Requirement: <Existing requirement name>

<Updated requirement text or scenario.>

## REMOVED Requirements

### Requirement: <Deprecated requirement name>

<Reason for removal and migration note if relevant.>
```

### Tasks Template

```markdown
# Tasks: <change-id>

## Execution note

Implement task-by-task. For behavioral code, use Red-Green-Refactor: write one failing test, verify it fails for the expected reason, implement the minimum code, verify green, then refactor while green.

## Slice breakdown

1. **<Slice title>**
   - **Type:** `[AFK]`
   - **Blocked by:** None
   - **User stories covered:** <PRD story numbers or "not specified">
   - **Acceptance criteria:**
     - [ ] <Observable result>

## Implementation

- [ ] `[HITL]` <Decision, approval, credential, or manual setup needed before dependent autonomous work>
- [ ] `[AFK]` Write one failing behavior test for <slice behavior>; expected failure: <why it should fail now>
- [ ] `[AFK]` Implement the minimum code for <slice behavior>
- [ ] `[AFK]` Verify the focused test passes, then run the relevant wider check
- [ ] `[AFK]` Refactor only after green, keeping tests green

## Verification

- [ ] <Command or manual check that proves the slice works>

## Documentation

- [ ] <PRD, handoff, project doc, or design doc update>
```

### Validation Before Apply

Before `/skillgrid-apply` starts implementation, confirm:

- PRD exists and links the OpenSpec change.
- Required OpenSpec artifacts are present.
- `openspec status` says apply-required artifacts are complete.
- `tasks.md` contains verifiable, ordered work.
- `tasks.md` uses tracer-bullet vertical slices rather than horizontal layer-only tasks.
- Slice granularity, dependencies, and HITL/AFK tags have been approved or explicitly recorded as assumptions.
- Human-blocked tasks are tagged `[HITL]`; autonomous tasks are tagged `[AFK]` when the repo uses those tags.
- Behavioral tasks name the first failing test or explain why TDD is not applicable.
- Tasks do not contain placeholders such as `TBD`, "add appropriate handling", or "write tests" without concrete behavior and verification.

## Commands

```bash
openspec new change "<change-id>"
openspec status --change "<change-id>" --json
openspec instructions tasks --change "<change-id>" --json
openspec validate --change "<change-id>"
```

## Resources

- PRD rules: `skillgrid-prd-artifacts`
- Slice rules: `skillgrid-vertical-slices`
- OpenSpec config overlay: `skillgrid-openspec-config`
- Command sources: `.cursor/commands/skillgrid-plan.md`, `.cursor/commands/skillgrid-breakdown.md`, `.cursor/commands/skillgrid-validate.md`
