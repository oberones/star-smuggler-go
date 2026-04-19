# Tasks: StarSmuggler Core Game

**Input**: Design documents from `/specs/001-space-smuggler-rpg/`  
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: Tests are required for this feature by the project constitution and plan. Each story includes executable test coverage tasks.

**Organization**: Tasks are grouped by user story so each story can be implemented, tested, and validated independently.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no blocking dependency on incomplete tasks)
- **[Story]**: Which user story this task belongs to (`[US1]`, `[US2]`, `[US3]`, `[US4]`)
- Every task includes an exact file path

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Bootstrap the Go + `go-dot` runtime and establish the new source layout without breaking the existing Godot project shell.

- [X] T001 Create the Go module and initial dependency manifest in `/Users/oberon/Projects/coding/godot/star-smuggler-go/go.mod`
- [X] T002 Create the application entrypoint scaffold in `/Users/oberon/Projects/coding/godot/star-smuggler-go/cmd/starsmuggler/main.go`
- [X] T003 [P] Create the internal package directory skeleton with package docs in `/Users/oberon/Projects/coding/godot/star-smuggler-go/internal/application/.gitkeep`, `/Users/oberon/Projects/coding/godot/star-smuggler-go/internal/content/.gitkeep`, `/Users/oberon/Projects/coding/godot/star-smuggler-go/internal/domain/.gitkeep`, `/Users/oberon/Projects/coding/godot/star-smuggler-go/internal/persistence/.gitkeep`, `/Users/oberon/Projects/coding/godot/star-smuggler-go/internal/presentation/godot/.gitkeep`, and `/Users/oberon/Projects/coding/godot/star-smuggler-go/internal/services/.gitkeep`
- [X] T004 [P] Create the Go test directory skeleton in `/Users/oberon/Projects/coding/godot/star-smuggler-go/tests/golden/.gitkeep`, `/Users/oberon/Projects/coding/godot/star-smuggler-go/tests/integration/.gitkeep`, and `/Users/oberon/Projects/coding/godot/star-smuggler-go/tests/smoke/.gitkeep`
- [X] T005 Document the Go + `go-dot` bootstrap and developer commands in `/Users/oberon/Projects/coding/godot/star-smuggler-go/README.md`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that must exist before any story can be completed.

**⚠️ CRITICAL**: No user story work should be considered complete until this phase is done.

- [X] T006 Implement content loader abstractions for ports, items, and events in `/Users/oberon/Projects/coding/godot/star-smuggler-go/internal/content/loader.go`
- [X] T007 [P] Define stable authored content structs in `/Users/oberon/Projects/coding/godot/star-smuggler-go/internal/domain/content_types.go`
- [X] T008 [P] Define runtime state structs in `/Users/oberon/Projects/coding/godot/star-smuggler-go/internal/domain/run_state.go`
- [X] T009 [P] Define save DTOs and save version constants in `/Users/oberon/Projects/coding/godot/star-smuggler-go/internal/persistence/save_types.go`
- [X] T010 Implement deterministic RNG and service wiring primitives in `/Users/oberon/Projects/coding/godot/star-smuggler-go/internal/services/runtime_context.go`
- [X] T011 Implement JSON content loading and validation for `/Users/oberon/Projects/coding/godot/star-smuggler-go/data/ports/ports.json`, `/Users/oberon/Projects/coding/godot/star-smuggler-go/data/items/items.json`, and `/Users/oberon/Projects/coding/godot/star-smuggler-go/data/events/events.json` in `/Users/oberon/Projects/coding/godot/star-smuggler-go/internal/content/json_repository.go`
- [X] T012 Implement save read/write and versioned hydration in `/Users/oberon/Projects/coding/godot/star-smuggler-go/internal/persistence/save_repository.go`
- [X] T013 Implement the Go application coordinator and route enum in `/Users/oberon/Projects/coding/godot/star-smuggler-go/internal/application/app.go`
- [X] T014 Implement the `go-dot` bridge shell for scene registration and route binding in `/Users/oberon/Projects/coding/godot/star-smuggler-go/internal/presentation/godot/app_bridge.go`
- [X] T014A Create the C# runtime retirement checklist in `/Users/oberon/Projects/coding/godot/star-smuggler-go/specs/001-space-smuggler-rpg/csharp-runtime-retirement.md`
- [X] T014B Audit and disable competing C# runtime authority in `/Users/oberon/Projects/coding/godot/star-smuggler-go/src/application/AppController.cs`, `/Users/oberon/Projects/coding/godot/star-smuggler-go/src/autoload/GameSession.cs`, `/Users/oberon/Projects/coding/godot/star-smuggler-go/src/autoload/DataRepository.cs`, and `/Users/oberon/Projects/coding/godot/star-smuggler-go/src/autoload/SaveService.cs`
- [X] T015 [P] Add golden tests for content loading in `/Users/oberon/Projects/coding/godot/star-smuggler-go/tests/golden/content_repository_test.go`
- [X] T016 [P] Add save round-trip and version tests in `/Users/oberon/Projects/coding/godot/star-smuggler-go/tests/integration/save_repository_test.go`
- [X] T017 [P] Add a headless smoke runner script for Godot boot validation in `/Users/oberon/Projects/coding/godot/star-smuggler-go/tests/smoke/run_headless_smoke.sh`

