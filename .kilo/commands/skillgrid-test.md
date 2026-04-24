---
name: /skillgrid-test
id: skillgrid-test
category: Workflow
description: Prove behavior: automated tests, E2E, browser DevTools
allowed-tools: Read, Write, Glob, Grep, Bash, Task
argument-hint: "[area, ticket id, or failing test]"
---

<objective>

You are executing **`/skillgrid-test`** (VERIFY phase) for the Skillgrid workflow.

</objective>

<process>

## Steps

1. **Automated tests** — Run or add tests that match the change; for bugs, prefer a failing test first, then the fix.
2. **E2E** — Cover critical user journeys; quarantine or stabilize flaky cases per team practice.
3. **Layers** — Use unit and integration tests with sensible mocks where the stack expects them.
4. **Browser** — For UI, use DevTools or a browser MCP: DOM, console, network, performance as needed.
5. **Debug** — Reproduce, localize, reduce, fix, then add a guard (test or assertion).

## Practices (inline)

- Tie tests to **success criteria** from the PRD or change artifacts; avoid speculative coverage.
- Prefer the repo’s existing test runner and patterns (`package.json`, `Makefile`, CI config).

## Notes

- Inspect the repo with tools; do not assume stack or layout.
- **Hybrid persistence** is the default (`openspec/` + Engram); align with **`/skillgrid-init`** if layout is unclear.

## Completion report (required)

End with a **Session wrap-up** the user can scan:

1. **What I did** — Bullets: which test commands ran, pass/fail summary, and artifacts (logs, coverage paths) if any.
2. **Token / usage** — If the product shows **input/output tokens**, **context used**, or **session cost** for this turn, report it. If not available, state **`Token usage: not shown in this environment`** (do not guess).
3. **Suggested next command** — **`/skillgrid-review`** for spec/implementation traceability; or **`/skillgrid-apply`** if tests failed and the change needs more work.

</process>
