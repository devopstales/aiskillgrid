#!/usr/bin/env bash

# Gemini CLI hook for tmux-agent-status.
# Tracks session/pane state via Gemini CLI lifecycle events.
#
# Setup: add to .gemini/settings.json or ~/.gemini/settings.json:
#   "hooks": {
#     "UserPromptSubmit": "bash <path>/gemini-hook.sh UserPromptSubmit",
#     "PreToolUse": "bash <path>/gemini-hook.sh PreToolUse",
#     "Stop": "bash <path>/gemini-hook.sh Stop",
#     "SessionStart": "bash <path>/gemini-hook.sh SessionStart"
#   }

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
# shellcheck source=lib.sh
source "$SCRIPT_DIR/lib.sh"

AGENT_NAME="gemini"
init_status_dirs
drain_stdin

TMUX_SESSION=$(get_tmux_session) || exit 0
HOOK_TYPE="${1:-}"

case "$HOOK_TYPE" in
  SessionStart)
    if [ ! -f "$WAIT_DIR/${TMUX_SESSION}.wait" ] && [ ! -f "$PARKED_DIR/${TMUX_SESSION}.parked" ]; then
      set_status "$TMUX_SESSION" "done"
      mark_refresh
    fi
    ;;
  UserPromptSubmit)
    clear_interaction_overrides "$TMUX_SESSION"
    set_status "$TMUX_SESSION" "working"
    mark_refresh
    ;;
  PreToolUse|PostToolUse)
    rm -f "$WAIT_DIR/${TMUX_SESSION}.wait"
    if [ ! -f "$PARKED_DIR/${TMUX_SESSION}.parked" ]; then
      set_status "$TMUX_SESSION" "working"
    fi
    mark_refresh
    ;;
  Stop)
    set_status "$TMUX_SESSION" "done"
    mark_refresh
    ;;
  Notification)
    set_status "$TMUX_SESSION" "done"
    mark_refresh
    ;;
esac

exit 0
