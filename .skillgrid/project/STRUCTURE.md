# AISkillGrid — Structure

> **Scope:** Repository layout for this hub (not consumer app trees).

## High-level tree

```text
aiskillgrid/
├── .agents/skills/          # SKILL.md trees (canonical skills)
├── .cursor/
│   ├── agents/              # IDE personas (source of truth)
│   ├── commands/            # skillgrid-* / opsx-* (source of truth)
│   └── rules/               # Cursor rules
├── .kilo/                   # Mirror: commands, agents, mcp (if used)
├── .opencode/               # Mirror: commands, agents, config
├── .github/
│   ├── prompts/             # Copilot: prompts derived from .cursor/commands
│   ├── agents/              # Mirror of .cursor/agents
│   └── workflows/           # e.g. hub-sync-check.yml
├── .configs/                # AGENTS.md template, opencode.json, etc.
├── .skillgrid/
│   └── project/             # This hub’s own exploration docs
├── docs/                    # wokflow, commands, skills, tools, agents
├── scripts/
│   └── sync-ide-assets.sh   # Cursor → other IDEs
├── install.sh
└── README.md
```

## Sync rules

- **Commands:** `skillgrid-*.md`, `opsx-*.md` copied from `.cursor/commands/` to `.kilo/commands/`, `.opencode/commands/`; GitHub prompts regenerate with `description`-only frontmatter.
- **Agents:** all `.md` under `.cursor/agents/` copied to `.kilo/agents/`, `.opencode/agents/`, `.github/agents/`.
- **Skills:** `install.sh` syncs `.agents/skills/` into per-IDE skill dirs for **target** projects (see `install.sh`); this hub keeps skills only under `.agents/skills/`.

## Runtime / deployment

This repository is **not** a deployed service; “runtime” is contributors’ machines and **GitHub Actions** (sync check).
