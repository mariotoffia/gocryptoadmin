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
	tx       common.TxGroup
}

type txgroupTxCache struct {
	cache     map[string]*txgroupTxCacheItem
	groupId   int64
	secwindow time.Duration
}

func (txg *txgroupTxCache) GetItem(tx common.Transaction) *txgroupTxCacheItem {

	if item, ok := txg.cache[tx.Exchange+tx.AssetPair.String()+string(tx.Side)]; ok {

		if 

		return item
	}

	item := &txgroupTxCacheItem{
		exchange: tx.Exchange,
		pair:     tx.AssetPair,
		side:     tx.Side,
		next:     tx.CreatedAt.Add(time.Second * txg.secwindow),
		tx: common.TxGroup{
			Transaction: tx,
			Tx:          []common.Transaction{},
		},
	}

	txg.groupId++

	item.tx.Fee = 0
	item.tx.ID = fmt.Sprint(txg.groupId)
	item.tx.PricePerUnit = 0
	item.tx.TotalPrice = 0
	item.tx.AssetSize = 0

	txg.cache[tx.Exchange+tx.AssetPair.String()+string(tx.Side)] = item

	return item

}
