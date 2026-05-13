#!/usr/bin/env bash

# Cursor hook for tmux-agent-status.
# Cursor does not have a native lifecycle hook system, so this hook
# is designed to be triggered via Cursor commands or manual invocation.
#
# Usage:
#   1. Add a Cursor command in .cursor/commands/agent-status.json that runs:
#      bash <path>/cursor-hook.sh UserPromptSubmit
#   2. Or wire into your workflow via a shell alias or keybinding.
#
# File-based fallback: if no TMUX session is detected, writes to
# ~/.cache/tmux-agent-status/cursor.status as a session-level marker.

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
# shellcheck source=lib.sh
source "$SCRIPT_DIR/lib.sh"

AGENT_NAME="cursor"
init_status_dirs
drain_stdin

TMUX_SESSION=$(get_tmux_session) || {
  # Fallback: use a named "cursor" session when not in tmux.
  TMUX_SESSION="cursor"
}
HOOK_TYPE="${1:-}"

case "$HOOK_TYPE" in
  UserPromptSubmit|CommandRun)
    # User submitted input or ran a command — cancel wait/park, mark working.
    clear_interaction_overrides "$TMUX_SESSION"
    set_status "$TMUX_SESSION" "working"
    mark_refresh
    ;;
  PreToolUse|ToolUse|BashRun)
    # Agent is executing — mark working but do NOT unpark.
    rm -f "$WAIT_DIR/${TMUX_SESSION}.wait"
    if [ ! -f "$PARKED_DIR/${TMUX_SESSION}.parked" ]; then
      set_status "$TMUX_SESSION" "working"
    fi
    mark_refresh
    ;;
  Stop|SessionEnd|Complete)
    # Cursor has finished.
    set_status "$TMUX_SESSION" "done"
    mark_refresh
    ;;
  SessionStart|Init)
    # Session started.
    if [ ! -f "$WAIT_DIR/${TMUX_SESSION}.wait" ] && [ ! -f "$PARKED_DIR/${TMUX_SESSION}.parked" ]; then
      set_status "$TMUX_SESSION" "done"
      mark_refresh
    fi
    ;;
  Notification|Waiting)
    # Cursor is waiting for user input.
    set_status "$TMUX_SESSION" "done"
    mark_refresh
    ;;
  *)
    # Default: treat unknown event as a working signal.
    set_status "$TMUX_SESSION" "working"
    mark_refresh
    ;;
esac

exit 0
