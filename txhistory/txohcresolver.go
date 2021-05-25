package txhistory

import (
	"fmt"
	"time"

	"github.com/mariotoffia/gocryptoadmin/common"
	"github.com/mariotoffia/gocryptoadmin/parsers"
)

// AssetTranslation is keyed with the `common.AssetType` and a translation
// expression.
type AssetTranslation map[string]parsers.ResolverExpression

// ExchangeAssetTranslation is keyed with exchange or _all_ and then
// the `common.AssetType` and the expression.
type ExchangeAssetTranslation map[string]AssetTranslation

// TxOHCResolver resolves a assetType to a another one
// possibly by one or more indirections.
//
// .Wanted behavior
// ====
// 1. Bittrex: BTC-USDT -> USDT-USD -> USD-EUR -> EUR (btx)
// 2. All: BTC-EUR -> EUR
// 3. All: LTC-BTC -> BTC-EUR -> EUR (cbx)
// 4. Bittrex: LTC-BTC -> BTC-USDT -> USDT-USD -> USD-EUR -> EUR (btx)
// ====
//
// .Resolver Expressions
// ====
// 1. btx:USDT = btx,all:USD -> ofx,all:EUR
// 2. N/A
// 3. all:BTC = cbx,all:EUR
// 4. btx:BTC = btx,all:USDT
// ====
//
// NOTE: That the above expressions expands, e..g. expression 1 will
// be expanded to btx:USDT = btx,all:USDT-USD -> ofx,all:USD-EUR.
//
// Since last expression (4) is translated to USDT and it do exist, it
// will continue to expand using (1) until a _FIAT_ is reached.
//
// Where prefix before equal sign means exchange. When it is _all_, it
// will try to use the submitted exchange first, if miss, it will try
// `common.ExchangeAll` in cache before giving up (if _all_ is submitted).
//
// TIP: If prefix is omitted, _all_ is automatically added e.g. it is possible
// to rewrite (3) from _all:BTC = cbx,all:EUR_ to _BTC = cbx,all:EUR_.
//
// Prefix Resolve Order
//
// 1. explicit first
// 2. all by the submitted exchange (if any)
// 3. all by `common.ExchangeAll`
//
// Prefix after equal sign is resolved in stated order where again _all_ is
// first tried with the submitted _exchange_ parameter, if fails, it will try
// `common.ExchangeAll` (if _all_ is submitted).
//
// All patterns are terminated with a new-line.
type TxOHCResolver struct {
	assets ExchangeAssetTranslation
	cache  *TxOHCCache
}

type ResolveAcceptResult int

const (
	// ResolveAcceptResultContinue denotes that it should continue it search and add the entry to
	// the search result.
	ResolveAcceptResultContinue ResolveAcceptResult = 0
	// ResolveAcceptResultFail specifies that the search is over and no result should be yielded
	ResolveAcceptResultFail ResolveAcceptResult = -1
	// ResolveAcceptResultAccept denotes that the search is over and the result is ready
	ResolveAcceptResultAccept ResolveAcceptResult = 1
)

// ResolveAcceptFunction is called each time it has found a match
// and query if it should continue or stop searching.
//
// NOTE: Both _path_ and _entry_ is *guaranteed* to be valid.
type ResolveAcceptFunction func(
	path parsers.ResolverExpressionPathItem,
	entry *common.TxOHCHistory,
) ResolveAcceptResult

func NewTxOHCResolver(cache *TxOHCCache) *TxOHCResolver {

	return &TxOHCResolver{
		assets: ExchangeAssetTranslation{},
		cache:  cache,
	}

}

func (resolver *TxOHCResolver) AddTranslations(expr ...parsers.ResolverExpression) *TxOHCResolver {

	for _, e := range expr {

		for _, prefix := range e.AssetPrefixes {

			ass := resolver.assets[prefix]
			if ass == nil {
				ass = AssetTranslation{}
				resolver.assets[prefix] = ass
			}

			if _, ok := ass[string(e.Asset)]; ok {

				panic(
					fmt.Sprintf(
						"asset: %s is already cached in exchange: %s", e.Asset, prefix,
					),
				)

			}

			ass[string(e.Asset)] = e
		}

	}

	return resolver
}

