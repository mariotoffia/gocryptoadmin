package transactions

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/ahmetb/go-linq/v3"
	"github.com/mariotoffia/gocryptoadmin/transactions/txprocessors"
	"github.com/mariotoffia/gocryptoadmin/transactions/txreaders/coinbasepro"
	"github.com/mariotoffia/gocryptoadmin/utils"
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

	for _, tx := range tx {

		fmt.Printf(
			"[%s:%d] %s %s %s \n S:%f  F:%f  T:%f P:%f\nGS:%f GF:%f GT:%f\n",
			tx.Exchange, tx.GroupID, tx.CreatedAt.String(), tx.Side, tx.AssetPair,
			tx.Size, tx.Fee, tx.Total, tx.Price,
			tx.GrpSize, tx.GrpFee, tx.GrpTotal,
		)

	}

	fmt.Println(len(tx))
}

func TestCoinbasedFileTxLogShallBeOrdered(t *testing.T) {

	entries := NewTxLogReader().
		UseDir("../data").
		IgnoreUnknownFiles().
		UseWindowSize(6*60*60 /*6h*/).
		RegisterReader("coinbasepro", coinbasepro.NewTransactionLogReader()).
		Read()

	for _, tx := range entries {

		fmt.Printf(
			"[%s:%d] %s %s %s  [ID:%s]\n S:%f  F:%f  T:%f P:%f\nGS:%f GF:%f GT:%f\n",
			tx.Exchange, tx.GroupID, tx.CreatedAt.String(), tx.Side, tx.AssetPair, tx.ID,
			tx.Size, tx.Fee, tx.Total, tx.Price,
			tx.GrpSize, tx.GrpFee, tx.GrpTotal,
		)

	}

	fmt.Println(len(entries))
}

func TestWeightedPrice(t *testing.T) {

	entries := NewTxLogReader().
		UseDir("../data").
		IgnoreUnknownFiles().
		UseWindowSize(6*60*60 /*6h*/).
		RegisterReader("coinbasepro", coinbasepro.NewTransactionLogReader()).
		Read()

		/*	var ltc []txcommon.Transaction
			linq.From(entries).
				Where(func(tx interface{}) bool {
					return tx.(txcommon.Transaction).Asset == "LTC-EUR"
				}).
				ToSlice(&ltc)
		*/

	weighted := txprocessors.WeightedPrice(entries)

	for _, tx := range weighted {

		fmt.Printf(
			"[%s:%d] %s %s %s Size:%f Price:%f Fee:%f Total:%f\n",
			tx.Exchange, tx.GroupID, tx.CreatedAt.String(), tx.Side, tx.AssetPair,
			tx.Size, tx.Price, tx.Fee, tx.Total,
		)

	}

	fmt.Printf("original: %d weighted: %d\n", len(entries), len(weighted))
}

