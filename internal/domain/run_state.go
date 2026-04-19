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

type EventResult struct {
	EventID             string
	Name                string
	ResolvedDescription string
	RolledValues        map[string]float64
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

type RunState struct {
	Player               PlayerShipState
	Cargo                CargoState
	MarketsByPortID      map[string]MarketSnapshot
	JumpsSinceLastUpdate int
	TotalJumps           int
	RecentEvent          *EventResult
	PendingRoute         *RouteState
}

func NewRunState() RunState {
	return RunState{
		Player: PlayerShipState{
			Credits:    StartingCredits,
			CargoLimit: StartingCargoLimit,
		},
		Cargo:           NewCargoState(),
		MarketsByPortID: make(map[string]MarketSnapshot),
	}
}
