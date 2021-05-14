package common

import (
	"fmt"
	"sort"
	"time"
)

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
}

type AccountLog struct {
	tx     TransactionEntry
	status map[AssetType]float64
	sorted []string
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

	acc := NewAccountLog(tx)

	if previous != nil {

		for k, v := range previous.GetAccountStatus() {
			acc.status[k] = v
		}

	}

	asset := tx.GetAssetPair().Asset
	costUnit := tx.GetAssetPair().CostUnit
	side := tx.GetSide()

	// Total is negative when buy and positive on sell
	acc.status[costUnit] = acc.status[costUnit] + tx.GetTotalPrice()

	if side == SideTypeSell || side == SideTypeTransfer {

		// Less asset since sold or transferred asset.
		acc.status[asset] = acc.status[asset] - tx.GetAssetSize()

	} else {

		// Get more of the asset since buy or have received the asset.
		acc.status[asset] = acc.status[asset] + tx.GetAssetSize()

	}

	acc.setSortedKeys()

	return acc

}

func (acc *AccountLog) EnsureAccounts(prototype map[AssetType]float64) AccountEntry {

	for k := range prototype {

		if _, ok := acc.status[k]; !ok {
			acc.status[k] = 0
		}

	}

	acc.setSortedKeys()

	return acc
}

func (acc *AccountLog) GetAccountStatus() map[AssetType]float64 {
	return acc.status
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

// ConsoleString implements the `ConsoleFormatter` interface
func (acc *AccountLog) ConsoleHeader() string {
	s := "Exchange\tSide\tDate\t\t\tPair\tSize\t\tPrice\t\tFee\t\tTotal"

	for _, v := range acc.sorted {
		s += fmt.Sprintf("\t\t%s", v)
	}

	s += "\n---------------------------------------------------------" +
		"-----------------------------------------------------" +
		"-----------------------------------------------------" +
		"-----------------------------------------------------"

	return s
}

// ConsoleString implements the `ConsoleFormatter` interface
func (acc *AccountLog) ConsoleString() string {

	s := fmt.Sprintf(
		"%s\t%s\t%s\t%s\t%f\t%f\t%f\t%f",
		acc.GetExchange(),
		acc.GetSide(),
		acc.GetCreatedAt().Format("2006-01-02 15:04:05.999999999"),
		acc.GetAssetPair().String(),
		acc.GetAssetSize(),
		acc.GetPricePerUnit(),
		acc.GetFee(),
		acc.GetTotalPrice(),
	)

	for _, v := range acc.sorted {

		s += fmt.Sprintf("\t%f", acc.status[AssetType(v)])

	}

	return s
}

func (acc *AccountLog) setSortedKeys() {

	keys := make([]string, len(acc.status))
	i := 0

	for k := range acc.status {
		keys[i] = string(k)
		i++
	}

	sort.Strings(keys)
	acc.sorted = keys

}
