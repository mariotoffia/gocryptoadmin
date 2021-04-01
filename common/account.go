package common

type Account interface {
	// Asset is the asset of which this account is tracking
	Asset() AssetType
	// CurrentAmount is the current amount of asset "right now"
	CurrentAmount() float64
	//Withdraw states that a transaction is widthdrawing from this account
	//
	// For example SELL LTC - EUR will call withdraw on LTC Account and Deposit on
	// EUR account. Or if
	Withdraw(tx Transaction)
	// Deposit is the reverse of `Withdraw`
	Deposit(tx Transaction)
	// Transactions get all transactions for this `Account`
	Transactions() []AccountTransaction
}

// Account keeps track on one asset type
type AccountImpl struct {
	Asset         AssetType
	CurrentAmount float64
}
