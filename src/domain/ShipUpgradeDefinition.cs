using System.Collections.Generic;

namespace StarSmugglerGo.Domain;

public sealed class ShipUpgradeDefinition
{
    public string Id { get; init; } = string.Empty;
    public string Name { get; init; } = string.Empty;
    public string Description { get; init; } = string.Empty;
    public UpgradeCategory Category { get; init; }
    public int CostCredits { get; init; }
    public string RequiredFactionId { get; init; } = string.Empty;
    public string MinimumStanding { get; init; } = string.Empty;
    public ShipSpecialization? Specialization { get; init; }
    public IReadOnlyList<UpgradeEffectDefinition> Effects { get; init; } = new List<UpgradeEffectDefinition>();
}
