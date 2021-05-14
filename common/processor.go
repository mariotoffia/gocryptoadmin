package common

type Processor interface {
	// Reset will clear any state in the processor, i.e. it enter its initialization state.
	Reset()
}

type TxAccountProcessor interface {
	Processor
	Process(tx TransactionEntry)
	Flush() []TransactionEntry
}

type TxEntryProcessor interface {
	Processor
	// Processes a transaction.
	//
	// It may hold the transaction in its cache if e.g. a merge operation is currently done. E.g.
	// using a group window.
	//
	// When `Transaction` instances is in this cache, they are said to be open. If `Flush` is invoked,
	// they are unconditionally merged and written to the underlying store.
	//
	// A `Processor` may merge `Transaction` instances as long as the following criteria is fulfilled
	//
	// 1. Within Group Window (if any)
	// 2. "Open `AssetPair` Transaction" - i.e. it is in a cache and not yet written to underlying _"store"_
	// 3. The new transaction, with same `AssetPair` do have same `SideType`
	// 4. The Asset part of the Open `Transaction` is not part of a `CostUnit` in the new `Transaction`
	//
	// If any of the above bullets fail, all _"Open"_ `Transaction` instances should be merged.
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

type TxGroupProcessor interface {
	Processor
	// Processes a transaction.
	//
	// It may hold the transaction in its cache if e.g. a merge operation is currently done. E.g.
	// using a group window.
	//
	// When `Transaction` instances is in this cache, they are said to be open. If `Flush` is invoked,
	// they are unconditionally merged and written to the underlying store.
	//
	// A `Processor` may merge `Transaction` instances as long as the following criteria is fulfilled
	//
	// 1. Within Group Window (if any)
	// 2. "Open `AssetPair` Transaction" - i.e. it is in a cache and not yet written to underlying _"store"_
	// 3. The new transaction, with same `AssetPair` do have same `SideType`
	// 4. The Asset part of the Open `Transaction` is not part of a `CostUnit` in the new `Transaction`
	//
	// If any of the above bullets fail, all _"Open"_ `Transaction` instances should be merged.
	Process(tx TxGroupEntry)

	// ProcessMany is *exactly* the same as `Process` but it will accept an array of entries instead.
	//
	// It is quite possible that the implementation just iterates these and calls `Process` underneath.
	ProcessMany(tx []TxGroupEntry)

	// Flush will make sure so any non committed to account `Transaction` instances is processed
	//
	// It will return all current flushed entries (including all earlier).
	Flush() []TxGroupEntry
}
