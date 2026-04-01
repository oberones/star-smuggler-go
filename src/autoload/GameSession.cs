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

    public RunState? CurrentRun { get; private set; }

    public bool HasActiveRun => CurrentRun is not null;

    public override void _Ready()
    {
        _dataRepository = GetNodeOrNull<DataRepository>("%DataRepository");
    }

    public void StartNewRun()
    {
        if (_dataRepository is null)
        {
            GD.PushError("GameSession could not find %DataRepository.");
            return;
        }

        CurrentRun = RunFactory.CreateNew(_dataRepository.Snapshot, _random);
        EmitSignal(SignalName.RunChanged);
    }

    public void ClearRun()
    {
        CurrentRun = null;
        EmitSignal(SignalName.RunChanged);
    }
}
