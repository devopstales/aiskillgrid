# Skillgrid workflow

Runnable steps live in the slash command files (for example `.cursor/commands/skillgrid-*.md`, mirrored under `.kilo/commands/`, `.opencode/commands/`, `.github/prompts/`). Commands own phase order, exit status, and completion reports. Reusable procedure details live in `.agents/skills/skillgrid-*/SKILL.md`; see `docs/skills.md` for the Skillgrid primitive catalog.

## Execution Model

Skillgrid is a file-first workflow layer. It borrows the useful state-machine discipline of heavier systems without adopting their runtime, database, cost ledger, or mandatory worktree model. Each active change should have:

* **Current state** — phase, status, active artifacts, blockers, and next recommended action in `.skillgrid/tasks/context_<change-id>.md`.
* **Stop condition** — what must be true before the next phase can start, usually written in the command completion report, PRD success criteria, OpenSpec scenarios, or `tasks.md`.
* **Evidence output** — tests, preview files, review reports, checkpoints, or research files that prove the state change.
* **Reassessment** — after each vertical slice, update the handoff with what changed, evidence, blockers, changed assumptions, and the next slice or command.

Think of the hierarchy as **PRD sequence -> PRD slice -> OpenSpec tasks**. A sequence of ordered PRDs can behave like a milestone or roadmap, each PRD should be a reviewable slice of product intent, and `openspec/changes/<id>/tasks.md` is the task-level implementation checklist for that slice.

## Session (optional)

* `/skillgrid-session`
* Run at the start of an agent session when you need charter, context budget, MCP selection, and checkpoints. Restores **`.skillgrid/config.json`** (ticketing and artifact store) with **`AGENTS.md`**, **OpenSpec** listing, and PRDs in flight.

## Help (optional)

* `/skillgrid-help`
* Explain command order, phase goals, artifact layout, and common next steps. With `current-state`, inspect **`.skillgrid/prd/INDEX.md`**, active **`context_*.md`** files, and **`openspec/changes/`** to recommend the next command without changing files.
* Local dashboard: run **`node .skillgrid/scripts/skillgrid-ui.mjs`**, then open **`http://127.0.0.1:8787`**. It includes **Board**, **Workflow**, and **Subagents** views, preview links, event-log status, and a **Graph** link when **`graphify-out/graph.html`** exists. Full UI docs: [`docs/dashboard.md`](dashboard.md).

## Checkpoint (optional)

* `/skillgrid-checkpoint`
* Create, verify, list, or clear named workflow checkpoints in **`.skillgrid/tasks/checkpoints.log`**. Checkpoints record branch, git SHA, dirty status, optional quick verification evidence, and active handoff files. They complement PRDs, OpenSpec changes, and **`context_<change-id>.md`**; they do **not** require git worktrees and do **not** create commits unless the user explicitly asks.

## Init

* `/skillgrid-init`
* **`.skillgrid/config.json`** (created or merged early):
  * **Question 1 — Ticketing:** `local` (PRDs + index, optional local Kanban script), or `github` / `gitlab` / `jira` for remote issue workflows in later commands.
  * **Question 2 — Artifact store:** `hybrid` (strongly recommended default), `openspec`, or `engram` — where spec and handoff data live: on-disk `openspec/`, Engram memory, or both. Same contract as the hub’s spec-driven init skills: `engram` does not add `openspec/` by default; `openspec` is disk-first; `hybrid` uses both.
* Create **`.skillgrid/`** tree: **`project/`**, **`prd/`**, **`tasks/`**, **`preview/`**, **`scripts/`** as needed; never put PRDs at the repository root.
* **Greenfield / brownfield** routing, **`DESIGN.md`**, **`.skillgrid/project/*.md`**, **root `AGENTS.md`**
* When artifact store includes OpenSpec: **`openspec init`** (baseline `openspec/`, **`config.yaml`**, changes tree). When it includes Engram: **`mem_save`** with a stable `topic_key` (for example `skillgrid-init/{project-name}`).
* **CocoIndex (`ccc`)** (optional): when the CLI is installed, from the repo root run **`ccc init`** if the tree is not yet a CocoIndex project, then **`ccc index`** so **`ccc search`** (and the **`cocoindex-code`** MCP) stay fresh; see [`.agents/skills/ccc/SKILL.md`](../.agents/skills/ccc/SKILL.md). If missing, skip and document; init still succeeds.
* **Graphify** (optional): when the CLI is installed (`uv tool install graphifyy`), run **`graphify .`** from the repository root during init (IDE: **`/graphify .`**) so **`graphify-out/`** exists; afterward use **`graphify update .`** after substantive code changes per **`AGENTS.md`**. If the CLI is missing, skip and document; init still succeeds.
* If brownfield, recommend **`/skillgrid-explore`** before large structural work.

## Explore (optional)

