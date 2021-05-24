package common

import (
	"time"
)

// TxOHCReader reads from it's datasource and returns the result.
type TxOHCReader interface {
	Read(pair AssetPair, since time.Time, interval time.Duration) []TxOHCHistory
	// SetExchangeName alters the default name.
	SetExchangeName(name string)
}

// TxOHCHistoryEntry represents a single historic transaction entry
// for a single exchange.
//
// If not bound to a exchage, the `ExchangeAll` is used.
type TxOHCHistoryEntry interface {
	GetID() string
	// GetResolution returns the number of minutes for entry
	//
	// 1, 5, 15, 30, 60 (1h), 240 (4h), 1440 (1D), 10080 (1W)
	GetResolution() int
	GetExchange() string
	GetDateTime() time.Time
	GetOpen() float64
	GetHigh() float64
	GetLow() float64
	GetClose() float64
	GetVolumeAsset() float64
	GetVolumeCostUnit() float64
	GetAssetPair() AssetPair
}

type TxOHCHistory struct {
	ID             string    `csv:"id"               json:"id"`
	Resolution     int       `csv:"resolution"       json:"resolution"`
	Exchange       string    `csv:"exchange"         json:"exchange"`
	DateTime       time.Time `csv:"time"             json:"time"`
	Open           float64   `csv:"open"             json:"open"`
	High           float64   `csv:"high"             json:"high"`
	Low            float64   `csv:"low"              json:"low"`
	Close          float64   `csv:"close"            json:"close"`
	AssetVolume    float64   `csv:"asset volume"     json:"assetvolume"`
	CostUnitVolume float64   `csv:"cost-unit volume" json:"cuvolume"`

	AssetPair
}

func (ohc *TxOHCHistory) GetID() string {
	return ohc.ID
}

func (ohc *TxOHCHistory) GetResolution() int {
	return ohc.Resolution
}

func (ohc *TxOHCHistory) GetExchange() string {
	return ohc.Exchange
}
func (ohc *TxOHCHistory) GetDateTime() time.Time {
	return ohc.DateTime
}
func (ohc *TxOHCHistory) GetOpen() float64 {
	return ohc.Open
}
func (ohc *TxOHCHistory) GetHigh() float64 {
	return ohc.High
}
func (ohc *TxOHCHistory) GetLow() float64 {
	return ohc.Low
}
func (ohc *TxOHCHistory) GetClose() float64 {
	return ohc.Close
}
func (ohc *TxOHCHistory) GetVolumeAsset() float64 {
	return ohc.AssetVolume
}
func (ohc *TxOHCHistory) GetVolumeCostUnit() float64 {
	return ohc.CostUnitVolume
}
func (ohc *TxOHCHistory) GetAssetPair() AssetPair {
	return ohc.AssetPair
}
