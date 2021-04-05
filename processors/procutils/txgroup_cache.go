package procutils

import (
	"fmt"
	"time"

	"github.com/mariotoffia/gocryptoadmin/common"
)

type TxGroupCacheItem struct {
	Next time.Time
	Tx   common.TxGroupEntry
}

type TxGroupCache struct {
	cache     map[string]*TxGroupCacheItem
	groupId   int64
	secwindow time.Duration
}

func NewTxGroupCache(secwindow time.Duration) *TxGroupCache {

	return &TxGroupCache{
		cache:     map[string]*TxGroupCacheItem{},
		secwindow: secwindow,
	}

}

// GetCache gets the cache item associated with the transaction type.
//
// It uses the _Exchange_, _AssetPair_, and _Side_ to key the cache item.
func (txg *TxGroupCache) GetCache(
	tx common.TransactionLog,
) (item *TxGroupCacheItem, found bool) {

	item, found = txg.cache[tx.Exchange+tx.AssetPair.String()+string(tx.Side)]
	return

}

// CreateCacheAddTx will create a chache item and add the current transaction to it.
//
// It will panic if the cache item is already created.
func (txg *TxGroupCache) CreateCacheAddTx(tx common.TransactionLog) *TxGroupCacheItem {

	if _, ok := txg.cache[tx.Exchange+tx.AssetPair.String()+string(tx.Side)]; ok {

		panic(
			fmt.Sprintf(
				"cache item already created for: %s",
				tx.Exchange+tx.AssetPair.String()+string(tx.Side),
			),
		)

	}

	item := &TxGroupCacheItem{
		Next: tx.CreatedAt.Add(time.Second * txg.secwindow),
	}

	txg.groupId++

	item.Tx.ID = fmt.Sprint(txg.groupId)
	item.Tx.AddTransactionEntry(tx)

	txg.cache[tx.Exchange+tx.AssetPair.String()+string(tx.Side)] = item

	return item

}

// AddTransactionToCache will add the _tx_ to an existing cache item.
//
// If the cache item do not exist, it will panic.
func (txg *TxGroupCache) AddTransactionToCache(tx common.TransactionLog) {

	if item, ok := txg.cache[tx.Exchange+tx.AssetPair.String()+string(tx.Side)]; ok {

		item.Tx.AddTransactionEntry(tx)
		return
	}

	panic(
		fmt.Sprintf(
			"no cache item created for: %s (cannot add transaction log)",
			tx.Exchange+tx.AssetPair.String()+string(tx.Side),
		),
	)
}
