# Tasks: 011-mnemonic-pdg-taint

> **STATUS:** `in-progress` (2026-09-08) — 0/2 steps PASS
>
> **For agentic workers:** REQUIRED SUB-SKILL: use subagent-execution (or simple-execution) to implement step-by-step. Steps use checkbox (`- [ ]`) syntax.

**Goal:** Add an opt-in statement-level analysis tier to Mnemonic's code graph: per-function control-flow graphs → program-dependence graphs → taint findings (source→sink data-flow).

**Architecture:** Two additive, opt-in passes gated behind `--pdg` on top of 005's graph: a CFG→PDG pass (basic blocks + control/data dependence, intraprocedural M1, built from the existing gotreesitter AST) and a taint solver (deterministic configurable source/sink sets, source→sink paths, confidence-labeled, stops-at-boundary-not-fabricated). A non-`--pdg` index is byte-for-byte unchanged. See `change.md` decisions.

**Tech Stack:** Go (`skillgrid-cli`), CGo-free SQLite (`modernc.org/sqlite`), existing gotreesitter AST from 005 (no new grammar), MCP (`mcp-go`), CLI.

**Spec:** `docs/skillgrid/changes/011-mnemonic-pdg-taint/change.md`

**Acceptance:** `docs/skillgrid/changes/011-mnemonic-pdg-taint/acceptance.feature` (`@step-NN`)

---

## Goal

An agent gets statement-level data-flow answers on demand — "does this user input reach that SQL write?", "what's the control path that enables this branch?", "trace this value source→sink" — without reading files or re-deriving the CFG by hand. It is opt-in so the common 005/008/010 path stays lean.

## Out of scope / Non-Goals

