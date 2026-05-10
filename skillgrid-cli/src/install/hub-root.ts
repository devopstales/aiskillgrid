import { existsSync, realpathSync } from "node:fs";
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";

function isBunBundledPath(dir: string): boolean {
  return dir.includes("$bunfs") || dir.startsWith("/$/");
}

function pushUnique(out: string[], p: string): void {
  const r = resolve(p);
  if (!out.includes(r)) out.push(r);
}

function pushHubCandidatesFromCliDir(out: string[], cliDir: string): void {
  const d = resolve(cliDir);
  pushUnique(out, resolve(d, "..", ".."));
  pushUnique(out, resolve(d, ".."));
}

/**
 * Resolve AISkillgrid hub root (directory containing `install.sh`) from the CLI entry.
 * Dev: `skillgrid-cli/dist/cli.js` → two levels up to hub.
 * Bun `--compile`: `import.meta.url` is under `/$bunfs` — derive hub from `process.execPath` / `argv[1]`.
 */
export function hubRootFromCliModule(importMetaUrl: string): string {
  const here = dirname(fileURLToPath(importMetaUrl));
  const fallback = resolve(here, "..", "..");

  if (!isBunBundledPath(here)) {
    return fallback;
  }

  const candidates: string[] = [];
  try {
    pushHubCandidatesFromCliDir(candidates, dirname(realpathSync(process.execPath)));
  } catch {
    /* empty */
  }
  const argv1 = process.argv[1];
  if (argv1 && !argv1.startsWith("-")) {
    try {
      pushHubCandidatesFromCliDir(candidates, dirname(realpathSync(argv1)));
    } catch {
      /* empty */
    }
  }

  for (const root of candidates) {
    if (existsSync(resolve(root, "install.sh"))) {
      return root;
    }
  }

  return fallback;
}
