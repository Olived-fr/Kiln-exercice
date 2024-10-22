package pg

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"kiln-exercice/internal/model"
	"kiln-exercice/pkg/pgtest"
)

func initPollingDeps(ctx context.Context, t *testing.T) (*pgtest.Helper, *PollingRepository) {
	t.Helper()

	c := pgtest.NewPostgresContainer(ctx, t)

	return pgtest.NewHelper(c.GetDB()), NewPollingRepository(c.GetDB())
}

func TestPollingRepository_GetLastPolling(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	h, repo := initPollingDeps(ctx, t)

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name    string
		args    args
		init    func(h *pgtest.Helper)
		want    model.Polling
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "no polling",
			args: args{
				ctx: ctx,
			},
			init: func(h *pgtest.Helper) {
				return
			},
			want:    model.Polling{},
			wantErr: assert.Error,
		},
		{
			name: "get last polling",
			args: args{
				ctx: ctx,
			},
			init: func(h *pgtest.Helper) {
				h.MustInject(
					ctx, t, []pgtest.RecordSet{
						{
							Table: "polling",
							Records: []pgtest.Record{
								{
									"last_polled_at": time.Date(2022, 5, 5, 6, 29, 14, 0, time.UTC),
								},
							},
						},
					},
				)
			},
			want: model.Polling{
				ID:           1,
				LastPolledAt: time.Date(2022, 5, 5, 6, 29, 14, 0, time.UTC),
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				t.Parallel()

				tt.init(h)

				got, err := repo.GetLastPolling(tt.args.ctx)
				if !tt.wantErr(t, err, fmt.Sprintf("GetLastPolling(%v)", tt.args.ctx)) {
					return
				}
				assert.Equalf(t, tt.want, got, "GetLastPolling(%v)", tt.args.ctx)
			},
		)
	}
}

func TestPollingRepository_UpsertPolling(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	h, repo := initPollingDeps(ctx, t)

	type args struct {
		ctx     context.Context
		polling model.Polling
	}

	tests := []struct {
		name    string
		args    args
		init    func(h *pgtest.Helper)
		want    []pgtest.RecordSet
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "insert polling",
			args: args{
				ctx: ctx,
				polling: model.Polling{
					LastPolledAt: time.Date(2022, 5, 5, 6, 29, 14, 0, time.UTC),
				},
			},
			init: func(h *pgtest.Helper) {
				return
			},
			want: []pgtest.RecordSet{
				{
					Table: "polling",
					Records: []pgtest.Record{
						{
							"last_polled_at": time.Date(2022, 5, 5, 6, 29, 14, 0, time.UTC),
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "update polling",
			args: args{
				ctx: ctx,
				polling: model.Polling{
					ID:           1,
					LastPolledAt: time.Date(2022, 5, 5, 6, 29, 14, 0, time.UTC),
				},
			},
			init: func(h *pgtest.Helper) {
				h.MustInject(
					ctx, t, []pgtest.RecordSet{
						{
							Table: "polling",
							Records: []pgtest.Record{
								{
									"last_polled_at": time.Date(2022, 5, 5, 6, 29, 14, 0, time.UTC),
								},
							},
						},
					},
				)
			},
			want: []pgtest.RecordSet{
				{
					Table: "polling",
					Records: []pgtest.Record{
						{
							"id":             1,
							"last_polled_at": time.Date(2022, 5, 5, 6, 29, 14, 0, time.UTC),
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				t.Parallel()

				tt.init(h)

				err := repo.UpsertPolling(ctx, tt.args.polling)
				if !tt.wantErr(t, err, fmt.Sprintf("InsertDelegations(%v)", tt.args.polling)) {
					return
				}

				if err == nil {
					h.MustCheck(ctx, t, tt.want)
				}
			},
		)
	}
}
