package finance

import (
        "encoding/csv"
        "io"
        "strconv"
        "time"
)

const type_Dividend = "Dividend Received"
const type_Lt_Cap_Gain = "Long-term capital gain"
const type_St_Cap_Gain = "Short-term capital gain"

type VanguardStatementReader struct {
        c *csv.Reader
}

func NewVanguardStatementReader(r io.Reader) *VanguardStatementReader {
        c := csv.NewReader(r)
        c.Comment = '#'
        return &VanguardStatementReader{c}
}

func (r VanguardStatementReader)Read() (Transaction, error) {
        var t Transaction
        record, err := r.c.Read()
        if err != nil {
                return t, err
        }
        date, err := time.Parse("2006-01-02", record[1])
        if err != nil {
                return t, err
        }
        assetName := record[2]
        amount, err := strconv.ParseFloat(record[6], 32)
        if err != nil {
                return t, err
        }
        numShares, err := strconv.ParseFloat(record[4], 32)
        if err != nil {
                return t, err
        }
        t = Transaction{
                date,
                "",
                assetName,
                float32(amount),
                isAcquisition(float32(numShares)),
                isIncome(record[3]),
                isCapGain(record[3]),
        }
        return t, nil
}

func isAcquisition(numShares float32) bool {
        return numShares != 0
}

func isIncome(transType string) bool {
        return transType == type_Dividend || transType == type_Lt_Cap_Gain || transType == type_St_Cap_Gain
}

func isCapGain(transType string) bool {
        return false
}
