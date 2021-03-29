package transactions

import (
	"fmt"
	"time"

	"github.com/ahmetb/go-linq/v3"
	"github.com/mariotoffia/gocryptoadmin/transactions/txcommon"
)

type TxLogGrouperImpl struct {
	secwindow time.Duration
}

func NewTxLogGrouperImpl() *TxLogGrouperImpl {
	return &TxLogGrouperImpl{
		secwindow: time.Duration(5 * 60),
	}
}

func (tg *TxLogGrouperImpl) GroupByDefault(tx []txcommon.Transaction) []txcommon.TransactionGroup {

	var ordered []txcommon.Transaction

	linq.From(tx).
		OrderBy(func(tx interface{}) interface{} {
			return tx.(txcommon.Transaction).Exchange
		}).
		ThenBy(func(tx interface{}) interface{} {
			return tx.(txcommon.Transaction).CreatedAt.Unix()
		}).
		ThenBy(func(tx interface{}) interface{} {
			return tx.(txcommon.Transaction).Product.Product
		}).
		ThenBy(func(tx interface{}) interface{} {
			return string(tx.(txcommon.Transaction).Side)
		}).
		ToSlice(&ordered)

	last := ordered[0].CreatedAt.Add(time.Second * tg.secwindow)
	side := ordered[0].Side
	product := ordered[0].Product.Product
	exchange := ordered[0].Exchange
	price := ordered[0].Cost.Price

	txg := []txcommon.TransactionGroup{}

	grp := txcommon.TransactionGroup{
		Transaction: txcommon.Transaction{
			Portfolio: ordered[0].Portfolio,
			ID:        ordered[0].CreatedAt.String(),
			Exchange:  ordered[0].Exchange,
			Side:      ordered[0].Side,
			CreatedAt: ordered[0].CreatedAt,

			Product: txcommon.Product{
				Product: ordered[0].Product.Product,
				Size:    0,
				Unit:    ordered[0].Product.Unit,
			},
			Cost: txcommon.Cost{
				Price:    ordered[0].Cost.Price,
				Fee:      0,
				Total:    0,
				CostUnit: ordered[0].Cost.CostUnit,
			},
		},
		Tx: []txcommon.Transaction{},
	}

	for _, txr := range ordered {

		if txr.CreatedAt.Before(last) &&
			side == txr.Side &&
			product == txr.Product.Product &&
			exchange == txr.Exchange &&
			price == txr.Cost.Price {

			grp.Cost.Total += txr.Cost.Total
			grp.Cost.Fee += txr.Cost.Fee
			grp.Product.Size += txr.Product.Size

			grp.Tx = append(grp.Tx, txr)

		} else {

			txg = append(txg, grp)

			last = txr.CreatedAt.Add(time.Second * tg.secwindow)
			side = txr.Side
			product = txr.Product.Product
			exchange = txr.Exchange
			price = txr.Cost.Price

			grp = txcommon.TransactionGroup{
				Transaction: txcommon.Transaction{
					Portfolio: txr.Portfolio,
					ID:        txr.CreatedAt.String(),
					Exchange:  txr.Exchange,
					Side:      txr.Side,
					CreatedAt: txr.CreatedAt,

					Product: txcommon.Product{
						Product: txr.Product.Product,
						Size:    txr.Product.Size,
						Unit:    txr.Product.Unit,
					},
					Cost: txcommon.Cost{
						Price:    txr.Cost.Price,
						Fee:      txr.Cost.Fee,
						Total:    txr.Cost.Total,
						CostUnit: txr.Cost.CostUnit,
					},
				},
				Tx: []txcommon.Transaction{},
			}

			grp.Tx = append(grp.Tx, txr)

		}

	}

	txg = append(txg, grp)

	fmt.Printf("count: %d\n", len(txg))

	for _, grp := range txg {

		fmt.Printf(
			"[%s] %s %s %f %s Price: %f Total: %f\n",
			grp.Exchange, grp.ID, grp.Side, grp.Product.Size, grp.Product.Product,
			grp.Cost.Price, grp.Cost.Total,
		)
		/*
			for _, txr := range grp.Tx {
				fmt.Printf(
					"\t%s %f %s Price: %f Total: %f\n",
					txr.CreatedAt.String(), txr.Product.Size, txr.Product.Product,
					txr.Cost.Price, txr.Total,
				)
			}*/
	}

	return txg
}
