
# Behavioral guidelines

## 1. Read CONTEXT.md first

**Always** begin any `/sdd-*` task (except `/sdd-init`) by reading `.skillgrid/project/CONTEXT.md`:
- If the file exists, note any relevant `Glossary` terms, `Assumptions`, or `Success Criteria`.
- If a term the user uses conflicts with the glossary, **flag the conflict immediately** using a code-fence block:

```marked
**Context Conflict**
> `CONTEXT.md` defines 'cancellation' as X, but you said Y — which one is correct?
```

- Never contradict the glossary without updating the file first. To update, propose the change and ask for confirmation before writing.
- If `CONTEXT.md` doesn’t exist, propose creating it after the first domain term is clarified.

## 2. DRY (Don’t Repeat Yourself)

Never ask the user to re‑explain a term that is already defined in `CONTEXT.md`. Instead:

- Quote the existing definition.
- Ask a precise follow‑up question that builds on that definition.

## 3. Think Before Coding

**Don't assume. Don't hide confusion. Surface tradeoffs.**

Before implementing:
- State your assumptions explicitly. If uncertain, ask.
- If multiple interpretations exist, present them - don't pick silently.
- If a simpler approach exists, say so. Push back when warranted.
- If something is unclear, stop. Name what's confusing. Ask.

## 4. Simplicity First

**Minimum code that solves the problem. Nothing speculative.**

- No features beyond what was asked.
- No abstractions for single-use code.
- No "flexibility" or "configurability" that wasn't requested.
- No error handling for impossible scenarios.
- If you write 200 lines and it could be 50, rewrite it.

Ask yourself: "Would a senior engineer say this is overcomplicated?" If yes, simplify.

## 5. Surgical Changes

**Touch only what you must. Clean up only your own mess.**

When editing existing code:
- Don't "improve" adjacent code, comments, or formatting.
- Don't refactor things that aren't broken.
- Match existing style, even if you'd do it differently.
- If you notice unrelated dead code, mention it - don't delete it.

When your changes create orphans:
- Remove imports/variables/functions that YOUR changes made unused.
- Don't remove pre-existing dead code unless asked.

The test: Every changed line should trace directly to the user's request.

## 6. Goal-Driven Execution

**Define success criteria. Loop until verified.**

Transform tasks into verifiable goals:
- "Add validation" → "Write tests for invalid inputs, then make them pass"
- "Fix the bug" → "Write a test that reproduces it, then make it pass"
- "Refactor X" → "Ensure tests pass before and after"

For multi-step tasks, state a brief plan:
```
1. [Step] → verify: [check]
2. [Step] → verify: [check]
3. [Step] → verify: [check]
```

Strong success criteria let you loop independently. Weak criteria ("make it work") require constant clarification.

## 7. Project Context (.skillgrid/project/CONTEXT.md)

The project’s permanent domain glossary, assumptions, and success criteria are stored in `.skillgrid/project/CONTEXT.md`. This file is version‑controlled and shared across all SDD phases.

### Format of CONTEXT.md

```marked
# Project Context

## Domain
[One‑line description of the main domain]

## Glossary
- **Term**: definition (as agreed by the user)

## Assumptions
- [Any non‑obvious assumptions the agent should remember]

## Success Criteria
- [What “done” looks like from a domain perspective]
```

### How CONTEXT.md is used

| Command | Reads CONTEXT.md? | Writes CONTEXT.md? |
| ------- | ----------------- | ------------------ |
| `/sdd-init` | No | No |
| `/sdd-clarify` | Yes (at start) | Yes (when terms, assumptions, or success criteria are resolved) |
| `/sdd-plan` | Yes | Only if user requests an update |
| `/sdd-apply` | Yes | Only if user requests an update |
| `/sdd-verify` | Yes | No (but may quote success criteria) |
| `/sdd-archive` | Optional (for final docs) | No |

### Updating CONTEXT.md

- When to update: Only when the user explicitly agrees to a change, or after a `/sdd-clarify` session resolves a new term.
- Who can update: The agent, after obtaining user confirmation.
- How to update: Use the configured file-write tool (e.g. `write_file`) to replace `.skillgrid/project/CONTEXT.md` with the updated content.
- Never delete existing entries without asking the user.

## 8. SDD Commands Overview

Your workflow uses the following commands. The orchestrator pattern (if enabled) delegates execution to sub‑agents, but the **CONTEXT.md rule applies to every agent** irrespective of delegation.

| Command | Purpose | Must read CONTEXT.md? |
|---------|---------|----------------------|
| `/sdd-init` | Initialise a new change | No |
| `/sdd-clarify` | Clarify requirements and update CONTEXT.md | Yes (and writes) |
| `/sdd-plan` | Create a plan from the proposal | Yes |
| `/sdd-apply` | Implement tasks | Yes |
| `/sdd-verify` | Verify against success criteria | Yes (reads success criteria) |
| `/sdd-archive` | Archive the change | No (but may reference) |

## 9. Sub-agent context protocol (if orchestrator)

- **Reading context**: The orchestrator (or the agent itself) reads `.skillgrid/project/CONTEXT.md` and passes relevant definitions to the sub‑agent in the prompt.
- **Writing context**: Sub‑agents MAY propose updates to `CONTEXT.md` but MUST return those proposals to the orchestrator for user confirmation. Sub‑agents do NOT write to `CONTEXT.md` directly.
- **Skill loading**: If you use skills, the orchestrator resolves skill paths once per session and passes them to sub‑agents. Sub‑agents do not search for the skill registry themselves.

## 10. Tool & Integration Notes

### GitNexus (if used)

The v1 `AGENTS.md` includes extensive GitNexus integration. For v2, you can keep or remove that section depending on your toolchain. If you keep it, **prefix every GitNexus rule with**: *“After reading CONTEXT.md, …”*

### Engram (if used)

If you use engram for cross‑session memory, **store the location of `CONTEXT.md`** as a key‑value pair:

```
mem_save(key: "context_path", value: ".skillgrid/project/CONTEXT.md", project: "<project>")
```

Then other agents can retrieve it without hard‑coding paths.

## 11. Conflict Resolution Workflow

When `CONTEXT.md` and a user instruction disagree:

1. **Pause execution**.
2. **Quote the conflicting definitions** from `CONTEXT.md`.
3. **Ask the user which one should win**.
4. **If the user chooses the instruction**, propose an update to `CONTEXT.md`.
5. **Wait for confirmation** before overwriting the file.
6. **Proceed** only after the conflict is resolved.

This ensures the glossary remains the single source of truth without silently overriding user intent.

## 12. Final Self-Check

Before responding to any user request that involves domain terms, run this mental check:

- [ ] Have I read `.skillgrid/project/CONTEXT.md`?
- [ ] Do any terms the user used conflict with the glossary?
- [ ] If yes, did I pause and ask for clarification?
- [ ] Am I about to implement something that violates an `Assumption` or `Success Criterion`?

If any answer is **no**, stop and fix it.
