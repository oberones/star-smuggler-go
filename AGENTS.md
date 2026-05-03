# AGENTS.md

## Project Summary

StarSmuggler Go is a Go + Godot port of the MonoGame StarSmuggler game at `/Users/oberon/Projects/coding/monogame/StarSmuggler`. The port preserves the existing trading, travel, event, save/load, audio, and terminal-style screen flow while moving gameplay authority into testable Go packages.

The current active direction is the Spec Kit plan in `specs/001-space-smuggler-rpg/`: Go 1.24 gameplay/runtime code, Godot 4.6.1 scenes/assets, and thin Godot integration. The older `docs/PORTING_PLAN.md` is still useful for reference inventory and architecture warnings, but parts of it describe an earlier C#-first target. Prefer the active spec and constitution when they disagree.

Core player loop:

1. Main Menu
2. Port Overview
3. Trade or Travel
4. Travel Animation
5. Arrival, event/story/progression updates, or Game Over

Current implemented gameplay areas include economy/trading, route-cost travel, random travel events, game-over and recovery handling, story/faction/mission state, and ship upgrades/specializations.

## Primary Sources Of Truth

- `.specify/memory/constitution.md`: project guardrails. Gameplay logic belongs in Go, engine boundaries stay thin, behavior must be testable, content/save formats need stable IDs and deterministic behavior.
- `specs/001-space-smuggler-rpg/spec.md`: player-facing goals and requirements.
- `specs/001-space-smuggler-rpg/plan.md`: current architecture and phase plan.
- `specs/001-space-smuggler-rpg/tasks.md`: implementation history and intended task ordering.
- `specs/001-space-smuggler-rpg/contracts/`: screen-flow and content/save contracts.
- `README.md` and `specs/001-space-smuggler-rpg/quickstart.md`: developer commands and workflow.
- `docs/ROADMAP.md`, if present: current product-direction notes. Verify against code and specs because roadmap docs can lag implementation.
- `docs/NOTES.md`, if present: rough worldbuilding and backlog scratchpad, not an implementation contract.

## Repository Layout

- `cmd/starsmuggler/`: small Go bootstrap entrypoint. It is not yet the whole game runtime.
- `internal/domain/`: plain Go data types for authored definitions and runtime state. Keep this package free of Godot concepts.
- `internal/content/`: JSON content loading, default content paths, and content validation into `domain.DataSnapshot`.
- `internal/services/`: pure gameplay services: economy, trade, travel, events, run evaluation, story, factions, missions, upgrades, and runtime context/RNG.
- `internal/application/`: app-level orchestration. This layer owns routes, active runs, command flow, autosave decisions, travel resolution, recovery, mission acceptance, and upgrade purchase integration.
- `internal/persistence/`: JSON save DTOs, save versioning, save/load, hydration, and dehydration.
- `internal/presentation/godot/`: thin presentation bridge and view-model presenters for Godot. This package shapes state for screens and resolves resource/audio identifiers, but should not own business rules.
- `src/`: existing Godot C# shell and transitional reference implementation. It still boots the current Godot project and mirrors many gameplay concepts, but new gameplay authority should target Go. Several C# runtime-owner nodes become passive when `STARSMUGGLER_GO_RUNTIME=1`.
- `scenes/`: Godot `.tscn` files. `scenes/app/App.tscn` is the main scene; `scenes/screens/` contains the six core screens.
- `assets/`: imported art, audio, UI, port, screen, and trade resources. Preserve these unless the spec explicitly changes them.
- `data/`: authored JSON content for ports, items, events, factions, missions, story arcs, and ship upgrades.
- `tests/golden/`: content snapshot and golden compatibility tests.
- `tests/integration/`: Go integration coverage for gameplay, app orchestration, persistence, story, upgrades, and presentation bridge behavior.
- `tests/smoke/`: headless Godot smoke scripts.
- `.github/agents/` and `.github/prompts/`: Spec Kit support files for Copilot-style agents.
- `docs/`: porting notes, roadmap, and scratchpad planning. Treat these as context unless they explicitly update the active spec.

