package processors

import (
	"time"

	"github.com/mariotoffia/gocryptoadmin/common"
	"github.com/mariotoffia/gocryptoadmin/processors/procutils"
)

type TxGroupProcessor struct {
	cache          *procutils.TxGroupCache
	transactions   []common.TxGroupEntry
	timewindow     time.Duration
	flushProcessor common.TxGroupProcessor
}

// NewTxGroupProcessor creates a new processor with windowsize of 5 minutes.
//
// This `Processor` _REQUIRES_ that the log entries are ordered in chronological
// order since it will not sort entries.
//
// By default it assigns a `ChronologicalProcessor`, this can be cleared by
// setting `UseFlushProcessor(nil)`.
//
// If _timewindow_ is equal or less than zero, it will default to 5 minutes.
func NewTxGroupProcessor(timewindow time.Duration) *TxGroupProcessor {

	if timewindow <= 0 {
		timewindow = time.Duration(time.Minute * 5)
	}

	return &TxGroupProcessor{
		cache:          procutils.NewTxGroupCache(timewindow),
		transactions:   []common.TxGroupEntry{},
		timewindow:     timewindow,
		flushProcessor: NewChronologicalGroupTxEntryProcessor(),
	}

}

// UseFlushProcessor sets a new flush processor (executed during `Flush` operation).
//
// If _flushProcessor_ is `nil` it will remove the processor and nothing is processed before
// the records are returned during `Flush`. This may be useful when external processing is much
// more efficient.
//
// Since, `TxGroupProcessor` uses `ChronologicalTxEntryProcessor` by default, if sorting in `Flush`
// operation is wanted, just set it to `nil`.
func (txg *TxGroupProcessor) UseFlushProcessor(
	flushProcessor common.TxGroupProcessor,
) *TxGroupProcessor {

	txg.flushProcessor = flushProcessor
	return txg

}

// Resets clears the processor to start from scratch.
func (txg *TxGroupProcessor) Reset() {

	txg.transactions = []common.TxGroupEntry{}
	txg.cache = procutils.NewTxGroupCache(txg.timewindow)

}

// Processes a single transaction entry.
//
// It will hold the transaction in its cache until using the following steps.
//
// 1. Within Group Window (if any)
// 2. "Open `AssetPair` Transaction" - i.e. it is in a cache and not yet written to underlying _"store"_
// 3. The new transaction, with same `AssetPair` do have same `SideType`
// 4. The Asset part of the Open `Transaction` is not part of a `CostUnit` in the new `Transaction`
//
// If any of the above bullets fail, all _"Open"_ `Transaction` instances should be merged.
func (txg *TxGroupProcessor) Process(tx common.TransactionLog) {

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

	// 4. The Asset part of the Open `Transaction` is not part of a `CostUnit` in the new `Transaction`
	//
	// Since _tx_ is not  yet in _txg_ cache, it will find "all others" with same CostUnit.
	others := txg.cache.GetByExchangeCostUnit(
		tx.GetExchange(), tx.GetAssetPair().Asset,
	)

	if len(others) > 0 {

		// Need to close them since this tx is buying or selling same asset as been used in a cost unit in others, i.e.
		// balance of this asset is either increasing or declining.
		for i := range others {
			txg.transactions = append(txg.transactions, txg.cache.FlushCache(others[i]))
		}

	}

	// Add the transaction to cache
	txg.cache.AddTransactionToCache(tx)
}

// ProcessMany is calling `Process` by iterating the _tx_ array
func (txg *TxGroupProcessor) ProcessMany(tx []common.TransactionLog) {

	for i := range tx {

		txg.Process(tx[i])

	}

}

// Flush will flush all caches and runs the flush processor (if any attached).
//
// The flushprocessor will be invoked by `TxGroupProcessor.ProcessMany`, after
// it has been `TxGroupProcessor.Reset`. Lastly it fill perform a `TxGroupProcessor.Flush`
// on it.
func (txg *TxGroupProcessor) Flush() []common.TxGroupEntry {

	txg.transactions = append(txg.transactions, txg.cache.FlushAllCaches()...)

	if txg.flushProcessor == nil {

		return txg.transactions

	}

	txg.flushProcessor.Reset()
	txg.flushProcessor.ProcessMany(txg.transactions)

	return txg.flushProcessor.Flush()

}
