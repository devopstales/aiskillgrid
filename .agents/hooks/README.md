# tmux-agent-status Hooks

These hooks integrate Skillgrid with [tmux-agent-status](https://github.com/samleeney/tmux-agent-status) to show at-a-glance which tmux sessions have AI agents working vs idle.

## Architecture

All hooks share a common library (`.agents/hooks/lib.sh`) that handles:
- tmux session detection (local + remote)
- Status file writing (`~/.cache/tmux-agent-status/<session>.status`)
- Pane-level status tracking
- Wait/park mode management
- Sidebar refresh signaling

Each agent-specific hook sources `lib.sh` and maps its native lifecycle events to `working`/`done`/`wait` states.

## Setup by IDE

### Claude Code

Add to `~/.claude/settings.json`:

```json
{
  "hooks": {
    "UserPromptSubmit": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "bash ~/.agents/hooks/claude-hook.sh UserPromptSubmit"
          }
        ]
      }
    ],
    "PreToolUse": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "bash ~/.agents/hooks/claude-hook.sh PreToolUse"
          }
        ]
      }
    ],
    "Stop": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "bash ~/.agents/hooks/claude-hook.sh Stop"
          }
        ]
      }
    ],
    "Notification": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "bash ~/.agents/hooks/claude-hook.sh Notification"
          }
        ]
      }
    ]
  }
}
```

### Kilo

Kilo uses the `.kilo/hook/` mechanism. Add hook triggers to `.kilo/hook/hooks.md` or wire via `kilo.json`:

```json
{
  "hooks": {
    "PostShell": [
      { "command": "bash .agents/hooks/kilo-hook.sh PostShell" }
    ],
    "PreToolUse": [
      { "command": "bash .agents/hooks/kilo-hook.sh PreToolUse" }
    ],
    "Stop": [
      { "command": "bash .agents/hooks/kilo-hook.sh Stop" }
    ]
  }
}
```

### OpenCode

Add to `.opencode/opencode.jsonc` or `~/.config/opencode/config.json`:

```json
{
  "hooks": {
    "UserPromptSubmit": "bash <path>/opencode-hook.sh UserPromptSubmit",
    "PreToolUse": "bash <path>/opencode-hook.sh PreToolUse",
    "Stop": "bash <path>/opencode-hook.sh Stop",
    "SessionStart": "bash <path>/opencode-hook.sh SessionStart"
  }
}
```

### Gemini CLI

Add to `.gemini/settings.json` or `~/.gemini/settings.json`:

```json
{
  "hooks": {
    "UserPromptSubmit": "bash <path>/gemini-hook.sh UserPromptSubmit",
    "PreToolUse": "bash <path>/gemini-hook.sh PreToolUse",
    "Stop": "bash <path>/gemini-hook.sh Stop",
    "SessionStart": "bash <path>/gemini-hook.sh SessionStart"
  }
}
```

### Codex CLI

Add to `~/.codex/config.toml`:

```toml
[features]
codex_hooks = true
```

Then add to `~/.codex/hooks.json` (or `.codex/hooks.json` per project):

```json
{
  "hooks": {
    "SessionStart": [
      {
        "matcher": "startup|resume",
        "hooks": [
          {
            "type": "command",
            "command": "bash <path>/codex-hook.sh SessionStart"
          }
        ]
      }
    ],
    "UserPromptSubmit": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "bash <path>/codex-hook.sh UserPromptSubmit"
          }
        ]
      }
    ],
    "PreToolUse": [
      {
        "matcher": "Bash",
        "hooks": [
          {
            "type": "command",
            "command": "bash <path>/codex-hook.sh PreToolUse"
          }
        ]
      }
    ],
    "Stop": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "bash <path>/codex-hook.sh Stop"
          }
        ]
      }
    ]
  }
}
```

### Pi Coding Agent

Add to `~/.pi/agent/settings.json` or `.pi/settings.json`:

```json
{
  "hooks": {
    "SessionStart": [
      { "matcher": "startup", "hooks": [{ "type": "command", "command": "bash <path>/pi-hook.sh SessionStart" }] },
      { "matcher": "resume", "hooks": [{ "type": "command", "command": "bash <path>/pi-hook.sh SessionStart" }] }
    ],
    "UserPromptSubmit": [
      { "hooks": [{ "type": "command", "command": "bash <path>/pi-hook.sh UserPromptSubmit" }] }
    ],
    "PreToolUse": [
      { "hooks": [{ "type": "command", "command": "bash <path>/pi-hook.sh PreToolUse" }] }
    ],
    "Stop": [
      { "hooks": [{ "type": "command", "command": "bash <path>/pi-hook.sh Stop" }] }
    ],
    "SessionEnd": [
      { "hooks": [{ "type": "command", "command": "bash <path>/pi-hook.sh Stop" }] }
    ]
  }
}
```

Pi hooks use the [pi-hooks](https://github.com/hsingjui/pi-hooks) compatibility layer, which maps Claude Code-style hook config to Pi's TypeScript extension event system.

### Cursor

Cursor does not have a native hook system. Use one of these approaches:

**Option 1: Cursor command** — Create `.cursor/commands/agent-status.json`:

```json
{
  "name": "agent-status",
  "description": "Update tmux-agent-status",
  "command": "bash <path>/cursor-hook.sh UserPromptSubmit"
}
```

**Option 2: Shell alias** — Add to your shell config:

```bash
alias cursor-working="bash ~/.agents/hooks/cursor-hook.sh UserPromptSubmit"
alias cursor-done="bash ~/.agents/hooks/cursor-hook.sh Stop"
```

### Copilot CLI

Copilot CLI does not have a native hook system. Use shell aliases:

```bash
alias copilot-working="bash ~/.agents/hooks/copilot-hook.sh UserPromptSubmit"
alias copilot-done="bash ~/.agents/hooks/copilot-hook.sh Stop"
```

## Status States

| State | Meaning |
|-------|---------|
| `working` | Agent is actively processing |
| `done` | Agent is idle, waiting for input |
| `wait` | Agent is in timed wait mode |
| `parked` | Agent is parked for later |

## File Layout

```
.agents/hooks/
  lib.sh              ← Shared library (sourced by all hooks)
  claude-hook.sh      ← Claude Code lifecycle hooks
  codex-hook.sh       ← Codex CLI lifecycle hooks
  kilo-hook.sh        ← Kilo lifecycle hooks
  opencode-hook.sh    ← OpenCode lifecycle hooks
  cursor-hook.sh      ← Cursor (manual/command-based)
  gemini-hook.sh      ← Gemini CLI lifecycle hooks
  pi-hook.sh          ← Pi Coding Agent lifecycle hooks
  copilot-hook.sh     ← Copilot CLI (manual/alias-based)
  refresh-indexes.sh  ← Existing: index refresh after merge/pull
  README.md           ← This file
```

## Custom Agents

Any tool can integrate by writing status files directly:

```bash
# Session-level status
echo "working" > ~/.cache/tmux-agent-status/<session>.status
echo "myagent" > ~/.cache/tmux-agent-status/<session>.agent

# Pane-level status
echo "working" > ~/.cache/tmux-agent-status/panes/<session>_<pane>.status
echo "myagent" > ~/.cache/tmux-agent-status/panes/<session>_<pane>.agent
```
