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

export function rsyncDir(srcDir: string, dstDir: string, dryRun: boolean, extraExcludes: string[] = []): void {
  if (dryRun) {
    console.log(`[DRY-RUN] Would rsync (with hub excludes) ${srcDir}/ -> ${dstDir}/`);
    return;
  }
  assertRsyncAvailable();
  const args = ["-av"];
  for (const e of [...RSYNC_EXCLUDES, ...extraExcludes]) {
    args.push("--exclude", e);
  }
  args.push(`${srcDir.endsWith("/") ? srcDir : `${srcDir}/`}`, `${dstDir.endsWith("/") ? dstDir : `${dstDir}/`}`);
  const r = spawnSync("rsync", args, { stdio: "inherit", encoding: "utf8" });
  if (r.status !== 0) {
    throw new Error(`rsync failed with exit ${r.status}`);
  }
}

/** Sync without IDE-folder excludes (e.g. .skillgrid, .github). */
export function rsyncPlain(srcDir: string, dstDir: string, dryRun: boolean): void {
  if (dryRun) {
    console.log(`[DRY-RUN] Would rsync ${srcDir}/ -> ${dstDir}/`);
    return;
  }
  assertRsyncAvailable();
  const args = ["-av", `${srcDir.endsWith("/") ? srcDir : `${srcDir}/`}`, `${dstDir.endsWith("/") ? dstDir : `${dstDir}/`}`];
  const r = spawnSync("rsync", args, { stdio: "inherit", encoding: "utf8" });
  if (r.status !== 0) {
    throw new Error(`rsync failed with exit ${r.status}`);
  }
}

export function rsyncDelete(srcDir: string, dstDir: string, dryRun: boolean): void {
  if (dryRun) {
    console.log(`[DRY-RUN] Would rsync --delete ${srcDir}/ -> ${dstDir}/`);
    return;
  }
  assertRsyncAvailable();
  const args = ["-av", "--delete", `${srcDir.endsWith("/") ? srcDir : `${srcDir}/`}`, `${dstDir.endsWith("/") ? dstDir : `${dstDir}/`}`];
  const r = spawnSync("rsync", args, { stdio: "inherit", encoding: "utf8" });
  if (r.status !== 0) {
    throw new Error(`rsync failed with exit ${r.status}`);
  }
}
