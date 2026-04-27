---
name: skillgrid-parallel-research
description: >
  Coordinates parallel Skillgrid research across repo exploration, web search, documentation lookup, and subagents.
  Trigger: External evidence, current documentation, comparative options, or broad information gathering is needed.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

## When to Use

Use this skill during `/skillgrid-explore`, `/skillgrid-brainstorm`, `/skillgrid-plan`, design work, or validation when independent research threads can run in parallel.

## Critical Patterns

### Split By Question

Good research splits:

- repo architecture
- current library documentation
- prior art or alternatives
- security implications
- UI inspiration
- test strategy

Bad splits:

- multiple agents searching the same query
- broad web crawling “just in case”
- research without a decision it will inform

### Tool Choice

| Need | Prefer |
|---|---|
| existing repo behavior | `skillgrid-codebase-map`, `rg`, semantic search |
| library/API docs | `documentation-lookup`, `context7` |
| broad web evidence | `search-first`, `deep-research`, `exa-search`, Brave |
| pages that need extraction | Firecrawl skills |
| browser runtime evidence | `browser-testing-with-devtools` |
| UI alternatives | `skillgrid-ui-design-artifacts`, `superdesign`, `impeccable` |

### Output Discipline

Long research belongs in:

```text
.skillgrid/tasks/research/<change-id>/
```

The parent chat should receive:

- short answer
- evidence paths
- source links when web claims are used
- recommendation
- unresolved risks

### Handoff Update

After research completes, update `.skillgrid/tasks/context_<change-id>.md` with:

- research file path
- one-line finding
- decision or open question
- whether the work is `[HITL]` or `[AFK]`

### Parallel Research Plan Template

Use before launching several research lanes:

```markdown
## Parallel Research Plan: <change-id>

| Lane | Question | Tool / agent | Output path | Decision informed |
|---|---|---|---|---|
| repo-map | <What exists locally?> | `skillgrid-codebase-map` | `.skillgrid/tasks/research/<change-id>/repo-map.md` | <decision> |
| docs | <What do current docs say?> | `documentation-lookup` / `context7` | `.skillgrid/tasks/research/<change-id>/docs.md` | <decision> |
| prior-art | <How do others solve this?> | `deep-research` / `exa-search` | `.skillgrid/tasks/research/<change-id>/prior-art.md` | <decision> |
```

### Research Synthesis Template

```markdown
# Research Synthesis: <change-id>

## Decision needed

<Question the research was meant to answer.>

## Evidence index

- `.skillgrid/tasks/research/<change-id>/<lane>.md` — <one-line finding>

## Recommendation

<Recommended path and why.>

## Confidence

High | Medium | Low — <reason>

## Follow-up

- [ ] `[HITL]` <human question if needed>
- [ ] `[AFK]` <next autonomous research or implementation step>
```

### Stop Rule

Stop researching once the evidence is enough to decide. Additional searches must have a new hypothesis.

## Commands

No single command owns this skill. Use the available search, documentation, browser, and subagent tools in parallel when independent.

## Resources

- Handoff: `skillgrid-filesystem-handoff`
- Subagents: `skillgrid-subagent-orchestration`
- Codebase mapping: `skillgrid-codebase-map`
- Web research rules: `.configs/AGENTS.md`
