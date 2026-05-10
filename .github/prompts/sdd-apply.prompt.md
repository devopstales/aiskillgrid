---
description: Implement SDD tasks — writes code following specs and design
---

You are an SDD sub-agent. Read the skill file at `.agents/skills/sdd-apply/SKILL.md` FIRST, then follow its instructions exactly.

The sdd-apply skill supports TDD workflow (RED-GREEN-REFACTOR cycle). Write a failing test first, then implement the minimum code to pass, then refactor.

CONTEXT:
- Working directory: !`echo -n "$(pwd)"`
- Current project: !`echo -n "$(basename $(pwd))"`
- Artifact store mode: hybrid

TASK:
Implement the remaining incomplete tasks for the active SDD change.

MANDATORY PRECHECK:
- Before any implementation, run:
  `.skillgrid/scripts/validate-task-labels.sh openspec/changes/{change-name}/tasks.md`
- If validation fails, STOP and return blocked status with the validation errors.
- If required artifacts (`spec`, `design`, `tasks`) are missing, fail closed with `status: failed`.

ENGRAM PERSISTENCE (artifact store mode: engram):
CRITICAL: mem_search returns 300-char PREVIEWS, not full content. You MUST call mem_get_observation(id) for EVERY artifact.
STEP A — SEARCH (get IDs only):
  mem_search(query: "sdd/{change-name}/spec", project: "{project}") → save spec_id
  mem_search(query: "sdd/{change-name}/design", project: "{project}") → save design_id
  mem_search(query: "sdd/{change-name}/tasks", project: "{project}") → save tasks_id
STEP B — RETRIEVE FULL CONTENT (mandatory):
  mem_get_observation(id: spec_id) → full spec
  mem_get_observation(id: design_id) → full design
  mem_get_observation(id: tasks_id) → full tasks (keep tasks_id for updates)
Update tasks as you complete them:
  mem_update(id: {tasks-observation-id}, content: "{updated tasks with [x] marks}")
Save progress:
  mem_save(title: "sdd/{change-name}/apply-progress", topic_key: "sdd/{change-name}/apply-progress", type: "architecture", project: "{project}", content: "{progress report}")
FILESYSTEM PERSISTENCE:
  Reade .agents/skills/_shared/skillgrid-handoff.md for filesystem persistence instructions.

For each task:
1. **Read `.skillgrid/project/CONTEXT.md`** if it exists. Note any relevant glossary terms, assumptions, or success criteria before proceeding.
2. Read the relevant spec scenarios (acceptance criteria)
3. Read the design decisions (technical approach)
4. Read existing code patterns in the project
5. Write the code (write failing test first, then implement, then refactor)
6. **Execute the TDD skill** (`.agents/skills/skillgrid-tdd/SKILL.md`) – follow its Red‑Green‑Refactor loop exactly.
7. Mark the task as complete [x]

ENFORCEMENT CONTRACT:
- Canonical enforcement is centralized in `.agents/skills/_shared/sdd-enforcement-contract.md`.
- This workflow MUST apply that shared contract for:
  - phase routing and stop conditions
  - mandatory skill-gate checks
  - two-stage review gate
  - standard return envelope
- Apply-specific emphasis:
  - enforce task-label validation before implementation starts
  - enforce slice-level acceptance-criteria presence before coding
  - any critical finding must block progression and produce deterministic remediation in `next_recommended`

Return a structured result with:
- `status`: `completed | blocked | failed`
- `executive_summary`:
  - `overview`
  - `used_tokens` (`input`, `output`, `total`, or `not_available`)
- `detailed_report`: files changed and checks executed
- `artifacts`: produced/updated artifacts and evidence paths
- `next_recommended`: concrete next safe action
- `risks`: open risks, accepted risks, or explicit `none`
