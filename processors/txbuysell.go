package processors

import (
	"fmt"

	"github.com/mariotoffia/gocryptoadmin/common"
	"github.com/mariotoffia/gocryptoadmin/utils"
)

// TxBuySellProcessor will pair `common.SideTypeSell` with
// earlier `common.SideTypeBuy` transactions.
type TxBuySellProcessor struct {
	queue    *common.TxAssetFIFOQueues
	entries  []common.TxBuySellEntry
	log      bool
	taxation bool
}

func NewTxBuySellProcessor() *TxBuySellProcessor {

	return &TxBuySellProcessor{
		queue:   common.NewTxAssetFIFOQueues(),
		entries: []common.TxBuySellEntry{},
	}

}

func (bs *TxBuySellProcessor) Reset() {
	bs.entries = []common.TxBuySellEntry{}
	bs.queue.Reset()
}

// UseTaxationMarking enables the taxation marking, reducing the
// double taxation that may occur when e.g. selling LTC to BTC and
// then BTC to EUR. When this is enabled it will mark the LTC-BTC
// SELL tx as taxed. Hence, the BTC-EUR will not be taxed (other)
// than the difference between LTC-BTC, BTC Value, and the value of
// BTC when BTC-EUR occurrence.
func (bs *TxBuySellProcessor) UseTaxationMarking() {
	bs.taxation = true
}

// UseLog enables Enqueue, Dequeue, ReEnqueue logging
func (bs *TxBuySellProcessor) UseLog() {
	bs.log = true
}

func (bs *TxBuySellProcessor) ProcessMany(tx []common.TransactionEntry) {

	for i := range tx {

		bs.Process(tx[i])

	}
}

func (bs *TxBuySellProcessor) Process(tx common.TransactionEntry) {

	tx = tx.Clone()
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

		// Got crypto as payment, we may sell it again!
		bs.queue.Enq(assetPair.CostUnit, tx)

		if bs.log {
			logSingle("Push", assetPair.CostUnit, tx, true /*price*/, false)
		}
	}

	entries, res, size := bs.drainBuys(assetPair.Asset, tx.GetAssetSize())

	// Split last entry and PutBack overflow into queue again.
	if res == common.DequeueUntilResultOverflow {

		// putback and keep is reversed in overflow
		putback, keep := splitEntryByOverflow(entries[len(entries)-1], -size)
		bs.queue.Enq(assetPair.Asset, putback)

		entries = append(entries[:len(entries)-1], keep)

		if bs.log {
			log("Pop", assetPair.Asset, entries, false)
			fmt.Print("[ queue <- ")
			logSingle("PushBack", assetPair.Asset, putback, false /*size*/, false)
			fmt.Println("]")
		}

	} else {

		if bs.log {
			log("Pop", assetPair.Asset, entries, true)
		}

	}

	// Entries are the BUY transactions that matches this single sell!
	// Create TxPair and assign buy and sell side -> bs.entries
	bs.entries = append(bs.entries, common.NewTxBuySellLog(tx, entries))
}

// ProcessBuy will process a _tx_ that reflects a BUY transaction.
func (bs *TxBuySellProcessor) ProcessBuy(tx common.TransactionEntry) {

	assetPair := tx.GetAssetPair()

	// Enqueue the BUY order to later match a SELL.
	bs.queue.Enq(assetPair.Asset, tx)

	if bs.log {
		logSingle("Push", assetPair.Asset, tx, false /*size*/, assetPair.CostUnit.IsFIAT())
	}

	if assetPair.CostUnit.IsFIAT() {
		return
	}

	// Need to remove BUY transaction(s) for CostUnit
	// by getting the total price, since crypto this will match
	// up to BUY tx GetAssetSize().
	//
	// It is negated since the buy in crypto will log entry as with fiat -> negative value.
	entries, res, size := bs.drainBuys(assetPair.CostUnit, -tx.GetTotalPrice())

	if bs.log {
		log(
			"Pop",
			assetPair.CostUnit,
			entries,
			res == common.DequeueUntilResultDone,
		)
	}

	if res == common.DequeueUntilResultDone {
		return // All is removed
	}

	// Extract overflow and put it back to FIFO queue
	_, putback := splitEntryByOverflow(entries[len(entries)-1], -size)

	bs.queue.Enq(assetPair.CostUnit, putback)

	if bs.log {
		logSingle("PushBack", assetPair.CostUnit, putback, false /*size*/, true)
	}

}

func logSingle(dir string, asset common.AssetType, entry common.TransactionEntry, price, cr bool) {

	fmt.Printf("%s(", dir)

	f := entry.GetAssetSize()

	if price {
		f = entry.GetTotalPrice()
	}

	fmt.Printf("%.8f %s)  ", f, asset)

	if cr {
		fmt.Println()
	}
}

func log(dir string, asset common.AssetType, entries []common.TransactionEntry, cr bool) {

	fmt.Printf("%s(", dir)

	f := float64(0)
	for _, entry := range entries {

		side := entry.GetSide()
		if side == common.SideTypeSell {
			f = utils.ToFixed(f+entry.GetTotalPrice(), 8)
		} else if side == common.SideTypeBuy {
			f = utils.ToFixed(f+entry.GetAssetSize(), 8)
		} else {
			panic("expecting BUY or SELL while logging")
		}

	}

	fmt.Printf("%.8f %s)  ", f, asset)

	if cr {
		fmt.Println()
	}
}

func (bs *TxBuySellProcessor) Flush() (entries []common.TxBuySellEntry, noPairing []common.TransactionEntry) {

	// Get the overflow
	noPairing = bs.queue.DequeueAll()
	// Copy the buy-sell entries
	entries = bs.entries

	bs.Reset()
	return
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

	fullSize := utils.ToFixed(size, 8)

	entries, res := bs.queue.DequeueUntil(
		asset,
		func(tx common.TransactionEntry) common.DequeueUntilResult {

			if tx.GetSide() == common.SideTypeBuy {
				size = utils.ToFixed((size - tx.GetAssetSize()), 8)
			} else if tx.GetSide() == common.SideTypeSell {
				size = utils.ToFixed((size - tx.GetTotalPrice()), 8)
			} else {
				panic("expecting BUY or SELL only")
			}

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
