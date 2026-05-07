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
* agent personas
  * [ ] sdd-orchestrator
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
  * agent persona
  * paralel search subagents
* questions
  * Assign different AI models to different SDD phases
  * Build Loop
  * Specialist persona board

## Components

* [ ] AGENTS.md
  * [ ] CONTEXT.md
  * [ ] Karpathy rules
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

Focus on enforceable pipeline behavior (not branding), with strict phase gates, model routing, and persona-board decisions.

### Source-specific imports to adopt

- [ ] From `obra/superpowers`:
  - [ ] mandatory skill activation before task execution (fail-closed if required skill is skipped)
  - [ ] true two-stage review gating (spec compliance -> code quality)
  - [ ] branch-finish protocol (`verify -> merge/PR/keep/discard -> cleanup`)
  - [ ] systematic-debugging protocol integrated into `/sdd-diagnose`
  - [ ] optional git-worktree execution mode for risky/parallel slices
- [ ] From `gsd-build/get-shit-done`:
  - [ ] strict 6-step loop variant mapped to SDD (`init/discuss/plan/execute/verify/ship`)
  - [ ] command to rebuild context/index for returning sessions (GSD-style re-map/re-hydrate)
  - [ ] durable planning artifact pack parity:
    - [ ] `PROJECT` equivalent (vision)
    - [ ] `REQUIREMENTS` equivalent (scope)
    - [ ] `ROADMAP` equivalent (milestones/phases)
    - [ ] `STATE` equivalent (current step + decisions)
    - [ ] phase-context file for implementation decisions
  - [ ] "next best step" command (`gsd-progress --next` analogue for `/sdd-*`)
  - [ ] parallel wave execution contract with atomic commit boundaries per task/slice
  - [ ] verifier-driven fix-plan generation on failure (diagnose -> re-execute loop)
  - [ ] runtime mode switch (`interactive` vs `yolo`) guarded by HITL policies
  - [ ] quality/balanced/budget model profile preset system
- [ ] From `code-yeongyu/oh-my-openagent`:
  - [ ] explicit "ultrawork" style execution profile in orchestrator:
    - [ ] aggressive delegation defaults
    - [ ] background parallel scouting before implementation
    - [ ] "continue until done" policy bounded by HITL stops
  - [ ] planner + plan-consultant split role (Prometheus/Metis pattern):
    - [ ] planning agent
    - [ ] plan-challenge/quality agent
  - [ ] category-based routing for specialization (`visual`, `business-logic`, custom categories)
  - [ ] built-in TODO continuation enforcer:
    - [ ] detect unfinished queued tasks
    - [ ] auto-resume unfinished AFK tasks
    - [ ] prevent premature "done"
  - [ ] comment-quality guard (comment checker):
    - [ ] block low-value generated comments
    - [ ] require justification for non-obvious comments
  - [ ] LSP/AST-first refactor mode:
    - [ ] prefer deterministic rename/refactor tools over text-only edits
    - [ ] fall back to text edits only when tool path unavailable
  - [ ] configurable background concurrency limits per provider/model
  - [ ] context-injection policy:
    - [ ] auto-inject `AGENTS.md` + project rule docs + active change artifacts
    - [ ] keep injection bounded to avoid context bloat
  - [ ] hook surface parity review:
    - [ ] pre-tool
    - [ ] post-tool
    - [ ] prompt-submit
    - [ ] stop hooks
  - [ ] session tooling parity:
    - [ ] list/search/read historical sessions
    - [ ] extract reusable decisions into durable artifacts

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

- [ ] Implement `/sdd-loop` command (Ralph-loop style):
  - [ ] pick one AFK-safe slice
  - [ ] execute
  - [ ] capture evidence
  - [ ] reassess
  - [ ] continue or stop
- [ ] Implement specialist persona board as executable workflow:
  - [ ] preset selection by decision type (arch/security/UX/go-no-go)
  - [ ] parallel persona reports
  - [ ] merge + accepted/rejected options
  - [ ] conflict -> HITL block
- [ ] Define hard block semantics:
  - [ ] `spec-verifier` critical => block
  - [ ] `security-auditor` critical => block
  - [ ] unresolved persona conflict => block

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

- [ ] Thematic/lore persona branding
- [ ] Multiplexer/UI novelty features
- [ ] Broad cross-harness packaging before core gates are stable