- Re-implementing 005's symbols/edges/extractors, 008's communities/processes/knowledge, or 010's routes/affected/rename/watcher
- Whole-program interprocedural taint (M1 is **intraprocedural** per-function; a call-boundary summary is a later pass)
- A full data-flow engine beyond CFG→PDG→taint (no type inference beyond what 005 already resolves, no alias analysis beyond 005's receiver resolution)
- Default-on indexing (PDG is a separate, opt-in pass — it must not slow the common path)
- New languages beyond what 005's gotreesitter already parses (the PDG pass reuses 005's AST; M1 ships the same language set as 005)
- A new tree-sitter grammar or any CGo (CFG is a re-read of 005's existing AST)
- LLM-suggested sources/sinks (source/sink sets are deterministic and configurable)

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

- Additive + **opt-in** — `--pdg` is the only trigger; a non-`--pdg` index never populates PDG/taint tables and is byte-for-byte the 005/008/010 graph
- CGo-free: CFG is built from the existing gotreesitter AST (no new grammar, no C); the PDG + taint solver are pure Go
- Intraprocedural (M1): CFG/PDG/taint are per-function; crossing a call boundary without a resolved callee is `AMBIGUOUS`/truncated, not a fabricated intraprocedural hop. Interprocedural summaries are a later pass
- Every PDG/taint edge carries a Confidence Label: `EXTRACTED | INFERRED | AMBIGUOUS`; a taint path is only `EXTRACTED` where every hop is a resolved data-dependence, else `INFERRED`/`AMBIGUOUS`
- Taint **sources** and **sinks** are a configurable, deterministic set (e.g. sources: request params / env / file reads; sinks: SQL exec / shell exec / template render / file write); a finding is source→sink, not a guess
- A taint finding with no source→sink path is **not reported** (never fabricated); a path that ends at an unresolved boundary is reported with a "stops at <boundary>" note (reuses 005's graph-stops philosophy)
- Existing 005/008/010 `code_*` tools keep name + required params; all new tools use distinct `code_*` names
- Migration id `014_pdg_taint.sql` — leave `011` (005), `012` (008), `013` (010) as-is
- Function with a malformed/unparseable CFG → `warn+continue`; skip that function's PDG; index the rest
- PDG construction exceeds depth/step cap on a large function → `warn+continue`; truncate with a "stops at <block>" note; never abort
- Taint path crosses an unresolved call boundary → `warn+continue`; path marked `AMBIGUOUS` / "stops at <boundary>"; not fabricated
- `code_pdg_query` on a statement with no PDG (non-`--pdg` index) → `warn+continue`; clear "run `--pdg`" message; empty result, not an error
- Unknown / missing symbol or statement → `warn+continue`; not-found; no fabricated blocks or dependences
- Bad / missing args on new `code_*` tools → `abort` with clear validation error; do not invent findings
- Existing 005/008/010 `code_*` tool call → unchanged name + required params; a non-`--pdg` index is byte-for-byte unchanged

---

## State

```yaml
phase: spec          # spec | apply | verify | archive
current_step: 01-cfg-pdg
status: in_progress  # in_progress | blocked | done
updated: 2026-09-08
```

## Step map

| NN | Step | Tag | Blocked by | Acceptance |
|----|------|-----|------------|------------|
| 01 | `cfg-pdg` | `@step-01` | — (005 done) | Feature tagged `@step-01` |
| 02 | `taint-solver` | `@step-02` | 01 | Feature tagged `@step-02` |

## Review workload (change-level)

| Field | Value |
|-------|-------|
| Estimated changed lines (change) | ~800–1200 |
| 400-line budget risk | Medium |
| Chained PRs recommended | No |
| Delivery strategy | single-pr |

Two vertical slices: step 01 (schema + CFG + PDG + opt-in hook + `code_pdg_query`), step 02 (taint solver + `code_taint` + CLI/CLI-flag parity). Each likely its own work-unit commit; a single PR is acceptable for the change as a whole given the opt-in blast radius is bounded behind `--pdg`.

---

## 01-cfg-pdg

### Goal

Opt-in per-function CFG + control/data-dependence PDG so statement-level structure is queryable.

### Out of scope / Non-Goals

- Taint (step 02); interprocedural flow; changing 005/008/010 tools
- A new tree-sitter grammar or CGo; default-on PDG indexing

### Definition of Done

This step is done only when:

- [ ] All `### Tasks` checkboxes below are `[x]`
- [ ] All `@step-01` scenarios in `acceptance.feature` pass
- [ ] `### Verification` Verdict is `PASS` or `PASS WITH WARNINGS`
- [ ] Produces contracts listed under Interfaces are available to dependents
- [ ] No Global Constraint violated

> Depends on: none (005 done)

**Files:**
- Create: `skillgrid-cli/internal/mnemonic/store/migrations/014_pdg_taint.sql`
- Create: `skillgrid-cli/internal/mnemonic/pdg/cfg.go`
- Create: `skillgrid-cli/internal/mnemonic/pdg/pdg.go`
- Create: `skillgrid-cli/internal/mnemonic/mcp/tools_code_pdg.go`
- Modify: `skillgrid-cli/internal/mnemonic/codeindex/indexer.go`
- Modify: `skillgrid-cli/internal/mnemonic/mcp/server.go`
- Modify: `skillgrid-cli/cmd/skillgrid/code_intel.go`
- Modify: `skillgrid-cli/cmd/skillgrid/main.go`
- Test: `skillgrid-cli/internal/mnemonic/pdg/...`, `skillgrid-cli/internal/mnemonic/store/...`, `skillgrid-cli/internal/mnemonic/codeindex/...`, `skillgrid-cli/internal/mnemonic/mcp/...`

**Interfaces:**
- Consumes: 005's gotreesitter AST + symbols/edges + existing store migration runner + indexer `Run`
- Produces: `014_pdg_taint.sql` (cfg_blocks, cfg_edges, pdg_edges, taint_findings); per-function CFG (basic blocks + `CFG` edges); control-dependence + data-dependence PDG edges (each confidence-labeled); opt-in `--pdg` hook in `Indexer.Run` (same incremental transaction); `code_pdg_query <symbol> <statement>` MCP/CLI tool; `skillgrid index --pdg` flag

### Tasks

- [ ] 01.1 `[RED]` Mnemonic tool surface + Opt-in isolation — opt-in gate: a non-`--pdg` index is byte-for-byte unchanged (Scenario: Non-opt-in index is byte-for-byte unchanged) — threat: Mnemonic tool surface; Opt-in isolation
  - [ ] 01.1.a Write failing test — (1) index a fixture repo without `--pdg` and assert the `cfg_blocks`/`cfg_edges`/`pdg_edges`/`taint_findings` tables are empty (created but unpopulated); (2) assert every 005/008/010 `code_*` tool output for a representative query is byte-for-byte identical to a pre-011 baseline; (3) assert `code_pdg_query` against a non-`--pdg` index returns a clear "run `--pdg`" message and an empty result, not an error
  - [ ] 01.1.b Run to confirm fail — `Run: go test ./skillgrid-cli/internal/mnemonic/pdg/... ./skillgrid-cli/internal/mnemonic/codeindex/... ./skillgrid-cli/internal/mnemonic/store/... ./skillgrid-cli/internal/mnemonic/mcp/... -run OptInIsolation -count=1` — Expected: FAIL
  - [ ] 01.1.c Minimal implementation — opt-in gate: `--pdg` flag on the indexer; PDG/taint pass runs only when the flag is set, in the same incremental transaction as 005's extraction; without it, no PDG/taint rows are written
  - [ ] 01.1.d Run to confirm pass — `Run: go test ./skillgrid-cli/internal/mnemonic/pdg/... ./skillgrid-cli/internal/mnemonic/codeindex/... ./skillgrid-cli/internal/mnemonic/store/... ./skillgrid-cli/internal/mnemonic/mcp/... -run OptInIsolation -count=1` — Expected: PASS
  - [ ] 01.1.e Commit — `feat(mnemonic): opt-in --pdg gate with byte-for-byte unchanged common index`
- [ ] 01.2 `[RED]` Opt-in isolation — per-function CFG + PDG populates only under `--pdg` (Scenario: Opt-in index builds per-function CFG and PDG) — threat: Opt-in isolation
  - [ ] 01.2.a Write failing test — index a fixture repo (a 005-supported language) with `--pdg`; assert `cfg_blocks`/`cfg_edges` rows exist per function (basic blocks from branch/loop/return structure) and `pdg_edges` rows exist for both control-dependence and data-dependence; assert a non-`--pdg` index of the same repo leaves those tables empty
  - [ ] 01.2.b Run to confirm fail — `Run: go test ./skillgrid-cli/internal/mnemonic/pdg/... ./skillgrid-cli/internal/mnemonic/codeindex/... ./skillgrid-cli/internal/mnemonic/store/... -run CfgPdgBuild -count=1` — Expected: FAIL
  - [ ] 01.2.c Minimal implementation — `014_pdg_taint.sql` migration (cfg_blocks, cfg_edges, pdg_edges, taint_findings); per-function CFG builder over the gotreesitter AST (basic blocks from branch/loop/return structure, no new grammar); control-dependence + data-dependence derivation into `pdg_edges`; indexer hook after 005 extraction
  - [ ] 01.2.d Run to confirm pass — `Run: go test ./skillgrid-cli/internal/mnemonic/pdg/... ./skillgrid-cli/internal/mnemonic/codeindex/... ./skillgrid-cli/internal/mnemonic/store/... -run CfgPdgBuild -count=1` — Expected: PASS
  - [ ] 01.2.e Commit — `feat(mnemonic): per-function CFG and control/data-dependence PDG`
- [ ] 01.3 `[RED]` Deterministic reproducibility — same input yields identical PDG across repeated runs (Scenario: Repeated PDG builds are reproducible)
  - [ ] 01.3.a Write failing test — index the same fixture repo with `--pdg` twice (fresh store each time); assert the full set of `cfg_blocks`/`cfg_edges`/`pdg_edges` rows (ids, types, confidence labels, ordering-independent) is byte-for-byte identical across the two runs; assert no wall-clock, pointer, or map-iteration-order nondeterminism leaks into persisted rows
  - [ ] 01.3.b Run to confirm fail — `Run: go test ./skillgrid-cli/internal/mnemonic/pdg/... ./skillgrid-cli/internal/mnemonic/store/... -run PdgReproducible -count=1` — Expected: FAIL
  - [ ] 01.3.c Minimal implementation — deterministic traversal (sorted iteration where the AST/PG yields unordered collections); stable block/statement ids keyed by (symbol, line range), not by allocation order
  - [ ] 01.3.d Run to confirm pass — `Run: go test ./skillgrid-cli/internal/mnemonic/pdg/... ./skillgrid-cli/internal/mnemonic/store/... -run PdgReproducible -count=1` — Expected: PASS
  - [ ] 01.3.e Commit — `feat(mnemonic): deterministic per-function PDG construction`
- [ ] 01.4 `[RED]` Confidence labels — every PDG edge is labeled; unresolved data-dependence is `AMBIGUOUS`/`INFERRED` not fabricated (Scenario: Every PDG edge carries a confidence label)
  - [ ] 01.4.a Write failing test — (1) assert every row in `pdg_edges` has a non-empty Confidence Label in `EXTRACTED | INFERRED | AMBIGUOUS`; (2) fixture with a resolvable data-dependence → that edge is `EXTRACTED`; (3) fixture with a data-dependence that cannot be resolved (e.g. value leaves through an unresolved call) → that edge is `AMBIGUOUS` or `INFERRED`, and no fabricated `EXTRACTED` edge is written for it
  - [ ] 01.4.b Run to confirm fail — `Run: go test ./skillgrid-cli/internal/mnemonic/pdg/... -run PdgConfidenceLabels -count=1` — Expected: FAIL
  - [ ] 01.4.c Minimal implementation — Confidence Label column on `pdg_edges`; derivation logic labels each edge by how it was resolved; only fully-resolved data-dependences are `EXTRACTED`
  - [ ] 01.4.d Run to confirm pass — `Run: go test ./skillgrid-cli/internal/mnemonic/pdg/... -run PdgConfidenceLabels -count=1` — Expected: PASS
  - [ ] 01.4.e Commit — `feat(mnemonic): confidence-labeled PDG edges`
- [ ] 01.5 `[AFK]` Malformed function CFG is skipped and index continues (Scenario: Malformed function CFG skips and index continues) — `Run: go test ./skillgrid-cli/internal/mnemonic/pdg/... ./skillgrid-cli/internal/mnemonic/codeindex/... -count=1` — Expected: PASS
- [ ] 01.6 `[AFK]` Function over depth/step cap truncates with a note, never aborts (Scenario: Over-cap function truncates with a note) — `Run: go test ./skillgrid-cli/internal/mnemonic/pdg/... -count=1` — Expected: PASS
- [ ] 01.7 `[AFK]` CGo-free: build succeeds without CGO_ENABLED=1 (Scenario: PDG package builds without CGo) — `Run: CGO_ENABLED=0 go build ./skillgrid-cli/internal/mnemonic/pdg/... ./skillgrid-cli/internal/mnemonic/codeindex/... && CGO_ENABLED=0 go test ./skillgrid-cli/internal/mnemonic/pdg/... -count=1` — Expected: PASS
- [ ] 01.8 `[RED]` Mnemonic tool surface — `code_pdg_query` registered + existing 005/008/010 tools stable + bad args rejected (Scenario: code_pdg_query registered and bad args fail) — threat: Mnemonic tool surface
  - [ ] 01.8.a Write failing test — (1) assert `code_pdg_query` is registered with distinct name + required `symbol`/`statement` params; (2) assert the full set of 005/008/010 `code_*` tool names + required params is unchanged; (3) assert `code_pdg_query` with missing/unknown args is rejected with a clear validation error (abort, not a fabricated empty result)
  - [ ] 01.8.b Run to confirm fail — `Run: go test ./skillgrid-cli/internal/mnemonic/mcp/... -run PdgQueryTool -count=1` — Expected: FAIL
  - [ ] 01.8.c Minimal implementation — `tools_code_pdg.go` `code_pdg_query` tool + server registration without dropping existing `code_*`; arg validation
  - [ ] 01.8.d Run to confirm pass — `Run: go test ./skillgrid-cli/internal/mnemonic/mcp/... -run PdgQueryTool -count=1` — Expected: PASS
  - [ ] 01.8.e Commit — `feat(mnemonic): register code_pdg_query tool with stable 005/008/010 surface`
- [ ] 01.9 `[AFK]` `code_pdg_query` returns control/data dependents of a statement (Scenario: PDG query returns statement dependents) — `Run: go test ./skillgrid-cli/internal/mnemonic/mcp/... ./skillgrid-cli/internal/mnemonic/service/... -count=1` — Expected: PASS
- [ ] 01.10 `[AFK]` Unknown symbol or statement returns not-found, no fabricated dependences (Scenario: Unknown PDG query symbol returns not-found) — `Run: go test ./skillgrid-cli/internal/mnemonic/mcp/... -count=1` — Expected: PASS
- [ ] 01.11 `[AFK]` `skillgrid index --pdg` CLI flag parity + `skillgrid search pdg` (Scenario: CLI pdg flag and search parity) — `Run: go test ./skillgrid-cli/cmd/skillgrid/... -count=1` — Expected: PASS

### Verification

Verdict: `PENDING`  <!-- PASS | PASS WITH WARNINGS | FAIL -->

Evidence:

| Check | Run | Expected | Result | Notes |
|-------|-----|----------|--------|-------|
| Focused test | `go test ./skillgrid-cli/internal/mnemonic/pdg/... ./skillgrid-cli/internal/mnemonic/store/... -count=1` | PASS | | |
| Focused test (codeindex + mcp) | `go test ./skillgrid-cli/internal/mnemonic/codeindex/... ./skillgrid-cli/internal/mnemonic/mcp/... -count=1` | PASS | | |
| Acceptance `@step-01` / `@p0` | BDD / mapped unit scenarios | PASS | | |
| Runtime harness | `skillgrid index --pdg` on fixture repo; query `cfg_blocks`/`pdg_edges` via `code_pdg_query` | PASS | | |
| Byte-for-byte gate | `skillgrid index` (no `--pdg`) on fixture repo; diff 005/008/010 tool outputs against pre-011 baseline | PASS | | |
| Rollback boundary | Drop `014_pdg_taint.sql` + `pdg/` + `code_pdg_query` + `--pdg` hook; re-run non-`--pdg` index | PASS | | |
| Global Constraints | — | held | | |

### Commit

When step DoD is met: `feat(mnemonic): opt-in per-function CFG and control/data-dependence PDG`

---

## 02-taint-solver

### Goal

Configurable source/sink sets + intraprocedural source→sink taint so "does untrusted input reach this sink" is a query.

### Out of scope / Non-Goals

- Interprocedural flow (later pass); a full data-flow engine; changing 005/008/010 tools
- LLM-suggested sources/sinks; cross-function call-boundary summaries

### Definition of Done

This step is done only when:

- [ ] All `### Tasks` checkboxes below are `[x]`
- [ ] All `@step-02` scenarios in `acceptance.feature` pass
- [ ] `### Verification` Verdict is `PASS` or `PASS WITH WARNINGS`
- [ ] Depends-on step(s) already PASS / PASS WITH WARNINGS
- [ ] No Global Constraint violated

> Depends on: 01-cfg-pdg

**Files:**
- Create: `skillgrid-cli/internal/mnemonic/pdg/taint.go`
- Create: `skillgrid-cli/internal/mnemonic/mcp/tools_code_pdg.go` (extend with `code_taint`)
- Modify: `skillgrid-cli/internal/mnemonic/mcp/server.go`
- Modify: `skillgrid-cli/cmd/skillgrid/code_intel.go`
- Modify: `skillgrid-cli/cmd/skillgrid/main.go`
- Test: `skillgrid-cli/internal/mnemonic/pdg/...`, `skillgrid-cli/internal/mnemonic/mcp/...`, `skillgrid-cli/cmd/skillgrid/...`

**Interfaces:**
- Consumes: PDG edges + CFG blocks from step 01; 005's symbols/edges; existing store + indexer
- Produces: deterministic configurable source/sink sets; intraprocedural source→sink solver over the PDG; persisted `taint_findings` (source kind, sink kind, hop-by-hop path, per-hop Confidence Label); `code_taint` MCP/CLI tool (`--symbol`/`--file`/`--json`); `skillgrid search taint` CLI; findings registered on the MCP surface

### Tasks

- [ ] 02.1 `[RED]` Opt-in isolation — `--pdg` index adds taint findings without altering 005/008/010 results (Scenario: Opt-in taint index leaves 005/008/010 results unchanged) — threat: Opt-in isolation
  - [ ] 02.1.a Write failing test — index a fixture repo with `--pdg` that has at least one known source→sink path; assert (1) `taint_findings` rows exist for the known flow; (2) every 005/008/010 `code_*` tool output is byte-for-byte identical to the pre-011 baseline for the same repo; (3) a non-`--pdg` index of the same repo leaves `taint_findings` empty
  - [ ] 02.1.b Run to confirm fail — `Run: go test ./skillgrid-cli/internal/mnemonic/pdg/... ./skillgrid-cli/internal/mnemonic/mcp/... -run TaintOptInIsolation -count=1` — Expected: FAIL
  - [ ] 02.1.c Minimal implementation — taint solver wired into the `--pdg` hook after PDG derivation; findings written to `taint_findings` only when `--pdg` is set
  - [ ] 02.1.d Run to confirm pass — `Run: go test ./skillgrid-cli/internal/mnemonic/pdg/... ./skillgrid-cli/internal/mnemonic/mcp/... -run TaintOptInIsolation -count=1` — Expected: PASS
  - [ ] 02.1.e Commit — `feat(mnemonic): opt-in taint findings without altering 005/008/010 results`
- [ ] 02.2 `[RED]` Taint core — source→sink path found and persisted (Scenario: Source to sink taint path found)
  - [ ] 02.2.a Write failing test — fixture with a known source (e.g. request param) that flows to a known sink (e.g. SQL exec) through a resolvable data-dependence chain; assert `code_taint` returns a finding with the source kind, sink kind, and the hop-by-hop path; assert the finding is persisted in `taint_findings`
  - [ ] 02.2.b Run to confirm fail — `Run: go test ./skillgrid-cli/internal/mnemonic/pdg/... ./skillgrid-cli/internal/mnemonic/mcp/... -run TaintSourceToSink -count=1` — Expected: FAIL
  - [ ] 02.2.c Minimal implementation — `pdg/taint.go` default deterministic source/sink sets (sources: request params, env vars, file reads; sinks: SQL exec, shell exec, template render, file write); source→sink reachability solver over the PDG data-dependence edges; `code_taint` tool
  - [ ] 02.2.d Run to confirm pass — `Run: go test ./skillgrid-cli/internal/mnemonic/pdg/... ./skillgrid-cli/internal/mnemonic/mcp/... -run TaintSourceToSink -count=1` — Expected: PASS
  - [ ] 02.2.e Commit — `feat(mnemonic): intraprocedural source-to-sink taint solver`
- [ ] 02.3 `[RED]` Boundary-not-fabricated — a path ending at an unresolved call boundary is `AMBIGUOUS`/truncated with "stops at" note, not fabricated (Scenario: Taint path stops at an unresolved boundary)
  - [ ] 02.3.a Write failing test — (1) fixture where a source flows to a sink only through an unresolved call boundary (callee not resolvable) → assert the finding is reported with the path truncated at the boundary, the boundary hop marked `AMBIGUOUS`, and a "stops at <boundary>" note; (2) fixture where a source has no path to any sink → assert **no** finding is produced (never fabricated); (3) assert no finding is ever reported with an empty or fabricated hop list
  - [ ] 02.3.b Run to confirm fail — `Run: go test ./skillgrid-cli/internal/mnemonic/pdg/... -run TaintBoundaryNotFabricated -count=1` — Expected: FAIL
  - [ ] 02.3.c Minimal implementation — solver truncates the path at the first unresolved boundary; marks the boundary hop `AMBIGUOUS`; attaches a "stops at <boundary>" note; suppresses findings with no path
  - [ ] 02.3.d Run to confirm pass — `Run: go test ./skillgrid-cli/internal/mnemonic/pdg/... -run TaintBoundaryNotFabricated -count=1` — Expected: PASS
  - [ ] 02.3.e Commit — `feat(mnemonic): taint boundary truncation with stops-at note, no fabricated findings`
- [ ] 02.4 `[RED]` Confidence labels — every taint edge is confidence-labeled; path is `EXTRACTED` only when every hop is resolved (Scenario: Every taint hop carries a confidence label)
  - [ ] 02.4.a Write failing test — (1) assert every hop in every `taint_findings` path has a non-empty Confidence Label in `EXTRACTED | INFERRED | AMBIGUOUS`; (2) fixture where every hop is a resolved data-dependence → path is `EXTRACTED`; (3) fixture where at least one hop is unresolved → path is `INFERRED` or `AMBIGUOUS`, never `EXTRACTED`
  - [ ] 02.4.b Run to confirm fail — `Run: go test ./skillgrid-cli/internal/mnemonic/pdg/... -run TaintConfidenceLabels -count=1` — Expected: FAIL
  - [ ] 02.4.c Minimal implementation — per-hop Confidence Label on the path; path-level label derived from the worst hop (`EXTRACTED` only when all hops `EXTRACTED`)
  - [ ] 02.4.d Run to confirm pass — `Run: go test ./skillgrid-cli/internal/mnemonic/pdg/... -run TaintConfidenceLabels -count=1` — Expected: PASS
  - [ ] 02.4.e Commit — `feat(mnemonic): confidence-labeled taint findings`
- [ ] 02.5 `[RED]` Deterministic reproducibility — same source/sink config yields identical findings across repeated runs (Scenario: Repeated taint runs are reproducible)
  - [ ] 02.5.a Write failing test — index the same fixture repo with `--pdg` twice (fresh store each time, same source/sink config); assert the full set of `taint_findings` (source kind, sink kind, path, per-hop labels, ordering-independent) is byte-for-byte identical across the two runs; assert no nondeterminism leaks into findings
  - [ ] 02.5.b Run to confirm fail — `Run: go test ./skillgrid-cli/internal/mnemonic/pdg/... -run TaintReproducible -count=1` — Expected: FAIL
  - [ ] 02.5.c Minimal implementation — deterministic source/sink matching + path enumeration (sorted traversal, stable ordering); findings keyed by (source, sink, path) not by run order
  - [ ] 02.5.d Run to confirm pass — `Run: go test ./skillgrid-cli/internal/mnemonic/pdg/... -run TaintReproducible -count=1` — Expected: PASS
  - [ ] 02.5.e Commit — `feat(mnemonic): deterministic taint findings`
- [ ] 02.6 `[RED]` Mnemonic tool surface — `code_taint` registered + 005/008/010 tools stable + bad args rejected (Scenario: code_taint registered and bad args fail) — threat: Mnemonic tool surface
  - [ ] 02.6.a Write failing test — (1) assert `code_taint` is registered with distinct name + optional `--symbol`/`--file`/`--json` params; (2) assert the full set of 005/008/010 `code_*` tool names + required params is unchanged; (3) assert `code_taint` with bad/missing args is rejected with a clear validation error (abort, not invented findings)
  - [ ] 02.6.b Run to confirm fail — `Run: go test ./skillgrid-cli/internal/mnemonic/mcp/... -run TaintTool -count=1` — Expected: FAIL
  - [ ] 02.6.c Minimal implementation — `code_taint` tool + server registration without dropping existing `code_*`; arg validation
  - [ ] 02.6.d Run to confirm pass — `Run: go test ./skillgrid-cli/internal/mnemonic/mcp/... -run TaintTool -count=1` — Expected: PASS
  - [ ] 02.6.e Commit — `feat(mnemonic): register code_taint tool with stable 005/008/010 surface`
- [ ] 02.7 `[AFK]` `--symbol` / `--file` filter findings; `--json` for CI (Scenario: Taint findings filter by symbol and file) — `Run: go test ./skillgrid-cli/internal/mnemonic/mcp/... -count=1` — Expected: PASS
- [ ] 02.8 `[AFK]` Non-`--pdg` index returns a clear "run `--pdg`" message, not an error (Scenario: Non-pdg taint query returns run-pdg hint) — `Run: go test ./skillgrid-cli/internal/mnemonic/mcp/... -count=1` — Expected: PASS
- [ ] 02.9 `[AFK]` Source/sink sets are deterministic and configurable (Scenario: Source and sink sets are configurable) — `Run: go test ./skillgrid-cli/internal/mnemonic/pdg/... -count=1` — Expected: PASS
- [ ] 02.10 `[AFK]` `skillgrid search taint` CLI parity + `skillgrid index --pdg` (Scenario: CLI taint search and pdg flag parity) — `Run: go test ./skillgrid-cli/cmd/skillgrid/... -count=1` — Expected: PASS

### Verification

Verdict: `PENDING`  <!-- PASS | PASS WITH WARNINGS | FAIL -->

Evidence:

| Check | Run | Expected | Result | Notes |
|-------|-----|----------|--------|-------|
| Focused test | `go test ./skillgrid-cli/internal/mnemonic/pdg/... ./skillgrid-cli/internal/mnemonic/mcp/... -count=1` | PASS | | |
| Focused test (CLI) | `go test ./skillgrid-cli/cmd/skillgrid/... -count=1` | PASS | | |
| Acceptance `@step-02` / `@p0` | BDD / mapped unit scenarios | PASS | | |
| Runtime harness | `skillgrid index --pdg` on fixture with known source→sink; `code_taint` returns the finding; cross-boundary finding is `AMBIGUOUS` + "stops at" | PASS | | |
| Byte-for-byte gate | `skillgrid index --pdg` on fixture; diff 005/008/010 tool outputs against pre-011 baseline | PASS | | |
| Rollback boundary | Drop `pdg/taint.go` + `code_taint` + `taint_findings` writes; re-run index | PASS | | |
| Global Constraints | — | held | | |

### Commit

When step DoD is met: `feat(mnemonic): intraprocedural source-to-sink taint with code_taint`

---

## Archive gate checklist

- [ ] Change-level **Definition of Done** fully checked
- [ ] No unchecked `- [ ]` under any `### Tasks`
- [ ] Every step Verdict is `PASS` or `PASS WITH WARNINGS`
- [ ] No Global Constraint violated
- [ ] `## State` status is `done` and phase is `archive` (set by verify/archive)
- [ ] STATUS banner updated to `complete`
