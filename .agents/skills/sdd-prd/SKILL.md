---
description: Generate the formal PRD by filling the project template
agent: sdd-prd
---

You are the `sdd-prd` sub-agent. Produce the PRD file using the existing template at `.skillgrid/templates/template-prd.md`.

## Inputs
- `proposal.md`
- `specifications.md` (from `sdd-spec`)
- `design.md` (from `sdd-design`)
- `clarifications.md` (optional)
- Change slug: `$ARGUMENTS`

## Steps

1. **Read the template**  
   Open `.skillgrid/templates/template-prd.md`. It contains the exact structure and metadata your output must have.

2. **Determine the PRD file number**  
   List files matching `.skillgrid/prd/PRD*.md`. Take the highest `NN`, add 1, zero-pad to two digits. If none exist, use `01`.  
   The output file will be `.skillgrid/prd/PRD<NN>_$ARGUMENTS.md`.

3. **Fill the template**  
   Use the input artifacts to complete every section of the template. Inject the advanced SDD concepts as follows:
   - **Vertical slices** → `Decomposition` section and `Implementation tasks` list.  
     Each implementation task is a slice tagged `[HITL]` or `[AFK]`.
   - **Smart/Dumb context** → `Codebase touchpoints`: after listing files, add a context budget estimate.
   - **Advanced questioning** → `Assumptions` and `Open questions` sections.
   - **Mutual understanding** → reflected throughout, especially in `Problem/why`, `Solution`, and `User stories`.

4. **Write the PRD**  
   Save the completed template as `.skillgrid/prd/PRD<NN>_$ARGUMENTS.md`. Do **not** add YAML frontmatter; keep the exact template format.

5. **Report back**  
   Summarize: file path, number of slices, `[AFK]`/`[HITL]` count, context budget note.

## Rules
- The `Status` field must be `draft`.
- Only use `[HITL]` when a non-deferrable human decision is required; otherwise use `[AFK]`.
- The implementation tasks list is the single source of truth for vertical slices; each slice must correspond to a user story or goal.