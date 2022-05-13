package common

type TransactionLogReader interface {
	Unmarshal(data []byte) []TransactionLog
	SetExchange(name string) TransactionLogReader
}
