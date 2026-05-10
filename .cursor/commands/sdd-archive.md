---
name: /sdd-archive
id: sdd-archive
category: Workflow
description: Archive a completed SDD change — syncs specs and closes the cycle
agent: odin
subtask: true
---

You are an SDD sub-agent. Read the skill file at `.agents/skills/sdd-archive/SKILL.md` FIRST, then follow its instructions exactly.

CONTEXT:
- Working directory: !`echo -n "$(pwd)"`
- Current project: !`echo -n "$(basename $(pwd))"`
- Artifact store mode: hybrid

TASK:
Archive the active SDD change. Read the verification report first to confirm the change is ready. Then:

MANDATORY PRECHECK (CI-ready default, fail closed):
- Confirm verification evidence exists for the active change before archive.
- If verification report/artifacts are missing or indicate unresolved critical findings, return `status: failed` and do not archive.
- If needed, route to `/sdd-verify` first and set `next_recommended` accordingly.

ENGRAM PERSISTENCE (artifact store mode: engram):
CRITICAL: mem_search returns 300-char PREVIEWS, not full content. You MUST call mem_get_observation(id) for EVERY artifact.
STEP A — SEARCH (get IDs only):
  mem_search(query: "sdd/{change-name}/proposal", project: "{project}") → save proposal_id
  mem_search(query: "sdd/{change-name}/spec", project: "{project}") → save spec_id
  mem_search(query: "sdd/{change-name}/design", project: "{project}") → save design_id
  mem_search(query: "sdd/{change-name}/tasks", project: "{project}") → save tasks_id
  mem_search(query: "sdd/{change-name}/verify-report", project: "{project}") → save verify_id
STEP B — RETRIEVE FULL CONTENT (mandatory):
  mem_get_observation(id: proposal_id) → full proposal
  mem_get_observation(id: spec_id) → full spec
  mem_get_observation(id: design_id) → full design
  mem_get_observation(id: tasks_id) → full tasks
  mem_get_observation(id: verify_id) → full verification report
Record all observation IDs in the archive report for traceability.
Save:
  mem_save(title: "sdd/{change-name}/archive-report", topic_key: "sdd/{change-name}/archive-report", type: "architecture", project: "{project}", content: "{archive report with observation IDs}")
FILESYSTEM PERSISTENCE:
  Read `.agents/skills/_shared/skillgrid-handoff.md` for filesystem persistence instructions.

Then:
1. **Read `.skillgrid/project/CONTEXT.md`** if it exists. Note any relevant glossary terms, assumptions, or success criteria before proceeding.
2. Sync delta specs into main specs (source of truth)
3. Move the change folder to archive with date prefix
4. Verify the archive is complete

Return a structured result with: status, executive_summary, artifacts, and next_recommended.
