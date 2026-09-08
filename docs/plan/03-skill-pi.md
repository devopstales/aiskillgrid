# Plan: skill-pi — Full-Featured Skillgrid AI Coding Agent

> **Status:** Planning
> **Change:** 012-skill-pi-coding-agent
> **Ticket:** TASK-006
> **Date:** 2026-09-08
> **Merged from:** 009-skillgrid-pi (skillgrid-old) + user decisions + ecosystem research

## Vision

skill-pi is a fat npm distribution that vendors `@earendil-works/pi-coding-agent` and layers Skillgrid SDD workflow, Mnemonic persistent memory, SDD phase sub-agents, a custom dashboard, robust permissions, local LLM inference (Ollama + vLLM), and rpiv companions on top.

The command is `skill-pi`. The home folder is `~/skill-pi/`. The theme is tokyonight. The logo is skillgrid. The style is omegon-pi's polished splash + dashboard aesthetic.

The thesis: pi is the engine. Skillgrid is the brain. skill-pi is the car.

## Identity Decisions (Locked)

| # | Topic | Decision | Source |
|---|-------|----------|--------|
| 1 | Packaging | **Fat distribution** — npm package depends on/bundles pi, ships own bin | 009 + user |
| 2 | Command | **`skill-pi`** (not `pi`, not `skillgrid-pi`) | user correction |
| 3 | Home folder | **`~/skill-pi/`** (not `~/.skillgrid-pi/`, not `~/.pi/agent/`) | user correction |
| 4 | Package name | **`skill-pi`** (npm) | matches command |
| 5 | Repo path | **`skill-pi/`** (repo root) | matches 009 pattern |
| 6 | Install | **Both**: `npm i -g skill-pi` + `skillgrid install --agents skill-pi` | 009 + user |
| 7 | Go installer | Adds `skill-pi` to `AvailableAgents()` in skillgrid-cli | 009 |
| 8 | Asset sync | Hash-manifest sync of prompts/skills/agents/chains to `~/skill-pi/` | 009 |
| 9 | Core identity | **Full-featured all-in-one**: SDD + memory + code index + sub-agents + dashboard + permissions + local LLM + web + todo | user |
| 10 | SDD scope | **Full cycle** as subagent workflows (explore, propose, spec, apply, verify, archive) | user |
| 11 | Mnemonic | **pi-mcp-adapter** npm package for MCP connection | user |
| 12 | Sub-agents | **6 SDD phase agents** + Result Contracts (from gentle-pi pattern) | user + gentle-pi |
| 13 | Dashboard | **Custom TUI dashboard** (omegon-pi style) | user |
| 14 | Permissions | **gotgenes gates + inobit ask/plan UX** + `~/skill-pi/permission.json` | 009 |
| 15 | Local LLM | **Ollama + vLLM auto-detect** in `/login`, no JSON editing | user + 009 |
| 16 | rpiv companions | **3 pinned**: `@juicesharp/rpiv-ask-user-question`, `rpiv-todo`, `rpiv-web-tools` | 009 |
| 17 | Theme | **Tokyonight** + skillgrid logo (reuse `plugins/*/skillgrid-logo.tsx`) | user |
| 18 | Splash | **omegon-style** startup banner with logo + status | user + 009 |
| 19 | Web | **pi-web-access** npm package | user |
| 20 | MCP | **pi-mcp-adapter** for Mnemonic connection | user |

## Distribution

```
npm install -g skill-pi
skillgrid install --agents skill-pi
```

- **Command:** `skill-pi`
- **Home folder:** `~/skill-pi/`
- **Pattern:** Fat npm package (009 pattern). `bin/skill-pi.mjs` spawns pi's CLI as a child process, injects `--agent-dir ~/skill-pi/` plus bundled extensions, suppresses pi's auto-discovery (`--no-skills --no-prompt-templates --no-themes`), and runs an exit-75 restart loop for `/update`/`/restart`.
- **pi is a dependency**, not a peer requirement. The user never installs pi separately.

### State Layout

