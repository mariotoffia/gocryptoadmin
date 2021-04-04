package common

type Processor interface {
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
	// 2. "Open `AssetPair` Transaction" - i.e. it is in a cache and not yet registered in the ledger
	// 3. The new transaction, with same `AssetPair` do have same `SideType`
	// 4. The Asset part of the Open `Transaction` is not part of a `CostUnit` in the new `Transaction`
	//
	// If any of the above bullets fail, all _"Open"_ `Transaction` instances should be merged.
	Process(tx TransactionEntry)
	// Flush will make sure so any non committed to account `Transaction` instances is processed
	Flush()
	// UseGroupWindow specifies that the window to group transactions is up to _s_ seconds.
	//
	// Default is 6 hours.
	UseGroupWindow(s int64)
}
