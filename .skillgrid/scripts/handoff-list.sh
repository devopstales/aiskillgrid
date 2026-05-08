#!/usr/bin/env bash
set -euo pipefail

if ! git rev-parse --show-toplevel >/dev/null 2>&1; then
  echo "Error: run inside a git repository."
  exit 1
fi

repo_root="$(git rev-parse --show-toplevel)"
handoff_dir="${repo_root}/.skillgrid/handoffs"

if [[ ! -d "${handoff_dir}" ]]; then
  echo "No handoff directory found at ${handoff_dir}"
  exit 0
fi

ls -1t "${handoff_dir}"/*.md 2>/dev/null || echo "No handoff files found."
