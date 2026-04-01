using System.Collections.Generic;

namespace StarSmugglerGo.Domain;

public sealed class EventDefinition
{
    public string Id { get; init; } = string.Empty;
    public string Name { get; init; } = string.Empty;
    public string DescriptionTemplate { get; init; } = string.Empty;
    public EventEffectType EffectType { get; init; }
    public int Weight { get; init; } = 1;
    public IReadOnlyDictionary<string, double> Parameters { get; init; } = new Dictionary<string, double>();
}
