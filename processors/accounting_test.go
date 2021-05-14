package processors

import (
	"fmt"
	"testing"
	"time"

	"github.com/mariotoffia/gocryptoadmin/common"
	"github.com/mariotoffia/gocryptoadmin/txlog"
	"github.com/mariotoffia/gocryptoadmin/txlog/coinbasepro"
)

func TestAccountingCoinbaseProFiles(t *testing.T) {

	tx := txlog.NewTxLogReader(NewChronologicalTxEntryProcessor()).
		UseDir("../data").
		IgnoreUnknownFiles().
		RegisterReader("coinbasepro", coinbasepro.NewTransactionLogReader()).
		Read()

	proc := NewTxGroupProcessor(time.Hour * 20 /*20h*/) // time.Duration(30 * 60)

	for _, tx := range tx {
		proc.Process(tx)
	}

	txg := proc.Flush()

	acc := NewAccountingProcessor()
	for i := range txg {
		acc.Process(&txg[i]) // Since accepting interface, use indexer
	}

	txa := acc.Flush()

	cfa := txa[0].(common.ConsoleFormatter)
	fmt.Println(cfa.ConsoleHeader())

	for i, tx := range txa {

		fmt.Println(
			tx.(common.ConsoleFormatter).ConsoleString(),
		)

		if i == -1 {
			break
		}

	}

	fmt.Printf("tx.len: %d txg.len: %d\n", len(tx), len(txg))
}
