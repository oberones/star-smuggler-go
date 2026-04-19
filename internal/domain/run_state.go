package domain

const (
	StartingCredits    = 500
	StartingCargoLimit = 30
)

type PlayerShipState struct {
	Credits       int
	CargoLimit    int
	CurrentPortID string
}

type CargoState struct {
	ItemQuantities map[string]int
}

func NewCargoState() CargoState {
	return CargoState{
		ItemQuantities: make(map[string]int),
	}
}

func (c CargoState) TotalUnits() int {
	total := 0
	for _, quantity := range c.ItemQuantities {
		total += quantity
	}
	return total
}

func (c CargoState) QuantityFor(itemID string) int {
	return c.ItemQuantities[itemID]
}

func (c *CargoState) SetQuantity(itemID string, quantity int) {
	if c.ItemQuantities == nil {
		c.ItemQuantities = make(map[string]int)
	}

	if quantity <= 0 {
		delete(c.ItemQuantities, itemID)
		return
	}

	c.ItemQuantities[itemID] = quantity
}

type MarketSnapshot struct {
	PortID           string
	AvailableItemIDs []string
	PricesByItemID   map[string]int
}

type RouteStatus string

const (
	RouteStatusPreviewed RouteStatus = "previewed"
	RouteStatusCommitted RouteStatus = "committed"
	RouteStatusAnimating RouteStatus = "animating"
	RouteStatusResolved  RouteStatus = "resolved"
)

type RouteState struct {
	OriginPortID      string
	DestinationPortID string
	TravelCost        int
	Status            RouteStatus
}

type FactionStanding struct {
	FactionID        string
	Score            int
	StandingTier     string
	LastChangeReason string
}

type MissionStatus string

const (
	MissionStatusAvailable  MissionStatus = "available"
	MissionStatusAccepted   MissionStatus = "accepted"
	MissionStatusInProgress MissionStatus = "in_progress"
	MissionStatusCompleted  MissionStatus = "completed"
	MissionStatusFailed     MissionStatus = "failed"
	MissionStatusExpired    MissionStatus = "expired"
)

type MissionState struct {
	MissionDefinitionID string
	Status              MissionStatus
	AcceptedAtJump      int
	DeadlineJump        int
	ProgressFlags       map[string]bool
	RewardClaimed       bool
}

type StoryState struct {
	ActiveStoryArcIDs    []string
	CompletedStoryArcIDs []string
	StoryFlags           map[string]bool
	NamedCharacterStates map[string]string
}

type RunState struct {
	Player                    PlayerShipState
	Cargo                     CargoState
	MarketsByPortID           map[string]MarketSnapshot
	RoutePressureByKey        map[string]int
	CommodityPressureByItemID map[string]int
	Progression               ShipProgressionState
	FactionStandings          map[string]FactionStanding
	ActiveMissions            map[string]MissionState
	CompletedMissionIDs       []string
	Story                     StoryState
	EmergencyRecoveryUsed     bool
	JumpsSinceLastUpdate      int
	TotalJumps                int
	RecentEvent               *EventResult
	PendingRoute              *RouteState
}

func NewRunState() RunState {
	return RunState{
		Player: PlayerShipState{
			Credits:    StartingCredits,
			CargoLimit: StartingCargoLimit,
		},
		Cargo:                     NewCargoState(),
		MarketsByPortID:           make(map[string]MarketSnapshot),
		RoutePressureByKey:        make(map[string]int),
		CommodityPressureByItemID: make(map[string]int),
		Progression:               NewShipProgressionState(),
		FactionStandings:          make(map[string]FactionStanding),
		ActiveMissions:            make(map[string]MissionState),
		CompletedMissionIDs:       []string{},
		Story: StoryState{
			ActiveStoryArcIDs:    []string{},
			CompletedStoryArcIDs: []string{},
			StoryFlags:           make(map[string]bool),
			NamedCharacterStates: make(map[string]string),
		},
	}
}
