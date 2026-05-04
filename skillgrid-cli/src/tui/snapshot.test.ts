import { promises as fs } from "node:fs";
import os from "node:os";
import path from "node:path";
import { afterEach, describe, expect, it } from "vitest";
import { buildObserveSnapshot } from "./snapshot.js";

const roots: string[] = [];

afterEach(async () => {
  await Promise.all(roots.map((r) => fs.rm(r, { recursive: true, force: true })));
  roots.length = 0;
});

async function tempRepo(): Promise<string> {
  const root = await fs.mkdtemp(path.join(os.tmpdir(), "skillgrid-tui-"));
  roots.push(root);
  await fs.mkdir(path.join(root, ".skillgrid", "prd"), { recursive: true });
  await fs.mkdir(path.join(root, ".skillgrid", "tasks", "events"), { recursive: true });
  await fs.writeFile(path.join(root, ".skillgrid", "prd", "PRD01_x.md"), "# x\n", "utf8");
  await fs.writeFile(
    path.join(root, ".skillgrid", "tasks", "events", "c1.jsonl"),
    '{"time":"2026-05-02T10:00:00Z","changeId":"c1","phase":"apply","status":"ok","summary":"done"}\n',
    "utf8"
  );
  await fs.writeFile(
    path.join(root, ".skillgrid", "tasks", "checkpoints.log"),
    "2026-05-01T12:00:00Z name=cp1 branch=main sha=abc dirty=no\n",
    "utf8"
  );
  await fs.writeFile(path.join(root, ".skillgrid", "tasks", "context_c1.md"), "handoff", "utf8");
  await fs.mkdir(path.join(root, "openspec", "changes", "my-change"), { recursive: true });
  return root;
}

describe("buildObserveSnapshot", () => {
  it("collects PRDs, changes, events, checkpoints, and handoffs", async () => {
    const root = await tempRepo();
    const s = await buildObserveSnapshot(root);
    expect(s.repoName).toBe(path.basename(root));
    expect(s.prdBasenames).toContain("PRD01_x.md");
    expect(s.changeIds).toContain("my-change");
    expect(s.eventLines.length).toBeGreaterThan(0);
    expect(s.checkpointLines.length).toBeGreaterThan(0);
    expect(s.handoffBasenames).toContain("context_c1.md");
  });
});
