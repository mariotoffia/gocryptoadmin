package ofx

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/mariotoffia/gocryptoadmin/common"
	"github.com/mariotoffia/gocryptoadmin/utils"
)

// https://api.forex.se/currency/historicalexchangerates/SWE/SEK/EUR

type Ofx struct {
	baseURL  string
	exchange string
}

type Point struct {
	Time        int64   `json:"PointInTime"`
	Rate        float64 `json:"InterbankRate"`
	InverseRate float64 `json:"InverseInterbankRate"`
}

type Response struct {
	CurrentRate        float64 `json:"CurrentInterbankRate"`
	CurrentInverseRate float64 `json:"CurrentInverseInterbankRate"`
	Average            float64 `json:"Average"`
	Historical         []Point `json:"HistoricalPoints"`
}

func New(baseURL string) *Ofx {

	if baseURL == "" {
		baseURL = "https://api.ofx.com/PublicSite.ApiService/SpotRateHistory"
	}

	return &Ofx{
		baseURL:  baseURL,
		exchange: "ofx",
	}
}

func (ofx *Ofx) SetExchangeName(name string) {
	ofx.exchange = name
}

func toReportingPeriod(since time.Time) string {

	dur := time.Since(since).Hours() / 24 /*days*/

	if dur <= 7 {
		return "week"
	}

	if dur <= 93 {
		return "3month"
	}

	if dur <= 186 {
		return "6month"
	}

	if dur <= 366 {
		return "year"
	}

	if dur <= 1096 {
		return "3year"
	}

	if dur <= 1826 {
		return "5year"
	}

	return "10year"

}

func toReportingInterval(interval time.Duration) string {

	dur := (interval / time.Hour) / 24

	if dur <= 1 {
		return "daily"
	}

	if dur <= 31 {
		return "monthly"
	}

	return "yearly"

}

func (ofx *Ofx) Read(
	pair common.AssetPair,
	since time.Time,
	interval time.Duration,
) []common.TxOHCHistory {

	list := []common.TxOHCHistory{}

	req := fmt.Sprintf(
		"%s/%s/%s/%s?DecimalPlaces=6&ReportingInterval=%s&format=json",
		ofx.baseURL, toReportingPeriod(since),
		pair.Asset, pair.CostUnit,
		toReportingInterval(interval),
	)

	response, err := http.Get(req)

	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	var result Response
	if err = json.Unmarshal(data, &result); err != nil {
		panic(err)
	}

	for _, v := range result.Historical {

		list = append(list, ofx.toEntry(&v, pair, interval))

	}

	return list
}

func (ofx *Ofx) toEntry(
	point *Point,
	pair common.AssetPair,
	interval time.Duration) common.TxOHCHistory {

	entry := common.TxOHCHistory{
		Exchange:   ofx.exchange,
		AssetPair:  pair,
		DateTime:   utils.ToUnixMillisFromTimeStamp(point.Time).UTC(),
		Resolution: int(interval / time.Minute),
		Open:       point.Rate,
		High:       point.Rate,
		Low:        point.Rate,
		Close:      point.Rate,
	}

	entry.ID = utils.ToString(utils.HashFromTime(entry.DateTime))

	return entry
}
