package processors

import (
	"github.com/mariotoffia/gocryptoadmin/common"
	"github.com/mariotoffia/gocryptoadmin/processors/procutils"
)

type TxBuySellGroupProcessor struct {
	cache          *procutils.TxGroupCache
	transactions   []common.TxGroupEntry
	flushProcessor common.TxGroupProcessor
}
