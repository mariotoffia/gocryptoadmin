package txhistory

import (
	"testing"
	"time"

	"github.com/mariotoffia/gocryptoadmin/common"
	"github.com/mariotoffia/gocryptoadmin/parsers"
	"github.com/stretchr/testify/assert"
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
	require.Equal(t, 3, len(result), "Should have USDT-USD, USD-EUR, and EUR-SEK")

	assert.Equal(t, "2018-08-31T00:00:00Z", result[0].Entry.GetDateTime().Format(time.RFC3339))
	assert.Equal(t, "USDT-USD", result[0].Entry.GetAssetPair().String())
	assert.Equal(t, "btx", result[0].Exchange)
	assert.Equal(t, float64(1), result[0].Entry.GetOpen())
	assert.Equal(t, float64(1), result[0].Entry.GetHigh())
	assert.Equal(t, float64(0.982), result[0].Entry.GetLow())
	assert.Equal(t, float64(0.996), result[0].Entry.GetClose())

	assert.Equal(t, "2018-08-31T00:00:00Z", result[1].Entry.GetDateTime().Format(time.RFC3339))
	assert.Equal(t, "USD-EUR", result[1].Entry.GetAssetPair().String())
	assert.Equal(t, "ofx", result[1].Exchange)
	assert.Equal(t, float64(0.863015), result[1].Entry.GetOpen())
	assert.Equal(t, float64(0.863015), result[1].Entry.GetHigh())
	assert.Equal(t, float64(0.863015), result[1].Entry.GetLow())
	assert.Equal(t, float64(0.863015), result[1].Entry.GetClose())

	assert.Equal(t, "2018-08-31T00:00:00Z", result[2].Entry.GetDateTime().Format(time.RFC3339))
	assert.Equal(t, "EUR-SEK", result[2].Entry.GetAssetPair().String())
	assert.Equal(t, "all", result[2].Exchange, "Since expression is EUR = SEK (default to all)")
	assert.Equal(t, float64(10.611722), result[2].Entry.GetOpen())
	assert.Equal(t, float64(10.611722), result[2].Entry.GetHigh())
	assert.Equal(t, float64(10.611722), result[2].Entry.GetLow())
	assert.Equal(t, float64(10.611722), result[2].Entry.GetClose())

}

func TestResolveToSEKFIATWillBeResolvedUsingTwoExpressions(t *testing.T) {

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

	result, ok := resolver.ResolveToTarget(
		at,
		common.AssetTypeUSDT,
		common.AssetTypeSvenskKrona,
		"btx",
		common.ExchangeAll,
	)

	require.Equal(t, true, ok)
	require.Equal(t, 3, len(result), "Should have USDT-USD, USD-EUR, and EUR-SEK")

	assert.Equal(t, "2018-08-31T00:00:00Z", result[0].Entry.GetDateTime().Format(time.RFC3339))
	assert.Equal(t, "USDT-USD", result[0].Entry.GetAssetPair().String())
	assert.Equal(t, "btx", result[0].Exchange)
	assert.Equal(t, float64(1), result[0].Entry.GetOpen())
	assert.Equal(t, float64(1), result[0].Entry.GetHigh())
	assert.Equal(t, float64(0.982), result[0].Entry.GetLow())
	assert.Equal(t, float64(0.996), result[0].Entry.GetClose())

	assert.Equal(t, "2018-08-31T00:00:00Z", result[1].Entry.GetDateTime().Format(time.RFC3339))
	assert.Equal(t, "USD-EUR", result[1].Entry.GetAssetPair().String())
	assert.Equal(t, "ofx", result[1].Exchange)
	assert.Equal(t, float64(0.863015), result[1].Entry.GetOpen())
	assert.Equal(t, float64(0.863015), result[1].Entry.GetHigh())
	assert.Equal(t, float64(0.863015), result[1].Entry.GetLow())
	assert.Equal(t, float64(0.863015), result[1].Entry.GetClose())

	assert.Equal(t, "2018-08-31T00:00:00Z", result[2].Entry.GetDateTime().Format(time.RFC3339))
	assert.Equal(t, "EUR-SEK", result[2].Entry.GetAssetPair().String())
	assert.Equal(t, "all", result[2].Exchange, "Since expression is EUR = SEK (default to all)")
	assert.Equal(t, float64(10.611722), result[2].Entry.GetOpen())
	assert.Equal(t, float64(10.611722), result[2].Entry.GetHigh())
	assert.Equal(t, float64(10.611722), result[2].Entry.GetLow())
	assert.Equal(t, float64(10.611722), result[2].Entry.GetClose())

}

func TestResolveToEURFIATWillBeResolvedUsingOneExpressions(t *testing.T) {

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

	result, ok := resolver.ResolveToTarget(
		at,
		common.AssetTypeUSDT,
		common.AssetTypeEuro,
		"btx",
		common.ExchangeAll,
	)

	require.Equal(t, true, ok)
	require.Equal(t, 2, len(result), "Should have USDT-USD, USD-EUR")

	assert.Equal(t, "2018-08-31T00:00:00Z", result[0].Entry.GetDateTime().Format(time.RFC3339))
	assert.Equal(t, "USDT-USD", result[0].Entry.GetAssetPair().String())
	assert.Equal(t, "btx", result[0].Exchange)
	assert.Equal(t, float64(1), result[0].Entry.GetOpen())
	assert.Equal(t, float64(1), result[0].Entry.GetHigh())
	assert.Equal(t, float64(0.982), result[0].Entry.GetLow())
	assert.Equal(t, float64(0.996), result[0].Entry.GetClose())

	assert.Equal(t, "2018-08-31T00:00:00Z", result[1].Entry.GetDateTime().Format(time.RFC3339))
	assert.Equal(t, "USD-EUR", result[1].Entry.GetAssetPair().String())
	assert.Equal(t, "ofx", result[1].Exchange)
	assert.Equal(t, float64(0.863015), result[1].Entry.GetOpen())
	assert.Equal(t, float64(0.863015), result[1].Entry.GetHigh())
	assert.Equal(t, float64(0.863015), result[1].Entry.GetLow())
	assert.Equal(t, float64(0.863015), result[1].Entry.GetClose())
}
