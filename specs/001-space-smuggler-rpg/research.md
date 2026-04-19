# Research: StarSmuggler Core Game

## Decision 1: Use the MonoGame project as the MVP parity reference

- **Decision**: Treat `/Users/oberon/Projects/coding/monogame/StarSmuggler` as the behavioral reference for the first playable Go MVP.
- **Rationale**: The old project already defines the essential loop, content set, screen flow, asset identity, audio behavior, and travel/event cadence. Re-deriving the MVP from scratch would introduce unnecessary product drift.
- **Alternatives considered**:
  - Use the current C# Godot port as the sole reference. Rejected because it is transitional and already diverges architecturally from the desired Go + `go-dot` target.
  - Redesign the loop immediately around deeper RPG systems. Rejected because it would bury parity issues under new scope.

## Decision 2: Keep Godot scenes/assets at the root and add a Go module beside them

- **Decision**: Preserve the existing Godot project layout at the repository root and add a root Go module with `internal/` packages for gameplay logic.
- **Rationale**: This keeps Godot editor workflows simple, preserves imported assets, and allows the Go runtime to grow without hiding scenes or art in a separate subproject.
- **Alternatives considered**:
  - Move all Godot content into a nested `godot/` directory. Rejected because it would create churn without solving a pressing technical problem.
  - Build a fully separate Go game runtime outside Godot. Rejected because the user explicitly wants `go-dot` as the main framework.

## Decision 3: Use `go-dot` only as a thin Godot-facing adapter layer

- **Decision**: Restrict `go-dot` code to node registration, scene lifecycle hooks, signal forwarding, resource binding, and presenter wiring.
- **Rationale**: This satisfies the constitution requirement that gameplay logic remains testable in pure Go and prevents the same scene-coupled problems that existed in the MonoGame/C# transition code.
- **Alternatives considered**:
  - Put trade, travel, and event rules directly in Godot scene callbacks. Rejected because it harms testability, determinism, and portability.
  - Build a heavier framework-style adapter abstraction up front. Rejected because it would add complexity before concrete repeated use cases exist.

## Decision 4: Preserve the MonoGame screen composition and terminal aesthetic for the MVP

- **Decision**: The MVP should visually mirror the old game as closely as practical: full-screen background art, centered terminal framing, button/icon placement hierarchy, and the same six-screen route.
- **Rationale**: The user explicitly wants the new game to look and feel like the old one within the limits of `go-dot`, and the old UI already communicates the game fantasy effectively.
- **Alternatives considered**:
  - Redesign around fully responsive modern UI containers from day one. Rejected for MVP because it would alter the game’s identity too early.
  - Port the old fixed-pixel coordinates verbatim. Rejected because Godot still needs scalable layout rules even when preserving the look.

## Decision 5: Preserve custom graphics, animation, sound, and music as shipped assets

- **Decision**: Use the imported custom PNG, WAV, and MP3 assets already present in this Godot repo as the initial MVP asset set.
- **Rationale**: The assets are already aligned with the MonoGame game’s look and feel, and retaining them minimizes art drift while the runtime architecture changes.
- **Alternatives considered**:
  - Replace assets during the runtime migration. Rejected because it couples presentation redesign to architecture risk.
  - Use placeholder Godot-native demo art/audio during MVP. Rejected because it breaks the parity requirement.

## Decision 6: Externalize authored content into JSON with stable IDs

- **Decision**: Ports, items, events, and future story/faction/mission definitions live in versioned JSON content files, not static Go code.
- **Rationale**: The MonoGame project hardcodes content in static classes. JSON content gives stable IDs, easier balancing, safer save references, and clearer authored-vs-runtime separation.
- **Alternatives considered**:
  - Recreate static databases in Go source. Rejected because it repeats the same rigidity and makes balancing/test fixtures harder.
  - Store authored content directly in scene metadata. Rejected because it tangles game rules with presentation files.

## Decision 7: Use deterministic services for economy, travel, and event resolution

- **Decision**: All random systems are routed through injectable RNG-aware services in pure Go.
- **Rationale**: The MonoGame code instantiates `Random` in many places, which makes balancing and tests inconsistent. Deterministic services are required for reproducible tests and safer save logic.
- **Alternatives considered**:
  - Keep ad hoc randomization in port/event code. Rejected because it prevents reliable regression testing.

## Decision 8: Defer narrative expansion until after the parity MVP

