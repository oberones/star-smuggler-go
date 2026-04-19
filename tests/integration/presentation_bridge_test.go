package integration_test

import (
	"testing"

	"github.com/oberones/star-smuggler-go/internal/application"
	"github.com/oberones/star-smuggler-go/internal/presentation/godot"
)

func TestResourceCacheCachesTextureMusicAndSfxLookups(t *testing.T) {
	t.Parallel()

	cache := godot.NewResourceCache()
	if got := cache.ResolveTexture(" res://assets/ui/cockpit.png "); got != "res://assets/ui/cockpit.png" {
		t.Fatalf("expected trimmed texture path, got %q", got)
	}
	if got := cache.ResolveMusic(" world_default "); got != "world_default" {
		t.Fatalf("expected trimmed music id, got %q", got)
	}
	if got := cache.ResolveSfx(" click "); got != "click" {
		t.Fatalf("expected trimmed sfx id, got %q", got)
	}

	cache.ResolveTexture("res://assets/ui/cockpit.png")
	cache.ResolveMusic("world_default")
	cache.ResolveSfx("click")

	textures, music, sfx := cache.Stats()
	if textures != 1 || music != 1 || sfx != 1 {
		t.Fatalf("expected one cached entry per resource class, got textures=%d music=%d sfx=%d", textures, music, sfx)
	}
}

func TestAudioBridgeMatchesMonoGameRouteMusicIntent(t *testing.T) {
	t.Parallel()

	snapshot := loadSnapshot(t)
	run := seededRunAtPort(snapshot, "mars")
	cache := godot.NewResourceCache()
	bridge := godot.AudioBridge{Cache: cache}

	if track := bridge.MusicForRoute(application.RouteMainMenu, &run, snapshot); track != godot.MainMenuMusicTrackID {
		t.Fatalf("expected main menu track %q, got %q", godot.MainMenuMusicTrackID, track)
	}
	if track := bridge.MusicForRoute(application.RouteTravel, &run, snapshot); track != godot.WorldMusicTrackID {
		t.Fatalf("expected travel route to use %q, got %q", godot.WorldMusicTrackID, track)
	}
	if track := bridge.MusicForRoute(application.RoutePortOverview, &run, snapshot); track != snapshot.PortsByID["mars"].MusicTrackID {
		t.Fatalf("expected current port music %q, got %q", snapshot.PortsByID["mars"].MusicTrackID, track)
	}
	if bridge.ShouldReplay(snapshot.PortsByID["mars"].MusicTrackID, snapshot.PortsByID["mars"].MusicTrackID) {
		t.Fatal("expected audio bridge to avoid replaying the same track")
	}
	if !bridge.ShouldReplay(godot.MainMenuMusicTrackID, godot.WorldMusicTrackID) {
		t.Fatal("expected audio bridge to switch tracks when route music changes")
	}
	if sfx := bridge.ClickSfxForAction("buy"); sfx != godot.DefaultClickSfxID {
		t.Fatalf("expected click sfx %q, got %q", godot.DefaultClickSfxID, sfx)
	}
}
