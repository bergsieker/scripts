package mssbactivityreport

import (
        "encoding/csv"
	"fmt"
	"path/filepath"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
	"acb/transactions"
)

type ActivityReport struct {
	AutosaleReleases []AutosaleRelease
	Releases []Release
	NetShares []NetShares
	WithdrawalWires []WithdrawalWire
	Withdrawals []Withdrawal
}

type AutosaleRelease struct {
	ReleaseDate time.Time
	SaleDate time.Time
	OrderNumber string
	Plan string
	ReleaseType string
	Status string
	SalePrice float32
	Quantity float32
	GrossProceeds float32
	NetProceeds float32
}

type Release struct {
	VestDate time.Time
	OrderNumber string
	Plan string
	ReleaseType string
	Status string
	Price float32
	Quantity float32
	NetCash float32
	NetShares float32
	TaxPaymentMethod string
}

type WithdrawalWire struct {
	Date time.Time
	OrderNumber string
	Symbol string
	WithdrawalType string
	Status string
	Price float32
	Quantity float32
	NetCash float32
	NetShares float32
	TaxPaymentMethod string
}

type Withdrawal struct {
	Date time.Time
	OrderNumber string
	Plan string
	WithdrawalType string
	Status string
	Price float32
	Quantity float32
	NetCash float32
	NetShares float32
	TaxPaymentMethod string
}

type NetShares struct {
	Date time.Time
	OrderNumber string
	Plan string
	ReleaseType string
	Status string
	Price float32
	Quantity float32
	NetSharesProceeds1 float32
	NetSharesProceeds2 float32
	TaxPaymentMethod string
}

func (ar *ActivityReport)ToTransactions() ([]transactions.Transaction, error) {
	ts := make([]transactions.Transaction, 0)
	for _, rel := range ar.Releases {
		if rel.ReleaseType != "Release" {
			return []transactions.Transaction{}, fmt.Errorf("unexpected release type %s", rel.ReleaseType)
		}
		if rel.Status != "Complete" {
			return []transactions.Transaction{}, fmt.Errorf("unexpected status %s", rel.Status)
		}
		shares := rel.NetShares
		if shares <= 0 {
			shares = rel.Quantity
		}
		t := transactions.Transaction{
			Date: rel.VestDate,
			AssetName: mapSymbol(rel.Plan),
			Type: transactions.TransactionType_Acquisition,
			LotRef: rel.OrderNumber,
			CashAmount: rel.Price * shares,
			ShareAmount: shares,
			Fees: 0,
		}
		ts = append(ts, t)
	}

	for _, rel := range ar.AutosaleReleases {
		if rel.ReleaseType != "Release" {
			return []transactions.Transaction{}, fmt.Errorf("unexpected release type %s", rel.ReleaseType)
		}
		if rel.Status != "Complete" {
			return []transactions.Transaction{}, fmt.Errorf("unexpected status %s", rel.Status)
		}
		t := transactions.Transaction{
			Date: rel.SaleDate,
			AssetName: mapSymbol(rel.Plan),
			Type: transactions.TransactionType_Disposition,
			LotRef: rel.OrderNumber,
			CashAmount: rel.GrossProceeds,
			ShareAmount: rel.Quantity,
			Fees: rel.GrossProceeds - rel.NetProceeds,
		}
		ts = append(ts, t)
	}

	for _, rel := range ar.NetShares {
		if rel.ReleaseType != "Released Shares" {
			return []transactions.Transaction{}, fmt.Errorf("unexpected release type %s", rel.ReleaseType)
		}
		if rel.Status != "Completed" {
			return []transactions.Transaction{}, fmt.Errorf("unexpected status %s", rel.Status)
		}
		if rel.NetSharesProceeds1 > 0 || rel.NetSharesProceeds2 > 0 {
			return []transactions.Transaction{}, fmt.Errorf("unexpected net shares proceeds: %s", rel)
		}
	}

	for _, w := range ar.Withdrawals {
		if w.Status != "Complete" {
			return []transactions.Transaction{}, fmt.Errorf("unexpected status %s", w.Status)
		}
		if w.NetShares != 0 {
			return []transactions.Transaction{}, fmt.Errorf("unexpected net shares")
		}
		if w.Plan == "Cash" {
			continue
		}
		switch w.WithdrawalType {
		case "Sale":
			t := transactions.Transaction{
				Date: w.Date,
				AssetName: mapSymbol(w.Plan),
				LotRef: w.OrderNumber,
				Type: transactions.TransactionType_Disposition,
				CashAmount: w.Price * -1 * w.Quantity,
				ShareAmount: -1 * w.Quantity,
				Fees: w.Price * -1 * w.Quantity - w.NetCash,
			}
			ts = append(ts, t)
		case "":
			// A blank type indicates a transfer in-kind (usually for a charitable donation).
			t := transactions.Transaction{
				Date: w.Date,
				AssetName: mapSymbol(w.Plan),
				LotRef: w.OrderNumber,
				Type: transactions.TransactionType_Disposition,
				CashAmount: w.Price * -1 * w.Quantity,
				ShareAmount: -1 * w.Quantity,
				Fees: 0,
			}
			ts = append(ts, t)
		case "Dividend":
			t := transactions.Transaction{
				Date: w.Date,
				AssetName: mapSymbol(w.Plan),
				LotRef: w.OrderNumber,
				Type: transactions.TransactionType_Dividend,
				CashAmount: w.Quantity * w.Price,
				ShareAmount: 0,
				Fees: 0,
			}
			ts = append(ts, t)
		case "Transfer":
			if mapSymbol(w.Plan) == "GOOGL" && w.Date.Format("2006-01-02") == "2014-04-09" {
				// This is how the class A/C split appears in the activity report.
				// We convert this to an initial tranche of class C stock, acquired at 0 cost.
				// That may sound weird, but it's how the CRA handles this particular split.
				t := transactions.Transaction{
					Date: w.Date,
					AssetName: "GOOG",
					LotRef: w.OrderNumber,
					Type: transactions.TransactionType_Acquisition,
					CashAmount: 0.0,
					ShareAmount: w.Quantity * -1,
					Fees: 0.0,
				}
				ts = append(ts, t)
			} else {
				return []transactions.Transaction{}, fmt.Errorf("unsupported type %s", w.WithdrawalType)
			}
		}
	}

	return ts, nil
}

