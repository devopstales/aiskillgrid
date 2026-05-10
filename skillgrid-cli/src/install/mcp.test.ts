import test from "node:test";
import assert from "node:assert/strict";
import { mkdtempSync, writeFileSync, rmSync } from "node:fs";
import { join } from "node:path";
import { tmpdir } from "node:os";
import { mergeMcpJsonFiles, filterMcpByKeys, type McpServersShape } from "./mcp.js";

test("mergeMcpJsonFiles merges mcpServers from files", () => {
  const d = mkdtempSync(join(tmpdir(), "sg-mcp-"));
  try {
    writeFileSync(join(d, "x.json"), JSON.stringify({ mcpServers: { a: { command: "c", args: [] } } }));
    writeFileSync(join(d, "y.json"), JSON.stringify({ mcpServers: { b: { command: "d", args: [] } } }));
    const o = mergeMcpJsonFiles([join(d, "x.json"), join(d, "y.json")]);
    assert.ok("a" in o.mcpServers && "b" in o.mcpServers);
  } finally {
    rmSync(d, { recursive: true, force: true });
  }
});

test("filterMcpByKeys keeps subset", () => {
  const m: McpServersShape = {
    mcpServers: {
      a: { command: "x" },
      b: { command: "y" },
    },
  };
  const f = filterMcpByKeys(m, ["b"]);
  assert.deepEqual(Object.keys(f.mcpServers), ["b"]);
});
