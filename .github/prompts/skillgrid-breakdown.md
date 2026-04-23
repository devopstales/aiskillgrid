---
description: PRD implementation checklist synced with OpenSpec tasks.md
allowed-tools: Read, Write, Glob, Grep, Bash, Task
argument-hint: "[change-id or PRD slug]"
---

<objective>

You are executing **`/skillgrid-breakdown`** (TASKS phase) for the Skillgrid workflow.

</objective>

<process>

## Steps

1. **PRD — Implementation tasks chapter** — In the change’s PRD (`prd/<slug>.md` or project path from `/skillgrid-plan`), add or replace the **Implementation tasks** section using the checklist format below. Every checkbox item must **trace** to product intent (goals or requirements above in the same PRD).
2. **OpenSpec `tasks.md` (canonical)** — For the active change, produce or update `openspec/changes/<change-id>/tasks.md` using `sdd-tasks` from proposal, specs, and design. **Treat this file as the canonical implementation checklist** when the project uses OpenSpec.
3. **Keep PRD and `tasks.md` identical** — After any edit to the task list, **update both** the PRD Implementation tasks block and `openspec/changes/<change-id>/tasks.md` so the numbered workstreams and `- [ ]` lines match. Do not leave one file ahead of the other. If only one location exists (no PRD folder), use `tasks.md` only and note that in the PRD index or proposal when applicable.
4. **Link** — In the PRD, next to the Implementation tasks heading (or in the title block), include a **markdown link** to the canonical file, e.g. `[tasks.md](openspec/changes/<change-id>/tasks.md)` (adjust path from repo root). The linked file and the in-PRD checklist must stay the same.
5. **Delta specs** — Confirm requirements and scenarios are complete enough to task out (`sdd-spec` / OpenSpec delta); close gaps before finalizing tasks.
6. **Testing stance** — For tasks that need automated proof, align test levels with `test-driven-development` and `tdd-guide` (what to test first, pyramid, acceptance tests). Testing workstreams often appear as their own numbered group in the checklist (see example below).
7. **Framework tasks** — When tasks touch framework behavior, ground them with `source-driven-development` (doc links or citations in task line text where helpful).
8. **Task quality** — Use `planning-and-task-breakdown`: small, verifiable items with clear acceptance, dependency ordering, parallel vs sequential work called out in workstream titles or notes.
9. **Index** — After large breakdown edits to code or structure, plan a `ccc` index refresh if the team uses CocoIndex.
10. **No orphan work** — Each task maps to a spec or agreed design; call out gaps before coding.

## PRD: Implementation tasks checklist format

Mirror the structure used in human-friendly PRDs with OpenSpec (example: [clustered-replica-mode PRD](https://github.com/devopstales/KubeDash/blob/feature/4.1-enhanced-log-terminal/project/prd/clustered-replica-mode.md)):

- Optional horizontal rule `---` before the section.
- Section title, e.g. `### Implementation tasks` or `### Implementation tasks (from OpenSpec)`.
- **Workstreams** as `#### <n>. <Workstream title>` (Redis backend, Leader election, Testing, Docs, etc.).
- **Sub-tasks** as github-style checklist items with **global numbering** across the PRD: `- [ ] 1.1 ...`, `- [ ] 1.2 ...`, then `#### 2. ...` with `- [ ] 2.1 ...`, and so on.
- Keep phrasing **action-oriented** and **verifiable** (specific files, behaviors, or artifacts).

Apply the **same** structure and **same** numbering inside `openspec/changes/<change-id>/tasks.md` when using OpenSpec.

**Minimal pattern:**

```markdown
---

### Implementation tasks

**Canonical checklist:** [tasks.md](openspec/changes/<change-id>/tasks.md) — keep this section in sync with that file.

#### 1. <First workstream>

- [ ] 1.1 …
- [ ] 1.2 …

#### 2. <Next workstream>

- [ ] 2.1 …
```

## Skills to read and follow

- `.agents/skills/karpathy-guidelines/SKILL.md` — verifiable tasks, no hidden scope.
- `.agents/skills/sdd-spec/SKILL.md` — delta specs before task lock-in.
- `.agents/skills/sdd-tasks/SKILL.md` — `tasks.md` from proposal, specs, and design.
- `.agents/skills/openspec-propose/SKILL.md` — OpenSpec change layout and artifacts under `openspec/changes/`.
- `.agents/skills/planning-and-task-breakdown/SKILL.md` — atomic tasks, acceptance criteria, dependencies.
- `.agents/skills/test-driven-development/SKILL.md` — pyramid, DAMP, browser testing where relevant.
- `.agents/skills/tdd-guide/SKILL.md` — TDD patterns for implementation tasks.
- `.agents/skills/source-driven-development/SKILL.md` — cite official docs for framework-dependent tasks.
- `.agents/skills/ccc/SKILL.md` — refresh semantic index after significant code or layout changes.

## Optional: IDE personas

After you draft or sync **`tasks.md` and the PRD checklist**, you can spawn **`skillgrid-task-breakdown-auditor`** ([`.cursor/agents/skillgrid-task-breakdown-auditor.md`](../../.cursor/agents/skillgrid-task-breakdown-auditor.md)) for a planning-only audit (AC, ordering, testability) without a full code read.

## Notes

- Inspect the repo with tools; do not assume stack, change id, or layout.
- If OpenSpec or SDD modes are unclear, ask once, then follow existing `openspec/` or repo persistence conventions.
- For PRD document sections outside Implementation tasks, follow `/skillgrid-plan` (problem, goals, requirements, success criteria, boundaries).

</process>
