package application

import (
	"context"
	"fmt"

	"github.com/oberones/star-smuggler-go/internal/domain"
	"github.com/oberones/star-smuggler-go/internal/services"
)

func (a *App) StartNewRun(ctx context.Context) error {
	run, err := a.runCommands().StartNewRun()
	if err != nil {
		return err
	}

	a.activeRun = &run
	_ = a.storyCommands().SyncAtRunStart(a.activeRun)
	route, err := a.saveCommands().Autosave(ctx, run)
	if err != nil {
		return err
	}

	a.route = route
	return nil
}

func (a *App) ContinueSavedRun(ctx context.Context) error {
	run, route, err := a.saveCommands().ContinueSavedRun(ctx)
	if err != nil {
		return err
	}

	a.activeRun = &run
	_ = a.storyCommands().SyncAtRunStart(a.activeRun)
	a.route = route
	return nil
}

func (a *App) OpenPortOverview() error {
	if err := a.requireActiveRun(); err != nil {
		return err
	}

	a.route = a.saveCommands().RouteForRun(*a.activeRun)
	return nil
}

func (a *App) OpenTrade() error {
	if err := a.requireActiveRun(); err != nil {
		return err
	}

	a.route = RouteTrade
	return nil
}

func (a *App) OpenTravel() error {
	if err := a.requireActiveRun(); err != nil {
		return err
	}

	a.route = RouteTravel
	return nil
}

func (a *App) PreviewTravel() ([]services.TravelQuote, error) {
	if err := a.requireActiveRun(); err != nil {
		return nil, err
	}

	return a.runCommands().PreviewTravel(*a.activeRun)
}

func (a *App) Buy(ctx context.Context, itemID string, quantity int) (services.TradeResult, error) {
	if err := a.requireActiveRun(); err != nil {
		return services.TradeResult{}, err
	}

	result, err := a.runCommands().Buy(a.activeRun, itemID, quantity)
	if err != nil {
		return services.TradeResult{}, err
	}
	storyUpdate := a.storyCommands().SyncAfterTrade(a.activeRun, itemID, quantity, true)
	if summary := storyUpdate.Summary(); summary != "" {
		result.Message = result.Message + "\n" + summary
	}

	if err := a.autosaveWithPreferredRoute(ctx, RouteTrade); err != nil {
		return services.TradeResult{}, err
	}

	return result, nil
}

func (a *App) Sell(ctx context.Context, itemID string, quantity int) (services.TradeResult, error) {
	if err := a.requireActiveRun(); err != nil {
		return services.TradeResult{}, err
	}

	result, err := a.runCommands().Sell(a.activeRun, itemID, quantity)
	if err != nil {
		return services.TradeResult{}, err
	}
	storyUpdate := a.storyCommands().SyncAfterTrade(a.activeRun, itemID, quantity, false)
	if summary := storyUpdate.Summary(); summary != "" {
		result.Message = result.Message + "\n" + summary
	}

	if err := a.autosaveWithPreferredRoute(ctx, RouteTrade); err != nil {
		return services.TradeResult{}, err
	}

	return result, nil
}

func (a *App) BeginTravel(destinationPortID string) (*domain.RouteState, error) {
	if err := a.requireActiveRun(); err != nil {
		return nil, err
	}

	route, err := a.travelCommands().BeginTravel(a.activeRun, destinationPortID)
	if err != nil {
		return nil, err
	}

	a.route = RouteTravelAnimation
	return route, nil
}

func (a *App) ResolveTravel(ctx context.Context) (domain.TravelResolution, error) {
	if err := a.requireActiveRun(); err != nil {
		return domain.TravelResolution{}, err
	}

	resolution, err := a.travelCommands().ResolvePendingTravel(a.activeRun)
	if err != nil {
		return domain.TravelResolution{}, err
	}
	storyUpdate := a.storyCommands().SyncAfterTravel(a.activeRun)
	if summary := storyUpdate.Summary(); summary != "" {
		resolution.Message = resolution.Message + "\n" + summary
	}

	if err := a.autosaveWithPreferredRoute(ctx, RoutePortOverview); err != nil {
		return domain.TravelResolution{}, err
	}

	return resolution, nil
}

