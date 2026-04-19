package integration_test

import (
	"strings"
	"testing"

	"github.com/oberones/star-smuggler-go/internal/application"
	"github.com/oberones/star-smuggler-go/internal/domain"
	"github.com/oberones/star-smuggler-go/internal/services"
)

func TestStoryCommandsIntegrateMissionAndStoryAcrossTradeAndTravel(t *testing.T) {
	t.Parallel()

	snapshot := loadSnapshot(t)
	run := domain.NewRunState()
	run.Player.CurrentPortID = "mars"
	run.FactionStandings = domain.DefaultFactionStandings(snapshot)
	run.ActiveMissions = map[string]domain.MissionState{}
	run.Story = domain.NewStoryState()

	commands := application.NewStoryCommands(snapshot)
	commands.Factions = services.FactionService{}
	commands.Missions = services.MissionService{}
	commands.Stories = services.StoryService{}

	update, err := commands.AcceptMission(&run, "guild_supply_run")
	if err != nil {
		t.Fatalf("accept mission: %v", err)
	}
	if len(update.Notices) == 0 || !strings.Contains(update.Notices[0], "Accepted mission") {
		t.Fatalf("expected accept-mission notice, got %#v", update.Notices)
	}

	run.Cargo.SetQuantity("synthspice", 2)
	tradeUpdate := commands.SyncAfterTrade(&run, "synthspice", 2, true)
	if len(tradeUpdate.Notices) == 0 || !strings.Contains(tradeUpdate.Notices[0], "Mission cargo secured") {
		t.Fatalf("expected cargo-loaded notice, got %#v", tradeUpdate.Notices)
	}

	if _, ok := commands.Factions.ApplyStandingDelta(&run, snapshot, "freeguild", 25, "Built trust"); !ok {
		t.Fatal("expected faction standing update")
	}

	run.Player.CurrentPortID = "titan"
	run.TotalJumps = 1
	travelUpdate := commands.SyncAfterTravel(&run)
	summary := travelUpdate.Summary()
	if !strings.Contains(summary, "Mission completed: Guild Supply Run") {
		t.Fatalf("expected mission completion notice, got %q", summary)
	}
	if !strings.Contains(summary, "Story unlocked: Whispers In The Cargo Bay") {
		t.Fatalf("expected story unlock notice, got %q", summary)
	}
}
