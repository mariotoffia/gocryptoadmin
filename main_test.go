package main

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/mariotoffia/gocryptoadmin/common"
	"github.com/mariotoffia/gocryptoadmin/output"
	"github.com/mariotoffia/gocryptoadmin/parsers"
	"github.com/mariotoffia/gocryptoadmin/processors"
	"github.com/mariotoffia/gocryptoadmin/txhistory"
	"github.com/mariotoffia/gocryptoadmin/txhistory/bittrex"
	"github.com/mariotoffia/gocryptoadmin/txhistory/coinbasepro"
	"github.com/mariotoffia/gocryptoadmin/txhistory/ofx"
	"github.com/mariotoffia/gocryptoadmin/txlog"
	txlcbp "github.com/mariotoffia/gocryptoadmin/txlog/coinbasepro"
)

// init populates the cache if missing.
func init() {

	if _, err := os.Stat("data/cost-unit/resolvers"); !os.IsNotExist(err) {
		return /* cache already populated*/
	}

	txr := txhistory.NewTxOHCReader().
		Register("cbx", coinbasepro.New("")).
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

	entriesLTCBTC := txr.Read(common.AssetPair{
		Asset:    common.AssetTypeLTC,
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

	os.MkdirAll("data/cost-unit/resolvers", 0700)

	cache := txhistory.NewTxOHCCache().
		Add(entriesBTCEUR, common.ExchangeAll).
		Add(entriesETHBTC, common.ExchangeAll).
		Add(entriesLTCBTC, common.ExchangeAll).
		Add(entriesEURUSD, common.ExchangeAll).
		Add(entriesUSDEUR, common.ExchangeAll).
		Add(entriesEURSEK, common.ExchangeAll).
		Add(entriesUSDTUSD, common.ExchangeAll)

	//defer cache.Clear("data/cost-unit/resolvers")

	cache.Store(
		"data/cost-unit/resolvers", cache.GetExchanges(common.ExchangeAll)...,
	)

}

func TestBuySell(t *testing.T) {

	// 0. If you need more exchange rate - delete the data/cost-unit/resolvers
	//    and add entries in the init() function (see above).

	// 1. Setup resolver
	expr := parsers.NewResolverParser().
		Parse("cbx:ETH = cbx:BTC").
		Parse("cbx:LTC = cbx:BTC").
		Parse("cbx:BTC = cbx,all:EUR").
		Parse("EUR = SEK").
		GetExpressions()

	cache := txhistory.NewTxOHCCache().Load(
		"data/cost-unit/resolvers",
		func(
			cache *txhistory.TxOHCCache, exchange string, entries []common.TxOHCHistory,
		) {

			cache.Add(entries, common.ExchangeAll) // make visible to all as well

		})

	resolver := txhistory.NewTxOHCResolver(cache).AddTranslations(expr...)

	// 2. Load tx log
	tx := txlog.NewTxLogReader(processors.NewChronologicalTxEntryProcessor()).
		UseDir("./data/redovisning-anders").
		IgnoreUnknownFiles().
		RegisterReader("cbx", txlcbp.NewTransactionLogReader()).
		Read()

	// 3. Apply cost unit tracking on tx log entries
	coproc := processors.NewCostUnitProcessor(resolver, nil /*default pricing*/)

	coproc.RegisterAsset(common.AssetTypeEuro, common.AssetTypeSvenskKrona)
	coproc.ProcessMany(tx)

	tx = coproc.Flush()

	// 4. Group the transaction entries
	proc := processors.NewTxGroupProcessor(time.Hour * 20 /*20h*/)

	for _, tx := range tx {
		proc.Process(&tx)
	}

	txg := proc.Flush()

	// 5. Apply accounting (single exchange)
	acc := processors.NewAccountingProcessor(common.ExchangeAll) // cbx
	for _, tx := range txg {
		acc.Process(&tx)
	}

	transactions := acc.Flush()

	// 6. Output tax report
	buysell := processors.NewTxBuySellProcessor()
	buysell.UseTaxationMarking()
	//buysell.UseLog()

	buysell.ProcessMany(transactions)

	txPairs, in_possession := buysell.Flush()

	op := output.NewStdPrinterDefaults(os.Stdout, "sek-default-buysell")

	for _, tx := range txPairs {

		op.Process(tx)

	}

	op.Flush()

	// 7. Output assets still in possession
	op = output.NewStdPrinterDefaults(os.Stdout, "default-no-account")

	for _, tx := range in_possession {

		op.Process(tx)
	}

	op.Flush()
	fmt.Println()

}
