# OpenCode Hook Policy

Use the shared refresh hook to enforce index refresh after merge/bootstrap shell commands.

## Trigger

- post-shell execution for:
  - `git merge`
  - `gh pr merge`
  - `git pull`
  - `ccc init`

## Command

Run:

`./.agents/hooks/refresh-indexes.sh`

The script is fail-open and logs refresh results for `ccc index` and `npx gitnexus analyze`.
