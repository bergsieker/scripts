package finance

import (
	"encoding/csv"
	"fmt"
	"io"
	"regexp"
	"strconv"
	// "strings"
	"time"
)

type VanguardOfxTransactionListReader struct {
	c    *csv.Reader
	acct string
}

func NewVanguardOfxTransactionListReader(r io.Reader, acct string) *VanguardOfxTransactionListReader {
	c := csv.NewReader(r)
	c.Comment = '#'
	return &VanguardOfxTransactionListReader{c, acct}
}

func (r VanguardOfxTransactionListReader) Read() (Transaction, error) {
	var t Transaction
	for {
		record, err := r.c.Read()
		if err != nil {
			return t, err
		}
		if record[0] != r.acct {
			continue
		}
		date, err := time.Parse("01/02/2006", record[1])
		if err != nil {
			fmt.Println("Error parsing date: ", err)
			return t, err
		}
		assetName := record[5]
		amt := record[8]
		amount, err := strconv.ParseFloat(amt, 32)
		if err != nil {
			fmt.Printf("Record: %+q, parsed: %+q", record[8], amt)
			return t, err
		}
		numShares, err := strconv.ParseFloat(record[7], 32)
		if err != nil {
			fmt.Println("Error parsing shares: ", err)
			return t, err
		}
		transType, err := parseVanguardOfxTransactionType(record[4])
		if err != nil {
			fmt.Println("Error parsing type: ", err)
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
}

func parseVanguardOfxTransactionType(t string) (TransactionType, error) {
	bt := []byte(t)
	m, err := regexp.Match("^INCOME DIVIDEND.*", bt)
	if err == nil && m {
		return TransactionType_Dividend, nil
	}
	m, err = regexp.Match("^LT CAP GAIN.*", bt)
	if err == nil && m {
		return TransactionType_LtCapGain, nil
	}
	m, err = regexp.Match("^ST CAP GAIN.*", bt)
	if err == nil && m {
		return TransactionType_StCapGain, nil
	}
	m, err = regexp.Match("^ADMIRAL CONVERSION.*", bt)
	if err == nil && m {
		return TransactionType_Transfer, nil
	}
	return TransactionType_Unknown, fmt.Errorf("Unknown transaction type: %s", t)
}
