package integration_test

import (
	"testing"

	"github.com/oberones/star-smuggler-go/internal/domain"
	"github.com/oberones/star-smuggler-go/internal/services"
)

func TestRunEvaluatorDistinguishesRecoverableAndStrandedStates(t *testing.T) {
	t.Parallel()

	snapshot := loadSnapshot(t)
	runtime := seededRuntime(23)
	economy := services.EconomyService{}
	travel := services.TravelService{}
	evaluator := services.RunEvaluator{}
	markets := economy.CreateInitialMarkets(snapshot, runtime.RNG)

	recoverable := domain.NewRunState()
	recoverable.Player.CurrentPortID = "mercury"
	recoverable.Player.Credits = 0
	recoverable.MarketsByPortID = markets
	recoverable.Cargo.SetQuantity("synthspice", 2)

	if evaluator.IsGameOver(recoverable, snapshot, economy, travel) {
		t.Fatal("expected run with sellable cargo to remain recoverable")
	}

	stranded := domain.NewRunState()
	stranded.Player.CurrentPortID = "mercury"
	stranded.Player.Credits = 0
	stranded.MarketsByPortID = markets

	if !evaluator.IsGameOver(stranded, snapshot, economy, travel) {
		t.Fatal("expected run without credits or sellable cargo to be game over")
	}
}
