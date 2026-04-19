package godot

import (
	"fmt"

	"github.com/oberones/star-smuggler-go/internal/domain"
	"github.com/oberones/star-smuggler-go/internal/services"
)

const cockpitBackgroundTexturePath = "res://assets/ui/cockpit.png"

type TravelDestinationViewModel struct {
	PortID             string
	Name               string
	ZoneName           string
	Description        string
	PreviewTexturePath string
	TravelCost         int
}

type TravelScreenViewModel struct {
	CurrentPortName       string
	BackgroundTexturePath string
	Credits               int
	Destinations          []TravelDestinationViewModel
	StatusMessage         string
}

type TravelPresenter struct {
	Data    domain.DataSnapshot
	Travel  services.TravelService
	Balance services.EconomyBalanceService
}

func (p TravelPresenter) Present(run domain.RunState, quotes []services.TravelQuote, statusOverride string) (TravelScreenViewModel, error) {
	currentPort, ok := p.Data.PortsByID[run.Player.CurrentPortID]
	if !ok {
		return TravelScreenViewModel{}, fmt.Errorf("current port %q was not found", run.Player.CurrentPortID)
	}

	destinations := make([]TravelDestinationViewModel, 0, len(quotes))
	for _, quote := range quotes {
		destinations = append(destinations, TravelDestinationViewModel{
			PortID:             quote.Destination.ID,
			Name:               quote.Destination.Name,
			ZoneName:           string(quote.Destination.Zone),
			Description:        quote.Destination.Description,
			PreviewTexturePath: quote.Destination.PreviewTexturePath,
			TravelCost:         quote.Cost,
		})
	}

	statusMessage := statusOverride
	if statusMessage == "" {
		statusMessage = "Choose your next destination."
	}

	return TravelScreenViewModel{
		CurrentPortName:       currentPort.Name,
		BackgroundTexturePath: cockpitBackgroundTexturePath,
		Credits:               run.Player.Credits,
		Destinations:          destinations,
		StatusMessage:         statusMessage,
	}, nil
}
