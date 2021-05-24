package txhistory

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/mariotoffia/gocryptoadmin/common"
	"github.com/mariotoffia/gocryptoadmin/txhistory/coinbasepro"
)

func TestWriteNewCache(t *testing.T) {

	txr := NewTxOHCReader().Register("cbx", coinbasepro.New(""))
	from, _ := time.Parse(time.RFC3339, "2017-08-31T00:00:00.000Z")

	entriesBTCEUR := txr.Read(common.AssetPair{
		Asset:    common.AssetTypeBTC,
		CostUnit: common.AssetTypeEuro,
	}, from, time.Hour*24, "cbx")

	entriesETHBTC := txr.Read(common.AssetPair{
		Asset:    common.AssetTypeETH,
		CostUnit: common.AssetTypeBTC,
	}, from, time.Hour*24, "cbx")

	os.MkdirAll("testfiles/cache-test", 0700)
	cache := NewTxOHCCache().
		Add(entriesBTCEUR).
		Add(entriesETHBTC).
		Store("testfiles/cache-test")

	fmt.Println(cache.GetExchanges())
}
