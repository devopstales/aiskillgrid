---
name: /skillgrid-design
id: skillgrid-design
category: Workflow
description: Design: UI, APIs, architecture, delta specs
---

You are executing **`/skillgrid-design`** (PLAN phase) for the Skillgrid workflow.

## Steps

1. **Plan UI** — User flows, key screens, states, and empty/error paths; use `frontend-design` for distinctive UI direction and `frontend-ui-engineering` for production-quality implementation guidance.
2. **Visual assets** — When the work needs generated or edited images (thumbnails, diagrams, etc.), follow `nano-banana` / team image workflow.
3. **DESIGN.md** — Create or update with UX and UI decisions (what and why).
4. **APIs and boundaries** — Use `api-and-interface-design` for contracts, error semantics, and module boundaries.
5. **Architecture doc** — Create or update **`.skillgrid/project/ARCHITECTURE.md`** (or team path) with components, data flow, and major dependencies.
6. **PRDs** — Create or update PRDs if scope or acceptance criteria changed during design.
7. **Proposals** — If intent or scope shifted, refresh the change via `openspec-propose` or `sdd-propose` (per project mode).
8. **Delta specs** — Write or update requirements and scenarios with `sdd-spec` (or equivalent OpenSpec delta spec steps).
9. **Data** — Apply `database-design` before migrations when schema or storage behavior is in scope.

## Skills to read and follow

- `.agents/skills/karpathy-guidelines/SKILL.md` — keep design decisions falsifiable and scoped.
- `.agents/skills/spec-driven-development/SKILL.md` — keep intent and acceptance criteria traceable while design evolves.
- `.agents/skills/sdd-spec/SKILL.md` — delta specs: requirements and scenarios.
- `.agents/skills/database-design/SKILL.md` — schema and data modeling before migrations.
- `.agents/skills/openspec-propose/SKILL.md` — refresh proposal when scope or intent changes.
- `.agents/skills/sdd-propose/SKILL.md` — same, under SDD artifact layout.
- `.agents/skills/frontend-design/SKILL.md` — distinctive, production-grade UI concepts.
- `.agents/skills/frontend-ui-engineering/SKILL.md` — components, design systems, accessibility (WCAG-oriented).
- `.agents/skills/api-and-interface-design/SKILL.md` — contract-first APIs, boundaries, error semantics.
- `.agents/skills/nano-banana/SKILL.md` — image generation/editing when visuals are in scope.

## Optional: IDE personas

For a second pass on **`DESIGN.md`** (flows, accessibility, API boundaries at design time—not code review), spawn **`skillgrid-design-critic`** ([`.cursor/agents/skillgrid-design-critic.md`](../../.cursor/agents/skillgrid-design-critic.md)).

## Notes

- Inspect the repo with tools; do not assume stack or layout.
- If OpenSpec or SDD modes are unclear, ask once, then follow existing `openspec/` or repo persistence conventions.
