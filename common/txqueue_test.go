package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUseTxFIFOQueue(t *testing.T) {

	q := NewTxFIFOQueue()

	q.Enq(&TransactionLog{ID: "1"}).
		Enq(&TransactionLog{ID: "2"}).
		Enq(&TransactionLog{ID: "3"})

	assert.Equal(t, 3, q.Len())
	assert.Equal(t, false, q.IsEmpty())
	assert.Equal(t, "1", q.Deq().GetID())
	assert.Equal(t, "2", q.Deq().GetID())
	assert.Equal(t, "3", q.Deq().GetID())
	assert.Equal(t, true, q.IsEmpty())
}

func TestTxFIFOQueueDequeueUntil(t *testing.T) {

	q := NewTxFIFOQueue()

	q.Enq(&TransactionLog{ID: "1"}).
		Enq(&TransactionLog{ID: "2"}).
		Enq(&TransactionLog{ID: "3"})

	assert.Equal(t, 3, q.Len())
	assert.Equal(t, false, q.IsEmpty())

	entries, res := q.DequeueUntil(func(tx TransactionEntry) DequeueUntilResult {

		if tx.GetID() == "3" {
			return DequeueUntilResultOverflow
		}

		return DequeueUntilResultContinue
	})

	assert.Equal(t, DequeueUntilResultOverflow, res)
	assert.Equal(t, true, q.IsEmpty())
	require.Equal(t, 3, len(entries))

	assert.Equal(t, "1", entries[0].GetID())
	assert.Equal(t, "2", entries[1].GetID())
	assert.Equal(t, "3", entries[2].GetID())

}

func TestUseTxAssetFIFOQueues(t *testing.T) {

	q := NewTxAssetFIFOQueues()

	q.Enq(AssetTypeBTC, &TransactionLog{ID: "1"}).
		Enq(AssetTypeLTC, &TransactionLog{ID: "2"}).
		Enq(AssetTypeBTC, &TransactionLog{ID: "3"})

	assert.Equal(t, 2, q.Len(AssetTypeBTC))
	assert.Equal(t, 1, q.Len(AssetTypeLTC))
	assert.Equal(t, false, q.IsEmpty(AssetTypeLTC))
	assert.Equal(t, false, q.IsEmpty(AssetTypeBTC))

	assert.Equal(t, "2", q.Deq(AssetTypeLTC).GetID())
	assert.Equal(t, "1", q.Deq(AssetTypeBTC).GetID())
	assert.Equal(t, "3", q.Deq(AssetTypeBTC).GetID())

	assert.Equal(t, true, q.IsEmpty(AssetTypeLTC))
	assert.Equal(t, true, q.IsEmpty(AssetTypeBTC))
}
