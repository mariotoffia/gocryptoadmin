package functions

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mariotoffia/gocryptoadmin/common"
	"github.com/mariotoffia/gocryptoadmin/utils"
)

func tax(value interface{}, command, text string, tax float64, assets ...string) string {

	entry := toFirstEntry(value)

	if entry == nil {

		panic(
			fmt.Sprintf(
				"expecting either array or scalar common.TransactionEntry, got: %T", value,
			),
		)

	}

	list := []common.AssetType{}
	for _, asset := range entry.GetTranslatedAssets() {

		if len(assets) == 0 || utils.StringContains(assets, string(asset)) {
			list = append(list, asset)
		}

	}

	sort.Slice(list, func(i, j int) bool {
		return list[i] < list[j]
	})

	if len(list) == 0 {
		return ""
	}

	if command == "header" {

		s := ""
		for _, asset := range list {
			s += fmt.Sprintf("%s %-5v", text, asset)
		}

		return s
	}

	if command == "separator" {

		l := len(text) + 22
		return strings.Repeat("-", len(list)*l)

	}

	var bs common.TxBuySellEntry
	if e, ok := entry.(common.TxBuySellEntry); ok {
		bs = e
	} else {
		panic(fmt.Sprintf("expecting TxBuySellEntry, found: %T", e))
	}

	buy := bs.GetBuy()
	sell := bs.GetSell()

	tax /= 100

	if command == "tax-all" {

		s := ""
		for _, asset := range list {

			buyPrice := buy.GetTranslatedTotalPrice(asset)
			sellPrice := sell.GetTranslatedTotalPrice(asset)
			taxed := utils.ToFixed((sellPrice+buyPrice)*tax, 8)

			s += fmt.Sprintf("% -13f|", taxed)

		}

		return s

	}

	return fmt.Sprintf("unknown command: %s", command)
}
