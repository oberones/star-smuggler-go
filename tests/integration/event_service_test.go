package integration_test

import (
	"testing"

	"github.com/oberones/star-smuggler-go/internal/domain"
	"github.com/oberones/star-smuggler-go/internal/services"
)

func TestEventServiceSelectsDeterministicWeightedTravelEvent(t *testing.T) {
	t.Parallel()

	ports := []domain.PortDefinition{{ID: "mars", Name: "Mars", Zone: domain.PortZoneInner}}
	items := []domain.ItemDefinition{{ID: "ore", Name: "Ore", Rarity: domain.ItemRarityCommon, BasePrice: 20}}
	events := []domain.EventDefinition{
		{
			ID:                  "light_fee",
			Name:                "Light Fee",
			DescriptionTemplate: "A port clerk takes {credits} credits.",
			EffectType:          domain.EventEffectCreditsLossScaled,
			Weight:              1,
			Parameters: map[string]float64{
				"minimum":    10,
				"maximum":    10,
				"minPercent": 0.01,
				"maxPercent": 0.01,
			},
		},
		{
			ID:                  "heavy_fee",
			Name:                "Heavy Fee",
			DescriptionTemplate: "Pirates take {credits} credits.",
			EffectType:          domain.EventEffectCreditsLossScaled,
			Weight:              3,
			Parameters: map[string]float64{
				"minimum":    20,
				"maximum":    20,
				"minPercent": 0.02,
				"maxPercent": 0.02,
			},
		},
	}

	snapshot := domain.NewDataSnapshot(ports, items, events)
	run := domain.NewRunState()
	run.Player.CurrentPortID = "mars"
	run.MarketsByPortID["mars"] = domain.MarketSnapshot{
		PortID:           "mars",
		AvailableItemIDs: []string{"ore"},
		PricesByItemID:   map[string]int{"ore": 20},
	}

	rng := &stubRNG{
		intValues:   []int{0, 3, 0},
		floatValues: []float64{0},
	}

	result := services.EventService{}.TryResolveTravelEvent(&run, snapshot, services.EconomyService{}, rng)
	if result == nil {
		t.Fatal("expected a travel event to resolve")
	}
	if result.EventID != "heavy_fee" {
		t.Fatalf("expected weighted event heavy_fee, got %q", result.EventID)
	}
}

func TestEventServiceCanSkipTravelEventRoll(t *testing.T) {
	t.Parallel()

	snapshot := loadSnapshot(t)
	run := domain.NewRunState()
	run.Player.CurrentPortID = snapshot.Ports[0].ID
	run.MarketsByPortID[run.Player.CurrentPortID] = domain.MarketSnapshot{
		PortID:           run.Player.CurrentPortID,
		AvailableItemIDs: []string{},
		PricesByItemID:   map[string]int{},
	}

	rng := &stubRNG{intValues: []int{99}}
	result := services.EventService{}.TryResolveTravelEvent(&run, snapshot, services.EconomyService{}, rng)
	if result != nil {
		t.Fatalf("expected no travel event when the chance roll misses, got %#v", result)
	}
}
