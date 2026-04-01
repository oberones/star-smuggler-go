using System.Collections.Generic;

namespace StarSmugglerGo.Domain;

public sealed class DataSnapshot
{
    public IReadOnlyList<PortDefinition> Ports { get; init; } = new List<PortDefinition>();
    public IReadOnlyList<ItemDefinition> Items { get; init; } = new List<ItemDefinition>();
    public IReadOnlyList<EventDefinition> Events { get; init; } = new List<EventDefinition>();
    public IReadOnlyDictionary<string, PortDefinition> PortsById { get; init; } = new Dictionary<string, PortDefinition>();
    public IReadOnlyDictionary<string, ItemDefinition> ItemsById { get; init; } = new Dictionary<string, ItemDefinition>();
    public IReadOnlyDictionary<string, EventDefinition> EventsById { get; init; } = new Dictionary<string, EventDefinition>();
}
