import path from "node:path";
import { listDirectories, listFilesRecursive, pathExists, readTextIfExists } from "./fs-read.js";

export type ObserveSnapshot = {
  repoRoot: string;
  repoName: string;
  prdBasenames: string[];
  changeIds: string[];
  eventLines: string[];
  checkpointLines: string[];
  handoffBasenames: string[];
  generatedAt: string;
};

const MAX_EVENT_LINES = 48;
const MAX_CHECKPOINT_LINES = 24;

function tailLines(text: string, max: number): string[] {
  const lines = text.split(/\r?\n/).filter((l) => l.trim().length > 0);
  return lines.slice(-max);
}

type ParsedEvent = { time: string; line: string };

function parseEventLine(raw: string, fileHint: string): ParsedEvent | undefined {
  const trimmed = raw.trim();
  if (!trimmed) return undefined;
  try {
    const obj = JSON.parse(trimmed) as Record<string, unknown>;
    const time = typeof obj.time === "string" ? obj.time : "";
    const changeId = typeof obj.changeId === "string" ? obj.changeId : "";
    const phase = typeof obj.phase === "string" ? obj.phase : "";
    const status = typeof obj.status === "string" ? obj.status : "";
    const summary = typeof obj.summary === "string" ? obj.summary : "";
    const node = typeof obj.node === "string" ? obj.node : "";
    const bits = [time, changeId || fileHint, phase, node, status, summary].filter(Boolean);
    return { time: time || "?", line: bits.join(" | ").slice(0, 500) };
  } catch {
    return { time: "?", line: `${fileHint}: ${trimmed.slice(0, 400)}` };
  }
}

export async function buildObserveSnapshot(repoRoot: string): Promise<ObserveSnapshot> {
  const root = path.resolve(repoRoot);
  const prdDir = path.join(root, ".skillgrid", "prd");
  const eventsDir = path.join(root, ".skillgrid", "tasks", "events");
  const tasksDir = path.join(root, ".skillgrid", "tasks");
  const checkpointsPath = path.join(tasksDir, "checkpoints.log");
  const changesRoot = path.join(root, "openspec", "changes");

  const prdFiles = (await listFilesRecursive(prdDir))
    .filter((f) => f.endsWith(".md") && path.basename(f) !== "INDEX.md")
    .map((f) => path.basename(f))
    .sort((a, b) => a.localeCompare(b));

  const activeDirs = await listDirectories(changesRoot);
  const archiveRoot = path.join(changesRoot, "archive");
  const archiveDirs = (await pathExists(archiveRoot)) ? await listDirectories(archiveRoot) : [];
  const changeIds = [
    ...activeDirs
      .map((d) => path.basename(d))
      .filter((name) => name !== "archive"),
    ...archiveDirs.map((d) => path.basename(d))
  ].sort((a, b) => a.localeCompare(b));

  const jsonlFiles = (await listFilesRecursive(eventsDir)).filter((f) => f.endsWith(".jsonl"));
  const parsed: ParsedEvent[] = [];
  for (const file of jsonlFiles) {
    const content = (await readTextIfExists(file)) ?? "";
    const hint = path.basename(file, ".jsonl");
    for (const line of content.split(/\r?\n/)) {
      const ev = parseEventLine(line, hint);
      if (ev) parsed.push(ev);
    }
  }
  parsed.sort((a, b) => (b.time ?? "").localeCompare(a.time ?? ""));
  const eventLines = parsed.slice(0, MAX_EVENT_LINES).map((p) => p.line);

  const checkpointRaw = (await readTextIfExists(checkpointsPath)) ?? "";
  const checkpointLines = tailLines(checkpointRaw, MAX_CHECKPOINT_LINES);

  const handoffFiles = (await listFilesRecursive(tasksDir)).filter((f) => /^context_.+\.md$/.test(path.basename(f)));
  const handoffBasenames = handoffFiles.map((f) => path.basename(f)).sort((a, b) => a.localeCompare(b));

  return {
    repoRoot: root,
    repoName: path.basename(root),
    prdBasenames: prdFiles,
    changeIds,
    eventLines,
    checkpointLines,
    handoffBasenames,
    generatedAt: new Date().toISOString()
  };
}