```
~/skill-pi/                    # home folder (state + config)
├── config.json                # global settings (model, theme, permissions)
├── sessions/                  # JSONL session files
├── extensions/                # global extensions (auto-discovered)
├── agents/                    # SDD phase subagent definitions (synced)
├── skills/                    # SDD skills (synced)
├── prompts/                   # prompt templates (synced)
├── chains/                    # SDD phase chains (synced)
├── themes/                    # themes (tokyonight default)
├── permission.json            # permission rules
├── mcp.json                   # MCP server configs (skillgrid)
└── npm/                       # installed pi packages (pi install)
```

## Repo Layout

```
skill-pi/                      # source in skillgrid monorepo (repo root)
├── bin/
│   └── skill-pi.mjs           # CLI entry: spawn pi, inject extensions, restart loop
├── extensions/
│   ├── 00-splash/             # startup banner, logo, tokyonight theme, project detection
│   ├── 01-mnemonic/           # pi-mcp-adapter bridge → Mnemonic tools
│   ├── 02-sdd/                # SDD workflow: routing, phase subagents, Result Contracts
│   ├── 03-permissions/        # gotgenes gates + inobit ask/plan UX
│   ├── 04-dashboard/          # custom TUI dashboard (omegon-pi style)
│   ├── 05-local-llm/          # Ollama + vLLM auto-detection, /login integration
│   └── lib/
│       ├── shared-state.ts    # in-memory state bus (cross-extension events)
│       ├── config.ts          # config.yaml discovery + parsing
│       ├── agent-home.ts      # resolve ~/skill-pi/, hash-manifest sync
│       └── result-contract.ts # SDD Result Contract type + validation
├── agents/                    # SDD phase subagent definitions (source for sync)
│   ├── sdd-explore.md
│   ├── sdd-propose.md
│   ├── sdd-spec.md
│   ├── sdd-apply.md
│   ├── sdd-verify.md
│   └── sdd-archive.md
├── skills/                    # SDD skills (source for sync, from .agents/skills/)
├── prompts/                   # SDD phase prompt templates (source for sync)
├── chains/                    # SDD phase chains (source for sync)
├── themes/
│   └── tokyonight.json        # tokyonight theme (pi built-in, customized)
├── test/
│   ├── scaffold.test.ts
│   ├── extensions.test.ts
│   ├── asset-sync.test.ts
│   ├── sdd-mnemonic.test.ts
│   ├── installer.test.ts
│   ├── branding.test.ts
│   ├── local-api.test.ts
│   ├── rpiv.test.ts
│   ├── permission.test.ts
│   ├── subagents.test.ts
│   ├── dashboard.test.ts
│   ├── docs.test.ts
│   └── e2e.test.ts
├── package.json               # npm global, bin: skill-pi, pi manifest
├── tsconfig.json
└── README.md
```

## Step Blueprint (12 Steps)

Merged from 009's 10 steps + 2 new steps for sub-agents and dashboard.

| NN | Step slug | Goal (one line) | Primary package / entry | Depends on |
|----|-----------|-----------------|-------------------------|------------|
| 01 | `package-scaffold` | Repo-root `skill-pi/` npm package with pinned Pi deps, bin `skill-pi` only, no `pi` alias | `skill-pi/package.json`, `skill-pi/bin` | — |
| 02 | `extensions-core` | Packaged extensions load on start (status/doctor hooks, sync trigger) | `skill-pi/extensions` | 01 |
| 03 | `asset-sync` | Prompts/skills/agents/chains sync into `~/skill-pi/` with hash-manifest policy | `skill-pi/lib/agent-home.ts` + `assets/` | 02 |
| 04 | `sdd-mnemonic` | SDD phase assets mapped; Mnemonic via pi-mcp-adapter; fail-soft when unavailable | `skill-pi/extensions/01-mnemonic`, `02-sdd` | 03 |
| 05 | `installer` | `skillgrid install --agents skill-pi` installs/links the package; Go `AvailableAgents` updated | `skillgrid-cli/internal/install` | 01 |
| 06 | `branding-look` | Tokyo Night theme + startup splash/logo (omegon-style) using skillgrid brand | `skill-pi/themes`, `skill-pi/extensions/00-splash` | 02 |
| 07 | `local-api-login` | `/login` Local API entry: Ollama/vLLM auto-detect, wizard, connection test, register | `skill-pi/extensions/05-local-llm`, `lib/local-backends/` | 02 |
| 08 | `rpiv-companions` | Pin/activate `@juicesharp/rpiv-ask-user-question`, `rpiv-todo`, `rpiv-web-tools` | `skill-pi/package.json`, Pi package resources | 01 |
| 09 | `permission` | gotgenes gates + inobit ask/plan UX; config at `~/skill-pi/permission.json` | `skill-pi/extensions/03-permissions`, `lib/permission/` | 02 |
| 10 | `sdd-subagents` | 6 SDD phase subagents + Result Contracts + `/sdd` routing + `sdd-full` chain | `skill-pi/agents/`, `extensions/02-sdd` | 04 |
| 11 | `dashboard` | Custom TUI dashboard: SDD status, memory, code index, session, agents | `skill-pi/extensions/04-dashboard` | 04, 10 |
| 12 | `docs-rollout` | User docs + plan harness list; install/update/coexistence; smoke checklist | `docs/user-manual`, `docs/plan` | 04, 05, 06, 07, 08, 09, 10, 11 |

