package output

import (
	"os"
	"testing"
	"time"

	"github.com/mariotoffia/gocryptoadmin/common"
)

func TestOutputDefaultFullTemplateSingleTxEntry(t *testing.T) {

	time, _ := time.Parse(time.RFC3339, "2017-12-06T13:00:00.000Z")
	tx := &common.TransactionLog{
		ID:             "1234",
		Exchange:       "cbx",
		Side:           common.SideTypeBuy,
		SideIdentifier: "cbx",
		CreatedAt:      time,
		AssetSize:      20,
		PricePerUnit:   87,
		Fee:            3,
		TotalPrice:     1743,
		AssetPair: common.AssetPair{
			Asset:    common.AssetTypeLTC,
			CostUnit: common.AssetTypeEuro,
		},
		TranslatedTotalPrice: map[string]float64{
			"EUR": 1743,
			"SEK": 17430,
		},
		TranslatedFee: map[string]float64{
			"EUR": 3,
			"SEK": 30,
		},
	}

	proc := NewStdPrinterDefaults(os.Stdout, "default")
	proc.Process(tx)
	proc.Flush()

}

func TestOutputDefaultFullTemplateSingleTxAccountEntry(t *testing.T) {

	time, _ := time.Parse(time.RFC3339, "2017-12-06T13:00:00.000Z")
	tx := &common.TransactionLog{
		ID:             "1234",
		Exchange:       "cbx",
		Side:           common.SideTypeBuy,
		SideIdentifier: "cbx",
		CreatedAt:      time,
		AssetSize:      20,
		PricePerUnit:   87,
		Fee:            3,
		TotalPrice:     1743,
		AssetPair: common.AssetPair{
			Asset:    common.AssetTypeLTC,
			CostUnit: common.AssetTypeEuro,
		},
		TranslatedTotalPrice: map[string]float64{
			"EUR": 1743,
			"SEK": 17430,
		},
		TranslatedFee: map[string]float64{
			"EUR": 3,
			"SEK": 30,
		},
	}

	proc := NewStdPrinterDefaults(os.Stdout, "default")
	proc.Process(common.NextAccountLog(nil, tx))
	proc.Flush()

}
