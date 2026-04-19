package integration_test

import (
	"testing"

	"github.com/oberones/star-smuggler-go/internal/application"
	"github.com/oberones/star-smuggler-go/internal/domain"
	"github.com/oberones/star-smuggler-go/internal/services"
)

func TestEmergencyRecoveryDistinguishesRecoverableSetbackFromTrueGameOver(t *testing.T) {
	t.Parallel()

	snapshot := loadSnapshot(t)
	economy := services.EconomyService{}
	travel := services.TravelService{}
	evaluator := services.RunEvaluator{}
	recovery := application.NewRecoveryCommands(snapshot)

	run := domain.NewRunState()
	run.Player.CurrentPortID = "mercury"
	run.Player.Credits = 0
	run.MarketsByPortID = economy.CreateInitialMarkets(snapshot, seededRuntime(37).RNG)

	if !evaluator.IsGameOver(run, snapshot, economy, travel) {
		t.Fatal("expected stranded run to start in game-over state")
	}

	message, recovered, err := recovery.TryEmergencyRecovery(&run)
	if err != nil {
		t.Fatalf("apply emergency recovery: %v", err)
	}
	if !recovered {
		t.Fatalf("expected recovery to succeed, got message %q", message)
	}
	if !run.EmergencyRecoveryUsed {
		t.Fatal("expected recovery flag to be set")
	}
	if evaluator.IsGameOver(run, snapshot, economy, travel) {
		t.Fatal("expected recovery to move run out of immediate game-over state")
	}

	run.Player.Credits = 0
	secondMessage, secondRecovered, err := recovery.TryEmergencyRecovery(&run)
	if err != nil {
		t.Fatalf("retry emergency recovery: %v", err)
	}
	if secondRecovered {
		t.Fatalf("expected second recovery to fail, got message %q", secondMessage)
	}
	if !evaluator.IsGameOver(run, snapshot, economy, travel) {
		t.Fatal("expected run to return to true game-over state after the one recovery is spent")
	}
}
