#!/usr/bin/env bash

# Claude Code hook for tmux-agent-status.
# Tracks session/pane state (working/done/wait) via Claude Code lifecycle hooks.
#
# Setup: add to ~/.claude/settings.json:
#   "hooks": {
#     "UserPromptSubmit": [{"hooks": [{"type": "command", "command": "bash <path>/claude-hook.sh UserPromptSubmit"}]}],
#     "PreToolUse":       [{"hooks": [{"type": "command", "command": "bash <path>/claude-hook.sh PreToolUse"}]}],
#     "Stop":             [{"hooks": [{"type": "command", "command": "bash <path>/claude-hook.sh Stop"}]}],
#     "Notification":     [{"hooks": [{"type": "command", "command": "bash <path>/claude-hook.sh Notification"}]}]
#   }

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
# shellcheck source=lib.sh
source "$SCRIPT_DIR/lib.sh"

AGENT_NAME="claude"
init_status_dirs
drain_stdin

TMUX_SESSION=$(get_tmux_session) || exit 0
HOOK_TYPE="${1:-}"

case "$HOOK_TYPE" in
  UserPromptSubmit)
    # User submitted a prompt — cancel wait/park, mark working.
    clear_interaction_overrides "$TMUX_SESSION"
    set_status "$TMUX_SESSION" "working"
    mark_refresh
    ;;
  PreToolUse)
    # Agent is calling a tool — mark working but do NOT unpark.
    rm -f "$WAIT_DIR/${TMUX_SESSION}.wait"
    if [ ! -f "$PARKED_DIR/${TMUX_SESSION}.parked" ]; then
      set_status "$TMUX_SESSION" "working"
    fi
    mark_refresh
    ;;
  Stop)
    # Claude has finished responding.
    set_status "$TMUX_SESSION" "done"
    mark_refresh
    ;;
  Notification)
    # Claude is waiting for user input.
    set_status "$TMUX_SESSION" "done"
    mark_refresh
    ;;
esac

exit 0
