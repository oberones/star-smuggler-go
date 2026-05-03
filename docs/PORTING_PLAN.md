# StarSmuggler Godot Port Plan

## Executive Summary

StarSmuggler already has a playable MonoGame core loop. The Godot port should not be treated as a straight file-for-file rewrite. It should preserve the existing game rules and content while replacing the custom screen/UI framework with Godot's scene tree, `Control`-based UI, signals, and built-in animation/audio systems.

The current Godot project is effectively a blank `4.6.1` shell, so the plan includes both project bootstrapping and the actual feature migration work.

Recommended target:

- Godot `4.6.1`
- Godot `C#` for gameplay/domain logic
- Godot scenes and `Control` nodes for screens and menus
- Externalized data for ports, items, and events instead of hardcoded static C# lists

Primary goal:

1. Reach feature parity with the current playable MonoGame game.
2. Make the Godot project the new home for future work such as contracts, upgrades, and factions.

## Current Source Inventory

The MonoGame project currently contains:

- 6 game screens: Main Menu, Port Overview, Trade, Travel, Travel Animation, Game Over
- 9 active ports across 3 zones
- 18 active items across 3 rarities
- 7 active random travel events
- Save/load via JSON
- Music/SFX support with song crossfading
- A custom terminal-style UI component set

### Current Gameplay Loop

1. Start a new run at a random Inner-zone port with 500 credits and 30 cargo space.
2. View the current port overview and any recent event.
3. Enter the trade screen to buy and sell goods.
4. Enter the travel screen to select a destination and pay travel costs.
5. Play a travel animation, then process arrival, market refreshes, and random events.
6. Return to the port overview or hit game over if stranded.

### Current Systems Mapped From Code

#### Core Runtime

- `GameManager` is the central authority for player state, pricing, travel, events, and save/load flow.
- `ScreenManager` swaps the active screen based on a `GameState` enum.
- `Game1` drives MonoGame initialization, input polling, screen updates, and drawing.

#### Data and Rules

- `PortsDatabase` stores the static port list.
- `ItemsDatabase` stores the static item list and rarity/zone selection helpers.
- `EventDatabase` stores random travel events and inline gameplay effects.
- `PlayerData` stores credits, cargo, current port, prices, event, and jump count.
- `SaveLoadManager` writes JSON save data to app data.

#### Presentation

- `MainMenuScreen`, `PortOverviewScreen`, `TradeScreen`, `TravelScreen`, `TravelAnimationScreen`, `GameOverScreen`
- Reusable UI widgets: `Button`, `BackButton`, `Terminal`, `InfoPanel`, `NumericInput`
- `AudioManager` for music fade and click SFX
- `AnimatedTexture` for spritesheet playback during travel

### MonoGame Pain Points To Fix During Port

- Heavy reliance on a singleton `GameManager`
- Screen flow tightly coupled to game logic
- Hardcoded content definitions in static C# classes
- Fixed pixel layouts tuned for `1536x1024`
- Custom button/input rendering and hit testing for basic menu behavior
- Runtime state keyed partly by object identity and display names
- Minimal automated validation around economy/save/event behavior
- Static definition objects are mutated at runtime
- Presentation screens currently trigger important gameplay side effects directly

### Architecture Changes To Make On Purpose

These are not optional cleanup items. They are the main architectural improvements the port should make:

#### 1. Split Immutable Definitions From Mutable Run State

Current issue:

- `Port` objects contain both authored data and per-visit runtime state.
- cargo is keyed by `Item` object references.
- save data relies on names instead of stable internal IDs.

Godot target:

- immutable definitions for ports, items, events, and future contracts/upgrades
- separate runtime state objects for player inventory, current market prices, available goods, recent event result, and progress flags

#### 2. Replace The God Object With Focused Services

Current issue:

- `GameManager` owns state, navigation, economy, travel, events, and game-over evaluation.

Godot target:

- `RunState` for the current play session
- `EconomyService` for price generation and trade validation
- `TravelService` for route cost and arrival processing
- `EventService` for random event selection and effect resolution
- `RunEvaluator` for game-over and other state checks
- `SaveService` for persistence

#### 3. Keep Presentation Passive

Current issue:

- the travel animation screen currently performs the actual arrival logic when the animation completes
- screens decide too much about state transitions

Godot target:

- scenes emit user intent and animation completion signals
- domain/application services perform the actual state changes
- UI never becomes the source of truth for business rules

#### 4. Redesign Events As Data + Resolved Results

Current issue:

- events are static objects with inline `Action<PlayerData>` logic
- event descriptions are mutated at runtime
- this invites shared-state bugs and blocks serialization/data tooling

Godot target:

- immutable `EventDefinition`
- typed effect descriptors and parameters
- runtime `EventResult` object containing rolled values and final player-facing text

