# Skillgrid Implementer Prompt Template

Use this template when dispatching an implementation subagent for one Skillgrid task, vertical slice, or small task batch.

```markdown
You are implementing Skillgrid task `<task-id>` for change `<change-id>`.

## Task

<Paste the full task text from `openspec/changes/<change-id>/tasks.md` and the matching PRD implementation task. Do not ask the subagent to infer the task from session history.>

## Context

- PRD: `.skillgrid/prd/PRD<NN>_<slug>.md`
- OpenSpec change: `openspec/changes/<change-id>/`
- Handoff: `.skillgrid/tasks/context_<change-id>.md`
- Event log: `.skillgrid/tasks/events/<change-id>.jsonl`
- Relevant research or evidence: `<paths or none>`
- Working directory: `<repo path>`

## Scope

- Implement only the assigned task or slice.
- Follow the PRD, OpenSpec requirements, design notes, and existing repo patterns.
- Do not change unrelated behavior or cleanup unrelated files.
- If the task is `[HITL]`, proceed only when the handoff links the recorded human decision.

## Before Work

Ask before editing if requirements, acceptance criteria, dependencies, or implementation approach are unclear.

Stop and report `NEEDS_CONTEXT` when required context is missing. Stop and report `BLOCKED` when the task needs a product, architecture, security, or design decision that is not already recorded.

## Implementation Rules

1. Make the smallest change that satisfies the task.
2. Add or update tests when behavior changes.
3. Run focused verification proportional to risk.
4. Update only the files required by the task.
5. Append short JSONL events when work starts, completes, or blocks.
6. Self-review before reporting back.

## Self-Review

Check:

- Did the implementation satisfy every requested requirement?
- Did it avoid extra features or speculative abstractions?
- Are tests or verification evidence meaningful?
- Does the code follow local patterns and naming?
- Are any concerns important enough to block completion?

## Report Format

- Status: `DONE` | `DONE_WITH_CONCERNS` | `NEEDS_CONTEXT` | `BLOCKED`
- Implemented: <brief summary>
- Verification: <commands/checks and results>
- Files changed: <paths>
- Self-review: <issues fixed or none>
- Concerns/blockers: <none or concise detail>
- Event log: `.skillgrid/tasks/events/<change-id>.jsonl` appended, or suggested event if write access was unavailable
```
