# Tasks: 012-skill-pi-coding-agent

> **STATUS:** `in-progress` — 0/12 steps PASS
>
> **For agentic workers:** REQUIRED SUB-SKILL: use subagent-execution (or simple-execution) to implement step-by-step. Steps use checkbox (`- [ ]`) syntax.

**Goal:** Ship a fat `skill-pi` CLI distribution at repo root that vendors Pi, delivers Skillgrid SDD with 6 phase sub-agents + Result Contracts, Mnemonic wiring via pi-mcp-adapter, a Tokyo Night Skillgrid look with omegon-style splash, `/login`-integrated local Ollama/vLLM setup, pinned rpiv companions, a Skillgrid permission extension, and a custom TUI dashboard.

**Architecture:** Repo-root npm package `skill-pi/` depends on pinned `@earendil-works/pi-coding-agent`, exposes bin `skill-pi` only (no `pi` alias), packages TS extensions (splash, mnemonic, sdd, permissions, dashboard, local-llm), ships a Tokyo Night theme + logo, pins and activates three `@juicesharp/rpiv-*` companions + `pi-mcp-adapter` + `pi-web-access` + `pi-subagents`, hash-syncs prompts/skills/agents/chains into `~/skill-pi`, and extends the Go installer with an additive `skill-pi` agent. SDD phase sub-agents run as `pi --mode rpc` children returning validated Result Contracts. See `change.md` Architecture decisions.

**Tech Stack:** Node/TypeScript (Bun runtime), npm-publishable package at repo root `skill-pi/`; Go (`skillgrid-cli/internal/install`); Mnemonic MCP via pi-mcp-adapter; Pi themes JSON + splash; local Ollama/vLLM via `/login` Local API; permission policy at `~/skill-pi/permission.json`.

**Spec:** `docs/skillgrid/changes/012-skill-pi-coding-agent/change.md`

**Acceptance:** `docs/skillgrid/changes/012-skill-pi-coding-agent/acceptance.feature` (`@step-NN`)

---

## Goal

Operators can install and run `skill-pi` (npm global and/or `skillgrid install`) and get a Pi-based coding agent that loads Skillgrid SDD workflow assets and Mnemonic wiring, presents a Skillgrid Tokyo Night splash/theme, runs SDD phases as sub-agents with Result Contracts, adds local Ollama/vLLM from the `/login` menu without hand-editing config, ships rpiv ask/todo/web-tools companions, enforces Skillgrid permissions (ask dialog + read-only `/plan`), and shows a live dashboard with SDD status, memory, code index, and session info.

## Out of scope / Non-Goals

- Thin launcher that only `exec`s a separately installed upstream `pi` (rejected depth)
- Shipping a `pi` PATH alias (rejected CLI identity)
- Sibling-repo publish for v1 (monorepo `skill-pi/` only for now)
- Gentle RDD / receipt-driven review / `gentle-ai` binary provisioning
- Full Gentle/omegon ambiance/QoL theme catalog beyond the single Skillgrid Tokyo Night look
- Forking or depending on oh-my-pi / `omp` as the runtime
- Replacing OpenCode/Kilo/Cursor install paths
- Shipping a separate slash command (`/localllm`, `/crossbar`, `/local`) as the primary UX — local setup lives under `/login` → Local API
- Full Crossbar multi-backend matrix in v1 (LM Studio, llama-swap, Anthropic cloud, LAN sweep) — Ollama + vLLM (+ optional generic OpenAI-compatible URL) only
- Enabling LAN/mDNS discovery by default (localhost-only unless operator opts in later)
- Installing the full `@juicesharp/rpiv-pi` pipeline / `/rpiv-setup` — only the three named companions
- Forking or rebranding rpiv-mono; Skillgrid pins published npm packages and activates them
- Using gotgenes/inobit default config paths under `~/.pi/agent/...` — Skillgrid permission config is `~/skill-pi/permission.json`
- Shipping a full OS sandbox / container isolation layer (permission is an in-process gate)
- Depending on `@gotgenes/pi-subagents` / model-judge kitchen sink as required v1 surface
- Sub-agent orchestration beyond the 6 SDD phase agents (no generic `task` tool, no Agent Hub TUI in v1)
- Web UI / remote dashboard (TUI-only dashboard in v1)

## Definition of Done

Change is done only when **all** of the following are true:

- [ ] Every success criterion / DoD checkbox in `change.md` is met
- [ ] Every `@step-NN` Feature in `acceptance.feature` has passing `@happy`, `@edge`, and `@failure` scenarios
- [ ] Every step below has Verdict `PASS` or `PASS WITH WARNINGS`
- [ ] No unchecked `- [ ]` under any `### Tasks`
- [ ] No **Global Constraint** violated
- [ ] Applicable threat-matrix rows have RED coverage that passed
- [ ] Testing strategy commands in `change.md` are green
- [ ] Rollback path in `change.md` is still valid (or N/A documented)
- [ ] Change archived under `docs/skillgrid/archive/012-skill-pi-coding-agent/`
- [ ] `## State` status is `done` (set at archive gate)

## Global Constraints

Copy verbatim from `change.md` (Error handling + Non-Goals + stack rules). Every step inherits these — do not restate per step.

- **Runtime/stack:** Node/TypeScript on the Bun runtime; npm-publishable package at repo-root `skill-pi/`; Go for `skillgrid-cli` install/setup.
- **CLI identity:** Package name and bin are `skill-pi` only; NEVER register a `pi` binary alias. Installer agent key/display is `skill-pi` (not `pi`).
- **Agent home:** Agent home for synced markdown is `~/skill-pi` (not `~/.pi/agent`, not `~/.skillgrid-pi`). Extensions ship inside the npm package; prompts/skills/agents/chains sync to `~/skill-pi`.
- **Fail-closed sync:** Agent-home sync cannot write `~/skill-pi/` → `abort` (fail closed before claiming ready). Managed file dirty vs package hash → `warn+continue` (skip that file; report in status/doctor). Path traversal / symlink escape under `~/skill-pi` is refused.
- **Mnemonic is fail-soft:** Missing/unavailable Mnemonic MCP or HTTP → `warn+continue`; SDD assets still load, memory tools degrade with a doctor hint. Never block SDD on Mnemonic.
- **Permission config location:** Policy is read from and written to `~/skill-pi/permission.json` ONLY — never `~/.pi/agent/...`. Missing config → `warn+continue` (seed safe defaults). Gate internal error → `abort` that tool call (fail-closed, surface `gate_error`).
- **Pi dependency:** Pi dependency missing / wrong version at runtime → `abort` with a clear reinstall/update message. Pin exact Pi versions.
- **Local API:** `/login` Local API hook unsupported on this Pi → `warn+continue` (doctor explains; cloud `/login` still works). Disallowed host/scheme → `abort` that add. Failed connection test → `abort` that add (do not register). API keys use Pi auth storage; NEVER write secrets into synced markdown under `~/skill-pi`.
- **rpiv companions:** Required rpiv companion missing or unloadable → `warn+continue` (doctor lists which of the three failed; core SDD still starts). Only the three named companions; no `/rpiv-setup` dependency.
- **SDD sub-agents:** Invalid Result Contract → `warn+continue` (log contract, continue with `executive_summary` only; do not block the chain). Subagent process crash → `abort` that phase (report which phase failed; suggest re-run; do not corrupt session state). User gate after spec is mandatory: no auto-apply.
- **Installer:** `skillgrid install` package install fails → `abort` (non-zero; do not mark agent configured). Unknown agent key typo (`pi`) → `abort` suggesting `skill-pi`. Existing agent keys must remain available (additive only).
- **Dashboard:** Dashboard state bus unavailable → `warn+continue` (show "state unavailable"; do not block agent).
- **Brand:** Reuse existing Skillgrid wordmark / Tokyo purple `#9854f1` from `plugins/opencode/skillgrid-logo.tsx` and `plugins/kilo/skillgrid-logo.tsx`; single Skillgrid Tokyo Night look.
- **No Windows support** in this change; no marketing site / gallery listing; TUI-only dashboard.

---

## State

```yaml
phase: spec          # spec | apply | verify | archive
current_step: 01-package-scaffold
status: in_progress  # in_progress | blocked | done
updated: 2026-09-09
```

## Step map

| NN | Step | Tag | Blocked by | Acceptance |
|----|------|-----|------------|------------|
| 01 | `package-scaffold` | `@step-01` | — | Feature tagged `@step-01` |
| 02 | `extensions-core` | `@step-02` | 01 | Feature tagged `@step-02` |
| 03 | `asset-sync` | `@step-03` | 02 | Feature tagged `@step-03` |
| 04 | `sdd-mnemonic` | `@step-04` | 03 | Feature tagged `@step-04` |
| 05 | `installer` | `@step-05` | 01 | Feature tagged `@step-05` |
| 06 | `branding-look` | `@step-06` | 02 | Feature tagged `@step-06` |
| 07 | `local-api-login` | `@step-07` | 02 | Feature tagged `@step-07` |
| 08 | `rpiv-companions` | `@step-08` | 01 | Feature tagged `@step-08` |
| 09 | `permission` | `@step-09` | 02 | Feature tagged `@step-09` |
| 10 | `sdd-subagents` | `@step-10` | 04 | Feature tagged `@step-10` |
| 11 | `dashboard` | `@step-11` | 04, 10 | Feature tagged `@step-11` |
| 12 | `docs-rollout` | `@step-12` | 04, 05, 06, 07, 08, 09, 10, 11 | Feature tagged `@step-12` |

