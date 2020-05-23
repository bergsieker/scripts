package finance

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"time"
)

type ITradeTransactionListParser struct {
	c *csv.Reader
}

func NewITradeTransactionListParser(r io.Reader) *ITradeTransactionListParser {
	c := csv.NewReader(r)
	c.Comment = '#'
	return &ITradeTransactionListParser{c}
}

func (r ITradeTransactionListParser) Read() (Transaction, error) {
	var t Transaction
	record, err := r.c.Read()
	if err != nil {
		return t, err
	}
	date, err := time.Parse("02-Jan-2006", record[2])
	if err != nil {
		return t, err
	}
	assetName := record[1]
	transType, err := parseITradeTransactionType(record[5])
	if err != nil {
		return t, err
	}
	amt := record[9]
	amount, err := strconv.ParseFloat(amt, 32)
	if err != nil {
		return t, err
	}
	if transType == TransactionType_Buy {
		amount *= -1
	}
	numShares, err := strconv.ParseFloat(record[6], 32)
	if err != nil {
		return t, err
	}
	t = Transaction{
		date,
		"",
		assetName,
		float32(amount),
		transType,
		float32(numShares),
	}
	return t, nil
}

func parseITradeTransactionType(t string) (TransactionType, error) {
	switch t {
	case "CASH DIV":
		return TransactionType_Dividend, nil
	case "L/T CAP GNS":
		return TransactionType_LtCapGain, nil
	case "BUY":
		return TransactionType_Buy, nil
	case "TRT":
		return TransactionType_Transfer, nil
	case "TRANSFER":
		return TransactionType_Transfer, nil
	default:
		return TransactionType_Unknown, fmt.Errorf("Unknown transaction type: %s", t)
	}
}
