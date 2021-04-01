package common

import (
	"time"
)

type SideType string

const (
	SideTypeBuy      SideType = "BUY"
	SideTypeSell     SideType = "SELL"
	SideTypeReceive  SideType = "RECEIVE"
	SideTypeTransfer SideType = "TRANSFER"
)

// Transaction represents a single transaction
//
type Transaction struct {
	ID           string    `csv:"id"       json:"id"`
	Exchange     string    `csv:"exchange" json:"exchange"`
	Side         SideType  `csv:"side"     json:"side"`
	CreatedAt    time.Time `csv:"created"  json:"created"`
	AssetSize    float64   `csv:"size"     json:"size"`
	PricePerUnit float64   `csv:"price"    json:"price"`
	Fee          float64   `csv:"fee"      json:"fee"`
	TotalPrice   float64   `csv:"total"    json:"total"`

	AssetPair
}
