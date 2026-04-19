package integration_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/oberones/star-smuggler-go/internal/content"
	"github.com/oberones/star-smuggler-go/internal/domain"
	"github.com/oberones/star-smuggler-go/internal/services"
)

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
