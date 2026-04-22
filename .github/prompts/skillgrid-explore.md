---
description: Explore the problem and repo: OpenSpec explore, AGENTS/PROJECT, semantic search
---

You are executing **`/skillgrid-explore`** for the aiskillgrid workflow.

**Canonical checklist:** `docs/wokflow.md` (section `/skillgrid-explore`). Prefer that document if this prompt and the doc diverge.

## Actions

1. Enter explore stance: investigate and clarify; do not implement production code unless the user explicitly leaves explore mode.
2. Run OpenSpec-oriented exploration where applicable (read `openspec-explore` skill).
3. Create or refresh `AGENTS.md` and `PROJECT.md` when missing or stale.
4. Use `documentation-and-adrs` norms when documenting decisions during exploration.

## Skills to read and follow

Load each file below before doing substantive work (read fully or skim per skill length):

- `.agents/skills/openspec-explore/SKILL.md`
- `.agents/skills/search-first/SKILL.md`
- `.agents/skills/deep-research/SKILL.md`
- `.agents/skills/ccc/SKILL.md`
- `.agents/skills/documentation-and-adrs/SKILL.md`

## Notes

- Use tools to inspect the repo; do not assume stack or layout.
- If OpenSpec or SDD modes are unclear, ask once, then pick the path consistent with existing `openspec/` or project docs.
