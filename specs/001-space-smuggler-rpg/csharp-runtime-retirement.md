# C# Runtime Retirement Checklist

## Purpose

Track the current C# runtime ownership points that must be retired or downgraded to passive reference status as the Go runtime becomes authoritative.

## Current Runtime Owners

- `src/application/AppController.cs`
- `src/autoload/GameSession.cs`
- `src/autoload/DataRepository.cs`
- `src/autoload/SaveService.cs`

## Retirement Steps

- [ ] Freeze the current behavior of each file as migration reference material only.
- [ ] Identify which Go package takes over each responsibility:
  - app orchestration
  - content loading
  - save persistence
  - run/session ownership
- [ ] Remove or disable gameplay mutations from the C# runtime path once the Go equivalent is active.
- [ ] Keep scene scripts and autoloads limited to presentation or engine-bridge concerns during the transition.
- [ ] Confirm Godot scene flow still works when gameplay authority comes from Go.

## Current Guardrail

- The C# runtime-owner files now check the `STARSMUGGLER_GO_RUNTIME` environment variable.
- When the variable is set to `1` or `true`, these nodes become passive and stop performing gameplay mutations or save/content authority work.
- This keeps the current C# path usable by default while giving the Go runtime a clean takeover switch during migration.

## Notes

- Until the Go bridge is active, these files remain useful for parity checks.
- The end state is a thin Godot-facing integration layer with pure Go gameplay logic.
