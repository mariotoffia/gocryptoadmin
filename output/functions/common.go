package functions

import "github.com/mariotoffia/gocryptoadmin/common"

func toFirstEntry(value interface{}) common.TransactionEntry {

	if e, ok := value.([]common.TransactionEntry); ok {

		if len(e) == 0 {
			return nil
		}

		return e[0]

	} else if e, ok := value.(common.AccountEntry); ok {
		return e
	} else if e, ok := value.([]common.AccountEntry); ok {

		if len(e) == 0 {
			return nil
		}

		return e[0]

	} else if e, ok := value.(common.TxBuySellEntry); ok {
		return e
	} else if e, ok := value.([]common.TxBuySellEntry); ok {

		if len(e) == 0 {
			return nil
		}

		return e[0]

	} else if e, ok := value.(*common.TxBuyGroupLog); ok {
		return e
	} else if e, ok := value.([]*common.TxBuyGroupLog); ok {
		if len(e) == 0 {
			return nil
		}

		return e[0]

	} else if e, ok := value.(*common.TransactionLog); ok {
		return e
	} else if e, ok := value.([]*common.TransactionLog); ok {

		if len(e) == 0 {
			return nil
		}

		return e[0]

	} else if e, ok := value.(*common.TxGroupEntry); ok {
		return e
	} else if e, ok := value.([]*common.TxGroupEntry); ok {

		if len(e) == 0 {
			return nil
		}

		return e[0]

	}

	return nil
}
