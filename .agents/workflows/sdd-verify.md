---
description: Validate implementation matches specs, design, and tasks
agent: sdd-orchestrator
subtask: true
---

You are an SDD sub-agent. Read the skill file at `.agents/skills/sdd-verify/SKILL.md` FIRST, then follow its instructions exactly.

CONTEXT:
- Working directory: !`echo -n "$(pwd)"`
- Current project: !`echo -n "$(basename $(pwd))"`
- Artifact store mode: hybrid

TASK:
Verify the active SDD change. Read the proposal, specs, design, and tasks artifacts. Then:

MANDATORY PRECHECK:
- Run `.skillgrid/scripts/validate-task-labels.sh openspec/changes/{change-name}/tasks.md` before verification.
- If validation fails, report a CRITICAL gate failure and return FAIL.
- If required artifacts (`proposal`, `spec`, `design`, `tasks`) are missing, fail closed with `status: failed`.

ENGRAM PERSISTENCE (artifact store mode: engram):
CRITICAL: mem_search returns 300-char PREVIEWS, not full content. You MUST call mem_get_observation(id) for EVERY artifact.
STEP A — SEARCH (get IDs only):
  mem_search(query: "sdd/{change-name}/spec", project: "{project}") → save spec_id
  mem_search(query: "sdd/{change-name}/design", project: "{project}") → save design_id
  mem_search(query: "sdd/{change-name}/tasks", project: "{project}") → save tasks_id
STEP B — RETRIEVE FULL CONTENT (mandatory):
  mem_get_observation(id: spec_id) → full spec
  mem_get_observation(id: design_id) → full design
  mem_get_observation(id: tasks_id) → full tasks
Save report:
  mem_save(title: "sdd/{change-name}/verify-report", topic_key: "sdd/{change-name}/verify-report", type: "architecture", project: "{project}", content: "{verification report}")
FILESYSTEM PERSISTENCE:
  Reade .agents/skills/_shared/skillgrid-handoff.md for filesystem persistence instructions.

Then:
1. **Read `.skillgrid/project/CONTEXT.md`** if it exists. Note any relevant glossary terms, assumptions, or success criteria before proceeding.
2. Check completeness — are all tasks done?
3. Check correctness — does code match specs?
4. Check coherence — were design decisions followed?
5. Run tests and build (real execution)
6. Build the spec compliance matrix

ENFORCEMENT CONTRACT:
- Canonical enforcement is centralized in `.agents/skills/_shared/sdd-enforcement-contract.md`.
- This workflow MUST apply that shared contract for:
  - phase routing and stop conditions
  - mandatory skill-gate checks
  - two-stage review gate
  - standard return envelope
- Verify-specific progression rule:
  - any critical finding in Stage 1 or Stage 2 => `status: failed` with explicit remediation in `next_recommended`
  - both stages pass => allow progression to `/sdd-archive`

Return a structured verification report with:
- `status`: `completed | blocked | failed`
- `executive_summary`:
  - `overview`
  - `used_tokens` (`input`, `output`, `total`, or `not_available`)
- `detailed_report`: verification matrix, executed checks, and critical findings
- `artifacts`: report paths and evidence paths
- `next_recommended`: concrete next command/action
- `risks`: open/accepted risks or explicit `none`