## Extensions

### 00-splash

**Purpose:** Brand identity, project detection, startup info.

**Functions:**
1. `renderBanner()` — skillgrid logo (ASCII/unicode art) + tokyonight theme
2. `detectProject()` — walk up from cwd, find `docs/skillgrid/config.yaml`, extract project name
3. `showStatus()` — version, project, SDD phase, Mnemonic status, code index status, model

**Splash output:**
```
┌─────────────────────────────────────────────┐
│           ⬢ SKILLGRID PI                    │
│                                              │
│  project:  skillgrid                         │
│  sdd:      012-skill-pi-coding-agent (spec) │
│  memory:   ✓ connected (64 observations)    │
│  index:    ✓ fresh (1,247 files)            │
│  model:    claude-sonnet-4-5                 │
│                                              │
│  /sdd /dash /help                            │
└─────────────────────────────────────────────┘
```

**Brand:** Reuses skillgrid wordmark + Tokyo purple `#9854f1` from `plugins/opencode/skillgrid-logo.tsx` and `plugins/kilo/skillgrid-logo.tsx`. Do not invent a second brand.

### 01-mnemonic

**Purpose:** Connect Mnemonic MCP server via pi-mcp-adapter, register tools as native pi tools.

**Functions:**
1. `writeMcpConfig()` — write `~/skill-pi/mcp.json` with `skillgrid: { command: "skillgrid", args: ["mcp"] }`
2. `registerMnemonicTools()` — via pi-mcp-adapter, register: `mem_save`, `mem_search`, `mem_context`, `mem_session_start`, `mem_session_end`, `mem_session_summary`, `code_status`, `code_index`, `code_search`, `code_read`, `code_orient`, `code_grep`, `web_cache_lookup`, `web_cache_save`, `web_cache_search`
3. `failSoft()` — if `skillgrid` not on PATH or MCP spawn fails → warn, continue without memory tools
4. `status()` — return connection state for dashboard

**Fail-soft behavior:**
```
warn: skillgrid not found on PATH — running without Mnemonic memory tools
warn: MCP server "skillgrid" failed to connect — running without memory tools
```

### 02-sdd

**Purpose:** SDD workflow routing, phase subagents, Result Contracts.

**Functions:**
1. `registerCommands()` — register `/sdd`, `/sdd explore`, `/sdd propose`, `/sdd spec`, `/sdd apply`, `/sdd verify`, `/sdd archive`, `/sdd full`, `/sdd status`
2. `detectPhase()` — read `docs/skillgrid/changes/` + `tasks.md` State to determine current phase
3. `routePhase()` — dispatch to the correct subagent with the right system prompt
4. `runPhaseAgent()` — spawn subagent, capture Result Contract
5. `validateResultContract()` — check required fields: `status`, `executive_summary`, `artifacts`, `next_recommended`, `risks`
6. `configGate()` — read `config.yaml` constraints (TDD mode, test commands, verify gates) before each phase
7. `userGate()` — after spec phase, prompt user: Implement or Revise. No auto-apply.
8. `sddFullChain()` — chain: explore → propose → spec → [gate] → apply → verify → archive

