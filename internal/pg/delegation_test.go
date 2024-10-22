package pg

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"kiln-exercice/internal/model"
	"kiln-exercice/pkg/pgtest"
)

func initDelegationDeps(ctx context.Context, t *testing.T) (*pgtest.Helper, *DelegationRepository) {
	t.Helper()

	c := pgtest.NewPostgresContainer(ctx, t)

	return pgtest.NewHelper(c.GetDB()), NewDelegationRepository(c.GetDB(), 1)
}

func TestInsertDelegations(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	h, repo := initDelegationDeps(ctx, t)

	type args struct {
		ctx         context.Context
		delegations []model.Delegation
	}

	tests := []struct {
		name    string
		args    args
		want    []pgtest.RecordSet
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "creation",
			args: args{
				ctx: ctx,
				delegations: []model.Delegation{
					{
						Datetime:  time.Date(2022, 5, 5, 6, 29, 14, 0, time.UTC),
						Amount:    decimal.RequireFromString("125896"),
						Delegator: "tz1a1SAaXRt9yoGMx29rh9FsBF4UzmvojdTL",
						Height:    2338084,
						TxHash:    "tx_hash_1",
					},
					{
						Datetime:  time.Date(2021, 5, 7, 14, 48, 7, 0, time.UTC),
						Amount:    decimal.RequireFromString("9856354"),
						Delegator: "KT1JejNYjmQYh8yw95u5kfQDRuxJcaUPjUnf",
						Height:    1461334,
						TxHash:    "tx_hash_2",
					},
				},
			},
			wantErr: assert.NoError,
			want: []pgtest.RecordSet{
				{
					Table: "delegation",
					Records: []pgtest.Record{
						{
							"datetime":  time.Date(2022, 5, 5, 6, 29, 14, 0, time.UTC),
							"delegator": "tz1a1SAaXRt9yoGMx29rh9FsBF4UzmvojdTL",
							"amount":    decimal.RequireFromString("125896"), "height": 2338084, "tx_hash": "tx_hash_1",
						},
						{
							"datetime":  time.Date(2021, 5, 7, 14, 48, 7, 0, time.UTC),
							"delegator": "KT1JejNYjmQYh8yw95u5kfQDRuxJcaUPjUnf",
							"amount":    decimal.RequireFromString("9856354"), "height": 1461334,
							"tx_hash": "tx_hash_2",
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				t.Parallel()

				err := repo.InsertDelegations(ctx, tt.args.delegations)
				if !tt.wantErr(t, err, fmt.Sprintf("InsertDelegations(%v)", tt.args.delegations)) {
					return
				}

				if err == nil {
					h.MustCheck(ctx, t, tt.want)
				}
			},
		)
	}
}

func TestListDelegations(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	h, repo := initDelegationDeps(ctx, t)

	h.MustInject(
		ctx, t, []pgtest.RecordSet{
			{
				Table: "delegation",
				Records: []pgtest.Record{
					{
						"datetime":  time.Date(2024, 5, 5, 6, 29, 14, 0, time.UTC),
						"delegator": "tz1a1SAaXRt9yoGMx29rh9FsBF4UzmvojdTL",
						"amount":    decimal.RequireFromString("125896"), "height": 2338084, "tx_hash": "tx_hash_1",
					},
					{
						"datetime":  time.Date(2024, 5, 7, 14, 48, 7, 0, time.UTC),
						"delegator": "KT1JejNYjmQYh8yw95u5kfQDRuxJcaUPjUnf",
						"amount":    decimal.RequireFromString("9856354"), "height": 1461334, "tx_hash": "tx_hash_2",
					},
				},
			},
		},
	)

	type args struct {
		ctx    context.Context
		year   int
		offset int
		limit  int
	}

	tests := []struct {
		name    string
		args    args
		want    []model.Delegation
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "list with year",
			args: args{
				ctx:    ctx,
				year:   2024,
				offset: 0,
				limit:  10,
			},
			wantErr: assert.NoError,
			want: []model.Delegation{
				{
					Datetime:  time.Date(2024, 5, 7, 14, 48, 7, 0, time.UTC),
					Amount:    decimal.RequireFromString("9856354"),
					Delegator: "KT1JejNYjmQYh8yw95u5kfQDRuxJcaUPjUnf",
					Height:    1461334,
				},
				{
					Datetime:  time.Date(2024, 5, 5, 6, 29, 14, 0, time.UTC),
					Amount:    decimal.RequireFromString("125896"),
					Delegator: "tz1a1SAaXRt9yoGMx29rh9FsBF4UzmvojdTL",
					Height:    2338084,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				t.Parallel()

				got, err := repo.ListDelegations(ctx, tt.args.year, tt.args.offset, tt.args.limit)
				if !tt.wantErr(t, err, fmt.Sprintf("ListDelegations(%v)", tt.args.year)) {
					return
				}

				assert.Equalf(t, tt.want, got, "ListDelegations(%v)", tt.args.year)
			},
		)
	}
}
