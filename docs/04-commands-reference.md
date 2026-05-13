# Skillgrid Commands & Workflows Reference

This document lists every slash command, what it does, and which skills it uses.

## Quick Reference

| Command | Phase | Description |
|---------|-------|-------------|
| `/sdd-init` | Setup | Initialize SDD context, detect stack, bootstrap persistence |
| `/sdd-brainstorm <name>` | Planning | Full planning pipeline: explore → propose → spec → design → tasks |
| `/sdd-explore <topic>` | Planning | Free-form codebase investigation (no code changes) |
| `/sdd-clarify <name>` | Planning | Interactive questioning to sharpen terminology, update CONTEXT.md |
| `/sdd-apply [name]` | Build | Implement tasks from a change (TDD: RED → GREEN → REFACTOR) |
| `/sdd-loop [name]` | Build | Controlled build loop for AFK-safe slices |
| `/sdd-verify [name]` | Verify | Validate implementation against specs with execution evidence |
| `/sdd-archive [name]` | Archive | Sync delta specs, move change to archive |
| `/sdd-adr [topic]` | Specialist | Author or review architecture decision records |
| `/sdd-design-ui [surface]` | Specialist | UI design workshop with high-fidelity skills |
| `/sdd-gherkin <input>` | Specialist | Draft, review, or tighten Gherkin/BDD scenarios |
| `/sdd-c4 <scope>` | Specialist | C4-style architecture diagrams (ASCII or Mermaid) |
| `/sdd-diagnose <bug>` | Specialist | Structured debugging session |
| `/sdd-openspec-git` | Gate | OpenSpec git discipline gates |
| `/sdd-persona-board <decision>` | Board | Multi-perspective decision-making with Norse personas |
| `/sdd-persona-route <type>` | Board | Select personas for a decision type |
| `/sdd-persona-report <id>` | Board | Merge and summarize persona verdicts |
| `/sdd-persona-resolve <id>` | Board | Record accepted decision and rejected options |
| `/sdd-persona-health` | Board | Validate persona prompt packs and readiness |
| `/sdd-persona-list` | Board | List all available personas and roles |

---

## SDD Core Workflow

### `/sdd-init`

**Phase:** Setup (run once per project)

**What it does:** Detects project stack, conventions, and testing capabilities. Bootstraps the active persistence backend (engram, openspec, hybrid, or none). Creates `.skillgrid/` and `openspec/` directory structures. Builds the skill registry. Refreshes semantic indexes.

**Skills used:** `sdd-init`, `skillgrid-skill-registry`, `ccc`, `gitnexus-cli`, `openspec-onboard`

**Input:** None (auto-detects project)

**Output artifacts:**
- `.skillgrid/config.json`
- `.skillgrid/project/ARCHITECTURE.md`
- `.skillgrid/project/PROJECT.md`
- `.skillgrid/project/STRUCTURE.md`
- `.skillgrid/project/SKILL_REGISTRY.md`
- `openspec/config.yaml`
- `openspec/specs/`
- `openspec/changes/`

---

### `/sdd-brainstorm <change-name>`

**Phase:** Planning (start of every new change)

**What it does:** Orchestrates the full planning pipeline. Delegates to sub-agents for each phase in sequence. Integrates UI design sub-flow and optional Beads sync.

**Skills used (sequential pipeline):**

| Step | Skill | Output |
|------|-------|--------|
| 1 | `sdd-explore` | Codebase investigation, approach comparison |
| 2 | `sdd-clarify` | Interactive questioning, CONTEXT.md updates |
| 3 | `sdd-propose` | `proposal.md` — intent, scope, approach, rollback plan |
| 4 | `sdd-spec` | `specs/<domain>/spec.md` — delta specs with Given/When/Then |
| 5 | `sdd-design` | `design.md` — architecture decisions, data flow, file changes |
| 6 | `sdd-adr` (conditional) | ADRs in `.skillgrid/adr/` if architectural decisions needed |
| 7 | `sdd-design-ui` (conditional) | UI design artifacts if user-facing scope |
| 8 | `sdd-prd` | `.skillgrid/prd/PRD<NN>_<slug>.md` — product requirements |
| 9 | `sdd-tasks` | `tasks.md` — implementation checklist with HITL/AFK labels |
| 10 | `beads-sync` (conditional) | Beads epic + issues if `beads_enable: true` |

**Input:** Change name (kebab-case)

