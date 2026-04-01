namespace StarSmugglerGo.Domain;

public sealed class PortDefinition
{
    public string Id { get; init; } = string.Empty;
    public string Name { get; init; } = string.Empty;
    public string Description { get; init; } = string.Empty;
    public PortZone Zone { get; init; }
    public string BackgroundTexturePath { get; init; } = string.Empty;
    public string PreviewTexturePath { get; init; } = string.Empty;
    public string TradeBackgroundPath { get; init; } = string.Empty;
    public string MusicTrackId { get; init; } = string.Empty;
}
