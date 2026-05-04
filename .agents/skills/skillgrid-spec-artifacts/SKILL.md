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
- ensure **one** `specs/<vertical-slice-slug>/spec.md` per named vertical slice (breakdown is incomplete if slices exist only in the PRD)
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
  -> openspec/changes/<change-id>/specs/<vertical-slice-slug>/spec.md   # one folder per vertical slice
  -> openspec/changes/<change-id>/tasks.md
  -> validate/apply/archive
```

**Per-slice layout (required for Skillgrid):** Each vertical slice gets `openspec/changes/<change-id>/specs/<slice-slug>/spec.md` — slice-scoped requirements, scenarios, and a checklist agents can load without the whole change. **`tasks.md`** is the cross-slice ordering and integration checklist. Per `docs/03-skillgrid-logic.md`, a shippable unit is **not** “tasks.md or slice spec”; it is **tasks.md + `specs/<slice>/spec.md`**. PRD slice write-ups must stay in lockstep and link to these paths; never stop at PRD-only slice bullets after `/skillgrid-breakdown`.

**Canonical blanks:** **`.skillgrid/templates/template-openspec-slice-spec.md`** and **`.skillgrid/templates/template-openspec-tasks.md`**. See **`docs/03-skillgrid-logic.md`**.

**Optional main-spec mirror:** `openspec/specs/<change-id>/spec.md` — umbrella requirements for the initiative when you want a stable top-level spec separate from delta layout. Not required for every change.

### Concept Mapping

Use this mapping when translating product intent into OpenSpec:

| Skillgrid / PRD concept | OpenSpec artifact | Ownership |
|---|---|---|
| Product problem, user goal, scope, success criteria | `proposal.md` | What changes and why |
| Technical approach, boundaries, tradeoffs | `design.md` | How the system changes |
| Observable behavior and edge cases | `specs/<slice>/spec.md` under the change (and optional `openspec/specs/<change-id>/spec.md`) | Verifiable requirements and scenarios per slice |
| Vertical-slice work checklist | `tasks.md` | Ordered implementation and verification |
| Current session state, blockers, evidence | `.skillgrid/tasks/context_<change-id>.md` | Rolling handoff and next action |

Do not let these drift. If one artifact changes the meaning of the work, refresh the others or record an explicit deferral.

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
- If the PRD, proposal, delta specs, and `tasks.md` disagree, stop and reconcile before `/skillgrid-apply`.
- External issues and Engram notes are indexes or summaries until their content is imported into PRD/OpenSpec artifacts.

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

### Breakdown completeness (before marking breakdown done)

- Every vertical slice named in the PRD implementation section has a corresponding directory `openspec/changes/<change-id>/specs/<vertical-slice-slug>/` with `spec.md` (use **`.skillgrid/templates/template-openspec-slice-spec.md`**).
- `tasks.md` lists the same slices (titles/slugs consistent with folder names); no orphan PRD-only slices.

### Validation Before Apply

Before `/skillgrid-apply` starts implementation, confirm:

- PRD exists and links the OpenSpec change.
- Required OpenSpec artifacts are present (including **all** slice `spec.md` files for slices in scope).
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

- Templates and planning logic: `docs/03-skillgrid-logic.md`
- Canonical blanks: `.skillgrid/templates/template-openspec-*.md`
- PRD rules: `skillgrid-prd-artifacts`
- Slice rules: `skillgrid-vertical-slices`
- OpenSpec config overlay: `skillgrid-openspec-config`
- Command sources: `.cursor/commands/skillgrid-plan.md`, `.cursor/commands/skillgrid-breakdown.md`, `.cursor/commands/skillgrid-validate.md`
