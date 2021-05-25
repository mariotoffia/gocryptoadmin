package txhistory

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

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
	addFunc func(cache *TxOHCCache, exchange string, entries []common.TxOHCHistory),
	exchange ...string) *TxOHCCache {

	if addFunc == nil {

		addFunc = func(
			cache *TxOHCCache, exchange string, entries []common.TxOHCHistory,
		) {
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
		if err = csvutil.Unmarshal(data, &entries); err != nil {
			return err
		}

		exchange := strings.Split(filepath.Base(path), "_")[0]

		sort.Slice(entries, func(i, j int) bool {

			return entries[i].GetDateTime().Before(
				entries[j].GetDateTime(),
			)

		})

		addFunc(cache, exchange, entries)

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

func (cache *TxOHCCache) GetEntryForAssset(
	assetPair common.AssetPair,
	at time.Time,
	exchange ...string,
) (*common.TxOHCHistory, string) {

	if len(exchange) == 0 {
		exchange = []string{common.ExchangeAll}
	}

	ap := assetPair.String()

	for _, ex := range exchange {

		if entries, ok := cache.entries[ex]; ok {

			if c, ok := entries.entries[ap]; ok {

				if entry, ok := cache.FindEntry(c, at); ok {

					return entry, ex

				}

			}

		}

	}

	return nil, ""
}

func (cache *TxOHCCache) FindEntry(
	entries []common.TxOHCHistory,
	at time.Time,
) (*common.TxOHCHistory, bool) {

	var found *common.TxOHCHistory

	for i, entry := range entries {

		if at.Equal(entry.DateTime) {
			found = &entries[i]
			break
		}

		// Before since ascending sorted order expected
		if at.Before(entry.DateTime) {

			if i == 0 {
				return nil, false
			}

			found = &entries[i-1]
			break
		}

	}

	return found, found != nil

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
