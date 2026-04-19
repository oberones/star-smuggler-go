using Godot;
using StarSmugglerGo.Domain;
using StarSmugglerGo.Persistence;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text.Json;

namespace StarSmugglerGo.Autoload;

public partial class SaveService : Node
{
    public const int CurrentSaveVersion = 1;

    public string SavePath => "user://save.json";

    public bool HasSave()
    {
        if (IsGoRuntimeAuthorityEnabled())
        {
            return false;
        }

        return Godot.FileAccess.FileExists(SavePath);
    }

    public void SaveRun(RunState run)
    {
        if (IsGoRuntimeAuthorityEnabled())
        {
            GD.Print("SaveService.SaveRun ignored because STARSMUGGLER_GO_RUNTIME is enabled.");
            return;
        }

        SaveData dto = ToSaveData(run);
        string json = JsonSerializer.Serialize(dto, new JsonSerializerOptions
        {
            WriteIndented = true,
        });

        using Godot.FileAccess file = Godot.FileAccess.Open(SavePath, Godot.FileAccess.ModeFlags.Write);
        file.StoreString(json);
    }

    public RunState? LoadRun()
    {
        if (IsGoRuntimeAuthorityEnabled())
        {
            GD.Print("SaveService.LoadRun ignored because STARSMUGGLER_GO_RUNTIME is enabled.");
            return null;
        }

        if (!HasSave())
        {
            return null;
        }

        using Godot.FileAccess file = Godot.FileAccess.Open(SavePath, Godot.FileAccess.ModeFlags.Read);
        string json = file.GetAsText();
        SaveData? dto = JsonSerializer.Deserialize<SaveData>(json, new JsonSerializerOptions
        {
            PropertyNameCaseInsensitive = true,
        });

        if (dto is null || dto.Version > CurrentSaveVersion)
        {
            return null;
        }

        return FromSaveData(dto);
    }

    public void DeleteSave()
    {
        if (IsGoRuntimeAuthorityEnabled())
        {
            GD.Print("SaveService.DeleteSave ignored because STARSMUGGLER_GO_RUNTIME is enabled.");
            return;
        }

        if (HasSave())
        {
            DirAccess.RemoveAbsolute(ProjectSettings.GlobalizePath(SavePath));
        }
    }

    private static bool IsGoRuntimeAuthorityEnabled()
    {
        string value = OS.GetEnvironment("STARSMUGGLER_GO_RUNTIME");
        return string.Equals(value, "1", StringComparison.Ordinal) ||
            string.Equals(value, "true", StringComparison.OrdinalIgnoreCase);
    }

    private static SaveData ToSaveData(RunState run)
    {
        return new SaveData
        {
            Version = CurrentSaveVersion,
            Player = new PlayerSaveData
            {
                Credits = run.Player.Credits,
                CargoLimit = run.Player.CargoLimit,
                CurrentPortId = run.Player.CurrentPortId,
            },
            CargoByItemId = new Dictionary<string, int>(run.Cargo.ItemQuantities, StringComparer.Ordinal),
            Markets = run.MarketsByPortId.Values.Select(market => new MarketSnapshotSaveData
            {
                PortId = market.PortId,
                AvailableItemIds = new List<string>(market.AvailableItemIds),
                PricesByItemId = new Dictionary<string, int>(market.PricesByItemId, StringComparer.Ordinal),
            }).ToList(),
            JumpsSinceLastUpdate = run.JumpsSinceLastUpdate,
            RecentEvent = run.RecentEvent is null
                ? null
                : new EventResultSaveData
                {
                    EventId = run.RecentEvent.EventId,
                    Name = run.RecentEvent.Name,
                    ResolvedDescription = run.RecentEvent.ResolvedDescription,
                    RolledValues = new Dictionary<string, double>(run.RecentEvent.RolledValues, StringComparer.Ordinal),
                },
        };
    }

    private static RunState FromSaveData(SaveData dto)
    {
        CargoState cargo = new();
        foreach ((string itemId, int quantity) in dto.CargoByItemId)
        {
            cargo.SetQuantity(itemId, quantity);
        }

        Dictionary<string, MarketSnapshot> markets = dto.Markets.ToDictionary(
            market => market.PortId,
            market => new MarketSnapshot
            {
                PortId = market.PortId,
                AvailableItemIds = new List<string>(market.AvailableItemIds),
                PricesByItemId = new Dictionary<string, int>(market.PricesByItemId, StringComparer.Ordinal),
            },
            StringComparer.Ordinal);

        return new RunState
        {
            Player = new PlayerState
            {
                Credits = dto.Player.Credits,
                CargoLimit = dto.Player.CargoLimit,
                CurrentPortId = dto.Player.CurrentPortId,
            },
            Cargo = cargo,
            MarketsByPortId = markets,
            JumpsSinceLastUpdate = dto.JumpsSinceLastUpdate,
            RecentEvent = dto.RecentEvent is null
                ? null
                : new EventResult
                {
                    EventId = dto.RecentEvent.EventId,
                    Name = dto.RecentEvent.Name,
                    ResolvedDescription = dto.RecentEvent.ResolvedDescription,
                    RolledValues = new Dictionary<string, double>(dto.RecentEvent.RolledValues, StringComparer.Ordinal),
                },
        };
    }
}