func mapSymbol(s string) string {
	switch (s) {
	case "GSU Class C":
		return "GOOG"
	case "GSU Class A", "GSU":
		return "GOOGL"
	}
	fmt.Println("Could not map to symbol: " + s)
	return ""
}

func Parse(path string) (*ActivityReport, error) {
	autosaleReleaseFile, err := os.Open(filepath.Join(path, "Autosale_Trading Plan - Releases.csv"))
	if err != nil {
		return nil, err
	}
	ar, err := ParseAutosaleReleases(autosaleReleaseFile)
	if err != nil {
		return nil, err
	}

	releasesFile, err := os.Open(filepath.Join(path, "Releases Report.csv"))
	if err != nil {
		return nil, err
	}
	r, err := ParseReleases(releasesFile)
	if err != nil {
		return nil, err
	}

	netSharesFile, err := os.Open(filepath.Join(path, "Releases Net Shares Report.csv"))
	if err != nil {
		return nil, err
	}
	ns, err := ParseNetShares(netSharesFile)
	if err != nil {
		return nil, err
	}

	withdrawalWiresFile, err := os.Open(filepath.Join(path, "Withdrawal Wire Report.csv"))
	if err != nil {
		return nil, err
	}
	wws, err := ParseWithdrawalWires(withdrawalWiresFile)
	if err != nil {
		return nil, err
	}

	withdrawalsFile, err := os.Open(filepath.Join(path, "Withdrawals Report.csv"))
	if err != nil {
		return nil, err
	}
	ws, err := ParseWithdrawals(withdrawalsFile)
	if err != nil {
		return nil, err
	}

	return &ActivityReport{
		AutosaleReleases: *ar,
		Releases: *r,
		NetShares: *ns,
		WithdrawalWires: *wws,
		Withdrawals: *ws,
	}, nil
}

func parseFloat(s string) (float32, error) {
	s = strings.ReplaceAll(s, "$", "")
	s = strings.ReplaceAll(s, ",", "")
	if len(s) > 0 {
		f, err := strconv.ParseFloat(s, 32)
		return float32(f), err
	}
	return 0, nil
}

