using StarSmugglerGo.Domain;
using System;
using System.Collections.Generic;
using System.Linq;

namespace StarSmugglerGo.Services;

public sealed class TravelService
{
    public int GetTravelCost(PortDefinition fromPort, PortDefinition toPort)
    {
        if (fromPort.Id == toPort.Id)
        {
            return 0;
        }

        int baseCost = 15;
        int zoneDifference = Math.Abs((int)fromPort.Zone - (int)toPort.Zone);
        int cost = baseCost + (zoneDifference * 2);

        if (zoneDifference >= 2)
        {
            cost *= 2;
        }

        return cost;
    }

    public int GetCheapestTravelCostFromPort(PortDefinition origin, IEnumerable<PortDefinition> ports)
    {
        int cheapest = int.MaxValue;

        foreach (PortDefinition destination in ports)
        {
            if (destination.Id == origin.Id)
            {
                continue;
            }

            cheapest = Math.Min(cheapest, GetTravelCost(origin, destination));
        }

        return cheapest == int.MaxValue ? 0 : cheapest;
    }

    public IReadOnlyList<PortDefinition> GetDestinationsFromPort(PortDefinition origin, IEnumerable<PortDefinition> ports)
    {
        return ports
            .Where(port => port.Id != origin.Id)
            .OrderBy(port => port.Zone)
            .ThenBy(port => port.Name, StringComparer.Ordinal)
            .ToList();
    }
}
