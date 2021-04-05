package coinbasepro

import (
	"fmt"
	"time"

	"github.com/jszwec/csvutil"
	"github.com/mariotoffia/gocryptoadmin/common"
)

// cbp implements the `TransactionLogReader` interface
type cbp struct {
}

func NewTransactionLogReader() common.TransactionLogReader {
	return &cbp{}
}

func (c *cbp) Unmarshal(data []byte) []common.TransactionLog {

	var v []CbpTransaction
	err := csvutil.Unmarshal(data, &v)

	if err != nil {
		panic(err)
	}

	tx := []common.TransactionLog{}

	for i := range v {
		tx = append(tx, c.Transform(v[i]))
	}

	return tx
}

// Transform will get the instance pointer returned from `Entry`
// and is expected to transform to a `Transaction`
func (c *cbp) Transform(entry interface{}) common.TransactionLog {

	if v, ok := entry.(CbpTransaction); ok {

		return common.TransactionLog{
			ID:           v.ID,
			Exchange:     "coinbase-pro",
			Side:         v.Side,
			CreatedAt:    v.CreatedAt,
			AssetSize:    v.Size,
			PricePerUnit: v.Price,
			Fee:          v.Fee,
			TotalPrice:   v.Total,
			AssetPair: common.AssetPair{
				Asset:    common.AssetType(v.Unit),
				CostUnit: common.AssetType(v.CostUnit),
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
