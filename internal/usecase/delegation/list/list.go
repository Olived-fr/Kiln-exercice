package list

import (
	"context"
	"fmt"
	"time"

	"kiln-exercice/internal/model"
	"kiln-exercice/pkg/api"
)

type DelegationRepository interface {
	ListDelegations(ctx context.Context, year, offset, limit int) ([]model.Delegation, error)
}

type UseCase struct {
	DelegationRepo DelegationRepository
}

func NewUseCase(delegationRepo DelegationRepository) *UseCase {
	return &UseCase{
		DelegationRepo: delegationRepo,
	}
}

// ListDelegations returns a list of delegations from the repository.
func (uc *UseCase) ListDelegations(ctx context.Context, input Input) (Output, error) {
	// This could be checked in an OpenAPI spec.
	// Arbitrary year range.
	if input.Year != 0 && (input.Year > time.Now().Year() || input.Year < 2018) {
		return Output{}, api.NewError(api.InvalidArgument, fmt.Sprintf("invalid year: %d", input.Year), nil)
	}

	delegations, err := uc.DelegationRepo.ListDelegations(ctx, input.Year, input.Offset(), input.Limit())
	if err != nil {
		return Output{}, api.NewError(api.Unknown, "error listing delegations", err)
	}

	return buildOutput(delegations), nil
}