#### 5. Introduce Save Versioning Early

Current issue:

- the current save model is small and workable, but brittle for future expansion

Godot target:

- explicit save version field
- migration path for future schema changes
- DTOs separate from domain models and scene objects

#### 6. Separate App Routing From In-Run State

Current issue:

- one enum currently covers multiple concepts at once: menu flow, gameplay screens, and terminal states

Godot target:

- `AppRoute` or scene route for navigation
- `RunPhase` or equivalent for gameplay state if needed
- screens should reflect state, not define it

#### 7. Centralize Randomness

Current issue:

- `new Random()` is created in multiple places

Godot target:

- one injected RNG source per run or per service boundary
- easier debugging, repeatability, and balancing

## Porting Goals

### Goal 1: Preserve Existing Playable Scope

The Godot port should preserve the current shipping slice before adding new mechanics:

- economy rules
- travel cost rules
- item rarity and zone behavior
- random event behavior
- save/load
- music and SFX behavior
- all current screens
- all current ports, items, and art/audio assets

### Goal 2: Rebuild UI As Godot-Native Scenes

The MonoGame UI layer is exactly where Godot should help most. The port should replace the current custom screen/widget framework with:

- `Control` scenes
- reusable Godot UI components
- theme-driven fonts/colors
- anchors and containers instead of fixed coordinates
- signal-driven interactions instead of manual polling everywhere

### Goal 3: Prepare For Future Expansion

The new project should be ready for:

- delivery contracts
- ship upgrades
- port condition modifiers
- faction systems
- narrative/quest content

That means the architecture should separate:

- static definitions
- runtime state
- UI scenes
- persistence
- audio/animation helpers
- application flow/orchestration

## Recommended Technical Direction

### Use Godot C# For Game Logic

Use Godot C# as the primary scripting language for the port.

Why:

- The existing game logic is already in C#.
- Travel, pricing, save/load, and event systems can be ported with less translation risk.
- Domain tests are easier to preserve when rules stay in the same language.
- The game is logic-heavy even though it is UI-driven.

GDScript can still be used for tiny scene-local scripts if useful, but the main domain layer should stay in C#.

One important constraint:

- gameplay rules should be runnable headlessly outside the scene tree so they can be tested without loading UI scenes.

### Use Godot Scenes For Screens, Not A Custom ScreenManager Clone

Do not recreate the MonoGame `IScreen` + `ScreenManager` model 1:1.

Instead:

- create one root application scene
- load/swap child screen scenes
- use signals for transitions and actions
- keep scene responsibilities narrow

Good candidate scenes:

- `MainMenuScreen.tscn`
- `PortOverviewScreen.tscn`
- `TradeScreen.tscn`
- `TravelScreen.tscn`
- `TravelAnimationScreen.tscn`
- `GameOverScreen.tscn`

Good candidate orchestration layer:

- one application coordinator that owns route changes
- one run/session coordinator that owns gameplay transitions

### Separate Definitions From Runtime State

Move hardcoded data out of static classes and into data files.

Recommended split:

- authored definitions as custom Godot `Resource` files when they need editor-assigned asset references
- JSON save data for runtime persistence
- runtime C# classes for `RunState`, `PlayerState`, `MarketSnapshot`, `CargoState`, `TravelRequest`, `TravelOutcome`, and `EventResult`

Important rule:

- Use stable IDs everywhere.
- Do not key save data by display names.
- Do not key cargo state by `Item` object references.
- Do not store mutable market state on definition resources.

Suggested definition model:

- `PortDefinition`: id, name, description, zone, background texture, preview texture, trade background, music
- `ItemDefinition`: id, name, description, rarity, base price
- `EventDefinition`: id, name, base description, effect type, effect parameters, weighting/tags

Suggested runtime model:

- `RunState`: player, active market snapshots, jump counters, recent event result, active modifiers
- `CargoState`: item ID to quantity
- `MarketSnapshot`: port ID, available item IDs, item prices, active local modifiers
- `EventResult`: event ID, resolved text, rolled values, applied consequences

### Use Autoload Services Sparingly

Recommended autoloads:

- `AppController`
- `GameSession`
- `SceneRouter`
- `SaveService`
- `AudioService`
- `DataRepository`

Keep them focused. Avoid building a new giant Godot singleton that repeats the current `GameManager` problem.

If two autoloads start depending on each other heavily, collapse responsibilities downward into plain C# services instead of upward into more global state.

### Introduce A Thin Application Layer

Use a small application layer between UI scenes and domain services.

Responsibilities:

- accept intents from UI
- call domain services
- update `RunState`
- publish state changes back to scenes

This is the layer that should answer questions like:

- what happens when the player clicks Travel
- what screen opens after arrival
- when should autosave occur
- when does a game-over check run

It should not contain rendering code, and scenes should not reimplement its rules.

