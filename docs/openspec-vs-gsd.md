# OpenSpec vs GSD-2: file layout and workflow

This document compares **OpenSpec** (change-centric specs in `openspec/`, driven by the `openspec` CLI — see [Fission-AI/OpenSpec](https://github.com/Fission-AI/OpenSpec)) and **GSD-2** (“Get Shit Done” v2: `gsd-pi` — see [gsd-build/GSD-2](https://github.com/gsd-build/GSD-2)). Both aim to keep AI-assisted development aligned with written intent; they differ in *what lives on disk*, *who orchestrates* (CLI vs a long-running agent app), and *how work is decomposed*.

## Antigravity Kit vs GSD-2 (name check)

Some lists used to label the wrong repo as “GSD.” These are **different systems**:

| | **Antigravity Kit** | **GSD-2 (Get Shit Done v2)** |
|--|--------------------|------------------------------|
| **Repo / package** | [vudovn/antigravity-kit](https://github.com/vudovn/antigravity-kit), `npx @vudovn/ag-kit init` | [gsd-build/GSD-2](https://github.com/gsd-build/GSD-2), `npm i -g gsd-pi` |
| **What you get** | A **template pack** into **`.agent/`**: agents, skills, and workflow-style slash commands (`/plan`, `/brainstorm`, …) | A **long-running GSD app** with **`.gsd/`**, SQLite `gsd.db`, milestones → slices → tasks, auto/step mode, reports |
| **Orchestrator** | Your editor/assistant + workflow markdown; no GSD state machine | `gsd` / `/gsd auto` state machine, crash recovery, cost tracking, etc. |
| **Spec story** | Workflows + skills; not the same as `openspec/` delta specs | Roadmaps, plans, and DB-backed state—not OpenSpec’s `changes/` + `specs/` merge model |

*Antigravity Kit* is an **ag-kit** installer for `.agent/`. *GSD-2* is the **successor to the old “GSD v1” prompt pack** (now a full `gsd-pi` app). This doc’s main comparison is **OpenSpec vs GSD-2**; Antigravity is separate unless you **choose** to combine it with OpenSpec in one repo.

**Antigravity on disk (typical):**

```text
project/
└── .agent/                 # from `ag-kit init` — agents, skills, workflows (see upstream README for layout)
```

Keep **`.agent/`** out of `.gitignore` if your IDE needs to index slash commands, or use `.git/info/exclude` for local-only ignores (per upstream docs).

---

## At a glance

| Dimension | OpenSpec | GSD-2 |
|-----------|----------|--------|
| **Primary artifact** | One **change** folder under `openspec/changes/<name>/` with proposal, design, **delta** specs, and tasks | A **.gsd/** workspace: milestones → slices → tasks, plus SQLite `gsd.db` and many derived markdown files |
| **“Source of truth” for behavior** | `openspec/specs/<domain>/spec.md` (cumulative). Changes carry **deltas** until archived | `PROJECT.md`, `DECISIONS.md`, per-milestone `M*-ROADMAP.md`, slice/task plans; knowledge graph and **memories** in DB |
| **Orchestration** | `openspec` CLI + your editor/assistant (`openspec new`, `openspec status`, `openspec instructions`, `openspec archive`, …) | `gsd` / `gsd-pi` **application**: step mode (`/gsd`) and **auto** (`/gsd auto`) state machine, fresh context per unit |
| **Work breakdown** | User/schema-defined artifacts (often: proposal → specs → design → tasks) | Fixed hierarchy: **Milestone** (4–10 **slices**) → slice (1–7 **tasks**); task = one context-window unit |
| **Typical flow** | **Propose → apply → verify → archive** (merge deltas into `specs/`) | **Discuss → plan → execute → complete → reassess** per slice; milestone validation gate |
| **Best when** | You want **git-friendly**, reviewable spec deltas and a **portable** `openspec/` tree across many IDEs; minimal runtime | You want a **batteries-included** agent runner (cost tracking, worktrees, crash recovery, HTML reports, verification commands) |
| **Commit / install** | Add `openspec/` + CLI; no long-running process | `npm i -g gsd-pi`; project gains `.gsd/` (and team policies for what to **gitignore**) |

---

## File structure

### OpenSpec (after `openspec init`)

Conceptually:

```text
openspec/
├── config.yaml                 # optional project config, schemas
├── specs/                      # current system behavior (by domain)
│   └── <domain>/
│       └── spec.md
└── changes/                    # one folder per in-flight change
    ├── <change-name>/
    │   ├── .openspec.yaml     # change metadata (schema, etc.)
    │   ├── proposal.md        # why / what (intent, scope)
    │   ├── design.md          # technical approach
    │   ├── tasks.md           # implementation checklist
    │   └── specs/             # delta specs (ADDED / MODIFIED / REMOVED)
    │       └── <domain>/
    │           └── spec.md
    └── archive/               # completed changes
        └── YYYY-MM-DD-<name>/
```

- **Deltas** live under the change until **archive** merges them into `openspec/specs/`.
- Multiple **changes** can exist in parallel; merge order is a team process.
- [AISkillGrid](../README.md) maps **slash commands** like `/opsx-*` and skills (`openspec-propose`, `openspec-apply-change`, …) to this layout.

### GSD-2 (project + engine)

GSD-2 is **not** “only markdown”: it uses a **SQLite** database (`gsd.db`) and generated files. Representative layout (from upstream docs; exact names follow milestone/slice IDs):

```text
project/
├── .gsd/
│   ├── gsd.db*                 # authoritative state, memories, tool ledger (do not hand-edit)
│   ├── STATE.md                # quick-glance “what’s next”
│   ├── PREFERENCES.md         # project prefs (optional; YAML frontmatter)
│   ├── reports/                # generated HTML reports per milestone
│   ├── milestones/             # milestone-scoped continue markers, etc.
│   ├── runtime/, activity/, worktrees/, parallel/, …
│   └── (lock files, metrics, event logs — operational)
├── PROJECT.md                 # living project description (managed)
├── DECISIONS.md              # decisions register
├── KNOWLEDGE.md              # patterns / lessons
├── M001-ROADMAP.md           # milestone plan, slices, dependencies
├── M001-CONTEXT.md           # discuss-phase capture
├── M001-RESEARCH.md         # research
├── S01-PLAN.md, T01-PLAN.md, T01-SUMMARY.md, …
└── (your app source)
```

- **Milestone / slice / task** files encode the same *kind* of intent as a big feature, but the **GSD app** dispatches work, not a thin CLI alone.
- Teams often **gitignore** large parts of `.gsd/` except selected files; see GSD-2’s “working in teams” guidance.

---

## Workflow comparison

### OpenSpec

1. **Init** — `openspec init` scaffolds `openspec/` and (often) editor integrations.
2. **New change** — `openspec new change "<name>"` (or skills/commands that wrap it) creates `openspec/changes/<name>/` with metadata.
3. **Fill artifacts** — Follow schema order (commonly: proposal, delta specs, design, tasks). The CLI can emit **instructions** and paths (`openspec instructions`, `openspec status`).
4. **Implement** — Code against `tasks.md` and the delta spec; your assistant uses the same files as humans.
5. **Verify** — Optional `openspec verify` / project tests; this hub’s [`openspec-verify-change`](../.agents/skills/openspec-verify-change/SKILL.md) skill aligns with that gate.
6. **Archive** — `openspec archive` merges approved deltas into `openspec/specs/` and moves the change to `changes/archive/…`.

*Philosophy:* **Agree in files**, then implement; the repo holds **auditable** spec history. Orchestration is **yours** (agent + tools), with the CLI as the source of structure.

### GSD-2

1. **Start** — `gsd` in repo; if no `.gsd/`, a **discuss** flow captures goals and constraints.
2. **Milestones & roadmap** — Roadmap markdown (`M…-ROADMAP.md`) defines slices; each slice is planned then executed.
3. **Per slice** — Integrated **research + plan** → **execute tasks** in **fresh** contexts → **complete** with summaries and optional UAT.
4. **Auto or step** — `/gsd auto` runs a **state machine** (read `STATE.md`, dispatch next unit, verify, commit) until done; `/gsd` does one unit at a time.
5. **Milestone end** — Validation gate vs roadmap; **HTML reports** in `.gsd/reports/`.
6. **Ongoing** — `LEARNINGS.md` extraction, knowledge graph, memories — cross-session **project memory** is a first-class feature.

*Philosophy:* **The engine runs the loop**; files are both human-readable and machine-state. Heavier runtime, **stronger** automation and observability.

---

## Mapping concepts (rough)

| OpenSpec | GSD-2 (approximate) |
|----------|----------------------|
| `proposal.md` | Discuss + roadmap intent / milestone charter |
| `specs/…` deltas | Requirements embedded in plans + roadmap + “must-haves” in task plans |
| `design.md` | Slice plan / architecture decisions (`DECISIONS.md`, `S*-PLAN.md`) |
| `tasks.md` | `T*-PLAN.md` and execution summaries |
| `openspec archive` + merge | Milestone complete + report + optional squash merge |
| (no first-party DB) | `gsd.db` + `memories` + reports |

---

## Choosing one (or both)

- Prefer **OpenSpec** if you want a **lightweight, spec-first** contract in-repo, **IDE-agnostic** tooling, and explicit **delta** specs that **merge** into a long-lived `specs/` tree.
- Prefer **GSD-2** if you want a **single CLI/TUI** that **manages** sessions, worktrees, budgets, verification, and **milestones** with minimal custom glue.
- **Together:** not mutually exclusive in theory (e.g. store OpenSpec under `openspec/` *and* use GSD to execute), but you should avoid **duplicating** the same requirements in two systems—pick a **source of truth** for product behavior.

---

## References

- OpenSpec: [Fission-AI/OpenSpec](https://github.com/Fission-AI/OpenSpec) (getting-started docs describe `openspec/` layout and archive behavior).
- GSD-2: [gsd-build/GSD-2](https://github.com/gsd-build/GSD-2) (README, user docs under `docs/user-docs/`).
- Antigravity Kit: [vudovn/antigravity-kit](https://github.com/vudovn/antigravity-kit) (`@vudovn/ag-kit`, **`.agent/`** templates).
- This hub: [commands.md](commands.md), [wokflow.md](wokflow.md), [tools.md](tools.md).
