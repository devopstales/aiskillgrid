# Tasks: 008-mnemonic-community-knowledge-graph

> **STATUS:** `in-progress` (2026-09-08) — 0/3 steps PASS
>
> **For agentic workers:** REQUIRED SUB-SKILL: use subagent-driven-development (or simple-execution) to implement step-by-step. Steps use checkbox (`- [ ]`) syntax.

**Goal:** Extend the 005 code-intelligence graph with Leiden community detection + god nodes (architectural orientation), a **precomputed process layer** (entry-point → execution flows, so agents see what a subsystem *does* end-to-end), and a knowledge-graph layer that maps docs, configs, and SQL schemas as nodes linked to code — so agents see subsystems, their flows, and the "why" beyond the call graph.

**Architecture:** Additive `012_*` schema on top of 005's symbols/edges. Layer 1 runs `bluuewhale/loom` Leiden (pure-Go, zero deps) over the edges table to produce communities + god nodes + LLM-free labels. Layer 2 is the **process pass**: traces execution flows from 010's entry points through 005 call edges (depth-capped, cross-community flag) into LLM-labeled `processes`/`process_steps`. Layer 3 adds deterministic doc/config/SQL extractors feeding the same graph. See `change.md` decisions.

**Tech Stack:** Go (`skillgrid-cli`), SQLite (`modernc.org/sqlite`, CGo-free), gotreesitter graph from 005, `bluuewhale/loom` (pure-Go Leiden/Louvain), MCP (`mcp-go`), CLI, optional LLM for process labels (cached by content-hash; not required for graph structure).

**Spec:** `docs/skillgrid/changes/008-mnemonic-community-knowledge-graph/change.md`

**Acceptance:** `docs/skillgrid/changes/008-mnemonic-community-knowledge-graph/acceptance.feature` (`@step-NN`)

---

## Goal

Coding agents and operators get subsystem-level orientation (communities, god nodes), **precomputed execution flows (processes)**, and a knowledge graph that connects code to its docs, configs, and data schema — answering "what are the core modules?", "what does this subsystem *do*, end to end?", "which flow does this symbol participate in?", and "which code reads/writes this table?" without reading files.

## Out of scope / Non-Goals

- Re-implementing 005's symbols/edges/extractors/hybrid search — this change is additive on top of 005's graph
- Git-diff / `code_affected` (CodeGraph `affected` / graphify `prs`) — pure edge-traversal on 005's `edges`; owned by `010-mnemonic-framework-routes-affected`, not this change
- Framework-aware routes (`route`/`navigates` edges) + fsnotify watcher + staleness banner — CodeGraph-derived; owned by `010-mnemonic-framework-routes-affected`, not this change
- Video/audio/image semantic extraction — docs + configs + SQL only
- Replacing the per-project SQLite store or the `code_*` tool surface
- Cloud sync; multi-project community merge

## Definition of Done

Change is done only when **all** of the following are true:

- [ ] `code_communities` returns Leiden-clustered subsystems with LLM-free labels
- [ ] `code_god_nodes` returns the most-connected symbols (with `--exclude-hubs` to suppress utility super-hubs)
- [ ] A community view explains what a subsystem contains and its key entry points
- [ ] `code_processes` returns precomputed execution flows traced from entry points through call chains; each process has named steps + a cross-community flag; `code_process <name>` returns the full trace; `code_explain_symbol` surfaces which processes a symbol participates in
- [ ] Doc files (`.md`) with `[text](./other.md)` / `[[wikilinks]]` become nodes with `references` edges
- [ ] Config files (`.yaml`/`.toml`/`.json`) become nodes with `configures` edges to the code they configure
- [ ] SQL schema (`.sql` DDL) becomes table/column nodes with `reads`/`writes` edges to code that references them
- [ ] Every new edge carries a Confidence Label (`EXTRACTED | INFERRED | AMBIGUOUS`)
- [ ] Existing 005 `code_*` tools are unchanged (name + required params); `go test ./...` passes for touched packages
- [ ] Every Step Blueprint entry has a matching section in `tasks.md` with Verdict `PASS` or `PASS WITH WARNINGS`
- [ ] Every `@step-NN` Feature in `acceptance.feature` has passing `@happy`, `@edge`, and `@failure` scenarios
- [ ] Applicable threat-matrix rows have RED coverage that passed
- [ ] Testing strategy commands below are green
- [ ] Rollback path below is still valid (or N/A documented)
- [ ] Change archived under `docs/skillgrid/archive/008-mnemonic-community-knowledge-graph/`

