---
description: Break work into ordered, verifiable tasks
---

You are executing **`/skillgrid-breakdown`** (PLAN phase) for the Skillgrid workflow.

## Steps

1. **PRDs** — Add or refine implementation tasks in PRDs so every task trace back to product intent.
2. **OpenSpec / SDD tasks** — For the active change, produce `tasks.md` (or equivalent) using `sdd-tasks` from proposal, specs, and design.
3. **Task quality** — Use `planning-and-task-breakdown`: small, verifiable items with clear acceptance criteria and dependency ordering; flag parallel vs sequential work.
4. **No orphan work** — Each task should map to a spec or agreed design; call out gaps before coding.

## Skills to read and follow

- `.agents/skills/sdd-tasks/SKILL.md` — tasks.md checklist from proposal, specs, and design.
- `.agents/skills/planning-and-task-breakdown/SKILL.md` — atomic tasks, acceptance criteria, dependencies.

## Notes

- Inspect the repo with tools; do not assume stack or layout.
- If OpenSpec or SDD modes are unclear, ask once, then follow existing `openspec/` or repo persistence conventions.
