---
name: /skillgrid-init
id: skillgrid-init
category: Workflow
description: Bootstrap workflow: structure, graphify, hybrid persistence (OpenSpec + Engram)
allowed-tools: Read, Write, Glob, Grep, Bash, Task
argument-hint: "[optional: app purpose, stack, or brownfield notes]"
---

<objective>

You are executing **`/skillgrid-init`** (DEFINE phase) for the Skillgrid workflow.

You detect the stack and repo layout, then apply a **persistence mode** the team chooses (see **Persistence modes**). **Skillgrid’s default recommendation is `hybrid`**: on-disk **`openspec/`** plus **Engram** so work survives compactions. Reconcile the `openspec/` tree with whatever your installed **OpenSpec CLI** scaffolds (`openspec init`, `openspec new change`, project templates)—do not assume a frozen layout beyond what the CLI and repo already use.

**PRD files (non-negotiable):** **Never** create **`prd/`** at the **repository root**—no **`prd/INDEX.md`**, no **`prd/PRD*.md`**, no root **`prd/`** folder. The **only** Skillgrid paths are **`.skillgrid/prd/INDEX.md`** and **`.skillgrid/prd/PRD<NN>_<slug>.md`**. If root `prd/` exists from legacy work, **do not** add new files there; create under **`.skillgrid/prd/`** and migrate or leave root as read-only until cleaned up.

</objective>

<process>

## Flow

```mermaid
flowchart TD
    START([User calls /skillgrid-init])
    GFB{Greenfield or brownfield?}
    START --> GFB
    GFB -->|Greenfield| ASK[Ask purpose, stack, tools]
    ASK --> DESIGN[Create root DESIGN.md from template]
    DESIGN --> ARCH[Create .skillgrid/project/ARCHITECTURE.md]
    ARCH --> STRUCT[Create .skillgrid/project/STRUCTURE.md]
    STRUCT --> PROJ[Create .skillgrid/project/PROJECT.md]
    PROJ --> AGENTS[Update root AGENTS.md]
    AGENTS --> OPS{openspec/ exists?}
    OPS -->|No| INIT[openspec init --tools none]
    INIT --> CONFIG[Populate openspec/config.yaml]
    CONFIG --> PRD_CREATE[Create .skillgrid/prd/ and .skillgrid/tasks/]
    PRD_CREATE --> ENGRAM[mem_save to Engram]
    ENGRAM --> GRAPH[graphify update .]
    GRAPH --> DONE([Handoff to /skillgrid-plan or /skillgrid-explore])
    GFB -->|Brownfield| EXPLORE[Recommend /skillgrid-explore]
    EXPLORE --> DONE

    OPS -->|Yes| PRD_CREATE
```

## Steps

1. **Greenfield vs brownfield** — Classify from **evidence**, not from “repo already has files.”

   - **Greenfield (default for init) when** any of these hold: there is **no** substantive **application** source tree yet—typical roots include `src/`, `app/`, `lib/`, `backend/`, `packages/<name>/` **with product code** (not only config). **OR** the tree is **only** workflow / IDE / agent setup: e.g. **`.cursor/`**, **`.kilo/`**, **`.opencode/`**, **`.github/`**, **`.agents/`**, **`.vscode/`**, **`.skillgrid/`**, **`openspec/`**, root **`AGENTS.md`**, **`docs/`** — that pattern is a **config hub or pre-code** project → **greenfield**, not brownfield. Having **`openspec/`** or **`.git/`** does **not** mean brownfield.
   - **Brownfield** when: real product code exists in the repo to map and change (services, app entrypoints, libraries beyond stubs), or the user **explicitly** says the codebase is existing / legacy / brownfield.
   - **Wrong:** calling a repo brownfield just because **`.skillgrid/`**, **`openspec/`**, or IDE folders exist.
   - **Greenfield path:** clarify **purpose**, **stack**, and **tools**; then create or update **`.skillgrid/project/*.md`** and the **Project** section of root **`AGENTS.md`** (see **`.skillgrid/scripts/skillgrid-workflow.md`**).
   - * **Design intent** — Ask: *“Do you have design preferences—colors, typography, aesthetic direction—or a design system you’d like to use?”*  
     - * Create (or update) **`DESIGN.md`** at the root.  
     - * Always scaffold the file using the template below, even if the user provides only a vague direction. Fill provided tokens; leave unspecified ones empty (`""`).  
     - * The front matter must match the structure used by Skillgrid’s design tokens: `name`, `colors`, `typography`, `rounded`.  
     - * If the user wants to defer all design decisions, still create the file with empty front matter and a placeholder `## Overview` (e.g. “To be defined during brainstorming.”).
   - **Brownfield path:** recommend **`/skillgrid-explore`** before big structural work so **`.skillgrid/project/`** and **`AGENTS.md`** are grounded in the real code.

