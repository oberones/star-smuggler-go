package integration_test

import (
	"testing"

	"github.com/oberones/star-smuggler-go/internal/application"
	"github.com/oberones/star-smuggler-go/internal/services"
)

func TestBaselineTravelConsumesCreditsAndAdvancesJumpProgression(t *testing.T) {
	t.Parallel()

	snapshot := loadSnapshot(t)
	commands := application.NewRunCommands(snapshot, nil, seededRuntime(7))

	run, err := commands.StartNewRun()
	if err != nil {
		t.Fatalf("start new run: %v", err)
	}

	originPortID := run.Player.CurrentPortID
	quotes, err := commands.PreviewTravel(run)
	if err != nil {
		t.Fatalf("preview travel: %v", err)
	}
	if len(quotes) == 0 {
		t.Fatal("expected travel destinations")
	}

	destination := quotes[0]
	startingCredits := run.Player.Credits

	message, err := commands.CommitBaselineTravel(&run, destination.Destination.ID)
	if err != nil {
		t.Fatalf("commit baseline travel: %v", err)
	}

	if run.Player.Credits != startingCredits-destination.Cost {
		t.Fatalf("expected credits to drop by %d, now have %d", destination.Cost, run.Player.Credits)
	}
	if run.Player.CurrentPortID != destination.Destination.ID {
		t.Fatalf("expected current port to change to %q, got %q", destination.Destination.ID, run.Player.CurrentPortID)
	}
	if run.TotalJumps != 1 || run.JumpsSinceLastUpdate != 1 {
		t.Fatalf("expected jump counters to advance to 1, got total=%d sinceRefresh=%d", run.TotalJumps, run.JumpsSinceLastUpdate)
	}
	if run.PendingRoute == nil || run.PendingRoute.OriginPortID != originPortID || run.PendingRoute.DestinationPortID != destination.Destination.ID {
		t.Fatalf("expected resolved route state to be recorded, got %#v", run.PendingRoute)
	}
	if run.RecentEvent != nil {
		t.Fatalf("expected baseline travel not to resolve an event yet, got %#v", run.RecentEvent)
	}
	if message == "" {
		t.Fatal("expected travel result message")
	}
}

func TestBaselineTravelRequiresCreditsInsteadOfFuelMeter(t *testing.T) {
	t.Parallel()

	snapshot := loadSnapshot(t)
	commands := application.NewRunCommands(snapshot, nil, seededRuntime(11))

	run, err := commands.StartNewRun()
	if err != nil {
		t.Fatalf("start new run: %v", err)
	}

	quotes, err := commands.PreviewTravel(run)
	if err != nil {
		t.Fatalf("preview travel: %v", err)
	}
	if len(quotes) == 0 {
		t.Fatal("expected travel destinations")
	}

	run.Player.Credits = 0
	if _, err := commands.CommitBaselineTravel(&run, quotes[0].Destination.ID); err == nil {
		t.Fatal("expected travel to fail when route cost cannot be paid")
	}

	_ = services.TravelService{}
}
