package finance

/*
import (
	"reflect"
	"testing"
	"time"
)

func TestAssetYearProcessTransaction(t *testing.T) {
        d := time.Date(2015, time.January, 4, 0, 0, 0, 0, time.UTC)
        a := AssetYear{0.0, 0.0, 0.0, 0.0}
        trans := Transaction{d, "account", "name", 10.0, true, false, false}

        a.ProcessTransaction(trans)
        e := AssetYear{10.0, 10.0, 0.0, 0.0}
        if a != e {
                t.Errorf("Incorrect AssetYear. Expected %v, got %v", e, a)
        }

        trans = Transaction{d, "account", "name", -4.5, true, false, false}
        a.ProcessTransaction(trans)
        e = AssetYear{10.0, 5.5, 0.0, 0.0}
        if a != e {
                t.Errorf("Incorrect AssetYear. Expected %v, got %v", e, a)
        }

        trans = Transaction{d, "account", "name", 1.5, true, true, false}
        a.ProcessTransaction(trans)
        e = AssetYear{10.0, 7.0, 1.5, 0.0}
        if a != e {
                t.Errorf("Incorrect AssetYear. Expected %v, got %v", e, a)
        }

        trans = Transaction{d, "account", "name", 4.5, true, true, true}
        a.ProcessTransaction(trans)
        e = AssetYear{11.5, 11.5, 6.0, 4.5}
        if a != e {
                t.Errorf("Incorrect AssetYear. Expected %v, got %v", e, a)
        }
}

func TestDumpDataForYears(t *testing.T) {
        years := []AssetYear{
                {2.0, 1.0, 3.0, 4.0},
                {6.0, 5.0, 7.0, 8.0},
                {9.0, 1.0, 2.0, 3.0},
        }
        a := Asset{2000, 15.5, years}

        e := []AssetYear{
                {17.5, 16.5, 3.0, 4.0},
                {22.5, 21.5, 7.0, 8.0},
                {30.5, 22.5, 2.0, 3.0},
        }
        d := a.DumpDataForYears(2000, 2000)
        es := e[0:1]
        if !reflect.DeepEqual(d, es) {
                t.Errorf("Incorrect dump. Expected '%v', got '%v'", es, d)
        }

        d = a.DumpDataForYears(2001, 2001)
        es = e[1:2]
        if !reflect.DeepEqual(d, es) {
                t.Errorf("Incorrect dump. Expected '%v', got '%v'", es, d)
        }

        d = a.DumpDataForYears(2001, 2002)
        es = e[1:3]
        if !reflect.DeepEqual(d, es) {
                t.Errorf("Incorrect dump. Expected '%v', got '%v'", es, d)
        }
}

func TestAssetProcessTransaction(t *testing.T) {
        d := time.Date(2002, time.January, 1, 0, 0, 0, 0, time.UTC)
        t1 := Transaction{d, "account", "name", 4.5, true, true, false}
        t2 := Transaction{d.AddDate(0, 1, 0), "account", "name", 1.2, true, true, false}
        t3 := Transaction{d.AddDate(1, 0, 0), "account", "name", 6.4, true, true, true}
        a := CreateAsset(2001, 2004)
        a.ProcessTransaction(t1)
        a.ProcessTransaction(t2)
        a.ProcessTransaction(t3)

        e := []AssetYear{
                {0.0, 0.0, 0.0, 0.0},
                {5.7, 5.7, 5.7, 0.0},
                {6.4, 6.4, 6.4, 6.4},
                {0.0, 0.0, 0.0, 0.0},
        }
        if !reflect.DeepEqual(a.years, e) {
                t.Errorf("Incorrect transaction processing. Expected '%v', got '%v'", e, a.years)
        }
}
*/
