# Source: docs/skillgrid/changes/012-skill-pi-coding-agent/change.md
# Template: .agents/skills/_shared/templates/template-acceptance.feature
# Trace: change.md ## Goal + ## Definition of Done; tasks.md @step-NN verify lines.
# Mapping: @p0 scenarios ↔ change.md DoD / Testing strategy; @p1 = important failure paths.
# One Feature per step; tag each Feature with @step-NN matching tasks.md.

@step-01
Feature: skill-pi package scaffold with pinned Pi deps and skill-pi bin only
  As an operator
  I want an installable npm package at repo root with bin `skill-pi` only
  So that I can install and run the Skillgrid Pi distribution without a `pi` alias

  @happy @p0
  Scenario: no pi alias
    Given a packaged `skill-pi/` directory with `package.json`
    When  the operator inspects the `bin` manifest
    Then  `bin` contains `skill-pi` and does NOT contain `pi`

  @happy
  Scenario: pinned pi dep
    Given a packaged `skill-pi/` directory with `package.json`
    When  the operator inspects `dependencies`
    Then  `@earendil-works/pi-coding-agent` is declared at a pinned (non-caret, non-star) version

  @happy
  Scenario: companion deps present
    Given a packaged `skill-pi/` directory with `package.json`
    When  the operator inspects `dependencies`
    Then  all six of `pi-mcp-adapter`, `pi-web-access`, `pi-subagents`, `@juicesharp/rpiv-ask-user-question`, `@juicesharp/rpiv-todo`, `@juicesharp/rpiv-web-tools` are declared

  @edge
  Scenario: bin exposed
    Given a locally installed `skill-pi` package
    When  the operator runs `skill-pi --help`
    Then  the `skill-pi` binary is on PATH and responds

  @failure @p1
  Scenario: pi dep missing aborts
    Given a runtime where `@earendil-works/pi-coding-agent` is missing or wrong version
    When  the operator starts `skill-pi`
    Then  the launcher aborts with a clear reinstall/update message

@step-02
Feature: Packaged extensions load on session start with status/doctor surface
  As an operator
  I want packaged extensions to load when skill-pi starts
  So that I can verify the distribution is healthy via a status/doctor command

  @happy @p0
  Scenario: extensions load on start
    Given a fresh `skill-pi` distribution with packaged extensions
    When  the operator starts a session
    Then  all packaged extensions are registered and loaded

  @happy
  Scenario: status fields
    Given a running `skill-pi` session
    When  the operator invokes the status/doctor surface
    Then  the output contains package name, version, extension count, Mnemonic status, and code index status

  @edge
  Scenario: doctor invocable
    Given a running `skill-pi` session
    When  the operator invokes the Skillgrid status/doctor command
    Then  the surface responds without crashing the session

@step-03
Feature: Hash-manifest asset sync into ~/skill-pi
  As an operator
  I want managed prompts/skills/agents/chains synced into ~/skill-pi
  So that my Skillgrid SDD assets are always current without clobbering my edits

  @happy @p0
  Scenario: fresh home synced
    Given a first-run `skill-pi` session with no existing `~/skill-pi/`
    When  asset sync runs on session start
    Then  `~/skill-pi/` is created with prompts/skills/agents/chains populated

  @happy
  Scenario: rerun only managed
    Given an existing `~/skill-pi/` with some user-edited managed files
    When  asset sync runs again
    Then  only hash-managed files are updated and user-edited files are not overwritten

  @edge
  Scenario: dirty file skipped
    Given a managed file under `~/skill-pi/` that the user has edited (hash mismatch)
    When  asset sync runs
    Then  the dirty file is skipped and the skipped path is reported in status

  @edge
  Scenario: path traversal refused
    Given an asset path containing `..` or a symlink escape
    When  asset sync attempts to write the file
    Then  the write is refused and nothing is written outside resolved `~/skill-pi/`

  @edge
  Scenario: executable bits ignored
    Given a markdown asset under a managed tree with the executable bit set
    When  asset sync runs
    Then  the file is synced as markdown and is not treated as runnable

  @failure @p1
  Scenario: unwritable home aborts
    Given a `~/skill-pi/` directory that is not writable
    When  asset sync runs
    Then  sync fails closed with a clear error before claiming ready