## Global Constraints

Copy verbatim from `change.md` (Error handling + Non-Goals + stack rules). Every step inherits these — do not restate per step.

- Additive on top of 005 — never rewrite 005's symbols/edges/extractors
- CGo-free: Leiden is a pure-Go implementation (`bluuewhale/loom`, no C dependency)
- Every new edge carries a Confidence Label: `EXTRACTED | INFERRED | AMBIGUOUS`
- Community labels are LLM-free (derived from top god-node names + file paths, not an API call)
- Processes are **precomputed at index time** — traced from entry points through 005 call edges, LLM-labeled once, cached by process content-hash (re-label only on change); a `code_processes` query returns complete flows in one call with no per-query traversal
- Process tracing is deterministic (graph traversal + LLM label only); no per-query LLM call
- Process label cache keyed to process content-hash; LLM down → flow is cached **unlabeled**, never fabricated
- Process trace is depth/step-capped; a trace that can't be completed is truncated with a "stops at <symbol> (<reason>)" note (reuses 005's graph-stops), never a silent cut
- Doc/config/SQL extraction is deterministic (AST/regex/link-parse); no LLM for graph structure
- Unresolvable config refs are `AMBIGUOUS`, never dropped
- Per-file extract failure → `warn+continue` (fallback); never abort the whole index run
- Existing 005 `code_*` tools keep name + required params; all new tools use distinct `code_*` names
- Migration id `012_community_knowledge_graph.sql` — leave `011` for 005
- Leiden on a graph with < 2 nodes → `warn+continue` (single trivial community; no crash)
- Community label resolution finds no god node → `warn+continue` (label = "community-N"; never fabricated)
- Doc file with unparseable links → `warn+continue` (skip bad links; index the rest)
- SQL DDL that fails to parse → `warn+continue` (skip that statement; index the rest)
- Entry point with no traceable call chain → `warn+continue` (single-step or skipped; not fabricated)
- Bad / missing args on new `code_*` tools → `abort` with clear validation error; do not invent communities or processes
- No git-PR-impact, no video/audio/image extraction, no cloud sync or multi-project merge
- Communities are cached by content-hash; communities are advisory, not load-bearing (seeded RNG + pinned resolution)

---

## State

```yaml
phase: spec          # spec | apply | verify | archive
current_step: 01-community-detection
status: in_progress  # in_progress | blocked | done
updated: 2026-09-08T10:00:00+02:00
```

## Step map

| NN | Step | Tag | Blocked by | Acceptance |
|----|------|-----|------------|------------|
| 01 | `community-detection` | `@step-01` | — (005 done) | Feature tagged `@step-01` |
| 02 | `process-flows` | `@step-02` | 01, 010 (entry points) | Feature tagged `@step-02` |
| 03 | `knowledge-graph-nodes` | `@step-03` | 02 | Feature tagged `@step-03` |

## Review workload (change-level)

| Field | Value |
|-------|-------|
| Estimated changed lines (change) | ~1200–1800 |
| 400-line budget risk | High |
| Chained PRs recommended | Yes |
| Delivery strategy | ask-on-risk |

Honest forecast: three vertical slices (communities → processes → knowledge nodes), each likely its own stacked PR. Step 01 alone is near/over the 400-line budget (~450–600: schema + loom adapter + god nodes + labels + 3 tools). Do not attempt a single-PR delivery without explicit exception.

Suggested split: PR1 communities → PR2 processes → PR3 knowledge nodes · Chain strategy: stacked-to-main

---

## 01-community-detection

### Goal

Additive `012_*` schema + pure-Go Leiden (`bluuewhale/loom`) + god nodes + LLM-free community labels + `code_communities` / `code_god_nodes` / `code_explain_community` MCP/CLI tools so agents see subsystems and the most-connected concepts.

### Out of scope / Non-Goals

- Doc/config/SQL knowledge nodes (step 03)
- Process flows (step 02)
- Hybrid search; changing 005 tools

### Definition of Done

This step is done only when:

