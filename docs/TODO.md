
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

# Plans and Bugs

* env bariables
* commands
  * `/skillgirid-plan` sas `/osx:apply` is the next
* mcp
* planning
  * database structure
  * mermaid graph in db plan
* ticketing integration
  * github
  * gitlab
  * jira cloud
  * web based kamban UI
* paralelism
  * subagent-driven-development or executing-plans
  * checkpont
    * git-worktrees (separate branch?)
    * checkpoint commands
    * create separate git branch for prds and work on that
    * at the and create PR
* generate user journey 
  * to Markdown Presentation Ecosystem format
  * to mermeade
* installer
  * Claude Code plugin
  * Cursor plugin
  * opencode plugin
  * node
  * skill.sh

## Sources

* [X] [andrej-karpathy-skills](https://github.com/forrestchang/andrej-karpathy-skills)
* [ ] [Everything Claude Code](https://github.com/affaan-m/everything-claude-code)
* [ ] [Antigravity Kit](https://github.com/vudovn/antigravity-kit)
* [ ] [Aemini Skills](https://github.com/google-gemini/gemini-skills/)
* [ ] [context-mode](https://github.com/mksglu/context-mode/tree/main/configs)
* [ ] [antigravity-awesome-skills](https://github.com/sickn33/antigravity-awesome-skills)
* [ ] [awesome-codex-subagents](https://github.com/VoltAgent/awesome-codex-subagents)
* [ ] [AI Gentle Stack](https://github.com/Gentleman-Programming/gentle-ai)
* [ ] [Awesome Copilot](https://github.com/github/awesome-copilot)
  * [ ] [Awesome Copilot Agents](https://awesome-copilot.github.com/agents/)
* [ ] [gstack](https://github.com/garrytan/gstack)
* Spec Driven Development
  * [ ] [superpowers](https://github.com/obra/superpowers)
  * [ ] [agent-skills](https://github.com/addyosmani/agent-skills)
  * [ ] [Archon](https://github.com/coleam00/Archon)
  * [ ] [omegon-pi](https://github.com/styrene-lab/omegon-pi)
* tools
  * [ ] [Cleave]