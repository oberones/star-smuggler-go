package godot

import (
	"github.com/oberones/star-smuggler-go/internal/application"
	"github.com/oberones/star-smuggler-go/internal/domain"
)

const (
	MainMenuMusicTrackID = "singularity"
	WorldMusicTrackID    = "world_default"
	DefaultClickSfxID    = "click"
)

type AudioBridge struct {
	Cache *ResourceCache
}

func (b AudioBridge) MusicForRoute(route application.Route, run *domain.RunState, snapshot domain.DataSnapshot) string {
	trackID := ""

	switch route {
	case application.RouteMainMenu, application.RouteGameOver:
		trackID = MainMenuMusicTrackID
	case application.RoutePortOverview, application.RouteTrade:
		trackID = b.portMusic(run, snapshot)
	case application.RouteTravel, application.RouteTravelAnimation:
		trackID = WorldMusicTrackID
	}

	return b.resolveMusic(trackID)
}

func (b AudioBridge) ClickSfxForAction(_ string) string {
	return b.resolveSfx(DefaultClickSfxID)
}

func (b AudioBridge) ShouldReplay(currentTrackID string, nextTrackID string) bool {
	return nextTrackID != "" && currentTrackID != nextTrackID
}

func (b AudioBridge) portMusic(run *domain.RunState, snapshot domain.DataSnapshot) string {
	if run == nil {
		return MainMenuMusicTrackID
	}

	port, ok := snapshot.PortsByID[run.Player.CurrentPortID]
	if !ok || port.MusicTrackID == "" {
		return WorldMusicTrackID
	}

	return port.MusicTrackID
}

func (b AudioBridge) resolveMusic(trackID string) string {
	if b.Cache == nil {
		return trackID
	}
	return b.Cache.ResolveMusic(trackID)
}

func (b AudioBridge) resolveSfx(sfxID string) string {
	if b.Cache == nil {
		return sfxID
	}
	return b.Cache.ResolveSfx(sfxID)
}
