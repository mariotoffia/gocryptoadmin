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

const (
	// ExchangeAll represents all exchanges
	ExchangeAll string = "all"
)

type CostUnitTranslations interface {
	// TranslatedTotalPrice returns zero or positive value
	// for translated values (if needed) to the specified
	// `AssetType` for `GetTotalPrice`.
	//
	// NOTE: It may need be processed by a cost unit processor
	// before any valid values may be returned.
	GetTranslatedTotalPrice(asset AssetType) float64
	// TranslatedFee is the same as `TranslatedTotalPrice`
	// but reflects the `GetFee`.
	GetTranslatedFee(asset AssetType) float64
	// GetTranslatedAssets returns all `AssetType`s that can be used in
	// `GetTranslatedTotalPrice` and `GetTranslatedFee`
	GetTranslatedAssets() []AssetType
}

type TransactionEntry interface {
	CostUnitTranslations
	GetID() string
	GetExchange() string
	GetSide() SideType
	GetSideIdentifier() string
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
	ID                   string             `csv:"id"       json:"id"`
	Exchange             string             `csv:"exchange" json:"exchange"`
	Side                 SideType           `csv:"side"     json:"side"`
	SideIdentifier       string             `csv:"sideid"   json:"sideid,omitempty"`
	CreatedAt            time.Time          `csv:"created"  json:"created"`
	AssetSize            float64            `csv:"size"     json:"size"`
	PricePerUnit         float64            `csv:"price"    json:"price"`
	Fee                  float64            `csv:"fee"      json:"fee"`
	TotalPrice           float64            `csv:"total"    json:"total"`
	TranslatedTotalPrice map[string]float64 `               json:"translatedprice"`
	TranslatedFee        map[string]float64 `               json:"translatedfee"`

	AssetPair
}

func (tx *TransactionLog) GetID() string {
	return tx.ID
}

func (tx *TransactionLog) GetTranslatedTotalPrice(asset AssetType) float64 {

	if tx.TranslatedTotalPrice != nil {

		if f, ok := tx.TranslatedTotalPrice[string(asset)]; ok {
			return f
		}

	}

	return -1

}

func (tx *TransactionLog) GetTranslatedFee(asset AssetType) float64 {

	if tx.TranslatedFee != nil {

		if f, ok := tx.TranslatedFee[string(asset)]; ok {
			return f
		}

	}

	return -1

}

func (tx *TransactionLog) GetTranslatedAssets() []AssetType {

	// Assume both are equal
	assets := []AssetType{}

	if tx.TranslatedTotalPrice != nil {

		for k := range tx.TranslatedTotalPrice {

			assets = append(assets, AssetType(k))

		}

		return assets
	}

	if tx.TranslatedFee != nil {

		for k := range tx.TranslatedFee {

			assets = append(assets, AssetType(k))

		}

	}

	return assets
}

func (tx *TransactionLog) GetExchange() string {
	return tx.Exchange
}

func (tx *TransactionLog) GetSide() SideType {
	return tx.Side
}

func (tx *TransactionLog) GetSideIdentifier() string {
	return tx.SideIdentifier
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
