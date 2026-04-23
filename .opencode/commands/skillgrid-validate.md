---
name: /skillgrid-validate
id: skillgrid-validate
category: Workflow
description: Combined gate — full review then full security (same as review + security)
---

You are executing **`/skillgrid-validate`** for the Skillgrid workflow.

This command is a **single-turn combined gate**: perform every step in **`/skillgrid-review`**, then every step in **`/skillgrid-security`**, in that order.

## Execution

1. Open and follow **`skillgrid-review.md`** in this repo’s commands directory (same folder as this file, or under `.kilo/commands/`, `.opencode/commands/`, `.github/prompts/`).
2. When review is complete, open and follow **`skillgrid-security.md`** in the same way.

If you cannot read those files, use the **Steps** and **Skills to read and follow** sections from both commands as the authoritative checklist.

## Notes

- Prefer **`/skillgrid-review`** then **`/skillgrid-security`** as separate invocations when you want clearer phase boundaries or human sign-off between them.
- Inspect the repo with tools; do not assume stack or layout.
