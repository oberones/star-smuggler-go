# Contract: MVP Screen Flow

## Purpose

Define the runtime contract between the Godot scenes and the Go application layer for the MVP screen loop. This contract exists so that screens remain passive and gameplay transitions are owned by the Go runtime.

## Screen Set

- Main Menu
- Port Overview
- Trade
- Travel
- Travel Animation
- Game Over

## Route Contract

### Main Menu

- **Inputs from runtime**:
  - `can_continue`
  - `status_message`
- **Outputs to runtime**:
  - `start_requested`
  - `continue_requested`
  - `quit_requested`

### Port Overview

- **Inputs from runtime**:
  - current port name
  - port description
  - zone label
  - background asset path
  - current credits
  - cargo load and limit
  - recent event text
  - available goods summary
  - route status summary
- **Outputs to runtime**:
  - `open_trade_requested`
  - `open_travel_requested`
  - `back_to_menu_requested`

### Trade

- **Inputs from runtime**:
  - current port name
  - trade background asset path
  - item list with prices and owned quantities
  - credits and cargo summary
  - status/error message
- **Outputs to runtime**:
  - `buy_requested(item_id, quantity)`
  - `sell_requested(item_id, quantity)`
  - `back_to_port_requested`

### Travel

- **Inputs from runtime**:
  - current port name
  - background asset path
  - destination list
  - selected destination preview
  - travel costs and route summary
  - status/error message
- **Outputs to runtime**:
  - `travel_requested(destination_port_id)`
  - `back_to_port_requested`
  - `destination_selected(destination_port_id)` optional

### Travel Animation

- **Inputs from runtime**:
  - origin name
  - destination name
  - travel cost
  - duration
  - background/animation asset references
  - status text
- **Outputs to runtime**:
  - `animation_finished`
  - `skip_requested`

### Game Over

- **Inputs from runtime**:
  - summary text
  - background asset path
- **Outputs to runtime**:
  - `restart_requested`
  - `return_to_menu_requested`

## Transition Rules

- Screens do not mutate gameplay state directly.
- A screen may only emit intent signals.
- The Go application coordinator resolves each emitted intent into:
  - state mutation
  - save action
  - route change
  - audio update
  - validation error message

## Forbidden Behavior

- A screen must not:
  - calculate trade validity
  - resolve travel costs
  - roll random events
  - write saves
  - decide game-over state
  - own story progression logic

## MVP Parity Expectations

- The route order and available actions must remain consistent with the MonoGame MVP.
- Button labels and decision language should remain familiar unless the spec explicitly changes them.
- The terminal-style composition should remain visually recognizable even if Godot containers replace fixed coordinates internally.
