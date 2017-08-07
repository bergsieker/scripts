package finance

import (
        "encoding/csv"
        "fmt"
        "io/ioutil"
        "os"
        "path/filepath"
        "strconv"
        "strings"
        "time"
)

type PacificRateLoader struct {
        Root string
}

func (p PacificRateLoader)GetRates() (map[string][]float32, error) {
        files, err := ioutil.ReadDir(p.Root)
        if err != nil {
                return nil, fmt.Errorf("Error loading directory: %s", err)
        }
        rates := make(map[string][]float32)
        for _, file := range files {
                if !file.IsDir() {
                        basename := file.Name()
                        fullname := filepath.Join(p.Root, basename)
                        key := strings.TrimSuffix(basename, filepath.Ext(basename))
                        _, err := strconv.Atoi(key)
                        if err != nil {
                                fmt.Println("Skipping file ", basename)
                        } else {
                                r, err := getRatesFromFile(fullname)
                                if err != nil {
                                        return nil, err
                                }
                                rates[key] = r
                        }
                }
        }
        return rates, nil
}

func getRatesFromFile(filename string) ([]float32, error) {
        f, err := os.Open(filename)
        if err != nil {
                return nil, fmt.Errorf("Error loading file %s: %s", filename, err)
        }
        defer f.Close()
        c := csv.NewReader(f)
        c.Comment = '#'
        records, err := c.ReadAll()
        if err != nil {
                return nil, err
        }
        rates := make([]float32, 366)
        for _, record := range records {
                date, err := time.Parse("2006/01/02", record[1])
                if err != nil {
                        return nil, err
                }
                rate, err := strconv.ParseFloat(record[3], 32)
                if err != nil {
                        return nil, err
                }
                rates[date.YearDay() - 1] = float32(rate)
        }
        return rates, nil
}
