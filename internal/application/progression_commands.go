package application

import (
	"github.com/oberones/star-smuggler-go/internal/domain"
	"github.com/oberones/star-smuggler-go/internal/services"
)

type ProgressionCommands struct {
	Data     domain.DataSnapshot
	Factions services.FactionService
	Upgrades services.UpgradeService
}

func NewProgressionCommands(data domain.DataSnapshot) ProgressionCommands {
	return ProgressionCommands{
		Data: data,
	}
}

func (c ProgressionCommands) AvailableUpgrades(run domain.RunState) []domain.ShipUpgradeDefinition {
	return c.Upgrades.AvailableUpgrades(run, c.Data, c.Factions)
}

func (c ProgressionCommands) PurchaseUpgrade(run *domain.RunState, upgradeID string) services.UpgradePurchaseResult {
	return c.Upgrades.PurchaseUpgrade(run, c.Data, c.Factions, upgradeID)
}
