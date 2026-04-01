using System.Collections.Generic;

namespace StarSmugglerGo.Application;

public sealed class TravelScreenViewModel
{
    public string CurrentPortName { get; init; } = string.Empty;
    public string BackgroundTexturePath { get; init; } = string.Empty;
    public int Credits { get; init; }
    public IReadOnlyList<TravelDestinationViewModel> Destinations { get; init; } = new List<TravelDestinationViewModel>();
    public string StatusMessage { get; init; } = string.Empty;
}
