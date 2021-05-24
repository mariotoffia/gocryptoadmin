package txhistory

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/mariotoffia/gocryptoadmin/common"
	"github.com/mariotoffia/gocryptoadmin/txhistory/kraken"
)

func TestReadFromKraken(t *testing.T) {

	txr := NewTxOHCReader().Register("kr", kraken.New(""))
	from, _ := time.Parse(time.RFC3339, "2020-01-01T00:00:00.000Z")

	entries := txr.Read(common.AssetPair{
		Asset:    common.AssetTypeBTC,
		CostUnit: common.AssetTypeEuro,
	}, from, time.Hour*24, "kr")

	data, _ := json.Marshal(entries)
	fmt.Println(string(data))
	fmt.Println(len(entries))
}
