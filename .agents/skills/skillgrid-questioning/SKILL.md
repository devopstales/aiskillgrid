---
name: skillgrid-questioning
description: >
  Guides concise user questioning during Skillgrid phases and records answers into durable artifacts.
  Trigger: Requirements, scope, design, ticketing, validation, or implementation choices are unclear.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

## When to Use

Use this skill when a Skillgrid command needs user input before it can safely plan, create artifacts, choose a path, validate work, or proceed with implementation.

## Critical Patterns

### Ask Only Blocking Questions

Ask when the answer changes:

- PRD scope
- OpenSpec requirements
- implementation boundary
- ticketing provider behavior
- UI direction
- security or data handling
- validation sign-off

Do not ask for trivia the agent can infer from the repo.

### Explore Before Asking

If a question can be answered by reading the repo, artifacts, lints, or existing docs, gather that evidence instead of asking the user. Ask only after exploration leaves a real decision for the user.

### Brainstorming Gate

Before a Skillgrid phase turns a vague idea into PRDs, specs, tasks, or implementation, shape it into an approved direction:

1. Explore project context: existing PRDs, OpenSpec changes, project docs, recent relevant files, and known constraints.
2. Assess scope. If the request spans multiple independent subsystems, propose decomposition before refining details.
3. Ask clarifying questions one at a time, focused on purpose, constraints, success criteria, and boundaries.
4. Propose 2-3 approaches with tradeoffs and a recommended default.
5. Present the selected design or plan in sections scaled to complexity.
6. Record approval, accepted assumptions, or `[HITL]` open questions before implementation begins.

Do not treat "simple" as permission to skip design. For small changes, the design may be only a few sentences, but assumptions still need to be explicit.

### Question Batch Size

- Ask 1-2 critical questions at a time.
- Prefer multiple-choice when the decision space is known.
- Include a recommended default when one is safe.
- If a later question depends on an earlier one, ask the earlier question first and name the dependency.
- For high-impact scope, design, data model, API, or implementation-boundary decisions, ask exactly one question at a time.
- If the user gives a vague answer, offer 2-3 concrete options and clearly label the recommended one instead of silently advancing.

### Decision-Tree Interview Mode

Use a stricter interview when the user asks to stress-test, challenge, or clarify a plan, or when a Skillgrid phase would otherwise proceed with major unknowns.

Walk the decision tree in dependency order:

1. Product goal and audience before scope.
2. Scope before PRD split or OpenSpec shape.
3. Data model or state ownership before API shape.
4. API or module boundary before UI behavior.
5. UI direction before implementation details.
6. Acceptance criteria before task breakdown.

For each blocking branch:

- Ask one question, then wait.
- Explain what the answer unblocks.
- Provide a recommended answer when there is a safe default.
- If the user says "I don't know", present options and mark one recommended.
- Do not move to implementation details until the major branches are decided or explicitly recorded as open questions.

The interview is complete when you can restate the design or plan without gaps. End with:

- decisions made, with brief rationale
- assumptions accepted
- open questions still requiring resolution

### Good Question Shape

```markdown
I need one decision before writing the PRD:

Should this change ship as:
- A. one PRD and one OpenSpec change (recommended)
- B. separate PRDs for each vertical slice
```

### Record Answers

After the user answers, write the decision to the right artifact:

| Answer type | Artifact |
|---|---|
| product goal or scope | PRD |
| technical approach | OpenSpec `design.md` |
| acceptance criteria | PRD and delta specs |
| implementation ordering | `tasks.md` |
| design preference | `DESIGN.md`, PRD, or OpenSpec `design.md` |
| durable cross-session decision | Engram summary when available |

### Avoid Interview Loops

Questioning should unblock progress, not become a ritual. If more than two rounds are needed for a normal Skillgrid phase, suggest `/skillgrid-brainstorm` or create an explicit `[HITL]` task.

Decision-tree interview mode is the exception: continue until each major branch has a decision, an accepted assumption, or an explicitly recorded open question.

### Never

- Never accept "I'll figure that out later" without recording it as an open question or `[HITL]` task.
- Never ask multiple unrelated questions in one turn.
- Never ask about implementation details before upstream product, data, API, or UI decisions they depend on.
- Never keep asking when existing repository evidence already answers the question.

## Commands

No shell command is required. Use the IDE’s structured question tool when available; otherwise ask directly in chat.

## Resources

- PRD artifacts: `skillgrid-prd-artifacts`
- Vertical slices: `skillgrid-vertical-slices`
- UI decisions: `skillgrid-ui-design-artifacts`
- Persistence: `skillgrid-hybrid-persistence`
