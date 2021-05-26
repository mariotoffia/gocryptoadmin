package txhistory

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/mariotoffia/gocryptoadmin/common"
	"github.com/mariotoffia/gocryptoadmin/txhistory/bittrex"
	"github.com/mariotoffia/gocryptoadmin/txhistory/coinbasepro"
	"github.com/mariotoffia/gocryptoadmin/txhistory/kraken"
	"github.com/mariotoffia/gocryptoadmin/txhistory/ofx"
	"github.com/stretchr/testify/assert"
)

// https://support.kraken.com/hc/en-us/articles/360047124832-Downloadable-historical-OHLCVT-Open-High-Low-Close-Volume-Trades-data

func TestReadFromKraken(t *testing.T) {

	txr := NewTxOHCReader().Register("kr", kraken.New(""))
	from, _ := time.Parse(time.RFC3339, "2020-12-31T00:00:00.000Z")

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

func TestGetEURtoUSDExchangeRate(t *testing.T) {

	txr := NewTxOHCReader().Register("ofx", ofx.New(""))
	from, _ := time.Parse(time.RFC3339, "2017-08-31T00:00:00.000Z")

	entries := txr.Read(common.AssetPair{
		Asset:    common.AssetTypeEuro,
		CostUnit: common.AssetTypeUsDollar,
	}, from, time.Hour*24*31, "ofx")

	data, err := json.MarshalIndent(entries, "", " ")

	assert.Equal(t, nil, err)

	fmt.Println(string(data))
	fmt.Println(len(entries))
}

func TestBTCUSDTBittrex(t *testing.T) {

	txr := NewTxOHCReader().Register("btx", bittrex.New(""))
	from, _ := time.Parse(time.RFC3339, "2017-08-31T00:00:00.000Z")

	entries := txr.Read(common.AssetPair{
		Asset:    common.AssetTypeUSDT,
		CostUnit: common.AssetTypeUsDollar,
	}, from, time.Hour*24, "btx")

	data, err := json.MarshalIndent(entries, "", " ")

	assert.Equal(t, nil, err)

	fmt.Println(string(data))
	fmt.Println(len(entries))
}
