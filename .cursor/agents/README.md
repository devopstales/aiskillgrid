# Agent Personas

Specialist personas that play a single role with a single perspective. Each persona is a Markdown file consumed as a system prompt by your harness (Claude Code, Cursor, Copilot, etc.).

| Persona | Role | Best for |
|---------|------|---------|
| [skillgrid-code-reviewer](skillgrid-code-reviewer.md) | Senior Staff Engineer | Five-axis code review before merge |
| [skillgrid-security-auditor](skillgrid-security-auditor.md) | Security Engineer | Vulnerability detection, OWASP-style audit |
| [skillgrid-test-engineer](skillgrid-test-engineer.md) | QA Engineer | Test strategy, coverage analysis, Prove-It pattern |
| [skillgrid-spec-verifier](skillgrid-spec-verifier.md) | Spec traceability | Delta specs / tasks / implementation alignment |
| [skillgrid-explore-architect](skillgrid-explore-architect.md) | Systems explorer | Brownfield docs: ARCHITECTURE, STRUCTURE, PROJECT |
| [skillgrid-task-breakdown-auditor](skillgrid-task-breakdown-auditor.md) | Planning reviewer | `tasks.md` quality before implementation |
| [skillgrid-design-critic](skillgrid-design-critic.md) | UX/design critic | DESIGN.md flows, a11y, API boundaries at design stage |
| [skillgrid-researcher](skillgrid-researcher.md) | Research analyst | Cited web/repo/docs research via Exa, Firecrawl, DeepWiki, Context7 |

## How personas relate to skills and commands

Three layers, each with a distinct job:

| Layer | What it is | Example | Composition role |
|-------|-----------|---------|------------------|
| **Skill** | A workflow with steps and exit criteria | `code-review-and-quality` | The *how* — invoked from inside a persona or command |
| **Persona** | A role with a perspective and an output format | `skillgrid-code-reviewer` | The *who* — adopts a viewpoint, produces a report |
| **Command** | A user-facing entry point | `/skillgrid-validate`, `/skillgrid-validate` | The *when* — composes workflows and skills |

The user (or a slash command) is the orchestrator. **Personas do not call other personas.** Skills are mandatory hops inside a persona's workflow.

Each `skillgrid-*` persona includes an **Indexing and memory** section: when Engram, graphify, or MCP `server-memory` are configured, use them per [`.agents/skills/references/indexing-and-memory.md`](../../.agents/skills/references/indexing-and-memory.md).

## Skillgrid mapping (this hub)

| Persona | Typical Skillgrid phase |
|---------|-------------------------|
| `skillgrid-explore-architect` | `/skillgrid-explore` |
| `skillgrid-researcher` | `/skillgrid-brainstorm` (heavy research), optional `/skillgrid-explore` depth |
| `skillgrid-design-critic` | `/skillgrid-plan` (with `DESIGN.md`) or direct spawn |
| `skillgrid-task-breakdown-auditor` | `/skillgrid-breakdown` |
| `skillgrid-spec-verifier` | `/skillgrid-validate` (spec pass) or pre-merge traceability |
| `skillgrid-code-reviewer` | `/skillgrid-validate` |
| `skillgrid-test-engineer` | `/skillgrid-test` |
| `skillgrid-security-auditor` | `/skillgrid-security` |

## When to use each

### Direct persona invocation
Pick this when you want one perspective on the current change and the user is in the loop.

- "Review this PR" → invoke `skillgrid-code-reviewer` directly
- "Does this match the delta spec and tasks?" → invoke `skillgrid-spec-verifier` directly
- "Map this repo for onboarding" → invoke `skillgrid-explore-architect` directly
- "Are these tasks ready to implement?" → invoke `skillgrid-task-breakdown-auditor` directly
- "Critique DESIGN.md before we build" → invoke `skillgrid-design-critic` directly
- "Research prior art / landscape with citations" → invoke `skillgrid-researcher` directly
- "Are there security issues in `auth.ts`?" → invoke `skillgrid-security-auditor` directly
- "What tests are missing for the checkout flow?" → invoke `skillgrid-test-engineer` directly

### Slash command (single perspective)
Pick this when the hub’s phase command already encodes skills and steps; spawn a persona when you want that phase in a **fresh subagent context**.

- `/skillgrid-validate` → often pairs with `skillgrid-code-reviewer` and/or `skillgrid-spec-verifier` (user chooses)
- `/skillgrid-test` → pairs with `skillgrid-test-engineer` when delegating test design
- `/skillgrid-explore` → pairs with `skillgrid-explore-architect` for large repos

### Slash command (orchestrator — fan-out)
Pick this only when **independent** investigations can run in parallel and produce reports that a single agent then merges.

