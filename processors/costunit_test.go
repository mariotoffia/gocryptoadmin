package processors

import (
	"fmt"
	"testing"
	"time"

	"github.com/mariotoffia/gocryptoadmin/common"
	"github.com/mariotoffia/gocryptoadmin/parsers"
	"github.com/mariotoffia/gocryptoadmin/txhistory"
	"github.com/mariotoffia/gocryptoadmin/txlog"
	"github.com/mariotoffia/gocryptoadmin/txlog/coinbasepro"
	"github.com/mariotoffia/gocryptoadmin/utils"
)

func TestApplyEURCostUnit(t *testing.T) {

	expr := parsers.NewResolverParser().
		Parse("cbx:ETH = cbx:BTC").
		Parse("cbx:BTC = cbx,all:EUR").
		Parse("EUR = SEK").
		GetExpressions()

	cache := txhistory.NewTxOHCCache().Load(
		"testfiles/cost-unit/resolvers",
		func(
			cache *txhistory.TxOHCCache, exchange string, entries []common.TxOHCHistory,
		) {

			cache.Add(entries, common.ExchangeAll) // make visible to all as well

		})

	resolver := txhistory.NewTxOHCResolver(cache).AddTranslations(expr...)

	txr := txlog.NewTxLogReader(NewChronologicalTxEntryProcessor()).
		RegisterReader("cbx", coinbasepro.NewTransactionLogReader())

	tx := txr.ReadBufferAsExchange(
		"cbx", utils.ReadFile("testfiles/cost-unit/cb.csv"),
	)

	coproc := NewCostUnitProcessor(resolver, nil /*default pricing*/)

	coproc.RegisterAsset(common.AssetTypeEuro, common.AssetTypeSvenskKrona)
	coproc.ProcessMany(tx)

	tx = coproc.Flush()

	proc := NewTxGroupProcessor(time.Hour * 20 /*20h*/)

	for _, tx := range tx {
		proc.Process(&tx)
	}

	txg := proc.Flush()

	acc := NewMultiExchangeAccountingProcessor()
	for i := range txg {
		acc.Process(&txg[i]) // Since accepting interface, use indexer
	}

	for exchange, transactions := range acc.Flush() {

		if exchange != "cbx" {
			continue
		}

		cfa := transactions[0].(common.ConsoleFormatter)

		fmt.Printf("\nExchange: %s\n\n", exchange)
		fmt.Print(cfa.ConsoleHeader())

		for _, asset := range transactions[0].GetTranslatedAssets() {

			fmt.Printf(
				"\t\tTotal_%s\tFee_%s", asset, asset)

		}

		fmt.Println()

		for i, tx := range transactions {

			fmt.Print(
				tx.(common.ConsoleFormatter).ConsoleString(),
			)

			for _, asset := range tx.GetTranslatedAssets() {

				fmt.Printf(
					"\t%f\t%f",
					tx.GetTranslatedTotalPrice(asset),
					tx.GetTranslatedFee(asset),
				)

			}

			fmt.Println("")

			if i == -1 {
				break
			}

		}

	}

}
