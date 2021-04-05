package processors

import (
	"fmt"
	"time"

	"github.com/mariotoffia/gocryptoadmin/common"
	"github.com/mariotoffia/gocryptoadmin/processors/procutils"
)

type TxGroupProcessor struct {
	cache     *procutils.TxGroupCache
	processed bool
}

func NewTxGroupProcessor() *TxGroupProcessor {
	return &TxGroupProcessor{
		cache: procutils.NewTxGroupCache(time.Duration(5 * 60)),
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
// 2. "Open `AssetPair` Transaction" - i.e. it is in a cache and not yet registered in the ledger
// 3. The new transaction, with same `AssetPair` do have same `SideType`
// 4. The Asset part of the Open `Transaction` is not part of a `CostUnit` in the new `Transaction`
//
// If any of the above bullets fail, all _"Open"_ `Transaction` instances should be merged.
func (txg *TxGroupProcessor) Process(tx common.TransactionLog) {

	txg.processed = true

	item, ok := txg.cache.GetCache(tx)

	if !ok {

		txg.cache.CreateCacheAddTx(tx)
		return

	}

	fmt.Println(item)

}

func (txg *TxGroupProcessor) Flush() {
	// TODO: need to do something here?
}

func (txg *TxGroupProcessor) UseGroupWindow(s int64) {

	if txg.processed {

		panic("cannot set group windows after process has been invoked!")

	}

	txg.cache = procutils.NewTxGroupCache(time.Duration(s))

}
