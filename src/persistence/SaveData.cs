using System.Collections.Generic;

namespace StarSmugglerGo.Persistence;

public sealed class SaveData
{
    public int Version { get; set; }
    public PlayerSaveData Player { get; set; } = new();
    public Dictionary<string, int> CargoByItemId { get; set; } = new();
    public List<MarketSnapshotSaveData> Markets { get; set; } = new();
    public bool EmergencyRecoveryUsed { get; set; }
    public int JumpsSinceLastUpdate { get; set; }
    public EventResultSaveData? RecentEvent { get; set; }
}
