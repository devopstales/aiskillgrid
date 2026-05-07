# TODO

* mcp
  * [X] engram
  * [X] gitnexus
  * [X] context7
  * [X] deepwiki
  * [X] exa
  * [X] fire-claw
  * [ ] playwright
  * [ ] ai-mcp-sequentialthinking
* tools
  * [X] engram
  * [X] gitnexus
  * [X] ccc
  * [ ] codemaps
* memory
  * [-] engram
  * [-] gitnexus
  * [-] ccc
  * [ ] codemaps
* skillgrid-cli
  * installer
  * tui
    * engram integration ?
    * engram tui
  * web ui
    * engram integration ?
* search
  * web search
    * use mcp
  * code search
    * use ccc gitnexus
  * [X] agent personas
  * paralel search subagents
* questions
  * Assign different AI models to different SDD phases
  * Build Loop
  * Specialist persona board

## Components

* [ ] AGENTS.md
  * [ ] CONTEXT.md
  * [ ] Karpathy rules
* rules
* command
  * sdd-init
    * [X] sdd-init skill
    * [X] store to engram
    * [X] `_shared/skillgrid-handoff.md`
    * [ ] status, executive_summary, artifacts, and next_recommended ???
  * sdd-explore
    * [X] sdd-explore skill
    * [X] store to engram
    * [X] `_shared/skillgrid-handoff.md`
    * [ ] status, executive_summary, artifacts, and next_recommended ???
  * sdd-brainstorm
    * [X] Advanced questioning -> sdd-clarify skill
    * [X] sdd-explore skill
    * [X] sdd-propose skill
    * [X] sdd-spec skill
    * [X] sdd-design skill
    * [X] sdd-prd skill
    * [X] sdd-task skill
    * [X] `_shared/skillgrid-handoff.md`
  * sdd-apply
    * [X] smart/dum side
    * [X] process one slice at a time
    * [X] respect `[HITL]` and `[AFK]` labels
    * [X] `_shared/skillgrid-handoff.md`
  * sdd-archive
    * [X] `_shared/skillgrid-handoff.md`
  * sdd-verify
    * [X] `_shared/skillgrid-handoff.md`
  * sdd-ui-design
* skill
  * [X] `_shared/engram-convention.md`
  * [-] `_shared/openspec-convention.md`
    * [ ] Artifact File Paths
  * [X] `_shared/sdd-phase-common.md`
  * [-] `_shared/skillgrid-convention.md`
  * [-] `_shared/skillgrid-handoff.md`
  * [-] sdd-init
    * [X] config backend mode
    * [X] persist context
      * [X] `_shared/engram-convention.md`
      * [X] `_shared/openspec-convention.md`
      * [ ] `_shared/skillgrid-handoff.md`
    * [X] detect files
    * [X] skillgrid folder structure
    * [X] penspec folder structure
    * [-] openspec/config.yaml
      * [ ] ticketing integration
    * [X] create skill-registry
    * [X] Persist Project Context
      * [X] `_shared/skillgrid-convention.md`
    * [X] initialize ccc and gitnexus
    * [X] Return Summary for hybrid mode
      * [X] `_shared/skillgrid-convention.md`
  * [-] sdd-explore
    * [ ] exploration.md 
    * [X] persist context
      * [X] `_shared/engram-convention.md`
      * [X] `_shared/openspec-convention.md`
      * [ ] `_shared/skillgrid-handoff.md`
    * [X] `_shared/sdd-phase-common.md`
    * [X] Persist Artifact
      * [X] `_shared/skillgrid-convention.md`
    * [X] Import PRD Artifacts
  * [-] sdd-propose
    * [!] PRDs vs proposal.md !!! 
  * [-] sdd-design
  * [-] sdd-spec
  * [-] sdd-task
    * [X] agents to group tasks into vertical slice units
    * [X] add `[HITL]` or `[AFK]` labels.
  * [-] sdd-apply
  * [-] sdd-verify
  * [-] sdd-archive
  * [-] sdd-ui-design
  * [ ] sdd-test - runs tests and captures evidence tied to success criteria.
  * [ ] sdd-security - performs a deeper security pass when needed.

  * [ ] orchestrator ???
  * [ ] skill-registry
  * [ ] skillgrid-import-artifacts ???
  * [ ] skillgrid-parallel-research ???
  * [X] gitnexus-*
  * [X] ccc
  * [X] context7
  * [ ] deepwiki
  * [X] exa-search
  * [ ] fire-claw
  * [ ] playwright

## Plan

### Objective

Implement the highest-impact workflow upgrades from:
- Gentle-AI
- oh-my-opencode-slim
- code-yeongyu/oh-my-openagent
- obra/superpowers
- gsd-build/get-shit-done

Focus on enforceable pipeline behavior with strict phase gates, model routing, and persona-board decisions, plus workflow-level Norse mythology persona branding where it improves operator clarity.

