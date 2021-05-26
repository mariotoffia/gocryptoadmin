package common

import "time"

type TxPair interface {
	Sell() TxGroupEntry
	Buy() TxGroupEntry
	SellCreatedAt() time.Time
	BuyCreatedAt() time.Time
	SellTotal() float64
	BuyTotal() float64
	// SellProminent gets the entry with the most `AssetSize`
	SellProminent() TransactionEntry
	// BuyProminent gets the entry with the most `AssetSize`
	BuyProminent() TransactionEntry
}

// TxPairEntry is a pair of transaction entries.
//
// This could be e.g. a BUY -> SELL pair. Since multiple entries may consitute
// zero, one or both sides both sides are represented as `TxGroupEntry`
// (even if it is just considered as a single transaction).
type TxPairEntry struct {
	SellTx TxGroupEntry
	BuyTx  TxGroupEntry
}

func (tx *TxPairEntry) Sell() TxGroupEntry {

	return tx.SellTx

}

func (tx *TxPairEntry) Buy() TxGroupEntry {

	return tx.BuyTx

}

func (tx *TxPairEntry) SellCreatedAt() time.Time {

	return tx.SellTx.GetCreatedAt()

}

func (tx *TxPairEntry) BuyCreatedAt() time.Time {

	return tx.BuyTx.GetCreatedAt()

}

func (tx *TxPairEntry) SellTotal() float64 {

	return tx.SellTx.GetTotalPrice()

}

func (tx *TxPairEntry) BuyTotal() float64 {

	return tx.BuyTx.GetTotalPrice()

}

// SellProminent gets the entry with the most `AssetSize`
func (tx *TxPairEntry) SellProminent() TransactionEntry {

	return tx.SellTx.GetMostProminentSizeTransactionLog()

}

// BuyProminent gets the entry with the most `AssetSize`
func (tx *TxPairEntry) BuyProminent() TransactionEntry {

	return tx.BuyTx.GetMostProminentSizeTransactionLog()

}
