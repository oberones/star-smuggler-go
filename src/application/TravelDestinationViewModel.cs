namespace StarSmugglerGo.Application;

public sealed class TravelDestinationViewModel
{
    public string PortId { get; init; } = string.Empty;
    public string Name { get; init; } = string.Empty;
    public string ZoneName { get; init; } = string.Empty;
    public string Description { get; init; } = string.Empty;
    public string PreviewTexturePath { get; init; } = string.Empty;
    public int TravelCost { get; init; }
}
