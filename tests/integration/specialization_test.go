package integration_test

import (
	"testing"

	"github.com/oberones/star-smuggler-go/internal/domain"
	"github.com/oberones/star-smuggler-go/internal/services"
)

func TestSpeedSpecializationReducesTravelCosts(t *testing.T) {
	t.Parallel()

	snapshot := loadSnapshot(t)
	run := domain.NewRunState()
	run.Player.CurrentPortID = "mars"
	run.Player.Credits = 900
	run.FactionStandings = domain.DefaultFactionStandings(snapshot)

	upgrades := services.UpgradeService{}
	factions := services.FactionService{}
	travel := services.TravelService{}

	result := upgrades.PurchaseUpgrade(&run, snapshot, factions, "slipstream_tuner")
	if !result.Succeeded {
		t.Fatalf("expected slipstream_tuner purchase to succeed, got %#v", result)
	}
	if !run.Progression.HasSpecialization(domain.ShipSpecializationSpeed) {
		t.Fatal("expected speed specialization to be active")
	}

	baseCost := travel.GetTravelCost(snapshot.PortsByID["mars"], snapshot.PortsByID["titan"])
	adjustedCost := upgrades.AdjustTravelCost(run, baseCost, snapshot)
	if adjustedCost >= baseCost {
		t.Fatalf("expected adjusted travel cost %d to be lower than base cost %d", adjustedCost, baseCost)
	}
	if adjustedCost != (baseCost*80)/100 {
		t.Fatalf("expected 20 percent travel discount, got %d from base %d", adjustedCost, baseCost)
	}
}

func TestInfluenceSpecializationBoostsMissionRewards(t *testing.T) {
	t.Parallel()

	snapshot := loadSnapshot(t)
	run := domain.NewRunState()
	run.Player.CurrentPortID = "mars"
	run.Player.Credits = 1200
	run.FactionStandings = domain.DefaultFactionStandings(snapshot)
	run.ActiveMissions = map[string]domain.MissionState{}
	run.Cargo.SetQuantity("synthspice", 2)

	upgrades := services.UpgradeService{}
	factions := services.FactionService{}
	missions := services.MissionService{}

	if _, ok := factions.ApplyStandingDelta(&run, snapshot, "freeguild", 25, "Trusted specialist"); !ok {
		t.Fatal("expected freeguild standing update to succeed")
	}

	purchase := upgrades.PurchaseUpgrade(&run, snapshot, factions, "guild_credential_spoof")
	if !purchase.Succeeded {
		t.Fatalf("expected influence upgrade to unlock and purchase, got %#v", purchase)
	}
	if !run.Progression.HasSpecialization(domain.ShipSpecializationInfluence) {
		t.Fatal("expected influence specialization to be active")
	}

	mission := snapshot.MissionsByID["guild_supply_run"]
	accepted := missions.AcceptMission(&run, mission)
	if accepted.Status != domain.MissionStatusInProgress {
		t.Fatalf("expected mission to become active, got %#v", accepted)
	}

	creditsBeforeDelivery := run.Player.Credits
	run.Player.CurrentPortID = mission.DestinationPortID
	completed, failed := missions.ResolveTravelArrival(&run, snapshot, factions, upgrades)
	if len(failed) != 0 {
		t.Fatalf("expected no failed missions, got %#v", failed)
	}
	if len(completed) != 1 {
		t.Fatalf("expected one completed mission, got %#v", completed)
	}

	rewardDelta := run.Player.Credits - creditsBeforeDelivery
	if rewardDelta != 225 {
		t.Fatalf("expected mission reward bonus to produce 225 credits, got %d", rewardDelta)
	}
}
