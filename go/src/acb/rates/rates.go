package rates

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"time"
)

type Rates struct {
	rates []dailyRate
}

type dailyRate struct {
	date time.Time
	rate float32
}

func (r *Rates) LoadRates(path string) (error) {
	r.rates = nil
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	csvr := csv.NewReader(f)
	if err != nil {
		return err
	}
	firstRow := true
	for {
		record, err := csvr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if firstRow {
			firstRow = false
			continue
		}

		date, err := time.Parse("2006-01-02", record[0])
		if err != nil {
			return err
		}
		rate, err := strconv.ParseFloat(record[1], 32)
		if err != nil {
			return err
		}
		r.rates = append(r.rates, dailyRate{
			date: date,
			rate: float32(rate),
		})
	}
	return nil
}

func (r *Rates) RateForDate(t time.Time) (float32, error) {
	if t.Compare(r.rates[0].date) < 0 || t.Compare(r.rates[len(r.rates)-1].date) > 0 {
		return 0.0, fmt.Errorf("exchange rate out of range")
	}
	index, found := sort.Find(len(r.rates), func (i int) int {
		return t.Compare(r.rates[i].date)
	})
	if index >= len(r.rates) {
		return 0.0, fmt.Errorf("exchange rate out of range")
	}
	if found {
		return r.rates[index].rate, nil
	}
	return r.rates[index-1].rate, nil
}