**Checkpoint**: Foundation ready. User stories can now proceed in priority order and remain independently testable.

---

## Phase 3: User Story 1 - Build Profit Through Risky Trading (Priority: P1) 🎯 MVP Part 1

**Goal**: Deliver the economic half of the MonoGame-parity MVP: new run, port overview, trade, baseline inter-port travel, save/continue, and stranded evaluation.

**Independent Test**: Start a new run, inspect the current port, buy and sell goods across at least two ports, save and continue, and confirm that credits, cargo, prices, and game-over checks behave like the MonoGame reference.

### Tests for User Story 1

- [X] T018 [P] [US1] Add travel cost parity tests against MonoGame rules in `/Users/oberon/Projects/coding/godot/star-smuggler-go/tests/integration/travel_cost_test.go`
- [X] T018A [P] [US1] Add tests locking the MVP travel-pressure model to route credit cost, jump-count progression, and event risk in `/Users/oberon/Projects/coding/godot/star-smuggler-go/tests/integration/travel_pressure_test.go`
- [X] T019 [P] [US1] Add market generation and pricing parity tests in `/Users/oberon/Projects/coding/godot/star-smuggler-go/tests/integration/economy_service_test.go`
- [X] T020 [P] [US1] Add trade transaction and cargo-cap validation tests in `/Users/oberon/Projects/coding/godot/star-smuggler-go/tests/integration/trade_service_test.go`
- [X] T021 [P] [US1] Add run evaluator tests for stranded and recoverable states in `/Users/oberon/Projects/coding/godot/star-smuggler-go/tests/integration/run_evaluator_test.go`

### Implementation for User Story 1

