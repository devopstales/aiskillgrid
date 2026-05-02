## GitNexus

Before answering architecture or codebase questions, prefer GitNexus MCP/resources when available for repository context, clusters, processes, impact analysis, and symbol context.
If `.gitnexus/` exists, treat it as local generated graph state and use it before broad raw file reads.
Run `npx -y gitnexus@1.3.11 analyze` from the repository root to build or update the knowledge graph.
