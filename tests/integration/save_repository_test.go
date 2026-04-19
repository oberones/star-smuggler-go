package integration_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/oberones/star-smuggler-go/internal/domain"
	"github.com/oberones/star-smuggler-go/internal/persistence"
)

func TestJSONSaveRepositoryRoundTrip(t *testing.T) {
	t.Parallel()

	savePath := filepath.Join(t.TempDir(), "save.json")
	repository := persistence.NewJSONSaveRepository(persistence.OSFileStore{}, savePath)

	run := domain.NewRunState()
	run.Player.Credits = 725
	run.Player.CargoLimit = 42
	run.Player.CurrentPortID = "mars"
	run.Cargo.SetQuantity("synthspice", 3)
	run.Cargo.SetQuantity("alloy", 5)
	run.MarketsByPortID["mars"] = domain.MarketSnapshot{
		PortID:           "mars",
		AvailableItemIDs: []string{"synthspice", "alloy"},
		PricesByItemID: map[string]int{
			"synthspice": 48,
			"alloy":      31,
		},
	}
	run.JumpsSinceLastUpdate = 2
	run.TotalJumps = 7
	run.RecentEvent = &domain.EventResult{
		EventID:             "merchant_strike",
		Name:                "Merchant Strike",
		ResolvedDescription: "Prices surge at the port.",
		RolledValues: map[string]float64{
			"multiplier": 2,
		},
	}

	if err := repository.Save(context.Background(), run); err != nil {
		t.Fatalf("save run: %v", err)
	}

	exists, err := repository.Exists()
	if err != nil {
		t.Fatalf("check save existence: %v", err)
	}
	if !exists {
		t.Fatal("expected save file to exist")
	}

	loaded, err := repository.Load(context.Background())
	if err != nil {
		t.Fatalf("load run: %v", err)
	}

	if loaded.Player != run.Player {
		t.Fatalf("player mismatch: %#v != %#v", loaded.Player, run.Player)
	}
	if !reflect.DeepEqual(loaded.Cargo.ItemQuantities, run.Cargo.ItemQuantities) {
		t.Fatalf("cargo mismatch: %#v != %#v", loaded.Cargo.ItemQuantities, run.Cargo.ItemQuantities)
	}
	if !reflect.DeepEqual(loaded.MarketsByPortID, run.MarketsByPortID) {
		t.Fatalf("markets mismatch: %#v != %#v", loaded.MarketsByPortID, run.MarketsByPortID)
	}
	if loaded.JumpsSinceLastUpdate != run.JumpsSinceLastUpdate {
		t.Fatalf("jump refresh mismatch: %d != %d", loaded.JumpsSinceLastUpdate, run.JumpsSinceLastUpdate)
	}
	if loaded.TotalJumps != run.TotalJumps {
		t.Fatalf("total jumps mismatch: %d != %d", loaded.TotalJumps, run.TotalJumps)
	}
	if !reflect.DeepEqual(loaded.RecentEvent, run.RecentEvent) {
		t.Fatalf("recent event mismatch: %#v != %#v", loaded.RecentEvent, run.RecentEvent)
	}
}

func TestJSONSaveRepositoryRejectsUnknownVersion(t *testing.T) {
	t.Parallel()

	savePath := filepath.Join(t.TempDir(), "save.json")
	rawSave, err := json.Marshal(persistence.SaveData{
		Version: 999,
	})
	if err != nil {
		t.Fatalf("marshal invalid save: %v", err)
	}

	if err := os.WriteFile(savePath, rawSave, 0o644); err != nil {
		t.Fatalf("write invalid save: %v", err)
	}

	repository := persistence.NewJSONSaveRepository(persistence.OSFileStore{}, savePath)
	if _, err := repository.Load(context.Background()); err == nil {
		t.Fatal("expected unknown-version error")
	}
}
