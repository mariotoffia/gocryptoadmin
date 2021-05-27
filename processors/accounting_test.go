package processors

import (
	"os"
	"testing"
	"time"

	"github.com/mariotoffia/gocryptoadmin/common"
	"github.com/mariotoffia/gocryptoadmin/output"
	"github.com/mariotoffia/gocryptoadmin/txlog"
	"github.com/mariotoffia/gocryptoadmin/txlog/coinbasepro"
	"github.com/mariotoffia/gocryptoadmin/utils"
	"github.com/stretchr/testify/assert"
)

func TestAccountingCoinbaseProFiles(t *testing.T) {

	tx := txlog.NewTxLogReader(NewChronologicalTxEntryProcessor()).
		UseDir("../data").
		IgnoreUnknownFiles().
		RegisterReader("cbx", coinbasepro.NewTransactionLogReader()).
		Read()

	proc := NewTxGroupProcessor(time.Hour * 20 /*20h*/)

	for _, tx := range tx {
		proc.Process(&tx)
	}

	txg := proc.Flush()

	acc := NewAccountingProcessor(common.ExchangeAll)
	for _, tx := range txg {
		acc.Process(&tx) // Since accepting interface, use indexer
	}

	txa := acc.Flush()

	op := output.NewStdPrinterDefaults(os.Stdout, "default")

	for i, tx := range txa {

		op.Process(tx)
		if i == -1 {
			break
		}

	}

	op.Flush()
}

func TestReceiveAndSellAllShallHaveOnlyEuroLeft(t *testing.T) {

	tx := txlog.NewTxLogReader(NewChronologicalTxEntryProcessor()).
		RegisterReader("cbp", coinbasepro.NewTransactionLogReader()).
		ReadBufferAsExchange("cbp", utils.ReadFile("testfiles/receivetest.csv"))

	proc := NewTxGroupProcessor(time.Hour * 20 /*20h*/)

	for _, tx := range tx {
		proc.Process(&tx)
	}

	txg := proc.Flush()

	acc := NewAccountingProcessor(common.ExchangeAll)
	for _, tx := range txg {
		acc.Process(&tx)
	}

	txa := acc.Flush()

	assert.Equal(t, float64(0.0), txa[1].(common.AccountEntry).GetAccountStatus()["LTC"])
	assert.Equal(t, float64(750.00135), txa[1].(common.AccountEntry).GetAccountStatus()["EUR"])
}

func TestWhenSideIdPresentItShallBeOnTxLog(t *testing.T) {

	tx := txlog.NewTxLogReader(NewChronologicalTxEntryProcessor()).
		RegisterReader("cbp", coinbasepro.NewTransactionLogReader()).
		ReadBufferAsExchange("cbp", utils.ReadFile("testfiles/receivetest.csv"))

	proc := NewTxGroupProcessor(time.Hour * 20 /*20h*/)

	for _, tx := range tx {
		proc.Process(&tx)
	}

	txg := proc.Flush()

	acc := NewAccountingProcessor(common.ExchangeAll)
	for _, tx := range txg {
		acc.Process(&tx)
	}

	txa := acc.Flush()

	assert.Equal(t, "kraken", txa[0].GetSideIdentifier())
	assert.Equal(t, "", txa[1].GetSideIdentifier())
}

func TestWhenSideIdNotPresentItShallNotBeOnTxLog(t *testing.T) {

	tx := txlog.NewTxLogReader(NewChronologicalTxEntryProcessor()).
		RegisterReader("cbp", coinbasepro.NewTransactionLogReader()).
		ReadBufferAsExchange("cbp", utils.ReadFile("testfiles/receivetest-iii.csv"))

	proc := NewTxGroupProcessor(time.Hour * 20 /*20h*/)

	for _, tx := range tx {
		proc.Process(&tx)
	}

	txg := proc.Flush()

	acc := NewAccountingProcessor(common.ExchangeAll)
	for _, tx := range txg {
		acc.Process(&tx)
	}

	txa := acc.Flush()

	assert.Equal(t, "", txa[0].GetSideIdentifier())
	assert.Equal(t, "", txa[1].GetSideIdentifier())
}

func TestReceiveAndSellReceive(t *testing.T) {

	tx := txlog.NewTxLogReader(NewChronologicalTxEntryProcessor()).
		RegisterReader("cbp", coinbasepro.NewTransactionLogReader()).
		ReadBufferAsExchange("cbp", utils.ReadFile("testfiles/receivetest-ii.csv"))

	proc := NewTxGroupProcessor(time.Hour * 20 /*20h*/)

	for _, tx := range tx {
		proc.Process(&tx)
	}

	txg := proc.Flush()

	acc := NewAccountingProcessor(common.ExchangeAll)
	for _, tx := range txg {
		acc.Process(&tx)
	}

	txa := acc.Flush()
	op := output.NewStdPrinterDefaults(os.Stdout, "default")

	for i, tx := range txa {

		op.Process(tx)
		if i == -1 {
			break
		}

	}

	op.Flush()

	assert.Equal(t, float64(0.0), txa[1].(common.AccountEntry).GetAccountStatus()["LTC"])
	assert.Equal(t, float64(750.00135), txa[1].(common.AccountEntry).GetAccountStatus()["EUR"])
}

func TestMultiExchangeSingleAccount(t *testing.T) {

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

	acc := NewAccountingProcessor(common.ExchangeAll)
	for i := range txg {
		acc.Process(&txg[i]) // Since accepting interface, use indexer
	}

	txa := acc.Flush()

	op := output.NewStdPrinterDefaults(os.Stdout, "default")

	for i, tx := range txa {

		op.Process(tx)
		if i == -1 {
			break
		}

	}

	op.Flush()
}
