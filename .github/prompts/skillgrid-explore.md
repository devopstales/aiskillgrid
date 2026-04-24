---
name: /skillgrid-explore
id: skillgrid-explore
category: Workflow
description: Explore the problem and repo: OpenSpec list, PRD backfill, .skillgrid/project, AGENTS, graphify
allowed-tools: Read, Write, Glob, Grep, Bash, Task
argument-hint: "[optional: topic, change id, or area to explore]"
---

<objective>

You are executing **`/skillgrid-explore`** (DEFINE phase) for the Skillgrid workflow.

</objective>

<process>

## Stance: explore, do not implement

**Explore mode is for thinking, not implementing.** You may read files, search code, and investigate the codebase, but you must **not** write application code or ship features. If the user asks you to implement something, remind them to leave explore mode and use **`/skillgrid-plan`** (or a later build phase). You **may** create or edit planning artifacts (PRDs, OpenSpec `proposal.md` / `design.md` text under `openspec/changes/`, ADRs) when that is capturing thinking—not production code.

**This is a stance, not a fixed script** — there is no mandatory sequence. Be curious, open-ended, patient, and grounded in the real repo. Use diagrams when they help. Do not auto-capture: offer to save insights, then let the user decide.

---

## Steps

1. **Hybrid persistence** — Use **`openspec list` / on-disk `openspec/`** as the file inventory, **and** offer to persist important decisions to **Engram** with a stable `topic_key` (see **`/skillgrid-init`**) so context survives compactions. If `openspec/` is still missing, align with init: bootstrap the tree, then run the steps below.

2. **OpenSpec inventory** — If the project uses OpenSpec and the CLI is on `PATH`, from the repository root run:

   ```bash
   openspec list --json
   ```

   This surfaces active changes, names, and status. If the CLI is missing or the command fails, still inspect **`openspec/changes/`** and **`openspec/specs/`** (if present) and summarize what exists. If the project uses `openspec list` without `--json`, that is acceptable; prefer `--json` when available for machine-readable status.

3. **PRD coverage for existing changes** — For each change surfaced by `openspec list` (or each directory under `openspec/changes/<change-id>/` when you inventoried manually), confirm there is a **PRD** that names that change. **Canonical path:** **`.skillgrid/prd/PRD<NN>_<slug>.md`** (two-digit **`<NN>`** = execution order). **First**, glob **`.skillgrid/prd/PRD*.md`** (and root `prd/PRD*.md` only if migrating legacy), **sort by `<NN>`**, and read **`.skillgrid/prd/INDEX.md`**. The PRD’s title block or body should point at the change path (e.g. `openspec/changes/<id>/`). **If execution order of open work is wrong**, renumber files (`PRD01_…` …) and update **`.skillgrid/prd/INDEX.md`**, cross-links, and any Engram or OpenSpec notes (same rules as **`/skillgrid-plan`**). **If a change has no PRD**, create one under **`.skillgrid/prd/`** using the **PRD template** below and assign the next appropriate `<NN>`. **Do not** create new PRD files at repo root `prd/`. `docs/PRD/` may mirror; keep numbering consistent with **`.skillgrid/prd/`**.

4. **Read existing change artifacts in context** — If the user names a change or one is relevant:

   - Read `openspec/changes/<name>/proposal.md`, `design.md`, `tasks.md`, and delta specs as needed.
   - Reference them in conversation; offer to capture new decisions in the right file when the user wants that recorded.

5. **Project docs (canonical)** — Create or refresh **`.skillgrid/project/ARCHITECTURE.md`**, **`STRUCTURE.md`**, and **`PROJECT.md`**. Use **Project document templates** in **`/skillgrid-init`** when helpful.

6. **AGENTS.md** — Create or refresh at repo root so agent behavior and project rules are current.

7. **Documentation** — When recording exploration outcomes, document the *why* (ADRs, API docs, inline standards) per team norms.

