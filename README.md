# AISkillGrid

A **configuration hub** for opinionated AI-assisted development: reusable **skills**, **slash commands**, merged **MCP server** definitions, and an **install script** that copies normalized settings into your application repo. The workflow uses **OpenSpec**-style change management (specs under `openspec/`, `openspec` CLI) and spec-driven development skills together with production-oriented practices (testing, security, documentation).

---

## Highlights

| Feature | What it does | Why it matters |
|---|---|---|
| Skillgrid workflow | Guides work through init, explore, brainstorm, plan, breakdown, apply, test, validate, and finish. | Keeps agent work tied to explicit phases, artifacts, and exit checks. |
| Multi-IDE command hub | Ships `/skillgrid-*` and `/opsx-*` commands for Cursor, Kilo, OpenCode, and GitHub Copilot prompts. | One workflow can travel across the IDEs and agents you use. |
| Agent skills catalog | Provides reusable skills for TDD, review, security, UI design, research, GitNexus, Engram, OpenSpec, and more. | Agents get focused operating procedures instead of ad hoc chat instructions. |
| Local Skillgrid dashboard | Run `skillgrid serve` from [skillgrid-cli](skillgrid-cli/) (build with Bun) for PRD Kanban, Workflow, Checkpoints, Subagents, previews, and GitNexus UI/context links. | Product intent, events, previews, and subagent activity are visible in one local web UI. |
| Terminal TUI (Phase 1b) | Run `skillgrid tui --repo <project>` for a **dmux-style** multi-pane TUI: PRDs, OpenSpec changes, handoff files, merged event tail, and `checkpoints.log` tail (read-only, auto-refresh). | Same file-first signals as the web UI without a browser; pairs with [dmux](https://github.com/standardagents/dmux)-style terminal workflows. |
| File-first handoff | Stores PRDs, OpenSpec changes, handoff files, event logs, previews, checkpoints, and research under the repo. | Work survives context resets without requiring a database or hosted service. |
| Intent-gated loop | Adds `/skillgrid-loop` for the next safe phase or `[AFK]` slice, with explicit HITL and verification stop conditions. | Long-running agent work stays bounded by artifacts, evidence, and user authority. |
| Installer sanity check | Runs `./install.sh --sanity-check` to verify expected tools, hub files, and script syntax. | Setup problems are caught before copying configs into a project. |

---

## What you get

- **Skills** — `.agents/skills/` (OpenSpec lifecycle, SDD-style phases, code review, security, TDD, built-in `playwright` and `git-master`, GitNexus, Engram memory protocol, and more). See [docs/04-skills.md](docs/04-skills.md).
- **Commands** — Phase commands (`/skillgrid-*`) and OpenSpec commands (`/opsx-*`) for Cursor, Kilo, OpenCode, and GitHub Copilot prompts. See [docs/03-commands.md](docs/03-commands.md).
- **Workflow** — End-to-end phases from init through finish in [docs/02-workflow-usage.md](docs/02-workflow-usage.md).
- **Installer** — [`install.sh`](install.sh) syncs IDE folders, merges MCP JSON, copies `AGENTS.md`, and optionally installs CLIs (OpenSpec, GitNexus, dmux, Engram, and other optional tools from [docs/01-installation.md](docs/01-installation.md)).
- **Shared agent rules** — [`.configs/AGENTS.md`](.configs/AGENTS.md) is copied to the target project and IDE config dirs when you install.

---

## Workflow

**Phase 0** (`/skillgrid-session`) plus the Skillgrid **`/skillgrid-*`** steps summarized in [docs/02-workflow-usage.md](docs/02-workflow-usage.md), including optional **`/skillgrid-loop`** for controlled AFK progression and **`/skillgrid-validate`** as a single review-and-security gate. The diagram below is a six-phase mental model; exact steps and templates live in each command file. **OpenSpec**-focused steps also have **`/opsx-*`** aliases (see [docs/03-commands.md](docs/03-commands.md)).

**Init** writes **`.skillgrid/config.json`**: **ticketing** (`local` or a remote issue backend), **artifact store** (`hybrid` is the strongly recommended default, or `openspec` / `engram` for constrained setups), and **PRD workflow** (`prdWorkflow.statuses` and phase mapping) — so later commands know whether to use **`openspec/`** on disk, **Engram** memory, or both, and which local Kanban/status lifecycle to follow. See **`/skillgrid-init`** and [docs/02-workflow-usage.md](docs/02-workflow-usage.md).

```
  PHASE 0 (optional at session start)
 /skillgrid-session — charter, context budget, MCP selection, checkpoints

  DEFINE          PLAN           BUILD                  VERIFY (split)                    SHIP
 ┌──────┐        ┌──────┐        ┌──────┐    ┌──────────┬──────────┬──────────┐        ┌──────┐
 │ Init │ ─────▶ │ Plan │ ─────▶ │Apply │ ─▶ │   Test   │  Review  │ Security │ ─────▶ │Finish│
 │Explor│        │      │        │ Code │    │   Prove  │  (specs) │  (SAST)  │        │  PR  │
 │Brain │        │Break │        │      │    │          │          │          │        │      │
 └──────┘        └──────┘        └──────┘    └──────────┴──────────┴──────────┘        └──────┘
 /skillgrid-init /skillgrid-plan /skillgrid-apply  /skillgrid-test /skillgrid-validate /skillgrid-security /skillgrid-finish
 /skillgrid-design optional in DEFINE for UI/UX direction, previews, and DESIGN.md
```

Also run **`/skillgrid-explore`**, **`/skillgrid-design`**, and **`/skillgrid-brainstorm`** in DEFINE as needed, and **`/skillgrid-breakdown`** in PLAN. Use **`/skillgrid-loop`** only after artifacts identify a safe next phase or `[AFK]` slice; use **`/skillgrid-validate`** when you want review and security in a single turn.

---

## Quick start

Hub maintainers: after cloning, run **`npm ci`** in this repository root once so [package-lock.json](package-lock.json) pins Node tools under `node_modules/`.

1. Verify this hub is ready:

```bash
./install.sh --sanity-check
```

2. From this repository, install Skillgrid into **your** project directory:

```bash
./install.sh -p /path/to/your/project
```

3. In your project, start with Skillgrid:

```text
/skillgrid-init
```

4. Open the local dashboard when you want PRDs, workflow events, previews, subagent actions, and GitNexus in one place. From the hub repo, build the CLI then run serve from your project (or pass `--repo`):

```bash
cd skillgrid-cli && npm ci && npm run build && ./bin/skillgrid serve --repo /path/to/your/project
```

Then open `http://127.0.0.1:8787` (or the host/port you set).

For a **terminal-only** live view of events and checkpoints (Phase 1b), from the same built binary:

```bash
./bin/skillgrid tui --repo /path/to/your/project
```

Use `Tab` to move between scrollable panes, `r` to refresh, `q` to quit. See `skillgrid tui --help`.

Common variants:

- **All supported IDEs and defaults:** `./install.sh -p /path/to/project -A -y`
- **Dry run:** `./install.sh -p /path/to/project -n`
- **Dependencies only:** `./install.sh -d`
- **Sanity check only:** `./install.sh --sanity-check`
- **Optional CLIs** (OpenSpec, GitNexus, dmux, Engram, Brave Search, CocoIndex): add `-t` in an interactive terminal

You need a working **bash**, **rsync**, and **jq** for a full install with MCP merge; **Node** and **`npm ci`** in this repo help with pinned **npx**-based tools. For Python CLIs, use **uv** (see [docs/01-installation.md](docs/01-installation.md)).

---

## Repository layout (high level)

| Path | Purpose |
|------|---------|
| `.agents/skills/` | Agent skill definitions (`SKILL.md` per skill); shared `references/` checklists |
| `.cursor/`, `.kilo/`, `.opencode/` | IDE commands, rules, agents |
| `.github/` | Copilot prompts, workflows, agents |
| `.configs/` | Hub `AGENTS.md`, MCP fragments, OpenCode template config |
| `docs/` | Workflow, skills index, commands index, tools, notes |
| `install.sh` | Copy hub config into a target project and merge MCP |

---

## Documentation

| Doc | Contents |
|-----|----------|
| [docs/06-subagent-personas.md](docs/06-subagent-personas.md) | IDE agent personas (`.cursor/agents/`, mirrors, format) |
| [docs/02-workflow-usage.md](docs/02-workflow-usage.md) | Skillgrid phases, **`.skillgrid/config.json`**, PRD and OpenSpec handoff |
| `.skillgrid/project/` | Project-local architecture, structure, and onboarding docs when initialized |
| [docs/04-skills.md](docs/04-skills.md) | Catalog of all skills with paths and summaries |
| [docs/03-commands.md](docs/03-commands.md) | Slash commands and where they live per IDE |
| [docs/01-installation.md](docs/01-installation.md) | Install toolchain and workflow CLIs / IDEs |
| [docs/TODO.md](docs/TODO.md) | IDE matrix and internal notes |
| [docs/00-start-here.md](docs/00-start-here.md) | Start-here overview and manifesto: human-in-the-loop pipelines, spec-driven guidance, harness engineering |

---

## License

[MIT](LICENSE)
