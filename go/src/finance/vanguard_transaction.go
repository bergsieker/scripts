package finance

import (
        "encoding/csv"
        "io"
        "strconv"
        "strings"
        "time"
)

const type_Distribution = "Distribution"
const desc_Dividend = "INCOME DIVIDEND"
const desc_Lt_Cap_Gain = "LT CAP GAIN"

type VanguardTransactionReader struct {
        c *csv.Reader
}

func NewVanguardTransactionReader(r io.Reader) *VanguardTransactionReader {
        c := csv.NewReader(r)
        c.Comment = '#'
        return &VanguardTransactionReader{c}
}

func (r VanguardTransactionReader)Read() (Transaction, error) {
        var t Transaction
        record, err := r.c.Read()
        if err != nil {
                return t, err
        }
        date, err := time.Parse("01/02/2006", record[1])
        if err != nil {
                return t, err
        }
        account := record[0]
        assetName := record[5]
        amount, err := strconv.ParseFloat(record[8], 32)
        if err != nil {
                return t, err
        }
        t = Transaction{
                date,
                account,
                assetName,
                float32(amount),
                r.isAcquisition(record[3]),
                r.isIncome(record[3], record[4]),
                r.isCapGain(record[3], record[4]),
        }
        return t, nil
}

func (r *VanguardTransactionReader)isAcquisition(transType string) bool {
        return transType == type_Distribution
}

func (r *VanguardTransactionReader)isIncome(transType string, transDesc string) bool {
        return transType == type_Distribution && (strings.HasPrefix(transDesc, desc_Dividend) || strings.HasPrefix(transDesc, desc_Lt_Cap_Gain))
}

func (r *VanguardTransactionReader)isCapGain(transType string, transDesc string) bool {
        return transType == type_Distribution && strings.HasPrefix(transDesc, desc_Lt_Cap_Gain)
}
