package processors

import (
	"os"
	"time"

	"github.com/mariotoffia/gocryptoadmin/common"
	"github.com/mariotoffia/gocryptoadmin/txhistory"
	"github.com/mariotoffia/gocryptoadmin/txhistory/bittrex"
	cbx "github.com/mariotoffia/gocryptoadmin/txhistory/coinbasepro"
	"github.com/mariotoffia/gocryptoadmin/txhistory/ofx"
)

// init populates the cache if missing.
func init() {

	if _, err := os.Stat("testfiles/cost-unit/resolvers"); !os.IsNotExist(err) {
		return /* cache already populated*/
	}

	txr := txhistory.NewTxOHCReader().
		Register("cbx", cbx.New("")).
		Register("ofx", ofx.New("")).
		Register("btx", bittrex.New(""))

	from, _ := time.Parse(time.RFC3339, "2017-08-31T00:00:00.000Z")

	entriesBTCEUR := txr.Read(common.AssetPair{
		Asset:    common.AssetTypeBTC,
		CostUnit: common.AssetTypeEuro,
	}, from, time.Hour*24, "cbx")

	entriesETHBTC := txr.Read(common.AssetPair{
		Asset:    common.AssetTypeETH,
		CostUnit: common.AssetTypeBTC,
	}, from, time.Hour*24, "cbx")

	entriesEURUSD := txr.Read(common.AssetPair{
		Asset:    common.AssetTypeEuro,
		CostUnit: common.AssetTypeUsDollar,
	}, from, time.Hour*24, "ofx")

	entriesUSDEUR := txr.Read(common.AssetPair{
		Asset:    common.AssetTypeUsDollar,
		CostUnit: common.AssetTypeEuro,
	}, from, time.Hour*24, "ofx")

	entriesEURSEK := txr.Read(common.AssetPair{
		Asset:    common.AssetTypeEuro,
		CostUnit: common.AssetTypeSvenskKrona,
	}, from, time.Hour*24, "ofx")

	entriesUSDTUSD := txr.Read(common.AssetPair{
		Asset:    common.AssetTypeUSDT,
		CostUnit: common.AssetTypeUsDollar,
	}, from, time.Hour*24, "btx")

	os.MkdirAll("testfiles/cost-unit/resolvers", 0700)

	cache := txhistory.NewTxOHCCache().
		Add(entriesBTCEUR, common.ExchangeAll).
		Add(entriesETHBTC, common.ExchangeAll).
		Add(entriesEURUSD, common.ExchangeAll).
		Add(entriesUSDEUR, common.ExchangeAll).
		Add(entriesEURSEK, common.ExchangeAll).
		Add(entriesUSDTUSD, common.ExchangeAll)

	//defer cache.Clear("testfiles/cache-test")

	cache.Store(
		"testfiles/cost-unit/resolvers", cache.GetExchanges(common.ExchangeAll)...,
	)

}
