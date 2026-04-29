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

### Fresh Context Only

Construct every subagent prompt from durable artifacts:

- PRD
- OpenSpec proposal, design, specs, and `tasks.md`
- `.skillgrid/tasks/context_<change-id>.md`
- cited files under `.skillgrid/tasks/research/<change-id>/`
- the specific task, scope, constraints, expected output path, and return format

Never paste session history or chain-of-thought into a subagent prompt. If a subagent needs more context, add the missing durable artifact path or a concise task-specific excerpt and re-dispatch.

### Prompt Contract

Every Skillgrid subagent prompt should include:

- goal
- phase
- PRD path
- OpenSpec change path when present
- `.skillgrid/tasks/context_<change-id>.md`
- expected output path under `.skillgrid/tasks/research/<change-id>/` when output is long
- exact return format: short summary plus file paths

Prompts should be:

- **Focused:** one clear task or problem domain
- **Self-contained:** enough artifact paths and concise excerpts to succeed without session history
- **Constrained:** explicit boundaries on what not to change or inspect
- **Specific about output:** exact summary, status, files, and blocker format

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

Bad prompt: "Fix all validation issues."

Good prompt: "Review `openspec/changes/<change-id>/tasks.md` against PRD success criteria 1-3. Write gaps to `<path>`. Do not edit files. Return pass/needs-changes plus the path."

### Local Prompt Templates

Use local templates when a Skillgrid phase repeatedly dispatches the same role:

- `skillgrid-subagent-orchestration/implementer-prompt.md`
- `skillgrid-subagent-orchestration/spec-reviewer-prompt.md`
- `skillgrid-subagent-orchestration/code-quality-reviewer-prompt.md`

These templates are role prompts, not shared memory. Fill them with fresh task context each time.

### Model Selection

Use the least capable model that can reliably complete the role:

| Work type | Model tier |
|---|---|
| Complete 1-2 file mechanical task with clear acceptance criteria | cheap / fast |
| Multi-file integration, debugging, or local pattern matching | standard |
| Design, architecture, spec review, code quality review, security judgment, broad repo reasoning | most capable |

If a subagent reports `NEEDS_CONTEXT` or `BLOCKED`, change the prompt context, split the task, or select a more capable model before retrying. Do not repeat the same failing dispatch unchanged.

### Parallelization

Parallelize independent work:

- repo mapping vs external docs
- UI options vs API constraints
- test strategy vs security review
- multiple design directions

Use one subagent per independent problem domain. Before parallel dispatch, verify:

- each subagent has a different scope and output path
- no shared mutable state is required
- subagents will not edit the same files
- one result is unlikely to invalidate the others
- the parent can review and integrate all returns

Do not parallelize:

- related failures that may share one root cause
- exploratory work where the parent does not yet know the domains
- implementation tasks that touch overlapping files or artifacts
- tasks that depend on a prior human decision
- multiple implementation subagents inside the same worktree unless their file ownership is explicit and non-overlapping

When agents return:

1. Read every summary and cited output file.
2. Check for conflicting findings, overlapping file edits, or inconsistent assumptions.
3. Run the relevant full verification, not only lane-specific checks.
4. Update the handoff with integrated findings, decisions, evidence, changed assumptions, and the next recommended action.
5. Append a short event to `.skillgrid/tasks/events/<change-id>.jsonl` when the return starts, completes, blocks, or changes workflow state.

Event entries should follow `skillgrid-filesystem-handoff` and include `time`, `changeId`, `node`, `phase`, `status`, `summary`, and useful `artifacts`. For delegated work, also include `agent` or `subagent`, `role`, `task`, and `output` when known so the dashboard Subagents view can show who did what.

### Apply Dispatch Loop

For `/skillgrid-apply` implementation delegation:

