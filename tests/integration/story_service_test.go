package integration_test

import (
	"testing"

	"github.com/oberones/star-smuggler-go/internal/domain"
	"github.com/oberones/star-smuggler-go/internal/services"
)

func TestStoryServiceActivatesAndAdvancesArcFromFactionStanding(t *testing.T) {
	t.Parallel()

	snapshot := loadSnapshot(t)
	run := domain.NewRunState()
	run.FactionStandings = domain.DefaultFactionStandings(snapshot)
	run.Story = domain.NewStoryState()

	factions := services.FactionService{}
	stories := services.StoryService{}

	if _, ok := factions.ApplyStandingDelta(&run, snapshot, "freeguild", 25, "Built trust"); !ok {
		t.Fatal("expected faction standing update")
	}

	activated := stories.ActivateEligibleArcs(&run, snapshot, factions)
	if len(activated) != 1 || activated[0] != "guild_intro_arc" {
		t.Fatalf("expected guild intro arc to activate, got %#v", activated)
	}

	firstBeat, ok := stories.AdvanceArc(&run, snapshot, "guild_intro_arc")
	if !ok || firstBeat.ID != "guild_intro_offer" {
		t.Fatalf("expected first story beat to resolve, got %#v", firstBeat)
	}

	secondBeat, ok := stories.AdvanceArc(&run, snapshot, "guild_intro_arc")
	if !ok || secondBeat.ID != "guild_intro_commitment" {
		t.Fatalf("expected second story beat to resolve, got %#v", secondBeat)
	}
	if !run.Story.StoryFlags["guild_intro_arc_complete"] {
		t.Fatal("expected story arc completion flag to be set")
	}
}
