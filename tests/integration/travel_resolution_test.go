package integration_test

import (
	"strings"
	"testing"

	"github.com/oberones/star-smuggler-go/internal/application"
	"github.com/oberones/star-smuggler-go/internal/domain"
	"github.com/oberones/star-smuggler-go/internal/services"
)

func TestTravelCommandsResolveArrivalRefreshJumpCountAndEvent(t *testing.T) {
	t.Parallel()

	baseSnapshot := loadSnapshot(t)
	creditsLossEvent := domain.EventDefinition{
		ID:                  "customs_shakedown",
		Name:                "Customs Shake-Down",
		DescriptionTemplate: "Local security stops you for a 'random inspection' and demands a bribe of {credits} credits.",
		EffectType:          domain.EventEffectCreditsLossScaled,
		Weight:              1,
		Parameters: map[string]float64{
			"minimum":    25,
			"maximum":    25,
			"minPercent": 0.05,
			"maxPercent": 0.05,
		},
	}
	snapshot := domain.NewDataSnapshot(
		baseSnapshot.Ports,
		baseSnapshot.Items,
		[]domain.EventDefinition{creditsLossEvent},
		baseSnapshot.Factions,
		baseSnapshot.Missions,
		baseSnapshot.StoryArcs,
		baseSnapshot.Upgrades,
	)

	runCommands := application.NewRunCommands(snapshot, nil, seededRuntime(17))
	run, err := runCommands.StartNewRun()
	if err != nil {
		t.Fatalf("start new run: %v", err)
	}

	travelRuntime := services.NewRuntimeContext(&stubRNG{
		intValues:   []int{0, 0, 0},
		floatValues: []float64{0},
	}, nil)
	travelCommands := application.NewTravelCommands(snapshot, travelRuntime)

	quotes, err := runCommands.PreviewTravel(run)
	if err != nil {
		t.Fatalf("preview travel: %v", err)
	}
	if len(quotes) == 0 {
		t.Fatal("expected travel destinations")
	}

	destination := quotes[0]
	run.MarketsByPortID[destination.Destination.ID] = domain.MarketSnapshot{
		PortID:           destination.Destination.ID,
		AvailableItemIDs: nil,
		PricesByItemID:   map[string]int{},
	}

	if _, err := travelCommands.BeginTravel(&run, destination.Destination.ID); err != nil {
		t.Fatalf("begin travel: %v", err)
	}
	if run.PendingRoute == nil || run.PendingRoute.Status != domain.RouteStatusAnimating {
		t.Fatalf("expected animating pending route, got %#v", run.PendingRoute)
	}

	startingCredits := run.Player.Credits
	resolution, err := travelCommands.ResolvePendingTravel(&run)
	if err != nil {
		t.Fatalf("resolve pending travel: %v", err)
	}

	expectedCredits := startingCredits - destination.Cost - 25
	if run.Player.Credits != expectedCredits {
		t.Fatalf("expected credits %d after fare and event, got %d", expectedCredits, run.Player.Credits)
	}
	if run.Player.CurrentPortID != destination.Destination.ID {
		t.Fatalf("expected arrival at %q, got %q", destination.Destination.ID, run.Player.CurrentPortID)
	}
	if run.TotalJumps != 1 || run.JumpsSinceLastUpdate != 1 {
		t.Fatalf("expected jump counters to advance to 1, got total=%d since=%d", run.TotalJumps, run.JumpsSinceLastUpdate)
	}
	if run.PendingRoute == nil || run.PendingRoute.Status != domain.RouteStatusResolved {
		t.Fatalf("expected resolved pending route, got %#v", run.PendingRoute)
	}
	if run.RecentEvent == nil || run.RecentEvent.EventID != "customs_shakedown" {
		t.Fatalf("expected customs_shakedown event, got %#v", run.RecentEvent)
	}
	if market := run.MarketsByPortID[destination.Destination.ID]; len(market.AvailableItemIDs) == 0 {
		t.Fatal("expected destination market to refresh available goods on arrival")
	}
	if resolution.AppliedEvent == nil || !strings.Contains(resolution.Message, resolution.AppliedEvent.ResolvedDescription) {
		t.Fatalf("expected travel resolution message to include event text, got %#v", resolution)
	}
}

func TestTravelAnimationDurationScalesWithZoneDifference(t *testing.T) {
	t.Parallel()

	snapshot := loadSnapshot(t)
	commands := application.NewTravelCommands(snapshot, seededRuntime(23))
	run := domain.NewRunState()
	run.PendingRoute = &domain.RouteState{
		OriginPortID:      "mars",
		DestinationPortID: "titan",
		TravelCost:        34,
		Status:            domain.RouteStatusAnimating,
	}

	duration := commands.AnimationDuration(run.PendingRoute)
	if duration <= 2.0 {
		t.Fatalf("expected a longer travel animation duration for cross-zone travel, got %f", duration)
	}
}
