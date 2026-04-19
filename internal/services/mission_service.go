package services

import "github.com/oberones/star-smuggler-go/internal/domain"

type MissionService struct{}

func (s MissionService) AvailableMissions(run domain.RunState, snapshot domain.DataSnapshot, factions FactionService) []domain.MissionDefinition {
	missions := make([]domain.MissionDefinition, 0)
	for _, mission := range snapshot.Missions {
		if mission.OriginPortID != run.Player.CurrentPortID {
			continue
		}
		if mission.UnlockConditions.FactionID != "" &&
			!factions.MeetsMinimumStanding(run, snapshot, mission.UnlockConditions.FactionID, mission.UnlockConditions.MinimumStanding) {
			continue
		}
		if _, active := run.ActiveMissions[mission.ID]; active {
			continue
		}
		if containsString(run.CompletedMissionIDs, mission.ID) {
			continue
		}
		missions = append(missions, mission)
	}
	return missions
}

func (s MissionService) AcceptMission(run *domain.RunState, definition domain.MissionDefinition) domain.MissionState {
	state := domain.NewMissionState(definition, run.TotalJumps)
	run.ActiveMissions[definition.ID] = state
	return state
}

func (s MissionService) ResolveTravelArrival(run *domain.RunState, snapshot domain.DataSnapshot, factions FactionService, upgrades UpgradeService) ([]domain.MissionState, []domain.MissionState) {
	completed := make([]domain.MissionState, 0)
	failed := make([]domain.MissionState, 0)

	for missionID, missionState := range run.ActiveMissions {
		definition, ok := snapshot.MissionsByID[missionID]
		if !ok {
			continue
		}

		if run.TotalJumps > missionState.DeadlineJump {
			missionState.Status = domain.MissionStatusExpired
			failed = append(failed, missionState)
			delete(run.ActiveMissions, missionID)
			if definition.FailureConsequences.FactionID != "" && definition.FailureConsequences.StandingDelta != 0 {
				factions.ApplyStandingDelta(run, snapshot, definition.FailureConsequences.FactionID, definition.FailureConsequences.StandingDelta, "Mission expired")
			}
			continue
		}

		if run.Player.CurrentPortID != definition.DestinationPortID {
			continue
		}
		if definition.RequiredCommodityID != "" && run.Cargo.QuantityFor(definition.RequiredCommodityID) < definition.RequiredQuantity {
			continue
		}

		if definition.RequiredCommodityID != "" {
			run.Cargo.SetQuantity(definition.RequiredCommodityID, run.Cargo.QuantityFor(definition.RequiredCommodityID)-definition.RequiredQuantity)
		}
		run.Player.Credits += upgrades.AdjustMissionReward(*run, definition.Reward.Credits, snapshot)
		missionState.Status = domain.MissionStatusCompleted
		missionState.RewardClaimed = true
		completed = append(completed, missionState)
		run.CompletedMissionIDs = append(run.CompletedMissionIDs, missionID)
		delete(run.ActiveMissions, missionID)
	}

	return completed, failed
}
