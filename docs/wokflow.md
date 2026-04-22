# Skillgrid workflow

Full indexes: [skills.md](skills.md), [commands.md](commands.md), [tools.md](tools.md).

```bash
/skillgrid-init
## Create folder structure
## graphify init
## CocoIndex init
### opencode init
### Skills (.agents/skills/)
- sdd-init — bootstrap SDD context (stack, conventions, persistence)
- openspec-onboard — first full OpenSpec cycle (requires OpenSpec CLI)
- search-first — research tools and patterns before building
- karpathy-guidelines — assumptions, simplicity, surgical edits (baseline)
- skill-creator — add or extend Agent Skills in this repo
- ccc — CocoIndex Code: `ccc init`, indexing, semantic search setup
- context-engineering — rules, context packing, MCP usage (from addyosmani/agent-skills)
- using-agent-skills — meta: how to use the agent-skills pack

/skillgrid-explore
## openspec-explore
## Create AGENTS.md
## Create PROJECT.md
## Skill: documentation - Architecture Decision Records, API docs, inline documentation standards - document the why
### Skills (.agents/skills/)
- openspec-explore — explore problem and codebase before a change
- search-first — research external options when discovery matters
- deep-research — deeper investigation when exploration needs breadth
- ccc — semantic codebase search (`ccc search`)

/skillgrid-brainstorm
## Skill: idea-refine -
### Skills (.agents/skills/)
- karpathy-guidelines — surface tradeoffs and alternatives before locking direction
- deep-research — validate assumptions and gather evidence
- idea-refine — divergent/convergent thinking to sharpen a vague idea

/skillgrid-plan
## Create PRDs
## Create openspec Change from PRD - openspec-propose
## Skill: spec-driven-development - Write a PRD covering objectives, commands, structure, code style, testing, and boundaries before any code
### Skills (.agents/skills/)
- openspec-propose — OpenSpec change + CLI-driven artifacts to apply-ready
- sdd-propose — proposal.md only (SDD orchestrator / Engram or openspec mode)
- spec-driven-development — PRD-style scope before code (objectives, style, testing, boundaries)

/skillgrid-design
# Plan UI
## Skill: nanobanana
## Create DESIGN.md
## Skill: frontend-ui-engineering
## Skill: api-and-interface-design
## Create ARCHITECTURE.md
## Create and Update PRDs
## Update openspec Change if scope shifts - openspec-propose / sdd-propose
## Write delta specs (requirements/scenarios) - sdd-spec
### Skills (.agents/skills/)
- sdd-spec — delta specs: requirements and scenarios for the change
- database-design — schema and data modeling before migrations
- openspec-propose — refresh proposal when scope or intent changes
- sdd-propose — same, under SDD artifact layout
- frontend-ui-engineering — components, design systems, a11y (WCAG-oriented)
- api-and-interface-design — contract-first APIs, boundaries, error semantics

/skillgrid-breakdown
## Add tasks to PRDs
## Create Tasks for openspec Change - sdd-tasks
## Skill: task-breakdown - Decompose specs into small, verifiable tasks with acceptance criteria and dependency ordering
### Skills (.agents/skills/)
- sdd-tasks — tasks.md checklist from proposal, specs, and design
- planning-and-task-breakdown — atomic tasks, acceptance criteria, dependencies

/skillgrid-apply
## using-git-worktrees
## Apply openspec Change - openspec-apply-change
## Go for Minimum Viable Change
## Skill: test-driven-development - Red-Green-Refactor
## Skill: source-driven-development - Ground every framework decision in official documentation
### Skills (.agents/skills/)
- openspec-apply-change — implement from OpenSpec change tasks
- karpathy-guidelines — surgical, verifiable steps
- incremental-implementation — vertical slices, verify, commit; flags and safe defaults
- tdd-workflow — structured red-green-refactor discipline
- tdd-guide — TDD guidance and patterns
- test-driven-development — red-green-refactor, pyramid, DAMP, Beyonce Rule
- source-driven-development — cite official docs for framework decisions
- clean-code — readability and maintainability while implementing
- database-migrations — apply schema/data migrations safely
- ccc — refresh semantic index after significant code changes

/skillgrid-test
## Skill: browser-testing-with-devtools - Chrome DevTools MCP for live runtime data - DOM inspection, console logs, network traces, performance profiling
### Skills (.agents/skills/)
- e2e-testing — end-to-end test design and implementation
- e2e-runner — run and troubleshoot E2E suites
- testing-patterns — general testing patterns beyond E2E
- browser-testing-with-devtools — DOM, console, network, performance via DevTools MCP
- debugging-and-error-recovery — reproduce, localize, reduce, fix, guard

/skillgrid-validate
## openspec-verify-change
## Skill: code-review - Review code for correctness, style, and adherence to specs
## Skill: security-review - Review code for security vulnerabilities and best practices
### Use trivy
## Skill: document - Architecture Decision Records, API docs, inline documentation standards - document the why
### Skills (.agents/skills/)
- openspec-verify-change — implementation matches change artifacts
- sdd-verify — SDD verification against specs, design, and tasks
- karpathy-guidelines — success criteria and scope discipline
- clean-code — review for clarity and coupling
- code-review-and-quality — five-axis review, change sizing, severity labels
- code-simplification — simplify without changing behavior (Chesterton\'s Fence, Rule of 500)
- security-review — code-focused security review
- security-and-hardening — OWASP-oriented patterns, auth, secrets, boundaries
- performance-optimization — measure first, Web Vitals, profiling, bundle checks
- semgrep-security — static analysis with Semgrep
- trivy-security — container/dep scanning with Trivy
- vulnerability-scanner — threat modeling and vuln prioritization (OWASP-oriented)
- security-scan — audit agent/IDE config (e.g. `.claude/`, MCP, hooks)
- database-reviewer — PostgreSQL schema, SQL, RLS, performance review
- documentation-and-adrs — ADRs, API docs, document the why

/skillgrid-finish
## Archive change - openspec-archive-change
## Optional: sync delta specs to main specs - openspec-sync-specs
## Create PR
### Skills (.agents/skills/)
- openspec-archive-change — complete and archive the change (per your merge process)
- openspec-sync-specs — promote delta specs without archiving, if your flow needs it
- git-workflow-and-versioning — trunk-style workflow, atomic commits, small changes
- ci-cd-and-automation — pipelines, feature flags, quality gates
- deprecation-and-migration — sunset paths, migrations, zombie code
- shipping-and-launch — pre-launch checks, rollouts, rollback, monitoring
```
