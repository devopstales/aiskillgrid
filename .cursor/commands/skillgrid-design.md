---
name: /skillgrid-design
id: skillgrid-design
category: Workflow
description: Design: UI, APIs, architecture, delta specs
---

You are executing **`/skillgrid-design`** for the aiskillgrid workflow.

**Canonical checklist:** `docs/wokflow.md` (section `/skillgrid-design`). Prefer that document if this prompt and the doc diverge.

## Actions

1. Produce or update `DESIGN.md` and `ARCHITECTURE.md` as appropriate.
2. Apply UI and API skills for user-facing and boundary design.
3. Write or update delta specs with `sdd-spec` when using SDD/OpenSpec artifacts.
4. If scope shifts, refresh proposal via `openspec-propose` / `sdd-propose`.

## Skills to read and follow

Load each file below before doing substantive work (read fully or skim per skill length):

- `.agents/skills/sdd-spec/SKILL.md`
- `.agents/skills/database-design/SKILL.md`
- `.agents/skills/openspec-propose/SKILL.md`
- `.agents/skills/sdd-propose/SKILL.md`
- `.agents/skills/frontend-ui-engineering/SKILL.md`
- `.agents/skills/api-and-interface-design/SKILL.md`

## Notes

- Use tools to inspect the repo; do not assume stack or layout.
- If OpenSpec or SDD modes are unclear, ask once, then pick the path consistent with existing `openspec/` or project docs.
