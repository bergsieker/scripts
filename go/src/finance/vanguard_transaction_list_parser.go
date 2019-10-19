package finance

import (
        "encoding/csv"
        "fmt"
        "io"
        "strings"
        "strconv"
        "time"
)

type VanguardTransactionListReader struct {
        c *csv.Reader
}

func NewVanguardTransactionListReader(r io.Reader) *VanguardTransactionListReader {
        c := csv.NewReader(r)
        c.Comment = '#'
        c.Comma = '\t'
        return &VanguardTransactionListReader{c}
}

func (r VanguardTransactionListReader)Read() (Transaction, error) {
        repl := strings.NewReplacer("$", "", " ", "", ",", "", "\u2013", "-")
        var t Transaction
        record, err := r.c.Read()
        if err != nil {
                return t, err
        }
        date, err := time.Parse("01/02/2006", record[0])
        if err != nil {
                return t, err
        }
        assetName := record[1]
        amt := record[5]
        amt = repl.Replace(amt)
        amount, err := strconv.ParseFloat(amt, 32)
        if err != nil {
                fmt.Printf("Record: %+q, sanitized: %+q, parsed: %+q\n", record[5], repl.Replace(record[5]), amt)
                return t, err
        }
        numShares, err := strconv.ParseFloat(repl.Replace(record[3]), 32)
        if err != nil {
                return t, err
        }
        transType, err := parseVanguardTransactionType(record[2])
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

func parseVanguardTransactionType(t string) (TransactionType, error) {
        switch t {
        case "Starting Balance":
                return TransactionType_StartingBalance, nil
        case "Dividend Received":
                return TransactionType_Dividend, nil
        case "Long-term capital gain":
                return TransactionType_LtCapGain, nil
        case "Short-term capital gain":
                return TransactionType_StCapGain, nil
        case "Sell":
                return TransactionType_Sale, nil
        case "Buy":
                return TransactionType_Buy, nil
        case "Transfer":
                return TransactionType_Transfer, nil
        default:
                return TransactionType_Unknown, fmt.Errorf("Unknown transaction type: %s", t)
        }
}
