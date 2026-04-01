using Godot;
using StarSmugglerGo.Autoload;
using StarSmugglerGo.Presentation.Screens;
using StarSmugglerGo.Domain;
using StarSmugglerGo.Services;
using System;
using System.Collections.Generic;
using System.Linq;

namespace StarSmugglerGo.Application;

public partial class AppController : Node
{
    private const string MainMenuScenePath = "res://scenes/screens/MainMenuScreen.tscn";
    private const string PortOverviewScenePath = "res://scenes/screens/PortOverviewScreen.tscn";
    private const string TradeScenePath = "res://scenes/screens/TradeScreen.tscn";
    private const string TravelScenePath = "res://scenes/screens/TravelScreen.tscn";
    private const string GameOverScenePath = "res://scenes/screens/GameOverScreen.tscn";
    private const string TravelAnimationScenePath = "res://scenes/screens/TravelAnimationScreen.tscn";

    private readonly EconomyService _economyService = new();
    private readonly TravelService _travelService = new();
    private readonly RunEvaluator _runEvaluator = new();
    private readonly TradeService _tradeService = new();
    private readonly EventService _eventService = new();
    private readonly Random _random = new();

    [Signal]
    public delegate void RouteChangedEventHandler(string routeName);

    private Control? _screenHost;
    private GameSession? _gameSession;
    private DataRepository? _dataRepository;
    private SaveService? _saveService;
    private AudioService? _audioService;
    private Control? _currentScreen;
    private PendingTravelRequest? _pendingTravelRequest;

    public AppRoute CurrentRoute { get; private set; } = AppRoute.None;

    public override void _Ready()
    {
        _screenHost = GetNodeOrNull<Control>("%ScreenHost");
        _gameSession = GetNodeOrNull<GameSession>("%GameSession");
        _dataRepository = GetNodeOrNull<DataRepository>("%DataRepository");
        _saveService = GetNodeOrNull<SaveService>("%SaveService");
        _audioService = GetNodeOrNull<AudioService>("%AudioService");

        if (_screenHost is null)
        {
            GD.PushError("AppController could not find %ScreenHost.");
            return;
        }

        NavigateTo(AppRoute.MainMenu);
    }

    public void NavigateTo(AppRoute route)
    {
        if (_screenHost is null)
        {
            GD.PushError("AppController cannot navigate without a screen host.");
            return;
        }

        Control? nextScreen = route switch
        {
            AppRoute.MainMenu => CreateMainMenuScreen(),
            AppRoute.PortOverview => CreatePortOverviewScreen(),
            AppRoute.Trade => CreateTradeScreen(),
            AppRoute.Travel => CreateTravelScreen(),
            AppRoute.GameOver => CreateGameOverScreen(),
            AppRoute.TravelAnimation => CreateTravelAnimationScreen(),
            _ => null,
        };

        if (nextScreen is null)
        {
            GD.PushError($"AppController could not build route '{route}'.");
            return;
        }

        _currentScreen?.QueueFree();
        _currentScreen = nextScreen;
        _screenHost.AddChild(_currentScreen);
        CurrentRoute = route;
        PlayRouteMusic(route);
        EmitSignal(SignalName.RouteChanged, route.ToString());
    }

    private Control? CreateMainMenuScreen()
    {
        PackedScene? scene = ResourceLoader.Load<PackedScene>(MainMenuScenePath);
        if (scene?.Instantiate() is not MainMenuScreen screen)
        {
            GD.PushError($"AppController failed to load '{MainMenuScenePath}'.");
            return null;
        }

        screen.StartRequested += OnStartRequested;
        screen.ContinueRequested += OnContinueRequested;
        screen.QuitRequested += OnQuitRequested;
        screen.SetCanContinue(_saveService?.HasSave() ?? false);
        screen.SetStatusMessage((_saveService?.HasSave() ?? false)
            ? "Resume a saved smuggling run or start a fresh route."
            : "Start a new run to begin rebuilding StarSmuggler in Godot.");

        return screen;
    }

