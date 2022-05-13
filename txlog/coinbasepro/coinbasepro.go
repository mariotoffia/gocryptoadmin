package coinbasepro

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/jszwec/csvutil"
	"github.com/mariotoffia/gocryptoadmin/common"
)

// cbp implements the `TransactionLogReader` interface
type cbp struct {
	exchange string
}

func NewTransactionLogReader() common.TransactionLogReader {
	return &cbp{exchange: "cbx"}
}

func (c *cbp) SetExchange(name string) common.TransactionLogReader {

	c.exchange = name

	return c
}

func (c *cbp) Unmarshal(data []byte) []common.TransactionLog {

	dec, err := csvutil.NewDecoder(csv.NewReader(bytes.NewReader(data)))

	if err != nil {
		panic(err)
	}

	header := dec.Header()

	tx := []common.TransactionLog{}

	for {
		cbp := CbpTransaction{}
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
					os.Stderr, "[cbx] Unknown field: %s = %s", header[i], dec.Record()[i],
				)

			}

		}

		tx = append(tx, c.Transform(&cbp, sideIdentifier))

	}

	return tx
}

// Transform will get the instance pointer returned from `Entry`
// and is expected to transform to a `Transaction`
func (c *cbp) Transform(v *CbpTransaction, sideIdentifier string) common.TransactionLog {

	return common.TransactionLog{
		ID:             v.ID,
		Exchange:       c.exchange,
		Side:           v.Side,
		SideIdentifier: sideIdentifier,
		CreatedAt:      v.CreatedAt,
		AssetSize:      v.Size,
		PricePerUnit:   v.Price,
		Fee:            v.Fee,
		TotalPrice:     v.Total,
		AssetPair: common.AssetPair{
			Asset:    common.AssetType(v.Unit),
			CostUnit: common.AssetType(v.CostUnit),
		},
	}
}

type CbpCost struct {
	Price    float64 `csv:"price"                json:"price"`
	Fee      float64 `csv:"fee"                  json:"fee"`
	Total    float64 `csv:"total"                json:"total"`
	CostUnit string  `csv:"price/fee/total unit" json:"priceFeeTotalUnit"`
}

type CbpProduct struct {
	Product string  `csv:"product"   json:"product"`
	Size    float64 `csv:"size"      json:"size"`
	Unit    string  `csv:"size unit" json:"unit"`
}

type CbpTransaction struct {
	Portfolio string          `csv:"portfolio"  json:"portfolio"`
	ID        string          `csv:"trade id"   json:"id"`
	Side      common.SideType `csv:"side"       json:"side"`
	CreatedAt time.Time       `csv:"created at" json:"created"`
	CbpProduct
	CbpCost
}
