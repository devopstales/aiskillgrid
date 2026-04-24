---
name: /skillgrid-apply
id: skillgrid-apply
category: Workflow
description: Implement from tasks; keep PRD and OpenSpec tasks.md checked off in sync
allowed-tools: Read, Write, Edit, Glob, Grep, Bash, Task
argument-hint: "[change-id; optional task focus]"
---

<objective>

You are executing **`/skillgrid-apply`** (BUILD phase) for the Skillgrid workflow.

</objective>

<process>

## Steps

1. **Worktrees** — Prefer git worktrees for parallel changes when the workflow benefits from isolation.
2. **Apply change** — Implement using `openspec-apply-change` from the change’s tasks and delta specs (when OpenSpec/SDD artifacts drive the work).
3. **Task bookkeeping (mandatory)** — When you **complete** a task, **immediately** update the checklist in **both** places so they never diverge:
   - `openspec/changes/<change-id>/tasks.md` — set the matching items to `- [x]` (or adjust wording if the task was split, deferred, or corrected; then mirror that edit in the PRD).
   - The change’s **PRD** — the **Implementation tasks** section (e.g. `prd/<slug>.md`) must show the **same** `- [x]` / `- [ ]` state and the same numbered lines as `tasks.md`.
   If the project has no PRD file for the change, update `tasks.md` only and add or restore a PRD pointer when the team expects one. Never mark done in one file and leave the other stale.
4. **Scope** — Ship a **minimum viable change**: smallest increment that meets the current tasks; avoid unrelated refactors.
5. **Contracts** — Preserve public APIs and error semantics agreed in design; use `api-and-interface-design` when changing boundaries.
6. **TDD** — When tests are part of the stack, use red–green–refactor (`tdd-workflow`, `tdd-guide`, `test-driven-development`).
7. **Framework choices** — Ground decisions in official documentation (`source-driven-development`); note citations in commits or code where helpful.
8. **Quality** — Keep edits readable and cohesive (`clean-code`, `karpathy-guidelines`); use `incremental-implementation` for vertical slices with verification and atomic commits.
9. **Migrations** — Use `database-migrations` for schema and data changes safely.
10. **Graph** — After substantial edits, run **`graphify update .`** when the project uses graphify (see `AGENTS.md`).

## Skills to read and follow

- `.agents/skills/openspec-apply-change/SKILL.md` — implement from OpenSpec change tasks.
- `.agents/skills/api-and-interface-design/SKILL.md` — contracts, Hyrum’s Law, error semantics at module edges.
- `.agents/skills/karpathy-guidelines/SKILL.md` — surgical, verifiable steps.
- `.agents/skills/incremental-implementation/SKILL.md` — vertical slices, verify, commit; safe defaults.
- `.agents/skills/tdd-workflow/SKILL.md` — structured red–green–refactor discipline.
- `.agents/skills/tdd-guide/SKILL.md` — TDD patterns.
- `.agents/skills/test-driven-development/SKILL.md` — red–green–refactor, pyramid, DAMP, Beyoncé rule.
- `.agents/skills/source-driven-development/SKILL.md` — cite official docs for framework decisions.
- `.agents/skills/clean-code/SKILL.md` — readability and maintainability while implementing.
- `.agents/skills/database-migrations/SKILL.md` — apply schema/data migrations safely.
- `.agents/skills/references/indexing-and-memory.md` — graphify refresh and structural search after significant code changes.

## Notes

- Inspect the repo with tools; do not assume stack or layout.
- If OpenSpec or SDD modes are unclear, ask once, then follow existing `openspec/` or repo persistence conventions.
- Checklist format and the PRD↔`tasks.md` link live in `/skillgrid-breakdown`; apply keeps them aligned as work lands.

</process>