## Review workload (change-level)

| Field | Value |
|-------|-------|
| Estimated changed lines (change) | ~2500 |
| 400-line budget risk | High |
| Chained PRs recommended | Yes |
| Delivery strategy | auto-chain |
<!-- single-pr = one PR for the change; each ## NN step still commits separately when DoD is met (see work-unit-commits / commits.md). -->

---

## 01-package-scaffold

### Goal

A runnable npm package skeleton exists at repo root `skill-pi/` with correct identity (name `skill-pi`, bin `skill-pi` only, no `pi` alias), pinned `@earendil-works/pi-coding-agent`, and the six required companion/plugin dependencies declared.

### Out of scope / Non-Goals

- Extension behavior (step 02 onward)
- Asset sync (step 03)
- Go installer agent option (step 05)

### Definition of Done

This step is done only when:

- [ ] All `### Tasks` checkboxes below are `[x]`
- [ ] All `@step-01` scenarios in `acceptance.feature` pass
- [ ] `### Verification` Verdict is `PASS` or `PASS WITH WARNINGS`
- [ ] Produces contracts listed under Interfaces are available to dependents
- [ ] No Global Constraint violated

> Depends on: none

**Files:**
- Create: `skill-pi/package.json`
- Create: `skill-pi/bin/skill-pi.mjs`
- Create: `skill-pi/tsconfig.json`
- Test: `skill-pi/test/scaffold.test.ts`

**Interfaces:**
- Consumes: none
- Produces: `skill-pi` npm package identity (name/bin), pinned Pi + companion deps, launcher entry — relied on by 02, 05, 08.

### Tasks

- [ ] 01.1 `[RED]` Accidental `pi` bin alias publish (THREAT) — pack/manifest asserts no `bin.pi` key and no `pi` alias
  - [ ] 01.1.a Write failing test: package `bin` must contain `skill-pi` and must NOT contain `pi`
  - [ ] 01.1.b Run to confirm fail — `Run: cd skill-pi && bun test test/scaffold.test.ts -t "no pi alias"` — Expected: FAIL
  - [ ] 01.1.c Minimal implementation: `package.json` `bin: { "skill-pi": "./bin/skill-pi.mjs" }` only
  - [ ] 01.1.d Run to confirm pass — `Run: cd skill-pi && bun test test/scaffold.test.ts -t "no pi alias"` — Expected: PASS
  - [ ] 01.1.e Commit — `feat(skill-pi): scaffold package with skill-pi bin only (no pi alias)`
- [ ] 01.2 `[RED]` Package declares pinned `@earendil-works/pi-coding-agent` dependency (WHAT: pinned Pi dep)
  - [ ] 01.2.a Write failing test: dependencies has `@earendil-works/pi-coding-agent` at a pinned (non-caret, non-star) version
  - [ ] 01.2.b Run to confirm fail — `Run: cd skill-pi && bun test test/scaffold.test.ts -t "pinned pi dep"` — Expected: FAIL
  - [ ] 01.2.c Minimal implementation: pin exact `@earendil-works/pi-coding-agent` version
  - [ ] 01.2.d Run to confirm pass — `Run: cd skill-pi && bun test test/scaffold.test.ts -t "pinned pi dep"` — Expected: PASS
  - [ ] 01.2.e Commit — `feat(skill-pi): pin pi-coding-agent dependency`
- [ ] 01.3 `[RED]` Unpinned / undeclared companion or plugin (THREAT: six deps present) — package identity fails if any of the six is missing
  - [ ] 01.3.a Write failing test: `pi-mcp-adapter`, `pi-web-access`, `pi-subagents`, `@juicesharp/rpiv-ask-user-question`, `@juicesharp/rpiv-todo`, `@juicesharp/rpiv-web-tools` all declared
  - [ ] 01.3.b Run to confirm fail — `Run: cd skill-pi && bun test test/scaffold.test.ts -t "companion deps present"` — Expected: FAIL
  - [ ] 01.3.c Minimal implementation: declare all six deps (pinned)
  - [ ] 01.3.d Run to confirm pass — `Run: cd skill-pi && bun test test/scaffold.test.ts -t "companion deps present"` — Expected: PASS
  - [ ] 01.3.e Commit — `feat(skill-pi): declare companion and plugin dependencies`
- [ ] 01.4 `[AFK]` Operator can npm install the local package and see bin `skill-pi` (WHAT) — `Run: cd skill-pi && bun test test/scaffold.test.ts -t "bin exposed"` — Expected: PASS
- [ ] 01.5 `[AFK]` Pi dependency missing or wrong version at runtime aborts with clear message (THREAT: shell/subprocess fail-closed on non-zero child) — `Run: cd skill-pi && bun test test/scaffold.test.ts -t "pi dep missing aborts"` — Expected: PASS

### Verification

Verdict: `PENDING`  <!-- PASS | PASS WITH WARNINGS | FAIL -->

Evidence:

| Check | Run | Expected | Result | Notes |
|-------|-----|----------|--------|-------|
| Focused test | `cd skill-pi && bun test test/scaffold.test.ts` | PASS | | |
| Acceptance `@step-01` / `@p0` | `cd skill-pi && bun test --integration -t "@step-01"` | PASS | | |
| Runtime harness | `cd skill-pi && bun x npm pack --dry-run` | PASS | | bin `skill-pi` only |
| Rollback boundary | remove `skill-pi/` package; no `pi` alias leaked | PASS | | |
| Global Constraints | — | held | | |

### Commit

When step DoD is met: `feat(skill-pi): scaffold repo-root package with pinned Pi deps and skill-pi bin only`

---

## 02-extensions-core

### Goal

Packaged extensions load when `skill-pi` starts, and a Skillgrid status/doctor surface reports package identity and extension load.

### Out of scope / Non-Goals

- Full SDD asset tree (step 04)
- Go installer (step 05)

### Definition of Done

This step is done only when:

- [ ] All `### Tasks` checkboxes below are `[x]`
- [ ] All `@step-02` scenarios in `acceptance.feature` pass
- [ ] `### Verification` Verdict is `PASS` or `PASS WITH WARNINGS`
- [ ] Depends-on step 01 already PASS / PASS WITH WARNINGS
- [ ] No Global Constraint violated

> Depends on: 01-package-scaffold

**Files:**
- Create: `skill-pi/extensions/*.ts`
- Test: `skill-pi/test/extensions.test.ts`

**Interfaces:**
- Consumes: `skill-pi` package identity + bin from 01
- Produces: session-start extension loader, status/doctor surface, sync trigger hook — relied on by 03, 06, 07, 09, 11.

### Tasks

- [ ] 02.1 `[RED]` Session start loads packaged extensions from the distribution (WHAT: extensions load on start)
  - [ ] 02.1.a Write failing test: starting a session registers the packaged extension set
  - [ ] 02.1.b Run to confirm fail — `Run: cd skill-pi && bun test test/extensions.test.ts -t "extensions load on start"` — Expected: FAIL
  - [ ] 02.1.c Minimal implementation: register extensions from package pi resources
  - [ ] 02.1.d Run to confirm pass — `Run: cd skill-pi && bun test test/extensions.test.ts -t "extensions load on start"` — Expected: PASS
  - [ ] 02.1.e Commit — `feat(skill-pi): load packaged extensions on session start`
- [ ] 02.2 `[RED]` Status reports package name, version, extension count, Mnemonic status, code index status (WHAT: status fields)
  - [ ] 02.2.a Write failing test: doctor output contains all five required fields
  - [ ] 02.2.b Run to confirm fail — `Run: cd skill-pi && bun test test/extensions.test.ts -t "status fields"` — Expected: FAIL
  - [ ] 02.2.c Minimal implementation: assemble status fields from package metadata + probe state
  - [ ] 02.2.d Run to confirm pass — `Run: cd skill-pi && bun test test/extensions.test.ts -t "status fields"` — Expected: PASS
  - [ ] 02.2.e Commit — `feat(skill-pi): status/doctor surface reports identity and extension load`
- [ ] 02.3 `[AFK]` Operator can invoke the Skillgrid status/doctor surface from the running agent (WHAT) — `Run: cd skill-pi && bun test test/extensions.test.ts -t "doctor invocable"` — Expected: PASS

### Verification

Verdict: `PENDING`

Evidence:

