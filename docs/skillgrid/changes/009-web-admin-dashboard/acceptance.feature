# Source: docs/skillgrid/changes/009-web-admin-dashboard/change.md
# Template: .agents/skills/_shared/templates/template-acceptance.feature
# Trace: change.md ## Goal + ## Definition of Done; tasks.md @step-NN verify lines.
# Mapping: @p0 scenarios ↔ change.md DoD / Testing strategy; @p1 = important failure paths.
# One Feature per step; tag each Feature with @step-NN matching tasks.md.
# Note: @step-04b renders 013's governance/layer data (asset library, layer drill-down,
# share, in-place edit, status, loadout) as a forward-compat placeholder (a labeled
# collapsed panel + the flat pre-013 view) when 013 is absent. 013 is a SOFT dependency —
# the Memory tab stays fully interactive with or without 013 (per-widget error isolation).

@step-01
Feature: HTTP API exposes memory + session read surfaces for the dashboard
  As an operator of a skillgrid-installed machine
  I want the HTTP API to expose observation detail, pin/unpin, session list, and session summary
  So that the dashboard's Memory and Sessions tabs have the data they need

  @happy @p0
  Scenario: GET /observations/{id} returns full untruncated content
    Given a saved observation with id 42 and content "the full body"
    When  I send GET /observations/42
    Then  I receive 200 with the full content "the full body" (same shape as mem_get_observation)

  @happy @p0
  Scenario: GET /sessions returns the session list
    Given two started sessions with titles "session A" and "session B"
    When  I send GET /sessions
    Then  I receive 200 with both sessions (id, title, started_at, status)

  @happy @p0
  Scenario: GET /sessions/{id}/summary returns the session summary
    Given a session that was ended with a summary
    When  I send GET /sessions/{id}/summary
    Then  I receive 200 with the session summary

  @edge
  Scenario: GET /observations/{id} returns 404 for unknown id
    Given no observation exists with id 999
    When  I send GET /observations/999
    Then  I receive 404 with a JSON error

  @edge
  Scenario: GET /sessions/{id}/summary returns 404 for unknown session
    Given no session exists with id "nope"
    When  I send GET /sessions/nope/summary
    Then  I receive 404

  @failure @p1
  Scenario: Pin without token returns 401 when SKILLGRID_HTTP_TOKEN is set
    Given SKILLGRID_HTTP_TOKEN is set to "secret"
    And  an observation with id 42
    When  I send POST /memory/observations/42/pin without a bearer token
    Then  I receive 401

  @failure @p1
  Scenario: Pin with token returns 200 and is idempotent
    Given SKILLGRID_HTTP_TOKEN is set to "secret"
    And  an observation with id 42
    When  I send POST /memory/observations/42/pin with bearer token "secret"
    Then  I receive 200
    And  I send POST /memory/observations/42/pin again with the same token
    Then  I receive 200 (idempotent)

  @happy @p0
  Scenario: GET /observations/{id} is open (no token required)
    Given an observation with id 42
    When  I send GET /observations/42 without a bearer token
    Then  I receive 200 (read routes are open)