* `/skillgrid-explore`
* Map the codebase; parallel discovery; optional research memos under **`.skillgrid/tasks/research/<change-id>/`**
* Refresh **`.skillgrid/project/`** (`ARCHITECTURE`, `STRUCTURE`, `PROJECT`) and **root `DESIGN.md`**
* Automatically import/backfill existing PRDs and **OpenSpec** changes that lack canonical **`.skillgrid/prd/`** coverage; align **`INDEX.md`**; optional **External** column when ticketing is not **`local`**

## Import (optional)

* `/skillgrid-import`
* Import legacy PRDs from **`prd/`**, **`docs/PRD/`**, **`docs/prd/`**, obvious **`PRD*.md`** paths, and existing **OpenSpec** changes into canonical **`.skillgrid/prd/PRD<NN>_<slug>.md`** files and **`.skillgrid/prd/INDEX.md`**.
* Use this directly when you already know the PRD path or change id, or let **`/skillgrid-explore`** invoke the same import/backfill rules during inventory.
* Import is conservative: it does not delete source PRDs, silently merge ambiguous matches, create remote issues, or implement product code.

## Design (optional)

* `/skillgrid-design`
* On-demand design workshop: extract, browse, tune, SuperDesign, Impeccable — converges on root **`DESIGN.md`**

## Brainstorm

* `/skillgrid-brainstorm`
* Explore project context first; clarify goals one decision at a time; propose **2-3 approaches** with tradeoffs before converging; optional previews in **`.skillgrid/preview/`**
* Present the selected direction for approval before it becomes PRD/spec/task input. Small changes can have short designs, but assumptions must still be explicit.
* Record architecture decisions in **`.skillgrid/project/ARCHITECTURE.md`**
* No remote issues unless **`.skillgrid/config.json`** says so (see init).

## Plan

* `/skillgrid-plan`
* PRD first: **`.skillgrid/prd/PRD<NN>_<slug>.md`**, **`.skillgrid/prd/INDEX.md`**
* **OpenSpec** change scaffold and artifact loop; hybrid **Engram** handoff
* `tasks.md` should be executable by a fresh agent: concrete scope, no placeholders, vertical slices, and test-first steps for behavioral work.
* Optional remote issues only per **`ticketing.provider`** in **`.skillgrid/config.json`**

## Breakdown

* `/skillgrid-breakdown`
* Sync PRD **Implementation tasks** with **`openspec/changes/<id>/tasks.md`**; tag tasks with **`[HITL]`** or **`[AFK]`** (see *HITL vs AFK slices* below); order HITL decisions before dependent AFK work; **prefer AFK** when scoping slices
* Optional vertical-slice issues per ticketing provider; **`local`** keeps work in PRD + `tasks.md` only

## Apply

* `/skillgrid-apply`
* Before every apply run that proceeds to implementation, automatically create a named **`/skillgrid-checkpoint create before-apply-<change-id>`** entry in **`.skillgrid/tasks/checkpoints.log`**.
* Critically review the selected task before editing; stop if instructions, verification, or scope are unclear.
* Implement from `tasks.md` / **OpenSpec** apply instructions; use Red-Green-Refactor for behavioral code unless an exception is explicit; **per-change handoff** (`.skillgrid/tasks/context_<change-id>.md`) in a **single working tree**—optional feature branch in that clone is fine; do **not** assume git worktrees
* After each slice or delegated task, reassess: update the handoff with completed work, evidence, blockers, changed assumptions, and the next recommended slice.
* Delegated implementation uses a double-review loop: spec compliance first, then code quality/security/maintainability. Accepted review fixes go back through the implementer and the same review stage repeats before the task is marked complete.

## Test

* `/skillgrid-test`
* Quality gates, functional tests, automated security baselines; tie to PRD success criteria. **Issue id** in arguments follows the configured ticketing provider.
* Confirm behavior is tested through public interfaces. Bug fixes need a reproduction test that failed before the fix unless a documented exception exists.

## Security (optional)

* `/skillgrid-security`
* Deeper review than the Test-phase scanners alone

## Validate

* `/skillgrid-validate`
* Single gate: review + security + user sign-off; **rollback** path if the user rejects
* Request independent review with fresh artifact context, then receive feedback skeptically: clarify unclear items, verify against the codebase and specs, push back on incorrect or out-of-scope suggestions, and fix accepted findings one at a time. Critical and important findings block completion until fixed, accepted as explicit risk, or converted into follow-up work.

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
│   │   ├── events/
│   │   │   └── <change-id>.jsonl
│   │   └── research/
│   │       └── <change-id>/
│   │           └── <topic-or-agent>_<optional-date>.md
│   ├── preview/
│   └── scripts/
│       ├── skillgrid-ui.mjs
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
| **`.skillgrid/tasks/events/<change-id>.jsonl`** | Append-only workflow timeline rendered by the local dashboard |
| **`.skillgrid/tasks/research/<change-id>/`** | Long research, scrapes, or subagent reports (one file per topic or run) |

**Subagent contract**

