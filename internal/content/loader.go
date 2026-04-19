package content

import (
	"context"
	"os"
	"path/filepath"

	"github.com/oberones/star-smuggler-go/internal/domain"
)

const (
	DefaultPortsPath  = "data/ports/ports.json"
	DefaultItemsPath  = "data/items/items.json"
	DefaultEventsPath = "data/events/events.json"
)

type Loader interface {
	LoadSnapshot(ctx context.Context) (domain.DataSnapshot, error)
}

type FileReader interface {
	ReadFile(name string) ([]byte, error)
}

type OSFileReader struct{}

func (OSFileReader) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

func ResolvePath(baseDir string, relativePath string) string {
	if baseDir == "" {
		return relativePath
	}
	return filepath.Join(baseDir, relativePath)
}
