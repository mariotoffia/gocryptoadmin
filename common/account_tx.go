package common

type AccountTransaction struct {
	Transaction
	AccountTxID int64 `csv:"account tx" json:"atx"`
}