**Output artifacts:**
- `openspec/changes/<name>/proposal.md`
- `openspec/changes/<name>/specs/<domain>/spec.md`
- `openspec/changes/<name>/design.md`
- `openspec/changes/<name>/tasks.md`
- `.skillgrid/prd/PRD<NN>_<name>.md`
- Optional: ADRs, Beads epic/issues, UI previews

---

### `/sdd-explore <topic>`

**Phase:** Planning (standalone investigation)

**What it does:** Free-form exploration of a topic or feature idea. Reads codebase, compares approaches, provides recommendation. Makes NO code changes.

**Skills used:** `sdd-explore`, `parallel-delegate`, `exa-search`

**Input:** Topic to explore

**Output:** Exploration findings (engram topic or `exploration.md`)

---

### `/sdd-clarify <change-name>`

**Phase:** Planning (interactive questioning)

**What it does:** Challenges the plan against the existing domain model. Sharpens terminology, resolves ambiguity, updates `.skillgrid/project/CONTEXT.md`.

**Skills used:** `sdd-clarify`

**Input:** Change name

**Output:** Updated `CONTEXT.md`, clarified terminology

---

### `/sdd-apply [change-name]`

**Phase:** Implementation

**What it does:** Implements tasks from an SDD change. Follows TDD (RED → GREEN → REFACTOR) when `strict_tdd: true`. Validates task labels before starting.

**Skills used:** `sdd-apply`, `skillgrid-tdd`

**Input:** Optional change name (picks active change if omitted)

**Preflight gate:**
1. `.skillgrid/scripts/validate-task-labels.sh` must pass
2. Checks `[Label: AFK|HITL]` and `[Budget: safe|RISK]` per slice
3. HITL slices stop for human decision; RISK slices warn and ask

**Output:** Updated codebase, tests, `[x]` marks on `tasks.md`

---

### `/sdd-loop [change-name]`

**Phase:** Implementation (controlled continuation)

**What it does:** Picks one AFK-safe slice per iteration, implements it, captures verification evidence, then reassesses risk. Stops deterministically on gate failures. Safer for long-running sessions where context might be lost.

**Skills used:** `sdd-apply` (delegated per slice), `sdd-persona-board` (escalation on risk)

**Input:** Optional change name

**Output artifacts:**
- Updated `.skillgrid/tasks/context_<change-id>.md`
- Event log `.skillgrid/tasks/events/<change-id>.jsonl`

---

### `/sdd-verify [change-name]`

**Phase:** Verification

**What it does:** Multi-stage verification that implementation matches specs, design, and tasks. Builds a spec compliance matrix, runs tests and build. Can escalate to persona board for critical findings.

**Skills used:** `sdd-verify`, `openspec-verify-change` (optional delegate), `beads-retrospective` (post-verify)

**Verification stages:**

| Stage | Check |
|-------|-------|
| 1 | Completeness — all tasks done? |
| 2 | Label contract — HITL/AFK valid? (`validate-task-labels.sh`) |
| 3 | Correctness — static spec matching |
| 4 | Coherence — design decisions followed? |
| 5 | Testing — run tests + build + typecheck |
| 6 | Spec compliance — behavioral validation (Given/When/Then) |
| 7 | Pre-done gate — tests, regression, docs, persistence |

**Verdict:** PASS / PASS WITH WARNINGS / FAIL

**Output:** Verification report with CRITICAL / WARNING / SUGGESTION findings

---

### `/sdd-archive [change-name]`

**Phase:** Archiving (after verification passes)

**What it does:** Syncs delta specs from the change to main specs. Moves the change folder to `archive/` with a date prefix.

**Skills used:** `sdd-archive`, `openspec-sync-specs`, `openspec-archive-change`

**Input:** Optional change name

**Output:**
- Updated `openspec/specs/` (main specs)
- `openspec/changes/archive/YYYY-MM-DD-<name>/`
- Archive report (engram)

---

## Specialist Commands

### `/sdd-adr [topic]`

**What it does:** Authors or reviews architectural decision records using MADR (Markdown Any Decision Record) format with hub templates.

**Skills used:** `architectural-decision-records`

**Input:** Decision topic or file path

**Output:** ADR file in `.skillgrid/adr/NNNN-kebab-title.md`

---

### `/sdd-design-ui [surface]`

**What it does:** UI design workshop. Generates visual direction options, compares across accessibility, responsive, cost, and design-system fit. Optionally produces high-fidelity previews.

