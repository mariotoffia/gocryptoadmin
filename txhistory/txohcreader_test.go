package txhistory

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/mariotoffia/gocryptoadmin/common"
	"github.com/mariotoffia/gocryptoadmin/txhistory/coinbasepro"
	"github.com/mariotoffia/gocryptoadmin/txhistory/kraken"
)

// https://support.kraken.com/hc/en-us/articles/360047124832-Downloadable-historical-OHLCVT-Open-High-Low-Close-Volume-Trades-data

func TestReadFromKraken(t *testing.T) {

	txr := NewTxOHCReader().Register("kr", kraken.New(""))
	from, _ := time.Parse(time.RFC3339, "2017-08-31T00:00:00.000Z")

	entries := txr.Read(common.AssetPair{
		Asset:    common.AssetTypeBTC,
		CostUnit: common.AssetTypeEuro,
	}, from, time.Hour*24*7, "kr")

	data, _ := json.MarshalIndent(entries, "", " ")
	fmt.Println(string(data))
	fmt.Println(len(entries))
}

func TestReadFromCbx(t *testing.T) {

	txr := NewTxOHCReader().Register("cbx", coinbasepro.New(""))
	from, _ := time.Parse(time.RFC3339, "2017-08-31T00:00:00.000Z")

	entries := txr.Read(common.AssetPair{
		Asset:    common.AssetTypeBTC,
		CostUnit: common.AssetTypeEuro,
	}, from, time.Hour*24, "cbx")

	data, _ := json.MarshalIndent(entries, "", " ")
	fmt.Println(string(data))
	fmt.Println(len(entries))
}
