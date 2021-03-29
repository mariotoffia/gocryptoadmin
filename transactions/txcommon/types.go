package txcommon

import "time"

type SideType string

const (
	SideTypeBuy  string = "BUY"
	SideTypeSell string = "SELL"
)

type Cost struct {
	Price    float64 `csv:"price" json:"price"`
	Fee      float64 `csv:"fee" json:"fee"`
	Total    float64 `csv:"total" json:"total"`
	CostUnit string  `csv:"constunit" json:"unit"`
}

type Product struct {
	Product string  `csv:"product" json:"product"`
	Size    float64 `csv:"size" json:"size"`
	Unit    string  `csv:"sizeunit" json:"sizeunit"`
}

type Transaction struct {
	Portfolio string    `csv:"portfolio" json:"portfolio"`
	ID        string    `csv:"id" json:"id"`
	Side      SideType  `csv:"side" json:"side"`
	CreatedAt time.Time `csv:"created" json:"created"`
	Product
	Cost
}

type TransactionGroup struct {
	Tx []Transaction      `json:"tx,omitempty"`
	Tg []TransactionGroup `json:"tg,omitempty"`
}
