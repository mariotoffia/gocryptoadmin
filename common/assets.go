package common

// AssetType is the name of an asset e.g. EUR, BTC, XLM etc.
type AssetType string

const (
	AssetTypeEuro        AssetType = "EUR"
	AssetTypeSvenskKrona AssetType = "SEK"
	AssetTypeUSDT        AssetType = "USDT"
	AssetTypeBTC         AssetType = "BTC"
	AssetTypeLTC         AssetType = "LTC"
	AssetTypeETH         AssetType = "ETH"
	AssetTypeETC         AssetType = "ETC"
	AssetTypeXRP         AssetType = "XRP"
	AssetTypeCVC         AssetType = "CVC"
	AssetTypeXLM         AssetType = "XLM"
)

// IsFIAT checks if the `AssetType` is plain FIAT or crypto currency
func (asset AssetType) IsFIAT() bool {
	return asset == AssetTypeEuro || asset == AssetTypeSvenskKrona
}

// IsCrypto returns `true` when the _asset_ is a non _FIAT_ currency.
func (asset AssetType) IsCrypto() bool {
	return !asset.IsFIAT()
}

// IsTether returns `true` if the _asset_ is a tether (fake crypto _FIAT_) such
// as _USDT_.
func (asset AssetType) IsTether() bool {
	return asset == AssetTypeUSDT
}

type AssetPair struct {
	Asset    AssetType `csv:"asset"    json:"asset"`
	CostUnit AssetType `csv:"costunit" json:"costunit"`
}

func (ap AssetPair) String() string {
	return string(ap.Asset) + "-" + string(ap.CostUnit)
}
