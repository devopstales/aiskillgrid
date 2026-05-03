import { spawnSync } from "node:child_process";
import type { IdeId, OptionalToolId } from "./types.js";
import { commandOnPath, shCheck } from "./exec.js";
import { toolIsSelected } from "./optional-tools-helpers.js";

export interface CoreDep {
  name: string;
  brew: string;
  pip: string;
  npm: string;
  checkCmd: string;
}

export const CORE_DEPENDENCIES: CoreDep[] = [
  { name: "rsync", brew: "rsync", pip: "", npm: "", checkCmd: "rsync --version" },
  { name: "jq", brew: "jq", pip: "", npm: "", checkCmd: "jq --version" },
  { name: "node", brew: "node", pip: "", npm: "", checkCmd: "node --version" },
  { name: "python3", brew: "python@3.12", pip: "", npm: "", checkCmd: "python3 --version" },
  { name: "pip3", brew: "", pip: "", npm: "", checkCmd: "pip3 --version" },
  { name: "npx", brew: "", pip: "", npm: "", checkCmd: "npx --version" },
];

export interface IdeDep {
  name: string;
  brew: string;
  pip: string;
  npm: string;
  checkCmd: string;
  ideFlag: string;
}

export const IDE_DEPENDENCIES: IdeDep[] = [
  { name: "opencode", brew: "", pip: "", npm: "", checkCmd: "opencode --version", ideFlag: "opencode" },
  { name: "openspec", brew: "openspec", pip: "", npm: "", checkCmd: "openspec --version", ideFlag: "openspec" },
  { name: "kilo", brew: "Kilo-Org/homebrew-tap/kilo", pip: "", npm: "", checkCmd: "kilo --version", ideFlag: "kilo" },
  { name: "semgrep", brew: "semgrep", pip: "", npm: "", checkCmd: "semgrep --version", ideFlag: "semgrep" },
  { name: "trivy", brew: "trivy", pip: "", npm: "", checkCmd: "trivy --version", ideFlag: "trivy" },
  { name: "trivy-mcp", brew: "trivy", pip: "", npm: "", checkCmd: "trivy plugin list", ideFlag: "trivy-mcp" },
];

function ideSelected(selectedIdes: IdeId[], ideFlag: string, allIdes: boolean, selectedTools: OptionalToolId[]): boolean {
  if (allIdes) return true;
  if (ideFlag === "openspec" && toolIsSelected(selectedTools, "openspec")) return true;
  return (selectedIdes as string[]).includes(ideFlag);
}

export function countMissingDeps(
  selectedIdes: IdeId[],
  allIdes: boolean,
  selectedTools: OptionalToolId[],
): { brew: string[]; pip: string[]; npm: string[] } {
  const brew: string[] = [];
  const pip: string[] = [];
  const npm: string[] = [];

  for (const dep of CORE_DEPENDENCIES) {
    if (!dep.checkCmd) continue;
    if (!shCheck(dep.checkCmd)) {
      if (dep.brew) brew.push(dep.brew);
      else if (dep.pip) pip.push(dep.pip);
      else if (dep.npm) npm.push(dep.npm);
    }
  }

  for (const dep of IDE_DEPENDENCIES) {
    if (!dep.checkCmd) continue;
    if (!ideSelected(selectedIdes, dep.ideFlag, allIdes, selectedTools)) continue;
    if (!shCheck(dep.checkCmd)) {
      if (dep.brew) brew.push(dep.brew);
      else if (dep.pip) pip.push(dep.pip);
      else if (dep.npm) npm.push(dep.npm);
    }
  }

  return { brew, pip, npm };
}

export function installDependencyPackages(brewPkgs: string[], pipPkgs: string[], npmPkgs: string[], dryRun: boolean): void {
  if (brewPkgs.length > 0) {
    console.log(`Installing with Homebrew: ${brewPkgs.join(" ")}`);
    if (dryRun) console.log(`[DRY-RUN] brew install ${brewPkgs.join(" ")}`);
    else if (commandOnPath("brew")) {
      spawnSync("brew", ["install", ...brewPkgs], { stdio: "inherit" });
    } else {
      console.log("  ⚠ Homebrew not found. Skipping brew packages.");
      console.log("  Install Homebrew: https://brew.sh");
    }
    console.log("");
  }

  if (pipPkgs.length > 0) {
    console.log(`Installing with pip: ${pipPkgs.join(" ")}`);
    if (dryRun) console.log(`[DRY-RUN] pip3 install ${pipPkgs.join(" ")}`);
    else if (commandOnPath("pip3")) {
      spawnSync("pip3", ["install", "--user", ...pipPkgs], { stdio: "inherit" });
    } else console.log("  ⚠ pip3 not found. Skipping pip packages.");
    console.log("");
  }

  if (npmPkgs.length > 0) {
    console.log(`Installing with npm: ${npmPkgs.join(" ")}`);
    if (dryRun) console.log(`[DRY-RUN] npm install -g ${npmPkgs.join(" ")}`);
    else if (commandOnPath("npm")) {
      spawnSync("npm", ["install", "-g", ...npmPkgs], { stdio: "inherit" });
    } else console.log("  ⚠ npm not found. Skipping npm packages.");
    console.log("");
  }

  if (brewPkgs.length === 0 && pipPkgs.length === 0 && npmPkgs.length === 0) {
    console.log("All dependencies are already installed!");
    console.log("");
  }
}

export function showDependencies(selectedIdes: IdeId[], allIdes: boolean, selectedTools: OptionalToolId[]): void {
  console.log("=== Dependency Check ===");
  console.log("");

  const installed: string[] = [];
  const missing: { name: string; brew: string; pip: string; npm: string }[] = [];

  console.log("Core Dependencies:");
  for (const dep of CORE_DEPENDENCIES) {
    if (!dep.checkCmd) continue;
    if (shCheck(dep.checkCmd)) {
      installed.push(dep.name);
      console.log(`  ✓ ${dep.name}`);
    } else {
      missing.push({ name: dep.name, brew: dep.brew, pip: dep.pip, npm: dep.npm });
      console.log(`  ✗ ${dep.name}`);
    }
  }

  console.log("");
  console.log("IDE-Specific Dependencies:");
  for (const dep of IDE_DEPENDENCIES) {
    if (!dep.checkCmd) continue;
    if (!ideSelected(selectedIdes, dep.ideFlag, allIdes, selectedTools)) {
      console.log(`  - ${dep.name} (skipped, not selected)`);
      continue;
    }
    if (shCheck(dep.checkCmd)) {
      installed.push(dep.name);
      console.log(`  ✓ ${dep.name}`);
    } else {
      missing.push({ name: dep.name, brew: dep.brew, pip: dep.pip, npm: dep.npm });
      console.log(`  ✗ ${dep.name}`);
    }
  }

  console.log("");
  console.log(`Installed: ${installed.length}`);
  console.log(`Missing: ${missing.length}`);
  console.log("");

  if (missing.length > 0) {
    console.log("Missing dependencies:");
    for (const m of missing) {
      let installCmd = "";
      if (m.brew) installCmd = `brew install ${m.brew}`;
      else if (m.npm) installCmd = `npm install -g ${m.npm}`;
      else if (m.pip) installCmd = `pip3 install ${m.pip}`;
      if (installCmd) console.log(`  - ${m.name} → ${installCmd}`);
      else console.log(`  - ${m.name} → (manual install required)`);
    }
    console.log("");
  }
}
