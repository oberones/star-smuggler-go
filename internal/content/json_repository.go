package content

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/oberones/star-smuggler-go/internal/domain"
)

type JSONRepository struct {
	reader       FileReader
	portsPath    string
	itemsPath    string
	eventsPath   string
	factionsPath string
	missionsPath string
	storyPath    string
	upgradesPath string
}

func NewJSONRepository(reader FileReader, portsPath string, itemsPath string, eventsPath string, factionsPath string, missionsPath string, storyPath string, upgradesPath string) *JSONRepository {
	if reader == nil {
		reader = OSFileReader{}
	}

	return &JSONRepository{
		reader:       reader,
		portsPath:    portsPath,
		itemsPath:    itemsPath,
		eventsPath:   eventsPath,
		factionsPath: factionsPath,
		missionsPath: missionsPath,
		storyPath:    storyPath,
		upgradesPath: upgradesPath,
	}
}

func NewDefaultJSONRepository(baseDir string) *JSONRepository {
	return NewJSONRepository(
		OSFileReader{},
		ResolvePath(baseDir, DefaultPortsPath),
		ResolvePath(baseDir, DefaultItemsPath),
		ResolvePath(baseDir, DefaultEventsPath),
		ResolvePath(baseDir, DefaultFactionsPath),
		ResolvePath(baseDir, DefaultMissionsPath),
		ResolvePath(baseDir, DefaultStoryPath),
		ResolvePath(baseDir, DefaultUpgradesPath),
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

	factions, err := loadList[domain.FactionDefinition](r.reader, r.factionsPath)
	if err != nil {
		return domain.DataSnapshot{}, fmt.Errorf("load factions: %w", err)
	}

	missions, err := loadList[domain.MissionDefinition](r.reader, r.missionsPath)
	if err != nil {
		return domain.DataSnapshot{}, fmt.Errorf("load missions: %w", err)
	}

	storyArcs, err := loadList[domain.StoryArcDefinition](r.reader, r.storyPath)
	if err != nil {
		return domain.DataSnapshot{}, fmt.Errorf("load story arcs: %w", err)
	}

	upgrades, err := loadList[domain.ShipUpgradeDefinition](r.reader, r.upgradesPath)
	if err != nil {
		return domain.DataSnapshot{}, fmt.Errorf("load upgrades: %w", err)
	}

	snapshot := domain.NewDataSnapshot(ports, items, events, factions, missions, storyArcs, upgrades)
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

	factionIDs := make(map[string]struct{}, len(snapshot.Factions))
	for _, faction := range snapshot.Factions {
		if faction.ID == "" || faction.Name == "" || faction.Description == "" {
			return fmt.Errorf("faction has empty required fields: %+v", faction)
		}
		if !faction.Alignment.IsValid() {
			return fmt.Errorf("faction %q uses invalid alignment %q", faction.ID, faction.Alignment)
		}
		if len(faction.StandingThresholds) == 0 {
			return fmt.Errorf("faction %q must define standing thresholds", faction.ID)
		}
		if _, exists := factionIDs[faction.ID]; exists {
			return fmt.Errorf("duplicate faction id %q", faction.ID)
		}
		lastMinimum := -1 << 30
		for _, threshold := range faction.StandingThresholds {
			if threshold.Tier == "" {
				return fmt.Errorf("faction %q has threshold with empty tier", faction.ID)
			}
			if threshold.MinimumScore < lastMinimum {
				return fmt.Errorf("faction %q thresholds must be sorted by minimum score", faction.ID)
			}
			lastMinimum = threshold.MinimumScore
		}
		for _, homePortID := range faction.HomePortIDs {
			if _, exists := snapshot.PortsByID[homePortID]; !exists {
				return fmt.Errorf("faction %q references unknown home port %q", faction.ID, homePortID)
			}
		}
		factionIDs[faction.ID] = struct{}{}
	}

	missionIDs := make(map[string]struct{}, len(snapshot.Missions))
	for _, mission := range snapshot.Missions {
		if mission.ID == "" || mission.Name == "" || mission.Briefing == "" {
			return fmt.Errorf("mission has empty required fields: %+v", mission)
		}
		if !mission.MissionType.IsValid() {
			return fmt.Errorf("mission %q uses invalid mission type %q", mission.ID, mission.MissionType)
		}
		if _, exists := snapshot.PortsByID[mission.OriginPortID]; !exists {
			return fmt.Errorf("mission %q references unknown origin port %q", mission.ID, mission.OriginPortID)
		}
		if _, exists := snapshot.PortsByID[mission.DestinationPortID]; !exists {
			return fmt.Errorf("mission %q references unknown destination port %q", mission.ID, mission.DestinationPortID)
		}
		if mission.RequiredCommodityID != "" {
			if _, exists := snapshot.ItemsByID[mission.RequiredCommodityID]; !exists {
				return fmt.Errorf("mission %q references unknown commodity %q", mission.ID, mission.RequiredCommodityID)
			}
			if mission.RequiredQuantity <= 0 {
				return fmt.Errorf("mission %q must require a positive quantity", mission.ID)
			}
		}
		if mission.Reward.Credits <= 0 {
			return fmt.Errorf("mission %q must grant positive reward credits", mission.ID)
		}
		if mission.UnlockConditions.FactionID != "" {
			if _, exists := factionIDs[mission.UnlockConditions.FactionID]; !exists {
				return fmt.Errorf("mission %q references unknown unlock faction %q", mission.ID, mission.UnlockConditions.FactionID)
			}
		}
		if mission.FailureConsequences.FactionID != "" {
			if _, exists := factionIDs[mission.FailureConsequences.FactionID]; !exists {
				return fmt.Errorf("mission %q references unknown failure faction %q", mission.ID, mission.FailureConsequences.FactionID)
			}
		}
		if _, exists := missionIDs[mission.ID]; exists {
			return fmt.Errorf("duplicate mission id %q", mission.ID)
		}
		missionIDs[mission.ID] = struct{}{}
	}

	storyArcIDs := make(map[string]struct{}, len(snapshot.StoryArcs))
	for _, storyArc := range snapshot.StoryArcs {
		if storyArc.ID == "" || storyArc.Name == "" {
			return fmt.Errorf("story arc has empty required fields: %+v", storyArc)
		}
		if storyArc.EntryFactionID != "" {
			if _, exists := factionIDs[storyArc.EntryFactionID]; !exists {
				return fmt.Errorf("story arc %q references unknown entry faction %q", storyArc.ID, storyArc.EntryFactionID)
			}
		}
		if len(storyArc.Beats) == 0 {
			return fmt.Errorf("story arc %q must define at least one beat", storyArc.ID)
		}
		for _, beat := range storyArc.Beats {
			if beat.ID == "" || beat.Text == "" {
				return fmt.Errorf("story arc %q has beat with empty required fields", storyArc.ID)
			}
		}
		if _, exists := storyArcIDs[storyArc.ID]; exists {
			return fmt.Errorf("duplicate story arc id %q", storyArc.ID)
		}
		storyArcIDs[storyArc.ID] = struct{}{}
	}

	upgradeIDs := make(map[string]struct{}, len(snapshot.Upgrades))
	for _, upgrade := range snapshot.Upgrades {
		if upgrade.ID == "" || upgrade.Name == "" || upgrade.Description == "" {
			return fmt.Errorf("upgrade has empty required fields: %+v", upgrade)
		}
		if !upgrade.Category.IsValid() {
			return fmt.Errorf("upgrade %q uses invalid category %q", upgrade.ID, upgrade.Category)
		}
		if upgrade.CostCredits <= 0 {
			return fmt.Errorf("upgrade %q must have positive credit cost", upgrade.ID)
		}
		if upgrade.Specialization != "" && !upgrade.Specialization.IsValid() {
			return fmt.Errorf("upgrade %q uses invalid specialization %q", upgrade.ID, upgrade.Specialization)
		}
		if upgrade.RequiredFactionID != "" {
			if _, exists := factionIDs[upgrade.RequiredFactionID]; !exists {
				return fmt.Errorf("upgrade %q references unknown required faction %q", upgrade.ID, upgrade.RequiredFactionID)
			}
			if upgrade.MinimumStanding == "" {
				return fmt.Errorf("upgrade %q must define minimum standing when a faction requirement exists", upgrade.ID)
			}
		}
		if len(upgrade.Effects) == 0 {
			return fmt.Errorf("upgrade %q must define at least one effect", upgrade.ID)
		}
		for _, effect := range upgrade.Effects {
			if !effect.Type.IsValid() {
				return fmt.Errorf("upgrade %q uses invalid effect type %q", upgrade.ID, effect.Type)
			}
			if effect.Value <= 0 {
				return fmt.Errorf("upgrade %q must define positive effect values", upgrade.ID)
			}
		}
		if _, exists := upgradeIDs[upgrade.ID]; exists {
			return fmt.Errorf("duplicate upgrade id %q", upgrade.ID)
		}
		upgradeIDs[upgrade.ID] = struct{}{}
	}

	return nil
}
