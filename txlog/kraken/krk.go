package kraken

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/jszwec/csvutil"
	"github.com/mariotoffia/gocryptoadmin/common"
)

// krk implements the `TransactionLogReader` interface
type krk struct {
	exchange string
}

func NewTransactionLogReader() common.TransactionLogReader {
	return &krk{
		exchange: "krk",
	}
}

func (c *krk) SetExchange(name string) common.TransactionLogReader {

	c.exchange = name

	return c
}

func (c *krk) Unmarshal(data []byte) []common.TransactionLog {

	dec, err := csvutil.NewDecoder(csv.NewReader(bytes.NewReader(data)))

	if err != nil {
		panic(err)
	}

	header := dec.Header()

	tx := []common.TransactionLog{}

	for {
		cbp := KrkTransaction{}
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
					os.Stderr, "[krk] Unknown field: %s = %s", header[i], dec.Record()[i],
				)

			}

		}

		tx = append(tx, c.Transform(&cbp, sideIdentifier))

	}

	return tx
}

// Transform will get the instance pointer returned from `Entry`
// and is expected to transform to a `Transaction`
func (c *krk) Transform(v *KrkTransaction, sideIdentifier string) common.TransactionLog {

	t, err := time.Parse("2006-01-02 15:04:05.0000", v.CreatedAt)
	if err != nil {

		t, err = time.Parse("2006-01-02 15:04:05.000", v.CreatedAt)

		if err != nil {
			panic(err)
		}
	}

	tx := common.TransactionLog{
		ID:             v.ID,
		Exchange:       c.exchange,
		Side:           toSide(v.Side),
		SideIdentifier: sideIdentifier,
		CreatedAt:      t.UTC(),
		AssetSize:      v.Size,
		PricePerUnit:   v.Price,
		Fee:            v.Fee,
		AssetPair:      toAssetPair(v.Pair),
	}

	tx.TotalPrice = toTotalPrice(v.Total, v.Fee, tx.Side)

	if tx.Side == common.SideTypeBuy || tx.Side == common.SideTypeTransfer {
		tx.TotalPrice = -tx.TotalPrice
	}

	return tx
}

// toTotalPrice recalculates to use fee included in price.
//
// Using the following calculations
//
// 1. Sell Fee: size * price - fee [example: (0,7 * 910,32) - 1,59306 = 635,63094]
// 2. Buy Fee: size * price + fee  [example: (1782 * 0,112815) - 0,301554495 = 201,337884495]
func toTotalPrice(total, fee float64, side common.SideType) float64 {

	switch side {
	case common.SideTypeBuy, common.SideTypeTransfer:
		return total + fee
	case common.SideTypeSell, common.SideTypeReceive:
		return total - fee
	}

	panic(fmt.Sprintf("unknown side type: %v", side))
}

func toSide(side string) common.SideType {

	return common.SideType(strings.ToUpper(side))

}

func toAssetPair(pair string) common.AssetPair {

	switch pair {
	case "XXBTZEUR", "XBTZEUR":
		return common.AssetPair{Asset: common.AssetTypeBTC, CostUnit: common.AssetTypeEuro}
	case "XXRPZEUR":
		return common.AssetPair{Asset: common.AssetTypeXRP, CostUnit: common.AssetTypeEuro}
	case "XETHZEUR":
		return common.AssetPair{Asset: common.AssetTypeETH, CostUnit: common.AssetTypeEuro}
	case "XLTCZEUR":
		return common.AssetPair{Asset: common.AssetTypeLTC, CostUnit: common.AssetTypeEuro}
	case "XXLMXXBT":
		return common.AssetPair{Asset: common.AssetTypeXLM, CostUnit: common.AssetTypeBTC}
	case "ZEURZEUR":
		return common.AssetPair{Asset: common.AssetTypeEuro, CostUnit: common.AssetTypeEuro}
	case "XLTCXLTC":
		return common.AssetPair{Asset: common.AssetTypeLTC, CostUnit: common.AssetTypeLTC}
	case "XETHXETH":
		return common.AssetPair{Asset: common.AssetTypeETH, CostUnit: common.AssetTypeETH}
	case "XBTCXBTC":
		return common.AssetPair{Asset: common.AssetTypeBTC, CostUnit: common.AssetTypeBTC}
	}

	panic(fmt.Sprintf("unknown pair - please add to kraken txlog: %s", pair))

}

type KrkTransaction struct {
	ID        string  `csv:"txid"      json:"id"`
	OrderTxID string  `csv:"ordertxid" json:"ordertxid"`
	Pair      string  `csv:"pair"      json:"assetpair"`
	CreatedAt string  `csv:"time"      json:"created"`
	Side      string  `csv:"type"      json:"side"`
	OrderType string  `csv:"ordertype" json:"ordertype"`
	Price     float64 `csv:"price"     json:"price"`
	Total     float64 `csv:"cost"      json:"cost"`
	Fee       float64 `csv:"fee"       json:"fee"`
	Size      float64 `csv:"vol"       json:"size"`
	Margin    float64 `csv:"margin"    json:"margin"`
	Misc      string  `csv:"misc"      json:"notes"`
	Ledgers   string  `csv:"ledgers"   json:"ledgers"`
}
