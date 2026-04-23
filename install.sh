#!/bin/bash
# =============================================================================
# install.sh — AI IDE Config Installer
# =============================================================================
#
# PURPOSE:
#   Copies AI assistant configuration folders (Cursor, Copilot, Kilo, OpenCode,
#   Antigravity) from this hub repo into a target project directory.
#   Merges MCP server configs from .configs/mcp/ and normalizes them per-IDE.
#
# USAGE:
#   ./install.sh [OPTIONS]
#
#   Common workflows:
#     # Interactive — pick IDEs and MCP servers via prompts
#     ./install.sh -p /path/to/project
#
#     # Non-interactive — all IDEs, all MCP servers
#     ./install.sh -p /path/to/project -A -y
#
#     # Dry run — preview changes without writing
#     ./install.sh -p /path/to/project -n
#
#     # Check/install dependencies only
#     ./install.sh -d
#
#     # Uninstall managed folders from a project
#     ./install.sh -p /path/to/project -u
#
# ARCHITECTURE:
#   1. Argument parsing  →  populate SELECTED_IDES, flags, PROJECT_PATH
#   2. Interactive prompts → IDE selection, MCP server selection (if eligible)
#   3. Dependency check  → optional install of missing tools
#   4. MCP merge         → jq-smerge .configs/mcp/*.json into one object
#   5. IDE setup         → per-IDE setup_*() writes normalized config files
#   6. Optional tools    → -t selects tools; CLIs via uv (Python) + hub npm (Node) + brew, then copy + openspec init
#
# DEPENDENCIES:
#   Runtime:  bash 3.2+ (incl. macOS /bin/bash), rsync, jq
#   Optional: node, npx, python3, pip3, uv (Python CLIs; hub Node deps via npm ci)
#   IDE CLIs: opencode, kilo, semgrep, trivy (installed on demand)
#
# VERSION: 1.0.0
# =============================================================================

set -e

VERSION="1.0.0"
PROJECT_PATH=""
SELECTED_IDES=()
SELECTED_TOOLS=()
TOOLS_INTERACTIVE=false
DRY_RUN=false
UNINSTALL=false
CHECK_DEPS=false
ALL_IDES=false
NON_INTERACTIVE=false
MERGE_MCP=true
MCP_KEY_FILTER_JSON=""

# =============================================================================
# GLOBALS: Colors, IDE mappings, dependency declarations
# =============================================================================

# Color support — only emit ANSI codes when stdout is a terminal
if [ -t 1 ]; then
    RED='\033[0;31m'
    GREEN='\033[0;32m'
    YELLOW='\033[1;33m'
    BLUE='\033[0;34m'
    CYAN='\033[0;36m'
    NC='\033[0m'
else
    RED=''
    GREEN=''
    YELLOW=''
    BLUE=''
    CYAN=''
    NC=''
fi

# IDE folder name mapping (case-based — bash 3.2 on macOS has no declare -A)
ide_folder_for() {
    case "$1" in
        cursor) echo ".cursor" ;;
        copilot) echo ".vscode" ;;
        kilo) echo ".kilo" ;;
        opencode) echo ".opencode" ;;
        antigravity) echo ".agents" ;;
        *) echo "" ;;
    esac
}

# Core dependencies — always checked regardless of IDE selection
# Format: "name:brew_package:pip_package:npm_package:check_command"
declare -a CORE_DEPENDENCIES=(
    # Core tools
    "rsync:rsync:::rsync --version"
    "jq:jq:::jq --version"
    "node:node:::node --version"
    "python3:python@3.12:::python3 --version"
    "pip3::::pip3 --version"
    "npx::::npx --version"
)

# IDE-specific dependencies — only checked when the corresponding IDE is selected
# Format: "name|brew_package|pip_package|npm_package|check_command|ide_flag"
declare -a IDE_DEPENDENCIES=(
    "opencode|opencode|||opencode --version|opencode"
    "openspec|openspec|||openspec --version|openspec"
    "kilo|Kilo-Org/homebrew-tap/kilo|||kilo --version|kilo"
    "semgrep|semgrep|||semgrep --version|semgrep"
    "trivy|trivy|||trivy --version|trivy"
    "trivy-mcp|trivy|||trivy plugin list|trivy-mcp"
)

# =============================================================================
# UTILITY FUNCTIONS: logging, help, version
# =============================================================================

show_help() {
    cat << EOF
Usage:
  $(basename "$0") [OPTIONS]

Options:
  -p, --path <dir>      Install to a specific project path
  -c, --cursor          Setup configuration for Cursor
  -C, --copilot         Setup configuration for Copilot
  -k, --kilo            Setup configuration for Kilocode
  -o, --opencode        Setup configuration for opencode
  -a, --antigravity     Setup configuration for Google Antigravity
  -A, --all             Setup for all supported IDEs (Default if none selected)
  -t, --tools           Interactive prompt to select optional tools (openspec, graphify, cocoindex-code, dmux, engram)
  -d, --deps            Check and install dependencies before install
  -y, --yes             Non-interactive mode (skip prompts)
  --no-mcp              Skip MCP server configuration
  -n, --dry-run         Show what would be installed without making changes
  -u, --uninstall       Remove .ai-config and managed IDE dirs from target
  -v, --version         Print Version
  -h, --help            Show this help message

Interactive mode: On TTY with no IDE flags, choose IDEs (1-5 or a=all) and MCP servers.
Use -t to pick optional tools interactively (openspec, graphifyy, cocoindex-code, dmux, engram — installed via brew / hub npm / uv when selected; see docs/tools.md).
EOF
}

show_version() {
    echo "install.sh version $VERSION"
}

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[OK]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# =============================================================================
# INTERACTIVE PROMPTS: IDE selection, MCP selection
# =============================================================================

