import { existsSync, readFileSync, readdirSync } from "node:fs";
import { join } from "node:path";
import type { Dirent } from "node:fs";

export type McpServersShape = { mcpServers: Record<string, unknown> };

function normalize(item: unknown): McpServersShape {
  if (item && typeof item === "object" && "mcpServers" in item) {
    const m = item as { mcpServers?: Record<string, unknown> };
    return { mcpServers: { ...(m.mcpServers ?? {}) } };
  }
  return { mcpServers: { ...(item as Record<string, unknown>) } };
}

export function mergeMcpJsonFiles(paths: string[]): McpServersShape {
  let acc: McpServersShape = { mcpServers: {} };
  for (const p of paths) {
    const raw = readFileSync(p, "utf8");
    const parsed = JSON.parse(raw) as unknown;
    const n = normalize(parsed);
    acc = { mcpServers: { ...acc.mcpServers, ...n.mcpServers } };
  }
  return acc;
}

export function collectMcpMergePaths(hubRoot: string): string[] {
  const mcpDir = join(hubRoot, ".configs", "mcp");
  const mainMcp = join(hubRoot, ".configs", "mcp.json");
  const paths: string[] = [];
  if (existsSync(mainMcp)) paths.push(mainMcp);
  if (existsSync(mcpDir)) {
    const files = listJsonRecursive(mcpDir).filter((p) => p !== mainMcp);
    paths.push(...files.sort());
  }
  return paths;
}

function listJsonRecursive(dir: string): string[] {
  const out: string[] = [];
  const walk = (d: string) => {
    let entries: Dirent[];
    try {
      entries = readdirSync(d, { withFileTypes: true });
    } catch {
      return;
    }
    for (const ent of entries) {
      const p = join(d, ent.name);
      if (ent.isDirectory()) walk(p);
      else if (ent.isFile() && ent.name.endsWith(".json")) out.push(p);
    }
  };
  walk(dir);
  return out;
}

export function getAvailableMcpServerKeys(merged: McpServersShape): string[] {
  const keys = Object.keys(merged.mcpServers ?? {});
  return [...new Set(keys)].sort().reverse();
}

export function filterMcpByKeys(merged: McpServersShape, keep: string[]): McpServersShape {
  const set = new Set(keep);
  const next: Record<string, unknown> = {};
  for (const k of Object.keys(merged.mcpServers)) {
    if (set.has(k)) next[k] = merged.mcpServers[k];
  }
  return { mcpServers: next };
}

/** Cursor .cursor/mcp.json */
export function emitMcpForCursor(merged: McpServersShape): McpServersShape {
  const mcpServers: Record<string, unknown> = {};
  for (const [k, v] of Object.entries(merged.mcpServers ?? {})) {
    mcpServers[k] = emitCursorServerValue(v);
  }
  return { mcpServers };
}

function emitCursorServerValue(v: unknown): unknown {
  if (!v || typeof v !== "object") return v;
  const o = v as Record<string, unknown>;
  if ("url" in o && !("command" in o)) {
    const out: Record<string, unknown> = { url: o.url };
    if ("headers" in o) out.headers = o.headers;
    if ("auth" in o) out.auth = o.auth;
    return out;
  }
  if ("command" in o) {
    const out: Record<string, unknown> = {
      type: "stdio",
      command: o.command,
      args: (o.args as unknown[]) ?? [],
    };
    if ("env" in o) out.env = o.env;
    if ("envFile" in o) out.envFile = o.envFile;
    if ("cwd" in o) out.cwd = o.cwd;
    return out;
  }
  return v;
}

/** VS Code / Copilot .vscode/mcp.json */
export function emitMcpForVscode(merged: McpServersShape): { servers: Record<string, unknown> } {
  const servers: Record<string, unknown> = {};
  for (const [k, v] of Object.entries(merged.mcpServers ?? {})) {
    servers[k] = emitVscodeServerValue(v);
  }
  return { servers };
}

function emitVscodeServerValue(v: unknown): unknown {
  if (!v || typeof v !== "object") return v;
  const o = v as Record<string, unknown>;
  if ("url" in o && !("command" in o)) {
    const t = o.type === "sse" ? "sse" : "http";
    const out: Record<string, unknown> = { type: t, url: o.url };
    if ("headers" in o) out.headers = o.headers;
    return out;
  }
  if ("command" in o) {
    const out: Record<string, unknown> = {
      type: "stdio",
      command: o.command,
      args: (o.args as unknown[]) ?? [],
    };
    if ("env" in o) out.env = o.env;
    if ("cwd" in o) out.cwd = o.cwd;
    return out;
  }
  return v;
}

/** OpenCode / Kilo mcp object */
export function emitOpencodeStyleMcpObject(merged: McpServersShape): Record<string, unknown> {
  const out: Record<string, unknown> = {};
  for (const [k, v] of Object.entries(merged.mcpServers ?? {})) {
    out[k] = emitOpencodeServerValue(v);
  }
  return out;
}

function emitOpencodeServerValue(v: unknown): unknown {
  if (!v || typeof v !== "object") return v;
  const o = v as Record<string, unknown>;
  if ("url" in o && !("command" in o)) {
    const out: Record<string, unknown> = { enabled: true, type: "remote", url: o.url };
    if ("headers" in o) out.headers = o.headers;
    return out;
  }
  if ("command" in o) {
    const cmd = o.command as string;
    const args = (o.args as string[]) ?? [];
    const out: Record<string, unknown> = {
      enabled: true,
      type: "local",
      command: [cmd, ...args],
    };
    if ("env" in o) out.environment = o.env;
    return out;
  }
  return v;
}

/** Antigravity ~/.gemini/antigravity/mcp_config.json */
export function emitMcpForAntigravity(merged: McpServersShape): McpServersShape {
  const mcpServers: Record<string, unknown> = {};
  for (const [k, v] of Object.entries(merged.mcpServers ?? {})) {
    mcpServers[k] = emitAntigravityServerValue(v);
  }
  return { mcpServers };
}

function emitAntigravityServerValue(v: unknown): unknown {
  if (!v || typeof v !== "object") return v;
  const o = v as Record<string, unknown>;
  if ("serverUrl" in o) return v;
  if ("url" in o && !("command" in o)) {
    const out: Record<string, unknown> = { serverUrl: o.url };
    if ("headers" in o) out.headers = o.headers;
    if ("env" in o) out.env = o.env;
    return out;
  }
  if ("command" in o) {
    const out: Record<string, unknown> = {
      command: o.command,
      args: (o.args as unknown[]) ?? [],
    };
    if ("env" in o) out.env = o.env;
    if ("cwd" in o) out.cwd = o.cwd;
    return out;
  }
  return v;
}
