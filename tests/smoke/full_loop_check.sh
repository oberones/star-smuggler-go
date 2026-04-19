#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$repo_root"

go test ./tests/golden
go test ./tests/integration -run 'TestAppMVPRouteLoopCanStartTravelAndResolveBackToPort|TestAppMVPRecoveryAndContinuePersistRunState|TestStoryCommandsIntegrateMissionAndStoryAcrossTradeAndTravel|TestUpgradeServicePurchasesCargoUpgradeAndEnforcesRules|TestAppPurchaseUpgradeAutosavesAndKeepsTradeRoute|TestSpeedSpecializationReducesTravelCosts|TestInfluenceSpecializationBoostsMissionRewards|TestAudioBridgeMatchesMonoGameRouteMusicIntent|TestResourceCacheCachesTextureMusicAndSfxLookups'
bash tests/smoke/travel_flow_check.sh
