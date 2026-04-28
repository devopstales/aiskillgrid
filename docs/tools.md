# Tools and external dependencies

## Hub maintenance: IDE mirrors

Cursor is the **source of truth** for slash commands ([`.cursor/commands/`](../.cursor/commands/)) and agent personas ([`.cursor/agents/`](../.cursor/agents/)). After editing those files, run:

```bash
./scripts/sync-ide-assets.sh
```

That copies `skillgrid-*.md` and `opsx-*.md` to **`.kilo/commands/`** and **`.opencode/commands/`**, regenerates **`.github/prompts/`** (Copilot: `description` + body only), and syncs **`.cursor/agents/*.md`** to **`.kilo/agents/`**, **`.opencode/agents/`**, and **`.github/agents/`**.

Check for drift without writing (e.g. in CI):

```bash
./scripts/sync-ide-assets.sh --check
```

Requires **Python 3** (for GitHub prompt generation) and **`cmp`** (POSIX).

If you **rename or remove** a persona file under **`.cursor/agents/`**, delete the old filename from **`.kilo/agents/`**, **`.opencode/agents/`**, and **`.github/agents/`** (the script overwrites matching names but does not delete orphans), then run the sync again.

## Install tools (canonical)

**Base (system or Homebrew):** `brew`, `rsync`, `jq`, **Node (LTS)** + `npx`, **Python 3**.

### Python: `uv` only

