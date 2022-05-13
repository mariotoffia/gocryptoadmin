package common

import (
	"time"

	"github.com/mariotoffia/gocryptoadmin/utils"
)

type SideType string

const (
	SideTypeUnknown  SideType = "UNKNOWN"
	SideTypeBuy      SideType = "BUY"
	SideTypeSell     SideType = "SELL"
	SideTypeReceive  SideType = "RECEIVE"
	SideTypeTransfer SideType = "TRANSFER"
	// SideTypeBuySell is a paired transaction where it will contain
	// both one or more transaction that represents a _BUY_ that
	// accumulates to one _SELL_ transaction.
	//
	// However, the buy side may contain both `SideTypeBuy` *and* `SideTypeSell`.
	// The latter is when a _SELL_ transaction resulted in some crypto currency
	// that later got sold for _FIAT_.
	SideTypeBuySell SideType = "BUYSELL"
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

	Clone() TransactionEntry
	// SplitSize will split the current `TransactionEntry` by creating one by _size_ and
	// the other _overflow_ with the rest. All data is recalculated on each side, _split_ and _overflow_
	// so adding up both will have the same sums as the current one.
	SplitSize(size float64) (sized TransactionEntry, overflow TransactionEntry)
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

// SplitSize will split the current `TransactionEntry` by creating one by _size_ and
// the other _overflow_ with the rest. All data is recalculated on each side, _split_ and _overflow_
// so adding up both will have the same sums as the current one.
func (tx *TransactionLog) SplitSize(
	size float64,
) (sized TransactionEntry, overflow TransactionEntry) {

	percent := size / tx.AssetSize

	szd := tx.Clone().(*TransactionLog)
	ofl := tx.Clone().(*TransactionLog)

	szd.AssetSize = utils.ToFixed(size, 8)
	ofl.AssetSize = utils.ToFixed((tx.AssetSize - size), 8)

	szd.Fee = utils.ToFixed(tx.Fee*percent, 8)
	ofl.Fee = utils.ToFixed(tx.Fee-szd.Fee, 8)
	szd.TotalPrice = utils.ToFixed(tx.TotalPrice*percent, 8)
	ofl.TotalPrice = utils.ToFixed(tx.TotalPrice-szd.TotalPrice, 8)

	if len(tx.TranslatedFee) > 0 {

		for k, v := range szd.TranslatedFee {

			fee := v * percent

			szd.TranslatedFee[k] = utils.ToFixed(fee, 8)
			ofl.TranslatedFee[k] = utils.ToFixed(v-fee, 8)

		}

	}

	if len(tx.TranslatedTotalPrice) > 0 {

		for k, v := range szd.TranslatedTotalPrice {

			fee := v * percent

			szd.TranslatedTotalPrice[k] = utils.ToFixed(fee, 8)
			ofl.TranslatedTotalPrice[k] = utils.ToFixed(v-fee, 8)

		}

	}

	return szd, ofl
}

func (tx *TransactionLog) Clone() TransactionEntry {

	l := &TransactionLog{
		ID:             tx.ID,
		Exchange:       tx.Exchange,
		Side:           tx.Side,
		SideIdentifier: tx.SideIdentifier,
		CreatedAt:      tx.CreatedAt,
		AssetSize:      tx.AssetSize,
		PricePerUnit:   tx.PricePerUnit,
		Fee:            tx.Fee,
		TotalPrice:     tx.TotalPrice,
		AssetPair: AssetPair{
			Asset:    tx.Asset,
			CostUnit: tx.CostUnit,
		},
	}

	if len(tx.TranslatedFee) > 0 {

		l.TranslatedFee = map[string]float64{}
		for k, v := range tx.TranslatedFee {
			l.TranslatedFee[k] = v
		}

	}

	if len(tx.TranslatedTotalPrice) > 0 {

		l.TranslatedTotalPrice = map[string]float64{}
		for k, v := range tx.TranslatedTotalPrice {
			l.TranslatedTotalPrice[k] = v
		}

	}

	return l
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
