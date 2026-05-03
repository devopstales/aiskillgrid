import type { TaskStats } from "../shared/types.js";

export const EMPTY_TASK_STATS: TaskStats = {
  total: 0,
  completed: 0,
  hitl: 0,
  afk: 0,
  blocked: 0
};

export type MarkdownMetadata = {
  title?: string;
  status?: string;
  changeId?: string;
  specPath?: string;
  contextPath?: string;
  owner?: string;
  agent?: string;
  previewUrl?: string;
  externalUrl?: string;
  summary?: string;
};

export function parseMarkdownMetadata(content: string): MarkdownMetadata {
  const metadata: MarkdownMetadata = {};
  const frontMatter = content.match(/^---\n([\s\S]*?)\n---/);
  if (frontMatter) {
    for (const line of frontMatter[1].split(/\r?\n/)) {
      const match = line.match(/^([A-Za-z][A-Za-z0-9_-]*):\s*(.*)$/);
      if (match) {
        assignMetadata(metadata, match[1], stripQuotes(match[2]));
      }
    }
  }

  for (const line of content.split(/\r?\n/).slice(0, 80)) {
    if (!metadata.title) {
      const heading = line.match(/^#\s+(.+)$/);
      if (heading) {
        metadata.title = heading[1].trim();
      }
    }

    const colonField =
      line.match(/^(Status|Change|Change ID|Spec \/ change|OpenSpec change|Session context|Owner|Agent|Preview|External|Summary):\s*(.+)$/i) ??
      line.match(/^\*\*(Status|Change|Change ID|Spec \/ change|OpenSpec change|Session context|Owner|Agent|Preview|External|Summary)\*\*:\s*(.+)$/i) ??
      line.match(/^\*\*(Status|Change|Change ID|Spec \/ change|OpenSpec change|Session context|Owner|Agent|Preview|External|Summary):\*\*\s*(.+)$/i);
    if (colonField) {
      assignMetadata(metadata, colonField[1], colonField[2].trim());
    }
  }

  return metadata;
}

export function parseTaskStats(content: string): TaskStats {
  const stats = { ...EMPTY_TASK_STATS };

  for (const line of content.split(/\r?\n/)) {
    const task = line.match(/^\s*[-*]\s+\[([ xX-])\]\s+(.+)$/);
    if (!task) {
      continue;
    }

    stats.total += 1;
    if (task[1].toLowerCase() === "x") {
      stats.completed += 1;
    }

    const text = task[2].toLowerCase();
    if (text.includes("[hitl]")) {
      stats.hitl += 1;
    }
    if (text.includes("[afk]")) {
      stats.afk += 1;
    }
    if (text.includes("blocked") || text.includes("[blocked]")) {
      stats.blocked += 1;
    }
  }

  return stats;
}

export function firstParagraph(content: string): string | undefined {
  const withoutFrontMatter = content.replace(/^---\n[\s\S]*?\n---/, "");
  const paragraph = withoutFrontMatter
    .split(/\n\s*\n/)
    .map((part) => part.trim())
    .find((part) => isNarrativeParagraph(part));

  return paragraph?.replace(/\s+/g, " ").slice(0, 280);
}

function assignMetadata(metadata: MarkdownMetadata, key: string, value: string): void {
  const normalized = key.toLowerCase().replace(/[\s_-]+/g, "");
  if (normalized === "title") metadata.title = value;
  if (normalized === "status" && !metadata.status) metadata.status = value;
  if ((normalized === "change" || normalized === "changeid") && !metadata.changeId) metadata.changeId = value;
  if (normalized === "spec/change" && metadata.specPath === undefined) metadata.specPath = cleanMarkdownValue(value);
  if (normalized === "openspecchange" && metadata.specPath === undefined) metadata.specPath = cleanMarkdownValue(value);
  if (normalized === "sessioncontext" && !metadata.contextPath) metadata.contextPath = cleanMarkdownValue(value);
  if (normalized === "owner" && !metadata.owner) metadata.owner = value;
  if (normalized === "agent" && !metadata.agent) metadata.agent = value;
  if (normalized === "preview" && !metadata.previewUrl) metadata.previewUrl = value;
  if (normalized === "external" && !metadata.externalUrl) metadata.externalUrl = value;
  if (normalized === "summary" && !metadata.summary) metadata.summary = value;
}

function stripQuotes(value: string): string {
  return value.trim().replace(/^["']|["']$/g, "");
}

function cleanMarkdownValue(value: string): string {
  const cleaned = value.trim().replace(/^`|`$/g, "").replace(/\s+$/g, "");
  if (
    !cleaned ||
    /^_?none_?$/i.test(cleaned) ||
    /not created/i.test(cleaned) ||
    /when added/i.test(cleaned)
  ) {
    return "";
  }
  return cleaned;
}

function isNarrativeParagraph(part: string): boolean {
  if (!part || part.startsWith("#") || part.startsWith("- [") || part.startsWith(">")) {
    return false;
  }

  const lines = part.split(/\r?\n/).map((line) => line.trim()).filter(Boolean);
  if (lines.length === 0) {
    return false;
  }

  return lines.some((line) => !isMetadataLine(line) && !line.startsWith("---"));
}

function isMetadataLine(line: string): boolean {
  return /^\*\*(File|Spec \/ change|OpenSpec change|Session context|Status|Imported from):?\*\*:?\s*/i.test(line)
    || /^(File|Spec \/ change|OpenSpec change|Session context|Status|Imported from):\s*/i.test(line);
}
