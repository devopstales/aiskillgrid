# IDEs

* Google Antigravity
* Cursor
* Github Copilot
* Kilo Code
  * kilo cli
* Opencode
  * OpenChamber

## IDE Configs

|           | Copilot | Cursor | Kilo | Opencode | Antigravity | Gemini |
| --------- | ------- | ------ | ---- | -------- | ----------- | ------ |
| Commands  | `.github/prompts/` | `.cursor/commands/`? | `.kilo/commands/` | `.opencode/commands/*.md` | `.agents/workflows` | `.gemini/commands/*.toml` |
| Skills    | `.agents/skills/`, `.claude/skills/`, `.github/skills`, `.copilot/skills` | `.agents/skills/`, `.cursor/skills/` | `.agents/skills/ `, `.claude/skills/ `, `.kilo/skills/` | `.agents/skills/`, `.claude/skills/`, `.opencode/skills/` | `.agents/skills/` | `.agents/skills/` |
| Rules     | `.github/instructions/*.instructions.md`, `.github/copilot-instructions.md`, `AGENTS.md` | `.cursor/rules/`, `AGENTS.md` | `.kilo/rules/`, `AGENTS.md` | `.opencode/rules/`, `AGENTS.md` | `.agents/rules/`, `AGENTS.md` | ? |
| Agents    | `.github/agents/` | `.cursor/agents/` | `.kilo/agents/` | `.opencode/agents/` | ? | ? |
| Hooks     | `.github/hooks/*.json` | `.cursor/hooks.json`, `.cursor/hooks/` | ? | `.opencode/hook/hooks.md` | `.agents/hooks/` | `.gemini/settings.json` |
| MCP       | `.vscode/mcp.json` | `.cursor/mcp.json` | `.kilo/kilo.jsonc` | `.opencode/opencode.jsonc` | `~/.gemini/antigravity/mcp_config.json` | `.gemini/settings.json` |
| plugins   | X | `.cursor/plugins/` | X | `.opencode/plugins/` | X | ? |

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