package transactions

import (
	"fmt"
	"time"

	"github.com/ahmetb/go-linq/v3"
	"github.com/mariotoffia/gocryptoadmin/transactions/txcommon"
)

type TxLogGrouperImpl struct {
	secwindow time.Duration
	tx        []txcommon.Transaction
	txg       []txcommon.TransactionGroup
}

func NewTxLogGrouperImpl(tx []txcommon.Transaction) *TxLogGrouperImpl {

	return (&TxLogGrouperImpl{
		secwindow: time.Duration(5 * 60),
		tx:        tx,
		txg:       []txcommon.TransactionGroup{},
	}).groupByExchangeCreatedProductSide()

}

func (tg *TxLogGrouperImpl) TransactionGroups() []txcommon.TransactionGroup {
	return tg.txg
}

func (tg *TxLogGrouperImpl) Transactions() []txcommon.Transaction {
	return tg.tx
}

func (tg *TxLogGrouperImpl) GroupViaTimeWindow() *TxLogGrouperImpl {

	tx := tg.tx

	last := tx[0].CreatedAt.Add(time.Second * tg.secwindow)
	side := tx[0].Side
	product := tx[0].Product.Product
	exchange := tx[0].Exchange
	price := tx[0].Cost.Price

	txg := []txcommon.TransactionGroup{}

	grp := txcommon.TransactionGroup{
		Transaction: txcommon.Transaction{
			Portfolio: tx[0].Portfolio,
			ID:        tx[0].CreatedAt.String(),
			Exchange:  tx[0].Exchange,
			Side:      tx[0].Side,
			CreatedAt: tx[0].CreatedAt,

			Product: txcommon.Product{
				Product: tx[0].Product.Product,
				Size:    0,
				Unit:    tx[0].Product.Unit,
			},
			Cost: txcommon.Cost{
				Price:    tx[0].Cost.Price,
				Fee:      0,
				Total:    0,
				CostUnit: tx[0].Cost.CostUnit,
			},
		},
		Tx: []txcommon.Transaction{},
	}

	for _, txr := range tx {

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

	tg.txg = append(txg, grp)
	return tg
}

func (tg *TxLogGrouperImpl) DumpGroup(details bool) *TxLogGrouperImpl {

	for _, grp := range tg.txg {

		fmt.Printf(
			"[%s] %s %s %f %s Price: %f Total: %f\n",
			grp.Exchange, grp.ID, grp.Side, grp.Product.Size, grp.Product.Product,
			grp.Cost.Price, grp.Cost.Total,
		)

		if details {

			for _, txr := range grp.Tx {
				fmt.Printf(
					"\t%s %f %s Price: %f Total: %f\n",
					txr.CreatedAt.String(), txr.Product.Size, txr.Product.Product,
					txr.Cost.Price, txr.Total,
				)
			}

		}
	}

	return tg
}

// groupByExchangeCreatedProductSide prepares the transactions to be grouped
func (tg *TxLogGrouperImpl) groupByExchangeCreatedProductSide() *TxLogGrouperImpl {

	var ordered []txcommon.Transaction

	linq.From(tg.tx).
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

	tg.tx = ordered
	return tg

}
