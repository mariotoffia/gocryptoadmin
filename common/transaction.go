package common

import (
	"time"
)

type SideType string

const (
	SideTypeUnknown  SideType = "UNKNOWN"
	SideTypeBuy      SideType = "BUY"
	SideTypeSell     SideType = "SELL"
	SideTypeReceive  SideType = "RECEIVE"
	SideTypeTransfer SideType = "TRANSFER"
)

type TransactionEntry interface {
	GetID() string
	GetExchange() string
	GetSide() SideType
	GetCreatedAt() time.Time
	GetAssetSize() float64
	GetPricePerUnit() float64
	GetFee() float64
	GetTotalPrice() float64

	GetAssetPair() AssetPair
}

// TransactionLog represents a single transaction
//
type TransactionLog struct {
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

func (tx *TransactionLog) GetID() string {
	return tx.ID
}

func (tx *TransactionLog) GetExchange() string {
	return tx.Exchange
}

func (tx *TransactionLog) GetSide() SideType {
	return tx.Side
}

func (tx *TransactionLog) GetCreatedAt() time.Time {
	return tx.CreatedAt
}

func (tx *TransactionLog) GetAssetSize() float64 {
	return tx.AssetSize
}

func (tx *TransactionLog) GetPricePerUnit() float64 {
	return tx.PricePerUnit
}

func (tx *TransactionLog) GetFee() float64 {
	return tx.Fee
}

func (tx *TransactionLog) GetTotalPrice() float64 {
	return tx.TotalPrice
}

func (tx *TransactionLog) GetAssetPair() AssetPair {
	return tx.AssetPair
}