| Check | Run | Expected | Result | Notes |
|-------|-----|----------|--------|-------|
| Focused test | `cd skill-pi && bun test test/extensions.test.ts` | PASS | | |
| Acceptance `@step-02` / `@p0` | `cd skill-pi && bun test --integration -t "@step-02"` | PASS | | |
| Runtime harness | `cd skill-pi && bun bin/skill-pi.mjs --doctor` | PASS | | status surface |
| Rollback boundary | extensions optional; start still works without them | PASS | | |
| Global Constraints | — | held | | |

### Commit

When step DoD is met: `feat(skill-pi): packaged extensions load with status/doctor surface`

---

## 03-asset-sync

### Goal

Managed prompts/skills/agents/chains appear under `~/skill-pi/` using a hash-manifest policy: fresh home gets synced trees, dirty managed files are skipped with a warning, unwritable home fails closed, and path traversal is refused.

### Out of scope / Non-Goals

- Semantic content completeness of every Skillgrid skill (mapping quality is step 04)

### Definition of Done

This step is done only when:

- [ ] All `### Tasks` checkboxes below are `[x]`
- [ ] All `@step-03` scenarios in `acceptance.feature` pass
- [ ] `### Verification` Verdict is `PASS` or `PASS WITH WARNINGS`
- [ ] Depends-on step 02 already PASS / PASS WITH WARNINGS
- [ ] No Global Constraint violated

> Depends on: 02-extensions-core

**Files:**
- Create: `skill-pi/lib/agent-home.ts`
- Create: `skill-pi/assets/{prompts,skills,agents,chains}/**`
- Test: `skill-pi/test/asset-sync.test.ts`

**Interfaces:**
- Consumes: sync trigger hook from 02
- Produces: resolved `~/skill-pi/` home + hash-manifest sync of managed trees — relied on by 04 (SDD asset mapping), 09 (permission home), 10 (agent defs sync).

### Tasks

- [ ] 03.1 `[RED]` Agent-home path traversal / symlink escape (THREAT) — sync refuses to write outside resolved `~/skill-pi/` when asset paths contain `..` or a symlink escape
  - [ ] 03.1.a Write failing test: asset path with `..`/symlink escape is refused; nothing written outside home
  - [ ] 03.1.b Run to confirm fail — `Run: cd skill-pi && bun test test/asset-sync.test.ts -t "path traversal refused"` — Expected: FAIL
  - [ ] 03.1.c Minimal implementation: resolve + compare realpaths against home root
  - [ ] 03.1.d Run to confirm pass — `Run: cd skill-pi && bun test test/asset-sync.test.ts -t "path traversal refused"` — Expected: PASS
  - [ ] 03.1.e Commit — `feat(skill-pi): refuse path traversal and symlink escape in agent-home sync`
- [ ] 03.2 `[RED]` Hash-manifest clobber of dirty user files (THREAT) — dirty managed file is skipped, clean hash match updates, status reports skipped paths
  - [ ] 03.2.a Write failing test: dirty managed file skipped + reported; clean hash match updates
  - [ ] 03.2.b Run to confirm fail — `Run: cd skill-pi && bun test test/asset-sync.test.ts -t "dirty file skipped"` — Expected: FAIL
  - [ ] 03.2.c Minimal implementation: hash manifest compare; skip dirty, record skipped
  - [ ] 03.2.d Run to confirm pass — `Run: cd skill-pi && bun test test/asset-sync.test.ts -t "dirty file skipped"` — Expected: PASS
  - [ ] 03.2.e Commit — `feat(skill-pi): hash-manifest policy skips dirty managed files`
- [ ] 03.3 `[RED]` Documentation-like paths / executable markdown on sync (THREAT) — sync refuses to treat unexpected executable bits under managed trees as runnable
  - [ ] 03.3.a Write failing test: executable bit under managed tree does not make a markdown asset runnable
  - [ ] 03.3.b Run to confirm fail — `Run: cd skill-pi && bun test test/asset-sync.test.ts -t "executable bits ignored"` — Expected: FAIL
  - [ ] 03.3.c Minimal implementation: ignore exec bit on synced markdown
  - [ ] 03.3.d Run to confirm pass — `Run: cd skill-pi && bun test test/asset-sync.test.ts -t "executable bits ignored"` — Expected: PASS
  - [ ] 03.3.e Commit — `feat(skill-pi): ignore executable bits on synced markdown assets`
- [ ] 03.4 `[RED]` First run creates `~/skill-pi/` with synced asset trees (WHAT: fresh home)
  - [ ] 03.4.a Write failing test: fresh home gets prompts/skills/agents/chains populated
  - [ ] 03.4.b Run to confirm fail — `Run: cd skill-pi && bun test test/asset-sync.test.ts -t "fresh home synced"` — Expected: FAIL
  - [ ] 03.4.c Minimal implementation: create home + write managed trees on first run
  - [ ] 03.4.d Run to confirm pass — `Run: cd skill-pi && bun test test/asset-sync.test.ts -t "fresh home synced"` — Expected: PASS
  - [ ] 03.4.e Commit — `feat(skill-pi): first-run sync populates ~/skill-pi trees`
- [ ] 03.5 `[AFK]` Re-run updates only hash-managed files; user-edited managed files not overwritten (WHAT) — `Run: cd skill-pi && bun test test/asset-sync.test.ts -t "rerun only managed"` — Expected: PASS
- [ ] 03.6 `[AFK]` Sync failure fails closed with clear error when home unwritable (WHAT / Error handling: abort) — `Run: cd skill-pi && bun test test/asset-sync.test.ts -t "unwritable home aborts"` — Expected: PASS

### Verification

Verdict: `PENDING`

Evidence:

| Check | Run | Expected | Result | Notes |
|-------|-----|----------|--------|-------|
| Focused test | `cd skill-pi && bun test test/asset-sync.test.ts` | PASS | | |
| Acceptance `@step-03` / `@p0` | `cd skill-pi && bun test --integration -t "@step-03"` | PASS | | |
| Runtime harness | `cd skill-pi && bun test test/asset-sync.test.ts -t "fresh home synced"` | PASS | | |
| Rollback boundary | delete `~/skill-pi`; next start re-syncs cleanly | PASS | | |
| Global Constraints | — | held | | |

### Commit

When step DoD is met: `feat(skill-pi): hash-manifest asset sync into ~/skill-pi`

---

## 04-sdd-mnemonic

### Goal

Skillgrid SDD pipeline is usable inside skill-pi; Mnemonic wiring works via pi-mcp-adapter with fail-soft degradation.

### Out of scope / Non-Goals

- Subagent orchestration (step 10)
- Dashboard (step 11)

### Definition of Done

This step is done only when:

- [ ] All `### Tasks` checkboxes below are `[x]`
- [ ] All `@step-04` scenarios in `acceptance.feature` pass
- [ ] `### Verification` Verdict is `PASS` or `PASS WITH WARNINGS`
- [ ] Depends-on step 03 already PASS / PASS WITH WARNINGS
- [ ] No Global Constraint violated

> Depends on: 03-asset-sync

**Files:**
- Create: `skill-pi/extensions/01-mnemonic/`
- Create: `skill-pi/extensions/02-sdd/` (base)
- Test: `skill-pi/test/sdd-mnemonic.test.ts`

**Interfaces:**
- Consumes: sync trigger hook from 02; resolved `~/skill-pi/` home from 03
- Produces: SDD command registration, Mnemonic MCP tool surface, config constraint loading — relied on by 10 (subagents), 11 (dashboard), 12 (docs).

### Tasks

- [ ] 04.1 `[RED]` Synced agents/skills cover SDD routing (WHAT: SDD commands registered) — SDD extension registers onboard → propose → spec → apply ⇄ verify → archive commands
  - [ ] 04.1.a Write failing test: starting a session with synced SDD assets registers commands for onboard, propose, spec, apply, verify, and archive
  - [ ] 04.1.b Run to confirm fail — `Run: cd skill-pi && bun test test/sdd-mnemonic.test.ts -t "SDD commands registered"` — Expected: FAIL
  - [ ] 04.1.c Minimal implementation: register SDD slash commands from synced asset manifest
  - [ ] 04.1.d Run to confirm pass — `Run: cd skill-pi && bun test test/sdd-mnemonic.test.ts -t "SDD commands registered"` — Expected: PASS
  - [ ] 04.1.e Commit — `feat(skill-pi): register SDD slash commands from synced assets`
- [ ] 04.2 `[RED]` pi-mcp-adapter connects to skillgrid mcp and registers tools (WHAT: Mnemonic MCP connected) — `mem_*`, `code_*`, `web_*` tools are callable when available
  - [ ] 04.2.a Write failing test: with `skillgrid mcp` running, pi-mcp-adapter registers `mem_*`, `code_*`, and `web_*` tools
  - [ ] 04.2.b Run to confirm fail — `Run: cd skill-pi && bun test test/sdd-mnemonic.test.ts -t "mnemonic mcp connected"` — Expected: FAIL
  - [ ] 04.2.c Minimal implementation: wire pi-mcp-adapter to `skillgrid mcp` and register tool surface
  - [ ] 04.2.d Run to confirm pass — `Run: cd skill-pi && bun test test/sdd-mnemonic.test.ts -t "mnemonic mcp connected"` — Expected: PASS
  - [ ] 04.2.e Commit — `feat(skill-pi): connect Mnemonic MCP via pi-mcp-adapter`
