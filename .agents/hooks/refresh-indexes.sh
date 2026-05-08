#!/usr/bin/env bash
set -u

# Shared index refresh hook for all agent surfaces.
# Intended trigger points: post-merge, post-pull, and init/bootstrap.

payload="$(cat)"

extract_command() {
  python3 - "$payload" <<'PY'
import json
import sys

raw = sys.argv[1]
cmd = ""
try:
    data = json.loads(raw)
    cmd = (
        data.get("command")
        or data.get("tool_input", {}).get("command")
        or data.get("input", {}).get("command")
        or ""
    )
except Exception:
    cmd = ""
print(cmd)
PY
}

command_text="$(extract_command)"

should_refresh=0
if [[ "$command_text" =~ git[[:space:]]+merge ]] || \
   [[ "$command_text" =~ gh[[:space:]]+pr[[:space:]]+merge ]] || \
   [[ "$command_text" =~ git[[:space:]]+pull ]] || \
   [[ "$command_text" =~ ccc[[:space:]]+init ]]; then
  should_refresh=1
fi

if [[ "$should_refresh" -eq 0 ]]; then
  exit 0
fi

echo "[hook] Refresh policy triggered by: $command_text" >&2

if command -v ccc >/dev/null 2>&1; then
  if ! ccc index >/dev/null 2>&1; then
    echo "[hook] ccc index failed (continuing, fail-open)." >&2
  else
    echo "[hook] ccc index refreshed." >&2
  fi
else
  echo "[hook] ccc not installed; skipping ccc index." >&2
fi

# For GitNexus, refresh only when git history changed or index is missing.
if [[ "$command_text" =~ git[[:space:]]+merge ]] || \
   [[ "$command_text" =~ gh[[:space:]]+pr[[:space:]]+merge ]] || \
   [[ "$command_text" =~ git[[:space:]]+pull ]]; then
  if command -v npx >/dev/null 2>&1; then
    current_head="$(git rev-parse HEAD 2>/dev/null || true)"

    gitnexus_meta_path=""
    if [[ -f ".gitnexus/meta.json" ]]; then
      gitnexus_meta_path=".gitnexus/meta.json"
    fi

    last_commit=""
    had_embeddings="0"
    if [[ -n "$gitnexus_meta_path" ]] && command -v python3 >/dev/null 2>&1; then
      parsed="$(python3 - "$gitnexus_meta_path" <<'PY'
import json
import sys

path = sys.argv[1]
last = ""
emb = "0"
try:
    with open(path, "r", encoding="utf-8") as f:
        data = json.load(f)
    last = data.get("lastCommit") or ""
    stats = data.get("stats") or {}
    emb_count = stats.get("embeddings") or 0
    emb = "1" if emb_count > 0 else "0"
except Exception:
    pass
print(f"{last}\n{emb}")
PY
)"
      last_commit="$(printf "%s" "$parsed" | sed -n '1p')"
      had_embeddings="$(printf "%s" "$parsed" | sed -n '2p')"
    fi

    if [[ -n "$current_head" ]] && [[ -n "$last_commit" ]] && [[ "$current_head" == "$last_commit" ]]; then
      echo "[hook] gitnexus already fresh for HEAD ${current_head:0:7}." >&2
    else
      analyze_args=("gitnexus" "analyze")
      if [[ "$had_embeddings" == "1" ]]; then
        analyze_args+=("--embeddings")
      fi

      if ! npx "${analyze_args[@]}" >/dev/null 2>&1; then
        echo "[hook] gitnexus analyze failed (continuing, fail-open)." >&2
      else
        if [[ "$had_embeddings" == "1" ]]; then
          echo "[hook] gitnexus index refreshed with embeddings." >&2
        else
          echo "[hook] gitnexus index refreshed." >&2
        fi
      fi
    fi
  else
    echo "[hook] npx not installed; skipping gitnexus analyze." >&2
  fi
fi

exit 0
