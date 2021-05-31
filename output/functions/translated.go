package functions

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mariotoffia/gocryptoadmin/common"
	"github.com/mariotoffia/gocryptoadmin/utils"
)

func translated(value interface{}, command string, assets ...string) string {

	var entry common.TransactionEntry

	if e, ok := value.([]common.TransactionEntry); ok {

		if len(e) == 0 {
			return ""
		}

		entry = e[0]

	} else if e, ok := value.(common.TransactionEntry); ok {
		entry = e
	} else {

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
			s += fmt.Sprintf("Total Price %-5v|Fee %-9v|", asset, asset)
		}

		return s
	}

	if command == "separator" {
		return strings.Repeat("-", len(list)*32)
	}

	if command == "total-and-fee" || command == "" {

		s := ""
		for _, asset := range list {

			s += fmt.Sprintf(
				"% -17f|% -13f|",
				entry.GetTranslatedTotalPrice(asset),
				entry.GetTranslatedFee(asset),
			)

		}

		return s

	}

	return fmt.Sprintf("unknown command: %s", command)
}
