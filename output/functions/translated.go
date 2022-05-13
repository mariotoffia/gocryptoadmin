package functions

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mariotoffia/gocryptoadmin/common"
	"github.com/mariotoffia/gocryptoadmin/utils"
)

func translated(value interface{}, command, text, fee string, assets ...string) string {

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

	if command == "header" || command == "csv-header" {

		csv := strings.HasPrefix(command, "csv-")

		s := ""
		for _, asset := range list {

			if csv {

				s += fmt.Sprintf("%s %v;%s %v;", text, asset, fee, asset)

			} else {

				assetLen := len(string(asset))

				s += fmt.Sprintf(
					"%s %v%s|%s %v%s|",
					text,
					asset,
					strings.Repeat(" ", 16-(len(text)+assetLen)),
					fee,
					asset,
					strings.Repeat(" ", 12-(len(fee)+assetLen)),
				)

			}

		}

		return s
	}

	if command == "separator" {

		return strings.Repeat("-", len(list)*32)

	}

	if command == "" {
		command = "total-and-fee"
	}

	if strings.HasPrefix(command, "total-and-fee") ||
		strings.HasPrefix(command, "csv-total-and-fee") {

		positive := strings.HasSuffix(command, "-positive")
		csv := strings.HasPrefix(command, "csv-")

		s := ""
		for _, asset := range list {

			tot := entry.GetTranslatedTotalPrice(asset)

			if positive && tot < 0 {
				tot = -tot
			}

			if csv {

				s += fmt.Sprintf(
					"%.2f;%.2f;",
					tot,
					entry.GetTranslatedFee(asset),
				)

			} else {

				s += fmt.Sprintf(
					"% -17.2f|% -13.2f|",
					tot,
					entry.GetTranslatedFee(asset),
				)

			}

		}

		return s

	}

	return fmt.Sprintf("unknown command: %s", command)
}
