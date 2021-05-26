package processors

import (
	"fmt"

	"github.com/mariotoffia/gocryptoadmin/common"
)

// TxBuySellProcessor will pair `common.SideTypeSell` with
// earlier `common.SideTypeBuy` transactions.
type TxBuySellProcessor struct {
	queue   *common.TxAssetFIFOQueues
	entries []common.TxPair
}

func NewTxBuySellProcessor() *TxBuySellProcessor {

	return &TxBuySellProcessor{
		queue:   common.NewTxAssetFIFOQueues(),
		entries: []common.TxPair{},
	}

}

func (bs *TxBuySellProcessor) Reset() {
	bs.entries = []common.TxPair{}
}

func (bs *TxBuySellProcessor) ProcessMany(tx []common.TransactionEntry) {

	for i := range tx {

		bs.Process(tx[i])

	}

}

func (bs *TxBuySellProcessor) Process(tx common.TransactionEntry) {

	assetPair := tx.GetAssetPair()
	side := tx.GetSide()

	if side == common.SideTypeBuy {

		bs.queue.Enq(assetPair.Asset, tx)

		if !assetPair.CostUnit.IsFIAT() {

			deq := bs.queue.Deq(assetPair.CostUnit)
			// TODO: keep deq until tx.GetTotalPrice() is satisfied
			// TODO: If any excess, split and PutBack overflow into queue again.
			fmt.Println(deq)
		}

		return
	}

	if side != common.SideTypeSell {
		return
	}

	if !assetPair.CostUnit.IsFIAT() {
		bs.queue.Enq(assetPair.CostUnit, tx)
	}

	deq := bs.queue.Deq(assetPair.Asset)
	// TODO: keep deq until tx.GetSize() is satisfied
	// TODO: If any excess, split and PutBack overflow into queue again.
	fmt.Println(deq)

	// TODO: Create TxPair and assign buy and sell side -> bs.entries
}

func (bs *TxBuySellProcessor) Flush() []common.TxPair {

	p := bs.entries
	bs.Reset()
	return p

}
