package processors

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/mariotoffia/gocryptoadmin/common"
	"github.com/mariotoffia/gocryptoadmin/output"
	"github.com/mariotoffia/gocryptoadmin/txlog"
	"github.com/mariotoffia/gocryptoadmin/txlog/bitstamp"
	"github.com/mariotoffia/gocryptoadmin/txlog/bittrex"
	"github.com/mariotoffia/gocryptoadmin/txlog/coinbasepro"
	"github.com/mariotoffia/gocryptoadmin/txlog/kraken"
)

func TestReadAllDataTxLog(t *testing.T) {

	t.SkipNow()

	tx := txlog.NewTxLogReader(NewChronologicalTxEntryProcessor()).
		UseDir("../data").
		IgnoreUnknownFiles().
		RegisterReader("cbx", coinbasepro.NewTransactionLogReader()).
		RegisterReader("bst", bitstamp.NewTransactionLogReader()).
		RegisterReader("btx", bittrex.NewTransactionLogReader()).
		RegisterReader("krk", kraken.NewTransactionLogReader()).
		Read()

	proc := output.NewStdPrinterDefaults(os.Stdout, "default")

	for _, tx := range tx {

		proc.Process(&tx)
	}

	proc.Flush()

	fmt.Println(len(tx))
}

func TestReadAllDataAccounting(t *testing.T) {

	tx := txlog.NewTxLogReader(NewChronologicalTxEntryProcessor()).
		UseDir("../data").
		IgnoreUnknownFiles().
		RegisterReader("cbx", coinbasepro.NewTransactionLogReader()).
		RegisterReader("bst", bitstamp.NewTransactionLogReader()).
		RegisterReader("btx", bittrex.NewTransactionLogReader()).
		RegisterReader("krk", kraken.NewTransactionLogReader()).
		RegisterReader("lf", coinbasepro.NewTransactionLogReader()).
		Read()

	proc := NewTxGroupProcessor(time.Hour * 20 /*20h*/)

	for _, tx := range tx {
		proc.Process(&tx)
	}

	txg := proc.Flush()

	acc := NewMultiExchangeAccountingProcessor()
	for i := range txg {
		acc.Process(&txg[i]) // Since accepting interface, use indexer
	}

	txa := acc.Flush()

	for exchange, txs := range txa {

		if exchange == "cbx" || exchange == common.ExchangeAll {
			continue
		}

		op := output.NewStdPrinterDefaults(os.Stdout, "default-sideid")

		for _, tx := range txs {
			op.Process(tx)
		}

		op.Flush()
	}

}
