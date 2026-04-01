using System.Collections.Generic;

namespace StarSmugglerGo.Persistence;

public sealed class EventResultSaveData
{
    public string EventId { get; set; } = string.Empty;
    public string Name { get; set; } = string.Empty;
    public string ResolvedDescription { get; set; } = string.Empty;
    public Dictionary<string, double> RolledValues { get; set; } = new();
}
