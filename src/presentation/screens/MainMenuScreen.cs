using Godot;

namespace StarSmugglerGo.Presentation.Screens;

public partial class MainMenuScreen : Control
{
    [Signal]
    public delegate void StartRequestedEventHandler();

    [Signal]
    public delegate void ContinueRequestedEventHandler();

    [Signal]
    public delegate void QuitRequestedEventHandler();

    private Button? _startButton;
    private Button? _continueButton;
    private Button? _quitButton;
    private Label? _statusLabel;
    private string _pendingStatusMessage = "Waiting for the first gameplay slice.";
    private bool _isViewReady;
    private bool _canContinue;

    public override void _Ready()
    {
        _startButton = GetNodeOrNull<Button>("%StartButton");
        _continueButton = GetNodeOrNull<Button>("%ContinueButton");
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

        if (_continueButton is not null)
        {
            _continueButton.Pressed += OnContinuePressed;
        }

        _isViewReady = true;
        ApplyViewState();
    }

    public void SetCanContinue(bool canContinue)
    {
        _canContinue = canContinue;
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

    private void OnContinuePressed()
    {
        EmitSignal(SignalName.ContinueRequested);
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

        if (_continueButton is not null)
        {
            _continueButton.Disabled = !_canContinue;
            _continueButton.Visible = true;
        }
    }
}
