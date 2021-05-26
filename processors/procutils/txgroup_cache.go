package procutils

import (
	"fmt"
	"time"

	"github.com/ahmetb/go-linq/v3"
	"github.com/mariotoffia/gocryptoadmin/common"
)

// CacheKeyRenderer renders a key to use when accessing the cache
type CacheKeyRenderer func(tx common.TransactionEntry, side common.SideType) string

type TxGroupCache struct {
	cache      map[string]*TxGroupCacheItem
	groupId    int64
	timewindow time.Duration
	render     CacheKeyRenderer
}

// NewTxGroupCache creates a new group with specified _timewindow_.
//
// The _timewindow_ specified the max time that this group would capture
// in terms of `common.TransactionLog.CreatedAt`.
func NewTxGroupCache(timewindow time.Duration, render CacheKeyRenderer) *TxGroupCache {

	if render == nil {

		render = func(tx common.TransactionEntry, side common.SideType) string {

			return fmt.Sprintf(
				"%s_%s_%s", tx.GetExchange(), tx.GetAssetPair(), side,
			)

		}

	}

	return &TxGroupCache{
		cache:      map[string]*TxGroupCacheItem{},
		timewindow: timewindow,
		render:     render,
	}

}

// GetCache gets the cache item associated with the transaction type.
//
// It uses submitted `CacheKeyRenderer` to key the cache item.
func (txg *TxGroupCache) GetCache(
	tx common.TransactionLog,
	override CacheKeyRenderer,
) (item *TxGroupCacheItem, found bool) {

	var key string

	if override != nil {
		key = override(&tx, tx.GetSide())
	} else {
		key = txg.render(&tx, tx.GetSide())
	}

	item, found = txg.cache[key]
	return

}

// FlushCache will extract the `common.TxGroupEntry` and remove the cache completely.
func (txg *TxGroupCache) FlushCache(item *TxGroupCacheItem) common.TxGroupEntry {

	tx := item.tx
	key := txg.render(&tx, tx.GetSide())

	delete(txg.cache, key)
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

// GetAssetPairWhenNonFIAT returns all cache items that matches either `common.AssetPair.Asset`
// or `common.AssetPair.CostUnit`.
//
// It checks if _Asset_ or _CostUnit_ is Crypto before checking each end of the cached item.
func (txg *TxGroupCache) GetAssetPairWhenNonFIAT(
	assetPair common.AssetPair,
) (items []*TxGroupCacheItem, found bool) {

	items = []*TxGroupCacheItem{}

	for _, itm := range txg.cache {

		ap := itm.tx.GetAssetPair()

		if assetPair.Asset.IsCrypto() && (ap.Asset == assetPair.Asset ||
			ap.CostUnit == assetPair.Asset) {

			items = append(items, itm)
			continue

		}

		if assetPair.CostUnit.IsCrypto() && (ap.Asset == assetPair.CostUnit ||
			ap.CostUnit == assetPair.CostUnit) {

			items = append(items, itm)

		}

	}

	return items, len(items) > 0

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

		key := txg.render(&tx, s)
		if item, found := txg.cache[key]; found {
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
	switch tx.GetSide() {
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

	key := txg.render(&tx, side)
	item, found = txg.cache[key]
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

	key := txg.render(&tx, tx.Side)

	if _, ok := txg.cache[key]; ok {

		panic(
			fmt.Sprintf("cache item already created for: %s", key),
		)

	}

	item := &TxGroupCacheItem{
		next: tx.CreatedAt.Add(txg.timewindow),
	}

	txg.groupId++

	item.tx.ID = fmt.Sprint(txg.groupId)
	item.tx.AddTransactionEntry(tx)

	txg.cache[key] = item

	return item

}

// AddTransactionToCache will add the _tx_ to an existing cache item.
//
// If the cache item do not exist, it will panic.
func (txg *TxGroupCache) AddTransactionToCache(tx common.TransactionLog) {

	key := txg.render(&tx, tx.Side)

	if item, ok := txg.cache[key]; ok {

		item.tx.AddTransactionEntry(tx)
		return
	}

	panic(
		fmt.Sprintf(
			"no cache item created for: %s (cannot add transaction log)", key,
		),
	)
}