@step-04
Feature: SDD commands registered with Mnemonic MCP wiring via pi-mcp-adapter
  As an agent inside skill-pi
  I want SDD workflow commands and Mnemonic tools available
  So that I can run the full SDD pipeline with persistent memory

  @happy @p0
  Scenario: SDD commands registered
    Given a `skill-pi` session with synced SDD assets
    When  the agent lists available SDD commands
    Then  commands for onboard, propose, spec, apply, verify, and archive are registered

  @happy
  Scenario: mnemonic mcp connected
    Given a `skill-pi` session with `skillgrid mcp` available
    When  pi-mcp-adapter connects to Mnemonic
    Then  `mem_*`, `code_*`, and `web_*` tools are registered and callable

  @happy
  Scenario: config constraints applied
    Given a project with `docs/skillgrid/config.yaml` declaring TDD mode and test commands
    When  SDD extensions load
    Then  the constraints from config.yaml are read and applied to the session

  @edge
  Scenario: mnemonic unavailable warn+continue
    Given a `skill-pi` session where `skillgrid mcp` is not running
    When  SDD extensions attempt Mnemonic connection
    Then  a warning is logged and SDD assets still load; memory tools are degraded with a doctor hint

@step-05
Feature: Go installer treats skill-pi as a first-class agent option
  As an operator
  I want `skillgrid install --agents skill-pi` to install the package
  So that I can install skill-pi from the hub without manual npm steps

  @happy @p0
  Scenario: available agents includes skill-pi
    Given the Go installer `AvailableAgents()` function
    When  the operator calls the agent list
    Then  `skill-pi` is present in the available agents

  @happy
  Scenario: install dry-run succeeds
    Given a system with npm available
    When  the operator runs `skillgrid install --agents skill-pi` in dry-run mode
    Then  the install path resolves and reports success without registering a `pi` alias

  @edge
  Scenario: existing agent keys unchanged
    Given the Go installer `AvailableAgents()` function
    When  the operator calls the agent list
    Then  OpenCode, Kilo, and Cursor are still present alongside `skill-pi`

  @failure @p1
  Scenario: invalid pi key fails with hint
    Given the Go installer
    When  the operator runs `skillgrid install --agents pi`
    Then  the install aborts with a message suggesting `skill-pi`

@step-06
Feature: Skillgrid Tokyo Night theme and omegon-style startup splash
  As an operator
  I want a Skillgrid-branded look with Tokyo Night theme and splash
  So that skill-pi presents a consistent brand identity on startup

  @happy @p0
  Scenario: tokyo night theme declared
    Given a `skill-pi` package with `themes/tokyonight.json`
    When  the `pi` manifest is inspected
    Then  the Tokyo Night theme is declared and loads with the distribution

  @happy
  Scenario: splash shows skillgrid wordmark
    Given a `skill-pi` session starting
    When  the splash extension renders
    Then  the Skillgrid wordmark is displayed using Tokyo purple `#9854f1`

  @happy
  Scenario: brand art matches identity
    Given the existing Skillgrid logo in `plugins/opencode/skillgrid-logo.tsx` and `plugins/kilo/skillgrid-logo.tsx`
    When  the splash logo is compared
    Then  the brand art matches the existing OpenCode/Kilo Skillgrid logo identity

  @edge
  Scenario: operator disables splash
    Given a `skill-pi` session with splash disabled via documented config
    When  the agent starts
    Then  no splash is shown and SDD functionality is unaffected

@step-07
Feature: Local API entry under /login for Ollama and vLLM setup
  As an operator
  I want to add local Ollama or vLLM from the /login menu
  So that I can use local models without hand-editing config files

  @happy @p0
  Scenario: login shows local api
    Given a `skill-pi` session
    When  the operator opens the `/login` menu
    Then  a Local API entry is visible

  @happy
  Scenario: ollama wizard adds provider
    Given a running Ollama instance on localhost:11434
    When  the operator runs the Local API wizard selecting Ollama and the connection test succeeds
    Then  a Pi provider is registered and models appear in the model picker

  @happy
  Scenario: vllm wizard adds provider
    Given a running vLLM instance on localhost:8000
    When  the operator runs the Local API wizard selecting vLLM and the connection test succeeds
    Then  a Pi provider is registered and models appear in the model picker

  @happy
  Scenario: auto-detect localhost
    Given Ollama running on localhost:11434 and vLLM running on localhost:8000
    When  the Local API wizard starts
    Then  both endpoints are auto-detected and suggested

  @edge
  Scenario: failed connection test aborts
    Given no Ollama instance running on localhost:11434
    When  the operator runs the Local API wizard and the connection test fails
    Then  no provider is registered and a reachable error is shown

  @edge
  Scenario: secrets in pi auth storage
    Given a successful Local API add with an API key
    When  the operator inspects `~/skill-pi/` synced markdown
    Then  the API key is stored in Pi auth storage and is NOT written into synced markdown

  @failure @p1
  Scenario: disallowed host rejected
    Given a Local API wizard with a non-localhost URL
    When  the operator enters a URL with a disallowed host or non-http(s) scheme
    Then  the add is aborted and no provider is registered

