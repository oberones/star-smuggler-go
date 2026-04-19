package integration_test

import (
	"testing"

	"github.com/oberones/star-smuggler-go/internal/domain"
	"github.com/oberones/star-smuggler-go/internal/services"
)

func TestFactionServiceAppliesStandingTransitions(t *testing.T) {
	t.Parallel()

	snapshot := loadSnapshot(t)
	run := domain.NewRunState()
	run.FactionStandings = domain.DefaultFactionStandings(snapshot)

	service := services.FactionService{}
	standing, ok := service.ApplyStandingDelta(&run, snapshot, "freeguild", 25, "Completed a guild run")
	if !ok {
		t.Fatal("expected faction standing update to succeed")
	}
	if standing.StandingTier != "Trusted" {
		t.Fatalf("expected Trusted tier, got %q", standing.StandingTier)
	}
	if standing.LastChangeReason != "Completed a guild run" {
		t.Fatalf("expected change reason to be recorded, got %q", standing.LastChangeReason)
	}
}
