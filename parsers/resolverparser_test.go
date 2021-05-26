package parsers

import (
	"testing"

	"github.com/mariotoffia/gocryptoadmin/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMultiPathExpression(t *testing.T) {

	expr := NewResolverParser().
		Parse("btx:USDT = btx,all:USD -> ofx,all:EUR").
		GetExpressions()

	assert.Equal(t, 1, len(expr))
	assert.Equal(t, 1, len(expr[0].AssetPrefixes))
	assert.Equal(t, "btx", expr[0].AssetPrefixes[0])
	assert.Equal(t, common.AssetTypeUSDT, expr[0].Asset)

	assert.Equal(t, 2, len(expr[0].Path))
	assert.Equal(t, 2, len(expr[0].Path[0].AssetPrefixes))
	assert.Equal(t, "btx", expr[0].Path[0].AssetPrefixes[0])
	assert.Equal(t, "all", expr[0].Path[0].AssetPrefixes[1])
	assert.Equal(t, "USDT-USD", expr[0].Path[0].AssetPair.String())

	assert.Equal(t, 2, len(expr[0].Path[1].AssetPrefixes))
	assert.Equal(t, "ofx", expr[0].Path[1].AssetPrefixes[0])
	assert.Equal(t, "all", expr[0].Path[1].AssetPrefixes[1])
	assert.Equal(t, "USD-EUR", expr[0].Path[1].AssetPair.String())
}

func TestNoPrefixMeansAll(t *testing.T) {

	expr := NewResolverParser().
		Parse("USDT = btx:USD -> all:EUR").
		GetExpressions()

	assert.Equal(t, 1, len(expr))
	require.Equal(t, 1, len(expr[0].AssetPrefixes))
	assert.Equal(t, "all", expr[0].AssetPrefixes[0])
	assert.Equal(t, common.AssetTypeUSDT, expr[0].Asset)

	require.Equal(t, 2, len(expr[0].Path))
	require.Equal(t, 1, len(expr[0].Path[0].AssetPrefixes))
	assert.Equal(t, "btx", expr[0].Path[0].AssetPrefixes[0])
	assert.Equal(t, "USDT-USD", expr[0].Path[0].AssetPair.String())

	require.Equal(t, 1, len(expr[0].Path[1].AssetPrefixes))
	assert.Equal(t, "all", expr[0].Path[1].AssetPrefixes[0])
	assert.Equal(t, "USD-EUR", expr[0].Path[1].AssetPair.String())
}
