# Persona Board And Norse Specialist Personas

This document describes the **specialist persona decision board** and the **Norse-themed specialist agent personas** used in AISkillGrid. It complements the general multi-agent guidance in `06-multi-agent-work.md` and the broader persona notes in `07-subagent-personas.md`.

## Purpose

A **persona board** is for decisions that benefit from **independent viewpoints** before the parent session or human operator commits: architecture trade-offs, security posture, UX and content clarity, go/no-go release, and explicit risk acceptance.

The board is **advisory**. It does not replace the user, PRD, OpenSpec artifacts, or orchestration policy. It produces **durable reports**, a **recorded decision**, and a **clear continue vs HITL** outcome.

## Source Of Truth

Canonical definitions (personas, routing matrix, required return fields, hard rules) live in:

- `.configs/norse-persona-contract.json`

Workflow prompts for board-related commands live under:

- `.agents/workflows/sdd-persona-*.md`
- `.agents/workflows/sdd-board.md` (compatibility alias for the Norse board flow)

## Specialist Agent Personas (Norse IDs)

Workflow-facing specialist IDs are lowercase: `odin`, `thor`, `tyr`, `heimdall`, `frigg`, `loki`. The contract maps each to a previous neutral role name for readability only; **routing and commands use the Norse IDs**.

| ID | Title | Role | Hard gate on critical? | Notes |
| --- | --- | --- | --- | --- |
| `odin` | Odin | Orchestrator and planner authority | No | Board routing, merge, HITL boundaries. |
| `thor` | Thor | Implementation enforcer | No | Delivery feasibility, execution quality, momentum. |
| `tyr` | Tyr | Spec and compliance verifier | **Yes** | Traceability, acceptance criteria, spec compliance. |
| `heimdall` | Heimdall | Security and release-gate sentinel | **Yes** | Threat model, exploitability, release risk. |
| `frigg` | Frigg | UX and product clarity reviewer | No | Flows, accessibility, content clarity. |
| `loki` | Loki | Adversarial critic | No | Counterexamples, assumption stress-test; can surface conflicts that require HITL. |

## Decision Routing Matrix

Presets map **decision types** to **which personas** participate. Aliases such as `arch`, `ux`, `release`, and `risk` normalize to the canonical keys below (see `sdd-persona-board.md` and `sdd-persona-route.md`).

| Decision type | Personas invoked (order is workflow-defined) |
| --- | --- |
| `architecture` | `odin`, `thor`, `tyr`, `loki` |
| `security` | `heimdall`, `tyr`, `thor`, `loki` |
| `ux-content` | `frigg`, `loki`, `thor` |
| `go-no-go-release` | `odin`, `tyr`, `heimdall`, `thor`, `frigg` |
| `risk-acceptance` | `odin`, `loki`, `tyr`, `heimdall` |

## Commands (Persona Board Family)

Primary entrypoints:

- `/sdd-persona-board --preset <arch|security|ux|release|risk>` — full board cycle: scope, route, parallel reports, merge, persist.
- `/sdd-persona-list` — personas, roles, and availability or limitations by surface when known.
- `/sdd-persona-route --decision <type>` — auto-select personas from the matrix; **fail closed** on unknown types.
- `/sdd-persona-report --id <decision-id>` — merge verdicts and conflict summary for one decision.
- `/sdd-persona-resolve --id <decision-id>` — record accepted decision and rejected options in handoff/events.
- `/sdd-persona-health` — check contract alignment, model tier mapping, and surface capability readiness where applicable.

Compatibility:

- `/sdd-board` is treated as an alias that follows the same Norse board flow (see `.agents/workflows/sdd-board.md`).

Command inventory and phase context: `04-commands.md`.

## Durable Artifacts

Board work must be visible outside chat. Typical paths:

- Reports: `.skillgrid/tasks/research/<change-id>/`
- Handoff and decisions: `.skillgrid/tasks/context_<change-id>.md`
- Timeline: `.skillgrid/tasks/events/<change-id>.jsonl`

Follow handoff and event contracts in `.agents/skills/_shared/skillgrid-handoff.md` when persisting board output.

## Return Contract (Persona Commands)

Persona-related commands share an extended envelope. Required fields (from the contract):

- `status`
- `executive_summary` (including `overview` and `used_tokens` where applicable)
- `artifacts`
- `next_recommended`
- `risks`
- `personas_invoked`
- `conflicts`
- `hitl_required`
- `accepted_decision`

Milestone-wide enforcement envelopes for other `sdd-*` phases remain defined in `.agents/skills/_shared/sdd-enforcement-contract.md`; persona commands **add** the board-specific keys above.

## Hard Gates And HITL Rules

From `.configs/norse-persona-contract.json`:

- No persona overrides **hard gates**.
- **`tyr`** or **`heimdall`** reporting **critical** severity blocks progression until resolved or explicitly handled per policy.
- **Unresolved critical conflict** between personas blocks progression (HITL).
- The **user** remains the final authority on **release** and **destructive** choices.

Unresolved conflicts escalated by `loki` or merged reports still flow through these rules.

## Where Agent Prompts Live By Surface

Specialist persona markdown agents use the shared `skillgrid-<persona>.md` filenames (same body conventions: identity, mandatory context, rules, composition):

- Cursor: `.cursor/agents/`
- GitHub Copilot agents mirror: `.github/agents/`
- Kilo and OpenCode mirrors (transformed frontmatter): `.kilo/agents/`, `.opencode/agents/`
- Antigravity hub: `.agent/agents/`

Regenerate cross-IDE persona manifests (optional tooling):

```bash
node scripts/render-multi-ide-personas.mjs
```

Mirror agent files from Cursor sources:

```bash
./scripts/sync-ide-assets.sh
```

Surface capability tiers and fallback behavior: `.configs/ide-persona-capabilities.json`, described alongside other IDE layout notes in `100-ide-configs.md`.

## Related Reading

- `04-commands.md` — command list and persona board entrypoints  
- `02-workflow-usage.md` — when boards appear in the main workflow  
- `06-multi-agent-work.md` — delegate work, not responsibility  
- `07-subagent-personas.md` — subagent operating model  
- `100-ide-configs.md` — Norse persona config set and IDE layout  