- [X] T022 [P] [US1] Implement the run factory and starter state rules in `/Users/oberon/Projects/coding/godot/star-smuggler-go/internal/domain/run_factory.go`
- [X] T023 [P] [US1] Implement economy rules and market refresh behavior in `/Users/oberon/Projects/coding/godot/star-smuggler-go/internal/services/economy_service.go`
- [X] T024 [P] [US1] Implement trade validation and transaction application in `/Users/oberon/Projects/coding/godot/star-smuggler-go/internal/services/trade_service.go`
- [X] T025 [P] [US1] Implement travel cost, route payment validation, and destination resolution in `/Users/oberon/Projects/coding/godot/star-smuggler-go/internal/services/travel_service.go`
- [X] T026 [P] [US1] Implement run viability and game-over evaluation in `/Users/oberon/Projects/coding/godot/star-smuggler-go/internal/services/run_evaluator.go`
- [X] T026A [P] [US1] Add regression tests for repeated short-route and single-commodity exploit patterns in `/Users/oberon/Projects/coding/godot/star-smuggler-go/tests/integration/economy_balance_test.go`
- [X] T026B [US1] Implement anti-exploit market pressure rules for repeated route farming in `/Users/oberon/Projects/coding/godot/star-smuggler-go/internal/services/economy_balance_service.go`
- [X] T026C [US1] Integrate exploit-pressure adjustments into market refresh and route resolution in `/Users/oberon/Projects/coding/godot/star-smuggler-go/internal/application/run_commands.go`
- [X] T027 [US1] Implement application commands for new run, continue, trade, travel preview, baseline travel commit, and save in `/Users/oberon/Projects/coding/godot/star-smuggler-go/internal/application/run_commands.go`
- [X] T028 [P] [US1] Implement main menu presenter and view model in `/Users/oberon/Projects/coding/godot/star-smuggler-go/internal/presentation/godot/main_menu_presenter.go`
- [X] T029 [P] [US1] Implement port overview presenter and view model in `/Users/oberon/Projects/coding/godot/star-smuggler-go/internal/presentation/godot/port_overview_presenter.go`
- [X] T030 [P] [US1] Implement trade presenter and view model in `/Users/oberon/Projects/coding/godot/star-smuggler-go/internal/presentation/godot/trade_presenter.go`
- [X] T031 [P] [US1] Implement travel presenter and route preview view model in `/Users/oberon/Projects/coding/godot/star-smuggler-go/internal/presentation/godot/travel_presenter.go`
- [X] T032 [US1] Bind the Go presenters to the existing Godot scenes in `/Users/oberon/Projects/coding/godot/star-smuggler-go/internal/presentation/godot/scene_bindings.go`
- [ ] T033 [US1] Preserve the MonoGame-faithful MVP route flow and button actions in `/Users/oberon/Projects/coding/godot/star-smuggler-go/scenes/screens/MainMenuScreen.tscn`, `/Users/oberon/Projects/coding/godot/star-smuggler-go/scenes/screens/PortOverviewScreen.tscn`, `/Users/oberon/Projects/coding/godot/star-smuggler-go/scenes/screens/TradeScreen.tscn`, and `/Users/oberon/Projects/coding/godot/star-smuggler-go/scenes/screens/TravelScreen.tscn`
- [X] T034 [US1] Implement autosave and continue integration for the MVP loop in `/Users/oberon/Projects/coding/godot/star-smuggler-go/internal/application/save_commands.go`
- [X] T034A [P] [US1] Add tests for recoverable setback vs. true game-over transitions in `/Users/oberon/Projects/coding/godot/star-smuggler-go/tests/integration/failure_state_test.go`
- [X] T034B [US1] Implement Game Over route presenter and summary mapping in `/Users/oberon/Projects/coding/godot/star-smuggler-go/internal/presentation/godot/game_over_presenter.go`
- [ ] T034C [US1] Bind the Go runtime to `/Users/oberon/Projects/coding/godot/star-smuggler-go/scenes/screens/GameOverScreen.tscn`
- [X] T034D [US1] Implement at least one explicit recovery mechanic short of game over in `/Users/oberon/Projects/coding/godot/star-smuggler-go/internal/application/recovery_commands.go`

**Checkpoint**: User Story 1 should now provide the core economy loop and baseline inter-port trading, but the MonoGame-parity MVP is not complete until User Story 2 is finished.

---

## Phase 4: User Story 2 - Explore A Living Star Map (Priority: P2) 🎯 MVP Part 2

**Goal**: Complete the MonoGame-parity MVP by porting travel animation timing, random encounter resolution, and the remaining exploration feel from the MonoGame game into the Go runtime.

**Independent Test**: From a working baseline trading run, choose multiple destinations, observe preview data, complete travel animation transitions, and verify random travel events mutate run state and return the player to a valid next screen.

### Tests for User Story 2

