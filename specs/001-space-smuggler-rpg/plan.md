# Implementation Plan: StarSmuggler Core Game

**Branch**: `001-space-smuggler-rpg` | **Date**: 2026-04-18 | **Spec**: [spec.md](./spec.md)  
**Input**: Feature specification from `/specs/001-space-smuggler-rpg/spec.md`

## Summary

Build the Go + `go-dot` version of StarSmuggler by using the existing MonoGame project at `/Users/oberon/Projects/coding/monogame/StarSmuggler` as the initial MVP reference. The first milestone is not a redesign. It is a faithful port of the current playable loop, visual composition, music behavior, screen flow, and content set into a maintainable Go architecture that keeps domain logic outside Godot scene scripts. The MVP is strict parity with the MonoGame loop; deeper story, faction, mission, and progression systems begin in the next milestone.

## Technical Context

**Language/Version**: Go 1.24.x for gameplay/runtime code, Godot 4.6.1 project assets and scenes  
**Primary Dependencies**: `go-dot` as the main Godot integration framework, Godot 4.6.1 editor/runtime, Go standard library for serialization/testing/tooling  
**Storage**: JSON files for authored content and save data, imported Godot resources for graphics/audio  
**Testing**: `go test ./...`, table-driven unit tests, JSON golden tests for content/save compatibility, headless Godot smoke runs for screen boot verification  
**Target Platform**: Desktop Godot game for macOS, Windows, and Linux  
**Project Type**: Single-player desktop game using Godot scenes plus internal Go gameplay packages  
**Performance Goals**: Stable 60 FPS on the MVP screens, no visible stutter on route changes, no avoidable resource reload hitches during menu navigation, travel animation, or audio transitions  
**Constraints**: Must preserve the MonoGame look and feel as closely as practical, must keep `go-dot` adapters thin, must keep core gameplay deterministic and testable outside scene callbacks, must continue using custom graphics, animations, sound, and music  
**Scale/Scope**: Initial MVP includes 6 core screens, 9 ports, 18 items, 7 travel events, save/load, music/SFX transitions, and the complete MonoGame functional loop with no new narrative systems

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **Idiomatic Go, Thin Engine Boundaries**: PASS
  The plan keeps gameplay, economy, travel, events, save migration, and story flags in Go packages. `go-dot` is limited to node registration, signal forwarding, scene state application, and resource hooks.
- **Testable By Default**: PASS
  The plan requires table-driven tests for economy, travel, event resolution, progression rules, and save compatibility before the corresponding slice is considered complete.
- **Consistent Player Experience**: PASS
  The MVP explicitly preserves the old MonoGame screen order, art direction, terminal framing, route presentation, audio cues, and information hierarchy.
- **Performance Budgets Are Product Requirements**: PASS
  The plan avoids per-frame logic rebuilds, centralizes resource caching, and treats screen-transition hitches and audio reload stutter as blocking regressions.
- **Data Integrity And Determinism**: PASS
  The plan uses stable IDs, immutable authored content, explicit runtime state, versioned saves, and injectable RNG for testable event/economy behavior.

**Post-Design Re-check**: PASS
No constitution violations are required for the MVP architecture. The only intentional complexity is a thin Godot bridge plus pure Go domain packages, which is mandated by the constitution rather than an exception.

## Project Structure

### Documentation (this feature)

```text
specs/001-space-smuggler-rpg/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   ├── mvp-screen-flow.md
│   └── content-save-contract.md
└── tasks.md
```

### Source Code (repository root)

```text
assets/
├── audio/
├── ports/
├── screens/
├── trade/
└── ui/

data/
├── events/
├── items/
└── ports/

scenes/
├── app/
├── components/
└── screens/

internal/
├── application/
├── content/
├── domain/
├── persistence/
├── presentation/
│   └── godot/
├── services/
└── testing/

cmd/
└── starsmuggler/

tests/
├── golden/
├── integration/
└── smoke/
```

**Structure Decision**: Keep the Godot project root, scenes, assets, and imported resources in place so the visual/audio pipeline remains editor-friendly. Introduce a root Go module and place all gameplay code in `internal/*` packages. `presentation/godot` owns `go-dot` registration and scene-to-view-model binding. The existing `src/*.cs` tree becomes temporary parity reference material until the Go equivalent replaces it slice by slice.

## Phase Plan

### Phase 0 - Reference Lock And Bootstrap

1. Freeze the MVP parity target using the MonoGame codebase as the source of truth for:
   - screen order: Main Menu -> Port Overview -> Trade -> Travel -> Travel Animation -> Game Over
   - starter values: 500 credits, 30 cargo, random inner-zone start
   - travel cost rules, item rarity distribution, event set, and market loop
   - baseline art/audio identity: terminal UI, port backgrounds, cockpit/travel screen, singularity/world music, click SFX
2. Add the root Go module and bootstrap `go-dot` integration without deleting current Godot assets or scenes.
3. Decide and document the Go runtime entrypoints and Godot bridge points.
4. Add headless smoke infrastructure:
   - `go test ./...`
   - Godot headless launch for scene boot
5. Mark the current C# Godot implementation as transitional reference, not the target architecture.
6. Disable the existing C# gameplay path as an active runtime authority once the Go bootstrap is ready:
   - retain `src/*.cs` only as migration reference material
   - ensure Godot scenes no longer depend on C# business logic for gameplay state changes
   - route all gameplay mutations through the Go application layer and `go-dot` bridge