- [ ] 04.3 `[RED]` Missing Mnemonic does not prevent SDD asset use (WHAT: fail-soft warn+continue) — SDD assets load and memory tools degrade with doctor hint
  - [ ] 04.3.a Write failing test: with `skillgrid mcp` unavailable, SDD commands still register and doctor reports Mnemonic degraded
  - [ ] 04.3.b Run to confirm fail — `Run: cd skill-pi && bun test test/sdd-mnemonic.test.ts -t "mnemonic unavailable warn+continue"` — Expected: FAIL
  - [ ] 04.3.c Minimal implementation: catch MCP connection failure; log warning; mark memory tools degraded in doctor
  - [ ] 04.3.d Run to confirm pass — `Run: cd skill-pi && bun test test/sdd-mnemonic.test.ts -t "mnemonic unavailable warn+continue"` — Expected: PASS
  - [ ] 04.3.e Commit — `feat(skill-pi): Mnemonic fail-soft degradation with doctor hint`
- [ ] 04.4 `[RED]` config.yaml constraints are read and applied (WHAT: TDD mode, test commands) — SDD extension reads `docs/skillgrid/config.yaml` and applies constraints
  - [ ] 04.4.a Write failing test: with a project config.yaml declaring TDD mode, the SDD session enforces it
  - [ ] 04.4.b Run to confirm fail — `Run: cd skill-pi && bun test test/sdd-mnemonic.test.ts -t "config constraints applied"` — Expected: FAIL
  - [ ] 04.4.c Minimal implementation: discover and parse `docs/skillgrid/config.yaml`; apply TDD mode and test commands to session
  - [ ] 04.4.d Run to confirm pass — `Run: cd skill-pi && bun test test/sdd-mnemonic.test.ts -t "config constraints applied"` — Expected: PASS
  - [ ] 04.4.e Commit — `feat(skill-pi): read and apply config.yaml constraints`

### Verification

Verdict: `PENDING`  <!-- PASS | PASS WITH WARNINGS | FAIL -->

Evidence:

| Check | Run | Expected | Result | Notes |
|-------|-----|----------|--------|-------|
| Focused test | `cd skill-pi && bun test test/sdd-mnemonic.test.ts` | PASS | | |
| Acceptance `@step-04` / `@p0` | `cd skill-pi && bun test --integration -t "@step-04"` | PASS | | |
| Runtime harness | `cd skill-pi && bun bin/skill-pi.mjs --doctor` | PASS | | SDD + Mnemonic status |
| Rollback boundary | remove `01-mnemonic/` + `02-sdd/`; SDD degrades but agent starts | PASS | | |
| Global Constraints | — | held | | Mnemonic fail-soft |

### Commit

When step DoD is met: `feat(skill-pi): SDD commands + Mnemonic MCP wiring with fail-soft`

---

## 05-installer

### Goal

Hub installer treats `skill-pi` as a first-class agent option alongside OpenCode/Kilo/Cursor.

### Out of scope / Non-Goals

- Changing OpenCode/Kilo/Cursor behavior beyond additive list
- npm package publishing pipeline (step 01)

### Definition of Done

This step is done only when:

- [ ] All `### Tasks` checkboxes below are `[x]`
- [ ] All `@step-05` scenarios in `acceptance.feature` pass
- [ ] `### Verification` Verdict is `PASS` or `PASS WITH WARNINGS`
- [ ] Depends-on step 01 already PASS / PASS WITH WARNINGS
- [ ] No Global Constraint violated

> Depends on: 01-package-scaffold

**Files:**
- Modify: `skillgrid-cli/internal/install/config.go`
- Modify: `skillgrid-cli/internal/install/*.go`
- Create: `skillgrid-cli/internal/install/*_test.go`

**Interfaces:**
- Consumes: `skill-pi` npm package identity from 01
- Produces: `skill-pi` entry in `AvailableAgents()`, install/link path — relied on by 12 (docs).

### Tasks

- [ ] 05.1 `[RED]` Installer agent-key confusion (`pi` vs `skill-pi`) (THREAT) — unknown key `pi` rejected with message naming `skill-pi`
  - [ ] 05.1.a Write failing test: `skillgrid install --agents pi` fails with a message suggesting `skill-pi`
  - [ ] 05.1.b Run to confirm fail — `Run: cd skillgrid-cli && go test ./internal/install/... -run "TestInvalidPiKey"` — Expected: FAIL
  - [ ] 05.1.c Minimal implementation: validate agent key against `AvailableAgents()`; return hint for `pi` → `skill-pi`
  - [ ] 05.1.d Run to confirm pass — `Run: cd skillgrid-cli && go test ./internal/install/... -run "TestInvalidPiKey"` — Expected: PASS
  - [ ] 05.1.e Commit — `feat(installer): reject pi agent key with skill-pi hint`
- [ ] 05.2 `[RED]` Shell/subprocess install-update (THREAT) — installer npm invocations fail closed on non-zero child exit
  - [ ] 05.2.a Write failing test: npm install subprocess returns non-zero → installer aborts, does not mark agent configured
  - [ ] 05.2.b Run to confirm fail — `Run: cd skillgrid-cli && go test ./internal/install/... -run "TestInstallFailsClosed"` — Expected: FAIL
  - [ ] 05.2.c Minimal implementation: check subprocess exit code; abort with clear error on non-zero
  - [ ] 05.2.d Run to confirm pass — `Run: cd skillgrid-cli && go test ./internal/install/... -run "TestInstallFailsClosed"` — Expected: PASS
  - [ ] 05.2.e Commit — `feat(installer): fail closed on npm subprocess non-zero exit`
- [ ] 05.3 `[RED]` AvailableAgents includes skill-pi (WHAT: additive agent key) — `skill-pi` present in the agent list
  - [ ] 05.3.a Write failing test: `AvailableAgents()` returns a list containing `skill-pi`
  - [ ] 05.3.b Run to confirm fail — `Run: cd skillgrid-cli && go test ./internal/install/... -run "TestAvailableAgents"` — Expected: FAIL
  - [ ] 05.3.c Minimal implementation: add `skill-pi` to `AvailableAgents()` in config.go
  - [ ] 05.3.d Run to confirm pass — `Run: cd skillgrid-cli && go test ./internal/install/... -run "TestAvailableAgents"` — Expected: PASS
  - [ ] 05.3.e Commit — `feat(installer): add skill-pi to AvailableAgents`
- [ ] 05.4 `[AFK]` Existing agent keys remain available (WHAT: additive only) — `Run: cd skillgrid-cli && go test ./internal/install/... -run "TestExistingAgents"` — Expected: PASS
- [ ] 05.5 `[AFK]` skillgrid install --agents skill-pi succeeds in dry-run (WHAT) — `Run: cd skillgrid-cli && go test ./internal/install/... -run "TestInstallDryRun"` — Expected: PASS

### Verification

Verdict: `PENDING`  <!-- PASS | PASS WITH WARNINGS | FAIL -->

Evidence:

| Check | Run | Expected | Result | Notes |
|-------|-----|----------|--------|-------|
| Focused test | `cd skillgrid-cli && go test ./internal/install/...` | PASS | | |
| Acceptance `@step-05` / `@p0` | `cd skillgrid-cli && go test ./internal/install/... -run "TestStep05"` | PASS | | |
| Runtime harness | `cd skillgrid-cli && go run . install --agents skill-pi --dry-run` | PASS | | |
| Rollback boundary | remove `skill-pi` from `AvailableAgents()`; existing agents still work | PASS | | |
| Global Constraints | — | held | | additive only |

### Commit

When step DoD is met: `feat(installer): add skill-pi agent with pi-key hint and fail-closed install`

---

## 06-branding-look

### Goal

skill-pi presents a Skillgrid look: Tokyo Night theme and startup splash/logo like omegon-pi, reusing existing brand identity.

### Out of scope / Non-Goals

- Dashboard (step 11)
- Ambiance catalog, inventing a second brand palette
- Full Gentle/omegon QoL theme set

### Definition of Done

This step is done only when:

- [ ] All `### Tasks` checkboxes below are `[x]`
- [ ] All `@step-06` scenarios in `acceptance.feature` pass
- [ ] `### Verification` Verdict is `PASS` or `PASS WITH WARNINGS`
- [ ] Depends-on step 02 already PASS / PASS WITH WARNINGS
- [ ] No Global Constraint violated

> Depends on: 02-extensions-core

