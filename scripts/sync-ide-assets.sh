#!/usr/bin/env bash
# Sync hub IDE mirrors from Cursor sources of truth.
# Usage:
#   ./scripts/sync-ide-assets.sh           # write .kilo, .opencode, .github/prompts, .github/agents
#   ./scripts/sync-ide-assets.sh --check   # exit 1 if any mirror differs (CI)
# Removes command/prompt/agent files under mirrors when the matching source file is deleted.
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

for dest in "$ROOT/.github/agents"; do
  for f in "$ROOT/.cursor/agents"/*.md; do
    sync_pair "$f" "$dest/$(basename "$f")"
  done
  # Drop mirror files removed from .cursor/agents (e.g. disabled personas)
  for f in "$dest"/*.md; do
    [[ -f "$f" ]] || continue
    base="$(basename "$f")"
    if [[ ! -f "$ROOT/.cursor/agents/$base" ]]; then
      if [[ $CHECK -eq 1 ]]; then
        echo "ORPHAN: $f (no longer in .cursor/agents)" >&2
        EXIT=1
      else
        rm -f "$f"
      fi
    fi
  done
done

# Kilo/OpenCode agents use OpenCode's Markdown agent schema:
# - the filename is the agent name
# - access is configured through permission, not comma-separated tools
if ! python3 - "$ROOT" "$CHECK" <<'PY'
import re
import sys
from pathlib import Path

root = Path(sys.argv[1])
check = sys.argv[2] == "1"
cur = root / ".cursor/agents"
dest_dirs = [root / ".kilo/agents", root / ".opencode/agents"]
drift = False


def frontmatter(text: str):
    match = re.match(r"^---\n(.*?)\n---\n?", text, re.DOTALL)
    if not match:
        return None, text
    fields = {}
    for line in match.group(1).splitlines():
        if ":" not in line:
            continue
        key, value = line.split(":", 1)
        fields[key.strip()] = value.strip()
    return fields, text[match.end() :]


def render_agent(src: Path) -> str:
    text = src.read_text(encoding="utf-8")
    fields, body = frontmatter(text)
    if fields is None:
        return text

    description = fields.get("description")
    if not description:
        print(f"sync-ide-assets: missing description: {src}", file=sys.stderr)
        sys.exit(2)

    tools = {tool.strip() for tool in fields.get("tools", "").split(",") if tool.strip()}
    lines = [
        "---",
        f"description: {description}",
        "mode: subagent",
        "permission:",
        "  read: allow",
        "  glob: allow",
        "  grep: allow",
        "  edit: deny",
        "  task: deny",
        f"  bash: {'allow' if 'Bash' in tools else 'deny'}",
    ]
    if "WebSearch" in tools:
        lines.append("  websearch: allow")
    if "WebFetch" in tools:
        lines.append("  webfetch: allow")
    if color := fields.get("color"):
        lines.append(f"color: {color}")
    lines.append("---")
    return "\n".join(lines) + "\n" + body


files = sorted(cur.glob("*.md"))
for dest_dir in dest_dirs:
    for src in files:
        out = render_agent(src)
        dest = dest_dir / src.name
        if check:
            if not dest.exists() or dest.read_text(encoding="utf-8") != out:
                print(f"DRIFT: {dest}", file=sys.stderr)
                drift = True
        else:
            dest.parent.mkdir(parents=True, exist_ok=True)
            dest.write_text(out, encoding="utf-8")

    # Drop mirror files removed from .cursor/agents (e.g. disabled personas)
    for dest in sorted(dest_dir.glob("*.md")):
        src = cur / dest.name
        if not src.exists():
            if check:
                print(f"ORPHAN: {dest} (no longer in .cursor/agents)", file=sys.stderr)
                drift = True
            else:
                dest.unlink()

sys.exit(1 if drift else 0)
PY
then
  EXIT=1
fi

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
