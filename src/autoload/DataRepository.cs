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

    public bool IsInitialized { get; private set; }
    public DataSnapshot Snapshot { get; private set; } = new();

    public override void _Ready()
    {
        Snapshot = LoadSnapshot();
        IsInitialized = true;
    }

    private static DataSnapshot LoadSnapshot()
    {
        List<PortRecord> portRecords = LoadJson<List<PortRecord>>(PortsDataPath) ?? new();
        List<ItemRecord> itemRecords = LoadJson<List<ItemRecord>>(ItemsDataPath) ?? new();
        List<EventRecord> eventRecords = LoadJson<List<EventRecord>>(EventsDataPath) ?? new();

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

        return new DataSnapshot
        {
            Ports = ports,
            Items = items,
            Events = events,
            PortsById = ports.ToDictionary(definition => definition.Id, StringComparer.Ordinal),
            ItemsById = items.ToDictionary(definition => definition.Id, StringComparer.Ordinal),
            EventsById = events.ToDictionary(definition => definition.Id, StringComparer.Ordinal),
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
}
