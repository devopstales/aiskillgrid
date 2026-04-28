# Agent commands

Slash-style commands are **Markdown prompts** checked into this repo. Each file tells the agent what to do for that command; Cursor, Kilo, and OpenCode use YAML frontmatter (`name`, `id`, `category`, `description`, plus optional **`allowed-tools`** and **`argument-hint`**). Bodies use **`<objective>`** / **`<process>`** to separate intent from procedure. GitHub Copilot uses **`.github/prompts/`** with `description`, the same optional keys when present, and the same body.

**Lifecycle overview:** [docs/workflow.md](workflow.md) summarizes phases; each **`skillgrid-*.md`** command is the source of truth for checklists. **Skillgrid** commands wrap those phases; **OpenSpec (`opsx`)** commands target the OpenSpec CLI workflow directly. **Tooling and installers:** [`tools.md`](tools.md).

**Mirroring:** Canonical slash commands and agent personas live under **`.cursor/commands/`** and **`.cursor/agents/`**. After editing, run **[`scripts/sync-ide-assets.sh`](../scripts/sync-ide-assets.sh)** (see [`tools.md`](tools.md)) to update **`.kilo/commands/`**, **`.opencode/commands/`**, **`.github/prompts/`**, and **`*/agents/`** mirrors. Use **`--check`** in CI to detect drift.

---

## Where files live

| Target | Directory | Notes |
|--------|-----------|--------|
| Cursor | [`.cursor/commands/`](../.cursor/commands/) | Full frontmatter; primary copy for **skillgrid-*** updates. |
| Kilo | [`.kilo/commands/`](../.kilo/commands/) | Same bodies as Cursor where mirrored. |
| OpenCode | [`.opencode/commands/`](../.opencode/commands/) | Same pattern as Cursor. |
| GitHub Copilot | [`.github/prompts/`](../.github/prompts/) | `description` + body only. |

Invocation in the UI may show **`/command-id`** (for example `/skillgrid-plan` or `/opsx-apply`). Exact behavior depends on the product; the file **`name`** field in Cursor-style commands is the canonical label.

---

## Skillgrid (full workflow)

These align with **[docs/workflow.md](workflow.md)**. **Commands orchestrate phases**: objective, order, status-on-exit, and completion report. Reusable procedures now live in **`.agents/skills/skillgrid-*/SKILL.md`** and are referenced by commands. Update **docs/workflow.md** when phase descriptions change and **docs/skills.md** when reusable Skillgrid primitives change.

