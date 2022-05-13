package processors

import (
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

func TestBuySell(t *testing.T) {

	// 1. Setup resolver
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

	// 2. Load tx log
	txr := txlog.NewTxLogReader(NewChronologicalTxEntryProcessor()).
		RegisterReader("cbx", coinbasepro.NewTransactionLogReader())

	tx := txr.ReadBufferAsExchange(
		"cbx", utils.ReadFile("testfiles/cost-unit/cb.csv"),
	)

	// 3. Apply cost unit tracking on tx log entries
	coproc := NewCostUnitProcessor(resolver, nil /*default pricing*/)

	coproc.RegisterAsset(common.AssetTypeEuro, common.AssetTypeSvenskKrona)
	coproc.ProcessMany(tx)

	tx = coproc.Flush()

	// 4. Group the transaction entries
	proc := NewTxGroupProcessor(time.Hour * 20 /*20h*/)

	for _, tx := range tx {
		proc.Process(&tx)
	}

	txg := proc.Flush()

	// 5. Apply accounting (single exchange)
	acc := NewAccountingProcessor("cbx")
	for _, tx := range txg {
		acc.Process(&tx)
	}

	transactions := acc.Flush()

	buysell := NewTxBuySellProcessor()
	//buysell.UseLog()
	buysell.UseTaxationMarking()

	buysell.ProcessMany(transactions)

	txPairs, _ := buysell.Flush()

	op := output.NewStdPrinterDefaults(os.Stdout, "default-buysell")

	for _, tx := range txPairs {

		op.Process(tx)

	}

	op.Flush()

}
