package bitstamp

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jszwec/csvutil"
	"github.com/mariotoffia/gocryptoadmin/common"
	"github.com/mariotoffia/gocryptoadmin/utils"
)

// bst implements the `TransactionLogReader` interface
type bst struct {
	exchange string
}

func NewTransactionLogReader() common.TransactionLogReader {
	return &bst{
		exchange: "btx",
	}
}

func (c *bst) SetExchange(name string) common.TransactionLogReader {

	c.exchange = name

	return c
}

func (c *bst) Unmarshal(data []byte) []common.TransactionLog {

	dec, err := csvutil.NewDecoder(csv.NewReader(bytes.NewReader(data)))

	if err != nil {
		panic(err)
	}

	header := dec.Header()

	tx := []common.TransactionLog{}

	for {
		cbp := BstTransaction{}
		sideIdentifier := ""

		if err = dec.Decode(&cbp); err == io.EOF {

			break

		} else if err != nil {

			panic(err)

		}

		for _, i := range dec.Unused() {

			if header[i] == "sideid" {

				sideIdentifier = dec.Record()[i]

			} else {

				fmt.Fprintf(
					os.Stderr, "[btx] Unknown field: %s = %s", header[i], dec.Record()[i],
				)

			}

		}

		tx = append(tx, c.Transform(&cbp, sideIdentifier))

	}

	return tx
}

// Transform will get the instance pointer returned from `Entry`
// and is expected to transform to a `Transaction`
func (c *bst) Transform(v *BstTransaction, sideIdentifier string) common.TransactionLog {

	t, err := time.Parse("Jan. 02, 2006, 15:04 PM", v.CreatedAt)
	if err != nil {
		panic(err)
	}

	tx := common.TransactionLog{
		ID: utils.ToString(
			utils.HashFromString(v.CreatedAt + v.Amount + v.Rate + v.Side),
		),
		Exchange:       c.exchange,
		Side:           toSide(v),
		SideIdentifier: sideIdentifier,
		CreatedAt:      t.UTC(),
	}

	if size, asset, pok := toAmountAndAssetType(v.Amount); pok {

		tx.AssetSize = size
		tx.Asset = asset

	} else {

		panic(
			fmt.Sprintf("missing size and asset type: %s", v.Amount),
		)

	}

	if tx.Side == common.SideTypeReceive || tx.Side == common.SideTypeTransfer {

		tx.CostUnit = tx.Asset
		tx.TotalPrice = tx.AssetSize
		tx.PricePerUnit = 1

	}

	if price, costunit, pok := toAmountAndAssetType(v.Rate); pok {

		tx.PricePerUnit = price
		tx.CostUnit = costunit

	}

	if totalprice, costunit, pok := toAmountAndAssetType(v.Value); pok {

		tx.TotalPrice = totalprice
		tx.CostUnit = costunit

	}

	if fee, costunit, pok := toAmountAndAssetType(v.Fee); pok {

		tx.Fee = fee
		tx.CostUnit = costunit

	}

	if tx.TotalPrice != 0 && tx.Fee != 0 {

		tx.TotalPrice = toTotalPrice(tx.TotalPrice, tx.Fee, tx.Side)

	}

	if tx.Side == common.SideTypeBuy || tx.Side == common.SideTypeTransfer {
		tx.TotalPrice = -tx.TotalPrice
	}

	return tx
}

// toTotalPrice recalculates to use fee included in price.
//
// Using the following calculations
//
// 1. Sell Fee: total - fee [example: (0,7 * 910,32) - 1,59306 = 635,63094]
// 2. Buy Fee: total + fee  [example: (1782 * 0,112815) - 0,301554495 = 201,337884495]
func toTotalPrice(total, fee float64, side common.SideType) float64 {

	switch side {
	case common.SideTypeBuy, common.SideTypeTransfer:
		return total + fee
	case common.SideTypeSell, common.SideTypeReceive:
		return total - fee
	}

	panic(fmt.Sprintf("unknown side type: %v", side))
}

func toSide(tx *BstTransaction) common.SideType {

	switch tx.Side {
	case "Buy":
		return common.SideTypeBuy
	case "Sell":
		return common.SideTypeSell
	}

	switch tx.Type {
	case "Deposit":
		return common.SideTypeReceive
	case "Withdrawal":
		return common.SideTypeTransfer
	}

	panic(
		fmt.Sprintf("unknown subtype: %s, type: %s", tx.Side, tx.Type),
	)

}

func toAmountAndAssetType(amountAndAsset string) (float64, common.AssetType, bool) {

	if amountAndAsset == "" {
		return 0, common.AssetTypeUnknown, false
	}

	c := strings.Split(amountAndAsset, " ")

	if len(c) != 2 {
		panic(fmt.Sprintf("incorrect amount and asset: %s", amountAndAsset))
	}

	f, err := strconv.ParseFloat(c[0], 64)
	if err != nil {
		panic(err)
	}

	return f, common.AssetType(c[1]), true

}

type BstTransaction struct {
	Type      string `csv:"Type"     json:"type"`
	CreatedAt string `csv:"Datetime" json:"time"`
	Account   string `csv:"Account"  json:"portfolio"`
	Amount    string `csv:"Amount"   json:"amount"`
	Value     string `csv:"Value"    json:"value"`
	Rate      string `csv:"Rate"     json:"rate"`
	Fee       string `csv:"Fee"      json:"fee"`
	Side      string `csv:"Sub Type" json:"side"`
}
