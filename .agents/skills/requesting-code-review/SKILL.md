---
name: requesting-code-review
description: >
  Author workflow: make a PR easy to review and signal readiness for human review (self-review, context, reviewers).
  Trigger: PR is open or opening soon; "ready for review", "request review", "prepare PR for reviewers", "self-review before review",
  draft → ready, add reviewers. Not for addressing existing comments (use receiving-code-review) or for conducting a review as reviewer (use gitnexus-pr-review / engram-pr-review-deep).
license: Apache-2.0
metadata:
  author: devopstales
  version: "1.0"
---

## Role

You are the **change author** optimizing for **incoming** review. You do not triage or fix review comments here. You do not perform the reviewer’s risk analysis here.

## Preconditions

- **Mechanical PR open** (branch, template, issue linkage, labels, CI): use project skills first — on Engram, **`engram-branch-pr`**. This skill starts once a PR exists or is about to exist and the goal is **reviewer experience**, not repo policy mechanics.

## Workflow

1. **Self-review the diff** — read as a stranger; note ambiguities, dead code, noisy formatting you can fix *before* humans spend time.
2. **Risk map for reviewers** — short bullets: intent, risky areas, what you did *not* change on purpose, how to exercise manually if needed.
3. **PR description** — ensure the body answers: what / why / how to test; link specs or issues already required by the project; add a **Reviewer notes** subsection if the template allows (non-duplicative of mandatory sections).
4. **Readiness** — if the PR is draft, confirm CI and local checks are green before marking ready; only undraft when the user wants visibility.
5. **Request reviewers** — if the team uses GitHub reviewers and the user asked: `gh pr view` / `gh pr edit --add-reviewer` per project norms; do not spam reviewers without instruction.

## Stop

Stop when the PR is in a **reviewable** state and reviewers are requested or the user declines. Do not implement review feedback in this skill.

## Out of scope

- **`receiving-code-review`** — comments already exist.
- **`gitnexus-pr-review`** / **`engram-pr-review-deep`** — you are not the external reviewer.
