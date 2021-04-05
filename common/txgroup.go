package common

import (
	"time"

	"github.com/ahmetb/go-linq/v3"
	"github.com/mariotoffia/gocryptoadmin/utils"
)

// TxLogGroup is a slice of the `TransactionLog` that have been grouped.
//
// IMPORTANT: The `ID` must be set explicitly on a entry, since it should
// represent a unique group _ID_.
//
// It is represented as a `TransactionEntry` where all properties are sum
// of all in the `TransactionLog` slice. The underlying `TransactionLog` instances
// may be obtained by the `GetTransactionEntries`.
//
// The only value that is not a sum, is the `GetPricePerUnit()`. Instead it is a weighted
// price where the each `TransactionLog` contributes with its price proportional to it's
// contributing `AssetSize`.
//
// .Two Logs Contributing with Different Price
// ====
// 1. Asset (1) Size: 0.3, Price: 87.00 EUR
// 2. Asset (2) Size: 0.6, Price: 90.00 EUR
// 3. Total Size: 0.9 =>
// 4. Asset 1 Contribution 0,3 / 0.9 = 1/3
// 5. Asset 2 Contribution 0,6 / 0.9 = 2/3
// 6  PricePer Unit (1,2) => 87 * 1/3 + 90 * 2/3 = 29 + 60 = 89
// 7. Verify Cost: 0,3 * 87 + 0,6 * 90 = 80,1 vs 0,9 * 89 = 80,1 i.e. match
// ====
type TxLogGroup interface {
	TransactionEntry
	// GetTransactionEntries returns all underlying `TransactionLog` instances that
	// is reflected in the `TransactionEntry` interface methods
	GetTransactionEntries() []TransactionLog
	// AddTransactionEntry adds a single entry to the `TransactionLog` array
	AddTransactionEntry(tx TransactionLog) *TxGroupEntry
	// GetMostProminentSizeTransactionLog get the `TransactionLog` of largest `AssetSize`
	// in the all of the `TransactionLog` instances.
	GetMostProminentSizeTransactionLog() TransactionLog
}
type TxGroupEntry struct {
	TransactionLog
	Tx []TransactionLog `csv:"-" json:"logs"`
}

func (txg *TxGroupEntry) GetTransactionEntries() []TransactionLog {

	return txg.Tx

}

func (txg *TxGroupEntry) AddTransactionEntry(tx TransactionLog) *TxGroupEntry {

	txg.Tx = append(txg.Tx, tx)
	return txg

}

func (txg *TxGroupEntry) GetExchange() string {

	if len(txg.Tx) == 0 {
		return ""
	}

	return txg.Tx[0].Exchange

}
func (txg *TxGroupEntry) GetSide() SideType {

	if len(txg.Tx) == 0 {
		return SideTypeUnknown
	}

	return txg.Tx[0].Side

}

func (txg *TxGroupEntry) GetCreatedAt() time.Time {

	if len(txg.Tx) == 0 {
		return time.Time{}
	}

	return txg.Tx[0].CreatedAt

}

func (txg *TxGroupEntry) GetAssetSize() float64 {

	if len(txg.Tx) == 0 {
		return 0
	}

	return linq.From(txg.Tx).
		Select(func(tx interface{}) interface{} {
			return tx.(TransactionLog).AssetSize
		}).
		SumFloats()

}

func (txg *TxGroupEntry) GetPricePerUnit() float64 {

	if len(txg.Tx) == 0 {
		return 0
	}

	price := float64(0)
	totalSize := txg.GetAssetSize()

	for _, b := range txg.Tx {

		weightedSize := b.AssetSize / totalSize
		price = utils.ToFixed(price+(b.PricePerUnit*weightedSize), 8)

	}

	return price

}

func (txg *TxGroupEntry) GetFee() float64 {

	if len(txg.Tx) == 0 {
		return 0
	}

	return linq.From(txg.Tx).
		Select(func(tx interface{}) interface{} {
			return tx.(TransactionLog).Fee
		}).
		SumFloats()

}

func (txg *TxGroupEntry) GetTotalPrice() float64 {

	if len(txg.Tx) == 0 {
		return 0
	}

	return linq.From(txg.Tx).
		Select(func(tx interface{}) interface{} {
			return tx.(TransactionLog).TotalPrice
		}).
		SumFloats()

}

func (txg *TxGroupEntry) GetAssetPair() AssetPair {
	if len(txg.Tx) == 0 {
		return AssetPair{}
	}

	return txg.Tx[0].AssetPair
}

func (txg *TxGroupEntry) GetMostProminentSizeTransactionLog() TransactionLog {

	if len(txg.Tx) == 0 {
		return TransactionLog{}
	}

	max := float64(0)
	found := 0

	for i := range txg.Tx {

		if txg.Tx[i].AssetSize > max {
			max = txg.Tx[i].AssetSize
			found = i
		}
	}

	return txg.Tx[found]

}
