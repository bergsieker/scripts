package finance

import (
        "testing"
        "time"
)

func TestGetRateForDate(t *testing.T) {
        c := CurrencyConverter{
                map[string][]float32{
                        "2016": []float32{1.5, 0.0, 0.0, 3.0, 0.0},
                        "2017": []float32{0.0, 4.5},
                },
        }

        d := time.Date(2015, time.January, 1, 0, 0, 0, 0, time.UTC)
        checkRateForDate(t, c, d, 0.0, true)

        d = time.Date(2016, time.January, 1, 0, 0, 0, 0, time.UTC)
        checkRateForDate(t, c, d, 1.5, false)

        d = time.Date(2016, time.January, 2, 0, 0, 0, 0, time.UTC)
        checkRateForDate(t, c, d, 3.0, false)

        d = time.Date(2016, time.January, 4, 0, 0, 0, 0, time.UTC)
        checkRateForDate(t, c, d, 3.0, false)

        d = time.Date(2016, time.January, 5, 0, 0, 0, 0, time.UTC)
        checkRateForDate(t, c, d, 4.5, false)

        d = time.Date(2016, time.January, 6, 0, 0, 0, 0, time.UTC)
        checkRateForDate(t, c, d, 4.5, false)

        d = time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC)
        checkRateForDate(t, c, d, 4.5, false)

        d = time.Date(2017, time.January, 7, 0, 0, 0, 0, time.UTC)
        checkRateForDate(t, c, d, 0.0, true)
}

func checkRateForDate(t *testing.T, c CurrencyConverter, d time.Time, r float32, e bool) {
        rate, err := c.GetRateForDate(d)
        if !e && err != nil {
                t.Errorf("Unexpected error on %s, %s", d, err)
        }
        if e && err == nil {
                t.Errorf("Unexpected non-error on %s", d)
        }
        if rate != r {
                t.Errorf("Wrong rate on %s: expected %f, got %f", d, r, rate)
        }
}
