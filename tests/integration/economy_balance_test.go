package integration_test

import (
	"testing"

	"github.com/oberones/star-smuggler-go/internal/application"
	"github.com/oberones/star-smuggler-go/internal/domain"
	"github.com/oberones/star-smuggler-go/internal/services"
)

func TestRepeatedShortRouteFarmingIncreasesQuotedTravelCost(t *testing.T) {
	t.Parallel()

	snapshot := loadSnapshot(t)
	commands := application.NewRunCommands(snapshot, nil, seededRuntime(31))
	run, err := commands.StartNewRun()
	if err != nil {
		t.Fatalf("start run: %v", err)
	}

	run.Player.CurrentPortID = "mercury"
	run.Player.Credits = 1000

	initialQuotes, err := commands.PreviewTravel(run)
	if err != nil {
		t.Fatalf("preview initial travel: %v", err)
	}

	baseCost := 0
	for _, quote := range initialQuotes {
		if quote.Destination.ID == "venus" {
			baseCost = quote.Cost
			break
		}
	}
	if baseCost == 0 {
		t.Fatal("expected Mercury -> Venus quote")
	}

	if _, err := commands.CommitBaselineTravel(&run, "venus"); err != nil {
		t.Fatalf("travel to venus: %v", err)
	}

	nextQuotes, err := commands.PreviewTravel(run)
	if err != nil {
		t.Fatalf("preview follow-up travel: %v", err)
	}

	returnCost := 0
	for _, quote := range nextQuotes {
		if quote.Destination.ID == "mercury" {
			returnCost = quote.Cost
			break
		}
	}
	if returnCost <= baseCost {
		t.Fatalf("expected return quote %d to exceed base cost %d after route farming pressure", returnCost, baseCost)
	}
}

func TestSingleCommodityPressureCompressesPriceSpread(t *testing.T) {
	t.Parallel()

	snapshot := loadSnapshot(t)
	balance := services.EconomyBalanceService{}
	run := domain.NewRunState()
	run.MarketsByPortID["mars"] = domain.MarketSnapshot{
		PortID:           "mars",
		AvailableItemIDs: []string{"synthspice"},
		PricesByItemID: map[string]int{
			"synthspice": 100,
		},
	}

	for i := 0; i < 4; i++ {
		balance.RecordCommodityTrade(&run, "synthspice", 1)
	}
	balance.ApplyMarketPressure(&run, snapshot)

	adjusted := run.MarketsByPortID["mars"].PricesByItemID["synthspice"]
	if adjusted >= 100 {
		t.Fatalf("expected commodity pressure to reduce the stretched price, got %d", adjusted)
	}
	if adjusted <= snapshot.ItemsByID["synthspice"].BasePrice {
		t.Fatalf("expected price compression, not collapse to or below base price, got %d", adjusted)
	}
}
