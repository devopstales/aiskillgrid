# Agent commands

Slash-style commands are **Markdown prompts** checked into this repo. Each file tells the agent what to do for that command; Cursor, Kilo, and OpenCode use YAML frontmatter (`name`, `id`, `category`, `description`). GitHub Copilot uses **`.github/prompts/`** with a short `description` only.

**Lifecycle overview:** [`wokflow.md`](wokflow.md) maps phases to skills. **Skillgrid** commands wrap those phases; **OpenSpec (`opsx`)** commands target the OpenSpec CLI workflow directly. **Tooling and installers:** [`tools.md`](tools.md).

---

## Where files live

| Target | Directory | Notes |
|--------|-----------|--------|
| Cursor | [`.cursor/commands/`](../.cursor/commands/) | Full frontmatter. |
| Kilo | [`.kilo/commands/`](../.kilo/commands/) | Same bodies as Cursor where mirrored. |
| OpenCode | [`.opencode/commands/`](../.opencode/commands/) | Same pattern as Cursor. |
| GitHub Copilot | [`.github/prompts/`](../.github/prompts/) | `description` + body only. |

Invocation in the UI may show **`/command-id`** (for example `/skillgrid-plan` or `/opsx-apply`). Exact behavior depends on the product; the file **`name`** field in Cursor-style commands is the canonical label.

---

## Skillgrid (full workflow)

These align with **`docs/wokflow.md`** sections. Each command lists skills to read under **`.agents/skills/`** and defers to the workflow doc if anything disagrees.

| Command | File (example) | Purpose |
|---------|----------------|---------|
| `/skillgrid-init` | [.cursor/commands/skillgrid-init.md](../.cursor/commands/skillgrid-init.md) | Bootstrap structure, graphify, CocoIndex (`ccc`), OpenCode, baseline skills. |
| `/skillgrid-explore` | [.cursor/commands/skillgrid-explore.md](../.cursor/commands/skillgrid-explore.md) | Explore problem and repo; AGENTS/PROJECT; semantic search. |
| `/skillgrid-brainstorm` | [.cursor/commands/skillgrid-brainstorm.md](../.cursor/commands/skillgrid-brainstorm.md) | Refine ideas before planning. |
| `/skillgrid-plan` | [.cursor/commands/skillgrid-plan.md](../.cursor/commands/skillgrid-plan.md) | PRDs and proposals (OpenSpec + SDD). |
| `/skillgrid-design` | [.cursor/commands/skillgrid-design.md](../.cursor/commands/skillgrid-design.md) | UI, APIs, architecture, delta specs. |
| `/skillgrid-breakdown` | [.cursor/commands/skillgrid-breakdown.md](../.cursor/commands/skillgrid-breakdown.md) | Task breakdown and ordering. |
| `/skillgrid-apply` | [.cursor/commands/skillgrid-apply.md](../.cursor/commands/skillgrid-apply.md) | Implement with OpenSpec apply, TDD, small batches. |
| `/skillgrid-test` | [.cursor/commands/skillgrid-test.md](../.cursor/commands/skillgrid-test.md) | Tests, E2E, browser tooling. |
| `/skillgrid-validate` | [.cursor/commands/skillgrid-validate.md](../.cursor/commands/skillgrid-validate.md) | Verify artifacts, review, security, performance, docs. |
| `/skillgrid-finish` | [.cursor/commands/skillgrid-finish.md](../.cursor/commands/skillgrid-finish.md) | Archive or sync specs, git/PR, ship checklist. |

The same filenames exist under **`.kilo/commands/`**, **`.opencode/commands/`**, and **`.github/prompts/`**.

---

## OpenSpec (`opsx`)

These map to OpenSpec operations and the matching skills under **`.agents/skills/openspec-*/`**.

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

- [`wokflow.md`](wokflow.md) — phase list and skill bullets per phase.
- [`skills.md`](skills.md) — full skill catalog and paths.
- [`tools.md`](tools.md) — CLIs, MCPs, and `install.sh` dependencies.
