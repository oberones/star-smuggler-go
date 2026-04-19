package godot

import (
	"errors"

	"github.com/oberones/star-smuggler-go/internal/application"
)

type RouteBinder interface {
	ShowRoute(route application.Route) error
}

// AppBridge is the thin shell the future go-dot integration layer will call into.
// It keeps route changes and scene presentation separated from gameplay logic.
type AppBridge struct {
	app    *application.App
	binder RouteBinder
}

func NewAppBridge(app *application.App, binder RouteBinder) *AppBridge {
	return &AppBridge{
		app:    app,
		binder: binder,
	}
}

func (b *AppBridge) SyncRoute() error {
	if b.app == nil {
		return errors.New("app bridge is missing application state")
	}
	if b.binder == nil {
		return errors.New("app bridge is missing route binder")
	}

	return b.binder.ShowRoute(b.app.CurrentRoute())
}

func (b *AppBridge) Navigate(route application.Route) error {
	if b.app == nil {
		return errors.New("app bridge is missing application state")
	}

	b.app.Navigate(route)
	return b.SyncRoute()
}
