package tzkt

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"
	"golang.org/x/time/rate"
)

const rateLimit = 10

type SDK struct {
	url    *url.URL
	client *resty.Client
}

func NewSDK(rawURL string) (*SDK, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	return &SDK{
		url: u,
		client: resty.New().SetRateLimiter(
			rate.NewLimiter(
				rate.Limit(rateLimit), rateLimit,
			),
		), // 10 requests per second with the free plan
	}, nil
}

func (s *SDK) GetDelegations(ctx context.Context, from, to time.Time) (delegations []Delegation, err error) {
	const path = "/v1/operations/delegations"
	var (
		result      []Delegation
		resultError error
		offset      = 0
		limit       = 10000
	)

	ctx, cancel := context.WithTimeout(ctx, 45*time.Second)
	defer cancel()

	for {
		resp, err := s.client.R().
			SetContext(ctx).
			SetResult(&result).
			SetError(&resultError).
			SetQueryParams(
				map[string]string{
					"timestamp.ge": from.Format(time.RFC3339),
					"timestamp.lt": to.Format(time.RFC3339),
					"offset":       strconv.Itoa(offset),
					"limit":        strconv.Itoa(limit),
				},
			).
			Get(s.url.String() + path)

		if err != nil {
			if errors.Is(err, resty.ErrRateLimitExceeded) {
				time.Sleep(rateLimit * time.Second)
				return s.GetDelegations(ctx, from, to)
			}

			return result, err
		}

		if !resp.IsSuccess() {
			return result, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
		}

		if len(result) == 0 {
			break
		}

		delegations = append(delegations, result...)
		offset += limit
	}

	return delegations, nil
}