| Phase | Command | File | Purpose |
|-------|---------|------|---------|
| 0 | `/skillgrid-help` | [skillgrid-help.md](../.cursor/commands/skillgrid-help.md) | Explain workflow commands, artifact layout, phases, and current-state next steps. |
| 0 | `/skillgrid-session` | [skillgrid-session.md](../.cursor/commands/skillgrid-session.md) | Session charter, context budget, MCP selection, checkpoints. |
| 0 | `/skillgrid-checkpoint` | [skillgrid-checkpoint.md](../.cursor/commands/skillgrid-checkpoint.md) | Named workflow checkpoints under `.skillgrid/tasks/checkpoints.log`; create, verify, list, or clear. |
| DEFINE | `/skillgrid-init` | [skillgrid-init.md](../.cursor/commands/skillgrid-init.md) | **`.skillgrid/config.json`** (ticketing + **artifact store** `hybrid` \| `openspec` \| `engram`), greenfield/brownfield, **OpenSpec** bootstrap when mode includes disk, Engram when mode includes memory, graphify, OpenCode, **`.skillgrid/`** tree. |
| DEFINE | `/skillgrid-explore` | [skillgrid-explore.md](../.cursor/commands/skillgrid-explore.md) | OpenSpec/PRD inventory, automatic import/backfill, **`.skillgrid/project/`** (`ARCHITECTURE`, `STRUCTURE`, `PROJECT`), root `AGENTS.md`, semantic search. |
| DEFINE | `/skillgrid-brainstorm` | [skillgrid-brainstorm.md](../.cursor/commands/skillgrid-brainstorm.md) | Clarify, research (`search-first`, `documentation-lookup`), refine ideas before planning. |
| DEFINE/PLAN | `/skillgrid-import` | [skillgrid-import.md](../.cursor/commands/skillgrid-import.md) | Import legacy PRDs and existing OpenSpec changes into canonical `.skillgrid/prd/` files and `INDEX.md`. |
| PLAN | `/skillgrid-plan` | [skillgrid-plan.md](../.cursor/commands/skillgrid-plan.md) | PRDs and OpenSpec CLI artifact loop (hybrid Engram). |
| PLAN | `/skillgrid-breakdown` | [skillgrid-breakdown.md](../.cursor/commands/skillgrid-breakdown.md) | Spec completeness, `tasks.md`, TDD stance, ordering, graphify refresh when needed. |
| BUILD | `/skillgrid-apply` | [skillgrid-apply.md](../.cursor/commands/skillgrid-apply.md) | OpenSpec apply, contracts (`api-and-interface-design`), TDD, small batches. |
| VERIFY | `/skillgrid-test` | [skillgrid-test.md](../.cursor/commands/skillgrid-test.md) | Tests, E2E, browser tooling, debugging. |
| REVIEW | `/skillgrid-security` | [skillgrid-security.md](../.cursor/commands/skillgrid-security.md) | Security review, SAST, scanning, threat framing, deprecation/attack-surface hygiene. |
| REVIEW | `/skillgrid-validate` | [skillgrid-validate.md](../.cursor/commands/skillgrid-validate.md) | **Combined gate:** OpenSpec/spec verification, code review, security review, evidence reconciliation, and sign-off. |
| SHIP | `/skillgrid-finish` | [skillgrid-finish.md](../.cursor/commands/skillgrid-finish.md) | Archive or sync specs, git/PR, CI, docs/ADRs, ship checklist. |

The same filenames exist under **`.kilo/commands/`**, **`.opencode/commands/`**, and **`.github/prompts/`**.

---

## OpenSpec (`opsx`)

These are **thin aliases** of **`skillgrid-*`** (same outcomes; canonical steps live in `skillgrid-*.md`). Optional **`.agents/skills/openspec-*/`** files remain for IDEs that discover skills separately.

| Command | File (example) | Purpose |
|---------|----------------|---------|
| `/opsx-onboard` | [.cursor/commands/opsx-onboard.md](../.cursor/commands/opsx-onboard.md) | Guided first full cycle (tutorial-style). |
| `/opsx-explore` | [.cursor/commands/opsx-explore.md](../.cursor/commands/opsx-explore.md) | Explore only: thinking partner, no implementation. |
| `/opsx-propose` | [.cursor/commands/opsx-propose.md](../.cursor/commands/opsx-propose.md) | New change and artifacts in one step. |
| `/opsx-apply` | [.cursor/commands/opsx-apply.md](../.cursor/commands/opsx-apply.md) | Implement change tasks. |
| `/opsx-verify` | [.cursor/commands/opsx-verify.md](../.cursor/commands/opsx-verify.md) | Verify code against change artifacts. |
| `/opsx-sync` | [.cursor/commands/opsx-sync.md](../.cursor/commands/opsx-sync.md) | Sync delta specs to main specs. |
| `/opsx-archive` | [.cursor/commands/opsx-archive.md](../.cursor/commands/opsx-archive.md) | Archive completed change. |

Mirrors: **`.kilo/commands/`**, **`.opencode/commands/`**, **`.github/prompts/`** with the same `opsx-*.md` names.

---

## Related documentation

- [`agents.md`](agents.md) — IDE agent personas (`.cursor/agents/` and mirrors).
- [docs/workflow.md](workflow.md) — phase list, **`.skillgrid/config.json`**, PRD and OpenSpec handoff.
- [`skills.md`](skills.md) — full skill catalog and paths.
- [`tools.md`](tools.md) — CLIs, MCPs, and `install.sh` dependencies.
