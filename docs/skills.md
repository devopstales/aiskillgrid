# Agent skills

This project keeps reusable agent instructions under **`.agents/skills/<skill-name>/SKILL.md`**. Tools discover them from `.agents/skills/` (and may also load skills from user-global paths, depending on the IDE).

For **IDE persona files** (`agents/` under each IDE config dir), see [`agents.md`](agents.md). For **how skills fit into the end-to-end flow**, see [`workflow.md`](workflow.md). For **CLIs, MCPs, and install.sh dependencies** referenced by skills, see [`tools.md`](tools.md). For **web search, scraping, doc MCPs, Brave API skills, and Context7**, see [`web-scraping-and-research.md`](web-scraping-and-research.md).

---

## Shared reference files

Some skills point to checklists under **`.agents/skills/references/`** (for example `testing-patterns.md`, `security-checklist.md`, `performance-checklist.md`, `accessibility-checklist.md`, `orchestration-patterns.md`). Open those from the skill that cites them.

---

## Skill index

Summaries are taken from each skill’s frontmatter `description` (trimmed). For full workflows, open the `SKILL.md` file.

### OpenSpec

| Skill | Path | Summary |
|--------|------|---------|
| `openspec-onboard` | [.agents/skills/openspec-onboard/SKILL.md](../.agents/skills/openspec-onboard/SKILL.md) | Guided first full OpenSpec cycle with narration. |
| `openspec-explore` | [.agents/skills/openspec-explore/SKILL.md](../.agents/skills/openspec-explore/SKILL.md) | Explore mode: ideas, investigation, requirements (no implementation stance). |
| `openspec-propose` | [.agents/skills/openspec-propose/SKILL.md](../.agents/skills/openspec-propose/SKILL.md) | New change: generate proposal, design, specs, tasks in one CLI-driven pass. |
| `openspec-apply-change` | [.agents/skills/openspec-apply-change/SKILL.md](../.agents/skills/openspec-apply-change/SKILL.md) | Implement tasks from an OpenSpec change. |
| `openspec-verify-change` | [.agents/skills/openspec-verify-change/SKILL.md](../.agents/skills/openspec-verify-change/SKILL.md) | Verify implementation matches change artifacts. |
| `openspec-sync-specs` | [.agents/skills/openspec-sync-specs/SKILL.md](../.agents/skills/openspec-sync-specs/SKILL.md) | Merge delta specs into main specs without archiving. |
| `openspec-archive-change` | [.agents/skills/openspec-archive-change/SKILL.md](../.agents/skills/openspec-archive-change/SKILL.md) | Archive a completed change. |

### Skillgrid reusable primitives

End-to-end phase orchestration still lives in **slash commands** under [`.cursor/commands/`](../.cursor/commands/) (mirrored to `.kilo/commands/`, `.opencode/commands/`, `.github/prompts/`). Reusable Skillgrid procedures live in the skills below and are referenced by the commands.

