using System.Collections.Generic;

namespace StarSmugglerGo.Application;

public sealed class PortOverviewViewModel
{
    public string PortName { get; init; } = string.Empty;
    public string PortDescription { get; init; } = string.Empty;
    public string ZoneName { get; init; } = string.Empty;
    public int Credits { get; init; }
    public int CargoLoad { get; init; }
    public int CargoLimit { get; init; }
    public int CheapestTravelCost { get; init; }
    public bool IsGameOver { get; init; }
    public string RecentEventText { get; init; } = string.Empty;
    public IReadOnlyList<string> AvailableGoods { get; init; } = new List<string>();
}
