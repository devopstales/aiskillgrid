#!/usr/bin/env bash

# OpenCode hook for tmux-agent-status.
# Tracks session/pane state via OpenCode lifecycle events.
#
# Setup: add to .opencode/opencode.jsonc:
#   "hooks": {
#     "UserPromptSubmit": "bash <path>/opencode-hook.sh UserPromptSubmit",
#     "PreToolUse": "bash <path>/opencode-hook.sh PreToolUse",
#     "Stop": "bash <path>/opencode-hook.sh Stop",
#     "SessionStart": "bash <path>/opencode-hook.sh SessionStart"
#   }

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
# shellcheck source=lib.sh
source "$SCRIPT_DIR/lib.sh"

AGENT_NAME="opencode"
init_status_dirs
drain_stdin

TMUX_SESSION=$(get_tmux_session) || exit 0
HOOK_TYPE="${1:-}"

case "$HOOK_TYPE" in
  SessionStart)
    # Session started — mark ready/done.
    if [ ! -f "$WAIT_DIR/${TMUX_SESSION}.wait" ] && [ ! -f "$PARKED_DIR/${TMUX_SESSION}.parked" ]; then
      set_status "$TMUX_SESSION" "done"
      mark_refresh
    fi
    ;;
  UserPromptSubmit)
    # User submitted a prompt — cancel wait/park, mark working.
    clear_interaction_overrides "$TMUX_SESSION"
    set_status "$TMUX_SESSION" "working"
    mark_refresh
    ;;
  PreToolUse|PostToolUse)
    # Agent is calling a tool — mark working but do NOT unpark.
    rm -f "$WAIT_DIR/${TMUX_SESSION}.wait"
    if [ ! -f "$PARKED_DIR/${TMUX_SESSION}.parked" ]; then
      set_status "$TMUX_SESSION" "working"
    fi
    mark_refresh
    ;;
  Stop)
    # OpenCode has finished responding.
    set_status "$TMUX_SESSION" "done"
    mark_refresh
    ;;
  Notification)
    # OpenCode is waiting for user input.
    set_status "$TMUX_SESSION" "done"
    mark_refresh
    ;;
esac

exit 0
