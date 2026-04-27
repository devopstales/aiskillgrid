---
name: skillgrid-subagent-orchestration
description: >
  Orchestrates Skillgrid subagents with handoff files, scoped prompts, short returns, and two-stage review.
  Trigger: Delegating Skillgrid exploration, research, design critique, task execution, review, or verification to subagents.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

## When to Use

Use this skill when a Skillgrid phase dispatches subagents for parallel discovery, research, design critique, implementation support, testing, security review, or validation.

## Critical Patterns

### Parent Owns The Change

The parent session owns workflow state and final decisions. Subagents provide bounded work products.

### Prompt Contract

Every Skillgrid subagent prompt should include:

- goal
- phase
- PRD path
- OpenSpec change path when present
- `.skillgrid/tasks/context_<change-id>.md`
- expected output path under `.skillgrid/tasks/research/<change-id>/` when output is long
- exact return format: short summary plus file paths

Use this prompt skeleton:

```markdown
You are working on Skillgrid change `<change-id>`.

Goal:
<bounded task>

Phase:
<explore|brainstorm|plan|breakdown|apply|test|security|validate|finish>

Read first:
- PRD: `.skillgrid/prd/PRD<NN>_<slug>.md`
- OpenSpec: `openspec/changes/<change-id>/`
- Handoff: `.skillgrid/tasks/context_<change-id>.md`

Write long output to:
`.skillgrid/tasks/research/<change-id>/<topic>.md`

Constraints:
- Do not modify product code unless explicitly assigned implementation work.
- Keep scope limited to the goal above.
- If blocked, write the blocker and stop.
- Return only a short summary with file paths.

Return format:
- Wrote: `<path>`
- Finding: <one sentence>
- Blocker: <none or short blocker>
```

### Parallelization

Parallelize independent work:

- repo mapping vs external docs
- UI options vs API constraints
- test strategy vs security review
- multiple design directions

Do not parallelize tasks that edit the same artifact unless one parent reconciles them.

### Two-Stage Review

For subagent-generated plans, code, or reports:

1. Check compliance against PRD, OpenSpec artifacts, and task scope.
2. Check quality: correctness, maintainability, security, performance, and evidence.

Critical issues block progress until fixed, accepted, or converted into explicit follow-up work.

### Return Discipline

Subagents should not return long raw output in chat. They should write long output to disk and return:

```markdown
Done.
- Wrote: `.skillgrid/tasks/research/<change-id>/<topic>.md`
- Finding: <one sentence>
- Blocker: <none or short blocker>
```

### Two-Stage Review Template

```markdown
## Stage 1: Spec compliance

- [ ] Matches PRD goals and scope.
- [ ] Matches OpenSpec requirements and scenarios.
- [ ] Completes only assigned task/slice.
- [ ] Updates required artifacts or explains why not.

## Stage 2: Quality

- [ ] Correct and maintainable.
- [ ] Tests or evidence are adequate.
- [ ] Security/privacy concerns are addressed.
- [ ] Performance/accessibility concerns are addressed when relevant.

## Verdict

Pass | Needs changes | Blocked

## Required fixes

- <Fix or none>
```

## Commands

No shell command is required. Use the IDE’s subagent/task facility when the work is independent enough to delegate.

## Resources

- Handoff: `skillgrid-filesystem-handoff`
- Research: `skillgrid-parallel-research`
- Validation: `skillgrid-spec-artifacts`
- Workflow overview: `docs/workflow.md`
