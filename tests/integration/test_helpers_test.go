package integration_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/oberones/star-smuggler-go/internal/content"
	"github.com/oberones/star-smuggler-go/internal/domain"
	"github.com/oberones/star-smuggler-go/internal/services"
)

type stubRNG struct {
	intValues   []int
	floatValues []float64
	intIndex    int
	floatIndex  int
}

func (r *stubRNG) Intn(n int) int {
	if n <= 0 {
		return 0
	}
	if r.intIndex >= len(r.intValues) {
		return 0
	}

	value := r.intValues[r.intIndex]
	r.intIndex++
	if value < 0 {
		value = -value
	}
	return value % n
}

func (r *stubRNG) Float64() float64 {
	if r.floatIndex >= len(r.floatValues) {
		return 0
	}

	value := r.floatValues[r.floatIndex]
	r.floatIndex++
	if value < 0 {
		return 0
	}
	if value > 1 {
		return 1
	}
	return value
}

func loadSnapshot(t *testing.T) domain.DataSnapshot {
	t.Helper()

	repoRoot := filepath.Clean(filepath.Join("..", ".."))
	repository := content.NewDefaultJSONRepository(repoRoot)
	snapshot, err := repository.LoadSnapshot(context.Background())
	if err != nil {
		t.Fatalf("load snapshot: %v", err)
	}
	return snapshot
}

func seededRuntime(seed int64) services.RuntimeContext {
	return services.NewSeededRuntimeContext(seed)
}

func seededRunAtPort(snapshot domain.DataSnapshot, portID string) domain.RunState {
	run := domain.NewRunState()
	run.Player.CurrentPortID = portID
	run.FactionStandings = domain.DefaultFactionStandings(snapshot)
	run.Story = domain.NewStoryState()
	run.ActiveMissions = map[string]domain.MissionState{}
	return run
}