- [ ] All `### Tasks` checkboxes below are `[x]`
- [ ] All `@step-01` scenarios in `acceptance.feature` pass
- [ ] `### Verification` Verdict is `PASS` or `PASS WITH WARNINGS`
- [ ] Produces contracts listed under Interfaces are available to dependents
- [ ] No Global Constraint violated

> Depends on: none (005 done; 010 not required for this step)

**Files:**
- Create: `skillgrid-cli/internal/mnemonic/store/migrations/012_community_knowledge_graph.sql`
- Create: `skillgrid-cli/internal/mnemonic/community/leiden.go`
- Create: `skillgrid-cli/internal/mnemonic/community/godnodes.go`
- Create: `skillgrid-cli/internal/mnemonic/community/labels.go`
- Create: `skillgrid-cli/internal/mnemonic/mcp/tools_code_community.go`
- Modify: `skillgrid-cli/internal/mnemonic/service/service.go` (community facade)
- Modify: `skillgrid-cli/cmd/skillgrid/code_intel.go` (CLI community commands)
- Modify: `skillgrid-cli/go.mod` (add `bluuewhale/loom`)
- Test: `skillgrid-cli/internal/mnemonic/community/...`, `skillgrid-cli/internal/mnemonic/store/...`, `skillgrid-cli/internal/mnemonic/mcp/...`, `skillgrid-cli/cmd/skillgrid/...`

**Interfaces:**
- Consumes: 005's symbols/edges tables, existing store migration runner, 005 tool registration baseline
- Produces: `communities` / `community_meta` tables + god-node ranking; LLM-free community labels; `code_communities`, `code_god_nodes` (`--exclude-hubs`), `code_explain_community` MCP tools + CLI parity; seeded Leiden partition (content-hash cache key)

### Tasks

- [ ] 01.1 `[RED]` Mnemonic tool surface — 005 `code_search` schema stable before community tools land (Scenario: code_communities returns labeled subsystems and 005 tools stay stable) — threat: Mnemonic tool surface
  - [ ] 01.1.a Write failing test — assert 005 `code_*` tool names + required params are unchanged (baseline lock) before community tools register
  - [ ] 01.1.b Run to confirm fail — `Run: go test ./skillgrid-cli/internal/mnemonic/mcp/... -run ToolSurfaceBaseline -count=1` — Expected: FAIL (or red until baseline lock exists)
  - [ ] 01.1.c Minimal implementation — lock 005 tool-surface baseline assertions
  - [ ] 01.1.d Run to confirm pass — `Run: go test ./skillgrid-cli/internal/mnemonic/mcp/... -run ToolSurfaceBaseline -count=1` — Expected: PASS
  - [ ] 01.1.e Commit — `test(mnemonic): lock 005 tool surface baseline`
- [ ] 01.2 `[RED]` Leiden over 005 edges produces labeled communities (Scenario: code_communities returns labeled subsystems and 005 tools stay stable)
  - [ ] 01.2.a Write failing test — fixture graph of symbols/edges; after the community pass, `communities` rows partition the graph coherently, each community carries an LLM-free label; `code_communities` returns the partition; 005 tools still registered
  - [ ] 01.2.b Run to confirm fail — `Run: go test ./skillgrid-cli/internal/mnemonic/community/... ./skillgrid-cli/internal/mnemonic/mcp/... -run Communities -count=1` — Expected: FAIL
  - [ ] 01.2.c Minimal implementation — `012_*` migration + `community/leiden.go` (loom `NodeRegistry` + `LeidenOptions{Seed, Resolution, MaxIterations, NumRuns}`, partition → community rows) + `labels.go` + `go.mod` dep + `tools_code_community.go` `code_communities` + service facade
  - [ ] 01.2.d Run to confirm pass — `Run: go test ./skillgrid-cli/internal/mnemonic/community/... ./skillgrid-cli/internal/mnemonic/mcp/... -run Communities -count=1` — Expected: PASS
  - [ ] 01.2.e Commit — `feat(mnemonic): leiden community detection with llm-free labels`
