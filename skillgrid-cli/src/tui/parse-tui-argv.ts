import path from "node:path";

export type TuiCliOptions = {
  repoRoot: string;
  help: boolean;
  intervalMs: number;
};

export function parseTuiArgv(argv: string[]): TuiCliOptions {
  let repoRoot = process.cwd();
  let help = false;
  let intervalMs = 2000;

  for (let i = 0; i < argv.length; i += 1) {
    const a = argv[i];
    const n = argv[i + 1];
    if (a === "--help" || a === "-h") help = true;
    else if (a === "--repo" && n) {
      repoRoot = path.resolve(n);
      i += 1;
    } else if (a === "--prd-dir" || a === "--prd_dir") {
      if (n) i += 1;
      throw new Error(
        "The --prd-dir option is not supported. Pass the repository root with --repo <path> (PRDs are read from <repo>/.skillgrid/prd)."
      );
    } else if (a === "--interval" && n) {
      intervalMs = Math.max(500, Number.parseInt(n, 10) || 2000);
      i += 1;
    }
  }

  return { repoRoot, help, intervalMs };
}

export function printTuiHelp(): void {
  console.log(`skillgrid tui — read-only dashboard in the terminal

PRD list, markdown body, event log, checkpoints, and a lane board aligned with the web UI.

Usage:
  skillgrid tui [options]

Options:
  --repo <path>       Repository root. Default: current directory.
  --interval <ms>     Refresh interval in milliseconds (min 500). Default: 2000.
  --help, -h          Show this help.

Keys (in the TUI):
  [  ]                Previous / next PRD in the list
  p                   Show selected PRD markdown in the main pane
  l                   Show event log in the main pane
  c                   Show checkpoints in the main pane
  Tab                 Focus next pane (scroll with arrows / PgUp / PgDn where supported)
  r                   Refresh now
  q, Ctrl+C           Quit
`);
}
