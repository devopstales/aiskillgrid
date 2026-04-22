---
name: /skillgrid-finish
id: skillgrid-finish
category: Workflow
description: Close the change: archive or sync specs, git hygiene, ship checklist, PR
---

You are executing **`/skillgrid-finish`** for the aiskillgrid workflow.

**Canonical checklist:** `docs/wokflow.md` (section `/skillgrid-finish`). Prefer that document if this prompt and the doc diverge.

## Actions

1. Archive OpenSpec changes with `openspec-archive-change` when done; use `openspec-sync-specs` if promoting specs without archiving.
2. Follow trunk-friendly git workflow; prepare PR with clear description and risk notes.
3. Use shipping checklist skills when deploying.

## Skills to read and follow

Load each file below before doing substantive work (read fully or skim per skill length):

- `.agents/skills/openspec-archive-change/SKILL.md`
- `.agents/skills/openspec-sync-specs/SKILL.md`
- `.agents/skills/git-workflow-and-versioning/SKILL.md`
- `.agents/skills/ci-cd-and-automation/SKILL.md`
- `.agents/skills/deprecation-and-migration/SKILL.md`
- `.agents/skills/shipping-and-launch/SKILL.md`

## Notes

- Use tools to inspect the repo; do not assume stack or layout.
- If OpenSpec or SDD modes are unclear, ask once, then pick the path consistent with existing `openspec/` or project docs.
