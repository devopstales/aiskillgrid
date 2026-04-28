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

### Independent Domain Check

Dispatch parallel lanes only when the domains are genuinely independent:

- each lane can understand its question without waiting on another lane
- lanes do not need to edit the same artifact
- one lane's answer is unlikely to invalidate another lane's work
- each lane has a different evidence target, subsystem, or decision

Use one lane per problem domain. If failures or research questions may share a root cause, investigate together first or run a small parent-led pass before dispatching.

Do not use parallel lanes when:

- the work requires a full-system synthesis before the questions are known
- lanes would compete for the same files, credentials, browser session, or external quota
- the user needs one sequential decision before later research makes sense
- a single root-cause fix might answer all questions

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

Each lane prompt should be focused and self-contained:

- **Specific scope:** one subsystem, source set, question, test area, or design direction
- **Clear goal:** the decision the lane must inform
- **Constraints:** what not to inspect, change, or assume
- **Expected output:** file path, evidence summary, recommendation, and blockers

Avoid broad prompts such as "research this feature". Prefer "compare current router docs against our existing routing layer and write the migration risks to `<path>`".

### Local Prompt Templates

Use local templates for repeated research roles:

- `skillgrid-parallel-research/research-lane-prompt.md` - one independent repo, docs, web, security, UI, or test-strategy lane
- `skillgrid-parallel-research/research-synthesis-prompt.md` - decision-oriented synthesis across completed lane outputs

Construct each prompt from durable artifacts: PRD, OpenSpec paths when present, handoff path, concrete research question, output path, and expected return format. Do not paste session history or chain-of-thought into research prompts.

### Handoff Update

After research completes, update `.skillgrid/tasks/context_<change-id>.md` with:

- research file path
- one-line finding
- decision or open question
- whether the work is `[HITL]` or `[AFK]`

### Parallel Research Plan Template

Use before launching several research lanes. Each lane should use `research-lane-prompt.md` and write a distinct file:

```markdown
## Parallel Research Plan: <change-id>

| Lane | Question | Tool / agent | Output path | Decision informed |
|---|---|---|---|---|
| repo-map | <What exists locally?> | `skillgrid-codebase-map` | `.skillgrid/tasks/research/<change-id>/repo-map.md` | <decision> |
| docs | <What do current docs say?> | `documentation-lookup` / `context7` | `.skillgrid/tasks/research/<change-id>/docs.md` | <decision> |
| prior-art | <How do others solve this?> | `deep-research` / `exa-search` | `.skillgrid/tasks/research/<change-id>/prior-art.md` | <decision> |
```

Before launch, sanity-check the plan:

- [ ] Every lane has a unique question.
- [ ] Every lane has a unique output path.
- [ ] No two lanes are expected to edit the same artifact.
- [ ] The parent knows how the lane results will be synthesized.
- [ ] The stop condition is clear.

### Research Synthesis Template

Use `research-synthesis-prompt.md` after lane reports are complete:

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
