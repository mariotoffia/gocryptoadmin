package processors

import "github.com/mariotoffia/gocryptoadmin/common"

// AccountingProcessor implements (ish) the `TxGroupProcessor` interface.
type AccountingProcessor struct {
	entries  []common.AccountEntry
	previous common.AccountEntry
	exchange string
}

// NewAccountingProcessor creates a new accounting processor that will
// keep book on each account status for both _Asset_ and _CostUnit_ to
// handle transaction log balance.
//
// The _exchange_ may be any of the exchanges, if it is a empty string
// it will accept *any* exchange and is alias for all.
func NewAccountingProcessor(exchange string) *AccountingProcessor {

	if exchange == "" {
		exchange = common.ExchangeAll
	}

	return &AccountingProcessor{
		entries:  []common.AccountEntry{},
		exchange: exchange,
	}

}

func (ap *AccountingProcessor) Reset() {
	ap.entries = []common.AccountEntry{}
}

func (ap *AccountingProcessor) ProcessMany(tx []common.TransactionEntry) {

	for i := range tx {

		ap.Process(tx[i])

	}

}

func (ap *AccountingProcessor) Process(tx common.TransactionEntry) {

	if ap.exchange != common.ExchangeAll &&
		tx.GetExchange() != ap.exchange {

		return

	}

	tx = tx.Clone()

	acc := common.NextAccountLog(ap.previous, tx)

	ap.entries = append(ap.entries, acc)
	ap.previous = acc

}

func (ap *AccountingProcessor) Flush() []common.TransactionEntry {

	list := make([]common.TransactionEntry, len(ap.entries))
	prototype := map[common.AssetType]float64{}

	for i := range ap.entries {

		for k := range ap.entries[i].GetAccountStatus() {

			prototype[k] = 0

		}

	}

	for i := range ap.entries {

		list[i] = ap.entries[i].EnsureAccounts(prototype)

	}

	ap.Reset()
	return list
}
