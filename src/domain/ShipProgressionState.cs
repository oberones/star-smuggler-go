using System.Collections.Generic;
using System.Linq;

namespace StarSmugglerGo.Domain;

public sealed class ShipProgressionState
{
    public List<string> PurchasedUpgradeIds { get; init; } = new();
    public Dictionary<string, bool> SpecializationFlags { get; init; } = new();

    public bool HasUpgrade(string upgradeId) => PurchasedUpgradeIds.Contains(upgradeId);

    public bool HasSpecialization(ShipSpecialization specialization)
        => SpecializationFlags.TryGetValue(specialization.ToString(), out bool isActive) && isActive;
}
