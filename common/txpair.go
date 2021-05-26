package common

import "time"

type TxPair interface {
	SideA() TxGroupEntry
	SideB() TxGroupEntry
	SideACreatedAt() time.Time
	SideBCreatedAt() time.Time
	SideATotal() float64
	SideBTotal() float64
	// SideAProminent gets the entry with the most `AssetSize`
	SideAProminent() TransactionLog
	// SideBProminent gets the entry with the most `AssetSize`
	SideBProminent() TransactionLog
}

// TxPairEntry is a pair of transaction entries.
//
// This could be e.g. a BUY -> SELL pair. Since multiple entries may consitute
// zero, one or both sides both sides are represented as `TxGroupEntry`
// (even if it is just considered as a single transaction).
type TxPairEntry struct {
	SideATx TxGroupEntry
	SideBTx TxGroupEntry
}

func (tx *TxPairEntry) SideA() TxGroupEntry {

	return tx.SideATx

}

func (tx *TxPairEntry) SideB() TxGroupEntry {

	return tx.SideBTx

}

func (tx *TxPairEntry) SideACreatedAt() time.Time {

	return tx.SideATx.GetCreatedAt()

}

func (tx *TxPairEntry) SideBCreatedAt() time.Time {

	return tx.SideBTx.GetCreatedAt()

}

func (tx *TxPairEntry) SideATotal() float64 {

	return tx.SideATx.GetTotalPrice()

}

func (tx *TxPairEntry) SideBTotal() float64 {

	return tx.SideBTx.GetTotalPrice()

}

// SideAProminent gets the entry with the most `AssetSize`
func (tx *TxPairEntry) SideAProminent() TransactionLog {

	return tx.SideATx.GetMostProminentSizeTransactionLog()

}

// SideBProminent gets the entry with the most `AssetSize`
func (tx *TxPairEntry) SideBProminent() TransactionLog {

	return tx.SideBTx.GetMostProminentSizeTransactionLog()

}
