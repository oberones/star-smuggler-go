package godot

import (
	"github.com/oberones/star-smuggler-go/internal/domain"
	"github.com/oberones/star-smuggler-go/internal/services"
)

type MainMenuScreen interface {
	SetCanContinue(bool)
	SetStatusMessage(string)
}

type PortOverviewScreen interface {
	Bind(PortOverviewViewModel)
}

type TradeScreen interface {
	Bind(TradeScreenViewModel)
}

type TravelScreen interface {
	Bind(TravelScreenViewModel)
}

type TravelAnimationScreen interface {
	Bind(TravelAnimationViewModel)
}

type GameOverScreen interface {
	SetSummary(string)
	SetRecoveryState(bool, string)
}

type SceneBindings struct {
	Resources    *ResourceCache
	Audio        AudioBridge
	MainMenu     MainMenuPresenter
	PortOverview PortOverviewPresenter
	Trade        TradePresenter
	Travel       TravelPresenter
	TravelAnim   TravelAnimationPresenter
	GameOver     GameOverPresenter
}

func NewSceneBindings(
	data domain.DataSnapshot,
	economy services.EconomyService,
	travel services.TravelService,
	balance services.EconomyBalanceService,
	runEval services.RunEvaluator,
) SceneBindings {
	resources := NewResourceCache()
	storyPresenter := StoryPresenter{
		Data:     data,
		Factions: services.FactionService{},
		Missions: services.MissionService{},
	}
	progressionPresenter := ProgressionPresenter{
		Data:     data,
		Factions: services.FactionService{},
		Upgrades: services.UpgradeService{},
	}

	return SceneBindings{
		Resources: resources,
		Audio: AudioBridge{
			Cache: resources,
		},
		MainMenu: MainMenuPresenter{},
		PortOverview: PortOverviewPresenter{
			Data:        data,
			Economy:     economy,
			Travel:      travel,
			RunEval:     runEval,
			Story:       storyPresenter,
			Progression: progressionPresenter,
			Resources:   resources,
		},
		Trade: TradePresenter{
			Data:        data,
			Economy:     economy,
			Story:       storyPresenter,
			Progression: progressionPresenter,
			Resources:   resources,
		},
		Travel: TravelPresenter{
			Data:      data,
			Travel:    travel,
			Balance:   balance,
			Resources: resources,
		},
		TravelAnim: TravelAnimationPresenter{
			Data:      data,
			Resources: resources,
		},
		GameOver: GameOverPresenter{
			Data:    data,
			Economy: economy,
			Travel:  travel,
		},
	}
}

func (b SceneBindings) BindMainMenu(screen MainMenuScreen, canContinue bool, statusOverride string) {
	viewModel := b.MainMenu.Present(canContinue, statusOverride)
	screen.SetCanContinue(viewModel.CanContinue)
	screen.SetStatusMessage(viewModel.StatusMessage)
}

func (b SceneBindings) BindPortOverview(screen PortOverviewScreen, run domain.RunState, statusOverride string) error {
	viewModel, err := b.PortOverview.Present(run, statusOverride)
	if err != nil {
		return err
	}

	screen.Bind(viewModel)
	return nil
}

func (b SceneBindings) BindTrade(screen TradeScreen, run domain.RunState, statusOverride string) error {
	viewModel, err := b.Trade.Present(run, statusOverride)
	if err != nil {
		return err
	}

	screen.Bind(viewModel)
	return nil
}

func (b SceneBindings) BindTravel(screen TravelScreen, run domain.RunState, quotes []services.TravelQuote, statusOverride string) error {
	viewModel, err := b.Travel.Present(run, quotes, statusOverride)
	if err != nil {
		return err
	}

	screen.Bind(viewModel)
	return nil
}

func (b SceneBindings) BindTravelAnimation(screen TravelAnimationScreen, run domain.RunState, durationSeconds float64, statusOverride string) error {
	viewModel, err := b.TravelAnim.Present(run, durationSeconds, statusOverride)
	if err != nil {
		return err
	}

	screen.Bind(viewModel)
	return nil
}

func (b SceneBindings) BindGameOver(screen GameOverScreen, run domain.RunState) error {
	viewModel, err := b.GameOver.Present(run)
	if err != nil {
		return err
	}

	screen.SetSummary(viewModel.Summary)
	screen.SetRecoveryState(viewModel.CanRecover, viewModel.RecoveryStatus)
	return nil
}
