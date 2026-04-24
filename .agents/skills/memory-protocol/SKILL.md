---
name: engram-memory-protocol
description: >
  Persistent memory discipline for Engram contributors.
  Trigger: Decisions, bugfixes, discoveries, preferences, or session closure.
license: Apache-2.0
metadata:
  author: gentleman-programming
  version: "1.0"
---

## Engram in this repository

This skill is vendored from [Engram `skills/memory-protocol`](https://github.com/Gentleman-Programming/engram/tree/main/skills/memory-protocol). It is the **canonical protocol** for **Engram MCP** tools (`mem_save`, `mem_search`, `mem_context`, `mem_session_summary`, …).

**Wire-up in this hub:** merge [`.configs/mcp/command/engram.json`](../../.configs/mcp/command/engram.json) into your IDE MCP config; see [`docs/tools.md`](../../docs/tools.md) and [Engram agent setup](https://github.com/Gentleman-Programming/engram/blob/main/docs/AGENT-SETUP.md). When Engram is enabled, follow the rules below as written.

## When to Use

Use this skill when:
- Making architecture or implementation decisions
- Fixing bugs with non-obvious root causes
- Discovering patterns, gotchas, or user preferences
- Closing a session or after compaction

---

## Save Rules

Call `mem_save` immediately after:
- decision
- bugfix
- pattern/discovery
- config/preference changes

Use structured content:
- What
- Why
- Where
- Learned

Use stable `topic_key` for evolving topics.

---

## Search Rules

- On recall requests: `mem_context` first, then `mem_search`.
- Before similar work: run proactive `mem_search`.
- On first message: if user references the project, a feature, or a problem, call `mem_search` with their keywords before responding.

---

## Session Close Rules

Before saying done/listo:
1. Call `mem_session_summary`.
2. Include goal, discoveries, accomplished, next steps, relevant files.

After compaction:
1. Save summary first.
2. Recover context.
3. Continue work.
