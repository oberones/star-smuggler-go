# Contract: Content And Save Stability

## Purpose

Define the authored-content and save-data stability rules for the MVP and early story expansion.

## Authored Content Contract

### General Rules

- Every authored entity uses a stable string ID.
- Display names are never used as primary references in saves or runtime relationships.
- Asset paths are data, not hardcoded scene assumptions.
- JSON files are the source of truth for authored ports, items, events, and later factions/missions/story arcs.

### Required Content Files

- `data/ports/*.json`
- `data/items/*.json`
- `data/events/*.json`
- future:
  - `data/factions/*.json`
  - `data/missions/*.json`
  - `data/story/*.json`

### Minimum Port Schema

- `id`
- `name`
- `description`
- `zone`
- `backgroundTexturePath`
- `previewTexturePath`
- `tradeBackgroundPath`
- `musicTrackId`

### Minimum Item Schema

- `id`
- `name`
- `description`
- `rarity`
- `basePrice`
- `legalClass` optional for MVP but reserved now

### Minimum Event Schema

- `id`
- `name`
- `baseDescription`
- `effectType`
- `effectParameters`
- `weight`

## Save Contract

### General Rules

- Saves must include an explicit `saveVersion`.
- Saves serialize runtime state only.
- Saves must not contain:
  - Godot node references
  - scene paths as state authority
  - object identity references

### Minimum MVP Save Fields

- `saveVersion`
- `player`
  - `credits`
  - `cargoLimit`
  - `currentPortId`
  - `upgradeIds`
- `cargo`
  - map of `commodityId -> quantity`
- `markets`
  - per-port market snapshots
- `recentEventResult`
- `jumpsSinceRefresh`
- `storyFlags`
- `factionStandings` reserved if the first story spine ships in MVP

## Compatibility Rules

- Existing IDs must not be renamed without a migration plan.
- Removing content referenced by a save requires:
  - migration logic, or
  - a documented compatibility break
- New optional fields may be added if defaults are defined during load.
- Any schema change must be covered by golden tests.

## Failure Behavior

- Unknown save versions must fail loudly with a migration-needed error.
- Invalid content references must surface a clear content-integrity failure during load or startup.
- Silent fallback to display-name matching is forbidden.