- [ ] 01.3 `[RED]` Mnemonic tool surface — community tools register and reject bad args (Scenario: community tools reject bad args clearly) — threat: Mnemonic tool surface
  - [ ] 01.3.a Write failing test — `code_communities` / `code_god_nodes` / `code_explain_community` registered with distinct `code_*` names; bad/missing args (e.g. non-existent community id) rejected clearly with a validation error, no invented communities
  - [ ] 01.3.b Run to confirm fail — `Run: go test ./skillgrid-cli/internal/mnemonic/mcp/... -run CommunityToolsArgs -count=1` — Expected: FAIL
  - [ ] 01.3.c Minimal implementation — arg validation in `tools_code_community.go` + service
  - [ ] 01.3.d Run to confirm pass — `Run: go test ./skillgrid-cli/internal/mnemonic/mcp/... -run CommunityToolsArgs -count=1` — Expected: PASS
  - [ ] 01.3.e Commit — `feat(mnemonic): community tool arg validation`
- [ ] 01.4 `[RED]` God nodes ranked by degree with exclude-hubs (Scenario: code_god_nodes ranks hubs and exclude-hubs suppresses them)
  - [ ] 01.4.a Write failing test — `code_god_nodes` returns symbols ranked by degree; `--exclude-hubs` suppresses utility super-hubs from the ranking
  - [ ] 01.4.b Run to confirm fail — `Run: go test ./skillgrid-cli/internal/mnemonic/community/... -run GodNodes -count=1` — Expected: FAIL
  - [ ] 01.4.c Minimal implementation — `community/godnodes.go` (degree ranking + `--exclude-hubs`) + `code_god_nodes` tool + CLI flag
  - [ ] 01.4.d Run to confirm pass — `Run: go test ./skillgrid-cli/internal/mnemonic/community/... -run GodNodes -count=1` — Expected: PASS
  - [ ] 01.4.e Commit — `feat(mnemonic): god-node ranking with exclude-hubs`
- [ ] 01.5 `[AFK]` code_explain_community returns members and entry points (Scenario: code_explain_community explains a subsystem) — `Run: go test ./skillgrid-cli/internal/mnemonic/mcp/... ./skillgrid-cli/internal/mnemonic/service/... -count=1` — Expected: PASS
- [ ] 01.6 `[AFK]` Tiny graph (< 2 nodes) yields one trivial community, no crash (Scenario: Tiny graph yields a single trivial community) — `Run: go test ./skillgrid-cli/internal/mnemonic/community/... -run TinyGraph -count=1` — Expected: PASS
- [ ] 01.7 `[AFK]` Community label falls back to community-N when no god node (Scenario: Community label falls back when no god node exists) — `Run: go test ./skillgrid-cli/internal/mnemonic/community/... -count=1` — Expected: PASS
- [ ] 01.8 `[AFK]` Communities are seeded, pinned, and cached by content-hash (Scenario: Community partition is reproducible and cached by content-hash) — `Run: go test ./skillgrid-cli/internal/mnemonic/community/... ./skillgrid-cli/internal/mnemonic/codeindex/... -count=1` — Expected: PASS
- [ ] 01.9 `[AFK]` CLI parity for community commands via `code_intel.go` (Scenario: CLI community commands return the same views) — `Run: go test ./skillgrid-cli/cmd/skillgrid/... -count=1` — Expected: PASS

### Verification

Verdict: `PENDING`  <!-- PASS | PASS WITH WARNINGS | FAIL -->

Evidence:

| Check | Run | Expected | Result | Notes |
|-------|-----|----------|--------|-------|
| Focused test | `go test ./skillgrid-cli/internal/mnemonic/community/... ./skillgrid-cli/internal/mnemonic/mcp/... ./skillgrid-cli/internal/mnemonic/service/... -count=1` | PASS | | |
| Acceptance `@step-01` / `@p0` | BDD / mapped unit scenarios | PASS | | |
| Runtime harness | `skillgrid index` on fixture repo; `code_communities` / `code_god_nodes` via MCP + CLI | PASS | | |
| Rollback boundary | Drop `012_*` + `community/` + community tools + `go.mod` dep | PASS | | |
| Global Constraints | — | held | | |

### Commit

When step DoD is met: `feat(mnemonic): leiden community detection, god nodes, and community tools`

---

## 02-process-flows

### Goal

Precomputed process layer — trace execution flows from 010's entry points (routes/handlers/CLI mains) through 005 call edges into named `processes`/`process_steps` with a cross-community flag + LLM labels cached by content-hash, exposed via `code_processes` / `code_process` and surfaced in 005's `code_explain_symbol`.