**Files:**
- Create: `skill-pi/themes/tokyonight.json`
- Create: `skill-pi/extensions/00-splash/`
- Test: `skill-pi/test/branding.test.ts`

**Interfaces:**
- Consumes: session-start extension loader from 02
- Produces: packaged theme + splash registration — relied on by 12 (docs).

### Tasks

- [ ] 06.1 `[RED]` Tokyo Night theme declared in pi manifest (WHAT: theme declared + loads) — `themes/tokyonight.json` is registered and loads with the distribution
  - [ ] 06.1.a Write failing test: the `pi` manifest includes `tokyonight` theme and it loads on session start
  - [ ] 06.1.b Run to confirm fail — `Run: cd skill-pi && bun test test/branding.test.ts -t "tokyo night theme declared"` — Expected: FAIL
  - [ ] 06.1.c Minimal implementation: create `themes/tokyonight.json`; register in package `pi` manifest
  - [ ] 06.1.d Run to confirm pass — `Run: cd skill-pi && bun test test/branding.test.ts -t "tokyo night theme declared"` — Expected: PASS
  - [ ] 06.1.e Commit — `feat(skill-pi): package Tokyo Night theme`
- [ ] 06.2 `[RED]` Startup splash shows Skillgrid wordmark with #9854f1 (WHAT: splash wordmark) — splash extension renders the Skillgrid logo using Tokyo purple
  - [ ] 06.2.a Write failing test: splash extension renders wordmark with `#9854f1` on session start
  - [ ] 06.2.b Run to confirm fail — `Run: cd skill-pi && bun test test/branding.test.ts -t "splash shows skillgrid wordmark"` — Expected: FAIL
  - [ ] 06.2.c Minimal implementation: create `00-splash/` extension; render wordmark art with `#9854f1`
  - [ ] 06.2.d Run to confirm pass — `Run: cd skill-pi && bun test test/branding.test.ts -t "splash shows skillgrid wordmark"` — Expected: PASS
  - [ ] 06.2.e Commit — `feat(skill-pi): startup splash with Skillgrid wordmark`
- [ ] 06.3 `[AFK]` Operator can disable splash via documented config (WHAT) — `Run: cd skill-pi && bun test test/branding.test.ts -t "operator disables splash"` — Expected: PASS
- [ ] 06.4 `[AFK]` Brand art matches existing OpenCode/Kilo Skillgrid logo identity (WHAT) — `Run: cd skill-pi && bun test test/branding.test.ts -t "brand art matches identity"` — Expected: PASS

### Verification

Verdict: `PENDING`  <!-- PASS | PASS WITH WARNINGS | FAIL -->

Evidence:

| Check | Run | Expected | Result | Notes |
|-------|-----|----------|--------|-------|
| Focused test | `cd skill-pi && bun test test/branding.test.ts` | PASS | | |
| Acceptance `@step-06` / `@p0` | `cd skill-pi && bun test --integration -t "@step-06"` | PASS | | |
| Runtime harness | `cd skill-pi && bun bin/skill-pi.mjs --splash-test` | PASS | | theme + splash render |
| Rollback boundary | remove `themes/` + `00-splash/`; agent starts with default theme | PASS | | |
| Global Constraints | — | held | | brand reuse |

### Commit

When step DoD is met: `feat(skill-pi): Tokyo Night theme + Skillgrid splash branding`

---

## 07-local-api-login

### Goal

Operators add local Ollama or vLLM from `/login` → Local API without editing config by hand.

### Out of scope / Non-Goals

- Separate `/localllm` or `/crossbar` commands
- Full Crossbar backend matrix (LM Studio, llama-swap, Anthropic cloud, LAN sweep)
- Default LAN/mDNS discovery

### Definition of Done

This step is done only when:

- [ ] All `### Tasks` checkboxes below are `[x]`
- [ ] All `@step-07` scenarios in `acceptance.feature` pass
- [ ] `### Verification` Verdict is `PASS` or `PASS WITH WARNINGS`
- [ ] Depends-on step 02 already PASS / PASS WITH WARNINGS
- [ ] No Global Constraint violated

> Depends on: 02-extensions-core

**Files:**
- Create: `skill-pi/extensions/05-local-llm/`
- Create: `skill-pi/lib/local-backends/`
- Test: `skill-pi/test/local-api.test.ts`

**Interfaces:**
- Consumes: session-start extension loader from 02; Pi `/login` provider-setup API
- Produces: Local API menu entry, Ollama/vLLM detect/probe/register helpers — relied on by 12 (docs).

### Tasks

- [ ] 07.1 `[RED]` Local API URL SSRF / unexpected host probe (THREAT) — reject disallowed schemes/hosts; localhost-default; no register on failed probe
  - [ ] 07.1.a Write failing test: non-http(s) scheme and disallowed host are rejected; no provider registered
  - [ ] 07.1.b Run to confirm fail — `Run: cd skill-pi && bun test test/local-api.test.ts -t "disallowed host rejected"` — Expected: FAIL
  - [ ] 07.1.c Minimal implementation: validate URL scheme (http/https only) and host (localhost or explicit non-localhost URL); reject others
  - [ ] 07.1.d Run to confirm pass — `Run: cd skill-pi && bun test test/local-api.test.ts -t "disallowed host rejected"` — Expected: PASS
  - [ ] 07.1.e Commit — `feat(skill-pi): validate Local API URL scheme and host`
- [ ] 07.2 `[RED]` Local API secret leakage into synced agent home (THREAT) — API keys never written under `~/skill-pi/` managed trees
  - [ ] 07.2.a Write failing test: after a successful Local API add with an API key, no key value appears in `~/skill-pi/` synced markdown
  - [ ] 07.2.b Run to confirm fail — `Run: cd skill-pi && bun test test/local-api.test.ts -t "secrets in pi auth storage"` — Expected: FAIL
  - [ ] 07.2.c Minimal implementation: store API keys via Pi auth storage; never write to synced markdown paths
  - [ ] 07.2.d Run to confirm pass — `Run: cd skill-pi && bun test test/local-api.test.ts -t "secrets in pi auth storage"` — Expected: PASS
  - [ ] 07.2.e Commit — `feat(skill-pi): store Local API secrets in Pi auth storage`
- [ ] 07.3 `[RED]` /login menu includes Local API entry (WHAT: menu entry) — Local API is visible in the `/login` menu
  - [ ] 07.3.a Write failing test: opening `/login` shows a Local API menu item
  - [ ] 07.3.b Run to confirm fail — `Run: cd skill-pi && bun test test/local-api.test.ts -t "login shows local api"` — Expected: FAIL
  - [ ] 07.3.c Minimal implementation: register Local API menu item in Pi `/login` extension hook
  - [ ] 07.3.d Run to confirm pass — `Run: cd skill-pi && bun test test/local-api.test.ts -t "login shows local api"` — Expected: PASS
  - [ ] 07.3.e Commit — `feat(skill-pi): add Local API entry to /login menu`
- [ ] 07.4 `[RED]` Failed connection test does not register a provider (WHAT: abort on fail) — failed probe → no registration
  - [ ] 07.4.a Write failing test: with no Ollama on localhost:11434, wizard connection test fails and no provider is registered
  - [ ] 07.4.b Run to confirm fail — `Run: cd skill-pi && bun test test/local-api.test.ts -t "failed connection test aborts"` — Expected: FAIL
  - [ ] 07.4.c Minimal implementation: run connection test; on failure abort the add and show reachable error
  - [ ] 07.4.d Run to confirm pass — `Run: cd skill-pi && bun test test/local-api.test.ts -t "failed connection test aborts"` — Expected: PASS
  - [ ] 07.4.e Commit — `feat(skill-pi): abort Local API add on failed connection test`
- [ ] 07.5 `[AFK]` Wizard supports Ollama and vLLM (+ optional generic OpenAI-compatible) (WHAT) — `Run: cd skill-pi && bun test test/local-api.test.ts -t "ollama wizard adds provider"` — Expected: PASS
- [ ] 07.6 `[AFK]` Auto-detection probes localhost:11434 and localhost:8000 (WHAT) — `Run: cd skill-pi && bun test test/local-api.test.ts -t "auto-detect localhost"` — Expected: PASS
- [ ] 07.7 `[AFK]` Successful add registers Pi provider/models without hand-edited settings.json (WHAT) — `Run: cd skill-pi && bun test test/local-api.test.ts -t "vllm wizard adds provider"` — Expected: PASS

### Verification

Verdict: `PENDING`  <!-- PASS | PASS WITH WARNINGS | FAIL -->

Evidence:

| Check | Run | Expected | Result | Notes |
|-------|-----|----------|--------|-------|
| Focused test | `cd skill-pi && bun test test/local-api.test.ts` | PASS | | |
| Acceptance `@step-07` / `@p0` | `cd skill-pi && bun test --integration -t "@step-07"` | PASS | | |
| Runtime harness | `cd skill-pi && bun bin/skill-pi.mjs --local-api-test` | PASS | | wizard + register |
| Rollback boundary | remove `05-local-llm/` + `lib/local-backends/`; cloud /login still works | PASS | | |
| Global Constraints | — | held | | secrets not synced |

