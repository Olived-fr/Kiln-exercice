package tzkt

import (
	"time"

	"github.com/shopspring/decimal"
)

type Delegation struct {
	Type                 string          `json:"type"`
	ID                   int             `json:"id"`
	Level                int             `json:"level"`
	Timestamp            time.Time       `json:"timestamp"`
	Block                string          `json:"block"`
	Hash                 string          `json:"hash"`
	Counter              int             `json:"counter"`
	Initiator            Delegate        `json:"initiator"`
	Sender               Sender          `json:"sender"`
	SenderCodeHash       int             `json:"senderCodeHash"`
	Nonce                int             `json:"nonce"`
	GasLimit             int             `json:"gasLimit"`
	GasUsed              int             `json:"gasUsed"`
	StorageLimit         int             `json:"storageLimit"`
	BakerFee             int             `json:"bakerFee"`
	Amount               decimal.Decimal `json:"amount"`
	StakingUpdatesCount  int             `json:"stakingUpdatesCount"`
	PrevDelegate         Delegate        `json:"prevDelegate"`
	NewDelegate          Delegate        `json:"newDelegate"`
	Status               string          `json:"status"`
	Errors               []Error         `json:"errors"`
	Quote                Quote           `json:"quote"`
	UnstakedPseudotokens int             `json:"unstakedPseudotokens"`
	UnstakedBalance      int             `json:"unstakedBalance"`
	UnstakedRewards      int             `json:"unstakedRewards"`
}

type Sender struct {
	Alias   string `json:"alias"`
	Address string `json:"address"`
}

type Delegate struct {
	Alias   string `json:"alias"`
	Address string `json:"address"`
}

type Error struct {
	Type string `json:"type"`
}

type Quote struct {
	Btc decimal.Decimal `json:"btc"`
	Eur decimal.Decimal `json:"eur"`
	Usd decimal.Decimal `json:"usd"`
	Cny decimal.Decimal `json:"cny"`
	Jpy decimal.Decimal `json:"jpy"`
	Krw decimal.Decimal `json:"krw"`
	Eth decimal.Decimal `json:"eth"`
	Gbp decimal.Decimal `json:"gbp"`
}
