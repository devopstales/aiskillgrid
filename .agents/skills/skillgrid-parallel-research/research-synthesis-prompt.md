# Skillgrid Research Synthesis Prompt Template

Use this template after independent research lanes have produced files under `.skillgrid/tasks/research/<change-id>/`.

```markdown
You are synthesizing research for Skillgrid change `<change-id>`.

## Decision Needed

<The concrete decision that the research should support.>

## Inputs

- PRD: `.skillgrid/prd/PRD<NN>_<slug>.md` or `<none yet>`
- OpenSpec change: `openspec/changes/<change-id>/` or `<none yet>`
- Handoff: `.skillgrid/tasks/context_<change-id>.md`
- Research files:
  - `.skillgrid/tasks/research/<change-id>/<lane-a>.md`
  - `.skillgrid/tasks/research/<change-id>/<lane-b>.md`

## Synthesis Rules

- Read each lane report before synthesizing.
- Prefer direct evidence over repeated claims.
- Identify agreement, conflict, and gaps across lanes.
- Do not invent new research; list missing evidence as follow-up.
- Produce a decision-oriented recommendation, not a generic summary.
- Update the handoff only if explicitly assigned by the parent.

## Output File Format

Write the synthesis to `.skillgrid/tasks/research/<change-id>/research-synthesis.md`:

# Research Synthesis: <change-id>

## Decision Needed

<decision>

## Evidence Index

- `.skillgrid/tasks/research/<change-id>/<lane>.md` - <one-line finding>

## Cross-Lane Findings

- <finding supported by one or more lanes>

## Conflicts Or Gaps

- <conflict, uncertainty, or `none`>

## Recommendation

<recommended path and rationale>

## Confidence

High | Medium | Low - <reason>

## Follow-up

- [ ] `[HITL]` <human question if needed>
- [ ] `[AFK]` <next autonomous research, planning, or implementation step>

## Return Format

- Wrote: `.skillgrid/tasks/research/<change-id>/research-synthesis.md`
- Recommendation: <one sentence>
- Confidence: High | Medium | Low
- Blocker: <none or short blocker>
```