**Skills used:** `sdd-ui-design`, `sdd-design`, `engram-ui-elements`, `engram-visual-language`, `design-taste-frontend`, `frontend-ui-engineering`, `high-end-visual-design`, `impeccable`, `superdesign`, `huashu-design`, `image-to-code`

**Input:** Surface or change name

**Output:** UI design direction, preview artifacts (via `preview.sh`)

---

### `/sdd-gherkin <input>`

**What it does:** Drafts, reviews, or tightens Gherkin/BDD scenarios and acceptance criteria.

**Skills used:** `gherkin-authoring`

**Input:** Feature file path, change name, or inline scenario text

**Output:** Updated `.feature` files or spec excerpts

---

### `/sdd-c4 <scope>`

**What it does:** Generates C4-style architecture diagrams (ASCII or Mermaid) for systems and codebases.

**Skills used:** `c4-diagrams`

**Input:** System or scope description

**Output:** ASCII or Mermaid diagrams

---

### `/sdd-diagnose <bug>`

**What it does:** Structured debugging session — reproduce, isolate, root-cause, fix, verify. For bugs, unexpected behavior, or test failures.

**Skills used:** Built-in diagnostic loop (no external skill)

**Input:** Bug report or issue description

**Output:** Root cause analysis, fix, verification evidence

---

### `/sdd-openspec-git`

**What it does:** Enforces OpenSpec git gates — ensures proposal/apply/archive cross main branch correctly, commits are explicit. Never auto-commits or merges.

**Skills used:** `openspec-git-discipline`

**Input:** Optional context or change name

**Output:** Git facts, gate outcomes, next recommended action

---

## Persona Board Commands

The Norse persona board provides multi-perspective decision-making with hard-gate enforcement.

### `/sdd-persona-board <decision> [--preset <alias>]`

**What it does:** Runs a complete persona board cycle — defines decision scope, resolves routing preset, dispatches selected personas in parallel, merges findings, persists artifacts.

**Skills used:** `sdd-persona-route` (routing), individual persona prompt packs

**Presets:**

| Preset | Personas | Use case |
|--------|----------|----------|
| `arch` | Tyr, Mimir | Architecture decisions |
| `security` | Heimdall, Vidar | Security review |
| `ux` | Frigg | UX/content decisions |
| `release` | Tyr, Heimdall | Go/no-go release gate |
| `risk` | Loki, Tyr | Risk acceptance |
| `debug` | Vidar, Mimir | Debugging strategy |
| `bootstrap` | Mimir | Init/readiness check |

**Input:** Decision ID/question, optional preset alias

**Output:** Persona reports in `.skillgrid/tasks/research/<change-id>/`, decision record

---

### `/sdd-persona-route <decision-type>`

**What it does:** Selects the right Norse personas for a given decision type. Maps decision types to persona sets.

**Input:** Decision type (architecture, security, ux-content, go-no-go-release, risk-acceptance, bootstrap-readiness, spec-quality, tasks-readiness, debugging)

**Output:** Selected personas with rationale

---

### `/sdd-persona-report <decision-id>`

**What it does:** Merges and summarizes persona verdicts for one decision. Detects conflicts and classifies severity.

**Input:** Decision ID

**Output:** Merged report, handoff, event log

---

### `/sdd-persona-resolve <decision-id>`

**What it does:** Records the accepted decision and rejected options from persona board output.

**Input:** Decision ID

**Output:** Updated handoff, event log

---

### `/sdd-persona-health`

**What it does:** Validates Norse persona prompt packs, model routing, and surface readiness.

**Input:** None

**Output:** Health report with blockers

---

### `/sdd-persona-list`

**What it does:** Lists all Norse personas, their mapped roles, and runtime availability.

**Input:** Optional filter

**Output:** Persona registry listing

---

## Utility Skills (Auto-Loaded by Context)

These skills are not invoked directly as commands but are auto-loaded when the orchestrator detects relevant context:

