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

	side := tx.GetSide()

	if side == common.SideTypeBuy {
		bs.ProcessBuy(tx)
		return
	}

	// Only process SELL
	if side != common.SideTypeSell {
		return
	}

	assetPair := tx.GetAssetPair()

	if !assetPair.CostUnit.IsFIAT() {
		// Got crypto as payment
		// TODO: Mark as "NO_TAX" since this entry is now already taxed
		// TODO: due to the _tx_ is sold and tax is paid for that.
		bs.queue.Enq(assetPair.CostUnit, tx)
	}

	entries, res, size := bs.drainBuys(assetPair.CostUnit, tx.GetAssetSize())

	// Split last entry and PutBack overflow into queue again.
	if res == common.DequeueUntilResultOverflow {

		keep, putback := splitEntryByOverflow(entries[len(entries)-1], -size)
		bs.queue.Enq(assetPair.CostUnit, putback)

		entries = append(entries[:len(entries)-1], keep)

	}

	// TODO: entries are the BUY transactions that matches this single sell!
	// TODO: Create TxPair and assign buy and sell side -> bs.entries

	fmt.Println(tx)      // SELL
	fmt.Println(entries) // BUY
}

// ProcessBuy will process a _tx_ that reflects a BUY transaction.
func (bs *TxBuySellProcessor) ProcessBuy(tx common.TransactionEntry) {

	assetPair := tx.GetAssetPair()

	// Enqueue the BUY order to later match a SELL.
	bs.queue.Enq(assetPair.Asset, tx)

	if assetPair.CostUnit.IsFIAT() {
		return
	}

	// Need to remove BUY transaction(s) for CostUnit
	// by getting the total price, since crypto this will match
	// up to BUY tx GetAssetSize().
	entries, res, size := bs.drainBuys(assetPair.CostUnit, tx.GetTotalPrice())

	if res == common.DequeueUntilResultDone {
		return // All is removed
	}

	// Extract overflow and put it back to FIFO queue
	_, putback := splitEntryByOverflow(entries[len(entries)-1], -size)
	bs.queue.Enq(assetPair.CostUnit, putback)
}

func (bs *TxBuySellProcessor) Flush() []common.TxPair {

	p := bs.entries
	bs.Reset()
	return p

}

// splitEntryByOverflow will split the _tx_ into the one to "keep" and the one
// overflow. The overflow is specified in _size_ and not total price. Hence,
// this is meant to split crypto BUY transactions that did not, exactly, match up a SELL.
func splitEntryByOverflow(
	tx common.TransactionEntry,
	overflow float64,
) (keep common.TransactionEntry, putback common.TransactionEntry) {

	return tx.SplitSize(overflow)

}

// drainBuys will remove BUYs from the queue until satisfied _size_.
//
// If `common.DequeueUntilResultUnderflow`, it will *panic*.
func (bs *TxBuySellProcessor) drainBuys(
	asset common.AssetType,
	size float64,
) ([]common.TransactionEntry, common.DequeueUntilResult, float64) {

	fullSize := size

	entries, res := bs.queue.DequeueUntil(
		asset,
		func(tx common.TransactionEntry) common.DequeueUntilResult {

			size -= tx.GetAssetSize()
			return bs.dequeueResultFromSize(size)

		},
	)

	if res == common.DequeueUntilResultUnderflow {

		panic(
			fmt.Sprintf(
				"Could not find all BUY entries for asset: %s size: %f, missing: %f",
				asset, fullSize, size,
			),
		)

	}

	return entries, res, size
}

// dequeueResultFromSize returns a proper `common.DequeueUntilResult` base on _size_.
func (bs *TxBuySellProcessor) dequeueResultFromSize(size float64) common.DequeueUntilResult {

	if size == 0 {
		return common.DequeueUntilResultDone
	}

	if size < 0 {
		return common.DequeueUntilResultOverflow
	}

	return common.DequeueUntilResultContinue

}
