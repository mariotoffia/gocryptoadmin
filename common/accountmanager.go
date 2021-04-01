package common

type AccountManager interface {
	// Processes a transaction.
	//
	// It may hold the transaction in its cache if e.g. a merge operation is currently done. E.g.
	// using a group window.
	Process(tx Transaction)
	// Flush will make sure so any non committed to account `Transaction` instances is processed
	Flush()
	// UseGroupWindow specifies that the window to group transactions is up to _s_ seconds.
	//
	// Default is 6 hours.
	UseGroupWindow(s int64)
}

// AccountManager is keeping track on all `Accounts`
type AccountManagerImpl struct {
	// Account is a map keyed with `AssetType` and value is `*Account`
	Accounts map[AssetType]*AccountImpl
}
