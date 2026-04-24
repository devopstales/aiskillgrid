# IDE agent personas (`agents/` folders)

This hub ships **specialist agent definitions** as Markdown files. Each IDE keeps them under an **`agents/`** directory inside that product’s config folder. The IDE (or harness) loads them as **personas**: focused system prompts for subagents or delegated turns—not the same thing as [Agent Skills](skills.md) under `.agents/skills/`.

Persona files in this hub use the **`skillgrid-`** prefix (frontmatter `name` and filename) so they are easy to distinguish from other agent definitions in a merged project.

---

## Where the files live (this repo)

| IDE / product | Directory | Notes |
|---------------|-----------|--------|
| **Cursor** | [`.cursor/agents/`](../.cursor/agents/) | **Editorial source of truth** in this hub for persona `.md` files. |
| **Kilo** | [`.kilo/agents/`](../.kilo/agents/) | Keep in sync with Cursor when you change a persona. |
| **OpenCode** | [`.opencode/agents/`](../.opencode/agents/) | Same. |
| **GitHub Copilot** | [`.github/agents/`](../.github/agents/) | Same. |

Each directory includes a **[README.md](../.cursor/agents/README.md)** that explains how **personas** relate to **skills** and **slash commands**, orchestration rules, and how to add a new persona.

`install.sh` **rsyncs the whole IDE folder** (e.g. `.cursor/`) into the target project, so **`agents/` is installed together** with `commands/`, `rules/`, and `mcp.json`. In **this hub repo**, run **[`scripts/sync-ide-assets.sh`](../scripts/sync-ide-assets.sh)** after editing `.cursor/agents/` so **`.kilo/`**, **`.opencode/`**, and **`.github/agents/`** stay aligned with Cursor.

---

## Personas shipped here

| File | Purpose |
|------|---------|
| [skillgrid-code-reviewer.md](../.cursor/agents/skillgrid-code-reviewer.md) | Five-axis code review (correctness, readability, architecture, security, performance). |
| [skillgrid-security-auditor.md](../.cursor/agents/skillgrid-security-auditor.md) | Security-focused audit and threat-oriented review. |
| [skillgrid-test-engineer.md](../.cursor/agents/skillgrid-test-engineer.md) | Test strategy, coverage, and test-quality review. |
| [skillgrid-spec-verifier.md](../.cursor/agents/skillgrid-spec-verifier.md) | Traceability: delta specs, `tasks.md`, and implementation alignment. |
| [skillgrid-explore-architect.md](../.cursor/agents/skillgrid-explore-architect.md) | Brownfield exploration; `.skillgrid/project/` architecture and onboarding docs. |
| [skillgrid-task-breakdown-auditor.md](../.cursor/agents/skillgrid-task-breakdown-auditor.md) | Audits task lists before implementation (AC, ordering, testability). |
| [skillgrid-design-critic.md](../.cursor/agents/skillgrid-design-critic.md) | Critiques `DESIGN.md` (flows, a11y, API boundaries)—not full code review. |
| [skillgrid-researcher.md](../.cursor/agents/skillgrid-researcher.md) | Cited research using hub MCPs: Exa, Firecrawl, DeepWiki, Context7. |

Names match the **frontmatter `name`** field (e.g. `name: skillgrid-code-reviewer`). The **`description`** is what UIs show when picking an agent.

---

## File format

Persona files use YAML frontmatter plus Markdown body, for example:

```yaml
---
name: skillgrid-code-reviewer
description: Senior code reviewer that evaluates changes across five dimensions ...
---
```

The body defines role, scope, output shape, and (per repo convention) a **Composition** section at the end. Do not rely on unsupported frontmatter fields for Claude Code–style plugins (see [.cursor/agents/README.md](../.cursor/agents/README.md)).

---

## Skills vs `agents/` vs commands

| Kind | Location | Role |
|------|----------|------|
| **Skill** | `.agents/skills/<name>/SKILL.md` | Reusable workflow and checklists; loaded when a skill is invoked. |
| **IDE skills mirror** | e.g. `.cursor/.agents/skills/` | Populated by `install.sh` from `.agents/skills/` (sync target, not the persona folder). |
| **Persona (this doc)** | `.cursor/agents/skillgrid-*.md` (and mirrors) | *Who* is answering: a single role and report format. |
| **Slash command** | `.cursor/commands/*.md` | *When* to run a phase or workflow; may point agents at specific skills. |

For the end-to-end Skillgrid phase list, see [skillgrid-workflow.md](../.skillgrid/scripts/skillgrid-workflow.md) and [commands.md](commands.md).

---

## Related documentation

| Doc | Use |
|-----|-----|
| [.cursor/agents/README.md](../.cursor/agents/README.md) | Persona philosophy, orchestration patterns, adding a persona. |
| [.agents/skills/references/indexing-and-memory.md](../.agents/skills/references/indexing-and-memory.md) | Shared checklist: Engram, graphify, structural search, MCP memory (used by personas + skills). |
| [skills.md](skills.md) | Skill catalog (`SKILL.md` trees). |
| [commands.md](commands.md) | Slash commands per IDE. |
| [tools.md](tools.md) | `install.sh` and toolchain. |
