package domain

func DefaultFactionStandings(snapshot DataSnapshot) map[string]FactionStanding {
	result := make(map[string]FactionStanding, len(snapshot.Factions))

	for _, faction := range snapshot.Factions {
		tier := "Neutral"
		for _, threshold := range faction.StandingThresholds {
			if 0 >= threshold.MinimumScore {
				tier = threshold.Tier
			}
		}

		result[faction.ID] = FactionStanding{
			FactionID:        faction.ID,
			Score:            0,
			StandingTier:     tier,
			LastChangeReason: "Starting standing",
		}
	}

	return result
}