    private void OnStartRequested()
    {
        _gameSession?.StartNewRun();
        _audioService?.PlaySfx("click");

        if (_gameSession?.CurrentRun is RunState)
        {
            NavigateTo(AppRoute.PortOverview);
        }
    }

    private void OnQuitRequested()
    {
        _audioService?.PlaySfx("click");
        GetTree().Quit();
    }

    private void OnContinueRequested()
    {
        _audioService?.PlaySfx("click");
        if (_gameSession?.TryLoadSavedRun() == true)
        {
            NavigateTo(RouteForCurrentRun());
        }
    }

    private Control? CreatePortOverviewScreen()
    {
        if (_gameSession?.CurrentRun is not RunState run || _dataRepository is null)
        {
            GD.PushError("AppController cannot open PortOverview without an active run and data repository.");
            return null;
        }

        PackedScene? scene = ResourceLoader.Load<PackedScene>(PortOverviewScenePath);
        if (scene?.Instantiate() is not PortOverviewScreen screen)
        {
            GD.PushError($"AppController failed to load '{PortOverviewScenePath}'.");
            return null;
        }

        PortOverviewViewModel? viewModel = BuildPortOverviewViewModel(run, _dataRepository.Snapshot);
        if (viewModel is null)
        {
            return null;
        }

        screen.Bind(viewModel);
        screen.BackRequested += () =>
        {
            _audioService?.PlaySfx("click");
            NavigateTo(AppRoute.MainMenu);
        };
        screen.TravelRequested += () =>
        {
            _audioService?.PlaySfx("click");
            NavigateTo(AppRoute.Travel);
        };
        screen.TradeRequested += () =>
        {
            _audioService?.PlaySfx("click");
            NavigateTo(AppRoute.Trade);
        };

        return screen;
    }

    private Control? CreateTradeScreen()
    {
        if (_gameSession?.CurrentRun is not RunState run || _dataRepository is null)
        {
            GD.PushError("AppController cannot open Trade without an active run and data repository.");
            return null;
        }

        PackedScene? scene = ResourceLoader.Load<PackedScene>(TradeScenePath);
        if (scene?.Instantiate() is not TradeScreen screen)
        {
            GD.PushError($"AppController failed to load '{TradeScenePath}'.");
            return null;
        }

        screen.Bind(BuildTradeScreenViewModel(run, _dataRepository.Snapshot, string.Empty));
        screen.BackRequested += () =>
        {
            _audioService?.PlaySfx("click");
            NavigateTo(AppRoute.PortOverview);
        };
        screen.BuyRequested += OnBuyRequested;
        screen.SellRequested += OnSellRequested;

        return screen;
    }

    private Control? CreateTravelScreen()
    {
        if (_gameSession?.CurrentRun is not RunState run || _dataRepository is null)
        {
            GD.PushError("AppController cannot open Travel without an active run and data repository.");
            return null;
        }

        PackedScene? scene = ResourceLoader.Load<PackedScene>(TravelScenePath);
        if (scene?.Instantiate() is not TravelScreen screen)
        {
            GD.PushError($"AppController failed to load '{TravelScenePath}'.");
            return null;
        }

        screen.Bind(BuildTravelScreenViewModel(run, _dataRepository.Snapshot, string.Empty));
        screen.BackRequested += () =>
        {
            _audioService?.PlaySfx("click");
            NavigateTo(AppRoute.PortOverview);
        };
        screen.TravelRequested += OnTravelRequested;

        return screen;
    }

    private Control? CreateGameOverScreen()
    {
        if (_gameSession?.CurrentRun is not RunState run || _dataRepository is null)
        {
            GD.PushError("AppController cannot open GameOver without an active run and data repository.");
            return null;
        }

        PackedScene? scene = ResourceLoader.Load<PackedScene>(GameOverScenePath);
        if (scene?.Instantiate() is not GameOverScreen screen)
        {
            GD.PushError($"AppController failed to load '{GameOverScenePath}'.");
            return null;
        }

        screen.SetSummary(BuildGameOverSummary(run, _dataRepository.Snapshot));
        screen.RestartRequested += OnStartRequested;
        screen.MenuRequested += () =>
        {
            _audioService?.PlaySfx("click");
            NavigateTo(AppRoute.MainMenu);
        };
        return screen;
    }

