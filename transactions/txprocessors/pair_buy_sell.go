package txprocessors

import (
	"github.com/ahmetb/go-linq/v3"
	"github.com/mariotoffia/gocryptoadmin/transactions/txcommon"
)

func PairBuySell(logs []txcommon.Transaction) []txcommon.Transaction {

	paired := []txcommon.Transaction{}

	exchanges := txcommon.Exchanges(logs)
	assets := txcommon.Assets(logs)

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

			for _, tx := range group {

				if tx.Side == txcommon.SideTypeBuy {

					buyqueue.Push(&tx)
					continue
				}

			}

		}

	}

	return paired
}
