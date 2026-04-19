package integration_test

import (
	"strings"
	"testing"

	"github.com/oberones/star-smuggler-go/internal/domain"
	godotpresentation "github.com/oberones/star-smuggler-go/internal/presentation/godot"
	"github.com/oberones/star-smuggler-go/internal/services"
)

func TestStoryPresenterSummarizesFactionMissionAndStoryNotices(t *testing.T) {
	t.Parallel()

	snapshot := loadSnapshot(t)
	run := domain.NewRunState()
	run.Player.CurrentPortID = "mars"
	run.FactionStandings = domain.DefaultFactionStandings(snapshot)
	run.ActiveMissions = map[string]domain.MissionState{
		"guild_supply_run": {
			MissionDefinitionID: "guild_supply_run",
			Status:              domain.MissionStatusInProgress,
			AcceptedAtJump:      0,
			DeadlineJump:        4,
			ProgressFlags:       map[string]bool{"cargo_loaded": true},
		},
	}
	run.Story = domain.NewStoryState()
	run.Story.ActiveStoryArcIDs = []string{"guild_intro_arc"}

	presenter := godotpresentation.StoryPresenter{
		Data:     snapshot,
		Factions: services.FactionService{},
		Missions: services.MissionService{},
	}

	summary := presenter.Summary(run)
	if !strings.Contains(summary, "Faction: Free Guild") {
		t.Fatalf("expected faction notice in summary, got %q", summary)
	}
	if !strings.Contains(summary, "Active mission: Guild Supply Run") {
		t.Fatalf("expected mission notice in summary, got %q", summary)
	}
	if !strings.Contains(summary, "Story arc: Whispers In The Cargo Bay") {
		t.Fatalf("expected story notice in summary, got %q", summary)
	}
}
