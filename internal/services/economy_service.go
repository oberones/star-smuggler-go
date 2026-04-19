package services

import (
	"sort"

	"github.com/oberones/star-smuggler-go/internal/domain"
)

type EconomyService struct{}

func (s EconomyService) CreateInitialMarkets(data domain.DataSnapshot, rng RNG) map[string]domain.MarketSnapshot {
	result := make(map[string]domain.MarketSnapshot, len(data.Ports))
	for _, port := range data.Ports {
		result[port.ID] = s.CreateMarketForPort(port, data.Items, rng)
	}
	return result
}

func (s EconomyService) CreateMarketForPort(port domain.PortDefinition, items []domain.ItemDefinition, rng RNG) domain.MarketSnapshot {
	availableItems := buildAvailableItemsForPort(items, port.Zone, rng)
	priceMap := make(map[string]int, len(items))

	for _, item := range items {
		priceMap[item.ID] = calculatePrice(item, port.Zone, rng)
	}

	return domain.MarketSnapshot{
		PortID:           port.ID,
		AvailableItemIDs: availableItems,
		PricesByItemID:   priceMap,
	}
}

func (s EconomyService) RefreshAvailableGoods(run *domain.RunState, data domain.DataSnapshot, portID string, rng RNG) {
	port, ok := data.PortsByID[portID]
	if !ok {
		return
	}

	market, ok := run.MarketsByPortID[portID]
	if !ok {
		run.MarketsByPortID[portID] = s.CreateMarketForPort(port, data.Items, rng)
		return
	}

	market.AvailableItemIDs = buildAvailableItemsForPort(data.Items, port.Zone, rng)
	run.MarketsByPortID[portID] = market
}

func (s EconomyService) RefreshAllPrices(run *domain.RunState, data domain.DataSnapshot, rng RNG) {
	run.MarketsByPortID = s.CreateInitialMarkets(data, rng)
	run.JumpsSinceLastUpdate = 0
}

func (s EconomyService) GetCurrentMarket(run domain.RunState) (domain.MarketSnapshot, bool) {
	market, ok := run.MarketsByPortID[run.Player.CurrentPortID]
	return market, ok
}

func (s EconomyService) GetCargoLoad(run domain.RunState) int {
	return run.Cargo.TotalUnits()
}

func (s EconomyService) GetAvailableGoodsForCurrentPort(run domain.RunState, data domain.DataSnapshot) []domain.ItemDefinition {
	market, ok := s.GetCurrentMarket(run)
	if !ok {
		return nil
	}

	goods := make([]domain.ItemDefinition, 0, len(market.AvailableItemIDs))
	for _, itemID := range market.AvailableItemIDs {
		if item, exists := data.ItemsByID[itemID]; exists {
			goods = append(goods, item)
		}
	}
	return goods
}

func (s EconomyService) GetSellableCargoValueAtCurrentPort(run domain.RunState, data domain.DataSnapshot) int {
	market, ok := s.GetCurrentMarket(run)
	if !ok {
		return 0
	}

	total := 0
	for itemID, quantity := range run.Cargo.ItemQuantities {
		if quantity <= 0 {
			continue
		}

		if price, exists := market.PricesByItemID[itemID]; exists {
			total += quantity * price
			continue
		}

		if item, exists := data.ItemsByID[itemID]; exists {
			total += quantity * item.BasePrice
		}
	}

	return total
}

func buildAvailableItemsForPort(allItems []domain.ItemDefinition, zone domain.PortZone, rng RNG) []string {
	zoneRarity := rarityForZone(zone)
	zoneItems := make([]domain.ItemDefinition, 0)
	otherItems := make([]domain.ItemDefinition, 0)

	for _, item := range allItems {
		if item.Rarity == zoneRarity {
			zoneItems = append(zoneItems, item)
		} else {
			otherItems = append(otherItems, item)
		}
	}

	selected := append(pickDistinct(zoneItems, 4, rng), pickDistinct(otherItems, 2, rng)...)
	sort.Slice(selected, func(i, j int) bool {
		return selected[i].Name < selected[j].Name
	})

	ids := make([]string, 0, len(selected))
	for _, item := range selected {
		ids = append(ids, item.ID)
	}
	return ids
}

func pickDistinct(candidates []domain.ItemDefinition, count int, rng RNG) []domain.ItemDefinition {
	remaining := append([]domain.ItemDefinition(nil), candidates...)
	selected := make([]domain.ItemDefinition, 0, count)

	for len(selected) < count && len(remaining) > 0 {
		index := rng.Intn(len(remaining))
		selected = append(selected, remaining[index])
		remaining = append(remaining[:index], remaining[index+1:]...)
	}

	return selected
}

func calculatePrice(item domain.ItemDefinition, zone domain.PortZone, rng RNG) int {
	const variance = 0.3
	markup := getItemMarkup(item.Rarity, zone)
	multiplier := 1.0 + ((rng.Float64()*2.0 - 1.0) * variance) + markup
	price := int(float64(item.BasePrice) * multiplier)
	if price < 1 {
		return 1
	}
	return price
}

func rarityForZone(zone domain.PortZone) domain.ItemRarity {
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

func getItemMarkup(rarity domain.ItemRarity, zone domain.PortZone) float64 {
	switch {
	case rarity == domain.ItemRarityCommon && zone == domain.PortZoneFringe:
		return 2.0
	case rarity == domain.ItemRarityCommon && zone == domain.PortZoneOuter:
		return 0.5
	case rarity == domain.ItemRarityMidTier && zone == domain.PortZoneFringe:
		return 1.0
	case rarity == domain.ItemRarityMidTier && zone == domain.PortZoneInner:
		return 0.5
	case rarity == domain.ItemRarityExotic && zone == domain.PortZoneInner:
		return 2.0
	case rarity == domain.ItemRarityExotic && zone == domain.PortZoneOuter:
		return 0.5
	default:
		return 0.0
	}
}
