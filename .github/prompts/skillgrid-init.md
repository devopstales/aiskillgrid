---
description: Bootstrap workflow: greenfield/brownfield, structure, graphify, CocoIndex, OpenCode, baseline skills
---

You are executing **`/skillgrid-init`** (DEFINE phase) for the Skillgrid workflow.

## Steps

1. **Greenfield vs brownfield** — Inspect the repo (existing app code, tests, production history). **Brownfield:** tell the user to run **`/skillgrid-explore`** before deep structural work so architecture and layout are grounded. **Greenfield:** proceed with bootstrap below.
2. **Graphify** — If this repo uses graphify, initialize it (e.g. `graphify .` from the project directory when applicable).
3. **CocoIndex Code** — From project root: `ccc init` if needed, then `ccc index` so semantic search is available.
4. **OpenCode** — If the project uses OpenCode, ensure `.opencode/` (or product-specific) config is present and consistent; if not, initialize with `opencode init`.
5. **OpenSpec** — Ensure the OpenSpec CLI is on PATH when the team uses OpenSpec; run `openspec-onboard`-style first cycle only if the user wants a full guided onboarding now.
6. **Skills (`.agents/skills/`)** — Treat the list below as required reading; they define bootstrap behavior (stack detect, persistence mode, quality baseline).
7. **Create folder structure** — Establish or verify project layout and conventions (source roots, `.skillgrid/project/` + `prd/` per [`docs/wokflow.md`](../../docs/wokflow.md), config dirs, `openspec/` or SDD layout if used).

## Skills to read and follow

Load each file before substantive work (read fully or skim by length):

- `.agents/skills/sdd-init/SKILL.md` — bootstrap SDD context (stack, conventions, persistence).
- `.agents/skills/openspec-onboard/SKILL.md` — first full OpenSpec cycle (requires OpenSpec CLI).
- `.agents/skills/search-first/SKILL.md` — research tools and patterns before building.
- `.agents/skills/karpathy-guidelines/SKILL.md` — assumptions, simplicity, surgical edits (baseline).
- `.agents/skills/skill-creator/SKILL.md` — add or extend Agent Skills in this repo.
- `.agents/skills/ccc/SKILL.md` — CocoIndex: `ccc init`, indexing, semantic search.
- `.agents/skills/context-engineering/SKILL.md` — rules, context packing, MCP usage.
- `.agents/skills/using-agent-skills/SKILL.md` — meta: how to use the agent-skills pack.

## Project structure (example)

Conventions for a **Skillgrid-initialized** repo (names may vary; `install.sh` syncs the IDE and hub bits into *your* tree). Application source lives under whatever fits the stack (e.g. `src/`, `app/`, `lib/`).

```text
project-root/
├── AGENTS.md                      # rules for agents; often merged from hub .configs/AGENTS.md
├── DESIGN.md                      # UX/UI decisions (/skillgrid-design); optional location
├── README.md
├── .skillgrid/
│   ├── project/                   # exploration outputs (/skillgrid-explore); see docs/wokflow.md
│   │   ├── ARCHITECTURE.md
│   │   ├── STRUCTURE.md
│   │   └── PROJECT.md
│   └── prd/
│       ├── INDEX.md
│       └── <change-or-feature>.md
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
├── .cocoindex_code/              # optional: ccc index — path may differ; check `ccc init` output
└── src/ or app/ or lib/           # your product code, tests, package manifests, etc.
```

- **Per-project:** adjust `src/` to your monorepo layout (`packages/`, `apps/`, …) as needed; keep `AGENTS.md`, `.skillgrid/`, and `.agents/skills/` easy to find.
- **Optional tools:** add `graphify-out/`, CocoIndex dirs, and similar only after the tool has run and you know its on-disk layout.

## Notes

- **Phase 0:** For a new agent session or after compaction, run **`/skillgrid-session`** before heavy work so context, MCPs, and checkpoints stay bounded.
- Inspect the repo with tools; do not assume stack or layout.
- If OpenSpec vs SDD persistence is unclear, ask once, then follow existing `openspec/` trees or established repo conventions.