### Commit

When step DoD is met: `feat(skill-pi): /login Local API wizard for Ollama and vLLM`

---

## 08-rpiv-companions

### Goal

skill-pi ships and activates rpiv companions + pi-mcp-adapter + pi-web-access + pi-subagents without a separate manual install.

### Out of scope / Non-Goals

- Full `@juicesharp/rpiv-pi` pipeline, `/rpiv-setup`, advisor/workflow/voice/warp/btw/telemetry
- Forking or rebranding rpiv-mono

### Definition of Done

This step is done only when:

- [ ] All `### Tasks` checkboxes below are `[x]`
- [ ] All `@step-08` scenarios in `acceptance.feature` pass
- [ ] `### Verification` Verdict is `PASS` or `PASS WITH WARNINGS`
- [ ] Depends-on step 01 already PASS / PASS WITH WARNINGS
- [ ] No Global Constraint violated

> Depends on: 01-package-scaffold

**Files:**
- Modify: `skill-pi/package.json` (rpiv/pi-mcp/pi-web/pi-subagents deps)
- Test: `skill-pi/test/rpiv.test.ts`

**Interfaces:**
- Consumes: `skill-pi` npm package identity + bin from 01
- Produces: activated companion/plugin tool surface — relied on by 04 (Mnemonic), 10 (subagents), 12 (docs).

### Tasks

- [ ] 08.1 `[RED]` Unpinned/undeclared rpiv companion or plugin (THREAT) — package identity fails if any of the six deps is missing
  - [ ] 08.1.a Write failing test: all six deps (`pi-mcp-adapter`, `pi-web-access`, `pi-subagents`, `@juicesharp/rpiv-ask-user-question`, `@juicesharp/rpiv-todo`, `@juicesharp/rpiv-web-tools`) are declared at pinned versions
  - [ ] 08.1.b Run to confirm fail — `Run: cd skill-pi && bun test test/rpiv.test.ts -t "six deps pinned and loadable"` — Expected: FAIL
  - [ ] 08.1.c Minimal implementation: pin exact versions for all six deps in package.json; declare in `pi` resources
  - [ ] 08.1.d Run to confirm pass — `Run: cd skill-pi && bun test test/rpiv.test.ts -t "six deps pinned and loadable"` — Expected: PASS
  - [ ] 08.1.e Commit — `feat(skill-pi): pin and activate rpiv companions + plugins`
- [ ] 08.2 `[AFK]` Doctor reports companions loaded (WHAT) — `Run: cd skill-pi && bun test test/rpiv.test.ts -t "doctor reports loaded"` — Expected: PASS
- [ ] 08.3 `[AFK]` Missing companion reported by doctor without crashing (WHAT: warn+continue) — `Run: cd skill-pi && bun test test/rpiv.test.ts -t "missing companion warn+continue"` — Expected: PASS

### Verification

Verdict: `PENDING`  <!-- PASS | PASS WITH WARNINGS | FAIL -->

Evidence:

| Check | Run | Expected | Result | Notes |
|-------|-----|----------|--------|-------|
| Focused test | `cd skill-pi && bun test test/rpiv.test.ts` | PASS | | |
| Acceptance `@step-08` / `@p0` | `cd skill-pi && bun test --integration -t "@step-08"` | PASS | | |
| Runtime harness | `cd skill-pi && bun bin/skill-pi.mjs --doctor` | PASS | | companion status |
| Rollback boundary | remove one companion dep; doctor warns, agent starts | PASS | | |
| Global Constraints | — | held | | three named + three plugins only |

### Commit

When step DoD is met: `feat(skill-pi): pinned rpiv companions and plugins activate on load`

---

## 09-permission

### Goal

skill-pi enforces permissions with Skillgrid-owned config path (`~/skill-pi/permission.json`) and inobit-like ask/plan UX.

### Out of scope / Non-Goals

- OS sandbox / container isolation
- gotgenes default `~/.pi/agent/...` config ownership
- Project-scoped overlay `.skill-pi/permission.json` (deferred)

### Definition of Done

This step is done only when:

- [ ] All `### Tasks` checkboxes below are `[x]`
- [ ] All `@step-09` scenarios in `acceptance.feature` pass
- [ ] `### Verification` Verdict is `PASS` or `PASS WITH WARNINGS`
- [ ] Depends-on step 02 already PASS / PASS WITH WARNINGS
- [ ] No Global Constraint violated

> Depends on: 02-extensions-core

**Files:**
- Create: `skill-pi/extensions/03-permissions/`
- Create: `skill-pi/lib/permission/`
- Test: `skill-pi/test/permission.test.ts`

**Interfaces:**
- Consumes: session-start extension loader from 02; resolved `~/skill-pi/` home from 03
- Produces: permission gate, plan/build mode, ask dialog — relied on by 10 (SDD phases), 11 (dashboard), 12 (docs).

### Tasks

- [ ] 09.1 `[RED]` Permission config loaded from wrong home (`~/.pi/agent`) (THREAT) — test asserts only `~/skill-pi/permission.json` is consulted
  - [ ] 09.1.a Write failing test: permission extension reads policy from `~/skill-pi/permission.json` and does NOT consult `~/.pi/agent/permission.json`
  - [ ] 09.1.b Run to confirm fail — `Run: cd skill-pi && bun test test/permission.test.ts -t "config at skill-pi home"` — Expected: FAIL
  - [ ] 09.1.c Minimal implementation: hard-code `~/skill-pi/permission.json` as the only config path
  - [ ] 09.1.d Run to confirm pass — `Run: cd skill-pi && bun test test/permission.test.ts -t "config at skill-pi home"` — Expected: PASS
  - [ ] 09.1.e Commit — `feat(skill-pi): permission config at ~/skill-pi/permission.json only`
- [ ] 09.2 `[RED]` Plan mode allows write/edit (THREAT) — write/edit in `/plan` is denied
  - [ ] 09.2.a Write failing test: in `/plan` mode, a write or edit tool call is denied
  - [ ] 09.2.b Run to confirm fail — `Run: cd skill-pi && bun test test/permission.test.ts -t "plan mode read-only"` — Expected: FAIL
  - [ ] 09.2.c Minimal implementation: in plan mode, deny all write/edit tool calls
  - [ ] 09.2.d Run to confirm pass — `Run: cd skill-pi && bun test test/permission.test.ts -t "plan mode read-only"` — Expected: PASS
  - [ ] 09.2.e Commit — `feat(skill-pi): plan mode denies writes`
- [ ] 09.3 `[RED]` Ask path silently allows without dialog (THREAT) — ask-policy tool call requires dialog decision before proceeding
  - [ ] 09.3.a Write failing test: an ask-policy tool call blocks until the operator responds to the dialog
  - [ ] 09.3.b Run to confirm fail — `Run: cd skill-pi && bun test test/permission.test.ts -t "ask dialog y s n r"` — Expected: FAIL
  - [ ] 09.3.c Minimal implementation: on ask-policy, show dialog with y/s/n/r options; block until decision
  - [ ] 09.3.d Run to confirm pass — `Run: cd skill-pi && bun test test/permission.test.ts -t "ask dialog y s n r"` — Expected: PASS
  - [ ] 09.3.e Commit — `feat(skill-pi): ask dialog with y/s/n/r options`
- [ ] 09.4 `[RED]` Gate errors fail closed (WHAT: abort on gate_error) — internal gate error blocks the tool call
  - [ ] 09.4.a Write failing test: when the gate throws an internal error, the tool call is blocked and gate_error surfaced
  - [ ] 09.4.b Run to confirm fail — `Run: cd skill-pi && bun test test/permission.test.ts -t "gate error fails closed"` — Expected: FAIL
  - [ ] 09.4.c Minimal implementation: catch gate errors; fail closed (deny); surface `gate_error` message
  - [ ] 09.4.d Run to confirm pass — `Run: cd skill-pi && bun test test/permission.test.ts -t "gate error fails closed"` — Expected: PASS
  - [ ] 09.4.e Commit — `feat(skill-pi): gate errors fail closed`
- [ ] 09.5 `[AFK]` Allow/ask/deny semantics cover tools and bash (WHAT) — `Run: cd skill-pi && bun test test/permission.test.ts -t "allow ask deny semantics"` — Expected: PASS
- [ ] 09.6 `[AFK]` /build returns to normal mode; mode is session-scoped (WHAT) — `Run: cd skill-pi && bun test test/permission.test.ts -t "build returns to normal"` — Expected: PASS
- [ ] 09.7 `[AFK]` Read-only phases (spec, verify) automatically deny writes (WHAT) — `Run: cd skill-pi && bun test test/permission.test.ts -t "read-only phases deny writes"` — Expected: PASS

### Verification

Verdict: `PENDING`  <!-- PASS | PASS WITH WARNINGS | FAIL -->

