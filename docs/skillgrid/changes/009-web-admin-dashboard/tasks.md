# Tasks: 009-web-admin-dashboard

> **STATUS:** `in-progress` (2026-09-08) — 0/5 steps PASS
>
> **For agentic workers:** REQUIRED SUB-SKILL: use subagent-execution (or simple-execution) to implement step-by-step. Steps use checkbox (`- [ ]`) syntax.

**Goal:** Replace the minimal embedded data viewer with a full admin dashboard served by `skillgrid serve` at `/` — memory browser, backlog browser, code-index status, and sessions view — backed by new HTTP endpoints where the service already has the method.

**Architecture:** The existing `internal/mnemonic/http` server (Go 1.22 mux, `embed.FS` static serving at `/`) is extended: new read endpoints expose service methods that are today MCP-only, a new `backlog` bridge package shells out to the `backlog` CLI (`--json`), and the embedded SPA under `internal/mnemonic/http/ui/` is rewritten in place (vanilla JS, no build step). See change.md Architecture decisions.

**Tech Stack:** Go 1.22+ (`skillgrid-cli`), existing SQLite/MCP service layer (untouched), embedded static SPA (vanilla HTML/JS/CSS, `embed.FS`), `backlog` CLI (Bun, shell-out, `--json` output).

**Spec:** `docs/skillgrid/changes/009-web-admin-dashboard/change.md`

**Acceptance:** `docs/skillgrid/changes/009-web-admin-dashboard/acceptance.feature` (`@step-NN`)

---

## Goal

Operators can open `http://127.0.0.1:7438/` from a running `skillgrid serve` and administer all four surfaces in one dashboard: browse/search/edit Mnemonic memory, view and manage Backlog tasks, check code-index health and search it, and review sessions + summaries — without leaving the browser and without a separate frontend process.

## Out of scope / Non-Goals

- Session relay / cleave handoff surfaces (change 006) — sessions tab shows what exists today; relay view lands when 006 ships
- New MCP tools — MCP tool names, signatures, and return shapes are frozen for this change
- Mutations other than what is listed in In scope — e.g. no backlog task *creation* (read + status change only)
- Authentication beyond the existing `SKILLGRID_HTTP_TOKEN` bearer gate on write routes; no login UI, no sessions/cookies
- CORS, remote hosting, or any non-127.0.0.1 serving story
- React/Next.js or any npm build pipeline in `skillgrid-cli` — the SPA stays build-less
- Changes to Mnemonic storage schema, indexing, or retrieval logic
- Redoing Swagger UI — `/swagger-ui` stays as-is
- Graph canvas visualization (vis.js/Sigma.js) — deferred to a follow-up change that depends on 005/008 edges; 009 ships a forward-compat placeholder in the Code tab only
- AI chat panel in the dashboard — the agent lives in the terminal; the dashboard is operator-facing, not an agent client
- **Multi-tenant teams / role layers / LLM proxy** — TencentDB's product-scale layer; 009 renders the single-operator governance (owner/visibility/usage) that 013 provides, not the multi-tenant machinery
- **Memory layering + governance data model** (L0–L3, owner, version, status, usage, visibility) — owned by `013-mnemonic-layered-memory-governance`; 009 **renders** those fields (asset library + layer drill-down + review + explicit share) and shows forward-compat placeholders where 013 has not landed yet

## Definition of Done

Change is done only when **all** of the following are true:

- [ ] Every success criterion / DoD checkbox in `change.md` is met
- [ ] Every `@step-NN` Feature in `acceptance.feature` has passing scenarios
- [ ] Every step below has Verdict `PASS` or `PASS WITH WARNINGS`
- [ ] No unchecked `- [ ]` under any `### Tasks`
- [ ] No **Global Constraint** violated
- [ ] Rollback path in `change.md` is still valid (or N/A documented)
- [ ] `## State` status is `done` (set at archive gate)

## Global Constraints

Copy verbatim from `change.md` (Error handling + Non-Goals + stack rules). Every step inherits these — do not restate per step.

- `backlog` binary not on PATH → `warn+continue`: `GET /backlog/*` → 503 `{error: "backlog CLI not found"}`; dashboard renders disabled Backlog tab
- `backlog` CLI exits non-zero or times out (>10s) → `warn+continue`: → 502 with CLI stderr excerpt (truncated 200 chars)
- `backlog --json` output not valid / unknown `schemaVersion` → `warn+continue`: → 502 with reason; tab shows error state, no partial render
- Unknown observation id on `GET /observations/{id}` → `abort`: 404 `{error: ...}`
- Unknown session id on `GET /sessions/{id}/summary` → `abort`: 404
- Invalid status on `POST /backlog/tasks/{id}/status` (not in `config.yml` statuses) → `abort`: 400 listing valid statuses (validated client-side too)
- Pin on already-pinned / unpin on non-pinned observation → `warn+continue`: idempotent 200 (service-level behavior preserved)
- Write route called without token when `SKILLGRID_HTTP_TOKEN` is set → `abort`: 401 (existing `requireWriteAuth` behavior, unchanged)
- MCP tool names, signatures, and return shapes are frozen — no `tools_*.go` file is touched
- No npm build pipeline, no React/Next.js, no external CDN assets — the SPA is vanilla HTML/JS/CSS served from `embed.FS`
- Server stays bound to `127.0.0.1` by default; no new env vars for serving
- Backlog mutations go through the `backlog` CLI, never direct `.backlog/*.md` file writes
- `GET /memory/project` keeps its current CWD-based behavior (quirk, not a regression — no change)
- **013 is a soft dependency**: the Memory governance/layer views (asset library metadata, layer drill-down, explicit share, in-place edit, review/status) render 013's fields when present and show a forward-compat placeholder (a labeled collapsed panel + the flat pre-013 view) when 013 has not landed — 009 must not block on 013; a failed/absent 013 field kills that widget only (per-widget error isolation), never the tab
- **No multi-tenant teams / role layers / LLM proxy** — 009 renders the single-operator governance that 013 provides, not the multi-tenant machinery
- **013 owns the governance/layer data model; 009 only renders it** — no change to 013's schema, layering, or governance logic; 04b is pure UI over the 013-backed `mem_*` HTTP surface

