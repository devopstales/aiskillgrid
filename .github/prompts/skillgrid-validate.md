---
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

If you cannot read those files, use the **Steps** and **Skills to read and follow** sections from both commands as the authoritative checklist.

## Skills (combined)

Load **`/skillgrid-review`** and **`/skillgrid-security`** skill lists together; at minimum keep **`karpathy-guidelines`** in scope for scope and assumption checks across both passes.

## Optional: IDE personas

This command runs **review then security sequentially** in one turn. If your IDE supports **parallel subagents**, you can instead fan out independent reports (e.g. **`skillgrid-spec-verifier`**, **`skillgrid-code-reviewer`**, **`skillgrid-security-auditor`**, and optionally **`skillgrid-test-engineer`**) and merge—see [`.cursor/agents/README.md`](../../.cursor/agents/README.md).

## Notes

- Prefer **`/skillgrid-review`** then **`/skillgrid-security`** as separate invocations when you want clearer phase boundaries or human sign-off between them.
- Inspect the repo with tools; do not assume stack or layout.

</process>
