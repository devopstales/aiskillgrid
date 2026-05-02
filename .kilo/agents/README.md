# Agent Personas

Specialist personas that play a single role with a single perspective. Each persona is a Markdown file consumed as a system prompt by your harness (Claude Code, Cursor, Copilot, etc.).

| Persona | Role | Best for |
|---------|------|---------|
| [code-reviewer](code-reviewer.md) | Senior Staff Engineer | Five-axis code review before merge |
| [security-auditor](security-auditor.md) | Security Engineer | Vulnerability detection, OWASP-style audit |
| [test-engineer](test-engineer.md) | QA Engineer | Test strategy, coverage analysis, Prove-It pattern |
| [spec-verifier](spec-verifier.md) | Spec traceability | Delta specs / tasks / implementation alignment |
| [explore-architect](explore-architect.md) | Systems explorer | Brownfield docs: ARCHITECTURE, STRUCTURE, PROJECT |
| [task-breakdown-auditor](task-breakdown-auditor.md) | Planning reviewer | `tasks.md` quality before implementation |
| [design-critic](design-critic.md) | UX/design critic | DESIGN.md flows, a11y, API boundaries at design stage |
| [researcher](researcher.md) | Research analyst | Cited web/repo/docs research via Exa, Firecrawl, DeepWiki, Context7 |

## How personas relate to skills and commands

Three layers, each with a distinct job:

| Layer | What it is | Example | Composition role |
|-------|-----------|---------|------------------|
| **Skill** | A workflow with steps and exit criteria | `code-review-and-quality` | The *how* — invoked from inside a persona or command |
| **Persona** | A role with a perspective and an output format | `code-reviewer` | The *who* — adopts a viewpoint, produces a report |
| **Command** | A user-facing entry point | `/skillgrid-validate`, `/skillgrid-apply` | The *when* — composes workflows and skills |

The user (or a slash command) is the orchestrator. **Personas do not call other personas.** Skills are mandatory hops inside a persona's workflow.

Each persona includes an **Indexing and memory** section: when Engram, GitNexus, or MCP `server-memory` are configured, use them per [`.agents/skills/references/indexing-and-memory.md`](../../.agents/skills/references/indexing-and-memory.md).

## Skillgrid mapping (this hub)

| Persona | Typical Skillgrid phase |
|---------|-------------------------|
| `explore-architect` | `/skillgrid-explore` |
| `researcher` | `/skillgrid-brainstorm` (heavy research), optional `/skillgrid-explore` depth |
| `design-critic` | `/skillgrid-plan` (with `DESIGN.md`) or direct spawn |
| `task-breakdown-auditor` | `/skillgrid-breakdown` |
| `spec-verifier` | `/skillgrid-validate` (spec pass) or pre-merge traceability |
| `code-reviewer` | `/skillgrid-validate` |
| `test-engineer` | `/skillgrid-test` |
| `security-auditor` | `/skillgrid-security` |

## When to use each

### Direct persona invocation
Pick this when you want one perspective on the current change and the user is in the loop.

- "Review this PR" → invoke `code-reviewer` directly
- "Does this match the delta spec and tasks?" → invoke `spec-verifier` directly
- "Map this repo for onboarding" → invoke `explore-architect` directly
- "Are these tasks ready to implement?" → invoke `task-breakdown-auditor` directly
- "Critique DESIGN.md before we build" → invoke `design-critic` directly
- "Research prior art / landscape with citations" → invoke `researcher` directly
- "Are there security issues in `auth.ts`?" → invoke `security-auditor` directly
- "What tests are missing for the checkout flow?" → invoke `test-engineer` directly

### Slash command (single perspective)
Pick this when the hub’s phase command already encodes skills and steps; spawn a persona when you want that phase in a **fresh subagent context**.

- `/skillgrid-validate` → often pairs with `code-reviewer` and/or `spec-verifier` (user chooses)
- `/skillgrid-test` → pairs with `test-engineer` when delegating test design
- `/skillgrid-explore` → pairs with `explore-architect` for large repos

### Slash command (orchestrator — fan-out)
Pick this only when **independent** investigations can run in parallel and produce reports that a single agent then merges.

- **`/skillgrid-validate`** (conceptually) — parallel **spec + code + security** when reports are independent, e.g. `spec-verifier` + `code-reviewer` + `security-auditor`; add `test-engineer` when coverage is the third axis. **Ordering:** this hub’s command runs review then security **sequentially** in one turn; parallel fan-out is for harnesses that support true concurrent subagents.
- Legacy/example: `/ship` → same idea: `code-reviewer` + `security-auditor` + `test-engineer` in parallel, then merge.

