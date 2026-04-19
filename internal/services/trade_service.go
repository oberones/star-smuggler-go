package services

import (
	"fmt"

	"github.com/oberones/star-smuggler-go/internal/domain"
)

type TradeService struct{}

func (s TradeService) Buy(run *domain.RunState, market domain.MarketSnapshot, item domain.ItemDefinition, quantity int) TradeResult {
	if quantity <= 0 {
		return FailedTrade("Quantity must be at least 1.")
	}

	if !containsString(market.AvailableItemIDs, item.ID) {
		return FailedTrade(fmt.Sprintf("%s is not available at this port.", item.Name))
	}

	price, ok := market.PricesByItemID[item.ID]
	if !ok {
		return FailedTrade(fmt.Sprintf("No current market price exists for %s.", item.Name))
	}

	totalCost := price * quantity
	if run.Player.Credits < totalCost {
		return FailedTrade(fmt.Sprintf("You need %d credits to buy %d %s.", totalCost, quantity, item.Name))
	}

	projectedCargo := run.Cargo.TotalUnits() + quantity
	if projectedCargo > run.Player.CargoLimit {
		return FailedTrade("Your cargo hold cannot fit that purchase.")
	}

	run.Player.Credits -= totalCost
	run.Cargo.SetQuantity(item.ID, run.Cargo.QuantityFor(item.ID)+quantity)
	return SuccessfulTrade(fmt.Sprintf("Bought %d %s for %d credits.", quantity, item.Name, totalCost))
}

func (s TradeService) Sell(run *domain.RunState, market domain.MarketSnapshot, item domain.ItemDefinition, quantity int) TradeResult {
	if quantity <= 0 {
		return FailedTrade("Quantity must be at least 1.")
	}

	price, ok := market.PricesByItemID[item.ID]
	if !ok {
		return FailedTrade(fmt.Sprintf("No current market price exists for %s.", item.Name))
	}

	owned := run.Cargo.QuantityFor(item.ID)
	if owned < quantity {
		return FailedTrade(fmt.Sprintf("You only own %d %s.", owned, item.Name))
	}

	totalValue := price * quantity
	run.Player.Credits += totalValue
	run.Cargo.SetQuantity(item.ID, owned-quantity)
	return SuccessfulTrade(fmt.Sprintf("Sold %d %s for %d credits.", quantity, item.Name, totalValue))
}

func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
