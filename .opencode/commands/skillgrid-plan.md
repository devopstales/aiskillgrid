---
name: /skillgrid-plan
id: skillgrid-plan
category: Workflow
description: Plan: PRDs and proposals (OpenSpec + SDD)
---

You are executing **`/skillgrid-plan`** for the aiskillgrid workflow.

**Canonical checklist:** `docs/wokflow.md` (section `/skillgrid-plan`). Prefer that document if this prompt and the doc diverge.

## Actions

1. Produce or update a PRD covering objectives, commands, structure, code style, testing, and boundaries.
2. Create or update an OpenSpec change proposal (`openspec-propose`) or SDD `proposal.md` (`sdd-propose`) per project mode.

## Skills to read and follow

Load each file below before doing substantive work (read fully or skim per skill length):

- `.agents/skills/openspec-propose/SKILL.md`
- `.agents/skills/sdd-propose/SKILL.md`
- `.agents/skills/spec-driven-development/SKILL.md`

## Notes

- Use tools to inspect the repo; do not assume stack or layout.
- If OpenSpec or SDD modes are unclear, ask once, then pick the path consistent with existing `openspec/` or project docs.