2. **On-disk OpenSpec** — For modes **`openspec`** and **`hybrid`** only: if there is no `openspec/` yet, run the **OpenSpec CLI** from the **project root** using the Skillgrid default:

   - **Default command:** `openspec init --tools none`  
   - **Do not** substitute other “init” flows for this step (e.g. `opencode init . --tools all --force` or any `init . --tools all --force` pattern). Those are different products; Skillgrid OpenSpec bootstrap is **`openspec`** with **`--tools none`** unless the user explicitly asks for more tools.  
   - If your installed `openspec --help` shows different flags, follow the CLI, but keep **`--tools none`** as the baseline to avoid unsolicited editor integrations.  
   - Typical output includes `openspec/config.yaml` (or equivalent), `openspec/specs/`, `openspec/changes/`, and `openspec/changes/archive/`; **match generated output**, do not hand-craft a conflicting tree. If `openspec/` already exists, summarize what is there and ask before overwriting. For mode **`engram`**, **do not** create `openspec/` unless the user explicitly asks to add it later.
   - **`openspec/config.yaml` after init:** As soon as `openspec init` succeeds, ensure **`openspec/config.yaml`** exists and is populated for the project. If the CLI already wrote a minimal file, **merge**: keep a valid **`schema`** (default `spec-driven` unless the project uses another schema); fill or expand **`context`** from stack detection and this session; add **`rules`** if missing. If the file is absent, create it from **OpenSpec project config (template)** below. **`context`** must stay under **50KB** (OpenSpec hard limit). If a non-empty config already exists with custom content, do not blindly overwrite—show a short diff summary and ask before replacing.

3. **Engram** — For modes **`engram`**, **`hybrid`**, and (optionally) **`openspec`** when the team wants cross-session recall: after you have detected context (stack, test tooling, layout), **`mem_save`** with a stable **`topic_key`** so reruns **update** instead of duplicating:

   - **Topic key:** `skillgrid-init/{project-name}` (derive `{project-name}` from the repo or package name; keep it short and unique).
   - **Type** (e.g. `architecture` or your convention) and **project** field as your Engram setup expects.
   - **Content:** markdown summary: stack, test runner, key paths, whether `openspec/` exists, and which **persistence mode** is active.

4. **Create folder structure** — Establish or verify **`.skillgrid/project/`**, **`.skillgrid/prd/`**, **`.skillgrid/tasks/`**, **`.skillgrid/tasks/research/`** . **PRD files** use **`.skillgrid/prd/PRD<NN>_<slug>.md`**: two-digit `<NN>` is **execution order**; add **`.skillgrid/prd/INDEX.md`** listing PRDs sorted by `<NN>`. **Repeat:** do **not** create **`prd/`** at repo root. Application source lives where the stack expects (`src/`, `app/`, `packages/`, etc.).

5. **Graphify** — When graphify is in use, run **`graphify update .`** from the project root. **Do not** block the whole turn on it: use a subagent/background run; the user can read **`graphify-out/`** when ready.

6. **OpenCode** — If the product uses OpenCode, ensure **`.opencode/`** (or product-specific) config exists; if not, `opencode init` per product docs. **Do not** treat this as the OpenSpec bootstrap in step 2 — that must stay **`openspec init --tools none`**.

7. **OpenSpec CLI check** — If the team uses OpenSpec, ensure `openspec` is on `PATH` (`openspec --version`). For a **full first-cycle walkthrough** (narrated propose → apply → archive), run **`/opsx-onboard`** or follow the same phases embedded there.

8. **Baseline practices** (no external skill list required) — Research tools and similar repos before big bespoke code; keep edits small and verifiable; save decisions to Engram that should survive compactions; align **`AGENTS.md`** and project rules with how your team works.

## Persistence modes (authoritative)

Ask once if unclear; otherwise infer from the repo (**existing `openspec/`** vs **Engram-only** conventions in `AGENTS.md`).

| Mode | `openspec/` on disk | Engram `mem_save` |
|------|---------------------|-------------------|
| **`hybrid`** (recommended) | Yes — bootstrap if missing | Yes — `topic_key` **`skillgrid-init/{project-name}`** |
| **`openspec`** | Yes — bootstrap if missing | Optional — only if the team wants memory in addition to disk |
| **`engram`** | **No** — do not create `openspec/` by default | Yes — same `topic_key`; planning artifacts live in memory / external docs per team |
| **`none`** | No | No — return detected context only; warn that multi-session spec work will not persist |

- **Skill registry (optional):** If your org uses a generated skill registry (e.g. `.atl/skill-registry.md`), follow the same policies as the rest of the repo; do not block init on it.

## OpenSpec project config (template)

