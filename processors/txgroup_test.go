package processors

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/mariotoffia/gocryptoadmin/output"
	"github.com/mariotoffia/gocryptoadmin/txlog"
	"github.com/mariotoffia/gocryptoadmin/txlog/coinbasepro"
)

func TestReadCoinbasedTxLogFiles(t *testing.T) {

	tx := txlog.NewTxLogReader(NewChronologicalTxEntryProcessor()).
		UseDir("../data").
		IgnoreUnknownFiles().
		RegisterReader("cbx", coinbasepro.NewTransactionLogReader()).
		Read()

	proc := NewTxGroupProcessor(time.Hour * 20 /*20h*/) // time.Duration(30 * 60)

	for _, tx := range tx {
		proc.Process(&tx)
	}

	txg := proc.Flush()

	op := output.NewStdPrinterDefaults(os.Stdout, "default")

	for _, tx := range txg {

		op.Process(&tx)
	}

	op.Flush()

	fmt.Printf("tx.len: %d txg.len: %d\n", len(tx), len(txg))
}
