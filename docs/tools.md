# Tools and external dependencies

This document is split into **(1) the install toolchain** used to run [`install.sh`](../install.sh) and install language runtimes, and **(2) workflow CLIs, optional stacks, and IDEs** this hub integrates with. It is a catalog, not a mandatory install list.

Related docs: [`skills.md`](skills.md), [`commands.md`](commands.md), [`wokflow.md`](wokflow.md).

---

## 1. Install toolchain (foundation)

These are **package managers, runtimes, and OS utilities** you install once on a machine (or in CI) so scripts and agents can run. [`install.sh`](../install.sh) assumes several of them when merging MCP configs, syncing files, or delegating to `brew` / `npm` / `pipx`.

| Tool | What it is | Typical use here |
|------|------------|-------------------|
| **bash** | Shell (4+ per `install.sh` header). | Run `install.sh`. |
| **rsync** | File sync. | Copy hub folders (`.cursor/`, skills, etc.) into the target project. |
| **jq** | JSON processor. | Merge [`.configs/mcp`](../.configs/mcp/) into per-IDE MCP files; patch OpenCode config. Omit only with `--no-mcp`. |
| **brew** (Homebrew) | macOS/Linux package manager. | Preferred path to install `jq`, `node`, `python@3.x`, `openspec`, `opencode`, `kilo`, `semgrep`, `trivy`, `engram` tap, etc., when `brew` is on `PATH`. |
| **node** | JavaScript runtime. | Run project scripts; required for `npx` / global npm CLIs. |
| **npm** | Node package manager. | Global installs used by the script: e.g. `@fission-ai/openspec`, `@kilocode/cli`, `dmux`. |
| **npx** | Run Node packages without a global install. | Invoked in skills for one-off tools (`tsc`, etc.). |
| **python3** | Python 3 interpreter. | Present on dev machines; some automation and pip-based installs. |
| **pip** / **pip3** | Python package installer. | `install.sh` may use `pip3 install --user` for declared pip deps; **`ensure_pipx`** can use `pip3 install --user pipx`. |
| **pipx** | Install Python CLI apps in isolated envs. | **`cocoindex-code[full]`**, **graphifyy** when selected via `install.sh -t`. |
| **uv** | Fast Python package/project tool. | Listed in `install.sh` optional dependencies for Python-centric workflows. |

**Minimal set to run `install.sh` with MCP merge:** `bash`, `rsync`, `jq`, plus **node** / **npx** if you rely on npm-based installs. **brew** is optional but unlocks one-command installs for many CLIs below.

---

## 2. Workflow tools and IDEs

### 2.1 Spec, index, and knowledge tools

| Tool | Role | How it usually arrives |
|------|------|------------------------|
| **OpenSpec** (`openspec` CLI) | Change workflow: `openspec new`, `status`, `instructions`, `archive`, etc. | `brew install openspec` or `npm install -g @fission-ai/openspec@latest`. Hub may copy `openspec/` and run `openspec init` when optional tool **openspec** is selected (`install.sh -t`). |
| **CocoIndex Code** | Semantic codebase index and search; **`ccc`** CLI in the [`ccc`](../.agents/skills/ccc/SKILL.md) skill (`ccc init`, `ccc index`, `ccc search`). | `pipx install 'cocoindex-code[full]'` (optional `-t`); config dir **`.cocoindex_code/`** may be copied from hub. |
| **Graphify** | Knowledge graph / AST-style tooling; CLI often **`graphifyy`** or **`graphify`**. | `pipx install graphifyy` (optional `-t`); hub may sync **`.graphify/`**. |

### 2.2 Security and agent-helper CLIs

| Tool | Role | How it usually arrives |
|------|------|------------------------|
| **Semgrep** | Static analysis (`semgrep-security` skill). | `brew install semgrep` when Kilocode/OpenCode path pulls IDE deps. |
| **Trivy** | Vulns, misconfigs, secrets, SBOM (`trivy-security` skill). | `brew install trivy`; **`trivy plugin install mcp`** for MCP integration. |
| **dmux** | Tmux-oriented multiplexer helper for agents. | `npm install -g dmux` (optional `-t`). |
| **Engram** | MCP / memory CLI used with SDD **engram** persistence mode. | `brew install gentleman-programming/tap/engram` (optional `-t`). |

