package domain

type PortZone string

const (
	PortZoneInner  PortZone = "Inner"
	PortZoneOuter  PortZone = "Outer"
	PortZoneFringe PortZone = "Fringe"
)

func (z PortZone) IsValid() bool {
	switch z {
	case PortZoneInner, PortZoneOuter, PortZoneFringe:
		return true
	default:
		return false
	}
}

type ItemRarity string

const (
	ItemRarityCommon  ItemRarity = "Common"
	ItemRarityMidTier ItemRarity = "MidTier"
	ItemRarityExotic  ItemRarity = "Exotic"
)

func (r ItemRarity) IsValid() bool {
	switch r {
	case ItemRarityCommon, ItemRarityMidTier, ItemRarityExotic:
		return true
	default:
		return false
	}
}

type EventEffectType string

const (
	EventEffectCreditsLossScaled    EventEffectType = "CreditsLossScaled"
	EventEffectPortPriceMultiplier  EventEffectType = "PortPriceMultiplier"
	EventEffectSingleItemMultiplier EventEffectType = "SingleItemPriceMultiplier"
	EventEffectLoseRandomCargo      EventEffectType = "LoseRandomCargo"
)

func (e EventEffectType) IsValid() bool {
	switch e {
	case EventEffectCreditsLossScaled, EventEffectPortPriceMultiplier, EventEffectSingleItemMultiplier, EventEffectLoseRandomCargo:
		return true
	default:
		return false
	}
}

type PortDefinition struct {
	ID                    string   `json:"id"`
	Name                  string   `json:"name"`
	Description           string   `json:"description"`
	Zone                  PortZone `json:"zone"`
	BackgroundTexturePath string   `json:"backgroundTexturePath"`
	PreviewTexturePath    string   `json:"previewTexturePath"`
	TradeBackgroundPath   string   `json:"tradeBackgroundPath"`
	MusicTrackID          string   `json:"musicTrackId"`
}

type ItemDefinition struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Rarity      ItemRarity `json:"rarity"`
	BasePrice   int        `json:"basePrice"`
}

type EventDefinition struct {
	ID                  string             `json:"id"`
	Name                string             `json:"name"`
	DescriptionTemplate string             `json:"descriptionTemplate"`
	EffectType          EventEffectType    `json:"effectType"`
	Weight              int                `json:"weight"`
	Parameters          map[string]float64 `json:"parameters"`
}

type DataSnapshot struct {
	Ports      []PortDefinition
	Items      []ItemDefinition
	Events     []EventDefinition
	PortsByID  map[string]PortDefinition
	ItemsByID  map[string]ItemDefinition
	EventsByID map[string]EventDefinition
}

func NewDataSnapshot(ports []PortDefinition, items []ItemDefinition, events []EventDefinition) DataSnapshot {
	snapshot := DataSnapshot{
		Ports:      append([]PortDefinition(nil), ports...),
		Items:      append([]ItemDefinition(nil), items...),
		Events:     append([]EventDefinition(nil), events...),
		PortsByID:  make(map[string]PortDefinition, len(ports)),
		ItemsByID:  make(map[string]ItemDefinition, len(items)),
		EventsByID: make(map[string]EventDefinition, len(events)),
	}

	for _, port := range ports {
		snapshot.PortsByID[port.ID] = port
	}

	for _, item := range items {
		snapshot.ItemsByID[item.ID] = item
	}

	for _, event := range events {
		snapshot.EventsByID[event.ID] = event
	}

	return snapshot
}
