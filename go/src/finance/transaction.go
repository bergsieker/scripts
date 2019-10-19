package finance

import (
        "time"
)

type TransactionType int

const (
        TransactionType_Unknown TransactionType = iota
        TransactionType_StartingBalance
        TransactionType_Dividend
        TransactionType_LtCapGain
        TransactionType_StCapGain
        TransactionType_Sale
        TransactionType_Buy
        TransactionType_Transfer
)

type Transaction struct {
        Date time.Time
        Account string
        AssetName string
        Amount float32
        Type TransactionType
        NumShares float32
}

func (t Transaction)IsAcquisition() bool {
        return t.NumShares > 0
}

func (t Transaction)IsIncome() bool {
        return t.Type == TransactionType_Dividend || t.Type == TransactionType_LtCapGain || t.Type == TransactionType_StCapGain
}

func (t Transaction)IsDisposition() bool {
        return t.Type == TransactionType_Sale
}
