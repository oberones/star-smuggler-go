using System;
using System.Linq;
using StarSmugglerGo.Services;

namespace StarSmugglerGo.Domain;

public static class RunFactory
{
    public static RunState CreateNew(DataSnapshot data, Random random)
    {
        if (data.Ports.Count == 0)
        {
            throw new InvalidOperationException("Cannot start a run without any port definitions.");
        }

        if (data.Items.Count == 0)
        {
            throw new InvalidOperationException("Cannot start a run without any item definitions.");
        }

        List<PortDefinition> innerPorts = data.Ports.Where(port => port.Zone == PortZone.Inner).ToList();
        if (innerPorts.Count == 0)
        {
            throw new InvalidOperationException("Cannot start a run without at least one Inner-zone port.");
        }

        PortDefinition startingPort = innerPorts[random.Next(innerPorts.Count)];
        var economyService = new EconomyService();
        Dictionary<string, MarketSnapshot> marketSnapshots = economyService.CreateInitialMarkets(data, random);

        return new RunState
        {
            Player = new PlayerState
            {
                Credits = PlayerState.StartingCredits,
                CargoLimit = PlayerState.StartingCargoLimit,
                CurrentPortId = startingPort.Id,
            },
            Cargo = new CargoState(),
            MarketsByPortId = marketSnapshots,
            Progression = new ShipProgressionState(),
            JumpsSinceLastUpdate = 0,
        };
    }
}
