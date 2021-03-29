package coinbasepro

import (
	"fmt"
	"time"

	"github.com/jszwec/csvutil"
	"github.com/mariotoffia/gocryptoadmin/transactions/txcommon"
)

// https://github.com/jszwec/csvutil

// cbp implements the `TransactionLogReader` interface
type cbp struct {
}

func NewTransactionLogReader() txcommon.TransactionLogReader {
	return &cbp{}
}

func (c *cbp) Unmarshal(data []byte) []txcommon.Transaction {

	var v []CbpTransaction
	err := csvutil.Unmarshal(data, &v)

	if err != nil {
		panic(err)
	}

	tx := []txcommon.Transaction{}

	for i := range v {
		tx = append(tx, c.Transform(v[i]))
	}

	return tx
}

// Transform will get the instance pointer returned from `Entry`
// and is expected to transform to a `Transaction`
func (c *cbp) Transform(entry interface{}) txcommon.Transaction {

	if v, ok := entry.(CbpTransaction); ok {

		return txcommon.Transaction{
			Portfolio: v.Portfolio,
			ID:        v.ID,
			Side:      v.Side,
			CreatedAt: v.CreatedAt,
			Product: txcommon.Product{
				Product: v.Product,
				Size:    v.Size,
				Unit:    v.Unit,
			},
			Cost: txcommon.Cost{
				Price:    v.Price,
				Fee:      v.Fee,
				Total:    v.Total,
				CostUnit: v.CostUnit,
			},
		}

	}

	panic(
		fmt.Sprintf(
			"Incorrect entry type: %T, expecting *coinbasepro.CbpTransaction", entry,
		),
	)

}

type CbpCost struct {
	Price    float64 `csv:"price" json:"price"`
	Fee      float64 `csv:"fee" json:"fee"`
	Total    float64 `csv:"total" json:"total"`
	CostUnit string  `csv:"price/fee/total unit" json:"priceFeeTotalUnit"`
}

type CbpProduct struct {
	Product string  `csv:"product" json:"product"`
	Size    float64 `csv:"size" json:"size"`
	Unit    string  `csv:"size unit" json:"unit"`
}

type CbpTransaction struct {
	Portfolio string            `csv:"portfolio" json:"portfolio"`
	ID        string            `csv:"trade id" json:"id"`
	Side      txcommon.SideType `csv:"side" json:"side"`
	CreatedAt time.Time         `csv:"created at" json:"created"`
	CbpProduct
	CbpCost
}
