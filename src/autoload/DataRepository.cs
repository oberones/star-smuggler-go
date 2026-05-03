using Godot;
using StarSmugglerGo.Domain;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text.Json;

namespace StarSmugglerGo.Autoload;

public partial class DataRepository : Node
{
    private const string PortsDataPath = "res://data/ports/ports.json";
    private const string ItemsDataPath = "res://data/items/items.json";
    private const string EventsDataPath = "res://data/events/events.json";
    private const string UpgradesDataPath = "res://data/upgrades/ship_upgrades.json";

    public bool IsInitialized { get; private set; }
    public DataSnapshot Snapshot { get; private set; } = new();

    public override void _Ready()
    {
        if (IsGoRuntimeAuthorityEnabled())
        {
            GD.Print("DataRepository is passive because STARSMUGGLER_GO_RUNTIME is enabled.");
            return;
        }

        Snapshot = LoadSnapshot();
        IsInitialized = true;
    }

    private static bool IsGoRuntimeAuthorityEnabled()
    {
        string value = OS.GetEnvironment("STARSMUGGLER_GO_RUNTIME");
        return string.Equals(value, "1", StringComparison.Ordinal) ||
            string.Equals(value, "true", StringComparison.OrdinalIgnoreCase);
    }

    private static DataSnapshot LoadSnapshot()
    {
        List<PortRecord> portRecords = LoadJson<List<PortRecord>>(PortsDataPath) ?? new();
        List<ItemRecord> itemRecords = LoadJson<List<ItemRecord>>(ItemsDataPath) ?? new();
        List<EventRecord> eventRecords = LoadJson<List<EventRecord>>(EventsDataPath) ?? new();
        List<UpgradeRecord> upgradeRecords = LoadJson<List<UpgradeRecord>>(UpgradesDataPath) ?? new();

        List<PortDefinition> ports = portRecords.Select(record => new PortDefinition
        {
            Id = record.Id,
            Name = record.Name,
            Description = record.Description,
            Zone = Enum.Parse<PortZone>(record.Zone),
            BackgroundTexturePath = record.BackgroundTexturePath,
            PreviewTexturePath = record.PreviewTexturePath,
            TradeBackgroundPath = record.TradeBackgroundPath,
            MusicTrackId = record.MusicTrackId,
        }).ToList();

        List<ItemDefinition> items = itemRecords.Select(record => new ItemDefinition
        {
            Id = record.Id,
            Name = record.Name,
            Description = record.Description,
            Rarity = Enum.Parse<ItemRarity>(record.Rarity),
            BasePrice = record.BasePrice,
        }).ToList();

        List<EventDefinition> events = eventRecords.Select(record => new EventDefinition
        {
            Id = record.Id,
            Name = record.Name,
            DescriptionTemplate = record.DescriptionTemplate,
            EffectType = Enum.Parse<EventEffectType>(record.EffectType),
            Weight = record.Weight,
            Parameters = new Dictionary<string, double>(record.Parameters, StringComparer.Ordinal),
        }).ToList();

        List<ShipUpgradeDefinition> upgrades = upgradeRecords.Select(record => new ShipUpgradeDefinition
        {
            Id = record.Id,
            Name = record.Name,
            Description = record.Description,
            Category = Enum.Parse<UpgradeCategory>(record.Category),
            CostCredits = record.CostCredits,
            RequiredFactionId = record.RequiredFactionId,
            MinimumStanding = record.MinimumStanding,
            Specialization = string.IsNullOrWhiteSpace(record.Specialization)
                ? null
                : Enum.Parse<ShipSpecialization>(record.Specialization),
            Effects = record.Effects.Select(effect => new UpgradeEffectDefinition
            {
                Type = Enum.Parse<UpgradeEffectType>(effect.Type),
                Value = effect.Value,
            }).ToList(),
        }).ToList();

        return new DataSnapshot
        {
            Ports = ports,
            Items = items,
            Events = events,
            Upgrades = upgrades,
            PortsById = ports.ToDictionary(definition => definition.Id, StringComparer.Ordinal),
            ItemsById = items.ToDictionary(definition => definition.Id, StringComparer.Ordinal),
            EventsById = events.ToDictionary(definition => definition.Id, StringComparer.Ordinal),
            UpgradesById = upgrades.ToDictionary(definition => definition.Id, StringComparer.Ordinal),
        };
    }

    private static T? LoadJson<T>(string resourcePath)
    {
        using var file = Godot.FileAccess.Open(resourcePath, Godot.FileAccess.ModeFlags.Read);
        if (file is null)
        {
            GD.PushError($"DataRepository failed to open '{resourcePath}'.");
            return default;
        }

        string json = file.GetAsText();

        return JsonSerializer.Deserialize<T>(json, new JsonSerializerOptions
        {
            PropertyNameCaseInsensitive = true,
        });
    }

    private sealed class PortRecord
    {
        public string Id { get; set; } = string.Empty;
        public string Name { get; set; } = string.Empty;
        public string Description { get; set; } = string.Empty;
        public string Zone { get; set; } = string.Empty;
        public string BackgroundTexturePath { get; set; } = string.Empty;
        public string PreviewTexturePath { get; set; } = string.Empty;
        public string TradeBackgroundPath { get; set; } = string.Empty;
        public string MusicTrackId { get; set; } = string.Empty;
    }

    private sealed class ItemRecord
    {
        public string Id { get; set; } = string.Empty;
        public string Name { get; set; } = string.Empty;
        public string Description { get; set; } = string.Empty;
        public string Rarity { get; set; } = string.Empty;
        public int BasePrice { get; set; }
    }

    private sealed class EventRecord
    {
        public string Id { get; set; } = string.Empty;
        public string Name { get; set; } = string.Empty;
        public string DescriptionTemplate { get; set; } = string.Empty;
        public string EffectType { get; set; } = string.Empty;
        public int Weight { get; set; } = 1;
        public Dictionary<string, double> Parameters { get; set; } = new(StringComparer.Ordinal);
    }

    private sealed class UpgradeRecord
    {
        public string Id { get; set; } = string.Empty;
        public string Name { get; set; } = string.Empty;
        public string Description { get; set; } = string.Empty;
        public string Category { get; set; } = string.Empty;
        public int CostCredits { get; set; }
        public string RequiredFactionId { get; set; } = string.Empty;
        public string MinimumStanding { get; set; } = string.Empty;
        public string Specialization { get; set; } = string.Empty;
        public List<UpgradeEffectRecord> Effects { get; set; } = new();
    }

    private sealed class UpgradeEffectRecord
    {
        public string Type { get; set; } = string.Empty;
        public int Value { get; set; }
    }
}
