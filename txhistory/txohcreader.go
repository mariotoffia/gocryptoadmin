package txhistory

import (
	"sort"
	"time"

	"github.com/mariotoffia/gocryptoadmin/common"
)

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
) []common.TxOHCHistory {

	list := []common.TxOHCHistory{}

	for i := range reader {

		entries := txr.readers[reader[i]].Read(pair, since, interval)
		list = append(list, entries...)

	}

	sort.Slice(list, func(i, j int) bool {

		return list[i].GetDateTime().Before(
			list[j].GetDateTime(),
		)

	})

	return list
}