func TestPairedBuySell(t *testing.T) {

	entries := NewTxLogReader().
		UseDir("../data").
		IgnoreUnknownFiles().
		UseWindowSize(6*60*60 /*6h*/).
		RegisterReader("coinbasepro", coinbasepro.NewTransactionLogReader()).
		Read()

		/*
			var ltc []txcommon.Transaction

			linq.From(entries).
				Where(func(tx interface{}) bool {
					return tx.(txcommon.Transaction).Asset == "LTC-EUR"
				}).
				ToSlice(&ltc)*/
	ltc := entries

	weighted := txprocessors.WeightedPrice(ltc)
	paired, unpaired := txprocessors.PairBuySell(weighted)

	/* TODO: this is incorrect!
			[coinbase-pro] BTC-EUR 1.000000 (Earned: 161.350000)
		2018-01-28 21:52:23.686 +0000 UTC BP:6431.140509 BF:15.799944 BT:-6335.777526
		2018-02-07 09:52:29.562 +0000 UTC SP:6497.130000 SF:0.000000 ST:6497.130000

	buy: 10818410
	buy: 11650766
	Sell: 11661389
	*/
	for _, tx := range paired {

		if tx.BoughtTotal == -6335.777526 {
			fmt.Println("")
		}

		fmt.Printf(
			"[%s] %s %f (Earned: %f)\n%s BP:%f BF:%f BT:%f\n%s SP:%f SF:%f ST:%f\n",
			tx.Exchange, tx.AssetPair, tx.Size, utils.ToFixed(tx.BoughtTotal+tx.Sell.Total, 2),
			tx.BoughtAt.String(), tx.BoughtPrice, tx.BoughtFee, tx.BoughtTotal,
			tx.SoldAt.String(), tx.SoldPrice, tx.SoldFee, tx.SoldTotal,
		)

	}

	fmt.Printf("\nUnparied\n--------------------------\n")

	for _, tx := range unpaired {

		fmt.Printf(
			"[%s] %s %s %s %s Size: %f Price: %f Fee: %f Total: %f \n",
			tx.Exchange, tx.CreatedAt, tx.Side, tx.AssetPair, tx.Unit, tx.Size,
			tx.Price, tx.Fee, tx.Total,
		)

	}

	fmt.Printf("original: %d weighted+paired: %d\n", len(entries), len(weighted))
}

func TestEarningsPerYear(t *testing.T) {

	entries := NewTxLogReader().
		UseDir("../data").
		IgnoreUnknownFiles().
		UseWindowSize(6*60*60 /*6h*/).
		RegisterReader("coinbasepro", coinbasepro.NewTransactionLogReader()).
		Read()

	weighted := txprocessors.WeightedPrice(entries)
	paired, unpaired := txprocessors.PairBuySell(weighted)

	const layout = "2006-01-02"
	const tax = 0.3
	for _, year := range []int{2017, 2018, 2019, 2020, 2021} {

		var data []txprocessors.PairedTransaction

		linq.From(paired).
			Where(func(tx interface{}) bool {
				return tx.(txprocessors.PairedTransaction).SoldAt.Year() == year
			}).
			ToSlice(&data)

		fmt.Printf("\nÅr: %d\n------------------------------\n", year)
		fmt.Println(
			"Antal		Namn		Inköpsdatum	Inköpsbelopp	Försäljningsdatum" +
				"	Försäljningsbelopp	Vinst		Skatt",
		)

		fmt.Println(
			"-------------------------------------------------------------------" +
				"--------------------------------------------------------------------------",
		)

		totTax := float64(0)
		tot := float64(0)
		for _, tx := range data {

			totTax += utils.ToFixed((tx.BoughtTotal+tx.Sell.Total)*tax, 8)
			tot += utils.ToFixed(tx.BoughtTotal+tx.Sell.Total, 8)

			fmt.Printf(
				"%f\t%s\t\t%s\t%f\t%s\t\t%f\t\t%f\t%f\n",
				tx.Size, tx.AssetPair,
				tx.BoughtAt.Format(layout), utils.ToFixed(-tx.BoughtTotal, 2),
				tx.SoldAt.Format(layout), utils.ToFixed(tx.SoldTotal, 2),
				utils.ToFixed(tx.BoughtTotal+tx.Sell.Total, 2),
				utils.ToFixed((tx.BoughtTotal+tx.Sell.Total)*tax, 2),
			)

		}

		fmt.Printf("\n-------------------------------------\nTax: %f Earned: %f\n-------------------------------------\n", totTax, tot)

	}

	fmt.Printf("\nUnparied\n--------------------------\n")

	for _, tx := range unpaired {

		fmt.Printf(
			"[%s] %s %s %s %s Size: %f Price: %f Fee: %f Total: %f [%s]\n",
			tx.Exchange, tx.CreatedAt, tx.Side, tx.AssetPair, tx.Unit, tx.Size,
			tx.Price, tx.Fee, tx.Total, tx.CostUnit,
		)

	}

	fmt.Printf("original: %d weighted+paired: %d\n", len(entries), len(weighted))

}
