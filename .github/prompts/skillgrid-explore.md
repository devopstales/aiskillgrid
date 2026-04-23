---
description: Explore the problem and repo: OpenSpec explore, .skillgrid/project docs, AGENTS, semantic search
---

You are executing **`/skillgrid-explore`** (DEFINE phase) for the Skillgrid workflow.

## Steps

1. **Stance** — Explore and clarify; do not implement production code unless the user explicitly leaves explore mode.
2. **OpenSpec explore** — Follow `openspec-explore`: investigate the problem space and relevant code paths before locking a change.
3. **Project docs (canonical)** — Create or refresh **`.skillgrid/project/ARCHITECTURE.md`**, **`STRUCTURE.md`**, and **`PROJECT.md`** (system design, repo layout and optional runtime topology, onboarding narrative). Use templates in [`docs/wokflow.md`](../../docs/wokflow.md) when helpful.
4. **AGENTS.md** — Create or refresh at repo root so agent behavior and project rules are current.
5. **Documentation** — When recording exploration outcomes, document the *why* (ADRs, API docs, inline standards) per team norms.
6. **Semantic search** — Use `ccc search` and repo navigation where CocoIndex is set up.
7. **Optional depth** — Use `deep-research` when the question needs external evidence or broader comparison.

## Skills to read and follow

- `.agents/skills/openspec-explore/SKILL.md`
- `.agents/skills/search-first/SKILL.md`
- `.agents/skills/deep-research/SKILL.md`
- `.agents/skills/ccc/SKILL.md`
- `.agents/skills/documentation-and-adrs/SKILL.md`

## Optional: IDE personas

When spawning a **subagent** for exploration-only work in a clean context, use **`skillgrid-explore-architect`** ([`.cursor/agents/skillgrid-explore-architect.md`](../../.cursor/agents/skillgrid-explore-architect.md)).

For **external / cited research** (landscape, competitors, docs via Exa, Firecrawl, DeepWiki, Context7) rather than in-repo mapping, use **`skillgrid-researcher`** ([`.cursor/agents/skillgrid-researcher.md`](../../.cursor/agents/skillgrid-researcher.md)).

## Notes

- Inspect the repo with tools; do not assume stack or layout.
- If OpenSpec or SDD modes are unclear, ask once, then align with existing `openspec/` or repo conventions.
