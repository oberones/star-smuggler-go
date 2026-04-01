using Godot;
using StarSmugglerGo.Application;

namespace StarSmugglerGo.Presentation.Screens;

public partial class TravelAnimationScreen : Control
{
    [Signal]
    public delegate void AnimationFinishedEventHandler();

    private Label? _routeLabel;
    private Label? _costLabel;
    private Label? _statusLabel;
    private ProgressBar? _progressBar;
    private Button? _skipButton;
    private ColorRect? _starsA;
    private ColorRect? _starsB;

    private TravelAnimationViewModel _viewModel = new();
    private bool _isViewReady;
    private double _elapsedSeconds;
    private bool _isComplete;

    public override void _Ready()
    {
        _routeLabel = GetNodeOrNull<Label>("%RouteLabel");
        _costLabel = GetNodeOrNull<Label>("%CostLabel");
        _statusLabel = GetNodeOrNull<Label>("%StatusLabel");
        _progressBar = GetNodeOrNull<ProgressBar>("%ProgressBar");
        _skipButton = GetNodeOrNull<Button>("%SkipButton");
        _starsA = GetNodeOrNull<ColorRect>("StarsA");
        _starsB = GetNodeOrNull<ColorRect>("StarsB");

        if (_skipButton is not null)
        {
            _skipButton.Pressed += CompleteImmediately;
        }

        _isViewReady = true;
        ApplyViewState();
    }

    public override void _Process(double delta)
    {
        if (!_isViewReady || _isComplete)
        {
            return;
        }

        _elapsedSeconds += delta;
        UpdateProgressVisuals();

        if (_elapsedSeconds >= _viewModel.DurationSeconds)
        {
            CompleteImmediately();
        }
    }

    public void Bind(TravelAnimationViewModel viewModel)
    {
        _viewModel = viewModel;
        _elapsedSeconds = 0;
        _isComplete = false;
        ApplyViewState();
    }

    private void ApplyViewState()
    {
        if (!_isViewReady)
        {
            return;
        }

        if (_routeLabel is not null)
        {
            _routeLabel.Text = $"{_viewModel.OriginName} -> {_viewModel.DestinationName}";
        }

        if (_costLabel is not null)
        {
            _costLabel.Text = $"Travel cost: {_viewModel.TravelCost} credits";
        }

        if (_statusLabel is not null)
        {
            _statusLabel.Text = string.IsNullOrWhiteSpace(_viewModel.StatusMessage)
                ? "Cruising through the black."
                : _viewModel.StatusMessage;
        }

        UpdateProgressVisuals();
    }

    private void UpdateProgressVisuals()
    {
        double duration = _viewModel.DurationSeconds <= 0 ? 0.1 : _viewModel.DurationSeconds;
        double progress = Mathf.Clamp((float)(_elapsedSeconds / duration), 0f, 1f);

        if (_progressBar is not null)
        {
            _progressBar.Value = progress * 100.0;
        }

        if (_starsA is not null)
        {
            _starsA.Position = new Vector2((float)(progress * 90.0), _starsA.Position.Y);
        }

        if (_starsB is not null)
        {
            _starsB.Position = new Vector2((float)(-120.0 + (progress * 150.0)), _starsB.Position.Y);
        }
    }

    private void CompleteImmediately()
    {
        if (_isComplete)
        {
            return;
        }

        _isComplete = true;
        EmitSignal(SignalName.AnimationFinished);
    }
}
