---
name: /skillgrid-help
id: skillgrid-help
category: Workflow
description: Explain Skillgrid workflow commands, artifacts, and next steps
allowed-tools: Read, Glob, Grep
argument-hint: "[optional: phase, command, artifact, or current-state]"
---

<objective>

You are executing **`/skillgrid-help`** for the Skillgrid workflow.

Explain how to use Skillgrid commands, phases, artifacts, and next steps in concise user-facing language.

**Status on exit:** no artifact or PRD status changes.

</objective>

<process>

## Read First

Use these files as the source of truth:

- `docs/workflow.md`
- `docs/commands.md`
- `docs/skills.md` when the user asks about reusable skills
- `.skillgrid/config.json` if present

## Modes

1. If no argument is provided, show the short path:
   - `/skillgrid-init`
   - `/skillgrid-explore`
   - `/skillgrid-brainstorm`
   - `/skillgrid-plan`
   - `/skillgrid-breakdown`
   - `/skillgrid-apply`
   - `/skillgrid-test`
   - `/skillgrid-security` optional
   - `/skillgrid-validate`
   - `/skillgrid-finish`
2. If the argument names a command, explain that command's purpose, inputs, outputs, artifacts, and likely next command.
3. If the argument names a phase, explain the phase goal, commands, status behavior, and artifacts.
4. If the argument names an artifact, explain where it lives, who owns it, and which commands update it.
5. If the argument is `current-state`, inspect `.skillgrid/prd/INDEX.md`, `.skillgrid/tasks/context_*.md`, and `openspec/changes/` if present, then recommend the next command.

## Response Style

- Keep output concise and user-facing.
- Prefer plain workflow guidance over internal implementation details.
- Mention missing setup only when it changes the recommended next command.
- Do not modify files.

## Completion Report

Give the requested explanation and, when possible, one recommended next command.

</process>
