using StarSmugglerGo.Domain;
using System;
using System.Collections.Generic;
using System.Linq;

namespace StarSmugglerGo.Services;

public sealed class EconomyService
{
    public Dictionary<string, MarketSnapshot> CreateInitialMarkets(DataSnapshot data, Random random)
    {
        var result = new Dictionary<string, MarketSnapshot>(StringComparer.Ordinal);

        foreach (PortDefinition port in data.Ports)
        {
            result[port.Id] = CreateMarketForPort(port, data.Items, random);
        }

        return result;
    }

    public MarketSnapshot CreateMarketForPort(PortDefinition port, IReadOnlyList<ItemDefinition> items, Random random)
    {
        List<ItemDefinition> availableItems = BuildAvailableItemsForPort(items, port.Zone, random);
        var priceMap = new Dictionary<string, int>(StringComparer.Ordinal);

        foreach (ItemDefinition item in items)
        {
            priceMap[item.Id] = CalculatePrice(item, port.Zone, random);
        }

        return new MarketSnapshot
        {
            PortId = port.Id,
            AvailableItemIds = availableItems.Select(item => item.Id).ToList(),
            PricesByItemId = priceMap,
        };
    }

    public void RefreshAvailableGoods(RunState run, DataSnapshot data, string portId, Random random)
    {
        if (!data.PortsById.TryGetValue(portId, out PortDefinition? port))
        {
            return;
        }

        if (!run.MarketsByPortId.TryGetValue(portId, out MarketSnapshot? market))
        {
            run.MarketsByPortId[portId] = CreateMarketForPort(port, data.Items, random);
            return;
        }

        List<ItemDefinition> availableItems = BuildAvailableItemsForPort(data.Items, port.Zone, random);
        market.AvailableItemIds.Clear();
        market.AvailableItemIds.AddRange(availableItems.Select(item => item.Id));
    }

    public void RefreshAllPrices(RunState run, DataSnapshot data, Random random)
    {
        run.MarketsByPortId.Clear();

        foreach (KeyValuePair<string, MarketSnapshot> pair in CreateInitialMarkets(data, random))
        {
            run.MarketsByPortId[pair.Key] = pair.Value;
        }

        run.JumpsSinceLastUpdate = 0;
    }

    public MarketSnapshot? GetCurrentMarket(RunState run)
    {
        return run.MarketsByPortId.TryGetValue(run.Player.CurrentPortId, out MarketSnapshot? market)
            ? market
            : null;
    }

    public int GetCargoLoad(RunState run)
    {
        return run.Cargo.TotalUnits;
    }

    public IReadOnlyList<ItemDefinition> GetAvailableGoodsForCurrentPort(RunState run, DataSnapshot data)
    {
        MarketSnapshot? market = GetCurrentMarket(run);
        if (market is null)
        {
            return Array.Empty<ItemDefinition>();
        }

        List<ItemDefinition> goods = market.AvailableItemIds
            .Select(itemId => data.ItemsById.TryGetValue(itemId, out ItemDefinition? item) ? item : null)
            .Where(item => item is not null)
            .Cast<ItemDefinition>()
            .ToList();

        return goods;
    }

    public int GetSellableCargoValueAtCurrentPort(RunState run, DataSnapshot data)
    {
        MarketSnapshot? market = GetCurrentMarket(run);
        if (market is null)
        {
            return 0;
        }

        int total = 0;

        foreach ((string itemId, int quantity) in run.Cargo.ItemQuantities)
        {
            if (quantity <= 0)
            {
                continue;
            }

            if (market.PricesByItemId.TryGetValue(itemId, out int price))
            {
                total += quantity * price;
            }
            else if (data.ItemsById.TryGetValue(itemId, out ItemDefinition? item))
            {
                total += quantity * item.BasePrice;
            }
        }

        return total;
    }

    private static List<ItemDefinition> BuildAvailableItemsForPort(
        IReadOnlyList<ItemDefinition> allItems,
        PortZone zone,
        Random random)
    {
        List<ItemDefinition> zoneItems = allItems.Where(item => RarityForZone(zone) == item.Rarity).ToList();
        List<ItemDefinition> otherItems = allItems.Where(item => RarityForZone(zone) != item.Rarity).ToList();

        return PickDistinct(zoneItems, 4, random)
            .Concat(PickDistinct(otherItems, 2, random))
            .OrderBy(item => item.Name, StringComparer.Ordinal)
            .ToList();
    }

    private static IEnumerable<ItemDefinition> PickDistinct(List<ItemDefinition> candidates, int count, Random random)
    {
        List<ItemDefinition> remaining = new(candidates);
        List<ItemDefinition> selected = new();

        while (selected.Count < count && remaining.Count > 0)
        {
            int index = random.Next(remaining.Count);
            selected.Add(remaining[index]);
            remaining.RemoveAt(index);
        }

        return selected;
    }

    private static int CalculatePrice(ItemDefinition item, PortZone zone, Random random)
    {
        double variance = 0.3;
        double markup = GetItemMarkup(item.Rarity, zone);
        double multiplier = 1.0 + ((random.NextDouble() * 2.0 - 1.0) * variance) + markup;
        return Math.Max(1, (int)(item.BasePrice * multiplier));
    }

    private static ItemRarity RarityForZone(PortZone zone)
    {
        return zone switch
        {
            PortZone.Inner => ItemRarity.Common,
            PortZone.Outer => ItemRarity.MidTier,
            PortZone.Fringe => ItemRarity.Exotic,
            _ => ItemRarity.Common,
        };
    }

    private static double GetItemMarkup(ItemRarity rarity, PortZone zone)
    {
        return rarity switch
        {
            ItemRarity.Common when zone == PortZone.Fringe => 2.0,
            ItemRarity.Common when zone == PortZone.Outer => 0.5,
            ItemRarity.Common => 0.0,
            ItemRarity.MidTier when zone == PortZone.Fringe => 1.0,
            ItemRarity.MidTier when zone == PortZone.Inner => 0.5,
            ItemRarity.MidTier => 0.0,
            ItemRarity.Exotic when zone == PortZone.Inner => 2.0,
            ItemRarity.Exotic when zone == PortZone.Outer => 0.5,
            ItemRarity.Exotic => 0.0,
            _ => 0.0,
        };
    }
}
