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

type FactionAlignment string

const (
	FactionAlignmentLawful   FactionAlignment = "Lawful"
	FactionAlignmentCriminal FactionAlignment = "Criminal"
	FactionAlignmentNeutral  FactionAlignment = "Neutral"
)

func (a FactionAlignment) IsValid() bool {
	switch a {
	case FactionAlignmentLawful, FactionAlignmentCriminal, FactionAlignmentNeutral:
		return true
	default:
		return false
	}
}

type StandingThresholdDefinition struct {
	Tier         string `json:"tier"`
	MinimumScore int    `json:"minimumScore"`
}

type FactionDefinition struct {
	ID                 string                        `json:"id"`
	Name               string                        `json:"name"`
	Description        string                        `json:"description"`
	Alignment          FactionAlignment              `json:"alignment"`
	HomePortIDs        []string                      `json:"homePortIds"`
	RivalFactionIDs    []string                      `json:"rivalFactionIds"`
	StandingThresholds []StandingThresholdDefinition `json:"standingThresholds"`
}

type MissionType string

const (
	MissionTypeDelivery MissionType = "Delivery"
)

func (m MissionType) IsValid() bool {
	switch m {
	case MissionTypeDelivery:
		return true
	default:
		return false
	}
}

type MissionRewardDefinition struct {
	Credits int `json:"credits"`
}

type MissionFailureConsequences struct {
	FactionID     string `json:"factionId"`
	StandingDelta int    `json:"standingDelta"`
}

type MissionUnlockConditions struct {
	FactionID       string `json:"factionId"`
	MinimumStanding string `json:"minimumStanding"`
}

type MissionDefinition struct {
	ID                  string                     `json:"id"`
	Name                string                     `json:"name"`
	Briefing            string                     `json:"briefing"`
	MissionType         MissionType                `json:"missionType"`
	OriginPortID        string                     `json:"originPortId"`
	DestinationPortID   string                     `json:"destinationPortId"`
	RequiredCommodityID string                     `json:"requiredCommodityId,omitempty"`
	RequiredQuantity    int                        `json:"requiredQuantity,omitempty"`
	DeadlineJumpLimit   int                        `json:"deadlineJumpLimit"`
	Reward              MissionRewardDefinition    `json:"reward"`
	FailureConsequences MissionFailureConsequences `json:"failureConsequences"`
	UnlockConditions    MissionUnlockConditions    `json:"unlockConditions"`
}

type StoryBeatDefinition struct {
	ID            string          `json:"id"`
	Text          string          `json:"text"`
	RequiredFlags map[string]bool `json:"requiredFlags,omitempty"`
	SetFlags      map[string]bool `json:"setFlags,omitempty"`
}

type StoryCompletionEffect struct {
	SetFlags map[string]bool `json:"setFlags,omitempty"`
}

type StoryArcDefinition struct {
	ID                string                `json:"id"`
	Name              string                `json:"name"`
	EntryFactionID    string                `json:"entryFactionId"`
	MinimumStanding   string                `json:"minimumStanding"`
	EntryFlags        map[string]bool       `json:"entryFlags,omitempty"`
	Beats             []StoryBeatDefinition `json:"beats"`
	CompletionEffects StoryCompletionEffect `json:"completionEffects"`
}

type DataSnapshot struct {
	Ports         []PortDefinition
	Items         []ItemDefinition
	Events        []EventDefinition
	Factions      []FactionDefinition
	Missions      []MissionDefinition
	StoryArcs     []StoryArcDefinition
	PortsByID     map[string]PortDefinition
	ItemsByID     map[string]ItemDefinition
	EventsByID    map[string]EventDefinition
	FactionsByID  map[string]FactionDefinition
	MissionsByID  map[string]MissionDefinition
	StoryArcsByID map[string]StoryArcDefinition
}

func NewDataSnapshot(
	ports []PortDefinition,
	items []ItemDefinition,
	events []EventDefinition,
	factions []FactionDefinition,
	missions []MissionDefinition,
	storyArcs []StoryArcDefinition,
) DataSnapshot {
	snapshot := DataSnapshot{
		Ports:         append([]PortDefinition(nil), ports...),
		Items:         append([]ItemDefinition(nil), items...),
		Events:        append([]EventDefinition(nil), events...),
		Factions:      append([]FactionDefinition(nil), factions...),
		Missions:      append([]MissionDefinition(nil), missions...),
		StoryArcs:     append([]StoryArcDefinition(nil), storyArcs...),
		PortsByID:     make(map[string]PortDefinition, len(ports)),
		ItemsByID:     make(map[string]ItemDefinition, len(items)),
		EventsByID:    make(map[string]EventDefinition, len(events)),
		FactionsByID:  make(map[string]FactionDefinition, len(factions)),
		MissionsByID:  make(map[string]MissionDefinition, len(missions)),
		StoryArcsByID: make(map[string]StoryArcDefinition, len(storyArcs)),
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

	for _, faction := range factions {
		snapshot.FactionsByID[faction.ID] = faction
	}

	for _, mission := range missions {
		snapshot.MissionsByID[mission.ID] = mission
	}

	for _, storyArc := range storyArcs {
		snapshot.StoryArcsByID[storyArc.ID] = storyArc
	}

	return snapshot
}