### 2.3 IDE CLIs (terminal companions)

| CLI | IDE / product | Notes |
|-----|----------------|-------|
| **opencode** | OpenCode | `brew install opencode`; config under **`.opencode/`** (`opencode.jsonc`, commands, skills). |
| **kilo** | Kilo Code | `brew install Kilo-Org/homebrew-tap/kilo`. |
| **@kilocode/cli** (`npm install -g`) | Kilocode | `install.sh` **`setup_kilo`** installs this if `kilo` is missing (npm global). |

### 2.4 IDEs this hub targets (config layout)

The installer copies or generates config for these environments (flags like `-c`, `-C`, `-k`, `-o`, `-a`, `-A`):

| IDE | Folders / files in target project | MCP output (when merged) |
|-----|-----------------------------------|---------------------------|
| **Cursor** | `.cursor/` (commands, rules, agents, plugins; skills → `.cursor/.agents/skills/`) | `.cursor/mcp.json` |
| **GitHub Copilot** (VS Code) | `.github/`, `.vscode/` | `.vscode/mcp.json` (`servers` shape) |
| **Kilo Code** | `.kilo/` (skills → `.kilo/skills/`) | `.kilo/mcp.json` (or `kilo.jsonc` per product) |
| **OpenCode** | `.opencode/` (skills → `.opencode/skills/`) | MCP merged into **`opencode.json`** / `opencode.jsonc` |
| **Google Antigravity** | Project **`.agents/`**; global **`~/.gemini/antigravity/mcp_config.json`** | Antigravity MCP file |
| **Kiro** | Not wired in `install.sh` today; product uses **`.kiro/skills/`** (and **Agents.md**) per [Kiro skills docs](https://kiro.dev/docs/skills/). | Configure per Kiro docs. |

Skills sync destinations from **`install.sh`**: Cursor → `.cursor/.agents/skills/`; Kilo → `.kilo/skills/`; OpenCode → `.opencode/skills/`; Antigravity → `.agents/skills/`. Copilot does not get a separate skills rsync path in the script (skills live under hub `.agents/skills/` copied via `.github` / project conventions).

### 2.5 MCP config merge

- **Inputs:** [`.configs/mcp.json`](../.configs/mcp.json) and [`.configs/mcp/**/*.json`](../.configs/mcp/).
- **Requires:** `jq` (section 1).
- **Servers** (examples referenced in skills): Exa, Firecrawl, Context7, browser/DevTools, Engram, etc. Exact keys depend on your merged JSON.

### 2.6 Other tools referenced by skills

| Area | Tools | Example skills |
|------|--------|----------------|
| Testing | `npm test`, Playwright, **agent-browser** (`npm i -g agent-browser`) | `e2e-testing`, `e2e-runner`, `test-driven-development` |
| DB | `psql`, Prisma CLI | `database-reviewer`, `database-migrations`, `shipping-and-launch` |
| Research MCPs | Exa, Firecrawl | `deep-research`, `exa-search` |
| Docs MCP | Context7 | `documentation-lookup` |
| Config audit | AgentShield (via `security-scan`) | `security-scan` |
| Images | Gemini / Nano Banana CLI | `nano-banana` |
| Optional script | `python scripts/security_scan.py` | `vulnerability-scanner` |

---

## Choosing what to install

1. **Hub installer only:** section 1 minimum + `jq` + Node/`npx` if you use npm globals.
2. **OpenSpec:** `openspec` CLI + optional hub `openspec/` tree.
3. **Semantic search / graph:** **cocoindex-code** + **`ccc`**, and/or **graphifyy**.
4. **Security automation:** **Trivy**, **Semgrep**, plus **Trivy MCP** plugin if needed.
5. **IDE:** install the matching CLI (section 2.3) and run `install.sh` for that IDE.

---

## Keeping this document accurate

When [`install.sh`](../install.sh) dependency arrays change or skills add a new hard dependency, update sections 1–2 and the tables in [`TODO.md`](TODO.md) if you track matrices there.
