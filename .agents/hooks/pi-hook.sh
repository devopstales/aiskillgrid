#!/usr/bin/env bash

# Pi Coding Agent hook for tmux-agent-status.
# Tracks session/pane state via Pi lifecycle hooks.
#
# Pi supports command hooks via the hsingjui/pi-hooks compatibility layer
# or via TypeScript extensions. For shell-based integration, use the
# Claude Code-compatible hook format that pi-hooks supports.
#
# Setup: add to ~/.pi/agent/settings.json or .pi/settings.json:
#   {
#     "hooks": {
#       "SessionStart": [
#         { "matcher": "startup", "hooks": [{ "type": "command", "command": "bash <path>/pi-hook.sh SessionStart" }] },
#         { "matcher": "resume", "hooks": [{ "type": "command", "command": "bash <path>/pi-hook.sh SessionStart" }] }
#       ],
#       "UserPromptSubmit": [
#         { "hooks": [{ "type": "command", "command": "bash <path>/pi-hook.sh UserPromptSubmit" }] }
#       ],
#       "PreToolUse": [
#         { "hooks": [{ "type": "command", "command": "bash <path>/pi-hook.sh PreToolUse" }] }
#       ],
#       "Stop": [
#         { "hooks": [{ "type": "command", "command": "bash <path>/pi-hook.sh Stop" }] }
#       ],
#       "SessionEnd": [
#         { "hooks": [{ "type": "command", "command": "bash <path>/pi-hook.sh Stop" }] }
#       ]
#     }
#   }

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
# shellcheck source=lib.sh
source "$SCRIPT_DIR/lib.sh"

AGENT_NAME="pi"
init_status_dirs
drain_stdin

TMUX_SESSION=$(get_tmux_session) || {
  TMUX_SESSION="pi"
}
HOOK_TYPE="${1:-}"

case "$HOOK_TYPE" in
  SessionStart)
    # Session started or resumed — mark ready/done.
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
  Stop|SessionEnd)
    # Pi has finished responding or session ended.
    set_status "$TMUX_SESSION" "done"
    mark_refresh
    ;;
  Notification)
    # Pi is waiting for user input.
    set_status "$TMUX_SESSION" "done"
    mark_refresh
    ;;
esac

exit 0
