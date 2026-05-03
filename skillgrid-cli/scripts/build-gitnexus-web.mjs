#!/usr/bin/env node
/**
 * Vendors a production build of GitNexus Web from the upstream monorepo into
 * dist/dashboard/gitnexus so the Skillgrid dashboard can serve it at /gitnexus/.
 *
 * Env:
 *   SKIP_GITNEXUS_WEB=1     — skip (no network / no Node 20+).
 *   GITNEXUS_REF=main       — git ref to fetch (shallow clone).
 *   GITNEXUS_FORCE_CI=1     — always run npm ci in gitnexus-web (default: ci if no node_modules).
 *   GITNEXUS_VITE_VERBOSE=1 — show full Vite logs (default: --logLevel error to hide upstream chunk / Rolldown noise).
 */
import { spawnSync } from "node:child_process";
import { cpSync, existsSync, mkdirSync, readFileSync, rmSync, writeFileSync } from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const ROOT = path.join(path.dirname(fileURLToPath(import.meta.url)), "..");
const CLONE_DIR = path.join(ROOT, "tmp", "gitnexus-src");
const WEB_DIR = path.join(CLONE_DIR, "gitnexus-web");
const VITE_CONFIG = path.join(WEB_DIR, "vite.config.ts");
const OUT_DIR = path.join(ROOT, "dist", "dashboard", "gitnexus");
const REF = process.env.GITNEXUS_REF || "main";

function run(cmd, args, opts) {
  const r = spawnSync(cmd, args, { stdio: "inherit", shell: false, ...opts });
  if (r.status !== 0) process.exit(r.status ?? 1);
}

function ensureBasePatch() {
  let s = readFileSync(VITE_CONFIG, "utf8");
  if (s.includes("base: '/gitnexus/'") || s.includes('base: "/gitnexus/"')) return;
  const needle = "export default defineConfig({";
  if (!s.includes(needle)) {
    console.error("build-gitnexus-web: unexpected vite.config.ts layout; cannot inject base.");
    process.exit(1);
  }
  s = s.replace(needle, `${needle}\n  base: '/gitnexus/',`);
  writeFileSync(VITE_CONFIG, s);
}

/** Upstream bundle is large; silences Vite's default 500 kB chunk warning without forking their app. */
function ensureChunkSizePatch() {
  let s = readFileSync(VITE_CONFIG, "utf8");
  if (s.includes("chunkSizeWarningLimit")) return;
  if (s.includes("base: '/gitnexus/',")) {
    s = s.replace("base: '/gitnexus/',", "base: '/gitnexus/',\n  build: { chunkSizeWarningLimit: 4096 },");
  } else {
    const needle = "export default defineConfig({";
    if (!s.includes(needle)) {
      console.error("build-gitnexus-web: unexpected vite.config.ts layout; cannot inject chunkSizeWarningLimit.");
      process.exit(1);
    }
    s = s.replace(needle, `${needle}\n  build: { chunkSizeWarningLimit: 4096 },`);
  }
  writeFileSync(VITE_CONFIG, s);
}

function main() {
  if (process.env.SKIP_GITNEXUS_WEB === "1") {
    console.warn("build-gitnexus-web: SKIP_GITNEXUS_WEB=1 — skipping GitNexus web bundle.");
    return;
  }

  const major = Number.parseInt(process.versions.node.split(".")[0] ?? "0", 10);
  if (major < 20) {
    console.error(
      "build-gitnexus-web: GitNexus web requires Node ^20.19 or >=22.12 (see upstream package.json). Upgrade Node or set SKIP_GITNEXUS_WEB=1."
    );
    process.exit(1);
  }

  mkdirSync(path.dirname(CLONE_DIR), { recursive: true });

  if (!existsSync(path.join(CLONE_DIR, ".git"))) {
    rmSync(CLONE_DIR, { recursive: true, force: true });
    run("git", ["clone", "--depth", "1", "--branch", REF, "https://github.com/abhigyanpatwari/GitNexus.git", CLONE_DIR], {
      cwd: ROOT
    });
  } else {
    console.log(
      "build-gitnexus-web: reusing tmp/gitnexus-src — delete skillgrid-cli/tmp/gitnexus-src to pull a fresh upstream checkout."
    );
  }

  ensureBasePatch();
  ensureChunkSizePatch();

  const nodeModules = path.join(WEB_DIR, "node_modules");
  if (!existsSync(nodeModules) || process.env.GITNEXUS_FORCE_CI === "1") {
    run("npm", ["ci"], { cwd: WEB_DIR, env: { ...process.env } });
  }

  /* Upstream `npm run build` runs `tsc -b` which does not resolve the Vite alias for gitnexus-shared; Vite build alone emits the static bundle we need. */
  const viteArgs = ["vite", "build"];
  if (process.env.GITNEXUS_VITE_VERBOSE !== "1") {
    viteArgs.push("--logLevel", "error");
  }
  run("npx", viteArgs, { cwd: WEB_DIR, env: { ...process.env } });

  const built = path.join(WEB_DIR, "dist");
  if (!existsSync(path.join(built, "index.html"))) {
    console.error("build-gitnexus-web: gitnexus-web build produced no dist/index.html");
    process.exit(1);
  }

  mkdirSync(path.dirname(OUT_DIR), { recursive: true });
  rmSync(OUT_DIR, { recursive: true, force: true });
  cpSync(built, OUT_DIR, { recursive: true });
  console.log(`build-gitnexus-web: copied GitNexus web to ${OUT_DIR}`);
}

main();
