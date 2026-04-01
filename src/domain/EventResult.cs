using System.Collections.Generic;

namespace StarSmugglerGo.Domain;

public sealed class EventResult
{
    public string EventId { get; init; } = string.Empty;
    public string Name { get; init; } = string.Empty;
    public string ResolvedDescription { get; init; } = string.Empty;
    public IReadOnlyDictionary<string, double> RolledValues { get; init; } = new Dictionary<string, double>();
}
