#!/usr/bin/env bash
set -euo pipefail

change_id="${1:-}"
if [[ -z "${change_id}" ]]; then
  echo "Usage: $0 <change-id>"
  exit 1
fi

if ! git rev-parse --show-toplevel >/dev/null 2>&1; then
  echo "Error: run inside a git repository."
  exit 1
fi

repo_root="$(git rev-parse --show-toplevel)"
template="${repo_root}/.skillgrid/templates/template-handoff-registry.md"
target="${repo_root}/.skillgrid/tasks/registry_${change_id}.md"

if [[ ! -f "${template}" ]]; then
  echo "Error: missing template ${template}"
  exit 1
fi

mkdir -p "${repo_root}/.skillgrid/tasks"
cp "${template}" "${target}"

today="$(date -u +%Y-%m-%d)"
sed -i '' "s|<change-id>|${change_id}|g" "${target}"
sed -i '' "s|<YYYY-MM-DD>|${today}|g" "${target}"

echo "${target}"
