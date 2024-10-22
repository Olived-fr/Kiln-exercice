package poll

import (
	"kiln-exercice/internal/model"
	"kiln-exercice/pkg/tzkt"
)

func convertToModelDelegations(delegations []tzkt.Delegation) []model.Delegation {
	var modelDelegations []model.Delegation
	for _, d := range delegations {
		modelDelegations = append(modelDelegations, model.Delegation{
			Datetime:  d.Timestamp,
			Amount:    d.Amount,
			Delegator: d.Sender.Address,
			Height:    d.Level,
			TxHash:    d.Hash,
		})
	}
	return modelDelegations
}
