namespace StarSmugglerGo.Application;

public sealed class UpgradeOptionViewModel
{
    public string UpgradeId { get; init; } = string.Empty;
    public string Name { get; init; } = string.Empty;
    public string Description { get; init; } = string.Empty;
    public string EffectSummary { get; init; } = string.Empty;
    public string AvailabilityText { get; init; } = string.Empty;
    public int CostCredits { get; init; }
    public bool CanPurchase { get; init; }
    public bool IsInstalled { get; init; }
}
