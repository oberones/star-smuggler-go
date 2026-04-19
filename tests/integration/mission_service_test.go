package integration_test

import (
	"testing"

	"github.com/oberones/star-smuggler-go/internal/domain"
	"github.com/oberones/star-smuggler-go/internal/services"
)

func TestMissionServiceAcceptsAndCompletesDeliveryMission(t *testing.T) {
	t.Parallel()

	snapshot := loadSnapshot(t)
	run := domain.NewRunState()
	run.Player.CurrentPortID = "mars"
	run.FactionStandings = domain.DefaultFactionStandings(snapshot)
	run.ActiveMissions = map[string]domain.MissionState{}
	run.CompletedMissionIDs = []string{}
	run.Cargo.SetQuantity("synthspice", 2)

	mission := snapshot.MissionsByID["guild_supply_run"]
	service := services.MissionService{}
	factions := services.FactionService{}
	upgrades := services.UpgradeService{}

	accepted := service.AcceptMission(&run, mission)
	if accepted.Status != domain.MissionStatusInProgress {
		t.Fatalf("expected mission to enter in-progress state, got %q", accepted.Status)
	}

	run.Player.CurrentPortID = mission.DestinationPortID
	completed, failed := service.ResolveTravelArrival(&run, snapshot, factions, upgrades)
	if len(failed) != 0 {
		t.Fatalf("expected no failed missions, got %#v", failed)
	}
	if len(completed) != 1 {
		t.Fatalf("expected one completed mission, got %#v", completed)
	}
	if run.Cargo.QuantityFor("synthspice") != 0 {
		t.Fatalf("expected cargo to be delivered, still have %d units", run.Cargo.QuantityFor("synthspice"))
	}
	if !sliceContainsString(run.CompletedMissionIDs, mission.ID) {
		t.Fatalf("expected mission %q to be marked completed", mission.ID)
	}
}

func sliceContainsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
