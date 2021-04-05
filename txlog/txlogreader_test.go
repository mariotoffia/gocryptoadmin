package txlog

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/mariotoffia/gocryptoadmin/processors"
	"github.com/mariotoffia/gocryptoadmin/txlog/coinbasepro"
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
		RegisterReader("coinbasepro", coinbasepro.NewTransactionLogReader()).
		ReadBuffer("coinbasepro", []byte(data))

	for _, tx := range tx {

		fmt.Printf(
			"[%s] %s %s %s \n S:%f  F:%f  T:%f P:%f\n",
			tx.Exchange, tx.CreatedAt.String(), tx.Side, tx.Asset,
			tx.AssetSize, tx.Fee, tx.TotalPrice, tx.PricePerUnit,
		)

	}

	fmt.Println(len(tx))
}

func TestReadCoinbasedTxLogFiles(t *testing.T) {

	tx := NewTxLogReader(processors.NewChronologicalTxEntryProcessor()).
		UseDir("../data").
		IgnoreUnknownFiles().
		RegisterReader("coinbasepro", coinbasepro.NewTransactionLogReader()).
		Read()

	for _, tx := range tx {

		fmt.Printf(
			"[%s] %s %s %s \n S:%f  F:%f  T:%f P:%f\n",
			tx.Exchange, tx.CreatedAt.String(), tx.Side, tx.Asset,
			tx.AssetSize, tx.Fee, tx.TotalPrice, tx.PricePerUnit,
		)

	}

	fmt.Println(len(tx))
}
