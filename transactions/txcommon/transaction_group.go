package txcommon

type TransactionGroup struct {
	Transaction
	Excluded bool          `csv:"excluded" json:"excluded"`
	Tx       []Transaction `csv:"-" json:"tx,omitempty"`
}
