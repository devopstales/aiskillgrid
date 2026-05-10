const TTY = process.stdout.isTTY;

const codes = TTY
  ? { red: "\u001b[0;31m", green: "\u001b[0;32m", yellow: "\u001b[1;33m", blue: "\u001b[0;34m", cyan: "\u001b[0;36m", nc: "\u001b[0m" }
  : { red: "", green: "", yellow: "", blue: "", cyan: "", nc: "" };

export function logInfo(msg: string) {
  console.error(`${codes.blue}[INFO]${codes.nc} ${msg}`);
}

export function logSuccess(msg: string) {
  console.error(`${codes.green}[OK]${codes.nc} ${msg}`);
}

export function logWarn(msg: string) {
  console.error(`${codes.yellow}[WARN]${codes.nc} ${msg}`);
}

export function logError(msg: string) {
  console.error(`${codes.red}[ERROR]${codes.nc} ${msg}`);
}
