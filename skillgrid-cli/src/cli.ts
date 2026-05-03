import { hubRootFromCliModule } from "./install/hub-root.js";
import { parseInstallArgv } from "./install/parse-install-argv.js";
import { runInstallCli } from "./install/run-install.js";

function printTopHelp() {
  console.log(`skillgrid — Skillgrid hub CLI

Usage:
  skillgrid install [OPTIONS]   Native installer (parity with ./install.sh)
  skillgrid --help              Show this message

Alternate installer (bash):
  ./install.sh [OPTIONS]

Build this CLI (requires Bun for the binary):
  cd skillgrid-cli && npm install && npm run build
  ./bin/skillgrid install --path /your/project -A -y
`);
}

async function main() {
  const argv = process.argv.slice(2);
  if (argv.length === 0 || argv[0] === "-h" || argv[0] === "--help") {
    printTopHelp();
    process.exit(0);
  }

  if (argv[0] === "install") {
    const hubRoot = hubRootFromCliModule(import.meta.url);
    let parsed;
    try {
      parsed = parseInstallArgv(argv.slice(1));
    } catch (e) {
      console.error((e as Error).message);
      process.exit(1);
    }
    const code = await runInstallCli(hubRoot, parsed);
    process.exit(code);
  }

  console.error(`Unknown command: ${argv[0]}`);
  printTopHelp();
  process.exit(1);
}

main().catch((e) => {
  console.error(e);
  process.exit(1);
});
