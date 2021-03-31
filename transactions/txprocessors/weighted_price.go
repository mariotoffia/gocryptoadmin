package txprocessors

import (
	"github.com/ahmetb/go-linq/v3"
	"github.com/mariotoffia/gocryptoadmin/transactions/txcommon"
	"github.com/mariotoffia/gocryptoadmin/utils"
)

func WeightedPrice(logs []txcommon.Transaction) []txcommon.Transaction {

	weighted := []txcommon.Transaction{}

	max := linq.From(logs).Select(func(tx interface{}) interface{} {
		return tx.(txcommon.Transaction).GroupID
	}).Max().(int64)

	for i := int64(1); i < max; i++ {

		group := []txcommon.Transaction{}

		linq.From(logs).
			Where(func(tx interface{}) bool {
				return tx.(txcommon.Transaction).GroupID == i
			}).
			ToSlice(&group)

		if len(group) == 0 {
			continue
		}

		size := float64(0)
		totalPrice := float64(0)
		fee := float64(0)

		for _, tx := range group {

			size += tx.Size
			totalPrice += tx.Price * tx.Size
			fee += tx.Fee

		}

		tx := group[0]
		tx.Size = utils.ToFixed(size, 8)
		tx.Fee = utils.ToFixed(fee, 8)
		tx.Price = utils.ToFixed(totalPrice/size, 8)

		if tx.Side == txcommon.SideTypeBuy {
			tx.Total = utils.ToFixed(-totalPrice-fee, 8)
		} else {
			tx.Total = utils.ToFixed(totalPrice-fee, 8)
		}

		tx.GrpFee = tx.Fee
		tx.GrpSize = tx.Size
		tx.GrpTotal = tx.Total

		weighted = append(weighted, tx)

	}

	return weighted
}
