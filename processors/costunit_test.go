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