Use **[uv](https://docs.astral.sh/uv/)** for everything Python: global CLIs and per-skill environments. Do not standardize on **pipx** in this repo (legacy docs may still mention it—prefer `uv`).

- **Global CLIs (replaces pipx):** `uv tool install <package>`, e.g. `uv tool install graphifyy`, `uv tool install --upgrade 'cocoindex-code[full]'` (**`ccc`**), `uv tool install tavily-cli`. Hub [`install.sh`](../install.sh) **`-t`** can install **cocoindex-code** (menu **6**) among optional tools.
- **Upgrade:** `uv tool upgrade <name>`
- **Per-skill venvs:** in a directory with `pyproject.toml` / `uv.lock`, run `uv sync` and `uv run …`

### JavaScript: one lockfile, `npx` to run

At the **repository root** of this hub, [package.json](../package.json) lists pinned **devDependencies** (OpenSpec CLI, `dmux`, and MCP packages aligned with [`.configs/mcp/node/`](../.configs/mcp/node/)). Install once:

```bash
npm ci
```

(Use `npm install` if you are bootstrapping and do not have a lockfile yet; prefer committing **package-lock.json**.)

Run tools via **`npx`**, which resolves to **`node_modules/.bin`** when dependencies are installed—no need for a separate global install for day-to-day hub work. **`npx` is the runner, not a replacement for `npm install`:** run `npm ci` in the hub first so versions match the lockfile.

`npm install -g` is reserved for exceptional cases; [`install.sh`](../install.sh) may fall back to global install when a hub-local install is not available.

## MCP: `npx -y` vs project `npx`

MCP server fragments under [`.configs/mcp/node/`](../.configs/mcp/node/) use **`npx`** (often with **`-y`**) so **any** machine can start a server without a prior `npm ci` in this repo. That is the **portable, zero-config** path (may fetch from the registry on first run).

For **reproducible, lockfile-pinned** runs (maintenance and CI in this repo), run **`npm ci`** in the hub, then use **`npx` without `-y`** so the same versions as [package-lock.json](../package-lock.json) are used from `node_modules`. Optional later improvement: set MCP `cwd` to the hub root so `npx` always prefers local `node_modules` (depends on your IDE’s MCP working directory).

## Context engineering

* [ ] [carl](https://github.com/ChristopherKahler/carl)
* [ ] [paul](https://github.com/ChristopherKahler/paul)
* [ ] [seed](https://github.com/ChristopherKahler/seed)

## Spec Driven Development

* [X] [OpenSpec](https://github.com/openspecs/openspec)
* [ ] [Antigravity Kit](https://github.com/vudovn/antigravity-kit)
* [-] [GSD](https://github.com/gsd-build/get-shit-done)
* [ ] [GSD-2](https://github.com/gsd-build/gsd-2)
* [ ] [superpowers](https://github.com/obra/superpowers)
* [ ] [Archon](https://github.com/coleam00/Archon)

## Web research and scraping (inventory)

Canonical tables and skills: **[web-scraping-and-research.md](web-scraping-and-research.md)**. Hub MCP fragments live under [`.configs/mcp/`](../.configs/mcp/) and are merged by [`install.sh`](../install.sh); **merged server keys** match the JSON top-level names below.

| Tool | Fragment | Key (example) |
|------|----------|----------------|
| [Context7](https://context7.com) | [.configs/mcp/http/context7.json](../.configs/mcp/http/context7.json) | `context7` |
| [Exa](https://exa.ai) | [.configs/mcp/http/exa.json](../.configs/mcp/http/exa.json) | `exa-http` |
| [DeepWiki](https://deepwiki.com) | [.configs/mcp/http/deepwiki.json](../.configs/mcp/http/deepwiki.json) | `deepwiki` |
| [Firecrawl](https://firecrawl.dev) | [.configs/mcp/node/fireclaw.json](../.configs/mcp/node/fireclaw.json) | `firecrawl` (`npx -y firecrawl-mcp`) |

**Brave Search** and **Tavily** are documented as **skills / CLIs** in [web-scraping-and-research.md](web-scraping-and-research.md); they are not default MCP fragments in this repo unless you add them.

Quick checklist:

* [X] Context7 — MCP fragment + [documentation-lookup](../.agents/skills/documentation-lookup/SKILL.md) / [context7](../.agents/skills/context7/SKILL.md)
* [X] Exa — MCP + [exa-search](../.agents/skills/exa-search/SKILL.md)
* [X] DeepWiki — MCP
* [X] Firecrawl — MCP + [deep-research](../.agents/skills/deep-research/SKILL.md)
* [ ] [Brave Search API](https://api.search.brave.com) — `brave-*` / `bx` skills (see web-scraping doc)
* [ ] Per-session install: [Tavily](https://tavily.com) CLI via `uv tool install tavily-cli` (optional)

## Memory and codebase index

Conceptual overview: **[memory.md](memory.md)** (what each layer is for).

* [X] [CocoIndex Code (`ccc`)](https://github.com/cocoindex-io/cocoindex-code) — MCP [cocoindex-code.json](../.configs/mcp/command/cocoindex-code.json); skill [ccc](../.agents/skills/ccc/SKILL.md); during **`/skillgrid-init`**: **`ccc init`** (if needed) then **`ccc index`**
* [X] [graphify](https://github.com/safishamsi/graphify) — PyPI [`graphifyy`](https://pypi.org/project/graphifyy/); MCP [graphify.json](../.configs/mcp/python/graphify.json); init **`graphify .`**, refresh **`graphify update .`**
* [X] [Engram](https://github.com/Gentleman-Programming/engram)
* [ ] [context-mode](https://github.com/mksglu/context-mode)

### Engram MCP (multi-agent)

Use the **same** Engram MCP entry everywhere you run agents (Cursor, OpenCode, Copilot MCP, etc.) so everyone reads/writes the same store (default: `~/.engram/`).

1. **Binary:** `brew install gentleman-programming/tap/engram` or run [`install.sh`](../install.sh) with **`-t`** and select **engram**.
2. **MCP fragment:** [`.configs/mcp/command/engram.json`](../.configs/mcp/command/engram.json) — `command`: `engram`, `args`: `["mcp", "--tools=agent"]`. Merge this into the IDE’s merged MCP config (this hub’s [`install.sh`](../install.sh) composes `.configs/mcp/**/*.json` for targets that use it).
3. **Agent discipline:** `.agents/skills/memory-protocol/SKILL.md` — when to `mem_save`, `mem_search`, `mem_session_summary`.
4. **Optional team sync:** [`engram sync`](https://github.com/Gentleman-Programming/engram) exports compressed chunks under `.engram/`; others run `engram sync --import`. See upstream [DOCS.md](https://github.com/Gentleman-Programming/engram/blob/main/DOCS.md) and [AGENT-SETUP.md](https://github.com/Gentleman-Programming/engram/blob/main/docs/AGENT-SETUP.md).

## Design Tools

* [ ] [taste-skill](https://github.com/Leonxlnx/taste-skill)
* [ ] [npxskillui](https://github.com/amaancoderx/npxskillui)
* [ ] [impeccable](https://github.com/pbakaus/impeccable)

## Browser tools

* [ ] [agent-browser](https://github.com/vercel-labs/agent-browser)
* [ ] [PlayWright-mcp](https://github.com/microsoft/playwright-mcp)
  * [ ] [Skill-1](https://github.com/testdino-hq/playwright-skill)
  * [ ] [Skill-2](https://github.com/lackeyjb/playwright-skill)


## Agent multiplexers

* [ ] [dmux](https://github.com/standardagents/dmux)
* [ ] [mux](https://github.com/coder/mux)
* [ ] [superset](https://github.com/superset-sh/superset)
* [ ] [toad](https://github.com/batrachianai/toad)
