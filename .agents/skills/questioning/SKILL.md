---
name: questioning
description: "Stress-test a plan, decision, or idea branch by branch before implementation. Use when a request is ambiguous, scope is unclear, or you need to clarify intent before designing or coding."
license: MIT
metadata:
  author: devopstales
  part-of: Skillgrid
  version: "2.0"
---

# Questioning

Clarify intent before acting. Map decisions as a tree, ask them in frontier order, stop only on user approval.

## Iron Law

```
NO IMPLEMENTATION UNTIL THE USER CONFIRMS SHARED UNDERSTANDING
```

**Violating the letter of this gate is violating its spirit.**

Implementation means: writing or editing code, scaffolding, creating branches or worktrees, editing specs / proposals / tasks, dispatching implementer sub-agents, running any write-shaped command. Preparing to implement IS implementing.

Allowed before approval: reading code, docs, specs, git history, searching memory or code index, asking questions. Reading ≠ implementing.

## Core Model

- **Design tree:** each decision branches into the decisions hanging off it.
- **Frontier:** decisions whose prerequisites are settled — the only questions you can honestly ask now.
- **Round:** one frontier, asked in full, answered in full.

Every question ships with a recommendation so the user can answer by number:

```
Q1 — **<title>**: <body, may include choices>

  Recommendation: <your recommended answer>
```

Rules: ask the whole frontier; never put two dependent questions in the same round.

## Classify First

State it out loud — "this looks **bounded**, so I'll present a short design in chat" — so the user can override:

- **Spike** — feasibility question ("can we…"). Output is an answer, not kept code. State the probe (2–3 sentences), get a nod, investigate cheaply.
- **Bounded** — small change to existing code (flag, endpoint, one-file fix). Ask what matters, present a short design in chat, STOP, wait for approval.
- **Architectural** — new subsystem, new project, interface change. Questions → 2–3 approaches with trade-offs → sectioned design → approval.

When in doubt, take the heavier path. Complexity found mid-task upgrades the path — say so, step up. Never downgrade mid-task.

## Facts vs Decisions

- **Decisions are the user's.** Never answer your own decisions. Answering your own decision is a gate violation — stop, strike it, re-ask it as a Q with a recommendation.
- **Facts are your job.** If the environment settles it (filesystem, git, docs, code), look it up yourself — don't ask.

## Workflow

1. **Classify** and confirm scope boundary (in/out).
2. **Ask in rounds** by frontier. Each question gets a recommendation.
3. **(Architectural only)** Propose 2–3 approaches with trade-offs; lead with your recommendation; apply YAGNI.
4. **Present design** — scale to complexity (sentences for bounded; ~200 words/section for architectural: architecture, components, data flow, errors, testing).
5. **Confirmation gate** — frontier empty ≠ done. Done = explicit user yes to the presented design (`yes`, `approved`, `proceed`). Answering round questions is NOT approval. Silence, "looks good but…", or a partial nod is NOT approval. A two-sentence design still needs a yes.

## Round Format

```
Subject: {plan / decision / idea}
Classification: {spike | bounded | architectural} — why
Scope: {in / out}

---

Round 1 — frontier is {N} decisions:

Q1 — **<title>**: <body>

  Recommendation: <answer>

---

Reply by number (e.g. "1 yes, 2 second option").
```

Recompute the frontier after each answer; ask the next round.

## Example

Input: *"Let customers save checkout drafts."*

```
Subject: Persist checkout drafts
Classification: architectural — new storage, lifecycle, conflicts
Scope: in — save/recover; out — reminders

Round 1 — frontier is 2 decisions:

Q1 — **Where does draft state live?**: backend-only or local + backend?
  Recommendation: backend-only. Local-only loses drafts on device switch.

Q2 — **When is a draft discarded?**: on completion, timeout, or manual only?
  Recommendation: timeout (30 days) + manual delete; keep after completion for re-entry.
```

Answer `1 backend-only, 2 timeout only` unblocks conflict-resolution next round.

## Red Flags — STOP and return to questioning

| Thought | Reality |
|---|---|
| "Too simple for a design" | Simple = short design, not no design. Present it, then STOP and wait. |
| "I'll start while they read" | Preparation is implementation. STOP. Wait for yes. |
| "Questions answered, so I have approval" | Answers ≠ approval. Approval is a yes to the design. STOP. Present design, wait. |
| "I'll set up the branch / spike while waiting" | Setup is implementation. STOP. |
| "User is slow, they can revert" | Revert ≠ approval. STOP. Wait. |
| "I know this domain, so it's bounded" | Bounded measures the repo, not your familiarity. STOP. Reclassify. |
| "Spike works, I'll keep the code" | Spike output is an answer. Keeping code is a new request — STOP, reclassify. |

## Verification checklist

Before claiming shared understanding:

- [ ] Classification stated out loud with reason
- [ ] Scope in/out confirmed
- [ ] Every frontier question shipped a recommendation
- [ ] Design presented at the right scale
- [ ] Explicit yes to the design observed (not inferred from answers)
- [ ] Zero implementation actions taken before that yes

Can't check all boxes? Not done. Keep questioning.

## Gotchas

- The frontier is judgment, not a computed graph. If you batched two questions and one should have changed the other, say so and reopen next round.
- If your recommendation disputes the question's framing, tell the user to answer the recommendation.
- Long session = scope problem. Split it and question the pieces.
