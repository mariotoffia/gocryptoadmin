package transactions

import "github.com/mariotoffia/gocryptoadmin/transactions/txcommon"

type TxLogGrouperImpl struct {
}

func NewTxLogGrouperImpl() *TxLogGrouperImpl {
	return &TxLogGrouperImpl{}
}

func (tg *TxLogGrouperImpl) GroupByDefault(tx []txcommon.Transaction) []txcommon.TransactionGroup {

	txg := []txcommon.TransactionGroup{}

	return txg
}
