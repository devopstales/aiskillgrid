#!/usr/bin/env bash
# Sync hub IDE mirrors from Cursor sources of truth.
# Usage:
#   ./scripts/sync-ide-assets.sh           # write .kilo, .opencode, .github/prompts, .github/agents
#   ./scripts/sync-ide-assets.sh --check   # exit 1 if any mirror differs (CI)
# Removes command/prompt files under mirrors when the matching file is deleted from .cursor/commands/.
set -uo pipefail
shopt -s nullglob

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
CHECK=0
if [[ "${1:-}" == "--check" ]]; then
  CHECK=1
fi

EXIT=0

sync_pair() {
  local src="$1" dst="$2"
  if [[ $CHECK -eq 1 ]]; then
    if [[ ! -f "$dst" ]] || ! cmp -s "$src" "$dst"; then
      echo "DRIFT: $dst (expected content from $src)" >&2
      EXIT=1
    fi
  else
    cp "$src" "$dst"
  fi
}

for dest in "$ROOT/.kilo/commands" "$ROOT/.opencode/commands"; do
  for f in "$ROOT/.cursor/commands"/skillgrid-*.md "$ROOT/.cursor/commands"/opsx-*.md; do
    sync_pair "$f" "$dest/$(basename "$f")"
  done
  # Drop mirror files removed from .cursor/commands (e.g. disabled commands)
  for f in "$dest"/skillgrid-*.md "$dest"/opsx-*.md; do
    [[ -f "$f" ]] || continue
    base="$(basename "$f")"
    if [[ ! -f "$ROOT/.cursor/commands/$base" ]]; then
      if [[ $CHECK -eq 1 ]]; then
        echo "ORPHAN: $f (no longer in .cursor/commands)" >&2
        EXIT=1
      else
        rm -f "$f"
      fi
    fi
  done
done

for dest in "$ROOT/.kilo/agents" "$ROOT/.opencode/agents" "$ROOT/.github/agents"; do
  for f in "$ROOT/.cursor/agents"/*.md; do
    sync_pair "$f" "$dest/$(basename "$f")"
  done
done

# GitHub Copilot prompts: description-only frontmatter + body
if ! python3 - "$ROOT" "$CHECK" <<'PY'
import re
import sys
from pathlib import Path

root = Path(sys.argv[1])
check = sys.argv[2] == "1"
cur = root / ".cursor/commands"
gh = root / ".github/prompts"
files = sorted(cur.glob("skillgrid-*.md")) + sorted(cur.glob("opsx-*.md"))
drift = False
for f in files:
    text = f.read_text(encoding="utf-8")
    m = re.match(r"^---\n(.*?)\n---\n", text, re.DOTALL)
    if not m:
        print(f"sync-ide-assets: missing frontmatter: {f}", file=sys.stderr)
        sys.exit(2)
    fm = m.group(1)
    dm = re.search(r"^description:\s*(.+)$", fm, re.MULTILINE)
    if not dm:
        print(f"sync-ide-assets: missing description: {f}", file=sys.stderr)
        sys.exit(2)
    desc = dm.group(1).strip()
    extra_lines = []
    for key in ("allowed-tools", "argument-hint"):
        km = re.search(rf"^{re.escape(key)}:\s*(.+)$", fm, re.MULTILINE)
        if km:
            extra_lines.append(f"{key}: {km.group(1).strip()}")
    body = text[m.end() :]
    extra = ("\n".join(extra_lines) + "\n") if extra_lines else ""
    out = f"---\ndescription: {desc}\n{extra}---\n{body}"
    dest = gh / f.name
    if check:
        if not dest.exists() or dest.read_text(encoding="utf-8") != out:
            print(f"DRIFT: {dest}", file=sys.stderr)
            drift = True
    else:
        dest.parent.mkdir(parents=True, exist_ok=True)
        dest.write_text(out, encoding="utf-8")
# Remove GitHub prompt files for commands deleted from .cursor/commands
for dest in sorted(gh.glob("skillgrid-*.md")) + sorted(gh.glob("opsx-*.md")):
    src = cur / dest.name
    if not src.exists():
        if check:
            print(f"ORPHAN: {dest}", file=sys.stderr)
            drift = True
        else:
            dest.unlink()
sys.exit(1 if drift else 0)
PY
then
  EXIT=1
fi

exit "$EXIT"