    private Control? CreateTravelAnimationScreen()
    {
        if (_pendingTravelRequest is null || _dataRepository is null)
        {
            GD.PushError("AppController cannot open TravelAnimation without a pending travel request.");
            return null;
        }

        if (!_dataRepository.Snapshot.PortsById.TryGetValue(_pendingTravelRequest.OriginPortId, out PortDefinition? origin) ||
            !_dataRepository.Snapshot.PortsById.TryGetValue(_pendingTravelRequest.DestinationPortId, out PortDefinition? destination))
        {
            GD.PushError("AppController could not resolve travel animation endpoints.");
            return null;
        }

        PackedScene? scene = ResourceLoader.Load<PackedScene>(TravelAnimationScenePath);
        if (scene?.Instantiate() is not TravelAnimationScreen screen)
        {
            GD.PushError($"AppController failed to load '{TravelAnimationScenePath}'.");
            return null;
        }

        screen.Bind(new TravelAnimationViewModel
        {
            OriginName = origin.Name,
            DestinationName = destination.Name,
            BackgroundTexturePath = "res://assets/screens/travel_background.png",
            TravelCost = _pendingTravelRequest.TravelCost,
            DurationSeconds = _pendingTravelRequest.DurationSeconds,
            StatusMessage = "Engines hot. Hold course or skip once you're ready.",
        });
        screen.AnimationFinished += OnTravelAnimationFinished;
        return screen;
    }

    private PortOverviewViewModel? BuildPortOverviewViewModel(RunState run, DataSnapshot data)
    {
        if (!data.PortsById.TryGetValue(run.Player.CurrentPortId, out PortDefinition? port))
        {
            GD.PushError($"AppController could not find port definition '{run.Player.CurrentPortId}'.");
            return null;
        }

        var goods = _economyService.GetAvailableGoodsForCurrentPort(run, data);
        bool isGameOver = _runEvaluator.IsGameOver(run, data, _economyService, _travelService);
        int cheapestTravelCost = _travelService.GetCheapestTravelCostFromPort(port, data.Ports);

        return new PortOverviewViewModel
        {
            PortName = port.Name,
            PortDescription = port.Description,
            ZoneName = port.Zone.ToString(),
            BackgroundTexturePath = port.BackgroundTexturePath,
            MusicTrackId = port.MusicTrackId,
            Credits = run.Player.Credits,
            CargoLoad = _economyService.GetCargoLoad(run),
            CargoLimit = run.Player.CargoLimit,
            CheapestTravelCost = cheapestTravelCost,
            IsGameOver = isGameOver,
            RecentEventText = run.RecentEvent?.ResolvedDescription ?? string.Empty,
            AvailableGoods = goods
                .Select(item => $"{item.Name} ({item.BasePrice} base)")
                .ToList(),
        };
    }

    private TradeScreenViewModel BuildTradeScreenViewModel(RunState run, DataSnapshot data, string statusMessage)
    {
        PortDefinition port = data.PortsById[run.Player.CurrentPortId];
        MarketSnapshot? market = _economyService.GetCurrentMarket(run);

        var items = market is null
            ? new List<TradeItemViewModel>()
            : market.AvailableItemIds
                .Where(itemId => data.ItemsById.ContainsKey(itemId))
                .Select(itemId =>
                {
                    ItemDefinition item = data.ItemsById[itemId];
                    int price = market.PricesByItemId.TryGetValue(itemId, out int currentPrice)
                        ? currentPrice
                        : item.BasePrice;

                    return new TradeItemViewModel
                    {
                        ItemId = item.Id,
                        Name = item.Name,
                        Description = item.Description,
                        Price = price,
                        OwnedQuantity = run.Cargo.GetQuantity(item.Id),
                    };
                })
                .ToList();

        return new TradeScreenViewModel
        {
            PortName = port.Name,
            BackgroundTexturePath = port.TradeBackgroundPath,
            MusicTrackId = port.MusicTrackId,
            Credits = run.Player.Credits,
            CargoLoad = _economyService.GetCargoLoad(run),
            CargoLimit = run.Player.CargoLimit,
            Items = items,
            StatusMessage = string.IsNullOrWhiteSpace(statusMessage)
                ? "Select a good, choose a quantity, and trade."
                : statusMessage,
        };
    }

