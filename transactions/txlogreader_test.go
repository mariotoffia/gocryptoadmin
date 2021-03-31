package transactions

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/ahmetb/go-linq/v3"
	"github.com/mariotoffia/gocryptoadmin/transactions/coinbasepro"
	"github.com/mariotoffia/gocryptoadmin/transactions/txcommon"
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

	tx := NewTxLogReader().
		UseDir(filepath.Dir(fp)).
		RegisterReader("coinbasepro", coinbasepro.NewTransactionLogReader()).
		Read()

	fmt.Println(tx[0].ToCSVEntry())
	fmt.Println(len(tx))

}

func TestProcessing(t *testing.T) {

	entries := NewTxLogReader().
		UseDir("../data").
		IgnoreUnknownFiles().
		UseWindowSize(6*60*60 /*6h*/).
		RegisterReader("coinbasepro", coinbasepro.NewTransactionLogReader()).
		Read()

	var ltc []txcommon.Transaction
	linq.From(entries).
		Where(func(tx interface{}) bool {
			return tx.(txcommon.Transaction).Asset == "LTC-EUR"
		}).
		ToSlice(&ltc)

	for _, tx := range ltc {

		fmt.Printf(
			"[%s:%d] %s %s %s \n S:%f  F:%f  T:%f P:%f\nGS:%f GF:%f GT:%f\n",
			tx.Exchange, tx.GroupID, tx.CreatedAt.String(), tx.Side, tx.Asset,
			tx.Size, tx.Fee, tx.Total, tx.Price,
			tx.GrpSize, tx.GrpFee, tx.GrpTotal,
		)

	}

	fmt.Println(len(entries))
}
