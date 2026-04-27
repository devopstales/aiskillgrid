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

### Question Batch Size

- Ask 1-2 critical questions at a time.
- Prefer multiple-choice when the decision space is known.
- Include a recommended default when one is safe.
- If the user gives a vague answer, restate the decision and proceed with the safest interpretation.

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

Questioning should unblock progress. If more than two rounds are needed, suggest `/skillgrid-brainstorm` or create an explicit `[HITL]` task.

## Commands

No shell command is required. Use the IDE’s structured question tool when available; otherwise ask directly in chat.

## Resources

- PRD artifacts: `skillgrid-prd-artifacts`
- Vertical slices: `skillgrid-vertical-slices`
- UI decisions: `skillgrid-ui-design-artifacts`
- Persistence: `skillgrid-hybrid-persistence`