- [ ] T035 [P] [US2] Add deterministic travel event selection tests in `/Users/oberon/Projects/coding/godot/star-smuggler-go/tests/integration/event_service_test.go`
- [ ] T036 [P] [US2] Add travel resolution tests covering arrival refresh, jump count, and event application in `/Users/oberon/Projects/coding/godot/star-smuggler-go/tests/integration/travel_resolution_test.go`
- [ ] T037 [P] [US2] Add smoke coverage for travel route transitions in `/Users/oberon/Projects/coding/godot/star-smuggler-go/tests/smoke/travel_flow_check.sh`

### Implementation for User Story 2

- [ ] T038 [P] [US2] Implement typed travel event definitions and resolved event results in `/Users/oberon/Projects/coding/godot/star-smuggler-go/internal/domain/event_types.go`
- [ ] T039 [P] [US2] Implement event rolling and effect resolution in `/Users/oberon/Projects/coding/godot/star-smuggler-go/internal/services/event_service.go`
- [ ] T040 [US2] Extend travel commands with event resolution, jump progression handling, and final arrival refresh behavior in `/Users/oberon/Projects/coding/godot/star-smuggler-go/internal/application/travel_commands.go`
- [ ] T041 [P] [US2] Implement the travel animation presenter and state bridge in `/Users/oberon/Projects/coding/godot/star-smuggler-go/internal/presentation/godot/travel_animation_presenter.go`
- [ ] T042 [US2] Preserve MonoGame-like travel animation pacing, route labels, and skip behavior in `/Users/oberon/Projects/coding/godot/star-smuggler-go/scenes/screens/TravelAnimationScreen.tscn`
- [ ] T043 [US2] Surface recent event outcomes back into the port overview contract in `/Users/oberon/Projects/coding/godot/star-smuggler-go/internal/presentation/godot/port_overview_presenter.go`
- [ ] T044 [US2] Ensure travel visuals keep the current cockpit, preview, and terminal composition in `/Users/oberon/Projects/coding/godot/star-smuggler-go/scenes/screens/TravelScreen.tscn` and `/Users/oberon/Projects/coding/godot/star-smuggler-go/scenes/screens/TravelAnimationScreen.tscn`

**Checkpoint**: User Stories 1 and 2 now form the full MonoGame-parity MVP and should be demoable end to end.

---

## Phase 5: User Story 3 - Shape The Story Through Smuggling Decisions (Priority: P3)

**Goal**: Introduce a first narrative spine driven by trading and exploration, including one reactive faction layer, one mission chain, and stateful story consequences.

**Independent Test**: Accept and resolve one story-driven smuggling job, make at least one branching decision, and verify that faction standing, mission outcomes, and subsequent story text differ based on player choices.

### Tests for User Story 3

- [ ] T045 [P] [US3] Add faction standing transition tests in `/Users/oberon/Projects/coding/godot/star-smuggler-go/tests/integration/faction_service_test.go`
- [ ] T046 [P] [US3] Add mission lifecycle tests in `/Users/oberon/Projects/coding/godot/star-smuggler-go/tests/integration/mission_service_test.go`
- [ ] T047 [P] [US3] Add story arc trigger tests in `/Users/oberon/Projects/coding/godot/star-smuggler-go/tests/integration/story_service_test.go`

### Implementation for User Story 3

