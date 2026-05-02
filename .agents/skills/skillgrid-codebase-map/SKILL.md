---
name: skillgrid-codebase-map
description: >
  Maps repo structure, conventions, tests, design tokens, CocoIndex (ccc), and GitNexus output before Skillgrid planning or implementation.
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

1. **Exploration order:** For understanding or mapping **application code**, use **GitNexus MCP** (`query`, `context`, `impact`, `READ gitnexus://…/clusters`) **before** grep/glob/shell file walks. Follow **`gitnexus-exploring`** for the step-by-step workflow. Subagents inherit this.
2. `.skillgrid/project/STRUCTURE.md`, `ARCHITECTURE.md`, `PROJECT.md` — team summaries; skim first, then confirm in the graph.
3. GitNexus MCP/resources when configured (`query`, `context`, `impact`, repository clusters/processes)
4. `.gitnexus/` status when present
5. **`ccc search`** for concept-style discovery when the index exists (after **`ccc init`** / **`ccc index`** per [ccc skill](../ccc/SKILL.md))
6. `rg`, glob, and IDE search for exact files and symbols — **after** GitNexus for behavioral questions, or when the index/MCP is unavailable (re-run **`npx -y gitnexus@1.3.11 analyze`**, then retry GitNexus)
7. Source files and tests

### CocoIndex (`ccc`)

Semantic search over the working tree—pairs with GitNexus and **`rg`**.

- Prefer **`ccc search "<natural language>"`** when an index exists for “where is X implemented?” in concept form.
- If **`ccc`** is installed but the project is not initialized, run **`ccc init`** then **`ccc index`** at the repo root (**`/skillgrid-init`** does this when the CLI is present). After large tree changes, run **`ccc index`** again (see **ccc** skill).
- If **ccc** is absent or unindexed, rely on GitNexus + **`rg`** and note the gap.

### GitNexus

GitNexus is the codebase knowledge graph, MCP, and graph UI provider for Skillgrid.

- Use GitNexus MCP/resources before broad raw file reads when available.
- Pair GitNexus graph context with exact search.
- If **`.gitnexus/`** is missing but GitNexus is available, run an **initial** index from the repo root: **`npx -y gitnexus@1.3.11 analyze`**, e.g. during **`/skillgrid-init`**. After substantive code changes, run **`npx -y gitnexus@1.3.11 analyze`** again. Use **`npx -y gitnexus@1.3.11 analyze --force`** only for full rebuilds.
- If GitNexus output is stale or absent, do not block the work; map with repo search and note the gap.

### Pattern Analysis

Capture:

- major directories and ownership boundaries
- test commands and test file conventions
- framework and package manager
- API, UI, data, and background job boundaries
- design tokens and styling system
- existing PRD/OpenSpec/Skillgrid artifacts
- known risky areas or god nodes

### Agent-Friendly Architecture

Planning quality depends on module shape. While mapping, call out whether the target area has deep, testable modules with small interfaces or shallow/tangled modules that will make AFK implementation risky.

Record:

- proposed modules or boundaries a change is likely to touch;
- public interfaces that should stay small and stable;
- test seams and feedback loops already available;
- missing tests, god modules, circular dependencies, or broad side effects that should route the work through exploration/refactoring before large apply loops.

Prefer updating the journey plan to include a small architecture-deepening or refactoring slice before implementation when unclear boundaries would otherwise push agents into the dumb zone.

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
ccc search "natural language query about behavior or architecture"
npx -y gitnexus@1.3.11 analyze
```

Use `rg`/IDE search over broad shell scans.

## Resources

- GitNexus exploration workflow: `gitnexus-exploring`
- Project docs: `skillgrid-project-docs`
- Parallel research: `skillgrid-parallel-research`
- Workflow overview: `docs/02-workflow-usage.md`
- Repo rules: `.configs/AGENTS.md`
