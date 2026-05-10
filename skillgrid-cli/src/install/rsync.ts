import { spawnSync } from "node:child_process";
import { commandOnPath } from "./exec.js";

const RSYNC_EXCLUDES = [".git", ".gitignore", "docs", "configs"];

export function assertRsyncAvailable(): void {
  if (!commandOnPath("rsync")) {
    throw new Error(
      "rsync is not on PATH. Install rsync (macOS: Xcode CLIs; Windows: WSL or cwRsync) or use the legacy ./install.sh from Git Bash.",
    );
  }
}

export interface RsyncHubMode {
  /** When true, destination files not present in source are removed (hub mirror). Matches install.sh `rsync --delete`. */
  delete: boolean;
  /** Path segments excluded from transfer (IDE hub folders). */
  excludes?: string[];
}

/**
 * Hub → project rsync. With `-a`, rsync only transfers **changed** files (size/mtime) when the destination
 * already exists — i.e. incremental “diff” sync, not a full blind copy. With {@link RsyncHubMode.delete},
 * extra files under the destination that are absent from the source are removed.
 */
export function rsyncHub(srcDir: string, dstDir: string, dryRun: boolean, mode: RsyncHubMode): void {
  const modeLabel = mode.delete ? "rsync -av --delete" : "rsync -av (merge, no --delete)";
  if (dryRun) {
    const ex = mode.excludes?.length ? ` excludes=${mode.excludes.join(",")}` : "";
    console.log(`[DRY-RUN] Would ${modeLabel}${ex} ${srcDir}/ -> ${dstDir}/`);
    return;
  }
  assertRsyncAvailable();
  const args = ["-av"];
  if (mode.delete) args.push("--delete");
  for (const e of mode.excludes ?? []) {
    args.push("--exclude", e);
  }
  args.push(`${srcDir.endsWith("/") ? srcDir : `${srcDir}/`}`, `${dstDir.endsWith("/") ? dstDir : `${dstDir}/`}`);
  const r = spawnSync("rsync", args, { stdio: "inherit", encoding: "utf8" });
  if (r.status !== 0) {
    throw new Error(`rsync failed with exit ${r.status}`);
  }
}

/** IDE config dirs from hub: `install.sh`-style excludes + `--delete` mirror. */
export function rsyncDir(srcDir: string, dstDir: string, dryRun: boolean, extraExcludes: string[] = []): void {
  rsyncHub(srcDir, dstDir, dryRun, { delete: true, excludes: [...RSYNC_EXCLUDES, ...extraExcludes] });
}

/**
 * `.skillgrid/` hub → project: **no --delete** (same as `install.sh`) so local `prd/`, `preview/`, etc. survive.
 * Still incremental: only changed hub files are copied when the tree already exists.
 */
export function rsyncPlain(srcDir: string, dstDir: string, dryRun: boolean): void {
  rsyncHub(srcDir, dstDir, dryRun, { delete: false });
}

/** Full mirror with `--delete` (`.agents/`, `.agent/`, skills, `.github/`, `.copilot/`). */
export function rsyncDelete(srcDir: string, dstDir: string, dryRun: boolean): void {
  rsyncHub(srcDir, dstDir, dryRun, { delete: true });
}
