package godot

import (
	"fmt"

	"github.com/oberones/star-smuggler-go/internal/domain"
	"github.com/oberones/star-smuggler-go/internal/services"
)

type TradeItemViewModel struct {
	ItemID        string
	Name          string
	Description   string
	Price         int
	OwnedQuantity int
}

type TradeScreenViewModel struct {
	PortName              string
	BackgroundTexturePath string
	MusicTrackID          string
	Credits               int
	CargoLoad             int
	CargoLimit            int
	Items                 []TradeItemViewModel
	StoryNotices          []string
	ProgressionNotices    []string
	StatusMessage         string
}

type TradePresenter struct {
	Data        domain.DataSnapshot
	Economy     services.EconomyService
	Story       StoryPresenter
	Progression ProgressionPresenter
	Resources   *ResourceCache
}

func (p TradePresenter) Present(run domain.RunState, statusOverride string) (TradeScreenViewModel, error) {
	port, ok := p.Data.PortsByID[run.Player.CurrentPortID]
	if !ok {
		return TradeScreenViewModel{}, fmt.Errorf("current port %q was not found", run.Player.CurrentPortID)
	}

	market, ok := p.Economy.GetCurrentMarket(run)
	if !ok {
		return TradeScreenViewModel{}, fmt.Errorf("current market is not available")
	}

	items := make([]TradeItemViewModel, 0, len(market.AvailableItemIDs))
	for _, itemID := range market.AvailableItemIDs {
		item, ok := p.Data.ItemsByID[itemID]
		if !ok {
			continue
		}

		price := item.BasePrice
		if currentPrice, exists := market.PricesByItemID[itemID]; exists {
			price = currentPrice
		}

		items = append(items, TradeItemViewModel{
			ItemID:        item.ID,
			Name:          item.Name,
			Description:   item.Description,
			Price:         price,
			OwnedQuantity: run.Cargo.QuantityFor(item.ID),
		})
	}

	statusMessage := statusOverride
	if statusMessage == "" {
		statusMessage = "Select a good, choose a quantity, and trade."
	}
	storyViewModel := p.Story.Present(run)
	progressionViewModel := p.Progression.Present(run)
	if statusOverride == "" && len(storyViewModel.MissionNotices) > 0 {
		statusMessage += "\n" + storyViewModel.MissionNotices[0]
	} else if statusOverride == "" && len(progressionViewModel.AvailableUpgradeNotices) > 0 {
		statusMessage += "\n" + progressionViewModel.AvailableUpgradeNotices[0]
	}

	return TradeScreenViewModel{
		PortName:              port.Name,
		BackgroundTexturePath: p.resolveTexture(port.TradeBackgroundPath),
		MusicTrackID:          p.resolveMusic(port.MusicTrackID),
		Credits:               run.Player.Credits,
		CargoLoad:             p.Economy.GetCargoLoad(run),
		CargoLimit:            run.Player.CargoLimit,
		Items:                 items,
		StoryNotices:          append(append([]string{}, storyViewModel.FactionNotices...), append(storyViewModel.MissionNotices, storyViewModel.StoryNotices...)...),
		ProgressionNotices:    append(append([]string{}, progressionViewModel.OwnedUpgradeNotices...), append(progressionViewModel.AvailableUpgradeNotices, progressionViewModel.SpecializationNotices...)...),
		StatusMessage:         statusMessage,
	}, nil
}

func (p TradePresenter) resolveTexture(path string) string {
	if p.Resources == nil {
		return path
	}
	return p.Resources.ResolveTexture(path)
}

func (p TradePresenter) resolveMusic(trackID string) string {
	if p.Resources == nil {
		return trackID
	}
	return p.Resources.ResolveMusic(trackID)
}
