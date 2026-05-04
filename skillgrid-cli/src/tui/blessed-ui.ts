import { createRequire } from "node:module";
import { existsSync } from "node:fs";
import path from "node:path";
import { fileURLToPath, pathToFileURL } from "node:url";
import type * as BlessedExports from "blessed";
import type { DashboardData } from "../dashboard/shared/types.js";
import {
  escapeBlessedTags,
  formatCheckpointLines,
  formatEventLines,
  formatHeaderStats,
  formatLaneColumn,
  formatPrdSidebar,
  groupIssuesByLane,
  truncatePrdBody
} from "./tui-view-format.js";

function loadBlessed(): typeof BlessedExports {
  const seeds = new Set<string>();
  try {
    seeds.add(path.dirname(fileURLToPath(import.meta.url)));
  } catch {
    /* empty */
  }
  seeds.add(path.dirname(process.execPath));

  for (const start of seeds) {
    let dir = path.resolve(start);
    for (let i = 0; i < 16; i++) {
      const blessedPkg = path.join(dir, "node_modules", "blessed", "package.json");
      if (existsSync(blessedPkg)) {
        const req = createRequire(pathToFileURL(blessedPkg).href);
        return req(".") as typeof BlessedExports;
      }
      const parent = path.dirname(dir);
      if (parent === dir) break;
      dir = parent;
    }
  }
  throw new Error(
    "Cannot load blessed: install skillgrid-cli dependencies (npm ci in skillgrid-cli) so node_modules/blessed exists."
  );
}

const blessed = loadBlessed();

const HEADER_H = 3;
const FOOTER_H = 2;
const KANBAN_H = 11;

type DetailMode = "prd" | "events" | "checkpoints";

