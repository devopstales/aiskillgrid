---
description: Start a new SDD change — runs exploration, clarification, proposal, specs, design (including UI), PRD, and tasks
agent: sdd-orchestrator
---

Follow the SDD orchestrator workflow for starting a new change named "$ARGUMENTS".

WORKFLOW:
1. **Read `.skillgrid/project/CONTEXT.md`** if it exists. Note any relevant glossary terms, assumptions, or success criteria before proceeding.
2. Launch `sdd-explore` sub-agent to investigate the codebase for this change
3. Present the exploration summary to the user
4. Launch `sdd-clarify` sub-agent to refine understanding:
   - It will ask questions to resolve ambiguity about goals, scope, and constraints.
   - It captures the mutual understanding in a `clarifications.md` artifact.
5. Present the clarifications summary and ask the user to confirm or iterate.
6. Once mutual understanding is confirmed, launch `sdd-propose` sub-agent
7. Present the proposal summary and ask the user if they want to continue with specs and design

### UI DESIGN INTEGRATION POINT

If the proposal involves user-facing changes (components, layouts, interactions, visual updates):

- Flag the change with `ui_scope: true` in the proposal metadata
- After `sdd-design`, launch the UI design sub-flow:
  a. `engram-ui-elements` — validate component composition & reusability
  b. `engram-visual-language` — apply spacing, typography, color, and accessibility rules
  c. `sdd-ui-design` (if available) — generate wireframes, interaction flows, or ASCII/Mermaid mockups
- Merge UI artifacts into the technical design before proceeding to PRD

### CONTINUE WITH STANDARD PHASES

Run these sub-agents in sequence after clarification (and UI design if applicable):

1. `sdd-propose` — create the proposal (using the clarifications)
2. `sdd-spec` — write functional & technical specifications
3. `sdd-design` — create technical design architecture
   → **If UI scope**: trigger UI design sub-flow above, then merge outputs
4. `sdd-prd` — consolidate into PRD (`.skillgrid/prd/PRD<NN>_$ARGUMENTS.md`)
   - Include UI decisions, wireframe references, and accessibility notes in the PRD
5. `sdd-tasks` — break down into vertically sliced implementation tasks
   - Tag UI-related tasks with `area: frontend` or `area: ui`

CONTEXT:
- Working directory: !`echo -n "$(pwd)"`
- Current project: !`echo -n "$(basename $(pwd))"`
- Change name: $ARGUMENTS
- Artifact store mode: hybrid
- UI scope detection: Analyze proposal for keywords: [ui, ux, frontend, component, layout, visual, wireframe, mockup, css, styling, accessibility]

ENGRAM NOTE:
Sub-agents handle persistence automatically. Each phase saves its artifact to engram with topic_key:
- `"sdd/$ARGUMENTS/{type}"` where type ∈ {explore, clarify, propose, spec, design, prd, tasks}
- UI-specific artifacts use: `"sdd/$ARGUMENTS/ui/{wireframes, decisions, tokens}"`

FILESYSTEM PERSISTENCE:
- Read `.agents/skills/_shared/skillgrid-handoff.md` for filesystem persistence instructions.
- UI artifacts in hybrid mode should also write to:

```
openspec/changes/$ARGUMENTS/
  ├── ui-wireframes.md      # ASCII/Mermaid diagrams
  ├── ui-decisions.md       # Design rationale & tradeoffs
  └── ui-tokens.md          # Color, spacing, typography references (if project uses design tokens)
```

ORCHESTRATOR RULES:
- Do NOT execute phase work inline — delegate to sub-agents.
- Before launching `sdd-tasks`, verify:
  - Technical design is complete
  - UI artifacts are merged (if ui_scope: true)
  - Accessibility considerations are documented
  - All artifacts are persisted per artifact_store mode
- If user interrupts or requests changes, use `mem_get_observation` to reload prior phase artifacts instead of re-running.

UI DESIGN HANDOFF CONTRACT:
When UI sub-agents complete, ensure the design artifact includes:

```yaml
ui_design_summary:
components_affected: [List<ComponentName>]
wireframe_refs: ["path/to/wireframe.md#section"]
accessibility_notes: ["contrast ratios", "keyboard nav", "ARIA labels"]
responsive_breakpoints: ["mobile", "tablet", "desktop"]
design_token_usage: ["color.primary", "spacing.md", "typography.body"]
risks: ["browser compatibility", "performance impact"]
```

Read the orchestrator instructions to coordinate this workflow. Do NOT execute phase work inline — delegate to sub-agents.