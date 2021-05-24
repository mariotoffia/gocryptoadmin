package processors

import (
	"fmt"
	"testing"
	"time"

	"github.com/mariotoffia/gocryptoadmin/common"
	"github.com/mariotoffia/gocryptoadmin/txlog"
	"github.com/mariotoffia/gocryptoadmin/txlog/coinbasepro"
	"github.com/mariotoffia/gocryptoadmin/utils"
)

// 1. BTC-USDT -> USDT-USD -> USD-EUR -> EUR (btx)
// 2. BTC-EUR -> EUR (cbx)
// 3. LTC-BTC -> BTC-EUR -> EUR (cbx)
// 4. LTC-BTC -> BTC-USDT -> USDT-USD -> USD-EUR -> EUR (btx)

// Resolver Pattern
// -----------------
// 1. USDT = btx,all:USDT-USD;ofx,all:USD-EUR
// 2. N/A
// 3. BTC = cbx,all:BTC-EUR
// 4. BTC = btx,all:BTC-USDT (it will find 1 since now USDT)

func TestApplyEURCostUnit(t *testing.T) {

	txr := txlog.NewTxLogReader(NewChronologicalTxEntryProcessor()).
		RegisterReader("cbx", coinbasepro.NewTransactionLogReader())

	tx := txr.ReadBufferAsExchange(
		"cbx", utils.ReadFile("testfiles/cost-unit/cb.csv"),
	)

	proc := NewTxGroupProcessor(time.Hour * 20 /*20h*/)

	for _, tx := range tx {
		proc.Process(tx)
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
		fmt.Println(cfa.ConsoleHeader())

		for i, tx := range transactions {

			fmt.Println(
				tx.(common.ConsoleFormatter).ConsoleString(),
			)

			if i == -1 {
				break
			}

		}

	}

}