func (a *App) Recover(ctx context.Context) (string, bool, error) {
	if err := a.requireActiveRun(); err != nil {
		return "", false, err
	}

	message, recovered, err := a.recoveryCommands().TryEmergencyRecovery(a.activeRun)
	if err != nil {
		return "", false, err
	}

	if recovered {
		if err := a.autosaveWithPreferredRoute(ctx, RoutePortOverview); err != nil {
			return "", false, err
		}
		return message, true, nil
	}

	a.route = a.saveCommands().RouteForRun(*a.activeRun)
	return message, false, nil
}

func (a *App) ClearActiveRun() {
	a.activeRun = nil
	a.route = RouteMainMenu
}

func (a *App) AvailableMissions() ([]domain.MissionDefinition, error) {
	if err := a.requireActiveRun(); err != nil {
		return nil, err
	}

	return a.storyCommands().AvailableMissions(*a.activeRun), nil
}

func (a *App) AvailableUpgrades() ([]domain.ShipUpgradeDefinition, error) {
	if err := a.requireActiveRun(); err != nil {
		return nil, err
	}

	return a.progressionCommands().AvailableUpgrades(*a.activeRun), nil
}

func (a *App) AcceptMission(ctx context.Context, missionID string) (StoryUpdate, error) {
	if err := a.requireActiveRun(); err != nil {
		return StoryUpdate{}, err
	}

	update, err := a.storyCommands().AcceptMission(a.activeRun, missionID)
	if err != nil {
		return StoryUpdate{}, err
	}

	if err := a.autosaveWithPreferredRoute(ctx, RoutePortOverview); err != nil {
		return StoryUpdate{}, err
	}

	return update, nil
}

func (a *App) PurchaseUpgrade(ctx context.Context, upgradeID string) (services.UpgradePurchaseResult, error) {
	if err := a.requireActiveRun(); err != nil {
		return services.UpgradePurchaseResult{}, err
	}

	result := a.progressionCommands().PurchaseUpgrade(a.activeRun, upgradeID)
	if result.Succeeded {
		if err := a.autosaveWithPreferredRoute(ctx, RouteTrade); err != nil {
			return services.UpgradePurchaseResult{}, err
		}
	}

	return result, nil
}

func (a *App) runCommands() RunCommands {
	commands := NewRunCommands(a.snapshot, a.saveRepository, a.runtime)
	commands.Economy = services.EconomyService{}
	commands.Balance = services.EconomyBalanceService{}
	commands.Trade = services.TradeService{}
	commands.Travel = services.TravelService{}
	commands.Upgrades = services.UpgradeService{}
	commands.RunEval = services.RunEvaluator{}
	return commands
}

func (a *App) travelCommands() TravelCommands {
	commands := NewTravelCommands(a.snapshot, a.runtime)
	commands.Economy = services.EconomyService{}
	commands.Balance = services.EconomyBalanceService{}
	commands.Travel = services.TravelService{}
	commands.Upgrades = services.UpgradeService{}
	commands.Events = services.EventService{}
	commands.RunEval = services.RunEvaluator{}
	return commands
}

func (a *App) saveCommands() SaveCommands {
	commands := NewSaveCommands(a.snapshot, a.saveRepository)
	commands.Economy = services.EconomyService{}
	commands.Travel = services.TravelService{}
	commands.RunEval = services.RunEvaluator{}
	return commands
}

func (a *App) recoveryCommands() RecoveryCommands {
	commands := NewRecoveryCommands(a.snapshot)
	commands.Economy = services.EconomyService{}
	commands.Travel = services.TravelService{}
	commands.RunEval = services.RunEvaluator{}
	return commands
}

func (a *App) storyCommands() StoryCommands {
	commands := NewStoryCommands(a.snapshot)
	commands.Factions = services.FactionService{}
	commands.Missions = services.MissionService{}
	commands.Stories = services.StoryService{}
	commands.Upgrades = services.UpgradeService{}
	return commands
}

func (a *App) progressionCommands() ProgressionCommands {
	commands := NewProgressionCommands(a.snapshot)
	commands.Factions = services.FactionService{}
	commands.Upgrades = services.UpgradeService{}
	return commands
}

func (a *App) autosaveWithPreferredRoute(ctx context.Context, preferred Route) error {
	if a.activeRun == nil {
		return fmt.Errorf("there is no active run to save")
	}

	route, err := a.saveCommands().Autosave(ctx, *a.activeRun)
	if err != nil {
		return err
	}

	if route == RouteGameOver {
		a.route = RouteGameOver
		return nil
	}

	a.route = preferred
	return nil
}

func (a *App) requireActiveRun() error {
	if a.activeRun == nil {
		return fmt.Errorf("there is no active run")
	}

	return nil
}
