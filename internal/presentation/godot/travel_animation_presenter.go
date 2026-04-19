package godot

import (
	"fmt"

	"github.com/oberones/star-smuggler-go/internal/domain"
)

const travelAnimationBackgroundTexturePath = "res://assets/screens/travel_background.png"

type TravelAnimationViewModel struct {
	OriginName            string
	DestinationName       string
	BackgroundTexturePath string
	TravelCost            int
	DurationSeconds       float64
	StatusMessage         string
}

type TravelAnimationPresenter struct {
	Data domain.DataSnapshot
}

func (p TravelAnimationPresenter) Present(run domain.RunState, durationSeconds float64, statusOverride string) (TravelAnimationViewModel, error) {
	if run.PendingRoute == nil {
		return TravelAnimationViewModel{}, fmt.Errorf("there is no pending route to animate")
	}

	origin, ok := p.Data.PortsByID[run.PendingRoute.OriginPortID]
	if !ok {
		return TravelAnimationViewModel{}, fmt.Errorf("origin port %q was not found", run.PendingRoute.OriginPortID)
	}

	destination, ok := p.Data.PortsByID[run.PendingRoute.DestinationPortID]
	if !ok {
		return TravelAnimationViewModel{}, fmt.Errorf("destination port %q was not found", run.PendingRoute.DestinationPortID)
	}

	statusMessage := statusOverride
	if statusMessage == "" {
		statusMessage = "Engines hot. Hold course or skip once you're ready."
	}

	if durationSeconds <= 0 {
		durationSeconds = 2.5
	}

	return TravelAnimationViewModel{
		OriginName:            origin.Name,
		DestinationName:       destination.Name,
		BackgroundTexturePath: travelAnimationBackgroundTexturePath,
		TravelCost:            run.PendingRoute.TravelCost,
		DurationSeconds:       durationSeconds,
		StatusMessage:         statusMessage,
	}, nil
}