- **`/skillgrid-validate`** (conceptually) — parallel **spec + code + security** when reports are independent, e.g. `skillgrid-spec-verifier` + `skillgrid-code-reviewer` + `skillgrid-security-auditor`; add `skillgrid-test-engineer` when coverage is the third axis. **Ordering:** this hub’s command runs review then security **sequentially** in one turn; parallel fan-out is for harnesses that support true concurrent subagents.
- Legacy/example: `/ship` → same idea: `skillgrid-code-reviewer` + `skillgrid-security-auditor` + `skillgrid-test-engineer` in parallel, then merge.

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
  ├── (parallel) skillgrid-spec-verifier      → traceability report
  ├── (parallel) skillgrid-code-reviewer      → code review report
  ├── (parallel) skillgrid-security-auditor   → audit report
  └── (parallel) skillgrid-test-engineer       → optional coverage report
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

**State channel (no extra checkout):** For Skillgrid + OpenSpec, the parent and `Task` subagents align on **`.skillgrid/tasks/context_<change-id>.md`** and **`openspec/changes/<id>/tasks.md`** in a **single working tree**—the same model as [obra/superpowers](https://github.com/obra/superpowers) subagent-driven work, but **without** git worktrees. Subagents read the handoff, spill long work to **`.skillgrid/tasks/research/<change-id>/`**, and return short paths; the orchestrator runs **`/skillgrid-apply`** with checkboxes and optional **`[HITL]`** / **`[AFK]`** tags (see `docs/workflow.md` — *Filesystem handoff* and *HITL vs AFK slices*).

## Worked example: invalid orchestration (do not build this)

A `meta-orchestrator` persona whose job is "decide which other persona to call":

```
/work-on-pr → meta-orchestrator
                  ↓ (decides "this needs a review")
              skillgrid-code-reviewer
                  ↓ (returns)
              meta-orchestrator (paraphrases result)
                  ↓
              user
```

Why this fails:
- Pure routing layer with no domain value
- Adds two paraphrasing hops → information loss + 2× token cost
- The user already knows they want a review; let them call `/skillgrid-validate` or spawn `skillgrid-code-reviewer` directly
- Replicates work that slash commands and `AGENTS.md` intent-mapping already do

## Rules for personas

1. A persona is a single role with a single output format. If you find yourself adding a second role, create a second persona.
2. **Personas do not invoke other personas.** Composition is the job of slash commands or the user. On Claude Code this is also a hard platform constraint — *"subagents cannot spawn other subagents"* — so the rule is enforced for you.
3. A persona may invoke skills (the *how*).
4. Every persona file ends with a "Composition" block stating where it fits.

## Claude Code interop

The personas in this repo are designed to work as Claude Code subagents and as Agent Teams teammates without modification:

- **As subagents:** auto-discovered when this plugin is enabled (no path config needed). Use the Agent tool with `subagent_type` matching the persona `name` in frontmatter (e.g. `skillgrid-code-reviewer`, `skillgrid-spec-verifier`, `skillgrid-explore-architect`, `skillgrid-researcher`). Parallel fan-out + merge is the canonical pattern when independent.
- **As Agent Teams teammates** (experimental, requires `CLAUDE_CODE_EXPERIMENTAL_AGENT_TEAMS=1`): reference the same persona name when spawning a teammate. The persona's body is **appended to** the teammate's system prompt as additional instructions (not a replacement), so your persona text sits on top of the team-coordination instructions the lead installs (SendMessage, task-list tools, etc.).

Subagents only report results back to the main agent. Agent Teams let teammates message each other directly. Use subagents when reports are enough; use Agent Teams when sub-agents need to challenge each other's findings (e.g. competing-hypothesis debugging). See [orchestration-patterns.md](../../.agents/skills_back/references/orchestration-patterns.md) for the full mapping.

Plugin agents do not support `hooks`, `mcpServers`, or `permissionMode` frontmatter — those fields are silently ignored. Avoid relying on them when authoring new personas here.

## Adding a new persona

1. Create `agents/skillgrid-<slug>.md` with the same frontmatter format used by existing personas; set `name: skillgrid-<slug>` in frontmatter (keep the `skillgrid-` prefix for hub personas).
2. Define the role, scope, output format, and rules.
3. Add a **Composition** block at the bottom (Invoke directly when / Invoke via / Do not invoke from another persona).
4. Add the persona to the table at the top of this file.
5. Mirror the file to `.kilo/agents/`, `.opencode/agents/`, and `.github/agents/`.
6. If the persona enables a new orchestration pattern, document it in `.agents/skills_back/references/orchestration-patterns.md` (or promote it into `docs/`) rather than inventing the pattern in the persona file itself.
