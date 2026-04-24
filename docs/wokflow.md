# Skillgrid workflow

**Runnable steps** live in the slash commands (e.g. `.cursor/commands/skillgrid-*.md`, mirrored under `.kilo/commands/`, `.opencode/commands/`, `.github/prompts/`). The sections below are a compact index; open the matching `skillgrid-*` file for the full checklist and skill paths.

## PRD / change `Status` (stages)

On each **`.skillgrid/prd/PRD<NN>_<slug>.md`** (and the **Status** column in **`.skillgrid/prd/INDEX.md`** or ticket tables), advance **`Status:`** in lockstep with the workflow commands:

| When this command completes (phase) | Set `Status` to |
|------------------------------------|-----------------|
| **`/skillgrid-plan`** | `draft` |
| **`/skillgrid-breakdown`** | `todo` |
| **`/skillgrid-apply`** | `inprogress` |
| **`/skillgrid-review`** | `devdone` |
| **`/skillgrid-finish`** | `done` |

- Single token values (`inprogress` has no space). The agent overwrites the previous value as the work moves forward.
- Authoritative rules and edge cases: **`/skillgrid-init`** — **PRD / change `Status` lifecycle**.

## Skillgrid layout (`.skillgrid/`)

Canonical on-disk tree for a project using Skillgrid (names may vary). Application code lives under `src/`, `app/`, or your stack’s layout elsewhere.

```text
project-root/
├── AGENTS.md
├── DESIGN.md
├── .skillgrid/
│   ├── project/                    # exploration / init: system & onboarding
│   │   ├── ARCHITECTURE.md
│   │   ├── STRUCTURE.md
│   │   └── PROJECT.md
│   ├── prd/                        # optional PRD index and per-change PRDs
│   │   ├── INDEX.md
│   │   ├── PRD01_<first-slug>.md
│   │   └── PRD02_<next-slug>.md
│   ├── tasks/
│   │   ├── context_<change-id>.md  # rolling orchestrator handoff (per OpenSpec change id)
│   │   └── research/
│   │       └── <change-id>/
│   │           └── <topic-or-agent>_<optional-date>.md  # long research / spill; not chat dumps
│   ├── preview/                    # brainstorm: ephemeral MD/HTML to compare and pick options
│   └── scripts/
│       ├── prd-kanban.mjs          # 
│       └── preview.sh              # scaffolds non-destructive stubs under preview/
├── .cursor/
│   └── commands/
├── openspec/                      # hybrid: disk + Engram
│   ├── config.yaml
│   ├── specs/
│   ├── changes/
│   └── changes/archive/
├── graphify-out/
└── src/ or app/ or lib/
```

- **`preview/`** — Short-lived brainstorm files (e.g. layout A vs B). This repo’s **`.skillgrid/preview/.gitignore`** may ignore `*.html`; commit policy is up to the team.  
- **`tasks/context_<change-id>.md`** and **`tasks/research/<change-id>/`** — **Filesystem handoff** for the main session and subagents: see **Filesystem handoff (session context)** below. The **`<change-id>`** is the OpenSpec change directory name (same as `openspec list`).  
- **`scripts/preview.sh`** — Run from project root: `.skillgrid/scripts/preview.sh [slug]` (optional `--md` for markdown). Not every IDE has an embedded browser; users can open files in the editor or a regular browser, or work from labeled A/B in chat.  
- A longer **example** (with IDE folders, `openspec/`, etc.) appears in **`/skillgrid-init`** under **Project structure (example)**. **Copy-paste templates** for `ARCHITECTURE.md`, `STRUCTURE.md`, `PROJECT.md`, and PRD index/skeleton are in the commands, not in this file (see **Formatting templates** below).

## Filesystem handoff (session context)

**Purpose:** Keep a **git-visible, change-scoped** handoff so the **parent** session and **`Task` / subagents** stay aligned without flooding chat or losing work after a subagent turn ends. Large tool output goes to **`tasks/research/…`** (“spill”); the handoff file stays a short index and state. This complements **Engram** (cross-session memory) and **OpenSpec** `contextFiles` (proposal, specs, `tasks.md`).

