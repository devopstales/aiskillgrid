---
name: /skillgrid-validate
id: skillgrid-validate
category: Workflow
description: Combined gate — full review then full security (same as review + security)
allowed-tools: Read, Glob, Grep, Bash, Task, Write
argument-hint: "[change-id or release scope]"
---

<objective>

You are executing **`/skillgrid-validate`** for the Skillgrid workflow.

This command is a **single-turn combined gate**: perform every step in **`/skillgrid-review`**, then every step in **`/skillgrid-security`**, in that order.

</objective>

<process>

## Execution

1. Open and follow **`skillgrid-review.md`** in this repo’s commands directory (same folder as this file, or under `.kilo/commands/`, `.opencode/commands/`, `.github/prompts/`).
2. When review is complete, open and follow **`skillgrid-security.md`** in the same way.

If you cannot read those files, use the **Steps** / **process** sections from both commands as the authoritative checklist (no external skill file paths).

## Practices (combined)

Across both passes: keep **scope and assumptions** explicit; do not treat this single gate as sufficient if the team requires a human sign-off between review and security.

## Verification discipline (completion claims)

**Principle:** Evidence before any success or completion claim. If you have not run the check that proves the claim in this turn, you cannot state that the claim holds.

**Before stating** that tests pass, the build is clean, a bug is fixed, requirements are met, or this gate is satisfied:

1. **Identify** the command or checklist step that would prove it.
2. **Run** it fully (fresh, not assumed from an earlier session or another agent).
3. **Read** the full output, exit code, and failure counts.
4. **Verify** that the output actually supports the claim; if not, report the real state with that evidence.
5. **Only then** make the claim, and cite what you ran.

**Weak claims to avoid** (they are not verification): “should pass,” “probably fine,” “looks correct,” satisfaction or “done” before reading output, trusting another agent’s success report without checking diffs and rerunning checks where needed, substituting linter green for a full build, or treating partial checks as proof of the whole gate.

**Regression / TDD:** A regression test is only credible after a red–green check (failure without the fix, pass with the fix) when that is feasible.

**Delegation:** If work was delegated, independently confirm repository state and verification output; a delegated “success” is not evidence by itself.

## Optional: IDE personas

This command runs **review then security sequentially** in one turn. If your IDE supports **parallel subagents**, you can instead fan out independent reports (e.g. **`skillgrid-spec-verifier`**, **`skillgrid-code-reviewer`**, **`skillgrid-security-auditor`**, and optionally **`skillgrid-test-engineer`**) and merge—see [`.cursor/agents/README.md`](../../.cursor/agents/README.md).

## Notes

- Prefer **`/skillgrid-review`** then **`/skillgrid-security`** as separate invocations when you want clearer phase boundaries or human sign-off between them.
- Inspect the repo with tools; do not assume stack or layout.

## Completion report (required)

End with a **Session wrap-up** the user can scan:

1. **What I did** — Bullets: review + security outcomes, reports written, and whether the change is merge-ready.
2. **Token / usage** — If the product shows **input/output tokens**, **context used**, or **session cost** for this turn, report it. If not available, state **`Token usage: not shown in this environment`** (do not guess).
3. **Suggested next command** — **`/skillgrid-finish`** to archive, sync specs, and prep PR/CI; or **`/skillgrid-apply`** if blockers need more implementation.

</process>
