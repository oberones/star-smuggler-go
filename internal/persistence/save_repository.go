package persistence

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/oberones/star-smuggler-go/internal/domain"
)

type FileStore interface {
	ReadFile(name string) ([]byte, error)
	WriteFile(name string, data []byte, perm os.FileMode) error
	Stat(name string) (os.FileInfo, error)
	MkdirAll(path string, perm os.FileMode) error
}

type OSFileStore struct{}

func (OSFileStore) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

func (OSFileStore) WriteFile(name string, data []byte, perm os.FileMode) error {
	return os.WriteFile(name, data, perm)
}

func (OSFileStore) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

func (OSFileStore) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

type Repository interface {
	Exists() (bool, error)
	Load(ctx context.Context) (domain.RunState, error)
	Save(ctx context.Context, run domain.RunState) error
}

type JSONSaveRepository struct {
	store FileStore
	path  string
}

func NewJSONSaveRepository(store FileStore, path string) *JSONSaveRepository {
	if store == nil {
		store = OSFileStore{}
	}

	return &JSONSaveRepository{
		store: store,
		path:  path,
	}
}

func (r *JSONSaveRepository) Exists() (bool, error) {
	if r.path == "" {
		return false, errors.New("save path is not configured")
	}

	_, err := r.store.Stat(r.path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func (r *JSONSaveRepository) Load(ctx context.Context) (domain.RunState, error) {
	if err := ctx.Err(); err != nil {
		return domain.RunState{}, err
	}

	bytes, err := r.store.ReadFile(r.path)
	if err != nil {
		return domain.RunState{}, err
	}

	var save SaveData
	if err := json.Unmarshal(bytes, &save); err != nil {
		return domain.RunState{}, fmt.Errorf("decode save file: %w", err)
	}

	if save.Version != CurrentSaveVersion {
		return domain.RunState{}, fmt.Errorf("unsupported save version %d", save.Version)
	}

	return hydrateRunState(save), nil
}

func (r *JSONSaveRepository) Save(ctx context.Context, run domain.RunState) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if r.path == "" {
		return errors.New("save path is not configured")
	}

	save := dehydrateRunState(run)
	bytes, err := json.MarshalIndent(save, "", "  ")
	if err != nil {
		return fmt.Errorf("encode save file: %w", err)
	}

	if dir := filepath.Dir(r.path); dir != "." && dir != "" {
		if err := r.store.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("create save directory: %w", err)
		}
	}

	if err := r.store.WriteFile(r.path, bytes, 0o644); err != nil {
		return fmt.Errorf("write save file: %w", err)
	}

	return nil
}

func dehydrateRunState(run domain.RunState) SaveData {
	cargo := make(map[string]int, len(run.Cargo.ItemQuantities))
	for itemID, quantity := range run.Cargo.ItemQuantities {
		cargo[itemID] = quantity
	}

	marketIDs := make([]string, 0, len(run.MarketsByPortID))
	for portID := range run.MarketsByPortID {
		marketIDs = append(marketIDs, portID)
	}
	sort.Strings(marketIDs)

	markets := make([]MarketSnapshotSaveData, 0, len(run.MarketsByPortID))
	for _, portID := range marketIDs {
		market := run.MarketsByPortID[portID]
		markets = append(markets, MarketSnapshotSaveData{
			PortID:           market.PortID,
			AvailableItemIDs: append([]string(nil), market.AvailableItemIDs...),
			PricesByItemID:   cloneIntMap(market.PricesByItemID),
		})
	}

	factionIDs := make([]string, 0, len(run.FactionStandings))
	for factionID := range run.FactionStandings {
		factionIDs = append(factionIDs, factionID)
	}
	sort.Strings(factionIDs)

	factionStandings := make([]FactionStandingSaveData, 0, len(run.FactionStandings))
	for _, factionID := range factionIDs {
		standing := run.FactionStandings[factionID]
		factionStandings = append(factionStandings, FactionStandingSaveData{
			FactionID:        standing.FactionID,
			Score:            standing.Score,
			StandingTier:     standing.StandingTier,
			LastChangeReason: standing.LastChangeReason,
		})
	}

	missionIDs := make([]string, 0, len(run.ActiveMissions))
	for missionID := range run.ActiveMissions {
		missionIDs = append(missionIDs, missionID)
	}
	sort.Strings(missionIDs)

	activeMissions := make([]MissionStateSaveData, 0, len(run.ActiveMissions))
	for _, missionID := range missionIDs {
		mission := run.ActiveMissions[missionID]
		activeMissions = append(activeMissions, MissionStateSaveData{
			MissionDefinitionID: mission.MissionDefinitionID,
			Status:              string(mission.Status),
			AcceptedAtJump:      mission.AcceptedAtJump,
			DeadlineJump:        mission.DeadlineJump,
			ProgressFlags:       cloneBoolMap(mission.ProgressFlags),
			RewardClaimed:       mission.RewardClaimed,
		})
	}

	return SaveData{
		Version: CurrentSaveVersion,
		Player: PlayerSaveData{
			Credits:       run.Player.Credits,
			CargoLimit:    run.Player.CargoLimit,
			CurrentPortID: run.Player.CurrentPortID,
		},
		CargoByItemID:             cargo,
		Markets:                   markets,
		RoutePressureByKey:        cloneIntMap(run.RoutePressureByKey),
		CommodityPressureByItemID: cloneIntMap(run.CommodityPressureByItemID),
		Progression: ShipProgressionSaveData{
			PurchasedUpgradeIDs: append([]string(nil), run.Progression.PurchasedUpgradeIDs...),
			SpecializationFlags: cloneBoolMap(run.Progression.SpecializationFlags),
		},
		FactionStandings:    factionStandings,
		ActiveMissions:      activeMissions,
		CompletedMissionIDs: append([]string(nil), run.CompletedMissionIDs...),
		Story: StoryStateSaveData{
			ActiveStoryArcIDs:    append([]string(nil), run.Story.ActiveStoryArcIDs...),
			CompletedStoryArcIDs: append([]string(nil), run.Story.CompletedStoryArcIDs...),
			StoryFlags:           cloneBoolMap(run.Story.StoryFlags),
			NamedCharacterStates: cloneStringMap(run.Story.NamedCharacterStates),
		},
		EmergencyRecoveryUsed: run.EmergencyRecoveryUsed,
		JumpsSinceLastUpdate:  run.JumpsSinceLastUpdate,
		TotalJumps:            run.TotalJumps,
		RecentEvent:           dehydrateEvent(run.RecentEvent),
	}
}

