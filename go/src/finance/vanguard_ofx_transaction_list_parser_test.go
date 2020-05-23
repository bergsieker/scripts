package finance

import (
	"strings"
	"testing"
	"time"
)

func TestVanguardOfxTransactionListReaderRead(t *testing.T) {
	testCases := []struct {
		input string
		trans Transaction
	}{
		{
			"12345678,12/23/2019,12/23/2019,Distribution,INCOME DIVIDEND    .4278,Small-Cap Index Fund Adm,79.44,5.626,446.94,446.94,",
			Transaction{
				Date:      time.Date(2019, time.December, 23, 0, 0, 0, 0, time.UTC),
				Account:   "",
				AssetName: "Small-Cap Index Fund Adm",
				Amount:    446.94,
				Type:      TransactionType_Dividend,
				NumShares: 5.626,
			},
		},
	}

	for _, testCase := range testCases {
		sr := strings.NewReader(testCase.input)
		r := NewVanguardOfxTransactionListReader(sr, "12345678")
		got, err := r.Read()
		if err != nil {
			t.Errorf("while reading %s: %v", testCase.input, err)
		}
		if testCase.trans != got {
			t.Errorf("Incorrect %s: want: %v, got: %v", testCase.input, testCase.trans, got)
		}
	}
}

func TestParseVanguardOfxTransactionType(t *testing.T) {
	tests := []struct {
		input     string
		transType TransactionType
		wantErr   bool
	}{
		{"not a valid transaction", TransactionType_Unknown, true},
		{"INCOME DIVIDEND", TransactionType_Dividend, false},
		{"INCOME DIVIDEND   0.1234", TransactionType_Dividend, false},
		{"INVALID INCOME DIVIDEND", TransactionType_Unknown, true},
		{"LT CAP GAIN", TransactionType_LtCapGain, false},
		{"LT CAP GAIN  0.1234", TransactionType_LtCapGain, false},
		{"INVALID LT CAP GAIN", TransactionType_Unknown, true},
		{"ST CAP GAIN", TransactionType_StCapGain, false},
		{"ST CAP GAIN  0.1234", TransactionType_StCapGain, false},
		{"INVALID ST CAP GAIN", TransactionType_Unknown, true},
		{"ADMIRAL CONVERSION TO", TransactionType_Transfer, false},
		{"ADMIRAL CONVERSION TO foo", TransactionType_Transfer, false},
		{"INVALID ADMIRAL CONVERSION TO", TransactionType_Unknown, true},
	}
	for _, test := range tests {
		tt, err := parseVanguardOfxTransactionType(test.input)
		if test.wantErr && err == nil || !test.wantErr && err != nil || test.transType != tt {
			t.Errorf("%v: %v (%v)", test, tt, err)
		}
	}
}
