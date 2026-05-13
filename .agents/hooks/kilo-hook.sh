#!/usr/bin/env bash

# Kilo hook for tmux-agent-status.
# Tracks session/pane state via Kilo lifecycle events.
#
# Kilo hooks are triggered via the .kilo/hook/ mechanism.
# This hook responds to shell command events to infer agent state.
#
# Setup: add to .kilo/hook/hooks.md or configure in kilo.json:
#   "hooks": {
#     "PostShell": [{"command": "bash <path>/kilo-hook.sh PostShell"}],
#     "PreToolUse": [{"command": "bash <path>/kilo-hook.sh PreToolUse"}],
#     "Stop": [{"command": "bash <path>/kilo-hook.sh Stop"}]
#   }

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
# shellcheck source=lib.sh
source "$SCRIPT_DIR/lib.sh"

AGENT_NAME="kilo"
init_status_dirs
drain_stdin

TMUX_SESSION=$(get_tmux_session) || exit 0
HOOK_TYPE="${1:-}"

case "$HOOK_TYPE" in
  PostShell|PreToolUse|ToolUse)
    # Kilo is executing a command or using a tool — mark working.
    rm -f "$WAIT_DIR/${TMUX_SESSION}.wait"
    if [ ! -f "$PARKED_DIR/${TMUX_SESSION}.parked" ]; then
      set_status "$TMUX_SESSION" "working"
    fi
    mark_refresh
    ;;
  Stop|SessionEnd)
    # Kilo has finished responding.
    set_status "$TMUX_SESSION" "done"
    mark_refresh
    ;;
  UserPrompt)
    # User submitted input — cancel wait/park, mark working.
    clear_interaction_overrides "$TMUX_SESSION"
    set_status "$TMUX_SESSION" "working"
    mark_refresh
    ;;
  Notification)
    # Kilo is waiting for user input.
    set_status "$TMUX_SESSION" "done"
    mark_refresh
    ;;
esac

exit 0
