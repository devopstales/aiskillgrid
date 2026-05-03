---
name: skillgrid-prd-artifacts
description: >
  Manages Skillgrid PRD artifacts, numbering, INDEX entries, status fields, and PRD links to OpenSpec changes.
  Trigger: Creating, backfilling, renumbering, splitting, or updating Skillgrid PRDs.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

## When to Use

Use this skill when a Skillgrid command needs to create or update files under `.skillgrid/prd/`, align a PRD with an OpenSpec change, update `.skillgrid/prd/INDEX.md`, or reason about PRD execution order and lifecycle status.

## Critical Patterns

### Canonical Location

- Canonical PRDs live under `.skillgrid/prd/`.
- Do not create new PRDs at repository root `prd/`.
- `docs/PRD/` may mirror or link to canonical PRDs, but `.skillgrid/prd/` owns the workflow state.

### Synthesize From Context First

When creating a PRD from the current conversation or an existing OpenSpec change, do not start with a broad interview. First synthesize what is already known:

1. Explore the repo and existing Skillgrid/OpenSpec artifacts if you have not already.
2. Summarize the user-facing problem, desired solution, constraints, and out-of-scope items from current context.
3. Sketch the major modules, surfaces, or system boundaries likely to be built or modified.
4. Look for deep-module opportunities: simple, stable interfaces that encapsulate meaningful behavior and can be tested in isolation.
5. Ask only for decisions that remain genuinely blocking. Use `skillgrid-questioning` for those.

If the user explicitly wants a GitHub/GitLab/Jira issue, hand off to `skillgrid-issue-creation` after the canonical PRD exists. Do not make the remote issue the only source of product intent.

### File Naming

Use:

```text
.skillgrid/prd/PRD<NN>_<descriptive-slug>.md
```

Rules:

- `<NN>` is a two-digit execution order: `01`, `02`, `03`.
- The slug is lowercase `kebab-case` or `lower_snake`.
- List existing `PRD*.md` files before choosing a number.
- Avoid duplicate numbers.
- If execution order changes, rename files and update `INDEX.md`, PRD cross-links, OpenSpec references, and Engram pointers.

### Work hierarchy (shared vocabulary)

Use one mental model across skills, AGENTS, and `docs/02-workflow-usage.md`:

| Level | Jira-style | GitHub-style | Skillgrid artifact |
|-------|------------|--------------|---------------------|
| Program / milestone | Epic | Milestone | `.skillgrid/prd/INDEX.md` — dependency-ordered PRDs **plus** execution snapshot (below) |
| Feature initiative | Task | Issue | `.skillgrid/prd/PRD<NN>_<slug>.md` + one `openspec/changes/<change-id>/` |
| Shippable unit | Sub-task | Checklist item | **Vertical slice** — `tasks.md` rows and/or `openspec/changes/<change-id>/specs/<slice-slug>/spec.md` |

There is **no** separate `.skillgrid/project/TASK.md`: progress-tracker and task-list behavior live in **INDEX snapshot sections** and in OpenSpec `tasks.md` / slice specs.

### INDEX Discipline

`.skillgrid/prd/INDEX.md` is the **single** ordered work index for PRDs. Keep PRD rows sorted by `PRD<NN>` **in dependency order** (blocked PRDs appear after their blockers in the table, or list `Depends on` per row).

Use the ordered PRD sequence as Skillgrid's lightweight roadmap or milestone view. A broad initiative should appear as several ordered PRDs, where each PRD is one independently reviewable product slice.

**Execution snapshot (optional but recommended):** At the **top** of `INDEX.md`, maintain short-lived fields so any session can see *where we are* without reading chat history. Refresh when phase, focus, or discovered work changes:

- **Current phase** — plan / breakdown / apply / validate / finish (or custom from `.skillgrid/config.json`).
- **Active change** — `openspec/changes/<change-id>/`.
- **Active slice** — `<slice-slug>` when implementing a vertical slice (`specs/<slice-slug>/spec.md`).
- **In progress** — bullets (PRD id, slice, or task id).
- **Recently completed** — last few shipped slices or PRDs.
- **Next up** — next slice or PRD after current.
- **Discovered during work** — backlog bullets to promote into PRD or `tasks.md`.
- **Open questions / blockers** — HITL reminders.
- **Session notes** — 2–5 lines for cold resume (complements Engram when enabled).

Then the **PRD index table** (see INDEX Template).

Each PRD row should include at least:

- PRD file
- title
- linked `openspec/changes/<change-id>/`
- `Status:`
- optional **Depends on** (other `PRD<NN>` or external blocker)
- optional external issue key or URL when ticketing is not `local`

### Title Block

Every PRD should start with:

