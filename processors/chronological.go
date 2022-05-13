package processors

import (
	"sort"
	"strconv"
	"strings"

	"github.com/mariotoffia/gocryptoadmin/common"
)

type ChronologicalTxEntryProcessor struct {
	tx []common.TransactionLog
}

type ChronologicalGroupTxEntryProcessor struct {
	tx []common.TxGroupEntry
}

func NewChronologicalTxEntryProcessor() *ChronologicalTxEntryProcessor {
	return &ChronologicalTxEntryProcessor{}
}

func NewChronologicalGroupTxEntryProcessor() *ChronologicalGroupTxEntryProcessor {
	return &ChronologicalGroupTxEntryProcessor{}
}

func (c *ChronologicalTxEntryProcessor) Reset() {
	c.tx = nil
}

func (c *ChronologicalTxEntryProcessor) ProcessMany(tx []common.TransactionLog) {

	if nil == c.tx {
		c.tx = tx
		return
	}

	c.tx = append(c.tx, tx...)
}

func (c *ChronologicalTxEntryProcessor) Process(tx common.TransactionLog) {
	c.tx = append(c.tx, tx)
}

func (c *ChronologicalTxEntryProcessor) Flush() []common.TransactionLog {

	sort.Slice(c.tx, func(i, j int) bool {

		if c.tx[i].CreatedAt.Equal(c.tx[j].CreatedAt) {

			if c.tx[i].Exchange == c.tx[j].Exchange {

				if c.tx[i].Exchange == "cbx" &&
					c.tx[i].AssetPair == c.tx[j].AssetPair &&
					!strings.HasPrefix(c.tx[i].ID, "M") &&
					!strings.HasPrefix(c.tx[j].ID, "M") {

					li, err := strconv.ParseInt(c.tx[i].ID, 10, 64)
					if err != nil {
						panic(err)
					}

					lj, err := strconv.ParseInt(c.tx[j].ID, 10, 64)
					if err != nil {
						panic(err)
					}

					return li < lj

				}
			}

		}

		return c.tx[i].CreatedAt.Before(c.tx[j].CreatedAt)
	})

	return c.tx
}

func (c *ChronologicalGroupTxEntryProcessor) Reset() {
	c.tx = nil
}

func (c *ChronologicalGroupTxEntryProcessor) ProcessMany(tx []common.TxGroupEntry) {

	if nil == c.tx {
		c.tx = tx
		return
	}

	c.tx = append(c.tx, tx...)
}

func (c *ChronologicalGroupTxEntryProcessor) Process(tx common.TxGroupEntry) {
	c.tx = append(c.tx, tx)
}

func (c *ChronologicalGroupTxEntryProcessor) Flush() []common.TxGroupEntry {

	sort.Slice(c.tx, func(i, j int) bool {

		if c.tx[i].GetCreatedAt().Equal(c.tx[j].GetCreatedAt()) {

			if c.tx[i].GetExchange() == c.tx[j].GetExchange() {

				if c.tx[i].GetExchange() == "cbx" &&
					c.tx[i].GetAssetPair().String() == c.tx[j].GetAssetPair().String() {

					li, err := strconv.ParseInt(c.tx[i].ID, 10, 64)
					if err != nil {
						panic(err)
					}

					lj, err := strconv.ParseInt(c.tx[j].ID, 10, 64)
					if err != nil {
						panic(err)
					}

					return li < lj

				}
			}

		}

		return c.tx[i].GetCreatedAt().Before(c.tx[j].GetCreatedAt())
	})

	return c.tx
}
