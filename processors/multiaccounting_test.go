package processors

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/mariotoffia/gocryptoadmin/output"
	"github.com/mariotoffia/gocryptoadmin/txlog"
	"github.com/mariotoffia/gocryptoadmin/txlog/coinbasepro"
	"github.com/mariotoffia/gocryptoadmin/utils"
)

func TestMultiExchangeMultiAccount(t *testing.T) {

	txr := txlog.NewTxLogReader(NewChronologicalTxEntryProcessor()).
		RegisterReader("lf", coinbasepro.NewTransactionLogReader()).
		RegisterReader("kr", coinbasepro.NewTransactionLogReader()).
		RegisterReader("cbx", coinbasepro.NewTransactionLogReader())

	tx := txr.ReadBufferAsExchange(
		"lf", utils.ReadFile("testfiles/multi-exchange/lf.csv"),
	)

	tx = append(
		tx,
		txr.ReadBufferAsExchange(
			"kr", utils.ReadFile("testfiles/multi-exchange/kr.csv"))...,
	)

	tx = append(
		tx,
		txr.ReadBufferAsExchange(
			"cbx", utils.ReadFile("testfiles/multi-exchange/cb.csv"))...,
	)

	proc := NewTxGroupProcessor(time.Hour * 20 /*20h*/)

	for _, tx := range tx {
		proc.Process(&tx)
	}

	txg := proc.Flush()

	acc := NewMultiExchangeAccountingProcessor()
	for i := range txg {
		acc.Process(&txg[i]) // Since accepting interface, use indexer
	}

	for _, transactions := range acc.Flush() {

		op := output.NewStdPrinterDefaults(os.Stdout, "default")

		for _, tx := range transactions {

			op.Process(tx)
		}

		op.Flush()
		fmt.Println()

	}

}
