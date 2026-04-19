package services

import "github.com/oberones/star-smuggler-go/internal/domain"

type StoryService struct{}

func (s StoryService) ActivateEligibleArcs(run *domain.RunState, snapshot domain.DataSnapshot, factions FactionService) []string {
	activated := make([]string, 0)
	for _, storyArc := range snapshot.StoryArcs {
		if containsString(run.Story.ActiveStoryArcIDs, storyArc.ID) || containsString(run.Story.CompletedStoryArcIDs, storyArc.ID) {
			continue
		}
		if storyArc.EntryFactionID != "" && !factions.MeetsMinimumStanding(*run, snapshot, storyArc.EntryFactionID, storyArc.MinimumStanding) {
			continue
		}
		if !flagsMatch(run.Story.StoryFlags, storyArc.EntryFlags) {
			continue
		}

		run.Story.ActiveStoryArcIDs = append(run.Story.ActiveStoryArcIDs, storyArc.ID)
		activated = append(activated, storyArc.ID)
	}

	return activated
}

func (s StoryService) AdvanceArc(run *domain.RunState, snapshot domain.DataSnapshot, storyArcID string) (domain.StoryBeatDefinition, bool) {
	storyArc, ok := snapshot.StoryArcsByID[storyArcID]
	if !ok {
		return domain.StoryBeatDefinition{}, false
	}

	for _, beat := range storyArc.Beats {
		if !flagsMatch(run.Story.StoryFlags, beat.RequiredFlags) {
			continue
		}
		if len(beat.SetFlags) > 0 && flagsMatch(run.Story.StoryFlags, beat.SetFlags) {
			continue
		}

		applyFlags(run.Story.StoryFlags, beat.SetFlags)

		allBeatsComplete := true
		for _, checkBeat := range storyArc.Beats {
			if !flagsMatch(run.Story.StoryFlags, checkBeat.SetFlags) {
				allBeatsComplete = false
				break
			}
		}
		if allBeatsComplete {
			run.Story.ActiveStoryArcIDs = removeString(run.Story.ActiveStoryArcIDs, storyArcID)
			run.Story.CompletedStoryArcIDs = append(run.Story.CompletedStoryArcIDs, storyArcID)
			applyFlags(run.Story.StoryFlags, storyArc.CompletionEffects.SetFlags)
		}

		return beat, true
	}

	return domain.StoryBeatDefinition{}, false
}

func flagsMatch(current map[string]bool, required map[string]bool) bool {
	for key, expected := range required {
		if current[key] != expected {
			return false
		}
	}
	return true
}

func applyFlags(current map[string]bool, updates map[string]bool) {
	for key, value := range updates {
		current[key] = value
	}
}

func removeString(values []string, target string) []string {
	filtered := make([]string, 0, len(values))
	for _, value := range values {
		if value != target {
			filtered = append(filtered, value)
		}
	}
	return filtered
}
