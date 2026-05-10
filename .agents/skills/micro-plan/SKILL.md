---
name: micro-plan
description: >
  Short operational plan and execution: 3–7 concrete steps with exit criteria for quick work, without a full SDD/OpenSpec change.
  Trigger: Quick fix or small feature; user wants a short plan before acting; "micro-plan", "few steps", "simple plan", "execute a small plan".
  Does not replace sdd-tasks, sdd-design, or openspec-continue-change for tracked changes.
license: Apache-2.0
metadata:
  author: devopstales
  version: "1.0"
---

## Role

You produce a **small, executable plan** (or follow one the user already gave), then drive execution stepwise. You stay lightweight: no full proposal/spec/design stack unless the user explicitly asks for SDD.

## When to use

- Bounded scope: hours to a day, few files or one subsystem.
- User wants clarity before edits: **plan first**, then implement.
- Spikes or chores that do not need `openspec/changes/...` lifecycle.

## When not to use

- Multi-sprint or cross-team change → **`sdd-tasks`**, **`sdd-design`**, **`openspec-continue-change`** (and related `sdd-*` / `openspec-*`) as appropriate.
- Need formal verification report / archive trail → **`sdd-verify`**, **`openspec-verify-change`**.

## Plan format (write this before coding)

Keep **3–7 numbered steps**. Each step must be **one action** (read, edit, run, commit) where possible.

```markdown
## Goal
<one sentence>

## Steps
1. …
2. …
(n ≤ 7)

## Exit criteria
- [ ] … (observable)
- [ ] …

## Risks / unknowns
- …
```

## Execution

1. Show the plan to the user if they did not supply one; get implicit OK or edits.
2. Execute steps in order; after each step, note **done / blocked** in one line.
3. Before closing, check **every exit criterion**; if any fail, say what is left.

## Stop

Stop when exit criteria pass or the user accepts documented residual risk. Do not silently expand scope—open a follow-up or switch to SDD if the work outgrows the micro-plan.

## Related

- **`sdd-tasks`** / **`openspec-continue-change`** — structured task lists for real changes.
- **`skillgrid-skill-registry`** — refresh the hub index after adding skills.
- **`.agents/skills/_shared/SKILL-authoring-template.md`** — canonical layout for new `SKILL.md` files.
