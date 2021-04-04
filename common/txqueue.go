package common

type TxQueue []*TransactionLog

func (q *TxQueue) Push(n *TransactionLog) {
	*q = append(*q, n)
}

func (q *TxQueue) PushFront(n *TransactionLog) {
	*q = append(TxQueue{n}, *q...)
}

func (q *TxQueue) Pop() (n *TransactionLog) {
	n = (*q)[0]
	*q = (*q)[1:]
	return
}

func (q *TxQueue) IsEmpty() bool {
	return len(*q) == 0
}

func (q *TxQueue) Len() int {
	return len(*q)
}
