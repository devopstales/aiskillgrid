---
name: /skillgrid-init
id: skillgrid-init
category: Workflow
description: Bootstrap workflow: greenfield/brownfield, structure, graphify, OpenCode, baseline skills
allowed-tools: Read, Write, Glob, Grep, Bash, Task
argument-hint: "[optional: app purpose, stack, or brownfield notes]"
---

<objective>

You are executing **`/skillgrid-init`** (DEFINE phase) for the Skillgrid workflow.

</objective>

<process>

## Steps

1. **Greenfield vs brownfield** — Inspect the repo (existing app code, tests, production history). **Brownfield:** tell the user to run **`/skillgrid-explore`** before deep structural work so **`.skillgrid/project/`** (`ARCHITECTURE`, `STRUCTURE`, `PROJECT`) and **AGENTS.md** are grounded. **Greenfield:** ask until you have a clear picture of the **app’s main purpose**, **target stack**, and **tools/technologies**; then create or update **`.skillgrid/project/ARCHITECTURE.md`**, **`STRUCTURE.md`**, **`PROJECT.md`**, and the **Project** chapter of root **`AGENTS.md`** (see [docs/wokflow.md](../../docs/wokflow.md)).
2. **Create folder structure** — Establish or verify project layout and conventions (source roots, **`.skillgrid/project/`** + **`prd/`** + **`preview/`** + **`scripts/preview.sh`** per [docs/wokflow.md](../../docs/wokflow.md), config dirs, `openspec/` or SDD layout if used).
3. **Graphify** — When graphify is part of the workflow, run **`graphify update .`** from the project root. **Do not block the main turn on graph output:** spawn a **subagent** (e.g. `Task` with the shell agent) or run the command in the **background** and continue init; the user can read **`graphify-out/`** when it finishes.
4. **OpenCode** — If the project uses OpenCode, ensure `.opencode/` (or product-specific) config is present and consistent; if not, initialize with `opencode init`.
5. **OpenSpec** — Ensure the OpenSpec CLI is on PATH when the team uses OpenSpec; run `openspec-onboard`-style first cycle only if the user wants a full guided onboarding now.
6. **Skills (`.agents/skills/`)** — Treat the list below as required reading; they define bootstrap behavior (stack detect, persistence mode, quality baseline).

## Skills to read and follow

Load each file before substantive work (read fully or skim by length):

- `.agents/skills/sdd-init/SKILL.md` — bootstrap SDD context (stack, conventions, persistence).
- `.agents/skills/openspec-onboard/SKILL.md` — first full OpenSpec cycle (requires OpenSpec CLI).
- `.agents/skills/search-first/SKILL.md` — research tools and patterns before building.
- `.agents/skills/karpathy-guidelines/SKILL.md` — assumptions, simplicity, surgical edits (baseline).
- `.agents/skills/skill-creator/SKILL.md` — add or extend Agent Skills in this repo.
- `.agents/skills/memory-protocol/SKILL.md` — Engram MCP: when to save, search, and close sessions.
- `.agents/skills/context-engineering/SKILL.md` — rules, context packing, MCP usage.
- `.agents/skills/using-agent-skills/SKILL.md` — meta: how to use the agent-skills pack.

## Project structure (example)

Conventions for a **Skillgrid-initialized** repo (names may vary; `install.sh` syncs the IDE and hub bits into *your* tree). Application source lives under whatever fits the stack (e.g. `src/`, `app/`, `lib/`).

```text
project-root/
├── AGENTS.md                      # rules for agents; often merged from hub .configs/AGENTS.md
├── DESIGN.md                      # UX/UI decisions; optional location
├── README.md
├── .skillgrid/
│   ├── project/                   # exploration outputs (/skillgrid-explore); see docs/wokflow.md
│   │   ├── ARCHITECTURE.md
│   │   ├── STRUCTURE.md
│   │   └── PROJECT.md
│   ├── prd/
│   │   ├── INDEX.md
│   │   └── <change-or-feature>.md
│   ├── preview/                  # /skillgrid-brainstorm: ephemeral MD/HTML (Preview picks)
│   └── scripts/
│       └── preview.sh            # scaffolds stubs under preview/
├── docs/                          # optional: extra PRD or architecture copies
│   ├── PRD/
│   └── ARCHITECTURE.md
├── .agents/
│   └── skills/                    # SKILL.md per skill; shared references/ checklists
├── .cursor/                       # from hub install: commands, rules, agents, mcp.json
│   ├── commands/                  # /skillgrid-*, /opsx-*
│   ├── rules/
│   └── mcp.json
├── .kilo/                         # same pattern when target uses Kilo Code
├── .opencode/                     # OpenCode: commands, opencode.json(c), skills mirror
├── .github/                       # Copilot prompts, agents, workflows
│   └── prompts/                 # e.g. skillgrid-*, opsx-*
├── .configs/                      # optional: hub-style MCP fragments (some installs keep these)
├── openspec/                      # after openspec init: changes/, specs/ (OpenSpec workflow)
├── .graphify/                     # optional: graphify project config
├── graphify-out/                  # optional: generated — graph.html, graph.json, report, cache/
└── src/ or app/ or lib/           # your product code, tests, package manifests, etc.
```

- **Per-project:** adjust `src/` to your monorepo layout (`packages/`, `apps/`, …) as needed; keep `AGENTS.md`, `.skillgrid/`, and `.agents/skills/` easy to find.
- **Optional tools:** add `graphify-out/` only after graphify has run and you know its on-disk layout.

## Notes

- **Phase 0:** For a new agent session or after compaction, run **`/skillgrid-session`** before heavy work so context, MCPs, and checkpoints stay bounded.
- Inspect the repo with tools; do not assume stack or layout.
- If OpenSpec vs SDD persistence is unclear, ask once, then follow existing `openspec/` trees or established repo conventions.

</process>
