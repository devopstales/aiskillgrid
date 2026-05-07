# SDD Enforcement Contract

Canonical runtime enforcement source for `sdd-orchestrator` and `sdd-*` phases.

## 1) Phase Routing And Stop Conditions

- Use explicit phase transitions; do not silently skip phases.
- `HITL` unresolved decision => `status: blocked`.
- `AFK` continuation only when scope, acceptance criteria, and checks are explicit.
- Missing mandatory artifact or gate output => fail closed with `status: failed`.

## 2) Mandatory Skill-Gate Matrix

### `sdd-brainstorm`
- Require planning chain artifacts:
  - explore, clarify, propose, spec, design, prd, tasks
- If `ui_scope: true`, also require `ui-wireframes.md` and `ui-decisions.md`.

### `sdd-apply`
- Label validator must pass.
- Require full `spec`, `design`, `tasks`.
- Require explicit slice acceptance criteria before coding.

### `sdd-verify`
- Label validator must pass.
- Require full `proposal`, `spec`, `design`, `tasks`.
- Require implementation evidence (tests/build/check output or equivalent).

Gate handling:
- Any gate failure must fail closed and include failing gate id + missing evidence.

## 3) Two-Stage Review Gate

All decision-ready reports must include:

1. Stage 1: spec compliance
2. Stage 2: code quality

Any critical finding in either stage blocks progression.

## 4) Standard Return Envelope

Use `skills/_shared/sdd-return-envelope.md` as canonical envelope contract.