This is the only orchestration pattern this repo endorses. See [orchestration-patterns.md](../../.agents/skills_back/references/orchestration-patterns.md) for the full pattern catalog and anti-patterns.

## Decision matrix

```
Is the work a single perspective on a single artifact?
├── Yes → Direct persona invocation
└── No  → Are the sub-tasks independent (no shared mutable state, no ordering)?
         ├── Yes → Slash command with parallel fan-out (e.g. /skillgrid-validate-style merge of reports)
         └── No  → Sequential slash commands run by the user (/skillgrid-plan → /skillgrid-breakdown → …)
```

## Worked example: valid orchestration

Parallel fan-out (subagent-capable harness), then merge:

```
(skillgrid-validate-style merge)
  ├── (parallel) spec-verifier      → traceability report
  ├── (parallel) code-reviewer      → code review report
  ├── (parallel) security-auditor   → audit report
  └── (parallel) test-engineer       → optional coverage report
                  ↓
        merge phase (main agent)
                  ↓
        go/no-go decision + rollback plan
```

Why this works:
- Each sub-agent operates on the same change but produces a **different perspective**
- They have no dependencies on each other → genuine parallelism, real wall-clock savings
- Each runs in a fresh context window → main session stays uncluttered
- The merge step is small and benefits from full context, so it stays in the main agent

**State channel (no extra checkout):** For Skillgrid + OpenSpec, the parent and `Task` subagents align on **`.skillgrid/tasks/context_<change-id>.md`** and **`openspec/changes/<id>/tasks.md`** in a **single working tree**—an artifact-backed subagent model with optional git worktree separation. Subagents read the handoff, spill long work to **`.skillgrid/tasks/research/<change-id>/`**, and return short paths; the orchestrator runs **`/skillgrid-apply`** with checkboxes and optional **`[HITL]`** / **`[AFK]`** tags (see `docs/02-workflow-usage.md` — *Filesystem handoff* and *HITL vs AFK slices*).

## Worked example: invalid orchestration (do not build this)

A `meta-orchestrator` persona whose job is "decide which other persona to call":

```
/work-on-pr → meta-orchestrator
                  ↓ (decides "this needs a review")
              code-reviewer
                  ↓ (returns)
              meta-orchestrator (paraphrases result)
                  ↓
              user
```

Why this fails:
- Pure routing layer with no domain value
- Adds two paraphrasing hops → information loss + 2× token cost
- The user already knows they want a review; let them call `/skillgrid-validate` or spawn `code-reviewer` directly
- Replicates work that slash commands and `AGENTS.md` intent-mapping already do

## Rules for personas

1. A persona is a single role with a single output format. If you find yourself adding a second role, create a second persona.
2. **Personas do not invoke other personas.** Composition is the job of slash commands or the user. On Claude Code this is also a hard platform constraint — *"subagents cannot spawn other subagents"* — so the rule is enforced for you.
3. A persona may invoke skills (the *how*).
4. Every persona file ends with a "Composition" block stating where it fits.

## Agent discipline

Skillgrid personas follow the lightweight discipline adapted from `oh-my-openagent/src/agents`, without copying its runtime or model fallback system:

1. **Identity is explicit:** each persona declares its designated identity and must stay in that role.
2. **Tool posture matches the role:** reviewer, verifier, auditor, critic, researcher, and explorer personas are report-first. They do not edit files unless the parent prompt explicitly assigns that work.
3. **No duplicate delegation:** when exploration or research has already been delegated for the same scope, do not repeat the same search. Continue only with non-overlapping work or wait for the delegated result.
4. **No speculative claims:** cite read artifacts, source URLs, commands, or evidence. Mark unknowns instead of inventing facts.
5. **Hard blocks:** no commits without explicit user request, no deleting or weakening tests to pass, no type/error suppression as a shortcut, and no hidden persona orchestration.

## Role coverage assessment

The current Skillgrid set covers the useful `oh-my-openagent` role families without adding new personas yet:

