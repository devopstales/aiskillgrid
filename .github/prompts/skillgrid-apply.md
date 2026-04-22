---
description: Implement a minimal viable change using OpenSpec apply and disciplined coding
---

You are executing **`/skillgrid-apply`** for the aiskillgrid workflow.

**Canonical checklist:** `docs/wokflow.md` (section `/skillgrid-apply`). Prefer that document if this prompt and the doc diverge.

## Actions

1. Prefer git worktrees for parallel changes when useful.
2. Follow `openspec-apply-change` for OpenSpec-backed changes.
3. Aim for minimum viable scope; red-green-refactor where tests exist.
4. Cite official docs for framework choices (`source-driven-development`).
5. Refresh `ccc index` after substantial edits if semantic search is used.

## Skills to read and follow

Load each file below before doing substantive work (read fully or skim per skill length):

- `.agents/skills/openspec-apply-change/SKILL.md`
- `.agents/skills/karpathy-guidelines/SKILL.md`
- `.agents/skills/incremental-implementation/SKILL.md`
- `.agents/skills/tdd-workflow/SKILL.md`
- `.agents/skills/tdd-guide/SKILL.md`
- `.agents/skills/test-driven-development/SKILL.md`
- `.agents/skills/source-driven-development/SKILL.md`
- `.agents/skills/clean-code/SKILL.md`
- `.agents/skills/database-migrations/SKILL.md`
- `.agents/skills/ccc/SKILL.md`

## Notes

- Use tools to inspect the repo; do not assume stack or layout.
- If OpenSpec or SDD modes are unclear, ask once, then pick the path consistent with existing `openspec/` or project docs.