@step-08
Feature: Pinned rpiv companions and plugins load with skill-pi
  As an operator
  I want rpiv companions and plugins bundled and activated
  So that I get ask/todo/web-tools/MCP/subagent infrastructure without manual install

  @happy @p0
  Scenario: six deps pinned and loadable
    Given a `skill-pi` package with all six companion/plugin deps pinned
    When  the operator starts a session
    Then  all six deps (`pi-mcp-adapter`, `pi-web-access`, `pi-subagents`, `@juicesharp/rpiv-ask-user-question`, `@juicesharp/rpiv-todo`, `@juicesharp/rpiv-web-tools`) are loaded

  @happy
  Scenario: doctor reports loaded
    Given a `skill-pi` session with all companions present
    When  the operator invokes the doctor surface
    Then  the doctor reports all companions as loaded

  @edge
  Scenario: missing companion warn+continue
    Given a `skill-pi` session where one rpiv companion is missing or unloadable
    When  the operator starts the session
    Then  the doctor lists which companion failed and core SDD still starts

@step-09
Feature: Skillgrid permission extension with ask dialog and plan mode
  As an operator
  I want permission gates with a Skillgrid-owned config path and plan/build modes
  So that I can control tool access and enforce read-only phases

  @happy @p0
  Scenario: config at skill-pi home
    Given a `skill-pi` session
    When  the permission extension loads policy
    Then  the config is read from `~/skill-pi/permission.json` only

  @happy
  Scenario: allow ask deny semantics
    Given a `skill-pi` session with a permission policy
    When  the agent calls a tool governed by the policy
    Then  allow permits, ask triggers the dialog, and deny blocks

  @happy
  Scenario: ask dialog y s n r
    Given an ask-policy tool call
    When  the operator sees the ask dialog
    Then  the dialog offers approve-once (y), approve-session (s), deny (n), and deny-with-reason (r)

  @happy
  Scenario: plan mode read-only
    Given a `skill-pi` session in `/plan` mode
    When  the agent attempts a write or edit
    Then  the write is denied

  @happy
  Scenario: build returns to normal
    Given a `skill-pi` session in `/plan` mode
    When  the operator invokes `/build`
    Then  the session returns to normal permission mode

  @edge
  Scenario: read-only phases deny writes
    Given a `skill-pi` session in a read-only SDD phase (spec or verify)
    When  the agent attempts a write
    Then  the write is automatically denied

  @failure @p1
  Scenario: gate error fails closed
    Given a `skill-pi` session where the permission gate throws an internal error
    When  the agent calls a governed tool
    Then  the tool call is blocked (fail-closed) and a gate_error is surfaced

