import path from "node:path";
import type {
  BoardIssue,
  BoardLane,
  Checkpoint,
  DashboardData,
  DashboardEvent,
  PrdCard
} from "../dashboard/shared/types.js";

const MAX_DETAIL_LINES = 4000;
const MAX_EVENTS = 80;
const MAX_CHECKPOINTS = 48;

export function escapeBlessedTags(line: string): string {
  return line.replace(/\{/g, "(").replace(/\}/g, ")");
}

export function formatHeaderStats(data: DashboardData): string {
  const counts =
    `{gray-fg}PRDs{/} ${data.prds.length}  ` +
    `{gray-fg}ch{/} ${data.changes.length}  ` +
    `{gray-fg}iss{/} ${data.issues.length}  ` +
    `{gray-fg}ev{/} ${data.events.length}  ` +
    `{gray-fg}cp{/} ${data.checkpoints.length}  ` +
    `{gray-fg}ho{/} ${data.handoffs.length}`;
  if (!data.warnings.length) return counts;
  return `${counts}\n{yellow-fg}⚠ ${escapeBlessedTags(data.warnings[0].slice(0, 140))}{/yellow-fg}`;
}

export function formatPrdSidebar(prds: PrdCard[], selectedIndex: number): string {
  if (!prds.length) return "{gray-fg}(no PRDs){/}";
  return prds
    .map((p, i) => {
      const base = path.basename(p.path);
      const status = escapeBlessedTags(p.status.slice(0, 18));
      const title = escapeBlessedTags(p.title.slice(0, 44));
      const block = `{cyan-fg}▸{/} ${escapeBlessedTags(base)}\n {gray-fg}${status}{/} ${title}`;
      return i === selectedIndex ? `{inverse} ${block} {/inverse}` : block;
    })
    .join("\n\n");
}

export function formatEventLines(events: DashboardEvent[]): string {
  const slice = events.slice(0, MAX_EVENTS);
  if (!slice.length) return "{gray-fg}(no events){/}";
  return slice
    .map((e) => {
      const bits = [e.time, e.changeId, e.phase, e.node, e.status, e.summary].filter(Boolean).join(" | ");
      return `{gray-fg}•{/} ${escapeBlessedTags(bits.slice(0, 500))}`;
    })
    .join("\n");
}

export function formatCheckpointLines(cps: Checkpoint[]): string {
  const slice = cps.slice(0, MAX_CHECKPOINTS);
  if (!slice.length) return "{gray-fg}(no checkpoints){/}";
  return slice.map((c) => escapeBlessedTags(c.raw)).join("\n");
}

export function truncatePrdBody(body: string): string {
  const lines = body.split(/\r?\n/);
  if (lines.length <= MAX_DETAIL_LINES) return escapeBlessedTags(body);
  return escapeBlessedTags(lines.slice(0, MAX_DETAIL_LINES).join("\n")) + "\n\n{gray-fg}… truncated{/}";
}

export function groupIssuesByLane(issues: BoardIssue[], lanes: BoardLane[]): Map<string, BoardIssue[]> {
  const map = new Map<string, BoardIssue[]>();
  for (const lane of lanes) {
    map.set(lane.id, []);
  }
  const fallback = lanes[0]?.id ?? "draft";
  for (const issue of issues) {
    const id = issue.lane && map.has(issue.lane) ? issue.lane : fallback;
    map.get(id)!.push(issue);
  }
  return map;
}

export function formatLaneColumn(issues: BoardIssue[], maxWidth: number): string {
  if (!issues.length) return "{gray-fg}— empty —{/}";
  const w = Math.max(12, maxWidth);
  const blocks: string[] = [];
  for (const issue of issues) {
    const key = escapeBlessedTags(issue.key.slice(0, w));
    const st = escapeBlessedTags(issue.status.slice(0, w - 2));
    const prog =
      issue.taskStats.total > 0
        ? `${issue.taskStats.completed}/${issue.taskStats.total}`
        : "—";
    const title = escapeBlessedTags(issue.title.slice(0, w));
    blocks.push(`{bold}${key}{/bold}\n {magenta-fg}${st}{/}\n {gray-fg}${prog}{/}\n ${title}\n{gray-fg}${"─".repeat(Math.min(w, 18))}{/}`);
  }
  return blocks.join("\n");
}
