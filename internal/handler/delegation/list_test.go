//go:generate mockery
package delegation

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"

	"kiln-exercice/internal/handler/delegation/mocks"
	delegationlist "kiln-exercice/internal/usecase/delegation/list"
	"kiln-exercice/pkg/api"
	"kiln-exercice/pkg/http/apitest"
)

func TestDelegationListHandler(t *testing.T) {
	t.Parallel()

	type env struct {
		useCase *mocks.DelegationUseCase
	}

	type args struct {
		w http.ResponseWriter
		r *http.Request
	}

	tests := []struct {
		name     string
		args     args
		init     func(*env)
		wantCode int
		wantBody string
		wantErr  bool
	}{
		{
			name: "happy path",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/?year=2021&page=1&page_size=10", nil),
			},
			init: func(e *env) {
				e.useCase.EXPECT().ListDelegations(
					mock.Anything,
					delegationlist.Input{Year: 2021, Pagination: api.Pagination{PageNumber: 1, PageSize: 10}},
				).Return(
					delegationlist.Output{
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
					}, nil,
				)
			},
			wantCode: http.StatusOK,
			wantBody: `
			{
			  "data": [ 
				{
					"timestamp": "2022-05-05T06:29:14Z",
					"amount": "125896",
					"delegator": "tz1a1SAaXRt9yoGMx29rh9FsBF4UzmvojdTL",
					"level": "2338084"
				},
				{
					"timestamp": "2021-05-07T14:48:07Z",
					"amount": "9856354",
					"delegator": "KT1JejNYjmQYh8yw95u5kfQDRuxJcaUPjUnf",
					"level": "1461334"
				}
			  ]
			}`,
			wantErr: false,
		},
		{
			name: "invalid year",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/?year=abc", nil),
			},
			init:     func(e *env) {},
			wantCode: http.StatusBadRequest,
			wantBody: `{"data": {"message": "invalid year format: abc"}}`,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				t.Parallel()

				e := env{
					useCase: mocks.NewDelegationUseCase(t),
				}

				tt.init(&e)

				apitest.TestHandler(
					t, tt.args.r, tt.wantCode, tt.wantBody, &DelegationListHandler{
						useCase: e.useCase,
					},
				)
			},
		)
	}
}
