# Source: docs/skillgrid/changes/011-mnemonic-pdg-taint/change.md
# Template: .agents/skills/_shared/templates/template-acceptance.feature
# Trace: change.md ## Goal + ## Definition of Done; tasks.md @step-NN verify lines.
# Mapping: @p0 scenarios ↔ change.md DoD / Testing strategy; @p1 = important failure paths.
# Threat: Mnemonic tool surface — owning steps 01, 02; Opt-in isolation — owning steps 01, 02
# One Feature per step; tag each Feature with @step-NN matching tasks.md.

@step-01
Feature: Opt-in per-function CFG and control/data-dependence PDG
  As a coding agent
  I want statement-level control and data dependence on demand
  So that I can answer "which statements depend on this variable" without re-deriving the CFG by hand

  @happy @p0
  Scenario: Opt-in index builds per-function CFG and PDG
    Given a project with functions in a supported language
    When the code index runs with the opt-in PDG flag
    Then per-function basic blocks and control-flow edges are queryable
    And control-dependence and data-dependence edges are queryable at statement level
    And a non-opt-in index of the same project leaves the PDG tables empty

  @edge
  Scenario: Non-opt-in index is byte-for-byte unchanged
    Given a project indexed without the opt-in PDG flag
    When the code index runs
    Then the PDG and taint tables are created but empty
    And every pre-existing code tool output is byte-for-byte identical to before this change
    And a PDG query on this index returns a clear run-PDG hint and an empty result, not an error

  @edge
  Scenario: Every PDG edge carries a confidence label
    Given an opt-in indexed project with both resolvable and unresolvable data dependences
    When a PDG query inspects the edges
    Then every edge carries a confidence label
    And a resolvable data dependence is labeled EXTRACTED
    And an unresolvable data dependence is labeled AMBIGUOUS or INFERRED, not fabricated as EXTRACTED

  @edge
  Scenario: Malformed function CFG skips and index continues
    Given one function whose CFG cannot be parsed
    When the code index runs with the opt-in PDG flag
    Then that function is skipped with a warning
    And the index completes for the remaining functions

  @edge
  Scenario: Over-cap function truncates with a note
    Given one function whose PDG exceeds the depth or step cap
    When the code index runs with the opt-in PDG flag
    Then that function is truncated with a stops-at note
    And the index never aborts

  @edge
  Scenario: Repeated PDG builds are reproducible
    Given the same project indexed twice with the opt-in PDG flag into fresh stores
    When the PDG is compared across the two runs
    Then the full set of basic blocks and PDG edges is byte-for-byte identical

  @failure @p1 @security
  Scenario: code_pdg_query registered and bad args fail
    Given the Mnemonic tool surface is registered
    When an agent inspects the code tool surface and calls PDG query with missing or unknown args
    Then PDG query is registered with a distinct name and required symbol and statement params
    And every pre-existing code tool name and required param is unchanged
    And the bad-args call is rejected clearly without an invented empty result

@step-02
Feature: Intraprocedural source-to-sink taint with confidence-labeled findings
  As a coding agent
  I want to ask "does untrusted input reach this sink" as a query
  So that I can judge data-flow security before edits without reading files

  @happy @p0
  Scenario: Source to sink taint path found
    Given an opt-in indexed project with a known untrusted source flowing to a known sink through a resolvable chain
    When an agent queries taint findings
    Then a finding is returned with the source kind, sink kind, and the hop-by-hop path
    And the finding is persisted

  @happy @p0
  Scenario: Every taint hop carries a confidence label
    Given taint findings for paths with both fully-resolved and partially-resolved hops
    When an agent inspects the findings
    Then every hop carries a confidence label
    And a path whose every hop is a resolved data dependence is EXTRACTED
    And a path with any unresolved hop is INFERRED or AMBIGUOUS, never EXTRACTED

  @happy @p0
  Scenario: Opt-in taint index leaves 005/008/010 results unchanged
    Given a project indexed with the opt-in PDG flag that has at least one known source to sink path
    When an agent queries taint and the pre-existing code tools
    Then taint findings exist for the known flow
    And every pre-existing code tool output is byte-for-byte identical to before this change
    And a non-opt-in index of the same project leaves the taint findings empty

  @edge
  Scenario: Taint path stops at an unresolved boundary
    Given an opt-in indexed project where a source reaches a sink only through an unresolved call boundary
    When an agent queries taint findings
    Then the finding is reported with the path truncated at the boundary
    And the boundary hop is marked AMBIGUOUS
    And the finding carries a stops-at note naming the boundary
    And no finding with an empty or fabricated hop list is ever reported

  @edge
  Scenario: No source to sink path means no finding
    Given an opt-in indexed project where a source has no path to any sink
    When an agent queries taint findings
    Then no finding is produced for that source
    And no fabricated path is reported

  @edge
  Scenario: Repeated taint runs are reproducible
    Given the same project indexed twice with the opt-in PDG flag into fresh stores with the same source and sink config
    When the taint findings are compared across the two runs
    Then the full set of findings is byte-for-byte identical

  @edge
  Scenario: Taint findings filter by symbol and file
    Given an opt-in indexed project with taint findings across multiple symbols and files
    When an agent queries taint findings filtered by symbol or file
    Then only matching findings are returned
    And the machine-readable output is available for CI

  @edge
  Scenario: Non-pdg taint query returns run-pdg hint
    Given a project indexed without the opt-in PDG flag
    When an agent queries taint findings
    Then a clear run-PDG hint is returned
    And the result is empty, not an error

  @edge
  Scenario: Source and sink sets are configurable
    Given the default deterministic source and sink sets
    When the source or sink set is configured differently
    Then taint findings follow the configured set deterministically

  @failure @p1 @security
  Scenario: code_taint registered and bad args fail
    Given the Mnemonic tool surface is registered
    When an agent inspects the code tool surface and calls taint with missing or unknown args
    Then taint is registered with a distinct name and optional symbol, file, and json params
    And every pre-existing code tool name and required param is unchanged
    And the bad-args call is rejected clearly without invented findings
