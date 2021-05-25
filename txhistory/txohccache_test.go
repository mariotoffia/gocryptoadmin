package txhistory

import (
	"os"
	"testing"
	"time"

	"github.com/mariotoffia/gocryptoadmin/common"
	"github.com/mariotoffia/gocryptoadmin/txhistory/bittrex"
	"github.com/mariotoffia/gocryptoadmin/txhistory/coinbasepro"
	"github.com/mariotoffia/gocryptoadmin/txhistory/ofx"
	"github.com/stretchr/testify/assert"
)

func TestRenderCorrectUTCUnixTime(t *testing.T) {
	ts := time.Unix(1530057600, 0).UTC()

	assert.Equal(t, "2018-06-27T00:00:00Z", ts.Format("2006-01-02T15:04:05Z"))

}

func TestWriteNewCache(t *testing.T) {

	// We're skipping this test - run this manually
	// to populate files to be used in other tests
	t.SkipNow()

	txr := NewTxOHCReader().
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

	os.MkdirAll("testfiles/cache-test", 0700)

	cache := NewTxOHCCache().
		Add(entriesBTCEUR, common.ExchangeAll).
		Add(entriesETHBTC, common.ExchangeAll).
		Add(entriesEURUSD, common.ExchangeAll).
		Add(entriesUSDEUR, common.ExchangeAll).
		Add(entriesEURSEK, common.ExchangeAll).
		Add(entriesUSDTUSD, common.ExchangeAll)

	//defer cache.Clear("testfiles/cache-test")

	cache.Store(
		"testfiles/cache-test", cache.GetExchanges(common.ExchangeAll)...,
	)

}

func TestResolveBTCEURPriceInMiddleOfEntry(t *testing.T) {

	cache := NewTxOHCCache().Load(
		"testfiles/cache-test",
		func(
			cache *TxOHCCache, exchange string, entries []common.TxOHCHistory,
		) {

			// Skip all
			if exchange != common.ExchangeAll {
				cache.Add(entries, common.ExchangeAll) // but make visible on all
			}

		})

	at, _ := time.Parse(time.RFC3339, "2018-08-31T13:00:00.000Z")

	found, _ := cache.GetEntryForAssset(common.AssetPair{
		Asset:    common.AssetTypeBTC,
		CostUnit: common.AssetTypeEuro,
	}, at, "cbx", common.ExchangeAll)

	assert.Equal(t, "2018-08-31T00:00:00", found.DateTime.Format("2006-01-02T15:04:05"))
}

func TestResolveBTCEURPriceExactMatch(t *testing.T) {

	cache := NewTxOHCCache().Load(
		"testfiles/cache-test",
		func(
			cache *TxOHCCache, exchange string, entries []common.TxOHCHistory,
		) {

			// Skip all
			if exchange != common.ExchangeAll {
				cache.Add(entries, common.ExchangeAll) // but make visible on all
			}

		})

	at, _ := time.Parse(time.RFC3339, "2018-08-31T00:00:00.000Z")

	found, _ := cache.GetEntryForAssset(common.AssetPair{
		Asset:    common.AssetTypeBTC,
		CostUnit: common.AssetTypeEuro,
	}, at, "cbx", common.ExchangeAll)

	assert.Equal(t, "2018-08-31T00:00:00", found.DateTime.Format("2006-01-02T15:04:05"))
}