### Redesign The Event System Instead Of Porting It Literally

Recommended shape:

- event selection based on weights, tags, and eligibility predicates
- effect handlers implemented as typed service logic
- descriptions templated from rolled values

Example improvement:

- instead of mutating a shared event description from `"Pirate Ambush"` into a rolled sentence, resolve a fresh `EventResult` like `"Pirate Ambush: lost 83 credits"` for this run only

This keeps authored data reusable and avoids state bleed across runs.

### Keep Domain Logic Out Of Animation And Audio Scripts

Animation and audio scripts should react to domain events, not own them.

Examples:

- travel animation completes and emits `travel_animation_finished`
- the application layer then applies `TravelOutcome`
- audio fades when the route changes, but the audio system does not decide which gameplay state comes next

### Prefer Theme And Component Composition Over One-Off UI Scripts

The MonoGame code needed bespoke widgets because the engine offered little UI support.

In Godot:

- use a project theme
- build a small reusable component library
- keep screen scripts focused on data binding and signals

Good shared components:

- `TerminalPanel`
- `ActionButton`
- `QuantityStepper`
- `PortSummaryCard`
- `EventBanner`
- `CreditsCargoBar`

### Suggested Repository Structure

```text
star-smuggler-go/
├── assets/
│   ├── audio/
│   ├── fonts/
│   ├── ports/
│   ├── trade/
│   └── ui/
├── data/
│   ├── events/
│   ├── items/
│   └── ports/
├── scenes/
│   ├── app/
│   ├── components/
│   └── screens/
├── src/
│   ├── autoload/
│   ├── application/
│   ├── domain/
│   ├── persistence/
│   ├── presentation/
│   └── services/
└── tests/
```

## What To Port Directly vs Rebuild

### Safe To Port Mostly As-Is

- travel cost formula
- rarity/zone markup logic
- price update cadence
- cargo capacity rules
- game over checks
- event outcome rules
- save/load field coverage

### Rebuild In Godot-Native Form

- screen lifecycle
- button/input handling
- layout and scaling
- screen transitions
- audio playback plumbing
- animated travel presentation
- font setup
- reusable UI components
- app-flow orchestration
- event definitions/runtime event resolution
- market/runtime state ownership

### Explicitly Do Not Port These Patterns

- one giant singleton owning every system
- mutable static content objects
- object-reference keyed cargo dictionaries
- save data keyed by names meant for display
- gameplay logic hidden inside screen scripts
- animation completion callbacks that directly mutate run state
- per-screen ad-hoc randomness

## Phase Plan

### Phase 0: Bootstrap The Godot Project

Goal: turn the blank Godot project into a usable C# project skeleton.

Tasks:

- enable and verify Godot C# support for this project
- create the scene/script/data/assets folder structure
- create a root app scene
- create a lightweight application coordinator and routing model
- define project display settings and target aspect strategy
- set up an input map for menu confirm/cancel/navigation
- create a first UI theme and font imports

Deliverables:

- app boots into a placeholder main menu scene
- C# scripts build successfully in Godot
- base theme is available project-wide
- the app has a clear place for routing and orchestration logic

### Phase 1: Extract And Stabilize The Domain Layer

Goal: lock down gameplay rules before UI work expands.

Tasks:

- create C# domain classes for definitions, runtime state, pricing, travel, and event resolution
- move static content from MonoGame classes into immutable definition files
- create application services with narrow responsibilities instead of a direct `GameManager` replacement
- centralize RNG usage instead of creating ad-hoc `Random` instances
- separate authored event definitions from resolved runtime event results
- write characterization tests for:
  - travel cost calculation
  - item markup rules
  - price generation/update cadence
  - cargo capacity checks
  - game over conditions
  - event application effects

Deliverables:

- data loads cleanly from files
- domain rules can run without scenes
- tests protect parity with the MonoGame behavior
- no runtime state is stored on definition objects

### Phase 2: Persistence And Save Migration

Goal: rebuild save/load cleanly in Godot and avoid data-model debt.

Tasks:

- create a new save schema using stable IDs
- store saves under `user://save.json`
- support credits, cargo, port, prices, jumps, and recent event state
- include a save version number from day one
- persist only DTOs, never scene nodes or Godot runtime objects
- decide whether to support importing MonoGame saves from the old JSON format
- add tests around save/load round-tripping

Recommendation:

- Support a one-time MonoGame save importer if continuity matters.
- If not, explicitly declare a save reset at the start of the Godot era.

Deliverables:

- save/load works in the Godot project
- schema is future-proof for contracts and upgrades

### Phase 3: Screen And Navigation Skeleton

Goal: establish the new screen flow without full polish.

Tasks:

- build the 6 core screen scenes
- create a root navigation flow between them
- wire placeholder data into each screen
- bind screens to application/domain state instead of embedding business rules in them
- add a transition helper for fades/slides if desired

