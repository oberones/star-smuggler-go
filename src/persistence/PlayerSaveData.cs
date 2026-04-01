namespace StarSmugglerGo.Persistence;

public sealed class PlayerSaveData
{
    public int Credits { get; set; }
    public int CargoLimit { get; set; }
    public string CurrentPortId { get; set; } = string.Empty;
}
