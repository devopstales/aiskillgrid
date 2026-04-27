---
name: skillgrid-ui-design-artifacts
description: >
  Places UI design decisions, previews, and design constraints into the correct Skillgrid artifacts.
  Trigger: Designing UI during Skillgrid brainstorm, design, plan, explore, or validation phases.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

## When to Use

Use this skill when UI/UX decisions need to become durable Skillgrid artifacts, especially when `/skillgrid-design`, `/skillgrid-brainstorm`, or `/skillgrid-plan` touches interface behavior.

## Critical Patterns

### Responsibilities

This skill owns where design decisions live:

- root `DESIGN.md`
- `.skillgrid/preview/`
- PRD success criteria
- OpenSpec `design.md`
- handoff files and research indexes

It does not own visual taste. Reuse UI skills for design quality:

- `superdesign`
- `impeccable`
- `frontend-design`
- `frontend-ui-engineering`
- high-end visual design skills when appropriate

### Artifact Placement

| Design fact | Artifact |
|---|---|
| stable design system token | `DESIGN.md` |
| draft or comparison | `.skillgrid/preview/` |
| user-facing acceptance criterion | PRD |
| implementation constraint | OpenSpec `design.md` |
| subagent or external reference | `.skillgrid/tasks/research/<change-id>/` |

### Parallel Design Options

When multiple UI directions are plausible, create or request distinct options and evaluate them against:

- PRD goals
- accessibility
- implementation cost
- design system fit
- responsive behavior
- validation criteria

Do not ship a visual direction just because it was generated first.

### Preview Discipline

Preview artifacts should be named by change and direction:

```text
.skillgrid/preview/<change-id>-option-a.md
.skillgrid/preview/<change-id>-option-b.md
```

Clean up or archive previews in `/skillgrid-finish` when they are no longer useful.

### Preview Template

```markdown
# Preview: <change-id> / <option-name>

- **Date:** <YYYY-MM-DD>
- **Surface:** <page, component, flow, or state>
- **Source PRD:** `.skillgrid/prd/PRD<NN>_<slug>.md`
- **OpenSpec design:** `openspec/changes/<change-id>/design.md`

## Direction

<One paragraph describing the visual or interaction direction.>

## Rationale

- <Why this supports the PRD goals>
- <Why this fits or intentionally changes DESIGN.md>

## Key states

- Default:
- Loading:
- Empty:
- Error:
- Success:
- Mobile:

## Accessibility notes

- <Keyboard, focus, contrast, labels, motion reduction>

## Implementation notes

- <Components, tokens, layout constraints, libraries>

## Decision

Accepted | Rejected | Needs iteration

## Follow-up

- [ ] <Action needed>
```

### OpenSpec Design Insert Template

When a UI decision affects implementation, add a concise section to OpenSpec `design.md`:

```markdown
## UI / UX constraints

- **Design source:** `.skillgrid/preview/<change-id>-option-a.md`
- **Tokens:** follow root `DESIGN.md`
- **States required:** default, loading, empty, error, success
- **Accessibility:** keyboard reachable, visible focus, reduced motion, contrast compliant
- **Responsive behavior:** <breakpoints or layout rules>
```

## Commands

Use the project’s chosen design tools or skills. For SuperDesign-based work, follow `superdesign`.

## Resources

- Project docs: `skillgrid-project-docs`
- Research: `skillgrid-parallel-research`
- UI skills: `superdesign`, `impeccable`, `frontend-design`, `frontend-ui-engineering`
