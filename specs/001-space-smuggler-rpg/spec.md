# Feature Specification: StarSmuggler Core Game

**Feature Branch**: `001-space-smuggler-rpg`  
**Created**: 2026-04-18  
**Status**: Draft  
**Input**: User description: "Build a game with the base mechanics of dope wars but set in space with deeper gameplay and an immersive story brought to life by the trading mechanics and exploration/travel aspects of the original game"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Build Profit Through Risky Trading (Priority: P1)

As a player starting with a small ship, limited credits, and incomplete market knowledge, I want to travel between spaceports, buy low, sell high, and stay solvent so that I can feel the same accessible pressure-and-reward loop that makes trading games compelling.

**Why this priority**: This is the minimum playable fantasy. If the game does not deliver a satisfying trading loop with meaningful economic decisions, it does not fulfill the core promise of a space-based dope-wars-style game.

**Independent Test**: Can be fully tested by starting a new run, visiting multiple ports, buying and selling goods, paying travel costs, and observing profit, loss, and survival pressure without needing story progression or advanced unlocks.

**Acceptance Scenarios**:

1. **Given** a player begins a new run with starter funds and a cargo limit, **When** they buy goods at one port and sell them at a second port with higher demand, **Then** their credits increase and cargo inventory updates correctly.
2. **Given** a player lacks enough credits, cargo space, or legal access to a trade, **When** they attempt the transaction, **Then** the game prevents the action and explains why.
3. **Given** a player is considering a route, **When** they inspect another destination, **Then** they can understand expected travel cost, market opportunity, and immediate risk well enough to choose their next move.

---

### User Story 2 - Explore A Living Star Map (Priority: P2)

As a player moving through the galaxy, I want travel to expose me to new locations, random encounters, hidden opportunities, and route-specific dangers so that movement feels like exploration rather than a menu shortcut between markets.

**Why this priority**: Space travel is the key differentiator from a plain dope-wars clone. It is where atmosphere, worldbuilding, and tension become part of the trading experience.

**Independent Test**: Can be tested by running multiple journeys between known and newly discovered locations, triggering travel events, resolving branching outcomes, and verifying that exploration changes future trading and navigation choices.

**Acceptance Scenarios**:

1. **Given** a player selects a destination, **When** travel begins, **Then** the trip consumes route cost, advances jump progression, and may trigger an event that changes credits, cargo, ship condition, reputation, or future access.
2. **Given** a player reaches new sectors over time, **When** they explore farther from the starting region, **Then** they encounter different market profiles, hazards, and opportunities than the safer starter routes.
3. **Given** a player receives a travel event offering a choice, **When** they choose a response, **Then** the result changes the run state in a way that is reflected in later travel, trade, or story content.

---

### User Story 3 - Shape The Story Through Smuggling Decisions (Priority: P3)

As a player trying to get rich in a dangerous frontier, I want factions, characters, and narrative arcs to react to the kinds of cargo I move, the places I visit, and the risks I take so that trading becomes the engine of an immersive story instead of a disconnected minigame.

**Why this priority**: The request explicitly asks for deeper gameplay and an immersive story brought to life by trading and exploration. This story layer is what turns the project into more than a mechanical homage.

**Independent Test**: Can be tested by progressing through one narrative arc, making different trading and travel choices, and confirming that story beats, missions, and faction reactions change based on prior decisions.

**Acceptance Scenarios**:

1. **Given** a player repeatedly trades with or against the interests of a faction, **When** they reach a narrative threshold, **Then** new missions, story scenes, or restrictions become available based on that relationship.
2. **Given** a player accepts a story-driven smuggling job, **When** they complete or fail it, **Then** the world state, rewards, and future narrative options reflect that outcome.
3. **Given** a player uncovers major information or reaches a milestone region, **When** the next story beat is triggered, **Then** the story references prior trade and travel behavior rather than feeling generic.

---

### User Story 4 - Grow From Scrappy Courier To Legendary Operator (Priority: P4)

As a player who survives the early game, I want to upgrade my ship, specialize my playstyle, and unlock higher-risk opportunities so that long-term progression keeps the trading loop fresh instead of becoming a grind.

**Why this priority**: Long-term progression supports replayability and gives players a strategic reason to pursue profit beyond a simple score chase.

**Independent Test**: Can be tested by earning enough money and reputation to unlock multiple upgrades, purchasing them, and confirming they meaningfully affect cargo capacity, survivability, travel range, market access, or event outcomes.

**Acceptance Scenarios**:

1. **Given** a player has enough resources, **When** they purchase a ship upgrade, **Then** the relevant run stats and available options change immediately.
2. **Given** a player specializes in stealth, cargo capacity, speed, or influence, **When** they continue playing, **Then** their chosen build creates noticeable advantages and tradeoffs.
3. **Given** a player reaches advanced progression, **When** they access dangerous sectors or premium contracts, **Then** those opportunities offer higher rewards while imposing proportionally greater risk.

---

### Edge Cases

