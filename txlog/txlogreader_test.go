package txlog

import (
	"fmt"
	"os"
	"testing"

	"github.com/mariotoffia/gocryptoadmin/output"
	"github.com/mariotoffia/gocryptoadmin/processors"
	"github.com/mariotoffia/gocryptoadmin/txlog/bittrex"
	"github.com/mariotoffia/gocryptoadmin/txlog/coinbasepro"
	"github.com/mariotoffia/gocryptoadmin/txlog/kraken"
)

func TestReadCoinbasedTxLog(t *testing.T) {

	tx := NewTxLogReader(processors.NewChronologicalTxEntryProcessor()).
		UseDir("testfiles/cbx").
		RegisterReader("cbx", coinbasepro.NewTransactionLogReader()).
		Read()

	proc := output.NewStdPrinterDefaults(os.Stdout, "default")

	for _, tx := range tx {

		proc.Process(&tx)
	}

	proc.Flush()

	fmt.Println(len(tx))
}

func TestReadCoinbasedTxLogFiles(t *testing.T) {

	tx := NewTxLogReader(processors.NewChronologicalTxEntryProcessor()).
		UseDir("../data").
		IgnoreUnknownFiles().
		RegisterReader("cbx", coinbasepro.NewTransactionLogReader()).
		Read()

	proc := output.NewStdPrinterDefaults(os.Stdout, "default")

	for _, tx := range tx {

		proc.Process(&tx)
	}

	proc.Flush()

	fmt.Println(len(tx))
}

func TestKrakenReadBuySell(t *testing.T) {

	tx := NewTxLogReader(processors.NewChronologicalTxEntryProcessor()).
		UseDir("testfiles/krk").
		RegisterReader("krk", kraken.NewTransactionLogReader()).
		Read()

	proc := output.NewStdPrinterDefaults(os.Stdout, "default")

	for _, tx := range tx {

		proc.Process(&tx)
	}

	proc.Flush()

	fmt.Println(len(tx))
}

func TestBittrexReadBuySell(t *testing.T) {

	tx := NewTxLogReader(processors.NewChronologicalTxEntryProcessor()).
		UseDir("testfiles/btx").
		RegisterReader("btx", bittrex.NewTransactionLogReader()).
		Read()

	proc := output.NewStdPrinterDefaults(os.Stdout, "default")

	for _, tx := range tx {

		proc.Process(&tx)
	}

	proc.Flush()

	fmt.Println(len(tx))
}
