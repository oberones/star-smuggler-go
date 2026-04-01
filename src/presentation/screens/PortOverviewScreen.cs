using Godot;
using StarSmugglerGo.Application;
using System.Linq;

namespace StarSmugglerGo.Presentation.Screens;

public partial class PortOverviewScreen : Control
{
    private TextureRect? _backdrop;
    [Signal]
    public delegate void BackRequestedEventHandler();

    [Signal]
    public delegate void TravelRequestedEventHandler();

    [Signal]
    public delegate void TradeRequestedEventHandler();

    private Label? _portNameLabel;
    private Label? _zoneLabel;
    private Label? _descriptionLabel;
    private Label? _statsLabel;
    private Label? _statusLabel;
    private Label? _goodsListLabel;
    private Button? _backButton;
    private Button? _travelButton;
    private Button? _tradeButton;
    private PortOverviewViewModel _viewModel = new();
    private string? _overrideStatusMessage;
    private bool _isViewReady;

    public override void _Ready()
    {
        _backdrop = GetNodeOrNull<TextureRect>("%Backdrop");
        _portNameLabel = GetNodeOrNull<Label>("%PortNameLabel");
        _zoneLabel = GetNodeOrNull<Label>("%ZoneLabel");
        _descriptionLabel = GetNodeOrNull<Label>("%DescriptionLabel");
        _statsLabel = GetNodeOrNull<Label>("%StatsLabel");
        _statusLabel = GetNodeOrNull<Label>("%StatusLabel");
        _goodsListLabel = GetNodeOrNull<Label>("%GoodsListLabel");
        _backButton = GetNodeOrNull<Button>("%BackButton");
        _travelButton = GetNodeOrNull<Button>("%TravelButton");
        _tradeButton = GetNodeOrNull<Button>("%TradeButton");

        if (_backButton is not null)
        {
            _backButton.Pressed += () => EmitSignal(SignalName.BackRequested);
        }

        if (_travelButton is not null)
        {
            _travelButton.Pressed += () => EmitSignal(SignalName.TravelRequested);
        }

        if (_tradeButton is not null)
        {
            _tradeButton.Pressed += () => EmitSignal(SignalName.TradeRequested);
        }

        _isViewReady = true;
        ApplyViewState();
    }

    public void Bind(PortOverviewViewModel viewModel)
    {
        _viewModel = viewModel;
        _overrideStatusMessage = null;
        ApplyViewState();
    }

    public void SetStatusMessage(string message)
    {
        _overrideStatusMessage = message;
        ApplyViewState();
    }

    private void ApplyViewState()
    {
        if (!_isViewReady)
        {
            return;
        }

        if (_portNameLabel is not null)
        {
            _portNameLabel.Text = _viewModel.PortName;
        }

        if (_backdrop is not null)
        {
            _backdrop.Texture = string.IsNullOrWhiteSpace(_viewModel.BackgroundTexturePath)
                ? null
                : ResourceLoader.Load<Texture2D>(_viewModel.BackgroundTexturePath);
        }

        if (_zoneLabel is not null)
        {
            _zoneLabel.Text = $"{_viewModel.ZoneName} Zone";
        }

        if (_descriptionLabel is not null)
        {
            _descriptionLabel.Text = _viewModel.PortDescription;
        }

        if (_statsLabel is not null)
        {
            _statsLabel.Text =
                $"Credits: {_viewModel.Credits}\n" +
                $"Cargo: {_viewModel.CargoLoad}/{_viewModel.CargoLimit}\n" +
                $"Cheapest route cost: {_viewModel.CheapestTravelCost}";
        }

        if (_statusLabel is not null)
        {
            string statusText = _viewModel.IsGameOver
                ? "Status: Stranded. This run would currently evaluate as game over."
                : "Status: Operational. The trading loop can continue from here.";

            if (!string.IsNullOrWhiteSpace(_viewModel.RecentEventText))
            {
                statusText += $"\nRecent event: {_viewModel.RecentEventText}";
            }

            _statusLabel.Text = _overrideStatusMessage ?? statusText;
        }

        if (_goodsListLabel is not null)
        {
            _goodsListLabel.Text = _viewModel.AvailableGoods.Count == 0
                ? "No goods are currently loaded for this port."
                : string.Join("\n", _viewModel.AvailableGoods.Select(good => $"• {good}"));
        }
    }
}
