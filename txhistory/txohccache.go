package txhistory

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/jszwec/csvutil"
	"github.com/mariotoffia/gocryptoadmin/common"
)

type ExchangeOHCEntries struct {
	entries map[string][]common.TxOHCHistory
}

// TxOHCCache keeps OHC entries in a cache.
//
// `common.ExchangeAll` can used when global rates
// is accepted.
type TxOHCCache struct {
	entries map[string]*ExchangeOHCEntries
}

func NewTxOHCCache() *TxOHCCache {

	return &TxOHCCache{
		entries: map[string]*ExchangeOHCEntries{},
	}

}

func (cache *TxOHCCache) GetExchanges(except ...string) []string {

	exchanges := []string{}

	for k := range cache.entries {

		skip := false
		for _, ex := range except {
			if ex == k {
				skip = true
				break
			}
		}

		if !skip {
			exchanges = append(exchanges, k)
		}

	}

	return exchanges
}

func (cache *TxOHCCache) Clear(path string, except ...string) *TxOHCCache {

	err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {

		if info.IsDir() || err != nil {
			return err
		}

		for _, ex := range except {

			file := filepath.Base(path)
			if strings.HasPrefix(file, ex+"_") && strings.HasSuffix(file, ".csv") {
				return nil // Skip
			}

		}

		return os.Remove(path)

	})

	if err != nil {
		panic(err)
	}

	return cache
}

func (cache *TxOHCCache) Load(
	path string,
	addFunc func(cache *TxOHCCache, entries []common.TxOHCHistory),
	exchange ...string) *TxOHCCache {

	if addFunc == nil {

		addFunc = func(cache *TxOHCCache, entries []common.TxOHCHistory) {
			cache.Add(entries)
		}

	}

	err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {

		if info.IsDir() || err != nil {
			return err
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		var entries []common.TxOHCHistory
		if err = csvutil.Unmarshal(data, entries); err != nil {
			return err
		}

		addFunc(cache, entries)

		return nil
	})

	if err != nil {
		panic(err)
	}

	return cache
}

func (cache *TxOHCCache) Store(path string, exchange ...string) *TxOHCCache {

	if len(exchange) == 0 {

		for k := range cache.entries {
			exchange = append(exchange, k)
		}

	}

	for _, ex := range exchange {

		for ap, entry := range cache.entries[ex].entries {

			fqFile := filepath.Join(path, renderFileName(ex, ap, entry))

			data, err := csvutil.Marshal(entry)
			if err != nil {
				panic(err)
			}

			if err = ioutil.WriteFile(fqFile, data, 0644); err != nil {
				panic(err)
			}

		}

	}

	return cache
}

func (cache *TxOHCCache) Add(entries []common.TxOHCHistory, exchange ...string) *TxOHCCache {

	for i := range entries {

		ap := entries[i].GetAssetPair().String()
		ex := entries[i].GetExchange()

		all := append(exchange, ex)

		for _, ex := range all {

			exchangeEntries := cache.entries[ex]

			if exchangeEntries == nil {
				exchangeEntries = &ExchangeOHCEntries{
					entries: map[string][]common.TxOHCHistory{},
				}

				cache.entries[ex] = exchangeEntries
			}

			c := exchangeEntries.entries[ap]
			c = append(c, entries[i])
			exchangeEntries.entries[ap] = c

		}

	}

	return cache
}

func renderFileName(exchange, assetPair string, entry []common.TxOHCHistory) string {

	if len(entry) == 0 {

		panic(
			fmt.Sprintf("zero entries not allowed for asset: %s", assetPair),
		)

	}

	start := entry[0].GetDateTime()
	end := start

	if len(entry) > 1 {
		end = entry[len(entry)-1].GetDateTime()
	}

	return fmt.Sprintf(
		"%s_%s_%s_%s.csv",
		exchange,
		assetPair,
		start.Format("2006-01-02"),
		end.Format("2006-01-02"),
	)

}
