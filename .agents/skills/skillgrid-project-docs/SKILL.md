---
name: skillgrid-project-docs
description: >
  Creates and refreshes Skillgrid project documentation, including DESIGN.md and .skillgrid/project files.
  Trigger: Initializing, exploring, mapping, or finishing Skillgrid work that changes durable project knowledge.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

## When to Use

Use this skill when Skillgrid needs durable project documentation outside a single PRD or OpenSpec change.

## Critical Patterns

### Canonical Files

```text
DESIGN.md
.skillgrid/project/ARCHITECTURE.md
.skillgrid/project/STRUCTURE.md
.skillgrid/project/PROJECT.md
```

### File Responsibilities

| File | Purpose |
|---|---|
| `DESIGN.md` | design system tokens, visual principles, UI constraints |
| `ARCHITECTURE.md` | major subsystems, boundaries, durable decisions |
| `STRUCTURE.md` | directory map and ownership conventions |
| `PROJECT.md` | stack, commands, tests, tooling, operational notes, shared language |

### Refresh Rules

Refresh docs when:

- initializing a brownfield project
- codebase mapping discovers stable facts
- a change adds durable architecture, tooling, or design conventions
- finishing a change would leave docs stale

Do not update project docs for transient implementation details.

### Shared Language And ADR Capture

Use `.skillgrid/project/PROJECT.md` for stable domain vocabulary, naming decisions, and operating language that future agents should reuse. If a decision is hard to rediscover from code alone, record it as a durable decision in `ARCHITECTURE.md`, `PROJECT.md`, or as a **MADR** file under **`.skillgrid/adr/`** (copy **`.skillgrid/templates/template-adr.md`**; see `documentation-and-adrs`). Link from `ARCHITECTURE.md` **Durable decisions** to the ADR file (e.g. `0001-chose-postgres.md`).

Capture:

- domain terms and preferred names;
- terms to avoid because they mean something else in the product;
- durable architectural decisions and rationale;
- decision links to PRDs, OpenSpec changes, issues, or review reports.

Do not let archived PRDs override current code, tests, or validated project docs. Finish/archive work should mark stale planning artifacts closed or link them to the implementation state.

### DESIGN.md

`DESIGN.md` may include machine-readable frontmatter and human-readable rationale. Extract stable tokens from:

- Tailwind config
- CSS custom properties
- theme files
- design system modules
- existing UI components

Use `skillgrid-ui-design-artifacts` for UI previews and design-option decisions.

### DESIGN.md Template

Canonical blank: **`.skillgrid/templates/template-design.md`** (repo root `DESIGN.md`). See **`docs/skillgrid-templates-and-logic.md`**.

### ARCHITECTURE.md Template

Canonical blank: **`.skillgrid/templates/template-architecture.md`** (path `.skillgrid/project/ARCHITECTURE.md`). MADR files live under `.skillgrid/adr/`; copy from **`.skillgrid/templates/template-adr.md`**.

### STRUCTURE.md Template

Canonical blank: **`.skillgrid/templates/template-structure.md`**.

### PROJECT.md Template

Canonical blank: **`.skillgrid/templates/template-project.md`**.

## Commands

```bash
rg "tailwind|theme|:root|--color|font-family"
```

## Resources

- Templates and planning logic: `docs/skillgrid-templates-and-logic.md`
- Canonical blanks: `.skillgrid/templates/template-*.md` (see `README.md` there)
- Codebase mapping: `gitnexus-exploring`, `skillgrid-codebase-map`
- UI design artifacts: `skillgrid-ui-design-artifacts`
- Workflow overview: `docs/02-workflow-usage.md`
