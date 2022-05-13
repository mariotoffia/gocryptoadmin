package processors

import (
	"fmt"
	"time"

	"github.com/mariotoffia/gocryptoadmin/common"
	"github.com/mariotoffia/gocryptoadmin/processors/procutils"
)

type TxGroupProcessor struct {
	cache          *procutils.TxGroupCache
	transactions   []common.TxGroupEntry
	timewindow     time.Duration
	flushProcessor common.TxGroupProcessor
	cachekeys      procutils.CacheKeyRenderer
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

	render := func(tx common.TransactionEntry, side common.SideType) string {
		return fmt.Sprintf(
			"%s_%s_%s", tx.GetExchange(), tx.GetAssetPair(), side,
		)
	}

	return &TxGroupProcessor{
		cache:          procutils.NewTxGroupCache(timewindow, render),
		transactions:   []common.TxGroupEntry{},
		timewindow:     timewindow,
		flushProcessor: NewChronologicalGroupTxEntryProcessor(),
		cachekeys:      render,
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
	txg.cache = procutils.NewTxGroupCache(txg.timewindow, txg.cachekeys)

}

// Processes a single transaction entry.
//
// It will hold the transaction in its cache until using the following steps.
//
// 1. Is New Tx -> Check if exists asset with either cost or asset type that is not the exact
// AssetPair as new tx. (e.g. new: LTC-EUR, in cache LTC-LTC, the latter will then close). It will check
// both asset and cost unit (only non _FIAT_!).
//
// 2. Is _TRANSFER_ or _RECEIVE_ and not new Tx, it will close any older of same kind and `AssetPair`. If
// _TRANSFER_ it will check if any `AssetType` same as _TRANSFER_ `AssetType`. It will terminate those.
// _RECEIVE_, it will check if any `AssetType` same as _RECEIVE_ `AssetType`. It will terminate those.
//
// 3. Within Group Window (if any)
//
// If any of the above bullets fail, all _"Open"_ `Transaction` instances should be merged.
func (txg *TxGroupProcessor) Process(tx common.TransactionEntry) {

	tx = tx.Clone()
	item, ok := txg.cache.GetCache(tx, nil /*no override*/)

	if !ok {

		// 1.
		if items, found := txg.cache.GetAssetPairWhenNonFIAT(tx.GetAssetPair()); found {

			for i := range items {

				if items[i].IsOpen() {
					// Close the other asset pairs, since containing a crypto assset / cost unit of this _tx_.
					txg.transactions = append(txg.transactions, txg.cache.FlushCache(items[i]))

				}

			}

		} else if tx.GetSide() == common.SideTypeReceive || tx.GetSide() == common.SideTypeTransfer {

			// 2. Is _TRANSFER_ or _RECEIVE_...
			txg.cache.CreateCacheAddTx(tx)

			if items, ok := txg.cache.GetByExchangeAssetType(tx.GetExchange(), tx.GetAssetPair().Asset); ok {

				for i := range items {

					if items[i].IsOpen() {
						txg.transactions = append(txg.transactions, txg.cache.FlushCache(items[i]))

					}

				}

			}

			return

		}

		txg.cache.CreateCacheAddTx(tx)
		return

	}

	if !item.WithinWindow(tx) {

		txg.transactions = append(txg.transactions, txg.cache.FlushCache(item))
		txg.cache.CreateCacheAddTx(tx)
		return

	}

	// 3. Within Group Window (if any)
	txg.cache.AddTransactionToCache(tx)
}

// ProcessMany is calling `Process` by iterating the _tx_ array
func (txg *TxGroupProcessor) ProcessMany(tx []common.TransactionEntry) {

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
