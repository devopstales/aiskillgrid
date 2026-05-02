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

**Status on exit:** set the PRD and `INDEX.md` entry to `.skillgrid/config.json` `prdWorkflow.phaseStatusMap.plan` after successful plan creation (default: `draft`).

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

- `skillgrid-questioning`
- `skillgrid-codebase-map`
- `ccc`
- `skillgrid-prd-artifacts`
- `skillgrid-spec-artifacts`
- `skillgrid-vertical-slices`
- `skillgrid-issue-creation`
- `skillgrid-hybrid-persistence`

## Event Log Rule

For any identified Skillgrid change id, create `.skillgrid/tasks/events/` if needed and append short JSONL events to `.skillgrid/tasks/events/<change-id>.jsonl` when this command starts, completes, blocks, skips, dispatches/receives subagents, or changes workflow state. If a delegated agent cannot write, require it to return a suggested event object and append that event from the parent session before advancing.

## Alignment Gate

Before creating or updating PRDs, run a shared-understanding gate for ambiguous work:

- Ask one blocking question at a time.
- Include a recommended answer when the codebase, domain, or user intent supports one.
- Record accepted answers as assumptions, decisions, non-goals, risks, or open questions in durable artifacts.
- Keep unresolved product, domain, visual, credential, destructive, or release choices `[HITL]`.
- Use discuss/assumptions mode when a full interview would be too slow: read the relevant artifacts, state assumptions, and ask the user to correct only what is wrong.

Do not send product alignment into an unattended Build Loop.

## Steps

1. Run the `skillgrid-questioning` intent gate. Continue only when the request is `plan` or a clearly approved continuation from `explore`/`brainstorm`.
2. If the gate finds `apply`, `test`, `validate`, or `finish`, stop and recommend the correct command instead of creating planning artifacts from execution intent.
3. Derive or confirm the change id and PRD title; ask only if scope is unclear.
4. Use decision-tree interview mode until goal, scope, non-goals, acceptance criteria, and HITL/AFK boundaries are explicit or recorded as open questions.
5. Create/update `.skillgrid/prd/PRD<NN>_<slug>.md` and `.skillgrid/prd/INDEX.md`.
6. Split oversized work into ordered PRDs when vertical slices can ship independently.
7. Create/update `.skillgrid/tasks/context_<change-id>.md` for the active change, including the intent gate result and accepted assumptions.
8. Use `skillgrid-parallel-research` local templates when planning needs independent repo, docs, or prior-art lanes before specs.
9. Create or refresh `openspec/changes/<change-id>/` artifacts using the OpenSpec CLI and `skillgrid-spec-artifacts`.
10. Ensure `tasks.md` is executable by a fresh agent: concrete file paths when known, no placeholders, one vertical slice at a time, `[HITL]` / `[AFK]` labels, and TDD-first steps for behavioral code.
11. Self-review PRD/OpenSpec/tasks for coverage, contradictions, missing verification, and placeholder text before reporting completion.
12. Create external issues only when `.skillgrid/config.json` configures a remote provider.
13. Save a concise Engram summary when hybrid persistence is active.

## Completion Report

Report PRD path, OpenSpec change path, issue links if any, configured status set for `plan`, memory saves, and recommended `/skillgrid-breakdown`.

</process>
