---
name: sdd-gherkin
description: Draft, review, or tighten Gherkin / BDD scenarios and acceptance criteria
usage: "/sdd-gherkin <feature-path-or-inline-topic>"
---

You are executing `/sdd-gherkin` for the SDD workflow.

## Skill Loading Contract

Load and follow first:
- `.agents/skills/gherkin-authoring/SKILL.md`

## Workflow

1. Identify the Gherkin scope from `$ARGUMENTS`: `.feature` file path, change name, fenced block, or free-text scenario.
2. Preserve Markdown or document wrappers unless the user asks to edit broader prose.
3. Keep scenarios short and observable; align domain language with project glossary if present (e.g. `.skillgrid/project/CONTEXT.md`).

## Return Contract

Return a structured result with:
- `status`: `completed | blocked | failed`
- `executive_summary`
- `detailed_report` (syntax notes, scenarios added/changed)
- `artifacts` (feature file paths or spec excerpts)
- `next_recommended`
- `risks`
