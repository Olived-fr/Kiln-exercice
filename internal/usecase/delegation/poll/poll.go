package poll

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"kiln-exercice/internal/model"
	"kiln-exercice/pkg/tzkt"
	"kiln-exercice/pkg/worker"
)

const pollingDaysByWorker = 100

type DelegationRepository interface {
	InsertDelegations(ctx context.Context, delegations []model.Delegation) error
}

type PollingRepository interface {
	GetLastPolling(ctx context.Context) (model.Polling, error)
	UpsertPolling(ctx context.Context, polling model.Polling) error
}

type XTZSDK interface {
	GetDelegations(ctx context.Context, from, to time.Time) ([]tzkt.Delegation, error)
}

type UseCase struct {
	DelegationRepo     DelegationRepository
	PollingRepo        PollingRepository
	XTZSDK             XTZSDK
	DefaultPollingFrom time.Time
	TimeNow            func() time.Time
}

func NewUseCase(delegationRepo DelegationRepository, pollingRepo PollingRepository, xtzSDK XTZSDK, pollingFrom time.Time, timeNow func() time.Time) *UseCase {
	return &UseCase{
		DelegationRepo:     delegationRepo,
		PollingRepo:        pollingRepo,
		XTZSDK:             xtzSDK,
		DefaultPollingFrom: pollingFrom,
		TimeNow:            timeNow,
	}
}

func (uc *UseCase) PollDelegations(ctx context.Context) error {
	from := uc.DefaultPollingFrom

	polling, err := uc.PollingRepo.GetLastPolling(ctx)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("get last polling: %w", err)
	}

	if !polling.LastPolledAt.IsZero() {
		from = polling.LastPolledAt
	}

	to := uc.TimeNow()

	numWorkers := int(uc.TimeNow().Sub(from).Hours() / 24 / pollingDaysByWorker)
	if numWorkers <= 1 {
		numWorkers = 1
	}
	fmt.Println("numWorkers", numWorkers)

	delegationsCh := make(chan []model.Delegation, numWorkers)
	defer close(delegationsCh)

	pool := worker.NewWorkerPool(ctx, numWorkers)
	pool.Start(uc.fetchDelegations)

	for i := 0; i < numWorkers; i++ {
		pool.Submit(from.Add(time.Duration(i) * time.Duration(pollingDaysByWorker) * 24 * time.Hour))
	}

	var (
		delegations []model.Delegation
		errs        []error
	)
	for i := 0; i < numWorkers; i++ {
		result := pool.GetResult()
		fmt.Println("result", len(result.Output.([]model.Delegation)), result.Err)
		if result.Err != nil {
			errs = append(errs, result.Err)
		}

		if result.Output != nil {
			delegations = append(delegations, result.Output.([]model.Delegation)...)
		}
	}

	pool.Stop()

	if len(errs) > 0 {
		errors.Join(errs...)
	}

	if len(delegations) == 0 {
		return nil
	}

	fmt.Println("inserting delegations", len(delegations))
	if err = uc.DelegationRepo.InsertDelegations(ctx, delegations); err != nil {
		return fmt.Errorf("insert delegations: %w", err)
	}

	polling.LastPolledAt = to
	if err = uc.PollingRepo.UpsertPolling(ctx, polling); err != nil {
		return fmt.Errorf("upsert polling: %w", err)
	}

	return nil
}

func (uc *UseCase) fetchDelegations(ctx context.Context, input any) (any, error) {
	from, ok := input.(time.Time)
	if !ok {
		return nil, fmt.Errorf("invalid input type, expected time.Time")
	}

	to := from.Add(time.Duration(pollingDaysByWorker) * 24 * time.Hour)

	if to.After(uc.TimeNow()) {
		to = uc.TimeNow()
	}
	delegations, err := uc.XTZSDK.GetDelegations(ctx, from, to)
	if err != nil {
		return []model.Delegation{}, fmt.Errorf("sdk get delegations: %w", err)
	}

	return convertToModelDelegations(delegations), nil
}