### Out of scope / Non-Goals

- Community detection (step 01)
- Knowledge nodes (step 03)
- PDG/taint (011)

### Definition of Done

This step is done only when:

- [ ] All `### Tasks` checkboxes below are `[x]`
- [ ] All `@step-02` scenarios in `acceptance.feature` pass
- [ ] `### Verification` Verdict is `PASS` or `PASS WITH WARNINGS`
- [ ] Depends-on step(s) already PASS / PASS WITH WARNINGS
- [ ] No Global Constraint violated

> Depends on: 01-community-detection (community + cross-community flag), 010 (entry points)

**Files:**
- Create: `skillgrid-cli/internal/mnemonic/process/trace.go`
- Create: `skillgrid-cli/internal/mnemonic/process/labels.go`
- Create: `skillgrid-cli/internal/mnemonic/mcp/tools_code_process.go`
- Modify: `skillgrid-cli/internal/mnemonic/service/service.go` (process facade + process participation in `code_explain_symbol`)
- Test: `skillgrid-cli/internal/mnemonic/process/...`, `skillgrid-cli/internal/mnemonic/mcp/...`, `skillgrid-cli/internal/mnemonic/service/...`

**Interfaces:**
- Consumes: 010's entry points (route/handler/CLI-main symbols), 005 call edges, 01's communities + `code_explain_symbol`
- Produces: `processes` / `process_steps` tables (cross-community flag + content-hash key); deterministic entry-point → call-chain trace (depth-capped, graph-stops); LLM labels cached by content-hash; `code_processes`, `code_process <name>` MCP tools; `code_explain_symbol` now surfaces process participation (step N/M)

### Tasks

- [ ] 02.1 `[RED]` Mnemonic tool surface — `code_explain_symbol` surfaces process participation + 005 tools stable (Scenario: code_explain_symbol surfaces process participation and 005 tools stay stable) — threat: Mnemonic tool surface
  - [ ] 02.1.a Write failing test — after process pass, `code_explain_symbol <sym>` (005) includes which processes the symbol participates in (step N/M); 005 tool names + required params unchanged; process tools registered with distinct names
  - [ ] 02.1.b Run to confirm fail — `Run: go test ./skillgrid-cli/internal/mnemonic/mcp/... ./skillgrid-cli/internal/mnemonic/service/... -run ExplainSymbolProcess -count=1` — Expected: FAIL
  - [ ] 02.1.c Minimal implementation — wire process participation into `code_explain_symbol` + register `code_processes` / `code_process`
  - [ ] 02.1.d Run to confirm pass — `Run: go test ./skillgrid-cli/internal/mnemonic/mcp/... ./skillgrid-cli/internal/mnemonic/service/... -run ExplainSymbolProcess -count=1` — Expected: PASS
  - [ ] 02.1.e Commit — `feat(mnemonic): surface process participation in code_explain_symbol`
- [ ] 02.2 `[RED]` Entry-point → call-chain trace builds precomputed flows (Scenario: code_processes returns precomputed flows from entry points)
  - [ ] 02.2.a Write failing test — seeded from 010 entry points (routes/handlers/CLI mains), trace through 005 call edges into `processes` + `process_steps`; `code_processes` returns complete flows in one call (no per-query traversal); each process has named steps + a cross-community flag
  - [ ] 02.2.b Run to confirm fail — `Run: go test ./skillgrid-cli/internal/mnemonic/process/... ./skillgrid-cli/internal/mnemonic/mcp/... -run ProcessTrace -count=1` — Expected: FAIL
  - [ ] 02.2.c Minimal implementation — `process/trace.go` (entry-point seeding, depth-capped BFS over call edges, cross-community flag from 01, persist `processes`/`process_steps`) + `tools_code_process.go` `code_processes` + service facade
  - [ ] 02.2.d Run to confirm pass — `Run: go test ./skillgrid-cli/internal/mnemonic/process/... ./skillgrid-cli/internal/mnemonic/mcp/... -run ProcessTrace -count=1` — Expected: PASS
  - [ ] 02.2.e Commit — `feat(mnemonic): precomputed process flow tracing`
