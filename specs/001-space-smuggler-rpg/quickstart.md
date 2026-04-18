# Quickstart: StarSmuggler Core Game

## Goal

Get the Go + `go-dot` MVP planning target into a runnable development loop while preserving the existing Godot project assets and using the MonoGame version as the parity reference.

## Prerequisites

- Godot 4.6.1 installed
- Go 1.24.x installed
- This repository checked out locally
- MonoGame reference project available at `/Users/oberon/Projects/coding/monogame/StarSmuggler`

## Reference Material To Keep Open

- `specs/001-space-smuggler-rpg/spec.md`
- `specs/001-space-smuggler-rpg/plan.md`
- `/Users/oberon/Projects/coding/monogame/StarSmuggler/GameManager.cs`
- `/Users/oberon/Projects/coding/monogame/StarSmuggler/Screens/*.cs`
- `/Users/oberon/Projects/coding/monogame/StarSmuggler/Content/Content.mgcb`

## Initial Developer Workflow

1. Open the Godot project from the repository root:

```bash
godot --path /Users/oberon/Projects/coding/godot/star-smuggler-go
```

2. Run Go tests from the repository root once the Go module bootstrap is in place:

```bash
go test ./...
```

3. Run a headless Godot smoke test to verify scene and asset boot:

```bash
godot --headless --path /Users/oberon/Projects/coding/godot/star-smuggler-go --quit-after 1
```

4. Compare the current behavior against the MonoGame reference project for the same screen or loop slice before signing off on a change.

## MVP Implementation Order

1. Bootstrap the Go module and `go-dot` bridge layer.
2. Port immutable content loading and runtime state models.
3. Port economy, travel, event, save, and game-over logic into pure Go.
4. Rebind the six MVP screens through `go-dot` presenters.
5. Verify parity with MonoGame behavior and visual composition.
6. Add the first narrative spine after parity is stable.

## Definition Of Ready For A New Slice

Before implementing any slice, confirm:

- the relevant behavior exists in the MonoGame reference
- the target Go package boundary is clear
- expected authored content IDs are defined
- tests can be written outside Godot scene code
- the screen/state contract for the slice is documented

## Definition Of Done For A Slice

A slice is done only when all of the following are true:

- gameplay logic lives in Go rather than scene callbacks
- `go test ./...` passes
- any new content/save schema has stable IDs and tests
- Godot smoke run still boots
- the screen matches the old game closely enough in layout, action flow, and feedback
- any intentional deviations from MonoGame behavior are documented

## MVP Parity Checklist

- Main menu matches button flow and music behavior
- Port overview matches information hierarchy and event display
- Trade screen matches terminal-driven buy/sell loop and pricing cues
- Travel screen matches preview/navigation/travel commitment flow
- Travel animation matches route pacing and transition role
- Game over matches stranded-state messaging

## Notes

- The existing C# Godot implementation is useful reference material during migration, but new feature work for this spec should target the Go runtime path.
- Preserve the imported custom art/audio unless the spec explicitly calls for a change.
