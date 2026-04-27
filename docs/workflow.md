# Skillgrid workflow

Runnable steps live in the slash command files (for example `.cursor/commands/skillgrid-*.md`, mirrored under `.kilo/commands/`, `.opencode/commands/`, `.github/prompts/`). The sections below are a compact index; open the matching `skillgrid-*` file for the full checklist and templates.

## Session (optional)

* `/skillgrid-session`
* Run at the start of an agent session when you need charter, context budget, MCP selection, and checkpoints. Restores **`.skillgrid/config.json`** (ticketing and artifact store) with **`AGENTS.md`**, **OpenSpec** listing, and PRDs in flight.

## Checkpoint (optional)

* `/skillgrid-checkpoint`
* Create, verify, list, or clear named workflow checkpoints in **`.skillgrid/tasks/checkpoints.log`**. Checkpoints record branch, git SHA, dirty status, optional quick verification evidence, and active handoff files. They complement PRDs, OpenSpec changes, and **`context_<change-id>.md`**; they do **not** require git worktrees and do **not** create commits unless the user explicitly asks.

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
* Sync PRD **Implementation tasks** with **`openspec/changes/<id>/tasks.md`**; tag tasks with **`[HITL]`** or **`[AFK]`** (see *HITL vs AFK slices* below); order HITL decisions before dependent AFK work; **prefer AFK** when scoping slices
* Optional vertical-slice issues per ticketing provider; **`local`** keeps work in PRD + `tasks.md` only

## Apply

* `/skillgrid-apply`
* Before every apply run that proceeds to implementation, automatically create a named **`/skillgrid-checkpoint create before-apply-<change-id>`** entry in **`.skillgrid/tasks/checkpoints.log`**.
* Implement from `tasks.md` / **OpenSpec** apply instructions; TDD when appropriate; **per-change handoff** (`.skillgrid/tasks/context_<change-id>.md`) in a **single working tree**—optional feature branch in that clone is fine; do **not** assume git worktrees

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
* Archive **OpenSpec** change, optional sync of delta specs to main **specs**, preview cleanup, change-scoped checkpoint cleanup, PR/merge, **PRD** **`Status: done`**
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
│   │   ├── checkpoints.log
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
└── src/ or app/ or lib/
```

Optional (not part of the default Skillgrid path): a team may use an extra **git worktree** or a **`.worktree/<slug>/`** directory for their own process—the hub does not require it; isolation is **change id** + handoff + OpenSpec, not a second checkout.

## Filesystem handoff (per-change)

Keep **change-scoped** state on disk so the **orchestrator session** and **`Task` subagents** stay aligned without pasting long tool output in chat. Full behavioral rules: [`.configs/AGENTS.md`](../.configs/AGENTS.md) (Skillgrid section).

| Path | Role |
|------|------|
| **`.skillgrid/tasks/context_<change-id>.md`** | Rolling handoff: goal, state, index of subagent work, HITL blockers (see below) |
| **`.skillgrid/tasks/research/<change-id>/`** | Long research, scrapes, or subagent reports (one file per topic or run) |

**Subagent contract**

1. **Before work:** Read `context_<change-id>.md` (and any cited `research/…` files).  
2. **During work:** Write lengthy output under `research/<change-id>/` with a **distinct** filename.  
3. **After work:** Update the handoff (state, index row, next actions).  
4. **Return to parent:** A short message with file paths to read before product code changes.

**Cross-link** — In `openspec/changes/<change-id>/proposal.md`, include a line the reader cannot miss, e.g. `**Skillgrid session context:** .skillgrid/tasks/context_<change-id>.md`.

### Handoff file skeleton (copy and trim)

```markdown
# Session context: <change-id>

## Change / links
- OpenSpec change: `openspec/changes/<change-id>/`
- PRD: `.skillgrid/prd/...`
- Branch (optional): …

## Current goal
…

## State
planning | research | todo | inprogress | blocked

## HITL / human gates
- … or *none — ready for AFK apply*

## Subagent / research index
| Topic | File |
|------|------|
| … | `.skillgrid/tasks/research/<change-id>/…` |

## Last checkpoint
…
```

## Parallel discovery

When **independent** investigations can run without shared mutable state, multiple subagents (or parallel tool use) may each produce a **different** file under `.skillgrid/tasks/research/<change-id>/`. The **parent** merges summaries into the handoff before the next phase. If subtasks are **not** independent, run them **sequentially**. See [`.agents/skills_back/references/orchestration-patterns.md`](../.agents/skills_back/references/orchestration-patterns.md).

## HITL vs AFK slices

Classify work in **`tasks.md`** (and the handoff) so humans and agents know what automation is allowed.

| Tag | Meaning |
|-----|---------|
| **`[HITL]`** | Needs a human before implementation or merge (architecture, design sign-off, product/security exception, unclear requirements). |
| **`[AFK]`** | Can be implemented, verified, and merged under your policy **without** a human gate, given a clear spec and design. |

**Rule:** **Prefer AFK** when scoping: narrow HITL to real decisions; use research + spec updates to turn uncertain work into a later **AFK** slice. Tag lines in `tasks.md` with **`[HITL]`** or **`[AFK]`**; order so HITL decisions that block work come **before** dependent AFK tasks. **`/skillgrid-apply`** should implement **AFK** items and **stop** (or require recorded approval in the handoff) when the next item is **HITL** unless a decision is already linked (e.g. ADR, approved note with date). **`/skillgrid-validate`** may still require explicit human sign-off for HITL-heavy changes, per project policy.

## Search vs spec (multi-agent)

- **Parallel search** — OK when research questions are **independent** (e.g. prior art vs API docs). Each subagent writes a **separate** file under `research/<change-id>/`; the parent updates the handoff, then spec work proceeds.  
- **Spec (proposal + delta specs)** — Usually **after** research: one coherent pass (or a single subagent) updates `proposal.md` and delta specs so requirements stay consistent. **Do not** blindly fan out spec writers in parallel unless spec areas are **truly** independent and inputs are frozen.  
- **After implementation** — **Parallel review** is still appropriate (e.g. spec verifier + code reviewer + security) when each report is an independent perspective.

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

**Project `openspec/config.yaml` —** Root **`context`** (injected into artifact instructions) should **mirror** **ticketing** and **artifact store** from **`.skillgrid/config.json`**, with the same merge discipline as **`/skillgrid-init`**, so agents running **`openspec instructions`** do not assume a different issue tracker. **`/skillgrid-plan`** and exploration passes refresh this when needed.

**CLI** — Use the **OpenSpec** CLI as the project documents (`openspec status`, `openspec instructions tasks` during breakdown). **Artifact store** for whether **`openspec/`** exists: **`/skillgrid-init`** and **`.skillgrid/config.json` → `artifactStore.mode`**.

**Persistence** — **`artifactStore.mode`**: **`hybrid`** (disk + Engram), **`openspec`**, or **`engram`**. The init command records this; **`/skillgrid-session`** loads it for cold starts.
