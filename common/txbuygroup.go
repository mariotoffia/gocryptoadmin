package common

import (
	"fmt"

	"github.com/ahmetb/go-linq/v3"
	"github.com/mariotoffia/gocryptoadmin/utils"
)

type TxLBuyGroupEntry interface {
	// IsMultiAsset returns `true` if any of the `AssetType`.
	//
	// Depending on entry `SideType.Buy` or `SideType.Sell` it will
	// use the asset or cost unit to determine equality of `AssetType`
	// amongst the entries.
	IsMultiAsset() bool
}

// TxBuyGroupLog is exactly the same as `TxGroupEntry` but functions
// has been overridden to cope with _BUY_ *and* _SELL_ transactions in
// same group.
//
// This is due to that a sell may include buys that was _really_ a sell
// from one crypto currency to another, and that _"another"_ is now sold.
type TxBuyGroupLog struct {
	TxGroupEntry
	multi bool
}

func NewTxBuyGroupLog(entries []TransactionEntry) *TxBuyGroupLog {

	l := &TxBuyGroupLog{
		TxGroupEntry: TxGroupEntry{
			Tx: entries,
		},
	}

	if len(entries) > 0 {
		l.ID = fmt.Sprintf("%s-buygroup", entries[0].GetID())
	}

	l.IsMultiAsset() // Update cache
	return l
}

func (txg *TxBuyGroupLog) AddTransactionEntry(tx TransactionEntry) *TxGroupEntry {

	txg.Tx = append(txg.Tx, tx)

	txg.IsMultiAsset() // Update cache
	return &txg.TxGroupEntry

}

// IsMultiAsset checks if all transactions do have the same
// `AssetType`. If not it will return `false`.
//
// NOTE: When _SELL_ transactions it will check _CostUnit_ instead of _Asset_.
func (txg *TxBuyGroupLog) IsMultiAsset() bool {

	txg.multi = false
	if len(txg.Tx) == 0 {
		return false
	}

	var asset AssetType
	if txg.Tx[0].GetSide() == SideTypeBuy {
		asset = txg.Tx[0].GetAssetPair().Asset
	} else {
		asset = txg.Tx[0].GetAssetPair().CostUnit // Sell
	}

	txg.iterate(
		func(entry TransactionEntry, side SideType, adjsize float64) bool {

			if side == SideTypeBuy && entry.GetAssetPair().Asset == asset {
				return true
			}

			if entry.GetAssetPair().CostUnit == asset {
				return true
			}

			txg.multi = true
			return false
		})

	return txg.multi

}

// GetAssetSize will iterate all entries and sum the asset sizes.
//
// If the direction is BUY it will use `GetAssetSize`, if it is
// sell it will use `GetTotalPrice` since it is a _BUY_ of a
// crypto currency that has been sold now (in `TxBuySellLog`).
func (txg *TxBuyGroupLog) GetAssetSize() float64 {

	if len(txg.Tx) == 0 {
		return 0
	}

	size := float64(0)

	txg.iterate(
		func(entry TransactionEntry, side SideType, adjsize float64) bool {
			size += adjsize
			return true
		})

	return size

}

// GetPricePerUnit returns the weighetd price per unit for all entires.
//
// Since _SELL_ transactions have incorrect cost unit, it will calculate
// the price using `GetAssetSize() / (GetTotalPrice() - GetFee())`.
//
// CAUTION: It does not support multi asset entries!
func (txg *TxBuyGroupLog) GetPricePerUnit() float64 {

	if len(txg.Tx) == 0 {
		return 0
	}

	if txg.multi {
		panic("cannot get price per unit on multi entry")
	}

	price := float64(0)
	totalSize := txg.GetAssetSize()

	txg.iterate(
		func(entry TransactionEntry, side SideType, adjsize float64) bool {

			weightedSize := adjsize / totalSize

			var ppe float64
			if side == SideTypeBuy {
				ppe = entry.GetPricePerUnit()
			} else {
				ppe = entry.GetAssetSize() / (entry.GetTotalPrice() - entry.GetFee())
			}

			price += utils.ToFixed(price+(ppe*weightedSize), 8)
			return true
		})

	return price

}