    private TravelScreenViewModel BuildTravelScreenViewModel(RunState run, DataSnapshot data, string statusMessage)
    {
        PortDefinition currentPort = data.PortsById[run.Player.CurrentPortId];
        IReadOnlyList<TravelDestinationViewModel> destinations = _travelService
            .GetDestinationsFromPort(currentPort, data.Ports)
            .Select(destination => new TravelDestinationViewModel
            {
                PortId = destination.Id,
                Name = destination.Name,
                ZoneName = destination.Zone.ToString(),
                Description = destination.Description,
                PreviewTexturePath = destination.PreviewTexturePath,
                TravelCost = _travelService.GetTravelCost(currentPort, destination),
            })
            .ToList();

        return new TravelScreenViewModel
        {
            CurrentPortName = currentPort.Name,
            BackgroundTexturePath = "res://assets/ui/cockpit.png",
            Credits = run.Player.Credits,
            Destinations = destinations,
            StatusMessage = string.IsNullOrWhiteSpace(statusMessage)
                ? "Choose your next destination."
                : statusMessage,
        };
    }

    private void OnBuyRequested(string itemId, int quantity)
    {
        ApplyTrade(itemId, quantity, isBuy: true);
    }

    private void OnSellRequested(string itemId, int quantity)
    {
        ApplyTrade(itemId, quantity, isBuy: false);
    }

    private void ApplyTrade(string itemId, int quantity, bool isBuy)
    {
        if (_gameSession?.CurrentRun is not RunState run || _dataRepository is null)
        {
            return;
        }

        if (!_dataRepository.Snapshot.ItemsById.TryGetValue(itemId, out ItemDefinition? item))
        {
            RefreshTradeScreen("That item no longer exists in the data repository.");
            return;
        }

        MarketSnapshot? market = _economyService.GetCurrentMarket(run);
        if (market is null)
        {
            RefreshTradeScreen("No market is loaded for the current port.");
            return;
        }

        TradeResult result = isBuy
            ? _tradeService.Buy(run, market, item, quantity)
            : _tradeService.Sell(run, market, item, quantity);

        _audioService?.PlaySfx("click");
        _gameSession.SaveCurrentRun();
        RefreshTradeScreen(result.Message);

        if (ShouldRouteToGameOver(run))
        {
            NavigateTo(AppRoute.GameOver);
        }
    }

    private void RefreshTradeScreen(string statusMessage)
    {
        if (_currentScreen is not TradeScreen tradeScreen ||
            _gameSession?.CurrentRun is not RunState run ||
            _dataRepository is null)
        {
            return;
        }

        tradeScreen.Bind(BuildTradeScreenViewModel(run, _dataRepository.Snapshot, statusMessage));
    }

    private void OnTravelRequested(string destinationPortId)
    {
        if (_gameSession?.CurrentRun is not RunState run || _dataRepository is null)
        {
            return;
        }

        DataSnapshot data = _dataRepository.Snapshot;
        if (!data.PortsById.TryGetValue(run.Player.CurrentPortId, out PortDefinition? origin) ||
            !data.PortsById.TryGetValue(destinationPortId, out PortDefinition? destination))
        {
            RefreshTravelScreen("That route is no longer valid.");
            return;
        }

        int cost = _travelService.GetTravelCost(origin, destination);
        if (run.Player.Credits < cost)
        {
            RefreshTravelScreen($"You need {cost} credits to reach {destination.Name}.");
            return;
        }

        int zoneDifference = Math.Abs((int)origin.Zone - (int)destination.Zone);
        _pendingTravelRequest = new PendingTravelRequest
        {
            OriginPortId = origin.Id,
            DestinationPortId = destination.Id,
            TravelCost = cost,
            DurationSeconds = 2.0 + (zoneDifference * 1.5),
        };

        _audioService?.PlaySfx("click");
        NavigateTo(AppRoute.TravelAnimation);
    }