**Conventions**

- **Handoff file:** **`.skillgrid/tasks/context_<change-id>.md`**
- **Spill directory:** **`.skillgrid/tasks/research/<change-id>/`**
- **ID:** Use the **OpenSpec change id** (kebab-case directory under `openspec/changes/<change-id>/`). Cross-link to **`.skillgrid/prd/PRD<NN>_<slug>.md`** inside the handoff; do not duplicate the full PRD.

**Minimum sections** (start from this skeleton when creating a new handoff):

```markdown
# Session context: <change-id>

## Change / PRD links
- OpenSpec: `openspec/changes/<change-id>/`
- PRD: `.skillgrid/prd/PRD<NN>_<slug>.md`

## Current goal
(One paragraph.)

## State
- **Phase:** planning | research | implementing | blocked
- **Owner:** (optional)

## Subagent / research index
| Date | Persona / agent | Report path | Outcome (one line) |
|------|-----------------|-------------|---------------------|

## Decisions
(Durable choices; point to Engram `mem_save` for cross-session recall if needed.)

## Next actions
(What the parent should do next, e.g. read `research/.../x.md` then run task N from `tasks.md`.)
```

**Link from OpenSpec:** `openspec/changes/<change-id>/proposal.md` should include one line: **Skillgrid session context:** `.skillgrid/tasks/context_<change-id>.md` (see **`/skillgrid-plan`**). PRD title block may repeat the same path for discoverability.

