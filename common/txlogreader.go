package common

type TransactionLogReader interface {
	Unmarshal(data []byte) []TransactionLog
}
