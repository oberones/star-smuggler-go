using Godot;

namespace StarSmugglerGo.Presentation.Screens;

public partial class GameOverScreen : Control
{
    [Signal]
    public delegate void RecoveryRequestedEventHandler();

    [Signal]
    public delegate void RestartRequestedEventHandler();

    [Signal]
    public delegate void MenuRequestedEventHandler();

    private Label? _summaryLabel;
    private Button? _recoveryButton;
    private Label? _recoveryStatusLabel;
    private Button? _restartButton;
    private Button? _menuButton;
    private string _pendingSummary = string.Empty;
    private string _pendingRecoveryStatus = string.Empty;
    private bool _canRecover;
    private bool _isViewReady;

    public override void _Ready()
    {
        _summaryLabel = GetNodeOrNull<Label>("%SummaryLabel");
        _recoveryButton = GetNodeOrNull<Button>("%RecoveryButton");
        _recoveryStatusLabel = GetNodeOrNull<Label>("%RecoveryStatusLabel");
        _restartButton = GetNodeOrNull<Button>("%RestartButton");
        _menuButton = GetNodeOrNull<Button>("%MenuButton");

        if (_recoveryButton is not null)
        {
            _recoveryButton.Pressed += () => EmitSignal(SignalName.RecoveryRequested);
        }

        if (_restartButton is not null)
        {
            _restartButton.Pressed += () => EmitSignal(SignalName.RestartRequested);
        }

        if (_menuButton is not null)
        {
            _menuButton.Pressed += () => EmitSignal(SignalName.MenuRequested);
        }

        _isViewReady = true;
        ApplyViewState();
    }

    public void SetSummary(string summary)
    {
        _pendingSummary = summary;
        ApplyViewState();
    }

    public void SetRecoveryState(bool canRecover, string status)
    {
        _canRecover = canRecover;
        _pendingRecoveryStatus = status;
        ApplyViewState();
    }

    private void ApplyViewState()
    {
        if (!_isViewReady)
        {
            return;
        }

        if (_summaryLabel is not null)
        {
            _summaryLabel.Text = _pendingSummary;
        }

        if (_recoveryButton is not null)
        {
            _recoveryButton.Disabled = !_canRecover;
        }

        if (_recoveryStatusLabel is not null)
        {
            _recoveryStatusLabel.Text = _pendingRecoveryStatus;
        }
    }
}
