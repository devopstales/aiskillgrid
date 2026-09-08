# Source: docs/skillgrid/changes/008-mnemonic-community-knowledge-graph/change.md
# Template: .agents/skills/_shared/templates/template-acceptance.feature
# Trace: change.md ## Goal + ## Definition of Done; tasks.md @step-NN verify lines.
# Mapping: @p0 scenarios ↔ change.md DoD / Testing strategy; @p1 = important failure paths.
# Threat: Mnemonic tool surface — owning steps 01, 02, 03 (005 tools must stay stable in each)
# One Feature per step; tag each Feature with @step-NN matching tasks.md.
# WHAT not HOW: no file paths or function names.

@step-01
Feature: Leiden community detection with god nodes and community tools
  As a coding agent
  I want subsystem-level orientation from clustered communities and hub symbols
  So that I know the core modules and the most-connected concepts before editing

  @happy @p0
  Scenario: code_communities returns labeled subsystems and 005 tools stay stable
    Given an indexed project whose symbols and edges form distinct subsystems
    When the community detection pass runs and the agent lists communities
    Then clustered subsystems are returned, each with an LLM-free label
    And the existing 005 code search tools keep their names and required parameters

  @happy @p0
  Scenario: code_god_nodes ranks hubs and exclude-hubs suppresses them
    Given an indexed project with a few very highly-connected utility symbols
    When the agent requests god nodes with hubs excluded
    Then the most-connected symbols are ranked by degree
    And the utility super-hubs are suppressed from the ranking

  @edge
  Scenario: code_explain_community explains a subsystem
    Given a community returned by community detection
    When the agent requests an explanation for that community
    Then the community's member symbols and its key entry points are returned

  @edge
  Scenario: Tiny graph yields a single trivial community
    Given a graph with fewer than two connected symbols
    When the community detection pass runs
    Then a single trivial community is returned
    And the pass does not crash

  @edge
  Scenario: Community label falls back when no god node exists
    Given a community whose members have no dominant hub symbol
    When a community label is resolved for it
    Then a neutral community-N style label is returned
    And no label is fabricated from absent names

  @edge
  Scenario: Community partition is reproducible and cached by content-hash
    Given an indexed project with a stable graph
    When the community pass runs twice over an unchanged graph
    Then the partition is reproducible under the pinned seed and resolution
    And the result is cached by content-hash so an unchanged re-index is not recomputed

  @failure @p1
  Scenario: Community tools reject bad args clearly
    Given the community tools are registered on the Mnemonic tool surface
    When the agent calls a community tool with a missing or unknown community id
    Then a clear validation error is returned
    And no community is invented

@step-02
Feature: Precomputed process flows from entry points through call chains
  As a coding agent
  I want to see what a subsystem does end to end
  So that I can answer which flow a symbol participates in without reading every file

  @happy @p0
  Scenario: code_processes returns precomputed flows from entry points
    Given an indexed project with entry points for routes handlers and CLI mains
    When the process pass runs and the agent lists processes
    Then complete execution flows are returned in a single call with no per-query traversal
    And each flow has named steps and a cross-community flag

  @happy @p0
  Scenario: code_process returns the full step-by-step trace
    Given a named process from the process layer
    When the agent requests that process by name
    Then the full step-by-step trace is returned
    And each hop carries a confidence label

  @edge
  Scenario: Process labels are cached by content-hash and re-labeled only on change
    Given a process flow that has already been labeled
    When the index re-runs over an unchanged flow
    Then the cached label is reused and the LLM is not called again
    When the flow structure changes
    Then a new label is produced for the new content-hash

  @edge
  Scenario: Cross-community process is flagged
    Given a flow that passes through symbols in more than one community
    When the process pass traces it
    Then the process carries a cross-community flag

  @edge
  Scenario: Trace stops at a dispatch boundary with a note
    Given a call chain that reaches an interface to implementation or message bus or callback boundary
    When the process trace reaches that boundary
    Then the trace is truncated with a stops-at note naming the symbol and reason
    And the trace is not silently cut

  @edge
  Scenario: code_explain_symbol surfaces process participation and 005 tools stay stable
    Given a symbol that participates in a precomputed process
    When the agent explains that symbol
    Then the processes the symbol participates in are surfaced with its step position
    And the existing 005 code tools keep their names and required parameters

  @edge
  Scenario: Untraceable entry point yields single-step or is skipped
    Given an entry point with no traceable call chain
    When the process pass seeds it
    Then it yields a single-step process or is skipped
    And no flow is fabricated

  @failure @p1
  Scenario: LLM down caches the flow unlabeled
    Given a process flow to label and the LLM is unavailable
    When the process pass labels it
    Then the flow structure is still cached
    And it is cached without a fabricated label

  @failure @p1
  Scenario: Process tools reject bad args clearly
    Given the process tools are registered on the Mnemonic tool surface
    When the agent requests an unknown process name with bad args
    Then a clear validation error is returned
    And no process is invented

@step-03
Feature: Knowledge-graph nodes for docs, configs, and SQL
  As a coding agent
  I want code linked to its docs, configs, and data schema in one graph
  So that I can see which code reads or writes a table and trace code to its documentation

  @happy @p0
  Scenario: Markdown links and wikilinks become references edges
    Given indexed markdown docs containing relative links and wikilinks
    When the knowledge extraction pass runs
    Then doc nodes are created with references edges between them
    And each new edge carries a confidence label

  @happy @p0
  Scenario: Config references become configures edges
    Given indexed config files in yaml toml and json
    When the knowledge extraction pass runs
    Then config nodes are created with configures edges to the code they configure
    And each new edge carries a confidence label

  @happy @p0
  Scenario: SQL DDL becomes table and column nodes with reads and writes
    Given indexed SQL schema defining tables and columns
    When the knowledge extraction pass runs
    Then table and column nodes are created
    And the code that references them gets reads and writes edges with confidence labels

  @edge
  Scenario: code_path traces code to doc to config to table
    Given a graph that now spans code, docs, configs, and tables
    When the agent requests a path from a code symbol to a data table
    Then a single query traces through doc, config, and table nodes
    And the existing 005 path tool is unchanged

  @edge
  Scenario: Unresolvable config ref is ambiguous not dropped
    Given a config reference that resolves to no known symbol
    When the knowledge extraction pass runs
    Then the edge is kept and marked ambiguous
    And it is not silently dropped

  @edge
  Scenario: Malformed doc file falls back and indexes the rest
    Given one doc file with unparseable links among valid docs
    When the knowledge extraction pass runs
    Then the bad links are skipped and the rest is indexed
    And the index does not abort

  @edge
  Scenario: Malformed SQL statement is skipped and the rest is indexed
    Given one SQL file with a statement that fails to parse among valid DDL
    When the knowledge extraction pass runs
    Then that statement is skipped and the rest is indexed
    And the index does not abort

  @edge
  Scenario: Indexer hook runs community, process, and knowledge in one transaction
    Given an index run over a project
    When the indexer finishes code extraction
    Then the community, process, and knowledge passes run in the same incremental transaction

  @failure @p1
  Scenario: Knowledge tools register and reject bad args
    Given the knowledge tools are registered on the Mnemonic tool surface
    When the agent calls a knowledge tool with missing args
    Then a clear validation error is returned
    And the existing 005 code tools keep their names and required parameters
