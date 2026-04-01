using StarSmugglerGo.Domain;
using System;

namespace StarSmugglerGo.Services;

public sealed class TradeService
{
    public TradeResult Buy(RunState run, MarketSnapshot market, ItemDefinition item, int quantity)
    {
        if (quantity <= 0)
        {
            return TradeResult.Failure("Quantity must be at least 1.");
        }

        if (!market.AvailableItemIds.Contains(item.Id))
        {
            return TradeResult.Failure($"{item.Name} is not available at this port.");
        }

        if (!market.PricesByItemId.TryGetValue(item.Id, out int price))
        {
            return TradeResult.Failure($"No current market price exists for {item.Name}.");
        }

        int totalCost = price * quantity;
        if (run.Player.Credits < totalCost)
        {
            return TradeResult.Failure($"You need {totalCost} credits to buy {quantity} {item.Name}.");
        }

        int projectedCargo = run.Cargo.TotalUnits + quantity;
        if (projectedCargo > run.Player.CargoLimit)
        {
            return TradeResult.Failure("Your cargo hold cannot fit that purchase.");
        }

        run.Player.Credits -= totalCost;
        run.Cargo.SetQuantity(item.Id, run.Cargo.GetQuantity(item.Id) + quantity);

        return TradeResult.Success($"Bought {quantity} {item.Name} for {totalCost} credits.");
    }

    public TradeResult Sell(RunState run, MarketSnapshot market, ItemDefinition item, int quantity)
    {
        if (quantity <= 0)
        {
            return TradeResult.Failure("Quantity must be at least 1.");
        }

        if (!market.PricesByItemId.TryGetValue(item.Id, out int price))
        {
            return TradeResult.Failure($"No current market price exists for {item.Name}.");
        }

        int owned = run.Cargo.GetQuantity(item.Id);
        if (owned < quantity)
        {
            return TradeResult.Failure($"You only own {owned} {item.Name}.");
        }

        int totalValue = price * quantity;
        run.Player.Credits += totalValue;
        run.Cargo.SetQuantity(item.Id, owned - quantity);

        return TradeResult.Success($"Sold {quantity} {item.Name} for {totalValue} credits.");
    }
}