---

## State

```yaml
phase: spec
current_step: 01-memory-endpoints
status: in_progress
updated: 2026-09-09T00:00:00Z
```

## Step map

| NN | Step | Tag | Blocked by | Acceptance |
|----|------|-----|------------|------------|
| 01 | `memory-endpoints` | `@step-01` | — | Feature tagged `@step-01` |
| 02 | `backlog-bridge` | `@step-02` | — | Feature tagged `@step-02` |
| 03 | `dashboard-shell` | `@step-03` | 01 | Feature tagged `@step-03` |
| 04 | `dashboard-backlog-sessions` | `@step-04` | 03, 02 | Feature tagged `@step-04` |
| 04b | `memory-governance-view` | `@step-04b` | 03 (soft: 013) | Feature tagged `@step-04b` |
| 05 | `openapi-and-polish` | `@step-05` | 04, 04b | Feature tagged `@step-05` |

## Review workload (change-level)

| Field | Value |
|-------|-------|
| Estimated changed lines (change) | ~4100 (SPA ~3100, Go ~500, tests ~300, openapi ~200) |
| 400-line budget risk | Medium (SPA is the bulk; each step commits separately) |
| Chained PRs recommended | No |
| Delivery strategy | single-pr |

---

## 01-memory-endpoints

### Goal

The HTTP API exposes every read surface the dashboard's Memory and Sessions tabs need: observation detail, pin/unpin, session list, session summary.

### Out of scope / Non-Goals

- Backlog, code, web routes (belong to step 02 or already exist)
- Any service-method signature change
- MCP surface (frozen)
- UI work (belongs to steps 03–04)

### Definition of Done

This step is done only when:

- [ ] All `### Tasks` checkboxes below are `[x]`
- [ ] All `@step-01` scenarios in `acceptance.feature` pass
- [ ] `### Verification` Verdict is `PASS` or `PASS WITH WARNINGS`
- [ ] Produces contracts listed under Interfaces are available to dependents
- [ ] No Global Constraint violated

> Depends on: none

**Files:**
- Modify: `skillgrid-cli/internal/mnemonic/http/server.go`
- Test: `skillgrid-cli/internal/mnemonic/integration/integration_test.go`

**Interfaces:**
- Consumes: existing `service.Service` methods (`GetObservation` at service.go:694, `PinObservation`/`UnpinObservation` at service.go:780/790, `SessionSummary` at service.go:190, session listing via store access)
- Produces: `GET /observations/{id}`, `GET /sessions`, `GET /sessions/{id}/summary`, `POST /memory/observations/{id}/pin`, `POST /memory/observations/{id}/unpin` — five new HTTP routes on the existing mux

### Tasks

- [ ] 01.1 `[RED]` Threat: Authz — pin/unpin write-gated, read routes open
  - [ ] 01.1.a Write failing test: with `SKILLGRID_HTTP_TOKEN` set, `POST /memory/observations/{id}/pin` without token → 401; with token → 200. `GET /observations/{id}` without token → 200 (open read).
  - [ ] 01.1.b Run to confirm fail — `Run: go test ./skillgrid-cli/internal/mnemonic/integration/... -run TestStep01_Authz` — Expected: FAIL
  - [ ] 01.1.c Minimal implementation: add `POST /memory/observations/{id}/pin` and `POST /memory/observations/{id}/unpin` routes with `requireWriteAuth`; add `GET /observations/{id}` route open.
  - [ ] 01.1.d Run to confirm pass — `Run: go test ./skillgrid-cli/internal/mnemonic/integration/... -run TestStep01_Authz` — Expected: PASS
  - [ ] 01.1.e Commit — `feat(http): add observation detail + pin/unpin routes with authz`
- [ ] 01.2 `[RED]` GET /observations/{id} returns full untruncated content
  - [ ] 01.2.a Write failing test: save an observation, `GET /observations/{id}` returns the full content (same shape as `mem_get_observation`).
  - [ ] 01.2.b Run to confirm fail — `Run: go test ./skillgrid-cli/internal/mnemonic/integration/... -run TestStep01_ObservationDetail` — Expected: FAIL
  - [ ] 01.2.c Minimal implementation: wire `service.GetObservation` to the handler.
  - [ ] 01.2.d Run to confirm pass — `Run: go test ./skillgrid-cli/internal/mnemonic/integration/... -run TestStep01_ObservationDetail` — Expected: PASS
  - [ ] 01.2.e Commit — `feat(http): GET /observations/{id} returns full content`
