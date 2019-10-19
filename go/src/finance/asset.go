package finance

type AssetYear struct {
        MaxDelta float32
        CurrentDelta float32
        Income float32
        CapGains float32
}

type Asset struct {
        startingYear int
        startingCost float32
        years []AssetYear
}

func (a *AssetYear)ProcessTransaction(t Transaction) {
        if t.IsAcquisition() {
                a.CurrentDelta += t.Amount
                if a.CurrentDelta > a.MaxDelta {
                        a.MaxDelta = a.CurrentDelta
                }
        }
        if t.IsIncome() {
                a.Income += t.Amount
        }
        if t.IsDisposition() {
                a.CapGains += t.Amount
        }
}

func CreateAsset(startingYear int, endingYear int) (Asset) {
        totalYears := endingYear - startingYear + 1
        years := make([]AssetYear, totalYears)
        return Asset{startingYear, 0.0, years}
}

func (a *Asset)ProcessTransaction(t Transaction) {
        a.years[a.getIndexForYear(t.Date.Year())].ProcessTransaction(t)
}

func (a Asset)DumpDataForYears(start int, end int) []AssetYear {
        var data []AssetYear
        val := a.startingCost
        for y := a.startingYear; y <= end; y++ {
                ay := a.years[a.getIndexForYear(y)]
                endVal := val + ay.CurrentDelta
                if y >= start {
                        d := AssetYear{val + ay.MaxDelta, endVal, ay.Income, ay.CapGains}
                        data = append(data, d)
                }
                val = endVal
        }
        return data
}

func (a *Asset)getIndexForYear(year int) int {
        return year - a.startingYear
}
