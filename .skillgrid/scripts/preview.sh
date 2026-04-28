#!/usr/bin/env bash
# Scaffold an HTML file under .skillgrid/preview/ for A/B/C "pick a preview" brainstorms.
#
# Usage:
#   .skillgrid/scripts/preview.sh [slug]
#   .skillgrid/scripts/preview.sh --md [slug]   # markdown stub instead of HTML
#
# Examples:
#   .skillgrid/scripts/preview.sh dashboard-layout
#   .skillgrid/scripts/preview.sh
#
# If slug is omitted, uses preview-YYYY-MM-DD. Writes under .skillgrid/preview/ next
# to this script's .skillgrid directory (works in any project that uses the Skillgrid tree).
set -euo pipefail

usage() {
  cat <<'EOF'
Scaffold a file under .skillgrid/preview/ for A/B/C "pick a preview" brainstorms.

Usage:
  .skillgrid/scripts/preview.sh [slug]
  .skillgrid/scripts/preview.sh --md [slug]   # markdown stub instead of HTML

Examples:
  .skillgrid/scripts/preview.sh dashboard-layout
  .skillgrid/scripts/preview.sh

If slug is omitted, uses preview-YYYY-MM-DD. Refuses to overwrite an existing file;
prints its path. Writes next to the .skillgrid/ directory that contains this script.
EOF
}

if [[ "${1:-}" == "-h" || "${1:-}" == "--help" ]]; then
  usage
  exit 0
fi

SKILLGRID_DIR="$(cd "$(dirname "$0")/.." && pwd)"
PREVIEW_DIR="${SKILLGRID_DIR}/preview"
MODE="html"

if [[ "${1:-}" == "--md" || "${1:-}" == "-m" ]]; then
  MODE="md"
  shift
fi

SLUG="${1:-preview-$(date +%Y-%m-%d)}"
SLUG="${SLUG//[^a-zA-Z0-9._-]/-}"

mkdir -p "$PREVIEW_DIR"

if [[ "$MODE" == "md" ]]; then
  OUT="${PREVIEW_DIR}/${SLUG}.md"
  if [[ -f "$OUT" ]]; then
    echo "Exists (not overwriting): $OUT" >&2
    echo "$OUT"
    exit 0
  fi
  cat >"$OUT" <<'EOF'
# Preview

Short label for this brainstorm preview.

## Option A
- …

## Option B
- …

## Option C
- …

_Pick A, B, C, or name a variant in chat._
EOF
  echo "Wrote: $OUT"
  echo "$OUT"
  exit 0
fi

OUT="${PREVIEW_DIR}/${SLUG}.html"
if [[ -f "$OUT" ]]; then
  echo "Exists (not overwriting): $OUT" >&2
  echo "$OUT"
  exit 0
fi

cat >"$OUT" <<EOF
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>Preview: ${SLUG}</title>
  <style>
    :root { font-family: system-ui, sans-serif; line-height: 1.4; }
    body { max-width: 48rem; margin: 1.5rem auto; padding: 0 1rem; }
    h1 { font-size: 1.15rem; }
    .opt { border: 1px solid #ccc; border-radius: 6px; padding: 0.75rem 1rem; margin: 0.75rem 0; }
    .opt h2 { margin: 0 0 0.35rem; font-size: 1rem; }
    .hint { color: #555; font-size: 0.9rem; }
  </style>
</head>
<body>
  <h1>Preview: ${SLUG}</h1>
  <p class="hint">Replace the placeholders; open this file in a browser. Reply in chat with A, B, C, or edits.</p>
  <div class="opt">
    <h2>Option A</h2>
    <p>…</p>
  </div>
  <div class="opt">
    <h2>Option B</h2>
    <p>…</p>
  </div>
  <div class="opt">
    <h2>Option C</h2>
    <p>…</p>
  </div>
</body>
</html>
EOF
echo "Wrote: $OUT"
echo "$OUT"