- **Decision**: The parity MVP will not add a new story spine. Narrative, faction, and mission systems begin immediately after the MonoGame feature set is fully reproduced in Go + `go-dot`.
- **Rationale**: This keeps the first milestone tightly aligned with the original game and avoids mixing architecture migration with product expansion.
- **Alternatives considered**:
  - Add one lightweight story spine during MVP. Rejected because it changes scope before parity is proven.
  - Attempt multiple factions and endings in the first pass. Rejected because it would swamp the migration effort.

## Decision 9: Adopt a Go-first testing strategy with Godot smoke verification

- **Decision**: Use `go test ./...` for unit and integration coverage, then layer headless Godot smoke runs on top for editor/runtime verification.
- **Rationale**: This matches the constitution and keeps most game behavior testable without scene boot. Godot smoke runs still catch asset resolution and wiring issues that pure Go tests cannot.
- **Alternatives considered**:
  - Rely only on manual playtesting in Godot. Rejected because it is too brittle for a migration.
  - Drive all tests through Godot scenes. Rejected because it would be slower, harder to debug, and more coupled to presentation.

## Decision 10: The current C# Godot implementation is reference material, not target architecture

- **Decision**: Existing `src/*.cs` code remains temporarily useful for parity checks and UI behavior comparison, but new MVP implementation work should land in Go rather than extending the C# path.
- **Rationale**: The user has explicitly selected Go + `go-dot` as the main framework. Continuing to deepen the C# architecture would create a second migration later.
- **Alternatives considered**:
  - Finish the C# implementation and port to Go later. Rejected because it doubles migration effort.

## Parity Audit: 2026-04-18

### Reference reviewed

- `/Users/oberon/Projects/coding/monogame/StarSmuggler/GameManager.cs`
- `/Users/oberon/Projects/coding/monogame/StarSmuggler/Game1.cs`
- `/Users/oberon/Projects/coding/monogame/StarSmuggler/Audio/AudioManager.cs`
- `/Users/oberon/Projects/coding/monogame/StarSmuggler/Screens/MainMenuScreen.cs`
- `/Users/oberon/Projects/coding/monogame/StarSmuggler/Screens/PortOverviewScreen.cs`
- `/Users/oberon/Projects/coding/monogame/StarSmuggler/Screens/TradeScreen.cs`
- `/Users/oberon/Projects/coding/monogame/StarSmuggler/Screens/TravelScreen.cs`
- `/Users/oberon/Projects/coding/monogame/StarSmuggler/Screens/TravelAnimationScreen.cs`
- `/Users/oberon/Projects/coding/monogame/StarSmuggler/Screens/GameOverScreen.cs`

### Parity confirmed

- The Go runtime preserves the same six-screen loop: main menu, port overview, trade, travel, travel animation, and game over.
- Starter economy parity remains aligned with MonoGame: 500 starting credits, 30 cargo capacity, random inner-port start, 15-credit base travel cost, zone-difference scaling, and every-fourth-jump market refresh.
- The current Godot presentation still keeps the old full-screen background art, centered terminal framing, cockpit travel screen, and travel-animation role.
- Main-menu music vs. route music intent matches the MonoGame reference: `singularity` for menu/game-over, current-port music for port/trade, and neutral travel music for route transitions.
- Click SFX remain bound to the same class of menu/trade/travel actions as the MonoGame flow.

### Intentional deviations

- The Go port externalizes content into JSON and isolates gameplay in pure Go services instead of using the MonoGame singleton/stateful-screen pattern. This is a deliberate architectural correction, not a gameplay change.
- The current trade flow returns to port overview after trading rather than forcing an immediate travel handoff via a `Done` button. This improves loop readability while preserving the same effective decision cadence.
- The travel screen currently exposes destination choices through a presenter-generated list while still preserving the cockpit/preview composition, instead of the exact MonoGame previous/next carousel. The destination information and commitment pressure remain equivalent.
- Story, faction, mission, and upgrade systems now exist beyond the MonoGame reference. They are additive layers on top of the parity loop rather than replacements for its baseline rules.
- The Go-side smoke/test harness provides deterministic coverage that the MonoGame project did not have. This is an engineering improvement rather than a player-facing deviation.

### Remaining polish notes

- The Go presentation helpers now model route-audio decisions and resource caching, but the final `go-dot` runtime still needs to consume those helpers directly once the scene authority fully shifts away from the C# shell.
- Headless Godot smoke still reports the known shutdown warning about leaked Godot resources on exit. It has not blocked the validated gameplay/test flow, but it remains worth revisiting after the full `go-dot` handoff.
