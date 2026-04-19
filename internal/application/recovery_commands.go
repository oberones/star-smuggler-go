package application

import (
	"fmt"

	"github.com/oberones/star-smuggler-go/internal/domain"
	"github.com/oberones/star-smuggler-go/internal/services"
)

type RecoveryCommands struct {
	Data    domain.DataSnapshot
	Economy services.EconomyService
	Travel  services.TravelService
	RunEval services.RunEvaluator
}

func NewRecoveryCommands(data domain.DataSnapshot) RecoveryCommands {
	return RecoveryCommands{
		Data: data,
	}
}

func (c RecoveryCommands) TryEmergencyRecovery(run *domain.RunState) (string, bool, error) {
	if run.EmergencyRecoveryUsed {
		return "Emergency recovery has already been used this run.", false, nil
	}

	if !c.RunEval.IsGameOver(*run, c.Data, c.Economy, c.Travel) {
		return "Emergency recovery is not needed right now.", false, nil
	}

	currentPort, ok := c.Data.PortsByID[run.Player.CurrentPortID]
	if !ok {
		return "", false, fmt.Errorf("current port %q was not found", run.Player.CurrentPortID)
	}

	cheapestTravelCost := c.Travel.GetCheapestTravelCostFromPort(currentPort, c.Data.Ports)
	recoveryGrant := cheapestTravelCost - run.Player.Credits
	if recoveryGrant < 25 {
		recoveryGrant = 25
	}

	run.Player.Credits += recoveryGrant
	run.EmergencyRecoveryUsed = true

	return fmt.Sprintf("Emergency recovery granted %d credits to get you moving again.", recoveryGrant), true, nil
}