// GetFee returns the sum of all fees.
//
// Since _SELL_ entries do have fee in incorrect cost unit,
// it is calculated using  `GetFee() / GetPricePerUnit()`.
//
// CAUTION: It does not support multi asset entries!
func (txg *TxBuyGroupLog) GetFee() float64 {

	if len(txg.Tx) == 0 {
		return 0
	}

	if txg.multi {
		panic("cannot get fee on multi entry")
	}

	var fee float64
	txg.iterate(
		func(entry TransactionEntry, side SideType, adjsize float64) bool {

			if side == SideTypeBuy {
				fee += entry.GetFee()
			} else {
				fee += entry.GetFee() / entry.GetPricePerUnit()
			}

			return true
		})

	return fee

}

// GetTotalPrice calculates the total price.
//
// On _BUY_ it uses the `GetTotalPrice` but on _SELL_,
// it uses the `GetAssetSize` since that denotes the
// price.
//
// CAUTION: It does not support multi asset entries!
func (txg *TxBuyGroupLog) GetTotalPrice() float64 {

	if len(txg.Tx) == 0 {
		return 0
	}

	if txg.multi {
		panic("cannot get total price on multi entry")
	}

	var price float64
	txg.iterate(
		func(entry TransactionEntry, side SideType, adjsize float64) bool {

			if side == SideTypeBuy {
				price += entry.GetTotalPrice()
			} else {
				price += entry.GetAssetSize()
			}

			return true
		})

	return price
}

// GetTranslatedTotalPrice is overloaded due to that we need to set the sell total
// price as negative since it is, in this `TxBuyGroup` counted as it where a sort of a _BUY_.
func (txg *TxBuyGroupLog) GetTranslatedTotalPrice(asset AssetType) float64 {

	if len(txg.Tx) == 0 {
		return 0
	}

	return linq.From(txg.Tx).
		Select(func(tx interface{}) interface{} {

			entry := tx.(TransactionEntry)

			if entry.GetSide() == SideTypeBuy {
				return entry.GetTranslatedTotalPrice(asset)
			}
			return -entry.GetTranslatedTotalPrice(asset)

		}).
		SumFloats()

}

func (txg *TxBuyGroupLog) GetMostProminentSizeTransactionLog() TransactionEntry {

	if len(txg.Tx) == 0 {
		return &TransactionLog{}
	}

	max := float64(0)
	var found TransactionEntry

	txg.iterate(
		func(entry TransactionEntry, side SideType, adjsize float64) bool {

			if adjsize > max {
				max = adjsize
				found = entry
			}

			return true
		})

	return found

}

func (txg *TxBuyGroupLog) Clone() TransactionEntry {

	e := &TxBuyGroupLog{
		TxGroupEntry: TxGroupEntry{
			TransactionLog: txg.TxGroupEntry.TransactionLog,
		},
		multi: txg.multi,
	}

	if len(txg.Tx) > 0 {

		for i := range txg.Tx {
			e.Tx = append(e.Tx, txg.Tx[i].Clone())
		}
	}

	return e

}

func (txg *TxBuyGroupLog) iterate(
	processor func(entry TransactionEntry, side SideType, adjsize float64) bool,
) {

	for i := range txg.Tx {

		entry := txg.Tx[i].(TransactionEntry)
		side := entry.GetSide()
		adjsize := float64(0)

		if side == SideTypeSell {
			adjsize = entry.GetTotalPrice()
		} else if side == SideTypeBuy {
			adjsize = entry.GetAssetSize()
		} else {
			panic("sell or buy expected")
		}

		if !processor(entry, side, adjsize) {
			break
		}
	}

}
