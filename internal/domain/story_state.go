package domain

func NewStoryState() StoryState {
	return StoryState{
		ActiveStoryArcIDs:    []string{},
		CompletedStoryArcIDs: []string{},
		StoryFlags:           map[string]bool{},
		NamedCharacterStates: map[string]string{},
	}
}
