# AISkillGrid — Project

**AISkillGrid** is a **configuration hub** for AI-assisted development: reusable **skills** (`.agents/skills/`), **slash commands** (`.cursor/commands/` and mirrors), **IDE agent personas** (`.cursor/agents/` and mirrors), merged **MCP** definitions, and **`install.sh`** to copy normalized settings into an application repository.

## Audience

- **Maintainers:** evolve the hub, keep Cursor sources of truth in sync via `scripts/sync-ide-assets.sh`, and run CI `hub-sync-check` on PRs.
- **Consumers:** run `install.sh` from a clone or release to wire skills, commands, agents, and MCP into a target project.
- **Agents:** follow root **`AGENTS.md`**, phase commands **`/skillgrid-*`**, and skills under `.agents/skills/`.

## Repository map (high level)

| Path | Role |
|------|------|
| `.agents/skills/` | Agent Skills (`SKILL.md` per skill) |
| `.cursor/commands/` | Canonical `skillgrid-*` and `opsx-*` commands |
| `.cursor/agents/` | Canonical IDE personas |
| `.configs/` | Hub `AGENTS.md` and shared JSON fragments |
| `docs/` | Workflow, skills index, commands, tools, agents |
| `scripts/sync-ide-assets.sh` | Sync commands + agents to Kilo, OpenCode, GitHub |
| `install.sh` | Install hub artifacts into a target repo |

## How to work in this repo

- **Phase 0:** `/skillgrid-session` when context resets.
- **Change the hub:** edit `.cursor/commands/` and `.cursor/agents/` first, then run `./scripts/sync-ide-assets.sh`.
- **Full lifecycle:** see [docs/wokflow.md](../docs/wokflow.md) and [docs/commands.md](../docs/commands.md).

## Where to read next

- [ARCHITECTURE.md](ARCHITECTURE.md) — components and flows
- [STRUCTURE.md](STRUCTURE.md) — directory layout
- [docs/agents.md](../docs/agents.md) — IDE personas vs skills