| Skill | Path | Summary |
|--------|------|---------|
| `skillgrid-questioning` | [.agents/skills/skillgrid-questioning/SKILL.md](../.agents/skills/skillgrid-questioning/SKILL.md) | Ask only blocking questions and record answers into artifacts. |
| `skillgrid-codebase-map` | [.agents/skills/skillgrid-codebase-map/SKILL.md](../.agents/skills/skillgrid-codebase-map/SKILL.md) | Map repo structure, graphify output, tests, design tokens, and conventions. |
| `skillgrid-parallel-research` | [.agents/skills/skillgrid-parallel-research/SKILL.md](../.agents/skills/skillgrid-parallel-research/SKILL.md) | Coordinate parallel research, docs lookup, web evidence, and research spill files. |
| `skillgrid-subagent-orchestration` | [.agents/skills/skillgrid-subagent-orchestration/SKILL.md](../.agents/skills/skillgrid-subagent-orchestration/SKILL.md) | Dispatch bounded subagents with handoff paths and two-stage review. |
| `skillgrid-import-artifacts` | [.agents/skills/skillgrid-import-artifacts/SKILL.md](../.agents/skills/skillgrid-import-artifacts/SKILL.md) | Import legacy PRDs and OpenSpec changes into canonical `.skillgrid/prd/` and `INDEX.md` artifacts. |
| `skillgrid-prd-artifacts` | [.agents/skills/skillgrid-prd-artifacts/SKILL.md](../.agents/skills/skillgrid-prd-artifacts/SKILL.md) | PRD naming, numbering, `INDEX.md`, title blocks, and lifecycle status. |
| `skillgrid-spec-artifacts` | [.agents/skills/skillgrid-spec-artifacts/SKILL.md](../.agents/skills/skillgrid-spec-artifacts/SKILL.md) | Bridge PRDs to OpenSpec proposal, design, delta specs, tasks, and validation. |
| `skillgrid-vertical-slices` | [.agents/skills/skillgrid-vertical-slices/SKILL.md](../.agents/skills/skillgrid-vertical-slices/SKILL.md) | Split PRDs/specs/tasks into testable, shippable slices with `[HITL]` / `[AFK]` tags. |
| `skillgrid-ui-design-artifacts` | [.agents/skills/skillgrid-ui-design-artifacts/SKILL.md](../.agents/skills/skillgrid-ui-design-artifacts/SKILL.md) | Place UI decisions in `DESIGN.md`, previews, PRDs, and OpenSpec design docs. |
| `skillgrid-issue-creation` | [.agents/skills/skillgrid-issue-creation/SKILL.md](../.agents/skills/skillgrid-issue-creation/SKILL.md) | Local/GitHub/GitLab/Jira issue behavior from `.skillgrid/config.json`. |
| `skillgrid-hybrid-persistence` | [.agents/skills/skillgrid-hybrid-persistence/SKILL.md](../.agents/skills/skillgrid-hybrid-persistence/SKILL.md) | Coordinate on-disk artifacts and Engram memory. |
| `skillgrid-filesystem-handoff` | [.agents/skills/skillgrid-filesystem-handoff/SKILL.md](../.agents/skills/skillgrid-filesystem-handoff/SKILL.md) | Maintain `context_<change-id>.md` and `research/<change-id>/` handoff files. |
| `skillgrid-openspec-config` | [.agents/skills/skillgrid-openspec-config/SKILL.md](../.agents/skills/skillgrid-openspec-config/SKILL.md) | Keep `openspec/config.yaml` aligned with Skillgrid config. |
| `skillgrid-project-docs` | [.agents/skills/skillgrid-project-docs/SKILL.md](../.agents/skills/skillgrid-project-docs/SKILL.md) | Refresh `DESIGN.md` and `.skillgrid/project/*` docs. |
| `skillgrid-checkpoints` | [.agents/skills/skillgrid-checkpoints/SKILL.md](../.agents/skills/skillgrid-checkpoints/SKILL.md) | Manage `.skillgrid/tasks/checkpoints.log`. |

Commands remain the source of truth for phase order, exit status, and final reports. Skills own reusable procedure details.

### Define and plan (product and tasks)

| Skill | Path | Summary |
|--------|------|---------|
| `idea-refine` | [.agents/skills/idea-refine/SKILL.md](../.agents/skills/idea-refine/SKILL.md) | Divergent and convergent thinking to sharpen ideas. |
| `spec-driven-development` | [.agents/skills/spec-driven-development/SKILL.md](../.agents/skills/spec-driven-development/SKILL.md) | Spec before code for new work. |
| `planning-and-task-breakdown` | [.agents/skills/planning-and-task-breakdown/SKILL.md](../.agents/skills/planning-and-task-breakdown/SKILL.md) | Decompose specs into ordered, verifiable tasks. |
| `incremental-implementation` | [.agents/skills/incremental-implementation/SKILL.md](../.agents/skills/incremental-implementation/SKILL.md) | Thin vertical slices: implement, verify, commit. |
| `karpathy-guidelines` | [.agents/skills/karpathy-guidelines/SKILL.md](../.agents/skills/karpathy-guidelines/SKILL.md) | Simplicity, surgical edits, explicit assumptions, verifiable goals. |
| `search-first` | [.agents/skills/search-first/SKILL.md](../.agents/skills/search-first/SKILL.md) | Research existing solutions before building. |
| `deep-research` | [.agents/skills/deep-research/SKILL.md](../.agents/skills/deep-research/SKILL.md) | Multi-source research with citations. |
| `exa-search` | [.agents/skills/exa-search/SKILL.md](../.agents/skills/exa-search/SKILL.md) | Neural search via Exa MCP. |
| `context-engineering` | [.agents/skills/context-engineering/SKILL.md](../.agents/skills/context-engineering/SKILL.md) | Rules, context packing, session quality. |
| `using-agent-skills` | [.agents/skills/using-agent-skills/SKILL.md](../.agents/skills/using-agent-skills/SKILL.md) | Meta: discover and invoke other skills. |

