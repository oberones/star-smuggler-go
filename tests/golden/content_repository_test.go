package golden_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/oberones/star-smuggler-go/internal/content"
	"github.com/oberones/star-smuggler-go/internal/domain"
)

type mapReader map[string][]byte

func (m mapReader) ReadFile(name string) ([]byte, error) {
	return m[name], nil
}

func TestJSONRepositoryMatchesGoldenSnapshot(t *testing.T) {
	t.Parallel()

	repoRoot := filepath.Clean(filepath.Join("..", ".."))
	repository := content.NewDefaultJSONRepository(repoRoot)

	snapshot, err := repository.LoadSnapshot(context.Background())
	if err != nil {
		t.Fatalf("load snapshot: %v", err)
	}

	actual, err := json.MarshalIndent(struct {
		Ports  []domain.PortDefinition  `json:"ports"`
		Items  []domain.ItemDefinition  `json:"items"`
		Events []domain.EventDefinition `json:"events"`
	}{
		Ports:  snapshot.Ports,
		Items:  snapshot.Items,
		Events: snapshot.Events,
	}, "", "  ")
	if err != nil {
		t.Fatalf("marshal snapshot: %v", err)
	}

	expectedPath := filepath.Join("testdata", "content_snapshot.golden.json")
	expected, err := os.ReadFile(expectedPath)
	if err != nil {
		t.Fatalf("read golden file: %v", err)
	}

	var expectedValue any
	if err := json.Unmarshal(expected, &expectedValue); err != nil {
		t.Fatalf("decode golden snapshot: %v", err)
	}

	var actualValue any
	if err := json.Unmarshal(actual, &actualValue); err != nil {
		t.Fatalf("decode actual snapshot: %v", err)
	}

	if !reflect.DeepEqual(expectedValue, actualValue) {
		t.Fatalf("golden snapshot mismatch\nexpected:\n%s\n\nactual:\n%s", expected, actual)
	}
}

func TestJSONRepositoryRejectsInvalidContent(t *testing.T) {
	t.Parallel()

	reader := mapReader{
		"ports.json":  []byte(`[{"id":"mars","name":"Mars","description":"Port","zone":"Inner","backgroundTexturePath":"res://bg","previewTexturePath":"res://preview","tradeBackgroundPath":"res://trade","musicTrackId":"music"}]`),
		"items.json":  []byte(`[{"id":"bad_item","name":"Bad Item","description":"Oops","rarity":"Legendary","basePrice":10}]`),
		"events.json": []byte(`[{"id":"event","name":"Event","descriptionTemplate":"Hello","effectType":"CreditsLossScaled","weight":1,"parameters":{}}]`),
	}

	repository := content.NewJSONRepository(reader, "ports.json", "items.json", "events.json")
	_, err := repository.LoadSnapshot(context.Background())
	if err == nil {
		t.Fatal("expected invalid content error")
	}
}
