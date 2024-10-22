//go:generate mockery
package poll

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"

	"kiln-exercice/internal/model"
	"kiln-exercice/internal/usecase/delegation/poll/mocks"
	"kiln-exercice/pkg/tzkt"
)

func TestNewUseCase(t *testing.T) {
	delegationRepo := mocks.NewDelegationRepository(t)
	pollingRepo := mocks.NewPollingRepository(t)
	xtzSDK := mocks.NewXTZSDK(t)
	defaultPollingFrom := time.Now()
	timeNow := func() time.Time { return time.Unix(0, 1).UTC() }

	useCase := NewUseCase(delegationRepo, pollingRepo, xtzSDK, defaultPollingFrom, timeNow)

	if useCase.DelegationRepo != delegationRepo {
		t.Errorf("expected %v, got %v", delegationRepo, useCase.DelegationRepo)
	}
	if useCase.XTZSDK != xtzSDK {
		t.Errorf("expected %v, got %v", xtzSDK, useCase.XTZSDK)
	}
	if useCase.PollingRepo != pollingRepo {
		t.Errorf("expected %v, got %v", pollingRepo, useCase.PollingRepo)
	}
	if !useCase.DefaultPollingFrom.Equal(defaultPollingFrom) {
		t.Errorf("expected %v, got %v", defaultPollingFrom, useCase.DefaultPollingFrom)
	}
	if useCase.TimeNow() != timeNow() {
		t.Errorf("expected %v, got %v", timeNow(), useCase.TimeNow())
	}
}

func TestUseCase_PollDelegations(t *testing.T) {

	type env struct {
		DelegationRepo *mocks.DelegationRepository
		PollingRepo    *mocks.PollingRepository
		XTZSDK         *mocks.XTZSDK
		PollingFrom    time.Time
		TimeNow        func() time.Time
	}

	tests := []struct {
		name    string
		env     env
		init    func(*env)
		wantErr bool
	}{
		{
			name: "happy path first polling",
			env: env{
				DelegationRepo: mocks.NewDelegationRepository(t),
				PollingRepo:    mocks.NewPollingRepository(t),
				XTZSDK:         mocks.NewXTZSDK(t),
				PollingFrom:    time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC).UTC(),
				TimeNow:        time.Date(2021, 1, 11, 0, 0, 0, 0, time.UTC).UTC,
			},
			init: func(e *env) {
				e.PollingRepo.EXPECT().GetLastPolling(
					mock.MatchedBy(
						func(_ context.Context) bool {
							return true
						},
					),
				).Return(model.Polling{}, sql.ErrNoRows)
				e.XTZSDK.EXPECT().GetDelegations(
					mock.MatchedBy(
						func(_ context.Context) bool {
							return true
						},
					),
					e.PollingFrom,
					e.PollingFrom.Add(pollingDaysByWorker*24*time.Hour),
				).Return(
					[]tzkt.Delegation{
						{
							Timestamp: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
							Amount:    decimal.RequireFromString("1000"),
							Sender:    tzkt.Sender{Address: "tz1SenderAddress"},
							Level:     1,
							Hash:      "txHash1",
						},
						{
							Timestamp: time.Date(2021, 1, 2, 1, 0, 0, 0, time.UTC),
							Amount:    decimal.RequireFromString("2000"),
							Sender:    tzkt.Sender{Address: "tz1SenderAddress2"},
							Level:     2,
							Hash:      "txHash2",
						},
					}, nil,
				)
				e.DelegationRepo.EXPECT().InsertDelegations(
					mock.MatchedBy(
						func(_ context.Context) bool {
							return true
						},
					), []model.Delegation{
						{
							Datetime:  time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
							Amount:    decimal.RequireFromString("1000"),
							Delegator: "tz1SenderAddress",
							Height:    1,
							TxHash:    "txHash1",
						},
						{
							Datetime:  time.Date(2021, 1, 2, 1, 0, 0, 0, time.UTC),
							Amount:    decimal.RequireFromString("2000"),
							Delegator: "tz1SenderAddress2",
							Height:    2,
							TxHash:    "txHash2",
						},
					},
				).Return(nil)
				e.PollingRepo.EXPECT().UpsertPolling(
					mock.MatchedBy(
						func(_ context.Context) bool { return true },
					),
					model.Polling{
						LastPolledAt: e.TimeNow(),
					},
				).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "happy path subsequent polling",
			env: env{
				DelegationRepo: mocks.NewDelegationRepository(t),
				PollingRepo:    mocks.NewPollingRepository(t),
				XTZSDK:         mocks.NewXTZSDK(t),
				PollingFrom:    time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC).UTC(),
				TimeNow:        time.Date(2021, 1, 11, 0, 0, 0, 0, time.UTC).UTC,
			},
			init: func(e *env) {
				lastPolling := model.Polling{
					ID:           1,
					LastPolledAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC).UTC(),
				}
				e.PollingRepo.EXPECT().GetLastPolling(
					mock.MatchedBy(
						func(_ context.Context) bool {
							return true
						},
					),
				).Return(lastPolling, nil)
				e.XTZSDK.EXPECT().GetDelegations(
					mock.MatchedBy(
						func(_ context.Context) bool {
							return true
						},
					),
					lastPolling.LastPolledAt,
					lastPolling.LastPolledAt.Add(pollingDaysByWorker*24*time.Hour),
				).Return(
					[]tzkt.Delegation{
						{
							Timestamp: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
							Amount:    decimal.RequireFromString("1000"),
							Sender:    tzkt.Sender{Address: "tz1SenderAddress"},
							Level:     1,
							Hash:      "txHash1",
						},
						{
							Timestamp: time.Date(2021, 1, 2, 1, 0, 0, 0, time.UTC),
							Amount:    decimal.RequireFromString("2000"),
							Sender:    tzkt.Sender{Address: "tz1SenderAddress2"},
							Level:     2,
							Hash:      "txHash2",
						},
					}, nil,
				)
				e.DelegationRepo.EXPECT().InsertDelegations(
					mock.MatchedBy(
						func(_ context.Context) bool {
							return true
						},
					), []model.Delegation{
						{
							Datetime:  time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
							Amount:    decimal.RequireFromString("1000"),
							Delegator: "tz1SenderAddress",
							Height:    1,
							TxHash:    "txHash1",
						},
						{
							Datetime:  time.Date(2021, 1, 2, 1, 0, 0, 0, time.UTC),
							Amount:    decimal.RequireFromString("2000"),
							Delegator: "tz1SenderAddress2",
							Height:    2,
							TxHash:    "txHash2",
						},
					},
				).Return(nil)
				e.PollingRepo.EXPECT().UpsertPolling(
					mock.MatchedBy(
						func(_ context.Context) bool { return true },
					),
					model.Polling{
						ID:           1,
						LastPolledAt: e.TimeNow(),
					},
				).Return(nil)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				t.Parallel()

				tt.init(&tt.env)

				ctx, cancel := context.WithCancel(context.Background())
				time.AfterFunc(time.Second, cancel)

				uc := &UseCase{
					DelegationRepo:     tt.env.DelegationRepo,
					PollingRepo:        tt.env.PollingRepo,
					XTZSDK:             tt.env.XTZSDK,
					DefaultPollingFrom: tt.env.PollingFrom,
					TimeNow:            tt.env.TimeNow,
				}

				err := uc.PollDelegations(ctx)
				if (err != nil) != tt.wantErr {
					t.Errorf("PollDelegations() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			},
		)
	}
}
