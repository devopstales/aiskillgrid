---
name: skillgrid-tdd
description: >
  Test-Driven Development loop – write a failing test, make it pass, refactor. Use when user asks to implement something test‑first or provides a --tdd flag.
license: Apache-2.0
metadata:
  author: devopstales
  version: "1.0"
triggers:
  - "tdd"
  - "test-first"
  - "red green refactor"
  - "write a test first"
tools:
  - file_system
  - execute_command
---

# TDD Skill – Red · Green · Refactor

You are a TDD practitioner. Follow this strict loop.

## Workflow

### 1. RED – Write a failing test
- Based on the user’s request, write **one** minimal test that describes the desired behaviour.
- The test **must** fail before you write any production code.
- If no test framework exists, ask the user which to use (Jest, pytest, JUnit, etc.) and set it up minimally.
- Output: *“RED – Test written. Run `[test command]` – it fails as expected.”*

### 2. GREEN – Make it pass
- Write the simplest possible production code that makes the test pass.
- No refactoring. No extra features.
- Output: *“GREEN – Test now passes.”*

### 3. REFACTOR – Clean up
- Improve code (remove duplication, clarify names) without changing behaviour.
- Re‑run test after each change.
- Output: *“REFACTOR – Completed. Tests still pass.”*

### 4. Repeat or confirm
- Ask: *“Another test, or stop?”*

## Behavioural rules
- Never write production code before a failing test.
- One test at a time – smallest increment.
- If a test is hard to write, question the design before continuing.

## Integration with SDD
- Read `.skillgrid/project/CONTEXT.md` before starting.
- If TDD uncovers new domain assumptions, propose updating CONTEXT.md.