func hydrateRunState(save SaveData) domain.RunState {
	run := domain.NewRunState()
	run.Player.Credits = save.Player.Credits
	run.Player.CargoLimit = save.Player.CargoLimit
	run.Player.CurrentPortID = save.Player.CurrentPortID
	run.Cargo.ItemQuantities = cloneIntMap(save.CargoByItemID)
	run.MarketsByPortID = make(map[string]domain.MarketSnapshot, len(save.Markets))
	run.RoutePressureByKey = cloneIntMap(save.RoutePressureByKey)
	run.CommodityPressureByItemID = cloneIntMap(save.CommodityPressureByItemID)
	run.Progression = domain.ShipProgressionState{
		PurchasedUpgradeIDs: append([]string(nil), save.Progression.PurchasedUpgradeIDs...),
		SpecializationFlags: cloneBoolMap(save.Progression.SpecializationFlags),
	}
	run.FactionStandings = make(map[string]domain.FactionStanding, len(save.FactionStandings))
	run.ActiveMissions = make(map[string]domain.MissionState, len(save.ActiveMissions))
	run.CompletedMissionIDs = append([]string(nil), save.CompletedMissionIDs...)
	run.Story = domain.StoryState{
		ActiveStoryArcIDs:    append([]string(nil), save.Story.ActiveStoryArcIDs...),
		CompletedStoryArcIDs: append([]string(nil), save.Story.CompletedStoryArcIDs...),
		StoryFlags:           cloneBoolMap(save.Story.StoryFlags),
		NamedCharacterStates: cloneStringMap(save.Story.NamedCharacterStates),
	}
	run.EmergencyRecoveryUsed = save.EmergencyRecoveryUsed
	run.JumpsSinceLastUpdate = save.JumpsSinceLastUpdate
	run.TotalJumps = save.TotalJumps
	run.RecentEvent = hydrateEvent(save.RecentEvent)

	for _, market := range save.Markets {
		run.MarketsByPortID[market.PortID] = domain.MarketSnapshot{
			PortID:           market.PortID,
			AvailableItemIDs: append([]string(nil), market.AvailableItemIDs...),
			PricesByItemID:   cloneIntMap(market.PricesByItemID),
		}
	}

	for _, standing := range save.FactionStandings {
		run.FactionStandings[standing.FactionID] = domain.FactionStanding{
			FactionID:        standing.FactionID,
			Score:            standing.Score,
			StandingTier:     standing.StandingTier,
			LastChangeReason: standing.LastChangeReason,
		}
	}

	for _, mission := range save.ActiveMissions {
		run.ActiveMissions[mission.MissionDefinitionID] = domain.MissionState{
			MissionDefinitionID: mission.MissionDefinitionID,
			Status:              domain.MissionStatus(mission.Status),
			AcceptedAtJump:      mission.AcceptedAtJump,
			DeadlineJump:        mission.DeadlineJump,
			ProgressFlags:       cloneBoolMap(mission.ProgressFlags),
			RewardClaimed:       mission.RewardClaimed,
		}
	}

	return run
}

func dehydrateEvent(event *domain.EventResult) *EventResultSaveData {
	if event == nil {
		return nil
	}

	return &EventResultSaveData{
		EventID:             event.EventID,
		Name:                event.Name,
		ResolvedDescription: event.ResolvedDescription,
		RolledValues:        cloneFloatMap(event.RolledValues),
	}
}

func hydrateEvent(event *EventResultSaveData) *domain.EventResult {
	if event == nil {
		return nil
	}

	return &domain.EventResult{
		EventID:             event.EventID,
		Name:                event.Name,
		ResolvedDescription: event.ResolvedDescription,
		RolledValues:        cloneFloatMap(event.RolledValues),
	}
}

func cloneIntMap(source map[string]int) map[string]int {
	if source == nil {
		return map[string]int{}
	}

	clone := make(map[string]int, len(source))
	for key, value := range source {
		clone[key] = value
	}
	return clone
}

func cloneFloatMap(source map[string]float64) map[string]float64 {
	if source == nil {
		return map[string]float64{}
	}

	clone := make(map[string]float64, len(source))
	for key, value := range source {
		clone[key] = value
	}
	return clone
}

func cloneBoolMap(source map[string]bool) map[string]bool {
	if source == nil {
		return map[string]bool{}
	}

	clone := make(map[string]bool, len(source))
	for key, value := range source {
		clone[key] = value
	}
	return clone
}

func cloneStringMap(source map[string]string) map[string]string {
	if source == nil {
		return map[string]string{}
	}

	clone := make(map[string]string, len(source))
	for key, value := range source {
		clone[key] = value
	}
	return clone
}
