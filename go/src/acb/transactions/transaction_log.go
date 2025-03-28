package transactions

import (
	"time"
)

type TransactionType int

const (
	TransactionType_Unknown = iota
	TransactionType_StartingBalance
	TransactionType_Acquisition
	TransactionType_Transfer
	TransactionType_Dividend
	TransactionType_LtCapGain
	TransactionType_StCapGain
	TransactionType_Disposition
)

type Transaction struct {
	Date time.Time
	AssetName string
	Type TransactionType
	LotRef string
	// Account string
	// Price float32
	CashAmount float32
	ShareAmount float32
	Fees float32
}
