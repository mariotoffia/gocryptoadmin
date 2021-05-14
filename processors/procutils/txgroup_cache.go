package procutils

import (
	"fmt"
	"time"

	"github.com/ahmetb/go-linq/v3"
	"github.com/mariotoffia/gocryptoadmin/common"
)

// NewTxGroupCache creates a new group with specified _timewindow_.
//
// The _timewindow_ specified the max time that this group would capture
// in terms of `common.TransactionLog.CreatedAt`.
func NewTxGroupCache(timewindow time.Duration) *TxGroupCache {

	return &TxGroupCache{
		cache:      map[string]*TxGroupCacheItem{},
		timewindow: timewindow,
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

// FlushAllCaches iterates the cache and returns all `TxGroupCacheItem.IsOpen` entries.
//
// Before returning, it will clear the complete cache, hence empty.
func (txg *TxGroupCache) FlushAllCaches() []common.TxGroupEntry {

	tx := []common.TxGroupEntry{}
	for _, cache := range txg.cache {

		// Only items that do have transactions left in it
		if cache.IsOpen() {
			tx = append(tx, cache.tx)
		}

	}

	txg.cache = map[string]*TxGroupCacheItem{}

	return tx
}

// GetAllOtherSide will return all sides found in cache, except, in param side.
func (txg *TxGroupCache) GetAllOtherSide(
	tx common.TransactionLog,
	side common.SideType,
) (items []*TxGroupCacheItem, found bool) {

	other := []common.SideType{
		common.SideTypeBuy,
		common.SideTypeSell,
		common.SideTypeReceive,
		common.SideTypeTransfer,
	}

	for i := 0; i < 4; i++ {

		if other[i] == side {

			other = append(other[:i], other[i+1:]...)
			break
		}

	}

	items = []*TxGroupCacheItem{}

	for _, s := range other {

		if item, found := txg.cache[tx.Exchange+tx.AssetPair.String()+string(s)]; found {
			items = append(items, item)
		}

	}

	return items, len(items) > 0

}

// GetOtherSide is same as `GetCache` except that it will get the inverse `SideType`
// of the transaction. It will *panic* if the `SideType` is unknown.
//
// .Other Side
// ====
// 1. tx side is _BUY_ -> it will look for _SELL_
// 2. tx side is _RECEIVE_ or _TRANSFER_ it will panic!
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
		fallthrough
	case common.SideTypeTransfer:
		fallthrough
	default:
		panic(
			fmt.Sprintf(
				"not supported sidetype in get other side operation: %s", string(tx.Side),
			),
		)
	}

	item, found = txg.cache[tx.Exchange+tx.AssetPair.String()+string(side)]
	return
}

// GetByExchangeCostUnit will return all `TxGroupCacheItem` instances that do have
// transactions with same `CostUnit` and exchange as the `tx.Asset`.
//
// Usually, the cost unit is in EUR or $ but it may be e.g. LTC, BTC and if e.g. _assetType_ is
// _BTC_, and there exist a few cache items with `CostUnit` of _BTC_ (e.g. _XRP-BTC_),
// those will be returned.
func (txg *TxGroupCache) GetByExchangeCostUnit(
	exchange string,
	assetType common.AssetType,
) []*TxGroupCacheItem {

	var res []*TxGroupCacheItem

	linq.From(txg.cache).
		Where(func(kv interface{}) bool {

			c := kv.(linq.KeyValue).Value.(*TxGroupCacheItem)
			e := c.tx

			return e.GetExchange() == exchange &&
				e.GetAssetPair().CostUnit == assetType &&
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
		next: tx.CreatedAt.Add(txg.timewindow),
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
