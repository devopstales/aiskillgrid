
skill-pi (Pi distribution)
    * [03-skill-pi.md](03-skill-pi.md) — full plan for skill-pi: fat Pi distribution with SDD subagents, Mnemonic, dashboard, permissions, local LLM
    * command: `skill-pi`, home: `~/skill-pi/`, install: `npm i -g skill-pi` or `skillgrid install --agents skill-pi`
    * extensions: splash, mnemonic, sdd, permissions, dashboard, local-llm
    * subagents: sdd-explore, sdd-propose, sdd-spec, sdd-apply, sdd-verify, sdd-archive
    * bundled plugins: pi-mcp-adapter, pi-web-access, pi-subagents, rpiv-ask-user-question, rpiv-todo, rpiv-web-tools
    * local LLM: Ollama + vLLM auto-detect in /login
    * permission: gotgenes gates + inobit ask/plan UX, config at ~/skill-pi/permission.json
    * theme: tokyonight + skillgrid logo (omegon-style splash)
opencode plugins
    * hooks plugin allowing the execution of `/.cursor/hooks/`
    * mnemonic memory plugin similare to engram
    * tui: skillgrid logo
vscode plugin
    * visualize mnemonic memory
    * proposed plugins
        * ysamlan.vscode-backlog-md
        * kilocode.Kilo-Code
        * zgy.opencode-vscode-ui
        * kwickramasekara.opencode-chat-unofficial
        * LanTingxin.opencode-enhanced-ui
hooks
  * 
agents
    * add hermes ???
        * hermes desktop
        * hermes webui
local olama
    * indexing and gemini for llmviki
obsidian integration
    * local obsidian vault vs docs folder
    * canvas integration
    * llmWiki
        * https://youtu.be/znj-WpMj1dI?si=O970kJy9YnnExIXq 
        * https://www.youtube.com/watch?v=_bieksxg6oY
        * https://www.patreon.com/wanderloots/posts/build-trusted-in-168046073
        * https://youtu.be/T33iI6izAKw?si=AwvqmXykSyY_uTWk
    * google Open Knowledge Format (OKF)
        * https://youtu.be/iod1M9Dc3HI?si=GPJKBRV6F3vXaifV
skillgrid-cli
    * init on project
    * add proposed vscode plugins
mnemonic
    * obsidian export
    * llmWiki
tui:
    * tmux
    * backlog tasks view
    * opensessions
webui:
    * nmemonic memory visualization
    * include opencode web
    * include backlog browser