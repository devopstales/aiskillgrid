# Skillgrid workflow

**Runnable steps** live in the slash commands (e.g. `.cursor/commands/skillgrid-*.md`, mirrored under `.kilo/commands/`, `.opencode/commands/`, `.github/prompts/`). The sections below are a compact index; open the matching `skillgrid-*` file for the full checklist and skill paths.

## Session (Optional)

* `/skillgrid-session`
* It’s a pre‑condition you run at the start of any agent session before entering a workflow phase.

## Init

* `/skillgrid-init`
* Create skillgrid folder structure
  * Initialize graphify
  * Initialize engram
  * Initialize openspec
* Detect if the project is greenfield or brownfield
  * If brownfield:
    * Propose Exploration
  * If greenfield:
    * Ask want is the goal of this project -> AGENTS.md
    * Ask want technologie shod be use -> ARCHITECTURE.md
    * Ask how shod the app lokks like -> DESIGN.md
      * 

## Explore (Optional)

* `/skillgrid-explore`
* Explore existing code
  * Parallel discovery
* Generate content for ARCHITECTURE.md, STRUCTURE.md, DESIGN.md
* Use openspec-explore
  * Geberate PRD from OpenSpec change and vica versa
  * Genarate Tasks on PRD from OpenSpec spec in change and vica gersa

## Design (Optional)

* `/skillgrid-design`
* It’s an on-demand design workshop that bootstraps, tunes, and audits the visual identity of the project.

## Brainstorm

* `/skillgrid-brainstorm`
* User describes what they want
* Ask back to clarify; diverge and converge
* Generate Preview to pick
* Document Architectural Decisions -> ARCHITECTURE.md

## Plan

* `/skillgrid-plan`
* Generate a PRD for the feature
  * `.skillgrid/prd/PRD01_<first-slug>.md`
* Use the PRD to drive the OpenSpec change
* Ask User to validate Plan

## Breakdown

* `/skillgrid-reakdown`
* Break down the PRD into tasks
* Create tasks under OpenSpec change.

## Aplly 

* `/skillgrid-apply`
* *subagent-driven-development* or *executing-plans* (Phase 2)
* Create workre in separate branch for paralel execution
  * `.worktree/<first-slug>/`
* reindex ater changes
* Test Driven Workflow

## Test

* `/skillgrid-test`
* Automated Code Quality Testing
* Automated Functional Testing
* Automated Security testing

## Security (Optional)

* `/skillgrid-security`
* dedicated, deeper security vlidate

## Validate

* `/skillgrid-validate`
* Ask User to validate
* roll back atmic change if nececary

## Finish

* `/skillgrid-finish`
* Archive change - openspec-archive-change
* Optional: sync delta specs to main specs - openspec-sync-specs
* cleanup preview
* Create PR

