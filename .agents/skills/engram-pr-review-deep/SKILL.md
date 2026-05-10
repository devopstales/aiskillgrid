---
name: engram-pr-review-deep
description: >
  Reviewer workflow: deep technical review protocol for Engram pull requests before merge.
  Trigger: Reviewing any external or internal contribution as a reviewer (not the author integrating review comments — use receiving-code-review).
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

## When to Use

Use this skill when:
- Evaluating PRs from contributors
- Reviewing risky refactors
- Deciding merge vs request-changes

Do **not** use this skill when **you** are the PR author responding to review threads — use **`receiving-code-review`**.

---

## Review Protocol

1. Read full diff, not only summary.
2. Run relevant tests locally.
3. Validate API/contracts and migration safety.
4. Check docs against implementation.
5. Flag commit hygiene violations.

---

## Merge Gate

Merge only when:
- checks are green
- risk is understood
- blockers are resolved
- scope is coherent

Otherwise request changes with actionable items.