- [ ] T048 [P] [US3] Add faction definitions and initial JSON content in `/Users/oberon/Projects/coding/godot/star-smuggler-go/data/factions/sol_factions.json`
- [ ] T049 [P] [US3] Add first mission chain content in `/Users/oberon/Projects/coding/godot/star-smuggler-go/data/missions/intro_smuggling_jobs.json`
- [ ] T050 [P] [US3] Add first story arc content in `/Users/oberon/Projects/coding/godot/star-smuggler-go/data/story/intro_arc.json`
- [ ] T051 [P] [US3] Implement faction runtime state and standing rules in `/Users/oberon/Projects/coding/godot/star-smuggler-go/internal/domain/faction_state.go`
- [ ] T052 [P] [US3] Implement mission runtime state in `/Users/oberon/Projects/coding/godot/star-smuggler-go/internal/domain/mission_state.go`
- [ ] T053 [P] [US3] Implement story state and progression flags in `/Users/oberon/Projects/coding/godot/star-smuggler-go/internal/domain/story_state.go`
- [ ] T054 [P] [US3] Implement faction standing logic in `/Users/oberon/Projects/coding/godot/star-smuggler-go/internal/services/faction_service.go`
- [ ] T055 [P] [US3] Implement mission acceptance, completion, and failure rules in `/Users/oberon/Projects/coding/godot/star-smuggler-go/internal/services/mission_service.go`
- [ ] T056 [P] [US3] Implement story arc progression logic in `/Users/oberon/Projects/coding/godot/star-smuggler-go/internal/services/story_service.go`
- [ ] T057 [US3] Integrate mission and story updates into trade and travel command flows in `/Users/oberon/Projects/coding/godot/star-smuggler-go/internal/application/story_commands.go`
- [ ] T058 [US3] Present faction, mission, and story notices through the existing terminal-style screens in `/Users/oberon/Projects/coding/godot/star-smuggler-go/internal/presentation/godot/story_presenter.go`

**Checkpoint**: User Story 3 should layer reactive narrative play onto the parity loop without needing new screen types.

---

## Phase 6: User Story 4 - Grow From Scrappy Courier To Legendary Operator (Priority: P4)

**Goal**: Add ship upgrades and specialization choices that create meaningful long-term progression without breaking the MonoGame-like accessibility of the core loop.

**Independent Test**: Earn enough resources to buy at least one upgrade, observe the ship state change immediately, and confirm that the upgrade changes subsequent trade, travel, or event outcomes.

### Tests for User Story 4

- [ ] T059 [P] [US4] Add ship upgrade rule tests in `/Users/oberon/Projects/coding/godot/star-smuggler-go/tests/integration/upgrade_service_test.go`
- [ ] T060 [P] [US4] Add specialization effect tests in `/Users/oberon/Projects/coding/godot/star-smuggler-go/tests/integration/specialization_test.go`

### Implementation for User Story 4

- [ ] T061 [P] [US4] Add upgrade definitions in `/Users/oberon/Projects/coding/godot/star-smuggler-go/data/upgrades/ship_upgrades.json`
- [ ] T062 [P] [US4] Implement ship upgrade and specialization state in `/Users/oberon/Projects/coding/godot/star-smuggler-go/internal/domain/ship_progression.go`
- [ ] T063 [P] [US4] Implement upgrade purchase and application rules in `/Users/oberon/Projects/coding/godot/star-smuggler-go/internal/services/upgrade_service.go`
- [ ] T064 [US4] Integrate upgrade availability and purchase flow into the application layer in `/Users/oberon/Projects/coding/godot/star-smuggler-go/internal/application/progression_commands.go`
- [ ] T065 [US4] Present ship progression state in the existing port overview and trade flow in `/Users/oberon/Projects/coding/godot/star-smuggler-go/internal/presentation/godot/progression_presenter.go`

**Checkpoint**: All user stories should now be independently functional and cumulatively form the full target experience.

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Final improvements that affect multiple stories and bring the experience closer to the old game’s presentation quality.

