# Quickstart: StarSmuggler Core Game

## Goal

Get the Go + `go-dot` planning target into a repeatable development loop while preserving the existing Godot assets and using the MonoGame version as the parity reference.

## Prerequisites

- Godot 4.6.1 .NET-enabled editor installed
- Go 1.24.x installed
- .NET SDK installed for the legacy Godot C# shell build
- This repository checked out locally
- MonoGame reference project available at `/Users/oberon/Projects/coding/monogame/StarSmuggler`

## Reference Material To Keep Open

- `specs/001-space-smuggler-rpg/spec.md`
- `specs/001-space-smuggler-rpg/plan.md`
- `/Users/oberon/Projects/coding/monogame/StarSmuggler/GameManager.cs`
- `/Users/oberon/Projects/coding/monogame/StarSmuggler/Screens/*.cs`
- `/Users/oberon/Projects/coding/monogame/StarSmuggler/Content/Content.mgcb`

## Initial Developer Workflow

1. Check the available helper commands:

```bash
make help
```

2. Run the Go suite:

```bash
make test
```

3. Run the smoke checks:

```bash
make smoke-travel
make smoke-full
```

4. Build the existing Godot C# shell when you need to verify the legacy scene layer still boots:

```bash
make build-dotnet
```

5. Open the project in the .NET-enabled Godot editor for visual/manual checks:

```bash
make godot-open
```

6. Compare the current behavior against the MonoGame reference project for the same screen or loop slice before signing off on a change.

## MVP Implementation Order

1. Bootstrap the Go module and `go-dot` bridge layer.
2. Port immutable content loading and runtime state models.
3. Port economy, travel, event, save, and game-over logic into pure Go.
4. Rebind the six MVP screens through `go-dot` presenters and bridge helpers.
5. Verify parity with MonoGame behavior and visual composition.
6. Add post-MVP story/progression layers after the parity loop is stable.

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
- `make test` passes
- any new content/save schema has stable IDs and tests
- `make smoke` and `make smoke-go` still boot
- the screen matches the old game closely enough in layout, action flow, and feedback
- any intentional deviations from MonoGame behavior are documented

## MVP Parity Checklist

- Main menu matches button flow and music behavior
- Port overview matches information hierarchy and event display
- Trade screen matches terminal-driven buy/sell loop and pricing cues
- Travel screen matches preview/navigation/travel commitment flow
- Travel animation matches route pacing and transition role
- Game over matches stranded-state messaging

## Validation Status

Validated on 2026-04-18 with:

- `make help`
- `make test`
- `make build-dotnet`
- `make smoke-travel`
- `make smoke-full`

## Notes

- The existing C# Godot implementation is useful reference material during migration, but new feature work for this spec should target the Go runtime path.
- Preserve the imported custom art/audio unless the spec explicitly calls for a change.