1. **Before work:** Read `context_<change-id>.md` (and any cited `research/…` files).  
2. **During work:** Write lengthy output under `research/<change-id>/` with a **distinct** filename.  
3. **After work:** Update the handoff (state, index row, next actions).  
4. **Return to parent:** A short message with file paths to read before product code changes.

**Cross-link** — In `openspec/changes/<change-id>/proposal.md`, include a line the reader cannot miss, e.g. `**Skillgrid session context:** .skillgrid/tasks/context_<change-id>.md`.

**Event log** — Append one JSON object per line under **`.skillgrid/tasks/events/<change-id>.jsonl`** for phase starts, progress, blockers, review outcomes, and completions. Required dashboard fields are **`time`**, **`changeId`**, **`node`** or **`phase`**, **`status`**, and **`summary`**. Optional fields include **`prd`**, **`blocker`**, **`artifacts`**, and **`external`**. The event log is a timeline; **`context_<change-id>.md`** remains the current-state summary.

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
- Optional title block fields **Preview** and **External** help the dashboard link generated preview HTMLs and tracker issues.

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

OpenSpec is the secondary technical spec backend, not the whole Skillgrid workflow. Skillgrid owns PRDs, the ordered work index, handoff state, checkpoints, previews, ticketing links, and phase/status flow around OpenSpec.

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

**Persistence** — **`artifactStore.mode`**: **`hybrid`** (disk + Engram), **`openspec`**, or **`engram`**. The init command records this; **`/skillgrid-session`** loads it for cold starts. Use **`hybrid`** unless a project explicitly cannot keep on-disk artifacts or cannot run Engram: disk gives reviewable, git-visible source files, while Engram gives compact recall across sessions and compactions.

| Mode | Use when | Artifact contract |
|------|----------|-------------------|
| **`hybrid`** | Recommended default for most teams and repos. | Create PRDs, OpenSpec changes, handoff, research, events, checkpoints, and previews on disk; save concise Engram summaries with stable topic keys that point back to those paths. |
| **`openspec`** | Engram is unavailable or the project wants repo-only state. | Same on-disk artifacts as `hybrid`; Engram saves are optional durable notes only. |
| **`engram`** | Disk specs are not allowed or a lightweight memory-only setup is required. | Do not assume `openspec/` exists. Store the equivalents below as Engram observations and include enough path or issue references to recover external context. |

Concrete Engram-mode equivalents:

| Disk artifact in hybrid/openspec | Engram-mode equivalent |
|----------------------------------|------------------------|
| `.skillgrid/prd/PRD<NN>_<slug>.md` | `mem_save` topic `skillgrid/<change-id>/prd` with problem, goals, scope, requirements, success criteria, status, and links to any external issue. |
| `.skillgrid/prd/INDEX.md` | `mem_save` topic `skillgrid/index` with ordered PRDs, statuses, active change ids, and external tracker links. |
| `openspec/changes/<change-id>/proposal.md` | `mem_save` topic `skillgrid/<change-id>/proposal` with intent, non-goals, design summary, and context references. |
| `openspec/changes/<change-id>/specs/*/spec.md` | `mem_save` topic `skillgrid/<change-id>/spec` with requirements and scenarios written in verifiable language. |
| `openspec/changes/<change-id>/tasks.md` | `mem_save` topic `skillgrid/<change-id>/tasks` with ordered `[HITL]` / `[AFK]` tasks and verification steps. |
| `.skillgrid/tasks/context_<change-id>.md` | `mem_save` topic `skillgrid/<change-id>/context` with current state, blockers, evidence, and next action. |
| `.skillgrid/tasks/research/<change-id>/...` | `mem_save` topic `skillgrid/<change-id>/research/<topic>` with concise findings and source links; avoid dumping large scrape output. |
| `.skillgrid/tasks/events/<change-id>.jsonl` | `mem_save` topic `skillgrid/<change-id>/events` with milestone events, status changes, blockers, and evidence pointers. |
| `.skillgrid/tasks/checkpoints.log` | `mem_save` topic `skillgrid/<change-id>/checkpoint` with checkpoint name, branch/SHA when available, dirty status, and verification evidence. |
| `.skillgrid/preview/*.html` and root `DESIGN.md` updates | `mem_save` topic `skillgrid/<change-id>/design` with selected direction, preview/export links, tokens, accessibility notes, and unresolved risks. |

## Source Of Truth

Avoid duplicating requirements as independent truth across PRDs, OpenSpec, external issues, and memory:

| Artifact | Owns |
|----------|------|
| **PRD** | Product intent, user-facing problem, scope, success criteria, and slice boundary. |
| **OpenSpec delta specs** | Verifiable technical behavior and scenarios for the active change. |
| **`tasks.md`** | Ordered implementation work and verification checklist for the slice. |
| **Handoff** | Current state, blockers, evidence index, and next action. |
| **External tracker** | Coordination mirror, assignment, and discussion; not canonical intent unless imported back. |
| **Engram** | Durable cross-session memory summary; not a replacement for repo artifacts. |
