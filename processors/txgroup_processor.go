package processors

import (
	"fmt"
	"time"

	"github.com/mariotoffia/gocryptoadmin/common"
)

type TxGroupProcessor struct {
	cache *txgroupTxCache
}

func NewTxGroupProcessor() *TxGroupProcessor {
	return &TxGroupProcessor{
		cache: &txgroupTxCache{
			cache:     map[string]*txgroupTxCacheItem{},
			secwindow: time.Duration(5 * 60),
		},
	}
}

// Processes a transaction.
//
// It may hold the transaction in its cache if e.g. a merge operation is currently done. E.g.
// using a group window.
//
// When `Transaction` instances is in this cache, they are said to be open. If `Flush` is invoked,
// they are unconditionally merged and written to the ledger.
//
// A `Ledger` may merge `Transaction` instances as long as the following criteria is fulfilled
//
// 1. Within Group Window (if any)
// 2. "Open `AssetPair` Transaction" - i.e. it is in a cache and not yet registered in the ledger
// 3. The new transaction, with same `AssetPair` do have same `SideType`
// 4. The Asset part of the Open `Transaction` is not part of a `CostUnit` in the new `Transaction`
//
// If any of the above bullets fail, all _"Open"_ `Transaction` instances should be merged.
//
// NOTE: This function _REQUIRES_ that the transactions are in chronological order!!
func (txg *TxGroupProcessor) Process(tx common.TransactionLog) {

	item := txg.cache.GetItem(tx)
	fmt.Println(item)
	// TODO: do the logic...

}

func (txg *TxGroupProcessor) Flush() {
}

func (txg *TxGroupProcessor) UseGroupWindow(s int64) {
	txg.cache.secwindow = time.Duration(s)
}
