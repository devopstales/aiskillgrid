# Agent skills

This project keeps reusable agent instructions under **`.agents/skills/<skill-name>/SKILL.md`**. Tools discover them from `.agents/skills/` (and may also load skills from user-global paths, depending on the IDE).

For **IDE persona files** (`agents/` under each IDE config dir), see [`agents.md`](agents.md). For **how skills fit into the end-to-end flow**, see [`wokflow.md`](wokflow.md). For **CLIs, MCPs, and install.sh dependencies** referenced by skills, see [`tools.md`](tools.md).

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

### Spec-Driven Development (SDD)

| Skill | Path | Summary |
|--------|------|---------|
| `sdd-init` | [.agents/skills/sdd-init/SKILL.md](../.agents/skills/sdd-init/SKILL.md) | Initialize SDD context: stack, conventions, persistence mode. |
| `sdd-propose` | [.agents/skills/sdd-propose/SKILL.md](../.agents/skills/sdd-propose/SKILL.md) | Create or update `proposal.md` for a named change. |
| `sdd-spec` | [.agents/skills/sdd-spec/SKILL.md](../.agents/skills/sdd-spec/SKILL.md) | Write delta specs (requirements and scenarios). |
| `sdd-tasks` | [.agents/skills/sdd-tasks/SKILL.md](../.agents/skills/sdd-tasks/SKILL.md) | Produce implementation task checklist (`tasks.md`). |
| `sdd-verify` | [.agents/skills/sdd-verify/SKILL.md](../.agents/skills/sdd-verify/SKILL.md) | Verify implementation against specs, design, and tasks. |

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
| `documentation-lookup` | [.agents/skills/documentation-lookup/SKILL.md](../.agents/skills/documentation-lookup/SKILL.md) | Framework docs via Context7 MCP. |
| `doc-updater` | [.agents/skills/doc-updater/SKILL.md](../.agents/skills/doc-updater/SKILL.md) | Codemaps and doc refresh workflows. |

### Tooling and codebase index

| Skill | Path | Summary |
|--------|------|---------|
| `ccc` | [.agents/skills/ccc/SKILL.md](../.agents/skills/ccc/SKILL.md) | CocoIndex Code: semantic search and index lifecycle. |
| `skill-creator` | [.agents/skills/skill-creator/SKILL.md](../.agents/skills/skill-creator/SKILL.md) | Author new skills to the Agent Skills layout. |
| `nano-banana` | [.agents/skills/nano-banana/SKILL.md](../.agents/skills/nano-banana/SKILL.md) | Image generation and editing (Gemini / Nano Banana). |

### Session learning (optional)

| Skill | Path | Summary |
|--------|------|---------|
| `continuous-learning` | [.agents/skills/continuous-learning/SKILL.md](../.agents/skills/continuous-learning/SKILL.md) | Extract patterns from sessions into learned skills. |
| `continuous-learning-v2` | [.agents/skills/continuous-learning-v2/SKILL.md](../.agents/skills/continuous-learning-v2/SKILL.md) | Hook-driven instincts with confidence and project scope. |

---

## Adding or updating a skill

Use **`skill-creator`** and keep one directory per skill with a root **`SKILL.md`**. Optional subfolders: `assets/`, `references/` (per-skill), or rely on shared **`.agents/skills/references/`** for checklists.

After adding skills, update this file (`docs/skills.md`) and the phase lists in [`wokflow.md`](wokflow.md) if the workflow should surface them. If a skill introduces a new CLI or MCP dependency, add it to [`tools.md`](tools.md).
