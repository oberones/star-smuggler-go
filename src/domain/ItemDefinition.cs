namespace StarSmugglerGo.Domain;

public sealed class ItemDefinition
{
    public string Id { get; init; } = string.Empty;
    public string Name { get; init; } = string.Empty;
    public string Description { get; init; } = string.Empty;
    public ItemRarity Rarity { get; init; }
    public int BasePrice { get; init; }
}
