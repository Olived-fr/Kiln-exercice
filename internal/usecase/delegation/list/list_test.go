//go:generate mockery
package list

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/shopspring/decimal"

	"kiln-exercice/internal/model"
	"kiln-exercice/internal/usecase/delegation/list/mocks"
	"kiln-exercice/pkg/api"
)

func TestNewUseCase(t *testing.T) {
	t.Parallel()

	mockRepo := mocks.NewDelegationRepository(t)
	uc := NewUseCase(mockRepo)

	if uc.DelegationRepo != mockRepo {
		t.Errorf("NewUseCase() DelegationRepo = %v, want %v", uc.DelegationRepo, mockRepo)
	}
}

func TestUseCase_ListDelegations(t *testing.T) {
	t.Parallel()

	type env struct {
		DelegationRepo *mocks.DelegationRepository
	}

	type args struct {
		ctx   context.Context
		input Input
	}

	tests := []struct {
		name    string
		env     env
		args    args
		init    func(*env)
		want    Output
		wantErr bool
	}{
		{
			name: "happy path",
			env: env{
				DelegationRepo: mocks.NewDelegationRepository(t),
			},
			args: args{
				ctx: context.Background(),
				input: Input{
					Year: 2021,
					Pagination: api.Pagination{
						PageNumber: 1,
						PageSize:   10,
					},
				},
			},
			init: func(e *env) {
				e.DelegationRepo.EXPECT().ListDelegations(
					context.Background(), 2021, 0, 10,
				).Return(
					[]model.Delegation{
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
					},
					nil,
				)
			},
			want: Output{
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
			},
			wantErr: false,
		},
		{
			name: "invalid year",
			args: args{
				ctx: context.Background(),
				input: Input{
					Year: 1900,
				},
			},
			init:    func(e *env) {},
			want:    Output{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				t.Parallel()

				tt.init(&tt.env)

				uc := UseCase{
					DelegationRepo: tt.env.DelegationRepo,
				}

				got, err := uc.ListDelegations(tt.args.ctx, tt.args.input)
				if (err != nil) != tt.wantErr {
					t.Errorf("ListDelegations() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("ListDelegations() got = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
