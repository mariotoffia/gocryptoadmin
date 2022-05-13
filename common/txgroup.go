package common

import (
	"math"
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
	GetTransactionEntries() []TransactionEntry
	// AddTransactionEntry adds a single entry to the `TransactionLog` array
	AddTransactionEntry(tx TransactionEntry) *TxGroupEntry
	// GetMostProminentSizeTransactionLog get the `TransactionLog` of largest `AssetSize`
	// in the all of the `TransactionLog` instances.
	GetMostProminentSizeTransactionLog() TransactionEntry
}
type TxGroupEntry struct {
	// Not used, just derivation
	TransactionLog
	Tx []TransactionEntry `json:"logs"`
}

func (txg *TxGroupEntry) GetTransactionEntries() []TransactionEntry {

	return txg.Tx

}

func (txg *TxGroupEntry) AddTransactionEntry(tx TransactionEntry) *TxGroupEntry {

	txg.Tx = append(txg.Tx, tx)
	return txg

}

func (txg *TxGroupEntry) GetExchange() string {

	if len(txg.Tx) == 0 {
		return ""
	}

	return txg.Tx[0].GetExchange()

}
func (txg *TxGroupEntry) GetSide() SideType {

	if len(txg.Tx) == 0 {
		return SideTypeUnknown
	}

	return txg.Tx[0].GetSide()

}

func (txg *TxGroupEntry) GetSideIdentifier() string {

	if len(txg.Tx) == 0 {
		return ""
	}

	return txg.Tx[0].GetSideIdentifier()

}

func (txg *TxGroupEntry) GetCreatedAt() time.Time {

	if len(txg.Tx) == 0 {
		return time.Time{}
	}

	return txg.Tx[0].GetCreatedAt()

}

func (txg *TxGroupEntry) GetAssetSize() float64 {

	if len(txg.Tx) == 0 {
		return 0
	}

	return linq.From(txg.Tx).
		Select(func(tx interface{}) interface{} {
			return tx.(TransactionEntry).GetAssetSize()
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

		weightedSize := b.GetAssetSize() / totalSize
		price = utils.ToFixed(price+(b.GetPricePerUnit()*weightedSize), 8)

	}

	return price

}

func (txg *TxGroupEntry) GetFee() float64 {

	if len(txg.Tx) == 0 {
		return 0
	}

	return linq.From(txg.Tx).
		Select(func(tx interface{}) interface{} {
			return tx.(TransactionEntry).GetFee()
		}).
		SumFloats()

}

func (txg *TxGroupEntry) GetTotalPrice() float64 {

	if len(txg.Tx) == 0 {
		return 0
	}

	return linq.From(txg.Tx).
		Select(func(tx interface{}) interface{} {
			return tx.(TransactionEntry).GetTotalPrice()
		}).
		SumFloats()

}

func (txg *TxGroupEntry) GetTranslatedTotalPrice(asset AssetType) float64 {

	if len(txg.Tx) == 0 {
		return 0
	}

	return linq.From(txg.Tx).
		Select(func(tx interface{}) interface{} {

			return tx.(TransactionEntry).GetTranslatedTotalPrice(asset)

		}).
		SumFloats()

}

func (txg *TxGroupEntry) GetTranslatedFee(asset AssetType) float64 {

	if len(txg.Tx) == 0 {
		return 0
	}

	return linq.From(txg.Tx).
		Select(func(tx interface{}) interface{} {

			return tx.(TransactionEntry).GetTranslatedFee(asset)

		}).
		SumFloats()

}

func (txg *TxGroupEntry) GetTranslatedAssets() []AssetType {

	if len(txg.Tx) == 0 {
		return []AssetType{}
	}

	return txg.Tx[0].GetTranslatedAssets()

}

func (txg *TxGroupEntry) GetAssetPair() AssetPair {

	if len(txg.Tx) == 0 {
		return AssetPair{}
	}

	return txg.Tx[0].GetAssetPair()

}

func (txg *TxGroupEntry) GetMostProminentSizeTransactionLog() TransactionEntry {

	if len(txg.Tx) == 0 {
		return &TransactionLog{}
	}

	max := float64(0)
	found := 0

	for i := range txg.Tx {

		if txg.Tx[i].GetAssetSize() > max {
			max = txg.Tx[i].GetAssetSize()
			found = i
		}
	}

	return txg.Tx[found]

}

func (txg *TxGroupEntry) Clone() TransactionEntry {

	e := &TxGroupEntry{TransactionLog: txg.TransactionLog}

	if len(txg.Tx) > 0 {

		for i := range txg.Tx {
			e.Tx = append(e.Tx, txg.Tx[i].Clone())
		}
	}

	return e

}

func (txg *TxGroupEntry) SplitSize(
	size float64,
) (sized TransactionEntry, overflow TransactionEntry) {

	szd := &TxGroupEntry{TransactionLog: txg.TransactionLog}
	ofl := &TxGroupEntry{TransactionLog: txg.TransactionLog}

	found := txg.FindBySize(size, true /*closest*/)

	if found == -1 {

		// Brute-force, add entries until sized (may need to split last entry) -> sized
		// Rest (including overflow on last split - if needed) -> overflow
		for i, entry := range txg.Tx {

			size = utils.ToFixed(size-entry.GetAssetSize(), 8)

			if size == 0 {

				// We're done
				szd.Tx = append(szd.Tx, entry)
				ofl.Tx = append(ofl.Tx, txg.Tx[i+1:]...)

				break

			}

			if size < 0 {
				// We're done, but need to split this entry
				entryoverflow, entrysized := entry.SplitSize(-size)

				szd.Tx = append(szd.Tx, entrysized.(TransactionEntry))

				ofl.Tx = append(ofl.Tx, entryoverflow.(TransactionEntry))
				ofl.Tx = append(ofl.Tx, txg.Tx[i+1:]...)

				break

			}

			szd.Tx = append(szd.Tx, entry)

		}

		return szd, ofl
	}

	entry := txg.Tx[found]

	if entry.GetAssetSize() == size {

		szd.Tx = []TransactionEntry{entry}
		ofl.Tx = append(ofl.Tx, txg.Tx[:found]...)
		ofl.Tx = append(ofl.Tx, txg.Tx[found+1:]...)
		return szd, ofl

	}

	// Not exact match -> split entry use sized -> sized
	// Use overflow on split + all others in tx -> overflow
	entrysized, entryoverflow := entry.SplitSize(size)

	szd.Tx = []TransactionEntry{entrysized.(TransactionEntry)}

	ofl.Tx = append(ofl.Tx, txg.Tx[:found]...)
	ofl.Tx = append(ofl.Tx, entryoverflow.(TransactionEntry))
	ofl.Tx = append(ofl.Tx, txg.Tx[found+1:]...)

	return szd, ofl

}

func (txg *TxGroupEntry) FindBySize(size float64, closest bool) int {

	idx := -1
	check := math.MaxFloat64

	for i, entry := range txg.Tx {

		if entry.GetAssetSize() == size {
			return i
		}

		if !closest {
			continue
		}

		approx := entry.GetAssetSize() - size

		// Only entries larger than size chan be closest
		if approx < 0 {
			continue
		}

		if approx < check {
			check = approx
			idx = i
		}

	}

	if !closest {
		return -1
	}

	return idx

}
