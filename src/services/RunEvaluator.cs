using StarSmugglerGo.Domain;

namespace StarSmugglerGo.Services;

public sealed class RunEvaluator
{
    public bool IsGameOver(RunState run, DataSnapshot data, EconomyService economyService, TravelService travelService, UpgradeService upgradeService)
    {
        if (!data.PortsById.TryGetValue(run.Player.CurrentPortId, out PortDefinition? currentPort))
        {
            return false;
        }

        int cheapestTravelCost = upgradeService.GetCheapestTravelCostFromPort(run, data, travelService, currentPort);

        if (run.Player.Credits >= cheapestTravelCost)
        {
            return false;
        }

        int sellableValue = economyService.GetSellableCargoValueAtCurrentPort(run, data);
        return run.Player.Credits + sellableValue < cheapestTravelCost;
    }
}