    private void OnTravelAnimationFinished()
    {
        if (_gameSession?.CurrentRun is not RunState run || _dataRepository is null || _pendingTravelRequest is null)
        {
            return;
        }

        DataSnapshot data = _dataRepository.Snapshot;
        if (!data.PortsById.TryGetValue(_pendingTravelRequest.DestinationPortId, out PortDefinition? destination))
        {
            _pendingTravelRequest = null;
            NavigateTo(AppRoute.Travel);
            return;
        }

        run.Player.Credits -= _pendingTravelRequest.TravelCost;
        run.Player.CurrentPortId = destination.Id;
        run.JumpsSinceLastUpdate++;

        if (run.JumpsSinceLastUpdate > 3)
        {
            _economyService.RefreshAllPrices(run, data, _random);
        }
        else
        {
            _economyService.RefreshAvailableGoods(run, data, destination.Id, _random);
        }

        run.RecentEvent = _eventService.TryResolveTravelEvent(run, data, _economyService, _random);
        _gameSession.SaveCurrentRun();
        _pendingTravelRequest = null;
        NavigateTo(RouteForCurrentRun());
    }

    private void RefreshTravelScreen(string statusMessage)
    {
        if (_currentScreen is not TravelScreen travelScreen ||
            _gameSession?.CurrentRun is not RunState run ||
            _dataRepository is null)
        {
            return;
        }

        travelScreen.Bind(BuildTravelScreenViewModel(run, _dataRepository.Snapshot, statusMessage));
    }

    private AppRoute RouteForCurrentRun()
    {
        if (_gameSession?.CurrentRun is not RunState run)
        {
            return AppRoute.MainMenu;
        }

        return ShouldRouteToGameOver(run) ? AppRoute.GameOver : AppRoute.PortOverview;
    }

    private bool ShouldRouteToGameOver(RunState run)
    {
        return _dataRepository is not null &&
               _runEvaluator.IsGameOver(run, _dataRepository.Snapshot, _economyService, _travelService);
    }

    private string BuildGameOverSummary(RunState run, DataSnapshot data)
    {
        if (!data.PortsById.TryGetValue(run.Player.CurrentPortId, out PortDefinition? port))
        {
            return "This run ended, but the current port could not be resolved.";
        }

        int cargoValue = _economyService.GetSellableCargoValueAtCurrentPort(run, data);
        int cheapestTravel = _travelService.GetCheapestTravelCostFromPort(port, data.Ports);

        return
            $"You are stranded at {port.Name}.\n\n" +
            $"Credits: {run.Player.Credits}\n" +
            $"Sellable cargo value: {cargoValue}\n" +
            $"Cheapest travel cost: {cheapestTravel}\n\n" +
            $"No route remains that your current cash and cargo can cover.";
    }

    private void PlayRouteMusic(AppRoute route)
    {
        if (_audioService is null)
        {
            return;
        }

        string? trackId = route switch
        {
            AppRoute.MainMenu => "singularity",
            AppRoute.GameOver => "singularity",
            AppRoute.PortOverview or AppRoute.Trade => ResolveCurrentRunMusic(),
            AppRoute.Travel or AppRoute.TravelAnimation => "world_default",
            _ => null,
        };

        if (!string.IsNullOrWhiteSpace(trackId))
        {
            _audioService.PlayMusic(trackId);
        }
    }

    private string ResolveCurrentRunMusic()
    {
        if (_gameSession?.CurrentRun is not RunState run || _dataRepository is null)
        {
            return "singularity";
        }

        return _dataRepository.Snapshot.PortsById.TryGetValue(run.Player.CurrentPortId, out PortDefinition? port) &&
               !string.IsNullOrWhiteSpace(port.MusicTrackId)
            ? port.MusicTrackId
            : "world_default";
    }
}
