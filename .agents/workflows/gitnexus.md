# Workflow: GitNexus
**Command:** gitnexus
**Description:** Index a repository with GitNexus and expose graph-aware code intelligence to agents

## Steps
Run `npx -y gitnexus@1.3.11 analyze` from the repository root to create or refresh the local `.gitnexus/` index.

Use `npx -y gitnexus@1.3.11 analyze --skills` when repo-specific generated agent skills are desired.

Use `npx -y gitnexus@1.3.11 serve` for local web UI bridge mode, and configure MCP with `npx -y gitnexus@1.3.11 mcp` so agents can use graph-aware tools.
