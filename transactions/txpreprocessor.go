package transactions

import (
	"sort"
	"time"

	"github.com/mariotoffia/gocryptoadmin/transactions/txcommon"
)

func (lr *TxLogReaderImpl) preProcess(logs []txcommon.Transaction) []txcommon.Transaction {

	sort.Slice(logs, func(i, j int) bool {
		return logs[i].CreatedAt.Before(logs[j].CreatedAt)
	})

	type processor struct {
		next        time.Time
		groupId     int64
		intervalAcc txcommon.Accumulate
	}

	groupId := int64(0)
	products := map[string]processor{}
	lastSide := map[string]txcommon.SideType{}
	accAsset := map[string]txcommon.Accumulate{}

	for i, tx := range logs {

		size := tx.Size
		if tx.Side == txcommon.SideTypeSell {
			size = -tx.Size
		}

		if _, ok := products[tx.Exchange+tx.Asset+string(tx.Side)]; !ok {

			lastSide[tx.Exchange+tx.Asset] = tx.Side

			accAsset[tx.Exchange+tx.Asset] = txcommon.Accumulate{
				AccSize:  size,
				AccFee:   tx.Fee,
				AccTotal: tx.Total,
			}

			groupId++

			products[tx.Exchange+tx.Asset+string(tx.Side)] = processor{
				groupId: groupId,
				next:    tx.CreatedAt.Add(time.Second * lr.secwindow),
				intervalAcc: txcommon.Accumulate{
					AccSize:  size,
					AccFee:   tx.Fee,
					AccTotal: tx.Total,
				},
			}

			logs[i].AccSize = size
			logs[i].AccTotal = tx.Total
			logs[i].AccFee = tx.Fee
			logs[i].GrpSize = size
			logs[i].GrpTotal = tx.Total
			logs[i].GrpFee = tx.Fee
			logs[i].GroupID = groupId

			continue
		}

		// Update asset accumulation properties
		acc := accAsset[tx.Exchange+tx.Asset]
		acc.AccSize = acc.AccSize + size
		acc.AccFee = acc.AccFee + tx.Fee
		acc.AccTotal = acc.AccTotal + tx.Total
		accAsset[tx.Exchange+tx.Asset] = acc

		proc := products[tx.Exchange+tx.Asset+string(tx.Side)]

		proc.intervalAcc.AccSize = proc.intervalAcc.AccSize + size
		proc.intervalAcc.AccTotal = proc.intervalAcc.AccTotal + tx.Cost.Total
		proc.intervalAcc.AccFee = proc.intervalAcc.AccFee + tx.Cost.Fee

		if lastSide[tx.Exchange+tx.Asset] == tx.Side &&
			tx.CreatedAt.Before(proc.next) {

		} else {

			// Next slice
			groupId++
			proc.groupId = groupId
			proc.intervalAcc.AccSize = size
			proc.intervalAcc.AccFee = tx.Fee
			proc.intervalAcc.AccTotal = tx.Total
			proc.next = tx.CreatedAt.Add(time.Second * lr.secwindow)

			lastSide[tx.Exchange+tx.Asset] = tx.Side // Might have changed

		}

		logs[i].GrpSize = proc.intervalAcc.AccSize
		logs[i].GrpFee = proc.intervalAcc.AccFee
		logs[i].GrpTotal = proc.intervalAcc.AccTotal
		logs[i].GroupID = proc.groupId
		logs[i].AccSize = acc.AccSize
		logs[i].AccTotal = acc.AccTotal
		logs[i].AccFee = acc.AccFee

		products[tx.Exchange+tx.Asset+string(tx.Side)] = proc

	}

	return logs
}
