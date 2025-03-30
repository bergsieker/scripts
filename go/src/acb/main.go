package main

import (
	"fmt"
	"log"
	"sort"
	"time"

	"acb/mssbactivityreport"
	"acb/rates"
	"acb/transactions"
)

type acb struct {
	cost float32
	shares float32
}

type stats struct {
	period string
	cost float32
	maxCost float32
	shares float32
	income float32
	gains float32
}

func main() {
	rates := rates.Rates{}
	if err := rates.LoadRates("rates/usdcadrates.csv"); err != nil {
		log.Fatalf("error loading rates: %v", err)
	}

	ar, err := mssbactivityreport.Parse("/Users/sbergsieker/Downloads/Activity Report - Bergsieker, Steven")
	if err != nil {
		log.Fatalf("%s", err)
	}
	ts, err := ar.ToTransactions()
	if err != nil {
		log.Fatalf("%s", err)
	}
	byAsset := make(map[string][]transactions.Transaction)
	for _, t := range ts {
		a := byAsset[t.AssetName]
		byAsset[t.AssetName] = append(a, t)
	}
	statsByAsset := make(map[string][]*stats)
	for k, v := range byAsset {
		fmt.Println("Processing transactions for asset " + k)
		stats, err := processAsset(v, rates)
		if err != nil {
			log.Fatalf("error processing %s: %v", k, err)
		}
		statsByAsset[k] = stats
	}

	for k, stats := range statsByAsset {
		fmt.Println(k)
		for _, period := range stats {
			fmt.Printf("%s YE Shares=%0.2f YE Cost=%0.2f MaxCost=%0.2f Income=%0.2f Gain=%0.2f\n", period.period, period.shares, period.cost, period.maxCost, period.income, period.gains)
		}
	}
}

func findNonidenticalLots(ts []transactions.Transaction) (map[string][]transactions.Transaction, error) {
	nonidenticalRefs := make(map[string]bool)
	nonidenticals := make(map[string][]transactions.Transaction)
	for i, t := range ts {
		switch t.Type {
		case transactions.TransactionType_Acquisition:
			maxDate := t.Date.Add(time.Hour * 24 * 30)
			for j := i+1; j < len(ts); j++{
				t2 := ts[j]
				if maxDate.Compare(t2.Date) <= 0 {
					break
				}
				if t.LotRef == t2.LotRef && t2.Type == transactions.TransactionType_Disposition {
					if _, found := nonidenticalRefs[t.LotRef]; found {
						return make(map[string][]transactions.Transaction), fmt.Errorf("found duplicate lotref for nonidentical transaction: %s", t.LotRef)
					}
					nonidenticals[t.LotRef] = []transactions.Transaction{t, t2}
					nonidenticalRefs[t.LotRef] = true
				}
			}
		}
	}
	return nonidenticals, nil
}

func processAsset(ts []transactions.Transaction, rates rates.Rates) ([]*stats, error) {
	// Make sure that input is sorted by date, with acquisitions preceeding dispositions.
	sort.Slice(ts, func(i, j int) bool {
			dc := ts[i].Date.Compare(ts[j].Date)
			if dc < 0 {
				return true
			} else if dc > 0 {
				return false
			}
			if ts[i].Type == ts[j].Type {
				// For transactions on the same day with the same type, just sort by ref number.
				return ts[i].LotRef < ts[j].LotRef
			}
			// Sort by type, so acquisitions come before dispositions.
			return ts[i].Type < ts[j].Type
	})

	nonidenticals, err := findNonidenticalLots(ts)
	if err != nil {
		return []*stats{}, err
	}

	years := make([]*stats, 0)
	nonidenticalCost := float32(0)
	var year *stats
	for _, t := range ts {
		rate, err := rates.RateForDate(t.Date)
		if err != nil {
			return []*stats{}, err
		}
		y := fmt.Sprintf("%d", t.Date.Year())
		if len(years) == 0 || y != years[len(years)-1].period {
			if len(years) == 0 {
				years = append(years, &stats{
					period: y,
				})
			} else {
				years = append(years, &stats{
					period: y,
					cost: year.cost,
					maxCost: year.cost,
					shares: year.shares,
				})
			}
			year = years[len(years)-1]
		}
		switch t.Type {
		case transactions.TransactionType_Acquisition:
			if _, found := nonidenticals[t.LotRef]; found {
				incCost := (t.CashAmount + t.Fees) * rate
				nonidenticalCost += incCost
			} else {
				year.shares += t.ShareAmount
				incCost := (t.CashAmount + t.Fees) * rate
				year.cost += incCost
				fmt.Printf("%v ACQ %.3f for %.2f CAD ((%.2f + %.2f) * %f), ACB=%.3f (%.2f / %.3f)\n", t.Date.Format("2006-01-02"), t.ShareAmount, incCost, t.CashAmount, t.Fees, rate, year.cost / year.shares, year.cost, year.shares)
			}

		case transactions.TransactionType_Disposition:
			if ni, found := nonidenticals[t.LotRef]; found {
				acq, dis := ni[0], ni[1]
				disRate := rate
				acqRate, err := rates.RateForDate(acq.Date)
				if err != nil {
					return []*stats{}, err
				}
				cost := (acq.CashAmount + acq.Fees) * acqRate
				proceeds := (dis.CashAmount - dis.Fees) * disRate
				gains := proceeds - cost
				fmt.Printf("%v Non-identical sale: %.2f gain (dis (%.2f - %.2f) * %.2f, acq (%.2f + %.2f) * %.2f\n", dis.Date.Format("2006-01-02"), gains, dis.CashAmount, dis.Fees, disRate, acq.CashAmount, acq.Fees, acqRate)
				year.gains += gains
				nonidenticalCost -= cost
			} else {
				acb := year.cost / year.shares
				incCost := acb * t.ShareAmount
				gains := (t.CashAmount - t.Fees) * rate - incCost
				year.shares -= t.ShareAmount
				year.cost -= incCost
				year.gains += gains
				fmt.Printf("%v DIS gains=%.2f, %.3f @ ((%.2f - %.2f) * %.2f - %.2f * %.3f, ACB=%.3f (%.2f / %.3f)\n", t.Date.Format("2006-01-02"), gains, t.ShareAmount, t.CashAmount, t.Fees, rate, acb, t.ShareAmount, year.cost / year.shares, year.cost, year.shares)
			}

		case transactions.TransactionType_Dividend:
			// TODO: was there a return of capital distribution?
			year.income += t.CashAmount * rate

		default:
			fmt.Println("WARNING: unhandled transaction type")
		}
		curCost := year.cost + nonidenticalCost
		if curCost > year.maxCost {
			year.maxCost = curCost
		}
	}

	if nonidenticalCost < -0.1 || nonidenticalCost > 0.1 {
		fmt.Println("WARNING: ended with non-zero cost for nonidentical shares: " + fmt.Sprintf("%.2f", nonidenticalCost))
	}

	return years, nil
}
