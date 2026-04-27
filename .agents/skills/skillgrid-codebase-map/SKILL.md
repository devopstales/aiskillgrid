---
name: skillgrid-codebase-map
description: >
  Maps repo structure, conventions, tests, design tokens, and graphify output before Skillgrid planning or implementation.
  Trigger: Exploring a brownfield repo, planning a change, refreshing project docs, or checking architectural context.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

## When to Use

Use this skill before large Skillgrid planning, brownfield exploration, task breakdown, or implementation that depends on existing architecture and conventions.

## Critical Patterns

### Mapping Sources

Use the cheapest accurate source first:

1. `.skillgrid/project/STRUCTURE.md`, `ARCHITECTURE.md`, `PROJECT.md`
2. `graphify-out/wiki/index.md` when present
3. `graphify-out/GRAPH_REPORT.md` when present
4. `rg`, glob, and semantic search for exact files and symbols
5. source files and tests

### Graphify

Graphify is a codebase-index tool, not a separate workflow.

- Read graphify reports before answering architecture questions when present.
- Pair graphify with exact search.
- After modifying code files, run `graphify update .` when the project policy requires it.
- If graphify output is stale or absent, do not block the work; map with repo search and note the gap.

### Pattern Analysis

Capture:

- major directories and ownership boundaries
- test commands and test file conventions
- framework and package manager
- API, UI, data, and background job boundaries
- design tokens and styling system
- existing PRD/OpenSpec/Skillgrid artifacts
- known risky areas or god nodes

### Durable Output

If mapping discovers stable facts, update or recommend updating:

- `.skillgrid/project/STRUCTURE.md`
- `.skillgrid/project/ARCHITECTURE.md`
- `.skillgrid/project/PROJECT.md`
- `DESIGN.md` for stable design tokens only

Do not dump raw search output into project docs. Summarize conventions and link paths.

## Commands

```bash
rg "<symbol-or-pattern>"
graphify update .
```

Use `rg`/IDE search over broad shell scans.

## Resources

- Project docs: `skillgrid-project-docs`
- Parallel research: `skillgrid-parallel-research`
- Workflow overview: `docs/workflow.md`
- Repo rules: `.configs/AGENTS.md`
