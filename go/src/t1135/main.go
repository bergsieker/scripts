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
	"Growth Equity Fund":    "U.S. Growth Fund Investor",
	"FTSE Social Index Inv": "FTSE Social Index",
	"FTSE Social Index Adm": "FTSE Social Index",
}

var (
	ratesRoot       = flag.String("rates_root", "", "")
	vanguardFile    = flag.String("vanguard_file", "", "")
	vanguardOfxFile = flag.String("vanguard_ofx_file", "", "")
	endingYear      = flag.Int("ending_year", 2012, "The final year to load data for")
	firstPrintYear  = flag.Int("first_print_year", 2012, "The first year to print in the output")
)

type TransactionReader interface {
	Read() (finance.Transaction, error)
}

func main() {
	flag.Parse()

	pr := finance.PacificRateLoader{*ratesRoot}
	rates, err := pr.GetRates()
	if err != nil {
		log.Fatal(err)
	}
	c := finance.CurrencyConverter{rates}

	usdAssets := make(map[string]finance.Asset)
	cadAssets := make(map[string]finance.Asset)

	vf, err := os.Open(*vanguardFile)
	if err != nil {
		log.Fatal(err)
	}
	defer vf.Close()
	var r TransactionReader
	r = finance.NewVanguardTransactionListReader(vf)
	err = processTransactions(r, c, usdAssets, cadAssets)
	if err != nil {
		log.Fatal(err)
	}

	vofxf, err := os.Open(*vanguardOfxFile)
	if err != nil {
		log.Fatal(err)
	}
	defer vofxf.Close()
	r = finance.NewVanguardOfxTransactionListReader(vofxf, "09936218284")
	err = processTransactions(r, c, usdAssets, cadAssets)
	if err != nil {
		log.Fatal(err)
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

func processTransactions(tr TransactionReader, c finance.CurrencyConverter, usdAssets map[string]finance.Asset, cadAssets map[string]finance.Asset) error {
	for {
		t, err := tr.Read()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
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
			return err
		}
		t.Amount = t.Amount * rate
		a, ok = cadAssets[name]
		if !ok {
			a = finance.CreateAsset(startingYear, *endingYear)
			cadAssets[name] = a
		}
		a.ProcessTransaction(t)
	}
}

func printAssetTable(w *tabwriter.Writer, name string, a finance.Asset, cur string, startingYear int, endingYear int) {
	d := a.DumpDataForYears(startingYear, endingYear)
	fmt.Println(name, cur)
	fmt.Fprintln(w, "Year\tMax Cost\tCurrent Cost\tIncome\tCap Gains\t")
	for i, y := range d {
		fmt.Fprintf(w, "%d\t%0.0f\t%0.0f\t%0.0f\t%0.0f\t\n", startingYear+i, y.MaxDelta, y.CurrentDelta, y.Income, y.CapGains)
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