# Check if interactive mode is eligible
interactive_eligible() {
    # Not interactive if: NON_INTERACTIVE, CI env, not a TTY, or IDE flags already set
    [ "$NON_INTERACTIVE" != true ] || return 1
    case "${CI:-}" in
        true|1|yes|YES) return 1 ;;
    esac
    [ -t 0 ] && [ -t 1 ] || return 1
    [ ${#SELECTED_IDES[@]} -eq 0 ] || return 1
    return 0
}

# Interactive IDE selection
interactive_ide_selection() {
    interactive_eligible || return 0

    echo ""
    echo -e "${CYAN}IDE integration${NC} — symlink the hub into which tools?"
    echo "  1) Cursor (.cursor/)"
    echo "  2) Copilot (.vscode/)"
    echo "  3) Kilocode (.kilo/)"
    echo "  4) OpenCode (.opencode/)"
    echo "  5) Antigravity (.agents/)"
    echo ""
    echo "  a — all five   |   e.g. 1,3,5 — only those numbers"
    echo ""

    local choice
    while true; do
        if ! read -r -p "IDE choice [a]: " choice; then
            echo ""
            ALL_IDES=true
            SELECTED_IDES=("cursor" "copilot" "kilo" "opencode" "antigravity")
            log_info "IDE: all five (default)"
            return 0
        fi
        choice=$(printf '%s' "$choice" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
        if [ -z "$choice" ]; then
            ALL_IDES=true
            SELECTED_IDES=("cursor" "copilot" "kilo" "opencode" "antigravity")
            log_info "IDE: all five (default)"
            return 0
        fi
        local lower
        lower=$(printf '%s' "$choice" | tr '[:upper:]' '[:lower:]')
        case "$lower" in
            a|all)
                ALL_IDES=true
                SELECTED_IDES=("cursor" "copilot" "kilo" "opencode" "antigravity")
                log_info "IDE: all five"
                return 0
                ;;
        esac

        SELECTED_IDES=()
        local -a parts
        local tok invalid=""
        IFS=',' read -ra parts <<< "$choice"
        for tok in "${parts[@]}"; do
            tok=$(printf '%s' "$tok" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
            [ -z "$tok" ] && continue
            case "$tok" in
                1) SELECTED_IDES+=("cursor") ;;
                2) SELECTED_IDES+=("copilot") ;;
                3) SELECTED_IDES+=("kilo") ;;
                4) SELECTED_IDES+=("opencode") ;;
                5) SELECTED_IDES+=("antigravity") ;;
                *) invalid="invalid index: $tok (use 1–5 or a)"; break ;;
            esac
        done

        if [ -n "$invalid" ]; then
            log_warn "$invalid"
            continue
        fi
        if [ ${#SELECTED_IDES[@]} -eq 0 ]; then
            log_warn "Pick at least one number (1–5) or a for all"
            continue
        fi
        log_info "IDE: selected ${#SELECTED_IDES[@]} tool(s)"
        return 0
    done
}

# Get available MCP servers from .configs/mcp/
get_available_mcp_servers() {
    local mcp_dir="$SCRIPT_DIR/.configs/mcp"
    local servers=()

    if [ -d "$mcp_dir" ]; then
        # Collect all JSON files recursively
        local -a merge_paths=()
        local main_mcp="$SCRIPT_DIR/.configs/mcp.json"
        [ -f "$main_mcp" ] && merge_paths+=("$main_mcp")

        local f
        while IFS= read -r -d '' f; do
            merge_paths+=("$f")
        done < <(find "$mcp_dir" -name '*.json' -type f -print0 2>/dev/null | sort -z)

        # Extract unique server keys from all files, one per line
        if [ ${#merge_paths[@]} -gt 0 ] && command -v jq &>/dev/null; then
            jq -s '[.[] | (if (type == "object") and has("mcpServers") then .mcpServers else . end) | keys[]] | unique | reverse | .[]' "${merge_paths[@]}" 2>/dev/null
        fi
    fi
}

# Interactive MCP server selection
interactive_mcp_selection() {
    # Skip if --no-mcp was passed
    [ "$MERGE_MCP" != false ] || return 0
    # Not interactive if: NON_INTERACTIVE, CI env, not a TTY
    [ "$NON_INTERACTIVE" != true ] || return 0
    case "${CI:-}" in
        true|1|yes|YES) return 0 ;;
    esac
    [ -t 0 ] && [ -t 1 ] || return 0

    local mcp_dir="$SCRIPT_DIR/.configs/mcp"
    local main_mcp="$SCRIPT_DIR/.configs/mcp.json"

    # Check if any MCP configs exist
    local has_configs=false
    [ -f "$main_mcp" ] && has_configs=true
    [ -d "$mcp_dir" ] && [ -n "$(ls -A "$mcp_dir" 2>/dev/null)" ] && has_configs=true

    if [ "$has_configs" = false ]; then
        log_info "MCP: no .configs/mcp/ directory found — skipping MCP selection"
        return 0
    fi

    # Check for jq
    if ! command -v jq &>/dev/null; then
        log_warn "MCP: jq not found — skipping MCP selection (install jq or use --no-mcp)"
        return 0
    fi

    echo ""
    echo -e "${CYAN}MCP Servers${NC} — which servers to enable?"
    echo "  a — all servers   |   n — skip MCP   |   e.g. 1,3,5 — subset"
    echo ""

    # List available servers
    local -a servers
    local i=1
    while IFS= read -r server; do
        if [ -n "$server" ] && [ "$server" != "null" ]; then
            servers+=("$server")
            echo "  $i) $server"
            i=$((i + 1))
        fi
    done < <(get_available_mcp_servers)

    if [ ${#servers[@]} -eq 0 ]; then
        log_info "MCP: no servers found — skipping"
        return 0
    fi

    echo ""

    local choice
    while true; do
        if ! read -r -p "MCP choice [a]: " choice; then
            echo ""
            log_info "MCP: all servers (default)"
            return 0
        fi
        choice=$(printf '%s' "$choice" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
        if [ -z "$choice" ]; then
            log_info "MCP: all servers (default)"
            return 0
        fi
        local lower
        lower=$(printf '%s' "$choice" | tr '[:upper:]' '[:lower:]')
        case "$lower" in
            a|all)
                log_info "MCP: all servers"
                return 0
                ;;
            n|no|skip)
                MERGE_MCP=false
                log_info "MCP: skipped"
                return 0
                ;;
        esac

        local -a selected_indices
        local -a parts
        local tok invalid=""
        IFS=',' read -ra parts <<< "$choice"
        for tok in "${parts[@]}"; do
            tok=$(printf '%s' "$tok" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
            [ -z "$tok" ] && continue
            if [[ "$tok" =~ ^[0-9]+$ ]] && [ "$tok" -ge 1 ] && [ "$tok" -le "${#servers[@]}" ]; then
                selected_indices+=("$tok")
            else
                invalid="invalid index: $tok (use 1–${#servers[@]}, a, or n)"; break
            fi
        done

        if [ -n "$invalid" ]; then
            log_warn "$invalid"
            continue
        fi
        if [ ${#selected_indices[@]} -eq 0 ]; then
            log_warn "Pick at least one number, a for all, or n to skip"
            continue
        fi

        # Build JSON filter for selected servers
        local -a selected_names
        for idx in "${selected_indices[@]}"; do
            local name="${servers[$((idx-1))]}"
            selected_names+=("$name")
        done

        # Create JSON array
        MCP_KEY_FILTER_JSON="["
        local first=true
        for name in "${selected_names[@]}"; do
            if [ "$first" = true ]; then
                MCP_KEY_FILTER_JSON+="\"$name\""
                first=false
            else
                MCP_KEY_FILTER_JSON+=",\"$name\""
            fi
        done
        MCP_KEY_FILTER_JSON+="]"

        log_info "MCP: selected ${#selected_names[@]} server(s)"
        return 0
    done
}

# Return 0 if optional tool id is in SELECTED_TOOLS
tool_is_selected() {
    local id="$1"
    local t
    for t in "${SELECTED_TOOLS[@]}"; do
        [ "$t" = "$id" ] && return 0
    done
    return 1
}

# Ensure uv is available (for cocoindex-code / graphifyy: uv tool install)
ensure_uv() {
    command -v uv &>/dev/null && return 0
    if [ "$DRY_RUN" = true ]; then
        echo "[DRY-RUN] brew install uv  (or: curl -LsSf https://astral.sh/uv/install.sh | sh)"
        return 0
    fi
    log_info "Installing uv..."
    if command -v brew &>/dev/null; then
        brew install uv || true
    elif command -v curl &>/dev/null; then
        curl -LsSf https://astral.sh/uv/install.sh | sh || true
    else
        log_warn "uv not found — install with: brew install uv (or https://docs.astral.sh/uv/)"
        return 1
    fi
    command -v uv &>/dev/null
}

# Run npm ci in the hub when package.json exists (pins openspec, dmux, MCP packages; see package.json)
hub_ensure_npm_install() {
    [ -n "${SCRIPT_DIR:-}" ] && [ -f "$SCRIPT_DIR/package.json" ] || return 1
    if [ "$DRY_RUN" = true ]; then
        echo "[DRY-RUN] (cd \"$SCRIPT_DIR\" && npm ci)"
        return 0
    fi
    log_info "Hub Node dependencies (npm ci in $SCRIPT_DIR)..."
    (cd "$SCRIPT_DIR" && npm ci) 2>/dev/null || (cd "$SCRIPT_DIR" && npm install) || {
        log_warn "hub npm ci failed — check Node, network, and package-lock.json (see docs/tools.md)"
        return 1
    }
    return 0
}

# Install openspec CLI: prefer PATH, else Homebrew, else hub node_modules, else npm -g
install_openspec_cli() {
    if command -v openspec &>/dev/null; then
        log_info "openspec CLI already on PATH"
        return 0
    fi
    if [ "$DRY_RUN" = true ]; then
        echo "[DRY-RUN] brew install openspec  ||  hub npm ci  ||  npm install -g @fission-ai/openspec@latest"
        return 0
    fi
    if command -v brew &>/dev/null; then
        if brew install openspec; then
            log_success "openspec installed (Homebrew)"
            return 0
        fi
    fi
    if hub_ensure_npm_install && [ -x "$SCRIPT_DIR/node_modules/.bin/openspec" ]; then
        log_success "openspec available: cd \"$SCRIPT_DIR\" && npx openspec (or add node_modules/.bin to PATH)"
        return 0
    fi
    if command -v npm &>/dev/null; then
        log_info "Installing openspec via npm global (@fission-ai/openspec)..."
        if npm install -g @fission-ai/openspec@latest; then
            log_success "openspec installed (npm -g)"
            return 0
        fi
    fi
    log_warn "openspec: install manually — brew install openspec  OR  npm ci in hub  OR  npm install -g @fission-ai/openspec@latest"
    return 1
}

# Install CLIs for SELECTED_TOOLS (uv / hub npm / brew)
install_optional_tool_clis() {
    [ ${#SELECTED_TOOLS[@]} -eq 0 ] && return 0

    echo ""
    echo "Optional tools — installing CLIs..."
    echo ""

    if tool_is_selected cocoindex-code || tool_is_selected graphify; then
        ensure_uv || true
    fi

    if tool_is_selected cocoindex-code; then
        if command -v ccc &>/dev/null 2>&1 || command -v cocoindex &>/dev/null 2>&1 || command -v cocoindex-code &>/dev/null 2>&1; then
            log_info "cocoindex-code (ccc) CLI already present"
        elif [ "$DRY_RUN" = true ]; then
            echo "[DRY-RUN] uv tool install 'cocoindex-code[full]'"
        elif command -v uv &>/dev/null; then
            log_info "Installing cocoindex-code (uv tool install)..."
            if uv tool install 'cocoindex-code[full]'; then
                log_success "cocoindex-code installed (ccc; ensure uv tool path is on PATH)"
            else
                log_warn "cocoindex-code: uv tool install failed"
            fi
        else
            log_warn "cocoindex-code: uv missing — run: uv tool install 'cocoindex-code[full]'"
        fi
    fi

    if tool_is_selected graphify; then
        if command -v graphifyy &>/dev/null 2>&1 || command -v graphify &>/dev/null 2>&1; then
            log_info "graphify CLI already present"
        elif [ "$DRY_RUN" = true ]; then
            echo "[DRY-RUN] uv tool install graphifyy"
        elif command -v uv &>/dev/null; then
            log_info "Installing graphifyy (uv tool install)..."
            if uv tool install graphifyy; then
                log_success "graphifyy installed"
            else
                log_warn "graphifyy: uv tool install failed"
            fi
        else
            log_warn "graphify: uv missing — run: uv tool install graphifyy"
        fi
    fi

    if tool_is_selected openspec; then
        install_openspec_cli || true
    fi

    if tool_is_selected dmux; then
        if command -v dmux &>/dev/null; then
            log_info "dmux CLI already on PATH"
        elif [ "$DRY_RUN" = true ]; then
            echo "[DRY-RUN] hub npm ci  ||  npm install -g dmux"
        elif hub_ensure_npm_install && [ -x "$SCRIPT_DIR/node_modules/.bin/dmux" ]; then
            log_success "dmux available: cd \"$SCRIPT_DIR\" && npx dmux (or add node_modules/.bin to PATH)"
        elif command -v npm &>/dev/null; then
            log_info "Installing dmux (npm -g fallback)..."
            if npm install -g dmux; then
                log_success "dmux installed (npm -g)"
            else
                log_warn "dmux: npm install -g failed"
            fi
        else
            log_warn "dmux: npm not found — install Node.js"
        fi
    fi

    if tool_is_selected engram; then
        if command -v engram &>/dev/null; then
            log_info "engram CLI already installed"
        elif [ "$DRY_RUN" = true ]; then
            echo "[DRY-RUN] brew install gentleman-programming/tap/engram"
        elif command -v brew &>/dev/null; then
            log_info "Installing engram (Homebrew)..."
            if brew install gentleman-programming/tap/engram; then
                log_success "engram installed"
            else
                log_warn "engram: brew install failed"
            fi
        else
            log_warn "engram: Homebrew not found — run: brew install gentleman-programming/tap/engram"
        fi
    fi

    echo ""
}

# Interactive optional tools (openspec, graphify, cocoindex-code, dmux, engram)
interactive_tools_selection() {
    [ "$TOOLS_INTERACTIVE" = true ] || return 0
    [ "$NON_INTERACTIVE" != true ] || {
        log_info "Optional tools: skipping prompt (--yes); none selected"
        return 0
    }
    case "${CI:-}" in
        true|1|yes|YES)
            log_info "Optional tools: skipping prompt (CI)"
            return 0
            ;;
    esac
    [ -t 0 ] && [ -t 1 ] || {
        log_warn "Optional tools: not a TTY — skipping tool selection (use a terminal for -t)"
        return 0
    }

    echo ""
    echo -e "${CYAN}Optional tools${NC} — CLIs via uv, hub npm ci, or brew (see docs/tools.md)"
    echo "  1) openspec — OpenSpec (brew, hub npx, or npm -g)"
    echo "  2) graphify — graphifyy (uv tool install graphifyy)"
    echo "  3) cocoindex-code — CocoIndex (uv tool install cocoindex-code[full] → ccc)"
    echo "  4) dmux — tmux pane manager (hub npx, or npm -g fallback)"
    echo "  5) engram — Engram MCP CLI (brew gentleman-programming/tap)"
    echo ""
    echo "  a — all five   |   n — none   |   e.g. 1,3 — pick by number"
    echo ""

    local choice
    while true; do
        if ! read -r -p "Tool choice [n]: " choice; then
            echo ""
            log_info "Optional tools: none (default)"
            return 0
        fi
        choice=$(printf '%s' "$choice" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
        if [ -z "$choice" ]; then
            log_info "Optional tools: none (default)"
            return 0
        fi
        local lower
        lower=$(printf '%s' "$choice" | tr '[:upper:]' '[:lower:]')
        case "$lower" in
            a|all)
                SELECTED_TOOLS=("openspec" "graphify" "cocoindex-code" "dmux" "engram")
                log_info "Optional tools: openspec, graphify, cocoindex-code, dmux, engram"
                return 0
                ;;
            n|no|none|skip)
                log_info "Optional tools: none"
                return 0
                ;;
        esac

        SELECTED_TOOLS=()
        local tok invalid=""
        local -a parts
        IFS=',' read -ra parts <<< "$choice"
        for tok in "${parts[@]}"; do
            tok=$(printf '%s' "$tok" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
            [ -z "$tok" ] && continue
            case "$tok" in
                1) SELECTED_TOOLS+=("openspec") ;;
                2) SELECTED_TOOLS+=("graphify") ;;
                3) SELECTED_TOOLS+=("cocoindex-code") ;;
                4) SELECTED_TOOLS+=("dmux") ;;
                5) SELECTED_TOOLS+=("engram") ;;
                *) invalid="invalid index: $tok (use 1–5, a, or n)"; break ;;
            esac
        done

        if [ -n "$invalid" ]; then
            log_warn "$invalid"
            continue
        fi
        if [ ${#SELECTED_TOOLS[@]} -eq 0 ]; then
            log_warn "Pick at least one number (1–5), a for all, or n for none"
            continue
        fi
        log_info "Optional tools: selected ${#SELECTED_TOOLS[@]} tool(s)"
        return 0
    done
}

# Merge all MCP configs from .configs/mcp.json and .configs/mcp/**/*.json
# Returns merged JSON on stdout (log messages go to stderr)
merge_mcp_configs() {
    local mcp_dir="$SCRIPT_DIR/.configs/mcp"
    local main_mcp="$SCRIPT_DIR/.configs/mcp.json"

    if [ "$MERGE_MCP" = false ]; then
        return 0
    fi

    # Collect all JSON files
    local -a merge_paths=()
    [ -f "$main_mcp" ] && merge_paths+=("$main_mcp")
    
    local f
    while IFS= read -r -d '' f; do
        merge_paths+=("$f")
    done < <(find "$mcp_dir" -name '*.json' -type f -print0 2>/dev/null | sort -z)

    if [ ${#merge_paths[@]} -eq 0 ]; then
        log_info "No MCP configs found — skip MCP merge" >&2
        return 0
    fi

    if ! command -v jq &>/dev/null; then
        log_error "jq is required to merge MCP configs (install jq or pass --no-mcp)" >&2
        return 1
    fi

    log_info "Merging MCP configs (${#merge_paths[@]} file(s))" >&2

    # Merge all configs
    local merged_json
    merged_json=$(jq -s '
      def normalize:
        if (type == "object") and has("mcpServers") then { mcpServers: .mcpServers }
        else { mcpServers: . }
        end;
      reduce .[] as $item ({"mcpServers": {}}; ($item | normalize) as $n | { mcpServers: (.mcpServers * $n.mcpServers) })
    ' "${merge_paths[@]}" 2>/dev/null) || {
        log_error "MCP merge failed (jq)" >&2
        return 1
    }

    # Apply filter if subset selected
    if [ -n "$MCP_KEY_FILTER_JSON" ]; then
        merged_json=$(echo "$merged_json" | jq --argjson keys "$MCP_KEY_FILTER_JSON" '{ mcpServers: (.mcpServers | with_entries(select(.key as $k | ($keys | index($k) != null)))) }' 2>/dev/null) || {
            log_error "MCP subset filter failed (jq)" >&2
            return 1
        }
        local count
        count=$(echo "$merged_json" | jq '.mcpServers | keys | length')
        log_info "Applied MCP server subset ($count server(s))" >&2
    fi

    local server_count
    server_count=$(echo "$merged_json" | jq '.mcpServers | keys | length')
    log_success "Merged MCP: $server_count server(s)" >&2

    # Return the merged JSON
    echo "$merged_json"
}

show_dependencies() {
    echo "=== Dependency Check ==="
    echo ""
    
    local missing=()
    local installed=()
    
    # Check core dependencies
    echo "Core Dependencies:"
    for dep in "${CORE_DEPENDENCIES[@]}"; do
        IFS=':' read -r name brew pip npm check_cmd <<< "$dep"
        
        # Skip empty check commands
        if [ -z "$check_cmd" ]; then
            continue
        fi
        
        # Check if dependency is installed
        if eval "$check_cmd" &>/dev/null; then
            installed+=("$name")
            echo "  ✓ $name"
        else
            missing+=("$name|$brew|$pip|$npm|core")
            echo "  ✗ $name"
        fi
    done
    
    echo ""
    
    # Check IDE-specific dependencies
    echo "IDE-Specific Dependencies:"
    for dep in "${IDE_DEPENDENCIES[@]}"; do
        IFS='|' read -r name brew pip npm check_cmd ide_flag <<< "$dep"
        
        # Skip if this IDE is not selected
        local is_selected=false
        for selected in "${SELECTED_IDES[@]}"; do
            if [ "$selected" = "$ide_flag" ]; then
                is_selected=true
                break
            fi
        done
        
        # Special handling for openspec (optional tools)
        if [ "$ide_flag" = "openspec" ] && tool_is_selected openspec; then
            is_selected=true
        fi

        # Check if running with -A (all)
        if [ "$ALL_IDES" = true ]; then
            is_selected=true
        fi

        if [ "$is_selected" = false ]; then
            echo "  - $name (skipped, not selected)"
            continue
        fi
        
        # Skip empty check commands
        if [ -z "$check_cmd" ]; then
            continue
        fi
        
        # Check if dependency is installed
        if eval "$check_cmd" &>/dev/null; then
            installed+=("$name")
            echo "  ✓ $name"
        else
            missing+=("$name|$brew|$pip|$npm|$ide_flag")
            echo "  ✗ $name"
        fi
    done
    
    echo ""
    echo "Installed: ${#installed[@]}"
    echo "Missing: ${#missing[@]}"
    echo ""
    
    if [ ${#missing[@]} -gt 0 ]; then
        echo "Missing dependencies:"
        for dep in "${missing[@]}"; do
            IFS='|' read -r name brew pip npm ide_flag <<< "$dep"
            local install_cmd=""
            
            if [ -n "$brew" ]; then
                install_cmd="brew install $brew"
            elif [ -n "$npm" ]; then
                install_cmd="npm install -g $npm"
            elif [ -n "$pip" ]; then
                install_cmd="pip3 install $pip"
            fi
            
            if [ -n "$install_cmd" ]; then
                echo "  - $name → $install_cmd"
            else
                echo "  - $name → (manual install required)"
            fi
        done
        echo ""
    fi
}

install_dependencies() {
    echo "=== Installing Dependencies ==="
    echo ""
    
    local brew_pkgs=()
    local pip_pkgs=()
    local npm_pkgs=()
    
    # Find missing core dependencies
    for dep in "${CORE_DEPENDENCIES[@]}"; do
        IFS=':' read -r name brew pip npm check_cmd <<< "$dep"
        
        # Skip empty check commands
        if [ -z "$check_cmd" ]; then
            continue
        fi
        
        # Check if dependency is installed
        if ! eval "$check_cmd" &>/dev/null; then
            if [ -n "$brew" ]; then
                brew_pkgs+=("$brew")
            elif [ -n "$pip" ]; then
                pip_pkgs+=("$pip")
            elif [ -n "$npm" ]; then
                npm_pkgs+=("$npm")
            fi
        fi
    done
    
    # Find missing IDE-specific dependencies
    for dep in "${IDE_DEPENDENCIES[@]}"; do
        IFS='|' read -r name brew pip npm check_cmd ide_flag <<< "$dep"
        
        # Skip if this IDE is not selected
        local is_selected=false
        for selected in "${SELECTED_IDES[@]}"; do
            if [ "$selected" = "$ide_flag" ]; then
                is_selected=true
                break
            fi
        done
        
        # Special handling for openspec (optional tools)
        if [ "$ide_flag" = "openspec" ] && tool_is_selected openspec; then
            is_selected=true
        fi

        # Check if running with -A (all)
        if [ "$ALL_IDES" = true ]; then
            is_selected=true
        fi

        if [ "$is_selected" = false ]; then
            continue
        fi

        # Skip empty check commands
        if [ -z "$check_cmd" ]; then
            continue
        fi

        # Check if dependency is installed
        if ! eval "$check_cmd" &>/dev/null; then
            if [ -n "$brew" ]; then
                brew_pkgs+=("$brew")
            elif [ -n "$pip" ]; then
                pip_pkgs+=("$pip")
            elif [ -n "$npm" ]; then
                npm_pkgs+=("$npm")
            fi
        fi
    done
    
    # Install brew packages
    if [ ${#brew_pkgs[@]} -gt 0 ]; then
        echo "Installing with Homebrew: ${brew_pkgs[*]}"
        if [ "$DRY_RUN" = true ]; then
            echo "[DRY-RUN] brew install ${brew_pkgs[*]}"
        else
            # Check if brew is available
            if command -v brew &>/dev/null; then
                brew install "${brew_pkgs[@]}"
            else
                echo "  ⚠ Homebrew not found. Skipping brew packages."
                echo "  Install Homebrew: https://brew.sh"
            fi
        fi
        echo ""
    fi

    # Install trivy MCP plugin if trivy was installed/updated
    if [ "$DRY_RUN" != true ]; then
        local trivy_was_installed=true
        for dep in "${IDE_DEPENDENCIES[@]}"; do
            IFS='|' read -r name brew pip npm check_cmd ide_flag <<< "$dep"
            if [ "$name" = "trivy" ]; then
                if ! eval "$check_cmd" &>/dev/null; then
                    trivy_was_installed=false
                fi
                break
            fi
        done
        if [ "$trivy_was_installed" = false ] && command -v trivy &>/dev/null; then
            echo "Installing trivy MCP plugin..."
            trivy plugin install mcp
            echo ""
        fi
    fi
    
    # Install pip packages
    if [ ${#pip_pkgs[@]} -gt 0 ]; then
        echo "Installing with pip: ${pip_pkgs[*]}"
        if [ "$DRY_RUN" = true ]; then
            echo "[DRY-RUN] pip3 install ${pip_pkgs[*]}"
        else
            if command -v pip3 &>/dev/null; then
                pip3 install --user "${pip_pkgs[@]}"
            else
                echo "  ⚠ pip3 not found. Skipping pip packages."
            fi
        fi
        echo ""
    fi
    
    # Install npm packages
    if [ ${#npm_pkgs[@]} -gt 0 ]; then
        echo "Installing with npm: ${npm_pkgs[*]}"
        if [ "$DRY_RUN" = true ]; then
            echo "[DRY-RUN] npm install -g ${npm_pkgs[*]}"
        else
            if command -v npm &>/dev/null; then
                npm install -g "${npm_pkgs[@]}"
            else
                echo "  ⚠ npm not found. Skipping npm packages."
            fi
        fi
        echo ""
    fi
    
    if [ ${#brew_pkgs[@]} -eq 0 ] && [ ${#pip_pkgs[@]} -eq 0 ] && [ ${#npm_pkgs[@]} -eq 0 ]; then
        echo "All dependencies are already installed!"
        echo ""
    fi
}

prompt_install_deps() {
    local missing_count=0

    # Count missing core dependencies
    for dep in "${CORE_DEPENDENCIES[@]}"; do
        IFS=':' read -r name brew pip npm check_cmd <<< "$dep"

        if [ -z "$check_cmd" ]; then
            continue
        fi

        if ! eval "$check_cmd" &>/dev/null; then
            missing_count=$((missing_count + 1))
        fi
    done
    
    # Count missing IDE-specific dependencies
    for dep in "${IDE_DEPENDENCIES[@]}"; do
        IFS='|' read -r name brew pip npm check_cmd ide_flag <<< "$dep"
        
        # Skip if this IDE is not selected
        local is_selected=false
        for selected in "${SELECTED_IDES[@]}"; do
            if [ "$selected" = "$ide_flag" ]; then
                is_selected=true
                break
            fi
        done
        
        # Special handling for openspec (optional tools)
        if [ "$ide_flag" = "openspec" ] && tool_is_selected openspec; then
            is_selected=true
        fi
        
        # Check if running with -A (all)
        if [ "$ALL_IDES" = true ]; then
            is_selected=true
        fi
        
        if [ "$is_selected" = false ]; then
            continue
        fi
        
        if [ -z "$check_cmd" ]; then
            continue
        fi

        if ! eval "$check_cmd" &>/dev/null; then
            missing_count=$((missing_count + 1))
        fi
    done

    if [ $missing_count -gt 0 ]; then
        echo "Found $missing_count missing dependencies."
        echo ""
        read -p "Would you like to install missing dependencies? [y/N] " -n 1 -r
        echo ""
        
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            install_dependencies
        else
            echo "Skipping dependency installation."
            echo ""
        fi
    fi
}

# =============================================================================
# ARGUMENT PARSING
# =============================================================================

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -p|--path)
            PROJECT_PATH="$2"
            shift 2
            ;;
        -c|--cursor)
            SELECTED_IDES+=("cursor")
            shift
            ;;
        -C|--copilot)
            SELECTED_IDES+=("copilot")
            shift
            ;;
        -k|--kilo)
            SELECTED_IDES+=("kilo")
            shift
            ;;
        -o|--opencode)
            SELECTED_IDES+=("opencode")
            shift
            ;;
        -a|--antigravity)
            SELECTED_IDES+=("antigravity")
            shift
            ;;
        -A|--all)
            SELECTED_IDES=("cursor" "copilot" "kilo" "opencode" "antigravity")
            ALL_IDES=true
            shift
            ;;
        -t|--tools)
            TOOLS_INTERACTIVE=true
            shift
            ;;
        -d|--deps)
            CHECK_DEPS=true
            shift
            ;;
        -y|--yes)
            NON_INTERACTIVE=true
            shift
            ;;
        --no-mcp)
            MERGE_MCP=false
            shift
            ;;
        -n|--dry-run)
            DRY_RUN=true
            shift
            ;;
        -u|--uninstall)
            UNINSTALL=true
            shift
            ;;
        -v|--version)
            show_version
            exit 0
            ;;
        -h|--help)
            show_help
            exit 0
            ;;
        *)
            echo "Error: Unknown option: $1"
            echo ""
            show_help
            exit 1
            ;;
    esac
