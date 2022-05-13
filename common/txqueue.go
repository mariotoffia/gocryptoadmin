package common

import "github.com/mariotoffia/gocryptoadmin/utils"

type DequeueUntilResult int

const (
	// DequeueUntilResultContinue denotes that it should continue it dequeue and add the entry to
	// the dequeue result.
	DequeueUntilResultContinue DequeueUntilResult = 0
	// DequeueUntilResultOverflow specifies that the last entry should be added, but it will overflow
	// the actual amount.
	//
	// This is useful, when partial `common.TransactionEntry` must be used by the client.
	DequeueUntilResultOverflow DequeueUntilResult = 1
	// DequeueUntilResultDone will make the `DequeueUntil` to stop with Success status. It will
	// add the last dequeue item onto the result.
	DequeueUntilResultDone DequeueUntilResult = 2
	// DequeueUntilResultUnderflow is when not enough entries where captured to satisfy the accept function.
	DequeueUntilResultUnderflow DequeueUntilResult = 3
)

type FIFOTxQueue interface {
	Enq(n TransactionEntry) FIFOTxQueue
	PutBack(n TransactionEntry) FIFOTxQueue
	Deq() TransactionEntry

	DequeueUntil(
		accept func(tx TransactionEntry) DequeueUntilResult,
	) ([]TransactionEntry, DequeueUntilResult)

	IsEmpty() bool
	Len() int
}

type TxFIFOQueue struct {
	queue *utils.Queue
}

type TxAssetFIFOQueues struct {
	queues map[AssetType]FIFOTxQueue
}

func NewTxFIFOQueue() *TxFIFOQueue {

	return &TxFIFOQueue{
		queue: utils.NewQueue(),
	}

}

// Enq will enqueue the `TransactionLog`
func (q *TxFIFOQueue) Enq(n TransactionEntry) FIFOTxQueue {

	q.queue.PushBack(n)
	return q

}

// PutBack will enqueue the `TransactionEntry` but make it available
// on subsequent `Deq` (ordinary it would be last, since last enqueued)
func (q *TxFIFOQueue) PutBack(n TransactionEntry) FIFOTxQueue {

	q.queue.PushFront(n)
	return q

}

// Deq will dequeue the first pushed `TransactionEntry`.
func (q *TxFIFOQueue) Deq() TransactionEntry {
	return q.queue.PopFront().(TransactionEntry)
}

func (q *TxFIFOQueue) DequeueUntil(
	accept func(tx TransactionEntry) DequeueUntilResult,
) ([]TransactionEntry, DequeueUntilResult) {

	entries := []TransactionEntry{}

	for !q.queue.Empty() {

		entry := q.Deq()
		entries = append(entries, entry)

		if res := accept(entry); res != DequeueUntilResultContinue {
			return entries, res
		}

	}

	if len(entries) != 0 {
		return entries, DequeueUntilResultUnderflow
	}

	return entries, DequeueUntilResultDone

}

func (q *TxFIFOQueue) IsEmpty() bool {
	return q.queue.Empty()
}

func (q *TxFIFOQueue) Len() int {
	return q.queue.Len()
}

func NewTxAssetFIFOQueues() *TxAssetFIFOQueues {

	return &TxAssetFIFOQueues{
		queues: map[AssetType]FIFOTxQueue{},
	}

}

func (q *TxAssetFIFOQueues) Reset() *TxAssetFIFOQueues {
	q.queues = map[AssetType]FIFOTxQueue{}
	return q
}

// Enq will enqueue the `TransactionLog` into the _asset_ queue.
func (q *TxAssetFIFOQueues) Enq(asset AssetType, n TransactionEntry) *TxAssetFIFOQueues {

	q.getQueue(asset).Enq(n)
	return q

}

// PutBack will enqueue the `TransactionEntry` but make it available
// on subsequent `Deq` (ordinary it would be last, since last enqueued)
func (q *TxAssetFIFOQueues) PutBack(asset AssetType, n TransactionEntry) *TxAssetFIFOQueues {

	q.getQueue(asset).PutBack(n)
	return q

}

// Deq will dequeue the first pushed `TransactionLog` onto the _asset_ queue.
func (q *TxAssetFIFOQueues) Deq(asset AssetType) TransactionEntry {
	return q.getQueue(asset).Deq().(TransactionEntry)
}

func (q *TxAssetFIFOQueues) DequeueUntil(
	asset AssetType,
	accept func(tx TransactionEntry) DequeueUntilResult,
) ([]TransactionEntry, DequeueUntilResult) {

	return q.getQueue(asset).DequeueUntil(accept)

}

func (q *TxAssetFIFOQueues) IsEmpty(asset AssetType) bool {
	return q.getQueue(asset).IsEmpty()
}

func (q *TxAssetFIFOQueues) Len(asset AssetType) int {
	return q.getQueue(asset).Len()
}

// DequeueAll will dequeue assets entries, from the queue
func (q *TxAssetFIFOQueues) DequeueAll() []TransactionEntry {

	entries := []TransactionEntry{}
	for _, queue := range q.queues {

		for !queue.IsEmpty() {
			entries = append(entries, queue.Deq())
		}
	}

	return entries
}
func (q *TxAssetFIFOQueues) TotalLen() int {

	cnt := 0

	for _, q := range q.queues {
		cnt += q.Len()
	}

	return cnt
}

func (q *TxAssetFIFOQueues) getQueue(asset AssetType) FIFOTxQueue {

	queue := q.queues[asset]

	if queue == nil {
		queue = NewTxFIFOQueue()
		q.queues[asset] = queue
	}

	return queue
}
