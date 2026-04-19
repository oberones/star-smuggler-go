# Data Model: StarSmuggler Core Game

## Overview

The data model separates authored content from mutable run state. Authored content is loaded from JSON using stable IDs. Mutable runtime state lives in pure Go domain structs and is the only layer that can change during a run. Presentation state is derived from runtime state and never becomes the source of truth.

## Authored Content Entities

### PortDefinition

- **Purpose**: Defines a starport the player can visit.
- **Key fields**:
  - `id`
  - `name`
  - `description`
  - `zone`
  - `background_texture_path`
  - `preview_texture_path`
  - `trade_background_path`
  - `music_track_id`
  - `neighbors` or route references
  - `faction_id` optional
  - `story_tags`
- **Validation**:
  - `id` must be unique and stable
  - visual/audio paths must resolve to imported project assets
  - `zone` must be one of the defined region tiers

### CommodityDefinition

- **Purpose**: Defines a tradable item.
- **Key fields**:
  - `id`
  - `name`
  - `description`
  - `rarity`
  - `base_price`
  - `legal_class`
  - `story_tags`
- **Validation**:
  - `base_price` must be positive
  - rarity must map to balancing rules
  - IDs must remain stable across save versions

### TravelEventDefinition

- **Purpose**: Defines a possible travel encounter.
- **Key fields**:
  - `id`
  - `name`
  - `base_description`
  - `effect_type`
  - `effect_parameters`
  - `eligibility_rules`
  - `weight`
  - `choice_set` optional
  - `story_tags`
- **Validation**:
  - effect parameters must be valid for the effect type
  - weighted events must have non-negative weights
  - events referenced by story or route data must exist

### FactionDefinition

- **Purpose**: Defines a world actor that reacts to player behavior.
- **Key fields**:
  - `id`
  - `name`
  - `description`
  - `alignment`
  - `home_ports`
  - `rival_faction_ids`
  - `standing_thresholds`
- **Validation**:
  - thresholds must be ordered and non-overlapping
  - referenced ports and factions must exist

### MissionDefinition

- **Purpose**: Defines authored trade/story missions.
- **Key fields**:
  - `id`
  - `name`
  - `briefing`
  - `mission_type`
  - `origin_port_id`
  - `destination_port_id`
  - `required_commodity_id` optional
  - `required_quantity` optional
  - `deadline_rule`
  - `reward_definition`
  - `failure_consequences`
  - `unlock_conditions`
- **Validation**:
  - referenced content IDs must exist
  - quantity must be positive when required
  - deadlines must be compatible with the game’s time model

### StoryArcDefinition

- **Purpose**: Defines milestone-driven narrative sequences.
- **Key fields**:
  - `id`
  - `name`
  - `entry_conditions`
  - `beats`
  - `completion_effects`
- **Validation**:
  - beats must be ordered
  - entry conditions must be expressible from runtime state

## Runtime State Entities

### RunState

- **Purpose**: Top-level mutable state for one run.
- **Key fields**:
  - `run_id`
  - `save_version`
  - `current_route`
  - `player_ship`
  - `cargo`
  - `markets_by_port_id`
  - `discovered_port_ids`
  - `recent_event_result`
  - `story_state` post-MVP expansion
  - `faction_standings` post-MVP expansion
  - `active_missions` post-MVP expansion
  - `completed_mission_ids` post-MVP expansion
  - `rng_seed`
  - `jumps_since_refresh`
  - `total_jumps`
- **State transitions**:
  - `new_run` -> starter state
  - `trade_applied` -> credits/cargo/market mutation
  - `travel_requested` -> pending route state
  - `travel_resolved` -> new port, event outcome, jump progression
  - `mission_updated` -> active/completed/failed mission changes
  - `game_over` -> final locked result state

### PlayerShipState

- **Purpose**: Mutable player vessel stats.
- **Key fields**:
  - `credits`
  - `cargo_limit`
  - `current_port_id`
  - `ship_hull_state`
  - `upgrade_ids` post-MVP expansion
  - `specialization_flags` post-MVP expansion
