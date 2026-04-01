using System.Collections.Generic;
using System.Linq;

namespace StarSmugglerGo.Domain;

public sealed class CargoState
{
    private readonly Dictionary<string, int> _itemQuantities = new();

    public IReadOnlyDictionary<string, int> ItemQuantities => _itemQuantities;

    public int TotalUnits => _itemQuantities.Values.Sum();

    public int GetQuantity(string itemId)
    {
        return _itemQuantities.TryGetValue(itemId, out int quantity) ? quantity : 0;
    }

    public void SetQuantity(string itemId, int quantity)
    {
        if (quantity <= 0)
        {
            _itemQuantities.Remove(itemId);
            return;
        }

        _itemQuantities[itemId] = quantity;
    }
}
