package domain

type EventResult struct {
	EventID             string
	Name                string
	ResolvedDescription string
	RolledValues        map[string]float64
}

type TravelResolution struct {
	Message      string
	AppliedEvent *EventResult
	Route        *RouteState
}
