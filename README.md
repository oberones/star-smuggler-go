# StarSmuggler Go

StarSmuggler Go is the in-progress Go + Godot port of the MonoGame game in `/Users/oberon/Projects/coding/monogame/StarSmuggler`. The current goal is full MonoGame feature parity for the MVP, followed by post-MVP narrative and progression expansion.

## Current State

- The Godot project shell, imported assets, and C# reference implementation remain in place.
- The Go runtime bootstrap now lives at the repository root and will become the single gameplay authority as implementation continues.
- `src/*.cs` remains useful reference material during the migration, but new gameplay/runtime work should target the Go path defined in `specs/001-space-smuggler-rpg/`.

## Repository Layout

- `cmd/starsmuggler/`: Go entrypoint scaffold
- `internal/`: future pure-Go application, domain, services, persistence, content, and Godot bridge packages
- `scenes/`: Godot scenes and visual shell
- `assets/`: imported art, audio, and UI resources
- `data/`: JSON-authored ports, items, and events
- `specs/001-space-smuggler-rpg/`: active Speckit specification, plan, contracts, and tasks

## Developer Commands

Run Go tests from the repository root:

```bash
go test ./...
```

Run the Go bootstrap entrypoint:

```bash
go run ./cmd/starsmuggler
```

Build the existing Godot C# reference project:

```bash
dotnet build StarSmugglerGo.sln
```

Open the Godot project with the .NET-enabled editor:

```bash
open -a "/Applications/Godot .NET.app" /Users/oberon/Projects/coding/godot/star-smuggler-go
```

Run a headless Godot smoke boot:

```bash
'/Applications/Godot .NET.app/Contents/MacOS/Godot' --headless --path /Users/oberon/Projects/coding/godot/star-smuggler-go --quit-after 1
```

## Implementation Notes

- Use the MonoGame project as the parity reference before changing behavior or screen flow.
- Keep gameplay logic in Go packages and keep Godot scene scripts passive.
- Treat `.specify/memory/constitution.md` as the implementation guardrail for architecture, testing, UX consistency, and performance expectations.
