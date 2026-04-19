using System.Collections.Generic;

namespace StarSmugglerGo.Domain;

public sealed class RunState
{
    public PlayerState Player { get; init; } = new();
    public CargoState Cargo { get; init; } = new();
    public Dictionary<string, MarketSnapshot> MarketsByPortId { get; init; } = new();
    public bool EmergencyRecoveryUsed { get; set; }
    public int JumpsSinceLastUpdate { get; set; }
    public EventResult? RecentEvent { get; set; }
}
