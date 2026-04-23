---
description: Brainstorm and refine ideas before committing to a plan
---

You are executing **`/skillgrid-brainstorm`** (DEFINE phase) for the Skillgrid workflow.

## Steps

1. **Clarify** — Ask targeted questions until goals, constraints, and success criteria are explicit.
2. **Diverge** — List options, alternatives, and tradeoffs; keep judgment light until the space is wide enough.
3. **Research** — Use `search-first` and the open web for prior art; use `documentation-lookup` (Context7) when the idea depends on a specific framework or library.
4. **Converge** — Rank approaches; state assumptions and risks (see `karpathy-guidelines`).
5. **Validate** — Use `deep-research` when external evidence or breadth should inform the choice.
6. **Refine** — Use `idea-refine` (divergent/convergent structure) to sharpen a vague idea into a defensible direction.
7. **Handoff** — Output should feed `/skillgrid-plan`, not a full locked spec. Do not over-commit to file artifacts unless the user asks.

## Skills to read and follow

- `.agents/skills/karpathy-guidelines/SKILL.md` — surface tradeoffs and alternatives before locking direction.
- `.agents/skills/search-first/SKILL.md` — research tools and patterns before building.
- `.agents/skills/documentation-lookup/SKILL.md` — authoritative library/framework docs via Context7 MCP.
- `.agents/skills/deep-research/SKILL.md` — validate assumptions and gather evidence.
- `.agents/skills/idea-refine/SKILL.md` — structured ideation to sharpen a vague idea.

## Optional: IDE personas

For a **dedicated research subagent** that leans on hub MCPs (**Exa**, **Firecrawl**, **DeepWiki**, **Context7**) and delivers a cited memo, spawn **`skillgrid-researcher`** ([`.cursor/agents/skillgrid-researcher.md`](../../.cursor/agents/skillgrid-researcher.md)).

## Notes

- Inspect the repo with tools when brainstorming touches implementation reality.
- If OpenSpec or SDD modes are unclear, ask once, then align with existing `openspec/` or repo conventions.
