#!/usr/bin/env bash
# Validate that tasks.md uses enforceable AFK/HITL labeling.
#
# Enforced rules:
# - Every actionable checkbox task ("- [ ]" or "- [x]") MUST include "[Label: AFK]" or "[Label: HITL]".
# - Every labeled task MUST include a short label reason using one of:
#   - "[Reason: ...]"
#   - "[Gate: ...]"
#   - "[HITL Reason: ...]"
#   - "[AFK Reason: ...]"
#
# Usage:
#   .skillgrid/scripts/validate-task-labels.sh [path/to/tasks.md]
#
# Exit codes:
#   0 = pass
#   2 = validation failed
#   3 = usage or file error
set -euo pipefail

usage() {
  cat <<'EOF'
Validate enforceable AFK/HITL labels in tasks.md.

Usage:
  .skillgrid/scripts/validate-task-labels.sh [path/to/tasks.md]

Defaults:
  path = openspec/changes/*/tasks.md (newest by mtime)
EOF
}

if [[ "${1:-}" == "-h" || "${1:-}" == "--help" ]]; then
  usage
  exit 0
fi

resolve_tasks_file() {
  if [[ -n "${1:-}" ]]; then
    printf '%s\n' "$1"
    return 0
  fi

  local candidate
  candidate="$(ls -t openspec/changes/*/tasks.md 2>/dev/null | head -n 1 || true)"
  if [[ -z "$candidate" ]]; then
    return 1
  fi
  printf '%s\n' "$candidate"
}

TASKS_FILE="$(resolve_tasks_file "${1:-}")" || {
  echo "ERROR: Could not resolve tasks.md. Pass a file path explicitly." >&2
  exit 3
}

if [[ ! -f "$TASKS_FILE" ]]; then
  echo "ERROR: File not found: $TASKS_FILE" >&2
  exit 3
fi

awk '
BEGIN {
  failures = 0
}

function is_actionable_task(line) {
  # Task checkboxes should be enforced; DECISION tasks are enforced too.
  return (line ~ /^[[:space:]]*-[[:space:]]*\[[ xX]\][[:space:]]+/)
}

function has_label(line) {
  return (line ~ /\[Label:[[:space:]]*(AFK|HITL)\]/)
}

function has_reason(line) {
  if (line ~ /\[Reason:[[:space:]]*[^]]+\]/) return 1
  if (line ~ /\[Gate:[[:space:]]*[^]]+\]/) return 1
  if (line ~ /\[HITL Reason:[[:space:]]*[^]]+\]/) return 1
  if (line ~ /\[AFK Reason:[[:space:]]*[^]]+\]/) return 1
  return 0
}

{
  if (!is_actionable_task($0)) {
    next
  }

  if (!has_label($0)) {
    printf("FAIL line %d: missing [Label: AFK|HITL]\n  %s\n", NR, $0) > "/dev/stderr"
    failures++
  }

  if (!has_reason($0)) {
    printf("FAIL line %d: missing label reason tag ([Reason:], [Gate:], [HITL Reason:], or [AFK Reason:])\n  %s\n", NR, $0) > "/dev/stderr"
    failures++
  }
}

END {
  if (failures > 0) {
    printf("\nValidation failed: %d issue(s) found in %s\n", failures, FILENAME) > "/dev/stderr"
    exit 2
  }

  printf("PASS: enforceable AFK/HITL labels validated in %s\n", FILENAME)
}
' "$TASKS_FILE"
