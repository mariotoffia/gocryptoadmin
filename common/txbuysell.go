package common

import (
	"fmt"
)

type TxBuySellEntry interface {
	TransactionEntry
	GetBuy() *TxBuyGroupLog
	GetSell() TransactionEntry
}

// TxBuySellLog is a _SELL_ that have corresponding _BUYs_
//
// All of the _BUY_ `TransactionEntry` instances are encapsulated
// in a `TxBuyGroup` where functions are overloaded to handle both
// _BUY_ and _SELL_ transactions.
//
// It may contain _SELL_ transactions that gained a certain `AssetType`
// that is now sold. For example _SELL_ 100LT -> 1BTC and then _SELL_
// 1BTC -> 20.0000EUR. Thus the former sell is part of the buy array.
//
// Since _SELL_ transaction is quite different, therefore be cautious when
// e.g. sum up all buy transactions etc. This has been accommodated in the
// `TxBuyGroupLog` overrides.
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
			TranslatedTotalPrice: map[string]float64{},
			TranslatedFee:        map[string]float64{},
			AssetPair:            sellTx.GetAssetPair(),
		},
		SellTx: sellTx,
		BuyTx:  *txg,
	}

	for _, asset := range sellTx.GetTranslatedAssets() {

		log.TranslatedTotalPrice[string(asset)] = sellTx.GetTranslatedTotalPrice(asset)
		log.TranslatedFee[string(asset)] = sellTx.GetTranslatedFee(asset)

	}

	return log
}

func (tx *TxBuySellLog) GetSell() TransactionEntry {
	return tx.SellTx
}

func (tx *TxBuySellLog) GetBuy() *TxBuyGroupLog {
	return &tx.BuyTx
}

func (tx *TxBuySellLog) Clone() TransactionEntry {

	buyTx := tx.BuyTx.Clone().(*TxBuyGroupLog)

	log := &TxBuySellLog{
		TransactionLog: tx.TransactionLog,
		SellTx:         tx.SellTx.Clone(),
		BuyTx:          *buyTx,
	}

	return log
}
