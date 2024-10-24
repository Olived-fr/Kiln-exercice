package tzkt

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSDK(t *testing.T) {
	rawURL := "https://example.com"
	sdk, err := NewSDK(rawURL)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	parsedURL, _ := url.Parse(rawURL)
	if sdk.url.String() != parsedURL.String() {
		t.Errorf("expected URL %v, got %v", parsedURL, sdk.url)
	}

	assert.NotNil(t, sdk.client)
}

func TestSDK_GetDelegations(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx  context.Context
		from time.Time
		to   time.Time
	}

	tests := []struct {
		name    string
		args    args
		handler http.Handler
		want    []Delegation
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "happy path",
			args: args{
				ctx:  context.Background(),
				from: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				to:   time.Date(2021, 1, 11, 0, 0, 0, 0, time.UTC),
			},
			handler: http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					assert.Equal(t, http.MethodGet, r.Method)
					assert.Equal(t, "/v1/operations/delegations", r.URL.Path)
					assert.Equal(t, "2021-01-01T00:00:00Z", r.URL.Query().Get("timestamp.ge"))
					assert.Equal(t, "2021-01-11T00:00:00Z", r.URL.Query().Get("timestamp.lt"))
					assert.Equal(t, "10000", r.URL.Query().Get("limit"))

					var body []byte
					switch r.URL.Query().Get("offset") {
					case "10000":
						body = []byte(`[]`)
					case "0":
						body = []byte(`[
					{
						"type": "delegation",
						"id": 1098907648,
						"level": 109,
						"timestamp": "2018-06-30T19:30:27Z",
						"block": "BLwRUPupdhP8TyWp9J6TbjLSCxPPW6tyhVPF2KmNAbLPt7thjPw",
						"hash": "ooP37LNma6DiWjVxDbS2XZu4PiNKy7fbHZWSn8Vj8FX1hWfkC3b",
						"counter": 23,
						"sender": {
							"address": "tz1Wit2PqodvPeuRRhdQXmkrtU8e8bRYZecd"
						},
						"gasLimit": 0,
						"gasUsed": 0,
						"storageLimit": 0,
						"bakerFee": 50000,
						"amount": 25079312620,
						"newDelegate": {
							"address": "tz1Wit2PqodvPeuRRhdQXmkrtU8e8bRYZecd"
						},
						"status": "applied"
					},
					{
						"type": "delegation",
						"id": 1649410048,
						"level": 167,
						"timestamp": "2018-06-30T20:29:42Z",
						"block": "BLzCkTwQGUf9ggfk24bW7YFeFzudfncp5zzJSkR6Lf4kSP923PK",
						"hash": "ooedoJWn6fFXaiCkNDftRVCvEbJ855M7fD7gzHryL6x6FXdejP4",
						"counter": 34,
						"sender": {
							"address": "tz1U2ufqFdVkN2RdYormwHtgm3ityYY1uqft"
						},
						"gasLimit": 0,
						"gasUsed": 0,
						"storageLimit": 0,
						"bakerFee": 100,
						"amount": 10199999690,
						"newDelegate": {
							"address": "tz1U2ufqFdVkN2RdYormwHtgm3ityYY1uqft"
						},
						"status": "applied"
					}
					]`)
					default:
						t.Fatalf("unexpected offset: %s", r.URL.Query().Get("offset"))
					}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)

					if _, err := w.Write(body); err != nil {
						t.Fatal(err)
					}
				},
			),
			want: []Delegation{
				{
					Type:      "delegation",
					ID:        1098907648,
					Level:     109,
					Timestamp: time.Date(2018, 6, 30, 19, 30, 27, 0, time.UTC),
					Block:     "BLwRUPupdhP8TyWp9J6TbjLSCxPPW6tyhVPF2KmNAbLPt7thjPw",
					Hash:      "ooP37LNma6DiWjVxDbS2XZu4PiNKy7fbHZWSn8Vj8FX1hWfkC3b",
					Counter:   23,
					Sender: Sender{
						Alias:   "",
						Address: "tz1Wit2PqodvPeuRRhdQXmkrtU8e8bRYZecd",
					},
					SenderCodeHash:      0,
					Nonce:               0,
					GasLimit:            0,
					GasUsed:             0,
					StorageLimit:        0,
					BakerFee:            50000,
					Amount:              decimal.RequireFromString("25079312620"),
					StakingUpdatesCount: 0,
					PrevDelegate: Delegate{
						Alias:   "",
						Address: "",
					},
					NewDelegate: Delegate{
						Alias:   "",
						Address: "tz1Wit2PqodvPeuRRhdQXmkrtU8e8bRYZecd",
					},
					Status: "applied",
				},
				{
					Type:      "delegation",
					ID:        1649410048,
					Level:     167,
					Timestamp: time.Date(2018, 6, 30, 20, 29, 42, 0, time.UTC),
					Block:     "BLzCkTwQGUf9ggfk24bW7YFeFzudfncp5zzJSkR6Lf4kSP923PK",
					Hash:      "ooedoJWn6fFXaiCkNDftRVCvEbJ855M7fD7gzHryL6x6FXdejP4",
					Counter:   34,
					Sender: Sender{
						Alias:   "",
						Address: "tz1U2ufqFdVkN2RdYormwHtgm3ityYY1uqft",
					},
					SenderCodeHash:      0,
					Nonce:               0,
					GasLimit:            0,
					GasUsed:             0,
					StorageLimit:        0,
					BakerFee:            100,
					Amount:              decimal.RequireFromString("10199999690"),
					StakingUpdatesCount: 0,
					PrevDelegate: Delegate{
						Alias:   "",
						Address: "",
					},
					NewDelegate: Delegate{
						Alias:   "",
						Address: "tz1U2ufqFdVkN2RdYormwHtgm3ityYY1uqft",
					},
					Status: "applied",
				},
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				t.Parallel()

				server := httptest.NewServer(tt.handler)
				defer server.Close()

				s, err := NewSDK(server.URL)
				require.NoError(t, err)

				got, err := s.GetDelegations(tt.args.ctx, tt.args.from, tt.args.to)
				if !tt.wantErr(t, err, fmt.Sprintf("GetDelegations(%v, %v)", tt.args.ctx, tt.args.from)) {
					return
				}

				assert.Equalf(t, tt.want, got, "GetDelegations(%v, %v)", tt.args.ctx, tt.args.from)
			},
		)
	}
}
