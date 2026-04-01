using System.Collections.Generic;

namespace StarSmugglerGo.Persistence;

public sealed class MarketSnapshotSaveData
{
    public string PortId { get; set; } = string.Empty;
    public List<string> AvailableItemIds { get; set; } = new();
    public Dictionary<string, int> PricesByItemId { get; set; } = new();
}
