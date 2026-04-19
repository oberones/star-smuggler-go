package services

import (
	"fmt"
	"sort"

	"github.com/oberones/star-smuggler-go/internal/domain"
)

type UpgradePurchaseResult struct {
	Succeeded bool
	Message   string
}

type UpgradeService struct{}

func SuccessfulUpgradePurchase(message string) UpgradePurchaseResult {
	return UpgradePurchaseResult{
		Succeeded: true,
		Message:   message,
	}
}

func FailedUpgradePurchase(message string) UpgradePurchaseResult {
	return UpgradePurchaseResult{
		Succeeded: false,
		Message:   message,
	}
}

func (s UpgradeService) AvailableUpgrades(run domain.RunState, snapshot domain.DataSnapshot, factions FactionService) []domain.ShipUpgradeDefinition {
	available := make([]domain.ShipUpgradeDefinition, 0)
	for _, upgrade := range snapshot.Upgrades {
		if run.Progression.HasUpgrade(upgrade.ID) {
			continue
		}
		if upgrade.RequiredFactionID != "" &&
			!factions.MeetsMinimumStanding(run, snapshot, upgrade.RequiredFactionID, upgrade.MinimumStanding) {
			continue
		}
		available = append(available, upgrade)
	}

	sort.SliceStable(available, func(i, j int) bool {
		if available[i].CostCredits == available[j].CostCredits {
			return available[i].Name < available[j].Name
		}
		return available[i].CostCredits < available[j].CostCredits
	})

	return available
}

func (s UpgradeService) PurchaseUpgrade(run *domain.RunState, snapshot domain.DataSnapshot, factions FactionService, upgradeID string) UpgradePurchaseResult {
	upgrade, ok := snapshot.UpgradesByID[upgradeID]
	if !ok {
		return FailedUpgradePurchase(fmt.Sprintf("Upgrade %q does not exist.", upgradeID))
	}
	if run.Progression.HasUpgrade(upgrade.ID) {
		return FailedUpgradePurchase(fmt.Sprintf("%s is already installed.", upgrade.Name))
	}
	if upgrade.RequiredFactionID != "" &&
		!factions.MeetsMinimumStanding(*run, snapshot, upgrade.RequiredFactionID, upgrade.MinimumStanding) {
		return FailedUpgradePurchase(fmt.Sprintf("%s requires %s standing with %s.", upgrade.Name, upgrade.MinimumStanding, upgrade.RequiredFactionID))
	}
	if run.Player.Credits < upgrade.CostCredits {
		return FailedUpgradePurchase(fmt.Sprintf("You need %d credits to install %s.", upgrade.CostCredits, upgrade.Name))
	}

	run.Player.Credits -= upgrade.CostCredits
	run.Progression.PurchasedUpgradeIDs = append(run.Progression.PurchasedUpgradeIDs, upgrade.ID)
	sort.Strings(run.Progression.PurchasedUpgradeIDs)
	if upgrade.Specialization != "" {
		if run.Progression.SpecializationFlags == nil {
			run.Progression.SpecializationFlags = make(map[string]bool)
		}
		run.Progression.SpecializationFlags[string(upgrade.Specialization)] = true
	}

	for _, effect := range upgrade.Effects {
		if effect.Type == domain.UpgradeEffectCargoLimitBonus {
			run.Player.CargoLimit += effect.Value
		}
	}

	return SuccessfulUpgradePurchase(fmt.Sprintf("Installed %s for %d credits.", upgrade.Name, upgrade.CostCredits))
}

func (s UpgradeService) AdjustTravelCost(run domain.RunState, baseCost int, snapshot domain.DataSnapshot) int {
	discountPercent := 0
	for _, upgradeID := range run.Progression.PurchasedUpgradeIDs {
		upgrade, ok := snapshot.UpgradesByID[upgradeID]
		if !ok {
			continue
		}
		for _, effect := range upgrade.Effects {
			if effect.Type == domain.UpgradeEffectTravelCostDiscountPct {
				discountPercent += effect.Value
			}
		}
	}

	if discountPercent <= 0 {
		return baseCost
	}
	if discountPercent > 90 {
		discountPercent = 90
	}

	adjustedCost := baseCost * (100 - discountPercent) / 100
	if baseCost > 0 && adjustedCost < 1 {
		return 1
	}
	return adjustedCost
}

func (s UpgradeService) AdjustMissionReward(run domain.RunState, baseReward int, snapshot domain.DataSnapshot) int {
	bonusPercent := 0
	for _, upgradeID := range run.Progression.PurchasedUpgradeIDs {
		upgrade, ok := snapshot.UpgradesByID[upgradeID]
		if !ok {
			continue
		}
		for _, effect := range upgrade.Effects {
			if effect.Type == domain.UpgradeEffectMissionRewardBonusPct {
				bonusPercent += effect.Value
			}
		}
	}

	if bonusPercent <= 0 {
		return baseReward
	}

	return baseReward + (baseReward*bonusPercent)/100
}
