using Godot;

namespace StarSmugglerGo.Presentation.Screens;

public partial class GameOverScreen : Control
{
    [Signal]
    public delegate void RestartRequestedEventHandler();

    [Signal]
    public delegate void MenuRequestedEventHandler();

    private Label? _summaryLabel;
    private Button? _restartButton;
    private Button? _menuButton;
    private string _pendingSummary = string.Empty;
    private bool _isViewReady;

    public override void _Ready()
    {
        _summaryLabel = GetNodeOrNull<Label>("%SummaryLabel");
        _restartButton = GetNodeOrNull<Button>("%RestartButton");
        _menuButton = GetNodeOrNull<Button>("%MenuButton");

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
    }
}
