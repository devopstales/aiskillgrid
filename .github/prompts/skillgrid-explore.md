---
description: Explore the problem and repo: OpenSpec list, PRD backfill, OpenSpec explore, .skillgrid/project, AGENTS, graphify
allowed-tools: Read, Write, Glob, Grep, Bash, Task
argument-hint: "[optional: topic, change id, or area to explore]"
---

<objective>

You are executing **`/skillgrid-explore`** (DEFINE phase) for the Skillgrid workflow.

</objective>

<process>

## Steps

1. **Stance** — Explore and clarify; do not implement production code unless the user explicitly leaves explore mode.
2. **OpenSpec inventory** — If the project uses OpenSpec and the CLI is on PATH, run **`openspec list`** from the repository root. If the CLI is missing or the command fails, still inspect **`openspec/changes/`** and **`openspec/specs/`** (if present) so you have a list of **active changes** and **main specs** to explore. Summarize what exists for the user.
3. **PRD coverage for existing changes** — For each change surfaced by `openspec list` (or each directory under `openspec/changes/<change-id>/` when you inventoried manually), confirm there is a **PRD** that names that change. Typical locations (use what the repo already has; do not invent a second tree): `prd/<slug>.md`, `.skillgrid/prd/<slug>.md`, `docs/PRD/<slug>.md`, and entries in **`prd/INDEX.md`** or **`.skillgrid/prd/INDEX.md`**. The PRD’s title block or body should point at the change path (e.g. `openspec/changes/<id>/`). **If a change has no PRD**, create one using the **PRD template** and **Implementation tasks** pattern below (same as [`/skillgrid-plan`](./skillgrid-plan.md) and [`/skillgrid-breakdown`](./skillgrid-breakdown.md)). Fill gaps from the change’s `proposal.md`, delta specs, or code until the PRD is reviewable. Update **`prd/INDEX.md`** (or the repo’s index) when the team uses it.
4. **OpenSpec explore** — Follow `openspec-explore`: investigate the problem space and relevant code paths before locking a change. Use the inventory and PRD map from steps 2–3 to prioritize files.
5. **Project docs (canonical)** — Create or refresh **`.skillgrid/project/ARCHITECTURE.md`**, **`STRUCTURE.md`**, and **`PROJECT.md`** (system design, repo layout and optional runtime topology, onboarding narrative). Use templates in [`docs/wokflow.md`](../../docs/wokflow.md) when helpful.
6. **AGENTS.md** — Create or refresh at repo root so agent behavior and project rules are current.
7. **Documentation** — When recording exploration outcomes, document the *why* (ADRs, API docs, inline standards) per team norms.
8. **Code discovery** — Use **`graphify-out/`** and **`AGENTS.md`** for orientation, then **`rg`/IDE search** and targeted file reads.
9. **Optional depth** — Use `deep-research` when the question needs external evidence or broader comparison.

## PRD template (use when a change has no PRD)

Adapt headings if the repo’s own template overrides. Default file locations match **`/skillgrid-plan`**: e.g. `prd/<slug>.md` and **`prd/INDEX.md`**.

### PRD document format (from `/skillgrid-plan`)

#### Title block

- Heading: `### PRD: <Title>`
- **Spec / change:** `<path>` — canonical source for status and technical artifacts (e.g. `openspec/changes/<id>/` or project equivalent)
- **Status:** `Proposed` | `In progress` | `Done` (or project convention)

#### Problem / why

What is wrong or missing, who is affected, and why it matters now.

#### Goals

Bullet list of measurable or clearly verifiable outcomes.

#### Assumptions (optional but recommended)

Surface assumptions the plan depends on; wrong assumptions should be corrected before design or implementation.

#### In scope / out of scope

What this change includes and what is explicitly not included (prevents scope creep).

#### User stories (optional)

Short “As a … I want … so that …” items when behavior is user-facing.

#### Functional requirements

Numbered or bulleted **must-haves** for behavior, APIs, UX, and data. Each item should be testable.

#### Non-functional requirements

Include as relevant: performance, security, privacy, accessibility, compatibility, observability, operational runbooks.

#### Success criteria

How reviewers will know the work is done (acceptance-level checks, not a task list).

#### Boundaries (agent / team guardrails)

- **Always do** — e.g. tests before merge, naming, validation
- **Ask first** — e.g. schema, new deps, CI
- **Never do** — e.g. secrets in repo, silent requirement changes

#### Project fit (when the change affects how work is done)

Concise notes on: **Commands** (real commands with flags), **structure** (paths for code, tests, docs), **code style** (one short illustrative pattern), **testing strategy** (levels and expectations). Skip subsections that are unchanged.

#### Implementation tasks (from `/skillgrid-breakdown`)

Add or update using the checklist format below. Every checkbox item must **trace** to goals or requirements above. **Keep PRD and `openspec/changes/<change-id>/tasks.md` identical** when that file exists (same numbering and `- [ ]` lines). Link: `[tasks.md](openspec/changes/<change-id>/tasks.md)`.

- Optional `---` before the section.
- Section title, e.g. `### Implementation tasks` or `### Implementation tasks (from OpenSpec)`.
- **Workstreams** as `#### <n>. <Workstream title>`.
- **Sub-tasks** with **global numbering**: `- [ ] 1.1 ...`, then `#### 2. ...` with `- [ ] 2.1 ...`, and so on.
- **Minimal pattern:**

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

If `tasks.md` does not exist yet, still include an **Implementation tasks** section in the new PRD with a reasonable draft checklist; the user can run **`/skillgrid-breakdown`** to sync to OpenSpec.

## Skills to read and follow

- `.agents/skills/openspec-explore/SKILL.md`
- `.agents/skills/openspec-propose/SKILL.md` — OpenSpec changes and `openspec/changes/<id>/` layout when backfilling PRDs.
- `.agents/skills/search-first/SKILL.md`
- `.agents/skills/deep-research/SKILL.md`
- `.agents/skills/references/indexing-and-memory.md`
- `.agents/skills/documentation-and-adrs/SKILL.md`
- `.agents/skills/spec-driven-development/SKILL.md` — when drafting PRD sections and boundaries.

## Optional: IDE personas

When spawning a **subagent** for exploration-only work in a clean context, use **`skillgrid-explore-architect`** ([`.cursor/agents/skillgrid-explore-architect.md`](../../.cursor/agents/skillgrid-explore-architect.md)).

For **external / cited research** (landscape, competitors, docs via Exa, Firecrawl, DeepWiki, Context7) rather than in-repo mapping, use **`skillgrid-researcher`** ([`.cursor/agents/skillgrid-researcher.md`](../../.cursor/agents/skillgrid-researcher.md)).

## Notes

- Inspect the repo with tools; do not assume stack or layout.
- If OpenSpec or SDD modes are unclear, ask once, then align with existing `openspec/` or repo conventions.

</process>
