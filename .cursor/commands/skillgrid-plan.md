---
name: /skillgrid-plan
id: skillgrid-plan
category: Workflow
description: Plan — PRD first, then OpenSpec change and CLI-driven artifacts
allowed-tools: Read, Write, Glob, Grep, Bash, Task
argument-hint: "[feature name, slug, or PRD title]"
---

<objective>

You are executing **`/skillgrid-plan`** for the Skillgrid workflow.

Create or update the PRD first, then create or refresh the OpenSpec change and required planning artifacts.

**Status on exit:** set the PRD and `INDEX.md` entry to `draft` after successful plan creation.

</objective>

<process>

## Shared Skillgrid Skills

Before acting, load only the skills needed for the phase:

- `.agents/skills/skillgrid-questioning/SKILL.md` — ask only blocking questions and record answers.
- `.agents/skills/skillgrid-codebase-map/SKILL.md` — map repo structure, graphify output, tests, and conventions.
- `.agents/skills/skillgrid-parallel-research/SKILL.md` — coordinate external research and long evidence capture.
- `.agents/skills/skillgrid-subagent-orchestration/SKILL.md` — dispatch bounded subagents with handoff paths and two-stage review.
- `.agents/skills/skillgrid-prd-artifacts/SKILL.md` — PRD numbering, `INDEX.md`, title blocks, and status lifecycle.
- `.agents/skills/skillgrid-spec-artifacts/SKILL.md` — PRD-to-OpenSpec artifacts and validation.
- `.agents/skills/skillgrid-vertical-slices/SKILL.md` — shippable slices, `[HITL]` / `[AFK]`, and testable increments.
- `.agents/skills/skillgrid-ui-design-artifacts/SKILL.md` — UI decisions, previews, `DESIGN.md`, and OpenSpec design constraints.
- `.agents/skills/skillgrid-issue-creation/SKILL.md` — local/GitHub/GitLab/Jira issue behavior from `.skillgrid/config.json`.
- `.agents/skills/skillgrid-hybrid-persistence/SKILL.md` — disk plus Engram persistence.
- `.agents/skills/skillgrid-filesystem-handoff/SKILL.md` — `context_<change-id>.md` and `research/<change-id>/`.
- `.agents/skills/skillgrid-openspec-config/SKILL.md` — `openspec/config.yaml` overlay rules.
- `.agents/skills/skillgrid-project-docs/SKILL.md` — `DESIGN.md` and `.skillgrid/project/*` docs.
- `.agents/skills/skillgrid-checkpoints/SKILL.md` — `.skillgrid/tasks/checkpoints.log`.


## Phase-Specific Skills

Load these first for this command:

- `skillgrid-questioning`
- `skillgrid-prd-artifacts`
- `skillgrid-spec-artifacts`
- `skillgrid-vertical-slices`
- `skillgrid-issue-creation`
- `skillgrid-hybrid-persistence`

## Steps

1. Derive or confirm the change id and PRD title; ask only if scope is unclear.
2. Create/update `.skillgrid/prd/PRD<NN>_<slug>.md` and `.skillgrid/prd/INDEX.md`.
3. Split oversized work into ordered PRDs when vertical slices can ship independently.
4. Create/update `.skillgrid/tasks/context_<change-id>.md` for the active change.
5. Use `skillgrid-parallel-research` local templates when planning needs independent repo, docs, or prior-art lanes before specs.
6. Create or refresh `openspec/changes/<change-id>/` artifacts using the OpenSpec CLI and `skillgrid-spec-artifacts`.
7. Ensure `tasks.md` is executable by a fresh agent: concrete file paths when known, no placeholders, one vertical slice at a time, and TDD-first steps for behavioral code.
8. Self-review PRD/OpenSpec/tasks for coverage, contradictions, missing verification, and placeholder text before reporting completion.
9. Create external issues only when `.skillgrid/config.json` configures a remote provider.
10. Save a concise Engram summary when hybrid persistence is active.

## Completion Report

Report PRD path, OpenSpec change path, issue links if any, status set to `draft`, memory saves, and recommended `/skillgrid-breakdown`.

</process>