Default body for **`openspec/config.yaml`** when creating or completing the file after **`openspec init`**. Replace `<placeholders>` with values from repo inspection and the init conversation. Fields: **`schema`** (required), **`context`** (optional, injected into artifact instructions), **`rules`** (optional, keyed by artifact id).

```yaml
# openspec/config.yaml — Skillgrid default; merge with CLI output if present
schema: spec-driven

context: |
  Project: <short name — purpose in one line>
  Stack: <languages, frameworks, package manager>
  Layout: <e.g. src/, app/, packages/>
  Testing: <runner; where tests live>
  Skillgrid: PRDs in .skillgrid/prd/; persistence per AGENTS.md (hybrid: openspec/ + Engram when used)
  Conventions: <AGENTS.md / team rules — keep brief>

rules:
  proposal:
    - Tie to the PRD under .skillgrid/prd/ when one exists; name scope and non-goals.
  specs:
    - Use clear requirements; scenarios (Given/When/Then) when they help acceptance.
  design:
    - Boundaries, data flow, migration or rollback when relevant.
  tasks:
    - Small, verifiable checkboxes; trace to specs and PRD execution order.
```

## Project structure (example)

Conventions for a **Skillgrid-initialized** repo. Application code lives under whatever fits the stack.

```text
project-root/
├── AGENTS.md
├── README.md
├── DESIGN.md
├── .skillgrid/
│   ├── tasks/
│   │   ├── context_<change-id>.md
│   │   └── research/
│   │       └── <change-id>/
│   │           └── <topic-or-agent>_<optional-date>.md
│   ├── project/
│   │   ├── ARCHITECTURE.md
│   │   ├── STRUCTURE.md
│   │   └── PROJECT.md
│   ├── prd/
│   │   ├── INDEX.md
│   │   ├── PRD01_<first-slug>.md
│   │   └── PRD02_<next-slug>.md
│   ├── preview/
│   └── scripts/
│       └── preview.sh
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

Adjust for monorepos (`packages/`, `apps/`) and tool-generated folders.

## PRD / change `Status` lifecycle (authoritative)

Use these **exact** values on the **PRD** `**Status:**` line (and in **`.skillgrid/prd/INDEX.md`** / ticket tables) so phases stay consistent across **`/skillgrid-plan`**, **`/skillgrid-breakdown`**, **`/skillgrid-apply`**, **`/skillgrid-review`**, and **`/skillgrid-finish`**:

| When this command completes successfully | Set `Status` to |
|------------------------------------------|-----------------|
| **`/skillgrid-init`** | `draft` |
| **`/skillgrid-plan`** | `draft` |
| **`/skillgrid-breakdown`** | `todo` |
| **`/skillgrid-apply`** | `inprogress` |
| **`/skillgrid-review`** | `devdone` |
| **`/skillgrid-finish`** | `done` |

- **Initial Status:** At creation the PRD’s **`Status:`** is `draft`
- **Single line of truth:** the PRD’s **`Status:`** (and any mirrored table column). The agent **updates** that field at the end of the matching phase; later phases may overwrite the previous value.
- **Optional metadata:** OpenSpec on-disk or Engram can carry parallel notes, but the PRD `Status` is what humans scan first.

## Project document templates

Use these when creating **`.skillgrid/project/ARCHITECTURE.md`**, **`.skillgrid/project/STRUCTURE.md`**, and **`.skillgrid/project/PROJECT.md`** (greenfield or after **`/skillgrid-explore`**). Prefer mermaid in architecture where it helps.

### `.skillgrid/project/ARCHITECTURE.md`

````markdown
# <Project Name> — System Architecture

> **Version:** <semver or date>
> **Last Updated:** <YYYY-MM-DD>
> **Status:** <Draft | ToDo | InProgress | DevDone | Done | Superseaded>

---

## Table of Contents

1. [Architecture overview](#1-architecture-overview)
2. [Application layers](#2-application-layers)
3. [<Domain-specific sections>](#) <!-- e.g. data, auth, plugins, API, frontend, security, multi-tenancy, errors -->

---

## 1. Architecture overview

<One paragraph: layered or modular mental model.>

### High-level system diagram

```mermaid
graph TB
  %% User / edge / app / data / external systems
```

### Design decisions

| Decision | Choice | Rationale |
| -------- | ------ | --------- |
| …        | …      | …         |

## 2. Application layers

<ASCII or mermaid: presentation → service → data + cross-cutting concerns.>

### Layer rules

| Rule | Description |
| ---- | ----------- |
| …    | …           |

## 3. <Next sections>

<!-- Add per-domain sections as needed, e.g. database, auth, external APIs, plugins, caching, UI, security. -->
````

### `.skillgrid/project/STRUCTURE.md`

````markdown
# <Project Name> — Structure

> **Scope:** Repository layout. Optional: runtime / deployment topology (not CI/CD unless needed).

## Repository tree (high level)

- `<path>/`: <what lives here; main packages, apps, services>
- …

### Layout diagram (optional)

```text
<tree or mermaid of packages>
```

## Runtime / deployment (optional)

<!-- Topology: mermaid, workloads, services, data stores, ingress, config/secrets, environments. -->

### High-level topology

```mermaid
graph TB
  %% users, edge, cluster/services, data stores
