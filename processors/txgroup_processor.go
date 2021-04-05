package processors

import (
	"time"

	"github.com/mariotoffia/gocryptoadmin/common"
	"github.com/mariotoffia/gocryptoadmin/processors/procutils"
)

type TxGroupProcessor struct {
	cache        *procutils.TxGroupCache
	transactions []common.TxGroupEntry
	processed    bool
}

func NewTxGroupProcessor() *TxGroupProcessor {
	return &TxGroupProcessor{
		cache:        procutils.NewTxGroupCache(time.Duration(5 * 60)),
		transactions: []common.TxGroupEntry{},
	}
}

// Processes a transaction.
//
// It may hold the transaction in its cache if e.g. a merge operation is currently done. E.g.
// using a group window.
//
// When `Transaction` instances is in this cache, they are said to be open. If `Flush` is invoked,
// they are unconditionally merged and written to the underlying store.
//
// A `Processor` may merge `Transaction` instances as long as the following criteria is fulfilled
//
// 1. Within Group Window (if any)
// 2. "Open `AssetPair` Transaction" - i.e. it is in a cache and not yet written to underlying _"store"_
// 3. The new transaction, with same `AssetPair` do have same `SideType`
// 4. The Asset part of the Open `Transaction` is not part of a `CostUnit` in the new `Transaction`
//
// If any of the above bullets fail, all _"Open"_ `Transaction` instances should be merged.
func (txg *TxGroupProcessor) Process(tx common.TransactionLog) {

	txg.processed = true

	item, ok := txg.cache.GetCache(tx)

	// 2. "Open `AssetPair` Transaction" - i.e. it is in a cache
	//     and not yet written to underlying _"store"_
	if !ok {

		txg.cache.CreateCacheAddTx(tx)

		// 3. The new transaction, with same `AssetPair` do have same `SideType`
		if other, ok := txg.cache.GetOtherSide(tx); ok {

			if other.IsOpen() {
				// Close the other since
				txg.transactions = append(txg.transactions, txg.cache.FlushCache(other))
			}

		}

		return

	}

	// 1. Within Group Window (if any)
	if !item.WithinWindow(tx) {

		txg.transactions = append(txg.transactions, txg.cache.FlushCache(item))
		txg.cache.CreateCacheAddTx(tx)
		return

	}

	// 4. The Asset part of the Open `Transaction` is not part of
	//    a `CostUnit` in the new `Transaction`
	others := txg.cache.GetOthersWithSameCostUnit(tx)
	if len(others) > 0 {

		// Need to close them since this tx is buying or selling
		// same asset as been used in a cost unit in others, i.e.
		// balance of this asset is either increasing or declining.
		for i := range others {
			txg.transactions = append(txg.transactions, txg.cache.FlushCache(others[i]))
		}

	}

	// Add the transaction to cache
	txg.cache.AddTransactionToCache(tx)
}

func (txg *TxGroupProcessor) Flush() {
	txg.transactions = append(txg.transactions, txg.cache.FlushAllCaches()...)
}

func (txg *TxGroupProcessor) UseGroupWindow(s int64) {

	if txg.processed {

		panic("cannot set group windows after process has been invoked!")

	}

	txg.cache = procutils.NewTxGroupCache(time.Duration(s))

}