### Source-specific imports to adopt

- [-] From `obra/superpowers` (partially adopted — see Milestone 1 enforcement + workflows):
  - [X] mandatory skill activation before task execution (fail-closed if required skill is skipped)
  - [X] true two-stage review gating (spec compliance -> code quality)
  - [ ] branch-finish protocol (`verify -> merge/PR/keep/discard -> cleanup`)
  - [X] systematic-debugging protocol integrated into `/sdd-diagnose`
  - [ ] optional git-worktree execution mode for risky/parallel slices
- [-] From `gsd-build/get-shit-done` (partially adopted — SDD/OpenSpec layering differs from GSD naming):
  - [-] strict 6-step loop variant mapped to SDD (`init/discuss/plan/execute/verify/ship`) — approximate mapping via `/sdd-init`, `/sdd-explore` + `/sdd-brainstorm`, `/sdd-loop` + `/sdd-apply`, `/sdd-verify`, `/sdd-archive`; not a branded GSD six-step runner
  - [ ] command to rebuild context/index for returning sessions (GSD-style re-map/re-hydrate)
  - [-] durable planning artifact pack parity:
    - [-] `PROJECT` equivalent (vision) — PRD layer (e.g. `.skillgrid/prd/`)
    - [-] `REQUIREMENTS` equivalent (scope) — proposal/specs/tasks under `openspec/changes/<change-id>/`
    - [-] `ROADMAP` equivalent (milestones/phases) — PRD index / execution snapshot patterns in docs
    - [-] `STATE` equivalent (current step + decisions) — handoff (`context_<change-id>.md`), events (`*.jsonl`), optional Engram `skillgrid/<change>/state`
    - [-] phase-context file for implementation decisions — `design.md`, slice specs, handoff entries
  - [X] "next best step" command (`gsd-progress --next` analogue for `/sdd-*`) — `next_recommended` on standard `sdd-*` return envelopes (not a standalone `/sdd-progress` CLI)
  - [-] parallel wave execution contract with atomic commit boundaries per task/slice — documented dependency-wave model (`06-multi-agent-work.md`, skills); no enforced atomic-commit automation
  - [-] verifier-driven fix-plan generation on failure (diagnose -> re-execute loop) — `/sdd-verify` and workflows route remediation via `next_recommended` → `/sdd-apply` / loops; no separate auto fix-plan agent
  - [-] runtime mode switch (`interactive` vs `yolo`) guarded by HITL policies — Interactive vs Automatic orchestration mode (`Automatic`/`Interactive` in `.configs/opencode.json` / Gentle-AI orchestrator docs); explicit `yolo` profile not present
  - [X] quality/balanced/budget model profile preset system — tier presets (`fast` / `balanced` / `deep`) in `.configs/norse-persona-contract.json`, `.configs/ide-model-mapping.json`
- [-] From `code-yeongyu/oh-my-openagent` (partially adopted):
  - [-] explicit "ultrawork" style execution profile in orchestrator:
    - [X] aggressive delegation defaults — orchestrator delegation rules (`sdd-orchestrator` prompt in `.configs/opencode.json`, Gentle-AI patterns)
    - [ ] background parallel scouting before implementation
    - [-] "continue until done" policy bounded by HITL stops — `/sdd-loop` bounded AFK continuation; full ultrawork-style autopilot out of scope
  - [ ] planner + plan-consultant split role (Prometheus/Metis pattern):
    - [ ] planning agent
    - [ ] plan-challenge/quality agent
  - [-] category-based routing for specialization (`visual`, `business-logic`, custom categories) — Norse persona **decision-type** routing (architecture/security/UX/release/risk); not arbitrary product categories
  - [-] built-in TODO continuation enforcer:
    - [-] detect unfinished queued tasks — task/spec/handovers plus verify gates
    - [X] auto-resume unfinished AFK tasks — `/sdd-loop` / `[AFK]` slice policy
    - [-] prevent premature "done" — gated by `/sdd-verify` + persona board hard blocks; not a dedicated TODO enforcer
  - [ ] comment-quality guard (comment checker):
    - [ ] block low-value generated comments
    - [ ] require justification for non-obvious comments
  - [ ] LSP/AST-first refactor mode:
    - [ ] prefer deterministic rename/refactor tools over text-only edits
    - [ ] fall back to text edits only when tool path unavailable
  - [ ] configurable background concurrency limits per provider/model
  - [-] context-injection policy:
    - [X] auto-inject `AGENTS.md` + project rule docs + active change artifacts — skill registry / compact rules injection protocol in orchestrator + `sdd-phase-common.md`
    - [-] keep injection bounded to avoid context bloat — policy described in skills; not a hard token budget automation
  - [-] hook surface parity review:
    - [X] pre-tool — e.g. Cursor `preToolUse` (see `100-ide-configs.md`)
    - [X] post-tool — e.g. Cursor `postToolUse` (see `100-ide-configs.md`)
    - [ ] prompt-submit
    - [ ] stop hooks
  - [-] session tooling parity:
    - [-] list/search/read historical sessions — depends on IDE/Engram; not uniform across surfaces
    - [X] extract reusable decisions into durable artifacts — handoff, events, Engram/mem_save patterns in skills/workflows

