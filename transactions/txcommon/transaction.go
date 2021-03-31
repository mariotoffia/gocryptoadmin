package txcommon

import (
	"bytes"
	"encoding/csv"
	"time"

	"github.com/ahmetb/go-linq/v3"
	"github.com/jszwec/csvutil"
)

type SideType string

const (
	SideTypeBuy  SideType = "BUY"
	SideTypeSell SideType = "SELL"
)

type Cost struct {
	Price    float64 `csv:"price" json:"price"`
	Fee      float64 `csv:"fee" json:"fee"`
	Total    float64 `csv:"total" json:"total"`
	CostUnit string  `csv:"constunit" json:"unit"`
}

type Product struct {
	Asset string  `csv:"product" json:"product"`
	Size  float64 `csv:"size" json:"size"`
	Unit  string  `csv:"sizeunit" json:"sizeunit"`
}

// GroupAccumulate is calculated by the `transaction.preProcess` function
// of a selection of `Transaction`
//
// It represents a accumulation within the same `GroupID` (hence local).
type GroupAccumulate struct {
	GrpSize  float64 `csv:"grpsize" json:"grpsize"`
	GrpTotal float64 `csv:"grptotal" json:"grptotal"`
	GrpFee   float64 `csv:"grpfee" json:"grpfee"`
}

// Transaction represents a single transaction
//
// NOTE: The `Total` and `GrpTotal` (`Price` * `Size`) - `Fee` on both buy and sell.
// the only difference is that the buy is (-`Price` * `Size`) - `Fee`.
type Transaction struct {
	Portfolio string    `csv:"portfolio" json:"portfolio"`
	ID        string    `csv:"id" json:"id"`
	GroupID   int64     `csv:"grpid" json:"grpid"`
	Exchange  string    `csv:"exchange" json:"exchange"`
	Side      SideType  `csv:"side" json:"side"`
	CreatedAt time.Time `csv:"created" json:"created"`
	Product
	Cost
	GroupAccumulate
}

// ToCSVEntry returns the CSV representation of the `Transaction`
func (tx *Transaction) ToCSVEntry() string {

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	enc := csvutil.NewEncoder(w)
	enc.AutoHeader = false

	if err := enc.Encode(tx); err != nil {
		panic(err)
	}

	w.Flush()

	s := buf.String()

	return s[:len(s)-1]
}

func Exchanges(tx []Transaction) []string {

	var exchanges []string

	linq.From(tx).
		Select(func(tx interface{}) interface{} {
			return tx.(Transaction).Exchange
		}).
		Distinct().
		ToSlice(&exchanges)

	return exchanges
}

func Products(tx []Transaction) []string {

	var products []string

	linq.From(tx).
		Select(func(tx interface{}) interface{} {
			return tx.(Transaction).Asset
		}).
		Distinct().
		ToSlice(&products)

	return products
}

// Ordered will order the _tx by _Exchange_, _Product_, _CreatedAt_
// so it is easy to process the buy sell in chronological order for a certain
// product at a exchange.
func Ordered(tx []Transaction) []Transaction {

	var ordered []Transaction

	linq.From(tx).
		OrderBy(func(tx interface{}) interface{} {
			return tx.(Transaction).Exchange
		}).
		ThenBy(func(tx interface{}) interface{} {
			return tx.(Transaction).Asset
		}).
		ThenBy(func(tx interface{}) interface{} {
			return tx.(Transaction).CreatedAt.Unix()
		}).
		ToSlice(&ordered)

	return ordered
}
