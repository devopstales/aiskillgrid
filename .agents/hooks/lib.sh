#!/usr/bin/env bash
# Shared library for tmux-agent-status hooks.
# Sourced by each agent-specific hook; never executed directly.

STATUS_DIR="${STATUS_DIR:-$HOME/.cache/tmux-agent-status}"
WAIT_DIR="$STATUS_DIR/wait"
PARKED_DIR="$STATUS_DIR/parked"
PANE_DIR="$STATUS_DIR/panes"
REFRESH_FILE="$STATUS_DIR/.sidebar-refresh"

# ── Initialization ──────────────────────────────────────────────

init_status_dirs() {
  mkdir -p "$STATUS_DIR" "$WAIT_DIR" "$PARKED_DIR" "$PANE_DIR"
  [ -f "$REFRESH_FILE" ] || : > "$REFRESH_FILE"
}

# ── Session / Pane Detection ────────────────────────────────────

in_remote_session() {
  [ -n "${SSH_CONNECTION:-}" ] || [ -n "${SSH_TTY:-}" ]
}

get_tmux_session() {
  local tmux_session=""

  if [ -n "${TMUX:-}" ] || in_remote_session; then
    tmux_session=$(tmux display-message -p '#{session_name}' 2>/dev/null)

    if [ -z "$tmux_session" ]; then
      if in_remote_session; then
        tmux_session=$(hostname -s 2>/dev/null || echo "unknown")
      elif [ -n "${TMUX:-}" ]; then
        local socket_path="${TMUX%%,*}"
        tmux_session=$(basename "$socket_path")
      fi
    fi
  fi

  [ -n "$tmux_session" ] || return 1
  printf '%s\n' "$tmux_session"
}

get_pane_id() {
  if [ -n "${TMUX_PANE:-}" ]; then
    printf '%s\n' "$TMUX_PANE"
    return 0
  fi
  return 1
}

# ── Status Writing ──────────────────────────────────────────────

set_status() {
  local tmux_session="$1"
  local requested_status="$2"
  local session_status="$requested_status"
  local status_file="$STATUS_DIR/${tmux_session}.status"
  local agent_file="$STATUS_DIR/${tmux_session}.agent"
  local remote_status_file="$STATUS_DIR/${tmux_session}-remote.status"

  # Write session-level status and agent identifier
  echo "$session_status" > "$status_file"
  echo "${AGENT_NAME:-unknown}" > "$agent_file"

  # Pane-level status
  if [ -n "${TMUX_PANE:-}" ]; then
    local pane_file="$PANE_DIR/${tmux_session}_${TMUX_PANE}.status"
    local pane_agent_file="$PANE_DIR/${tmux_session}_${TMUX_PANE}.agent"
    echo "$requested_status" > "$pane_file"
    echo "${AGENT_NAME:-unknown}" > "$pane_agent_file"

    # Aggregate: session is "working" if any pane is working
    session_status="done"
    local existing_pane_file=""
    for existing_pane_file in "$PANE_DIR/${tmux_session}_"*.status; do
      [ -f "$existing_pane_file" ] || continue
      local pane_status=""
      pane_status=$(cat "$existing_pane_file" 2>/dev/null || echo "")
      case "$pane_status" in
        working)
          session_status="working"
          break
          ;;
        wait)
          if [ "$session_status" != "working" ]; then
            session_status="wait"
          fi
          ;;
      esac
    done
    echo "$session_status" > "$status_file"
  fi

  # Remote session status
  if in_remote_session; then
    echo "$session_status" > "$remote_status_file" 2>/dev/null
  fi
}

# ── Wait / Park Management ──────────────────────────────────────

clear_interaction_overrides() {
  local tmux_session="$1"
  local session_wait_file="$WAIT_DIR/${tmux_session}.wait"
  local session_parked_file="$PARKED_DIR/${tmux_session}.parked"

  if [ -f "$session_wait_file" ]; then
    rm -f "$session_wait_file" "$WAIT_DIR/${tmux_session}_"*.wait 2>/dev/null
  elif [ -n "${TMUX_PANE:-}" ]; then
    rm -f "$WAIT_DIR/${tmux_session}_${TMUX_PANE}.wait"
  fi

  if [ -f "$session_parked_file" ]; then
    rm -f "$session_parked_file" "$PARKED_DIR/${tmux_session}_"*.parked 2>/dev/null
  elif [ -n "${TMUX_PANE:-}" ]; then
    rm -f "$PARKED_DIR/${tmux_session}_${TMUX_PANE}.parked"
  fi
}

# ── Refresh Signal ──────────────────────────────────────────────

mark_refresh() {
  touch "$REFRESH_FILE" 2>/dev/null || true
}

# ── Stdin Drain ─────────────────────────────────────────────────
# All hooks must drain stdin so the host agent can close cleanly.

drain_stdin() {
  cat >/dev/null 2>&1 || true
}
