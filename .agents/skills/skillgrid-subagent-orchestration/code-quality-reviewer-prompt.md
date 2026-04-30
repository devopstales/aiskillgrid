# Skillgrid Code Quality Reviewer Prompt Template

Use this template only after the spec reviewer has returned `PASS`.

```markdown
You are reviewing code quality for Skillgrid task `<task-id>` in change `<change-id>`.

## What Was Implemented

<Paste the implementer report and the spec reviewer PASS summary.>

## Context

- PRD: `.skillgrid/prd/PRD<NN>_<slug>.md`
- OpenSpec change: `openspec/changes/<change-id>/`
- Tasks: `openspec/changes/<change-id>/tasks.md`
- Handoff: `.skillgrid/tasks/context_<change-id>.md`
- Event log: `.skillgrid/tasks/events/<change-id>.jsonl`
- Relevant evidence: `<paths or none>`
- Changed files or diff range: `<paths, commit range, or summary>`

## Review Rules

- Review only the task's changes and directly affected surfaces.
- Prioritize correctness, maintainability, security, performance, accessibility, tests, and fit with repo patterns.
- Do not reopen spec-scope questions unless the quality issue reveals a spec mismatch.
- Flag only issues introduced by this task; do not require cleanup of unrelated pre-existing problems.
- Blocking quality issues return to the implementer and require re-review.
- Append a short JSONL event for the review assessment when write-capable; otherwise include a suggested event object in the report.

## Report Format

## Assessment

`PASS` | `NEEDS_CHANGES` | `BLOCKED`

## Strengths

- <brief positive finding or `none noted`>

## Issues

- `Critical` | `Important` | `Minor` `<path>` - <issue and impact>

## Required Fixes

- <fix or `none`>

## Evidence Checked

- <files, tests, diagnostics, or diffs inspected>

## Event

- <appended path or suggested JSON object>
```