@step-02
Feature: Backlog bridge provides /backlog/* routes with graceful degradation
  As an operator of a skillgrid-installed machine
  I want the dashboard to expose Backlog data via /backlog/* routes
  So that I can view and manage tasks without leaving the browser

  @happy @p0
  Scenario: GET /backlog/config returns statuses, types, and priorities
    Given the backlog CLI is available and configured
    When  I send GET /backlog/config
    Then  I receive 200 with statuses, types, and priorities from the CLI

  @happy @p0
  Scenario: GET /backlog/tasks returns the versioned task-list JSON
    Given the backlog CLI is available and there are 3 tasks
    When  I send GET /backlog/tasks
    Then  I receive 200 with the task-list JSON (id, title, status, priority, assignees, AC progress, isReady)

  @happy @p0
  Scenario: GET /backlog/tasks/{id} returns one task
    Given the backlog CLI is available and task TASK-001 exists
    When  I send GET /backlog/tasks/TASK-001
    Then  I receive 200 with the task details

  @failure @p1
  Scenario: GET /backlog/tasks/{id} returns 404 when task not found
    Given the backlog CLI is available
    And  no task with id TASK-999 exists
    When  I send GET /backlog/tasks/TASK-999
    Then  I receive 404

  @happy @p0
  Scenario: POST /backlog/tasks/{id}/status changes status via CLI (write-gated)
    Given SKILLGRID_HTTP_TOKEN is set to "secret"
    And  the backlog CLI is available and task TASK-001 has status "needs-triage"
    When  I send POST /backlog/tasks/TASK-001/status with body {"status": "ready-for-agent"} and bearer token "secret"
    Then  I receive 200
    And  the task status is now "ready-for-agent"

  @failure @p1
  Scenario: POST /backlog/tasks/{id}/status without token returns 401
    Given SKILLGRID_HTTP_TOKEN is set to "secret"
    And  the backlog CLI is available
    When  I send POST /backlog/tasks/TASK-001/status without a bearer token
    Then  I receive 401

  @failure @p1
  Scenario: POST /backlog/tasks/{id}/status with invalid status returns 400
    Given SKILLGRID_HTTP_TOKEN is set to "secret"
    And  the backlog CLI is available
    When  I send POST /backlog/tasks/TASK-001/status with body {"status": "not-a-real-status"} and bearer token "secret"
    Then  I receive 400 listing the valid statuses

  @edge
  Scenario: All /backlog/* routes return 503 when backlog CLI is missing
    Given the backlog CLI is not on PATH
    When  I send GET /backlog/config
    Then  I receive 503 with JSON error "backlog CLI not found"
    And  I send GET /backlog/tasks
    Then  I receive 503 with JSON error "backlog CLI not found"
    And  I send GET /backlog/tasks/TASK-001
    Then  I receive 503 with JSON error "backlog CLI not found"
    And  I send POST /backlog/tasks/TASK-001/status
    Then  I receive 503 with JSON error "backlog CLI not found"

  @failure @p1
  Scenario: GET /backlog/tasks returns 502 with stderr excerpt when CLI exits non-zero
    Given the backlog CLI exits 1 with stderr "boom: something broke"
    When  I send GET /backlog/tasks
    Then  I receive 502 with "boom: something broke" in the JSON error (truncated to 200 chars)

  @failure @p1
  Scenario: GET /backlog/tasks returns 502 when CLI times out
    Given the backlog CLI takes 11 seconds to respond
    When  I send GET /backlog/tasks
    Then  I receive 502 (timeout)

  @failure @p1
  Scenario: GET /backlog/tasks returns 502 when CLI output is not valid JSON
    Given the backlog CLI prints "not json"
    When  I send GET /backlog/tasks
    Then  I receive 502 with a reason

  @edge
  Scenario: GET /backlog/tasks returns 502 when schemaVersion is unknown
    Given the backlog CLI prints JSON with schemaVersion 99
    When  I send GET /backlog/tasks
    Then  I receive 502 with a reason

@step-03
Feature: Dashboard shell with working Memory and Code tabs replaces the old viewer
  As an operator of a skillgrid-installed machine
  I want the dashboard at / to have working Memory and Code tabs
  So that I can browse memory and check code-index health without leaving the browser

  @happy @p0
  Scenario: GET / serves the new dashboard shell with hash routing
    Given a running skillgrid serve
    When  I open http://127.0.0.1:7438/
    Then  I see the dashboard shell with 4 tabs (Memory, Backlog, Code, Sessions)
    And  the URL hash is #/memory (default tab)

  @happy @p0
  Scenario: Project selector persists to localStorage
    Given I have opened the dashboard
    When  I select project "my-project" in the selector
    Then  the selection is persisted to localStorage
    And  when I reload the page, the same project is selected

  @happy @p0
  Scenario: Memory tab search shows results and click opens detail pane
    Given I am on the Memory tab
    And  there is an observation matching "auth"
    When  I type "auth" in the search box (debounced)
    Then  I see a results list with the matching observation
    And  I click the observation
    Then  I see a detail pane with the full content

  @happy @p0
  Scenario: Memory tab pin/unpin and soft-delete actions work
    Given I am viewing an observation in the detail pane
    When  I click "Pin"
    Then  the observation is pinned (confirmed by re-fetching)
    And  I click "Unpin"
    Then  the observation is unpinned
    And  I click "Delete" and confirm
    Then  the observation is soft-deleted and removed from the list

  @happy @p0
  Scenario: Memory tab shows relation drill-down with confidence badges
    Given I am viewing an observation that has 2 relations (1 EXTRACTED, 1 INFERRED)
    When  I view the detail pane
    Then  I see a relations list with confidence badges (EXTRACTED, INFERRED)
    And  I click the related observation
    Then  I navigate to its detail pane

  @edge
  Scenario: Memory tab empty state shows suggested prompts
    Given I am on the Memory tab with no search results
    And  there are recent session prompts in the context
    When  I view the empty state
    Then  I see recent session prompts as clickable suggestions (not a blank "no results")

  @happy @p0
  Scenario: Code tab shows freshness banner with last-indexed and stale flag
    Given I am on the Code tab
    And  the code index was last indexed at 2026-09-08T08:00:00Z
    When  I view the Code tab
    Then  I see a top banner with "Last indexed: 2026-09-08T08:00:00Z" and the stale flag
    And  I see a "Re-index" action in the banner

  @happy @p0
  Scenario: Code tab search shows results and click opens source view
    Given I am on the Code tab
    And  the code index has a file matching "handler"
    When  I type "handler" in the search box
    Then  I see a results list with the matching chunk
    And  I click the result
    Then  I see a source view with the code

  @happy @p0
  Scenario: Code tab forward-compat graph placeholder is collapsed
    Given I am on the Code tab
    When  I view the Code tab
    Then  I see a collapsed panel labeled "Code graph (coming in 010)"
    And  expanding it shows a file-list fallback

  @happy @p0
  Scenario: Every data widget has a show-numbers table twin
    Given I am viewing any data widget (status card, search results, relation list)
    When  I toggle "show numbers"
    Then  I see the raw JSON/table behind the visual

  @edge
  Scenario: Web-cache view is preserved in Memory tab
    Given I am on the Memory tab
    When  I navigate to the web-cache sub-section
    Then  I see the web-cache search and results (equivalent to the old viewer)

  @edge
  Scenario: Dashboard works offline (no external CDN assets)
    Given I am on the dashboard
    When  I inspect the page source
    Then  there are no external CDN links (all assets are served from the binary)

  @failure @p1
  Scenario: A failed observation detail endpoint shows an error in the detail pane, not a blank page
    Given I am on the Memory tab
    And  GET /observations/{id} returns 500
    When  I click an observation in the results list
    Then  I see an error message in the detail pane (not a blank pane or console-only failure)
    And  the rest of the dashboard (search, Code tab) is still interactive

@step-04
Feature: Backlog and Sessions tabs are fully functional on the dashboard
  As an operator of a skillgrid-installed machine
  I want the Backlog and Sessions tabs to be fully functional
  So that I can manage tasks and review sessions without leaving the browser

  @happy @p0
  Scenario: Backlog tab renders a board grouped by status
    Given I am on the Backlog tab
    And  the backlog CLI is available with 5 tasks in 3 different statuses
    When  I view the Backlog tab
    Then  I see columns for each status from GET /backlog/config
    And  I see task cards (id, title, priority, AC progress) in the correct columns

  @happy @p0
  Scenario: Backlog tab click opens detail pane and status change works
    Given I am on the Backlog tab viewing a task card
    When  I click the task card
    Then  I see a detail pane (references, assignees, dates)
    And  I select "Move to… ready-for-agent"
    Then  the status changes (round-trip via the CLI) and the board re-renders

  @happy @p0
  Scenario: Backlog tab shows completeness banner
    Given I am on the Backlog tab
    And  the backlog CLI is available (version 1.2.3, schemaVersion 1)
    When  I view the Backlog tab
    Then  I see a banner with "backlog CLI v1.2.3 (schema v1)"

  @edge
  Scenario: Backlog tab warns on unknown schemaVersion
    Given I am on the Backlog tab
    And  the backlog CLI responds with schemaVersion 99
    When  I view the Backlog tab
    Then  I see a warning banner (not a hard error)

  @edge
  Scenario: Backlog tab degraded state when CLI is missing
    Given I am on the Backlog tab
    And  the backlog CLI is not on PATH (503)
    When  I view the Backlog tab
    Then  I see "backlog CLI not found — install with `skillgrid install`"
    And  interactions are disabled

  @happy @p0
  Scenario: Sessions tab shows session list, context, and summary
    Given I am on the Sessions tab
    And  there are 3 sessions with summaries
    When  I view the Sessions tab
    Then  I see a session list (title, started_at, status)
    And  I see a recent context section
    And  I click a session
    Then  I see a summary pane

  @failure @p1
  Scenario: Per-widget error isolation — a 500 on /backlog/tasks leaves other tabs interactive
    Given I am on the Backlog tab
    And  GET /backlog/tasks returns 500
    When  I view the Backlog tab
    Then  I see an error state in the Backlog tab
    And  I switch to the Memory tab
    Then  the Memory tab is fully interactive (search, detail, actions all work)

  @edge
  Scenario: Sessions tab 404 on unknown session shows in-tab error
    Given I am on the Sessions tab
    And  GET /sessions/{id}/summary returns 404
    When  I click a session with an unknown id
    Then  I see an error message in the summary pane (not a blank pane)

@step-04b
Feature: Memory governance control-panel views render 013 data with forward-compat placeholders
  As an operator of a skillgrid-installed machine
  I want the Memory tab to be a control panel (asset library, layer drill-down, explicit share, in-place edit, review/status)
  So that I can govern and correct memory — not just browse it — and 009 stays fully usable even before 013 lands

  # --- 013-present (happy path) ---

  @happy @p0
  Scenario: Asset library shows owner, version, status, usage, and visibility
    Given I am on the Memory tab against a 013-provisioned store
    And  there is an observation with owner "op-1", 3 versions, status "active", retrieval usage 12, visibility "team"
    When  I view the Memory list
    Then  each row shows owner, version count, status, retrieval usage, and a visibility badge
    And  toggling "show numbers" on the asset list shows the raw JSON/table

  @happy @p0
  Scenario: Layer drill-down renders the L0 to L3 chain lazy-loaded per layer
    Given I am viewing an observation that has a distilled L1 atom
    When  I expand the layer drill-down
    Then  I see the L0 -> L1 -> L2 -> L3 chain, loaded one layer at a time on expand
    And  each layer shows its provenance link
    And  the distilled L1 atom links back to its L0 source

  @happy @p0
  Scenario: Explicit share widens visibility by an explicit click
    Given I am viewing an observation with visibility "private"
    When  I select visibility "team" and confirm the share
    Then  the observation is now shared to "team" (an explicit action, never a default)
    And  a confirmation dialog is shown before the visibility change

  @happy @p0
  Scenario: In-place edit corrects an atom and appends a version
    Given I am viewing an L1 atom
    When  I edit it in place and save
    Then  the new content is current
    And  the prior content is recoverable in the version history view

  @happy @p0
  Scenario: Review status is visible and changeable
    Given I am viewing an observation with status "active"
    When  I change the status to "superseded"
    Then  the status is now "superseded" (the personal -> shared gate is a UI action)

  @happy @p0
  Scenario: Read-only agent loadout shows equipped assets
    Given I am on the Memory tab
    And  agent "alpha" is equipped with 2 assets (visibility "agent")
    When  I open the agent loadout panel for "alpha"
    Then  I see the 2 assets equipped to "alpha" (read-only)

  # --- forward-compat (013 absent — soft dependency) ---

  @edge
  Scenario: With 013 absent every governance and layer view renders a forward-compat placeholder
    Given I am on the Memory tab against a pre-013 store (no owner/version/status/usage/visibility or layer data)
    When  I view the asset library, layer drill-down, share, in-place edit, review/status, and loadout
    Then  each renders a labeled collapsed placeholder plus the flat pre-013 view
    And  the Memory tab's search, detail, and actions stay fully interactive

  # --- threat: Governance mutation / soft-dep (write-gated + per-widget isolation) ---

  @failure @p1
  Scenario: Governance mutations (edit, share, status) are write-gated
    Given SKILLGRID_HTTP_TOKEN is set to "secret"
    And  I am viewing an observation in the Memory tab
    When  I edit it in place, share it, or change its status without a bearer token
    Then  each governance mutation returns 401
    And  with the bearer token "secret", the same mutations return 200

  @failure @p1
  Scenario: In-place edit appends a 013 version rather than overwriting
    Given SKILLGRID_HTTP_TOKEN is set to "secret"
    And  I am viewing an L1 atom with prior content "v1"
    When  I edit it in place to "v2" with the bearer token
    Then  the current content is "v2"
    And  the prior content "v1" is still re-readable in the version history view

  @failure @p1
  Scenario: Share is idempotent and returns 400 on an unknown target
    Given SKILLGRID_HTTP_TOKEN is set to "secret"
    And  I am viewing an observation
    When  I share it to "team" with the token, then share to "team" again
    Then  both return 200 and the visibility stays "team" (idempotent)
    And  when I share it to an unknown target
    Then  I receive 400 and the visibility is unchanged

  @failure @p1
  Scenario: A failed or absent 013 field kills that widget only (per-widget isolation)
    Given I am on the Memory tab against a 013-provisioned store
    And  the layer-drill-down field for one observation is missing or returns 500
    When  I view that observation
    Then  the layer drill-down widget shows an error or placeholder
    And  the asset library, share, in-place edit, review/status, and loadout widgets remain fully interactive

@step-05
Feature: OpenAPI spec and polish match the shipped surface
  As an operator or agent consuming the HTTP API
  I want the OpenAPI spec and swagger-ui to document all new routes
  So that I can discover and exercise the API without reading source

  @happy @p0
  Scenario: openapi.yaml documents all new routes with examples
    Given a running skillgrid serve
    When  I fetch GET /openapi.yaml
    Then  I see all step-01 routes (GET /observations/{id}, GET /sessions, GET /sessions/{id}/summary, POST .../pin, POST .../unpin) with request/response examples
    And  I see all step-02 routes (GET /backlog/config, GET /backlog/tasks, GET /backlog/tasks/{id}, POST /backlog/tasks/{id}/status) with request/response examples

  @happy @p0
  Scenario: swagger-ui loads and can exercise each new route
    Given a running skillgrid serve
    When  I open /swagger-ui
    Then  I see all new routes listed
    And  I can try GET /observations/{id} with a valid id
    Then  I receive a 200 response in the swagger-ui try panel

  @edge
  Scenario: User-manual documents the dashboard tabs and backlog-CLI dependency
    Given I read the user-manual serve section
    When  I look for the dashboard documentation
    Then  I see all 4 tabs described (Memory, Backlog, Code, Sessions)
    And  I see the backlog-CLI dependency documented (what happens when it's missing)

  @happy @p0
  Scenario: Full DoD smoke passes
    Given a running skillgrid serve with the dashboard deployed
    When  I run go test ./...
    Then  all tests pass
    And  I run go vet ./...
    Then  no issues reported
    And  I open the dashboard in a browser
    Then  all 4 tabs render with live data

  @failure @p1
  Scenario: swagger-ui still serves when new routes are added
    Given a running skillgrid serve
    When  I open /swagger-ui
    Then  it loads (no 404 or JS error)
    And  the old routes are still documented

# Rules:
# - ≥1 @happy + @edge + @failure per step
# - Every change.md per-step WHAT bullet → a Scenario
# - Every applicable threat-matrix row → a Scenario in its owning @step-NN
# - Scenario names unique and referenceable from tasks.md `Run:` / Expected lines
# - @p0 must cover change.md Definition of Done user-visible criteria