```bash
/skillgrid-init
# Create skillgrid folder structure
# Detect if the project is greenfield or brownfield
## If greenfield: Ask the user abaute the man porpose of the app, what tools, technologys want to use. Ask until it is clear. Then generate content for ARCHITECTURE.md, STRUCTURE.md, PROJECT.md and Update AGENTS.md Prooject chapter
## If brownfield: Tell the user to use /skillgrid-explore because it is a brownfield project.
# index memory
## graphify — `graphify update .` from project root (see AGENTS.md); code discovery via `rg`/IDE search
### Skills (.agents/skills/)
- karpathy-guidelines — assumptions, simplicity, surgical edits 
- skillgrid-init (command) — bootstrap `.skillgrid/`, OpenSpec tree, persistence modes (`hybrid` / `openspec` / `engram` / `none`); see `.cursor/commands/skillgrid-init.md`
- memory-protocol — Engram MCP: when to save, search, close sessions
- context-engineering — Feed agents the right information at the right time - rules files, context packing, MCP integrations

# Optional only if greenfield project
/skillgrid-explore
# Explore existing code
## Use openspec-explore
## Generate content for ARCHITECTURE.md, STRUCTURE.md, PROJECT.md
## Update AGENTS.md Prooject chapter
### Skills (.agents/skills/)
- openspec-explore — explore problem and codebase before a change
- graphify-out + `rg`/IDE search — orient in repo; refresh graph after structural changes
- documentation-and-adrs - Architecture Decision Records, API docs, inline documentation standards - document the why

# Tempolary disabled
/skillgrid-ui-design
# User describes how the page should look
## Generate DESIGN.md
### Skills (.agents/skills/)
- karpathy-guidelines — assumptions, simplicity, surgical edits 
- frontend-design - 
- frontend-ui-engineering — Component architecture, design systems, state management, responsive design, WCAG 2.1 AA accessibility
- skillgrid-plan (command) — when UI work needs a scoped PRD + OpenSpec change
- spec-driven-development — Write a PRD covering objectives, commands, structure, code style, testing, and boundaries before any code

/skillgrid-brainstorm
# User describes what they want
## Ask back to clarify; diverge and converge (see .cursor/commands/skillgrid-brainstorm.md)
## Preview picks (optional): put compare-and-choose MD/HTML in .skillgrid/preview/; run .skillgrid/scripts/preview.sh to scaffold; works in any IDE (open file or use A/B in chat if no browser)
## Search the internet for good examples when it helps
### Skills (.agents/skills/)
- karpathy-guidelines — assumptions, simplicity, surgical edits 
- search-first — research tools and patterns before building
- deep-research — deeper investigation when exploration needs breadth
- idea-refine — Structured divergent/convergent thinking to turn vague ideas into concrete proposals
- documentation-lookup - Use up-to-date library and framework docs via Context7 MCP instead of training data. Activates for setup questions, API references, code examples, or when the user names a framework e.g. React, Next.js, Prisma.

/skillgrid-plan
## Generate prd for the function
## Use prd to generate openspec change (skillgrid-plan Part B)
### Skills (.agents/skills/)
- karpathy-guidelines — assumptions, simplicity, surgical edits 
- spec-driven-development — Write a PRD covering objectives, commands, structure, code style, testing, and boundaries before any code


/skillgrid-breakdown
# Break down the PRD into tasks
# Create tasks under openspec change.
### Skills (.agents/skills/)
- karpathy-guidelines — assumptions, simplicity, surgical edits 
- skillgrid-breakdown (command) — `openspec status` / `openspec instructions tasks`; PRD ↔ `tasks.md` sync
- planning-and-task-breakdown — Decompose specs into small, verifiable tasks with acceptance criteria and dependency ordering
- test-driven-development - Red-Green-Refactor, test pyramid 80/15/5, test sizes, DAMP over DRY, Beyonce Rule, browser testing
- tdd-guide — TDD guidance and patterns
- source-driven-development — Ground every framework decision in official documentation - verify, cite sources, flag what is unverified
- graphify — `graphify update .` after significant code or layout changes when graphify is in use
# - clean-code — readability and maintainability while implementing

/skillgrid-apply
### Skills (.agents/skills/)
- karpathy-guidelines — assumptions, simplicity, surgical edits 
- openspec-apply-change — implement from OpenSpec change tasks
- api-and-interface-design - Contract-first design, Hyrum Law, One-Version Rule, error semantics, boundary validation

/skillgrid-test
### Skills (.agents/skills/)
- karpathy-guidelines — assumptions, simplicity, surgical edits 
- browser-testing-with-devtools - Chrome DevTools MCP for live runtime data - DOM inspection, console logs, network traces, performance profiling
- debugging-and-error-recovery — Five-step triage: reproduce, localize, reduce, fix, guard. Stop-the-line rule, safe fallbacks
- e2e-testing — end-to-end test design and implementation
- e2e-runner — run and troubleshoot E2E suites
- testing-patterns — general testing patterns beyond E2E

/skillgrid-review
# Automated QA testing
# Automated functional testing
# Ask user to validate code really works
### Skills (.agents/skills/)
- karpathy-guidelines — assumptions, simplicity, surgical edits 
- skillgrid-review (command) — OpenSpec verify + Engram + code/product review; see `.cursor/commands/skillgrid-review.md`
- code-review-and-quality — Five-axis review, change sizing 100 lines, severity labels Nit/Optional/FYI, review speed norms, splitting strategies
- code-simplification - Chesterton Fence, Rule of 500, reduce complexity while preserving exact behavior
# clean-code — review for clarity and coupling
- documentation-and-adrs - Architecture Decision Records, API docs, inline documentation standards - document the why
- performance-optimization — Measure-first approach - Core Web Vitals targets, profiling workflows, bundle analysis, anti-pattern detection
#- database-reviewer — PostgreSQL schema, SQL, RLS, performance review

/skillgrid-security
# Security testing
### Skills (.agents/skills/)
- karpathy-guidelines — assumptions, simplicity, surgical edits 
- security-and-hardening — OWASP Top 10 prevention, auth patterns, secrets management, dependency auditing, three-tier boundary system
- security-review - Review code for security vulnerabilities and best practices
- semgrep-security — static analysis with Semgrep
- trivy-security — container/dep scanning with Trivy
- vulnerability-scanner — threat modeling and vuln prioritization (OWASP-oriented)
- security-scan — audit agent/IDE config (e.g. `.claude/`, MCP, hooks)
- deprecation-and-migration - Code-as-liability mindset, compulsory vs advisory deprecation, migration patterns, zombie code removal

/skillgrid-validate
# Combined gate: run the full /skillgrid-review checklist, then the full /skillgrid-security checklist (see skillgrid-validate command)
### Skills (.agents/skills/)
- karpathy-guidelines — assumptions, simplicity, surgical edits
- (plus every skill listed under /skillgrid-review and /skillgrid-security for that session)

/skillgrid-finish
# Archive change - openspec-archive-change
# Optional: sync delta specs to main specs - openspec-sync-specs
# Create PR
### Skills (.agents/skills/)
- karpathy-guidelines — assumptions, simplicity, surgical edits 
- openspec-archive-change — complete and archive the change (per your merge process)
- openspec-sync-specs — promote delta specs without archiving, if your flow needs it
- git-workflow-and-versioning — trunk-style workflow, atomic commits, small changes
- documentation-and-adrs - Architecture Decision Records, API docs, inline documentation standards - document the why
```

