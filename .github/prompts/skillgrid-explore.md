---
description: Explore the problem and repo: OpenSpec list, PRD backfill, .skillgrid/project, AGENTS, GitNexus
allowed-tools: Read, Write, Glob, Grep, Bash, Task
argument-hint: "[optional: topic, change id, or area to explore]"
---

<objective>

You are executing **`/skillgrid-explore`** for the Skillgrid workflow.

Map the problem and repo without implementing product behavior. Exploration may update planning artifacts, PRDs, and project docs.

**Status on exit:** no mandatory status change; backfilled PRDs use `skillgrid-prd-artifacts` status rules.

</objective>

<process>

## Shared Skillgrid Skills

Before acting, load only the skills needed for the phase:

- `.agents/skills/skillgrid-questioning/SKILL.md` — ask only blocking questions and record answers.
- `.agents/skills/skillgrid-codebase-map/SKILL.md` — map repo structure, GitNexus output, tests, and conventions.
- `.agents/skills/ccc/SKILL.md` — CocoIndex Code: `ccc init`, `ccc index`, `ccc search`; optional MCP `cocoindex-code`.
- `.agents/skills/references/indexing-and-memory.md` — Engram, GitNexus (`.gitnexus/`, `npx -y gitnexus@1.3.11 analyze`), **ccc**, and MCP memory ordering.
- `.agents/skills/skillgrid-parallel-research/SKILL.md` — coordinate external research and long evidence capture.
- `.agents/skills/skillgrid-subagent-orchestration/SKILL.md` — dispatch bounded subagents with handoff paths and two-stage review.
- `.agents/skills/skillgrid-import-artifacts/SKILL.md` — import existing PRDs and OpenSpec changes into canonical PRDs.
- `.agents/skills/skillgrid-prd-artifacts/SKILL.md` — PRD numbering, `INDEX.md`, title blocks, and status lifecycle.
- `.agents/skills/skillgrid-spec-artifacts/SKILL.md` — PRD-to-OpenSpec artifacts and validation.
- `.agents/skills/skillgrid-vertical-slices/SKILL.md` — shippable slices, `[HITL]` / `[AFK]`, and testable increments.
- `.agents/skills/skillgrid-ui-design-artifacts/SKILL.md` — UI decisions, previews, `DESIGN.md`, and OpenSpec design constraints.
- `.agents/skills/skillgrid-issue-creation/SKILL.md` — local/GitHub/GitLab/Jira issue behavior from `.skillgrid/config.json`.
- `.agents/skills/skillgrid-hybrid-persistence/SKILL.md` — disk plus Engram persistence.
- `.agents/skills/skillgrid-filesystem-handoff/SKILL.md` — `context_<change-id>.md`, `events/<change-id>.jsonl`, and `research/<change-id>/`.
- `.agents/skills/skillgrid-openspec-config/SKILL.md` — `openspec/config.yaml` overlay rules.
- `.agents/skills/skillgrid-project-docs/SKILL.md` — `DESIGN.md` and `.skillgrid/project/*` docs.
- `.agents/skills/skillgrid-checkpoints/SKILL.md` — `.skillgrid/tasks/checkpoints.log`.


## Phase-Specific Skills

Load these first for this command:

- `skillgrid-codebase-map`
- `ccc`
- `skillgrid-parallel-research`
- `skillgrid-import-artifacts`
- `skillgrid-prd-artifacts`
- `skillgrid-project-docs`
- `skillgrid-filesystem-handoff`

## Event Log Rule

For any identified Skillgrid change id, create `.skillgrid/tasks/events/` if needed and append short JSONL events to `.skillgrid/tasks/events/<change-id>.jsonl` when this command starts, completes, blocks, skips, dispatches/receives subagents, or changes workflow state. If a delegated agent cannot write, require it to return a suggested event object and append that event from the parent session before advancing.

## Steps

1. Select the topic, change, or repo area to explore; ask only if the target is unclear.
2. Inventory OpenSpec changes, legacy PRDs, canonical PRDs, project docs, handoff files, tests, **GitNexus / `.gitnexus/`** when present, and use **`ccc search`** when the CocoIndex index exists—follow **`ccc`**, **`skillgrid-codebase-map`**, and **`indexing-and-memory`** skills.
3. If existing PRDs or OpenSpec changes lack canonical `.skillgrid/prd/` coverage, automatically invoke `skillgrid-import-artifacts` import/backfill behavior.
4. Import root `prd/`, `docs/PRD/`, or `docs/prd/` PRDs into `.skillgrid/prd/` when unambiguous; report ambiguous matches instead of silently merging them.
5. Use `skillgrid-parallel-research` local templates for independent repo, docs, or web research lanes; spill long output to `.skillgrid/tasks/research/<change-id>/`.
6. Refresh `.skillgrid/project/*` and `DESIGN.md` only for stable discovered facts.
7. Stop when the scope is clear enough to recommend `/skillgrid-brainstorm` or `/skillgrid-plan`.

## Completion Report

Summarize explored areas, artifacts inspected or updated, important findings, open questions, and the recommended next command.

</process>
