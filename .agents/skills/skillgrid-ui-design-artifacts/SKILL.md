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
- `huashu-design`
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

When multiple UI directions or component/module interface shapes are plausible, create or request distinct options and evaluate them against:

- PRD goals
- accessibility
- implementation cost
- design system fit
- responsive behavior
- validation criteria
- interface simplicity and ease of correct use
- what complexity is hidden inside the component or module

Do not ship a visual direction just because it was generated first.

### Design It Twice For Interfaces

Use this pattern when the UI decision is also an interface decision: component props, page/module boundaries, API-to-UI contracts, reusable design-system primitives, or a workflow surface that other modules will call.

Before designing, gather:

- problem the interface solves
- callers or users: humans, components, modules, tests, or external consumers
- key operations and states
- constraints: accessibility, performance, compatibility, existing patterns, design system fit
- what should be hidden inside versus exposed to callers

Generate three or more radically different options. Use parallel subagents when available; otherwise produce the variants sequentially. Assign each option a distinct constraint, for example:

- Option A: minimize surface area, with the fewest props/actions/screens
- Option B: maximize flexibility for future variants
- Option C: optimize for the most common user path
- Option D: follow a specific local pattern, framework idiom, or design-system primitive

Each option should include:

- interface or surface shape: props, actions, routes, states, or screen structure
- usage example: how a caller or user interacts with it
- what it hides internally
- tradeoffs and likely failure modes

Compare options in prose, not a table, when the tradeoffs are nuanced. Cover:

- interface simplicity
- general-purpose versus specialized shape
- implementation efficiency
- depth: small external surface hiding meaningful complexity
- ease of correct use versus ease of misuse

Synthesize the final direction by asking which option best fits the primary use case and whether any elements from other options should be merged. Record only the accepted direction as implementation guidance; keep rejected options in `.skillgrid/preview/` or research notes when useful.

### Preview Discipline

For visual A/B/C options, prefer browser-viewable HTML previews scaffolded by the project helper:

```bash
.skillgrid/scripts/preview.sh <change-id>-options
```

Use Markdown previews only when a text-only comparison is enough:

```bash
.skillgrid/scripts/preview.sh --md <change-id>-options
```

Preview artifacts should be named by change and direction:

```text
.skillgrid/preview/<change-id>-options.html
.skillgrid/preview/<change-id>-option-notes.md
```

Link previews from the PRD title block and handoff so the local dashboard can find them reliably:

```markdown
- **Preview:** `.skillgrid/preview/<change-id>-options.html`
```

The dashboard also discovers `.skillgrid/preview/<change-id>*.html` and `.skillgrid/preview/<prd-slug>*.html`, but explicit links are preferred when the preview name is not obvious.

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

## Options

- Option A:
- Option B:
- Option C:

## Rationale

- <Why this supports the PRD goals>
- <Why this fits or intentionally changes DESIGN.md>

## Interface shape

- Surface / props / actions:
- Example usage:
- Hidden internals:
- Tradeoffs:

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

- **Design source:** `.skillgrid/preview/<change-id>-options.html`
- **Tokens:** follow root `DESIGN.md`
- **States required:** default, loading, empty, error, success
- **Accessibility:** keyboard reachable, visible focus, reduced motion, contrast compliant
- **Responsive behavior:** <breakpoints or layout rules>
- **Interface shape:** <accepted component/page/API-to-UI contract and what it hides>
```

## Commands

Use the project’s chosen design tools or skills. For HTML-based design variants or direction-advisor work, follow `huashu-design`. For SuperDesign-based work, follow `superdesign`. For multiple UI directions, start with `.skillgrid/scripts/preview.sh <slug>` so the user has a local HTML preview to open and compare.

## Resources

- Interface design: `api-and-interface-design`
- Project docs: `skillgrid-project-docs`
- Research: `skillgrid-parallel-research`
- UI skills: `huashu-design`, `superdesign`, `impeccable`, `frontend-design`, `frontend-ui-engineering`
