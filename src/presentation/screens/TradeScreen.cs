using Godot;
using StarSmugglerGo.Application;

namespace StarSmugglerGo.Presentation.Screens;

public partial class TradeScreen : Control
{
    private TextureRect? _backdrop;
    [Signal]
    public delegate void BackRequestedEventHandler();

    [Signal]
    public delegate void BuyRequestedEventHandler(string itemId, int quantity);

    [Signal]
    public delegate void SellRequestedEventHandler(string itemId, int quantity);

    private Label? _titleLabel;
    private Label? _summaryLabel;
    private ItemList? _itemList;
    private Label? _selectedItemLabel;
    private Label? _descriptionLabel;
    private Label? _priceLabel;
    private Label? _ownedLabel;
    private SpinBox? _quantitySpinBox;
    private Label? _statusLabel;
    private Button? _buyButton;
    private Button? _sellButton;
    private Button? _backButton;

    private TradeScreenViewModel _viewModel = new();
    private int _selectedIndex = -1;
    private bool _isViewReady;

    public override void _Ready()
    {
        _backdrop = GetNodeOrNull<TextureRect>("%Backdrop");
        _titleLabel = GetNodeOrNull<Label>("%TitleLabel");
        _summaryLabel = GetNodeOrNull<Label>("%SummaryLabel");
        _itemList = GetNodeOrNull<ItemList>("%ItemList");
        _selectedItemLabel = GetNodeOrNull<Label>("%SelectedItemLabel");
        _descriptionLabel = GetNodeOrNull<Label>("%DescriptionLabel");
        _priceLabel = GetNodeOrNull<Label>("%PriceLabel");
        _ownedLabel = GetNodeOrNull<Label>("%OwnedLabel");
        _quantitySpinBox = GetNodeOrNull<SpinBox>("%QuantitySpinBox");
        _statusLabel = GetNodeOrNull<Label>("%StatusLabel");
        _buyButton = GetNodeOrNull<Button>("%BuyButton");
        _sellButton = GetNodeOrNull<Button>("%SellButton");
        _backButton = GetNodeOrNull<Button>("%BackButton");

        if (_itemList is not null)
        {
            _itemList.ItemSelected += OnItemSelected;
        }

        if (_buyButton is not null)
        {
            _buyButton.Pressed += OnBuyPressed;
        }

        if (_sellButton is not null)
        {
            _sellButton.Pressed += OnSellPressed;
        }

        if (_backButton is not null)
        {
            _backButton.Pressed += () => EmitSignal(SignalName.BackRequested);
        }

        _isViewReady = true;
        ApplyViewState();
    }

    public void Bind(TradeScreenViewModel viewModel)
    {
        _viewModel = viewModel;
        if (_selectedIndex < 0 || _selectedIndex >= viewModel.Items.Count)
        {
            _selectedIndex = viewModel.Items.Count > 0 ? 0 : -1;
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
            _titleLabel.Text = $"Trading At {_viewModel.PortName}";
        }

        if (_backdrop is not null)
        {
            _backdrop.Texture = string.IsNullOrWhiteSpace(_viewModel.BackgroundTexturePath)
                ? null
                : ResourceLoader.Load<Texture2D>(_viewModel.BackgroundTexturePath);
        }

        if (_summaryLabel is not null)
        {
            _summaryLabel.Text =
                $"Credits: {_viewModel.Credits}    Cargo: {_viewModel.CargoLoad}/{_viewModel.CargoLimit}";
        }

        if (_itemList is not null)
        {
            _itemList.Clear();
            foreach (TradeItemViewModel item in _viewModel.Items)
            {
                _itemList.AddItem($"{item.Name}  |  {item.Price} cr  |  owned {item.OwnedQuantity}");
            }

            if (_selectedIndex >= 0)
            {
                _itemList.Select(_selectedIndex);
            }
        }

        if (_statusLabel is not null)
        {
            _statusLabel.Text = _viewModel.StatusMessage;
        }

        RefreshSelectionDetails();
    }

    private void OnItemSelected(long index)
    {
        _selectedIndex = (int)index;
        RefreshSelectionDetails();
    }

    private void OnBuyPressed()
    {
        TradeItemViewModel? item = GetSelectedItem();
        if (item is null || _quantitySpinBox is null)
        {
            return;
        }

        EmitSignal(SignalName.BuyRequested, item.ItemId, (int)_quantitySpinBox.Value);
    }

    private void OnSellPressed()
    {
        TradeItemViewModel? item = GetSelectedItem();
        if (item is null || _quantitySpinBox is null)
        {
            return;
        }

        EmitSignal(SignalName.SellRequested, item.ItemId, (int)_quantitySpinBox.Value);
    }

    private void RefreshSelectionDetails()
    {
        TradeItemViewModel? item = GetSelectedItem();
        if (item is null)
        {
            if (_selectedItemLabel is not null)
            {
                _selectedItemLabel.Text = "No market items available";
            }

            if (_descriptionLabel is not null)
            {
                _descriptionLabel.Text = string.Empty;
            }

            if (_priceLabel is not null)
            {
                _priceLabel.Text = string.Empty;
            }

            if (_ownedLabel is not null)
            {
                _ownedLabel.Text = string.Empty;
            }

            return;
        }

        if (_selectedItemLabel is not null)
        {
            _selectedItemLabel.Text = item.Name;
        }

        if (_descriptionLabel is not null)
        {
            _descriptionLabel.Text = item.Description;
        }

        if (_priceLabel is not null)
        {
            _priceLabel.Text = $"Current price: {item.Price} credits";
        }

        if (_ownedLabel is not null)
        {
            _ownedLabel.Text = $"Owned: {item.OwnedQuantity}";
        }
    }

    private TradeItemViewModel? GetSelectedItem()
    {
        return _selectedIndex >= 0 && _selectedIndex < _viewModel.Items.Count
            ? _viewModel.Items[_selectedIndex]
            : null;
    }
}
