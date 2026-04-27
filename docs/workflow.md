# Skillgrid workflow

Runnable steps live in the slash command files (for example `.cursor/commands/skillgrid-*.md`, mirrored under `.kilo/commands/`, `.opencode/commands/`, `.github/prompts/`). The sections below are a compact index; open the matching `skillgrid-*` file for the full checklist and templates.

## Session (optional)

* `/skillgrid-session`
* Run at the start of an agent session when you need charter, context budget, MCP selection, and checkpoints. Restores **`.skillgrid/config.json`** (ticketing and artifact store) with **`AGENTS.md`**, **OpenSpec** listing, and PRDs in flight.

## Init

* `/skillgrid-init`
* **`.skillgrid/config.json`** (created or merged early):
  * **Question 1 — Ticketing:** `local` (PRDs + index, optional local Kanban script), or `github` / `gitlab` / `jira` for remote issue workflows in later commands.
  * **Question 2 — Artifact store:** `hybrid` (default), `openspec`, or `engram` — where spec and handoff data live: on-disk `openspec/`, Engram memory, or both. Same contract as the hub’s spec-driven init skills: `engram` does not add `openspec/` by default; `openspec` is disk-first; `hybrid` uses both.
* Create **`.skillgrid/`** tree: **`project/`**, **`prd/`**, **`tasks/`**, **`preview/`**, **`scripts/`** as needed; never put PRDs at the repository root.
* **Greenfield / brownfield** routing, **`DESIGN.md`**, **`.skillgrid/project/*.md`**, **root `AGENTS.md`**
* When artifact store includes OpenSpec: **`openspec init`** (baseline `openspec/`, **`config.yaml`**, changes tree). When it includes Engram: **`mem_save`** with a stable `topic_key` (for example `skillgrid-init/{project-name}`).
* **Graphify** (optional, project policy)
* If brownfield, recommend **`/skillgrid-explore`** before large structural work.

## Explore (optional)

* `/skillgrid-explore`
* Map the codebase; parallel discovery; optional research memos under **`.skillgrid/tasks/research/<change-id>/`**
* Refresh **`.skillgrid/project/`** (`ARCHITECTURE`, `STRUCTURE`, `PROJECT`) and **root `DESIGN.md`**
* PRD backfill and coverage vs **OpenSpec** changes; align **`INDEX.md`**; optional **External** column when ticketing is not **`local`**

## Design (optional)

* `/skillgrid-design`
* On-demand design workshop: extract, browse, tune, SuperDesign, Impeccable — converges on root **`DESIGN.md`**

## Brainstorm

* `/skillgrid-brainstorm`
* Clarify goals; diverge and converge; optional previews in **`.skillgrid/preview/`**
* Record architecture decisions in **`.skillgrid/project/ARCHITECTURE.md`**
* No remote issues unless **`.skillgrid/config.json`** says so (see init).

## Plan

* `/skillgrid-plan`
* PRD first: **`.skillgrid/prd/PRD<NN>_<slug>.md`**, **`.skillgrid/prd/INDEX.md`**
* **OpenSpec** change scaffold and artifact loop; hybrid **Engram** handoff
* Optional remote issues only per **`ticketing.provider`** in **`.skillgrid/config.json`**

## Breakdown

* `/skillgrid-breakdown`
* Sync PRD **Implementation tasks** with **`openspec/changes/<id>/tasks.md`**
* Optional vertical-slice issues per ticketing provider; **`local`** keeps work in PRD + `tasks.md` only

## Apply

* `/skillgrid-apply`
* Implement from `tasks.md` / **OpenSpec** apply instructions; TDD when appropriate; worktrees (for example **`.worktree/<slug>/`**) if the project uses them

## Test

* `/skillgrid-test`
* Quality gates, functional tests, automated security baselines; tie to PRD success criteria. **Issue id** in arguments follows the configured ticketing provider.

## Security (optional)

* `/skillgrid-security`
* Deeper review than the Test-phase scanners alone

## Validate

* `/skillgrid-validate`
* Single gate: review + security + user sign-off; **rollback** path if the user rejects

## Finish

* `/skillgrid-finish`
* Archive **OpenSpec** change, optional sync of delta specs to main **specs**, preview cleanup, PR/merge, **PRD** **`Status: done`**
* Remote tracker hints (Git **issue** / merge keywords / Jira transition) only when **`artifactStore`** and ticketing imply an external system; **`local`** stays file-first.

---

# Skillgrid layout (`.skillgrid/`)

Canonical on-disk tree for a project using Skillgrid. Application code lives under `src/`, `app/`, or your stack’s layout elsewhere.

