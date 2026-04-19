package godot

import "strings"

type MainMenuViewModel struct {
	CanContinue   bool
	StatusMessage string
}

type MainMenuPresenter struct{}

func (p MainMenuPresenter) Present(canContinue bool, statusOverride string) MainMenuViewModel {
	statusMessage := strings.TrimSpace(statusOverride)
	if statusMessage == "" {
		if canContinue {
			statusMessage = "Resume a saved smuggling run or start a fresh route."
		} else {
			statusMessage = "Start a new run to begin rebuilding StarSmuggler in Go."
		}
	}

	return MainMenuViewModel{
		CanContinue:   canContinue,
		StatusMessage: statusMessage,
	}
}
