# Skillgrid templates and planning logic

This document is the **human-oriented** reference for how **PRDs** (product requirements documents), the **INDEX** (the `.skillgrid/prd/INDEX.md` file: dependency-ordered PRD list plus optional **execution snapshot** — current phase, active change or slice, discovered work, and session notes), **OpenSpec** (the repo’s spec-driven change layout under `openspec/`: one folder per initiative with proposal, design, tasks, and slice specs), **vertical slices** (thin, shippable units of work), **ADRs** (architecture decision records), usually stored as **MADR** (Markdown Any Decision Record) files under `.skillgrid/adr/`, and project docs fit together. **Blank file shapes** live under **`.skillgrid/templates/`** as **`template-<kebab-case>.md`** files (plus `README.md` for the index; synced by the hub `install.sh`); skills point there so templates stay editable without hunting through long skill markdown.

## How work flows through the artifacts

The diagram below is the Skillgrid mental model: product intent lives in **PRD** and **INDEX** files, execution intent is refined in an **OpenSpec** change, engineering work runs against **slice specs** and **tasks**, and feedback updates the catalog and snapshot. (Other tools—for example a local issue DAG or “beads”-style trackers—can mirror the same graph; they are optional and not required by this hub.)

```mermaid
flowchart TB
  subgraph Intent["Product intent on disk"]
    Idea([Initiative / idea])
    PRDNode["PRD file<br/>product requirements document<br/>.skillgrid/prd/ (see naming below)"]
    IndexNode["INDEX<br/>PRD catalog + execution snapshot<br/>.skillgrid/prd/INDEX.md"]
  end

  subgraph Change["OpenSpec change folder openspec/changes/change-id/"]
    Prop[proposal.md]
    Des[design.md]
    Task[tasks.md]
    Slice["Slice spec<br/>openspec ... specs slice spec.md"]
  end

  subgraph Loop["Engineering and proof"]
    Ctx["Handoff context stub<br/>.skillgrid/tasks per change"]
    Impl[Implement / apply]
    TV["Test + validate"]
  end

  subgraph Decisions["Durable architecture"]
    Madr["MADR file<br/>Markdown Any Decision Record<br/>.skillgrid/adr/NNNN_slug.md"]
    Arch["ARCHITECTURE.md<br/>links + summaries"]
  end

  Idea --> PRDNode
  PRDNode --> IndexNode
  PRDNode --> Prop
  Prop --> Des
  Des --> Task
  Task --> Slice
  Slice --> Ctx
  Ctx --> Impl
  Impl --> TV
  TV -->|"next slice, status, evidence"| IndexNode
  Des -.->|"significant tradeoffs"| Madr
  Madr -.-> Arch
  IndexNode -.-> Arch
```

## Where templates live

| Path | Purpose |
|------|---------|
| [`.skillgrid/templates/README.md`](../.skillgrid/templates/README.md) | Naming convention + index of all `template-*.md` files |
| [`.skillgrid/templates/template-adr.md`](../.skillgrid/templates/template-adr.md) | **MADR** (Markdown Any Decision Record) template for new **ADRs** (architecture decision records) in `.skillgrid/adr/` |
| [`.skillgrid/templates/template-prd.md`](../.skillgrid/templates/template-prd.md) | New **PRD** (product requirements document) body |
| [`.skillgrid/templates/template-index.md`](../.skillgrid/templates/template-index.md) | `.skillgrid/prd/INDEX.md` |
| [`.skillgrid/templates/template-openspec-tasks.md`](../.skillgrid/templates/template-openspec-tasks.md) | `openspec/changes/<change-id>/tasks.md` |
| [`.skillgrid/templates/template-openspec-slice-spec.md`](../.skillgrid/templates/template-openspec-slice-spec.md) | `openspec/changes/<change-id>/specs/<slice>/spec.md` |
| [`.skillgrid/templates/template-handoff-context.md`](../.skillgrid/templates/template-handoff-context.md) | `.skillgrid/tasks/context_<id>.md` |
| [`.skillgrid/templates/template-project.md`](../.skillgrid/templates/template-project.md) | `.skillgrid/project/PROJECT.md` |
| [`.skillgrid/templates/template-architecture.md`](../.skillgrid/templates/template-architecture.md) | `.skillgrid/project/ARCHITECTURE.md` |
| [`.skillgrid/templates/template-structure.md`](../.skillgrid/templates/template-structure.md) | `.skillgrid/project/STRUCTURE.md` |
| [`.skillgrid/templates/template-design.md`](../.skillgrid/templates/template-design.md) | Repo root `DESIGN.md` |

Skills (`skillgrid-prd-artifacts`, `skillgrid-project-docs`, `skillgrid-filesystem-handoff`, `documentation-and-adrs`) **describe behavior** and link these files; they should not drift into a second full copy of the template without updating `.skillgrid/templates/`.

## Planning logic: hierarchy

| Level | Jira-style | GitHub-style | Artifact |
|-------|------------|--------------|----------|
| Milestone / program | Epic | Milestone | `.skillgrid/prd/INDEX.md` — dependency-ordered **PRDs** (product requirements documents) + **execution snapshot** (phase, active change/slice, discovered work, session notes) |
| Feature | Task | Issue | `.skillgrid/prd/PRD<NN>_<slug>.md` + one `openspec/changes/<change-id>/` |
| Shippable unit | Sub-task | Checklist item | Vertical slice — `tasks.md` + `specs/<slice>/spec.md` |

There is **no** `.skillgrid/project/TASK.md`. “Where we are” lives in the **INDEX** snapshot and **OpenSpec** task artifacts.

## OpenSpec layout (per change)

```text
openspec/changes/<change-id>/
  proposal.md
  design.md
  tasks.md
  specs/<vertical-slice-slug>/
    spec.md
openspec/specs/<change-id>/spec.md   # optional umbrella
```

- **`tasks.md`** — cross-slice ordering and integration checklist.
- **`specs/<slice>/spec.md`** — bounded requirements and checklist for one vertical slice (preferred context for a single apply session).
- **`openspec/specs/<change-id>/spec.md`** — optional cross-cutting spec for the initiative.

## ADRs (MADR)

Repo-wide architectural decisions use **MADR** (Markdown Any Decision Record) files under **`.skillgrid/adr/`** (named `NNNN-kebab-title.md`). That directory holds **only** those **ADR** (architecture decision record) files — no README or other metadata there. Copy **`.skillgrid/templates/template-adr.md`**. Summaries and links belong in **`.skillgrid/project/ARCHITECTURE.md`** under **Durable decisions** so others can discover decisions. See **`documentation-and-adrs`** skill.

## Session bootstrap

For the read order agents should use at session start, see **`.configs/AGENTS.md`** (Skillgrid session bootstrap) and **`docs/02-workflow-usage.md`**.

## Related docs

- [Workflow usage](02-workflow-usage.md) — commands, slices, handoff, smart zone
- [Skills](05-skills.md) — how skills load and when to use registry
