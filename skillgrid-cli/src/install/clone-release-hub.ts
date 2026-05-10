import { existsSync, rmSync } from "node:fs";
import { join } from "node:path";
import { tmpdir } from "node:os";
import { spawnSync } from "node:child_process";
import { commandOnPath } from "./exec.js";
import { logInfo, logWarn } from "./log.js";

/**
 * Base dir for temp-like paths:
 * - **Unix** (macOS, Linux, …): `/tmp` (avoids macOS `TMPDIR` under `/var/folders/…`).
 * - **WSL**: `process.platform` is `linux` → `/tmp`.
 * - **Native win32** only: `os.tmpdir()` — prefer WSL for install (bash, rsync, git).
 */
function hubCloneTempBase(): string {
  if (process.platform === "win32") return tmpdir();
  return "/tmp";
}

/**
 * Canonical hub source for `skillgrid install` (matches public hub branch).
 * Web URL for humans: https://github.com/devopstales/aiskillgrid/tree/release/2
 */
export const AISKILLGRID_UPSTREAM_REPO = "https://github.com/devopstales/aiskillgrid.git";
export const AISKILLGRID_RELEASE_BRANCH = "release/2";

/** Fixed cache directory name under {@link hubCloneTempBase} (reused across runs; refresh with `git fetch` + `git reset`). */
export const RELEASE_HUB_CACHE_DIR = "skillgrid-aiskillgrid-release";

export interface ReleaseHubCache {
  /** Repo root containing `install.sh` and `.git`. */
  hubRoot: string;
}

function hubCacheRoot(): string {
  return join(hubCloneTempBase(), RELEASE_HUB_CACHE_DIR);
}

function cloneFresh(hubRoot: string): void {
  if (existsSync(hubRoot)) {
    rmSync(hubRoot, { recursive: true, force: true });
  }
  logInfo(`Cloning AISkillgrid hub (${AISKILLGRID_RELEASE_BRANCH}) → ${hubRoot} …`);
  const r = spawnSync(
    "git",
    ["clone", "--depth", "1", "-b", AISKILLGRID_RELEASE_BRANCH, AISKILLGRID_UPSTREAM_REPO, hubRoot],
    { stdio: "inherit" },
  );
  if (r.status !== 0) {
    throw new Error(
      `git clone failed (exit ${r.status ?? "unknown"}) — check network and branch ${AISKILLGRID_RELEASE_BRANCH}`,
    );
  }
  if (!existsSync(join(hubRoot, "install.sh"))) {
    throw new Error(`Clone completed but install.sh missing at ${hubRoot}`);
  }
}

function tryGitPullShallow(hubRoot: string): boolean {
  logInfo(`Updating AISkillgrid hub cache at ${hubRoot} (git fetch + reset)…`);
  const fetch = spawnSync(
    "git",
    ["-C", hubRoot, "fetch", "--depth", "1", "origin", AISKILLGRID_RELEASE_BRANCH],
    { stdio: "inherit" },
  );
  if (fetch.status !== 0) return false;
  const remoteRef = `origin/${AISKILLGRID_RELEASE_BRANCH}`;
  const reset = spawnSync("git", ["-C", hubRoot, "reset", "--hard", remoteRef], { stdio: "inherit" });
  if (reset.status !== 0) return false;
  return existsSync(join(hubRoot, "install.sh"));
}

/**
 * Ensures a shared shallow clone of the release hub exists under `/tmp/skillgrid-aiskillgrid-release` (Unix/WSL)
 * or `%TEMP%/skillgrid-aiskillgrid-release` (native Windows). Reuses the directory across runs: existing repo is
 * updated with `git fetch` + `git reset --hard` instead of cloning again.
 */
export function ensureReleaseHubCache(): ReleaseHubCache {
  if (!commandOnPath("git")) {
    throw new Error("git is not on PATH — required to clone or update the AISkillgrid release hub for install");
  }

  const hubRoot = hubCacheRoot();
  const gitDir = join(hubRoot, ".git");

  if (existsSync(gitDir)) {
    if (!tryGitPullShallow(hubRoot)) {
      logWarn("Hub cache update failed — recloning from scratch…");
      cloneFresh(hubRoot);
    }
  } else {
    if (existsSync(hubRoot)) {
      rmSync(hubRoot, { recursive: true, force: true });
    }
    cloneFresh(hubRoot);
  }

  if (!existsSync(join(hubRoot, "install.sh"))) {
    throw new Error(`Hub cache invalid (missing install.sh): ${hubRoot}`);
  }

  return { hubRoot };
}
