using StarSmugglerGo.Domain;
using System;
using System.Collections.Generic;
using System.Linq;

namespace StarSmugglerGo.Services;

public sealed class EventService
{
    public EventResult? TryResolveTravelEvent(RunState run, DataSnapshot data, EconomyService economyService, Random random)
    {
        int roll = random.Next(1, 101);
        if (roll > 30 || data.Events.Count == 0)
        {
            return null;
        }

        EventDefinition definition = SelectWeightedEvent(data.Events, random);
        MarketSnapshot? currentMarket = economyService.GetCurrentMarket(run);

        return definition.EffectType switch
        {
            EventEffectType.CreditsLossScaled => ResolveCreditsLoss(run, definition, random),
            EventEffectType.PortPriceMultiplier => ResolvePortPriceMultiplier(currentMarket, definition),
            EventEffectType.SingleItemPriceMultiplier => ResolveSingleItemPriceMultiplier(run, data, currentMarket, definition, random),
            EventEffectType.LoseRandomCargo => ResolveLoseRandomCargo(run, data, definition, random),
            _ => null,
        };
    }

    private static EventDefinition SelectWeightedEvent(IReadOnlyList<EventDefinition> events, Random random)
    {
        int totalWeight = events.Sum(evt => Math.Max(1, evt.Weight));
        int roll = random.Next(totalWeight);
        int cumulative = 0;

        foreach (EventDefinition definition in events)
        {
            cumulative += Math.Max(1, definition.Weight);
            if (roll < cumulative)
            {
                return definition;
            }
        }

        return events[^1];
    }

    private static EventResult ResolveCreditsLoss(RunState run, EventDefinition definition, Random random)
    {
        int minimum = (int)GetParameter(definition, "minimum", 25);
        int maximum = (int)GetParameter(definition, "maximum", 100);
        double minPercent = GetParameter(definition, "minPercent", 0.05);
        double maxPercent = GetParameter(definition, "maxPercent", 0.15);

        int baseLoss = random.Next(minimum, maximum + 1);
        double percentage = minPercent + (random.NextDouble() * (maxPercent - minPercent));
        int scaledLoss = Math.Max(baseLoss, (int)(run.Player.Credits * percentage));

        run.Player.Credits = Math.Max(0, run.Player.Credits - scaledLoss);

        return new EventResult
        {
            EventId = definition.Id,
            Name = definition.Name,
            ResolvedDescription = definition.DescriptionTemplate.Replace("{credits}", scaledLoss.ToString()),
            RolledValues = new Dictionary<string, double>
            {
                ["credits"] = scaledLoss,
                ["percentage"] = percentage,
            },
        };
    }

    private static EventResult ResolvePortPriceMultiplier(MarketSnapshot? market, EventDefinition definition)
    {
        double multiplier = GetParameter(definition, "multiplier", 2.0);

        if (market is not null)
        {
            foreach (string itemId in market.AvailableItemIds)
            {
                if (market.PricesByItemId.TryGetValue(itemId, out int currentPrice))
                {
                    market.PricesByItemId[itemId] = Math.Max(1, (int)(currentPrice * multiplier));
                }
            }
        }

        return new EventResult
        {
            EventId = definition.Id,
            Name = definition.Name,
            ResolvedDescription = definition.DescriptionTemplate,
            RolledValues = new Dictionary<string, double>
            {
                ["multiplier"] = multiplier,
            },
        };
    }

    private static EventResult ResolveSingleItemPriceMultiplier(
        RunState run,
        DataSnapshot data,
        MarketSnapshot? market,
        EventDefinition definition,
        Random random)
    {
        double multiplier = GetParameter(definition, "multiplier", 0.5);
        string itemName = "an item";

        if (market is not null && market.AvailableItemIds.Count > 0)
        {
            string itemId = market.AvailableItemIds[random.Next(market.AvailableItemIds.Count)];
            itemName = data.ItemsById.TryGetValue(itemId, out ItemDefinition? item) ? item.Name : itemId;

            if (market.PricesByItemId.TryGetValue(itemId, out int currentPrice))
            {
                market.PricesByItemId[itemId] = Math.Max(1, (int)(currentPrice * multiplier));
            }
        }

        return new EventResult
        {
            EventId = definition.Id,
            Name = definition.Name,
            ResolvedDescription = definition.DescriptionTemplate.Replace("{itemName}", itemName),
            RolledValues = new Dictionary<string, double>
            {
                ["multiplier"] = multiplier,
            },
        };
    }

    private static EventResult ResolveLoseRandomCargo(RunState run, DataSnapshot data, EventDefinition definition, Random random)
    {
        string itemName = "nothing";
        List<string> cargoItems = run.Cargo.ItemQuantities
            .Where(pair => pair.Value > 0)
            .Select(pair => pair.Key)
            .ToList();

        if (cargoItems.Count > 0)
        {
            string itemId = cargoItems[random.Next(cargoItems.Count)];
            int owned = run.Cargo.GetQuantity(itemId);
            run.Cargo.SetQuantity(itemId, owned - 1);
            itemName = data.ItemsById.TryGetValue(itemId, out ItemDefinition? item) ? item.Name : itemId;
        }

        return new EventResult
        {
            EventId = definition.Id,
            Name = definition.Name,
            ResolvedDescription = definition.DescriptionTemplate.Replace("{itemName}", itemName),
            RolledValues = new Dictionary<string, double>(),
        };
    }

    private static double GetParameter(EventDefinition definition, string key, double fallback)
    {
        return definition.Parameters.TryGetValue(key, out double value) ? value : fallback;
    }
}
