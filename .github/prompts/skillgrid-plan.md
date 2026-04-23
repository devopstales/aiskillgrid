---
description: Plan: PRDs and proposals (OpenSpec + SDD)
---

You are executing **`/skillgrid-plan`** (PLAN phase) for the Skillgrid workflow.

## Steps

1. **PRDs** — Produce or update PRD-style documents covering objectives, commands, structure, code style, testing, and boundaries *before* implementation work.
2. **OpenSpec change** — Create or update a change from the PRD using `openspec-propose` (OpenSpec layout and CLI) when the project uses OpenSpec.
3. **SDD proposal** — If the project uses SDD/Engram-style artifacts instead, create or update `proposal.md` per `sdd-propose` and the active persistence mode.
4. **Spec discipline** — Follow `spec-driven-development`: the PRD/plan is the single source of intent until superseded by design or delta specs.

## Skills to read and follow

- `.agents/skills/openspec-propose/SKILL.md` — OpenSpec change + CLI-driven artifacts to apply-ready.
- `.agents/skills/sdd-propose/SKILL.md` — proposal when SDD orchestration or Engram/openspec file modes apply.
- `.agents/skills/spec-driven-development/SKILL.md` — PRD-style scope before code.

## Notes

- Inspect the repo with tools; do not assume stack or layout.
- If OpenSpec or SDD modes are unclear, ask once, then follow existing `openspec/` or repo persistence conventions.