## Parallel discovery 

Parallel **subagents** for codebase mapping and domain research, you can **fan out** independent work, then **merge** in the main session:

- **Safe to run in parallel:** read-only **explore** passes on disjoint areas (e.g. different packages), **cited** landscape or prior-art research with non-overlapping briefs (stack vs competitors vs API docs), using personas such as `skillgrid-explore-architect` and `skillgrid-researcher` in separate subagent contexts when your harness allows concurrent `Task` / subagents.
- **Keep sequential:** `/skillgrid-plan` → `/skillgrid-breakdown` so intent stays a single chain; then `/skillgrid-apply` and later gates follow their ordered phases.

**Subagent contract (filesystem handoff)** — Same rules for `skillgrid-researcher`, `skillgrid-design-critic`, and `skillgrid-explore-architect` when spawned as subagents for a **change**:

1. **Before work:** Read **`.skillgrid/tasks/context_<change-id>.md`** (create with the parent if missing).
2. **Scope:** **Research, critique, or mapping only** unless the user explicitly asked this persona to implement. Implementation stays on the **parent** so it retains full repo context.
3. **Spill:** Write long output to **`.skillgrid/tasks/research/<change-id>/*.md`**; keep the chat return to a **short summary plus file paths**.
4. **After work:** Update the handoff (research index row, state, next actions).
5. **Return message to parent:** e.g. “Updated `context_<change-id>.md`; primary report: `…`; read those before continuing.”

When delegating, the **orchestrator must** pass the handoff path in the `Task` prompt. After a subagent returns, the **parent reads** the handoff and cited research files before writing product code. See also **`/skillgrid-plan`** and **`/skillgrid-apply`**.

Parallel fan-out and merge is the same orchestration idea as the **slash command (orchestrator — fan-out)** section in [`.cursor/agents/README.md`](../.cursor/agents/README.md): only when sub-tasks are **independent** (no shared mutable state, no required ordering). The hub’s `/skillgrid-validate` run may still be **sequential in one turn**; true wall-clock parallelism requires a harness with concurrent subagents.

## Formatting templates

Markdown **skeletons** for project and PRD files are maintained in the slash commands (single source of truth), not duplicated in this file.

| What | Where |
|------|--------|
| **PRD / change `Status` stages** (`draft` → `todo` → `inprogress` → `devdone` → `done`) | **PRD / change `Status` (stages)** (above) and **`/skillgrid-init`** — **PRD / change `Status` lifecycle** |
| **`.skillgrid/project/ARCHITECTURE.md`**, **STRUCTURE.md**, **PROJECT.md** | **`/skillgrid-init`** — **Project document templates** |
| **`.skillgrid/prd/INDEX.md`**, PRD skeleton, and full **PRD document format** | **`/skillgrid-plan`** — **Part A** and **PRD file templates (formatting)** |
| **`.skillgrid/tasks/context_<change-id>.md`**, **research/** spill layout | **Filesystem handoff (session context)** (above) and **`/skillgrid-plan`** / **`/skillgrid-apply`** |

**`.skillgrid/preview/`** and **`.skillgrid/scripts/preview.sh`** support **Preview picks** during **`/skillgrid-brainstorm`**. In Skillgrid, **STRUCTURE.md** holds repo tree and, if you use it, deployment topology in an optional section.
