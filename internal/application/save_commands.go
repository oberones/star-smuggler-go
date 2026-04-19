package application

import (
	"context"
	"fmt"

	"github.com/oberones/star-smuggler-go/internal/domain"
	"github.com/oberones/star-smuggler-go/internal/services"
)

type SaveCommands struct {
	SaveRepo SaveRepository
	Data     domain.DataSnapshot
	Economy  services.EconomyService
	Travel   services.TravelService
	RunEval  services.RunEvaluator
}

func NewSaveCommands(data domain.DataSnapshot, saveRepo SaveRepository) SaveCommands {
	return SaveCommands{
		SaveRepo: saveRepo,
		Data:     data,
	}
}

func (c SaveCommands) HasContinue() (bool, error) {
	if c.SaveRepo == nil {
		return false, fmt.Errorf("save repository is not configured")
	}

	return c.SaveRepo.Exists()
}

func (c SaveCommands) ContinueSavedRun(ctx context.Context) (domain.RunState, Route, error) {
	if c.SaveRepo == nil {
		return domain.RunState{}, RouteNone, fmt.Errorf("save repository is not configured")
	}

	run, err := c.SaveRepo.Load(ctx)
	if err != nil {
		return domain.RunState{}, RouteNone, err
	}

	return run, c.RouteForRun(run), nil
}

func (c SaveCommands) Autosave(ctx context.Context, run domain.RunState) (Route, error) {
	if c.SaveRepo == nil {
		return RouteNone, fmt.Errorf("save repository is not configured")
	}

	if err := c.SaveRepo.Save(ctx, run); err != nil {
		return RouteNone, err
	}

	return c.RouteForRun(run), nil
}

func (c SaveCommands) RouteForRun(run domain.RunState) Route {
	if c.RunEval.IsGameOver(run, c.Data, c.Economy, c.Travel) {
		return RouteGameOver
	}

	return RoutePortOverview
}
