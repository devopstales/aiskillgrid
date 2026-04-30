---
name: git-master
description: Built-in git hygiene skill for atomic commits, branch discipline, PR-ready diffs, conflict handling, and safe release preparation. Use when preparing commits, reviewing git state, splitting work, or making repository history reviewable.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

# Git Master

## When to Use

Use this skill when the work involves:

- checking repository state before edits, review, or handoff
- splitting changes into reviewable units
- creating commits after the user explicitly asks
- preparing a PR or release branch
- resolving merge conflicts or understanding branch divergence
- keeping generated mirrors separate from behavior changes when practical

For general git principles, load `git-workflow-and-versioning`.

## Hard Safety Rules

- Never create a commit unless the user explicitly asks.
- Never push, force-push, reset hard, checkout away changes, delete branches, or rewrite history unless the user explicitly asks and the risk is understood.
- Never hide unrelated dirty changes. Work around them and preserve them.
- Never commit secrets, `.env` files, credentials, tokens, private keys, or generated build outputs that the repo ignores.
- Never skip hooks unless the user explicitly asks.

## Atomic Commit Workflow

When the user asks for a commit:

1. Inspect git status, full diff, staged diff, and recent commit style.
2. Identify unrelated or user-owned changes and leave them unstaged unless the user says otherwise.
3. Stage only files belonging to the requested logical change.
4. Run relevant verification or clearly report why it was not run.
5. Commit with a concise message that explains why.
6. Re-check status and report remaining uncommitted work.

## Reviewable Change Discipline

Keep changes small:

- One logical purpose per commit or PR.
- Separate generated mirrors from source edits when the project expects reviewers to inspect them independently.
- Separate formatting-only changes from behavior, docs, or tooling changes.
- Prefer Skillgrid checkpoints for intermediate save points when the user has not asked for commits.

## Skillgrid Contract

For Skillgrid:

- `/skillgrid-apply` may create checkpoints but must not commit by default.
- `/skillgrid-finish` owns PR/merge/release preparation.
- `git-master` can advise on branch, diff, and commit hygiene at any phase.
- Record commit hashes, PR links, or rollback notes in `.skillgrid/tasks/context_<change-id>.md` when they become part of release evidence.

## Commands

Use the repository's normal git commands. For commits, follow the platform's safety protocol and pass messages with a heredoc.

```bash
git status --short
git diff
git diff --staged
git log --oneline -5
```

## References

- Core workflow: `git-workflow-and-versioning`
- PR flow: `/skillgrid-finish`
- Checkpoints without commits: `skillgrid-checkpoints`
