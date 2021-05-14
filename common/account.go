package common

import "time"

// AccountEntry is a single line in the accounting.
//
// Each entry may contain one or more `AssetType` with
// it's current value. E.g. _EUR 100, LTC 150.02, BTC 1.00_.
//
// The `TransactionEntry` that was the reason for this account
// entry (a change in the account) is always included. Hence
// a set of `TransactionEntry` will have a corresponding `AccountEntry`.
type AccountEntry interface {
	// Derives from TransactionEntry since it is attached to it.
	TransactionEntry
	// GetAccountStatus returns the current value in the account.
	GetAccountStatus() map[AssetType]float64
	// EnsureAccountsjust makes sure that all `AssetType`(s) do exist
	// in account status based on the _prototype_.
	EnsureAccounts(prototype map[AssetType]float64) AccountEntry
	// HasTransaction checks if it has a valid transaction.
	//
	// This function may return `false` when the `AccountEntry` is
	// considered as a prototype entry where `GetAccountStatus` contains
	// e.g. all `AssetType`(s) and hence used for root construction.
	HasTransaction() bool
}

type AccountLog struct {
	tx     TransactionEntry
	status map[AssetType]float64
}

// NewAccountLog creates a new `AccountLog` with `AccountLog.Status` initialied.
func NewAccountLog(tx TransactionEntry) *AccountLog {

	return &AccountLog{
		tx:     tx,
		status: map[AssetType]float64{},
	}

}

// NextAccountLog creates a new `AccountLog` based on _tx_ but copies all
// `Status` values into the new one and substract or add to `Status` based
// on the _tx_.
//
// If _previous_ is `nil`, it will invoke the `NewAccountLog` instead. If the
// `HasTransaction` returns `false` on _previous_, it will only be used to copy
// the `AssetPair`(s).
func NextAccountLog(previous AccountEntry, tx TransactionEntry) *AccountLog {

	ac := NewAccountLog(tx)

	if previous == nil {
		return ac
	}

	for k, v := range previous.GetAccountStatus() {
		ac.status[k] = v
	}

	if !previous.HasTransaction() {
		return ac
	}

	asset := tx.GetAssetPair().Asset
	costUnit := tx.GetAssetPair().CostUnit
	side := tx.GetSide()

	// Total is negative when buy and positive on sell
	ac.status[costUnit] = ac.status[costUnit] + tx.GetTotalPrice()

	if side == SideTypeSell || side == SideTypeTransfer {

		// Less asset since sold or transferred asset.
		ac.status[asset] = ac.status[asset] - tx.GetAssetSize()

	} else {

		// Get more of the asset since buy or have received the asset.
		ac.status[asset] = ac.status[asset] + tx.GetAssetSize()

	}

	return ac

}

func (acc *AccountLog) EnsureAccounts(prototype map[AssetType]float64) AccountEntry {

	for k := range prototype {

		if _, ok := acc.status[k]; !ok {
			acc.status[k] = 0
		}

	}

	return acc
}

func (acc *AccountLog) GetAccountStatus() map[AssetType]float64 {
	return acc.status
}

func (acc *AccountLog) HasTransaction() bool {
	return acc.tx != nil
}

func (acc *AccountLog) GetID() string {
	return acc.tx.GetID()
}

func (acc *AccountLog) GetExchange() string {
	return acc.tx.GetExchange()
}

func (acc *AccountLog) GetSide() SideType {
	return acc.tx.GetSide()
}

func (acc *AccountLog) GetCreatedAt() time.Time {
	return acc.tx.GetCreatedAt()
}

func (acc *AccountLog) GetAssetSize() float64 {
	return acc.tx.GetAssetSize()
}

func (acc *AccountLog) GetPricePerUnit() float64 {
	return acc.tx.GetPricePerUnit()
}

func (acc *AccountLog) GetFee() float64 {
	return acc.tx.GetFee()
}

func (acc *AccountLog) GetTotalPrice() float64 {
	return acc.tx.GetTotalPrice()
}

func (acc *AccountLog) GetAssetPair() AssetPair {
	return acc.tx.GetAssetPair()
}
