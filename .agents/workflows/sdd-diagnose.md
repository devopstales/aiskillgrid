---
name: sdd-diagnose
description: Structured debugging session to reproduce, isolate, hypothesize, fix, and verify issues. Use when user reports a bug, unexpected behavior, or test failure.
---

<context-awareness>

## Read CONTEXT.md

Before starting, read `.skillgrid/project/CONTEXT.md` if it exists. Note any relevant domain terms, assumptions, or success criteria that might relate to the issue.

If the issue involves a term that conflicts with the glossary, flag it immediately (see `.agents/rules/skillgrid-context-contract.mdc`).

</context-awareness>

<workflow>

## The Diagnostic Loop

Execute the following steps **in order**. Do not skip steps. Ask for user input at each stage before proceeding.

### 1. REPRODUCE
- Get exact steps or code to trigger the issue.
- Attempt to reproduce locally or ask the user for logs/error messages.
- If reproduction fails, ask for more details and loop until success.
- Output: *“Reproduced: [conditions that trigger the issue]”*

### 2. ISOLATE
- Narrow down the cause:
  - Which component/function/module?
  - Which input values cause vs avoid the issue?
  - Can you write a minimal test or script that fails?
- Use binary search (comment out code, split inputs) if possible.
- Output: *“Isolated to [file/function/line range]. Minimal repro: [code snippet]”*

### 3. HYPOTHESIZE
- Propose a root cause hypothesis. Include:
  - What is incorrect (logic, state, assumption, external call).
  - Why it happens (e.g., off-by-one, race condition, missing guard).
  - What would prove the hypothesis (e.g., changing this line fixes it).
- If multiple hypotheses, list them and ask the user which to test first.
- Output: *“Hypothesis: … To confirm, change X to Y and see if issue disappears.”*

### 4. FIX (with verification)
- Implement the minimal fix that resolves the hypothesis.
- **Before writing code**, state the planned change and get user approval.
- Apply the change.
- Re-run the reproduction steps to confirm the issue is gone.
- Output: *“Fix applied. Verified: issue no longer occurs.”*

### 5. REFLECT (optional, but recommended)
- Ask: *“Should this insight be added to CONTEXT.md (e.g., an assumption or gotcha) or captured as an ADR?”*
- If yes, propose the update and wait for user confirmation before writing.

</workflow>

<guidelines>

## Behavioural rules

- **Don't guess** – If you cannot reproduce, ask for more data. Do not propose fixes until you have a confirmed reproduction.
- **One change at a time** – When isolating, change only one variable per attempt.
- **Preserve evidence** – Keep a log of steps tried and outputs.
- **Use existing tests** – If the project has a test suite, try to add a failing test before fixing (Red-Green style). Ask the user if they want that.

</guidelines>

<example>

## Example interaction

User: “The login button does nothing when clicked.”

Agent (Reproduce):
> Can you open the browser console and click the button? What error appears?

User: “`Uncaught TypeError: submitForm is not a function`”

Agent (Isolate):
> That suggests the function is undefined. Let me check `login.js`. On line 42, I see `onclick="submitForm()"` but the function is named `submitLoginForm`. Is that the only button on the page?

User: “Yes.”

Agent (Hypothesize):
> Hypothesis: The function name in the HTML attribute does not match the actual function defined. To confirm, change `submitForm` to `submitLoginForm` in the HTML.

User: “Go ahead.”

Agent (Fix):
> Changed line 42. Now clicking the button shows no error and submits. Verified.
> 
> Reflect: Would you like to add an assumption to `CONTEXT.md`? E.g., *“All event handler names must exactly match declared function names.”`

</example>

<error-handling>

## When the loop fails

If after 5 isolation attempts the cause is still unclear:
- **Pause** and summarise what you’ve eliminated.
- **Ask the user** to provide additional logs, a debugger session, or a different reproduction environment.
- **Do not** propose random fixes.

</error-handling>

<integration>

## Connection to other SDD commands

- After a successful diagnose, the user may run `/sdd-apply` to implement the fix more thoroughly (if the fix was a temporary change).
- If the diagnosis reveals a missing success criterion or domain assumption, update `CONTEXT.md` via `/sdd-clarify` or ask the user.
- If the issue requires an architectural decision, offer to create an ADR.

</integration>