func ParseNetShares(r io.Reader) (*[]NetShares, error) {
	nsr := csv.NewReader(r)
	nsr.FieldsPerRecord = 10
	ns := make([]NetShares, 0)
	firstRow := true
	for {
		record, err := nsr.Read()
		if err == io.EOF {
			return &ns, nil
		}
		if err != nil {
			return nil, err
		}
		if firstRow {
			firstRow = false
			hr := csv.NewReader(strings.NewReader("Date,Order Number,Plan,Type,Order Status,Price,Quantity,Net Share Proceeds,Net Share Proceeds,Tax Payment Method"))
			head, err := hr.Read()
			if err != nil {
				return nil, err
			}
			if !slices.Equal(record, head) {
				return nil, fmt.Errorf("header row not in expected format: %s", record)
			}
			continue
		}
		date, err := time.Parse("02-Jan-2006", record[0])
		if err != nil {
			return nil, err
		}
		price, err := parseFloat(record[5])
		if err != nil {
			return nil, err
		}
		quantity, err := parseFloat(record[6])
		if err != nil {
			return nil, err
		}
		nsp1, err := parseFloat(record[7])
		if err != nil {
			return nil, err
		}
		nsp2, err := parseFloat(record[8])
		if err != nil {
			return nil, err
		}
		ns = append(ns, NetShares{
			Date: date,
			OrderNumber: record[1],
			Plan: record[2],
			ReleaseType: record[3],
			Status: record [4],
			Price: price,
			Quantity: quantity,
			NetSharesProceeds1: nsp1,
			NetSharesProceeds2: nsp2,
			TaxPaymentMethod: record[9],
		})

	}
}

func ParseReleases(r io.Reader) (*[]Release, error) {
	rr := csv.NewReader(r)
	rr.FieldsPerRecord = 10
	rel := make([]Release, 0)
	firstRow := true
	for {
		record, err := rr.Read()
		if err == io.EOF {
			return &rel, nil
		}
		if err != nil {
			return nil, err
		}
		if firstRow {
			firstRow = false
			hr := csv.NewReader(strings.NewReader("Vest Date,Order Number,Plan,Type,Status,Price,Quantity,Net Cash Proceeds,Net Share Proceeds,Tax Payment Method"))
			head, err := hr.Read()
			if err != nil {
				return nil, err
			}
			if !slices.Equal(record, head) {
				return nil, fmt.Errorf("header row not in expected format: %s", record)
			}
			continue
		}
		vDate, err := time.Parse("02-Jan-2006", record[0])
		if err != nil {
			return nil, err
		}
		orderNumber := record[1]
		plan := record[2]
		releaseType := record[3]
		status := record[4]
		price, err := parseFloat(record[5])
		if err != nil {
			return nil, err
		}
		quantity, err := parseFloat(record[6])
		if err != nil {
			return nil, err
		}
		netCash, err := parseFloat(record[7])
		if err != nil {
			return nil, err
		}
		netShares, err := parseFloat(record[8])
		if err != nil {
			return nil, err
		}
		taxPaymentMethod := record[9]
		rel = append(rel, Release{
			VestDate: vDate,
			OrderNumber: orderNumber,
			Plan: plan,
			ReleaseType: releaseType,
			Status: status,
			Price: price,
			Quantity: quantity,
			NetCash: netCash,
			NetShares: netShares,
			TaxPaymentMethod: taxPaymentMethod,
		})

	}
}

func ParseAutosaleReleases(r io.Reader) (*[]AutosaleRelease, error) {
	cr := csv.NewReader(r)
	cr.FieldsPerRecord = 10
	ar := make([]AutosaleRelease, 0)
	firstRow := true
	for {
		record, err := cr.Read()
		if err == io.EOF {
			return &ar, nil
		}
		if err != nil {
			return nil, err
		}
		if firstRow {
			firstRow = false
			hr := csv.NewReader(strings.NewReader("Release Date,Sale Date,Order Number,Plan,Type,Status,Sale Price,Quantity,Gross Proceeds,Net Proceeds After Fees"))
			head, err := hr.Read()
			if err != nil {
				return nil, err
			}
			if !slices.Equal(record, head) {
				return nil, fmt.Errorf("header row not in expected format: %s", record)
			}
			continue
		}
		rDate, err := time.Parse("02-Jan-2006", record[0])
		if err != nil {
			return nil, err
		}
		sDate, err := time.Parse("02-Jan-2006", record[1])
		if err != nil {
			return nil, err
		}
		orderNumber := record[2]
		plan := record[3]
		releaseType := record[4]
		status := record[5]
		salePrice, err := parseFloat(record[6])
		if err != nil {
			return nil, err
		}
		quantity, err := parseFloat(record[7])
		if err != nil {
			return nil, err
		}
		gross, err := parseFloat(record[8])
		if err != nil {
			return nil, err
		}
		net, err := parseFloat(record[9])
		if err != nil {
			return nil, err
		}
		ar = append(ar, AutosaleRelease{
			ReleaseDate: rDate,
			SaleDate: sDate,
			OrderNumber: orderNumber,
			Plan: plan,
			ReleaseType: releaseType,
			Status: status,
			SalePrice: salePrice,
			Quantity: quantity,
			GrossProceeds: gross,
			NetProceeds: net,
		})
	}
}

