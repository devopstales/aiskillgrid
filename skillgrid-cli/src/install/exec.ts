import { spawnSync } from "node:child_process";

export function shCheck(cmd: string): boolean {
  if (process.platform === "win32") {
    const r = spawnSync("cmd.exe", ["/d", "/s", "/c", cmd], { stdio: "ignore" });
    return r.status === 0;
  }
  const r = spawnSync("sh", ["-c", cmd], { stdio: "ignore" });
  return r.status === 0;
}

export function commandOnPath(name: string): boolean {
  const which = process.platform === "win32" ? `where ${name} 2>nul` : `command -v ${name}`;
  return shCheck(which);
}
