package persistence

const CurrentSaveVersion = 1

type SaveData struct {
	Version                   int                       `json:"version"`
	Player                    PlayerSaveData            `json:"player"`
	CargoByItemID             map[string]int            `json:"cargoByItemId"`
	Markets                   []MarketSnapshotSaveData  `json:"markets"`
	RoutePressureByKey        map[string]int            `json:"routePressureByKey,omitempty"`
	CommodityPressureByItemID map[string]int            `json:"commodityPressureByItemId,omitempty"`
	FactionStandings          []FactionStandingSaveData `json:"factionStandings,omitempty"`
	ActiveMissions            []MissionStateSaveData    `json:"activeMissions,omitempty"`
	CompletedMissionIDs       []string                  `json:"completedMissionIds,omitempty"`
	Story                     StoryStateSaveData        `json:"story"`
	EmergencyRecoveryUsed     bool                      `json:"emergencyRecoveryUsed"`
	JumpsSinceLastUpdate      int                       `json:"jumpsSinceLastUpdate"`
	TotalJumps                int                       `json:"totalJumps"`
	RecentEvent               *EventResultSaveData      `json:"recentEvent,omitempty"`
}

type PlayerSaveData struct {
	Credits       int    `json:"credits"`
	CargoLimit    int    `json:"cargoLimit"`
	CurrentPortID string `json:"currentPortId"`
}

type MarketSnapshotSaveData struct {
	PortID           string         `json:"portId"`
	AvailableItemIDs []string       `json:"availableItemIds"`
	PricesByItemID   map[string]int `json:"pricesByItemId"`
}

type EventResultSaveData struct {
	EventID             string             `json:"eventId"`
	Name                string             `json:"name"`
	ResolvedDescription string             `json:"resolvedDescription"`
	RolledValues        map[string]float64 `json:"rolledValues"`
}

type FactionStandingSaveData struct {
	FactionID        string `json:"factionId"`
	Score            int    `json:"score"`
	StandingTier     string `json:"standingTier"`
	LastChangeReason string `json:"lastChangeReason"`
}

type MissionStateSaveData struct {
	MissionDefinitionID string          `json:"missionDefinitionId"`
	Status              string          `json:"status"`
	AcceptedAtJump      int             `json:"acceptedAtJump"`
	DeadlineJump        int             `json:"deadlineJump"`
	ProgressFlags       map[string]bool `json:"progressFlags"`
	RewardClaimed       bool            `json:"rewardClaimed"`
}

type StoryStateSaveData struct {
	ActiveStoryArcIDs    []string          `json:"activeStoryArcIds"`
	CompletedStoryArcIDs []string          `json:"completedStoryArcIds"`
	StoryFlags           map[string]bool   `json:"storyFlags"`
	NamedCharacterStates map[string]string `json:"namedCharacterStates"`
}