**Result Contract (each phase returns):**
```typescript
interface ResultContract {
  status: "success" | "partial" | "blocked";
  executive_summary: string;
  artifacts: Array<{ path: string; description: string }>;
  next_recommended: string;
  risks: string[];
}
```

### 03-permissions

**Purpose:** Robust permission system. gotgenes allow/ask/deny gates + inobit ask dialog + `/plan`/`/build` modes.

**Functions:**
1. `toolGate()` — intercept `tool_call` events. Check tool against allow/ask/deny list.
2. `pathProtection()` — block writes to protected paths: `.env`, `.git/`, `node_modules/`, `~/skill-pi/config.json`, `docs/skillgrid/config.yaml`
3. `bashGate()` — confirm before destructive bash: `rm -rf`, `sudo`, `git push`, `git commit --amend`, `npm publish`
4. `planMode()` — `/plan` is read-only (writes denied); `/build` returns to normal. Session-scoped.
5. `loadConfig()` — read `~/skill-pi/permission.json` (global) + `.pi/permission.json` (project). Project overrides global.
6. `askDialog()` — `ctx.ui.confirm()` with y/s/n/r (approve-once / approve-session / deny / deny-with-reason)

**Permission config:**
```json
{
  "tools": {
    "bash": "ask",
    "write": "allow",
    "edit": "allow",
    "read": "allow"
  },
  "protected_paths": [
    ".env", ".env.*", ".git/", "node_modules/",
    "~/skill-pi/config.json", "docs/skillgrid/config.yaml"
  ],
  "bash_danger_patterns": [
    "rm -rf", "sudo ", "git push", "git commit --amend",
    "npm publish", "bun publish", "docker push"
  ],
  "read_only_phases": ["spec", "verify"]
}
```

**Risk levels:**
| Level | Trigger | Behavior |
|-------|---------|----------|
| `allow` | read, grep, find, ls | Allow silently |
| `ask` | bash with danger patterns, write to sensitive paths | Ask dialog (y/s/n/r) |
| `deny` | tool in deny list, write to protected paths | Block, show reason |
| `plan` | `/plan` mode active | All writes denied |

**Fail-closed:** Gate errors block the tool call. Never silently ungated.

### 04-dashboard

**Purpose:** Custom TUI dashboard, omegon-pi style. Live view of agent state.

**Functions:**
1. `registerCommand()` — `/dash` and `/dashboard` commands + `Ctrl+Shift+B` shortcut
2. `renderDashboard()` — full TUI component below editor (pi-tui custom component)
3. `renderSddSection()` — active change, current phase, task progress (x/y), next action
4. `renderMemorySection()` — Mnemonic status (connected/disconnected), observation count, recent saves (last 3)
5. `renderIndexSection()` — code index status (fresh/stale), file count, symbol count
6. `renderSessionSection()` — tokens used, cost, model, thinking level, session ID
7. `renderAgentSection()` — active subagents, status, progress

**Dashboard layout:**
```
┌─ DASHBOARD ─────────────────────────────────────────┐
│ SDD    012-skill-pi-coding-agent · spec · 3/5 tasks │
│        next: /sdd apply                              │
│─────────────────────────────────────────────────────│
│ MEMORY ✓ connected · 64 obs · last: 2m ago          │
│        • Proposed 012-skill-pi-coding-agent change  │
│        • Researched earendil-works/pi architecture  │
│─────────────────────────────────────────────────────│
│ INDEX  ✓ fresh · 1,247 files · 8,432 symbols        │
│─────────────────────────────────────────────────────│
│ AGENTS sdd-spec running · 2/3 artifacts             │
│─────────────────────────────────────────────────────│
│ SESSION claude-sonnet-4-5 · thinking: high          │
│        tokens: 45.2k in / 12.8k out · cost: $0.42  │
│        session: 3a9c816d                            │
└─────────────────────────────────────────────────────┘
```

**Reactive updates:** Uses `shared-state.ts` event bus. Extensions push state changes, dashboard re-renders without polling.

### 05-local-llm

**Purpose:** Auto-detect and register local LLM providers (Ollama, vLLM). Integrated into `/login` as "Local API" menu element. No JSON editing.

