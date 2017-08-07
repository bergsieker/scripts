package finance

import (
        "time"
)

type Transaction struct {
        Date time.Time
        Account string
        AssetName string
        Amount float32
        IsAcquisition bool
        IsIncome bool
        IsCapGain bool
}
