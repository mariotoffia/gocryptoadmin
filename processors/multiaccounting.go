package processors

import "github.com/mariotoffia/gocryptoadmin/common"

// MultiExchangeAccountingProcessor is a multi account processor
// that processes all and each exchange individual and hence
// it is possible to get a complete view and a view on each exchange
// locally.
//
// It implements the `common.MultiAccountTxProcessor` interface.
type MultiExchangeAccountingProcessor struct {
	processors map[string]*AccountingProcessor
}

// NewMultiExchangeAccountingProcessor creates a new `MultiExchangeAccountingProcessor`.
//
// It will automatically add each exchange automatically if not yet exist.
func NewMultiExchangeAccountingProcessor() *MultiExchangeAccountingProcessor {

	return &MultiExchangeAccountingProcessor{
		processors: map[string]*AccountingProcessor{},
	}
}

func (m *MultiExchangeAccountingProcessor) Reset() {
	m.processors = map[string]*AccountingProcessor{}
}

func (m *MultiExchangeAccountingProcessor) ProcessMany(tx []common.TransactionEntry) {

	for i := range tx {

		m.Process(tx[i])

	}

}

func (m *MultiExchangeAccountingProcessor) Process(tx common.TransactionEntry) {

	all := m.processors[common.ExchangeAll]
	if all == nil {

		all = NewAccountingProcessor(common.ExchangeAll)
		m.processors[common.ExchangeAll] = all

	}

	local := m.processors[tx.GetExchange()]
	if local == nil {

		local = NewAccountingProcessor(tx.GetExchange())
		m.processors[tx.GetExchange()] = local

	}

	all.Process(tx)
	local.Process(tx)

}

func (m *MultiExchangeAccountingProcessor) Flush() map[string][]common.TransactionEntry {

	r := map[string][]common.TransactionEntry{}

	for k, v := range m.processors {

		r[k] = v.Flush()

	}

	return r

}
