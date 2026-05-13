---
description: C4-style architecture diagrams (ASCII or Mermaid) for systems and codebases
---

You are executing `/sdd-c4` for the SDD workflow.

## Skill Loading Contract

Load and follow first:
- `.agents/skills/c4-diagrams/SKILL.md`
- `.agents/skills/c4-diagrams/templates.md` when the skill references it

## Workflow

1. Resolve purpose (existing code vs new system), format (ASCII vs Mermaid), and rigor per the skill’s gate table; ask one combined question if ambiguous.
2. Explore the codebase when mapping existing architecture.
3. Produce only diagram levels that add value; record explicit assumptions if the user chose speed over full rigor.

## Return Contract

Return a structured result with:
- `status`: `completed | blocked | failed`
- `executive_summary`
- `detailed_report` (diagrams, assumptions, files referenced)
- `artifacts` (paths if diagrams were written to disk; otherwise inline body in report)
- `next_recommended`
- `risks`
