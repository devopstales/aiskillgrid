# AISkillGrid

A **configuration hub** for opinionated AI-assisted development: reusable **skills**, **slash commands**, merged **MCP server** definitions, and an **install script** that copies normalized settings into your application repo. The workflow combines **OpenSpec**-style change management (CLI: `@fission-ai/openspec`) with **Spec-Driven Development** skills and production-oriented practices (testing, security, documentation).

---

## What you get

- **Skills** — `.agents/skills/` (OpenSpec lifecycle, SDD phases, code review, security, TDD, CocoIndex Code `ccc`, and more). See [docs/skills.md](docs/skills.md).
- **Commands** — Phase commands (`/skillgrid-*`) and OpenSpec commands (`/opsx-*`) for Cursor, Kilo, OpenCode, and GitHub Copilot prompts. See [docs/commands.md](docs/commands.md).
- **Workflow** — End-to-end phases from init through finish in [docs/wokflow.md](docs/wokflow.md).
- **Installer** — [`install.sh`](install.sh) syncs IDE folders, merges MCP JSON, copies `AGENTS.md`, and optionally installs CLIs (OpenSpec, graphifyy, CocoIndex Code, etc.). See [docs/tools.md](docs/tools.md).
- **Shared agent rules** — [`.configs/AGENTS.md`](.configs/AGENTS.md) is copied to the target project and IDE config dirs when you install.

---

## Workflow

Ten **`/skillgrid-*`** steps from [docs/wokflow.md](docs/wokflow.md), grouped like the six-phase diagram in [addyosmani/agent-skills](https://github.com/addyosmani/agent-skills). OpenSpec-focused steps also use **`/opsx-*`** ([docs/commands.md](docs/commands.md)).

```
  DEFINE          PLAN           BUILD             VERIFY           REVIEW                SHIP
 ┌──────┐        ┌──────┐        ┌──────┐          ┌──────┐        ┌──────┐            ┌──────┐
 │ Init │ ─────▶ │ Plan │ ─────▶ │Apply │ ───────▶ │ Test │ ─────▶ │Valid.│ ─────────▶ │Finish│
 │Explor│        │Design│        │ Code │          │ Prove│        │ Gate │            │  PR  │
 │Brain │        │Break │        │      │          │      │        │      │            │      │
 └──────┘        └──────┘        └──────┘          └──────┘        └──────┘            └──────┘
 /skillgrid-init /skillgrid-plan /skillgrid-apply  /skillgrid-test /skillgrid-validate /skillgrid-finish
```

Also run **`/skillgrid-explore`** and **`/skillgrid-brainstorm`** (DEFINE), **`/skillgrid-design`** and **`/skillgrid-breakdown`** (PLAN).

---

## Quick start

From this repository, run the installer against **your** project directory:

```bash
./install.sh -p /path/to/your/project
```

Common variants:

- **All supported IDEs and defaults:** `./install.sh -p /path/to/project -A -y`
- **Dry run:** `./install.sh -p /path/to/project -n`
- **Dependencies only:** `./install.sh -d`
- **Optional CLIs** (OpenSpec, graphify, CocoIndex Code, dmux, Engram): add `-t` in an interactive terminal

You need a working **bash**, **rsync**, and **jq** for a full install with MCP merge; **Node** / **npx** helps for npm-based global tools. Details: [docs/tools.md](docs/tools.md).

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
| [docs/wokflow.md](docs/wokflow.md) | Skillgrid phases and which skills apply per phase |
| [docs/skills.md](docs/skills.md) | Catalog of all skills with paths and summaries |
| [docs/commands.md](docs/commands.md) | Slash commands and where they live per IDE |
| [docs/tools.md](docs/tools.md) | Install toolchain (brew, node, pipx, …) and workflow CLIs / IDEs |
| [docs/TODO.md](docs/TODO.md) | IDE matrix and internal notes |

---

## License

[MIT](LICENSE)
