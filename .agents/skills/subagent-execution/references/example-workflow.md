# Example workflow

A full two-task run: Task 1 reviews clean, Task 2 gets a Spec ❌ with one fix round, then the final whole-branch review.

```
You: I'm using subagent-execution to execute this plan.

[Setup: scripts/sdd-workspace docs/plans/feature-plan.md — no ledger inside, fresh start]
[Read plan file once: docs/plans/feature-plan.md]
[Create todos for all tasks]

Task 1: hook-install script

[BASE = git rev-parse HEAD]
[task-brief for Task 1; dispatch implementer with brief + report paths]
Implementer: "Before I begin — should the hook be installed at user or system level?"
You: "User level (~/.config/hooks/)"
Implementer:
  - Implemented the hook-install command
  - Wrote + ran 3 tests: 3/3 passing
  - Self-review: clean
  - Committed
  - Report written to task-1-report.md: status DONE
[review-package PLAN_FILE BASE HEAD; dispatch task reviewer with the printed path]
Task reviewer: Spec ✅ — all requirements met, nothing extra. Task quality: Approved.
[Ledger: Task 1: complete (commits a1b2c3d..d4e5f6a, review clean)]

Task 2: recovery modes

[task-brief for Task 2; dispatch implementer]
Implementer:
  - Implemented verify + repair modes
  - Wrote + ran 5 tests: 5/5 passing
  - Committed
  - Report: status DONE
[review-package; dispatch task reviewer]
Task reviewer: Spec ❌ — Missing: progress-reporting (spec says "report at 10% intervals").
              Task quality: Approved.
[Ledger: Task 2: fix round 1/5 starting — 1 finding open]
[Fix round 1: resume the implementer with the finding verbatim]
Implementer: Added progress-reporting per 10% interval. Re-ran recovery.test.js — 6/6 passing.
             Fix report appended to task-2-report.md.
[review-package PLAN_FILE FIX_BASE HEAD; dispatch scoped re-review]
Re-reviewer: Missing progress-reporting — ADDRESSED (src/recovery.js:41). New breakage: none.
             Verdict: all findings addressed.
[Ledger: Task 2: fix round 1/5 (1 addressed, 0 open; commits d4e5f6a..b7c8d9e)]
[Ledger: Task 2: complete (commits d4e5f6a..b7c8d9e, review clean)]

[After all tasks: review-package PLAN_FILE MERGE_BASE HEAD; dispatch whole-branch reviewer, most capable model]
Reviewer: No issues.
[Ledger: feature-plan: complete. Rulings: none]

[Finish: no Ruling: lines to surface; rm -rf $WORKSPACE; mem_save the ledger]

Done. Ready for finishing.
```
