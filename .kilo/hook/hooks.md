# Kilo Hook Policy

## Trigger: Index Refresh

- post-shell execution for:
  - `git merge`
  - `gh pr merge`
  - `git pull`
  - `ccc init`

## Command: Index Refresh

Run:

`./.agents/hooks/refresh-indexes.sh`

The script is fail-open and logs refresh results for `ccc index` and `npx gitnexus analyze`.

## Trigger: Agent Status (tmux-agent-status)

- post-shell execution (all commands)
- pre-tool-use
- session stop

## Command: Agent Status

Run:

`./.agents/hooks/kilo-hook.sh PostShell`
`./.agents/hooks/kilo-hook.sh PreToolUse`
`./.agents/hooks/kilo-hook.sh Stop`
