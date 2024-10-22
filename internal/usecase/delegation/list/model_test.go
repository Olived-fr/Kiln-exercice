package list

import (
	"reflect"
	"testing"
	"time"

	"github.com/shopspring/decimal"

	"kiln-exercice/internal/model"
)

func TestBuildOutput(t *testing.T) {
	delegations := []model.Delegation{
		{
			Datetime:  time.Date(2022, 5, 5, 6, 29, 14, 0, time.UTC),
			Amount:    decimal.RequireFromString("125896"),
			Delegator: "tz1a1SAaXRt9yoGMx29rh9FsBF4UzmvojdTL",
			Height:    2338084,
		},
		{
			Datetime:  time.Date(2021, 5, 7, 14, 48, 7, 0, time.UTC),
			Amount:    decimal.RequireFromString("9856354"),
			Delegator: "KT1JejNYjmQYh8yw95u5kfQDRuxJcaUPjUnf",
			Height:    1461334,
		},
	}

	expected := Output{
		{
			Timestamp: time.Date(2022, 5, 5, 6, 29, 14, 0, time.UTC),
			Amount:    decimal.RequireFromString("125896"),
			Delegator: "tz1a1SAaXRt9yoGMx29rh9FsBF4UzmvojdTL",
			Level:     "2338084",
		},
		{
			Timestamp: time.Date(2021, 5, 7, 14, 48, 7, 0, time.UTC),
			Amount:    decimal.RequireFromString("9856354"),
			Delegator: "KT1JejNYjmQYh8yw95u5kfQDRuxJcaUPjUnf",
			Level:     "1461334",
		},
	}

	result := buildOutput(delegations)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("buildOutput() = %v, want %v", result, expected)
	}
}