Evidence:

| Check | Run | Expected | Result | Notes |
|-------|-----|----------|--------|-------|
| Focused test | `cd skill-pi && bun test test/permission.test.ts` | PASS | | |
| Acceptance `@step-09` / `@p0` | `cd skill-pi && bun test --integration -t "@step-09"` | PASS | | |
| Runtime harness | `cd skill-pi && bun bin/skill-pi.mjs --permission-test` | PASS | | gate + plan |
| Rollback boundary | remove `03-permissions/` + `lib/permission/`; agent starts with no gate | PASS | | |
| Global Constraints | — | held | | config path only |

### Commit

When step DoD is met: `feat(skill-pi): permission extension with ask dialog and plan/build modes`

---

## 10-sdd-subagents

### Goal

6 SDD phase sub-agents with Result Contracts, `/sdd` routing, `sdd-full` chain, and user gate after spec.

### Out of scope / Non-Goals

- Generic `task` tool, Agent Hub TUI, model-judge kitchen sink
- Sub-agent orchestration beyond the 6 SDD phase agents

### Definition of Done

This step is done only when:

- [ ] All `### Tasks` checkboxes below are `[x]`
- [ ] All `@step-10` scenarios in `acceptance.feature` pass
- [ ] `### Verification` Verdict is `PASS` or `PASS WITH WARNINGS`
- [ ] Depends-on step 04 already PASS / PASS WITH WARNINGS
- [ ] No Global Constraint violated

> Depends on: 04-sdd-mnemonic

**Files:**
- Create: `skill-pi/agents/sdd-explore.md`
- Create: `skill-pi/agents/sdd-propose.md`
- Create: `skill-pi/agents/sdd-spec.md`
- Create: `skill-pi/agents/sdd-apply.md`
- Create: `skill-pi/agents/sdd-verify.md`
- Create: `skill-pi/agents/sdd-archive.md`
- Modify: `skill-pi/extensions/02-sdd/` (subagents)
- Create: `skill-pi/lib/result-contract.ts`
- Test: `skill-pi/test/subagents.test.ts`

**Interfaces:**
- Consumes: SDD command registration from 04; Mnemonic tools from 04
- Produces: 6 phase agent definitions, Result Contract validation, sdd-full chain, user gate — relied on by 11 (dashboard), 12 (docs).

### Tasks

- [ ] 10.1 `[RED]` SDD phase agent returns invalid Result Contract (THREAT) — missing `status` or `artifacts` field → warn, continue with degraded contract
  - [ ] 10.1.a Write failing test: a phase agent returning a contract missing `status` triggers a warning and the chain continues with `executive_summary` only
  - [ ] 10.1.b Run to confirm fail — `Run: cd skill-pi && bun test test/subagents.test.ts -t "invalid result contract warn+continue"` — Expected: FAIL
  - [ ] 10.1.c Minimal implementation: validate Result Contract fields; on missing field warn and continue with degraded contract
  - [ ] 10.1.d Run to confirm pass — `Run: cd skill-pi && bun test test/subagents.test.ts -t "invalid result contract warn+continue"` — Expected: PASS
  - [ ] 10.1.e Commit — `feat(skill-pi): validate Result Contract with degraded fallback`
- [ ] 10.2 `[RED]` Subagent process crash mid-phase (THREAT) — child process exits non-zero → phase aborts, error reported
  - [ ] 10.2.a Write failing test: a phase sub-agent child process exiting non-zero aborts the phase and reports the error
  - [ ] 10.2.b Run to confirm fail — `Run: cd skill-pi && bun test test/subagents.test.ts -t "subagent crash aborts phase"` — Expected: FAIL
  - [ ] 10.2.c Minimal implementation: catch child process non-zero exit; abort phase; report which phase failed; suggest re-run
  - [ ] 10.2.d Run to confirm pass — `Run: cd skill-pi && bun test test/subagents.test.ts -t "subagent crash aborts phase"` — Expected: PASS
  - [ ] 10.2.e Commit — `feat(skill-pi): abort phase on subagent crash`
- [ ] 10.3 `[RED]` /sdd detects current phase and suggests next action (WHAT: routing) — `/sdd` inspects state and routes
  - [ ] 10.3.a Write failing test: invoking `/sdd` in a change directory detects the current phase and suggests the next action
  - [ ] 10.3.b Run to confirm fail — `Run: cd skill-pi && bun test test/subagents.test.ts -t "sdd routes to correct phase"` — Expected: FAIL
  - [ ] 10.3.c Minimal implementation: parse change state; detect phase; suggest next action
  - [ ] 10.3.d Run to confirm pass — `Run: cd skill-pi && bun test test/subagents.test.ts -t "sdd routes to correct phase"` — Expected: PASS
  - [ ] 10.3.e Commit — `feat(skill-pi): /sdd phase detection and routing`
- [ ] 10.4 `[RED]` /sdd full chain runs with user gate after spec (WHAT: sdd-full) — explore → propose → spec → [USER GATE] → apply → verify → archive
  - [ ] 10.4.a Write failing test: `/sdd full` runs the chain and pauses after spec for "Implement or Revise?"
  - [ ] 10.4.b Run to confirm fail — `Run: cd skill-pi && bun test test/subagents.test.ts -t "sdd full chain runs"` — Expected: FAIL
  - [ ] 10.4.c Minimal implementation: implement sdd-full chain with sequential phase dispatch and user gate after spec
  - [ ] 10.4.d Run to confirm pass — `Run: cd skill-pi && bun test test/subagents.test.ts -t "sdd full chain runs"` — Expected: PASS
  - [ ] 10.4.e Commit — `feat(skill-pi): sdd-full chain with user gate`
- [ ] 10.5 `[AFK]` Each phase agent returns valid Result Contract with all five fields (WHAT) — `Run: cd skill-pi && bun test test/subagents.test.ts -t "result contract valid"` — Expected: PASS
- [ ] 10.6 `[AFK]` Per-phase model routing: Haiku for explore/archive, Sonnet for others (WHAT) — `Run: cd skill-pi && bun test test/subagents.test.ts -t "per-phase model routing"` — Expected: PASS
- [ ] 10.7 `[AFK]` Six phase agent definitions exist with YAML frontmatter (WHAT) — `Run: cd skill-pi && bun test test/subagents.test.ts -t "six agent definitions"` — Expected: PASS

### Verification

Verdict: `PENDING`  <!-- PASS | PASS WITH WARNINGS | FAIL -->

Evidence:

| Check | Run | Expected | Result | Notes |
|-------|-----|----------|--------|-------|
| Focused test | `cd skill-pi && bun test test/subagents.test.ts` | PASS | | |
| Acceptance `@step-10` / `@p0` | `cd skill-pi && bun test --integration -t "@step-10"` | PASS | | |
| Runtime harness | `cd skill-pi && bun bin/skill-pi.mjs --sdd-test` | PASS | | phase dispatch |
| Rollback boundary | remove `agents/` + subagent dispatch; SDD commands still listed (no chain) | PASS | | |
| Global Constraints | — | held | | user gate mandatory |

### Commit

When step DoD is met: `feat(skill-pi): 6 SDD phase sub-agents with Result Contracts and sdd-full chain`

---

## 11-dashboard

### Goal

Custom TUI dashboard showing SDD status, memory, code index, session, and active agents with reactive updates.

### Out of scope / Non-Goals

- Web UI, remote dashboard
- Custom themes beyond tokyonight

### Definition of Done

This step is done only when:

- [ ] All `### Tasks` checkboxes below are `[x]`
- [ ] All `@step-11` scenarios in `acceptance.feature` pass
- [ ] `### Verification` Verdict is `PASS` or `PASS WITH WARNINGS`
- [ ] Depends-on steps 04 and 10 already PASS / PASS WITH WARNINGS
- [ ] No Global Constraint violated

> Depends on: 04-sdd-mnemonic, 10-sdd-subagents

**Files:**
- Create: `skill-pi/extensions/04-dashboard/`
- Test: `skill-pi/test/dashboard.test.ts`

**Interfaces:**
- Consumes: shared-state event bus from 02; SDD state from 10; Mnemonic state from 04
- Produces: dashboard TUI panel — relied on by 12 (docs).

### Tasks

- [ ] 11.1 `[RED]` Dashboard renders stale state (THREAT) — state bus event triggers re-render within 1s
  - [ ] 11.1.a Write failing test: publishing a state change event triggers a dashboard re-render within 1 second
  - [ ] 11.1.b Run to confirm fail — `Run: cd skill-pi && bun test test/dashboard.test.ts -t "reactive updates within 1s"` — Expected: FAIL
  - [ ] 11.1.c Minimal implementation: subscribe to shared-state event bus; on event, re-render dashboard within 1s
  - [ ] 11.1.d Run to confirm pass — `Run: cd skill-pi && bun test test/dashboard.test.ts -t "reactive updates within 1s"` — Expected: PASS
  - [ ] 11.1.e Commit — `feat(skill-pi): reactive dashboard re-render on state change`