done

# Get the directory where this script is located (needed for interactive MCP selection)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Interactive IDE selection (if eligible)
interactive_ide_selection

# If no IDEs selected after interactive, default to all
if [ ${#SELECTED_IDES[@]} -eq 0 ]; then
    SELECTED_IDES=("cursor" "copilot" "kilo" "opencode" "antigravity")
    ALL_IDES=true
fi

# Interactive MCP selection (if eligible)
interactive_mcp_selection

# Optional tools — must run before --deps counts (openspec / graphify / cocoindex-code / dmux / engram)
interactive_tools_selection

# Handle --deps flag (check and optionally install)
if [ "$CHECK_DEPS" = true ]; then
    show_dependencies

    # Count missing core dependencies
    missing_count=0
    for dep in "${CORE_DEPENDENCIES[@]}"; do
        IFS=':' read -r name brew pip npm check_cmd <<< "$dep"
        if [ -z "$check_cmd" ]; then continue; fi
        if ! eval "$check_cmd" &>/dev/null; then
            missing_count=$((missing_count + 1))
        fi
    done

    # Count missing IDE-specific dependencies
    for dep in "${IDE_DEPENDENCIES[@]}"; do
        IFS='|' read -r name brew pip npm check_cmd ide_flag <<< "$dep"

        # Skip if this IDE is not selected
        is_selected=false
        for selected in "${SELECTED_IDES[@]}"; do
            if [ "$selected" = "$ide_flag" ]; then
                is_selected=true
                break
            fi
        done

        # Special handling for openspec (optional tools)
        if [ "$ide_flag" = "openspec" ] && tool_is_selected openspec; then
            is_selected=true
        fi

        # Check if running with -A (all)
        if [ "$ALL_IDES" = true ]; then
            is_selected=true
        fi

        if [ "$is_selected" = false ]; then continue; fi
        if [ -z "$check_cmd" ]; then continue; fi

        if ! eval "$check_cmd" &>/dev/null; then
            missing_count=$((missing_count + 1))
        fi
    done

    if [ $missing_count -gt 0 ]; then
        read -p "Would you like to install missing dependencies? [y/N] " -n 1 -r
        echo ""
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            install_dependencies
        fi
    fi

    # If only checking deps, exit here
    if [ -z "$PROJECT_PATH" ]; then
        exit 0
    fi
fi

# Validate project path
if [ -z "$PROJECT_PATH" ]; then
    echo "Error: Project path is required."
    echo ""
    show_help
    exit 1
fi

if [ ! -d "$PROJECT_PATH" ]; then
    echo "Error: Directory '$PROJECT_PATH' does not exist"
    exit 1
fi

# Resolve to absolute path
PROJECT_PATH=$(cd "$PROJECT_PATH" && pwd)

# =============================================================================
# IDE SETUP FUNCTIONS: per-IDE config generation
# =============================================================================

# Setup Cursor: write mcp.json directly
setup_cursor() {
    local target="$1"
    local merged_mcp="$2"
    local tool_dir="$target/.cursor"
    local mcp_file="$tool_dir/mcp.json"

    log_info "Setting up Cursor configuration..."
    ensure_dir "$tool_dir"

    if [ "$MERGE_MCP" = true ] && [ -n "$merged_mcp" ]; then
        echo "$merged_mcp" | jq '{mcpServers: .mcpServers}' > "$mcp_file"
        local count
        count=$(jq '.mcpServers | keys | length' "$mcp_file")
        log_success "Generated: $mcp_file ($count server(s))"
    else
        log_info "Skipping Cursor mcp.json (--no-mcp or no merged data)"
    fi

    log_success "Cursor setup complete"
}

# Setup Copilot (VSCode): write mcp.json directly
setup_copilot() {
    local target="$1"
    local merged_mcp="$2"
    local tool_dir="$target/.vscode"
    local mcp_file="$tool_dir/mcp.json"

    log_info "Setting up Copilot configuration..."
    ensure_dir "$tool_dir"

    if [ "$MERGE_MCP" = true ] && [ -n "$merged_mcp" ]; then
        echo "$merged_mcp" | jq '{servers: .mcpServers}' > "$mcp_file"
        local count
        count=$(jq '.servers | keys | length' "$mcp_file")
        log_success "Generated: $mcp_file ($count server(s))"
    else
        log_info "Skipping Copilot mcp.json (--no-mcp or no merged data)"
    fi

    log_success "Copilot setup complete"
}

# Setup Kilo: write mcp.json directly
setup_kilo() {
    local target="$1"
    local merged_mcp="$2"
    local tool_dir="$target/.kilo"
    local mcp_file="$tool_dir/mcp.json"

    log_info "Setting up Kilocode configuration..."
    ensure_dir "$tool_dir"

    # Install KiloCode CLI if not already installed
    if ! command -v kilo &>/dev/null; then
        log_info "Installing KiloCode CLI (@kilocode/cli)..."
        if [ "$DRY_RUN" = true ]; then
            echo "[DRY-RUN] npm install -g @kilocode/cli"
        else
            if command -v npm &>/dev/null; then
                npm install -g @kilocode/cli
                if [ $? -eq 0 ]; then
                    log_success "KiloCode CLI installed successfully"
                else
                    log_warn "KiloCode CLI installation failed — you may need to install it manually: npm install -g @kilocode/cli"
                fi
            else
                log_warn "npm not found — skipping KiloCode CLI installation"
            fi
        fi
    else
        log_info "KiloCode CLI already installed"
    fi

    if [ "$MERGE_MCP" = true ] && [ -n "$merged_mcp" ]; then
        echo "$merged_mcp" | jq '{mcpServers: .mcpServers}' > "$mcp_file"
        local count
        count=$(jq '.mcpServers | keys | length' "$mcp_file")
        log_success "Generated: $mcp_file ($count server(s))"
    else
        log_info "Skipping Kilo mcp.json (--no-mcp or no merged data)"
    fi

    log_success "Kilocode setup complete"
}

# Setup OpenCode: copy .configs/opencode.json then patch MCP servers into opencode.json
setup_opencode() {
    local target="$1"
    local merged_mcp="$2"
    local tool_dir="$target/.opencode"
    local mcp_file="$tool_dir/opencode.json"
    local hub_opencode="$SCRIPT_DIR/.configs/opencode.json"

    ensure_dir "$tool_dir"

    if [ -f "$hub_opencode" ]; then
        cp "$hub_opencode" "$mcp_file"
        log_success "Copied hub config -> $mcp_file"
    fi

    if [ ! -f "$mcp_file" ]; then
        log_info "Skipping opencode: $mcp_file missing (add .configs/opencode.json to hub or .opencode/)"
        return 0
    fi

    if [ "$MERGE_MCP" != true ] || [ -z "$merged_mcp" ]; then
        log_info "Skipping opencode MCP (--no-mcp or no merged data)"
        return 0
    fi

    if ! command -v jq &>/dev/null; then
        log_warn "jq not found - cannot patch $mcp_file"
        return 0
    fi

    local tmp_file=$(mktemp)
    # Convert to opencode format and patch
    local opencode_mcp=$(echo "$merged_mcp" | jq '.mcpServers | to_entries | map({
      key: .key,
      value: (
        if .value.type == "streamable-http" or .value.type == "http" then
          {enabled: true, type: "remote", url: .value.url}
        else
          {enabled: true, type: "local", command: ([.value.command] + .value.args)}
        end
      )
    }) | from_entries')

    if jq --argjson mcp "$opencode_mcp" '.mcp = $mcp' "$mcp_file" > "$tmp_file" 2>/dev/null; then
        mv "$tmp_file" "$mcp_file"
        local count
        count=$(jq '.mcp | keys | length' "$mcp_file")
        log_success "Updated MCP config: $mcp_file ($count server(s))"
    else
        rm -f "$tmp_file"
        log_warn "Could not update $mcp_file (jq failed)"
    fi
}

# Setup Antigravity: write mcp_config.json directly
setup_antigravity() {
    local target="$1"
    local merged_mcp="$2"
    local mcp_dir="$HOME/.gemini/antigravity"
    local mcp_file="$mcp_dir/mcp_config.json"

    log_info "Setting up Google Antigravity configuration..."
    ensure_dir "$mcp_dir"

    if [ "$MERGE_MCP" = true ] && [ -n "$merged_mcp" ]; then
        echo "$merged_mcp" | jq '{mcpServers: .mcpServers}' > "$mcp_file"
        local count
        count=$(jq '.mcpServers | keys | length' "$mcp_file")
        log_success "Generated: $mcp_file ($count server(s))"
    else
        log_info "Skipping Antigravity mcp_config.json (--no-mcp or no merged data)"
    fi

    log_success "Antigravity setup complete"
}

# Helper function to ensure directory exists
ensure_dir() {
    local dir="$1"
    if [ ! -d "$dir" ]; then
        if [ "$DRY_RUN" = true ]; then
            echo "[DRY-RUN] Would create directory: $dir"
        else
            mkdir -p "$dir"
        fi
    fi
}

# =============================================================================
# MAIN: orchestration — wiring, validation, execution
# =============================================================================

main() {
    # Items to exclude from copying (within each config folder)
    EXCLUDES=(
        ".git"
        ".gitignore"
        "docs"
        "configs"
    )

    # Build rsync exclude arguments
    RSYNC_EXCLUDES=""
    for exclude in "${EXCLUDES[@]}"; do
        RSYNC_EXCLUDES="$RSYNC_EXCLUDES --exclude '$exclude'"
    done

    # Uninstall mode
    if [ "$UNINSTALL" = true ]; then
        echo "Uninstalling AI config folders from: $PROJECT_PATH"
        echo ""

        # Remove managed IDE folders
        for ide in "${SELECTED_IDES[@]}"; do
            folder=$(ide_folder_for "$ide")
            target="$PROJECT_PATH/$folder"
            if [ -d "$target" ]; then
                if [ "$DRY_RUN" = true ]; then
                    echo "[DRY-RUN] Would remove: $folder"
                else
                    echo "Removing: $folder"
                    rm -rf "$target"
                fi
            else
                echo "Skipping: $folder (not found)"
            fi
        done

        # Remove optional tool paths when selected via -t (same session)
        if tool_is_selected openspec; then
            target="$PROJECT_PATH/openspec"
            if [ -d "$target" ]; then
                if [ "$DRY_RUN" = true ]; then
                    echo "[DRY-RUN] Would remove: openspec"
                else
                    echo "Removing: openspec"
                    rm -rf "$target"
                fi
            fi
        fi

        if tool_is_selected cocoindex-code; then
            target="$PROJECT_PATH/.cocoindex_code"
            if [ -d "$target" ]; then
                if [ "$DRY_RUN" = true ]; then
                    echo "[DRY-RUN] Would remove: .cocoindex_code"
                else
                    echo "Removing: .cocoindex_code"
                    rm -rf "$target"
                fi
            fi
        fi

        if tool_is_selected graphify; then
            target="$PROJECT_PATH/.graphify"
            if [ -d "$target" ]; then
                if [ "$DRY_RUN" = true ]; then
                    echo "[DRY-RUN] Would remove: .graphify"
                else
                    echo "Removing: .graphify"
                    rm -rf "$target"
                fi
            fi
        fi

        # Remove .github / .copilot if copilot was selected
        if [[ " ${SELECTED_IDES[*]} " =~ " copilot " ]]; then
            for extra in ".github" ".copilot"; do
                target="$PROJECT_PATH/$extra"
                if [ -d "$target" ]; then
                    if [ "$DRY_RUN" = true ]; then
                        echo "[DRY-RUN] Would remove: $extra"
                    else
                        echo "Removing: $extra"
                        rm -rf "$target"
                    fi
                fi
            done
        fi

        echo ""
        echo "Done!"
        exit 0
    fi

    # Check dependencies first (interactive prompt)
    if [ "$CHECK_DEPS" = false ]; then
        # Always check deps but don't auto-install
        echo "Checking dependencies..."
        echo ""

        missing_count=0
        for dep in "${CORE_DEPENDENCIES[@]}"; do
            IFS=':' read -r name brew pip npm check_cmd <<< "$dep"
            if [ -z "$check_cmd" ]; then continue; fi
            if ! eval "$check_cmd" &>/dev/null; then
                missing_count=$((missing_count + 1))
            fi
        done

        # Count missing IDE-specific dependencies
        for dep in "${IDE_DEPENDENCIES[@]}"; do
            IFS='|' read -r name brew pip npm check_cmd ide_flag <<< "$dep"

            # Skip if this IDE is not selected
            is_selected=false
            for selected in "${SELECTED_IDES[@]}"; do
                if [ "$selected" = "$ide_flag" ]; then
                    is_selected=true
                    break
                fi
            done

            # Special handling for openspec (optional tools)
            if [ "$ide_flag" = "openspec" ] && tool_is_selected openspec; then
                is_selected=true
            fi

            # Check if running with -A (all)
            if [ "$ALL_IDES" = true ]; then
                is_selected=true
            fi

            if [ "$is_selected" = false ]; then continue; fi
            if [ -z "$check_cmd" ]; then continue; fi

            if ! eval "$check_cmd" &>/dev/null; then
                missing_count=$((missing_count + 1))
            fi
        done

        if [ $missing_count -gt 0 ]; then
            echo "Found $missing_count missing dependencies."
            if [ "$NON_INTERACTIVE" = true ]; then
                echo "Continuing without installing dependencies (non-interactive mode)..."
                echo ""
            else
                read -p "Would you like to install missing dependencies? [y/N] " -n 1 -r
                echo ""
                if [[ $REPLY =~ ^[Yy]$ ]]; then
                    install_dependencies
                else
                    echo "Continuing without installing dependencies..."
                    echo ""
                fi
            fi
        else
            echo "All dependencies are installed!"
            echo ""
        fi
    fi

    # Optional tool CLIs (uv / hub npm / brew) — after dep prompt, before copying configs
    install_optional_tool_clis

    # Dry run info
    if [ "$DRY_RUN" = true ]; then
        echo "=== DRY RUN MODE ==="
        echo ""
    fi

    echo "Installing AI config folders to: $PROJECT_PATH"
    echo ""

    # Copy selected IDE folders
    for ide in "${SELECTED_IDES[@]}"; do
        folder=$(ide_folder_for "$ide")
        src="$SCRIPT_DIR/$folder"
        dst="$PROJECT_PATH/$folder"

        if [ -d "$src" ]; then
            if [ -d "$dst" ]; then
                if [ "$DRY_RUN" = true ]; then
                    echo "[DRY-RUN] Would update: $folder"
                else
                    echo "Updating: $folder"
                    eval rsync -av $RSYNC_EXCLUDES "$src/" "$dst/"
                fi
            else
                if [ "$DRY_RUN" = true ]; then
                    echo "[DRY-RUN] Would create: $folder"
                else
                    echo "Creating: $folder"
                    eval rsync -av $RSYNC_EXCLUDES "$src/" "$dst/"
                fi
            fi
        else
            echo "Skipping: $folder (not found in source)"
        fi
    done

    # Copy hub .configs/AGENTS.md → project root and each selected IDE config dir
    HUB_AGENTS="$SCRIPT_DIR/.configs/AGENTS.md"
    if [ -f "$HUB_AGENTS" ] && [ ${#SELECTED_IDES[@]} -gt 0 ]; then
        echo ""
        if [ "$DRY_RUN" = true ]; then
            echo "[DRY-RUN] Would copy .configs/AGENTS.md -> $PROJECT_PATH/AGENTS.md"
            for ide in "${SELECTED_IDES[@]}"; do
                folder=$(ide_folder_for "$ide")
                echo "[DRY-RUN] Would copy .configs/AGENTS.md -> $PROJECT_PATH/$folder/AGENTS.md"
            done
        else
            echo "Copying AGENTS.md from hub (.configs/AGENTS.md)..."
            cp "$HUB_AGENTS" "$PROJECT_PATH/AGENTS.md"
            log_success "Wrote AGENTS.md (project root)"
            for ide in "${SELECTED_IDES[@]}"; do
                folder=$(ide_folder_for "$ide")
                dst_dir="$PROJECT_PATH/$folder"
                ensure_dir "$dst_dir"
                cp "$HUB_AGENTS" "$dst_dir/AGENTS.md"
                log_success "Wrote $folder/AGENTS.md"
            done
        fi
    elif [ ! -f "$HUB_AGENTS" ] && [ ${#SELECTED_IDES[@]} -gt 0 ]; then
        log_info ".configs/AGENTS.md not found in hub — skipping AGENTS.md copy"
    fi

    # Sync .agents/skills/ to each IDE's skills directory
    echo ""
    echo "Syncing skills configurations..."
    SKILLS_SRC="$SCRIPT_DIR/.agents/skills"
    if [ -d "$SKILLS_SRC" ]; then
        for ide in "${SELECTED_IDES[@]}"; do
            case "$ide" in
                antigravity)
                    # Antigravity: sync .agents/skills/ → .agents/skills/
                    target="$PROJECT_PATH/.agents/skills"
                    if [ "$DRY_RUN" = true ]; then
                        echo "[DRY-RUN] Would sync skills to: .agents/skills"
                    else
                        ensure_dir "$target"
                        echo "Syncing skills to: .agents/skills"
                        rsync -av --delete "$SKILLS_SRC/" "$target/"
                    fi
                    ;;
                cursor)
                    # Cursor: sync .agents/skills/ → .cursor/.agents/skills/
                    target="$PROJECT_PATH/.cursor/.agents/skills"
                    if [ "$DRY_RUN" = true ]; then
                        echo "[DRY-RUN] Would sync skills to: .cursor/.agents/skills"
                    else
                        ensure_dir "$target"
                        echo "Syncing skills to: .cursor/.agents/skills"
                        rsync -av --delete "$SKILLS_SRC/" "$target/"
                    fi
                    ;;
                kilo)
                    # Kilo: sync .agents/skills/ → .kilo/skills/
                    target="$PROJECT_PATH/.kilo/skills"
                    if [ "$DRY_RUN" = true ]; then
                        echo "[DRY-RUN] Would sync skills to: .kilo/skills"
                    else
                        ensure_dir "$target"
                        echo "Syncing skills to: .kilo/skills"
                        rsync -av --delete "$SKILLS_SRC/" "$target/"
                    fi
                    ;;
                opencode)
                    # OpenCode: sync .agents/skills/ → .opencode/skills/
                    target="$PROJECT_PATH/.opencode/skills"
                    if [ "$DRY_RUN" = true ]; then
                        echo "[DRY-RUN] Would sync skills to: .opencode/skills"
                    else
                        ensure_dir "$target"
                        echo "Syncing skills to: .opencode/skills"
                        rsync -av --delete "$SKILLS_SRC/" "$target/"
                    fi
                    ;;
            esac
        done
    else
        log_info "Skills source not found: $SKILLS_SRC — skipping skills sync"
    fi

    # Sync .github/ and .copilot/ when Copilot is selected
    if [[ " ${SELECTED_IDES[*]} " =~ " copilot " ]]; then
        for extra in ".github" ".copilot"; do
            src="$SCRIPT_DIR/$extra"
            dst="$PROJECT_PATH/$extra"
            if [ -d "$src" ]; then
                if [ -d "$dst" ]; then
                    if [ "$DRY_RUN" = true ]; then
                        echo "[DRY-RUN] Would update: $extra"
                    else
                        echo "Updating: $extra"
                        rsync -av "$src/" "$dst/"
                    fi
                else
                    if [ "$DRY_RUN" = true ]; then
                        echo "[DRY-RUN] Would create: $extra"
                    else
                        echo "Creating: $extra"
                        rsync -av "$src/" "$dst/"
                    fi
                fi
            else
                echo "Skipping: $extra (not found in source)"
            fi
        done
    fi

    # Merge MCP configs
    echo ""
    echo "Merging MCP configurations..."
    MERGED_MCP=""
    if [ "$MERGE_MCP" = true ]; then
        MERGED_MCP=$(merge_mcp_configs)
    fi

    # Setup selected IDEs
    echo ""
    echo "Setting up IDE configurations..."
    for ide in "${SELECTED_IDES[@]}"; do
        case "$ide" in
            cursor)
                if [ "$DRY_RUN" = true ]; then
                    echo "[DRY-RUN] Would setup: Cursor"
                else
                    setup_cursor "$PROJECT_PATH" "$MERGED_MCP"
                fi
                ;;
            copilot)
                if [ "$DRY_RUN" = true ]; then
                    echo "[DRY-RUN] Would setup: Copilot"
                else
                    setup_copilot "$PROJECT_PATH" "$MERGED_MCP"
                fi
                ;;
            kilo)
                if [ "$DRY_RUN" = true ]; then
                    echo "[DRY-RUN] Would setup: Kilocode"
                else
                    setup_kilo "$PROJECT_PATH" "$MERGED_MCP"
                fi
                ;;
            opencode)
                if [ "$DRY_RUN" = true ]; then
                    echo "[DRY-RUN] Would setup: OpenCode"
                else
                    setup_opencode "$PROJECT_PATH" "$MERGED_MCP"
                fi
                ;;
            antigravity)
                if [ "$DRY_RUN" = true ]; then
                    echo "[DRY-RUN] Would setup: Antigravity"
                else
                    setup_antigravity "$PROJECT_PATH" "$MERGED_MCP"
                fi
                ;;
        esac
    done

    # Optional tool: openspec (hub → project)
    if tool_is_selected openspec && [ -d "$SCRIPT_DIR/openspec" ]; then
        target="$PROJECT_PATH/openspec"
        if [ -d "$target" ]; then
            if [ "$DRY_RUN" = true ]; then
                echo "[DRY-RUN] Would update: openspec"
            else
                echo "Updating: openspec"
                eval rsync -av $RSYNC_EXCLUDES "$SCRIPT_DIR/openspec/" "$target/"
            fi
        else
            if [ "$DRY_RUN" = true ]; then
                echo "[DRY-RUN] Would create: openspec"
            else
                echo "Creating: openspec"
                eval rsync -av $RSYNC_EXCLUDES "$SCRIPT_DIR/openspec/" "$target/"
            fi
        fi
    elif tool_is_selected openspec && [ ! -d "$SCRIPT_DIR/openspec" ]; then
        log_info "openspec: hub has no openspec/ — skip copy"
    fi

    # Optional tool: cocoindex-code
    if tool_is_selected cocoindex-code && [ -d "$SCRIPT_DIR/.cocoindex_code" ]; then
        target="$PROJECT_PATH/.cocoindex_code"
        if [ -d "$target" ]; then
            if [ "$DRY_RUN" = true ]; then
                echo "[DRY-RUN] Would update: .cocoindex_code"
            else
                echo "Updating: .cocoindex_code"
                rsync -av "$SCRIPT_DIR/.cocoindex_code/" "$target/"
            fi
        else
            if [ "$DRY_RUN" = true ]; then
                echo "[DRY-RUN] Would create: .cocoindex_code"
            else
                echo "Creating: .cocoindex_code"
                rsync -av "$SCRIPT_DIR/.cocoindex_code/" "$target/"
            fi
        fi
    elif tool_is_selected cocoindex-code && [ ! -d "$SCRIPT_DIR/.cocoindex_code" ]; then
        log_info "cocoindex-code: hub has no .cocoindex_code/ — skip copy"
    fi

    # Run openspec init when openspec tool selected
    if tool_is_selected openspec; then
        echo ""
        if [ "$DRY_RUN" = true ]; then
            echo "[DRY-RUN] Would run: openspec init in $PROJECT_PATH/openspec"
        else
            echo "Running openspec init..."
            if [ -d "$PROJECT_PATH/openspec" ]; then
                cd "$PROJECT_PATH/openspec" && openspec init
            fi
        fi
    fi

    echo ""
    if [ "$DRY_RUN" = true ]; then
        echo "=== DRY RUN COMPLETE ==="
        echo "(No changes were made)"
    else
        echo "Done! IDE config folders have been copied to $PROJECT_PATH"
    fi
}

# =============================================================================
# ENTRY POINT
# =============================================================================
main "$@"
