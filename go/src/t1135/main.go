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

var aliases = map[string]string{
        "Growth Equity Fund" : "U.S. Growth Fund Investor",
}

var (
        ratesRoot = flag.String("rates_root", "", "")
        vanguardFile = flag.String("vanguard_file", "", "")
        endingYear = flag.Int("ending_year", 2012, "The final year to load data for")
        firstPrintYear = flag.Int("first_print_year", 2012, "The first year to print in the output")
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
        vtr := finance.NewVanguardTransactionListReader(vf)
        // vtr := finance.NewITradeTransactionListParser(vf)
        for {
                t, err := vtr.Read()
                if err == io.EOF {
                        break
                }
                if err != nil {
                        log.Fatal(err)
                }
                if t.Type == finance.TransactionType_Transfer {
                        // Skip transfers because none of the transfers I've
                        // done have been relevant to the t1135.
                        continue
                }
                name := assetNameForTransaction(t)
                a, ok := usdAssets[name]
                if !ok {
                        a = finance.CreateAsset(startingYear, *endingYear)
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
                        a = finance.CreateAsset(startingYear, *endingYear)
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
        for _, name := range names {
                // usdAsset := usdAssets[name]
                // printAssetTable(w, name, usdAsset, "USD", *firstPrintYear, *endingYear)
                cadAsset := cadAssets[name]
                printAssetTable(w, name, cadAsset, "CAD", *firstPrintYear, *endingYear)
        }
}

func printAssetTable(w *tabwriter.Writer, name string, a finance.Asset, cur string, startingYear int, endingYear int) {
        d := a.DumpDataForYears(startingYear, endingYear)
        fmt.Println(name, cur)
        fmt.Fprintln(w, "Year\tMax Cost\tCurrent Cost\tIncome\tCap Gains\t")
        for i, y := range d {
                fmt.Fprintf(w, "%d\t%0.0f\t%0.0f\t%0.0f\t%0.0f\t\n", startingYear + i, y.MaxDelta, y.CurrentDelta, y.Income, y.CapGains)
        }
        w.Flush()
        fmt.Println()
}

func assetNameForTransaction(t finance.Transaction) string {
        alias, ok := aliases[t.AssetName]
        if ok {
                return alias
        }
        return t.AssetName
}
