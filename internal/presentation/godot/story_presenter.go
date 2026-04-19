package godot

import (
	"fmt"
	"sort"
	"strings"

	"github.com/oberones/star-smuggler-go/internal/domain"
	"github.com/oberones/star-smuggler-go/internal/services"
)

type StoryNoticeViewModel struct {
	FactionNotices []string
	MissionNotices []string
	StoryNotices   []string
}

type StoryPresenter struct {
	Data     domain.DataSnapshot
	Factions services.FactionService
	Missions services.MissionService
}

func (p StoryPresenter) Present(run domain.RunState) StoryNoticeViewModel {
	return StoryNoticeViewModel{
		FactionNotices: p.presentFactionNotices(run),
		MissionNotices: p.presentMissionNotices(run),
		StoryNotices:   p.presentStoryNotices(run),
	}
}

func (p StoryPresenter) Summary(run domain.RunState) string {
	viewModel := p.Present(run)
	lines := append([]string{}, viewModel.FactionNotices...)
	lines = append(lines, viewModel.MissionNotices...)
	lines = append(lines, viewModel.StoryNotices...)
	return strings.Join(lines, "\n")
}

func (p StoryPresenter) presentFactionNotices(run domain.RunState) []string {
	if len(run.FactionStandings) == 0 {
		return nil
	}

	factionIDs := make([]string, 0, len(run.FactionStandings))
	for factionID := range run.FactionStandings {
		factionIDs = append(factionIDs, factionID)
	}
	sort.Strings(factionIDs)

	notices := make([]string, 0, len(factionIDs))
	for _, factionID := range factionIDs {
		standing := run.FactionStandings[factionID]
		factionName := factionID
		if definition, ok := p.Data.FactionsByID[factionID]; ok {
			factionName = definition.Name
		}

		notices = append(notices, fmt.Sprintf("Faction: %s (%s, %d)", factionName, standing.StandingTier, standing.Score))
	}

	return notices
}

func (p StoryPresenter) presentMissionNotices(run domain.RunState) []string {
	notices := make([]string, 0)

	available := p.Missions.AvailableMissions(run, p.Data, p.Factions)
	for _, mission := range available {
		notices = append(notices, fmt.Sprintf("Available mission: %s", mission.Name))
	}

	missionIDs := make([]string, 0, len(run.ActiveMissions))
	for missionID := range run.ActiveMissions {
		missionIDs = append(missionIDs, missionID)
	}
	sort.Strings(missionIDs)

	for _, missionID := range missionIDs {
		missionState := run.ActiveMissions[missionID]
		missionName := missionID
		if definition, ok := p.Data.MissionsByID[missionID]; ok {
			missionName = definition.Name
		}

		notices = append(notices, fmt.Sprintf("Active mission: %s (%s)", missionName, missionState.Status))
	}

	return notices
}

func (p StoryPresenter) presentStoryNotices(run domain.RunState) []string {
	notices := make([]string, 0, len(run.Story.ActiveStoryArcIDs))

	for _, storyArcID := range run.Story.ActiveStoryArcIDs {
		storyArcName := storyArcID
		if storyArc, ok := p.Data.StoryArcsByID[storyArcID]; ok {
			storyArcName = storyArc.Name
		}
		notices = append(notices, fmt.Sprintf("Story arc: %s", storyArcName))
	}

	if len(run.Story.CompletedStoryArcIDs) > 0 {
		notices = append(notices, fmt.Sprintf("Completed arcs: %d", len(run.Story.CompletedStoryArcIDs)))
	}

	return notices
}
