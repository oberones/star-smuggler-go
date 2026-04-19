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
	StatusMessage         string
}

type TradePresenter struct {
	Data    domain.DataSnapshot
	Economy services.EconomyService
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

	return TradeScreenViewModel{
		PortName:              port.Name,
		BackgroundTexturePath: port.TradeBackgroundPath,
		MusicTrackID:          port.MusicTrackID,
		Credits:               run.Player.Credits,
		CargoLoad:             p.Economy.GetCargoLoad(run),
		CargoLimit:            run.Player.CargoLimit,
		Items:                 items,
		StatusMessage:         statusMessage,
	}, nil
}
