---
name: beads-sync
description: >
  Synchronizes AISkillGrid markdown artifacts (PRDs, tasks.md, slice specs) with Beads graph database for dependency-aware task scheduling and agent memory augmentation.
license: Apache-2.0
metadata:
  author: devopstales
  version: "1.0"
---


## Trigger Conditions
- After `/skillgrid-breakdown` completes successfully
- Before `/skillgrid-apply` starts work on a new slice (optional pre-flight check)
- Manual invocation: `/skillgrid-run beads-sync --change <change-id>`

## Inputs
| Parameter | Required | Default | Description |
|-----------|----------|---------|-------------|
| `--change` | Yes | - | OpenSpec change ID (e.g., `auth-mfa`) |
| `--dry-run` | No | `false` | Preview sync actions without executing |
| `--force` | No | `false` | Overwrite existing Beads IDs in slice headers |
| `--sync-back` | No | `false` | Update OpenSpec tasks.md from Beads state |

## Environment
The skill requires the following environment variables:
- `BEADS_ENABLED` – Set to `true` to enable Beads integration.
- `BD_CLI` – Path to Beads CLI (default: `bd`).
- `BEADS_DB_PATH` – Path to Beads database (default: `./.beads`).

## Process

### Phase 1: Parse Markdown Artifacts
```bash
CHANGE_ID="${1}"
TASKS_FILE="openspec/changes/${CHANGE_ID}/tasks.md"
PRD_FILE=".skillgrid/prd/$(grep -l "change: ${CHANGE_ID}" .skillgrid/prd/*.md | head -1)"

# Extract epic ID from INDEX.md if available
EPIC_ID=$(grep -A5 "## ${CHANGE_ID}" .skillgrid/prd/INDEX.md | grep "beads:" | sed 's/.*beads:\(bd-[a-z0-9]*\).*/\1/')

# If no epic exists in Beads, create one
if [ -z "$EPIC_ID" ] && [ "$BEADS_ENABLED" = "true" ]; then
  EPIC_ID=$($BD_CLI create "Epic: $(basename "$PRD_FILE" .md)" --type epic --json | jq -r '.id')
  # Write back to INDEX.md for idempotency
  sed -i "/## ${CHANGE_ID}/a <!-- beads:${EPIC_ID} -->" .skillgrid/prd/INDEX.md
fi
```

### Phase 2: Map OpenSpec `tasks.md` to Beads Tasks
For each top-level task in `tasks.md` (e.g., `- [ ] 1. Implement API endpoint`), create or update a corresponding Beads task.

```bash
# Read tasks.md line by line
# Use a state machine to detect task list items and their descriptions
# For each task, extract description and any existing Beads ID comment

while IFS= read -r line; do
  if [[ "$line" =~ ^[[:space:]]*-[[:space:]]*\[[[:space:]]*\][[:space:]]*(.*)$ ]]; then
    TASK_DESC="${BASH_REMATCH[1]}"
    
    # Check if line already has a beads: comment
    if [[ "$line" =~ beads:\(bd-[a-z0-9]+\) ]]; then
      TASK_ID="${BASH_REMATCH[1]}"
      # Update existing task
      $BD_CLI update "$TASK_ID" --description "$TASK_DESC" --priority p1
    else
      # Create new task
      TASK_ID=$($BD_CLI create "$TASK_DESC" \
        --type task \
        --priority p1 \
        --parent "$EPIC_ID" \
        --label "openspec:$CHANGE_ID" \
        --label "source:tasks.md" \
        --description "$(generate_spec_context "$CHANGE_ID" "$TASK_DESC")" \
        --json | jq -r '.id')
      
      # Append beads ID as HTML comment for idempotency
      sed -i "s/^- \[ \] $TASK_DESC/- \[ \] $TASK_DESC <!-- beads:$TASK_ID -->/" "$TASKS_FILE"
    fi
  fi
done < "$TASKS_FILE"
```

### Phase 3: Sync Slice Specs to Beads Checklists
For each slice spec under `openspec/changes/$CHANGE_ID/specs/*/spec.md`, extract requirements and acceptance criteria, then attach them as checklists to the appropriate Beads task.