| Skill | Trigger Context |
|-------|-----------------|
| `skillgrid-tdd` | Writing tests, `--tdd` flag, RED/GREEN/REFACTOR |
| `skillgrid-vertical-slices` | Breaking down work, creating issues, planning implementation |
| `skillgrid-skill-registry` | Initializing Skillgrid, refreshing project context |
| `micro-plan` | Quick fix, "few steps", "simple plan", small operational plan |
| `parallel-delegate` | Multiple independent files, parallel research passes |
| `exa-search` | Web search, code research, company intel |
| `ccc` | Semantic code search, indexing after changes |
| `gitnexus-cli` | Analyzing/indexing repos, generating wikis |
| `full-output-enforcement` | Tasks requiring exhaustive, unabridged output |
| `markdown-converter` | Converting PDF, DOCX, PPTX, XLSX to markdown |
| `karpathy-guidelines` | Writing/reviewing/refactoring code to avoid LLM pitfalls |

## Engram Guardrail Skills (Auto-Loaded by Project)

These enforce project-specific conventions for the Engram project. They are auto-loaded when relevant files are touched:

| Skill | Trigger |
|-------|---------|
| `engram-project-structure` | Creating files, packages, handlers, templates, tests |
| `engram-commit-hygiene` | Commit creation, review, branch cleanup |
| `engram-testing-coverage` | Implementing behavior changes |
| `engram-server-api` | Route, handler, payload, status code changes |
| `engram-dashboard-htmx` | HTMX attributes, partial updates, forms |
| `engram-ui-elements` | Dashboard UI components, connected browsing flows |
| `engram-visual-language` | Typography, spacing, color, visual identity changes |
| `engram-tui-quality` | Bubbletea/Lipgloss model, update, view changes |
| `engram-gentleman-bubbletea` | Go files in `installer/internal/tui/` |
| `engram-architecture-guardrails` | System boundaries, ownership, state flow changes |
| `engram-business-rules` | Sync behavior, project controls, permissions |
| `engram-plugin-thin` | Plugin scripts/hooks changes |
| `engram-docs-alignment` | Code changes affecting user/contributor behavior |
| `engram-cultural-norms` | Starting work, reviewing changes, defining conventions |
| `engram-memory-protocol` | Decisions, bugfixes, discoveries, session closure |
| `engram-backlog-triage` | Auditing issues/PRs, triaging backlog |
| `engram-issue-creation` | Creating GitHub issues, reporting bugs |
| `engram-branch-pr` | Creating PRs, opening PRs, preparing branches |
| `engram-pr-review-deep` | Reviewing PRs as reviewer |
| `engram-sdd-flow` | SDD or multi-phase implementation planning |

---

## Workflow Diagram

```
/sdd-init
    │
    ▼
/sdd-brainstorm <name>
    ├── sdd-explore     → codebase investigation
    ├── sdd-clarify     → terminology, CONTEXT.md
    ├── sdd-propose     → proposal.md
    ├── sdd-spec        → specs/<domain>/spec.md
    ├── sdd-design      → design.md
    ├── sdd-adr         → ADRs (if architectural)
    ├── sdd-design-ui   → UI artifacts (if user-facing)
    ├── sdd-prd         → PRD<NN>_<name>.md
    ├── sdd-tasks       → tasks.md (HITL/AFK labels)
    └── beads-sync      → Beads epic + issues (if enabled)
    │
    ▼
/sdd-loop <name>  ──or──  /sdd-apply <name>
    │   (controlled              (direct
    │    continuation)           implementation)
    │
    ▼
/sdd-verify <name>
    ├── completeness check
    ├── label contract (validate-task-labels.sh)
    ├── static spec matching
    ├── design coherence
    ├── tests + build + typecheck
    ├── spec compliance matrix
    └── pre-done gate
    │
    ▼  (PASS)
/sdd-test
    │
    ▼  (PASS)
/sdd-archive <name>
    ├── sync delta specs → openspec/specs/
    └── move → openspec/changes/archive/
    │
    ▼
beads-retrospective
    ├── analyze patterns, tech debt
    └── suggest new OpenSpec proposals
```

---

## Utility Scripts

Located in `.skillgrid/scripts/`:

| Script | Purpose |
|--------|---------|
| `validate-task-labels.sh` | Validates `[Label: AFK|HITL]` and `[Budget: safe|RISK]` in tasks.md |
| `preview.sh` | Scaffolds preview artifacts for UI design |
| `handoff-create.sh` | Creates a new handoff record |
| `handoff-resume.sh` | Resumes from a handoff record |
| `handoff-validate.sh` | Validates handoff record integrity |
| `handoff-check-staleness.sh` | Checks if handoff is stale |
| `handoff-list.sh` | Lists available handoff records |
| `handoff-registry-init.sh` | Initializes handoff registry |
