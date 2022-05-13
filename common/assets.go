package common

import (
	"fmt"
	"strings"
)

type AssetPair struct {
	Asset    AssetType `csv:"asset"    json:"asset"`
	CostUnit AssetType `csv:"costunit" json:"costunit"`
}

func (ap AssetPair) String() string {
	return string(ap.Asset) + "-" + string(ap.CostUnit)
}

func ParseAssetPair(ap string) (AssetPair, error) {

	if ap == "" {

		return AssetPair{},
			fmt.Errorf(
				"cannot parse to assetpair from empty string",
			)

	}

	c := strings.Split(ap, "-")
	if len(c) != 2 {

		return AssetPair{},
			fmt.Errorf(
				"cannot parse to assetpair not on form Asset-CostUnit",
			)

	}

	return AssetPair{
		Asset:    AssetType(c[0]),
		CostUnit: AssetType(c[1]),
	}, nil

}

// AssetType is the name of an asset e.g. EUR, BTC, XLM etc.
type AssetType string

const (
	AssetTypeUnknown     AssetType = "UNK"
	AssetTypeEuro        AssetType = "EUR"
	AssetTypeSvenskKrona AssetType = "SEK"
	AssetTypeUsDollar    AssetType = "USD"
	AssetTypeUSDT        AssetType = "USDT"
	AssetTypeBTC         AssetType = "BTC"
	AssetTypeLTC         AssetType = "LTC"
	AssetTypeETH         AssetType = "ETH"
	AssetTypeETC         AssetType = "ETC"
	AssetTypeXRP         AssetType = "XRP"
	AssetTypeCVC         AssetType = "CVC"
	AssetTypeXLM         AssetType = "XLM"
	AssetTypeDASH        AssetType = "DASH"
	AssetTypeLSK         AssetType = "LSK"
	AssetTypeXVG         AssetType = "XVG"
	AssetTypePOWR        AssetType = "POWR"
	AssetTypeBCH         AssetType = "BCH"
	AssetTypeSALT        AssetType = "SALT"
)

// IsFIAT checks if the `AssetType` is plain FIAT or crypto currency
func (asset AssetType) IsFIAT() bool {

	return asset == AssetTypeEuro ||
		asset == AssetTypeSvenskKrona ||
		asset == AssetTypeUsDollar

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

// ExistsIn checks if _asset_ is part of _assets_.
func (asset AssetType) ExistsIn(assets ...AssetType) bool {

	for i := range assets {
		if asset == assets[i] {
			return true
		}
	}

	return false
}

// Normalize parses the name and makes sure that is matches a asset type.
//
// For example _XBT_ is translated to `AssetTypeBTC` or _ZEUR_ is translated
// to `AssetTypeEuro`.
func (asset AssetType) Normalize() AssetType {

	switch asset {
	case "XXBT", "XBT":
		return AssetTypeBTC
	case "XETH":
		return AssetTypeETH
	case "XLTC":
		return AssetTypeLTC
	case "XXRP":
		return AssetTypeXRP
	case "XXLM":
		return AssetTypeXLM
	case "ZEUR":
		return AssetTypeEuro
	case "ZUSD":
		return AssetTypeUsDollar
	case "ZSEK":
		return AssetTypeSvenskKrona
	}

	return asset
}

// ToISO translates a asset type into it's _ISO_ correspondence, if exists.
//
// For example: _BTC_ is translated to _XBT_.
func (asset AssetType) ToISO() AssetType {

	switch asset {
	case AssetTypeBTC:
		return AssetType("XBT")
	case AssetTypeETH:
		return AssetType("XETH")
	case AssetTypeLTC:
		return AssetType("XLTC")
	case AssetTypeXRP:
		return AssetType("XXRP")
	case AssetTypeXLM:
		return AssetType("XXLM")
	case AssetTypeEuro:
		return AssetType("ZEUR")
	case AssetTypeSvenskKrona:
		return AssetType("ZSEK")
	case AssetTypeUsDollar:
		return AssetType("ZUSD")
	}

	return asset
}
