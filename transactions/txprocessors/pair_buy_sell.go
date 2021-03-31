package txprocessors

import (
	"fmt"
	"time"

	"github.com/ahmetb/go-linq/v3"
	"github.com/mariotoffia/gocryptoadmin/transactions/txcommon"
	"github.com/mariotoffia/gocryptoadmin/utils"
)

type PairedTransaction struct {
	Sell txcommon.Transaction
	Buy  []txcommon.Transaction

	Exchange    string    `csv:"exchange" json:"exchange"`
	Asset       string    `csv:"product" json:"product"`
	Unit        string    `csv:"sizeunit" json:"sizeunit"`
	Size        float64   `csv:"size" json:"size"`
	SoldAt      time.Time `csv:"sold" json:"sold"`
	SoldPrice   float64   `csv:"sold price" json:"sold price"`
	SoldFee     float64   `csv:"sold fee" json:"sold fee"`
	SoldTotal   float64   `csv:"sold total" json:"sold total"`
	BoughtAt    time.Time `csv:"bought" json:"bought"`
	BoughtPrice float64   `csv:"bought price" json:"bought price"`
	BoughtFee   float64   `csv:"bought fee" json:"bought fee"`
	BoughtTotal float64   `csv:"bought total" json:"bought total"`
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
						tx.(txcommon.Transaction).Asset == asset
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
					pushme, pairme := splitBuy(tx, *buy)

					buyqueue.PushFront(&pushme)
					paired = append(paired, createPairedTx(tx, []txcommon.Transaction{pairme}))

					continue
				}

				if buy.Size < tx.Size {
					return
					// TODO: get all needed buys
					/*for buyqueue.IsEmpty() {
					}*/
				}

			} // for _, tx := range group

		}

	}

	return
}

func splitBuy(
	sell txcommon.Transaction,
	buy txcommon.Transaction) (pushme txcommon.Transaction, pairme txcommon.Transaction) {

	pushme = buy
	pairme = buy
	factor := utils.ToFixed(sell.Size/buy.Size, 8)

	fmt.Printf("must match %f = %f", utils.ToFixed(factor*buy.Size, 8), utils.ToFixed(buy.Size-sell.Size, 8))

	pushme.Size = utils.ToFixed(buy.Size-sell.Size, 8)
	pushme.Fee = utils.ToFixed(pushme.Fee*(1-factor), 8)
	pushme.Price = utils.ToFixed(pushme.Price*(1-factor), 8)
	pushme.Total = utils.ToFixed(pushme.Total*(1-factor), 8)
	pushme.GrpFee = pushme.Fee
	pushme.GrpSize = pushme.Size
	pushme.GrpTotal = pushme.Total

	pairme.Size = sell.Size
	pairme.Fee = utils.ToFixed(pairme.Fee*factor, 8)
	pairme.Price *= utils.ToFixed(pairme.Price*factor, 8)
	pairme.Total *= utils.ToFixed(pairme.Total*factor, 8)
	pairme.GrpFee = pairme.Fee
	pairme.GrpSize = pairme.Size
	pairme.GrpTotal = pairme.Total

	return
}

func createPairedTx(sell txcommon.Transaction, buy []txcommon.Transaction) PairedTransaction {

	pt := PairedTransaction{
		Exchange:  sell.Exchange,
		Asset:     sell.Asset,
		Sell:      sell,
		Unit:      sell.Unit,
		Buy:       buy,
		SoldAt:    sell.CreatedAt,
		SoldPrice: sell.Price,
		Size:      sell.Size,
		SoldTotal: sell.Total,
		SoldFee:   sell.Fee,
	}

	pt.BoughtAt = buy[0].CreatedAt

	for _, b := range buy {

		pt.BoughtFee = utils.ToFixed(pt.BoughtFee+b.Fee, 8)
		pt.BoughtPrice = utils.ToFixed(pt.BoughtPrice+b.Price, 8)
		pt.BoughtTotal = utils.ToFixed(pt.BoughtTotal+b.Total, 8)

	}

	return pt
}
