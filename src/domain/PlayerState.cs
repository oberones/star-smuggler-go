namespace StarSmugglerGo.Domain;

public sealed class PlayerState
{
    public const int StartingCredits = 500;
    public const int StartingCargoLimit = 30;

    public int Credits { get; set; } = StartingCredits;
    public int CargoLimit { get; set; } = StartingCargoLimit;
    public string CurrentPortId { get; set; } = string.Empty;
}
