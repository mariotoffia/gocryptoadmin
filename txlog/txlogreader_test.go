package txlog

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/mariotoffia/gocryptoadmin/output"
	"github.com/mariotoffia/gocryptoadmin/processors"
	"github.com/mariotoffia/gocryptoadmin/txlog/coinbasepro"
	"github.com/mariotoffia/gocryptoadmin/txlog/kraken"
	"github.com/stretchr/testify/require"
)

func writeFile(filename, data string) (string, error) {

	dir, err := ioutil.TempDir("", "gocryptoadmin")

	if err != nil {
		return "", err
	}

	f := filepath.Join(dir, filename)
	err = ioutil.WriteFile(f, []byte(data), 0644)

	return f, err
}

func TestReadCoinbasedTxLog(t *testing.T) {

	data := `portfolio,trade id,product,side,created at,size,size unit,price,fee,total,price/fee/total unit
default,381617,XLM-EUR,BUY,2019-06-26T09:43:18.503Z,1782.00000000,XLM,0.112815,0.301554495,-201.337884495,EUR
default,382592,XLM-EUR,SELL,2019-06-26T13:35:21.940Z,131.00000000,XLM,0.11375,0.022351875,14.878898125,EUR
default,382593,XLM-EUR,SELL,2019-06-26T13:35:46.772Z,439.00000000,XLM,0.11375,0.074904375,49.861345625,EUR`

	fp, err := writeFile("coinbasepro_tx_xyz.csv", data)
	defer os.Remove(fp)

	require.Equal(t, nil, err)

	tx := NewTxLogReader(processors.NewChronologicalTxEntryProcessor()).
		UseDir(filepath.Dir(fp)).
		RegisterReader("cbx", coinbasepro.NewTransactionLogReader()).
		ReadBuffer("cbx", []byte(data))

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
