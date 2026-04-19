#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$repo_root"

go test ./tests/integration -run 'TestEventService|TestTravelCommandsResolveArrivalRefreshJumpCountAndEvent|TestTravelAnimationDurationScalesWithZoneDifference'
bash tests/smoke/run_headless_smoke.sh
STARSMUGGLER_GO_RUNTIME=1 bash tests/smoke/run_headless_smoke.sh
