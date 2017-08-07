package finance

import (
        "io"
        "strings"
        "testing"
        "time"
)

func TestVanguardStatementReader_IsAcquisition(t *testing.T) {
        in := `
,2012-06-24,500 Index Fund Adm,Starting Balance,1,,"62292.96"
,2012-06-29,GNMA Fund Investor Shares,Dividend Received,0,$0.00,16.68
,2016-12-30,Prime Money Mkt Fund,Dividend Received,30.18,$1.00,30.18
`
        exp := []bool {
                true,
                false,
                true,
        }

        r := NewVanguardStatementReader(strings.NewReader(in))
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

func TestVanguardStatementReader_IsIncome(t *testing.T) {
        in := `,2012-06-24,500 Index Fund Adm,Starting Balance,1,,62292.96
,2012-06-29,GNMA Fund Investor Shares,Dividend Received,0,0.00,16.68
,2012-12-28,GNMA Fund Investor Shares,Long-term capital gain,0,0.00,24.57
,2012-12-28,GNMA Fund Investor Shares,Short-term capital gain,0,0.00,31.16
,2012-12-28,Prime Money Mkt Fund,Buy,57.51,1.00,57.51
,2014-02-21,Growth Equity Fund,Transfer,"-2861.47",16.62,-47557.58
`
        exp := []bool {
                false,
                true,
                true,
                true,
                false,
                false,
        }

        r := NewVanguardStatementReader(strings.NewReader(in))
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

func TestVanguardStatementReader_IsCapGain(t *testing.T) {
        in := `,2012-06-24,500 Index Fund Adm,Starting Balance,1,,62292.96
,2012-06-29,GNMA Fund Investor Shares,Dividend Received,0,0.00,16.68
,2012-12-28,GNMA Fund Investor Shares,Long-term capital gain,0,0.00,24.57
,2012-12-28,GNMA Fund Investor Shares,Short-term capital gain,0,0.00,31.16
,2012-12-28,Prime Money Mkt Fund,Buy,57.51,1.00,57.51
,2014-02-21,Growth Equity Fund,Transfer,"-2861.47",16.62,-47557.58
`
        exp := []bool {
                false,
                false,
                false,
                false,
                false,
                false,
        }

        r := NewVanguardStatementReader(strings.NewReader(in))
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

func TestVanguardStatementReader_ReadMulti(t *testing.T) {
        in := `,2012-06-24,500 Index Fund Adm,Starting Balance,1,,62292.96
,2012-06-29,GNMA Fund Investor Shares,Dividend Received,0,0.00,16.68
`
        exp := []Transaction{
                Transaction{
                        time.Date(2012, time.June, 24, 0, 0, 0, 0, time.UTC),
                        "",
                        "500 Index Fund Adm",
                        62292.96,
                        true,
                        false,
                        false,
                },
                Transaction{
                        time.Date(2012, time.June, 29, 0, 0, 0, 0, time.UTC),
                        "",
                        "GNMA Fund Investor Shares",
                        16.68,
                        false,
                        true,
                        false,
                },
        }

        r := NewVanguardStatementReader(strings.NewReader(in))
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
