package integration_test

import (
	"testing"

	"github.com/oberones/star-smuggler-go/internal/domain"
	"github.com/oberones/star-smuggler-go/internal/services"
)

func TestTradeServiceBuyAndSellParity(t *testing.T) {
	t.Parallel()

	snapshot := loadSnapshot(t)
	economy := services.EconomyService{}
	trade := services.TradeService{}
	run := domain.NewRunState()
	run.Player.CurrentPortID = "mars"
	run.MarketsByPortID = economy.CreateInitialMarkets(snapshot, seededRuntime(17).RNG)

	market := run.MarketsByPortID["mars"]
	item := snapshot.ItemsByID[market.AvailableItemIDs[0]]
	price := market.PricesByItemID[item.ID]

	buyResult := trade.Buy(&run, market, item, 1)
	if !buyResult.Succeeded {
		t.Fatalf("expected buy to succeed, got %q", buyResult.Message)
	}
	if run.Cargo.QuantityFor(item.ID) != 1 {
		t.Fatalf("expected to own 1 %s, got %d", item.Name, run.Cargo.QuantityFor(item.ID))
	}
	if run.Player.Credits != domain.StartingCredits-price {
		t.Fatalf("unexpected credits after buy: %d", run.Player.Credits)
	}

	sellResult := trade.Sell(&run, market, item, 1)
	if !sellResult.Succeeded {
		t.Fatalf("expected sell to succeed, got %q", sellResult.Message)
	}
	if run.Cargo.QuantityFor(item.ID) != 0 {
		t.Fatalf("expected to own 0 %s after sell, got %d", item.Name, run.Cargo.QuantityFor(item.ID))
	}
	if run.Player.Credits != domain.StartingCredits {
		t.Fatalf("unexpected credits after sell: %d", run.Player.Credits)
	}
}

func TestTradeServiceRejectsInvalidTransactions(t *testing.T) {
	t.Parallel()

	snapshot := loadSnapshot(t)
	economy := services.EconomyService{}
	trade := services.TradeService{}
	run := domain.NewRunState()
	run.Player.CurrentPortID = "mars"
	run.MarketsByPortID = economy.CreateInitialMarkets(snapshot, seededRuntime(19).RNG)

	market := run.MarketsByPortID["mars"]
	item := snapshot.ItemsByID[market.AvailableItemIDs[0]]

	run.Player.Credits = 0
	if result := trade.Buy(&run, market, item, 1); result.Succeeded {
		t.Fatal("expected buy without credits to fail")
	}

	run.Player.Credits = domain.StartingCredits
	run.Player.CargoLimit = 0
	if result := trade.Buy(&run, market, item, 1); result.Succeeded {
		t.Fatal("expected buy beyond cargo limit to fail")
	}

	if result := trade.Sell(&run, market, item, 1); result.Succeeded {
		t.Fatal("expected sell without owned cargo to fail")
	}
}
