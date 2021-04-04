package processors

import (
	"fmt"
	"time"

	"github.com/mariotoffia/gocryptoadmin/common"
)

type txgroupTxCacheItem struct {
	exchange string
	pair     common.AssetPair
	side     common.SideType
	next     time.Time
	tx       common.TxGroupEntry
}

type txgroupTxCache struct {
	cache     map[string]*txgroupTxCacheItem
	groupId   int64
	secwindow time.Duration
}

func (txg *txgroupTxCache) GetItem(tx common.TransactionLog) *txgroupTxCacheItem {

	if item, ok := txg.cache[tx.Exchange+tx.AssetPair.String()+string(tx.Side)]; ok {

		item.tx.AddTransactionEntry(tx)
		return item
	}

	item := &txgroupTxCacheItem{
		exchange: tx.Exchange,
		pair:     tx.AssetPair,
		side:     tx.Side,
		next:     tx.CreatedAt.Add(time.Second * txg.secwindow),
		tx: common.TxGroupEntry{
			TransactionLog: tx,
			Tx:             []common.TransactionLog{},
		},
	}

	txg.groupId++

	item.tx.ID = fmt.Sprint(txg.groupId)
	item.tx.AddTransactionEntry(tx)

	txg.cache[tx.Exchange+tx.AssetPair.String()+string(tx.Side)] = item

	return item

}
