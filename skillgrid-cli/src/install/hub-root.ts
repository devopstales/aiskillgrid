import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";

/** Resolve AISkillgrid hub root from `cli` entry (skillgrid-cli/src/cli.ts or skillgrid-cli/dist/cli.js). */
export function hubRootFromCliModule(importMetaUrl: string): string {
  const here = dirname(fileURLToPath(importMetaUrl));
  return resolve(here, "..", "..");
}
