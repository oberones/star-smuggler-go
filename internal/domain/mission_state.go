package domain

func NewMissionState(definition MissionDefinition, acceptedAtJump int) MissionState {
	return MissionState{
		MissionDefinitionID: definition.ID,
		Status:              MissionStatusInProgress,
		AcceptedAtJump:      acceptedAtJump,
		DeadlineJump:        acceptedAtJump + definition.DeadlineJumpLimit,
		ProgressFlags:       map[string]bool{},
		RewardClaimed:       false,
	}
}
