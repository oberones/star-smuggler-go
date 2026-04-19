package application

import (
	"fmt"
	"sort"
	"strings"

	"github.com/oberones/star-smuggler-go/internal/domain"
	"github.com/oberones/star-smuggler-go/internal/services"
)

type StoryUpdate struct {
	Notices []string
}

type StoryCommands struct {
	Data     domain.DataSnapshot
	Factions services.FactionService
	Missions services.MissionService
	Stories  services.StoryService
	Upgrades services.UpgradeService
}

func NewStoryCommands(data domain.DataSnapshot) StoryCommands {
	return StoryCommands{
		Data: data,
	}
}

func (c StoryCommands) SyncAtRunStart(run *domain.RunState) StoryUpdate {
	c.Factions.EnsureDefaults(run, c.Data)
	if run.Story.StoryFlags == nil {
		run.Story = domain.NewStoryState()
	}
	return StoryUpdate{}
}

func (c StoryCommands) AvailableMissions(run domain.RunState) []domain.MissionDefinition {
	return c.Missions.AvailableMissions(run, c.Data, c.Factions)
}

func (c StoryCommands) AcceptMission(run *domain.RunState, missionID string) (StoryUpdate, error) {
	available := c.AvailableMissions(*run)
	for _, mission := range available {
		if mission.ID != missionID {
			continue
		}

		state := c.Missions.AcceptMission(run, mission)
		notices := []string{
			fmt.Sprintf("Accepted mission: %s", mission.Name),
			fmt.Sprintf("Deadline: deliver within %d jumps.", state.DeadlineJump-state.AcceptedAtJump),
		}
		if mission.RequiredCommodityID != "" && run.Cargo.QuantityFor(mission.RequiredCommodityID) >= mission.RequiredQuantity {
			state.ProgressFlags["cargo_loaded"] = true
			run.ActiveMissions[mission.ID] = state
			notices = append(notices, fmt.Sprintf("Cargo ready: %d %s already in the hold.", mission.RequiredQuantity, mission.RequiredCommodityID))
		}

		return StoryUpdate{Notices: notices}, nil
	}

	return StoryUpdate{}, fmt.Errorf("mission %q is not currently available", missionID)
}

func (c StoryCommands) SyncAfterTrade(run *domain.RunState, itemID string, quantity int, isBuy bool) StoryUpdate {
	if run == nil || quantity <= 0 {
		return StoryUpdate{}
	}

	notices := make([]string, 0)
	missionIDs := make([]string, 0, len(run.ActiveMissions))
	for missionID := range run.ActiveMissions {
		missionIDs = append(missionIDs, missionID)
	}
	sort.Strings(missionIDs)

	for _, missionID := range missionIDs {
		missionState := run.ActiveMissions[missionID]
		definition, ok := c.Data.MissionsByID[missionState.MissionDefinitionID]
		if !ok || definition.RequiredCommodityID == "" || definition.RequiredCommodityID != itemID {
			continue
		}

		haveEnoughCargo := run.Cargo.QuantityFor(itemID) >= definition.RequiredQuantity
		if isBuy && haveEnoughCargo && !missionState.ProgressFlags["cargo_loaded"] {
			missionState.ProgressFlags["cargo_loaded"] = true
			run.ActiveMissions[missionID] = missionState
			notices = append(notices, fmt.Sprintf("Mission cargo secured for %s.", definition.Name))
			continue
		}

		if !isBuy && !haveEnoughCargo && missionState.ProgressFlags["cargo_loaded"] {
			missionState.ProgressFlags["cargo_loaded"] = false
			run.ActiveMissions[missionID] = missionState
			notices = append(notices, fmt.Sprintf("Mission cargo shortfall for %s.", definition.Name))
		}
	}

	return StoryUpdate{Notices: notices}
}

func (c StoryCommands) SyncAfterTravel(run *domain.RunState) StoryUpdate {
	if run == nil {
		return StoryUpdate{}
	}

	notices := make([]string, 0)

	completed, failed := c.Missions.ResolveTravelArrival(run, c.Data, c.Factions, c.Upgrades)
	for _, missionState := range completed {
		if definition, ok := c.Data.MissionsByID[missionState.MissionDefinitionID]; ok {
			notices = append(notices, fmt.Sprintf("Mission completed: %s", definition.Name))
		}
	}
	for _, missionState := range failed {
		if definition, ok := c.Data.MissionsByID[missionState.MissionDefinitionID]; ok {
			notices = append(notices, fmt.Sprintf("Mission failed: %s", definition.Name))
		}
	}

	activated := c.Stories.ActivateEligibleArcs(run, c.Data, c.Factions)
	for _, storyArcID := range activated {
		storyArc := c.Data.StoryArcsByID[storyArcID]
		notices = append(notices, fmt.Sprintf("Story unlocked: %s", storyArc.Name))

		if beat, ok := c.Stories.AdvanceArc(run, c.Data, storyArcID); ok {
			notices = append(notices, beat.Text)
		}
	}

	return StoryUpdate{Notices: notices}
}

func (u StoryUpdate) Summary() string {
	return strings.Join(u.Notices, "\n")
}
