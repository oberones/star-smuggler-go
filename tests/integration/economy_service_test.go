package integration_test

import (
	"testing"

	"github.com/oberones/star-smuggler-go/internal/domain"
	"github.com/oberones/star-smuggler-go/internal/services"
)

func TestEconomyServiceCreatesParityMarkets(t *testing.T) {
	t.Parallel()

	snapshot := loadSnapshot(t)
	runtime := seededRuntime(5)
	economy := services.EconomyService{}

	markets := economy.CreateInitialMarkets(snapshot, runtime.RNG)
	if len(markets) != len(snapshot.Ports) {
		t.Fatalf("expected %d markets, got %d", len(snapshot.Ports), len(markets))
	}

	for _, port := range snapshot.Ports {
		market, ok := markets[port.ID]
		if !ok {
			t.Fatalf("missing market for port %q", port.ID)
		}
		if len(market.AvailableItemIDs) != 6 {
			t.Fatalf("expected 6 available items for %q, got %d", port.ID, len(market.AvailableItemIDs))
		}
		if len(market.PricesByItemID) != len(snapshot.Items) {
			t.Fatalf("expected prices for all items at %q", port.ID)
		}

		matchingRarityCount := 0
		for _, itemID := range market.AvailableItemIDs {
			item := snapshot.ItemsByID[itemID]
			if item.Rarity == expectedZoneRarity(port.Zone) {
				matchingRarityCount++
			}
		}
		if matchingRarityCount != 4 {
			t.Fatalf("expected 4 zone-matching items for %q, got %d", port.ID, matchingRarityCount)
		}
	}
}

func TestEconomyServiceCalculatesSellableCargoValue(t *testing.T) {
	t.Parallel()

	snapshot := loadSnapshot(t)
	runtime := seededRuntime(13)
	economy := services.EconomyService{}
	markets := economy.CreateInitialMarkets(snapshot, runtime.RNG)

	run := domain.NewRunState()
	run.Player.CurrentPortID = "mars"
	run.MarketsByPortID = markets
	run.Cargo.SetQuantity("synthspice", 2)
	run.Cargo.SetQuantity("alloy", 1)

	market := run.MarketsByPortID["mars"]
	expected := (2 * market.PricesByItemID["synthspice"]) + market.PricesByItemID["alloy"]
	if value := economy.GetSellableCargoValueAtCurrentPort(run, snapshot); value != expected {
		t.Fatalf("expected sellable cargo value %d, got %d", expected, value)
	}
}

func expectedZoneRarity(zone domain.PortZone) domain.ItemRarity {
	switch zone {
	case domain.PortZoneInner:
		return domain.ItemRarityCommon
	case domain.PortZoneOuter:
		return domain.ItemRarityMidTier
	case domain.PortZoneFringe:
		return domain.ItemRarityExotic
	default:
		return domain.ItemRarityCommon
	}
}
