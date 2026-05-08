#!/usr/bin/env bash
set -euo pipefail

handoff_path="${1:-}"
if [[ -z "${handoff_path}" ]]; then
  echo "Usage: $0 <handoff-file-path>"
  exit 1
fi

if [[ ! -f "${handoff_path}" ]]; then
  echo "Error: file not found: ${handoff_path}"
  exit 1
fi

failures=0

contains_required_heading() {
  local heading="$1"
  if ! rg --fixed-strings --quiet "${heading}" "${handoff_path}"; then
    echo "Missing section: ${heading}"
    failures=$((failures + 1))
  fi
}

contains_required_heading "## Goal"
contains_required_heading "## Failed approaches"
if ! rg --fixed-strings --quiet "## Next steps" "${handoff_path}" && \
   ! rg --fixed-strings --quiet "## Immediate next steps" "${handoff_path}"; then
  echo "Missing section: ## Next steps (or ## Immediate next steps)"
  failures=$((failures + 1))
fi

if rg --fixed-strings --quiet "<" "${handoff_path}"; then
  echo "Template placeholders remain (found '<')."
  failures=$((failures + 1))
fi

line_count="$(wc -l < "${handoff_path}" | tr -d ' ')"
if [[ "${line_count}" -lt 15 ]]; then
  echo "Handoff seems too short (${line_count} lines)."
  failures=$((failures + 1))
fi

if rg -i --quiet "(api[_-]?key|secret|password|token\s*=|private[_-]?key)" "${handoff_path}"; then
  echo "Potential secret-like content detected. Review required."
  failures=$((failures + 1))
fi

if [[ "${failures}" -gt 0 ]]; then
  echo "VALIDATION: FAIL (${failures} issue(s))"
  exit 2
fi

echo "VALIDATION: PASS"
