package bittrex

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/mariotoffia/gocryptoadmin/common"
	"github.com/mariotoffia/gocryptoadmin/utils"
)

type Bittrex struct {
	baseURL  string
	exchange string
}

type Point struct {
	Time        string `json:"startsAt"`
	Open        string `json:"open"`
	High        string `json:"high"`
	Low         string `json:"low"`
	Close       string `json:"close"`
	Volume      string `json:"volume"`
	QuoteVolume string `json:"quoteVolume"`
}

func New(baseURL string) *Bittrex {

	if baseURL == "" {
		baseURL = "https://api.bittrex.com/v3/markets"
	}

	return &Bittrex{
		baseURL:  baseURL,
		exchange: "btx",
	}
}

func (btx *Bittrex) SetExchangeName(name string) {
	btx.exchange = name
}

func (btx *Bittrex) Read(
	pair common.AssetPair,
	since time.Time,
	interval time.Duration,
) []common.TxOHCHistory {

	list := []common.TxOHCHistory{}

	candleInterval := toCandleInterval(interval)

	for _, period := range toReportingPeriod(candleInterval, since) {

		req := fmt.Sprintf(
			"%s/%s-%s/candles/TRADE/%s/historical/%s",
			btx.baseURL, pair.Asset, pair.CostUnit,
			candleInterval,
			period,
		)

		list = append(
			list, btx.processRequest(req, pair, interval)...,
		)

	}

	return list

}

func (btx *Bittrex) processRequest(
	req string,
	pair common.AssetPair,
	interval time.Duration,
) []common.TxOHCHistory {

	list := []common.TxOHCHistory{}

	response, err := http.Get(req)

	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	var result []Point
	if err = json.Unmarshal(data, &result); err != nil {
		panic(fmt.Sprintf("data: %s err: %s", string(data), err.Error()))
	}

	for _, v := range result {

		list = append(list, btx.toEntry(&v, pair, interval))

	}

	return list

}
func (btx *Bittrex) toEntry(
	point *Point,
	pair common.AssetPair,
	interval time.Duration) common.TxOHCHistory {

	t, err := time.Parse(time.RFC3339, point.Time)

	if err != nil {
		panic(err)
	}

	entry := common.TxOHCHistory{
		Exchange:   btx.exchange,
		AssetPair:  pair,
		DateTime:   t.UTC(),
		Resolution: int(interval / time.Minute),
		Open:       utils.Float64FromString(point.Open),
		High:       utils.Float64FromString(point.High),
		Low:        utils.Float64FromString(point.Low),
		Close:      utils.Float64FromString(point.Close),
	}

	entry.ID = utils.ToString(utils.HashFromTime(entry.DateTime))

	return entry
}

func toReportingPeriod(candleInterval string, since time.Time) []string {

	now := time.Now().UTC()
	ret := []string{}

	for since.Before(now) {

		switch candleInterval {
		case "HOUR_1":
			if since.Month() == now.Month() &&
				since.Year() == now.Year() {

				return ret

			}

			ret = append(ret, fmt.Sprintf("%d/%d", since.Year(), since.Month()))
			since = since.Add(31 * time.Hour * 24)

		case "MINUTE_5", "MINUTE_1":

			if since.Month() == now.Month() &&
				since.Year() == now.Year() &&
				since.Day() == now.Day() {

				return ret

			}

			ret = append(ret, fmt.Sprintf("%d/%d/%d", since.Year(), since.Month(), since.Day()))
			since = since.Add(time.Hour * 24)

		case "DAY_1":

			if since.Year() == now.Year() {

				return ret

			}

			ret = append(ret, fmt.Sprintf("%d", since.Year()))
			since = since.Add(366 * time.Hour * 24)

		}

	}

	return ret
}

func toCandleInterval(interval time.Duration) string {

	dur := (interval / time.Minute)

	if dur <= 5 {
		return "MINUTE_5"
	}

	if dur <= 60 {
		return "HOUR_1"
	}

	return "DAY_1"

}