- [ ] 01.3 `[AFK]` GET /observations/{id} → 404 for unknown id — `Run: go test ./skillgrid-cli/internal/mnemonic/integration/... -run TestStep01_ObservationDetail_404` — Expected: PASS
- [ ] 01.4 `[RED]` GET /sessions returns session list
  - [ ] 01.4.a Write failing test: start two sessions, `GET /sessions` returns both (id, title, started_at, status).
  - [ ] 01.4.b Run to confirm fail — `Run: go test ./skillgrid-cli/internal/mnemonic/integration/... -run TestStep01_SessionList` — Expected: FAIL
  - [ ] 01.4.c Minimal implementation: wire session listing to the handler.
  - [ ] 01.4.d Run to confirm pass — `Run: go test ./skillgrid-cli/internal/mnemonic/integration/... -run TestStep01_SessionList` — Expected: PASS
  - [ ] 01.4.e Commit — `feat(http): GET /sessions returns session list`
- [ ] 01.5 `[RED]` GET /sessions/{id}/summary returns summary; 404 for unknown session
  - [ ] 01.5.a Write failing test: start a session, end with summary, `GET /sessions/{id}/summary` returns the summary. Unknown session → 404.
  - [ ] 01.5.b Run to confirm fail — `Run: go test ./skillgrid-cli/internal/mnemonic/integration/... -run TestStep01_SessionSummary` — Expected: FAIL
  - [ ] 01.5.c Minimal implementation: wire `service.SessionSummary` to the handler.
  - [ ] 01.5.d Run to confirm pass — `Run: go test ./skillgrid-cli/internal/mnemonic/integration/... -run TestStep01_SessionSummary` — Expected: PASS
  - [ ] 01.5.e Commit — `feat(http): GET /sessions/{id}/summary + 404 for unknown`
- [ ] 01.6 `[AFK]` Pin/unpin idempotent and reflected in GET /observations ordering — `Run: go test ./skillgrid-cli/internal/mnemonic/integration/... -run TestStep01_PinUnpin` — Expected: PASS
- [ ] 01.7 `[AFK]` Existing routes and their auth behavior unchanged (regression) — `Run: go test ./skillgrid-cli/internal/mnemonic/integration/...` — Expected: PASS

### Verification

Verdict: `PENDING`

Evidence:

| Check | Run | Expected | Result | Notes |
|-------|-----|----------|--------|-------|
| Focused test | `go test ./skillgrid-cli/internal/mnemonic/integration/... -run TestStep01_` | PASS | | |
| Acceptance `@step-01` / `@p0` | manual smoke: `skillgrid serve` + curl each new route | PASS | | |
| Runtime harness | `go test ./skillgrid-cli/internal/mnemonic/http/...` | PASS | | |
| Rollback boundary | `git revert` + `go test ./...` | PASS | | |
| Global Constraints | — | held | | |

### Commit

When step DoD is met: `feat(http): memory + session read endpoints for dashboard`

---

## 02-backlog-bridge

### Goal

`/backlog/*` routes give the dashboard a JSON view of Backlog that degrades gracefully without the CLI.

### Out of scope / Non-Goals

- Backlog create/update/delete (read + status change only)
- Parsing `.backlog` files directly (shell-out to CLI)
- Auth model change
- UI work (belongs to step 04)

### Definition of Done

This step is done only when:

- [ ] All `### Tasks` checkboxes below are `[x]`
- [ ] All `@step-02` scenarios in `acceptance.feature` pass
- [ ] `### Verification` Verdict is `PASS` or `PASS WITH WARNINGS`
- [ ] No Global Constraint violated

> Depends on: none

**Files:**
- Create: `skillgrid-cli/internal/mnemonic/http/backlog/bridge.go`
- Create: `skillgrid-cli/internal/mnemonic/http/backlog/bridge_test.go`
- Modify: `skillgrid-cli/internal/mnemonic/http/server.go`

**Interfaces:**
- Consumes: `backlog` CLI on PATH (`--json` output); existing `requireWriteAuth` for the status route
- Produces: `GET /backlog/config`, `GET /backlog/tasks`, `GET /backlog/tasks/{id}`, `POST /backlog/tasks/{id}/status` — four new HTTP routes on the existing mux

### Tasks

