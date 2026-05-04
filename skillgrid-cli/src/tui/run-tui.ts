import { existsSync } from "node:fs";
import path from "node:path";
import { buildDashboardData } from "../dashboard/server/adapters.js";
import { runBlessedTuiLoop } from "./blessed-ui.js";
import { parseTuiArgv, printTuiHelp } from "./parse-tui-argv.js";

export async function runTuiCommand(argv: string[]): Promise<void> {
  const opts = parseTuiArgv(argv);
  if (opts.help) {
    printTuiHelp();
    return;
  }

  const prdDir = path.join(opts.repoRoot, ".skillgrid", "prd");
  if (!existsSync(prdDir)) {
    console.error(`PRD directory not found: ${prdDir}`);
    console.error("Use --repo <path> pointing at a repository root with .skillgrid/prd.");
    process.exit(1);
  }

  const refresh = () =>
    buildDashboardData({
      repoRoot: opts.repoRoot,
      dashboardOrigin: "http://127.0.0.1:0",
      skipToolHealthChecks: true
    });
  runBlessedTuiLoop({ intervalMs: opts.intervalMs, refresh });
}
