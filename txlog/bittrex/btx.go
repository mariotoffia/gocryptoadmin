package bittrex

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

// btx implements the `TransactionLogReader` interface
type btx struct {
	exchange string
}

func NewTransactionLogReader() common.TransactionLogReader {
	return &btx{
		exchange: "btx",
	}
}

func (c *btx) SetExchange(name string) common.TransactionLogReader {

	c.exchange = name

	return c
}

func (c *btx) Unmarshal(data []byte) []common.TransactionLog {

	dec, err := csvutil.NewDecoder(csv.NewReader(bytes.NewReader(data)))

	if err != nil {
		panic(err)
	}

	header := dec.Header()

	tx := []common.TransactionLog{}

	for {
		cbp := BtxTransaction{}
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
func (c *btx) Transform(v *BtxTransaction, sideIdentifier string) common.TransactionLog {

	t, err := time.Parse("1/2/2006 15:04:05 PM", v.CreatedAt)
	if err != nil {

		t, err = time.Parse("2006/01/02 15:04:05", v.CreatedAt)

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
		PricePerUnit:   v.PricePerUnit,
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

func toSide(side string) common.SideType {

	switch side {
	case "LIMIT_BUY":
		return common.SideTypeBuy
	case "LIMIT_SELL":
		return common.SideTypeSell
	case "RECEIVE":
		return common.SideTypeReceive
	case "TRANSFER":
		return common.SideTypeTransfer
	}

	panic(fmt.Sprintf("unknown side: %s", side))

}

func toAssetPair(pair string) common.AssetPair {

	c := strings.Split(pair, "-")

	if len(c) != 2 {
		panic(fmt.Sprintf("incorrect assetpair: %s", pair))
	}

	return common.AssetPair{
		Asset:    common.AssetType(c[1]),
		CostUnit: common.AssetType(c[0]),
	}

}

type BtxTransaction struct {
	ID                string  `csv:"Uuid"              json:"id"`
	Pair              string  `csv:"Exchange"          json:"assetpair"`
	CreatedAt         string  `csv:"TimeStamp"         json:"created"`
	Side              string  `csv:"OrderType"         json:"side"`
	Limit             float64 `csv:"Limit"             json:"ordertype"`
	Size              float64 `csv:"Quantity"          json:"size"`
	SizeRemaining     float64 `csv:"QuantityRemaining" json:"remainingsize"`
	Fee               float64 `csv:"Commission"        json:"fee"`
	Total             float64 `csv:"Price"             json:"cost"`
	PricePerUnit      float64 `csv:"PricePerUnit"      json:"priceperunit"`
	IsConditional     bool    `csv:"IsConditional"     json:"conditional"`
	Condition         string  `csv:"Condition"         json:"condition"`
	ConditionTarget   float64 `csv:"ConditionTarget"   json:"conditiontarget"`
	ImmediateOrCancel bool    `csv:"ImmediateOrCancel" json:"ioc"`
	Closed            string  `csv:"Closed"            json:"closed"`
	TimeInForceTypeId int     `csv:"TimeInForceTypeId" json:"tiftid"`
	TimeInForce       string  `csv:"TimeInForce"       json:"timeinforce"`
}
