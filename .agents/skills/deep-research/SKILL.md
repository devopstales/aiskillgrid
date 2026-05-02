---
name: deep-research
description: Multi-source deep research using available search and retrieval tools (Firecrawl, Exa, Brave, Tavily, Context7, DeepWiki, or built-in web search). Synthesizes findings and delivers cited reports with source attribution. Use when the user wants thorough research on any topic with evidence and citations.
origin: ECC
---

# Deep Research

Produce thorough, cited research reports from multiple web sources. Prefer **two or more independent tools** when possible (e.g. Exa + Firecrawl, or Brave web search + Tavily) so coverage and ranking bias do not come from a single provider.

## When Used From Skillgrid

For Skillgrid work, use `skillgrid-parallel-research` to split research questions and manage handoff updates. Store long reports under `.skillgrid/tasks/research/<change-id>/`, cite web sources, and return a short recommendation with artifact paths instead of raw output in chat.

## When to Activate

- User asks to research any topic in depth
- Competitive analysis, technology evaluation, or market sizing
- Due diligence on companies, investors, or technologies
- Any question requiring synthesis from multiple sources
- User says "research", "deep dive", "investigate", or "what's the current state of"

## Search and research tooling (use what is configured)

You need **at least one** general web search path. Combine layers when the question needs it: web results + full-page read + official docs or repo knowledge.

### Primary web search and page extraction (MCP)

| Tool | Role | Typical entry points |
|------|------|----------------------|
| **Firecrawl** | Search, scrape single URLs, crawl sites | `firecrawl_search`, `firecrawl_scrape`, `firecrawl_crawl` |
| **Exa** | Neural/semantic search; URL crawl for reading | `web_search_exa`, `web_search_advanced_exa`, `crawling_exa` |

Hub MCP config fragments: `.configs/mcp/node/fireclaw.json`, `.configs/mcp/http/exa.json` (from repo root). For Exa usage details, read the sibling skill `exa-search` (`SKILL.md` next to this folder).

**Firecrawl + Exa together** usually give the strongest web coverage. Configure in `~/.claude.json`, IDE MCP, or `~/.codex/config.toml`.

### Brave Search API (HTTP; skills under `brave-*`)

Use when the user has **`BRAVE_SEARCH_API_KEY`**. All-in-one and specialized endpoints:

| Skill name | Use for |
|--------|---------|
| `brave-web-search` | Primary web search (snippets, URLs, freshness, SafeSearch) |
| `brave-bx` | Broad “omni” / agentic search (pre-extracted content, many modes) |
| `brave-news-search` | News and time-bounded current events |
| `brave-videos-search` / `brave-images-search` | Video or image surface research |
| `brave-llm-context` | Pre-extracted page text for RAG/grounding (token budget) |
| `brave-answers` | AI-grounded answers with citations (slower, thorough) |
| `brave-local-pois` / `brave-local-descriptions` | Places and POIs (needs location flow from web search) |
| `brave-suggest` / `brave-spellcheck` | Query refine before main search |

### Tavily (REST)

**Tavily** — sibling skill `tavily`: `POST /search`, `/extract`, mapping/crawl per that skill. Requires **`TAVILY_API_KEY`**. Strong when you want **structured, LLM-oriented** results (optional `include_answer`, raw content, topic/time filters). Use alongside Exa or Brave to diversify sources.

### Official library and framework docs (not open web)

**documentation-lookup** (sibling skill) + **Context7** MCP — `resolve-library-id`, `query-docs`. Use when research is **“how does X API work in 2026?”** for a named stack (React, Prisma, cloud SDKs). Complement web search; do not use Context7 as a replacement for general news or market research.

### GitHub repository knowledge (not generic web)

**DeepWiki** MCP — `read_wiki_structure`, `read_wiki_contents`, `ask_question` for **documentation-style Q&A about repos**. Use when the topic is a specific open-source project or “how is this implemented upstream?”

### Built-in (Cursor and similar IDEs)

When no research MCP is available:

- **WebSearch** — quick discovery and recent facts (always verify in the final report with primary URLs you cite)
- **WebFetch** — fetch a **known** `https://` URL for a short page or announcement

Treat built-in search as a **fallback** for breadth; prefer MCP + full scrape/read for deep dives.

### Other skills that pair with this workflow

- **`search-first`** — before large custom implementations, look for prior art (may route through researcher-style tools).
- Optional **CLIs** (Firecrawl, Tavily via `uv tool install`, etc.) when MCP is not used — only if your runbook or operator instructions say to use a CLI instead of MCP.

## Workflow

### Step 1: Understand the Goal

Ask 1-2 quick clarifying questions:
- "What's your goal — learning, making a decision, or writing something?"
- "Any specific angle or depth you want?"

If the user says "just research it" — skip ahead with reasonable defaults.

### Step 2: Plan the Research

