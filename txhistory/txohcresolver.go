package txhistory

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
// .Resolver Pattern
// ====
// 1. btx:USDT = btx,all:USDT-USD -> ofx,all:USD-EUR
// 2. N/A
// 3. all:BTC = cbx,all:BTC-EUR
// 4. btx:BTC = btx,all:BTC-USDT (it will find 1 since now USDT)
// ====
//
// Where prefix before equal sign means exchange. When it is _all_, it
// will try to use the submitted exchange first, if miss, it will try
// `common.ExchangeAll` in cache before giving up.
//
// Prefix Resolve Order
//
// 1. explicit first
// 2. all by the submitted exchange (if any)
// 3. all by `common.ExchangeAll`
//
// Prefix after equal sign is resolved in stated order where again _all_ is
// first tried with the submitted _exchange_ parameter, if fails, it will try
// `common.ExchangeAll`.
//
// All patterns are terminated with a new-line.
type TxOHCResolver struct {
}
