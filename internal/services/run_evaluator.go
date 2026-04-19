package services

import "github.com/oberones/star-smuggler-go/internal/domain"

type RunEvaluator struct{}

func (e RunEvaluator) IsGameOver(run domain.RunState, data domain.DataSnapshot, economy EconomyService, travel TravelService) bool {
	currentPort, ok := data.PortsByID[run.Player.CurrentPortID]
	if !ok {
		return false
	}

	cheapestTravelCost := travel.GetCheapestTravelCostFromPort(currentPort, data.Ports)
	if run.Player.Credits >= cheapestTravelCost {
		return false
	}

	sellableValue := economy.GetSellableCargoValueAtCurrentPort(run, data)
	return run.Player.Credits+sellableValue < cheapestTravelCost
}
