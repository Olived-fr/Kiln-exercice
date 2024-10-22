package list

import (
	"strconv"
	"time"

	"github.com/shopspring/decimal"

	"kiln-exercice/internal/model"
	"kiln-exercice/pkg/api"
)

type Input struct {
	Year int
	api.Pagination
}

type Output = []DelegationData

type DelegationData struct {
	Timestamp time.Time       `json:"timestamp"`
	Amount    decimal.Decimal `json:"amount"`
	Delegator string          `json:"delegator"`
	Level     string          `json:"level"`
}

func buildOutput(delegations []model.Delegation) Output {
	out := make([]DelegationData, len(delegations))

	for i, d := range delegations {
		out[i] = DelegationData{
			Timestamp: d.Datetime,
			Amount:    d.Amount,
			Delegator: d.Delegator,
			Level:     strconv.Itoa(d.Height),
		}
	}

	return out
}
