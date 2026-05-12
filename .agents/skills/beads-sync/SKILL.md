---
name: beads-sync
description: PROACTIVELY convert approved OpenSpec specs into Beads issues when user applies a change. Creates trackable work with dependencies, discovers gaps, and maintains audit trail.
---

# OpenSpec to Beads Bridge

## Inputs

| Variable | Required | Description |
|----------|----------|-------------|
| `CHANGE_ID` | Yes | OpenSpec change name (e.g., "add-payment") |

You are a workflow orchestrator that transforms planning artifacts (OpenSpec) into executable work streams (Beads), maintaining traceability and discovering implementation gaps proactively.

## When This Skill Activates

- After `/sdd-brainstorm` completes
- When user runs `openspec apply <change>`
- When user says "implement this spec" or "let's go"

## Core Philosophy

OpenSpec defines WHAT to build (immutable contract).
Beads tracks execution STATE (mutable reality).

Bridge the gap intelligently, not mechanically.

## Execution Process

### 1. Context Gathering

Read and UNDERSTAND:
- `openspec/changes/<change>/proposal.md` → WHY we're doing this
- `openspec/changes/<change>/tasks.md` → WHAT needs to happen
- `openspec/changes/<change>/design.md` (if exists) → HOW it's architected

### 2. Critical Analysis

**Ask yourself:**
1. Are tasks well-scoped? If too broad, flag it
2. Are dependencies obvious? DB migrations before API endpoints
3. What's missing? Testing? Documentation? Deployment?
4. What will go wrong? Identify tasks that will discover more work

**If you find problems, TELL THE USER BEFORE CONVERTING**

### 3. Smart Conversion

#### Priority Logic:
- **P0**: Infrastructure/setup, migrations, config, schema
- **P1**: Core business logic, implementation, tests, docs
- **P2**: UI polish, nice-to-haves

#### Type Detection:
- "setup", "configure", "refactor" → `task`
- "implement", "add", "create" → `feature`
- "document", "cleanup" → `chore`

#### Complexity:
- "setup", "simple" → `complexity:low`
- "implement", "integrate" → `complexity:medium`
- "refactor", "migrate" → `complexity:high`

### 4. Dependency Chain Construction

**Auto-detect dependencies:**
- Sequential within category (Task 1.2 blocks on 1.1)
- DB/migrations → block → API/business logic
- Config/setup → block → implementation
- Implementation → related → tests (NOT blocks - parallel work)

### 5. Proactive Gap Detection

**Auto-detect missing:**
- Rollback plans for migrations
- Rate limiting for API endpoints
- Error handling for external services
- Monitoring/metrics for features
- Tests for implementations

### 6. Create Tracking Infrastructure

```bash
# Read beads_enable from .skillgrid/config.json
BEADS_ENABLED=$(jq -r '.beads_enable // false' .skillgrid/config.json 2>/dev/null || echo "false")

if [ "${BEADS_ENABLED}" != "true" ]; then
    echo "⚠️ beads_enable not set or false in .skillgrid/config.json, skipping Beads sync"
    exit 0
fi

CHANGE_PATH="openspec/changes/${CHANGE_ID}"

if [ ! -d "$CHANGE_PATH" ]; then
    echo "❌ Change '${CHANGE_ID}' not found"
    exit 1
fi

# Create epic issue
EPIC_ID=$(bd create "🎯 [EPIC] ${CHANGE_ID}" \
    --type epic \
    --priority 0 \
    --labels "spec:${CHANGE_ID},epic,openspec" \
    --description "$(cat "${CHANGE_PATH}/proposal.md" 2>/dev/null || echo 'No proposal.md')" \
    --print-id 2>/dev/null)

echo "✅ Created epic: ${EPIC_ID}"

# Parse tasks and convert with intelligence
if [ -f "${CHANGE_PATH}/tasks.md" ]; then
    while IFS= read -r line; do
        if [[ "$line" =~ ^[[:space:]]*([*-]|[0-9]+\.)[[:space:]]+(.+) ]]; then
            TASK_DESC="${BASH_REMATCH[2]}"
            
            # Smart priority detection
            if [[ "$TASK_DESC" =~ (migration|schema|config|setup|infrastructure) ]]; then
                PRIORITY=0
            elif [[ "$TASK_DESC" =~ (test|document) ]]; then
                PRIORITY=1
            else
                PRIORITY=1
            fi
            
            # Smart type detection
            if [[ "$TASK_DESC" =~ (setup|configure|refactor|test) ]]; then
                TYPE="task"
            elif [[ "$TASK_DESC" =~ (implement|add|create) ]]; then
                TYPE="feature"
            else
                TYPE="task"
            fi
            
            # Complexity detection
            if [[ "$TASK_DESC" =~ (simple|setup) ]]; then
                COMPLEXITY="low"
            elif [[ "$TASK_DESC" =~ (refactor|migrate) ]]; then
                COMPLEXITY="high"
            else
                COMPLEXITY="medium"
            fi
            
            # Create task with labels
            TASK_ID=$(bd create "${TASK_DESC}" \
                --type "${TYPE}" \
                --priority "${PRIORITY}" \
                --labels "spec:${CHANGE_ID},complexity:${COMPLEXITY}" \
                --print-id 2>/dev/null)
            
            # Link to epic
            bd dep add "${TASK_ID}" "${EPIC_ID}" --type parent-child 2>/dev/null
            
            echo "  ➕ Task: ${TASK_ID} - ${TASK_DESC} (p${PRIORITY}, ${TYPE})"
        fi
    done < "${CHANGE_PATH}/tasks.md"
else
    echo "⚠️ No tasks.md found in ${CHANGE_PATH}"
fi

# Proactive gap detection
echo ""
echo "🔍 Checking for common gaps..."

# Check for migration without rollback
if grep -qi "migration" "${CHANGE_PATH}/tasks.md" 2>/dev/null; then
    if ! grep -qi "rollback" "${CHANGE_PATH}/tasks.md" 2>/dev/null; then
        GAP_ID=$(bd create "[${CHANGE_ID}] DISCOVERED: Add rollback migration" \
            --type task \
            --priority 1 \
            --labels "spec:${CHANGE_ID},discovered,proactive" \
            --print-id 2>/dev/null)
        echo "  🔍 Gap detected: missing rollback for migration"
    fi
fi

# Check for API without tests
if grep -qi "api\|endpoint" "${CHANGE_PATH}/tasks.md" 2>/dev/null; then
    if ! grep -qi "test\|spec" "${CHANGE_PATH}/tasks.md" 2>/dev/null; then
        GAP_ID=$(bd create "[${CHANGE_ID}] DISCOVERED: Add API tests" \
            --type task \
            --priority 1 \
            --labels "spec:${CHANGE_ID},discovered,proactive" \
            --print-id 2>/dev/null)
        echo "  🔍 Gap detected: missing tests for API"
    fi
fi

echo ""
echo "✅ Beads sync complete. Run 'bd ready' to see available work."
```

## Anti-Patterns to AVOID

- Don't blindly copy tasks 1:1
- Don't ignore missing test tasks
- Don't make everything high priority
- Don't create 50 issues for 5 tasks
- Don't lose the "why" - reference proposal.md in epic description

## Reference

For advanced patterns and templates, see the original implementation:
- https://github.com/lucastamoios/celeiro/tree/master/.claude/skills/openspec-to-beads