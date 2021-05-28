package common

import (
	"sort"
	"time"
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

	sized, overflow = acc.tx.SplitSize(size)

	szd := &AccountLog{
		tx:     sized,
		sorted: acc.sorted,
	}

	ofl := &AccountLog{
		tx:     overflow,
		sorted: acc.sorted,
	}

	if len(acc.status) > 0 {

		szd.status = AccountStatus{}
		ofl.status = AccountStatus{}

		for k, v := range acc.status {

			szd.status[k] = v
			ofl.status[k] = v

		}

	}

	percent := size / acc.GetAssetSize()

	asset := acc.GetAssetPair().Asset
	costUnit := acc.GetAssetPair().CostUnit
	side := acc.GetSide()

	// Belows is inverting the _percent_ amount that was processed
	// using `NextAccountLog`. The _ofl_ is not affected, only _szd_ will
	// have more or less in the account.
	szd.status[costUnit] = szd.status[costUnit] - acc.GetTotalPrice()*percent

	if asset != costUnit && acc.GetPricePerUnit() != 1.0 {

		if side == SideTypeSell || side == SideTypeTransfer {
			szd.status[asset] = szd.status[asset] + szd.GetAssetSize()
		} else {

			// Get more of the asset since buy or have received the asset.
			szd.status[asset] = szd.status[asset] - szd.GetAssetSize()

		}

	}

	return szd, ofl
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
