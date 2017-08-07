package main

import (
        "finance"
        "flag"
        "fmt"
        "io"
        "log"
        "os"
        "sort"
        "text/tabwriter"
)

const startingYear = 2012
const endingYear = 2016

var aliases = map[string]string{
        "Growth Equity Fund" : "U.S. Growth Fund Investor",
}

var (
        ratesRoot = flag.String("rates_root", "", "")
        vanguardFile = flag.String("vanguard_file", "", "")
)

func main() {
        flag.Parse()

        pr := finance.PacificRateLoader{*ratesRoot}
        r, err := pr.GetRates()
        if err != nil {
                log.Fatal(err)
        }
        c := finance.CurrencyConverter{r}

        usdAssets := make(map[string]finance.Asset)
        cadAssets := make(map[string]finance.Asset)

        vf, err := os.Open(*vanguardFile)
        if err != nil {
                log.Fatal(err)
        }
        vtr := finance.NewVanguardStatementReader(vf)
        for {
                t, err := vtr.Read()
                if err == io.EOF {
                        break
                }
                if err != nil {
                        log.Fatal(err)
                }
                name := assetNameForTransaction(t)
                a, ok := usdAssets[name]
                if !ok {
                        a = finance.CreateAsset(startingYear, endingYear)
                        usdAssets[name] = a
                }
                a.ProcessTransaction(t)

                rate, err := c.GetRateForDate(t.Date)
                if err != nil {
                        log.Fatal(err)
                }
                t.Amount = t.Amount * rate
                a, ok = cadAssets[name]
                if !ok {
                        a = finance.CreateAsset(startingYear, endingYear)
                        cadAssets[name] = a
                }
                a.ProcessTransaction(t)
        }

        names := []string{}
        for name, _ := range usdAssets {
                names = append(names, name)
        }
        sort.Strings(names)

        w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.AlignRight)
        printStartingYear := 2016
        for _, name := range names {
                /*
                a := usdAssets[name]
                d := a.DumpDataForYears(printStartingYear, endingYear)
                fmt.Println(name, "USD")
                fmt.Fprintln(w, "Year\tMax Cost\tCurrent Cost\tIncome\tCap Gains\t")
                for i, y := range d {
                        fmt.Fprintf(w, "%d\t%.2f\t%.2f\t%.2f\t%.2f\t\n", printStartingYear + i, y.MaxDelta, y.CurrentDelta, y.Income, y.CapGains)
                }
                w.Flush()
                fmt.Println()
                */

                a := cadAssets[name]
                d := a.DumpDataForYears(printStartingYear, endingYear)
                fmt.Println(name)
                fmt.Fprintln(w, "Year\tMax Cost\tCurrent Cost\tIncome\tCap Gains\t")
                for i, y := range d {
                        fmt.Fprintf(w, "%d\t%.2f\t%.2f\t%.2f\t%.2f\t\n", printStartingYear + i, y.MaxDelta, y.CurrentDelta, y.Income, y.CapGains)
                }
                w.Flush()
                fmt.Println()
        }
}

func assetNameForTransaction(t finance.Transaction) string {
        alias, ok := aliases[t.AssetName]
        if ok {
                return alias
        }
        return t.AssetName
}
