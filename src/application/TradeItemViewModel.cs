namespace StarSmugglerGo.Application;

public sealed class TradeItemViewModel
{
    public string ItemId { get; init; } = string.Empty;
    public string Name { get; init; } = string.Empty;
    public string Description { get; init; } = string.Empty;
    public int Price { get; init; }
    public int OwnedQuantity { get; init; }
}
