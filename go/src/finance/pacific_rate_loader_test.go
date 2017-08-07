package finance

import (
        "path/filepath"
        "testing"
)

func TestLoadData(t *testing.T) {
        p := PacificRateLoader{filepath.Join("rates", "USD", "CAD")}
        rates, err := p.GetRates()
        if err != nil {
                t.Errorf("Couldn't load: %s", err)
        }
        checkDataForDate(t, 0.0, rates["2016"][0], "2016/01/01")
        checkDataForDate(t, 1.3969, rates["2016"][3], "2016/01/04")
        checkDataForDate(t, 1.3427, rates["2016"][364], "2016/12/30")
}

func checkDataForDate(t *testing.T, exp float32, actual float32, date string) {
        if exp != actual {
                t.Errorf("%s exp=%f, actual=%f", date, exp, actual)
        }
}