**Functions:**
1. `detectOllama()` — `curl localhost:11434/api/tags` → list models
2. `detectVllm()` — `curl localhost:8000/v1/models` → list models
3. `registerModels()` — register detected models in pi's model registry via `pi.registerProvider()`
4. `loginHandler()` — `/login` → Local API → wizard: backend type (Ollama/vLLM/generic OpenAI-compat), base URL, optional API key, connection test, model pick
5. `effortHandler()` — `/effort <name>` → switch model tier. `local` → cheapest local model. `cloud` → best cloud model.

**Auto-detection on startup:**
```
info: Ollama detected at localhost:11434 (3 models: llama3.2, qwen2.5, mistral)
info: vllm detected at localhost:8000 (1 model: qwen3-8b)
```

**/login flow:**
```
/login
  > Anthropic
  > OpenAI
  > Google
  > Local API          ← new entry
    > Ollama
      Detected 3 models at localhost:11434:
      > llama3.2 (8B)
        qwen2.5 (14B)
        mistral (7B)
      Base URL [localhost:11434]:
      API Key [optional]:
      Connection test: ✓ (12ms)
      Select model [enter=first]:
      ✓ Registered llama3.2 (Ollama)
    > vLLM
      Detected 1 model at localhost:8000:
      > qwen3-8b
      Base URL [localhost:8000]:
      API Key [optional]:
      Connection test: ✓ (8ms)
      ✓ Registered qwen3-8b (vLLM)
    > Generic OpenAI-compatible
      Base URL:
      API Key:
      Model name:
      Connection test: ✓
      ✓ Registered
```

**Constraints:**
- Localhost-only by default. Non-localhost requires explicit URL.
- Connection test must succeed before registering.
- Secrets use Pi auth storage. Never written into `~/skill-pi/` synced markdown.
- Reject non-http(s) schemes. No LAN sweep in v1.

## SDD Phase Subagents

Six subagents, one per SDD phase. Each is a markdown file with YAML frontmatter (agent definition) + system prompt. Synced to `~/skill-pi/agents/` via hash-manifest.

| Agent | Model | Tools | Purpose |
|-------|-------|-------|---------|
| `sdd-explore` | Haiku (cheap) | read, grep, find, ls, bash (read-only) | Research external APIs, rare docs → research.md |
| `sdd-propose` | Sonnet | read, grep, find, ls, write | Write change.md from template |
| `sdd-spec` | Sonnet | read, grep, find, ls, write | Write tasks.md + acceptance.feature |
| `sdd-apply` | Sonnet | read, write, edit, bash, grep, find | Implement tasks, mark [x] |
| `sdd-verify` | Sonnet | read, grep, find, ls, bash (test commands) | Run verification, per-step verdicts |
| `sdd-archive` | Haiku (cheap) | read, write, bash (git mv) | Move to archive/, update changelog |

**Orchestration:** Uses `pi-subagents` npm package for the subagent infrastructure (spawn, stream, collect). Each phase agent runs as a `pi --mode rpc` child process with its own system prompt and tool surface.

**Result Contract:** Every phase agent must return a structured Result Contract. The orchestrator validates it before proceeding to the next phase.

**User Gate:** After `sdd-spec`, the agent pauses and asks: "Implement or Revise?" No auto-apply.

**SDD Full Chain:**
```
/sdd full
  → sdd-explore (if research needed)
  → sdd-propose → change.md
  → sdd-spec → tasks.md + acceptance.feature
  → [USER GATE: Implement or Revise?]
  → sdd-apply → implement, mark [x]
  → sdd-verify → per-step verdicts
  → [if PASS/WARNINGS + QA accepted]
  → sdd-archive → move to archive/
```

## Bundled Plugins

Auto-installed by `npm install -g skill-pi` (listed in `dependencies`):

