namespace StarSmugglerGo.Application;

public sealed class TravelAnimationViewModel
{
    public string OriginName { get; init; } = string.Empty;
    public string DestinationName { get; init; } = string.Empty;
    public string BackgroundTexturePath { get; init; } = string.Empty;
    public int TravelCost { get; init; }
    public double DurationSeconds { get; init; } = 2.5;
    public string StatusMessage { get; init; } = string.Empty;
}
