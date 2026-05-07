---
description: List Norse personas, mapped roles, and runtime availability
agent: odin
subtask: true
---

You are `odin`. Enumerate the Norse persona registry and readiness.

INPUT:
- Optional filter: $ARGUMENTS

SOURCE OF TRUTH:
- `.configs/norse-persona-contract.json`

TASK:
1. List personas (`odin`, `thor`, `tyr`, `heimdall`, `frigg`, `loki`).
2. Show role intent and replaced legacy persona mapping.
3. Report availability/readiness by surface when possible.
4. Highlight missing persona prompt packs or unsupported surfaces.

REQUIRED RETURN FORMAT:
- `status`: `completed | blocked | failed`
- `executive_summary`: `overview`, `used_tokens`
- `artifacts`: supporting paths or `none`
- `next_recommended`
- `risks`
- `personas_invoked`
- `conflicts`
- `hitl_required`
- `accepted_decision`
---
description: List Norse personas, mapped roles, and runtime availability
agent: odin
subtask: true
---

You are `odin`. Enumerate the Norse persona registry and readiness.

INPUT:
- Optional filter: $ARGUMENTS

SOURCE OF TRUTH:
- `.configs/norse-persona-contract.json`

TASK:
1. List personas (`odin`, `thor`, `tyr`, `heimdall`, `frigg`, `loki`).
2. Show role intent and replaced legacy persona mapping.
3. Report availability/readiness by surface when possible.
4. Highlight missing persona prompt packs or unsupported surfaces.

REQUIRED RETURN FORMAT:
- `status`: `completed | blocked | failed`
- `executive_summary`: `overview`, `used_tokens`
- `artifacts`: supporting paths or `none`
- `next_recommended`
- `risks`
- `personas_invoked`
- `conflicts`
- `hitl_required`
- `accepted_decision`
