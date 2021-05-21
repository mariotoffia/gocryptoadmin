package processors

import (
	"fmt"
	"testing"
	"time"

	"github.com/mariotoffia/gocryptoadmin/common"
	"github.com/mariotoffia/gocryptoadmin/txlog"
	"github.com/mariotoffia/gocryptoadmin/txlog/coinbasepro"
	"github.com/mariotoffia/gocryptoadmin/utils"
	"github.com/stretchr/testify/assert"
)

func TestAccountingCoinbaseProFiles(t *testing.T) {

	tx := txlog.NewTxLogReader(NewChronologicalTxEntryProcessor()).
		UseDir("../data").
		IgnoreUnknownFiles().
		RegisterReader("coinbasepro", coinbasepro.NewTransactionLogReader()).
		Read()

	proc := NewTxGroupProcessor(time.Hour * 20 /*20h*/)

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
}

func TestReceiveAndSellAllShallHaveOnlyEuroLeft(t *testing.T) {

	tx := txlog.NewTxLogReader(NewChronologicalTxEntryProcessor()).
		RegisterReader("coinbasepro", coinbasepro.NewTransactionLogReader()).
		ReadBuffer("coinbasepro", utils.ReadFile("testfiles/receivetest.csv"))

	proc := NewTxGroupProcessor(time.Hour * 20 /*20h*/)

	for _, tx := range tx {
		proc.Process(tx)
	}

	txg := proc.Flush()

	acc := NewAccountingProcessor()
	for i := range txg {
		acc.Process(&txg[i]) // Since accepting interface, use indexer
	}

	txa := acc.Flush()

	assert.Equal(t, float64(0.0), txa[1].(common.AccountEntry).GetAccountStatus()["LTC"])
	assert.Equal(t, float64(750.00135), txa[1].(common.AccountEntry).GetAccountStatus()["EUR"])
}

func TestReceiveAndSellReceive(t *testing.T) {

	tx := txlog.NewTxLogReader(NewChronologicalTxEntryProcessor()).
		RegisterReader("coinbasepro", coinbasepro.NewTransactionLogReader()).
		ReadBuffer("coinbasepro", utils.ReadFile("testfiles/receivetest-ii.csv"))

	proc := NewTxGroupProcessor(time.Hour * 20 /*20h*/)

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

		if i == 5 {
			break
		}

	}

	assert.Equal(t, float64(0.0), txa[1].(common.AccountEntry).GetAccountStatus()["LTC"])
	assert.Equal(t, float64(750.00135), txa[1].(common.AccountEntry).GetAccountStatus()["EUR"])
}
