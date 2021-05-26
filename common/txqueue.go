package common

import "github.com/mariotoffia/gocryptoadmin/utils"

type TxFIFOQueue struct {
	queue *utils.Queue
}

func NewTxFIFOQueue() *TxFIFOQueue {

	return &TxFIFOQueue{
		queue: utils.NewQueue(),
	}

}

// Enq will enqueue the `TransactionLog`
func (q *TxFIFOQueue) Enq(n *TransactionLog) *TxFIFOQueue {

	q.queue.PushBack(n)
	return q

}

// Deq will dequeue the first pushed `TransactionLog`.
func (q *TxFIFOQueue) Deq() *TransactionLog {
	return q.queue.PopFront().(*TransactionLog)
}

func (q *TxFIFOQueue) IsEmpty() bool {
	return q.queue.Empty()
}

func (q *TxFIFOQueue) Len() int {
	return q.queue.Len()
}
