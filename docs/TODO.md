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
  * [ ] agent-browser
* tools
  * [X] engram
  * [X] gitnexus
  * [X] ccc
  * [ ] codemaps
* memory
  * [-] engram
    * [X] memory versioning (`version.id`, `version.previous_id`, `timestamps.*`)
    * [X] deterministic conflict resolution rules
    * [X] cross-source ranking strategy with weighted score + freshness decay
    * [X] vector index mapping schema (versioned fields/dimensions/metric)
    * [X] HNSW tuning guidance (`M`, `efConstruction`, memory-optimized variants)
    * [X] pre-allocation + `sync.Pool` policy for hot paths
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
  * [X] Karpathy rules (`.agents/rules/skillgrid-karpathy-coding.mdc`)
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

Ship high-impact workflow upgrades: enforceable pipelines, strict phase gates, model routing, persona-board decisions, and Norse-themed operator clarity where it helps.

### Workflow backlog

- [-] Enforcement & review (partial — see Milestone 1):
  - [-] branch-finish protocol (`verify -> merge/PR/keep/discard -> cleanup`)
    - [X] explicit post-merge index refresh (`ccc index`, `npx gitnexus analyze`)
  - [ ] optional git-worktree execution mode for risky/parallel slices

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
