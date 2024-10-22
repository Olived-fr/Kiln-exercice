package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Delegation struct {
	ID        int             `db:"id"`
	Datetime  time.Time       `db:"datetime"`
	Amount    decimal.Decimal `db:"amount"`
	Delegator string          `db:"delegator"`
	Height    int             `db:"height"`
	TxHash    string          `db:"tx_hash"`
}
