namespace StarSmugglerGo.Application;

public sealed class PendingTravelRequest
{
    public string OriginPortId { get; init; } = string.Empty;
    public string DestinationPortId { get; init; } = string.Empty;
    public int TravelCost { get; init; }
    public double DurationSeconds { get; init; }
}