### Phase 1 - Domain Extraction And Data Contracts

1. Convert the MonoGame static databases into stable JSON-backed authored content with immutable IDs:
   - ports
   - items
   - events
   - future story/faction/mission definitions
2. Build pure Go domain types for:
   - run state
   - player ship state
   - cargo/inventory
   - market snapshot
   - travel route calculation
   - event definition and resolved event result
   - faction standing and story flags
3. Build persistence DTOs with explicit save versioning from the start.
4. Introduce deterministic RNG injection for price generation, event rolling, and content selection.
5. Add tests for:
   - travel cost parity with MonoGame
   - market generation parity with MonoGame
   - event effects and message formation
   - save round-trips and forward migration hooks

### Phase 2 - MVP Parity Loop In Go

1. Replace the core C# gameplay orchestration with Go services and app coordinators:
   - economy service
   - trade service
   - travel service
   - event service
   - run evaluator
   - save service
2. Keep the current Godot scenes as the visual shell, but rebind them through `go-dot` presenters/view-model binders.
3. Port the six core screens to the Go runtime contract:
   - Main Menu
   - Port Overview
   - Trade
   - Travel
   - Travel Animation
   - Game Over
4. Preserve the MonoGame look-and-feel targets:
   - centered terminal framing
   - current background art and port previews
   - travel animation timing based on route distance
   - audio crossfade behavior where practical
5. Reach parity on the original loop before deepening it:
   - new game
   - continue/save
   - buy/sell
   - route selection
   - travel event resolution
   - stranded/game-over evaluation

### Phase 3 - Story And World Reactivity

1. Add a first narrative spine that grows directly out of trade and travel:
   - one faction conflict
   - one named smuggling contact chain
   - one story cargo/job arc
2. Track world-state consequences in Go runtime state:
   - faction standing
   - unlocked ports or route restrictions
   - mission success/failure flags
   - story milestone triggers
3. Keep narrative delivery within the established UI language:
   - port overview notices
   - travel event choices
   - mission/event text in terminal panels
4. Test story gating and branching as state transitions, not scene-only behavior.

### Phase 4 - Progression, Polish, And Expansion Hooks

1. Add ship upgrades and specialization paths in ways that do not break MVP clarity.
2. Expand the event system to support choice-driven events and mission-linked encounters.
3. Improve accessibility and consistency:
   - keyboard/controller navigation
   - consistent button focus and feedback
   - text readability over custom backgrounds
4. Optimize asset loading and screen transitions:
   - cache common textures/audio
   - avoid redundant resource loads on route changes
5. Prepare clean extension points for additional factions, contracts, sectors, and endings.

## MVP Slice Definition

The initial implementation target for this feature is:

- Go + `go-dot` powered runtime
- one playable campaign/run mode
- full parity with the MonoGame economy/travel loop
- imported custom art, music, SFX, and travel animation assets
- save/continue support
- the same six-screen flow as the MonoGame game
- for task-planning purposes, MVP completion requires both the core economy loop and the full travel/event loop from the original MonoGame game
- no required narrative spine in the MVP; story expansion begins immediately after parity is stable

The following are intentionally deferred until parity is stable:

- the first faction-driven narrative spine
- multiple campaigns or endings
- combat-heavy systems
- advanced modding
- online or multiplayer features
- large faction matrices beyond the first reactive story layer

## Validation Strategy

### Automated Validation

- `go test ./...` MUST cover:
  - price generation and markup logic
  - cargo constraints
  - travel cost calculation
  - event resolution
  - game-over evaluation
  - save encoding/decoding and migration
- Golden tests MUST verify authored JSON content and save schemas remain stable.
- Integration tests MUST verify the application coordinator can drive the MVP loop without requiring scene-owned business logic.

### Runtime Validation

- Headless Godot smoke runs MUST confirm that:
  - the app scene boots
  - required assets resolve
  - route transitions do not crash
  - save/load bootstrap works on a clean project state

### Parity Validation

- Each MVP screen MUST be compared against the MonoGame reference for:
  - visible information hierarchy
  - available actions
  - audio behavior
  - background/terminal composition
  - progression timing and route feedback

## Risks And Mitigations

- **Risk**: The current repo is still C#-centric while the target architecture is Go + `go-dot`.  
  **Mitigation**: Treat the C# code as a temporary reference and replace logic slice by slice behind stable scenes and data files.

- **Risk**: `go-dot` integration may tempt gameplay logic into node callbacks.  
  **Mitigation**: Enforce a strict presenter/adapter boundary and reject scene-owned business rules in review.

- **Risk**: A story-first expansion could destabilize MVP parity.  
  **Mitigation**: Lock the initial slice to MonoGame parity first and add only one narrative spine after the economy/travel loop is stable.

- **Risk**: Custom art/audio fidelity may drift during migration.  
  **Mitigation**: Preserve asset names, route music intent, terminal framing, and screen composition as explicit parity checkpoints.

- **Risk**: Save compatibility becomes brittle as the story layer grows.  
  **Mitigation**: Version saves immediately and keep authored content IDs stable from the beginning.

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| None | N/A | N/A |