### Brave Search API (`brave-*`)

[Brave Search API](https://api.search.brave.com) via API key and `curl` (and optional **`bx`** CLI). Not shipped as default MCP fragments in [`.configs/mcp/`](../.configs/mcp/); pair with [`web-scraping-and-research.md`](web-scraping-and-research.md) for how this fits Exa / Firecrawl / Context7.

| Skill (directory) | Path | Summary |
|-------------------|------|---------|
| `brave-bx` (`bx`) | [.agents/skills/brave-bx/SKILL.md](../.agents/skills/brave-bx/SKILL.md) | All-in-one agentic search / grounding (`bx` CLI); pre-extracted, token-budgeted content. |
| `brave-web-search` | [.agents/skills/brave-web-search/SKILL.md](../.agents/skills/brave-web-search/SKILL.md) | Primary web search: snippets, URLs, thumbnails, Goggles, pagination. |
| `brave-llm-context` | [.agents/skills/brave-llm-context/SKILL.md](../.agents/skills/brave-llm-context/SKILL.md) | RAG / LLM grounding: pre-extracted web text, tables, code. |
| `brave-answers` | [.agents/skills/brave-answers/SKILL.md](../.agents/skills/brave-answers/SKILL.md) | AI answers via OpenAI-compatible `/chat/completions` with Brave grounding. |
| `brave-news-search` | [.agents/skills/brave-news-search/SKILL.md](../.agents/skills/brave-news-search/SKILL.md) | News articles with metadata, freshness and date filters. |
| `brave-images-search` | [.agents/skills/brave-images-search/SKILL.md](../.agents/skills/brave-images-search/SKILL.md) | Image search with thumbnails (up to 200 results). |
| `brave-videos-search` | [.agents/skills/brave-videos-search/SKILL.md](../.agents/skills/brave-videos-search/SKILL.md) | Video search with duration, views, creator. |
| `brave-suggest` | [.agents/skills/brave-suggest/SKILL.md](../.agents/skills/brave-suggest/SKILL.md) | Fast query autocomplete / suggestions. |
| `brave-spellcheck` | [.agents/skills/brave-spellcheck/SKILL.md](../.agents/skills/brave-spellcheck/SKILL.md) | Spell correction and “did you mean” (often redundant with search endpoints). |
| `brave-local-pois` | [.agents/skills/brave-local-pois/SKILL.md](../.agents/skills/brave-local-pois/SKILL.md) | Local POI details (requires POI IDs from location-filtered web search). |
| `brave-local-descriptions` | [.agents/skills/brave-local-descriptions/SKILL.md](../.agents/skills/brave-local-descriptions/SKILL.md) | AI POI descriptions from POI IDs (max 20 per request). |

### Build: code, UI, APIs, sources

| Skill | Path | Summary |
|--------|------|---------|
| `clean-code` | [.agents/skills/clean-code/SKILL.md](../.agents/skills/clean-code/SKILL.md) | Pragmatic standards: concise, direct, no over-engineering. |
| `source-driven-development` | [.agents/skills/source-driven-development/SKILL.md](../.agents/skills/source-driven-development/SKILL.md) | Ground decisions in official documentation. |
| `frontend-ui-engineering` | [.agents/skills/frontend-ui-engineering/SKILL.md](../.agents/skills/frontend-ui-engineering/SKILL.md) | Production-quality UI components and state. |
| `frontend-design` | [.agents/skills/frontend-design/SKILL.md](../.agents/skills/frontend-design/SKILL.md) | Distinctive, high-quality frontend visuals. |
| `api-and-interface-design` | [.agents/skills/api-and-interface-design/SKILL.md](../.agents/skills/api-and-interface-design/SKILL.md) | Stable APIs and module boundaries. |

### Test and debug

| Skill | Path | Summary |
|--------|------|---------|
| `tdd-workflow` | [.agents/skills/tdd-workflow/SKILL.md](../.agents/skills/tdd-workflow/SKILL.md) | Red–green–refactor workflow. |
| `tdd-guide` | [.agents/skills/tdd-guide/SKILL.md](../.agents/skills/tdd-guide/SKILL.md) | TDD specialist: tests-first, coverage discipline. |
| `test-driven-development` | [.agents/skills/test-driven-development/SKILL.md](../.agents/skills/test-driven-development/SKILL.md) | Tests as proof for logic and behavior changes. |
| `testing-patterns` | [.agents/skills/testing-patterns/SKILL.md](../.agents/skills/testing-patterns/SKILL.md) | Unit, integration, mocking patterns. |
| `e2e-testing` | [.agents/skills/e2e-testing/SKILL.md](../.agents/skills/e2e-testing/SKILL.md) | Playwright E2E patterns and CI. |
| `e2e-runner` | [.agents/skills/e2e-runner/SKILL.md](../.agents/skills/e2e-runner/SKILL.md) | Run and maintain E2E suites (browser tooling). |
| `browser-testing-with-devtools` | [.agents/skills/browser-testing-with-devtools/SKILL.md](../.agents/skills/browser-testing-with-devtools/SKILL.md) | Browser verification with DevTools-style runtime data. |
| `debugging-and-error-recovery` | [.agents/skills/debugging-and-error-recovery/SKILL.md](../.agents/skills/debugging-and-error-recovery/SKILL.md) | Systematic triage: reproduce, localize, fix, guard. |

### Review, security, performance

| Skill | Path | Summary |
|--------|------|---------|
| `code-review-and-quality` | [.agents/skills/code-review-and-quality/SKILL.md](../.agents/skills/code-review-and-quality/SKILL.md) | Multi-axis review before merge. |
| `code-simplification` | [.agents/skills/code-simplification/SKILL.md](../.agents/skills/code-simplification/SKILL.md) | Simplify while preserving behavior. |
| `security-review` | [.agents/skills/security-review/SKILL.md](../.agents/skills/security-review/SKILL.md) | Security patterns for auth, input, secrets, APIs. |
| `security-and-hardening` | [.agents/skills/security-and-hardening/SKILL.md](../.agents/skills/security-and-hardening/SKILL.md) | OWASP-oriented hardening and boundaries. |
| `semgrep-security` | [.agents/skills/semgrep-security/SKILL.md](../.agents/skills/semgrep-security/SKILL.md) | Semgrep static analysis. |
| `trivy-security` | [.agents/skills/trivy-security/SKILL.md](../.agents/skills/trivy-security/SKILL.md) | Trivy: vulns, misconfigs, secrets, SBOM. |
| `vulnerability-scanner` | [.agents/skills/vulnerability-scanner/SKILL.md](../.agents/skills/vulnerability-scanner/SKILL.md) | Threat modeling and prioritization. |
| `security-scan` | [.agents/skills/security-scan/SKILL.md](../.agents/skills/security-scan/SKILL.md) | Audit Claude Code / agent configuration layout. |
| `performance-optimization` | [.agents/skills/performance-optimization/SKILL.md](../.agents/skills/performance-optimization/SKILL.md) | Measure-first performance work. |

### Data and persistence

| Skill | Path | Summary |
|--------|------|---------|
| `database-design` | [.agents/skills/database-design/SKILL.md](../.agents/skills/database-design/SKILL.md) | Schema and persistence decisions. |
| `database-migrations` | [.agents/skills/database-migrations/SKILL.md](../.agents/skills/database-migrations/SKILL.md) | Safe migrations and rollbacks. |
| `database-reviewer` | [.agents/skills/database-reviewer/SKILL.md](../.agents/skills/database-reviewer/SKILL.md) | PostgreSQL performance, RLS, SQL review. |

### Ship, git, CI/CD, docs

| Skill | Path | Summary |
|--------|------|---------|
| `git-workflow-and-versioning` | [.agents/skills/git-workflow-and-versioning/SKILL.md](../.agents/skills/git-workflow-and-versioning/SKILL.md) | Trunk-friendly workflow and small changes. |
| `ci-cd-and-automation` | [.agents/skills/ci-cd-and-automation/SKILL.md](../.agents/skills/ci-cd-and-automation/SKILL.md) | Pipelines and quality gates. |
| `deprecation-and-migration` | [.agents/skills/deprecation-and-migration/SKILL.md](../.agents/skills/deprecation-and-migration/SKILL.md) | Deprecations and user migrations. |
| `shipping-and-launch` | [.agents/skills/shipping-and-launch/SKILL.md](../.agents/skills/shipping-and-launch/SKILL.md) | Launch checklists, rollouts, rollback. |
| `documentation-and-adrs` | [.agents/skills/documentation-and-adrs/SKILL.md](../.agents/skills/documentation-and-adrs/SKILL.md) | ADRs and durable documentation. |
| `documentation-templates` | [.agents/skills/documentation-templates/SKILL.md](../.agents/skills/documentation-templates/SKILL.md) | Templates for READMEs and API docs. |
| `documentation-lookup` | [.agents/skills/documentation-lookup/SKILL.md](../.agents/skills/documentation-lookup/SKILL.md) | Framework / SDK docs via **Context7 MCP** (`resolve-library-id`, `query-docs`). |
| `context7` | [.agents/skills/context7/SKILL.md](../.agents/skills/context7/SKILL.md) | Same Context7 catalog via **HTTP API** and `curl` when MCP is unavailable. |
| `doc-updater` | [.agents/skills/doc-updater/SKILL.md](../.agents/skills/doc-updater/SKILL.md) | Codemaps and doc refresh workflows. |

### Tooling and codebase index

| Skill | Path | Summary |
|--------|------|---------|
| `indexing-and-memory` | [.agents/skills/references/indexing-and-memory.md](../.agents/skills/references/indexing-and-memory.md) | Hub checklist: Engram, graphify, structural search, MCP memory. |
| `memory-protocol` | [.agents/skills/memory-protocol/SKILL.md](../.agents/skills/memory-protocol/SKILL.md) | Engram MCP discipline (`mem_save`, search, session close); includes **Engram in this repository** pretext—**prefer this path** over `engram-memory-protocol/`. |
| `skill-creator` | [.agents/skills/skill-creator/SKILL.md](../.agents/skills/skill-creator/SKILL.md) | Author new skills to the Agent Skills layout. |
| `nano-banana` | [.agents/skills/nano-banana/SKILL.md](../.agents/skills/nano-banana/SKILL.md) | Image generation and editing (Gemini / Nano Banana). |

### Engram skills (vendored)

These directories mirror (or align with) skills from **[Gentleman-Programming/engram / `skills/`](https://github.com/Gentleman-Programming/engram/tree/main/skills)**. Summaries are from each skill’s frontmatter; open `SKILL.md` for full rules. Several bodies reference the **Engram product** codebase (Go, SQLite, dashboard, TUI)—when a file has **Engram in this repository**, apply the **principles** here and concrete paths from **this** hub (see also [`memory.md`](memory.md)).

| Directory (`name` in frontmatter) | Path | Summary |
|-----------------------------------|------|---------|
| `engram-architecture-guardrails` | [.agents/skills/engram-architecture-guardrails/SKILL.md](../.agents/skills/engram-architecture-guardrails/SKILL.md) | System boundaries: local store, cloud, dashboard, plugins. |
| `engram-backlog-triage` | [.agents/skills/engram-backlog-triage/SKILL.md](../.agents/skills/engram-backlog-triage/SKILL.md) | Issue/PR audit, disposition report, maintainer stance. |
| `engram-branch-pr` | [.agents/skills/engram-branch-pr/SKILL.md](../.agents/skills/engram-branch-pr/SKILL.md) | PR workflow, issue-first enforcement. |
| `engram-business-rules` | [.agents/skills/engram-business-rules/SKILL.md](../.agents/skills/engram-business-rules/SKILL.md) | Product rules: sync, admin, permissions, memory semantics. |
| `engram-commit-hygiene` | [.agents/skills/engram-commit-hygiene/SKILL.md](../.agents/skills/engram-commit-hygiene/SKILL.md) | Conventional commits and branch naming (upstream rulesets). |
| `engram-cultural-norms` | [.agents/skills/engram-cultural-norms/SKILL.md](../.agents/skills/engram-cultural-norms/SKILL.md) | Collaboration and quality norms for contributors and agents. |
| `engram-dashboard-htmx` | [.agents/skills/engram-dashboard-htmx/SKILL.md](../.agents/skills/engram-dashboard-htmx/SKILL.md) | HTMX / server-rendered dashboard interaction rules. |
| `engram-docs-alignment` | [.agents/skills/engram-docs-alignment/SKILL.md](../.agents/skills/engram-docs-alignment/SKILL.md) | Docs must match current behavior and examples. |
| `gentleman-bubbletea` | [.agents/skills/engram-gentleman-bubbletea/SKILL.md](../.agents/skills/engram-gentleman-bubbletea/SKILL.md) | Bubbletea TUI patterns (installer / Go TUI context). |
| `engram-issue-creation` | [.agents/skills/engram-issue-creation/SKILL.md](../.agents/skills/engram-issue-creation/SKILL.md) | Issue workflow, issue-first system. |
| `engram-memory-protocol` | [.agents/skills/engram-memory-protocol/SKILL.md](../.agents/skills/engram-memory-protocol/SKILL.md) | Same MCP discipline as `memory-protocol`; **prefer** [.agents/skills/memory-protocol/SKILL.md](../.agents/skills/memory-protocol/SKILL.md) in this hub. |
| `engram-plugin-thin` | [.agents/skills/engram-plugin-thin/SKILL.md](../.agents/skills/engram-plugin-thin/SKILL.md) | Thin Claude/OpenCode/Gemini/Codex plugin adapters vs Go core. |
| `engram-pr-review-deep` | [.agents/skills/engram-pr-review-deep/SKILL.md](../.agents/skills/engram-pr-review-deep/SKILL.md) | Deep technical review before merge. |
| `engram-project-structure` | [.agents/skills/engram-project-structure/SKILL.md](../.agents/skills/engram-project-structure/SKILL.md) | Where files, handlers, templates, and tests belong (Engram repo layout). |
| `engram-sdd-flow` | [.agents/skills/engram-sdd-flow/SKILL.md](../.agents/skills/engram-sdd-flow/SKILL.md) | Short phase-order reference (Engram product context); hub workflow is **`skillgrid-*`** + **`openspec-*`** skills optional. |
| `engram-server-api` | [.agents/skills/engram-server-api/SKILL.md](../.agents/skills/engram-server-api/SKILL.md) | HTTP API contract guardrails for server changes. |
| `engram-testing-coverage` | [.agents/skills/engram-testing-coverage/SKILL.md](../.agents/skills/engram-testing-coverage/SKILL.md) | TDD and coverage standards (Engram packages). |
| `engram-tui-quality` | [.agents/skills/engram-tui-quality/SKILL.md](../.agents/skills/engram-tui-quality/SKILL.md) | Bubbletea/Lipgloss TUI quality and navigation. |
| `engram-ui-elements` | [.agents/skills/engram-ui-elements/SKILL.md](../.agents/skills/engram-ui-elements/SKILL.md) | Dashboard pages, cards, metrics, detail flows. |
| `engram-visual-language` | [.agents/skills/engram-visual-language/SKILL.md](../.agents/skills/engram-visual-language/SKILL.md) | Dashboard styling, typography, spacing, visual identity. |

Upstream catalog: [engram `skills/catalog.md`](https://github.com/Gentleman-Programming/engram/blob/main/skills/catalog.md).

### Session learning (optional)

| Skill | Path | Summary |
|--------|------|---------|
| `continuous-learning` | [.agents/skills/continuous-learning/SKILL.md](../.agents/skills/continuous-learning/SKILL.md) | Extract patterns from sessions into learned skills. |
| `continuous-learning-v2` | [.agents/skills/continuous-learning-v2/SKILL.md](../.agents/skills/continuous-learning-v2/SKILL.md) | Hook-driven instincts with confidence and project scope. |

---

## Adding or updating a skill

Use **`skill-creator`** and keep one directory per skill with a root **`SKILL.md`**. Optional subfolders: `assets/`, `references/` (per-skill), or rely on shared **`.agents/skills/references/`** for checklists.

After adding skills, update this file (`docs/skills.md`) and the phase lists in [`workflow.md`](workflow.md) if the workflow should surface them. **Markdown file templates** for project/PRD live in **`/skillgrid-init`** and **`/skillgrid-plan`**, not in `workflow.md`. If a skill introduces a new CLI or MCP dependency, add it to [`tools.md`](tools.md). For how **Engram**, **graphify**, and search fit together, see [`memory.md`](memory.md). **Workflow phases** (`skillgrid-init`, `skillgrid-plan`, …) are documented in command files, not as `sdd-*` skills.
