package kraken

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/mariotoffia/gocryptoadmin/common"
	"github.com/mariotoffia/gocryptoadmin/utils"
)

// Kraken reads from the public OHLC API.
type Kraken struct {
	baseURL  string
	exchange string
}

func New(baseURL string) *Kraken {

	if baseURL == "" {
		baseURL = "https://api.kraken.com/0/public/"
	}

	return &Kraken{
		baseURL:  baseURL,
		exchange: "kraken",
	}
}

func (k *Kraken) SetExchangeName(name string) {
	k.exchange = name
}

func (k *Kraken) Read(
	pair common.AssetPair,
	since time.Time,
	interval time.Duration,
) []common.TxOHCHistory {

	list := []common.TxOHCHistory{}

	req := fmt.Sprintf(
		"%s/OHLC?pair=%s%s&interval=%d&since=%d",
		k.baseURL, pair.Asset, pair.CostUnit,
		interval/time.Minute,
		since.Unix(),
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

	var m map[string]interface{}
	if err = json.Unmarshal(data, &m); err != nil {
		panic(err)
	}

	results := m["result"].(map[string]interface{})

	for name, v := range results {

		if name == "last" {
			continue
		}

		ohlc := v.([]interface{})
		for i := range ohlc {

			arr := ohlc[i].([]interface{})
			list = append(list, k.toEntry(arr, pair, interval))

		}
	}

	return list
}

func (k *Kraken) toEntry(
	arr []interface{},
	pair common.AssetPair,
	interval time.Duration) common.TxOHCHistory {

	entry := common.TxOHCHistory{
		Exchange:   k.exchange,
		AssetPair:  pair,
		DateTime:   time.Unix(int64(arr[0].(float64)), 0).UTC(),
		Resolution: int(interval / time.Minute),
		Open:       utils.Float64FromString(arr[1].(string)),
		High:       utils.Float64FromString(arr[2].(string)),
		Low:        utils.Float64FromString(arr[3].(string)),
		Close:      utils.Float64FromString(arr[4].(string)),
		/* arr[5] == string <vwap>*/
		CostUnitVolume: utils.Float64FromString(arr[6].(string)),
		AssetVolume:    arr[7].(float64),
	}

	entry.ID = utils.ToString(utils.HashFromTime(entry.DateTime))

	return entry
}
