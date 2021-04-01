package common

// AssetType is the name of an asset e.g. EUR, BTC, XLM etc.
type AssetType string

const (
	AssetTypeEuro AssetType = "EUR"
	AssetTypeUSDT AssetType = "USDT"
	AssetTypeBTC  AssetType = "BTC"
	AssetTypeLTC  AssetType = "LTC"
	AssetTypeETH  AssetType = "ETH"
	AssetTypeETC  AssetType = "ETC"
	AssetTypeXRP  AssetType = "XRP"
	AssetTypeCVC  AssetType = "CVC"
	AssetTypeXLM  AssetType = "XLM"
)

type AssetPair struct {
	Asset    AssetType `csv:"asset"    json:"asset"`
	CostUnit AssetType `csv:"costunit" json:"costunit"`
}

func (ap AssetPair) String() string {
	return string(ap.Asset) + "-" + string(ap.CostUnit)
}
