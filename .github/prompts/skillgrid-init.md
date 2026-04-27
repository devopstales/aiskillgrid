---
description: Bootstrap workflow: structure, graphify, hybrid persistence (OpenSpec + Engram)
allowed-tools: Read, Write, Glob, Grep, Bash, Task
argument-hint: "[optional: app purpose, stack, or brownfield notes]"
---

<objective>

You are executing **`/skillgrid-init`** for the Skillgrid workflow.

Bootstrap or align a repository for Skillgrid: config, artifact store, OpenSpec, project docs, design docs, and optional indexing.

**Status on exit:** no change PRD status unless creating initial planning artifacts.

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
- `skillgrid-codebase-map`
- `skillgrid-openspec-config`
- `skillgrid-project-docs`
- `skillgrid-hybrid-persistence`

## Steps

1. Ask only for blocking setup choices: ticketing provider and artifact-store mode.
2. Create or merge `.skillgrid/config.json` with local-first defaults when missing.
3. If artifact store includes OpenSpec, initialize or align `openspec/` and `openspec/config.yaml`.
4. Create or refresh `.skillgrid/project/` docs and root `DESIGN.md` using discovered repo facts.
5. Record durable setup decisions through Engram when available and appropriate.
6. Recommend `/skillgrid-explore` for brownfield mapping or `/skillgrid-plan` for a clear first change.

## Completion Report

Report config choices, created/aligned artifacts, project docs touched, memory saves, and recommended next command.

</process>