- [ ] T066 [P] Run a full parity audit against the MonoGame reference and record any intentional deviations in `/Users/oberon/Projects/coding/godot/star-smuggler-go/specs/001-space-smuggler-rpg/research.md`
- [ ] T067 Optimize texture and audio resource caching across presenters in `/Users/oberon/Projects/coding/godot/star-smuggler-go/internal/presentation/godot/resource_cache.go`
- [ ] T068 Improve music transition and click SFX behavior to match the MonoGame feel in `/Users/oberon/Projects/coding/godot/star-smuggler-go/internal/presentation/godot/audio_bridge.go`
- [ ] T069 [P] Add final smoke coverage for the complete MVP loop in `/Users/oberon/Projects/coding/godot/star-smuggler-go/tests/smoke/full_loop_check.sh`
- [ ] T070 Validate the documented developer flow and commands in `/Users/oberon/Projects/coding/godot/star-smuggler-go/specs/001-space-smuggler-rpg/quickstart.md`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1: Setup**: No dependencies, can start immediately
- **Phase 2: Foundational**: Depends on Setup completion and blocks all user story work
- **Phase 3: User Story 1**: Depends on Foundational completion and delivers the economy and baseline inter-port half of the MVP
- **Phase 4: User Story 2**: Depends on User Story 1 runtime foundation and, together with User Story 1, defines the full MVP
- **Phase 5: User Story 3**: Depends on User Stories 1 and 2 because missions and story consequences must attach to trade and travel outcomes
- **Phase 6: User Story 4**: Depends on MVP completion and should follow User Story 3 if upgrades affect mission or faction outcomes
- **Phase 7: Polish**: Depends on all desired stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational with no dependency on other stories
- **User Story 2 (P2)**: Depends on User Story 1’s route, market, save, and baseline travel infrastructure
- **User Story 3 (P3)**: Depends on User Story 1 and 2 state transitions
- **User Story 4 (P4)**: Depends on User Story 1 and can be integrated after core story/travel behavior is stable

### Within Each User Story

- Tests should be written before implementation for the corresponding slice
- Domain types and service rules come before presenter wiring
- Presenter wiring comes before scene-specific polish
- Story checkpoints should be validated before moving to the next priority

### Parallel Opportunities

- Setup tasks marked `[P]` can run in parallel
- Foundational tasks T007-T010 and tests T015-T017 can run in parallel once the module exists
- In US1, service implementations T022-T026 and presenter scaffolds T028-T031 can be split across contributors
- In US3, authored content tasks T048-T050 and runtime service tasks T051-T056 are good parallel lanes
- In US4, upgrade content and runtime rules can be built in parallel before integration

---

## Parallel Example: User Story 1

```bash
# Tests first for US1
Task: "T018 Add travel cost parity tests in tests/integration/travel_cost_test.go"
Task: "T019 Add market generation and pricing parity tests in tests/integration/economy_service_test.go"
Task: "T020 Add trade transaction and cargo-cap validation tests in tests/integration/trade_service_test.go"
Task: "T021 Add run evaluator tests in tests/integration/run_evaluator_test.go"

# Core runtime implementation for US1
Task: "T022 Implement the run factory in internal/domain/run_factory.go"
Task: "T023 Implement economy rules in internal/services/economy_service.go"
Task: "T024 Implement trade rules in internal/services/trade_service.go"
Task: "T025 Implement travel cost rules in internal/services/travel_service.go"
Task: "T026 Implement run evaluation in internal/services/run_evaluator.go"
```

---

## Implementation Strategy

### MVP First (User Stories 1 + 2)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational
3. Complete Phase 3: User Story 1
4. Complete Phase 4: User Story 2
5. **STOP and VALIDATE**: Compare the Go MVP against the complete MonoGame six-screen loop
6. Demo the MVP before adding story or progression work

### Incremental Delivery

1. Setup + Foundational -> pure Go runtime and Godot bridge ready
2. User Story 1 -> economy loop and baseline inter-port trading
3. User Story 2 -> full MonoGame-parity MVP
4. User Story 3 -> immersive narrative layer
5. User Story 4 -> long-term progression
6. Polish -> parity refinement and performance cleanup

### Parallel Team Strategy

With multiple contributors:

1. One contributor owns Go domain/services
2. One contributor owns `go-dot` presenters and Godot scene bindings
3. One contributor owns content JSON, narrative definitions, and parity checks
4. Shared validation happens at each story checkpoint

---

## Notes

- `[P]` means the task touches distinct files and should not require unfinished work in the same phase
- Every story is traced back to the feature spec and can be validated independently
- User Stories 1 and 2 together are the recommended MVP scope
- Preserve the MonoGame look and feel by default; any deviation should be deliberate and documented
