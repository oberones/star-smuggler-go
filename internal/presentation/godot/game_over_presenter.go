package godot

import (
	"fmt"

	"github.com/oberones/star-smuggler-go/internal/domain"
	"github.com/oberones/star-smuggler-go/internal/services"
)

type GameOverViewModel struct {
	Summary string
}

type GameOverPresenter struct {
	Data    domain.DataSnapshot
	Economy services.EconomyService
	Travel  services.TravelService
}

func (p GameOverPresenter) Present(run domain.RunState) (GameOverViewModel, error) {
	port, ok := p.Data.PortsByID[run.Player.CurrentPortID]
	if !ok {
		return GameOverViewModel{
			Summary: "This run ended, but the current port could not be resolved.",
		}, nil
	}

	cargoValue := p.Economy.GetSellableCargoValueAtCurrentPort(run, p.Data)
	cheapestTravel := p.Travel.GetCheapestTravelCostFromPort(port, p.Data.Ports)

	return GameOverViewModel{
		Summary: fmt.Sprintf(
			"You are stranded at %s.\n\nCredits: %d\nSellable cargo value: %d\nCheapest travel cost: %d\n\nNo route remains that your current cash and cargo can cover.",
			port.Name,
			run.Player.Credits,
			cargoValue,
			cheapestTravel,
		),
	}, nil
}