| Plugin | Purpose | Why |
|--------|---------|-----|
| `@earendil-works/pi-coding-agent` | pi CLI + SDK | Core runtime |
| `@earendil-works/pi-agent-core` | Agent loop, state management | Core runtime |
| `@earendil-works/pi-ai` | LLM provider abstraction | Core runtime |
| `@earendil-works/pi-tui` | Terminal UI | Core runtime |
| `pi-mcp-adapter` | MCP connection layer | Mnemonic integration |
| `pi-web-access` | Web fetch/search | Fills web gap |
| `pi-subagents` | Sub-agent infrastructure | SDD phase agents |
| `@juicesharp/rpiv-ask-user-question` | Structured user questions | QoL |
| `@juicesharp/rpiv-todo` | Live task overlay | QoL + todo tracking |
| `@juicesharp/rpiv-web-tools` | Web search/fetch | QoL + web |

## Asset Sync (Hash-Manifest Policy)

From 009. On session start, extensions sync managed markdown from the npm package to `~/skill-pi/`:

1. **First run:** Create `~/skill-pi/` with synced asset trees (prompts, skills, agents, chains)
2. **Re-run:** Update only hash-managed files. If a managed file was edited by the user (hash mismatch), skip it with a warning
3. **New files:** Add files that exist in the package but not in the home
4. **Removed files:** Do NOT delete files from home (user may have added their own)
5. **Fail-closed:** If home is unwritable, abort with clear error

**Hash manifest:** `~/skill-pi/.sync-manifest.json` maps each managed file to its content hash. On sync, compare current hash vs manifest vs package hash.

## Error Handling

| Failure | Behavior | Notes |
|---------|----------|-------|
| Pi dependency missing / wrong version at runtime | `abort` | Clear message to reinstall/update `skill-pi` |
| Agent-home sync cannot write `~/skill-pi/` | `abort` | Fail closed before claiming ready |
| Managed file dirty vs package hash | `warn+continue` | Skip that file; report in status/doctor |
| Mnemonic MCP unavailable | `warn+continue` | SDD assets still load; memory tools degraded with doctor hint |
| `skillgrid install` package install fails | `abort` | Non-zero; do not mark agent configured |
| Unknown agent key typo (`pi` instead of `skill-pi`) | `abort` | Suggest `skill-pi` |
| Local API connection test fails | `abort` (of that add) | Do not register provider; show reachable error |
| `/login` Local API hook unsupported on this Pi | `warn+continue` | Doctor explains; cloud `/login` entries still work |
| Disallowed host / scheme on Local API URL | `abort` (of that add) | Reject non-http(s) or non-allowed hosts |
| Required rpiv companion missing or unloadable | `warn+continue` | Doctor lists which of the three failed; core SDD still starts |
| Permission config missing | `warn+continue` | Seed safe defaults into `~/skill-pi/permission.json` |
| Permission gate internal error | `abort` (of that tool call) | Fail-closed: block the tool; surface gate_error |
| Operator denies in ask dialog | `abort` (of that tool call) | Model continues unless Esc hard-terminate |
| SDD phase agent returns invalid Result Contract | `warn+continue` | Log the contract, continue with executive_summary only |
| Subagent process crashes | `abort` (of that phase) | Report which phase failed; suggest re-run |

## Threat Matrix

| Boundary / threat | Applicable? | Owning step | Planned RED coverage |
|-------------------|-------------|-------------|----------------------|
| Accidental `pi` bin alias publish | Applicable | 01 | Package.json / pack listing asserts no `bin.pi` |
| Shell / subprocess install-update | Applicable | 01, 05 | Launcher/update and installer npm invocations fail closed on non-zero child exit |
| Agent-home path traversal / symlink escape | Applicable | 03 | Sync refuses to write outside resolved `~/skill-pi/` when asset paths contain `..` or symlink escape |
| Hash-manifest clobber of dirty user files | Applicable | 03 | Dirty managed file is skipped; clean hash match updates; status reports skipped paths |
| Documentation-like paths (executable markdown on sync) | Applicable | 03 | Sync refuses to treat unexpected executable bits under managed trees as runnable |
| Installer agent-key confusion (`pi` vs `skill-pi`) | Applicable | 05 | Unknown key `pi` rejected with message naming `skill-pi` |
| Local API URL SSRF / unexpected host probe | Applicable | 07 | Reject disallowed schemes/hosts; localhost-default; connection timeouts; no register on failed probe |
| Local API secret leakage into synced agent home | Applicable | 07 | Fixture asserts API keys never written under `~/skill-pi/` managed trees |
| Unpinned / undeclared rpiv companion | Applicable | 08 | Package identity test fails if any of the three deps is missing |
| Permission config loaded from wrong home (`~/.pi/agent`) | Applicable | 09 | Test asserts only `~/skill-pi/permission.json` is consulted |
| Plan mode allows write/edit | Applicable | 09 | Fixture: write/edit in `/plan` is denied |
| Ask path silently allows without dialog | Applicable | 09 | ask-policy tool call requires dialog decision before proceed |
| SDD phase agent returns invalid Result Contract | Applicable | 10 | Fixture: missing `status` or `artifacts` field → warn, continue with degraded contract |
| Subagent process crash mid-phase | Applicable | 10 | Fixture: child process exits non-zero → phase aborts, error reported |
| Dashboard renders stale state | Applicable | 11 | Fixture: state bus event triggers re-render within 1s |

