package domain

import "fmt"

type RandomIndexSource interface {
	Intn(n int) int
}

func CreateNewRun(data DataSnapshot, marketSnapshots map[string]MarketSnapshot, rng RandomIndexSource) (RunState, error) {
	if len(data.Ports) == 0 {
		return RunState{}, fmt.Errorf("cannot start a run without any port definitions")
	}
	if len(data.Items) == 0 {
		return RunState{}, fmt.Errorf("cannot start a run without any item definitions")
	}
	if rng == nil {
		return RunState{}, fmt.Errorf("random source is required")
	}

	innerPorts := make([]PortDefinition, 0)
	for _, port := range data.Ports {
		if port.Zone == PortZoneInner {
			innerPorts = append(innerPorts, port)
		}
	}

	if len(innerPorts) == 0 {
		return RunState{}, fmt.Errorf("cannot start a run without at least one inner-zone port")
	}

	startingPort := innerPorts[rng.Intn(len(innerPorts))]
	run := NewRunState()
	run.Player.CurrentPortID = startingPort.ID
	run.MarketsByPortID = cloneMarkets(marketSnapshots)
	run.FactionStandings = DefaultFactionStandings(data)
	run.Story = NewStoryState()
	return run, nil
}

func cloneMarkets(source map[string]MarketSnapshot) map[string]MarketSnapshot {
	if source == nil {
		return map[string]MarketSnapshot{}
	}

	result := make(map[string]MarketSnapshot, len(source))
	for portID, market := range source {
		result[portID] = MarketSnapshot{
			PortID:           market.PortID,
			AvailableItemIDs: append([]string(nil), market.AvailableItemIDs...),
			PricesByItemID:   clonePriceMap(market.PricesByItemID),
		}
	}
	return result
}

func clonePriceMap(source map[string]int) map[string]int {
	if source == nil {
		return map[string]int{}
	}

	result := make(map[string]int, len(source))
	for itemID, price := range source {
		result[itemID] = price
	}
	return result
}