8. **Code discovery** — Use **`graphify-out/`** and **`AGENTS.md`** for orientation, then **`rg` / IDE search** and targeted file reads. Optional: deeper external research when the question needs off-repo evidence (document sources).

9. **Ending** — There is no required ending. You may offer a short summary or suggest moving to **`/skillgrid-plan`** when the idea is ready to formalize.

### What you might do (from conversation)

- Explore the problem space, compare options, surface risks, sketch ASCII architecture.
- Map integration points and hidden complexity in the codebase.
- **When a change exists:** map insights to the table below and offer to update files—only if the user agrees.

| Insight type | Where to capture |
|--------------|------------------|
| New or changed requirement | `openspec/.../specs/` or delta specs (per project) |
| Design decision | `design.md` |
| Scope change | `proposal.md` |
| New work | `tasks.md` |
| Invalidated assumption | Relevant artifact |

### What you do not have to do

Follow a single script, force a single artifact, or rush to a conclusion. Long tangents are fine if they are valuable.

### Guardrails

- **Do not implement** product behavior in this phase. Planning markdown under `openspec/` or **`.skillgrid/prd/`** is allowed.
- **Do not** copy meta-blocks meant for the agent (e.g. raw `<context>` / `<rules>`) into user-facing files—use them as constraints only.
- **Do** question assumptions and visualize when it helps.
- **Do** ground discussion in the repo when relevant.

## PRD template (use when a change has no PRD)

Adapt headings if the repo’s own template overrides. **Filename** matches **`/skillgrid-plan`:** `.skillgrid/prd/PRD<NN>_<slug>.md` (execution order = `<NN>`) and **`.skillgrid/prd/INDEX.md`** (sorted by `<NN>`). Check existing **`.skillgrid/prd/PRD*.md`** before choosing `<NN>`.

#### Title block

- Heading: `### PRD: <Title>`
- **File:** `.skillgrid/prd/PRD<NN>_<slug>.md` — execution order = `<NN>`
- **Spec / change:** `<path>` — canonical source for status and technical artifacts (e.g. `openspec/changes/<id>/` or project equivalent)
- **Status:** Use the Skillgrid lifecycle: `draft` → `todo` → `inprogress` → `devdone` → `done` (set by each phase per **`/skillgrid-init`**) when the PRD is created or updated here
- **Depends on (optional):** other `PRDNN_` files that must land first
- **Tech / stack (optional):** one line — see full PRD outline in **`/skillgrid-plan`** (**Decomposition**, **Codebase touchpoints**, **Quality bar**, **Author self-review**)

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

## Optional: IDE personas

When spawning a **subagent** for exploration-only work in a clean context, use **`skillgrid-explore-architect`** ([`.cursor/agents/skillgrid-explore-architect.md`](../../.cursor/agents/skillgrid-explore-architect.md)).

For **external / cited research** rather than in-repo mapping, use **`skillgrid-researcher`** ([`.cursor/agents/skillgrid-researcher.md`](../../.cursor/agents/skillgrid-researcher.md)).

## Notes

- Inspect the repo with tools; do not assume stack or layout.
- **Hybrid** is the default: **`openspec/`** on disk plus Engram for durable summaries. If something is still missing, run **`/skillgrid-init`** or align with existing `openspec/` and `AGENTS.md`.

## Completion report (required)

End with a **Session wrap-up** the user can scan:

1. **What I did** — Bullets: what was explored, which changes/PRDs inspected, which **`.skillgrid/project/`** files were created or updated, and key findings.
2. **Token / usage** — If the product shows **input/output tokens**, **context used**, or **session cost** for this turn, report it. If not available, state **`Token usage: not shown in this environment`** (do not guess).
3. **Suggested next command** — **`/skillgrid-plan`** to formalize scope and open an OpenSpec change (or **`/skillgrid-init`** if the repo is still unbootstrapped).

</process>
