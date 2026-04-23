---
name: /skillgrid-test
id: skillgrid-test
category: Workflow
description: Prove behavior: automated tests, E2E, browser DevTools
---

You are executing **`/skillgrid-test`** (VERIFY phase) for the Skillgrid workflow.

## Steps

1. **Automated tests** — Run or add tests that match the change; for bugs, prefer a failing test first, then the fix.
2. **E2E** — Use `e2e-testing` and `e2e-runner` for journeys and critical paths; quarantine or stabilize flaky cases per team practice.
3. **Patterns** — Apply `testing-patterns` for layers beyond E2E (unit, integration, mocks) as appropriate to the stack.
4. **Browser** — For UI, use `browser-testing-with-devtools` (or your DevTools/browser MCP): DOM, console, network, performance as needed.
5. **Debug** — When something fails, follow `debugging-and-error-recovery`: reproduce, localize, reduce, fix, add a guard (test or assertion).

## Skills to read and follow

- `.agents/skills/e2e-testing/SKILL.md` — end-to-end test design and implementation.
- `.agents/skills/e2e-runner/SKILL.md` — run and troubleshoot E2E suites.
- `.agents/skills/testing-patterns/SKILL.md` — general testing patterns beyond E2E.
- `.agents/skills/browser-testing-with-devtools/SKILL.md` — DOM, console, network, performance.
- `.agents/skills/debugging-and-error-recovery/SKILL.md` — systematic debugging.

## Notes

- Inspect the repo with tools; do not assume stack or layout.
- If OpenSpec or SDD modes are unclear, ask once, then follow existing `openspec/` or repo persistence conventions.