```text
project-root/
├── AGENTS.md
├── DESIGN.md
├── .skillgrid/
│   ├── config.json                 # ticketing.provider; artifactStore.mode (hybrid | openspec | engram)
│   ├── project/
│   │   ├── ARCHITECTURE.md
│   │   ├── STRUCTURE.md
│   │   └── PROJECT.md
│   ├── prd/
│   │   ├── INDEX.md
│   │   ├── PRD01_<first-slug>.md
│   │   └── PRD02_<next-slug>.md
│   ├── tasks/
│   │   ├── context_<change-id>.md
│   │   └── research/
│   │       └── <change-id>/
│   │           └── <topic-or-agent>_<optional-date>.md
│   ├── preview/
│   └── scripts/
│       ├── prd-kanban.mjs
│       └── preview.sh
├── .cursor/
│   └── commands/
├── openspec/                       # when artifact store is hybrid or openspec
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

In Skillgrid, a **PRD** is the **human-facing slice of intent**: problem, goals, in/out of scope, requirements, and success criteria. It does not replace **`openspec/changes/<id>/tasks.md`** (the implementable checklist from **`/skillgrid-breakdown`**).

A PRD may map to one or more **OpenSpec** changes (commonly 1:1). If scope expands, add references under a “Related changes” heading. The PRD’s **`Status:`** can move backward if the plan is revised; update both the PRD and **`INDEX.md`**.

**Location and naming**

- **Canonical path:** **`.skillgrid/prd/PRD<NN>_<slug>.md`** — **`<NN>`** is a two-digit execution order (`01`, `02`, …).  
- **Index:** **`.skillgrid/prd/INDEX.md`** — one row or bullet per PRD, sorted by `NN`, with links to the matching **OpenSpec** change when it exists.  
- **Do not** create new PRDs under a root `prd/` folder; use **`.skillgrid/prd/`** only (see **`/skillgrid-init`**).  
- Title block: **file**, **Spec / change** (path under `openspec/changes/…`), optional **Session context** (`.skillgrid/tasks/context_<change-id>.md`), and **`Status:`** (table below).

**PRD** stays the product intent source until superseded: keep it **consistent** with `openspec/changes/<id>/proposal.md` and delta specs when both exist. Detailed file-by-file steps belong in **`tasks.md`**, not in the PRD body.

**Status lifecycle (`Status:` on the PRD and in `INDEX.md`)** — align with the workflow commands (authoritative list in **`/skillgrid-init`**):

| When this command completes (phase) | Set `Status` to |
|------------------------------------|-----------------|
| **`/skillgrid-plan`** | `draft` |
| **`/skillgrid-breakdown`** | `todo` |
| **`/skillgrid-apply`** | `inprogress` |
| **`/skillgrid-validate`** | `devdone` |
| **`/skillgrid-finish`** | `done` |

- Use **single-token** values (`inprogress` has no space). Overwrite the previous value as work advances.  
- Markdown **skeletons** for the PRD live in **`/skillgrid-plan`** (Part A — *PRD file templates*), not duplicated here.

## OpenSpec

**OpenSpec** in this workflow is the **on-disk** spec system for a change: `openspec/changes/<change-id>/` holds **proposal**, **delta specs**, **`tasks.md`**, and related artifacts; `openspec/specs/` holds **main** specs where the project uses them; `openspec/changes/archive/` stores **completed** changes after finish.

**Change id** — The directory name **`openspec/changes/<change-id>/`** (kebab-case) matches **`.skillgrid/tasks/context_<change-id>.md`** and `openspec list` output. Tie each Skillgrid handoff to **one** change id per active slice of work.

**How PRD and OpenSpec connect**

1. **`/skillgrid-plan`** — PRD first, then create or refresh the **OpenSpec** change. **Status** → `draft`. Optional **Engram** `mem_save` with a stable `topic_key` (hybrid or engram **artifact** store).  
2. **`/skillgrid-breakdown`** — **`tasks.md`**; PRD **Status** → `todo`.  
3. **`/skillgrid-apply`** — **Status** → `inprogress`.  
4. **`/skillgrid-validate`** — **Status** → `devdone`.  
5. **`/skillgrid-finish`** — archive change, optional spec sync, PR; **Status** → `done`.

**Context for agents** — The **OpenSpec** proposal may list `contextFiles`. The Skillgrid handoff file is **additional** context: one line in **`proposal.md`**: *Skillgrid session context:* `.skillgrid/tasks/context_<change-id>.md` (see **`/skillgrid-plan`**).

**CLI** — Use the **OpenSpec** CLI as the project documents (`openspec status`, `openspec instructions tasks` during breakdown). **Artifact store** for whether **`openspec/`** exists: **`/skillgrid-init`** and **`.skillgrid/config.json` → `artifactStore.mode`**.

**Persistence** — **`artifactStore.mode`**: **`hybrid`** (disk + Engram), **`openspec`**, or **`engram`**. The init command records this; **`/skillgrid-session`** loads it for cold starts.
