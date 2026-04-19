package content

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/oberones/star-smuggler-go/internal/domain"
)

type JSONRepository struct {
	reader     FileReader
	portsPath  string
	itemsPath  string
	eventsPath string
}

func NewJSONRepository(reader FileReader, portsPath string, itemsPath string, eventsPath string) *JSONRepository {
	if reader == nil {
		reader = OSFileReader{}
	}

	return &JSONRepository{
		reader:     reader,
		portsPath:  portsPath,
		itemsPath:  itemsPath,
		eventsPath: eventsPath,
	}
}

func NewDefaultJSONRepository(baseDir string) *JSONRepository {
	return NewJSONRepository(
		OSFileReader{},
		ResolvePath(baseDir, DefaultPortsPath),
		ResolvePath(baseDir, DefaultItemsPath),
		ResolvePath(baseDir, DefaultEventsPath),
	)
}

func (r *JSONRepository) LoadSnapshot(ctx context.Context) (domain.DataSnapshot, error) {
	if err := ctx.Err(); err != nil {
		return domain.DataSnapshot{}, err
	}

	ports, err := loadList[domain.PortDefinition](r.reader, r.portsPath)
	if err != nil {
		return domain.DataSnapshot{}, fmt.Errorf("load ports: %w", err)
	}

	items, err := loadList[domain.ItemDefinition](r.reader, r.itemsPath)
	if err != nil {
		return domain.DataSnapshot{}, fmt.Errorf("load items: %w", err)
	}

	events, err := loadList[domain.EventDefinition](r.reader, r.eventsPath)
	if err != nil {
		return domain.DataSnapshot{}, fmt.Errorf("load events: %w", err)
	}

	snapshot := domain.NewDataSnapshot(ports, items, events)
	if err := validateSnapshot(snapshot); err != nil {
		return domain.DataSnapshot{}, err
	}

	return snapshot, nil
}

func loadList[T any](reader FileReader, path string) ([]T, error) {
	bytes, err := reader.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var values []T
	if err := json.Unmarshal(bytes, &values); err != nil {
		return nil, err
	}

	return values, nil
}

func validateSnapshot(snapshot domain.DataSnapshot) error {
	portIDs := make(map[string]struct{}, len(snapshot.Ports))
	for _, port := range snapshot.Ports {
		if port.ID == "" || port.Name == "" || port.Description == "" {
			return fmt.Errorf("port has empty required fields: %+v", port)
		}
		if !port.Zone.IsValid() {
			return fmt.Errorf("port %q uses invalid zone %q", port.ID, port.Zone)
		}
		if port.BackgroundTexturePath == "" || port.PreviewTexturePath == "" || port.TradeBackgroundPath == "" || port.MusicTrackID == "" {
			return fmt.Errorf("port %q is missing asset or music metadata", port.ID)
		}
		if _, exists := portIDs[port.ID]; exists {
			return fmt.Errorf("duplicate port id %q", port.ID)
		}
		portIDs[port.ID] = struct{}{}
	}

	itemIDs := make(map[string]struct{}, len(snapshot.Items))
	for _, item := range snapshot.Items {
		if item.ID == "" || item.Name == "" || item.Description == "" {
			return fmt.Errorf("item has empty required fields: %+v", item)
		}
		if !item.Rarity.IsValid() {
			return fmt.Errorf("item %q uses invalid rarity %q", item.ID, item.Rarity)
		}
		if item.BasePrice <= 0 {
			return fmt.Errorf("item %q must have positive base price", item.ID)
		}
		if _, exists := itemIDs[item.ID]; exists {
			return fmt.Errorf("duplicate item id %q", item.ID)
		}
		itemIDs[item.ID] = struct{}{}
	}

	eventIDs := make(map[string]struct{}, len(snapshot.Events))
	for _, event := range snapshot.Events {
		if event.ID == "" || event.Name == "" || event.DescriptionTemplate == "" {
			return fmt.Errorf("event has empty required fields: %+v", event)
		}
		if !event.EffectType.IsValid() {
			return fmt.Errorf("event %q uses invalid effect type %q", event.ID, event.EffectType)
		}
		if event.Weight <= 0 {
			return fmt.Errorf("event %q must have positive weight", event.ID)
		}
		if event.Parameters == nil {
			event.Parameters = map[string]float64{}
		}
		if _, exists := eventIDs[event.ID]; exists {
			return fmt.Errorf("duplicate event id %q", event.ID)
		}
		eventIDs[event.ID] = struct{}{}
	}

	if len(snapshot.Ports) == 0 || len(snapshot.Items) == 0 || len(snapshot.Events) == 0 {
		return errors.New("snapshot must include ports, items, and events")
	}

	return nil
}
