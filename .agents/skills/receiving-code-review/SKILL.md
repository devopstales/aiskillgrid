---
name: receiving-code-review
description: >
  Author workflow: triage and integrate review feedback on your PR (threads, fixes, pushes, re-request review).
  Trigger: Review comments exist; "address review", "fix PR comments", "resolve review threads", "respond to reviewer".
  Not for opening a PR (engram-branch-pr / project template) or for preparing the first review pass (requesting-code-review) or for reviewing someone else's PR (gitnexus-pr-review / engram-pr-review-deep).
license: Apache-2.0
metadata:
  author: devopstales
  version: "1.0"
---

## Role

You are the **change author** acting on **existing** review feedback. You do not re-run the reviewer’s full merge assessment unless the user asks; you **close the loop** on comments.

## Workflow

1. **Ingest** — list open review threads (GitHub UI, `gh pr view --comments`, or API); note blockers vs nits vs questions.
2. **Triage** — must-fix before optional polish; ask the user only when a comment is ambiguous or conflicts with product intent.
3. **Implement** — smallest commits that map to threads when helpful; preserve project commit rules (e.g. **`engram-commit-hygiene`** on Engram).
4. **Verify** — run the same checks the project expects before push (tests, linters).
5. **Respond** — reply on each thread: what changed or why not; link commits when useful.
6. **Resolve** — mark conversations resolved when addressed; leave open if still under discussion.
7. **Re-request** — if policy requires: re-request review from prior reviewers after substantive fixes.

## Stop

Stop when blocking comments are addressed or explicitly deferred with written rationale on the thread. Do not open a new PR scope here; spin a follow-up issue if the reviewer uncovered new work.

## Out of scope

- **`requesting-code-review`** — no comments yet; you are still preparing the first review pass.
- **`gitnexus-pr-review`** / **`engram-pr-review-deep`** — those are **reviewer** lenses on the same PR, not the author’s integration pass.