Break the topic into 3-5 research sub-questions. Example:
- Topic: "Impact of AI on healthcare"
  - What are the main AI applications in healthcare today?
  - What clinical outcomes have been measured?
  - What are the regulatory challenges?
  - What companies are leading this space?
  - What's the market size and growth trajectory?

### Step 3: Execute Multi-Source Search

For EACH sub-question, run queries through **whatever is configured** in the “Search and research tooling” section above. Mix providers when you can (e.g. Exa + Brave, or Firecrawl + Tavily).

**With firecrawl:**
```
firecrawl_search(query: "<sub-question keywords>", limit: 8)
```

**With exa:**
```
web_search_exa(query: "<sub-question keywords>", numResults: 8)
web_search_advanced_exa(query: "<keywords>", numResults: 5, startPublishedDate: "2026-01-01")
```

**With Brave** — follow the `brave-web-search` skill (or `brave-bx` / `brave-news-search` when the sub-question is news- or media-specific). Use `brave-llm-context` or `brave-answers` when you need pre-extracted text or a cited narrative.

**With Tavily** — `POST /search` (and `POST /extract` for chosen URLs) per the `tavily` skill.

**Library-accurate research** — call **Context7** (`resolve-library-id` then `query-docs`) per the `documentation-lookup` skill for sub-questions that depend on a specific framework or SDK.

**Open-source on GitHub** — use **DeepWiki** for repo-specific “how does this project do X?” sub-questions.

**Fallback** — **WebSearch** / **WebFetch** when no research MCP is available; still list URLs in the Sources section.

**Search strategy:**
- Use 2-3 different keyword variations per sub-question
- Mix general and news-focused queries
- Aim for 15-30 unique sources total
- Prioritize: academic, official, reputable news > blogs > forums

### Step 4: Deep-Read Key Sources

For the most promising URLs, fetch full content:

**With firecrawl:**
```
firecrawl_scrape(url: "<url>")
```

**With exa:**
```
crawling_exa(url: "<url>", tokensNum: 5000)
```

**Alternatives:** **Tavily** `POST /extract` (see `tavily` skill) or **Brave** `brave-llm-context` for pre-extracted body text when Firecrawl/Exa are not the path you use for that URL.

Read 3-5 key sources in full for depth. Do not rely only on search snippets.

### Step 5: Synthesize and Write Report

Structure the report:

```markdown
# [Topic]: Research Report
*Generated: [date] | Sources: [N] | Confidence: [High/Medium/Low]*

## Executive Summary
[3-5 sentence overview of key findings]

## 1. [First Major Theme]
[Findings with inline citations: source title and full URL per claim]
- Key point — source title, full URL
- Supporting data — source title, full URL

## 2. [Second Major Theme]
...

## 3. [Third Major Theme]
...

## Key Takeaways
- [Actionable insight 1]
- [Actionable insight 2]
- [Actionable insight 3]

## Sources
1. Title — full URL — one-line summary
2. ...

## Methodology
Tools used: [e.g. Firecrawl, Exa, Brave web search, Tavily, Context7, DeepWiki, WebSearch]. Searched [N] queries across web and/or news. Analyzed [M] sources.
Sub-questions investigated: [list]
```

### Step 6: Deliver

- **Short topics**: Post the full report in chat
- **Long reports**: Post the executive summary + key takeaways, save full report to a file

### Step 7: Optional — persist for later sessions

If **Engram** (or another memory MCP) is available and the user wants recall across sessions, save a short anchored note with **`mem_save`** (title/topic_key scoped to the research topic, e.g. `research/<slug>`). Include the same source URLs you cited in the report. If the project uses the hub memory stack, follow `../references/indexing-and-memory.md` and `../../../docs/memory.md` from this file for ordering with GitNexus and structural search.

## Parallel Research with Subagents

For broad topics, use Claude Code's Task tool to parallelize:

```
Launch 3 research agents in parallel:
1. Agent 1: Research sub-questions 1-2
2. Agent 2: Research sub-questions 3-4
3. Agent 3: Research sub-question 5 + cross-cutting themes
```

Each agent searches, reads sources, and returns findings. The main session synthesizes into the final report.

## Quality Rules

1. **Every claim needs a source.** No unsourced assertions.
2. **Cross-reference.** If only one source says it, flag it as unverified.
3. **Recency matters.** Prefer sources from the last 12 months.
4. **Acknowledge gaps.** If you couldn't find good info on a sub-question, say so.
5. **No hallucination.** If you don't know, say "insufficient data found."
6. **Separate fact from inference.** Label estimates, projections, and opinions clearly.

## Examples

```
"Research the current state of nuclear fusion energy"
"Deep dive into Rust vs Go for backend services in 2026"
"Research the best strategies for bootstrapping a SaaS business"
"What's happening with the US housing market right now?"
"Investigate the competitive landscape for AI code editors"
```