- [ ] 02.3 `[RED]` LLM labels cached by content-hash (Scenario: Process labels are cached by content-hash and re-labeled only on change)
  - [ ] 02.3.a Write failing test — a flow is LLM-labeled once and the label is cached keyed to the process content-hash; an unchanged re-index does NOT re-call the LLM (label reused); a changed flow (new content-hash) re-labels
  - [ ] 02.3.b Run to confirm fail — `Run: go test ./skillgrid-cli/internal/mnemonic/process/... -run ProcessLabels -count=1` — Expected: FAIL
  - [ ] 02.3.c Minimal implementation — `process/labels.go` (content-hash key, LLM call, cache store)
  - [ ] 02.3.d Run to confirm pass — `Run: go test ./skillgrid-cli/internal/mnemonic/process/... -run ProcessLabels -count=1` — Expected: PASS
  - [ ] 02.3.e Commit — `feat(mnemonic): llm process labels cached by content-hash`
- [ ] 02.4 `[RED]` Dispatch-boundary truncation with "stops at" note (Scenario: Trace stops at a dispatch boundary with a note)
  - [ ] 02.4.a Write failing test — a trace that hits a dispatch boundary (interface→impl, message bus, callback) is truncated with a "stops at <symbol> (<reason>)" note (reusing 005 graph-stops), not silently cut; `code_process <name>` shows each hop's Confidence Label + the stop note
  - [ ] 02.4.b Run to confirm fail — `Run: go test ./skillgrid-cli/internal/mnemonic/process/... -run DispatchStop -count=1` — Expected: FAIL
  - [ ] 02.4.c Minimal implementation — depth/step cap + dispatch-boundary detection in `process/trace.go`
  - [ ] 02.4.d Run to confirm pass — `Run: go test ./skillgrid-cli/internal/mnemonic/process/... -run DispatchStop -count=1` — Expected: PASS
  - [ ] 02.4.e Commit — `feat(mnemonic): dispatch-boundary truncation in process trace`
- [ ] 02.5 `[AFK]` LLM down → flow cached unlabeled, never fabricated (Scenario: LLM down caches the flow unlabeled) — `Run: go test ./skillgrid-cli/internal/mnemonic/process/... -count=1` — Expected: PASS
- [ ] 02.6 `[AFK]` Entry point with no traceable chain → single-step or skipped (Scenario: Untraceable entry point yields single-step or is skipped) — `Run: go test ./skillgrid-cli/internal/mnemonic/process/... -count=1` — Expected: PASS
- [ ] 02.7 `[AFK]` Cross-community process is flagged (Scenario: Cross-community process is flagged) — `Run: go test ./skillgrid-cli/internal/mnemonic/process/... -count=1` — Expected: PASS
- [ ] 02.8 `[AFK]` code_process returns the full trace with confidence (Scenario: code_process returns the full step-by-step trace) — `Run: go test ./skillgrid-cli/internal/mnemonic/mcp/... ./skillgrid-cli/internal/mnemonic/process/... -count=1` — Expected: PASS

### Verification

Verdict: `PENDING`  <!-- PASS | PASS WITH WARNINGS | FAIL -->

Evidence:

| Check | Run | Expected | Result | Notes |
|-------|-----|----------|--------|-------|
| Focused test | `go test ./skillgrid-cli/internal/mnemonic/process/... ./skillgrid-cli/internal/mnemonic/mcp/... ./skillgrid-cli/internal/mnemonic/service/... -count=1` | PASS | | |
| Acceptance `@step-02` / `@p0` | BDD / mapped unit scenarios | PASS | | |
| Runtime harness | `code_processes` / `code_process` on a fixture with 010 entry points | PASS | | |
| Rollback boundary | Drop `process/` + process tools + `code_explain_symbol` participation | PASS | | |
| Global Constraints | — | held | | |

### Commit

When step DoD is met: `feat(mnemonic): precomputed process flows with cached llm labels`

---

## 03-knowledge-graph-nodes

### Goal

Doc/config/SQL extractors + non-code nodes + `references`/`configures`/`reads`/`writes` edges + indexer hook + knowledge MCP/CLI tools so the graph spans code and its knowledge, and `code_path` can trace code→doc→config→table in one query.

### Out of scope / Non-Goals