export function runBlessedTuiLoop(options: {
  intervalMs: number;
  refresh: () => Promise<DashboardData>;
}): void {
  const screen = blessed.screen({
    smartCSR: true,
    title: "skillgrid tui",
    fullUnicode: true
  });

  const headerLeft = blessed.box({
    parent: screen,
    top: 0,
    left: 0,
    width: "58%",
    height: HEADER_H,
    tags: true,
    border: { type: "line" },
    style: { border: { fg: "cyan" }, fg: "white" }
  });

  const headerStats = blessed.box({
    parent: screen,
    top: 0,
    left: "58%",
    width: "42%",
    height: HEADER_H,
    tags: true,
    align: "right",
    valign: "middle",
    border: { type: "line" },
    style: { border: { fg: "cyan" }, fg: "white" }
  });

  const prdList = blessed.box({
    parent: screen,
    top: HEADER_H,
    left: 0,
    width: "26%",
    bottom: FOOTER_H + KANBAN_H,
    label: " PRDs ([ ] cycle) ",
    tags: true,
    keys: true,
    mouse: true,
    scrollable: true,
    alwaysScroll: true,
    scrollbar: { ch: " ", track: { bg: "gray" }, style: { inverse: true } },
    border: { type: "line" },
    style: { border: { fg: "blue" }, fg: "white" }
  });

  const detail = blessed.box({
    parent: screen,
    top: HEADER_H,
    left: "26%",
    right: 0,
    bottom: FOOTER_H + KANBAN_H,
    label: " PRD ",
    tags: true,
    keys: true,
    mouse: true,
    scrollable: true,
    alwaysScroll: true,
    scrollbar: { ch: " ", track: { bg: "gray" }, style: { inverse: true } },
    border: { type: "line" },
    style: { border: { fg: "green" }, fg: "white" }
  });

  const kanban = blessed.box({
    parent: screen,
    bottom: FOOTER_H,
    left: 0,
    right: 0,
    height: KANBAN_H,
    label: " Board (same lanes as web dashboard) ",
    tags: true,
    border: { type: "line" },
    style: { border: { fg: "magenta" }, fg: "white" }
  });

  const footer = blessed.box({
    parent: screen,
    bottom: 0,
    left: 0,
    width: "100%",
    height: FOOTER_H,
    tags: true,
    border: { type: "line" },
    style: { border: { fg: "gray" }, fg: "gray" }
  });

  const laneBoxes: BlessedExports.Widgets.BoxElement[] = [];
  let laneCount = 0;

  function layoutLaneBoxes(lanes: { id: string; label: string }[]): void {
    for (const b of laneBoxes) {
      b.detach();
    }
    laneBoxes.length = 0;
    laneCount = Math.max(1, lanes.length);
    const pct = 100 / laneCount;
    lanes.forEach((lane, i) => {
      const isLast = i === laneCount - 1;
      const box = blessed.box({
        parent: kanban,
        top: 0,
        left: `${(pct * i).toFixed(2)}%`,
        width: isLast ? undefined : `${pct.toFixed(2)}%`,
        right: isLast ? 0 : undefined,
        height: "100%-2",
        label: ` ${lane.label} `,
        tags: true,
        keys: true,
        mouse: true,
        scrollable: true,
        alwaysScroll: true,
        scrollbar: { ch: " ", track: { bg: "gray" }, style: { inverse: true } },
        border: { type: "line" },
        style: { border: { fg: "yellow" }, fg: "white" }
      });
      laneBoxes.push(box);
    });
  }

  const focusables: BlessedExports.Widgets.BoxElement[] = [prdList, detail, ...laneBoxes];
  let focusIndex = 0;

  let lastData: DashboardData | null = null;
  let prdIndex = 0;
  let detailMode: DetailMode = "prd";

  function laneInnerWidth(): number {
    const w = typeof screen.width === "number" ? screen.width : Number(screen.width) || 120;
    const inner = Math.max(8, w - 2);
    return Math.floor(inner / laneCount) - 2;
  }

  function applyData(data: DashboardData): void {
    lastData = data;
    if (data.prds.length) {
      prdIndex = Math.min(prdIndex, data.prds.length - 1);
      prdIndex = Math.max(0, prdIndex);
    } else {
      prdIndex = 0;
    }

    headerLeft.setContent(
      ` {bold}skillgrid tui{/bold} {gray-fg}Phase 2{/gray-fg}\n` +
        ` {cyan-fg}${escapeBlessedTags(data.repoName)}{/cyan-fg}\n` +
        ` {gray-fg}updated ${escapeBlessedTags(data.generatedAt)}{/}`
    );
    headerStats.setContent(` ${formatHeaderStats(data)} `);

    prdList.setContent(formatPrdSidebar(data.prds, prdIndex));

    const prd = data.prds[prdIndex];
    if (detailMode === "events") {
      detail.setLabel(" Event log (l) ");
      detail.setContent(formatEventLines(data.events));
    } else if (detailMode === "checkpoints") {
      detail.setLabel(" Checkpoints (c) ");
      detail.setContent(formatCheckpointLines(data.checkpoints));
    } else {
      detail.setLabel(` PRD (${prdIndex + 1}/${Math.max(1, data.prds.length)}) — p · l · c `);
      detail.setContent(prd ? truncatePrdBody(prd.body ?? "") : "{gray-fg}(no PRD selected){/}");
    }

    if (laneBoxes.length !== data.lanes.length) {
      layoutLaneBoxes(data.lanes);
      focusables.length = 0;
      focusables.push(prdList, detail, ...laneBoxes);
      focusIndex = Math.min(focusIndex, focusables.length - 1);
    }

    const byLane = groupIssuesByLane(data.issues, data.lanes);
    const colW = laneInnerWidth();
    data.lanes.forEach((lane, i) => {
      const box = laneBoxes[i];
      if (box) {
        box.setLabel(` ${lane.label} (${(byLane.get(lane.id) ?? []).length}) `);
        box.setContent(formatLaneColumn(byLane.get(lane.id) ?? [], colW));
      }
    });

    footer.setContent(
      " {gray-fg}[ ]{/} prev/next PRD   {cyan-fg}p{/} PRD text   {cyan-fg}l{/} event log   {cyan-fg}c{/} checkpoints   " +
        "{cyan-fg}Tab{/} focus panes   {cyan-fg}r{/} refresh   {cyan-fg}q{/} quit "
    );
    screen.render();
  }

  async function paint(): Promise<void> {
    const snap = await options.refresh();
    applyData(snap);
  }

  const timer = setInterval(() => {
    paint().catch((e) => {
      footer.setContent(` {red-fg}${escapeBlessedTags(String((e as Error).message))}{/} `);
      screen.render();
    });
  }, options.intervalMs);

  function focusNext(): void {
    focusIndex = (focusIndex + 1) % focusables.length;
    focusables[focusIndex].focus();
    screen.render();
  }

  function bumpPrd(delta: number): void {
    if (!lastData?.prds.length) return;
    const n = lastData.prds.length;
    prdIndex = (prdIndex + delta + n) % n;
    if (lastData) applyData(lastData);
  }

  screen.key(["q", "C-c"], () => {
    clearInterval(timer);
    process.exit(0);
  });

  screen.key(["r"], () => {
    paint().catch(() => {});
  });

  screen.key(["tab"], () => {
    focusNext();
  });

  screen.key(["p"], () => {
    detailMode = "prd";
    if (lastData) applyData(lastData);
  });

  screen.key(["l"], () => {
    detailMode = "events";
    if (lastData) applyData(lastData);
  });

  screen.key(["c"], () => {
    detailMode = "checkpoints";
    if (lastData) applyData(lastData);
  });

  screen.key(["]"], () => {
    bumpPrd(1);
  });

  screen.key(["["], () => {
    bumpPrd(-1);
  });

  screen.on("resize", () => {
    if (lastData) applyData(lastData);
    else screen.render();
  });

  paint()
    .then(() => {
      focusables[0].focus();
      screen.render();
    })
    .catch((e) => {
      clearInterval(timer);
      console.error((e as Error).message);
      process.exit(1);
    });
}
