---
name: isolated-workspace
description: >
  Creates and manages isolated git worktrees for each SDD change.
  Ensures clean baseline, parallel work safety, and automatic cleanup.
  Trigger: Automatically invoked at start of sdd-apply.
license: Apache-2.0
metadata:
  author: aiskillgrid-integration
  version: "1.0"
  enforcement: mandatory
  triggers:
    - "sdd-apply_start"
    - "workspace_required"
---

# Isolated Workspace Management

## Overview

Automatically creates isolated git worktrees for each SDD change. Provides clean, disposable, parallel-safe execution environments without manual git management.

**Core principle:** Every `sdd-apply` session works in its own worktree; cleanup is automatic.

## When to Use

Automatically triggered:
- At the start of `sdd-apply` execution
- When `sdd-parallel-execute` spawns parallel slices
- Before any implementation task begins

Manual usage (if needed):
- `/workspace-create <change-id>`
- `/workspace-status`
- `/workspace-cleanup`

## Workspace Naming

```
Branch:     sdd/<change-id>/<slice-slug>
Worktree:   sdd-worktree-<change-id>-<slice-slug>
Directory:  .worktrees/<change-id>/<slice-slug>  (or configured path)
```

**Examples:**
- Change: `user-auth` → branch: `sdd/user-auth/login-flow`
- Worktree directory: `.worktrees/user-auth/login-flow`

## Lifecycle

### 1. Creation

At `sdd-apply` start:

```bash
# Determine change-id and slice from context
CHANGE_ID=$(get_active_change_id)          # from .skillgrid/state or prompt
SLICE=$(get_active_slice)                  # from task context or "full"

# Check if worktree already exists (resume scenario)
if git worktree list | grep -q "sdd-worktree-${CHANGE_ID}-${SLICE}"; then
  echo "Resuming existing worktree"
  cd "$(git worktree list --porcelain | grep "sdd-worktree-${CHANGE_ID}-${SLICE}" | head -1 | awk '{print $2}')"
else
  # Create fresh worktree from main
  git worktree add ".worktrees/${CHANGE_ID}/${SLICE}" "main"
  cd ".worktrees/${CHANGE_ID}/${SLICE}"

  # Record worktree metadata
  echo "change: ${CHANGE_ID}" > .worktree-metadata
  echo "slice: ${SLICE}" >> .worktree-metadata
  echo "created: $(date -Iseconds)" >> .worktree-metadata
  echo "branch: sdd/${CHANGE_ID}/${SLICE}" >> .worktree-metadata
fi
```

### 2. Baseline Verification

Before starting tasks:

```bash
# Run full test suite to ensure clean baseline
if [ -f package.json ]; then
  npm test 2>&1 | tee /tmp/baseline-test.log
  if [ ${PIPESTATUS[0]} -ne 0 ]; then
    echo "ERROR: Baseline tests failing. Fix main branch before proceeding."
    exit 1
  fi
elif [ -f Makefile ]; then
  make test 2>&1 | tee /tmp/baseline-test.log
  if [ ${PIPESTATUS[0]} -ne 0 ]; then
    echo "ERROR: Baseline tests failing. Fix main branch before proceeding."
    exit 1
  fi
else
  echo "WARNING: No test command detected. Skipping baseline check."
fi
```

### 3. Execution

All task execution happens within the worktree. The parent session maintains references to:
- Worktree path
- Branch name
- Original HEAD commit

### 4. Pre-Merge Verification

Before `sdd-archive`:

```bash
# Re-run full test suite in worktree
# Ensure no stray uncommitted changes
# Verify all tasks completed
```

### 5. Cleanup

Worktree teardown options (user choice via `sdd-archive`):

**Merge to main:**
```bash
git checkout main
git merge --no-ff "sdd/${CHANGE_ID}/${SLICE}" -m "Merge SDD change: ${CHANGE_ID}"
# Optionally: delete branch after merge
git branch -d "sdd/${CHANGE_ID}/${SLICE}"
```

**Open PR:**
```bash
git push origin "sdd/${CHANGE_ID}/${SLICE}"
# Create PR via gh CLI or platform API
```

**Keep branch:**
```bash
# Worktree removed from local, branch remains on remote
echo "Branch preserved: sdd/${CHANGE_ID}/${SLICE}"
```

**Discard:**
```bash
git worktree remove ".worktrees/${CHANGE_ID}/${SLICE}" --force
git branch -D "sdd/${CHANGE_ID}/${SLICE}" 2>/dev/null || true
```

## Configuration

`.skillgrid/config.json`:
```json
{
  "workspace": {
    "enabled": true,
    "auto_create": true,
    "auto_cleanup": "prompt",           // "always" | "never" | "prompt"
    "worktree_root": ".worktrees",       // relative to repo root
    "baseline_tests_required": true,
    "resume_existing": true
  }
}
```

## Parallel Work Safety

Multiple `sdd-apply` sessions can run concurrently if their slices are independent:
- Each gets its own worktree
- No file conflicts
- Independent test runs
- Merge coordination happens at archive time

**Detection:** Check for overlapping file modifications across parallel slices via `git diff --name-only` comparison during pre-merge verification.

## Offline / Submodule Scenarios

If worktrees unavailable (e.g., submodule, bare repo):
- Fall back to branch switching in main worktree
- Stash uncommitted changes before switch
- Warn user about limitations
- Document constraint in `.worktree-metadata`

## Error Recovery

**Stuck worktree:**
```bash
# If worktree process is dead but lock remains
git worktree prune
git worktree list --porcelain | grep "lock" | awk '{print $2}' | xargs -r rm -rf
```

**Orphaned branches:**
```bash
# Clean up branches without worktrees
git branch -vv | grep "\[gone\]" | awk '{print $1}' | xargs git branch -D
```

## Evidence Capture

Worktree metadata includes:
- `.worktree-metadata` — creation time, branch, slice, change-id
- `.worktree-session-log` — session start/stop timestamps, task completion markers
- `tdd-evidence/` — per-task test logs and code diffs

These are included in `pre-merge-verification` report.

## Integration

This skill is invoked by:
- `sdd-apply` command — mandatory pre-step
- `sdd-parallel-execute` — one worktree per parallel slice

It prepares the environment for:
- `sequential-agent-executor` — runs tasks inside worktree
- `granular-planning` — may inspect codebase from clean worktree

## See Also

- Git worktree documentation: `git help worktree`
- Related: `sdd-apply` execution orchestration
- Cleanup hook: `pre-merge-verification` triggers removal
