# Research: earendil-works/pi Architecture

> Source: https://github.com/earendil-works/pi (DeepWiki analysis, 2026-09-08)

## What is pi?

Pi is an extensible AI coding agent designed as a minimal terminal coding harness. It provides an interactive CLI for LLM-assisted development with tool calling, session persistence, and a rich extension system. It is a TypeScript monorepo using npm workspaces.

## Package Architecture

| Package | Purpose |
|---------|---------|
| `@earendil-works/pi-ai` | Unified LLM provider abstraction (OpenAI, Anthropic, Google, etc.) |
| `@earendil-works/pi-agent-core` | Agent loop, state management, tool calling, AgentHarness |
| `@earendil-works/pi-tui` | Differential terminal rendering engine |
| `@earendil-works/pi-coding-agent` | Main CLI, coding tools, SDK, session management |
| `@earendil-works/chord` | Distributed runtime facets/services (experimental) |
| `@earendil-works/pi-protocol` | Wire protocol, CBOR framing |
| `@earendil-works/pi-client` / `pi-server` | Client-server session transport |
| `@earendil-works/pi-telemetry` | Typed AI-request telemetry |
| `@earendil-works/pi-evals` | Behavioral evaluation framework |

## Core Architecture

### Layered Design

```
Interface (pi-tui, RPC, Print modes)
    ↓
Application (pi-coding-agent: AgentSession, SessionManager, ModelRuntime)
    ↓
Core Logic (pi-agent-core: Agent, agentLoop, AgentHarness)
    ↓
AI Abstraction (pi-ai: Models, streamSimple, CredentialStore)
    ↓
LLM Providers
```

### Agent Loop (pi-agent-core)

- `Agent` class: stateful wrapper with `AgentState` (messages, model, tools, streaming status)
- `runAgentLoop()`: iterative cycle — transform context → stream LLM → execute tools → check stop
- `AgentHarness`: durable runtime with lanes, checkpoints, compaction, crash recovery
- Event stream: `agent_start`, `turn_start`, `message_update`, `tool_execution_*`, `turn_end`, `agent_end`
- Tool execution: parallel (default) or sequential
- Steering: user can inject messages mid-turn

### Session Management

- JSONL-based persistence (append-only)
- Tree structure with `id` + `parentId` for branching/forking
- Automatic context compaction
- SQLite session backend (optional)
- Session recovery from crash

### LLM Provider Abstraction (pi-ai)

- `ModelRuntime`: centralizes model config, auth, provider catalogs
- `Models.streamSimple()`: unified streaming API
- `CredentialStore`: API keys + OAuth subscriptions
- Model resolution with thinking levels (off → max)
- Dynamic provider discovery

### Built-in Tools

- `read`, `write`, `edit`, `bash` (core)
- `find`, `grep`, `ls` (search)
- `powershell` (Windows, optional)
- Tools are extensible via the Extension API

### Extension System

- TypeScript modules loaded via `jiti`
- `ExtensionAPI`: register tools, commands, keybindings, event handlers, UI components
- Skills: on-demand capability packages (Agent Skills standard)
- Prompt templates: reusable markdown prompts
- Project trust: `ProjectTrustStore` security boundary

### TUI (pi-tui)

- Differential rendering (efficient terminal updates)
- Multi-line editor with grapheme-aware cursor
- Kitty graphics support
- Fullscreen UI mode
- Keybinding manager

### Execution Modes

1. **Interactive** (default): TUI with full agent loop
2. **Print** (`-p`): one-shot stdout output
3. **JSON**: structured JSON output
4. **RPC**: JSON-RPC 2.0 over stdin/stdout

### Build & Distribution

- TypeScript compiled with `tsgo`
- Standalone binary via `bun build --compile`
- Cross-platform: darwin/linux/windows × arm64/x64
- Lockstep versioning across packages

## Key Design Principles

1. **Extensibility over built-in features**: no sub-agents, plan mode, etc. baked in
2. **Modular packages**: strict layer separation
3. **Durable execution**: crash recovery, parallel lanes
4. **Provider-agnostic**: unified API across LLM providers
5. **Minimal core**: small agent loop, rich extension surface

## Implications for skill-pi

- Can build on top of `pi-agent-core` + `pi-ai` as the foundation
- `pi-coding-agent` provides the reference implementation for CLI/TUI/session
- Extension system is the primary customization surface (tools, commands, skills)
- Skillgrid skills can be loaded as pi skills via the Agent Skills standard
- Mnemonic integration would be a custom extension or tool