func ParseWithdrawalWires(r io.Reader) (*[]WithdrawalWire, error) {
	wr := csv.NewReader(r)
	wr.FieldsPerRecord = 10
	ws := make([]WithdrawalWire, 0)
	firstRow := true
	for {
		record, err := wr.Read()
		if err == io.EOF {
			return &ws, nil
		}
		if err != nil {
			return nil, err
		}
		if firstRow {
			firstRow = false
			hr := csv.NewReader(strings.NewReader("Date,Order Number,Symbol,Type,Order Status,Price,Quantity,Net Cash Proceeds,Net Share Proceeds,Tax Payment Method"))
			head, err := hr.Read()
			if err != nil {
				return nil, err
			}
			if !slices.Equal(record, head) {
				return nil, fmt.Errorf("header row not in expected format: %s", record)
			}
			continue
		}
		date, err := time.Parse("02-Jan-2006", record[0])
		if err != nil {
			return nil, err
		}
		price, err := parseFloat(record[5])
		if err != nil {
			return nil, err
		}
		quantity, err := parseFloat(record[6])
		if err != nil {
			return nil, err
		}
		netCash, err := parseFloat(record[7])
		if err != nil {
			return nil, err
		}
		netShares, err := parseFloat(record[8])
		if err != nil {
			return nil, err
		}
		ws = append(ws, WithdrawalWire{
			Date: date,
			OrderNumber: record[1],
			Symbol: record[2],
			WithdrawalType: record[3],
			Status: record[4],
			Price: price,
			Quantity: quantity,
			NetCash: netCash,
			NetShares: netShares,
			TaxPaymentMethod: record[9],
		})
	}
}

func ParseWithdrawals(r io.Reader) (*[]Withdrawal, error) {
	wr := csv.NewReader(r)
	wr.FieldsPerRecord = -1
	ws := make([]Withdrawal, 0)
	firstRow := true
	for {
		record, err := wr.Read()
		if err == io.EOF {
			return &ws, nil
		}
		if err != nil {
			return nil, err
		}
		if firstRow {
			firstRow = false
			hr := csv.NewReader(strings.NewReader("Execution Date,Order Number,Plan,Type,Order Status,Price,Quantity,Net Amount,Net Share Proceeds,Tax Payment Method"))
			head, err := hr.Read()
			if err != nil {
				return nil, err
			}
			if !slices.Equal(record, head) {
				return nil, fmt.Errorf("header row not in expected format: %s", record)
			}
			continue
		}
		if strings.HasPrefix(record[0], "Please note") {
			continue
		}
		if len(record) != 10 {
			return nil, fmt.Errorf("csv entry has %d fields, expected 10", len(record))
		}
		date, err := time.Parse("02-Jan-2006", record[0])
		if err != nil {
			return nil, err
		}
		preSplitValues := date.Compare(time.Date(2022, 7, 16, 0, 0, 0, 0, time.Local)) < 0
		price, err := parseFloat(record[5])
		if err != nil {
			return nil, err
		}
		if preSplitValues {
			price /= 20
		}
		quantity, err := parseFloat(record[6])
		if err != nil {
			return nil, err
		}
		if preSplitValues {
			quantity *= 20
		}
		netCash, err := parseFloat(record[7])
		if err != nil {
			return nil, err
		}
		netShares, err := parseFloat(record[8])
		if err != nil {
			return nil, err
		}
		ws = append(ws, Withdrawal{
			Date: date,
			OrderNumber: record[1],
			Plan: record[2],
			WithdrawalType: record[3],
			Status: record[4],
			Price: price,
			Quantity: quantity,
			NetCash: netCash,
			NetShares: netShares,
			TaxPaymentMethod: record[9],
		})
	}
}
