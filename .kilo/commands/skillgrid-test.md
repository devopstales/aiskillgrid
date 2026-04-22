---
name: /skillgrid-test
id: skillgrid-test
category: Workflow
description: Prove behavior: automated tests, E2E, browser DevTools
---

You are executing **`/skillgrid-test`** for the aiskillgrid workflow.

**Canonical checklist:** `docs/wokflow.md` (section `/skillgrid-test`). Prefer that document if this prompt and the doc diverge.

## Actions

1. Run or add tests that match the change; prefer failing test first when fixing bugs.
2. Use browser/DevTools MCP flows per `browser-testing-with-devtools` when validating UI.

## Skills to read and follow

Load each file below before doing substantive work (read fully or skim per skill length):

- `.agents/skills/e2e-testing/SKILL.md`
- `.agents/skills/e2e-runner/SKILL.md`
- `.agents/skills/testing-patterns/SKILL.md`
- `.agents/skills/browser-testing-with-devtools/SKILL.md`
- `.agents/skills/debugging-and-error-recovery/SKILL.md`

## Notes

- Use tools to inspect the repo; do not assume stack or layout.
- If OpenSpec or SDD modes are unclear, ask once, then pick the path consistent with existing `openspec/` or project docs.
