using System.Collections.Generic;

namespace StarSmugglerGo.Domain;

public sealed class MarketSnapshot
{
    public string PortId { get; init; } = string.Empty;
    public List<string> AvailableItemIds { get; init; } = new();
    public Dictionary<string, int> PricesByItemId { get; init; } = new();
}
