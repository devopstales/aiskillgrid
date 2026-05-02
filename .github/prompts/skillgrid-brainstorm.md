---
description: Brainstorm and refine ideas before committing to a plan
allowed-tools: Read, Write, Glob, Grep, WebSearch, WebFetch, Task, Bash
argument-hint: "[idea, question, or problem to refine]"
---

<objective>

You are executing **`/skillgrid-brainstorm`** for the Skillgrid workflow.

Clarify and shape an idea before formal planning, using research, alternatives, and concise user questions.

**Status on exit:** no required PRD status change.

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
- `skillgrid-parallel-research`
- `skillgrid-ui-design-artifacts`
- `skillgrid-project-docs`
- `skillgrid-hybrid-persistence`

## Event Log Rule

For any identified Skillgrid change id, create `.skillgrid/tasks/events/` if needed and append short JSONL events to `.skillgrid/tasks/events/<change-id>.jsonl` when this command starts, completes, blocks, skips, dispatches/receives subagents, or changes workflow state. If a delegated agent cannot write, require it to return a suggested event object and append that event from the parent session before advancing.

## Steps

1. Run the `skillgrid-questioning` intent gate: classify the request as `explore`, `plan`, `apply`, `test`, `validate`, `finish`, or `blocked` before shaping the idea.
2. If the request is actually apply/test/validate/finish work, stop and recommend the correct command unless the user explicitly wants brainstorming first.
3. State the current idea and identify the few decisions that matter.
4. Explore project context before asking questions; if the idea spans independent subsystems, propose decomposition before refinement.
5. Ask one blocking question at a time when decisions depend on each other, offering defaults when safe.
6. Explore 2-3 approaches with tradeoffs and a recommendation; use `skillgrid-parallel-research` local templates when independent research questions exist.
7. For UI work, create or compare directions and record durable choices through `skillgrid-ui-design-artifacts`.
8. Present the accepted direction or design and get approval, even if the final design is short.
9. Capture accepted decisions in PRD drafts, OpenSpec design notes, project docs, or Engram as appropriate.
10. Recommend `/skillgrid-plan` when the scope is ready.

## Completion Report

Report decisions made, alternatives rejected, artifacts updated, unresolved questions, and whether to proceed to `/skillgrid-plan`.

</process>
