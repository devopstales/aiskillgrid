---
description: Explore the repo—current architecture, DESIGN.md from implementation, user Q&A, ADRs—plus OpenSpec, PRDs, project docs, GitNexus
allowed-tools: Read, Write, Glob, Grep, Bash, Task
argument-hint: "[optional: topic, change id, or area to explore]"
---

<objective>

You are executing **`/skillgrid-explore`** for the Skillgrid workflow.

Map the problem and repo without implementing product behavior. Exploration may update planning artifacts, PRDs, and project docs.

Ground root **`DESIGN.md` in the current implemented design** (themes, tokens, components, CSS variables)—**not** from **`/skillgrid-init`** blank scaffolding or template-only placeholders. Treat **`/skillgrid-init`** output as provisional until explore replaces it with code-backed facts.

Reconcile **`.skillgrid/project/ARCHITECTURE.md`** (and related project docs) with the **actual** system: boundaries, dependencies, and tacit decisions visible in code and tests. Use **`skillgrid-questioning`** to ask the user about ambiguity, history, and trade-offs the repo cannot answer. **ADR gate:** for durable architectural decisions, follow **`documentation-and-adrs`**, but **do not** create new **`.skillgrid/adr/NNNN-*.md`** files or add new **Durable decisions** links in **`ARCHITECTURE.md`** until the user **explicitly accepts** a concrete ADR draft (show the draft in chat or in the handoff file first). After acceptance, write the MADR and links. If the user defers or rejects, record candidate ADRs only in the completion report or **`.skillgrid/tasks/research/<change-id>/`**—no new ADR files.

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
- `.agents/skills/documentation-and-adrs/SKILL.md` — MADR ADRs under `.skillgrid/adr/`, decision lifecycle, links from `ARCHITECTURE.md`.
- `.agents/skills/skillgrid-checkpoints/SKILL.md` — `.skillgrid/tasks/checkpoints.log`.


## Phase-Specific Skills

Load these first for this command:

- `skillgrid-codebase-map`
- `ccc`
- `skillgrid-parallel-research`
- `skillgrid-import-artifacts`
- `skillgrid-prd-artifacts`
- `skillgrid-project-docs`
- `documentation-and-adrs`
- `skillgrid-questioning`
- `skillgrid-filesystem-handoff`

## Event Log Rule

For any identified Skillgrid change id, create `.skillgrid/tasks/events/` if needed and append short JSONL events to `.skillgrid/tasks/events/<change-id>.jsonl` when this command starts, completes, blocks, skips, dispatches/receives subagents, or changes workflow state. If a delegated agent cannot write, require it to return a suggested event object and append that event from the parent session before advancing.

## Steps

1. Select the topic, change, or repo area to explore; ask only if the target is unclear.
2. Inventory OpenSpec changes, legacy PRDs, canonical PRDs, project docs, handoff files, tests, **GitNexus / `.gitnexus/`** when present, and use **`ccc search`** when the CocoIndex index exists—follow **`ccc`**, **`skillgrid-codebase-map`**, and **`indexing-and-memory`** skills.
3. If existing PRDs or OpenSpec changes lack canonical `.skillgrid/prd/` coverage, automatically invoke `skillgrid-import-artifacts` import/backfill behavior.
4. Import root `prd/`, `docs/PRD/`, or `docs/prd/` PRDs into `.skillgrid/prd/` when unambiguous; report ambiguous matches instead of silently merging them.
5. Use `skillgrid-parallel-research` local templates for independent repo, docs, or web research lanes; spill long output to `.skillgrid/tasks/research/<change-id>/`.
6. **Current architecture:** compare **`.skillgrid/project/ARCHITECTURE.md`** (and stack notes in **`PROJECT.md`**) to the real codebase—modules, boundaries, data flow, and infra. Record drift, missing subsystems, and decisions implied only by code; the doc must reflect **today’s** system, not init-time guesses.
7. **`DESIGN.md` from current design:** follow **`skillgrid-project-docs`** (and **`skillgrid-ui-design-artifacts`** when UI-heavy). **Extract** tokens, typography, color, spacing, and component conventions from **implemented** sources (Tailwind config, CSS variables, theme files, design-system packages). **Do not** copy empty **`/skillgrid-init`** or template-only **`DESIGN.md`** as the final word—replace boilerplate with evidence from the repo.
8. **User questions and ADRs:** use **`skillgrid-questioning`** for blocking or high-impact unknowns (ownership, deprecated paths, product constraints). Persist answers in handoff or **`.skillgrid/project/*`** as appropriate. **ADR files:** present each candidate ADR as a **draft** (context, options, decision) and **stop for explicit user acceptance** before any write. **Only after** the user clearly approves (e.g. “yes, add that ADR”), create **`NNNN-slug.md`** under **`.skillgrid/adr/`** from **`.skillgrid/templates/template-adr.md`**, set MADR **`status: accepted`**, and add the **Durable decisions** link in **`ARCHITECTURE.md`**. If the user has not accepted, **do not** create ADRs—list deferred candidates in the completion report or **`.skillgrid/tasks/research/<change-id>/`** only.
9. Refresh **`.skillgrid/project/STRUCTURE.md`**, **`PROJECT.md`**, and other project files for other stable, non-design facts discovered during mapping.
10. Stop when the scope is clear enough to recommend `/skillgrid-brainstorm` or `/skillgrid-plan`.

## Completion Report

Summarize explored areas, artifacts inspected or updated, important findings, open questions, and the recommended next command.

</process>