## Testing Strategy

- **Unit:** `cd skill-pi && bun test` — each extension tested in isolation with mocked pi ExtensionAPI
- **Integration:** `bun test --integration` — Mnemonic MCP connection, SDD phase routing, permission gates, asset sync
- **E2E:** `bun test test/e2e.test.ts` — full SDD cycle in a temp repo: propose → spec → apply → verify → archive
- **Go:** `go test ./skillgrid-cli/internal/install/...` — installer agent key tests
- **Smoke:** `skill-pi --version`, `skill-pi -p "hello"`, `skill-pi /sdd status`

## Rollback

- `npm uninstall -g skill-pi` removes the CLI
- `skillgrid install --remove skill-pi` removes the installer agent
- `rm -rf ~/skill-pi/` removes all state
- No changes to skillgrid-cli (Go) beyond additive agent key
- Upstream Pi installs untouched

## Open Questions

1. Exact Pi `/login` extension hook API — confirm against pinned `@earendil-works/pi-coding-agent` during step 07 apply
2. Whether optional generic OpenAI-compatible URL is exposed in the same Local API wizard in v1 (default: include as third choice)
3. Exact pin versions for `@juicesharp/rpiv-*` packages — chosen at apply time against current npm + Pi peer compatibility
4. Whether permission implementation depends on `@gotgenes/pi-permission-system` with a config-path adapter, or ports the model into Skillgrid-owned extension code
5. Whether project-scoped overlay `.skill-pi/permission.json` is v1 or deferred (default: global `~/skill-pi/permission.json` only in v1)
6. Pi version floor to pin — need live release check at propose time

## References

- [earendil-works/pi](https://github.com/earendil-works/pi) — pi agent framework
- [gentle-pi](https://github.com/Gentleman-Programming/gentle-pi) — SDD pi package (phase agents, Result Contracts, native authority)
- [oh-my-pi](https://github.com/can1357/oh-my-pi) — full-featured pi fork (typed subagent yields, dashboard)
- [omegon-pi](https://github.com/styrene-lab/omegon-pi) — pi distribution (wrapper pattern, dashboard, mcp-bridge, local LLM)
- [pi.dev/packages](https://pi.dev/packages) — pi package catalog
- [pi-mcp-adapter](https://www.npmjs.com/package/pi-mcp-adapter) — MCP connection for pi
- [pi-subagents](https://github.com/nicobailon/pi-subagents) — sub-agent infrastructure
- [pi-web-access](https://www.npmjs.com/package/pi-web-access) — web access
- [pi-localllm-provider](https://github.com/freeyoung/pi-localllm-provider) — local LLM wizard pattern
- [Crossbar](https://github.com/Hypabolic/Crossbar) — local LLM discovery/register pattern
- [rpiv-mono](https://github.com/juicesharp/rpiv-mono) — rpiv companion packages
- [gotgenes/pi-permission-system](https://pi.dev/packages/@gotgenes/pi-permission-system) — permission gates
- [inobit/pi-permission](https://github.com/inobit/pi-packages) — ask dialog + plan/build UX
- [009-skillgrid-pi](../skillgrid-old/docs/skillgrid/changes/009-skillgrid-pi/) — previous plan (merged into this)
