#!/usr/bin/env bash
set -euo pipefail

handoff_path="${1:-}"
max_age_days="${2:-3}"

if [[ -z "${handoff_path}" ]]; then
  echo "Usage: $0 <handoff-file-path> [max-age-days]"
  exit 1
fi

if [[ ! -f "${handoff_path}" ]]; then
  echo "Error: file not found: ${handoff_path}"
  exit 1
fi

if ! git rev-parse --show-toplevel >/dev/null 2>&1; then
  echo "Error: run inside a git repository."
  exit 1
fi

if ! [[ "${max_age_days}" =~ ^[0-9]+$ ]]; then
  echo "Error: max-age-days must be an integer."
  exit 1
fi

repo_root="$(git rev-parse --show-toplevel)"
handoff_abs="$(cd "$(dirname "${handoff_path}")" && pwd)/$(basename "${handoff_path}")"

if [[ "${handoff_abs}" != "${repo_root}"* ]]; then
  echo "Error: handoff file must be inside the current repo."
  exit 1
fi

now_epoch="$(date -u +%s)"
file_epoch="$(python3 -c 'import os, sys; print(int(os.path.getmtime(sys.argv[1])))' "${handoff_abs}")"
age_days="$(( (now_epoch - file_epoch) / 86400 ))"

relative_handoff="${handoff_abs#${repo_root}/}"
commits_since="$(git rev-list --count --since="@${file_epoch}" HEAD)"
changed_since_count="$(git status --porcelain | wc -l | tr -d ' ')"

stale_reasons=()
if [[ "${age_days}" -gt "${max_age_days}" ]]; then
  stale_reasons+=("age ${age_days}d > ${max_age_days}d")
fi
if [[ "${commits_since}" -gt 0 ]]; then
  stale_reasons+=("${commits_since} commit(s) since handoff timestamp")
fi
if [[ "${changed_since_count}" -gt 0 ]]; then
  stale_reasons+=("working tree drift detected")
fi

echo "Handoff: ${relative_handoff}"
echo "Age days: ${age_days}"
echo "Commits since handoff timestamp: ${commits_since}"
echo "Working tree drift count: ${changed_since_count}"

if [[ "${#stale_reasons[@]}" -eq 0 ]]; then
  echo "STALENESS: FRESH"
  exit 0
fi

echo "STALENESS: STALE"
for reason in "${stale_reasons[@]}"; do
  echo "- ${reason}"
done
exit 3
