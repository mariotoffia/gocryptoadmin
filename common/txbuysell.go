package common

import (
	"fmt"
)

type TxBuySellEntry interface {
	TransactionEntry
	GetBuy() *TxBuyGroupLog
	GetSell() TransactionEntry
}

// TxBuySellLog is a pair of transaction entries.
//
// This could be e.g. a BUY -> SELL pair. Since multiple entries may consitute
// zero, one or both sides both sides are represented as `TxGroupEntry`
// (even if it is just considered as a single transaction).
type TxBuySellLog struct {
	TransactionLog
	SellTx TransactionEntry
	BuyTx  TxBuyGroupLog
}

func NewTxBuySellLog(
	sellTx TransactionEntry,
	buyTx []TransactionEntry,
) *TxBuySellLog {

	if len(buyTx) == 0 {
		panic("trying to create a buysell entry without any buys!")
	}

	txg := NewTxBuyGroupLog(buyTx)

	log := &TxBuySellLog{
		TransactionLog: TransactionLog{
			ID:                   fmt.Sprintf("%s-buysell", buyTx[0].GetID()),
			Exchange:             sellTx.GetExchange(),
			Side:                 SideTypeBuySell,
			SideIdentifier:       sellTx.GetSideIdentifier(),
			CreatedAt:            txg.GetCreatedAt(),
			AssetSize:            sellTx.GetAssetSize(),
			PricePerUnit:         sellTx.GetPricePerUnit(),
			Fee:                  sellTx.GetFee(),
			TotalPrice:           sellTx.GetTotalPrice(),
			TranslatedTotalPrice: nil,
			TranslatedFee:        nil,
			AssetPair:            sellTx.GetAssetPair(),
		},
		SellTx: sellTx,
		BuyTx:  *txg,
	}

	return log
}

func (tx *TxBuySellLog) GetSell() TransactionEntry {
	return tx.SellTx
}

func (tx *TxBuySellLog) GetBuy() *TxBuyGroupLog {
	return &tx.BuyTx
}
