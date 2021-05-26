package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUseTxFIFOQueue(t *testing.T) {

	q := NewTxFIFOQueue()

	q.Enq(&TransactionLog{ID: "1"}).
		Enq(&TransactionLog{ID: "2"}).
		Enq(&TransactionLog{ID: "3"})

	assert.Equal(t, 3, q.Len())
	assert.Equal(t, false, q.IsEmpty())
	assert.Equal(t, "1", q.Deq().ID)
	assert.Equal(t, "2", q.Deq().ID)
	assert.Equal(t, "3", q.Deq().ID)
	assert.Equal(t, true, q.IsEmpty())
}