| Inspiration | Skillgrid coverage | Decision |
|-------------|--------------------|----------|
| Librarian / external reference search | `researcher` | Covered. Keep cited research in one persona. |
| Explore / contextual grep | `explore-architect` plus platform explore subagents | Covered. Use direct explore subagents for broad read-only search, and this persona for durable architecture docs. |
| Oracle / Metis / Momus consultation and plan critique | `spec-verifier`, `task-breakdown-auditor`, `code-reviewer`, and slash-command review gates | Covered as review modes for now. Add a new consultant only if repeated workflows need a distinct report format. |
| Sisyphus / Atlas orchestrators | `/skillgrid-*` slash commands and parent session | Do not add a meta-orchestrator persona; orchestration remains command-owned. |

## Discipline-agent fan-out

Skillgrid treats discipline agents as independent specialist perspectives that the parent command merges. The parent owns scope, state, and final decisions; personas provide bounded reports.

| Phase | Primary personas | Output |
|-------|------------------|--------|
| Explore / Brainstorm | `explore-architect`, `researcher` | Repo map, cited research, alternatives, unresolved questions. |
| Plan / Design | `design-critic`, `task-breakdown-auditor` | Design critique, task readiness, HITL/AFK quality. |
| Apply | Implementer subagent plus `spec-verifier` and `code-reviewer` when delegated | One AFK slice, spec compliance report, code-quality report. |
| Test | `test-engineer` | Test strategy, missing coverage, Playwright/browser evidence. |
| Security / Validate | `security-auditor`, `spec-verifier`, `code-reviewer`, optional `test-engineer` | Independent go/no-go reports merged by the parent. |

Fan-out is valid only when each persona has a different perspective and no shared mutable output. The parent must read every report, reconcile conflicts, update the handoff, and run final verification before advancing status.

## Consistency contract

All agents in this directory follow the same focused single-role contract:

1. **Frontmatter:** include `name`, `description`, `tools`, and `color`.
2. **Single job:** state the persona's role in the first paragraph and stay inside it.
3. **Identity and discipline:** include an "Identity and discipline" section before "Mandatory Context".
4. **Mandatory context:** read the relevant request, artifacts, and project rules before producing findings or edits.
5. **Scope discipline:** say what is out of scope and recommend another persona or command instead of doing that work.
6. **Structured output:** use the persona's report template and classify findings consistently.
7. **No hidden orchestration:** do not call or impersonate another persona.
8. **No source edits unless assigned:** reviewer, verifier, critic, researcher, and auditor personas produce reports; implementation remains with the parent session or `/skillgrid-apply`.
9. **Sync check enforcement:** `scripts/sync-ide-assets.sh --check` fails when personas are missing required frontmatter or required sections.

## Claude Code interop

The personas in this repo are designed to work as Claude Code subagents and as Agent Teams teammates without modification:

- **As subagents:** auto-discovered when this plugin is enabled (no path config needed). Use the Agent tool with `subagent_type` matching the persona `name` in frontmatter (e.g. `code-reviewer`, `spec-verifier`, `explore-architect`, `researcher`). Parallel fan-out + merge is the canonical pattern when independent.
- **As Agent Teams teammates** (experimental, requires `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1`): reference the same persona name when spawning a teammate. The persona's body is **appended to** the teammate's system prompt as additional instructions (not a replacement), so your persona text sits on top of the team-coordination instructions the lead installs (SendMessage, task-list tools, etc.).

Subagents only report results back to the main agent. Agent Teams let teammates message each other directly. Use subagents when reports are enough; use Agent Teams when sub-agents need to challenge each other's findings (e.g. competing-hypothesis debugging). See [orchestration-patterns.md](../../.agents/skills_back/references/orchestration-patterns.md) for the full mapping.

Plugin agents do not support `hooks`, `mcpServers`, or `permissionMode` frontmatter — those fields are silently ignored. Avoid relying on them when authoring new personas here.

## Adding a new persona

1. Create `agents/<slug>.md` with the same frontmatter format used by existing personas; set `name: <slug>` in frontmatter.
2. Define the role, scope, output format, and rules.
3. Add an **Identity and discipline** block before **Mandatory Context**.
4. Add a **Composition** block at the bottom (Invoke directly when / Invoke via / Do not invoke from another persona).
5. Add the persona to the table at the top of this file.
6. Mirror the file to `.kilo/agents/`, `.opencode/agents/`, and `.github/agents/`.
7. If the persona enables a new orchestration pattern, document it in `.agents/skills_back/references/orchestration-patterns.md` (or promote it into `docs/`) rather than inventing the pattern in the persona file itself.
