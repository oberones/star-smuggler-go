package services

import (
	"sort"

	"github.com/oberones/star-smuggler-go/internal/domain"
)

const (
	routePressureCap     = 4
	commodityPressureCap = 5
)

type EconomyBalanceService struct{}

func (s EconomyBalanceService) RecordRoute(run *domain.RunState, originID string, destinationID string) {
	if run.RoutePressureByKey == nil {
		run.RoutePressureByKey = make(map[string]int)
	}

	key := canonicalRouteKey(originID, destinationID)
	for existingKey, pressure := range run.RoutePressureByKey {
		if existingKey == key {
			continue
		}
		if pressure <= 1 {
			delete(run.RoutePressureByKey, existingKey)
		} else {
			run.RoutePressureByKey[existingKey] = pressure - 1
		}
	}

	next := run.RoutePressureByKey[key] + 1
	if next > routePressureCap {
		next = routePressureCap
	}
	run.RoutePressureByKey[key] = next
}

func (s EconomyBalanceService) RecordCommodityTrade(run *domain.RunState, itemID string, quantity int) {
	if run.CommodityPressureByItemID == nil {
		run.CommodityPressureByItemID = make(map[string]int)
	}

	for existingID, pressure := range run.CommodityPressureByItemID {
		if existingID == itemID {
			continue
		}
		if pressure <= 1 {
			delete(run.CommodityPressureByItemID, existingID)
		} else {
			run.CommodityPressureByItemID[existingID] = pressure - 1
		}
	}

	increase := 1
	if quantity >= 4 {
		increase = 2
	}

	next := run.CommodityPressureByItemID[itemID] + increase
	if next > commodityPressureCap {
		next = commodityPressureCap
	}
	run.CommodityPressureByItemID[itemID] = next
}

func (s EconomyBalanceService) AdditionalRouteCost(run domain.RunState, originID string, destinationID string) int {
	key := canonicalRouteKey(originID, destinationID)
	return run.RoutePressureByKey[key] * 2
}

func (s EconomyBalanceService) ApplyMarketPressure(run *domain.RunState, data domain.DataSnapshot) {
	for portID, market := range run.MarketsByPortID {
		prices := clonePriceMap(market.PricesByItemID)
		for itemID, pressure := range run.CommodityPressureByItemID {
			if pressure <= 0 {
				continue
			}

			item, ok := data.ItemsByID[itemID]
			if !ok {
				continue
			}

			currentPrice, ok := prices[itemID]
			if !ok {
				continue
			}

			prices[itemID] = compressPriceTowardBase(currentPrice, item.BasePrice, pressure)
		}

		market.PricesByItemID = prices
		run.MarketsByPortID[portID] = market
	}
}

func canonicalRouteKey(left string, right string) string {
	parts := []string{left, right}
	sort.Strings(parts)
	return parts[0] + "->" + parts[1]
}

func compressPriceTowardBase(current int, base int, pressure int) int {
	if current == base || pressure <= 0 {
		return current
	}

	factor := 1.0 - (0.15 * float64(pressure))
	if factor < 0.25 {
		factor = 0.25
	}

	delta := current - base
	adjusted := base + int(float64(delta)*factor)
	if adjusted == current {
		if delta > 0 {
			adjusted--
		} else {
			adjusted++
		}
	}

	if adjusted < 1 {
		return 1
	}
	return adjusted
}

func clonePriceMap(source map[string]int) map[string]int {
	if source == nil {
		return map[string]int{}
	}

	clone := make(map[string]int, len(source))
	for key, value := range source {
		clone[key] = value
	}
	return clone
}
