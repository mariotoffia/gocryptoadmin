package common

import (
	"fmt"
	"sort"
	"time"

	"github.com/mariotoffia/gocryptoadmin/utils"
)

// AccountStatus contains all assets and their current status.
type AccountStatus map[AssetType]float64

// ExchangeAccountStatus contains the account status for each
// exchange. The _"all"_, the the complete account status.
type ExchangeAccountStatus map[string]AccountStatus

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
	GetAccountStatus() AccountStatus
	// EnsureAccountsjust makes sure that all `AssetType`(s) do exist
	// in account status based on the _prototype_.
	EnsureAccounts(prototype AccountStatus) AccountEntry
}

type AccountLog struct {
	tx     TransactionEntry
	status AccountStatus
	sorted []string
}

// NewAccountLog creates a new `AccountLog` with `AccountLog.Status` initialied.
func NewAccountLog(tx TransactionEntry) *AccountLog {

	return &AccountLog{
		tx:     tx,
		status: AccountStatus{},
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

	if asset != costUnit && tx.GetPricePerUnit() != 1.0 {

		if side == SideTypeSell || side == SideTypeTransfer {

			// Less asset since sold or transferred asset.
			acc.status[asset] = acc.status[asset] - tx.GetAssetSize()

		} else {

			// Get more of the asset since buy or have received the asset.
			acc.status[asset] = acc.status[asset] + tx.GetAssetSize()

		}

	}

	acc.setSortedKeys()

	return acc

}

func (acc *AccountLog) EnsureAccounts(prototype AccountStatus) AccountEntry {

	for k := range prototype {

		if _, ok := acc.status[k]; !ok {
			acc.status[k] = 0
		}

	}

	acc.setSortedKeys()

	return acc
}

func (acc *AccountLog) GetAccountStatus() AccountStatus {
	return acc.status
}

func (acc *AccountLog) GetID() string {
	return acc.tx.GetID()
}

func (acc *AccountLog) GetTranslatedTotalPrice(asset AssetType) float64 {
	return acc.tx.GetTranslatedTotalPrice(asset)
}

func (acc *AccountLog) GetTranslatedFee(asset AssetType) float64 {
	return acc.tx.GetTranslatedFee(asset)
}

func (acc *AccountLog) GetTranslatedAssets() []AssetType {
	return acc.tx.GetTranslatedAssets()
}

func (acc *AccountLog) GetExchange() string {
	return acc.tx.GetExchange()
}

func (acc *AccountLog) GetSide() SideType {
	return acc.tx.GetSide()
}

func (acc *AccountLog) GetSideIdentifier() string {
	return acc.tx.GetSideIdentifier()
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

func (acc *AccountLog) Clone() TransactionEntry {

	a := &AccountLog{
		tx:     acc.tx.Clone(),
		sorted: acc.sorted,
	}

	if len(acc.status) > 0 {

		a.status = AccountStatus{}
		for k, v := range acc.status {
			a.status[k] = v
		}
	}

	return a

}

func (acc *AccountLog) SplitSize(
	size float64,
) (sized TransactionEntry, overflow TransactionEntry) {

	szd := acc.Clone().(*AccountLog)
	ofl := acc.Clone().(*AccountLog)

	sized, overflow = acc.tx.SplitSize(size)
	percent := sized.GetAssetSize() / size

	szd.tx = sized
	ofl.tx = overflow

	for k, v := range szd.status {

		status := v * percent

		szd.status[k] = utils.ToFixed(status, 8)
		ofl.status[k] = utils.ToFixed(v-status, 8)

	}

	return szd, ofl
}

// ConsoleString implements the `ConsoleFormatter` interface
func (acc *AccountLog) ConsoleHeader() string {
	s := "Exchange\tSide\t\tSide Identifier\t\tDate\t\t\tPair\tSize\t\tPrice\t\tFee\t\tTotal"

	for _, v := range acc.sorted {
		s += fmt.Sprintf("\t\t%s", v)
	}

	/*
		s += "\n-----------------------------------------------------------------------------------" +
			"---------------------------------------------------------------------"

		for i := 0; i < len(acc.sorted); i++ {
			s += "----------------"
		}*/

	return s
}

// ConsoleString implements the `ConsoleFormatter` interface
func (acc *AccountLog) ConsoleString() string {

	var s string
	if acc.GetSide() == SideTypeTransfer {

		s = fmt.Sprintf(
			"%s\t\t%s\t%s\t\t\t%s\t%s\t%f\t%f\t%f\t%f",
			acc.GetExchange(),
			acc.GetSide(),
			acc.GetSideIdentifier(),
			acc.GetCreatedAt().Format("2006-01-02 15:04:05.999999999"),
			acc.GetAssetPair().String(),
			acc.GetAssetSize(),
			acc.GetPricePerUnit(),
			acc.GetFee(),
			acc.GetTotalPrice(),
		)

	} else {

		s = fmt.Sprintf(
			"%s\t\t%s\t\t%s\t\t\t%s\t%s\t%f\t%f\t%f\t%f",
			acc.GetExchange(),
			acc.GetSide(),
			acc.GetSideIdentifier(),
			acc.GetCreatedAt().Format("2006-01-02 15:04:05.999999999"),
			acc.GetAssetPair().String(),
			acc.GetAssetSize(),
			acc.GetPricePerUnit(),
			acc.GetFee(),
			acc.GetTotalPrice(),
		)

	}

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
