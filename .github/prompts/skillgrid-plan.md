---
description: Plan — structured PRDs and proposals (OpenSpec + SDD)
allowed-tools: Read, Write, Glob, Grep, Bash, Task
argument-hint: "[feature name, slug, or PRD title]"
---

<objective>

You are executing **`/skillgrid-plan`** (PLAN phase) for the Skillgrid workflow.

</objective>

<process>

## Steps

1. **PRD (required format)** — Before implementation, produce or update a PRD that follows the structure below. Default locations: `prd/<slug>.md` for the change; add or update an entry in `prd/INDEX.md` when the repo uses that layout. If the project keeps PRDs under `docs/PRD/` or another path, follow existing conventions.
2. **OpenSpec change** — Create or update a change from the PRD using `openspec-propose` (OpenSpec layout and CLI) when the project uses OpenSpec.
3. **SDD proposal** — If the project uses SDD/Engram-style artifacts instead, create or update `proposal.md` per `sdd-propose` and the active persistence mode.
4. **Spec discipline** — Follow `spec-driven-development`: the PRD/plan is the single source of intent until superseded by design or delta specs.

## PRD document format

Use this outline for every PRD. Adapt headings only if the repo’s own PRD template says otherwise.

### Title block

- Heading: `### PRD: <Title>`
- **Spec / change:** `<path>` — canonical source for status and technical artifacts (e.g. `openspec/changes/<id>/` or project equivalent)
- **Status:** `Proposed` | `In progress` | `Done` (or project convention)

### Problem / why

What is wrong or missing, who is affected, and why it matters now.

### Goals

Bullet list of measurable or clearly verifiable outcomes.

### Assumptions (optional but recommended)

Surface assumptions the plan depends on; wrong assumptions should be corrected before design or implementation.

### In scope / out of scope

What this change includes and what is explicitly not included (prevents scope creep).

### User stories (optional)

Short “As a … I want … so that …” items when behavior is user-facing.

### Functional requirements

Numbered or bulleted **must-haves** for behavior, APIs, UX, and data. Each item should be testable.

### Non-functional requirements

Include as relevant: performance, security, privacy, accessibility, compatibility, observability, operational runbooks.

### Success criteria

How reviewers will know the work is done (acceptance-level checks, not a task list).

### Boundaries (agent / team guardrails)

Align with `spec-driven-development` where useful:

- **Always do** — e.g. tests before merge, naming, validation
- **Ask first** — e.g. schema, new deps, CI
- **Never do** — e.g. secrets in repo, silent requirement changes

### Project fit (when the change affects how work is done)

Concise notes on: **Commands** (real commands with flags), **structure** (paths for code, tests, docs), **code style** (one short illustrative pattern), **testing strategy** (levels and expectations). Skip subsections that are unchanged.

### Implementation tasks (optional in PRD; often refined in `/skillgrid-breakdown`)

If present: numbered workstreams and checkboxes; should trace to requirements above. Prefer linking to `tasks.md` for the full breakdown when it exists.

## Skills to read and follow

- `.agents/skills/karpathy-guidelines/SKILL.md` — assumptions, simplicity, surgical scope in the plan.
- `.agents/skills/openspec-propose/SKILL.md` — OpenSpec change + CLI-driven artifacts to apply-ready.
- `.agents/skills/sdd-propose/SKILL.md` — proposal when SDD orchestration or Engram/openspec file modes apply.
- `.agents/skills/spec-driven-development/SKILL.md` — objective, commands, structure, code style, testing, boundaries.

## Notes

- Inspect the repo with tools; do not assume stack or layout.
- If OpenSpec or SDD modes are unclear, ask once, then follow existing `openspec/` or repo persistence conventions.
- The repo’s narrative templates for `prd/INDEX.md` and per-PRD files also appear in `docs/wokflow.md` (PRD index and single-feature PRD examples).

</process>
