package services

import (
	"sort"

	"github.com/oberones/star-smuggler-go/internal/domain"
)

type TravelService struct{}

type TravelQuote struct {
	Destination domain.PortDefinition
	Cost        int
}

func (s TravelService) GetTravelCost(fromPort domain.PortDefinition, toPort domain.PortDefinition) int {
	if fromPort.ID == toPort.ID {
		return 0
	}

	baseCost := 15
	zoneDifference := abs(zoneRank(fromPort.Zone) - zoneRank(toPort.Zone))
	cost := baseCost + (zoneDifference * 2)
	if zoneDifference >= 2 {
		cost *= 2
	}
	return cost
}

func (s TravelService) GetCheapestTravelCostFromPort(origin domain.PortDefinition, ports []domain.PortDefinition) int {
	cheapest := 0
	for _, destination := range ports {
		if destination.ID == origin.ID {
			continue
		}

		cost := s.GetTravelCost(origin, destination)
		if cheapest == 0 || cost < cheapest {
			cheapest = cost
		}
	}
	return cheapest
}

func (s TravelService) GetDestinationsFromPort(origin domain.PortDefinition, ports []domain.PortDefinition) []domain.PortDefinition {
	destinations := make([]domain.PortDefinition, 0, len(ports))
	for _, port := range ports {
		if port.ID != origin.ID {
			destinations = append(destinations, port)
		}
	}

	sort.Slice(destinations, func(i, j int) bool {
		left := destinations[i]
		right := destinations[j]
		if zoneRank(left.Zone) == zoneRank(right.Zone) {
			return left.Name < right.Name
		}
		return zoneRank(left.Zone) < zoneRank(right.Zone)
	})

	return destinations
}

func zoneRank(zone domain.PortZone) int {
	switch zone {
	case domain.PortZoneInner:
		return 0
	case domain.PortZoneOuter:
		return 1
	case domain.PortZoneFringe:
		return 2
	default:
		return 0
	}
}

func abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}
