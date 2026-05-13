#!/usr/bin/env node

const fs = require("fs");
const path = require("path");
const { execFileSync } = require("child_process");

const ext = process.platform === "win32" ? ".exe" : "";
const binaryPath = path.join(__dirname, "..", "bin", `md2wechat-new${ext}`);

if (!fs.existsSync(binaryPath)) {
  console.error(
    "md2wechat-new binary is missing. Reinstall with `npm install -g @lexiaowen/md2wechat-new`."
  );
  process.exit(1);
}

try {
  execFileSync(binaryPath, process.argv.slice(2), { stdio: "inherit" });
} catch (error) {
  if (typeof error.status === "number") {
    process.exit(error.status);
  }

  console.error(`Failed to launch md2wechat-new: ${error.message}`);
  process.exit(1);
}