- [ ] 02.1 `[RED]` Threat: Subprocess — CLI missing → 503 on all /backlog/*
  - [ ] 02.1.a Write failing test: remove `backlog` from PATH (or point to empty dir), all four `/backlog/*` routes → 503 `{error: "backlog CLI not found"}`.
  - [ ] 02.1.b Run to confirm fail — `Run: go test ./skillgrid-cli/internal/mnemonic/http/backlog/... -run TestStep02_CLI_Missing` — Expected: FAIL
  - [ ] 02.1.c Minimal implementation: `exec.LookPath("backlog")` check in the bridge; return 503 with JSON reason.
  - [ ] 02.1.d Run to confirm pass — `Run: go test ./skillgrid-cli/internal/mnemonic/http/backlog/... -run TestStep02_CLI_Missing` — Expected: PASS
  - [ ] 02.1.e Commit — `feat(backlog): bridge returns 503 when CLI missing`
- [ ] 02.2 `[RED]` Threat: Subprocess — CLI non-zero exit → 502 with stderr excerpt
  - [ ] 02.2.a Write failing test: fixture CLI that exits 1 with stderr "boom", `GET /backlog/tasks` → 502 with "boom" in the JSON error (truncated 200 chars).
  - [ ] 02.2.b Run to confirm fail — `Run: go test ./skillgrid-cli/internal/mnemonic/http/backlog/... -run TestStep02_CLI_Exit1` — Expected: FAIL
  - [ ] 02.2.c Minimal implementation: capture stderr, truncate to 200 chars, return 502.
  - [ ] 02.2.d Run to confirm pass — `Run: go test ./skillgrid-cli/internal/mnemonic/http/backlog/... -run TestStep02_CLI_Exit1` — Expected: PASS
  - [ ] 02.2.e Commit — `feat(backlog): bridge returns 502 with stderr on CLI failure`
- [ ] 02.3 `[RED]` Threat: Subprocess — CLI timeout >10s → 502
  - [ ] 02.3.a Write failing test: fixture CLI that sleeps 11s, `GET /backlog/tasks` → 502 timeout. Use `exec.CommandContext` with 10s deadline.
  - [ ] 02.3.b Run to confirm fail — `Run: go test ./skillgrid-cli/internal/mnemonic/http/backlog/... -run TestStep02_CLI_Timeout` — Expected: FAIL
  - [ ] 02.3.c Minimal implementation: `exec.CommandContext` with 10s timeout.
  - [ ] 02.3.d Run to confirm pass — `Run: go test ./skillgrid-cli/internal/mnemonic/http/backlog/... -run TestStep02_CLI_Timeout` — Expected: PASS
  - [ ] 02.3.e Commit — `feat(backlog): bridge times out at 10s with 502`
- [ ] 02.4 `[RED]` Threat: Subprocess — CLI stdout not valid JSON → 502
  - [ ] 02.4.a Write failing test: fixture CLI that prints "not json", `GET /backlog/tasks` → 502 with reason.
  - [ ] 02.4.b Run to confirm fail — `Run: go test ./skillgrid-cli/internal/mnemonic/http/backlog/... -run TestStep02_CLI_BadJSON` — Expected: FAIL
  - [ ] 02.4.c Minimal implementation: `json.Unmarshal` check; unknown `schemaVersion` → 502.
  - [ ] 02.4.d Run to confirm pass — `Run: go test ./skillgrid-cli/internal/mnemonic/http/backlog/... -run TestStep02_CLI_BadJSON` — Expected: PASS
  - [ ] 02.4.e Commit — `feat(backlog): bridge validates JSON + schemaVersion`
- [ ] 02.5 `[RED]` GET /backlog/config returns statuses/types/priorities
  - [ ] 02.5.a Write failing test: fixture CLI that emits config JSON, `GET /backlog/config` returns the parsed config.
  - [ ] 02.5.b Run to confirm fail — `Run: go test ./skillgrid-cli/internal/mnemonic/http/backlog/... -run TestStep02_Config` — Expected: FAIL
  - [ ] 02.5.c Minimal implementation: shell out to `backlog config list` or parse `.backlog/config.yml` via CLI.
  - [ ] 02.5.d Run to confirm pass — `Run: go test ./skillgrid-cli/internal/mnemonic/http/backlog/... -run TestStep02_Config` — Expected: PASS
  - [ ] 02.5.e Commit — `feat(backlog): GET /backlog/config`
- [ ] 02.6 `[RED]` GET /backlog/tasks returns versioned task-list JSON
  - [ ] 02.6.a Write failing test: fixture CLI that emits `task-list` JSON, `GET /backlog/tasks` returns the parsed list.
  - [ ] 02.6.b Run to confirm fail — `Run: go test ./skillgrid-cli/internal/mnemonic/http/backlog/... -run TestStep02_Tasks` — Expected: FAIL
  - [ ] 02.6.c Minimal implementation: shell out to `backlog task list --json`.
  - [ ] 02.6.d Run to confirm pass — `Run: go test ./skillgrid-cli/internal/mnemonic/http/backlog/... -run TestStep02_Tasks` — Expected: PASS
  - [ ] 02.6.e Commit — `feat(backlog): GET /backlog/tasks`
- [ ] 02.7 `[RED]` GET /backlog/tasks/{id} returns one task; 404 when not found
  - [ ] 02.7.a Write failing test: fixture CLI that emits a single task JSON for a known id; unknown id → 404.
  - [ ] 02.7.b Run to confirm fail — `Run: go test ./skillgrid-cli/internal/mnemonic/http/backlog/... -run TestStep02_TaskDetail` — Expected: FAIL
  - [ ] 02.7.c Minimal implementation: shell out to `backlog task view {id} --json`.
  - [ ] 02.7.d Run to confirm pass — `Run: go test ./skillgrid-cli/internal/mnemonic/http/backlog/... -run TestStep02_TaskDetail` — Expected: PASS
  - [ ] 02.7.e Commit — `feat(backlog): GET /backlog/tasks/{id} + 404`
- [ ] 02.8 `[RED]` POST /backlog/tasks/{id}/status (write-gated) changes status via CLI
  - [ ] 02.8.a Write failing test: fixture CLI, `POST /backlog/tasks/{id}/status` with token → status changes (verified by re-reading via fixture); without token → 401; invalid status → 400 listing valid statuses.
  - [ ] 02.8.b Run to confirm fail — `Run: go test ./skillgrid-cli/internal/mnemonic/http/backlog/... -run TestStep02_StatusChange` — Expected: FAIL
  - [ ] 02.8.c Minimal implementation: shell out to `backlog task {id} --status {new}`; validate status against config.
  - [ ] 02.8.d Run to confirm pass — `Run: go test ./skillgrid-cli/internal/mnemonic/http/backlog/... -run TestStep02_StatusChange` — Expected: PASS
  - [ ] 02.8.e Commit — `feat(backlog): POST /backlog/tasks/{id}/status (write-gated)`
- [ ] 02.9 `[AFK]` Mount /backlog/* routes on the existing mux — `Run: go test ./skillgrid-cli/internal/mnemonic/http/... -run TestStep02_Routes` — Expected: PASS

### Verification

Verdict: `PENDING`

Evidence:

| Check | Run | Expected | Result | Notes |
|-------|-----|----------|--------|-------|
| Focused test | `go test ./skillgrid-cli/internal/mnemonic/http/backlog/...` | PASS | | |
| Acceptance `@step-02` / `@p0` | manual smoke: `skillgrid serve` + curl each /backlog/* route | PASS | | |
| Runtime harness | `go test ./skillgrid-cli/internal/mnemonic/http/...` | PASS | | |
| Rollback boundary | `git revert` + `go test ./...` | PASS | | |
| Global Constraints | — | held | | |

### Commit

When step DoD is met: `feat(backlog): bridge + /backlog/* routes with graceful degradation`

---

## 03-dashboard-shell

### Goal

`GET /` serves the new dashboard shell with working Memory and Code tabs — the old viewer is replaced.

### Out of scope / Non-Goals

- Backlog and Sessions tab behavior (markup stubs may exist; belongs to step 04)
- Any backend change (routes land in steps 01–02)
- OpenAPI spec (belongs to step 05)

### Definition of Done

This step is done only when:

- [ ] All `### Tasks` checkboxes below are `[x]`
- [ ] All `@step-03` scenarios in `acceptance.feature` pass
- [ ] `### Verification` Verdict is `PASS` or `PASS WITH WARNINGS`
- [ ] Depends-on step(s) already PASS / PASS WITH WARNINGS
- [ ] No Global Constraint violated

> Depends on: 01-memory-endpoints

**Files:**
- Modify: `skillgrid-cli/internal/mnemonic/http/ui/index.html`
- Modify: `skillgrid-cli/internal/mnemonic/http/ui/app.js`
- Create: `skillgrid-cli/internal/mnemonic/http/ui/app.css`

**Interfaces:**
- Consumes: all step-01 routes (`GET /observations/{id}`, `GET /sessions`, `GET /sessions/{id}/summary`, `POST .../pin`, `POST .../unpin`), existing routes (`GET /search`, `GET /observations`, `GET /code/status`, `GET /code/search`, `GET /code/read`, `POST /code/index`, `GET /web/search`, `GET /relations/{id}`, `GET /context`)
- Produces: the dashboard shell (hash router, tab layout, project selector) that step 04 extends with Backlog + Sessions tabs

### Tasks

- [ ] 03.1 `[AFK]` Tab shell with hash routing + project selector (localStorage) — `Run: go test ./skillgrid-cli/internal/mnemonic/http/... -run TestStep03_Shell` — Expected: PASS
- [ ] 03.2 `[AFK]` Memory tab: search box (debounced), results list, click → detail pane with full content — `Run: go test ./skillgrid-cli/internal/mnemonic/http/... -run TestStep03_MemorySearch` — Expected: PASS
- [ ] 03.3 `[AFK]` Memory tab: pin/unpin + soft-delete actions with confirmation — `Run: go test ./skillgrid-cli/internal/mnemonic/http/... -run TestStep03_MemoryActions` — Expected: PASS
- [ ] 03.4 `[AFK]` Memory tab: relation drill-down with confidence badges — `Run: go test ./skillgrid-cli/internal/mnemonic/http/... -run TestStep03_Relations` — Expected: PASS
- [ ] 03.5 `[AFK]` Memory tab: query-first with suggested prompts (empty state) — `Run: go test ./skillgrid-cli/internal/mnemonic/http/... -run TestStep03_SuggestedPrompts` — Expected: PASS
- [ ] 03.6 `[AFK]` Code tab: freshness banner (last-indexed + stale) with Re-index action — `Run: go test ./skillgrid-cli/internal/mnemonic/http/... -run TestStep03_CodeFreshness` — Expected: PASS
- [ ] 03.7 `[AFK]` Code tab: status card, search, source view — `Run: go test ./skillgrid-cli/internal/mnemonic/http/... -run TestStep03_CodeSearch` — Expected: PASS
- [ ] 03.8 `[AFK]` Code tab: forward-compat graph placeholder (collapsed, file-list fallback) — `Run: go test ./skillgrid-cli/internal/mnemonic/http/... -run TestStep03_GraphPlaceholder` — Expected: PASS
- [ ] 03.9 `[AFK]` Show-numbers table twin on every data widget — `Run: go test ./skillgrid-cli/internal/mnemonic/http/... -run TestStep03_ShowNumbers` — Expected: PASS
- [ ] 03.10 `[AFK]` Web-cache view preserved (moved into Memory tab sub-section) — `Run: go test ./skillgrid-cli/internal/mnemonic/http/... -run TestStep03_WebCache` — Expected: PASS
- [ ] 03.11 `[AFK]` No external CDN assets (offline check) — `Run: go test ./skillgrid-cli/internal/mnemonic/http/... -run TestStep03_NoCDN` — Expected: PASS

### Verification

Verdict: `PENDING`

Evidence:

| Check | Run | Expected | Result | Notes |
|-------|-----|----------|--------|-------|
| Focused test | `go test ./skillgrid-cli/internal/mnemonic/http/... -run TestStep03_` | PASS | | |
| Acceptance `@step-03` / `@p0` | manual smoke: `skillgrid serve` + browser — Memory + Code tabs fully usable | PASS | | |
| Runtime harness | `go test ./skillgrid-cli/internal/mnemonic/http/...` | PASS | | |
| Rollback boundary | `git revert` + `skillgrid serve` + browser — old viewer still works | PASS | | |
| Global Constraints | — | held | | |

### Commit

When step DoD is met: `feat(ui): dashboard shell with Memory + Code tabs (replaces old viewer)`

---

## 04-dashboard-backlog-sessions

### Goal

Backlog and Sessions tabs are fully functional on the dashboard shell from step 03.

### Out of scope / Non-Goals

- Backlog mutations other than status (create/update/delete stay CLI-only)
- Relay/cleave surfaces (006)
- Any backend change (routes land in steps 01–02)
- Memory governance/layer views (belong to step 04b)

### Definition of Done

This step is done only when:

- [ ] All `### Tasks` checkboxes below are `[x]`
- [ ] All `@step-04` scenarios in `acceptance.feature` pass
- [ ] `### Verification` Verdict is `PASS` or `PASS WITH WARNINGS`
- [ ] Depends-on step(s) already PASS / PASS WITH WARNINGS
- [ ] No Global Constraint violated

> Depends on: 03-dashboard-shell, 02-backlog-bridge

**Files:**
- Modify: `skillgrid-cli/internal/mnemonic/http/ui/index.html`
- Modify: `skillgrid-cli/internal/mnemonic/http/ui/app.js`

**Interfaces:**
- Consumes: step-02 routes (`GET /backlog/config`, `GET /backlog/tasks`, `GET /backlog/tasks/{id}`, `POST /backlog/tasks/{id}/status`), step-01 routes (`GET /sessions`, `GET /sessions/{id}/summary`), existing `GET /context`
- Produces: the complete four-tab dashboard

### Tasks

- [ ] 04.1 `[AFK]` Backlog tab: board grouped by status, card shows id/title/priority/AC-progress — `Run: go test ./skillgrid-cli/internal/mnemonic/http/... -run TestStep04_BacklogBoard` — Expected: PASS
- [ ] 04.2 `[AFK]` Backlog tab: click → detail pane (references, assignees, dates); "Move to…" status action — `Run: go test ./skillgrid-cli/internal/mnemonic/http/... -run TestStep04_BacklogDetail` — Expected: PASS
- [ ] 04.3 `[AFK]` Backlog tab: completeness banner (CLI version + schemaVersion; warning on unknown schema) — `Run: go test ./skillgrid-cli/internal/mnemonic/http/... -run TestStep04_BacklogBanner` — Expected: PASS
- [ ] 04.4 `[AFK]` Backlog tab degraded state: on 503 shows "backlog CLI not found" + disables interactions — `Run: go test ./skillgrid-cli/internal/mnemonic/http/... -run TestStep04_BacklogDegraded` — Expected: PASS
- [ ] 04.5 `[AFK]` Sessions tab: session list (title, started_at, status), recent context, click → summary pane — `Run: go test ./skillgrid-cli/internal/mnemonic/http/... -run TestStep04_Sessions` — Expected: PASS
- [ ] 04.6 `[AFK]` Per-widget error isolation: a 500 on /backlog/tasks leaves Memory, Code, Sessions fully interactive — `Run: go test ./skillgrid-cli/internal/mnemonic/http/... -run TestStep04_ErrorIsolation` — Expected: PASS

### Verification

Verdict: `PENDING`

Evidence:

| Check | Run | Expected | Result | Notes |
|-------|-----|----------|--------|-------|
| Focused test | `go test ./skillgrid-cli/internal/mnemonic/http/... -run TestStep04_` | PASS | | |
| Acceptance `@step-04` / `@p0` | manual smoke: `skillgrid serve` + browser — Backlog board renders, status change round-trips, Sessions list + summaries render | PASS | | |
| Runtime harness | `go test ./skillgrid-cli/internal/mnemonic/http/...` | PASS | | |
| Rollback boundary | `git revert` + `skillgrid serve` + browser — step-03 tabs still work | PASS | | |
| Global Constraints | — | held | | |

### Commit

When step DoD is met: `feat(ui): Backlog + Sessions tabs on dashboard`

---

## 04b-memory-governance-view

### Goal

The Memory tab becomes a control panel (not just a browser) — asset library + layer drill-down + explicit share + in-place edit + review/status — rendering 013's governance/layer data with forward-compat placeholders where 013 has not landed.

### Out of scope / Non-Goals

- The agent-loadout *binding engine* (a later change; 04b shows a read-only equipping view)
- Multi-tenant teams/roles/proxy
- Any change to 013's data model (013 owns it; 04b renders it)
- Backlog, Code, Sessions tab behavior (belong to step 04)
- Any new Go HTTP route for governance (04b is pure UI over the existing/013-backed `mem_*` HTTP surface)

### Definition of Done

This step is done only when:

- [ ] All `### Tasks` checkboxes below are `[x]`
- [ ] All `@step-04b` scenarios in `acceptance.feature` pass
- [ ] `### Verification` Verdict is `PASS` or `PASS WITH WARNINGS`
- [ ] Depends-on step(s) already PASS / PASS WITH WARNINGS (03; 013 is a soft dep — its absence is a forward-compat state, not a block)
- [ ] No Global Constraint violated

> Depends on: 03-dashboard-shell, (soft: 013-mnemonic-layered-memory-governance)

**Files:**
- Modify: `skillgrid-cli/internal/mnemonic/http/ui/index.html`
- Modify: `skillgrid-cli/internal/mnemonic/http/ui/app.js`

**Interfaces:**
- Consumes: the 013-backed `mem_*` HTTP surface — additive governance fields on `GET /observations/{id}` / `GET /observations` (owner, version count, status, retrieval usage, visibility) from 013's `mem_governance`; the L0→L3 layer chain + per-layer provenance from 013's `mem_layers`; write-gated governance mutations for explicit share (`mem_share`), in-place edit (013 version-appending update), and status change; a read-only agent-loadout view (assets with visibility=`agent`). All same-origin JSON; absent 013 fields are a forward-compat placeholder, not an error.
- Produces: the Memory control-panel views (asset library, layer drill-down, share/edit/status, loadout) that step 05 documents in OpenAPI; 013 is a soft dependency — 04b ships and is fully usable with the pre-013 flat view + placeholders when 013 is absent

### Tasks

- [ ] 04b.1 `[RED]` Threat: Governance mutation / soft-dep — write-gated edit/share/status
  - [ ] 04b.1.a Write failing test: with `SKILLGRID_HTTP_TOKEN` set, the in-place edit / explicit share / status-change calls send the bearer token and receive 200; the same calls without a token receive 401 (write-gated). Assert in the SPA fetch layer that governance mutations attach the token.
  - [ ] 04b.1.b Run to confirm fail — `Run: go test ./skillgrid-cli/internal/mnemonic/http/... -run TestStep04b_GovernanceWriteGated` — Expected: FAIL
  - [ ] 04b.1.c Minimal implementation: governance mutation handlers in `app.js` attach the bearer token (existing write-auth path) before calling the 013-backed mutation surface.
  - [ ] 04b.1.d Run to confirm pass — `Run: go test ./skillgrid-cli/internal/mnemonic/http/... -run TestStep04b_GovernanceWriteGated` — Expected: PASS
  - [ ] 04b.1.e Commit — `feat(ui): governance edit/share/status are write-gated`
- [ ] 04b.2 `[RED]` Threat: Governance mutation — in-place edit appends a 013 version (re-readable, not overwritten)
  - [ ] 04b.2.a Write failing test: in-place edit an L1–L3 atom; re-fetch → the new content is current AND the prior content is recoverable in the version history view (the edit appended a 013 version, not an overwrite).
  - [ ] 04b.2.b Run to confirm fail — `Run: go test ./skillgrid-cli/internal/mnemonic/http/... -run TestStep04b_EditAppendsVersion` — Expected: FAIL
  - [ ] 04b.2.c Minimal implementation: in-place edit calls the 013 version-appending update; the detail pane renders the version history.
  - [ ] 04b.2.d Run to confirm pass — `Run: go test ./skillgrid-cli/internal/mnemonic/http/... -run TestStep04b_EditAppendsVersion` — Expected: PASS
  - [ ] 04b.2.e Commit — `feat(ui): in-place edit appends a recoverable 013 version`
- [ ] 04b.3 `[RED]` Threat: Governance mutation — share idempotent + 400 on unknown target
  - [ ] 04b.3.a Write failing test: explicit share to a valid visibility (`private`→`team`) returns 200 and re-shares idempotently (second call → 200, visibility unchanged); share to an unknown target → 400.
  - [ ] 04b.3.b Run to confirm fail — `Run: go test ./skillgrid-cli/internal/mnemonic/http/... -run TestStep04b_Share` — Expected: FAIL
  - [ ] 04b.3.c Minimal implementation: explicit share control calls `mem_share` (write-gated); render 400 reason for unknown targets.
  - [ ] 04b.3.d Run to confirm pass — `Run: go test ./skillgrid-cli/internal/mnemonic/http/... -run TestStep04b_Share` — Expected: PASS
  - [ ] 04b.3.e Commit — `feat(ui): explicit share (idempotent, 400 on unknown target)`
- [ ] 04b.4 `[RED]` Threat: soft-dep — with 013 absent, every governance/layer view renders a forward-compat placeholder + the flat view; the tab stays interactive
  - [ ] 04b.4.a Write failing test: against a pre-013 store (governance/layer fields absent), the asset library, layer drill-down, share, in-place edit, review/status, and loadout all render a labeled collapsed placeholder + the flat pre-013 view; the Memory tab's search/detail/actions remain fully interactive; a missing 013 field kills that widget only (per-widget isolation).
  - [ ] 04b.4.b Run to confirm fail — `Run: go test ./skillgrid-cli/internal/mnemonic/http/... -run TestStep04b_ForwardCompatPlaceholder` — Expected: FAIL
  - [ ] 04b.4.c Minimal implementation: detect absent 013 fields (not just HTTP error); render the labeled collapsed placeholder + flat view; keep the tab interactive.
  - [ ] 04b.4.d Run to confirm pass — `Run: go test ./skillgrid-cli/internal/mnemonic/http/... -run TestStep04b_ForwardCompatPlaceholder` — Expected: PASS
  - [ ] 04b.4.e Commit — `feat(ui): forward-compat placeholder when 013 absent`
- [ ] 04b.5 `[AFK]` Asset library — Memory list is an asset registry (owner, version count, status, retrieval usage, visibility badge) with search + show-numbers table twin — `Run: go test ./skillgrid-cli/internal/mnemonic/http/... -run TestStep04b_AssetLibrary` — Expected: PASS
- [ ] 04b.6 `[AFK]` Layer drill-down — L0→L1→L2→L3 chain with per-layer provenance, lazy-loaded per layer (a `mem_layers` call on expand); a distilled atom links back to its L0 source — `Run: go test ./skillgrid-cli/internal/mnemonic/http/... -run TestStep04b_LayerDrilldown` — Expected: PASS
- [ ] 04b.7 `[AFK]` Explicit share control — visibility (`private`→`team`/`restricted`/`agent` + an ACL editor for `restricted`) is an explicit click, never a default — `Run: go test ./skillgrid-cli/internal/mnemonic/http/... -run TestStep04b_ShareControl` — Expected: PASS
- [ ] 04b.8 `[AFK]` Review/status — status (`active`/`superseded`/`archived`) is visible and changeable; the personal→shared gate is a UI action — `Run: go test ./skillgrid-cli/internal/mnemonic/http/... -run TestStep04b_ReviewStatus` — Expected: PASS
- [ ] 04b.9 `[AFK]` Read-only agent loadout — a panel shows which assets a named agent is equipped with (visibility=`agent` bindings); the binding engine is out of scope — `Run: go test ./skillgrid-cli/internal/mnemonic/http/... -run TestStep04b_Loadout` — Expected: PASS
- [ ] 04b.10 `[AFK]` No external CDN assets — all 04b governance/layer views are same-origin JSON, binary works offline — `Run: go test ./skillgrid-cli/internal/mnemonic/http/... -run TestStep04b_NoCDN` — Expected: PASS

### Verification

Verdict: `PENDING`

Evidence:

| Check | Run | Expected | Result | Notes |
|-------|-----|----------|--------|-------|
| Focused test | `go test ./skillgrid-cli/internal/mnemonic/http/... -run TestStep04b_` | PASS | | |
| Acceptance `@step-04b` / `@p0` | manual smoke: `skillgrid serve` + browser — 013-provisioned store: asset library metadata + layer drill-down + share/edit/status round-trip; pre-013 store: placeholders + tab interactive | PASS | | |
| Runtime harness | `go test ./skillgrid-cli/internal/mnemonic/http/...` | PASS | | |
| Rollback boundary | `git revert` + `skillgrid serve` + browser — the flat Memory tab is intact (04b views are pure UI over 013's API) | PASS | | |
| Global Constraints | — | held | | |

### Commit

When step DoD is met: `feat(ui): Memory governance control-panel views (asset library, layer drill-down, share, edit, status, loadout)`

---

## 05-openapi-and-polish

### Goal

Machine-facing docs and polish match the shipped surface: OpenAPI spec updated, swagger-ui verified, user-manual updated, full DoD smoke passes.

### Out of scope / Non-Goals

- Feature work (bugs found here are fixed, nothing new)
- Any backend change

### Definition of Done

This step is done only when:

- [ ] All `### Tasks` checkboxes below are `[x]`
- [ ] All `@step-05` scenarios in `acceptance.feature` pass
- [ ] `### Verification` Verdict is `PASS` or `PASS WITH WARNINGS`
- [ ] Depends-on step(s) already PASS / PASS WITH WARNINGS
- [ ] No Global Constraint violated

> Depends on: 04-dashboard-backlog-sessions, 04b-memory-governance-view

**Files:**
- Modify: `skillgrid-cli/internal/mnemonic/http/ui/openapi.yaml`
- Modify: `docs/skillgrid/user-manual/` (serve page)

**Interfaces:**
- Consumes: all routes from steps 01–04, and the 013-backed governance/layer surface rendered by step 04b
- Produces: the complete, documented, smoke-tested dashboard

### Tasks

- [ ] 05.1 `[AFK]` openapi.yaml documents all step-01/02 routes with request/response examples — `Run: go test ./skillgrid-cli/internal/mnemonic/http/... -run TestStep05_OpenAPI` — Expected: PASS
- [ ] 05.2 `[AFK]` /swagger-ui loads and can exercise each new route — `Run: go test ./skillgrid-cli/internal/mnemonic/http/... -run TestStep05_SwaggerUI` — Expected: PASS
- [ ] 05.3 `[AFK]` User-manual serve section documents dashboard tabs + backlog-CLI dependency — `Run: go test ./skillgrid-cli/internal/mnemonic/http/... -run TestStep05_UserManual` — Expected: PASS
- [ ] 05.4 `[AFK]` Full DoD smoke: `go test ./...` + `go vet ./...` + manual browser pass — `Run: go test ./...` — Expected: PASS

### Verification

Verdict: `PENDING`

Evidence:

| Check | Run | Expected | Result | Notes |
|-------|-----|----------|--------|-------|
| Focused test | `go test ./skillgrid-cli/internal/mnemonic/http/... -run TestStep05_` | PASS | | |
| Acceptance `@step-05` / `@p0` | manual smoke: swagger-ui + user-manual + full DoD checklist | PASS | | |
| Runtime harness | `go test ./...` | PASS | | |
| Rollback boundary | `git revert` + `go test ./...` | PASS | | |
| Global Constraints | — | held | | |

### Commit

When step DoD is met: `docs(ui): OpenAPI spec + user-manual + DoD smoke`

---

## Archive gate checklist

- [ ] Change-level **Definition of Done** fully checked
- [ ] No unchecked `- [ ]` under any `### Tasks`
- [ ] Every step Verdict is `PASS` or `PASS WITH WARNINGS`
- [ ] No Global Constraint violated
- [ ] `## State` status is `done` and phase is `archive` (set by verify/archive)
- [ ] STATUS banner updated to `complete`
