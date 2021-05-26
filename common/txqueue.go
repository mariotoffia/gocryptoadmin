package common

import "github.com/mariotoffia/gocryptoadmin/utils"

type FIFOTxQueue interface {
	Enq(n TransactionEntry) FIFOTxQueue
	Deq() TransactionEntry
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

// Deq will dequeue the first pushed `TransactionLog`.
func (q *TxFIFOQueue) Deq() TransactionEntry {
	return q.queue.PopFront().(TransactionEntry)
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

// Enq will enqueue the `TransactionLog` into the _asset_ queue.
func (q *TxAssetFIFOQueues) Enq(asset AssetType, n TransactionEntry) *TxAssetFIFOQueues {

	q.getQueue(asset).Enq(n)
	return q

}

// Deq will dequeue the first pushed `TransactionLog` onto the _asset_ queue.
func (q *TxAssetFIFOQueues) Deq(asset AssetType) TransactionEntry {
	return q.getQueue(asset).Deq().(TransactionEntry)
}

func (q *TxAssetFIFOQueues) IsEmpty(asset AssetType) bool {
	return q.getQueue(asset).IsEmpty()
}

func (q *TxAssetFIFOQueues) Len(asset AssetType) int {
	return q.getQueue(asset).Len()
}

func (q *TxAssetFIFOQueues) getQueue(asset AssetType) FIFOTxQueue {

	queue := q.queues[asset]

	if queue == nil {
		queue = NewTxFIFOQueue()
		q.queues[asset] = queue
	}

	return queue
}
