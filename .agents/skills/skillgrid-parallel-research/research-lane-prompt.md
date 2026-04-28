# Skillgrid Research Lane Prompt Template

Use this template when dispatching one independent research lane for a Skillgrid change.

```markdown
You are running research lane `<lane-name>` for Skillgrid change `<change-id>`.

## Research Question

<One specific question this lane must answer.>

## Decision This Informs

<The planning, design, implementation, testing, or validation decision that depends on this evidence.>

## Context

- PRD: `.skillgrid/prd/PRD<NN>_<slug>.md` or `<none yet>`
- OpenSpec change: `openspec/changes/<change-id>/` or `<none yet>`
- Handoff: `.skillgrid/tasks/context_<change-id>.md`
- Related files or docs: `<paths or none>`
- Output path: `.skillgrid/tasks/research/<change-id>/<lane-name>.md`

## Scope

- Stay within this lane's question.
- Do not duplicate another lane's work.
- Use the most authoritative source available for the question.
- Cite web or external documentation claims with source links.
- Stop once the evidence is sufficient to answer the decision question.

## Output File Format

Write the long result to `.skillgrid/tasks/research/<change-id>/<lane-name>.md`:

# Research Lane: <lane-name>

## Question

<question>

## Short Answer

<answer in 2-5 sentences>

## Evidence

- <source or file path> - <finding>

## Recommendation

<recommended action and why>

## Confidence

High | Medium | Low - <reason>

## Follow-up

- [ ] `[HITL]` <human decision if needed>
- [ ] `[AFK]` <next autonomous step if useful>

## Return Format

- Wrote: `.skillgrid/tasks/research/<change-id>/<lane-name>.md`
- Finding: <one sentence>
- Recommendation: <one sentence>
- Blocker: <none or short blocker>
```