Deliverables:

- user can move through the full scene graph
- no screen depends on hardcoded MonoGame-style drawing code
- navigation and gameplay orchestration live outside the scenes

### Phase 4: Port Overview And Trade Parity

Goal: recreate the game's heaviest UI screens first.

Tasks:

- build the port overview scene with background art, event text, and action buttons
- build reusable Godot equivalents of:
  - info panel
  - terminal panel
  - primary/secondary buttons
  - quantity stepper
- implement buy/sell logic with immediate UI refresh
- add disabled/error states for invalid actions
- keep validation and transaction rules inside domain/application services
- preserve current flow into travel and back

Deliverables:

- player can trade end-to-end in Godot
- cargo and credit changes are visible immediately
- event text and port flavor are displayed cleanly

### Phase 5: Travel Selection And Travel Animation Parity

Goal: port the travel loop and arrival processing.

Tasks:

- build the travel selection screen with destination preview and cost display
- prevent invalid travel choices cleanly
- implement the travel animation scene using `AnimationPlayer`, `Tween`, or `AnimatedSprite2D`
- let the animation scene signal completion, then trigger actual travel state changes in the application/domain layer
- on arrival:
  - deduct cost
  - move to the new port
  - refresh goods
  - increment jumps
  - update prices if needed
  - trigger a random event
  - save the run
  - check for game over

Deliverables:

- travel is fully functional
- animation and skip behavior work
- arrival processing matches MonoGame behavior
- travel presentation is decoupled from travel state mutation

### Phase 6: Audio, Assets, And UI Polish

Goal: match or exceed the current game's feel.

Tasks:

- import all current art/audio assets into Godot folders
- recreate the terminal-inspired style with Godot themes
- rebuild music crossfading using Godot audio players and buses
- add hover, disabled, pressed, and success/failure feedback states
- make layouts responsive at common desktop sizes
- refine travel animation presentation

Deliverables:

- assets are fully wired
- UI feels intentional rather than placeholder
- audio transitions are smooth

### Phase 7: Parity Validation And Cutover

Goal: make the Godot build the primary branch for future development.

Tasks:

- perform screen-by-screen parity checks against the MonoGame version
- verify all current content is present and assigned
- compare economic behavior on representative runs
- test save/load repeatedly across sessions
- confirm game over behavior matches expectations
- verify that no mutable runtime state leaks back into content definitions
- create desktop export profiles

Deliverables:

- Godot build is feature-complete for the current game
- MonoGame version is no longer required for day-to-day development

### Phase 8: Post-Parity Expansion

Only after parity is complete:

- contracts/jobs board
- ship upgrades
- expanded port/item/event content
- port condition modifiers
- faction/reputation scaffolding
- broader map and exploration work

## Content Migration Checklist

### Ports

- import all port backgrounds
- import all preview images
- preserve zone assignments and IDs
- store display text outside scene scripts

### Trade Art

- import all port-specific trade backgrounds
- wire them by port ID

### UI Art

- logo
- buttons
- terminal frames
- cockpit background
- info panel
- icons
- game over art

### Fonts

- replace MonoGame spritefonts with imported font files and theme variants
- preserve the current terminal-style identity

### Audio

- menu music
- port/world music
- click SFX
- optional new invalid-action/error SFX during polish

## Risks And Mitigations

### Risk: Recreating MonoGame Architecture In Godot

Mitigation:

- port rules, not plumbing
- use scenes/signals/theme/layout containers instead of custom draw/update systems

### Risk: Logic Drift During Rewrite

Mitigation:

- add characterization tests before changing formulas
- keep the first parity milestone conservative

### Risk: Save Data Fragility

Mitigation:

- use stable IDs
- keep save DTOs separate from scene nodes and resources

### Risk: UI Scope Creep

Mitigation:

- reach functional parity first
- schedule polish after screen flow and gameplay correctness are complete

### Risk: Hardcoded Data Coming Back

Mitigation:

- move ports/items/events to external data before content expansion starts

## Immediate Next Steps

1. Enable Godot C# in this project and create the base app scene plus application coordinator.
2. Define immutable content definitions and separate runtime state models before moving any UI over.
3. Add characterization tests for economy, travel, events, and save/load.
4. Implement narrow services for economy, travel, event resolution, and game-over checks instead of cloning `GameManager`.
5. Create the 6 screen scenes with placeholder layout and bind them to the new application/domain layer.
6. Port the trade loop first, then the travel loop, then audio/polish.

## Recommended First Milestone

The best first milestone is not "port everything visually."

It is:

"A Godot build that can start a run, load the current data set, move between placeholder versions of all six screens, and execute the trading/travel/game-over loop with correct rules."

Once that exists, the rest of the port becomes a UI/content refinement project rather than a risky rewrite.
