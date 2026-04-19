package application

import (
	"context"
	"fmt"

	"github.com/oberones/star-smuggler-go/internal/domain"
	"github.com/oberones/star-smuggler-go/internal/services"
)

type RunCommands struct {
	Data     domain.DataSnapshot
	Runtime  services.RuntimeContext
	SaveRepo SaveRepository
	Economy  services.EconomyService
	Trade    services.TradeService
	Travel   services.TravelService
	RunEval  services.RunEvaluator
}

func NewRunCommands(data domain.DataSnapshot, saveRepo SaveRepository, runtime services.RuntimeContext) RunCommands {
	return RunCommands{
		Data:     data,
		Runtime:  runtime,
		SaveRepo: saveRepo,
	}
}

func (c RunCommands) StartNewRun() (domain.RunState, error) {
	markets := c.Economy.CreateInitialMarkets(c.Data, c.Runtime.RNG)
	return domain.CreateNewRun(c.Data, markets, c.Runtime.RNG)
}

func (c RunCommands) Continue(ctx context.Context) (domain.RunState, error) {
	if c.SaveRepo == nil {
		return domain.RunState{}, fmt.Errorf("save repository is not configured")
	}
	return c.SaveRepo.Load(ctx)
}

func (c RunCommands) Save(ctx context.Context, run domain.RunState) error {
	if c.SaveRepo == nil {
		return fmt.Errorf("save repository is not configured")
	}
	return c.SaveRepo.Save(ctx, run)
}

func (c RunCommands) Buy(run *domain.RunState, itemID string, quantity int) (services.TradeResult, error) {
	item, ok := c.Data.ItemsByID[itemID]
	if !ok {
		return services.TradeResult{}, fmt.Errorf("item %q was not found", itemID)
	}

	market, ok := c.Economy.GetCurrentMarket(*run)
	if !ok {
		return services.TradeResult{}, fmt.Errorf("current market is not available")
	}

	return c.Trade.Buy(run, market, item, quantity), nil
}

func (c RunCommands) Sell(run *domain.RunState, itemID string, quantity int) (services.TradeResult, error) {
	item, ok := c.Data.ItemsByID[itemID]
	if !ok {
		return services.TradeResult{}, fmt.Errorf("item %q was not found", itemID)
	}

	market, ok := c.Economy.GetCurrentMarket(*run)
	if !ok {
		return services.TradeResult{}, fmt.Errorf("current market is not available")
	}

	return c.Trade.Sell(run, market, item, quantity), nil
}

func (c RunCommands) PreviewTravel(run domain.RunState) ([]services.TravelQuote, error) {
	origin, ok := c.Data.PortsByID[run.Player.CurrentPortID]
	if !ok {
		return nil, fmt.Errorf("current port %q was not found", run.Player.CurrentPortID)
	}

	destinations := c.Travel.GetDestinationsFromPort(origin, c.Data.Ports)
	quotes := make([]services.TravelQuote, 0, len(destinations))
	for _, destination := range destinations {
		quotes = append(quotes, services.TravelQuote{
			Destination: destination,
			Cost:        c.Travel.GetTravelCost(origin, destination),
		})
	}
	return quotes, nil
}

func (c RunCommands) CommitBaselineTravel(run *domain.RunState, destinationPortID string) (string, error) {
	origin, ok := c.Data.PortsByID[run.Player.CurrentPortID]
	if !ok {
		return "", fmt.Errorf("current port %q was not found", run.Player.CurrentPortID)
	}

	destination, ok := c.Data.PortsByID[destinationPortID]
	if !ok {
		return "", fmt.Errorf("destination port %q was not found", destinationPortID)
	}

	cost := c.Travel.GetTravelCost(origin, destination)
	if run.Player.Credits < cost {
		return "", fmt.Errorf("you need %d credits to reach %s", cost, destination.Name)
	}

	run.Player.Credits -= cost
	run.Player.CurrentPortID = destination.ID
	run.TotalJumps++
	run.JumpsSinceLastUpdate++
	run.PendingRoute = &domain.RouteState{
		OriginPortID:      origin.ID,
		DestinationPortID: destination.ID,
		TravelCost:        cost,
		Status:            domain.RouteStatusResolved,
	}
	run.RecentEvent = nil

	if run.JumpsSinceLastUpdate > 3 {
		c.Economy.RefreshAllPrices(run, c.Data, c.Runtime.RNG)
	} else {
		c.Economy.RefreshAvailableGoods(run, c.Data, destination.ID, c.Runtime.RNG)
	}

	return fmt.Sprintf("Traveled to %s for %d credits.", destination.Name, cost), nil
}

func (c RunCommands) RouteForRun(run domain.RunState) Route {
	if c.RunEval.IsGameOver(run, c.Data, c.Economy, c.Travel) {
		return RouteGameOver
	}
	return RoutePortOverview
}
