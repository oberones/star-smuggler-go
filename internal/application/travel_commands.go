package application

import (
	"fmt"
	"strings"

	"github.com/oberones/star-smuggler-go/internal/domain"
	"github.com/oberones/star-smuggler-go/internal/services"
)

type TravelCommands struct {
	Data    domain.DataSnapshot
	Runtime services.RuntimeContext
	Economy services.EconomyService
	Balance services.EconomyBalanceService
	Travel  services.TravelService
	Events  services.EventService
	RunEval services.RunEvaluator
}

func NewTravelCommands(data domain.DataSnapshot, runtime services.RuntimeContext) TravelCommands {
	return TravelCommands{
		Data:    data,
		Runtime: runtime,
	}
}

func (c TravelCommands) BeginTravel(run *domain.RunState, destinationPortID string) (*domain.RouteState, error) {
	origin, ok := c.Data.PortsByID[run.Player.CurrentPortID]
	if !ok {
		return nil, fmt.Errorf("current port %q was not found", run.Player.CurrentPortID)
	}

	destination, ok := c.Data.PortsByID[destinationPortID]
	if !ok {
		return nil, fmt.Errorf("destination port %q was not found", destinationPortID)
	}

	cost := c.Travel.GetTravelCost(origin, destination) + c.Balance.AdditionalRouteCost(*run, origin.ID, destination.ID)
	if run.Player.Credits < cost {
		return nil, fmt.Errorf("you need %d credits to reach %s", cost, destination.Name)
	}

	run.PendingRoute = &domain.RouteState{
		OriginPortID:      origin.ID,
		DestinationPortID: destination.ID,
		TravelCost:        cost,
		Status:            domain.RouteStatusAnimating,
	}

	return run.PendingRoute, nil
}

func (c TravelCommands) ResolvePendingTravel(run *domain.RunState) (domain.TravelResolution, error) {
	if run.PendingRoute == nil {
		return domain.TravelResolution{}, fmt.Errorf("there is no pending travel route to resolve")
	}

	route := run.PendingRoute
	origin, ok := c.Data.PortsByID[route.OriginPortID]
	if !ok {
		return domain.TravelResolution{}, fmt.Errorf("origin port %q was not found", route.OriginPortID)
	}

	destination, ok := c.Data.PortsByID[route.DestinationPortID]
	if !ok {
		return domain.TravelResolution{}, fmt.Errorf("destination port %q was not found", route.DestinationPortID)
	}

	if run.Player.Credits < route.TravelCost {
		return domain.TravelResolution{}, fmt.Errorf("you need %d credits to reach %s", route.TravelCost, destination.Name)
	}

	run.Player.Credits -= route.TravelCost
	run.Player.CurrentPortID = destination.ID
	run.TotalJumps++
	run.JumpsSinceLastUpdate++
	run.RecentEvent = nil
	c.Balance.RecordRoute(run, origin.ID, destination.ID)

	if run.JumpsSinceLastUpdate > 3 {
		c.Economy.RefreshAllPrices(run, c.Data, c.Runtime.RNG)
	} else {
		c.Economy.RefreshAvailableGoods(run, c.Data, destination.ID, c.Runtime.RNG)
	}
	c.Balance.ApplyMarketPressure(run, c.Data)

	run.RecentEvent = c.Events.TryResolveTravelEvent(run, c.Data, c.Economy, c.Runtime.RNG)
	run.PendingRoute.Status = domain.RouteStatusResolved

	message := fmt.Sprintf("Traveled to %s for %d credits.", destination.Name, route.TravelCost)
	if run.RecentEvent != nil && strings.TrimSpace(run.RecentEvent.ResolvedDescription) != "" {
		message += " " + run.RecentEvent.ResolvedDescription
	}

	return domain.TravelResolution{
		Message:      message,
		AppliedEvent: run.RecentEvent,
		Route:        run.PendingRoute,
	}, nil
}

func (c TravelCommands) AnimationDuration(route *domain.RouteState) float64 {
	if route == nil {
		return 2.5
	}

	origin, originOK := c.Data.PortsByID[route.OriginPortID]
	destination, destinationOK := c.Data.PortsByID[route.DestinationPortID]
	if !originOK || !destinationOK {
		return 2.5
	}

	zoneDifference := abs(zoneRank(origin.Zone) - zoneRank(destination.Zone))
	return 2.0 + (float64(zoneDifference) * 1.5)
}

func zoneRank(zone domain.PortZone) int {
	switch zone {
	case domain.PortZoneInner:
		return 0
	case domain.PortZoneOuter:
		return 1
	case domain.PortZoneFringe:
		return 2
	default:
		return 0
	}
}

func abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}