### Milestone 1 — Core Enforcement (highest ROI)

- [X] Build `sdd-orchestrator` as first-class runtime persona (not doc-only):
  - [X] explicit phase routing and stop conditions
  - [X] HITL hard-stop policy
  - [X] AFK continuation policy
- [X] Add mandatory skill-gate checks per phase (`sdd-brainstorm`, `sdd-apply`, `sdd-verify`):
  - [X] fail closed if required artifact or gate is missing
- [X] Enforce two-stage review gate everywhere:
  - [X] Stage 1: spec compliance
  - [X] Stage 2: code quality
  - [X] block progression on critical findings
- [X] Standardize return envelope for all `sdd-*` commands:
  - [X] `status`
  - [X] `executive_summary`
    - used tokens
  - [X] `artifacts`
  - [X] `next_recommended`
  - [X] `risks`

### Milestone 2 — Build Loop + Persona Board

- [X] Implement `/sdd-loop` command (Ralph-loop style):
  - [X] pick one AFK-safe slice
  - [X] execute
  - [X] capture evidence
  - [X] reassess
  - [X] continue or stop
- [X] Implement specialist persona board as executable workflow:
  - [X] preset selection by decision type (arch/security/UX/go-no-go)
  - [X] parallel persona reports
  - [X] merge + accepted/rejected options
  - [X] conflict -> HITL block
- [X] Define hard block semantics:
  - [X] `spec-verifier` critical => block
  - [X] `security-auditor` critical => block
  - [X] unresolved persona conflict => block

### Milestone 2.5 — Norse Persona Theming + Command Functions

- [X] Plan specialist agent personas with Norse mythology branding (workflow-facing, not decorative-only):
  - [X] `Odin` -> orchestrator/planner authority
  - [X] `Thor` -> implementation enforcer (delivery and momentum)
  - [X] `Tyr` -> spec/compliance verifier
  - [X] `Heimdall` -> security and gate sentinel
  - [X] `Frigg` -> UX/product clarity reviewer
  - [X] `Loki` -> adversarial critic/challenge persona
- [X] Define persona selection matrix by decision type:
  - [X] architecture
  - [X] security
  - [X] UX/content
  - [X] go/no-go release
  - [X] risk acceptance
- [X] Planned command functions for Norse persona workflow:
  - [X] `/sdd-persona-board --preset <arch|security|ux|release>`
  - [X] `/sdd-persona-list` (show personas, roles, and availability)
  - [X] `/sdd-persona-route --decision <type>` (auto-select personas)
  - [X] `/sdd-persona-report --id <decision-id>` (merge verdicts and conflicts)
  - [X] `/sdd-persona-resolve --id <decision-id>` (record accepted/rejected options)
  - [X] `/sdd-persona-health` (check persona prompt packs and tool readiness)
- [X] Define command return contracts for persona commands:
  - [X] `status`, `executive_summary`, `artifacts`, `next_recommended`, `risks`
  - [X] `personas_invoked`, `conflicts`, `hitl_required`, `accepted_decision`
- [X] Add HITL safety rules for themed personas:
  - [X] no persona can override hard gates
  - [X] unresolved critical conflict -> block
  - [X] user remains final authority on release/destructive choices

### Milestone 3 — Model Routing + Session Efficiency

- [ ] Add per-phase model routing in config:
  - [ ] `explore` fast/cheap
  - [ ] `apply` balanced
  - [ ] `verify/design` stronger reasoning
- [ ] Add runtime preset switching (without reinstall/restart)
- [ ] Add subagent session reuse for repeated slices/boards
- [ ] Add safe auto-continuation:
  - [ ] cooldown
  - [ ] AFK-only advancement
  - [ ] fresh verification precondition

### Workflow Hardening Backlog

- [ ] Optional worktree mode for risky or parallel implementation lanes
- [ ] Agent health-check command (ping all required personas + MCP readiness)
- [ ] Persona report contract template:
  - [ ] severity
  - [ ] evidence path
  - [ ] impacted artifact(s)
  - [ ] disposition (`must-fix | accept-risk | follow-up`)
- [ ] CI guards:
  - [ ] tasks label validator
  - [ ] spec matrix presence
  - [ ] gate result must exist before archive

### Out of Scope (for now)

- [ ] Full narrative/lore UI storytelling layers beyond command-level persona mapping
- [ ] Multiplexer/UI novelty features
- [ ] Broad cross-harness packaging before core gates are stable
