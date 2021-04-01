package txprocessors

import (
	"fmt"
	"sort"
	"time"

	"github.com/ahmetb/go-linq/v3"
	"github.com/mariotoffia/gocryptoadmin/transactions/txcommon"
	"github.com/mariotoffia/gocryptoadmin/utils"
)

type PairedTransaction struct {
	Sell txcommon.Transaction
	Buy  []txcommon.Transaction

	Exchange    string    `csv:"exchange" json:"exchange"`
	AssetPair   string    `csv:"asset pair" json:"assetpair"`
	Unit        string    `csv:"sizeunit" json:"sizeunit"`
	Size        float64   `csv:"size" json:"size"`
	SoldAt      time.Time `csv:"sold" json:"sold"`
	SoldPrice   float64   `csv:"sold price" json:"soldprice"`
	SoldFee     float64   `csv:"sold fee" json:"soldfee"`
	SoldTotal   float64   `csv:"sold total" json:"soldtotal"`
	BoughtAt    time.Time `csv:"bought" json:"bought"`
	BoughtPrice float64   `csv:"bought price" json:"boughtprice"`
	BoughtFee   float64   `csv:"bought fee" json:"boughtfee"`
	BoughtTotal float64   `csv:"bought total" json:"boughttotal"`
}

func PairBuySell(
	logs []txcommon.Transaction) (paired []PairedTransaction, unpaired []txcommon.Transaction) {

	paired = []PairedTransaction{}

	exchanges := txcommon.Exchanges(logs)
	assets := txcommon.Assets(logs)
	unpaired = []txcommon.Transaction{}

	for _, exchange := range exchanges {

		for _, asset := range assets {

			group := []txcommon.Transaction{}

			linq.From(logs).
				Where(func(tx interface{}) bool {
					return tx.(txcommon.Transaction).Exchange == exchange &&
						tx.(txcommon.Transaction).AssetPair == asset
				}).
				ToSlice(&group)

			buyqueue := txcommon.TxQueue{}

			for i, tx := range group {

				if tx.Side == txcommon.SideTypeBuy {
					buyqueue.Push(&group[i])
					continue
				}

				// Side: SELL --> Match any buyers

				if buyqueue.IsEmpty() {
					unpaired = append(unpaired, tx)
					continue
				}

				buy := buyqueue.Pop()

				if buy.Size == tx.Size {
					paired = append(paired, createPairedTx(tx, []txcommon.Transaction{*buy}))
					continue
				}

				if buy.Size > tx.Size {
					// Split buy into two
					pushme, pairme := splitTxBySize(tx.Size, *buy)

					buyqueue.PushFront(&pushme)
					paired = append(paired, createPairedTx(tx, []txcommon.Transaction{pairme}))
					continue
				}

				// buy.Size < tx.Size -> iterate all buys
				add, reminder := matchBuyWithSell(tx, buy, &buyqueue)

				if reminder != nil {
					unpaired = append(unpaired, *reminder)
				}

				paired = append(paired, add)

			} // for _, tx := range group

			for !buyqueue.IsEmpty() {

				unpaired = append(unpaired, *buyqueue.Pop())

			}
		}

	}

	sort.Slice(paired, func(i, j int) bool {
		return paired[i].BoughtAt.Before(paired[j].BoughtAt)
	})

	return
}

func matchBuyWithSell(
	sell txcommon.Transaction,
	buy *txcommon.Transaction,
	buyqueue *txcommon.TxQueue) ( /*paired*/ PairedTransaction /*remainder*/, *txcommon.Transaction) {

	buys := []txcommon.Transaction{*buy}
	size := buy.Size

	// get all needed buys
	for !buyqueue.IsEmpty() && size < sell.Size {

		buy = buyqueue.Pop()
		size = utils.ToFixed(size+buy.Size, 8)

		if size <= sell.Size {
			buys = append(buys, *buy)
		}

	}

	if size == sell.Size {
		return createPairedTx(sell, buys), nil
	}

	if size > sell.Size {
		// Split last buy into two
		pushme, pairme := splitTxBySize(sell.Size-(size-buy.Size), *buy)

		buyqueue.PushFront(&pushme)

		buys = append(buys, pairme)
		return createPairedTx(sell, buys), nil

	}

	// sell is still larger than buy -> split sell and add it to unpaired
	fmt.Println("WARNING: This is probably an error!!!??? TODO: need to return matches and reason!")
	remainder, split := splitTxBySize(size, sell)

	return createPairedTx(split, buys), &remainder

}

// splitTxBySize will split the in param _tx_ to two parts where the _split_ is the _splitSize_
// and the _remainder_ is the left over the _split_.
func splitTxBySize(
	splitSize float64,
	tx txcommon.Transaction) (remainder txcommon.Transaction, split txcommon.Transaction) {

	remainder = tx
	split = tx
	factor := utils.ToFixed(splitSize/tx.Size, 8)

	remainder.Size = utils.ToFixed(tx.Size-splitSize, 8)
	remainder.Fee = utils.ToFixed(remainder.Fee*(1-factor), 8)
	//remainder.Price = utils.ToFixed(remainder.Price*(1-factor), 8)
	remainder.Total = utils.ToFixed(remainder.Total*(1-factor), 8)
	remainder.GrpFee = remainder.Fee
	remainder.GrpSize = remainder.Size
	remainder.GrpTotal = remainder.Total

	split.Size = splitSize
	split.Fee = utils.ToFixed(split.Fee*factor, 8)
	//split.Price = utils.ToFixed(split.Price*factor, 8)
	split.Total = utils.ToFixed(split.Total*factor, 8)
	split.GrpFee = split.Fee
	split.GrpSize = split.Size
	split.GrpTotal = split.Total

	return
}

func createPairedTx(sell txcommon.Transaction, buy []txcommon.Transaction) PairedTransaction {

	pt := PairedTransaction{
		Exchange:  sell.Exchange,
		AssetPair: sell.AssetPair,
		Sell:      sell,
		Unit:      sell.Unit,
		Buy:       buy,
		SoldAt:    sell.CreatedAt,
		SoldPrice: sell.Price,
		Size:      sell.Size,
		SoldTotal: sell.Total,
		SoldFee:   sell.Fee,
	}

	price := float64(0)
	prominent := 0
	maxsize := float64(0)

	for i, b := range buy {

		sizeinsell := b.Size / sell.Size

		if sizeinsell > maxsize {
			maxsize = sizeinsell
			prominent = i
		}

		price = utils.ToFixed(price+(b.Price*sizeinsell), 8)

		pt.BoughtFee = utils.ToFixed(pt.BoughtFee+b.Fee, 8)
		pt.BoughtTotal = utils.ToFixed(pt.BoughtTotal+b.Total, 8)

	}

	pt.BoughtPrice = price
	pt.BoughtAt = buy[prominent].CreatedAt

	return pt
}
