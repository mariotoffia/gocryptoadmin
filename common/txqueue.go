package common

type TxFIFOQueue struct {
	queue *Queue
}

func NewTxFIFOQueue() *TxFIFOQueue {

	return &TxFIFOQueue{
		queue: NewQueue(),
	}

}

// Enq will enqueue the `TransactionLog`
func (q *TxFIFOQueue) Enq(n *TransactionLog) *TxFIFOQueue {

	q.queue.PushBack(n)
	return q

}

func (q *TxFIFOQueue) Deq() *TransactionLog {
	return q.queue.PopFront().(*TransactionLog)
}

func (q *TxFIFOQueue) IsEmpty() bool {
	return q.queue.empty()
}

func (q *TxFIFOQueue) Len() int {
	return q.queue.Len()
}
