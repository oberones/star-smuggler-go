using Godot;
using StarSmugglerGo.Application;

namespace StarSmugglerGo.Presentation.Screens;

public partial class TravelScreen : Control
{
    [Signal]
    public delegate void BackRequestedEventHandler();

    [Signal]
    public delegate void TravelRequestedEventHandler(string destinationPortId);

    private Label? _titleLabel;
    private Label? _summaryLabel;
    private ItemList? _destinationList;
    private Label? _selectedPortLabel;
    private Label? _zoneLabel;
    private Label? _descriptionLabel;
    private Label? _costLabel;
    private Label? _statusLabel;
    private TextureRect? _backdrop;
    private TextureRect? _previewTexture;
    private Button? _travelButton;
    private Button? _backButton;

    private TravelScreenViewModel _viewModel = new();
    private int _selectedIndex = -1;
    private bool _isViewReady;

    public override void _Ready()
    {
        _backdrop = GetNodeOrNull<TextureRect>("%Backdrop");
        _titleLabel = GetNodeOrNull<Label>("%TitleLabel");
        _summaryLabel = GetNodeOrNull<Label>("%SummaryLabel");
        _destinationList = GetNodeOrNull<ItemList>("%DestinationList");
        _selectedPortLabel = GetNodeOrNull<Label>("%SelectedPortLabel");
        _zoneLabel = GetNodeOrNull<Label>("%ZoneLabel");
        _descriptionLabel = GetNodeOrNull<Label>("%DescriptionLabel");
        _costLabel = GetNodeOrNull<Label>("%CostLabel");
        _statusLabel = GetNodeOrNull<Label>("%StatusLabel");
        _previewTexture = GetNodeOrNull<TextureRect>("%PreviewTexture");
        _travelButton = GetNodeOrNull<Button>("%TravelButton");
        _backButton = GetNodeOrNull<Button>("%BackButton");

        if (_destinationList is not null)
        {
            _destinationList.ItemSelected += OnDestinationSelected;
        }

        if (_travelButton is not null)
        {
            _travelButton.Pressed += OnTravelPressed;
        }

        if (_backButton is not null)
        {
            _backButton.Pressed += () => EmitSignal(SignalName.BackRequested);
        }

        _isViewReady = true;
        ApplyViewState();
    }

    public void Bind(TravelScreenViewModel viewModel)
    {
        _viewModel = viewModel;
        if (_selectedIndex < 0 || _selectedIndex >= viewModel.Destinations.Count)
        {
            _selectedIndex = viewModel.Destinations.Count > 0 ? 0 : -1;
        }
        ApplyViewState();
    }

    private void ApplyViewState()
    {
        if (!_isViewReady)
        {
            return;
        }

        if (_titleLabel is not null)
        {
            _titleLabel.Text = $"Travel From {_viewModel.CurrentPortName}";
        }

        if (_backdrop is not null)
        {
            _backdrop.Texture = string.IsNullOrWhiteSpace(_viewModel.BackgroundTexturePath)
                ? null
                : ResourceLoader.Load<Texture2D>(_viewModel.BackgroundTexturePath);
        }

        if (_summaryLabel is not null)
        {
            _summaryLabel.Text = $"Available credits: {_viewModel.Credits}";
        }

        if (_destinationList is not null)
        {
            _destinationList.Clear();
            foreach (TravelDestinationViewModel destination in _viewModel.Destinations)
            {
                _destinationList.AddItem($"{destination.Name}  |  {destination.TravelCost} cr");
            }

            if (_selectedIndex >= 0)
            {
                _destinationList.Select(_selectedIndex);
            }
        }

        if (_statusLabel is not null)
        {
            _statusLabel.Text = _viewModel.StatusMessage;
        }

        RefreshSelectionDetails();
    }

    public void SetStatusMessage(string message)
    {
        _viewModel = new TravelScreenViewModel
        {
            CurrentPortName = _viewModel.CurrentPortName,
            BackgroundTexturePath = _viewModel.BackgroundTexturePath,
            Credits = _viewModel.Credits,
            Destinations = _viewModel.Destinations,
            StatusMessage = message,
        };
        ApplyViewState();
    }

    private void OnDestinationSelected(long index)
    {
        _selectedIndex = (int)index;
        RefreshSelectionDetails();
    }

    private void OnTravelPressed()
    {
        TravelDestinationViewModel? destination = GetSelectedDestination();
        if (destination is null)
        {
            return;
        }

        EmitSignal(SignalName.TravelRequested, destination.PortId);
    }

    private void RefreshSelectionDetails()
    {
        TravelDestinationViewModel? destination = GetSelectedDestination();
        if (destination is null)
        {
            if (_selectedPortLabel is not null) _selectedPortLabel.Text = "No destinations available";
            if (_zoneLabel is not null) _zoneLabel.Text = string.Empty;
            if (_descriptionLabel is not null) _descriptionLabel.Text = string.Empty;
            if (_costLabel is not null) _costLabel.Text = string.Empty;
            if (_previewTexture is not null) _previewTexture.Texture = null;
            return;
        }

        if (_selectedPortLabel is not null) _selectedPortLabel.Text = destination.Name;
        if (_zoneLabel is not null) _zoneLabel.Text = $"{destination.ZoneName} Zone";
        if (_descriptionLabel is not null) _descriptionLabel.Text = destination.Description;
        if (_costLabel is not null) _costLabel.Text = $"Travel cost: {destination.TravelCost} credits";
        if (_previewTexture is not null)
        {
            _previewTexture.Texture = string.IsNullOrWhiteSpace(destination.PreviewTexturePath)
                ? null
                : ResourceLoader.Load<Texture2D>(destination.PreviewTexturePath);
        }
    }

    private TravelDestinationViewModel? GetSelectedDestination()
    {
        return _selectedIndex >= 0 && _selectedIndex < _viewModel.Destinations.Count
            ? _viewModel.Destinations[_selectedIndex]
            : null;
    }
}
