package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"log"
	"strconv"
	"strings"
)

type dailyRate struct {
	date string
	rate float32
}

func main() {
	rates, err := parseLegacyFile("LEGACY_CLOSING_RATES.csv")
	if err != nil {
		log.Fatalf("%v", err)
	}
	modernRates, err := parseModernFile("FX_RATES_DAILY-sd-2017-01-03.csv")
	if err != nil {
		log.Fatalf("%v", err)
	}

	for _, r := range modernRates {
		rates = append(rates, r)
	}

	o, err := os.Create("usdcadrates.csv")
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer o.Close()

	w := csv.NewWriter(o)
	err = w.Write([]string{"date", "rate"})
	if err != nil {
		log.Fatalf("%v", err)
	}
	for _, r := range rates {
		err = w.Write([]string{r.date, fmt.Sprintf("%f", r.rate)})
		if err != nil {
			log.Fatalf("%v", err)
		}
	}
	w.Flush()
	if w.Error() != nil {
		log.Fatalf("%v", err)
	}
}

func parseLegacyFile(path string) ([]dailyRate, error) {
	f, err := os.Open(path)
	if err != nil {
		return []dailyRate{}, err
	}
	defer f.Close()

	r, err := csvReader(f, "OBSERVATIONS")
	if err != nil {
		return []dailyRate{}, fmt.Errorf("error creating csv reader: %v", err)
	}

	rates := make([]dailyRate, 0)
	firstRow := true
	datePos := 0
	ratePos := 12
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return []dailyRate{}, nil
		}
		if firstRow {
			if record[datePos] != "date" || record[ratePos] != "IEXE0102" {
				return []dailyRate{}, fmt.Errorf("unexpected header format: %s", record)
			}
			firstRow = false
			continue
		}
		if strings.HasPrefix(record[datePos], "2017-") {
			continue
		}
		rate, err := strconv.ParseFloat(record[ratePos], 32)
		if err != nil {
				return []dailyRate{}, fmt.Errorf("error parsing rate: %v", err)
		}
		rates = append(rates, dailyRate{
			date: record[datePos],
			rate: float32(rate),
		})
	}
	return rates, nil
}

func parseModernFile(path string) ([]dailyRate, error) {
	f, err := os.Open(path)
	if err != nil {
		return []dailyRate{}, err
	}
	defer f.Close()

	r, err := csvReader(f, "OBSERVATIONS")
	if err != nil {
		return []dailyRate{}, fmt.Errorf("error creating csv reader: %v", err)
	}

	rates := make([]dailyRate, 0)
	firstRow := true
	datePos := 0
	ratePos := 25
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return []dailyRate{}, nil
		}
		if firstRow {
			if record[datePos] != "date" || record[ratePos] != "FXUSDCAD" {
				return []dailyRate{}, fmt.Errorf("unexpected header format: %s", record)
			}
			firstRow = false
			continue
		}
		rate, err := strconv.ParseFloat(record[ratePos], 32)
		if err != nil {
				return []dailyRate{}, fmt.Errorf("error parsing rate: %v", err)
		}
		rates = append(rates, dailyRate{
			date: record[datePos],
			rate: float32(rate),
		})
	}
	return rates, nil
}

func csvReader(r io.Reader, skipToLineContents string) (*csv.Reader, error) {
	csvr := csv.NewReader(r)
	csvr.LazyQuotes = true
	for {
		fields, err := csvr.Read()
		if err != nil {
			return nil, err
		}
		csvr.FieldsPerRecord = 0
		line := strings.Join(fields, ", ")
		if line == skipToLineContents {
			break
		}
	}
	return csvr, nil
}
