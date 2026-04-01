using Godot;
using StarSmugglerGo.Domain;
using System;

namespace StarSmugglerGo.Autoload;

public partial class GameSession : Node
{
    [Signal]
    public delegate void RunChangedEventHandler();

    private readonly Random _random = new();
    private DataRepository? _dataRepository;
    private SaveService? _saveService;

    public RunState? CurrentRun { get; private set; }

    public bool HasActiveRun => CurrentRun is not null;

    public override void _Ready()
    {
        _dataRepository = GetNodeOrNull<DataRepository>("%DataRepository");
        _saveService = GetNodeOrNull<SaveService>("%SaveService");
    }

    public void StartNewRun()
    {
        if (_dataRepository is null)
        {
            GD.PushError("GameSession could not find %DataRepository.");
            return;
        }

        CurrentRun = RunFactory.CreateNew(_dataRepository.Snapshot, _random);
        _saveService?.SaveRun(CurrentRun);
        EmitSignal(SignalName.RunChanged);
    }

    public bool TryLoadSavedRun()
    {
        if (_saveService is null || !_saveService.HasSave())
        {
            return false;
        }

        RunState? loadedRun = _saveService.LoadRun();
        if (loadedRun is null)
        {
            return false;
        }

        CurrentRun = loadedRun;
        EmitSignal(SignalName.RunChanged);
        return true;
    }

    public void SaveCurrentRun()
    {
        if (CurrentRun is not null)
        {
            _saveService?.SaveRun(CurrentRun);
        }
    }

    public void ClearRun()
    {
        CurrentRun = null;
        EmitSignal(SignalName.RunChanged);
    }
}