# Skillgrid layout (`.skillgrid/`)

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
│       ├── prd-kanban.mjs          # optional PRD / kanban helper
│       └── preview.sh              # scaffolds non-destructive stubs under preview/
├── .cursor/
│   └── commands/
├── openspec/                      # hybrid: disk + Engram
│   ├── config.yaml
│   ├── specs/
│   ├── changes/
│   └── changes/archive/
├── graphify-out/
├── .worktree/
│   └── <first-slug>/
└── src/ or app/ or lib/
```

## PRD (product requirements)

In Skillgrid, a **PRD** is the **human-facing slice of intent**: problem, goals, in/out of scope, requirements, and success criteria. It **does not** replace **`openspec/changes/<id>/tasks.md`** (that is the implementable checklist from **`/skillgrid-breakdown`**).

A PRD may map to one or more OpenSpec changes (commonly 1:1).
If the scope expands, add references under a “Related changes” heading.
The PRD’s `Status` field can be reset to an earlier stage if the plan is revised; update both the PRD and `INDEX.md`.

**Location and naming**

- **Canonical path:** **`.skillgrid/prd/PRD<NN>_<slug>.md`** — **`<NN>`** is a two-digit execution order (`01`, `02`, …).  
- **Index:** **`.skillgrid/prd/INDEX.md`** — one row or bullet per PRD, sorted by `NN`, with links to the matching OpenSpec change when it exists.  
- **Do not** create new PRDs under a root `prd/` folder; use **`.skillgrid/prd/`** only (see **`/skillgrid-init`**).  
- Title block in the PRD should list **file**, **Spec / change** (path under `openspec/changes/…`), optional **Session context** (`.skillgrid/tasks/context_<change-id>.md`), and **`Status:`** (see table below).

**PRD as source of product intent** until superseded: keep the PRD **consistent** with `openspec/changes/<id>/proposal.md` and delta specs when both exist. Detailed file-by-file steps belong in **`tasks.md`**, not in the PRD body.

**Status lifecycle (`Status:` on the PRD and in `INDEX.md`)** — advance in lockstep with the workflow commands:

| When this command completes (phase) | Set `Status` to |
|------------------------------------|-----------------|
| **`/skillgrid-plan`** | `draft` |
| **`/skillgrid-breakdown`** | `todo` |
| **`/skillgrid-apply`** | `inprogress` |
| **`/skillgrid-vlidate`** | `devdone` |
| **`/skillgrid-finish`** | `done` |

- Use **single-token** values (`inprogress` has no space). Overwrite the previous value as the work moves forward.  
- Authoritative rules and edge cases: **`/skillgrid-init`** — **PRD / change `Status` lifecycle**.  
- Markdown **skeletons** for the PRD title block and sections live in **`/skillgrid-plan`** (Part A — *PRD file templates*), not duplicated here; see *Formatting templates* at the end of this file.

## OpenSpec

**OpenSpec** in this workflow is the **on-disk spec system** for a change: `openspec/changes/<change-id>/` holds **proposal**, **delta specs**, **`tasks.md`**, and related artifacts; `openspec/specs/` holds **main** (accumulated) specs where your project uses them; `openspec/changes/archive/` stores **completed** changes after finish/archive.

**Change id** — The directory name **`openspec/changes/<change-id>/`** (kebab-case) is the same **`<change-id>`** used in **`.skillgrid/tasks/context_<change-id>.md`** and in `openspec list`. Tie every Skillgrid handoff to **one** change id per active slice of work.

**How PRD and OpenSpec connect**

1. **`/skillgrid-plan`** — Write or update the **PRD** first, then create or refresh the **OpenSpec change** (Part B in the command). Set PRD **Status** to `draft`. **`/skillgrid-plan`** may also **`mem_save`** to Engram with a stable `topic_key` pointing at the change id and PRD path (**hybrid persistence**: disk + Engram).  
2. **`/skillgrid-breakdown`** — Decompose into **`tasks.md`**; align with PRD; PRD **Status** → `todo`.  
3. **`/skillgrid-apply`** — Implement from **`tasks.md`**; **Status** → `inprogress`.  
4. **`/skillgrid-vlidate` / `/skillgrid-validate`** — Verify against specs and tasks; **Status** → `devdone`.  
5. **`/skillgrid-finish`** — **`openspec-archive-change`** (and optionally **`openspec-sync-specs`**); open PR; **Status** → `done`.

**Context for agents** — OpenSpec **proposal** may list **`contextFiles`**. The Skillgrid handoff file is **additional** filesystem context: add a line in **`proposal.md`**: *Skillgrid session context:* `.skillgrid/tasks/context_<change-id>.md` (see **`/skillgrid-plan`**). The orchestrator and subagents should **read the handoff** when it exists, not only `contextFiles`.

**CLI** — Use the OpenSpec CLI as your project documents it (e.g. `openspec status`, `openspec instructions tasks` during breakdown). Optional skills: **`openspec-*`** under `.agents/skills/` (e.g. `openspec-apply-change`, `openspec-verify-change`).

**Persistence modes** (repo-dependent) — e.g. **hybrid** (disk + Engram), **openspec**-only, **engram**-only, **none**; see **`/skillgrid-init`** when bootstrapping `openspec/` and `.skillgrid/`.