```markdown
### PRD: <Title>

- **File:** `.skillgrid/prd/PRD<NN>_<slug>.md`
- **Spec / change:** `openspec/changes/<change-id>/`
- **Session context:** `.skillgrid/tasks/context_<change-id>.md`
- **Status:** `draft`
- **Preview:** optional `.skillgrid/preview/<change-id>-options.html`
- **External:** optional issue key or URL
- **Depends on:** optional PRD dependencies
- **Tech / stack:** optional one-line summary
```

Adapt only when the repository has an explicit stronger template.

### Status Lifecycle

Read `.skillgrid/config.json` before changing PRD `Status:` values. The source of truth is:

```json
{
  "prdWorkflow": {
    "fallbackStatus": "draft",
    "statuses": [
      { "id": "draft", "label": "Draft" },
      { "id": "todo", "label": "Todo" },
      { "id": "inprogress", "label": "In Progress" },
      { "id": "devdone", "label": "Dev Done" },
      { "id": "done", "label": "Done" },
      { "id": "archived", "label": "Archived" }
    ],
    "phaseStatusMap": {
      "plan": "draft",
      "breakdown": "todo",
      "apply": "inprogress",
      "validate": "devdone",
      "finish": "archived"
    }
  }
}
```

When `prdWorkflow` is missing, use the default lifecycle:

```text
draft -> todo -> inprogress -> devdone -> done -> archived
```

Commands own phase transitions through `prdWorkflow.phaseStatusMap`:

| Phase | Command | Status on successful exit |
|---|---|---|
| Plan | `/skillgrid-plan` | `phaseStatusMap.plan` (default `draft`) |
| Breakdown | `/skillgrid-breakdown` | `phaseStatusMap.breakdown` (default `todo`) |
| Apply | `/skillgrid-apply` | `phaseStatusMap.apply` (default `inprogress`) |
| Validate | `/skillgrid-validate` | `phaseStatusMap.validate` (default `devdone`) |
| Finish | `/skillgrid-finish` | `phaseStatusMap.finish` (default `archived` after OpenSpec archive) |

When backfilling an existing change, map evidence to the configured workflow:

- archived OpenSpec change -> `phaseStatusMap.finish`
- tasks with completed checkboxes -> `phaseStatusMap.apply`
- proposal-only or planning-only -> `phaseStatusMap.plan`

### Required PRD Sections

Copy the canonical blank from **`.skillgrid/templates/template-prd.md`** when creating a new PRD file. Human-readable planning logic and template index: **`docs/skillgrid-templates-and-logic.md`**.

Keep product intent in the PRD. Detailed CLI instructions and exact code steps belong in OpenSpec artifacts or `tasks.md`.

### INDEX Template

Use **`.skillgrid/templates/template-index.md`** when creating or refreshing `.skillgrid/prd/INDEX.md` (trim snapshot bullets if unused). See **`docs/skillgrid-templates-and-logic.md`** for hierarchy.

If ticketing is `local`, `External` may be omitted or set to `local`. Omit the **Depends on** column if every PRD is independent.

### PRD And OpenSpec Alignment

- PRD is the product intent source.
- OpenSpec `proposal.md`, `design.md`, delta specs, and `tasks.md` are the technical contract.
- **Vertical slices** should have scoped requirements under `openspec/changes/<change-id>/specs/<vertical-slice-slug>/spec.md` (see `skillgrid-spec-artifacts`). Optional umbrella: `openspec/specs/<change-id>/spec.md` for cross-cutting requirements.
- Every PRD should link to exactly one primary `openspec/changes/<change-id>/` unless it is intentionally an umbrella PRD.
- If a PRD becomes too broad, split it into ordered PRDs instead of expanding one mega-PRD.

### Source-Of-Truth Rules

- PRD owns product intent: user-facing problem, goals, scope, success criteria, and slice boundary.
- OpenSpec delta specs own verifiable technical behavior and scenarios.
- OpenSpec `tasks.md` owns the implementable checklist; detailed file-by-file steps do not belong in the PRD body.
- External tracker issues mirror PRD slices or blockers for coordination. If an external issue changes product intent, import that decision back into the PRD or OpenSpec artifacts.
- Engram memory summarizes durable decisions and pointers. It must not become the only place a requirement lives.
- When artifacts disagree, stop and reconcile before creating more PRDs, issues, or tasks.

## Commands

```bash
ls .skillgrid/prd/PRD*.md
openspec list --json
```

## Resources

- Templates and planning logic: `docs/skillgrid-templates-and-logic.md`
- Canonical blanks: `.skillgrid/templates/template-*.md` (see `README.md` there)
- Full workflow overview: `docs/02-workflow-usage.md`
- Command sources: `.cursor/commands/skillgrid-plan.md`, `.cursor/commands/skillgrid-explore.md`, `.cursor/commands/skillgrid-breakdown.md`
- Related skills: `skillgrid-spec-artifacts`, `skillgrid-vertical-slices`, `skillgrid-filesystem-handoff`