```bash
for SPEC_FILE in openspec/changes/$CHANGE_ID/specs/*/spec.md; do
  SLICE=$(basename "$(dirname "$SPEC_FILE")")
  
  # Find which Beads task corresponds to this slice (by label or parent task)
  TASK_ID=$($BD_CLI list --label "openspec:$CHANGE_ID" --label "slice:$SLICE" --json | jq -r '.[0].id // empty')
  
  if [ -z "$TASK_ID" ]; then
    # Create a wrapper task for the slice
    TASK_ID=$($BD_CLI create "Slice: $SLICE" \
      --type task \
      --priority p1 \
      --parent "$EPIC_ID" \
      --label "openspec:$CHANGE_ID" \
      --label "slice:$SLICE" \
      --json | jq -r '.id')
  fi
  
  # Extract requirements (lines starting with - or ### Requirement:)
  REQUIREMENTS=$(grep -E '^- |### Requirement:' "$SPEC_FILE" | sed 's/^/- /')
  
  # Add as checklist items
  echo "$REQUIREMENTS" | while IFS= read -r req; do
    $BD_CLI checklist add "$TASK_ID" "$req"
  done
done
```

### Phase 4: Bidirectional Sync (Optional)
When `--sync-back` is provided, update OpenSpec `tasks.md` checkboxes (`[ ]` → `[x]`) based on Beads task completion status.

```bash
if [ "$SYNC_BACK" = "true" ]; then
  # Get all tasks for this change that are closed
  CLOSED_TASKS=$($BD_CLI list --label "openspec:$CHANGE_ID" --status closed --json | jq -r '.[].id')
  
  for TASK_ID in $CLOSED_TASKS; do
    # Find the matching line in tasks.md (by beads: comment or description)
    TASK_DESC=$($BD_CLI get "$TASK_ID" --json | jq -r '.title')
    sed -i "s/^- \[ \] $TASK_DESC <!-- beads:$TASK_ID -->/- \[x\] $TASK_DESC <!-- beads:$TASK_ID -->/" "$TASKS_FILE"
  done
fi
```

### Helper Function: Generate Spec Context
Used in Phase 2 to embed full acceptance criteria into Beads task description.

```bash
generate_spec_context() {
  local CHANGE_ID="$1"
  local TASK_DESC="$2"
  
  # Find the slice spec that contains this task description
  local SPEC_FILE=$(grep -l "$TASK_DESC" openspec/changes/$CHANGE_ID/specs/*/spec.md | head -1)
  
  if [ -n "$SPEC_FILE" ]; then
    cat <<EOF
## Spec Reference
$SPEC_FILE

## Requirements
$(grep -E '^- |### Requirement:' "$SPEC_FILE")

## Acceptance Criteria
$(awk '/## Acceptance Criteria/{flag=1; next} /##/{flag=0} flag' "$SPEC_FILE")
EOF
  else
    echo "No spec reference found for this task."
  fi
}
```

## Usage Examples

### After creating a new OpenSpec change

```bash
/skillgrid-run beads-sync --change add-user-description
```

### Preview changes without modifying Beads

```bash
/skillgrid-run beads-sync --change add-user-description --dry-run
```

### Update OpenSpec tasks.md from Beads (e.g., after completing work)

```bash
/skillgrid-run beads-sync --change add-user-description --sync-back
```


## Exit Codes
| Code | Meaning |
|------|---------|
| 0 | Success, no errors |
| 1 | Beads disabled or CLI not found |
| 2 | Missing required input (`--change`) |
| 3 | OpenSpec change directory not found |
| 4 | Beads command failed (check database) |

## Dependencies
- **Beads CLI** (`bd`) – must be installed and accessible in `$PATH` or via `$BD_CLI`
- **jq** – for JSON parsing
- **bash** 4.0+ – for associative arrays (if needed)

## Notes
- The skill is idempotent: running it multiple times will update existing Beads tasks rather than duplicating them.
- Beads task IDs are stored as HTML comments inside `tasks.md` and `INDEX.md` to survive human edits and git commits.
- The `--dry-run` flag only prints what would be changed; it does not touch Beads or markdown files.