type ResolvedOHCEntry struct {
	Entry             common.TxOHCHistoryEntry
	Exchange          string
	exchangeSelection []string
	AssetPair         common.AssetPair
}

// ResolveTarget will search beginning with _asset_ and path down to _target_ by
// adding resolved entries in path (excluding _asset_) _at_ the time.
//
// If succeeds to find _target_ (as _CostUnit_ on resolved path elements) it will
// return `true`, otherwise `false`.
func (resolver *TxOHCResolver) ResolveToTarget(
	at time.Time,
	asset common.AssetType,
	target common.AssetType,
	exchange ...string,
) ([]ResolvedOHCEntry, bool) {

	return resolver.Resolve(
		at, asset,
		func(
			path parsers.ResolverExpressionPathItem, entry *common.TxOHCHistory,
		) ResolveAcceptResult {

			if path.AssetPair.CostUnit == target {
				return ResolveAcceptResultAccept
			}

			return ResolveAcceptResultContinue

		}, exchange...)

}

// ResolveTarget will search beginning with _asset_ and path down to _fiat_ by
// adding resolved entries in path (excluding _asset_) _at_ the time.
//
// If succeeds to find _fiat_ (as _CostUnit_ on resolved path elements) it will
// return `true`, otherwise `false`.
//
// The _fiat_ parameter specifies how may times it should encounter a _FIAT_ as _CostUnit_
// until declaring success.
//
// This could be used to resolve the following:
// 1. USDT = USD -> EUR (fiat = 2)
// 2. EUR = SEK (fiat 1)
//
// Or if _fiat_ is set to 3 and _USDT_ it could resolve to UST -> SEK.
func (resolver *TxOHCResolver) ResolveToFIAT(
	at time.Time,
	asset common.AssetType,
	fiat int,
	exchange ...string,
) ([]ResolvedOHCEntry, bool) {

	return resolver.Resolve(
		at, asset,
		func(
			path parsers.ResolverExpressionPathItem, entry *common.TxOHCHistory,
		) ResolveAcceptResult {

			if path.AssetPair.CostUnit.IsFIAT() {
				fiat--
			}

			if fiat == 0 {
				return ResolveAcceptResultAccept
			}

			return ResolveAcceptResultContinue

		}, exchange...)

}

// ResolveTarget will search beginning with _asset_ and path down until _accept_ function
// is stating `ResolveAcceptResultContinue` by adding resolved entries in path
// (excluding _asset_) _at_ the time.
//
// If succeeds to find it will return `true`, otherwise `false`.
func (resolver *TxOHCResolver) Resolve(
	at time.Time,
	asset common.AssetType,
	accept ResolveAcceptFunction,
	exchange ...string,
) ([]ResolvedOHCEntry, bool) {

	var res []ResolvedOHCEntry

	walkedPath := false
	for _, ex := range exchange {

		if walkedPath {
			break
		}

		if ass, ok := resolver.assets[ex]; ok {

			if expr, ok := ass[string(asset)]; ok {

				walkedPath = true
				for _, path := range expr.Path {

					entry, foundExchange := resolver.cache.GetEntryForAssset(
						path.AssetPair,
						at,
						path.AssetPrefixes...)

					if entry == nil {
						return nil, false
					}

					res = append(res, ResolvedOHCEntry{
						Entry:             entry,
						Exchange:          foundExchange,
						AssetPair:         path.AssetPair,
						exchangeSelection: path.AssetPrefixes,
					})

					switch accept(path, entry) {
					case ResolveAcceptResultFail:
						return nil, false
					case ResolveAcceptResultAccept:
						return res, true
					}

				}

			}

		}

	}

	if len(res) == 0 {
		return res, false
	}

	last := res[len(res)-1]

	if cont, ok := resolver.Resolve(
		at, last.AssetPair.CostUnit, accept, last.exchangeSelection...,
	); ok {

		return append(res, cont...), true

	}

	return nil, false
}
