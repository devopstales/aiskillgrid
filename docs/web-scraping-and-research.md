# Web research, scraping, and external documentation

This hub does **not** ship a single “web scraper” product. Instead it combines **MCP servers** (search, scrape, crawl, docs APIs) with **Agent Skills** and the **`skillgrid-researcher`** persona so agents can gather **cited** evidence from the web and from official library docs.

For **install and MCP merging**, see [`tools.md`](tools.md) and [`install.sh`](../install.sh). Fragment paths below are under [`.configs/mcp/`](../.configs/mcp/).

---

## MCP servers in this repo (research-relevant)

These JSON fragments are merged into your IDE `mcp.json` when you use the hub installer. **Server keys** are the top-level keys inside each file (names may vary slightly in your merged config).

| Fragment | Typical key | Role |
|----------|-------------|------|
| [`http/exa.json`](../.configs/mcp/http/exa.json) | `exa-http` | [Exa](https://exa.ai) MCP over HTTP — neural / semantic web search and related retrieval. |
| [`node/fireclaw.json`](../.configs/mcp/node/fireclaw.json) | `firecrawl` | [Firecrawl](https://firecrawl.dev) via `npx` / `firecrawl-mcp` — search, scrape, crawl (requires API key where the product requires it). |
| [`http/deepwiki.json`](../.configs/mcp/http/deepwiki.json) | `deepwiki` | [DeepWiki](https://deepwiki.com) MCP — documentation-style Q&A over GitHub repos. |
| [`http/context7.json`](../.configs/mcp/http/context7.json) | `context7` | [Context7](https://context7.com) MCP — **library and framework docs** (`resolve-library-id`, `query-docs`), not general web scraping. |
| [`docker/fetch.json`](../.configs/mcp/docker/fetch.json) | `fetch-docker` | [MCP Fetch](https://github.com/modelcontextprotocol/servers) in Docker — simple **URL → content** fetch for allowed URLs. |

---

## Agent skills (workflows)

| Skill | Path | Use when |
|--------|------|----------|
| `deep-research` | [.agents/skills/deep-research/SKILL.md](../.agents/skills/deep-research/SKILL.md) | Multi-step research with **citations**; expects **Firecrawl** and/or **Exa** MCP tools. |
| `exa-search` | [.agents/skills/exa-search/SKILL.md](../.agents/skills/exa-search/SKILL.md) | Exa-specific queries — web, code examples, companies, people. |
| `documentation-lookup` | [.agents/skills/documentation-lookup/SKILL.md](../.agents/skills/documentation-lookup/SKILL.md) | **Framework / SDK** questions — Context7 **MCP** (`resolve-library-id`, `query-docs`). |
| `context7` | [.agents/skills/context7/SKILL.md](../.agents/skills/context7/SKILL.md) | Same **Context7** catalog via **HTTP API** and `curl` when MCP is not available. |
| `search-first` | [.agents/skills/search-first/SKILL.md](../.agents/skills/search-first/SKILL.md) | **Before coding** — look for existing libraries, patterns, and prior art (may delegate to researcher-style tooling). |

### Brave Search API (`brave-*`)

These skills document the **[Brave Search API](https://api.search.brave.com)** (API key + `curl`; some endpoints need specific subscription plans). They are **not** merged as MCP fragments in [`.configs/mcp/`](../.configs/mcp/) by default—use them from the shell or embed patterns in agent workflows. Prefer **`brave-bx`** (`bx` CLI) for an all-in-one agent-oriented search when you have the binary installed.

| Directory | `name` (frontmatter) | Role |
|-----------|----------------------|------|
| `brave-bx` | `bx` | Brave Search **CLI** — default agentic search / grounding. |
| `brave-web-search` | `web-search` | Primary web search (snippets, URLs, Goggles, pagination). |
| `brave-llm-context` | `llm-context` | Pre-extracted web content for RAG / LLM grounding. |
| `brave-answers` | `answers` | OpenAI-compatible `/chat/completions` with Brave grounding. |
| `brave-news-search` | `news-search` | News search. |
| `brave-images-search` | `images-search` | Image search. |
| `brave-videos-search` | `videos-search` | Video search. |
| `brave-suggest` | `suggest` | Query autocomplete / suggestions. |
| `brave-spellcheck` | `spellcheck` | Spell correction / “did you mean”. |
| `brave-local-pois` | `local-pois` | Local POI details (needs location IDs from web search). |
| `brave-local-descriptions` | `local-descriptions` | POI text descriptions from POI IDs. |

Full skill index with summaries: [`skills.md`](skills.md) (sections **Brave Search API** and **Ship, git…** for `documentation-lookup` / `context7`).

---

## IDE persona

**`skillgrid-researcher`** ([`.cursor/agents/skillgrid-researcher.md`](../.cursor/agents/skillgrid-researcher.md)) is the hub’s **research specialist**: it documents the same MCP table as above, ties together `deep-research`, `exa-search`, `documentation-lookup`, and `search-first`, and defines **citation** and output format. For **Brave**, use the `brave-*` skills or `bx` when your session has API keys / CLI; for **Context7** without MCP, use the `context7` skill. Synced copies live under `.kilo/agents/`, `.opencode/agents/`, `.github/agents/` when you run [`scripts/sync-ide-assets.sh`](../scripts/sync-ide-assets.sh).

---

## Optional CLIs (outside MCP)

[`tools.md`](tools.md) mentions **`uv tool install tavily-cli`** as an example of a Python-installed CLI. That is **not** wired as a default MCP fragment in `.configs/mcp/` here; add a fragment or use it from the terminal if you standardize on [Tavily](https://tavily.com) for your team.

---

## `tools.md` checklist vs this repo

The **Web scrapper** bullet list in [`tools.md`](tools.md) (Context7, Exa, DeepWiki, Brave Search, firecrawl-cli) mixes **MCP** tracking with **CLI/API** ideas. **MCP fragments** that ship in this repo are the JSON files under [`.configs/mcp/`](../.configs/mcp/) in the first table. **Brave** and **Context7** also have **Agent Skills** (`brave-*`, `context7`, `documentation-lookup`) documented above—those work independently of whether a Brave MCP fragment exists in your merge.

---

## Related

- [`memory.md`](memory.md) — persisting research summaries with **Engram** (`mem_save`) when you want cross-session recall.
- [`agents.md`](agents.md) — where persona files live per IDE.