- Community detection (step 01); process flows (step 02)
- Video/audio/image; LLM semantic doc↔code links

### Definition of Done

This step is done only when:

- [ ] All `### Tasks` checkboxes below are `[x]`
- [ ] All `@step-03` scenarios in `acceptance.feature` pass
- [ ] `### Verification` Verdict is `PASS` or `PASS WITH WARNINGS`
- [ ] Depends-on step(s) already PASS / PASS WITH WARNINGS
- [ ] No Global Constraint violated

> Depends on: 02-process-flows

**Files:**
- Create: `skillgrid-cli/internal/mnemonic/knowledge/doc.go`
- Create: `skillgrid-cli/internal/mnemonic/knowledge/config.go`
- Create: `skillgrid-cli/internal/mnemonic/knowledge/sql.go`
- Create: `skillgrid-cli/internal/mnemonic/mcp/tools_code_knowledge.go`
- Modify: `skillgrid-cli/internal/mnemonic/codeindex/indexer.go` (hook community + process + knowledge passes after 005 extraction, same tx)
- Modify: `skillgrid-cli/internal/mnemonic/mcp/server.go` (register new tool sets)
- Modify: `skillgrid-cli/cmd/skillgrid/main.go` (CLI dispatch)
- Test: `skillgrid-cli/internal/mnemonic/knowledge/...`, `skillgrid-cli/internal/mnemonic/codeindex/...`, `skillgrid-cli/internal/mnemonic/mcp/...`, `skillgrid-cli/cmd/skillgrid/...`

**Interfaces:**
- Consumes: 005 symbols/edges + `code_path`, 01 communities, 02 processes, indexer hook seam
- Produces: `doc_nodes` / `config_nodes` / `sql_schema_nodes` tables; `references` (doc→doc), `configures` (config→code), `reads`/`writes` (code→table) edges; indexer hook running community + process + knowledge passes in the same incremental tx; knowledge query tools; `code_path` spanning code→doc→config→table

### Tasks

- [ ] 03.1 `[RED]` Mnemonic tool surface — knowledge tools register + 005 tools stable + bad args rejected (Scenario: Knowledge tools register and reject bad args) — threat: Mnemonic tool surface
  - [ ] 03.1.a Write failing test — knowledge query tools registered with distinct `code_*` names; 005 tools still name/param-stable; bad/missing args rejected clearly
  - [ ] 03.1.b Run to confirm fail — `Run: go test ./skillgrid-cli/internal/mnemonic/mcp/... -run KnowledgeTools -count=1` — Expected: FAIL
  - [ ] 03.1.c Minimal implementation — `tools_code_knowledge.go` + `server.go` registration + `main.go` dispatch
  - [ ] 03.1.d Run to confirm pass — `Run: go test ./skillgrid-cli/internal/mnemonic/mcp/... -run KnowledgeTools -count=1` — Expected: PASS
  - [ ] 03.1.e Commit — `feat(mnemonic): register knowledge-graph query tools`
- [ ] 03.2 `[RED]` Markdown links/wikilinks become references edges (Scenario: Markdown links and wikilinks become references edges)
  - [ ] 03.2.a Write failing test — `.md` files with `[text](./other.md)` and `[[wikilinks]]` produce `doc_nodes` + `references` edges between doc nodes, each confidence-labeled
  - [ ] 03.2.b Run to confirm fail — `Run: go test ./skillgrid-cli/internal/mnemonic/knowledge/... -run DocLinks -count=1` — Expected: FAIL
  - [ ] 03.2.c Minimal implementation — `knowledge/doc.go` (markdown link/wikilink parser → `references`)
  - [ ] 03.2.d Run to confirm pass — `Run: go test ./skillgrid-cli/internal/mnemonic/knowledge/... -run DocLinks -count=1` — Expected: PASS
  - [ ] 03.2.e Commit — `feat(mnemonic): markdown doc-link extractor`
