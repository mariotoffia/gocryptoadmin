package txhistory

import (
	"time"

	"github.com/mariotoffia/gocryptoadmin/common"
)

// TODO: need to handle cache

type TxOHCReader struct {
	readers map[string]common.TxOHCReader
}

func NewTxOHCReader() *TxOHCReader {

	return &TxOHCReader{
		readers: map[string]common.TxOHCReader{},
	}

}

func (txr *TxOHCReader) Register(name string, reader common.TxOHCReader) *TxOHCReader {

	reader.SetExchangeName(name)
	txr.readers[name] = reader

	return txr

}

func (txr *TxOHCReader) Read(
	pair common.AssetPair,
	since time.Time,
	interval time.Duration,
	reader ...string,
) []common.TxOHCHistoryEntry {

	if len(reader) == 1 {
		return txr.readers[reader[0]].Read(pair, since, interval)
	}

	list := []common.TxOHCHistoryEntry{}

	for i := range reader {

		entries := txr.readers[reader[i]].Read(pair, since, interval)
		list = append(list, entries...)

	}

	return list
}
