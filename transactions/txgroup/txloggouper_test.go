package txgroup

import (
	"fmt"
	"testing"

	"github.com/mariotoffia/gocryptoadmin/transactions"
	"github.com/mariotoffia/gocryptoadmin/transactions/coinbasepro"
)

func TestGroupTxLog(t *testing.T) {

	tx := transactions.NewTxLogReader().
		UseDir("../../data").
		SortRead().
		IgnoreUnknownFiles().
		RegisterReader("coinbasepro", coinbasepro.NewTransactionLogReader()).
		Read()

	tg := NewTxLogGrouperImpl(tx).
		GroupViaTimeWindow().
		DumpGroup(true).
		TransactionGroups()

	fmt.Println(len(tg))
}
