package common

type Processor interface {
	// Reset will clear any state in the processor, i.e. it enter its initialization state.
	Reset()
}

type Flushable interface {
	Flush() []TransactionEntry
}

type TxEntryProcessor interface {
	Processor
	Process(tx TransactionEntry)
	ProcessMany(tx []TransactionEntry)
}

type TxFlushableEntryProcessor interface {
	TxEntryProcessor
	Flushable
}

type TxLogProcessor interface {
	Processor
	// Processes a transaction.
	Process(tx TransactionLog)

	// ProcessMany is *exactly* the same as `Process` but it will accept an array of entries instead.
	//
	// It is quite possible that the implementation just iterates these and calls `Process` underneath.
	ProcessMany(tx []TransactionLog)

	// Flush will make sure so any non committed to account `Transaction` instances is processed
	//
	// It will return all current flushed entries (including all earlier).
	Flush() []TransactionLog
}

// TxGroupProcessor is same as `TxEntryProcessor` except that it handles `TxGroupEntry` instances
// instead.
type TxGroupProcessor interface {
	Processor
	Process(tx TxGroupEntry)
	ProcessMany(tx []TxGroupEntry)
	Flush() []TxGroupEntry
}

// MultiAccountTxProcessor will keep account records for each exchange
// and a set of global `ExchangeAll` accounts.
type MultiAccountTxProcessor interface {
	Processor
	Process(tx TransactionEntry)
	ProcessMany(tx []TransactionEntry)
	Flush() map[string][]TransactionEntry
}