1. Parent reads PRD, OpenSpec artifacts, `tasks.md`, handoff, and relevant research.
2. Parent reviews the plan critically before edits. If instructions are unclear, verification is missing, or scope conflicts with artifacts, stop and resolve the plan before implementation.
3. Parent selects one `[AFK]` task or small vertical slice. Stop on `[HITL]` unless the handoff links the recorded decision.
4. Dispatch a fresh implementer with `implementer-prompt.md` and the full task text plus required context. Do not make the implementer read the plan and infer scope.
5. Require TDD for behavioral code unless the task explicitly records a non-TDD exception: failing test, expected failure, minimal implementation, green verification, then refactor.
6. Handle implementer status:
   - `DONE`: proceed to spec review.
   - `DONE_WITH_CONCERNS`: read concerns and resolve correctness or scope doubts before review.
   - `NEEDS_CONTEXT`: provide missing context and re-dispatch.
   - `BLOCKED`: change something before retrying: add context, split the task, use a more capable model, or escalate to the user.
7. Dispatch spec reviewer with `spec-reviewer-prompt.md`.
8. If spec review returns `NEEDS_CHANGES`, evaluate the feedback, send required fixes back to the implementer, and repeat spec review.
9. Dispatch code quality reviewer with `code-quality-reviewer-prompt.md` only after spec review passes.
10. If quality review returns `NEEDS_CHANGES`, evaluate the feedback, send required fixes back to the implementer, and repeat quality review.
11. Mark the task complete and update handoff only after both review stages pass.
12. Add a slice completion summary to the handoff: result, evidence, blockers, changed assumptions, and next recommended slice or command.
13. Append workflow events for apply start, blocker, implementation done, spec review outcome, quality review outcome, and final task completion.

Do not skip review loops. A reviewer finding means the implementer fixes it and the same stage reviews again. Spec compliance always comes before code quality.

### Work Unit Reassessment

After any delegated work unit, reassess before dispatching the next one:

- Did the result still match the PRD and OpenSpec scope?
- Did implementation reveal a missing requirement, hidden dependency, or broader slice?
- Are HITL blockers now cleared or newly introduced?
- Is the next `[AFK]` task still safe to run from existing artifacts?
- Does the handoff point to the latest evidence and next action?

If the answer changes the plan, pause implementation and update PRD/OpenSpec/tasks before continuing.

### Receiving Review Feedback

Treat reviewer output as technical evidence, not orders:

1. Read all feedback before reacting.
2. Restate or clarify unclear items before implementation.
3. Verify each item against the codebase, PRD, OpenSpec artifacts, and tests.
4. Push back with concise technical reasoning when feedback is wrong, breaks existing behavior, violates YAGNI, or conflicts with recorded decisions.
5. Implement accepted items one at a time, testing each fix before moving on.

For external reviewer feedback, be especially skeptical of suggestions that add unused "professional" features. Check whether the capability is actually required by the PRD, OpenSpec scenarios, or existing callers before expanding scope.

### Two-Stage Review

For subagent-generated plans, code, or reports:

1. Check compliance against PRD, OpenSpec artifacts, and task scope.
2. Check quality: correctness, maintainability, security, performance, and evidence.

Critical issues block progress until fixed, accepted, or converted into explicit follow-up work.

Request review early and often:

- after each delegated implementation slice
- after completing a major feature
- before validating or finishing a change
- when stuck and a fresh perspective can reduce risk

The reviewer prompt must include the implemented scope, requirements or task text, relevant artifact paths, and the diff range or file set to review. Never ask a reviewer to infer scope from chat history.

### Red Flags

Never:

- dispatch a subagent with vague scope or missing artifact paths
- let subagents inherit unstated session context
- ignore subagent questions or blockers
- retry the same failed prompt unchanged
- accept "close enough" when spec compliance found gaps
- accept reviewer feedback without verifying it
- proceed with unfixed critical or important review findings
- start code quality review before spec compliance passes
- move to the next implementation slice while either review has open issues
- use parallel implementation agents as a substitute for clear task boundaries

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
