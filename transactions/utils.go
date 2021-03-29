package transactions

import (
	"sort"

	"github.com/mariotoffia/gocryptoadmin/transactions/txcommon"
)

// MergeLogs into _dst_ by sorting the `txcommon.Transaction.CreatedAt`.
func MergeLogs(src []txcommon.Transaction, dst []txcommon.Transaction) []txcommon.Transaction {
	return SortLogs(append(dst, src...))
}

func SortLogs(logs []txcommon.Transaction) []txcommon.Transaction {

	sort.Slice(logs, func(i, j int) bool {
		return logs[i].CreatedAt.Before(logs[j].CreatedAt)
	})

	return logs
}
