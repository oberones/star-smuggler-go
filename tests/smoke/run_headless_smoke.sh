#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
GODOT_BIN="${GODOT_BIN:-/Applications/Godot .NET.app/Contents/MacOS/Godot}"
LOG_FILE="${TMPDIR:-/tmp}/starsmuggler-headless-$$-${RANDOM:-0}.log"

"${GODOT_BIN}" --headless --path "${ROOT_DIR}" --log-file "${LOG_FILE}" --quit-after 1
