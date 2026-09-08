# Interview — 009-skillgrid-pi

Date: 2026-09-07
Source: questioning rounds after `sdd-explore` (TASK-006)

## Classification

Architectural — new fat Pi distribution + installer agent.

## Locked decisions

| # | Topic | Choice |
|---|--------|--------|
| 1 | Packaging approach | **B** fat distribution wrapper (not A compose-only, not C omp fork) |
| 2 | Wrapper depth | **2** fat (depend/bundle `@earendil-works/pi-coding-agent`) |
| 3 | CLI identity | **1** `skillgrid-pi` only — no `pi` alias |
| 4 | v1 scope | **1** Core (SDD + Mnemonic + install/update) |
| 5 | Repo home | Monorepo root **`skillgrid-pi/`** |
| 6 | Install path | Both npm `-g` and `skillgrid install`; agent option **`skillgrid-pi`** |
| 7 | Asset shipping | Extensions packaged; prompts/skills/agents/chains sync to **`~/.skillgrid-pi`** |
| 8 | npm name | **`skillgrid-pi`** |
| Confirm | Shared understanding | **Yes** → proceed to propose |
| 9 | Branding revise (2026-09-07) | Add step **07-branding-look**: Tokyo Night theme + Skillgrid splash/logo (omegon-style); reuse OpenCode/Kilo logo identity |
| 10 | Local API revise (2026-09-07) | Add step **08-local-api-login**: `/login` → Local API for Ollama/vLLM (inspired by Crossbar / pi-localllm-provider, not a separate slash command) |
| 11 | rpiv companions revise (2026-09-07) | Add step **09-rpiv-companions**: pin/activate `@juicesharp/rpiv-ask-user-question`, `@juicesharp/rpiv-todo`, `@juicesharp/rpiv-web-tools` from [rpiv-mono](https://github.com/juicesharp/rpiv-mono); not full `rpiv-pi` |
| 12 | Permission revise (2026-09-07) | Add step **10-permission**: gotgenes-style gates ([pi-permission-system](https://pi.dev/packages/@gotgenes/pi-permission-system)), config **`~/.skillgrid-pi/permission.json`** (user wrote `~/skillgrid-pi/permission.json`; locked to agent home), ask dialog + `/plan`/`/build` like [@inobit/pi-permission](https://github.com/inobit/pi-packages) |

## Explicit non-goals (user)

Dashboard/cleave/inference kitchen-sink, RDD, subagent-orchestration package in v1, `pi` alias, sibling repo for v1, full `@juicesharp/rpiv-pi` pipeline, gotgenes/inobit `~/.pi/agent` config ownership, OS sandbox. Splash/theme branding, `/login` Local API, three rpiv companions, and Skillgrid permission (ask + plan mode) are **in** scope after revise.
