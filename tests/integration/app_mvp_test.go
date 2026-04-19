package integration_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/oberones/star-smuggler-go/internal/application"
	"github.com/oberones/star-smuggler-go/internal/domain"
)

type stubContentLoader struct {
	snapshot domain.DataSnapshot
	err      error
}

func (l stubContentLoader) LoadSnapshot(_ context.Context) (domain.DataSnapshot, error) {
	return l.snapshot, l.err
}

type memorySaveRepository struct {
	exists bool
	run    domain.RunState
}

func (r *memorySaveRepository) Exists() (bool, error) {
	return r.exists, nil
}

func (r *memorySaveRepository) Load(_ context.Context) (domain.RunState, error) {
	return cloneRunState(r.run), nil
}

func (r *memorySaveRepository) Save(_ context.Context, run domain.RunState) error {
	r.exists = true
	r.run = cloneRunState(run)
	return nil
}

func cloneRunState(run domain.RunState) domain.RunState {
	bytes, err := json.Marshal(run)
	if err != nil {
		panic(err)
	}

	var clone domain.RunState
	if err := json.Unmarshal(bytes, &clone); err != nil {
		panic(err)
	}

	if clone.Cargo.ItemQuantities == nil {
		clone.Cargo.ItemQuantities = map[string]int{}
	}
	if clone.MarketsByPortID == nil {
		clone.MarketsByPortID = map[string]domain.MarketSnapshot{}
	}
	if clone.RoutePressureByKey == nil {
		clone.RoutePressureByKey = map[string]int{}
	}
	if clone.CommodityPressureByItemID == nil {
		clone.CommodityPressureByItemID = map[string]int{}
	}
	if clone.Progression.SpecializationFlags == nil {
		clone.Progression.SpecializationFlags = map[string]bool{}
	}
	if clone.Progression.PurchasedUpgradeIDs == nil {
		clone.Progression.PurchasedUpgradeIDs = []string{}
	}

	return clone
}

func TestAppMVPRouteLoopCanStartTravelAndResolveBackToPort(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	snapshot := loadSnapshot(t)
	app := application.NewApp(stubContentLoader{snapshot: snapshot}, &memorySaveRepository{}, seededRuntime(17))

	if err := app.Bootstrap(ctx); err != nil {
		t.Fatalf("bootstrap app: %v", err)
	}
	if route := app.CurrentRoute(); route != application.RouteMainMenu {
		t.Fatalf("expected main menu route after bootstrap, got %v", route)
	}

	if err := app.StartNewRun(ctx); err != nil {
		t.Fatalf("start new run: %v", err)
	}
	if app.ActiveRun() == nil {
		t.Fatal("expected active run after starting a new game")
	}
	if route := app.CurrentRoute(); route != application.RoutePortOverview {
		t.Fatalf("expected port overview after new run, got %v", route)
	}

	if err := app.OpenTravel(); err != nil {
		t.Fatalf("open travel: %v", err)
	}
	if route := app.CurrentRoute(); route != application.RouteTravel {
		t.Fatalf("expected travel route, got %v", route)
	}

	quotes, err := app.PreviewTravel()
	if err != nil {
		t.Fatalf("preview travel: %v", err)
	}
	if len(quotes) == 0 {
		t.Fatal("expected at least one travel quote")
	}

	if _, err := app.BeginTravel(quotes[0].Destination.ID); err != nil {
		t.Fatalf("begin travel: %v", err)
	}
	if route := app.CurrentRoute(); route != application.RouteTravelAnimation {
		t.Fatalf("expected travel animation route, got %v", route)
	}
	if app.ActiveRun().PendingRoute == nil || app.ActiveRun().PendingRoute.Status != domain.RouteStatusAnimating {
		t.Fatalf("expected animating pending route, got %#v", app.ActiveRun().PendingRoute)
	}

	resolution, err := app.ResolveTravel(ctx)
	if err != nil {
		t.Fatalf("resolve travel: %v", err)
	}
	if resolution.Route == nil || resolution.Route.Status != domain.RouteStatusResolved {
		t.Fatalf("expected resolved route in travel resolution, got %#v", resolution.Route)
	}
	if route := app.CurrentRoute(); route != application.RoutePortOverview && route != application.RouteGameOver {
		t.Fatalf("expected route to land on port overview or game over, got %v", route)
	}
	if app.ActiveRun().Player.CurrentPortID != quotes[0].Destination.ID {
		t.Fatalf("expected current port to be %q, got %q", quotes[0].Destination.ID, app.ActiveRun().Player.CurrentPortID)
	}
}

func TestAppMVPRecoveryAndContinuePersistRunState(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	snapshot := loadSnapshot(t)
	saveRepo := &memorySaveRepository{}
	app := application.NewApp(stubContentLoader{snapshot: snapshot}, saveRepo, seededRuntime(23))

	if err := app.Bootstrap(ctx); err != nil {
		t.Fatalf("bootstrap app: %v", err)
	}
	if err := app.StartNewRun(ctx); err != nil {
		t.Fatalf("start new run: %v", err)
	}

	run := app.ActiveRun()
	run.Player.Credits = 0
	run.Cargo = domain.NewCargoState()

	message, recovered, err := app.Recover(ctx)
	if err != nil {
		t.Fatalf("recover: %v", err)
	}
	if !recovered {
		t.Fatalf("expected recovery to succeed, got message %q", message)
	}
	if !app.ActiveRun().EmergencyRecoveryUsed {
		t.Fatal("expected recovery usage to be recorded on the active run")
	}
	if route := app.CurrentRoute(); route != application.RoutePortOverview {
		t.Fatalf("expected recovery to return to port overview, got %v", route)
	}

	reloaded := application.NewApp(stubContentLoader{snapshot: snapshot}, saveRepo, seededRuntime(29))
	if err := reloaded.Bootstrap(ctx); err != nil {
		t.Fatalf("bootstrap second app: %v", err)
	}
	if err := reloaded.ContinueSavedRun(ctx); err != nil {
		t.Fatalf("continue saved run: %v", err)
	}
	if reloaded.ActiveRun() == nil || !reloaded.ActiveRun().EmergencyRecoveryUsed {
		t.Fatal("expected continued run to preserve recovery usage")
	}
}