## Architecture Rules

- Put gameplay rules in Go, not in Godot scenes, C# screen scripts, or `_Ready()` side effects.
- Keep `domain` independent of Godot, file systems, clocks, and random globals.
- Use `services.RuntimeContext` and injected RNG for randomized systems so tests can be deterministic.
- Treat `application.App` and command structs as the place for cross-screen flow, route changes, autosaves, and multi-service coordination.
- Keep `internal/presentation/godot` as a presenter/adapter boundary. It may compose view models and resolve resources, but it should not calculate trade validity, route costs, game-over state, mission outcomes, or event effects.
- Keep screens passive. The screen contract is intent signals out, explicit view models in.
- Do not recreate a broad singleton/game-manager architecture. The port intentionally splits immutable definitions, mutable run state, focused services, persistence DTOs, and presentation adapters.
- Do not add a separate MVP fuel mechanic. Travel pressure is route credit cost, jump-count progression, event risk, and economy pressure.

## Content And Save Rules

- JSON files under `data/` are the source of truth for authored content.
- Authored entities must use stable string IDs. Do not use display names as save keys or runtime references.
- Do not rename or remove content IDs without a migration plan and tests.
- Save data lives in `internal/persistence` DTOs and currently uses `CurrentSaveVersion = 1`.
- Saves must contain runtime state only. Never serialize Godot nodes, object identity, scene paths as authority, or display-name fallbacks.
- Schema or content-contract changes need focused tests, and golden snapshots should be updated intentionally.

## C# And Godot Notes

- `project.godot` is a Godot 4.6.1 .NET project with `scenes/app/App.tscn` as the main scene and a 1536x1024 reference viewport.
- `StarSmugglerGo.csproj` targets `net8.0` using `Godot.NET.Sdk/4.6.1`.
- The C# tree in `src/` is useful for current scene boot behavior and parity checks, but it is not the preferred place for new gameplay/runtime logic.
- Runtime-owner C# nodes check `STARSMUGGLER_GO_RUNTIME`; when set to `1` or `true`, they become passive so Go can take over authority during migration.
- If a scene change requires C# glue, keep the glue narrow and mirror the Go screen-flow contract.

## Common Commands

Run these from the repository root:

```bash
make help
make test
make test-integration
make test-golden
make smoke
make smoke-go
make smoke-travel
make smoke-full
make build-dotnet
make run-go
make fmt
```

Useful direct commands:

```bash
go test ./...
go run ./cmd/starsmuggler
dotnet build StarSmugglerGo.sln
STARSMUGGLER_GO_RUNTIME=1 bash tests/smoke/run_headless_smoke.sh
```

The smoke runner writes Godot logs to an explicit temp `--log-file`; keep that behavior because headless Godot can crash on the default macOS `user://logs` path.

## Testing Expectations

- Run `make test` for normal Go validation.
- Run `make smoke-travel` or `make smoke-full` after changes that affect screen flow, travel, resources, audio, or the Godot bridge.
- Run `make build-dotnet` when touching `src/`, `.tscn` script bindings, `project.godot`, or C# project files.
- Add or update tests at the same layer as the behavior: service tests for rules, integration tests for app/persistence/contracts, golden tests for content/save compatibility, and smoke tests for Godot boot/flow confidence.

## Working Guidance For Future Agents

- Start with `git status --short --branch` and do not overwrite user changes.
- Read the active spec, constitution, and nearby tests before making architectural changes.
- Prefer small, focused Go services and explicit command methods over broad managers or utility packages.
- When adding content, update `data/`, validation rules if needed, and tests together.
- When changing player-visible behavior, check whether the MonoGame parity target or spec requires the old behavior.
- When docs disagree, update the stale doc rather than quietly following it.
- Keep generated/import metadata such as Godot `.import` and `.uid` files only when they are actually produced by the intended editor/tooling change.
