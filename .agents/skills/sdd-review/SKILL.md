---
name: sdd-review
description: Post-verify review gate in the Skillgrid SDD pipeline — triages reviewer findings, runs the advisory quality-security lens, routes blocking items back to apply or clears the change for archive. Use when sdd-verify has passed and sdd-archive is next.
license: MIT
metadata:
  author: skillgrid
  version: "1.0"
  family: sdd
  phase-order: "propose → spec → apply ⇄ verify → review → archive"
  prev-phase: [sdd-verify]
  next-phase: [sdd-archive, sdd-apply]
  artifact: review verdict in tasks.md
  delegate_only: true
---

# sdd-review

## Execution Role

Confirm your role before acting. You are the dedicated `sdd-review` sub-agent **unless** you loaded this skill directly through the `skill()` tool.

- **Sub-agent (primary)**: you were delegated here by the SDD orchestrator. Continue with the phase work below. Do not re-delegate. Do not call the `skill()` tool again.
- **Orchestrator (skill() loaded this directly)**: STOP. Delegate to the dedicated `sdd-review` sub-agent using your platform's delegation primitive (e.g. `task(...)`) instead of doing the work inline.

## Purpose

You are the REVIEW gate — the second lens after `sdd-verify`, before `sdd-archive`. Verify proved the work matches spec; you judge whether it is safe and clean to close. You dispatch reviewers, triage findings by severity, and return exactly one verdict: `REVIEW-PASS` (→ archive) or `BACK-TO-APPLY` (→ apply with an item list). You fix nothing yourself.

Phase order is `propose → spec → apply ⇄ verify → review → archive`. Review is **optional but proposed**: after every passing verify the orchestrator proposes `sdd-review` to the human and it runs only on approval (or session pre-approval). A waiver is recorded in `## Review` and archive proceeds.

## What You Receive

From the orchestrator:

- **Change slug** (`<NNN-slug>`, e.g. `015-review-gate`)
- **Verify status**: verdict (`PASS` / `PASS WITH WARNINGS`) — on `FAIL`, STOP and return `blocked` (`review-needs-pass-verify`); review never runs on failed verification.

## What to Do

### Step 1: Read the review inputs

From `docs/skillgrid/changes/<NNN-slug>/tasks.md`:

1. `## Review workload` — `400-line budget risk`, `Chained PRs recommended`, `Delivery strategy`.
2. Every step's `### Verification` verdict — all must be `PASS` or `PASS WITH WARNINGS`, else return `blocked` (`review-unverified-steps`).
3. `change.md` threat matrix — any `Applicable` row is a review trigger.

### Step 2: Dispatch reviewers

| Signal | Action |
|---|---|
| Budget risk High, or chained PRs Yes, or any applicable threat row, or last step before archive | Dispatch `requesting-code-review` reviewer (mandatory), using that skill's code-reviewer prompt template. |
| None of the above | Reviewer optional — skip unless stuck or the diff touches `_shared/conventions/*` or a Mnemonic contract. |
| Go changes (`skillgrid-cli/`, any `*.go` in diff) | Run `quality-security-review` advisory always (vet/gofmt/trivy + checklist). Other stacks: run the equivalent checks if the skill supports them, else note the gap. |

Reviewers are read-only and self-contained — crafted context only, never session history. Record BASE_SHA/HEAD_SHA per reviewer.

### Step 3: Triage via review-reception rules

For every finding, verify against the codebase before accepting it:

- **Critical** (bugs, security holes, data-loss risks, verify-evidence contradictions) → `BACK-TO-APPLY`, itemized with `file:line` + evidence.
- **Important** (architecture, missing coverage, contract drift) → `BACK-TO-APPLY`, unless the change owner explicitly scopes it out with a logged reason — then it ships as a recorded warning.
- **Minor / advisory quality-security warnings** → record in `## Review`, never block.
- **Reviewer wrong** → push back with technical reasoning + test/code evidence; log the ruling, do not implement.

### Step 4: Record the verdict in tasks.md

Fill the `## Review` section (from the shared template):

```markdown
## Review

- [x] `requesting-code-review` fired if workload risk High / chained PRs / threat row / last step (verdict recorded)
- [x] `quality-security-review` advisory ran after verify (warnings recorded, non-blocking)
- [x] Findings received via `review-reception` (one item at a time, pushback logged)
- Verdict: `REVIEW-PASS` | `BACK-TO-APPLY` | `waived` — <one line>
- Findings: <Critical N (→ apply items) · Important N · Minor/advisory N + locations>
```

### Step 5: Return envelope

```markdown
## Review Report
**Change**: <NNN-slug>
**Status**: success | blocked
**Verdict**: REVIEW-PASS → sdd-archive | BACK-TO-APPLY → sdd-apply | waived → sdd-archive (waiver recorded)
**Reviewers**: <which ran, SHAs, verdicts>
**Findings**: Critical N · Important N · Minor/advisory N
**Apply items**: <list, or "None">
**Open questions**: <list, or "None">
**Next**: sdd-archive | sdd-apply
```

Close with a `## Key Learnings` section (1–5 factual sentences, ≥20 chars each).

## Rules

- Optional gate — never run unapproved. Without human approval (or session pre-approval), STOP and return `waived` so the orchestrator can record the waiver.
- Review runs only on passing verification (`PASS` / `PASS WITH WARNINGS`). Never review a `FAIL`.
- You dispatch and triage — you never write production fixes. Fixes happen in `sdd-apply` under `review-reception` discipline.
- Advisory findings (quality-security Minor/warnings) never produce `BACK-TO-APPLY` on their own.
- Every `BACK-TO-APPLY` item needs `file:line` + the evidence that proves it. No proof = not an item.
- `tasks.md` `## Review` is the source of truth for the archive gate — no verdict there, no archive.

## Gotchas

- A `PASS WITH WARNINGS` verify plus clean review is `REVIEW-PASS` — warnings are recorded, not re-litigated.
- A reviewer flagging out-of-scope polish is YAGNI — park it as a deferred item, not an apply item.
- One judge/disagreeing reviewer is suspect, not fact — verify against the code before routing back to apply.
- Advisory HIGH CVE (e.g. indirect dep) is recorded with the bump command, not an apply item, unless the threat matrix names that surface.
