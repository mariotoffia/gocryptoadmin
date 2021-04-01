package common

type TxGroup struct {
	Transaction
	Tx []Transaction `csv:"-" json:"tx"`
}
