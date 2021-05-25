package txhistory

import (
	"fmt"
	"testing"
	"time"

	"github.com/mariotoffia/gocryptoadmin/common"
	"github.com/mariotoffia/gocryptoadmin/parsers"
	"github.com/stretchr/testify/require"
)

func TestResolveUSDTtoUSDtoEURtoSEK(t *testing.T) {

	expr := parsers.NewResolverParser().
		Parse("btx:USDT = btx,all:USD -> ofx,all:EUR").
		Parse("EUR = SEK").
		GetExpressions()

	cache := NewTxOHCCache().Load(
		"testfiles/cache-test",
		func(
			cache *TxOHCCache, exchange string, entries []common.TxOHCHistory,
		) {

			cache.Add(entries, common.ExchangeAll) // make visible to all as well

		})

	resolver := NewTxOHCResolver(cache).AddTranslations(expr...)

	at, _ := time.Parse(time.RFC3339, "2018-08-31T00:00:00.000Z")

	result, ok := resolver.ResolveToFIAT(at, common.AssetTypeUSDT, 3, "btx", common.ExchangeAll)

	require.Equal(t, true, ok)

	fmt.Println(result)

}
