using System.Collections.Generic;

namespace StarSmugglerGo.Application;

public sealed class TradeScreenViewModel
{
    public string PortName { get; init; } = string.Empty;
    public int Credits { get; init; }
    public int CargoLoad { get; init; }
    public int CargoLimit { get; init; }
    public IReadOnlyList<TradeItemViewModel> Items { get; init; } = new List<TradeItemViewModel>();
    public string StatusMessage { get; init; } = string.Empty;
}
