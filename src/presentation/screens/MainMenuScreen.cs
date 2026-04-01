using Godot;

namespace StarSmugglerGo.Presentation.Screens;

public partial class MainMenuScreen : Control
{
    [Signal]
    public delegate void StartRequestedEventHandler();

    [Signal]
    public delegate void QuitRequestedEventHandler();

    private Button? _startButton;
    private Button? _quitButton;
    private Label? _statusLabel;
    private string _pendingStatusMessage = "Waiting for the first gameplay slice.";
    private bool _isViewReady;

    public override void _Ready()
    {
        _startButton = GetNodeOrNull<Button>("%StartButton");
        _quitButton = GetNodeOrNull<Button>("%QuitButton");
        _statusLabel = GetNodeOrNull<Label>("%StatusLabel");

        if (_startButton is not null)
        {
            _startButton.Pressed += OnStartPressed;
            _startButton.GrabFocus();
        }

        if (_quitButton is not null)
        {
            _quitButton.Pressed += OnQuitPressed;
        }

        _isViewReady = true;
        ApplyViewState();
    }

    public void SetStatusMessage(string message)
    {
        _pendingStatusMessage = message;
        ApplyViewState();
    }

    private void OnStartPressed()
    {
        EmitSignal(SignalName.StartRequested);
    }

    private void OnQuitPressed()
    {
        EmitSignal(SignalName.QuitRequested);
    }

    private void ApplyViewState()
    {
        if (!_isViewReady)
        {
            return;
        }

        if (_statusLabel is not null)
        {
            _statusLabel.Text = _pendingStatusMessage;
        }
    }
}
