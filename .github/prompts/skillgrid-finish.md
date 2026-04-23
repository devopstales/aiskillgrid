---
description: Close the change: archive or sync specs, git hygiene, ship checklist, PR
allowed-tools: Read, Write, Glob, Grep, Bash, Task
argument-hint: "[change-id or branch name]"
---

<objective>

You are executing **`/skillgrid-finish`** (SHIP phase) for the Skillgrid workflow.

</objective>

<process>

## Steps

1. **Archive** — Complete and archive the OpenSpec change with `openspec-archive-change` per your merge and branch policy.
2. **Sync specs (optional)** — If the team promotes delta specs to mainline specs *without* archiving, use `openspec-sync-specs` instead of or in addition to archive.
3. **Pull request** — Open or update a PR with a clear description, risk notes, and links to the change; keep commits atomic and reviewable (`git-workflow-and-versioning`).
4. **CI / gates** — Align with `ci-cd-and-automation`: required checks, feature flags, deployment hooks as used by the project.
5. **Deprecation** — If this change supersedes old paths, use `deprecation-and-migration` for timelines and cleanup.
6. **Ship** — When deploying, run through `shipping-and-launch` (rollout, monitoring, rollback).
7. **Documentation** — Update ADRs, API docs, and inline *why* via `documentation-and-adrs` when behavior or contracts shipped with this change.

## Skills to read and follow

- `.agents/skills/karpathy-guidelines/SKILL.md` — crisp PR text and honest risk notes.
- `.agents/skills/openspec-archive-change/SKILL.md` — complete and archive the change.
- `.agents/skills/openspec-sync-specs/SKILL.md` — promote delta specs without archiving, if the flow needs it.
- `.agents/skills/git-workflow-and-versioning/SKILL.md` — trunk-style workflow, atomic commits, small changes.
- `.agents/skills/documentation-and-adrs/SKILL.md` — record decisions and public surface changes for the next reader.
- `.agents/skills/ci-cd-and-automation/SKILL.md` — pipelines, feature flags, quality gates.
- `.agents/skills/deprecation-and-migration/SKILL.md` — sunset paths, migrations, dead code.
- `.agents/skills/shipping-and-launch/SKILL.md` — pre-launch checks, rollouts, rollback, monitoring.

## Notes

- Inspect the repo with tools; do not assume stack or layout.
- If OpenSpec or SDD modes are unclear, ask once, then follow existing `openspec/` or repo persistence conventions.

</process>
