---
description: Implement a minimal viable change using OpenSpec apply and disciplined coding
---

You are executing **`/skillgrid-apply`** (BUILD phase) for the Skillgrid workflow.

## Steps

1. **Worktrees** — Prefer git worktrees for parallel changes when the workflow benefits from isolation.
2. **Apply change** — Implement using `openspec-apply-change` from the change’s tasks and delta specs (when OpenSpec/SDD artifacts drive the work).
3. **Scope** — Ship a **minimum viable change**: smallest increment that meets the current tasks; avoid unrelated refactors.
4. **TDD** — When tests are part of the stack, use red–green–refactor (`tdd-workflow`, `tdd-guide`, `test-driven-development`).
5. **Framework choices** — Ground decisions in official documentation (`source-driven-development`); note citations in commits or code where helpful.
6. **Quality** — Keep edits readable and cohesive (`clean-code`, `karpathy-guidelines`); use `incremental-implementation` for vertical slices with verification and atomic commits.
7. **Migrations** — Use `database-migrations` for schema and data changes safely.
8. **Index** — After substantial edits, refresh CocoIndex: `ccc index` when the project relies on `ccc` search.

## Skills to read and follow

- `.agents/skills/openspec-apply-change/SKILL.md` — implement from OpenSpec change tasks.
- `.agents/skills/karpathy-guidelines/SKILL.md` — surgical, verifiable steps.
- `.agents/skills/incremental-implementation/SKILL.md` — vertical slices, verify, commit; safe defaults.
- `.agents/skills/tdd-workflow/SKILL.md` — structured red–green–refactor discipline.
- `.agents/skills/tdd-guide/SKILL.md` — TDD patterns.
- `.agents/skills/test-driven-development/SKILL.md` — red–green–refactor, pyramid, DAMP, Beyoncé rule.
- `.agents/skills/source-driven-development/SKILL.md` — cite official docs for framework decisions.
- `.agents/skills/clean-code/SKILL.md` — readability and maintainability while implementing.
- `.agents/skills/database-migrations/SKILL.md` — apply schema/data migrations safely.
- `.agents/skills/ccc/SKILL.md` — refresh semantic index after significant code changes.

## Notes

- Inspect the repo with tools; do not assume stack or layout.
- If OpenSpec or SDD modes are unclear, ask once, then follow existing `openspec/` or repo persistence conventions.
