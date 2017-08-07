package finance

import (
        "fmt"
        "errors"
        "strconv"
        "time"
)

type CurrencyConverter struct {
        Rates map[string][]float32
}

func (c CurrencyConverter) GetRateForDate(d time.Time) (float32, error) {
        yearDay := d.YearDay() - 1
        for year := d.Year(); ; year++ {
		key := strconv.Itoa(year)
		table, ok := c.Rates[key]
		if ok != true {
			return 0.0, fmt.Errorf("Missing data for key %q", key)
		}
                for ; yearDay < len(table); yearDay++ {
                        rate := table[yearDay]
                        if rate > 0 {
                                return rate, nil
                        }
                }
                yearDay = 0
	}
        return 0.0, errors.New("Could not find data")
}
