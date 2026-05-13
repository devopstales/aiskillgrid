#!/usr/bin/env bash

# Codex CLI hook for tmux-agent-status.
# Tracks session/pane state via Codex lifecycle hooks.
#
# Setup: enable hooks in ~/.codex/config.toml:
#   [features]
#   codex_hooks = true
#
# Then add to ~/.codex/hooks.json (or .codex/hooks.json per project):
#   {
#     "hooks": {
#       "SessionStart": [
#         { "matcher": "startup|resume", "hooks": [{ "type": "command", "command": "bash <path>/codex-hook.sh SessionStart" }] }
#       ],
#       "UserPromptSubmit": [
#         { "hooks": [{ "type": "command", "command": "bash <path>/codex-hook.sh UserPromptSubmit" }] }
#       ],
#       "PreToolUse": [
#         { "matcher": "Bash", "hooks": [{ "type": "command", "command": "bash <path>/codex-hook.sh PreToolUse" }] }
#       ],
#       "Stop": [
#         { "hooks": [{ "type": "command", "command": "bash <path>/codex-hook.sh Stop" }] }
#       ]
#     }
#   }

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
# shellcheck source=lib.sh
source "$SCRIPT_DIR/lib.sh"

AGENT_NAME="codex"
init_status_dirs
drain_stdin

TMUX_SESSION=$(get_tmux_session) || exit 0
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
  Stop)
    # Codex has finished responding.
    set_status "$TMUX_SESSION" "done"
    mark_refresh
    ;;
  Notification)
    # Codex is waiting for user input.
    set_status "$TMUX_SESSION" "done"
    mark_refresh
    ;;
esac

exit 0
