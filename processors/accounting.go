package processors

import "github.com/mariotoffia/gocryptoadmin/common"

type AccountingProcessor struct {
	entries  []common.AccountEntry
	previous common.AccountEntry
}

func NewAccountingProcessor() *AccountingProcessor {

	return &AccountingProcessor{
		entries: []common.AccountEntry{},
	}

}

func (ap *AccountingProcessor) Reset() {
	ap.entries = []common.AccountEntry{}
}

func (ap *AccountingProcessor) Process(tx common.TransactionEntry) {

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
