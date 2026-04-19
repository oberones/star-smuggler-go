package services

import (
	"fmt"
	"strings"

	"github.com/oberones/star-smuggler-go/internal/domain"
)

const travelEventChancePercent = 30

type EventService struct{}

func (s EventService) TryResolveTravelEvent(run *domain.RunState, data domain.DataSnapshot, economy EconomyService, rng RNG) *domain.EventResult {
	if rng == nil || len(data.Events) == 0 {
		return nil
	}

	if rng.Intn(100)+1 > travelEventChancePercent {
		return nil
	}

	definition := selectWeightedEvent(data.Events, rng)
	market, _ := economy.GetCurrentMarket(*run)

	switch definition.EffectType {
	case domain.EventEffectCreditsLossScaled:
		return resolveCreditsLoss(run, definition, rng)
	case domain.EventEffectPortPriceMultiplier:
		return resolvePortPriceMultiplier(market, definition)
	case domain.EventEffectSingleItemMultiplier:
		return resolveSingleItemPriceMultiplier(run, data, market, definition, rng)
	case domain.EventEffectLoseRandomCargo:
		return resolveLoseRandomCargo(run, data, definition, rng)
	default:
		return nil
	}
}

func selectWeightedEvent(events []domain.EventDefinition, rng RNG) domain.EventDefinition {
	totalWeight := 0
	for _, event := range events {
		weight := event.Weight
		if weight < 1 {
			weight = 1
		}
		totalWeight += weight
	}

	roll := rng.Intn(totalWeight)
	cumulative := 0
	for _, event := range events {
		weight := event.Weight
		if weight < 1 {
			weight = 1
		}
		cumulative += weight
		if roll < cumulative {
			return event
		}
	}

	return events[len(events)-1]
}

func resolveCreditsLoss(run *domain.RunState, definition domain.EventDefinition, rng RNG) *domain.EventResult {
	minimum := int(getEventParameter(definition, "minimum", 25))
	maximum := int(getEventParameter(definition, "maximum", 100))
	minPercent := getEventParameter(definition, "minPercent", 0.05)
	maxPercent := getEventParameter(definition, "maxPercent", 0.15)

	baseLoss := minimum
	if maximum > minimum {
		baseLoss += rng.Intn((maximum - minimum) + 1)
	}

	percentage := minPercent
	if maxPercent > minPercent {
		percentage += rng.Float64() * (maxPercent - minPercent)
	}

	scaledLoss := maxInt(baseLoss, int(float64(run.Player.Credits)*percentage))
	run.Player.Credits -= scaledLoss
	if run.Player.Credits < 0 {
		run.Player.Credits = 0
	}

	return &domain.EventResult{
		EventID:             definition.ID,
		Name:                definition.Name,
		ResolvedDescription: replacePlaceholder(definition.DescriptionTemplate, "{credits}", fmt.Sprintf("%d", scaledLoss)),
		RolledValues: map[string]float64{
			"credits":    float64(scaledLoss),
			"percentage": percentage,
		},
	}
}

func resolvePortPriceMultiplier(market domain.MarketSnapshot, definition domain.EventDefinition) *domain.EventResult {
	multiplier := getEventParameter(definition, "multiplier", 2.0)
	if market.PricesByItemID != nil {
		for _, itemID := range market.AvailableItemIDs {
			if currentPrice, ok := market.PricesByItemID[itemID]; ok {
				market.PricesByItemID[itemID] = maxInt(1, int(float64(currentPrice)*multiplier))
			}
		}
	}

	return &domain.EventResult{
		EventID:             definition.ID,
		Name:                definition.Name,
		ResolvedDescription: definition.DescriptionTemplate,
		RolledValues: map[string]float64{
			"multiplier": multiplier,
		},
	}
}

func resolveSingleItemPriceMultiplier(_ *domain.RunState, data domain.DataSnapshot, market domain.MarketSnapshot, definition domain.EventDefinition, rng RNG) *domain.EventResult {
	multiplier := getEventParameter(definition, "multiplier", 0.5)
	itemName := "an item"

	if len(market.AvailableItemIDs) > 0 {
		itemID := market.AvailableItemIDs[rng.Intn(len(market.AvailableItemIDs))]
		if item, ok := data.ItemsByID[itemID]; ok {
			itemName = item.Name
		} else {
			itemName = itemID
		}

		if currentPrice, ok := market.PricesByItemID[itemID]; ok {
			market.PricesByItemID[itemID] = maxInt(1, int(float64(currentPrice)*multiplier))
		}
	}

	return &domain.EventResult{
		EventID:             definition.ID,
		Name:                definition.Name,
		ResolvedDescription: replacePlaceholder(definition.DescriptionTemplate, "{itemName}", itemName),
		RolledValues: map[string]float64{
			"multiplier": multiplier,
		},
	}
}

func resolveLoseRandomCargo(run *domain.RunState, data domain.DataSnapshot, definition domain.EventDefinition, rng RNG) *domain.EventResult {
	itemName := "nothing"
	cargoItems := make([]string, 0, len(run.Cargo.ItemQuantities))
	for itemID, quantity := range run.Cargo.ItemQuantities {
		if quantity > 0 {
			cargoItems = append(cargoItems, itemID)
		}
	}

	if len(cargoItems) > 0 {
		itemID := cargoItems[rng.Intn(len(cargoItems))]
		run.Cargo.SetQuantity(itemID, run.Cargo.QuantityFor(itemID)-1)

		if item, ok := data.ItemsByID[itemID]; ok {
			itemName = item.Name
		} else {
			itemName = itemID
		}
	}

	return &domain.EventResult{
		EventID:             definition.ID,
		Name:                definition.Name,
		ResolvedDescription: replacePlaceholder(definition.DescriptionTemplate, "{itemName}", itemName),
		RolledValues:        map[string]float64{},
	}
}

func getEventParameter(definition domain.EventDefinition, key string, fallback float64) float64 {
	if value, ok := definition.Parameters[key]; ok {
		return value
	}

	return fallback
}

func replacePlaceholder(value string, placeholder string, replacement string) string {
	return strings.ReplaceAll(value, placeholder, replacement)
}

func maxInt(left int, right int) int {
	if left > right {
		return left
	}
	return right
}
