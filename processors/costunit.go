package processors

import (
	"fmt"
	"time"

	"github.com/mariotoffia/gocryptoadmin/common"
	"github.com/mariotoffia/gocryptoadmin/txhistory"
)

// PriceEntryCalculator calculates the price for the _entry_ and returns the value.
type PriceEntryCalculator func(side common.SideType, entry *txhistory.ResolvedOHCEntry) float64

// CostUnitProcessor implements interface `TxEntryProcessor`
// and should be executed on raw imported transactions.
type CostUnitProcessor struct {
	transactions []common.TransactionLog
	resolver     *txhistory.TxOHCResolver
	tracked      []common.AssetType
	priceCalc    PriceEntryCalculator
}

func NewCostUnitProcessor(
	resolver *txhistory.TxOHCResolver,
	priceCalc PriceEntryCalculator,
) *CostUnitProcessor {

	if priceCalc == nil {

		priceCalc = func(side common.SideType, entry *txhistory.ResolvedOHCEntry) float64 {

			return (entry.Entry.GetLow() + entry.Entry.GetHigh()) / 2

		}

	}

	return &CostUnitProcessor{
		transactions: []common.TransactionLog{},
		resolver:     resolver,
		tracked:      []common.AssetType{},
		priceCalc:    priceCalc,
	}

}

func (proc *CostUnitProcessor) RegisterAsset(asset ...common.AssetType) {

	proc.tracked = append(proc.tracked, asset...)

}

func (proc *CostUnitProcessor) Reset() {
	proc.transactions = []common.TransactionLog{}
}

func (proc *CostUnitProcessor) ProcessMany(tx []common.TransactionLog) {

	for i := range tx {

		proc.Process(tx[i])

	}

}

func (proc *CostUnitProcessor) Process(tx common.TransactionLog) {

	for idx, asset := range proc.tracked {

		if tx.TranslatedTotalPrice == nil {
			tx.TranslatedTotalPrice = map[string]float64{}
		}

		if tx.TranslatedFee == nil {
			tx.TranslatedFee = map[string]float64{}
		}

		if tx.CostUnit == asset {

			tx.TranslatedTotalPrice[string(asset)] = tx.TotalPrice
			tx.TranslatedFee[string(asset)] = tx.Fee

			continue

		}

		entries, ok := proc.resolver.ResolveToTarget(
			tx.CreatedAt,
			tx.CostUnit,
			asset,
			tx.Exchange,
			common.ExchangeAll,
		)

		if !ok {

			panic(
				fmt.Sprintf(
					"[index: %d] - could not get asset: %s, via asset-pair: %s at %s",
					idx,
					asset,
					tx.AssetPair.String(),
					tx.CreatedAt.Format(time.RFC3339),
				),
			)

		}

		tot := tx.TotalPrice
		fee := tx.Fee

		for _, entry := range entries {

			price := proc.priceCalc(tx.Side, &entry)
			tot *= price
			fee *= price

		}

		tx.TranslatedTotalPrice[string(asset)] = tot
		tx.TranslatedFee[string(asset)] = fee

	}

	proc.transactions = append(proc.transactions, tx)
}

func (proc *CostUnitProcessor) Flush() []common.TransactionLog {

	tx := proc.transactions
	proc.transactions = []common.TransactionLog{}

	return tx

}
