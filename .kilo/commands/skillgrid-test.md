---
name: /skillgrid-test
id: skillgrid-test
category: Workflow
description: Prove behavior: automated tests, E2E, browser DevTools
allowed-tools: Read, Write, Glob, Grep, Bash, Task
argument-hint: "[area, issue id/key (per .skillgrid/config.json), or failing test]"
---

<objective>

You are executing **`/skillgrid-test`** for the Skillgrid workflow.

Run quality gates and collect verification evidence for the active PRD/spec slice.

**Status on exit:** no automatic PRD status change; validation owns `devdone`.

</objective>

<process>

## Shared Skillgrid Skills

Before acting, load only the skills needed for the phase:

- `.agents/skills/skillgrid-questioning/SKILL.md` — ask only blocking questions and record answers.
- `.agents/skills/skillgrid-codebase-map/SKILL.md` — map repo structure, graphify output, tests, and conventions.
- `.agents/skills/ccc/SKILL.md` — CocoIndex Code: `ccc init`, `ccc index`, `ccc search`; optional MCP `cocoindex-code`.
- `.agents/skills/references/indexing-and-memory.md` — Engram, graphify (`graphify-out/`, `graphify update .`), **ccc**, and MCP memory ordering.
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

- `skillgrid-codebase-map`
- `ccc`
- `skillgrid-vertical-slices`
- `skillgrid-filesystem-handoff`
- test skills such as `test-driven-development`, `testing-patterns`, `e2e-testing`, `browser-testing-with-devtools`
- security scanners such as `semgrep-security` or `trivy-security` when relevant

## Steps

1. Select the PRD/change or infer it from the handoff.
2. Map PRD success criteria and OpenSpec scenarios to test evidence.
3. Confirm new or changed behavior has tests through public interfaces; bug fixes should include a reproduction test that failed before the fix.
4. Run focused unit, integration, E2E, browser, build, lint, type, and security checks as appropriate.
5. Triage failures by root cause; do not mask failures or accept skipped tests as success.
6. Store long evidence under `.skillgrid/tasks/research/<change-id>/` and summarize in the handoff.
7. Recommend `/skillgrid-apply` for fixes or `/skillgrid-validate` when ready.

## Completion Report

Report checks run, pass/fail status, evidence paths, blockers, and next command.

</process>
