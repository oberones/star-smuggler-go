package integration_test

import (
	"testing"

	"github.com/oberones/star-smuggler-go/internal/services"
)

func TestTravelCostParity(t *testing.T) {
	t.Parallel()

	snapshot := loadSnapshot(t)
	travel := services.TravelService{}

	mercury := snapshot.PortsByID["mercury"]
	mars := snapshot.PortsByID["mars"]
	ceres := snapshot.PortsByID["ceres"]
	pluto := snapshot.PortsByID["pluto"]

	if cost := travel.GetTravelCost(mercury, mercury); cost != 0 {
		t.Fatalf("expected same-port travel to cost 0, got %d", cost)
	}
	if cost := travel.GetTravelCost(mercury, mars); cost != 15 {
		t.Fatalf("expected inner-to-inner travel to cost 15, got %d", cost)
	}
	if cost := travel.GetTravelCost(mercury, ceres); cost != 17 {
		t.Fatalf("expected inner-to-outer travel to cost 17, got %d", cost)
	}
	if cost := travel.GetTravelCost(mercury, pluto); cost != 38 {
		t.Fatalf("expected inner-to-fringe travel to cost 38, got %d", cost)
	}
}
