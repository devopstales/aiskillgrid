# IDEs

* Embeddid
  * Google Antigravity
  * Cursor
* Plugins
  * Github Copilot
  * Kilo Code
  * OpenChamber
* Agents
  * kilo cli
  * Opencode
  * Gemini
  * Claude Code
  * Codex
  * pi

# IDE Configs

## Embedded

- Google Antigravity
- Cursor

| Config Area | Google Antigravity | Cursor |
| ----------- | ------------------ | ------ |
| Commands | `.agents/workflows` | `.cursor/commands/`? |
| Skills | `.agents/skills/` | `.agents/skills/`, `.cursor/skills/` |
| Rules | `.agents/rules/`, `AGENTS.md` | `.cursor/rules/`, `AGENTS.md` |
| Agents | `.agents/agents.md`, `.agent/agents/` | `.cursor/agents/` |
| Hooks | `.agents/hooks/` | `.cursor/hooks.json`, `.cursor/hooks/` |
| MCP | `~/.gemini/antigravity/mcp_config.json`, `.agent/mcp_config.json` | `.cursor/mcp.json` |
| Plugins | X | `.cursor/plugins/` |

## Plugins

- Github Copilot
- Kilo Code
- OpenChamber

| Config Area | Github Copilot | Kilo Code | OpenChamber |
| ----------- | -------------- | --------- | ----------- |
| Commands | `.github/prompts/` | `.kilo/commands/` | `.opencode/commands/*.md` |
| Skills | `.agents/skills/`, `.claude/skills/`, `.github/skills`, `.copilot/skills` | `.agents/skills/`, `.claude/skills/`, `.kilo/skills/` | `.agents/skills/`, `.claude/skills/`, `.opencode/skills/` |
| Rules | `.github/instructions/*.instructions.md`, `.github/copilot-instructions.md`, `AGENTS.md` | `.kilo/rules/`, `AGENTS.md` | `.opencode/rules/`, `AGENTS.md` |
| Agents | `.github/agents/` | `.kilo/agents/` | `.opencode/agents/` |
| Hooks | `.github/hooks/*.json` | ? | `.opencode/hook/hooks.md` |
| MCP | `.vscode/mcp.json` | `.kilo/kilo.jsonc` | `.opencode/opencode.jsonc` |
| Plugins | X | X | `.opencode/plugins/` |

## Agents

- kilo cli
- Opencode
- Gemini
- Claude Code
- Codex
- pi

| Config Area | kilo cli | Opencode | Gemini | Claude Code | Codex | pi |
| ----------- | -------- | -------- | ------ | ----------- | ----- | -- |
| Commands | `.kilo/commands/` | `.opencode/commands/*.md` | `.gemini/commands/*.toml` | ? | ? | ? |
| Skills | `.agents/skills/`, `.claude/skills/`, `.kilo/skills/` | `.agents/skills/`, `.claude/skills/`, `.opencode/skills/` | `.agents/skills/` | `.agents/skills/`, `.claude/skills/` | `.agents/skills/` | ? |
| Rules | `.kilo/rules/`, `AGENTS.md` | `.opencode/rules/`, `AGENTS.md` | ? | `AGENTS.md` | `AGENTS.md` | ? |
| Agents | `.kilo/agents/` | `.opencode/agents/` | ? | ? | ? | ? |
| Hooks | ? | `.opencode/hook/hooks.md` | `.gemini/settings.json` | ? | ? | ? |
| MCP | `.kilo/kilo.jsonc` | `.opencode/opencode.jsonc` | `.gemini/settings.json` | ? | ? | ? |
| Plugins | X | `.opencode/plugins/` | ? | ? | ? | ? |

## Norse Persona Config Set

Norse persona theming and routing are centralized in shared config files:

- `.configs/norse-persona-contract.json` (persona roles, routing matrix, hard rules, return contract)
- `.configs/ide-model-mapping.json` (surface model mapping by `fast|balanced|deep`)
- `.configs/ide-persona-capabilities.json` (surface capability tiers and fallback policy)

Render cross-IDE persona manifests with:

```bash
node scripts/render-multi-ide-personas.mjs
```

This writes `norse-personas.json` files into each surface (`.cursor`, `.kilo`, `.opencode`, `.github`, `.gemini`, `.codex`, `.pi`, `.antigravity`, `.vscode`, `.claude`) so command routing can stay consistent across toolchains.

### Google Antigravity

### Cursor

```json
.cursor/hooks.json
{
  "version": 1,
  "hooks": {
    "preToolUse": [
      {
        "command": "context-mode hook cursor pretooluse",
        "matcher": "Shell|Read|Grep|WebFetch|Task|MCP:ctx_execute|MCP:ctx_execute_file|MCP:ctx_batch_execute"
      }
    ],
    "postToolUse": [
      {
        "command": "context-mode hook cursor posttooluse"
      }
    ]
  }
}
```

```json
.cursor/mcp.json
{
  "mcpServers": {
    "context7": {
      "args": [
        "-y",
        "@upstash/context7-mcp"
      ],
      "command": "npx"
    },
    "engram": {
      "args": [
        "mcp",
        "--tools=agent"
      ],
      "command": "engram"
    }
  }
}
```


### VSCode / Github Copilot

```json
.vscode/mcp.json
{
  "servers": {
    "context7": {
      "type": "http",
      "url": "https://mcp.context7.com/mcp"
    },
    "engram": {
      "args": [
        "mcp",
        "--tools=agent"
      ],
      "command": "engram"
    }
  }
}
```

### Kilo Code / CLI


### Opencode

```json
.opencode/opencode.json
{
  "$schema": "https://opencode.ai/config.json",
  "instructions": ["CONTRIBUTING.md", "docs/guidelines.md", ".opencode/rules/*.md"],
  "plugin": [
    "./plugins"
  ]
}
```

```json
.opencode/opencode.json
{
  "$schema": "https://opencode.ai/config.json",
  "mcp": {
    "context7": {
      "enabled": true,
      "type": "remote",
      "url": "https://mcp.context7.com/mcp"
    },
    "engram": {
      "command": [
        "engram",
        "mcp",
        "--tools=agent"
      ],
      "enabled": true,
      "type": "local"
    }
  }
}
```

### pi

```bash
npm install -g @mariozechner/pi-coding-agent
```
