using System.Collections.Generic;

namespace StarSmugglerGo.Application;

public sealed class TradeScreenViewModel
{
    public string PortName { get; init; } = string.Empty;
    public string BackgroundTexturePath { get; init; } = string.Empty;
    public string MusicTrackId { get; init; } = string.Empty;
    public int Credits { get; init; }
    public int CargoLoad { get; init; }
    public int CargoLimit { get; init; }
    public IReadOnlyList<TradeItemViewModel> Items { get; init; } = new List<TradeItemViewModel>();
    public IReadOnlyList<UpgradeOptionViewModel> Upgrades { get; init; } = new List<UpgradeOptionViewModel>();
    public string StatusMessage { get; init; } = string.Empty;
}
