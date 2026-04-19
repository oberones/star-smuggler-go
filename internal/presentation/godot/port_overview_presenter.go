package godot

import (
	"fmt"

	"github.com/oberones/star-smuggler-go/internal/domain"
	"github.com/oberones/star-smuggler-go/internal/services"
)

type PortOverviewViewModel struct {
	PortName              string
	PortDescription       string
	ZoneName              string
	BackgroundTexturePath string
	MusicTrackID          string
	Credits               int
	CargoLoad             int
	CargoLimit            int
	CheapestTravelCost    int
	IsGameOver            bool
	RecentEventText       string
	AvailableGoods        []string
	StoryNotices          []string
	ProgressionNotices    []string
	StatusMessage         string
}

type PortOverviewPresenter struct {
	Data        domain.DataSnapshot
	Economy     services.EconomyService
	Travel      services.TravelService
	RunEval     services.RunEvaluator
	Story       StoryPresenter
	Progression ProgressionPresenter
}

func (p PortOverviewPresenter) Present(run domain.RunState, statusOverride string) (PortOverviewViewModel, error) {
	port, ok := p.Data.PortsByID[run.Player.CurrentPortID]
	if !ok {
		return PortOverviewViewModel{}, fmt.Errorf("current port %q was not found", run.Player.CurrentPortID)
	}

	goods := p.Economy.GetAvailableGoodsForCurrentPort(run, p.Data)
	availableGoods := make([]string, 0, len(goods))
	for _, item := range goods {
		availableGoods = append(availableGoods, fmt.Sprintf("%s (%d base)", item.Name, item.BasePrice))
	}

	isGameOver := p.RunEval.IsGameOver(run, p.Data, p.Economy, p.Travel)
	cheapestTravelCost := p.Travel.GetCheapestTravelCostFromPort(port, p.Data.Ports)
	storyViewModel := p.Story.Present(run)
	progressionViewModel := p.Progression.Present(run)
	statusMessage := statusOverride
	if statusMessage == "" {
		if isGameOver {
			statusMessage = "Status: Stranded. This run would currently evaluate as game over."
		} else {
			statusMessage = "Status: Operational. The trading loop can continue from here."
		}

		if run.RecentEvent != nil && run.RecentEvent.ResolvedDescription != "" {
			statusMessage += "\nRecent event outcome: " + run.RecentEvent.ResolvedDescription
		}
		if len(storyViewModel.MissionNotices) > 0 {
			statusMessage += "\n" + storyViewModel.MissionNotices[0]
		} else if len(progressionViewModel.AvailableUpgradeNotices) > 0 {
			statusMessage += "\n" + progressionViewModel.AvailableUpgradeNotices[0]
		}
	}

	return PortOverviewViewModel{
		PortName:              port.Name,
		PortDescription:       port.Description,
		ZoneName:              string(port.Zone),
		BackgroundTexturePath: port.BackgroundTexturePath,
		MusicTrackID:          port.MusicTrackID,
		Credits:               run.Player.Credits,
		CargoLoad:             p.Economy.GetCargoLoad(run),
		CargoLimit:            run.Player.CargoLimit,
		CheapestTravelCost:    cheapestTravelCost,
		IsGameOver:            isGameOver,
		RecentEventText:       recentEventDescription(run.RecentEvent),
		AvailableGoods:        availableGoods,
		StoryNotices:          append(append([]string{}, storyViewModel.FactionNotices...), append(storyViewModel.MissionNotices, storyViewModel.StoryNotices...)...),
		ProgressionNotices:    append(append([]string{}, progressionViewModel.OwnedUpgradeNotices...), append(progressionViewModel.AvailableUpgradeNotices, progressionViewModel.SpecializationNotices...)...),
		StatusMessage:         statusMessage,
	}, nil
}

func recentEventDescription(event *domain.EventResult) string {
	if event == nil {
		return ""
	}

	return event.ResolvedDescription
}
