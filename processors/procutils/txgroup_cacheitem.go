package procutils

import (
	"time"

	"github.com/mariotoffia/gocryptoadmin/common"
)

// TxGroupCacheItem is a single item in the `TxGroupCache` that holds a set of transactions logs.
type TxGroupCacheItem struct {
	next time.Time
	tx   common.TxGroupEntry
}

// WithinWindow checks if the _tx_ is within the current window in the cache item.
func (txi *TxGroupCacheItem) WithinWindow(tx common.TransactionEntry) bool {

	return tx.GetCreatedAt().Before(txi.next)

}

// IsOpen returns `true` if there are any entries in the `Tx.Tx` field, i.e. any transactions in the
// group.
func (txi *TxGroupCacheItem) IsOpen() bool {

	return len(txi.tx.Tx) != 0

}
