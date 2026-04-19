package persistence

const CurrentSaveVersion = 1

type SaveData struct {
	Version              int                      `json:"version"`
	Player               PlayerSaveData           `json:"player"`
	CargoByItemID        map[string]int           `json:"cargoByItemId"`
	Markets              []MarketSnapshotSaveData `json:"markets"`
	JumpsSinceLastUpdate int                      `json:"jumpsSinceLastUpdate"`
	TotalJumps           int                      `json:"totalJumps"`
	RecentEvent          *EventResultSaveData     `json:"recentEvent,omitempty"`
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
