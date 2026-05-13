#!/usr/bin/env bash

# GitHub Copilot CLI hook for tmux-agent-status.
# Copilot CLI does not have a native hook system; this is designed for
# manual invocation or integration via wrapper scripts/aliases.
#
# Usage:
#   bash <path>/copilot-hook.sh UserPromptSubmit
#   bash <path>/copilot-hook.sh Stop
#
# Can be wired via shell aliases or wrapper scripts around copilot commands.

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
# shellcheck source=lib.sh
source "$SCRIPT_DIR/lib.sh"

AGENT_NAME="copilot"
init_status_dirs
drain_stdin

TMUX_SESSION=$(get_tmux_session) || {
  TMUX_SESSION="copilot"
}
HOOK_TYPE="${1:-}"

case "$HOOK_TYPE" in
  UserPromptSubmit|CommandRun)
    clear_interaction_overrides "$TMUX_SESSION"
    set_status "$TMUX_SESSION" "working"
    mark_refresh
    ;;
  PreToolUse|ToolUse|BashRun)
    rm -f "$WAIT_DIR/${TMUX_SESSION}.wait"
    if [ ! -f "$PARKED_DIR/${TMUX_SESSION}.parked" ]; then
      set_status "$TMUX_SESSION" "working"
    fi
    mark_refresh
    ;;
  Stop|SessionEnd|Complete)
    set_status "$TMUX_SESSION" "done"
    mark_refresh
    ;;
  SessionStart|Init)
    if [ ! -f "$WAIT_DIR/${TMUX_SESSION}.wait" ] && [ ! -f "$PARKED_DIR/${TMUX_SESSION}.parked" ]; then
      set_status "$TMUX_SESSION" "done"
      mark_refresh
    fi
    ;;
  Notification|Waiting)
    set_status "$TMUX_SESSION" "done"
    mark_refresh
    ;;
  *)
    set_status "$TMUX_SESSION" "working"
    mark_refresh
    ;;
esac

exit 0
