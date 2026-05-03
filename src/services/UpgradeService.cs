using StarSmugglerGo.Domain;
using System.Collections.Generic;
using System.Linq;

namespace StarSmugglerGo.Services;

public sealed class UpgradeService
{
    public IReadOnlyList<ShipUpgradeDefinition> GetVisibleUpgrades(DataSnapshot data)
        => data.Upgrades
            .Where(upgrade => string.IsNullOrWhiteSpace(upgrade.RequiredFactionId))
            .OrderBy(upgrade => upgrade.CostCredits)
            .ThenBy(upgrade => upgrade.Name)
            .ToList();

    public bool CanPurchase(RunState run, ShipUpgradeDefinition upgrade)
    {
        if (run.Progression.HasUpgrade(upgrade.Id))
        {
            return false;
        }

        if (!string.IsNullOrWhiteSpace(upgrade.RequiredFactionId))
        {
            return false;
        }

        return run.Player.Credits >= upgrade.CostCredits;
    }

    public string DescribeAvailability(RunState run, ShipUpgradeDefinition upgrade)
    {
        if (run.Progression.HasUpgrade(upgrade.Id))
        {
            return "Installed";
        }

        if (!string.IsNullOrWhiteSpace(upgrade.RequiredFactionId))
        {
            return $"Locked: requires {upgrade.MinimumStanding} standing with {upgrade.RequiredFactionId}.";
        }

        if (run.Player.Credits < upgrade.CostCredits)
        {
            return $"Need {upgrade.CostCredits} credits.";
        }

        return "Ready to install.";
    }

    public UpgradePurchaseResult PurchaseUpgrade(RunState run, DataSnapshot data, string upgradeId)
    {
        if (!data.UpgradesById.TryGetValue(upgradeId, out ShipUpgradeDefinition? upgrade))
        {
            return UpgradePurchaseResult.Failure($"Upgrade '{upgradeId}' no longer exists.");
        }

        if (run.Progression.HasUpgrade(upgrade.Id))
        {
            return UpgradePurchaseResult.Failure($"{upgrade.Name} is already installed.");
        }

        if (!string.IsNullOrWhiteSpace(upgrade.RequiredFactionId))
        {
            return UpgradePurchaseResult.Failure(
                $"{upgrade.Name} requires {upgrade.MinimumStanding} standing with {upgrade.RequiredFactionId}.");
        }

        if (run.Player.Credits < upgrade.CostCredits)
        {
            return UpgradePurchaseResult.Failure($"You need {upgrade.CostCredits} credits to install {upgrade.Name}.");
        }

        run.Player.Credits -= upgrade.CostCredits;
        run.Progression.PurchasedUpgradeIds.Add(upgrade.Id);
        run.Progression.PurchasedUpgradeIds.Sort(System.StringComparer.Ordinal);

        if (upgrade.Specialization is ShipSpecialization specialization)
        {
            run.Progression.SpecializationFlags[specialization.ToString()] = true;
        }

        foreach (UpgradeEffectDefinition effect in upgrade.Effects)
        {
            if (effect.Type == UpgradeEffectType.CargoLimitBonus)
            {
                run.Player.CargoLimit += effect.Value;
            }
        }

        return UpgradePurchaseResult.Success($"Installed {upgrade.Name} for {upgrade.CostCredits} credits.");
    }

    public int AdjustTravelCost(RunState run, DataSnapshot data, int baseCost)
    {
        int discountPercent = 0;
        foreach (string upgradeId in run.Progression.PurchasedUpgradeIds)
        {
            if (!data.UpgradesById.TryGetValue(upgradeId, out ShipUpgradeDefinition? upgrade))
            {
                continue;
            }

            foreach (UpgradeEffectDefinition effect in upgrade.Effects)
            {
                if (effect.Type == UpgradeEffectType.TravelCostDiscountPercent)
                {
                    discountPercent += effect.Value;
                }
            }
        }

        if (discountPercent <= 0)
        {
            return baseCost;
        }

        discountPercent = System.Math.Min(discountPercent, 90);
        int adjustedCost = baseCost * (100 - discountPercent) / 100;
        return baseCost > 0 ? System.Math.Max(1, adjustedCost) : adjustedCost;
    }

    public int GetCheapestTravelCostFromPort(RunState run, DataSnapshot data, TravelService travelService, PortDefinition origin)
    {
        int baseCost = travelService.GetCheapestTravelCostFromPort(origin, data.Ports);
        return AdjustTravelCost(run, data, baseCost);
    }

    public string SummarizeEffects(ShipUpgradeDefinition upgrade)
    {
        return string.Join(", ", upgrade.Effects.Select(effect => effect.Type switch
        {
            UpgradeEffectType.CargoLimitBonus => $"+{effect.Value} cargo",
            UpgradeEffectType.TravelCostDiscountPercent => $"-{effect.Value}% travel cost",
            UpgradeEffectType.MissionRewardBonusPercent => $"+{effect.Value}% mission rewards",
            _ => $"{effect.Type} {effect.Value}",
        }));
    }
}
