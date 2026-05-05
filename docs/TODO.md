
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

* installer
  * [X] install.sh
  * [X] skillgrid-cli
  * [ ] brew
  * [ ] env variables
* plugins
  * [ ] Claude Code plugin
  * [ ] Cursor plugin
  * [ ] opencode plugin
* skillgrid-cli
  * [ ] installer
  * [ ] tui
  * [ ] webui
* components
  * mcp
    * [ ] only list nececary tools
    * agent-browser
  * skill/command rethink
    * based on sdd
    * sdd-orchestrator
      * https://github.com/Gentleman-Programming/gentle-ai
  * AGENTS.md
* agent personas
  * [ ] create planned agent personas
    * https://github.com/alvinunreal/oh-my-opencode-slim
    * https://github.com/opensoft/oh-my-opencode
* Sub agents
  * [-] paralel search
  * [-] persona board
  * [ ] Agent personas
  * [ ] context share
  * [ ] git worktrees
  * [ ] chackpoint
    * https://github.com/DeL-TaiseiOzaki/claude-code-orchestra/tree/main/.claude/skills/checkpointing
* Agentic workflo
  * [ ] agentic forkflo wit steps like Archon in skillgrid config.json
  * [ ] configure agent type, persona, modell for stap
* web ui
  * [X] web based kamban UI
  * [ ] multi agent wiev
  * [-] gitnexus ui
  * [?] openspec ui

# Tools

## Memory

* [ ] [Cleave](https://cleave.dev/)
* [X] [CocoIndex-Code](https://github.com/cocoindex-io/cocoindex-code)
* [ ] [GitNexus](https://github.com/abhigyanpatwari/GitNexus)
* [?] [GitNexus](https://github.com/abhigyanpatwari/GitNexus) 
* [X] [Engram](https://github.com/Gentleman-Programming/engram)
* [ ] [context-mode](https://github.com/mksglu/context-mode)
* [] [codemap](https://github.com/JordanCoin/codemap)

## Design Tools

* [ ] [taste-skill](https://github.com/Leonxlnx/taste-skill)
* [ ] [npxskillui](https://github.com/amaancoderx/npxskillui)
* [ ] [impeccable](https://github.com/pbakaus/impeccable)
* [ ] [open-design](https://github.com/nexu-io/open-design)

## Browser tools

* [ ] [agent-browser](https://github.com/vercel-labs/agent-browser)
* [ ] [PlayWright-mcp](https://github.com/microsoft/playwright-mcp)
  * [ ] [Skill-1](https://github.com/testdino-hq/playwright-skill)
  * [ ] [Skill-2](https://github.com/lackeyjb/playwright-skill)
* [X] [Context7—MCP]()
* [X] [exa-MCP]()
* [X] [DeepWiki—MCP]
* [X] [Firecrawl—MCP]()
* [ ] [Brave Search API](https://api.search.brave.com)
* [ ] [Tavily](https://tavily.com)

## Skills

* [ ] [opencode-conductor](https://github.com/derekbar90/opencode-conductor)
* [ ] [Antigravity Kit](https://github.com/vudovn/antigravity-kit)
* [X] [andrej-karpathy-skills](https://github.com/forrestchang/andrej-karpathy-skills)
* [ ] [Everything Claude Code](https://github.com/affaan-m/everything-claude-code)

## Agents (cli)

* claude
* codex
* opencode
  * kilo cli
* gemini
* copilot

## Agent Orchectration

* [ ] [claude-code-orchestra](https://github.com/DeL-TaiseiOzaki/claude-code-orchestra) - 
* [ ] [opencode-4hol](https://dev.to/uenyioha/porting-claude-codes-agent-teams-to-opencode-4hol) - 
* [ ] [oh-my-openagent](https://github.com/code-yeongyu/oh-my-openagent) - 
* [ ] [opencode-orchestrator](https://github.com/agnusdei1207/opencode-orchestrator) - 

## Agent multiplexers

* [ ] [dmux](https://github.com/standardagents/dmux) - Terminal
* [ ] [toad](https://github.com/batrachianai/toad) - Terminal in it's own window 
* [ ] [mux](https://github.com/coder/mux) - Separate App for macOS and Linux.
* [ ] [superset](https://github.com/superset-sh/superset) - Separate App for macOS
* [ ] [conductor](https://www.conductor.build/) - Separate App for macOS
* [ ] [Tmux-Orchestrator](https://github.com/Jedward23/Tmux-Orchestrator) - tmux based orcastrtator
* [ ] [ai-party](https://github.com/alexivison/ai-party) - ???
* [ ] [Codeman](https://github.com/Ark0N/Codeman) - tmux based multi session webui

## Agentix workflow
* [ ] [Archon](https://github.com/coleam00/archon) - Server Based

## Modell selector

* [ ] [cc-switch](https://github.com/farion1231/cc-switch)

## Sources

* [ ] [Aemini Skills](https://github.com/google-gemini/gemini-skills/)
* [ ] [context-mode](https://github.com/mksglu/context-mode/tree/main/configs)
* [ ] [antigravity-awesome-skills](https://github.com/sickn33/antigravity-awesome-skills)
* [ ] [awesome-codex-subagents](https://github.com/VoltAgent/awesome-codex-subagents)
* [ ] [AI Gentle Stack](https://github.com/Gentleman-Programming/gentle-ai)
* [ ] [Awesome Copilot](https://github.com/github/awesome-copilot)
  * [ ] [Awesome Copilot Agents](https://awesome-copilot.github.com/agents/)
* [ ] Role-based decision-board patterns
* Spec Driven Development
  * [ ] Isolated worktree, TDD, and double-review patterns
  * [ ] [agent-skills](https://github.com/addyosmani/agent-skills)
  * [ ] [Archon](https://github.com/coleam00/Archon)
  * [ ] [omegon-pi](https://github.com/styrene-lab/omegon-pi)
