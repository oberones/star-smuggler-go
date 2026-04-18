# StarSmugglerGo Constitution

## Core Principles

### I. Idiomatic Go, Thin Engine Boundaries
All gameplay rules, economy logic, event resolution, save transformations, and view-model shaping MUST live in idiomatic Go packages with clear ownership and explicit APIs. `go-dot` integration code MUST remain thin and be limited to scene wiring, signal bridging, rendering concerns, and engine lifecycle hooks. Engine-facing code MUST depend on domain code; domain code MUST NOT depend on Godot types, nodes, scene trees, or frame callbacks. Shared logic MUST favor plain structs, interfaces only when they remove real coupling, and allocation-conscious code over framework-style indirection.

### II. Testable By Default
Every non-trivial behavior change MUST arrive with automated tests at the same layer as the change. Pure domain logic MUST have table-driven unit tests. Serialization, save compatibility, content loading, and Godot boundary adapters MUST have focused integration or golden tests where appropriate. Uncovered logic is acceptable only for trivial glue code, and that exemption MUST be explained in review notes. A feature is not complete when it only works in-editor; it is complete when its behavior is reproducible in tests.

### III. Consistent Player Experience
The port MUST preserve a coherent StarSmuggler experience across screens, input methods, and resolutions. Navigation language, button placement, status messaging, and visual hierarchy MUST remain consistent unless a spec explicitly defines a deliberate redesign. UI state MUST never depend on incidental scene timing such as `_Ready()` ordering or hidden side effects. Screens SHOULD be passive presenters of explicit state. If behavior spans scenes, the transition contract MUST be owned by application code rather than duplicated in UI scripts.

### IV. Performance Budgets Are Product Requirements
Performance is a feature, not cleanup work. Core loops such as trade refresh, travel resolution, event rolls, content lookup, and save/load MUST avoid per-frame reflection, unnecessary allocations, and repeated resource loading. Screen updates SHOULD be event-driven and diff-aware instead of rebuilt every frame. Resource-heavy assets such as textures and audio SHOULD be loaded predictably and reused where practical. Any new feature that risks stutter, input delay, or scene-transition hitching MUST define its mitigation strategy before implementation.

### V. Data Integrity And Determinism
Authored content, runtime state, and presentation state MUST remain distinct. Immutable game content MUST have stable identifiers and versionable formats. Runtime mutations MUST be explicit, traceable, and safe to serialize. Randomized systems MUST support deterministic testing through injectable RNG or controlled seeds. Save migrations MUST preserve player progress whenever practical and MUST fail loudly rather than silently corrupting state.

## Engineering Standards

- Go packages MUST follow single-purpose boundaries such as `domain`, `application`, `services`, `content`, `persistence`, and `presentation` where those boundaries reduce coupling.
- Public APIs MUST prefer clear nouns and verbs over engine-inspired naming. Avoid vague managers, god objects, and grab-bag utility packages.
- Errors MUST be returned with context and handled intentionally. Panics are reserved for unrecoverable programmer errors or impossible bootstrap states.
- Concurrency MUST be opt-in and justified. Do not introduce goroutines, channels, or locks into gameplay code unless they solve a measured problem without making determinism or debugging materially worse.
- `go-dot` scene scripts MUST not own business rules, persistence decisions, or cross-screen navigation policy.
- New dependencies MUST be justified by capability that the standard library or existing project code cannot reasonably provide.

## Quality Gates And Workflow

- Every spec and implementation plan MUST explain where logic lives: pure Go package, `go-dot` adapter, or scene-only presentation.
- Before merge, contributors MUST verify:
  - `go test ./...` passes
  - new and changed behavior has matching tests
  - save/load implications are documented for state or schema changes
  - UX changes are reviewed for consistency with adjacent screens
  - performance-sensitive paths are checked for avoidable allocations or redundant loads
- Reviews MUST reject:
  - gameplay logic embedded in scene scripts
  - untested domain behavior
  - UI changes that create inconsistent navigation or feedback
  - broad abstractions without a concrete second use case
  - silent breaking changes to content IDs, saves, or input behavior

## Governance

This constitution supersedes local habits, ad hoc engine patterns, and convenience-driven shortcuts. All future plans, specs, and task lists in this repository MUST align with these principles. Amendments require:

1. a written rationale,
2. an explicit description of what existing rule is changing,
3. a migration plan for affected code or content, and
4. approval in the project record before the amendment is considered active.

Compliance is verified during planning, implementation, and review. If a proposal conflicts with this constitution, the proposal MUST be updated before work proceeds.

**Version**: 1.0.0 | **Ratified**: 2026-04-18 | **Last Amended**: 2026-04-18
