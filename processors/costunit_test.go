package processors

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/mariotoffia/gocryptoadmin/common"
	"github.com/mariotoffia/gocryptoadmin/output"
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

		op := output.NewStdPrinterDefaults(os.Stdout, "default")
		for i, tx := range transactions {

			op.Process(tx)

			if i == -1 {
				break
			}

		}

		op.Flush()
		fmt.Println()
	}

}
