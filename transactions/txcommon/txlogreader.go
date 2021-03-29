package txcommon

type TransactionLogReader interface {
	Unmarshal(data []byte) []Transaction
}