- [ ] 11.2 `[RED]` /dash and /dashboard commands open the dashboard panel (WHAT: open command) — panel renders below the editor
  - [ ] 11.2.a Write failing test: invoking `/dash` or `/dashboard` opens the panel below the editor
  - [ ] 11.2.b Run to confirm fail — `Run: cd skill-pi && bun test test/dashboard.test.ts -t "dash command opens dashboard"` — Expected: FAIL
  - [ ] 11.2.c Minimal implementation: register `/dash` and `/dashboard` commands; render panel below editor
  - [ ] 11.2.d Run to confirm pass — `Run: cd skill-pi && bun test test/dashboard.test.ts -t "dash command opens dashboard"` — Expected: PASS
  - [ ] 11.2.e Commit — `feat(skill-pi): /dash command opens dashboard panel`
- [ ] 11.3 `[RED]` State bus unavailable shows "state unavailable" without blocking agent (WHAT: warn+continue) — dashboard degrades gracefully
  - [ ] 11.3.a Write failing test: with the shared-state event bus unavailable, the dashboard shows "state unavailable" and the agent continues
  - [ ] 11.3.b Run to confirm fail — `Run: cd skill-pi && bun test test/dashboard.test.ts -t "state bus unavailable"` — Expected: FAIL
  - [ ] 11.3.c Minimal implementation: catch bus unavailability; render "state unavailable"; do not block agent
  - [ ] 11.3.d Run to confirm pass — `Run: cd skill-pi && bun test test/dashboard.test.ts -t "state bus unavailable"` — Expected: PASS
  - [ ] 11.3.e Commit — `feat(skill-pi): dashboard degrades gracefully when state bus unavailable`
- [ ] 11.4 `[AFK]` Five sections rendered: SDD status, memory, code index, session, agents (WHAT) — `Run: cd skill-pi && bun test test/dashboard.test.ts -t "five sections rendered"` — Expected: PASS
- [ ] 11.5 `[AFK]` Ctrl+Shift+B shortcut opens the dashboard (WHAT) — `Run: cd skill-pi && bun test test/dashboard.test.ts -t "ctrl shift b opens dashboard"` — Expected: PASS

### Verification

Verdict: `PENDING`  <!-- PASS | PASS WITH WARNINGS | FAIL -->

Evidence:

| Check | Run | Expected | Result | Notes |
|-------|-----|----------|--------|-------|
| Focused test | `cd skill-pi && bun test test/dashboard.test.ts` | PASS | | |
| Acceptance `@step-11` / `@p0` | `cd skill-pi && bun test --integration -t "@step-11"` | PASS | | |
| Runtime harness | `cd skill-pi && bun bin/skill-pi.mjs --dashboard-test` | PASS | | render + sections |
| Rollback boundary | remove `04-dashboard/`; agent starts without dashboard | PASS | | |
| Global Constraints | — | held | | TUI-only |

### Commit

When step DoD is met: `feat(skill-pi): custom TUI dashboard with reactive state updates`

---

## 12-docs-rollout

### Goal

Operators can discover install, update, and coexistence from docs; plan harness lists include skill-pi.

### Out of scope / Non-Goals

- Marketing site / gallery listing
- Full API reference or tutorial series

### Definition of Done

This step is done only when:

- [ ] All `### Tasks` checkboxes below are `[x]`
- [ ] All `@step-12` scenarios in `acceptance.feature` pass
- [ ] `### Verification` Verdict is `PASS` or `PASS WITH WARNINGS`
- [ ] Depends-on steps 04, 05, 06, 07, 08, 09, 10, 11 already PASS / PASS WITH WARNINGS
- [ ] No Global Constraint violated

> Depends on: 04-sdd-mnemonic, 05-installer, 06-branding-look, 07-local-api-login, 08-rpiv-companions, 09-permission, 10-sdd-subagents, 11-dashboard

**Files:**
- Modify: `docs/user-manual/01-installation.md`
- Modify: `docs/user-manual/09-plugins.md`
- Modify: `docs/plan/01-workflow-new.md`
- Modify: `docs/plan/02-future.md`
- Modify: `docs/plan/03-skill-pi.md`
- Modify: `AGENTS.md`

**Interfaces:**
- Consumes: all prior step outputs (install path, look, Local API, companions, permission, subagents, dashboard)
- Produces: operator-facing documentation — terminal step, no dependents.

### Tasks

- [ ] 12.1 `[RED]` User manual documents npm and skillgrid install paths (WHAT: install docs) — both paths documented
  - [ ] 12.1.a Write failing test: `docs/user-manual/01-installation.md` contains `npm install -g skill-pi` and `skillgrid install --agents skill-pi`
  - [ ] 12.1.b Run to confirm fail — `Run: cd skill-pi && bun test test/docs.test.ts -t "docs cover install paths"` — Expected: FAIL
  - [ ] 12.1.c Minimal implementation: add skill-pi install section to `01-installation.md`
  - [ ] 12.1.d Run to confirm pass — `Run: cd skill-pi && bun test test/docs.test.ts -t "docs cover install paths"` — Expected: PASS
  - [ ] 12.1.e Commit — `docs: add skill-pi installation paths to user manual`
- [ ] 12.2 `[RED]` Plan harness lists include skill-pi (WHAT: harness targets) — skill-pi appears in plan docs
  - [ ] 12.2.a Write failing test: `docs/plan/01-workflow-new.md`, `docs/plan/02-future.md`, and `docs/plan/03-skill-pi.md` all mention `skill-pi`
  - [ ] 12.2.b Run to confirm fail — `Run: cd skill-pi && bun test test/docs.test.ts -t "plan harness lists include skill-pi"` — Expected: FAIL
  - [ ] 12.2.c Minimal implementation: add skill-pi entries to all three plan docs
  - [ ] 12.2.d Run to confirm pass — `Run: cd skill-pi && bun test test/docs.test.ts -t "plan harness lists include skill-pi"` — Expected: PASS
  - [ ] 12.2.e Commit — `docs: add skill-pi to plan harness lists`
- [ ] 12.3 `[AFK]` Coexistence with standalone pi is documented (WHAT) — `Run: cd skill-pi && bun test test/docs.test.ts -t "docs cover coexistence"` — Expected: PASS
- [ ] 12.4 `[AFK]` Docs mention Skillgrid Tokyo Night look / splash briefly (WHAT) — `Run: cd skill-pi && bun test test/docs.test.ts -t "docs cover branding"` — Expected: PASS
- [ ] 12.5 `[AFK]` Docs describe /login Local API for Ollama/vLLM (WHAT) — `Run: cd skill-pi && bun test test/docs.test.ts -t "docs cover local api"` — Expected: PASS
- [ ] 12.6 `[AFK]` Docs list bundled companions (WHAT) — `Run: cd skill-pi && bun test test/docs.test.ts -t "docs mention bundled companions"` — Expected: PASS
- [ ] 12.7 `[AFK]` Docs describe permission config, ask dialog, /plan//build (WHAT) — `Run: cd skill-pi && bun test test/docs.test.ts -t "docs cover permission"` — Expected: PASS
- [ ] 12.8 `[AFK]` Docs describe SDD subagent usage: /sdd, Result Contracts, sdd-full, user gate (WHAT) — `Run: cd skill-pi && bun test test/docs.test.ts -t "docs cover sdd subagents"` — Expected: PASS
- [ ] 12.9 `[AFK]` Docs describe dashboard: /dash, Ctrl+Shift+B, sections (WHAT) — `Run: cd skill-pi && bun test test/docs.test.ts -t "docs cover dashboard"` — Expected: PASS
- [ ] 12.10 `[AFK]` AGENTS.md mentions skill-pi in agent list (WHAT) — `Run: cd skill-pi && bun test test/docs.test.ts -t "agents md mentions skill-pi"` — Expected: PASS
- [ ] 12.11 `[AFK]` Pi/skill-pi plugin notes in 09-plugins.md (WHAT) — `Run: cd skill-pi && bun test test/docs.test.ts -t "plugins doc updated"` — Expected: PASS

### Verification

Verdict: `PENDING`  <!-- PASS | PASS WITH WARNINGS | FAIL -->

Evidence:

| Check | Run | Expected | Result | Notes |
|-------|-----|----------|--------|-------|
| Focused test | `cd skill-pi && bun test test/docs.test.ts` | PASS | | |
| Acceptance `@step-12` / `@p0` | `cd skill-pi && bun test --integration -t "@step-12"` | PASS | | |
| Runtime harness | `cd skill-pi && bun test test/docs.test.ts -t "docs cover install paths"` | PASS | | |
| Rollback boundary | revert docs changes; package still installs and runs | PASS | | |
| Global Constraints | — | held | | no marketing |

### Commit

When step DoD is met: `docs: skill-pi rollout — install, coexistence, features, and plan harness`
