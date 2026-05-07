# SDD Label Gate Contract

Applies to implementation and verification phases using task labels.

## Validator Command

Run before implementation or deep verification:

```bash
.skillgrid/scripts/validate-task-labels.sh openspec/changes/{change-name}/tasks.md
```

## Rules

- Validation failure is a hard gate failure.
- Missing or invalid `[Label: AFK|HITL]` metadata is blocking.
- Record validator output in report artifacts.

## Status Mapping

- Apply phase: return `status: blocked` on validation failure.
- Verify phase: return `status: failed` with CRITICAL gate finding.

## Required Output Fields on Failure

- `status`
- `executive_summary`
- `artifacts` (include validator evidence path/output)
- `next_recommended` (explicit remediation)
- `risks` (include workflow-quality risk)
