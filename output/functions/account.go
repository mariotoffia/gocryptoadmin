package functions

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mariotoffia/gocryptoadmin/common"
	"github.com/mariotoffia/gocryptoadmin/utils"
)

func account(value interface{}, command string, assets ...string) string {

	optional := false
	if len(assets) > 0 {
		if strings.ToLower(assets[0]) == "optional" {
			optional = true
			assets = assets[1:]
		}
	}

	entry := toFirstEntry(value)
	if entry == nil {

		if optional {
			return ""
		}

		panic(
			fmt.Sprintf(
				"expecting either array or scalar common.TransactionEntry, got: %T", value,
			),
		)

	}

	var accentry common.AccountEntry
	if e, ok := entry.(common.AccountEntry); ok {
		accentry = e
	} else {

		if optional {
			return ""
		}

		panic(fmt.Sprintf("expecting account entry, found: %T", e))
	}

	list := []common.AssetType{}
	status := accentry.GetAccountStatus()
	for asset := range status {

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
			s += fmt.Sprintf("Account %-9v|", asset)
		}

		return s
	}

	if command == "separator" {
		return strings.Repeat("-", len(list)*18)
	}

	s := ""
	for _, asset := range list {
		s += fmt.Sprintf("% -17.8f|", status[asset])
	}

	return s
}