- **Validation**:
  - credits cannot drop below zero unless explicitly modeled as debt
  - cargo limit must be positive
  - current port must exist
  - the parity MVP models travel pressure through credits, jump progression, and event risk rather than a separate fuel meter

### CargoState

- **Purpose**: Player inventory keyed by commodity ID.
- **Key fields**:
  - `quantities_by_commodity_id`
- **Validation**:
  - quantity values must be non-negative integers
  - total cargo cannot exceed ship cargo limit

### MarketSnapshot

- **Purpose**: Current local trading state for a port.
- **Key fields**:
  - `port_id`
  - `available_commodity_ids`
  - `prices_by_commodity_id`
  - `local_modifiers`
  - `last_refresh_tick`
- **Validation**:
  - every price must be positive
  - every commodity must exist in authored content

### RouteState

- **Purpose**: Represents a pending or recent travel route.
- **Key fields**:
  - `origin_port_id`
  - `destination_port_id`
  - `travel_cost`
  - `travel_duration`
  - `route_risk`
  - `status`
- **States**:
  - `previewed`
  - `committed`
  - `animating`
  - `resolved`

### EventResult

- **Purpose**: Stores the concrete outcome of a resolved event.
- **Key fields**:
  - `event_definition_id`
  - `resolved_text`
  - `credits_delta`
  - `cargo_delta`
  - `ship_delta`
  - `standing_delta`
  - `story_flag_changes`
  - `choice_id` optional
- **Validation**:
  - result must be serializable without referencing scene/node objects

### FactionStanding

- **Purpose**: Tracks player reputation with a faction.
- **Key fields**:
  - `faction_id`
  - `score`
  - `standing_tier`
  - `last_change_reason`
- **Validation**:
  - faction must exist
  - tier must be derivable from score and thresholds

### MissionState

- **Purpose**: Stores the mutable runtime state of an accepted mission.
- **Key fields**:
  - `mission_definition_id`
  - `status`
  - `accepted_at_tick`
  - `deadline_tick`
  - `progress_flags`
  - `reward_claimed`
- **States**:
  - `available`
  - `accepted`
  - `in_progress`
  - `completed`
  - `failed`
  - `expired`

### StoryState

- **Purpose**: Tracks narrative progression at run level.
- **Key fields**:
  - `active_story_arc_ids`
  - `completed_story_arc_ids`
  - `story_flags`
  - `named_character_states`
- **Validation**:
  - flags must use stable keys
  - story arcs must reference valid authored arcs

## Persistence Model

### SaveGame

- **Purpose**: Versioned serialized representation of `RunState`.
- **Key fields**:
  - `save_version`
  - `created_at`
  - `updated_at`
  - `run_state`
- **Rules**:
  - version is required
  - migrations operate on DTOs before domain hydration
  - unknown versions fail loudly with a migration-required error

## Derived Presentation Models

These are not persisted and should be derived from runtime state:

- `MainMenuViewModel`
- `PortOverviewViewModel`
- `TradeScreenViewModel`
- `TravelScreenViewModel`
- `TravelAnimationViewModel`
- `GameOverViewModel`

They may contain:

- localized strings
- resolved asset paths
- formatted prices/cargo summaries
- button enabled/disabled state
- recent event/status text

They MUST NOT contain business logic or become writeable run state.

## Core Relationships

- `RunState.player_ship.current_port_id` -> `PortDefinition.id`
- `CargoState.quantities_by_commodity_id.*` -> `CommodityDefinition.id`
- `MarketSnapshot.port_id` -> `PortDefinition.id`
- `MarketSnapshot.prices_by_commodity_id.*` -> `CommodityDefinition.id`
- `EventResult.event_definition_id` -> `TravelEventDefinition.id`
- `FactionStanding.faction_id` -> `FactionDefinition.id`
- `MissionState.mission_definition_id` -> `MissionDefinition.id`
- `StoryState.active_story_arc_ids.*` -> `StoryArcDefinition.id`

## Invariants

- Authored content is immutable at runtime.
- Runtime state never stores direct Godot node references.
- All gameplay-affecting references use stable IDs, not display names or object identity.
- Saves can be reconstructed without scene files.
- Any random outcome must be replayable in tests through controlled RNG input.