```

### Environments

- **Local / dev:** …
- **Staging:** …
- **Production:** …

### Where to read next

- Link to `ARCHITECTURE.md`, `PROJECT.md`, and operator or developer docs.
````

### `.skillgrid/project/PROJECT.md`

````markdown
### Overview

**<Name>** is <one tight paragraph: stack, main capabilities, and notable companion tools or modules>.

### Audience

- **New developers:** …
- **Ops / SRE:** …
- **AI assistants / agents:** …

### Repository structure (high level)

- `<path>/`: <role>
- …

### Conceptual architecture

```mermaid
graph LR
  %% major boxes and trust boundaries
```

- **<Component A>:** <2–4 bullets>
- **<Component B>:** …

### How <main parts> fit together

- **Typical flow:** …

### Development setup (local)

- **Prerequisites:** …
- **Run:** …
- **Configuration:** …

### Testing

- **Unit / integration / e2e:** where tests live and how to run them.

### <Optional: product themes or roadmap at a glance>

```mermaid
graph TD
  CORE[Core] --> F1[Feature or theme 1]
  CORE --> F2[Feature or theme 2]
```

### Where to read next

- Architecture: `.skillgrid/project/ARCHITECTURE.md`
- Structure / infra: `.skillgrid/project/STRUCTURE.md`
- PRDs: `.skillgrid/prd/INDEX.md`
- …
````

### `DESIGN.md`

```markdown
---
name: <Project Name>
colors:
  primary: ""
  secondary: ""
  surface: ""
  on-surface: ""
  error: ""
typography:
  body-md:
    fontFamily: ""
    fontSize: ""
    fontWeight: ""
rounded:
  md: ""
# taste-skill dials (1–10)
design_variance: 5
motion_intensity: 3
visual_density: 5
# provenance
design_sources: []
---

# Design System

## Overview
<User-provided one-liner; else “To be defined during brainstorming.”>

## Design sources
<!-- List of sources: getdesign.md slug, skillui URL, Figma link, etc. -->
- None yet — defined in this file.

## Colors
- **Primary** (``): CTAs, active states, key interactive elements
- **Secondary** (``): Supporting UI, chips, secondary actions
- **Surface** (``): Page backgrounds
- **On-surface** (``): Primary text on dark/light backgrounds
- **Error** (``): Validation errors, destructive actions

## Typography
- **Headlines**: To be defined
- **Body**: To be defined
- **Labels**: To be defined

## Components
- **Buttons**: To be defined
- **Inputs**: To be defined
- **Cards**: To be defined

## Do's and Don'ts
- Do …
- Don't …
```

## Notes

- **Session hygiene:** For a new or compacted session, run **`/skillgrid-session`** before heavy work.
- If Engram is unavailable in a given session, still bootstrap **`openspec/`** and record in **`AGENTS.md`** that a **`mem_save`** to `skillgrid-init/{project-name}` is pending.

## Anti-patterns

- **Root `prd/` folder** – Never create `prd/INDEX.md` or new `prd/*.md` files at the repo root; all PRDs live under `.skillgrid/prd/`.
- **Config‑only = brownfield** – Don’t treat a repo that has only IDE/agent setup (`.cursor`, `openspec/`, `.skillgrid/`, etc.) as brownfield; it’s still greenfield unless there’s product code.
- **Skipping DESIGN.md** – Never defer the design document because “we’ll do it later”; at minimum create a stub with the token skeleton.
- **Ignoring persistence** – Don’t skip Engram if you’ve chosen hybrid mode; the `mem_save` is the backup for session compactions.
- **Blind overwrites** – Never overwrite an existing `openspec/config.yaml` without showing a diff and asking.

## Completion report (required)

End with a **Session wrap-up** the user can scan:

1. **What I did** — Bullets: files/dirs created or edited, CLI commands run, `openspec` / Engram / graphify actions, and classification (greenfield vs brownfield).
2. **Token / usage** — If the product shows **input/output tokens**, **context used**, or **session cost** for this turn, report it. If not available, state **`Token usage: not shown in this environment`** (do not guess).
3. **Suggested next command** — **Greenfield:** **`/skillgrid-plan`** (or **`/skillgrid-brainstorm`** if ideation first). **Brownfield:** **`/skillgrid-explore`** before planning. If context is cold: **`/skillgrid-session`**.

</process>
