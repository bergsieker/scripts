package main

import (
	"fmt"
	"log"
	"sort"
	"strings"
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
		stats, err := processAsset(k, v, rates)
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

func findNonidenticalLots(ts []transactions.Transaction) ([]transactions.Transaction, [][]transactions.Transaction, error) {
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

	nonidenticalRefs := make(map[string]bool)
	identicals := make([]transactions.Transaction, 0)
	nonidenticals := make([][]transactions.Transaction, 0)
	for i, t := range ts {
		switch t.Type {
		case transactions.TransactionType_Acquisition:
			isIdentical := true
			maxDate := t.Date.Add(time.Hour * 24 * 30)
			for j := i+1; j < len(ts); j++{
				t2 := ts[j]
				if maxDate.Compare(t2.Date) <= 0 {
					break
				}
				if t.LotRef == t2.LotRef && t2.Type == transactions.TransactionType_Disposition {
					if _, found := nonidenticalRefs[t.LotRef]; found {
						return []transactions.Transaction{}, [][]transactions.Transaction{}, fmt.Errorf("found duplicate lotref for nonidentical transaction: %s", t.LotRef)
					}
					nonidenticals = append(nonidenticals, []transactions.Transaction{t, t2})
					nonidenticalRefs[t.LotRef] = true
					isIdentical = false
				}
			}
			if isIdentical {
				identicals = append(identicals, t)
			}

		case transactions.TransactionType_Disposition:
			if _, found := nonidenticalRefs[t.LotRef]; !found {
				identicals = append(identicals, t)
			}

		default:
			identicals = append(identicals, t)
		}
	}
	return identicals, nonidenticals, nil
}

func processAsset(name string, ts []transactions.Transaction, rates rates.Rates) ([]*stats, error) {
	identicals, nonidenticals, err := findNonidenticalLots(ts)
	if err != nil {
		return []*stats{}, err
	}

	fmt.Printf("Processing identical lots for %s\n", name)
	years := make([]*stats, 0)
	var year *stats
	for _, t := range identicals {
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
			year.shares += t.ShareAmount
			incCost := (t.CashAmount + t.Fees) * rate
			year.cost += incCost
			fmt.Printf("%v ACQ %.3f for %.2f CAD ((%.2f + %.2f) * %f), ACB=%.3f (%.2f / %.3f)\n", t.Date.Format("2006-01-02"), t.ShareAmount, incCost, t.CashAmount, t.Fees, rate, year.cost / year.shares, year.cost, year.shares)

		case transactions.TransactionType_Disposition:
			acb := year.cost / year.shares
			incCost := acb * t.ShareAmount
			gains := (t.CashAmount - t.Fees) * rate - incCost
			year.shares -= t.ShareAmount
			year.cost -= incCost
			year.gains += gains
			fmt.Printf("%v DIS gains=%.2f, %.3f @ ((%.2f - %.2f) * %.2f - %.2f * %.3f, ACB=%.3f (%.2f / %.3f)\n", t.Date.Format("2006-01-02"), gains, t.ShareAmount, t.CashAmount, t.Fees, rate, acb, t.ShareAmount, year.cost / year.shares, year.cost, year.shares)

		case transactions.TransactionType_Dividend:
			// TODO: was there a return of capital distribution?
			year.income += t.CashAmount * rate

		default:
			fmt.Println("WARNING: unhandled transaction type")
		}
		if year.cost > year.maxCost {
			year.maxCost = year.cost
		}
	}

	fmt.Printf("Processing nonidentical lots for %s\n", name)
	for _, t := range nonidenticals {
		acq, dis := t[0], t[1]

		if acq.ShareAmount != dis.ShareAmount || acq.LotRef != acq.LotRef {
			return []*stats{}, fmt.Errorf("nonidentical lot does not match")
		}

		var year *stats
		y := fmt.Sprintf("%d", dis.Date.Year())
		yearIndex, found := sort.Find(len(years), func (i int) int {
			return strings.Compare(y, years[i].period)
		})
		if !found {
			if yearIndex != 0 {
				year = &stats{
					period: y,
					cost: years[yearIndex-1].cost,
					shares: years[yearIndex-1].shares,
					maxCost: years[yearIndex-1].cost,
				}
			} else {
				year = &stats{ period: y, }
			}
			years = append(years, year)
			sort.Slice(years, func (i, j int) bool {
				return years[i].period < years[j].period
			})
		} else {
			year = years[yearIndex]
		}

		acqRate, err := rates.RateForDate(acq.Date)
		if err != nil {
			return []*stats{}, err
		}
		disRate, err := rates.RateForDate(dis.Date)
		if err != nil {
			return []*stats{}, err
		}

		gains := (dis.CashAmount - dis.Fees) * disRate - (acq.CashAmount - acq.Fees) * acqRate
		fmt.Printf("%v Non-identical sale: %f gain (dis (%f - %f) * %f, acq (%f - %f) * %f\n", dis.Date, gains, dis.CashAmount, dis.Fees, disRate, acq.CashAmount, acq.Fees, acqRate)
		year.gains += gains
	}

	return years, nil
}
