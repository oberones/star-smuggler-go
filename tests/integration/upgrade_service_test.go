package integration_test

import (
	"context"
	"strings"
	"testing"

	"github.com/oberones/star-smuggler-go/internal/application"
	"github.com/oberones/star-smuggler-go/internal/domain"
	"github.com/oberones/star-smuggler-go/internal/services"
)

func TestUpgradeServicePurchasesCargoUpgradeAndEnforcesRules(t *testing.T) {
	t.Parallel()

	snapshot := loadSnapshot(t)
	run := domain.NewRunState()
	run.Player.CurrentPortID = "mars"
	run.Player.Credits = 900
	run.FactionStandings = domain.DefaultFactionStandings(snapshot)

	service := services.UpgradeService{}
	factions := services.FactionService{}

	available := service.AvailableUpgrades(run, snapshot, factions)
	if len(available) < 2 {
		t.Fatalf("expected at least two starter upgrades, got %#v", available)
	}
	for _, upgrade := range available {
		if upgrade.ID == "guild_credential_spoof" {
			t.Fatalf("expected guild_credential_spoof to stay locked at neutral standing, got %#v", available)
		}
	}

	result := service.PurchaseUpgrade(&run, snapshot, factions, "expanded_cargo_pods")
	if !result.Succeeded {
		t.Fatalf("expected purchase to succeed, got %#v", result)
	}
	if run.Player.Credits != 575 {
		t.Fatalf("expected credits to drop to 575, got %d", run.Player.Credits)
	}
	if run.Player.CargoLimit != domain.StartingCargoLimit+12 {
		t.Fatalf("expected cargo limit %d, got %d", domain.StartingCargoLimit+12, run.Player.CargoLimit)
	}
	if !run.Progression.HasUpgrade("expanded_cargo_pods") {
		t.Fatal("expected cargo upgrade to be recorded in progression state")
	}
	if !run.Progression.HasSpecialization(domain.ShipSpecializationCargo) {
		t.Fatal("expected cargo specialization to activate after purchase")
	}

	duplicate := service.PurchaseUpgrade(&run, snapshot, factions, "expanded_cargo_pods")
	if duplicate.Succeeded || !strings.Contains(duplicate.Message, "already installed") {
		t.Fatalf("expected duplicate purchase to fail, got %#v", duplicate)
	}

	locked := service.PurchaseUpgrade(&run, snapshot, factions, "guild_credential_spoof")
	if locked.Succeeded || !strings.Contains(locked.Message, "requires Trusted standing") {
		t.Fatalf("expected faction-gated upgrade to fail while locked, got %#v", locked)
	}
}

func TestAppPurchaseUpgradeAutosavesAndKeepsTradeRoute(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	snapshot := loadSnapshot(t)
	saveRepo := &memorySaveRepository{}
	app := application.NewApp(stubContentLoader{snapshot: snapshot}, saveRepo, seededRuntime(31))

	if err := app.Bootstrap(ctx); err != nil {
		t.Fatalf("bootstrap app: %v", err)
	}
	if err := app.StartNewRun(ctx); err != nil {
		t.Fatalf("start new run: %v", err)
	}
	if err := app.OpenTrade(); err != nil {
		t.Fatalf("open trade: %v", err)
	}

	result, err := app.PurchaseUpgrade(ctx, "expanded_cargo_pods")
	if err != nil {
		t.Fatalf("purchase upgrade through app: %v", err)
	}
	if !result.Succeeded {
		t.Fatalf("expected app purchase to succeed, got %#v", result)
	}
	if route := app.CurrentRoute(); route != application.RouteTrade {
		t.Fatalf("expected successful upgrade purchase to keep the player in trade flow, got %v", route)
	}
	if saveRepo.run.Player.CargoLimit != domain.StartingCargoLimit+12 {
		t.Fatalf("expected autosave to persist upgraded cargo limit, got %d", saveRepo.run.Player.CargoLimit)
	}
	if !saveRepo.run.Progression.HasUpgrade("expanded_cargo_pods") {
		t.Fatal("expected autosave to persist purchased upgrade")
	}
}
