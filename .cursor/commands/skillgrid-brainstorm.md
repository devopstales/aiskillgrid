---
name: /skillgrid-brainstorm
id: skillgrid-brainstorm
category: Workflow
description: Brainstorm and refine ideas before committing to a plan
allowed-tools: Read, Glob, Grep, WebSearch, WebFetch, Task
argument-hint: "[idea, question, or problem to refine]"
---

<objective>

You are executing **`/skillgrid-brainstorm`** (DEFINE phase) for the Skillgrid workflow.

Turn ideas into a **reviewable direction** through collaborative dialogue: understand context, ask in a disciplined way, explore options, then hand off to planning—without implementing.

</objective>

<process>

## Implementation gate

Do **not** write production code, scaffold apps, or drive **`/skillgrid-apply`** from this command. Brainstorming ends in a **clear enough** problem statement and preferred approach for **`/skillgrid-plan`** (use **`/skillgrid-explore`** first if repo or OpenSpec context is still thin).

## Questioning discipline

Base how you ask on [obra/superpowers **brainstorming**](https://github.com/obra/superpowers/blob/main/skills/brainstorming/SKILL.md) (dialogue patterns); stay inside Skillgrid phases and skills below.

| Practice | What to do |
|----------|------------|
| **One question per message** | Avoid bundling multiple unrelated questions in one turn. If a topic needs depth, split into a sequence of single questions. |
| **Multiple choice when it fits** | Prefer A/B/C (or similar) when it narrows intent faster than open-ended text; use open-ended when exploration is the point. |
| **Context before grilling** | When the idea touches the codebase, skim relevant files, docs, or recent commits before detailed questions so questions are grounded. |
| **Scope check early** | If the request bundles several independent subsystems (e.g. chat + billing + analytics in one breath), **name that** and help **decompose** before refining low-level details. Brainstorm the **first slice** here; other slices get their own plan cycles later. |
| **“Too simple” is still worth clarifying** | Short or “obvious” ideas still need explicit goals, constraints, and success criteria—keep the design proportionate (a few sentences vs a longer outline), but don’t skip alignment. |
| **Alternatives before commitment** | Before settling, surface **2–3 approaches** with tradeoffs and a recommendation (ties to **Diverge** / **Converge** steps). |
| **Incremental buy-in** | When you present a emerging design, chunk it by section (architecture, data flow, errors, testing mindset) and **check** “does this still match what you want?” as you go; revise if not. |
| **Visual vs text** | Not every UI-themed question needs pixels. Conceptual questions stay in chat; mockups, layout comparisons, or diagrams belong in a **visual** surface (e.g. Canvas or browser) **when** seeing beats reading. If you use a visual workflow, **offer it once in its own message** (consent, token cost)—**no** combining that offer with a clarifying question in the same turn. |

## Steps

1. **Clarify** — Follow **Questioning discipline** until goals, constraints, and success criteria are explicit.
2. **Diverge** — List options, alternatives, and tradeoffs; keep judgment light until the space is wide enough.
3. **Research** — Use `search-first` and the open web for prior art; use `documentation-lookup` (Context7) when the idea depends on a specific framework or library.
4. **Converge** — Rank approaches; state assumptions and risks (see `karpathy-guidelines`).
5. **Validate** — Use `deep-research` when external evidence or breadth should inform the choice.
6. **Refine** — Use `idea-refine` (divergent/convergent structure) to sharpen a vague idea into a defensible direction.
7. **Handoff** — Output should feed **`/skillgrid-plan`**, not a full locked spec. 

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
- Prior art for question style and gates

</process>
