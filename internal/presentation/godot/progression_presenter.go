package godot

import (
	"fmt"
	"sort"
	"strings"

	"github.com/oberones/star-smuggler-go/internal/domain"
	"github.com/oberones/star-smuggler-go/internal/services"
)

type ProgressionViewModel struct {
	OwnedUpgradeNotices     []string
	AvailableUpgradeNotices []string
	SpecializationNotices   []string
}

type ProgressionPresenter struct {
	Data     domain.DataSnapshot
	Factions services.FactionService
	Upgrades services.UpgradeService
}

func (p ProgressionPresenter) Present(run domain.RunState) ProgressionViewModel {
	owned := make([]string, 0, len(run.Progression.PurchasedUpgradeIDs))
	for _, upgradeID := range run.Progression.PurchasedUpgradeIDs {
		upgradeName := upgradeID
		if upgrade, ok := p.Data.UpgradesByID[upgradeID]; ok {
			upgradeName = upgrade.Name
		}
		owned = append(owned, fmt.Sprintf("Installed upgrade: %s", upgradeName))
	}

	availableDefs := p.Upgrades.AvailableUpgrades(run, p.Data, p.Factions)
	available := make([]string, 0, len(availableDefs))
	for _, upgrade := range availableDefs {
		available = append(available, fmt.Sprintf("Available upgrade: %s (%d cr)", upgrade.Name, upgrade.CostCredits))
	}

	specializations := make([]string, 0, len(run.Progression.SpecializationFlags))
	for specialization, active := range run.Progression.SpecializationFlags {
		if !active {
			continue
		}
		specializations = append(specializations, fmt.Sprintf("Specialization active: %s", specialization))
	}
	sort.Strings(specializations)

	return ProgressionViewModel{
		OwnedUpgradeNotices:     owned,
		AvailableUpgradeNotices: available,
		SpecializationNotices:   specializations,
	}
}

func (p ProgressionPresenter) Summary(run domain.RunState) string {
	viewModel := p.Present(run)
	lines := append([]string{}, viewModel.OwnedUpgradeNotices...)
	lines = append(lines, viewModel.AvailableUpgradeNotices...)
	lines = append(lines, viewModel.SpecializationNotices...)
	return strings.Join(lines, "\n")
}
