package common

type TransactionLogReader interface {
	Unmarshal(data []byte) []Transaction
}
