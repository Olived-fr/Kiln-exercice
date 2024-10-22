package poll

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"

	"kiln-exercice/internal/model"
	"kiln-exercice/pkg/tzkt"
)

func TestConvertToModelDelegations(t *testing.T) {
	delegations := []tzkt.Delegation{
		{
			Timestamp: time.Now(),
			Amount:    decimal.RequireFromString("1000"),
			Sender:    tzkt.Sender{Address: "tz1SenderAddress"},
			Level:     1,
			Hash:      "txHash1",
		},
		{
			Timestamp: time.Now().Add(time.Hour),
			Amount:    decimal.RequireFromString("2000"),
			Sender:    tzkt.Sender{Address: "tz1SenderAddress2"},
			Level:     2,
			Hash:      "txHash2",
		},
	}

	expected := []model.Delegation{
		{
			Datetime:  delegations[0].Timestamp,
			Amount:    delegations[0].Amount,
			Delegator: delegations[0].Sender.Address,
			Height:    delegations[0].Level,
			TxHash:    delegations[0].Hash,
		},
		{
			Datetime:  delegations[1].Timestamp,
			Amount:    delegations[1].Amount,
			Delegator: delegations[1].Sender.Address,
			Height:    delegations[1].Level,
			TxHash:    delegations[1].Hash,
		},
	}

	result := convertToModelDelegations(delegations)

	if len(result) != len(expected) {
		t.Fatalf("expected %d delegations, got %d", len(expected), len(result))
	}

	for i, r := range result {
		if r != expected[i] {
			t.Errorf("expected %v, got %v", expected[i], r)
		}
	}
}
