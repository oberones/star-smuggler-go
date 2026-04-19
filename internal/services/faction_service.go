package services

import "github.com/oberones/star-smuggler-go/internal/domain"

type FactionService struct{}

func (s FactionService) EnsureDefaults(run *domain.RunState, snapshot domain.DataSnapshot) {
	if len(run.FactionStandings) == 0 {
		run.FactionStandings = domain.DefaultFactionStandings(snapshot)
	}
}

func (s FactionService) ApplyStandingDelta(run *domain.RunState, snapshot domain.DataSnapshot, factionID string, delta int, reason string) (domain.FactionStanding, bool) {
	definition, ok := snapshot.FactionsByID[factionID]
	if !ok {
		return domain.FactionStanding{}, false
	}

	s.EnsureDefaults(run, snapshot)
	standing, ok := run.FactionStandings[factionID]
	if !ok {
		standing = domain.FactionStanding{
			FactionID:        factionID,
			Score:            0,
			StandingTier:     s.ResolveTier(definition, 0),
			LastChangeReason: "Starting standing",
		}
	}

	standing.Score += delta
	standing.StandingTier = s.ResolveTier(definition, standing.Score)
	standing.LastChangeReason = reason
	run.FactionStandings[factionID] = standing
	return standing, true
}

func (s FactionService) ResolveTier(definition domain.FactionDefinition, score int) string {
	tier := "Neutral"
	for _, threshold := range definition.StandingThresholds {
		if score >= threshold.MinimumScore {
			tier = threshold.Tier
		}
	}
	return tier
}

func (s FactionService) MeetsMinimumStanding(run domain.RunState, snapshot domain.DataSnapshot, factionID string, tier string) bool {
	definition, ok := snapshot.FactionsByID[factionID]
	if !ok {
		return false
	}

	requiredScore := -1 << 30
	for _, threshold := range definition.StandingThresholds {
		if threshold.Tier == tier {
			requiredScore = threshold.MinimumScore
			break
		}
	}
	if requiredScore == -1<<30 {
		return false
	}

	standing, ok := run.FactionStandings[factionID]
	if !ok {
		return 0 >= requiredScore
	}

	return standing.Score >= requiredScore
}