- [ ] 03.3 `[RED]` Config refs become configures edges (Scenario: Config references become configures edges)
  - [ ] 03.3.a Write failing test — `.yaml`/`.toml`/`.json` files produce `config_nodes` + `configures` edges to the code they configure; explicit syntax is `EXTRACTED`, resolved-but-inferred refs are `INFERRED`
  - [ ] 03.3.b Run to confirm fail — `Run: go test ./skillgrid-cli/internal/mnemonic/knowledge/... -run ConfigRefs -count=1` — Expected: FAIL
  - [ ] 03.3.c Minimal implementation — `knowledge/config.go` (key→symbol resolver → `configures`)
  - [ ] 03.3.d Run to confirm pass — `Run: go test ./skillgrid-cli/internal/mnemonic/knowledge/... -run ConfigRefs -count=1` — Expected: PASS
  - [ ] 03.3.e Commit — `feat(mnemonic): config-ref extractor`
- [ ] 03.4 `[RED]` SQL DDL becomes table/column nodes with reads/writes (Scenario: SQL DDL becomes table and column nodes with reads and writes)
  - [ ] 03.4.a Write failing test — `.sql` DDL produces `sql_schema_nodes` (tables + columns) and code that references them gets `reads`/`writes` edges, confidence-labeled
  - [ ] 03.4.b Run to confirm fail — `Run: go test ./skillgrid-cli/internal/mnemonic/knowledge/... -run SqlSchema -count=1` — Expected: FAIL
  - [ ] 03.4.c Minimal implementation — `knowledge/sql.go` (DDL parser + identifier-ref scan → `reads`/`writes`)
  - [ ] 03.4.d Run to confirm pass — `Run: go test ./skillgrid-cli/internal/mnemonic/knowledge/... -run SqlSchema -count=1` — Expected: PASS
  - [ ] 03.4.e Commit — `feat(mnemonic): sql-schema extractor`
- [ ] 03.5 `[AFK]` Unresolvable config ref is AMBIGUOUS, not dropped (Scenario: Unresolvable config ref is ambiguous not dropped) — `Run: go test ./skillgrid-cli/internal/mnemonic/knowledge/... -count=1` — Expected: PASS
- [ ] 03.6 `[AFK]` Malformed doc file: skip bad links, index the rest (Scenario: Malformed doc file falls back and indexes the rest) — `Run: go test ./skillgrid-cli/internal/mnemonic/knowledge/... -count=1` — Expected: PASS
- [ ] 03.7 `[AFK]` Malformed SQL statement: skip it, index the rest (Scenario: Malformed SQL statement is skipped and the rest is indexed) — `Run: go test ./skillgrid-cli/internal/mnemonic/knowledge/... -count=1` — Expected: PASS
- [ ] 03.8 `[AFK]` Indexer hook runs all three passes in the same transaction (Scenario: Indexer hook runs community, process, and knowledge in one transaction) — `Run: go test ./skillgrid-cli/internal/mnemonic/codeindex/... -count=1` — Expected: PASS
- [ ] 03.9 `[AFK]` code_path traces code to doc to config to table (Scenario: code_path traces code to doc to config to table) — `Run: go test ./skillgrid-cli/internal/mnemonic/mcp/... ./skillgrid-cli/internal/mnemonic/knowledge/... -count=1` — Expected: PASS

### Verification

Verdict: `PENDING`  <!-- PASS | PASS WITH WARNINGS | FAIL -->

Evidence:

| Check | Run | Expected | Result | Notes |
|-------|-----|----------|--------|-------|
| Focused test | `go test ./skillgrid-cli/internal/mnemonic/knowledge/... ./skillgrid-cli/internal/mnemonic/codeindex/... ./skillgrid-cli/internal/mnemonic/mcp/... -count=1` | PASS | | |
| Acceptance `@step-03` / `@p0` | BDD / mapped unit scenarios | PASS | | |
| Runtime harness | `skillgrid index` on fixture with docs/configs/SQL; `code_path` code→doc→config→table | PASS | | |
| Rollback boundary | Drop `knowledge/` + knowledge tools + indexer hook | PASS | | |
| Global Constraints | — | held | | |

### Commit

When step DoD is met: `feat(mnemonic): doc, config, and sql knowledge-graph nodes`

---

## Archive gate checklist

- [ ] Change-level **Definition of Done** fully checked
- [ ] No unchecked `- [ ]` under any `### Tasks`
- [ ] Every step Verdict is `PASS` or `PASS WITH WARNINGS`
- [ ] No Global Constraint violated
- [ ] `## State` status is `done` and phase is `archive` (set by verify/archive)
- [ ] STATUS banner updated to `complete`
