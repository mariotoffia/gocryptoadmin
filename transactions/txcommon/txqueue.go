package txcommon

type TxQueue []*Transaction

func (q *TxQueue) Push(n *Transaction) {
	*q = append(*q, n)
}

func (q *TxQueue) PushFront(n *Transaction) {
	*q = append(TxQueue{n}, *q...)
}

func (q *TxQueue) Pop() (n *Transaction) {
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