@step-10
Feature: Six SDD phase sub-agents with Result Contracts and sdd-full chain
  As an agent inside skill-pi
  I want SDD phases to run as sub-agents returning structured Result Contracts
  So that the full SDD pipeline runs with validated outputs and a user gate

  @happy @p0
  Scenario: sdd routes to correct phase
    Given a `skill-pi` session with SDD assets synced
    When  the agent invokes `/sdd`
    Then  the current phase is detected and the next action is suggested

  @happy
  Scenario: sdd full chain runs
    Given a `skill-pi` session in a fresh SDD change
    When  the agent invokes `/sdd full`
    Then  the chain runs explore → propose → spec → [USER GATE] → apply → verify → archive

  @happy
  Scenario: result contract valid
    Given a SDD phase sub-agent completing its phase
    When  the orchestrator receives the return
    Then  the Result Contract contains `status`, `executive_summary`, `artifacts`, `next_recommended`, and `risks`

  @happy
  Scenario: user gate after spec
    Given a `sdd-full` chain that has completed the spec phase
    When  the chain reaches the user gate
    Then  the agent presents "Implement or Revise?" and does not auto-apply

  @happy
  Scenario: per-phase model routing
    Given a `sdd-full` chain running
    When  the explore and archive phases dispatch
    Then  they use the Haiku model while other phases use Sonnet

  @failure @p1
  Scenario: invalid result contract warn+continue
    Given a SDD phase sub-agent returning a contract missing `status` or `artifacts`
    When  the orchestrator validates the contract
    Then  a warning is logged and the chain continues with `executive_summary` only

  @failure @p1
  Scenario: subagent crash aborts phase
    Given a SDD phase sub-agent whose child process exits non-zero
    When  the orchestrator detects the crash
    Then  the phase is aborted, the error is reported, and a re-run is suggested

@step-11
Feature: Custom TUI dashboard showing SDD status, memory, code index, session, and agents
  As an operator
  I want a live dashboard with SDD status, memory, code index, session, and agent info
  So that I can monitor the agent at a glance without leaving the TUI

  @happy @p0
  Scenario: dash command opens dashboard
    Given a `skill-pi` session
    When  the operator invokes `/dash`
    Then  the dashboard panel opens below the editor

  @happy
  Scenario: ctrl shift b opens dashboard
    Given a `skill-pi` session
    When  the operator presses Ctrl+Shift+B
    Then  the dashboard panel opens below the editor

  @happy
  Scenario: five sections rendered
    Given an open dashboard
    When  the operator inspects the panel
    Then  SDD status, memory, code index, session, and agents sections are all rendered

  @happy
  Scenario: reactive updates within 1s
    Given an open dashboard subscribed to the shared-state event bus
    When  a state change event is published
    Then  the dashboard re-renders within 1 second

  @edge
  Scenario: state bus unavailable
    Given a `skill-pi` session where the shared-state event bus is unavailable
    When  the operator opens the dashboard
    Then  the dashboard shows "state unavailable" without blocking the agent

@step-12
Feature: Documentation and rollout covering install, update, and coexistence
  As an operator
  I want docs covering npm install, skillgrid install, update, and coexistence
  So that I can discover and use skill-pi without tribal knowledge

  @happy @p0
  Scenario: docs cover install paths
    Given the `docs/user-manual/01-installation.md` file
    When  the operator reads the installation docs
    Then  both `npm install -g skill-pi` and `skillgrid install --agents skill-pi` paths are documented

  @happy
  Scenario: docs cover local api
    Given the user manual
    When  the operator reads the Local API section
    Then  the `/login` → Local API flow for Ollama/vLLM is described

  @happy
  Scenario: docs cover sdd subagents
    Given the user manual
    When  the operator reads the SDD subagent section
    Then  `/sdd` commands, Result Contracts, `sdd-full` chain, and user gate are described

  @happy
  Scenario: docs cover dashboard
    Given the user manual
    When  the operator reads the dashboard section
    Then  `/dash`, `Ctrl+Shift+B`, and all five sections are described

  @happy
  Scenario: docs cover permission
    Given the user manual
    When  the operator reads the permission section
    Then  config at `~/skill-pi/permission.json`, ask dialog, and `/plan`/`/build` are described

  @happy
  Scenario: docs cover coexistence
    Given the user manual
    When  the operator reads the coexistence section
    Then  coexistence with standalone `pi` is documented (no alias)

  @happy
  Scenario: plan harness lists include skill-pi
    Given `docs/plan/01-workflow-new.md`, `docs/plan/02-future.md`, and `docs/plan/03-skill-pi.md`
    When  the maintainer reviews the plan harness lists
    Then  `skill-pi` is included as a harness target

  @edge
  Scenario: agents md mentions skill-pi
    Given the `AGENTS.md` file at repo root
    When  the operator reads the agent list
    Then  `skill-pi` is mentioned in the agent list

  @failure @p1
  Scenario: docs mention bundled companions
    Given the user manual
    When  the operator reads the bundled components section
    Then  all six companions (rpiv ask/todo/web-tools + pi-mcp-adapter + pi-web-access + pi-subagents) are listed