- What happens when a player becomes stranded with cargo but no affordable route or legal market access?
- How does the game handle a travel event that would destroy or confiscate the player’s last remaining source of progression?
- What happens when story progression depends on a location or faction the player has alienated or lost access to?
- How does the economy behave if the player repeatedly farms the same short route or attempts to exploit a single profitable commodity?
- What happens when a player declines or abandons time-sensitive story cargo in transit?
- How does the game distinguish between a recoverable setback and a true game-over state?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST allow players to start a new run with a starter ship, starter funds, starting location, and a clear immediate objective.
- **FR-002**: The system MUST provide multiple trade goods with distinct legality, rarity, price volatility, and regional demand profiles.
- **FR-003**: The system MUST allow players to buy, sell, hold, and transport goods subject to credits, cargo capacity, location access, and other applicable constraints.
- **FR-004**: The system MUST provide a star map with multiple ports or regions that differ in market identity, travel cost, danger level, and narrative relevance.
- **FR-005**: The system MUST require travel between locations to consume the same meaningful progression pressures as the MonoGame game: route-based credit cost, jump-count progression, and event risk, with no new fuel or time-management subsystem introduced in the MVP.
- **FR-006**: The system MUST resolve travel through interactive or outcome-driven events that can affect credits, cargo, ship condition, crew state, faction standing, or story flags.
- **FR-007**: The system MUST communicate current route options, likely costs, and immediate constraints before the player commits to travel.
- **FR-008**: The system MUST include at least one persistent story structure composed of factions, named characters, missions, and milestone events tied directly to trading and exploration behavior.
- **FR-009**: The system MUST allow missions or contracts that require transporting specific goods, meeting deadlines, reaching sectors, or making morally significant decisions.
- **FR-010**: The system MUST track player relationships or standing with at least one reactive world system such as factions, law enforcement, smugglers, or local authorities.
- **FR-011**: The system MUST allow progression through ship upgrades, unlocks, or specialization choices that materially change gameplay.
- **FR-012**: The system MUST define both true game-over states and at least one recoverable setback path for insolvency, stranding, confiscation, or catastrophic travel outcomes.
- **FR-013**: The system MUST preserve game state across sessions, including economy, inventory, ship progression, and story progression.
- **FR-014**: The system MUST present market, travel, and story information in a consistent interface language so the player can understand consequences without relying on hidden rules.
- **FR-015**: The system MUST support replayable runs by varying market conditions, travel risks, mission availability, or story permutations while preserving a coherent world arc.
- **FR-016**: The system MUST prevent dominant low-risk trading exploits from trivializing progression over repeated loops.
- **FR-017**: The system MUST surface short-term goals and long-term aspirations so players always understand both what to do next and why it matters.

### Key Entities *(include if feature involves data)*

- **RunState**: The complete state of one playthrough, including credits, cargo, ship status, discovered locations, jump progression, and story progress.
- **PlayerShip**: The player’s vessel, including cargo capacity, travel range, survivability, equipment, and specialization upgrades.
- **Commodity**: A tradable good with legality, rarity, price behavior, narrative associations, and region-specific desirability.
- **Port**: A market location with local economy traits, services, faction presence, narrative hooks, and route connections.
- **Route**: A travel option between locations with cost, danger, travel time, discovery state, and possible event pools.
- **TravelEvent**: A triggered encounter or branching incident that can alter resources, relationships, or story state.
- **Faction**: A world actor with goals, territory, attitudes toward the player, and associated missions or consequences.
- **Mission**: A structured objective linked to cargo, destinations, time limits, rewards, and narrative outcomes.
- **StoryArc**: A set of related narrative beats driven by faction thresholds, discoveries, player choices, or mission outcomes.
- **MarketSnapshot**: The current local pricing, availability, restrictions, and special conditions for a port at a specific point in a run.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: New players can complete a first profitable trade cycle, including travel, in under 15 minutes without external instructions.
- **SC-002**: Players can clearly identify at least three distinct viable strategies by the midgame, such as safe arbitrage, high-risk smuggling, contract-driven play, or exploration-led play.
- **SC-003**: At least 90% of attempted trade actions fail only for understandable reasons that are communicated in the UI at the moment of decision.
- **SC-004**: A full run consistently produces at least five meaningful decision points where the player must choose between profit, safety, reputation, or story outcome.
- **SC-005**: Replaying the first hour of the game produces materially different combinations of prices, route pressure, events, or mission opportunities while preserving the same overall fantasy.
- **SC-006**: Players can describe the game as both a trading game and a story-driven space smuggling adventure in playtest feedback, indicating that neither layer feels bolted on.

## Assumptions

- The initial release is a single-player game focused on one core campaign or run structure rather than multiplayer or live-service features.
- The player experience is primarily menu-and-scene driven rather than twitch-action combat focused; danger is expressed mainly through strategic travel and encounter decisions.
- The economy is systemic enough to reward route planning but authored enough to preserve pacing, story beats, and readable player choice.
- Combat, if present, is secondary to trading, travel, and narrative consequence in the first major milestone.
- The first playable version targets full functional parity with the MonoGame game before adding narrative, faction, mission, or progression expansion.
- Accessibility, controller support, and save/continue behavior are in scope for the shipped product, but advanced mod support and online features are out of scope for this feature.
