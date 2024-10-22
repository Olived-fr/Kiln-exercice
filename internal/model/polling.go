package model

import "time"

type Polling struct {
	ID           int64     `db:"id"`
	LastPolledAt time.Time `db:"last_polled_at"`
}
