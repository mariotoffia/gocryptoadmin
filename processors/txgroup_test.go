package processors

import (
	"fmt"
	"testing"
	"time"

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

	fmt.Println("Exchange\tSide\tDate\t\t\tPair\tSize\t\tPrice\t\tFee\t\tTotal")
	fmt.Println(
		"---------------------------------------------------------" +
			"-----------------------------------------------------",
	)

	for _, tx := range txg {

		fmt.Printf(
			"%s\t%s\t%s\t%s\t%f\t%f\t%f\t%f\n",
			tx.GetExchange(),
			tx.GetSide(),
			tx.GetCreatedAt().Format("2006-01-02 15:04:05.999999999"),
			tx.GetAssetPair().String(),
			tx.GetAssetSize(),
			tx.GetPricePerUnit(),
			tx.GetFee(),
			tx.GetTotalPrice(),
		)

	}

	fmt.Printf("tx.len: %d txg.len: %d\n", len(tx), len(txg))
}
