package procutils

import (
	"fmt"
	"time"

	"github.com/ahmetb/go-linq/v3"
	"github.com/mariotoffia/gocryptoadmin/common"
)

type TxGroupCacheItem struct {
	next time.Time
	tx   common.TxGroupEntry
}

// WithinWindow checks if the _tx_ is within the current window in the cache item.
func (txi *TxGroupCacheItem) WithinWindow(tx common.TransactionLog) bool {

	return tx.CreatedAt.Before(txi.next)

}

// IsOpen returns `true` if there are any entries in the `Tx.Tx` field, i.e. any transactions in the
// group.
func (txi *TxGroupCacheItem) IsOpen() bool {

	return len(txi.tx.Tx) != 0

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

// FlushCache will extract the `common.TxGroupEntry` and remove the cache completely.
func (txg *TxGroupCache) FlushCache(item *TxGroupCacheItem) common.TxGroupEntry {

	tx := item.tx
	delete(txg.cache, tx.GetExchange()+tx.GetAssetPair().String()+string(tx.GetSide()))
	return tx
}

func (txg *TxGroupCache) FlushAllCaches() []common.TxGroupEntry {

	tx := []common.TxGroupEntry{}
	for _, cache := range txg.cache {

		if cache.IsOpen() {
			tx = append(tx, cache.tx)
		}

	}

	txg.cache = map[string]*TxGroupCacheItem{}

	return tx
}

// GetOtherSide is same as `GetCache` except that it will get the inverse `SideType`
// of the transaction. It will *panic* if the `SideType` is unknown.
//
// .Other Side
// ====
// 1. tx side is _BUY_ -> it will look for _SELL_
// 2. tx side is _RECEIVE_ -> it will look for _TRANSFER_.
// ====
func (txg *TxGroupCache) GetOtherSide(
	tx common.TransactionLog,
) (item *TxGroupCacheItem, found bool) {

	var side common.SideType
	switch tx.Side {
	case common.SideTypeBuy:
		side = common.SideTypeSell
	case common.SideTypeSell:
		side = common.SideTypeBuy
	case common.SideTypeReceive:
		side = common.SideTypeTransfer
	case common.SideTypeTransfer:
		side = common.SideTypeReceive
	default:
		panic(
			fmt.Sprintf("not supported sidetype in get other side operation: %s", string(tx.Side)),
		)
	}

	item, found = txg.cache[tx.Exchange+tx.AssetPair.String()+string(side)]
	return
}

// GetOthersWithSameCostUnit will return all caches that do have transactions with same
// `CostUnit` as the `tx.Asset`.
//
// Usually, the cost unit is in EUR or $ but it may be e.g. LTC, BTC and if e.g. _tx_ do have
// `Asset` of _BTC_, and there exist a few cache items with `CostUnit` of _BTC_ (e.g. _XRP-BTC_),
// those will be returned.
func (txg *TxGroupCache) GetOthersWithSameCostUnit(tx common.TransactionLog) []*TxGroupCacheItem {

	var res []*TxGroupCacheItem

	linq.From(txg.cache).
		Where(func(kv interface{}) bool {

			c := kv.(linq.KeyValue).Value.(*TxGroupCacheItem)
			e := c.tx

			return e.GetExchange() == tx.GetExchange() &&
				e.GetAssetPair().CostUnit == tx.GetAssetPair().Asset &&
				c.IsOpen()

		}).
		Select(func(kv interface{}) interface{} {
			return kv.(linq.KeyValue).Value
		}).
		ToSlice(&res)

	return res
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
		next: tx.CreatedAt.Add(time.Second * txg.secwindow),
	}

	txg.groupId++

	item.tx.ID = fmt.Sprint(txg.groupId)
	item.tx.AddTransactionEntry(tx)

	txg.cache[tx.Exchange+tx.AssetPair.String()+string(tx.Side)] = item

	return item

}

// AddTransactionToCache will add the _tx_ to an existing cache item.
//
// If the cache item do not exist, it will panic.
func (txg *TxGroupCache) AddTransactionToCache(tx common.TransactionLog) {

	if item, ok := txg.cache[tx.Exchange+tx.AssetPair.String()+string(tx.Side)]; ok {

		item.tx.AddTransactionEntry(tx)
		return
	}

	panic(
		fmt.Sprintf(
			"no cache item created for: %s (cannot add transaction log)",
			tx.Exchange+tx.AssetPair.String()+string(tx.Side),
		),
	)
}
