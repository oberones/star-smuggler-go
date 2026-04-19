package application

import (
	"context"
	"errors"

	"github.com/oberones/star-smuggler-go/internal/domain"
	"github.com/oberones/star-smuggler-go/internal/services"
)

type Route int

const (
	RouteNone Route = iota
	RouteMainMenu
	RoutePortOverview
	RouteTrade
	RouteTravel
	RouteGameOver
	RouteTravelAnimation
)

type ContentLoader interface {
	LoadSnapshot(ctx context.Context) (domain.DataSnapshot, error)
}

type SaveRepository interface {
	Exists() (bool, error)
	Load(ctx context.Context) (domain.RunState, error)
	Save(ctx context.Context, run domain.RunState) error
}

type App struct {
	contentLoader  ContentLoader
	saveRepository SaveRepository
	runtime        services.RuntimeContext
	route          Route
	snapshot       domain.DataSnapshot
}

func NewApp(contentLoader ContentLoader, saveRepository SaveRepository, runtime services.RuntimeContext) *App {
	return &App{
		contentLoader:  contentLoader,
		saveRepository: saveRepository,
		runtime:        runtime,
		route:          RouteNone,
	}
}

func (a *App) Bootstrap(ctx context.Context) error {
	if a.contentLoader == nil {
		return errors.New("content loader is not configured")
	}

	snapshot, err := a.contentLoader.LoadSnapshot(ctx)
	if err != nil {
		return err
	}

	a.snapshot = snapshot
	a.route = RouteMainMenu
	return nil
}

func (a *App) CurrentRoute() Route {
	return a.route
}

func (a *App) Navigate(route Route) {
	a.route = route
}

func (a *App) Snapshot() domain.DataSnapshot {
	return a.snapshot
}

func (a *App) Save(ctx context.Context, run domain.RunState) error {
	if a.saveRepository == nil {
		return errors.New("save repository is not configured")
	}
	return a.saveRepository.Save(ctx, run)
}

func (a *App) LoadSavedRun(ctx context.Context) (domain.RunState, error) {
	if a.saveRepository == nil {
		return domain.RunState{}, errors.New("save repository is not configured")
	}
	return a.saveRepository.Load(ctx)
}

func (a *App) HasSave() (bool, error) {
	if a.saveRepository == nil {
		return false, errors.New("save repository is not configured")
	}
	return a.saveRepository.Exists()
}

func (a *App) Runtime() services.RuntimeContext {
	return a.runtime
}
