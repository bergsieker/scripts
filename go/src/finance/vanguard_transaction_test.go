package finance

import (
        "io"
        "strings"
        "testing"
        "time"
)

func TestVanguardTransactionReader_IsAcquisition(t *testing.T) {
        in := `1234567890,10/30/2015,10/30/2015,Distribution,INCOME DIVIDEND,GNMA Fund Investor Shares,10.7,1.254,13.42,13.42,
1234567890,12/22/2015,12/22/2015,Distribution,LT CAP GAIN    .007,Total Bond Mkt Index Adm,10.66,1.888,20.13,20.13,
1234567890,08/11/2016,08/11/2016,Conversion,ADMIRAL CONVERSION TO 1234-1234567890,Total Stock Mkt Idx Inv,54.55,-188.674,-10292.17,-10292.17,
`
        exp := []bool {
                true,
                true,
                false,
        }

        r := NewVanguardTransactionReader(strings.NewReader(in))
        for i := 0; i < len(exp); i++ {
                trans, err := r.Read()
                if err != nil {
                        t.Errorf("Unexpected error reading Vanguard input: %s", err)
                }
                actual := trans.IsAcquisition
                if actual != exp[i] {
                        t.Errorf("Expected %v, got %v (%v)", exp[i], actual, trans)
                }
        }
}

func TestVanguardTransactionReader_IsIncome(t *testing.T) {
        in := `1234567890,10/30/2015,10/30/2015,Distribution,INCOME DIVIDEND,GNMA Fund Investor Shares,10.7,1.254,13.42,13.42,
1234567890,12/22/2015,12/22/2015,Distribution,LT CAP GAIN    .007,Total Bond Mkt Index Adm,10.66,1.888,20.13,20.13,
1234567890,08/11/2016,08/11/2016,Conversion,ADMIRAL CONVERSION TO 1234-1234567890,Total Stock Mkt Idx Inv,54.55,-188.674,-10292.17,-10292.17,
`
        exp := []bool {
                true,
                true,
                false,
        }

        r := NewVanguardTransactionReader(strings.NewReader(in))
        for i := 0; i < len(exp); i++ {
                trans, err := r.Read()
                if err != nil {
                        t.Errorf("Unexpected error reading Vanguard input: %s", err)
                }
                actual := trans.IsIncome
                if actual != exp[i] {
                        t.Errorf("Expected %v, got %v (%v)", exp[i], actual, trans)
                }
        }
}

func TestVanguardTransactionReader_IsCapGain(t *testing.T) {
        in := `1234567890,10/30/2015,10/30/2015,Distribution,INCOME DIVIDEND,GNMA Fund Investor Shares,10.7,1.254,13.42,13.42,
1234567890,12/22/2015,12/22/2015,Distribution,LT CAP GAIN    .007,Total Bond Mkt Index Adm,10.66,1.888,20.13,20.13,
1234567890,08/11/2016,08/11/2016,Conversion,ADMIRAL CONVERSION TO 1234-1234567890,Total Stock Mkt Idx Inv,54.55,-188.674,-10292.17,-10292.17,
`
        exp := []bool {
                false,
                true,
                false,
        }

        r := NewVanguardTransactionReader(strings.NewReader(in))
        for i := 0; i < len(exp); i++ {
                trans, err := r.Read()
                if err != nil {
                        t.Errorf("Unexpected error reading Vanguard input: %s", err)
                }
                actual := trans.IsCapGain
                if actual != exp[i] {
                        t.Errorf("Expected %v, got %v (%v)", exp[i], actual, trans)
                }
        }
}

func TestVanguardTransactionReader_ReadMulti(t *testing.T) {
        in := `1234567890,10/30/2015,10/30/2015,Distribution,INCOME DIVIDEND,GNMA Fund Investor Shares,10.7,1.254,13.42,13.42,
1234567890,11/30/2016,11/30/2016,Distribution,INCOME DIVIDEND,Total Bond Mkt Index Adm,10.76,5.976,64.3,64.3,
`
        exp := []Transaction{
                Transaction{
                        time.Date(2015, time.October, 30, 0, 0, 0, 0, time.UTC),
                        "1234567890",
                        "GNMA Fund Investor Shares",
                        13.42,
                        true,
                        true,
                        false,
                },
                Transaction{
                        time.Date(2016, time.November, 30, 0, 0, 0, 0, time.UTC),
                        "1234567890",
                        "Total Bond Mkt Index Adm",
                        64.3,
                        true,
                        true,
                        false,
                },
        }

        r := NewVanguardTransactionReader(strings.NewReader(in))
        for i := 0; i < len(exp); i++ {
                trans, err := r.Read()
                if err != nil {
                        t.Errorf("Unexpected error reading Vanguard input: %s", err)
                }
                if trans != exp[i] {
                        t.Errorf("Unexpected input exp: %v, actual: %v", exp[i], trans)
                }
        }
        _, err := r.Read()
        if err != io.EOF {
                t.Errorf("Able to read extra line")
        }
}
