package coinbasepro

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/mariotoffia/gocryptoadmin/common"
	"github.com/mariotoffia/gocryptoadmin/utils"
)

// granularity=86400s => 1440min (1D)
// start, end
// https://api.pro.coinbase.com/products/BTC-USD/candles?start=2017-08-31T00:00:00&granularity=86400

type Coinbase struct {
	baseURL  string
	format   string
	exchange string
}

type QueryRange struct {
	start       time.Time
	end         time.Time
	granularity int
}

func New(baseURL string) *Coinbase {

	if baseURL == "" {
		baseURL = "https://api.pro.coinbase.com/products"
	}

	return &Coinbase{
		baseURL:  baseURL,
		exchange: "cbx",
		format:   "2006-01-02T15:04:05Z",
	}
}

func (k *Coinbase) SetExchangeName(name string) {
	k.exchange = name
}

func (cbx *Coinbase) Read(
	pair common.AssetPair,
	since time.Time,
	interval time.Duration,
) []common.TxOHCHistory {

	granularity := interval / time.Second
	if granularity > 86400 || granularity < 60 {

		panic(
			"interval must be the following seconds:" +
				"{60, 300, 900, 3600, 21600, 86400}",
		)

	}

	list := []common.TxOHCHistory{}

	for _, qr := range cbx.calcRanges(since, int(granularity)) {

		list = append(
			list, cbx.getRange(pair, interval, &qr)...,
		)

	}

	return list
}

func (cbx *Coinbase) getRange(
	pair common.AssetPair,
	interval time.Duration,
	qr *QueryRange,
) []common.TxOHCHistory {

	list := []common.TxOHCHistory{}

	req := fmt.Sprintf(
		"%s/%s/candles?start=%s&end=%s&granularity=%d",
		cbx.baseURL, pair.String(),
		qr.start.Format(cbx.format),
		qr.end.Format(cbx.format),
		qr.granularity,
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

	var results []interface{}
	if err = json.Unmarshal(data, &results); err != nil {
		panic(fmt.Sprintf("data: %s, err: %s", string(data), err.Error()))
	}

	for _, v := range results {

		ohlc := v.([]interface{})
		list = append(list, cbx.toEntry(ohlc, pair, interval))

	}

	return list
}

func (cbx *Coinbase) toEntry(
	arr []interface{},
	pair common.AssetPair,
	interval time.Duration) common.TxOHCHistory {

	entry := common.TxOHCHistory{
		Exchange:    cbx.exchange,
		AssetPair:   pair,
		DateTime:    time.Unix(int64(arr[0].(float64)), 0).UTC(),
		Resolution:  int(interval / time.Minute),
		Low:         arr[1].(float64),
		High:        arr[2].(float64),
		Open:        arr[3].(float64),
		Close:       arr[4].(float64),
		AssetVolume: arr[5].(float64),
	}

	entry.ID = utils.ToString(utils.HashFromTime(entry.DateTime))

	return entry
}

func (cbx *Coinbase) calcRanges(since time.Time, granularity int) []QueryRange {

	now := time.Now()
	batches := int(now.Sub(since).Seconds()/float64(granularity*300)) + 1

	ranges := []QueryRange{}
	start := since

	for i := 0; i < batches; i++ {

		end := start.Add(time.Second * time.Duration(300*granularity))

		if end.After(now) {
			end = now
		}

		ranges = append(ranges, QueryRange{
			start:       start,
			end:         end,
			granularity: granularity,
		})

		start = end

		if start.After(now) {
			break
		}
	}

	return ranges
}
